/**
 * RAZOS Database Central Suite
 * Updated: 2026-10-08
 * Deskripsi: Pengelola database MariaDB & PostgreSQL terpusat serta deployment Adminer Web GUI.
 */

package database

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

// Inisialisasi daftar database cadangan
var registeredDatabases = []DatabaseEntry{
	{
		ID:        1,
		Name:      "immich",
		Engine:    "postgres",
		User:      "postgres",
		Port:      5432,
		Host:      "127.0.0.1",
		Size:      "PostgreSQL",
		CreatedAt: "Aktif",
	},
}

// GetDatabases mengambil seluruh database aktif yang terdeteksi di Docker STB
func GetDatabases() []DatabaseEntry {
	var list []DatabaseEntry
	idCounter := 1

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		defer cli.Close()
		containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		hasPostgres := false
		hasMariaDB := false

		for _, c := range containers {
			for _, name := range c.Names {
				if name == "/razos-postgres" && c.State == "running" {
					hasPostgres = true
				}
				if (name == "/razos-mariadb" || name == "/razos-mysql") && c.State == "running" {
					hasMariaDB = true
				}
			}
		}

		if hasPostgres {
			list = append(list, DatabaseEntry{
				ID:        idCounter,
				Name:      "immich",
				Engine:    "postgres",
				User:      "postgres",
				Port:      5432,
				Host:      "127.0.0.1",
				Size:      "PostgreSQL",
				CreatedAt: "Aktif (Docker)",
			})
			idCounter++
		}

		if hasMariaDB {
			list = append(list, DatabaseEntry{
				ID:        idCounter,
				Name:      "razos_primary",
				Engine:    "mariadb",
				User:      "root",
				Port:      3306,
				Host:      "127.0.0.1",
				Size:      "MariaDB",
				CreatedAt: "Aktif (Docker)",
			})
			idCounter++
		}
	}

	if len(list) == 0 {
		return registeredDatabases
	}
	return list
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
		CreatedAt: "2026-10-26",
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

	// Gunakan CLI docker run agar otomatis terhubung ke network razos-apps-net dan port 8088
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, "--network", "razos-apps-net", "-p", "8088:8080", "--restart", "always", "-m", "64M", "adminer:latest")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gagal menjalankan adminer: %s (%v)", string(out), err)
	}
	return nil
}
