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
	"io"
	"log"
	"net"
	"os"
	"os/exec"
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

	// Konfigurasi GitHub & Git Sync
	GitRepo       string `json:"git_repo,omitempty"`
	GitBranch     string `json:"git_branch,omitempty"`
	GitToken      string `json:"git_token,omitempty"` // GitHub PAT
	GitLastPull   string `json:"git_last_pull,omitempty"`
	GitLastCommit string `json:"git_last_commit,omitempty"`

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

// GetWebsites mengambil daftar seluruh situs web terdaftar dan menyinkronkan status kontainer
func GetWebsites() []Website {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	// Sinkronkan status realtime jika Docker aktif
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		defer cli.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		changed := false
		for i := range sites {
			inspect, inspectErr := cli.ContainerInspect(ctx, sites[i].ContainerName)
			if inspectErr == nil && inspect.State != nil {
				if inspect.State.Running {
					if sites[i].Status != "running" {
						sites[i].Status = "running"
						changed = true
					}
				} else {
					if sites[i].Status != "stopped" {
						sites[i].Status = "stopped"
						changed = true
					}
				}
			}
		}
		if changed {
			saveSitesToDisk()
		}
	}

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

	// Konfigurasi GitHub Repository & PAT
	GitRepo   string `json:"git_repo"`
	GitBranch string `json:"git_branch"`
	GitToken  string `json:"git_token"`
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

		GitRepo:     strings.TrimSpace(req.GitRepo),
		GitBranch:   strings.TrimSpace(req.GitBranch),
		GitToken:    strings.TrimSpace(req.GitToken),
		GitLastPull: "",
		GitLastCommit: "",

		Status:    "running",
		CreatedAt: time.Now().Format("2006-01-02 15:04"),
	}
	if newSite.GitBranch == "" && newSite.GitRepo != "" {
		newSite.GitBranch = "main"
	}
	nextID++

	// Jika git repo disediakan, lakukan auto-clone/pull ke root folder
	if newSite.GitRepo != "" {
		lastCommit, _, err := syncGitRepo(newSite.RootFolder, newSite.GitRepo, newSite.GitBranch, newSite.GitToken)
		if err == nil {
			newSite.GitLastPull = time.Now().Format("2006-01-02 15:04:05")
			newSite.GitLastCommit = lastCommit
		}
	}

	// Coba deploy ke Docker engine jika daemon aktif
	if err := deployDockerContainer(&newSite); err != nil {
		log.Printf("[RAZWebsite] Gagal deploy container: %v", err)
		newSite.Status = "stopped"
	} else {
		newSite.Status = "running"
	}

	// Daftarkan ke Reverse Proxy internal
	if newSite.Domain != "" {
		proxy.GetRouter().AddRoute(newSite.Domain, port)
	}

	sites = append(sites, newSite)
	saveSitesToDisk()

	return &newSite, nil
}

// createStarterFiles membuat file awal sesuai runtime stack yang dipilih pengguna
func createStarterFiles(dir, name, runtime, dbName, dbHost string, dbPort int, dbUser string) {
	switch runtime {
	case "golang":
		mainGoPath := filepath.Join(dir, "main.go")
		if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`package main

// RAZOS High-Performance Micro-Service: %s
// Sangat hemat memori untuk STB ARM64 (~20-40 MB RAM)
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type StatusResponse struct {
	App       string    `+"`json:\"app\"`"+`
	Status    string    `+"`json:\"status\"`"+`
	Timestamp time.Time `+"`json:\"timestamp\"`"+`
	Arch      string    `+"`json:\"arch\"`"+`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatusResponse{
			App:       "%s",
			Status:    "online",
			Timestamp: time.Now(),
			Arch:      "ARM64 STB Optimized",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `+"`"+`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>%s - RAZOS Go Micro-Service</title>
    <style>
        body { font-family: system-ui, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(56, 189, 248, 0.3); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; }
        .badge { background: rgba(0, 173, 216, 0.2); color: #00add8; padding: 0.3rem 0.8rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
    </style>
</head>
<body>
    <div class="card">
        <span class="badge">GOLANG MICRO-SERVICE • ⚡ 20-40 MB RAM</span>
        <h1 style="margin-top: 1rem;">%s</h1>
        <p style="color: #94a3b8; font-size: 0.875rem;">Service Go native aktif dan sangat responsif di STB Anda.</p>
        <p style="font-family: monospace; font-size: 0.75rem; color: #38bdf8;">GET /api/health</p>
    </div>
</body>
</html>`+"`"+`)
	})

	fmt.Println("Server Go berjalan pada port " + port)
	_ = http.ListenAndServe(":"+port, nil)
}
`, name, name, name, name)
			_ = os.WriteFile(mainGoPath, []byte(content), 0644)
		}

	case "php":
		indexPath := filepath.Join(dir, "index.php")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<?php
/**
 * RAZOS Powered PHP Website: %s
 * Dibuat secara otomatis oleh RAZOS Website Manager (STB Optimized)
 */
$dbName = '%s';
$dbHost = '%s';
$dbPort = %d;
$dbUser = '%s';
$isSqlite = (strpos($dbName, '.sqlite') !== false || empty($dbHost) || $dbHost === 'localhost');
?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - RAZOS PHP Web</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2.5rem; max-width: 520px; text-align: center; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); }
        .badge { background: rgba(14, 165, 233, 0.2); color: #38bdf8; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
        h1 { margin: 1rem 0 0.5rem; font-size: 1.5rem; }
        p { color: #94a3b8; font-size: 0.875rem; line-height: 1.5; }
        .db-info { background: #0f172a; padding: 1rem; border-radius: 0.75rem; margin-top: 1.5rem; font-size: 0.75rem; text-align: left; font-family: monospace; border: 1px solid rgba(251, 191, 36, 0.2); }
    </style>
</head>
<body>
    <div class="card">
        <span class="badge">PHP 8.2 • RAZOS CONTAINER (40-80 MB RAM)</span>
        <h1>Selamat Datang di %s</h1>
        <p>Website PHP ini aktif dan berjalan di server RAZOS dengan aman.</p>
        <?php if (!empty($dbName)): ?>
        <div class="db-info">
            <strong style="color: #fbbf24;">Koneksi Basis Data:</strong><br>
            Tipe: <?php echo $isSqlite ? 'SQLite (Berbasis Berkas - 0 MB Service RAM)' : 'MariaDB / MySQL Server'; ?><br>
            Database: <?php echo htmlspecialchars($dbName); ?><br>
            <?php if (!$isSqlite): ?>
            Host: <?php echo htmlspecialchars($dbHost); ?>:<?php echo $dbPort; ?><br>
            User: <?php echo htmlspecialchars($dbUser); ?>
            <?php endif; ?>
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
			content := fmt.Sprintf(`// RAZOS Node.js Website: %s (STB Light API)
const http = require('http');
const port = process.env.PORT || 80;

const htmlContent = '<!DOCTYPE html>' +
'<html lang="id">' +
'<head>' +
'    <meta charset="UTF-8">' +
'    <title>%s - RAZOS Node.js Web</title>' +
'    <style>' +
'        body { font-family: system-ui, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }' +
'        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; }' +
'        .badge { background: rgba(34, 197, 94, 0.2); color: #4ade80; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }' +
'    </style>' +
'</head>' +
'<body>' +
'    <div class="card">' +
'        <span class="badge">NODE.JS LIGHT RUNTIME • (70-120 MB RAM)</span>' +
'        <h1>%s</h1>' +
'        <p style="color: #94a3b8;">Node.js Express/Fastify service berjalan dengan lancar pada port ' + port + '.</p>' +
'    </div>' +
'</body>' +
'</html>';

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    res.end(htmlContent);
});

server.listen(port, () => {
    console.log('Server berjalan pada port ' + port);
});
`, name, name, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}

	case "python":
		appPyPath := filepath.Join(dir, "app.py")
		if _, err := os.Stat(appPyPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`# RAZOS Python Micro-service: %s
import http.server
import socketserver
import os

PORT = int(os.environ.get("PORT", 80))

HTML_CONTENT = """<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>%s - RAZOS Python Micro-service</title>
    <style>
        body { font-family: system-ui, sans-serif; background: #090d16; color: #f1f5f9; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: rgba(30, 41, 59, 0.7); border: 1px solid rgba(59, 130, 246, 0.3); border-radius: 1.25rem; padding: 2.5rem; max-width: 500px; text-align: center; }
        .badge { background: rgba(59, 130, 246, 0.2); color: #60a5fa; padding: 0.3rem 0.8rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
    </style>
</head>
<body>
    <div class="card">
        <span class="badge">PYTHON MICRO-SERVICE • (80-150 MB RAM)</span>
        <h1 style="margin-top: 1rem;">%s</h1>
        <p style="color: #94a3b8;">FastAPI / Flask micro-service siap melayani permintaan.</p>
    </div>
</body>
</html>"""

class Handler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html; charset=utf-8")
        self.end_headers()
        self.wfile.write(HTML_CONTENT.encode("utf-8"))

with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print(f"Python server berjalan pada port {PORT}")
    httpd.serve_forever()
`, name, name, name)
			_ = os.WriteFile(appPyPath, []byte(content), 0644)
		}

	case "flatcms":
		indexPath := filepath.Join(dir, "index.php")
		contentDir := filepath.Join(dir, "content")
		_ = os.MkdirAll(contentDir, 0755)
		homeMdPath := filepath.Join(contentDir, "home.md")
		if _, err := os.Stat(homeMdPath); os.IsNotExist(err) {
			_ = os.WriteFile(homeMdPath, []byte(fmt.Sprintf("# Selamat Datang di %s\n\nIni adalah website berbasis Flat-File CMS. Anda tidak membutuhkan database MySQL/MariaDB! Seluruh artikel disimpan langsung di dalam file Markdown.", name)), 0644)
		}
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<?php
/**
 * RAZOS Flat-File CMS (Grav / Kirby Style)
 * Bebas Database SQL • Ramah RAM STB (~35-50 MB RAM)
 */
$contentFile = __DIR__ . '/content/home.md';
$content = file_exists($contentFile) ? file_get_contents($contentFile) : 'Halaman belum tersedia.';
?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>%s - Flat-File CMS</title>
    <style>
        body { font-family: system-ui, sans-serif; background: #090d16; color: #f1f5f9; max-width: 680px; margin: 4rem auto; padding: 0 1rem; line-height: 1.6; }
        .badge { background: rgba(99, 102, 241, 0.2); color: #818cf8; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: bold; }
        .post { background: rgba(30, 41, 59, 0.6); border: 1px solid rgba(51, 65, 85, 0.8); border-radius: 1.25rem; padding: 2rem; margin-top: 1.5rem; }
    </style>
</head>
<body>
    <span class="badge">FLAT-FILE CMS • TANPA MYSQL (35-50 MB RAM)</span>
    <div class="post">
        <?php echo nl2br(htmlspecialchars($content)); ?>
    </div>
</body>
</html>`, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}

	case "dashboard":
		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Portal & Bookmark</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-950 text-slate-100 min-h-screen flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-slate-900/80 border border-slate-800 rounded-2xl p-6 text-center shadow-2xl backdrop-blur-md">
        <div class="w-16 h-16 mx-auto rounded-2xl bg-amber-500/20 border border-amber-500/30 flex items-center justify-center mb-4">
            <span class="text-2xl">⚡</span>
        </div>
        <h1 class="text-xl font-bold text-slate-100">%s</h1>
        <p class="text-xs text-slate-400 mt-1 mb-6">Linktree Pribadi & Portal Utilitas Host</p>
        <div class="space-y-2.5">
            <a href="http://localhost:3000" target="_blank" class="block w-full py-2.5 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition">🖥️ RAZOS Desktop OS</a>
            <a href="http://localhost:8088" target="_blank" class="block w-full py-2.5 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition">🗄️ Database Web GUI</a>
            <a href="https://github.com" target="_blank" class="block w-full py-2.5 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition">🐙 GitHub Repositories</a>
        </div>
        <div class="mt-6 pt-4 border-t border-slate-800/80 text-[10px] text-slate-500">
            Self-hosted di STB ARM64 • Konsumsi RAM &lt; 40 MB
        </div>
    </div>
</body>
</html>`, name, name)
			_ = os.WriteFile(indexPath, []byte(content), 0644)
		}

	default: // Static HTML (Nginx Alpine)
		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - RAZOS Landing Page</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-[#090d16] text-slate-100 min-h-screen flex items-center justify-center p-4">
    <div class="max-w-xl w-full bg-slate-900/70 border border-slate-800 rounded-2xl p-8 text-center shadow-2xl backdrop-blur-md">
        <span class="inline-block px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-bold tracking-wide uppercase">
            🔥 Paling Sempurna • 15–40 MB RAM
        </span>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-white mt-4 mb-2 tracking-tight">%s</h1>
        <p class="text-sm text-slate-400 max-w-md mx-auto leading-relaxed">
            Landing page modern ultra ringan yang melayani berkas HTML, CSS, JavaScript, dan aset gambar melalui Nginx Alpine.
        </p>
        <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
            <span class="px-3 py-1.5 rounded-xl bg-slate-800/90 border border-slate-700/80 text-xs text-slate-300 font-mono">STB S905X ARM64</span>
            <span class="px-3 py-1.5 rounded-xl bg-slate-800/90 border border-slate-700/80 text-xs text-slate-300 font-mono">Zero Database Overhead</span>
        </div>
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

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// Tentukan Image & Target Port Container sesuai spesifikasi STB
	imageName := "nginx:alpine"
	containerTargetPort := "80"
	mountDest := "/usr/share/nginx/html"

	switch site.Runtime {
	case "golang":
		imageName = "golang:1.21-alpine"
		mountDest = "/app"
		containerTargetPort = "80"
	case "php", "flatcms":
		imageName = "php:8.2-apache"
		mountDest = "/var/www/html"
		containerTargetPort = "80"
	case "wordpress":
		imageName = "wordpress:php8.2-apache"
		mountDest = "/var/www/html"
		containerTargetPort = "80"
	case "nodejs":
		imageName = "node:18-alpine"
		mountDest = "/app"
		containerTargetPort = "80"
	case "python":
		imageName = "python:3.11-alpine"
		mountDest = "/app"
		containerTargetPort = "80"
	case "dashboard", "filebrowser", "static":
	default:
		imageName = "nginx:alpine"
		mountDest = "/usr/share/nginx/html"
		containerTargetPort = "80"
	}

	// Pastikan Docker Image sudah tersedia di STB atau unduh (ImagePull)
	_, _, inspectErr := cli.ImageInspectWithRaw(ctx, imageName)
	if inspectErr != nil {
		log.Printf("[RAZWebsite] Mengunduh Docker image: %s ...", imageName)
		out, pullErr := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if pullErr != nil {
			log.Printf("[RAZWebsite] Gagal mengunduh image %s: %v", imageName, pullErr)
			return pullErr
		}
		defer out.Close()
		_, _ = io.Copy(io.Discard, out)
	}

	// Hapus jika ada container lama dengan nama sama
	_ = cli.ContainerRemove(ctx, site.ContainerName, types.ContainerRemoveOptions{Force: true})

	natPort, err := nat.NewPort("tcp", containerTargetPort)
	if err != nil {
		natPort = nat.Port(containerTargetPort + "/tcp")
	}
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

	// Command custom jika ada, atau default command sesuai runtime agar container tetap running
	var cmd []string
	if site.Command != "" {
		cmd = strings.Fields(site.Command)
	} else {
		switch site.Runtime {
		case "golang":
			cmd = []string{"go", "run", "main.go"}
		case "nodejs":
			cmd = []string{"node", "server.js"}
		case "python":
			cmd = []string{"python", "app.py"}
		}
	}

	// Pastikan direktori root dan starter file sudah dibuat
	_ = os.MkdirAll(site.RootFolder, 0755)
	createStarterFiles(site.RootFolder, site.Name, site.Runtime, site.DbName, site.DbHost, site.DbPort, site.DbUser)

	// Format mount path yang aman
	bindMount := fmt.Sprintf("%s:%s:rw", filepath.ToSlash(site.RootFolder), mountDest)

	// Buat container dengan ExposedPorts dan WorkingDir
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Env:          envList,
		Cmd:          cmd,
		WorkingDir:   mountDest,
		ExposedPorts: nat.PortSet{natPort: struct{}{}},
	}, &container.HostConfig{
		Binds: []string{
			bindMount,
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

// ToggleWebsiteStatus mengubah status aktif/berhenti container website dengan self-healing redeploy
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
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		defer cli.Close()

		switch action {
		case "start":
			inspect, inspectErr := cli.ContainerInspect(ctx, target.ContainerName)
			if inspectErr != nil {
				// Kontainer belum ada atau terhapus, lakukan redeploy otomatis
				if deployErr := deployDockerContainer(target); deployErr != nil {
					target.Status = "stopped"
					saveSitesToDisk()
					return target, deployErr
				}
				target.Status = "running"
			} else {
				if !inspect.State.Running {
					if startErr := cli.ContainerStart(ctx, target.ContainerName, types.ContainerStartOptions{}); startErr != nil {
						// Jika start gagal karena miskonfigurasi port/image lama, redeploy ulang
						if deployErr := deployDockerContainer(target); deployErr != nil {
							target.Status = "stopped"
							saveSitesToDisk()
							return target, startErr
						}
					}
				}
				target.Status = "running"
			}
		case "stop":
			stopTimeout := 5
			_ = cli.ContainerStop(ctx, target.ContainerName, container.StopOptions{Timeout: &stopTimeout})
			target.Status = "stopped"
		case "restart":
			stopTimeout := 5
			_ = cli.ContainerRestart(ctx, target.ContainerName, container.StopOptions{Timeout: &stopTimeout})
			inspect, inspectErr := cli.ContainerInspect(ctx, target.ContainerName)
			if inspectErr != nil || (inspect.State != nil && !inspect.State.Running) {
				_ = deployDockerContainer(target)
			}
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

// UpdateWebsiteRequest payload untuk memperbarui konfigurasi website yang ada
type UpdateWebsiteRequest struct {
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Port        int    `json:"port"`
	Runtime     string `json:"runtime"`
	RootFolder  string `json:"root_folder"`
	Icon        string `json:"icon"`
	Command     string `json:"command"`
	MemoryLimit string `json:"memory_limit"`
	CPUShares   int    `json:"cpu_shares"`
	EnvVars     string `json:"env_vars"`

	// Koneksi Database
	DatabaseID *int   `json:"database_id"`
	DbEngine   string `json:"db_engine"`
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`

	// Konfigurasi Git
	GitRepo   string `json:"git_repo"`
	GitBranch string `json:"git_branch"`
	GitToken  string `json:"git_token"`
}

// UpdateWebsite memperbarui konfigurasi website yang tersimpan di sistem
func UpdateWebsite(id int, req UpdateWebsiteRequest) (*Website, error) {
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

	name := strings.TrimSpace(req.Name)
	if name != "" {
		target.Name = name
	}

	oldDomain := target.Domain
	if req.Domain != "" {
		target.Domain = strings.TrimSpace(req.Domain)
	}
	if req.Port > 0 {
		target.HostPort = req.Port
	}
	if req.Runtime != "" {
		target.Runtime = req.Runtime
	}
	if req.RootFolder != "" {
		target.RootFolder = strings.TrimSpace(req.RootFolder)
	}
	if req.Icon != "" {
		target.Icon = strings.TrimSpace(req.Icon)
	}
	target.Command = strings.TrimSpace(req.Command)
	if req.MemoryLimit != "" {
		target.MemoryLimit = strings.TrimSpace(req.MemoryLimit)
	}
	if req.CPUShares > 0 {
		target.CPUShares = req.CPUShares
	}
	target.EnvVars = req.EnvVars

	// Update Database
	target.DatabaseID = req.DatabaseID
	target.DbEngine = req.DbEngine
	target.DbHost = req.DbHost
	target.DbPort = req.DbPort
	target.DbName = req.DbName
	target.DbUser = req.DbUser
	if req.DbPassword != "" {
		target.DbPassword = req.DbPassword
	}

	// Update Git
	target.GitRepo = strings.TrimSpace(req.GitRepo)
	target.GitBranch = strings.TrimSpace(req.GitBranch)
	if target.GitBranch == "" && target.GitRepo != "" {
		target.GitBranch = "main"
	}
	if req.GitToken != "" {
		target.GitToken = strings.TrimSpace(req.GitToken)
	}

	// Sinkronkan rute reverse proxy jika domain berubah
	if oldDomain != target.Domain {
		if oldDomain != "" {
			proxy.GetRouter().RemoveRoute(oldDomain)
		}
		if target.Domain != "" {
			proxy.GetRouter().AddRoute(target.Domain, target.HostPort)
		}
	}

	saveSitesToDisk()
	if target.Status == "running" {
		_ = deployDockerContainer(target)
	}
	return target, nil
}


// buildGitAuthURL menyisipkan GitHub Personal Access Token (PAT) ke URL HTTPS jika tersedia
func buildGitAuthURL(repoURL, token string) string {
	repoURL = strings.TrimSpace(repoURL)
	token = strings.TrimSpace(token)
	if token == "" {
		return repoURL
	}

	// Tangani format https://github.com/owner/repo
	if strings.HasPrefix(repoURL, "https://") {
		cleanURL := strings.TrimPrefix(repoURL, "https://")
		// Buang kredensial lama jika ada
		if atIdx := strings.Index(cleanURL, "@"); atIdx != -1 {
			cleanURL = cleanURL[atIdx+1:]
		}
		return fmt.Sprintf("https://oauth2:%s@%s", token, cleanURL)
	}
	return repoURL
}

// sanitizeGitOutput membersihkan token dari pesan log agar tidak terekspos di UI
func sanitizeGitOutput(output, token string) string {
	if token != "" {
		output = strings.ReplaceAll(output, token, "********")
	}
	return output
}

// syncGitRepo mengeksekusi fetch/clone/pull Git ke dalam root directory website
func syncGitRepo(dir, repoURL, branch, token string) (string, string, error) {
	if strings.TrimSpace(repoURL) == "" {
		return "", "", fmt.Errorf("URL GitHub repository belum dikonfigurasi")
	}
	if strings.TrimSpace(branch) == "" {
		branch = "main"
	}

	// Pastikan direktori target ada
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", fmt.Errorf("gagal membuat direktori root: %w", err)
	}

	authURL := buildGitAuthURL(repoURL, token)
	var logBuilder strings.Builder

	gitDir := filepath.Join(dir, ".git")
	isNewGit := false
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		isNewGit = true
	}

	if isNewGit {
		// Inisialisasi Git pada direktori
		initCmd := exec.Command("git", "init")
		initCmd.Dir = dir
		if out, err := initCmd.CombinedOutput(); err != nil {
			return "", "", fmt.Errorf("git init gagal: %s (%w)", sanitizeGitOutput(string(out), token), err)
		} else {
			logBuilder.WriteString(string(out) + "\n")
		}

		// Tambahkan remote origin
		remoteCmd := exec.Command("git", "remote", "add", "origin", authURL)
		remoteCmd.Dir = dir
		_ = remoteCmd.Run()
	} else {
		// Update URL remote origin jika ada perubahan token atau URL
		setRemoteCmd := exec.Command("git", "remote", "set-url", "origin", authURL)
		setRemoteCmd.Dir = dir
		_ = setRemoteCmd.Run()
	}

	// Fetch branch spesifik
	fetchCmd := exec.Command("git", "fetch", "--depth=1", "origin", branch)
	fetchCmd.Dir = dir
	fetchOut, err := fetchCmd.CombinedOutput()
	logBuilder.WriteString(string(fetchOut) + "\n")
	if err != nil {
		// Coba fetch tanpa depth jika depth 1 gagal
		fetchCmd2 := exec.Command("git", "fetch", "origin", branch)
		fetchCmd2.Dir = dir
		fetchOut2, err2 := fetchCmd2.CombinedOutput()
		logBuilder.WriteString(string(fetchOut2) + "\n")
		if err2 != nil {
			return "", sanitizeGitOutput(logBuilder.String(), token), fmt.Errorf("git fetch gagal: %s", sanitizeGitOutput(string(fetchOut2), token))
		}
	}

	// Checkout / Reset hard ke origin/branch agar sinkron sempurna dengan root folder
	resetCmd := exec.Command("git", "reset", "--hard", fmt.Sprintf("origin/%s", branch))
	resetCmd.Dir = dir
	resetOut, err := resetCmd.CombinedOutput()
	logBuilder.WriteString(string(resetOut) + "\n")
	if err != nil {
		// Coba checkout
		checkoutCmd := exec.Command("git", "checkout", "-B", branch, fmt.Sprintf("origin/%s", branch))
		checkoutCmd.Dir = dir
		coOut, coErr := checkoutCmd.CombinedOutput()
		logBuilder.WriteString(string(coOut) + "\n")
		if coErr != nil {
			return "", sanitizeGitOutput(logBuilder.String(), token), fmt.Errorf("git reset/checkout gagal: %s", sanitizeGitOutput(string(resetOut), token))
		}
	}

	// Ambil ringkasan commit terakhir
	logCmd := exec.Command("git", "log", "-1", "--pretty=format:%h - %s (%cr)")
	logCmd.Dir = dir
	lastCommitBytes, _ := logCmd.Output()
	lastCommit := strings.TrimSpace(string(lastCommitBytes))

	return lastCommit, sanitizeGitOutput(logBuilder.String(), token), nil
}

// PullGitRepository menarik pembaruan kode terbaru dari GitHub ke root folder website
func PullGitRepository(siteID int) (*Website, string, error) {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	var target *Website
	for i := range sites {
		if sites[i].ID == siteID {
			target = &sites[i]
			break
		}
	}

	if target == nil {
		return nil, "", fmt.Errorf("website dengan ID %d tidak ditemukan", siteID)
	}

	if strings.TrimSpace(target.GitRepo) == "" {
		return nil, "", fmt.Errorf("URL GitHub repository belum dikonfigurasi untuk website ini")
	}

	lastCommit, logs, err := syncGitRepo(target.RootFolder, target.GitRepo, target.GitBranch, target.GitToken)
	if err != nil {
		return target, logs, err
	}

	target.GitLastPull = time.Now().Format("2006-01-02 15:04:05")
	if lastCommit != "" {
		target.GitLastCommit = lastCommit
	}

	saveSitesToDisk()
	return target, logs, nil
}

// UpdateGitConfig memperbarui konfigurasi GitHub repository, branch, dan token untuk website
func UpdateGitConfig(siteID int, repo, branch, token string) (*Website, error) {
	siteMu.Lock()
	defer siteMu.Unlock()
	loadSitesFromDisk()

	var target *Website
	for i := range sites {
		if sites[i].ID == siteID {
			target = &sites[i]
			break
		}
	}

	if target == nil {
		return nil, fmt.Errorf("website dengan ID %d tidak ditemukan", siteID)
	}

	target.GitRepo = strings.TrimSpace(repo)
	target.GitBranch = strings.TrimSpace(branch)
	if target.GitBranch == "" && target.GitRepo != "" {
		target.GitBranch = "main"
	}
	target.GitToken = strings.TrimSpace(token)

	saveSitesToDisk()
	if target.Status == "running" {
		_ = deployDockerContainer(target)
	}
	return target, nil
}

