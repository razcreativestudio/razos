/**
 * RAZOS Main Server Entrypoint
 * Updated: 2026-10-08
 * Deskripsi: Server inti RAZOS berbasis Go dengan konsumsi memori hemat (<30 MB)
 * yang menjalankan router REST API dan menyajikan SPA frontend statis.
 */

package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/razos/razos/internal/api"
)

func main() {
	// Jalankan Gin dalam mode Release untuk menghemat CPU dan memori di STB
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// Konfigurasi CORS agar frontend dapat berkomunikasi lancar saat dev mode
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Daftarkan seluruh rute REST API
	api.RegisterRoutes(r)

	// Fallback route untuk root info
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     "RAZOS - RAZ Operating System",
			"author":  "RAZ CREATIVE STUDIO",
			"status":  "running",
			"port":    8080,
		})
	})

	log.Println("[RAZOS] Server berjalan pada http://0.0.0.0:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[RAZOS] Gagal memulai server: %v", err)
	}
}
