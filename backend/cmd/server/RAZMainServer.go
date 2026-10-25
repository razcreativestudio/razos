/**
 * RAZOS Main Server Entrypoint
 * Updated: 2026-10-25
 * Deskripsi: Server inti RAZOS berbasis Go dengan konsumsi memori hemat (<30 MB)
 * yang menjalankan router REST API dan menyajikan SPA frontend Desktop statis.
 */

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/razos/razos/internal/api"
)

func main() {
	// Jalankan Gin dalam mode Release untuk menghemat CPU dan memori di STB
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// Konfigurasi CORS agar frontend dapat berkomunikasi lancar
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Daftarkan seluruh rute REST API
	api.RegisterRoutes(r)

	// Deteksi direktori frontend dist yang tersedia di STB atau local
	distCandidates := []string{
		"/DATA/razos/dist",
		"/DATA/razos/src/frontend/dist",
		"../../frontend/dist",
		"frontend/dist",
		"dist",
	}

	var activeDist string
	for _, p := range distCandidates {
		if stat, err := os.Stat(p); err == nil && stat.IsDir() {
			activeDist = p
			break
		}
	}

	if activeDist != "" {
		log.Printf("[RAZOS] Menyajikan Web Desktop UI dari: %s", activeDist)
		r.Static("/assets", filepath.Join(activeDist, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(activeDist, "favicon.ico"))
		r.StaticFile("/manifest.json", filepath.Join(activeDist, "manifest.json"))
		r.StaticFile("/index.html", filepath.Join(activeDist, "index.html"))

		// Route utama: sajikan langsung index.html
		r.GET("/", func(c *gin.Context) {
			c.File(filepath.Join(activeDist, "index.html"))
		})

		// SPA Fallback: semua rute non-API dan non-WS diarahkan ke index.html
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/ws") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint API tidak ditemukan"})
				return
			}
			c.File(filepath.Join(activeDist, "index.html"))
		})
	} else {
		// Fallback route jika dist belum terpasang
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"app":     "RAZOS - RAZ Operating System",
				"author":  "RAZ CREATIVE STUDIO",
				"status":  "running",
				"port":    8080,
				"warning": "Frontend dist belum ditemukan di /DATA/razos/src/frontend/dist",
			})
		})
	}

	// Jalankan listener port 8080 di goroutine (kompatibilitas backward)
	go func() {
		srv8080 := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		log.Println("[RAZOS] Server sekunder listening pada port 8080 (http://0.0.0.0:8080)")
		if err := srv8080.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[RAZOS] Info port 8080: %v", err)
		}
	}()

	// Jalankan listener port 80 utama (akses langsung IP tanpa memasukkan port di URL)
	srv80 := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	log.Println("[RAZOS] Server utama listening pada port 80 (http://0.0.0.0 / Langsung IP)")
	if err := srv80.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[RAZOS] Gagal menjalankan server pada port 80: %v", err)
	}
}
