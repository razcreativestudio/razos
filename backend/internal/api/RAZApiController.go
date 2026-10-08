/**
 * RAZOS REST API Controller
 * Updated: 2026-10-08
 * Deskripsi: Handler endpoint HTTP REST API v1 untuk metrik, pemeliharaan sistem, dan status host.
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razos/razos/internal/system"
)

// RegisterRoutes mendaftarkan seluruh endpoint publik dan privat API ke Gin Engine
func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Metrik & Pemeliharaan Sistem
		v1.GET("/system/metrics", handleGetMetrics)
		v1.POST("/system/clean-ram", handleCleanRAM)
		v1.POST("/system/clean-storage", handleCleanStorage)

		// Health Check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"system":  "RAZOS Core",
				"version": "1.0.0",
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
