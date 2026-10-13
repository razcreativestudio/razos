/**
 * RAZOS Website Management Engine
 * Updated: 2026-10-13
 * Deskripsi: Mesin pengelolaan multi-website terisolasi berbasis Docker Container.
 * Mengatur root folder, alokasi port bebas (8100-8200), runtime stack (HTML Statis,
 * PHP dengan Database, Node.js, Next.js, Custom), integrasi Database Manager (Adminer),
 * resource limits (RAM & CPU), custom command, serta link icon untuk Desktop launcher.
 */

package website

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/razos/razos/internal/proxy"
)

// Website merepresentasikan data konfigurasi dan status sebuah situs web terisolasi
type Website struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Domain        string `json:"domain"`
	HostPort      int    `json:"port"`
	Runtime       string `json:"runtime"` // 'static' | 'php' | 'nodejs' | 'nextjs' | 'custom'
	ContainerName string `json:"container_name"`
	RootFolder    string `json:"root_folder"`
	Icon          string `json:"icon"`
	Command       string `json:"command,omitempty"`
	MemoryLimit   string `json:"memory_limit,omitempty"` // e.g. "128M", "256M"
	CPUShares     int    `json:"cpu_shares,omitempty"`   // e.g. 512
	EnvVars       string `json:"env_vars,omitempty"`     // Format: KEY=VAL\nKEY2=VAL2

	// Konfigurasi Database Terhubung
	DatabaseID *int   `json:"database_id,omitempty"`
	DbEngine   string `json:"db_engine,omitempty"` // 'mariadb' | 'postgres'
	DbHost     string `json:"db_host,omitempty"`
	DbPort     int    `json:"db_port,omitempty"`
	DbName     string `json:"db_name,omitempty"`
	DbUser     string `json:"db_user,omitempty"`
	DbPassword string `json:"db_password,omitempty"`

	// Status Operasional
	Status    string `json:"status"` // 'running' | 'stopped'
	CreatedAt string `json:"created_at"`
}

var (
	siteMu sync.Mutex
	sites  []Website
	nextID = 1
)

func getStoreFilePath() string {
	paths := []string{
		"backend/internal/store/RAZWebsites.json",
		"../backend/internal/store/RAZWebsites.json",
		"internal/store/RAZWebsites.json",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "backend/internal/store/RAZWebsites.json"
}

func init() {
	loadSitesFromDisk()
}

// loadSitesFromDisk membaca data tersimpan dari file JSON persisten
func loadSitesFromDisk() {
	filePath := getStoreFilePath()
	data, err := os.ReadFile(filePath)
	if err == nil && len(data) > 0 {
		var list []Website
		if err := json.Unmarshal(data, &list); err == nil {
			sites = list
			for _, s := range sites {
				if s.ID >= nextID {
					nextID = s.ID + 1
				}
			}
			return
		}
	}
	sites = []Website{}
}

// saveSitesToDisk menulis seluruh data website ke file JSON persisten
func saveSitesToDisk() {
	filePath := getStoreFilePath()
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)
	data, err := json.MarshalIndent(sites, "", "  ")
	if err == nil {
		_ = os.WriteFile(filePath, data, 0644)
	}
}

// Slugify mengubah nama menjadi string alfanumerik aman untuk container & folder
func Slugify(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := strings.ToLower(reg.ReplaceAllString(input, "-"))
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("site-%d", time.Now().Unix())
	}
	return slug
}

// GetWebsites mengambil daftar seluruh situs web terdaftar
func GetWebsites() []Website {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()
	return sites
}

// FindFreePort mencari port TCP yang bebas di host (prioritas port yang diminta jika tersedia)
func FindFreePort(preferred int) (int, error) {
	if preferred >= 80 && preferred <= 65535 {
		ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", preferred))
		if err == nil {
			_ = ln.Close()
			return preferred, nil
		}
	}

	for port := 8100; port <= 8250; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err == nil {
			_ = ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("tidak ada port bebas tersedia di rentang 8100-8250")
}

// CreateWebsiteRequest payload pembuatan website baru
type CreateWebsiteRequest struct {
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	Port          int    `json:"port"`
	Runtime       string `json:"runtime"` // 'static' | 'php' | 'nodejs' | 'nextjs' | 'custom'
	ContainerName string `json:"container_name"`
	RootFolder    string `json:"root_folder"`
	Icon          string `json:"icon"`
	Command       string `json:"command"`
	MemoryLimit   string `json:"memory_limit"`
	CPUShares     int    `json:"cpu_shares"`
	EnvVars       string `json:"env_vars"`

	// Koneksi Database
	DatabaseID *int   `json:"database_id"`
	DbEngine   string `json:"db_engine"`
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

// CreateWebsite memproses pembuatan container website, direktori root, dan starter files
func CreateWebsite(req CreateWebsiteRequest) (*Website, error) {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("nama website wajib diisi")
	}

	slug := Slugify(req.Name)

	// Pastikan container name unik agar tidak terjadi tabrakan
	containerName := strings.TrimSpace(req.ContainerName)
	if containerName == "" {
		containerName = fmt.Sprintf("site-%s-web", slug)
	}

	for _, s := range sites {
		if s.ContainerName == containerName {
			containerName = fmt.Sprintf("site-%s-%d", slug, time.Now().Unix()%10000)
			break
		}
	}

	// Alokasi Port Bebas
	port, err := FindFreePort(req.Port)
	if err != nil {
		return nil, err
	}

	// Penentuan Root Folder di Host
	rootFolder := strings.TrimSpace(req.RootFolder)
	if rootFolder == "" {
		rootFolder = filepath.Join("/DATA/apps/websites", slug)
	}

	// Buat direktori dan starter template file jika folder baru
	_ = os.MkdirAll(rootFolder, 0755)
	createStarterFiles(rootFolder, req.Name, req.Runtime, req.DbName, req.DbHost, req.DbPort, req.DbUser)

	// Default icon jika kosong
	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		switch req.Runtime {
		case "php":
			icon = "https://cdn.simpleicons.org/php/777BB4"
		case "nodejs":
			icon = "https://cdn.simpleicons.org/nodedotjs/5FA04E"
		case "nextjs":
			icon = "https://cdn.simpleicons.org/nextdotjs/white"
		default:
			icon = "/assets/icons/RAZIconWebsites.svg"
		}
	}

	// Default Memory Limit & CPU Shares
	memLimit := strings.TrimSpace(req.MemoryLimit)
	if memLimit == "" {
		memLimit = "128M" // Default aman untuk STB ARM64
	}
	cpuShares := req.CPUShares
	if cpuShares <= 0 {
		cpuShares = 512
	}

	newSite := Website{
		ID:            nextID,
		Name:          req.Name,
		Slug:          slug,
		Domain:        req.Domain,
		HostPort:      port,
		Runtime:       req.Runtime,
		ContainerName: containerName,
		RootFolder:    rootFolder,
		Icon:          icon,
		Command:       req.Command,
		MemoryLimit:   memLimit,
		CPUShares:     cpuShares,
		EnvVars:       req.EnvVars,

		DatabaseID: req.DatabaseID,
		DbEngine:   req.DbEngine,
		DbHost:     req.DbHost,
		DbPort:     req.DbPort,
		DbName:     req.DbName,
		DbUser:     req.DbUser,
		DbPassword: req.DbPassword,

		Status:    "running",
		CreatedAt: time.Now().Format("2006-01-02 15:04"),
	}
	nextID++

	// Coba deploy ke Docker engine jika daemon aktif
	_ = deployDockerContainer(&newSite)

	// Daftarkan ke Reverse Proxy internal
	if newSite.Domain != "" {
		proxy.GetRouter().AddRoute(newSite.Domain, port)
	}

	sites = append(sites, newSite)
	saveSitesToDisk()

	return &newSite, nil
}

// createStarterFiles membuat file awal (index.html/index.php/index.js) agar web langsung dapat dibuka
func createStarterFiles(dir, name, runtime, dbName, dbHost string, dbPort int, dbUser string) {
	switch runtime {
	case "php":
		indexPath := filepath.Join(dir, "index.php")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<?php
/**
 * RAZOS Powered PHP Website: %s
 * Dibuat secara otomatis oleh RAZOS Website Manager
 */
$dbName = '%s';
$dbHost = '%s';
$dbPort = %d;
$dbUser = '%s';
?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - RAZOS PHP Web</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); }
        .badge { background: rgba(14, 165, 233, 0.2); color: #38bdf8; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
        h1 { margin: 1rem 0 0.5rem; font-size: 1.5rem; }
        p { color: #94a3b8; font-size: 0.875rem; line-height: 1.5; }
        .db-info { background: #0f172a; padding: 1rem; border-radius: 0.75rem; margin-top: 1.5rem; font-size: 0.75rem; text-align: left; font-family: monospace; }
    </style>
</head>
<body>
    <div class="card">
        <span class="badge">PHP 8.2 • RAZOS ISOLATED CONTAINER</span>
        <h1>Selamat Datang di %s</h1>
        <p>Website PHP ini aktif dan berjalan di server RAZOS dengan aman.</p>
        <?php if (!empty($dbName)): ?>
        <div class="db-info">
            <strong style="color: #fbbf24;">Koneksi Database:</strong><br>
            Host: <?php echo htmlspecialchars($dbHost); ?>:<?php echo $dbPort; ?><br>
            Database: <?php echo htmlspecialchars($dbName); ?><br>
            User: <?php echo htmlspecialchars($dbUser); ?>
        </div>
        <?php endif; ?>
    </div>
</body>
</html>`, name, dbName, dbHost, dbPort, dbUser, name, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}

	case "nodejs":
		indexPath := filepath.Join(dir, "server.js")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`// RAZOS Node.js Website: %s
const http = require('http');
const port = process.env.PORT || 80;

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    res.end(\`
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <title>%s - RAZOS Node.js Web</title>
        <style>
            body { font-family: system-ui, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
            .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; }
            .badge { background: rgba(34, 197, 94, 0.2); color: #4ade80; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
        </style>
    </head>
    <body>
        <div class="card">
            <span class="badge">NODE.JS RUNTIME • RAZOS CONTAINER</span>
            <h1>%s</h1>
            <p style="color: #94a3b8;">Node.js service berjalan dengan lancar pada port \${port}.</p>
        </div>
    </body>
    </html>
    \`);
});

server.listen(port, () => {
    console.log('Server berjalan pada port ' + port);
});
`, name, name, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}

	default: // Static HTML
		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - RAZOS Static Website</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); }
        .badge { background: rgba(56, 189, 248, 0.2); color: #38bdf8; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
        h1 { margin: 1rem 0 0.5rem; font-size: 1.5rem; }
        p { color: #94a3b8; font-size: 0.875rem; line-height: 1.5; }
    </style>
</head>
<body>
    <div class="card">
        <span class="badge">NGINX STATIS • RAZOS ISOLATED CONTAINER</span>
        <h1>%s</h1>
        <p>Website statis Anda siap diedit. Masukkan file HTML, CSS, dan aset gambar ke direktori root website untuk mulai membangun proyek.</p>
    </div>
</body>
</html>`, name, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}
	}
}

// parseMemoryBytes mengonversi format "128M", "256M", "1G" menjadi bytes
func parseMemoryBytes(memStr string) int64 {
	memStr = strings.TrimSpace(strings.ToUpper(memStr))
	if memStr == "" {
		return 128 * 1024 * 1024
	}
	if strings.HasSuffix(memStr, "M") || strings.HasSuffix(memStr, "MB") {
		numStr := strings.TrimRight(memStr, "MB")
		num, _ := strconv.ParseInt(numStr, 10, 64)
		return num * 1024 * 1024
	}
	if strings.HasSuffix(memStr, "G") || strings.HasSuffix(memStr, "GB") {
		numStr := strings.TrimRight(memStr, "GB")
		num, _ := strconv.ParseInt(numStr, 10, 64)
		return num * 1024 * 1024 * 1024
	}
	num, _ := strconv.ParseInt(memStr, 10, 64)
	if num > 0 {
		return num
	}
	return 128 * 1024 * 1024
}

// deployDockerContainer menjalankan container Docker yang terisolasi
func deployDockerContainer(site *Website) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err // Docker mungkin belum aktif di lingkungan dev lokal
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tentukan Image & Target Port Container
	imageName := "nginx:alpine"
	containerTargetPort := "80"
	mountDest := "/usr/share/nginx/html"

	switch site.Runtime {
	case "php":
		imageName = "php:8.2-apache"
		mountDest = "/var/www/html"
	case "nodejs":
		imageName = "node:18-alpine"
		mountDest = "/app"
		containerTargetPort = "80"
	case "nextjs":
		imageName = "node:20-alpine"
		mountDest = "/app"
		containerTargetPort = "3000"
	}

	natPort, _ := nat.NewPort("tcp", containerTargetPort)
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: fmt.Sprintf("%d", site.HostPort),
	}

	// Environment variables
	var envList []string
	if site.EnvVars != "" {
		lines := strings.Split(site.EnvVars, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				envList = append(envList, line)
			}
		}
	}
	if site.DbName != "" {
		envList = append(envList,
			fmt.Sprintf("DB_HOST=%s", site.DbHost),
			fmt.Sprintf("DB_PORT=%d", site.DbPort),
			fmt.Sprintf("DB_NAME=%s", site.DbName),
			fmt.Sprintf("DB_USER=%s", site.DbUser),
			fmt.Sprintf("DB_PASS=%s", site.DbPassword),
		)
	}

	// Command custom jika ada
	var cmd []string
	if site.Command != "" {
		cmd = strings.Fields(site.Command)
	}

	// Buat container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   envList,
		Cmd:   cmd,
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", site.RootFolder, mountDest),
		},
		PortBindings: nat.PortMap{
			natPort: []nat.PortBinding{hostBinding},
		},
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		Resources: container.Resources{
			Memory:    parseMemoryBytes(site.MemoryLimit),
			CPUShares: int64(site.CPUShares),
		},
	}, nil, nil, site.ContainerName)

	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}

// ToggleWebsiteStatus mengubah status aktif/berhenti container website
func ToggleWebsiteStatus(id int, action string) (*Website, error) {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	var target *Website
	for i := range sites {
		if sites[i].ID == id {
			target = &sites[i]
			break
		}
	}

	if target == nil {
		return nil, fmt.Errorf("website dengan ID %d tidak ditemukan", id)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer cli.Close()

		switch action {
		case "start":
			_ = cli.ContainerStart(ctx, target.ContainerName, types.ContainerStartOptions{})
			target.Status = "running"
		case "stop":
			stopTimeout := 5
			_ = cli.ContainerStop(ctx, target.ContainerName, container.StopOptions{Timeout: &stopTimeout})
			target.Status = "stopped"
		case "restart":
			stopTimeout := 5
			_ = cli.ContainerRestart(ctx, target.ContainerName, container.StopOptions{Timeout: &stopTimeout})
			target.Status = "running"
		}
	} else {
		// Fallback status toggle jika di mock / host tanpa docker
		if action == "stop" {
			target.Status = "stopped"
		} else {
			target.Status = "running"
		}
	}

	saveSitesToDisk()
	return target, nil
}

// DeleteWebsite menghapus website dari sistem
func DeleteWebsite(id int, purgeFiles bool) error {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	var updated []Website
	var target *Website
	for i := range sites {
		if sites[i].ID == id {
			target = &sites[i]
		} else {
			updated = append(updated, sites[i])
		}
	}

	if target == nil {
		return fmt.Errorf("website dengan ID %d tidak ditemukan", id)
	}

	// Hapus container docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer cli.Close()

		_ = cli.ContainerRemove(ctx, target.ContainerName, types.ContainerRemoveOptions{
			Force: true,
		})
	}

	// Lepas rute reverse proxy
	if target.Domain != "" {
		proxy.GetRouter().RemoveRoute(target.Domain)
	}

	// Purge root folder jika diminta
	if purgeFiles && target.RootFolder != "" && strings.HasPrefix(target.RootFolder, "/DATA/apps/websites") {
		_ = os.RemoveAll(target.RootFolder)
	}

	sites = updated
	saveSitesToDisk()
	return nil
}
