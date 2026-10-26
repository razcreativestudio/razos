/**
 * RAZOS Native DevAPI Bridge Module
 * File: RAZDevBridge.go
 * Updated: 2026-10-26 (v4.0.0)
 * Deskripsi: Modul jembatan pengembang (DevBridge) native berbasis REST API
 * untuk menghubungkan AI Assistant / MCP Server dengan STB Armbian secara aman.
 * Fitur:
 * 1. Token Secret Authentication (X-RAZOS-Dev-Token)
 * 2. Remote Terminal Command Execution (/api/v1/dev/exec)
 * 3. Hot-Reload Automation (Git Pull + Go Build + Service Restart) (/api/v1/dev/hot-reload)
 * 4. Fast File Sync (/api/v1/dev/sync-file)
 * 5. Real-Time Journalctl & Docker Logs (/api/v1/dev/logs)
 * 6. Live STB Status (/api/v1/dev/status)
 */

package dev

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/mem"
)

// DevSecretToken token otentikasi rahasia untuk DevBridge MCP
const DevSecretToken = "razos-super-secret-dev-bridge-2026"

// AuthMiddleware memvalidasi token rahasia sebelum mengakses API pengembang
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-RAZOS-Dev-Token")
		if token == "" {
			token = c.Query("token")
		}
		if token != DevSecretToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Akses ditolak: X-RAZOS-Dev-Token tidak valid atau belum disertakan",
			})
			return
		}
		c.Next()
	}
}

// RegisterDevRoutes mendaftarkan seluruh endpoint DevBridge ke router Gin
func RegisterDevRoutes(rg *gin.RouterGroup) {
	rg.Use(AuthMiddleware())
	rg.GET("/status", handleDevStatus)
	rg.POST("/exec", handleDevExec)
	rg.POST("/hot-reload", handleDevHotReload)
	rg.POST("/sync-file", handleDevSyncFile)
	rg.GET("/logs", handleDevLogs)
}

// handleDevStatus mengembalikan status kesehatan, versi git, dan telemetri STB
func handleDevStatus(c *gin.Context) {
	// Ambil commit hash git terakhir
	gitCommit := "unknown"
	cmdCommit := exec.Command("git", "-C", "/DATA/razos/src", "rev-parse", "HEAD")
	if out, err := cmdCommit.Output(); err == nil {
		gitCommit = strings.TrimSpace(string(out))
	}

	gitBranch := "unknown"
	cmdBranch := exec.Command("git", "-C", "/DATA/razos/src", "branch", "--show-current")
	if out, err := cmdBranch.Output(); err == nil {
		gitBranch = strings.TrimSpace(string(out))
	}

	// Cek sisa RAM
	var ramFreeMB uint64
	if vMem, err := mem.VirtualMemory(); err == nil {
		ramFreeMB = vMem.Available / 1024 / 1024
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "online",
		"device":      "RAZOS ARM64 STB (Amlogic S905X)",
		"git_branch":  gitBranch,
		"git_commit":  gitCommit,
		"ram_free_mb": ramFreeMB,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	})
}

// ExecRequest payload untuk eksekusi perintah bash di STB
type ExecRequest struct {
	Command    string `json:"command" binding:"required"`
	WorkDir    string `json:"workdir"`
	TimeoutSec int    `json:"timeout_sec"`
}

// handleDevExec mengeksekusi perintah shell Linux di STB dan mengembalikan stdout/stderr
func handleDevExec(c *gin.Context) {
	var req ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid: " + err.Error()})
		return
	}

	timeout := 60
	if req.TimeoutSec > 0 {
		timeout = req.TimeoutSec
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", req.Command)

	// Tentukan direktori kerja (default ke /DATA/razos/src)
	workDir := strings.TrimSpace(req.WorkDir)
	if workDir == "" {
		workDir = "/DATA/razos/src"
	}
	cmd.Dir = workDir

	start := time.Now()
	outBytes, cmdErr := cmd.CombinedOutput()
	duration := time.Since(start).Milliseconds()

	exitCode := 0
	if cmdErr != nil {
		if exitError, ok := cmdErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"command":     req.Command,
		"output":      string(outBytes),
		"exit_code":   exitCode,
		"duration_ms": duration,
		"success":     exitCode == 0,
	})
}

// HotReloadRequest opsi hot-reload
type HotReloadRequest struct {
	PullGit    bool `json:"pull_git"`
	BuildGo    bool `json:"build_go"`
	RestartSvc bool `json:"restart_svc"`
}

// handleDevHotReload mengeksekusi git reset/pull, go build, dan restart service razos
func handleDevHotReload(c *gin.Context) {
	var req HotReloadRequest
	// Default nilai true jika tidak dikirimkan
	if c.Request.ContentLength > 0 {
		_ = c.ShouldBindJSON(&req)
	} else {
		req.PullGit = true
		req.BuildGo = true
		req.RestartSvc = true
	}

	var logSteps []string

	// 1. Tarik perubahan terbaru dari GitHub
	if req.PullGit {
		cmdGit := exec.Command("bash", "-c", "git fetch origin && git reset --hard origin/master")
		cmdGit.Dir = "/DATA/razos/src"
		outGit, errGit := cmdGit.CombinedOutput()
		logSteps = append(logSteps, fmt.Sprintf("=== [1. GIT PULL] ===\n%s", strings.TrimSpace(string(outGit))))
		if errGit != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal sinkronisasi git: " + errGit.Error(),
				"logs":  logSteps,
			})
			return
		}
	}

	// 2. Kompilasi Go Binary RAZOS
	if req.BuildGo {
		cmdBuild := exec.Command("bash", "-c", "GOTOOLCHAIN=local GOPATH=/DATA/razos/gopath go build -mod=mod -ldflags=\"-s -w\" -o /DATA/razos/bin/razos ./cmd/server/RAZMainServer.go")
		cmdBuild.Dir = "/DATA/razos/src/backend"
		outBuild, errBuild := cmdBuild.CombinedOutput()
		logSteps = append(logSteps, fmt.Sprintf("=== [2. GO BUILD] ===\n%s", strings.TrimSpace(string(outBuild))))
		if errBuild != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal kompilasi binary Go: " + errBuild.Error(),
				"logs":  logSteps,
			})
			return
		}
	}

	// 3. Restart Service RAZOS secara Asinkron (agar response HTTP sempat dikirim ke klien)
	if req.RestartSvc {
		logSteps = append(logSteps, "=== [3. SERVICE RESTART] ===\nService razos.service sedang di-restart via background nohup...")
		go func() {
			time.Sleep(800 * time.Millisecond)
			_ = exec.Command("bash", "-c", "systemctl restart razos").Run()
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hot-Reload berhasil dieksekusi!",
		"logs":    logSteps,
	})
}

// SyncFileRequest payload pengiriman file langsung ke STB
type SyncFileRequest struct {
	Path    string `json:"path" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// handleDevSyncFile menulis berkas langsung ke disk STB tanpa perlu menunggu git push
func handleDevSyncFile(c *gin.Context) {
	var req SyncFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid: " + err.Error()})
		return
	}

	// Bersihkan dan tentukan target path
	relPath := strings.TrimPrefix(req.Path, "/")
	relPath = strings.TrimPrefix(relPath, "\\")
	var targetPath string

	if filepath.IsAbs(req.Path) {
		targetPath = req.Path
	} else {
		targetPath = filepath.Join("/DATA/razos/src", relPath)
	}

	// Batasi agar tidak menimpa sistem root di luar /DATA/
	if !strings.HasPrefix(targetPath, "/DATA/") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: File sync hanya diperbolehkan di dalam partisi /DATA/"})
		return
	}

	// Buat folder induk jika belum ada
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori: " + err.Error()})
		return
	}

	// Tulis berkas ke disk
	if err := os.WriteFile(targetPath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis berkas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"path":          targetPath,
		"bytes_written": len(req.Content),
	})
}

// handleDevLogs membaca log journalctl atau docker di STB
func handleDevLogs(c *gin.Context) {
	service := c.DefaultQuery("service", "razos")
	lines := c.DefaultQuery("lines", "50")

	var cmd *exec.Cmd
	if service == "razos" {
		cmd = exec.Command("journalctl", "-u", "razos", "-n", lines, "--no-pager")
	} else if strings.HasPrefix(service, "docker:") {
		containerName := strings.TrimPrefix(service, "docker:")
		cmd = exec.Command("docker", "logs", "--tail", lines, containerName)
	} else {
		cmd = exec.Command("journalctl", "-u", service, "-n", lines, "--no-pager")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"output": string(out),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": service,
		"lines":   lines,
		"output":  string(out),
	})
}
