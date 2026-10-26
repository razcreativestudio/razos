/**
 * RAZOS App Store Engine
 * Updated: 2026-10-26 (v3.9.4)
 * Deskripsi: Engine deploy dan uninstall container Docker ke partisi /DATA/apps/
 * dengan validasi memori RAM minimum (> 150 MB free), multi-layer Docker CLI + Go SDK
 * pull engine dengan real-time terminal stdout scanner logging, robust container port exposure & clean uninstall.
 */

package store

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/shirou/gopsutil/v3/mem"
)

// LogItem struktur satu baris catatan aktivitas instalasi di terminal
type LogItem struct {
	Time string `json:"time"`
	Type string `json:"type"` // "info" | "docker" | "pull" | "success" | "error"
	Tag  string `json:"tag"`  // "[INFO]", "[DOCKER]", "[PULL]", etc.
	Text string `json:"text"`
}

var (
	installLogsMu sync.RWMutex
	installLogs   = make(map[string][]LogItem) // slug -> logs
)

// AddInstallLog mencatat satu baris aktivitas instalasi secara thread-safe
func AddInstallLog(slug, logType, tag, text string) {
	installLogsMu.Lock()
	defer installLogsMu.Unlock()
	item := LogItem{
		Time: time.Now().Format("15:04:05"),
		Type: logType,
		Tag:  tag,
		Text: text,
	}
	installLogs[slug] = append(installLogs[slug], item)
	// Batasi maksimal 250 baris per slug agar hemat RAM STB
	if len(installLogs[slug]) > 250 {
		installLogs[slug] = installLogs[slug][len(installLogs[slug])-250:]
	}
}

// GetInstallLogs mengambil riwayat log instalasi aplikasi tertentu
func GetInstallLogs(slug string) []LogItem {
	installLogsMu.RLock()
	defer installLogsMu.RUnlock()
	items, exists := installLogs[slug]
	if !exists {
		return []LogItem{}
	}
	copied := make([]LogItem, len(items))
	copy(copied, items)
	return copied
}

// ClearInstallLogs membersihkan riwayat log untuk melegakan memori
func ClearInstallLogs(slug string) {
	installLogsMu.Lock()
	defer installLogsMu.Unlock()
	delete(installLogs, slug)
}

//go:embed RAZApps.json
var embeddedAppsJSON []byte

// AppManifest struktur skema aplikasi di App Store
type AppManifest struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Category      string   `json:"category"`
	Description   string   `json:"description"`
	Image         string   `json:"image"`
	DefaultPort   int      `json:"default_port"`
	Port          int      `json:"port,omitempty"`
	DefaultVolume string   `json:"default_volume"`
	Env           []string `json:"env"`
	Installed     bool     `json:"installed"`
	Readme        string   `json:"readme,omitempty"`
	Icon          string   `json:"icon,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
}

// DeployRequest payload instalasi aplikasi dengan custom port/volume atau manual custom image
type DeployRequest struct {
	AppID         string   `json:"app_id"`
	Name          string   `json:"name"`
	Image         string   `json:"image,omitempty"`
	CustomImage   string   `json:"custom_image"`
	ContainerName string   `json:"container_name"`
	Port          int      `json:"port"`
	DataPath      string   `json:"data_path"`
	ExtraEnvs     []string `json:"extra_envs"`
}

// GetCatalog membaca katalog aplikasi dan mencocokkan status instalasi container
func GetCatalog() ([]AppManifest, error) {
	candidatePaths := []string{
		"/DATA/razos/src/backend/internal/store/RAZApps.json",
		"/DATA/razos/store/RAZApps.json",
		"backend/internal/store/RAZApps.json",
		"internal/store/RAZApps.json",
	}

	var data []byte
	for _, p := range candidatePaths {
		if content, err := os.ReadFile(p); err == nil && len(content) > 0 {
			data = content
			break
		}
	}

	if len(data) == 0 {
		data = embeddedAppsJSON
	}

	var catalog []AppManifest
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}

	// Cek status ke Docker Engine
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		defer cli.Close()
		containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		for i := range catalog {
			for _, c := range containers {
				for _, name := range c.Names {
					if name == "/razos-"+catalog[i].Slug {
						catalog[i].Installed = true
						// Ekstraksi port publik yang sedang dipetakan oleh container Docker
						for _, p := range c.Ports {
							if p.PublicPort > 0 {
								catalog[i].Port = int(p.PublicPort)
								break
							}
						}
						break
					}
				}
			}
			// Jika port belum terisi atau 0, gunakan DefaultPort sebagai fallback dinamis
			if catalog[i].Port == 0 {
				catalog[i].Port = catalog[i].DefaultPort
			}
		}
	} else {
		for i := range catalog {
			if catalog[i].Port == 0 {
				catalog[i].Port = catalog[i].DefaultPort
			}
		}
	}

	return catalog, nil
}

// DeployApp melakukan pull image dan menjalankan container ke /DATA/apps/ dengan live terminal logging
func DeployApp(req DeployRequest) error {
	var targetImage string
	var containerName string
	var defaultVol string
	var defaultPort int
	var envs []string
	var slug string

	// Fallback jika req.CustomImage kosong tetapi req.Image diberikan
	if req.CustomImage == "" && req.Image != "" {
		req.CustomImage = req.Image
	}

	// Jika instalasi manual kustom (Custom Docker Image ala CasaOS)
	if req.CustomImage != "" {
		targetImage = req.CustomImage
		slug = req.AppID
		if slug == "" {
			slug = req.ContainerName
		}
		if slug == "" {
			slug = "custom-app"
		}
		containerName = "razos-" + slug
		defaultVol = req.DataPath
		if defaultVol == "" {
			defaultVol = filepath.Join("/DATA/apps", slug)
		}
		defaultPort = req.Port
		if defaultPort <= 0 {
			defaultPort = 8080
		}
		envs = req.ExtraEnvs
	} else {
		// Instalasi dari katalog
		catalog, err := GetCatalog()
		if err != nil {
			return err
		}

		var targetApp *AppManifest
		for _, a := range catalog {
			if a.ID == req.AppID {
				targetApp = &a
				break
			}
		}
		if targetApp == nil {
			return fmt.Errorf("aplikasi dengan ID %s tidak ditemukan", req.AppID)
		}

		targetImage = targetApp.Image
		slug = targetApp.Slug
		containerName = "razos-" + targetApp.Slug
		defaultVol = targetApp.DefaultVolume
		defaultPort = targetApp.DefaultPort
		envs = append(targetApp.Env, req.ExtraEnvs...)
	}

	// Inisialisasi terminal log untuk slug ini
	ClearInstallLogs(slug)
	AddInstallLog(slug, "info", "[INIT]", fmt.Sprintf("Memulai instalasi %s (%s)...", slug, targetImage))

	// 1. Validasi Keamanan Hardware STB: Cek kapasitas RAM bebas minimal 150 MB
	vMem, err := mem.VirtualMemory()
	if err == nil {
		freeMB := vMem.Available / 1024 / 1024
		AddInstallLog(slug, "info", "[SYS]", fmt.Sprintf("Verifikasi RAM: %d MB bebas tersedia (Batas aman: 150 MB)", freeMB))
		if vMem.Available < 150*1024*1024 {
			AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("RAM bebas tidak mencukupi: hanya tersisa %d MB", freeMB))
			return fmt.Errorf("sumber daya tidak mencukupi: RAM bebas tersisa kurang dari 150 MB")
		}
	}

	// 2. Tentukan dan siapkan folder data di /DATA/apps/
	dataPath := req.DataPath
	if dataPath == "" {
		dataPath = filepath.Join("/DATA/apps", slug)
	}
	AddInstallLog(slug, "info", "[STORAGE]", fmt.Sprintf("Menyiapkan volume penyimpanan di %s...", dataPath))
	if mkErr := os.MkdirAll(dataPath, 0755); mkErr != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("Gagal membuat direktori %s: %v", dataPath, mkErr))
		return mkErr
	}
	AddInstallLog(slug, "success", "[STORAGE]", fmt.Sprintf("Direktori %s siap (izin 0755)", dataPath))

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("Koneksi Docker socket gagal: %v", err))
		return err
	}
	defer cli.Close()

	// 3. Setup Port Binding & Volume Mount Target
	targetPort := req.Port
	if targetPort <= 0 {
		targetPort = defaultPort
	}
	// Cegah bentrok dengan port inti RAZOS (80 & 8080)
	if targetPort == 80 || targetPort == 8080 {
		targetPort = defaultPort
		if targetPort == 80 || targetPort == 8080 {
			targetPort = 8090
		}
	}
	natPort, _ := nat.NewPort("tcp", fmt.Sprintf("%d", defaultPort))

	mountTarget := defaultVol
	if mountTarget == "" || strings.HasPrefix(mountTarget, "/DATA/") || strings.HasPrefix(mountTarget, "/data/") {
		mountTarget = "/data"
	}
	AddInstallLog(slug, "info", "[NET]", fmt.Sprintf("Alokasi Port Host :%d -> Container :%d", targetPort, defaultPort))
	AddInstallLog(slug, "info", "[MOUNT]", fmt.Sprintf("Binding Volume %s -> %s", dataPath, mountTarget))

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			natPort: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", targetPort)},
			},
		},
		Binds: []string{
			fmt.Sprintf("%s:%s", dataPath, mountTarget),
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		// Batasi memori per container di STB agar tidak memicu OOM Killer
		Resources: container.Resources{
			Memory: 256 * 1024 * 1024, // Limit 256 MB per container
		},
	}

	// 4. Unduh Docker image ke host secara tangguh (dengan timeout 10 menit untuk image besar)
	pullCtx, pullCancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer pullCancel()

	if err := pullDockerImage(pullCtx, cli, targetImage, slug); err != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("Gagal mengunduh image: %v", err))
		return err
	}

	// 5. Hapus jika ada container lama dengan nama sama
	AddInstallLog(slug, "docker", "[CLEAN]", fmt.Sprintf("Memeriksa & membersihkan container lama '%s'...", containerName))
	_ = cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{Force: true})

	// 6. Buat container baru dengan ExposedPorts agar binding port aktif sempurna
	AddInstallLog(slug, "docker", "[CREATE]", fmt.Sprintf("Membuat container '%s' dari image %s...", containerName, targetImage))
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        targetImage,
		Env:          envs,
		ExposedPorts: nat.PortSet{natPort: struct{}{}},
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("Gagal membuat container: %v", err))
		return fmt.Errorf("gagal membuat container %s: %v", containerName, err)
	}

	// 7. Jalankan container
	AddInstallLog(slug, "docker", "[START]", fmt.Sprintf("Menjalankan container (ID: %s)...", resp.ID[:12]))
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("Gagal start container: %v", err))
		return fmt.Errorf("gagal menjalankan container %s: %v", containerName, err)
	}

	AddInstallLog(slug, "success", "[ONLINE]", fmt.Sprintf("Container '%s' aktif dan siap diakses di port :%d!", containerName, targetPort))
	return nil
}

// pullDockerImage mengunduh Docker image ke host dengan live terminal line scanner
func pullDockerImage(ctx context.Context, cli *client.Client, imageName, slug string) error {
	// 1. Cek apakah image sudah tersedia secara lokal di host
	_, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		AddInstallLog(slug, "success", "[CACHE]", fmt.Sprintf("Image %s sudah tersedia di lokal host (skip download)", imageName))
		return nil
	}

	log.Printf("[RAZAppStore] Mengunduh Docker image: %s ...", imageName)
	AddInstallLog(slug, "docker", "[PULL]", fmt.Sprintf("docker pull %s", imageName))

	// 2. Prioritas 1: Gunakan Docker CLI bawaan sistem (/usr/bin/docker pull) dengan pipe scanner
	cmd := exec.CommandContext(ctx, "docker", "pull", imageName)
	stdout, outErr := cmd.StdoutPipe()
	stderr, errPipeErr := cmd.StderrPipe()

	if outErr == nil && errPipeErr == nil {
		multiReader := io.MultiReader(stdout, stderr)
		if startErr := cmd.Start(); startErr == nil {
			scanner := bufio.NewScanner(multiReader)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					AddInstallLog(slug, "pull", "[PULL]", line)
				}
			}
			if waitErr := cmd.Wait(); waitErr == nil {
				AddInstallLog(slug, "success", "[PULL]", fmt.Sprintf("Selesai mengunduh seluruh layer %s", imageName))
				return nil
			}
		}
	}

	AddInstallLog(slug, "docker", "[FALLBACK]", "Mencoba fallback unduh via Docker Go SDK...")

	// 3. Fallback: Gunakan Docker Go SDK jika CLI gagal atau terkendala
	out, sdkErr := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if sdkErr != nil {
		AddInstallLog(slug, "error", "[ERR]", fmt.Sprintf("SDK ImagePull gagal: %v", sdkErr))
		return fmt.Errorf("gagal mengunduh docker image %s: %v", imageName, sdkErr)
	}
	defer out.Close()

	// Parse JSON event stream dari Docker daemon untuk menangkap status & progress layer
	type pullEvent struct {
		Status         string `json:"status"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int64 `json:"current"`
			Total   int64 `json:"total"`
		} `json:"progressDetail"`
		Error       string `json:"error"`
		ErrorDetail struct {
			Message string `json:"message"`
		} `json:"errorDetail"`
	}

	dec := json.NewDecoder(out)
	var streamErr string
	for {
		var ev pullEvent
		if err := dec.Decode(&ev); err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		if ev.Status != "" {
			msg := ev.Status
			if ev.Progress != "" {
				msg += " " + ev.Progress
			}
			AddInstallLog(slug, "pull", "[PULL]", msg)
		}
		if ev.Error != "" {
			streamErr = ev.Error
			AddInstallLog(slug, "error", "[ERR]", ev.Error)
		} else if ev.ErrorDetail.Message != "" {
			streamErr = ev.ErrorDetail.Message
			AddInstallLog(slug, "error", "[ERR]", ev.ErrorDetail.Message)
		}
	}

	if streamErr != "" {
		return fmt.Errorf("gagal mengunduh docker image %s: %s", imageName, streamErr)
	}

	// 4. Verifikasi final bahwa image benar-benar terdaftar di daemon
	_, _, verifyErr := cli.ImageInspectWithRaw(ctx, imageName)
	if verifyErr != nil {
		return fmt.Errorf("image %s tidak ditemukan di daemon setelah proses pull selesai", imageName)
	}

	AddInstallLog(slug, "success", "[PULL]", fmt.Sprintf("Verifikasi image %s di daemon Docker berhasil!", imageName))
	return nil
}

// CleanUninstall menghapus container dan opsi pembersihan direktori data /DATA/apps/ dengan live logging
func CleanUninstall(appSlug string, purgeData bool) error {
	ClearInstallLogs(appSlug)
	AddInstallLog(appSlug, "info", "[INIT]", fmt.Sprintf("Inisialisasi proses uninstall aplikasi '%s'...", appSlug))

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		AddInstallLog(appSlug, "error", "[ERR]", fmt.Sprintf("Gagal inisialisasi Docker client: %v", err))
		return err
	}
	defer cli.Close()

	containerName := "razos-" + appSlug
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Menghentikan container (timeout 5 detik)
	AddInstallLog(appSlug, "docker", "[STOP]", fmt.Sprintf("Mengirim sinyal stop ke container '%s'...", containerName))
	AddInstallLog(appSlug, "docker", "[CMD]", fmt.Sprintf("docker stop -t 5 %s", containerName))
	stopTimeout := 5
	_ = cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &stopTimeout})
	AddInstallLog(appSlug, "success", "[STOP]", fmt.Sprintf("Container '%s' telah berhenti", containerName))

	// 2. Menghapus container dan volumenya
	AddInstallLog(appSlug, "docker", "[REMOVE]", fmt.Sprintf("Menghapus container dan volume: docker rm -v -f %s", containerName))
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if client.IsErrNotFound(err) || strings.Contains(errMsg, "no such") || strings.Contains(errMsg, "not found") {
			err = nil
			AddInstallLog(appSlug, "info", "[CLEAN]", fmt.Sprintf("Container %s sudah tidak ada di sistem (bersih)", containerName))
		} else {
			AddInstallLog(appSlug, "error", "[ERR]", fmt.Sprintf("Peringatan remove container: %v", err))
		}
	} else {
		AddInstallLog(appSlug, "success", "[REMOVE]", fmt.Sprintf("Container '%s' berhasil dihapus dari Docker daemon", containerName))
	}

	// 3. Jika opsi purgeData aktif, musnahkan folder data di /DATA/apps/<slug>
	if purgeData {
		folderPath := filepath.Join("/DATA/apps", appSlug)
		AddInstallLog(appSlug, "info", "[STORAGE]", fmt.Sprintf("Membersihkan direktori permanen: rm -rf %s", folderPath))
		_ = os.RemoveAll(folderPath)
		AddInstallLog(appSlug, "success", "[STORAGE]", fmt.Sprintf("Direktori %s berhasil dibersihkan tuntas", folderPath))
	} else {
		AddInstallLog(appSlug, "info", "[STORAGE]", "Opsi simpan data aktif: direktori /DATA/apps/ dipertahankan")
	}

	AddInstallLog(appSlug, "success", "[DONE]", fmt.Sprintf("Proses uninstall '%s' selesai 100%%. Memori dan resource sistem telah dilepaskan.", appSlug))
	return err
}
