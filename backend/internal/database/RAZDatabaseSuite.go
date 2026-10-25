/**
 * RAZOS Database Central Suite
 * Updated: 2026-10-08
 * Deskripsi: Pengelola database MariaDB & PostgreSQL terpusat serta deployment Adminer Web GUI.
 */

package database

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// DatabaseEntry menyimpan metadata database yang terdaftar
type DatabaseEntry struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Engine    string `json:"engine"` // 'mariadb' | 'postgres'
	User      string `json:"user"`
	Port      int    `json:"port"`
	Host      string `json:"host"`
	Size      string `json:"size"`
	CreatedAt string `json:"created_at"`
}

// Inisialisasi daftar database dalam memori/SQLite
var registeredDatabases = []DatabaseEntry{
	{
		ID:        1,
		Name:      "razos_primary",
		Engine:    "mariadb",
		User:      "raz_admin",
		Port:      3306,
		Host:      "127.0.0.1",
		Size:      "12.4 MB",
		CreatedAt: "2026-10-08",
	},
	{
		ID:        2,
		Name:      "razos_pg_master",
		Engine:    "postgres",
		User:      "postgres",
		Port:      5432,
		Host:      "127.0.0.1",
		Size:      "28.1 MB",
		CreatedAt: "2026-10-08",
	},
}

// GetDatabases mengambil seluruh database terdaftar
func GetDatabases() []DatabaseEntry {
	return registeredDatabases
}

// CreateDatabaseRequest payload pembuatan database baru
type CreateDatabaseRequest struct {
	Name     string `json:"name"`
	Engine   string `json:"engine"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// CreateDatabase mendaftarkan instance database baru
func CreateDatabase(req CreateDatabaseRequest) (*DatabaseEntry, error) {
	port := 3306
	if req.Engine == "postgres" {
		port = 5432
	}

	newDB := DatabaseEntry{
		ID:        len(registeredDatabases) + 1,
		Name:      req.Name,
		Engine:    req.Engine,
		User:      req.User,
		Port:      port,
		Host:      "127.0.0.1",
		Size:      "0.5 MB",
		CreatedAt: "2026-10-08",
	}

	registeredDatabases = append(registeredDatabases, newDB)
	return &newDB, nil
}

// EnsureAdminerContainer memastikan container Adminer aktif pada port 8088
func EnsureAdminerContainer() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()
	containerName := "razos-adminer"

	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				if c.State != "running" {
					return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
				}
				return nil
			}
		}
	}

	natPort, _ := nat.NewPort("tcp", "8080")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "adminer:latest",
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			natPort: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "8088"},
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Resources: container.Resources{
			Memory: 64 * 1024 * 1024, // Adminer ultra ringan 64 MB RAM
		},
	}, nil, nil, containerName)

	if err != nil {
		return err
	}
	return cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}
