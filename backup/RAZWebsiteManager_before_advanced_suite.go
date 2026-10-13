/**
 * RAZOS Website Management Engine
 * Updated: 2026-10-08
 * Deskripsi: Mesin pembuatan website dengan wizard otomatis, auto-assign port bebas (8100-8200),
 * pemilihan runtime (PHP, Node.js, Static HTML), dan koneksi database instan.
 */

package website

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/razos/razos/internal/proxy"
)

// Website merepresentasikan data situs web yang aktif
type Website struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	HostPort   int    `json:"port"`
	Runtime    string `json:"runtime"` // 'php' | 'nodejs' | 'static'
	DatabaseID *int   `json:"database_id"`
	DataPath   string `json:"data_path"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

var (
	siteMu sync.Mutex
	sites  = []Website{}
	nextID = 1
)

// GetWebsites mengambil daftar situs web yang terdaftar
func GetWebsites() []Website {
	siteMu.Lock()
	defer siteMu.Unlock()
	return sites
}

// FindFreePort mencari port bebas di rentang 8100-8200
func FindFreePort() (int, error) {
	for port := 8100; port <= 8200; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			_ = ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("semua alokasi port website (8100-8200) penuh")
}

// CreateWebsiteRequest payload pembuatan situs web baru
type CreateWebsiteRequest struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Runtime    string `json:"runtime"`
	DatabaseID *int   `json:"database_id"`
}

// CreateWebsite membuat situs web baru, membuat direktori /DATA/apps/<nama>, dan mendaftarkan rute reverse proxy
func CreateWebsite(req CreateWebsiteRequest) (*Website, error) {
	siteMu.Lock()
	defer siteMu.Unlock()

	port, err := FindFreePort()
	if err != nil {
		return nil, err
	}

	siteDir := filepath.Join("/DATA/apps/websites", req.Name)
	_ = os.MkdirAll(siteDir, 0755)

	newSite := Website{
		ID:         nextID,
		Name:       req.Name,
		Domain:     req.Domain,
		HostPort:   port,
		Runtime:    req.Runtime,
		DatabaseID: req.DatabaseID,
		DataPath:   siteDir,
		Status:     "Active",
		CreatedAt:  "2026-10-08",
	}
	nextID++

	// Daftarkan rute ke internal proxy
	proxy.GetRouter().AddRoute(req.Domain, port)

	sites = append(sites, newSite)
	return &newSite, nil
}

// DeleteWebsite menghapus situs web dan melepaskan rute proxy
func DeleteWebsite(id int) error {
	siteMu.Lock()
	defer siteMu.Unlock()

	var updated []Website
	var domainToRemove string
	for _, s := range sites {
		if s.ID == id {
			domainToRemove = s.Domain
		} else {
			updated = append(updated, s)
		}
	}

	if domainToRemove != "" {
		proxy.GetRouter().RemoveRoute(domainToRemove)
	}

	sites = updated
	return nil
}
