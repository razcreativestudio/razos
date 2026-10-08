/**
 * RAZOS REST API Controller
 * Updated: 2026-10-08
 * Deskripsi: Handler endpoint HTTP REST API v1 untuk metrik, container, proses, dan WebSocket terminal.
 */

package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razos/razos/internal/docker"
	"github.com/razos/razos/internal/system"
	"github.com/razos/razos/internal/terminal"
)

// RegisterRoutes mendaftarkan seluruh endpoint publik dan privat API ke Gin Engine
func RegisterRoutes(r *gin.Engine) {
	// WebSocket Terminal Route
	r.GET("/ws/terminal", terminal.HandleTerminalWebSocket)

	v1 := r.Group("/api/v1")
	{
		// Metrik & Pemeliharaan Sistem
		v1.GET("/system/metrics", handleGetMetrics)
		v1.POST("/system/clean-ram", handleCleanRAM)
		v1.POST("/system/clean-storage", handleCleanStorage)

		// Task Manager (Container & Processes)
		v1.GET("/containers", handleGetContainers)
		v1.POST("/containers/:id/stop", handleStopContainer)
		v1.POST("/containers/:id/restart", handleRestartContainer)
		v1.GET("/processes", handleGetProcesses)
		v1.POST("/processes/:pid/kill", handleKillProcess)

		// Health Check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"system":  "RAZOS Core",
				"version": "1.1.0",
			})
		})
	}
}

func handleGetMetrics(c *gin.Context) {
	metrics, err := system.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func handleCleanRAM(c *gin.Context) {
	err := system.CleanRAM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "RAM cache berhasil dibersihkan",
	})
}

func handleCleanStorage(c *gin.Context) {
	output, err := system.CleanStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Storage docker prune berhasil dieksekusi",
		"output":  output,
	})
}

func handleGetContainers(c *gin.Context) {
	containers, err := docker.GetContainers()
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, containers)
}

func handleStopContainer(c *gin.Context) {
	id := c.Param("id")
	if err := docker.StopContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stopped", "id": id})
}

func handleRestartContainer(c *gin.Context) {
	id := c.Param("id")
	if err := docker.RestartContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "restarted", "id": id})
}

func handleGetProcesses(c *gin.Context) {
	procs, err := docker.GetTopProcesses()
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, procs)
}

func handleKillProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}
	if err := docker.KillProcess(int32(pid)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "killed", "pid": pid})
}
