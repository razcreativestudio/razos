/**
 * RAZOS Over-The-Air (OTA) Updater Engine
 * Updated: 2026-10-08
 * Deskripsi: Modul pengecekan rilis versi baru ke GitHub Releases API
 * dan penggantian binary statis mandiri tanpa mengganggu service runtime host.
 */

package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

// CurrentVersion versi saat ini dari RAZOS Core
const CurrentVersion = "1.5.0"

// ReleaseInfo struktur response rilis dari GitHub Releases API
type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateStatus ringkasan status update untuk frontend UI
type UpdateStatus struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	HasUpdate       bool   `json:"has_update"`
	ReleaseNotes    string `json:"release_notes"`
	ReleaseDate     string `json:"release_date"`
	DownloadURL     string `json:"download_url"`
	Architecture    string `json:"architecture"`
}

// CheckUpdate memeriksa versi terbaru dari repositori GitHub resmi
func CheckUpdate() (*UpdateStatus, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/razcreative/razos/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "RAZOS-Updater/"+CurrentVersion)

	res, err := client.Do(req)
	if err != nil {
		// Fallback status jika STB dalam mode offline / tanpa internet
		return &UpdateStatus{
			CurrentVersion: CurrentVersion,
			LatestVersion:  CurrentVersion,
			HasUpdate:      false,
			ReleaseNotes:   "Sistem dalam versi terbaru (atau offline).",
			Architecture:   runtime.GOARCH,
		}, nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &UpdateStatus{
			CurrentVersion: CurrentVersion,
			LatestVersion:  CurrentVersion,
			HasUpdate:      false,
			ReleaseNotes:   "Tidak ada rilis baru yang tersedia saat ini.",
			Architecture:   runtime.GOARCH,
		}, nil
	}

	var rel ReleaseInfo
	if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
		return nil, err
	}

	targetBinary := fmt.Sprintf("razos-linux-%s", runtime.GOARCH)
	downloadURL := ""
	for _, a := range rel.Assets {
		if a.Name == targetBinary {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}

	hasUpdate := rel.TagName != "" && rel.TagName != CurrentVersion && rel.TagName != ("v"+CurrentVersion)

	return &UpdateStatus{
		CurrentVersion: CurrentVersion,
		LatestVersion:  rel.TagName,
		HasUpdate:      hasUpdate,
		ReleaseNotes:   rel.Body,
		ReleaseDate:    rel.PublishedAt,
		DownloadURL:    downloadURL,
		Architecture:   runtime.GOARCH,
	}, nil
}

// PerformSelfUpdate mengunduh binary rilis baru dan menggantikan binary lama di /DATA/razos/bin/razos
func PerformSelfUpdate(downloadURL string) error {
	if downloadURL == "" {
		return fmt.Errorf("URL binary unduhan tidak valid")
	}

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal mengunduh binary, status: %d", res.StatusCode)
	}

	// Buat file binary sementara
	tmpPath := "/DATA/razos/bin/razos.tmp"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}

	// Ganti file binary aktif secara atomik
	finalPath := "/DATA/razos/bin/razos"
	return os.Rename(tmpPath, finalPath)
}
