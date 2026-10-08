/**
 * RAZOS Advanced Storage & File Manager Engine
 * Updated: 2026-10-08
 * Deskripsi: Engine backend untuk membaca partisi disk/HDD/USB fisik (gopsutil/disk),
 * penjelajahan direktori, pembuatan berkas/folder, upload, download, rename, dan hapus berkas.
 */

package filemanager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// GetStorageDevices mendeteksi seluruh perangkat penyimpanan (eMMC, MicroSD, HDD eksternal, USB)
func GetStorageDevices() ([]StorageDevice, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var devices []StorageDevice
	for _, p := range partitions {
		// Abaikan mountpoint loopback dan pseudofilesystems di Linux
		if strings.HasPrefix(p.Mountpoint, "/sys") ||
			strings.HasPrefix(p.Mountpoint, "/proc") ||
			strings.HasPrefix(p.Mountpoint, "/dev") ||
			strings.HasPrefix(p.Mountpoint, "/run") {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		label := p.Mountpoint
		isUSB := false
		if strings.Contains(p.Mountpoint, "media") || strings.Contains(p.Mountpoint, "usb") || strings.Contains(p.Device, "sd") {
			isUSB = true
			label = "External Storage (" + p.Mountpoint + ")"
		} else if p.Mountpoint == "/DATA" {
			label = "MicroSD Primary (/DATA)"
		} else if p.Mountpoint == "/" {
			label = "System Root eMMC (/)"
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
	if targetPath == "" {
		targetPath = "/DATA"
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			targetPath = "."
		}
	}

	cleanPath := filepath.Clean(targetPath)
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

		ext := ""
		if !entry.IsDir() {
			ext = strings.ToLower(filepath.Ext(entry.Name()))
		}

		items = append(items, FileItem{
			Name:      entry.Name(),
			Path:      filepath.Join(cleanPath, entry.Name()),
			IsDir:     entry.IsDir(),
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
