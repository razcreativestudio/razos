/**
 * RAZOS App Store Engine
 * Updated: 2026-10-08
 * Deskripsi: Engine deploy dan uninstall container Docker ke partisi /DATA/apps/
 * dengan validasi memori RAM minimum (> 150 MB free) dan opsi Clean Uninstall.
 */

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/shirou/gopsutil/v3/mem"
)

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
}

// DeployRequest payload instalasi aplikasi dengan custom port/volume
type DeployRequest struct {
	AppID     string   `json:"app_id"`
	Port      int      `json:"port"`
	DataPath  string   `json:"data_path"`
	ExtraEnvs []string `json:"extra_envs"`
}

// GetCatalog membaca katalog aplikasi dan mencocokkan status instalasi container
func GetCatalog() ([]AppManifest, error) {
	data, err := os.ReadFile("internal/store/RAZApps.json")
	if err != nil {
		// Fallback jika dieksekusi dari direktori backend/
		data, err = os.ReadFile("backend/internal/store/RAZApps.json")
		if err != nil {
			return nil, err
		}
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

	// Tentukan folder data di /DATA/apps/
	dataPath := req.DataPath
	if dataPath == "" {
		dataPath = filepath.Join("/DATA/apps", targetApp.Slug)
	}
	_ = os.MkdirAll(dataPath, 0755)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()

	// Setup Port Binding
	portStr := fmt.Sprintf("%d/tcp", req.Port)
	if req.Port <= 0 {
		portStr = fmt.Sprintf("%d/tcp", targetApp.DefaultPort)
	}
	natPort, _ := nat.NewPort("tcp", fmt.Sprintf("%d", targetApp.DefaultPort))

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			natPort: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", req.Port)},
			},
		},
		Binds: []string{
			fmt.Sprintf("%s:%s", dataPath, targetApp.DefaultVolume),
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		// Batasi memori per container di STB agar tidak memicu OOM Killer
		Resources: container.Resources{
			Memory: 256 * 1024 * 1024, // Limit 256 MB per container
		},
	}

	containerName := "razos-" + targetApp.Slug
	// Hapus jika ada container lama dengan nama sama
	_ = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true})

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: targetApp.Image,
		Env:   append(targetApp.Env, req.ExtraEnvs...),
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}

// CleanUninstall menghapus container dan opsi pembersihan direktori data /DATA/apps/
func CleanUninstall(appSlug string, purgeData bool) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	containerName := "razos-" + appSlug
	ctx := context.Background()

	// Stop dan remove container
	_ = cli.ContainerStop(ctx, containerName, nil)
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})

	// Jika opsi purgeData aktif, musnahkan folder data di /DATA/apps/<slug>
	if purgeData {
		folderPath := filepath.Join("/DATA/apps", appSlug)
		_ = os.RemoveAll(folderPath)
	}

	return err
}
