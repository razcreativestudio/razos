/**
 * RAZOS Advanced Storage & File Manager Engine
 * Updated: 2026-10-08
 * Deskripsi: Engine backend untuk membaca partisi disk/HDD/USB fisik (gopsutil/disk),
 * penjelajahan direktori, pembuatan berkas/folder, upload, download, rename, dan hapus berkas.
 */

package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// StorageDevice mendeskripsikan disk fisik atau partisi yang terpasang pada STB
type StorageDevice struct {
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

		label := p.Mountpoint
		isUSB := false

		if p.Mountpoint == "/DATA" || p.Mountpoint == "C:\\DATA" {
			label = "Primary Storage (/DATA - MicroSD)"
		} else if p.Mountpoint == "/" || strings.HasPrefix(strings.ToUpper(p.Mountpoint), "C:") {
			label = "System Storage (/ - eMMC)"
		} else if strings.Contains(strings.ToLower(p.Mountpoint), "media") || strings.Contains(strings.ToLower(p.Mountpoint), "usb") || strings.HasPrefix(p.Mountpoint, "/mnt") {
			isUSB = true
			label = "External Storage (" + filepath.Base(p.Mountpoint) + ")"
		}

		devices = append(devices, StorageDevice{
			Device:     p.Device,
			Mountpoint: p.Mountpoint,
			Fstype:     p.Fstype,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			Percent:    usage.UsedPercent,
			Label:      label,
			IsUSB:      isUSB,
		})
	}

	return devices, nil
}

// ListDirectory membaca isi direktori target dengan proteksi path sanitization
func ListDirectory(targetPath string) ([]FileItem, error) {
	// Normalisasi path jika frontend Windows dev mengirimkan C:/DATA atau C:\DATA saat backend berjalan di Linux
	if runtime.GOOS != "windows" {
		tpLower := strings.ToLower(filepath.ToSlash(targetPath))
		if strings.HasPrefix(tpLower, "c:/data") {
			targetPath = "/DATA" + targetPath[7:]
		} else if strings.HasPrefix(tpLower, "c:/") {
			targetPath = "/" + targetPath[3:]
		} else if tpLower == "c:" {
			targetPath = "/"
		}
	}

	if targetPath == "" || targetPath == "." || targetPath == "/" || targetPath == "/DATA" {
		if _, err := os.Stat("/DATA"); err == nil {
			targetPath = "/DATA"
		} else if _, err := os.Stat("C:\\DATA"); err == nil {
			targetPath = "C:\\DATA"
		} else {
			targetPath = "."
		}
	}

	// Normalisasi Windows root drive letter misal "C:" menjadi "C:\"
	if len(targetPath) == 2 && targetPath[1] == ':' {
		targetPath += "\\"
	}

	cleanPath := filepath.Clean(targetPath)
	// Resolusi Windows Junction dan symbolic links secara transparan
	if evalPath, evalErr := filepath.EvalSymlinks(cleanPath); evalErr == nil {
		cleanPath = evalPath
	}
	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		return nil, err
	}

	var items []FileItem
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

// CreateFolder membuat folder baru
func CreateFolder(parentPath, folderName string) error {
	fullPath := filepath.Join(parentPath, folderName)
	return os.MkdirAll(fullPath, 0755)
}

// DeleteFileOrFolder menghapus berkas atau folder secara rekursif
func DeleteFileOrFolder(targetPath string) error {
	return os.RemoveAll(targetPath)
}

// RenameFileOrFolder mengubah nama atau memindahkan berkas
func RenameFileOrFolder(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// ReadFileContent membaca konten berkas teks (misal: log, config, json, yaml, env)
func ReadFileContent(targetPath string) (string, error) {
	stat, err := os.Stat(targetPath)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", fmt.Errorf("target adalah direktori, bukan berkas")
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		return "", err
	}
	// Batasi pembacaan maksimal 500 KB agar hemat RAM
	if len(data) > 500*1024 {
		return string(data[:500*1024]) + "\n... [Konten dipotong: berkas terlalu besar] ...", nil
	}
	return string(data), nil
}

// SaveFileContent menyimpan editan berkas teks
func SaveFileContent(targetPath, content string) error {
	return os.WriteFile(targetPath, []byte(content), 0644)
}
