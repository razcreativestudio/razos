/**
 * RAZOS Advanced Storage & File Manager Engine
 * Updated: 2026-10-25 (v3.8.7 - Robust Cross-Platform Path Normalization & File CRUD Suite)
 * Deskripsi: Engine backend untuk membaca partisi disk/HDD/USB fisik (gopsutil/disk),
 * penjelajahan direktori, pembacaan & penyimpanan berkas teks, upload, download, rename, dan hapus berkas.
 */

package filemanager

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// StorageDevice mendeskripsikan disk fisik atau partisi yang terpasang pada STB
type StorageDevice struct {
	Name       string  `json:"name"`       // Nama aslinya (misal: /DATA, /, C:, /media/usb)
	Device     string  `json:"device"`     // contoh: /dev/sda1, C:
	Mountpoint string  `json:"mountpoint"` // contoh: /DATA, /media/hdd, C:\
	Fstype     string  `json:"fstype"`     // contoh: ext4, ntfs, btrfs
	Total      uint64  `json:"total"`      // dalam bytes
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	Percent    float64 `json:"percent"`
	Label      string  `json:"label"`
	IsUSB      bool    `json:"is_usb"`
}

// FileItem mendeskripsikan satu entri berkas atau direktori
type FileItem struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	IsDir     bool      `json:"is_dir"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

// NormalizeStoragePath menstandarkan path berkas/direktori secara konsisten lintas sistem operasi (Linux vs Windows)
func NormalizeStoragePath(targetPath string) string {
	raw := strings.TrimSpace(targetPath)
	if raw == "" || raw == "." {
		if runtime.GOOS != "windows" {
			if _, err := os.Stat("/DATA"); err == nil {
				return "/DATA"
			}
			return "/"
		}
		if _, err := os.Stat("C:\\DATA"); err == nil {
			return "C:\\DATA"
		}
		return "."
	}

	// Ubah seluruh backslash ke forward slash
	clean := filepath.ToSlash(raw)

	// Skenario runtime Linux tapi request frontend membawa format path Windows
	if runtime.GOOS != "windows" {
		tpLower := strings.ToLower(clean)
		if strings.HasPrefix(tpLower, "c:/data") {
			clean = "/DATA" + clean[7:]
		} else if strings.HasPrefix(tpLower, "c:/") {
			clean = "/" + clean[3:]
		} else if tpLower == "c:" {
			clean = "/"
		}
		if clean == "" {
			clean = "/"
		}
	} else {
		// Skenario runtime Windows dev tapi request frontend membawa format path Linux /DATA
		tpLower := strings.ToLower(clean)
		if strings.HasPrefix(tpLower, "/data") {
			clean = "C:/DATA" + clean[5:]
		} else if strings.HasPrefix(clean, "/") && len(clean) > 1 && clean[1] != ':' {
			clean = "C:" + clean
		}
		if len(clean) == 2 && clean[1] == ':' {
			clean += "/"
		}
	}

	cleaned := filepath.Clean(filepath.FromSlash(clean))
	// Resolusi Windows Junction dan symbolic links jika target ada di sistem berkas
	if evalPath, evalErr := filepath.EvalSymlinks(cleaned); evalErr == nil {
		cleaned = evalPath
	}
	return filepath.ToSlash(cleaned)
}

// GetStorageDevices mendeteksi seluruh perangkat penyimpanan fisik yang valid (eMMC, MicroSD, HDD eksternal, USB)
func GetStorageDevices() ([]StorageDevice, error) {
	// false = hanya ambil physical mounts, bukan virtual/overlay mounts
	partitions, err := disk.Partitions(false)
	if err != nil {
		partitions, _ = disk.Partitions(true)
	}

	var devices []StorageDevice
	seenMounts := make(map[string]bool)

	for _, p := range partitions {
		// 1. Abaikan pseudofilesystem dan virtual devices (loop, zram, tmpfs, overlay)
		if strings.HasPrefix(p.Device, "/dev/loop") ||
			strings.HasPrefix(p.Device, "/dev/zram") ||
			strings.HasPrefix(p.Device, "overlay") ||
			strings.HasPrefix(p.Device, "shm") ||
			strings.HasPrefix(p.Device, "tmpfs") {
			continue
		}

		// 2. Abaikan filesystem type virtual / container
		fsLower := strings.ToLower(p.Fstype)
		if fsLower == "squashfs" || fsLower == "overlay" || fsLower == "tmpfs" || fsLower == "devtmpfs" || fsLower == "iso9660" {
			continue
		}

		// 3. Abaikan mountpoint internal sistem Linux, Docker, & boot
		mp := filepath.Clean(p.Mountpoint)
		if strings.HasPrefix(mp, "/sys") ||
			strings.HasPrefix(mp, "/proc") ||
			strings.HasPrefix(mp, "/dev") ||
			strings.HasPrefix(mp, "/run") ||
			strings.HasPrefix(mp, "/var/lib/docker") ||
			strings.HasPrefix(mp, "/var/lib/containerd") ||
			strings.HasPrefix(mp, "/snap") ||
			strings.HasPrefix(mp, "/boot") {
			continue
		}

		// Cegah duplikasi mountpoint
		if seenMounts[mp] {
			continue
		}
		seenMounts[mp] = true

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		// Gunakan nama asli dari mount point (misal /DATA, /, /media/hdd, C:)
		cleanLabel := filepath.ToSlash(p.Mountpoint)
		if strings.HasSuffix(cleanLabel, "/") && len(cleanLabel) > 1 {
			cleanLabel = strings.TrimSuffix(cleanLabel, "/")
		} else if strings.HasSuffix(cleanLabel, "\\") && len(cleanLabel) == 3 && cleanLabel[1] == ':' {
			cleanLabel = cleanLabel[:2]
		}

		isUSB := strings.Contains(strings.ToLower(p.Mountpoint), "media") ||
			strings.Contains(strings.ToLower(p.Mountpoint), "usb") ||
			strings.HasPrefix(p.Mountpoint, "/mnt")

		devices = append(devices, StorageDevice{
			Name:       cleanLabel,
			Device:     p.Device,
			Mountpoint: p.Mountpoint,
			Fstype:     p.Fstype,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			Percent:    usage.UsedPercent,
			Label:      cleanLabel,
			IsUSB:      isUSB,
		})
	}

	return devices, nil
}

// ListDirectory membaca isi direktori target dengan proteksi path sanitization dan resolusi lintas OS
func ListDirectory(targetPath string) ([]FileItem, error) {
	normPath := NormalizeStoragePath(targetPath)
	cleanPath := filepath.FromSlash(normPath)

	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		return nil, err
	}

	// Inisialisasi slice kosong non-nil agar direktori kosong menghasilkan JSON [] bukan null
	items := make([]FileItem, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		isDir := entry.IsDir()
		fullPath := filepath.Join(cleanPath, entry.Name())
		if !isDir && (entry.Type()&os.ModeSymlink != 0 || !strings.Contains(entry.Name(), ".")) {
			if stat, err := os.Stat(fullPath); err == nil {
				isDir = stat.IsDir()
			} else if entry.Type()&os.ModeSymlink != 0 {
				isDir = true
			}
		}

		ext := ""
		if !isDir {
			ext = strings.ToLower(filepath.Ext(entry.Name()))
		}

		// Gunakan forward slash seragam agar tidak terjadi konflik path di browser
		itemPath := filepath.ToSlash(fullPath)

		items = append(items, FileItem{
			Name:      entry.Name(),
			Path:      itemPath,
			IsDir:     isDir,
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Extension: ext,
		})
	}
	return items, nil
}

// CreateFolder membuat folder baru dengan normalisasi path lintas OS
func CreateFolder(parentPath, folderName string) error {
	normParent := NormalizeStoragePath(parentPath)
	fullPath := filepath.Join(filepath.FromSlash(normParent), folderName)
	return os.MkdirAll(fullPath, 0755)
}

// CreateFile membuat berkas teks atau skrip baru di direktori tertentu
func CreateFile(parentPath, fileName, content string) (FileItem, error) {
	normParent := NormalizeStoragePath(parentPath)
	cleanDir := filepath.FromSlash(normParent)
	if err := os.MkdirAll(cleanDir, 0755); err != nil {
		return FileItem{}, err
	}
	fullPath := filepath.Join(cleanDir, fileName)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return FileItem{}, err
	}
	stat, err := os.Stat(fullPath)
	if err != nil {
		return FileItem{}, err
	}
	return FileItem{
		Name:      fileName,
		Path:      filepath.ToSlash(fullPath),
		IsDir:     false,
		Size:      stat.Size(),
		ModTime:   stat.ModTime(),
		Extension: strings.ToLower(filepath.Ext(fileName)),
	}, nil
}

// DeleteFileOrFolder menghapus berkas atau folder secara rekursif
func DeleteFileOrFolder(targetPath string) error {
	normPath := NormalizeStoragePath(targetPath)
	return os.RemoveAll(filepath.FromSlash(normPath))
}

// RenameFileOrFolder mengubah nama atau memindahkan berkas
func RenameFileOrFolder(oldPath, newPath string) error {
	normOld := NormalizeStoragePath(oldPath)
	normNew := NormalizeStoragePath(newPath)
	return os.Rename(filepath.FromSlash(normOld), filepath.FromSlash(normNew))
}

// ReadFileContent membaca konten berkas teks (misal: log, config, json, yaml, env, txt, code)
func ReadFileContent(targetPath string) (string, error) {
	normPath := NormalizeStoragePath(targetPath)
	cleanPath := filepath.FromSlash(normPath)

	stat, err := os.Stat(cleanPath)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", fmt.Errorf("target adalah direktori, bukan berkas")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", err
	}
	// Batasi pembacaan maksimal 500 KB agar hemat RAM pada STB ARM64
	if len(data) > 500*1024 {
		return string(data[:500*1024]) + "\n... [Konten dipotong: berkas terlalu besar] ...", nil
	}
	return string(data), nil
}

// SaveFileContent menyimpan editan berkas teks ke disk
func SaveFileContent(targetPath, content string) error {
	normPath := NormalizeStoragePath(targetPath)
	cleanPath := filepath.FromSlash(normPath)
	_ = os.MkdirAll(filepath.Dir(cleanPath), 0755)
	return os.WriteFile(cleanPath, []byte(content), 0644)
}

// ExtractZip mengekstrak berkas arsip zip ke direktori tujuan secara aman (proteksi zip-slip)
func ExtractZip(srcZip, destDir string) error {
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		// Proteksi Zip Slip vulnerability
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(destDir)) {
			continue
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
