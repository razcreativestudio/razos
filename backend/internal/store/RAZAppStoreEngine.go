/**
 * RAZOS App Store Engine
 * Updated: 2026-10-26 (v3.9.0)
 * Deskripsi: Engine deploy dan uninstall container Docker ke partisi /DATA/apps/
 * dengan validasi memori RAM minimum (> 150 MB free), multi-layer Docker CLI + Go SDK
 * pull engine dengan error streaming decoder, serta robust container port exposure & clean uninstall.
 */

package store

import (
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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/shirou/gopsutil/v3/mem"
)

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
						break
					}
				}
			}
		}
	}

	return catalog, nil
}

// DeployApp melakukan pull image dan menjalankan container ke /DATA/apps/
func DeployApp(req DeployRequest) error {
	// 1. Validasi Keamanan Hardware STB: Cek kapasitas RAM bebas minimal 150 MB
	vMem, err := mem.VirtualMemory()
	if err == nil && vMem.Available < 150*1024*1024 {
		return fmt.Errorf("sumber daya tidak mencukupi: RAM bebas tersisa kurang dari 150 MB")
	}

	var targetImage string
	var containerName string
	var defaultVol string
	var defaultPort int
	var envs []string
	var slug string

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

	// Tentukan folder data di /DATA/apps/
	dataPath := req.DataPath
	if dataPath == "" {
		dataPath = filepath.Join("/DATA/apps", slug)
	}
	_ = os.MkdirAll(dataPath, 0755)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Setup Port Binding
	targetPort := req.Port
	if targetPort <= 0 {
		targetPort = defaultPort
	}
	natPort, _ := nat.NewPort("tcp", fmt.Sprintf("%d", defaultPort))

	// Pastikan mount target volume di dalam container valid (/data jika format bawaan adalah path host /DATA/...)
	mountTarget := defaultVol
	if mountTarget == "" || strings.HasPrefix(mountTarget, "/DATA/") || strings.HasPrefix(mountTarget, "/data/") {
		mountTarget = "/data"
	}

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

	// Unduh Docker image ke host secara tangguh (dengan timeout 10 menit untuk image besar)
	pullCtx, pullCancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer pullCancel()

	if err := pullDockerImage(pullCtx, cli, targetImage); err != nil {
		return err
	}

	// Hapus jika ada container lama dengan nama sama
	_ = cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{Force: true})

	// Buat container baru dengan ExposedPorts agar binding port aktif sempurna
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        targetImage,
		Env:          envs,
		ExposedPorts: nat.PortSet{natPort: struct{}{}},
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("gagal membuat container %s: %v", containerName, err)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("gagal menjalankan container %s: %v", containerName, err)
	}

	return nil
}

// pullDockerImage mengunduh Docker image ke host dengan prioritas Docker CLI dan fallback Go SDK
func pullDockerImage(ctx context.Context, cli *client.Client, imageName string) error {
	// 1. Cek apakah image sudah tersedia secara lokal di host
	_, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		return nil
	}

	log.Printf("[RAZAppStore] Mengunduh Docker image: %s ...", imageName)

	// 2. Prioritas 1: Gunakan Docker CLI bawaan sistem (/usr/bin/docker pull)
	// Docker CLI memiliki penanganan layer, DNS, arsitektur CPU (ARM64), dan kredensial paling stabil di Linux
	cmd := exec.CommandContext(ctx, "docker", "pull", imageName)
	outBytes, cmdErr := cmd.CombinedOutput()
	if cmdErr == nil {
		log.Printf("[RAZAppStore] Sukses mengunduh image via Docker CLI: %s", imageName)
		return nil
	}

	log.Printf("[RAZAppStore] Docker CLI pull gagal (%v: %s), mencoba via Docker Go SDK...", cmdErr, strings.TrimSpace(string(outBytes)))

	// 3. Fallback: Gunakan Docker Go SDK jika CLI gagal atau terkendala
	out, sdkErr := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if sdkErr != nil {
		if len(outBytes) > 0 {
			return fmt.Errorf("gagal mengunduh docker image %s: %s", imageName, strings.TrimSpace(string(outBytes)))
		}
		return fmt.Errorf("gagal mengunduh docker image %s: %v", imageName, sdkErr)
	}
	defer out.Close()

	// Parse JSON event stream dari Docker daemon untuk menangkap pesan error asli
	type pullEvent struct {
		Status      string `json:"status"`
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
		if ev.Error != "" {
			streamErr = ev.Error
		} else if ev.ErrorDetail.Message != "" {
			streamErr = ev.ErrorDetail.Message
		}
	}

	if streamErr != "" {
		return fmt.Errorf("gagal mengunduh docker image %s: %s", imageName, streamErr)
	}

	// 4. Verifikasi final bahwa image benar-benar terdaftar di daemon
	_, _, verifyErr := cli.ImageInspectWithRaw(ctx, imageName)
	if verifyErr != nil {
		if len(outBytes) > 0 {
			return fmt.Errorf("gagal mengunduh docker image %s: %s", imageName, strings.TrimSpace(string(outBytes)))
		}
		return fmt.Errorf("image %s tidak ditemukan di daemon setelah proses pull selesai", imageName)
	}

	return nil
}

// CleanUninstall menghapus container dan opsi pembersihan direktori data /DATA/apps/
func CleanUninstall(appSlug string, purgeData bool) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	containerName := "razos-" + appSlug
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop container (timeout 5 detik)
	stopTimeout := 5
	_ = cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &stopTimeout})

	// Hapus container dan volumenya
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		// Abaikan error jika container memang belum ada / sudah terhapus
		if client.IsErrNotFound(err) || strings.Contains(errMsg, "no such") || strings.Contains(errMsg, "not found") {
			err = nil
		}
	}

	// Jika opsi purgeData aktif, musnahkan folder data di /DATA/apps/<slug>
	if purgeData {
		folderPath := filepath.Join("/DATA/apps", appSlug)
		_ = os.RemoveAll(folderPath)
	}

	return err
}
