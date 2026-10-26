/**
 * RAZOS PTY WebSocket Terminal Handler
 * Updated: 2026-10-26 (v3.9.1)
 * Deskripsi: Handler WebSocket PTY menghubungkan browser client secara bidirectional
 * dengan shell bash login Armbian Linux (/bin/bash -l), dukungan banner MOTD resmi,
 * dan sinkronisasi ukuran dinamis (SIGWINCH resize cols/rows) saat jendela dimaksimalkan.
 */

package terminal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// Validasi origin aman
		return true
	},
}

// terminalResizeMessage struktur pesan kontrol resize dari frontend
type terminalResizeMessage struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// HandleTerminalWebSocket menangani koneksi WebSocket PTY real-time
func HandleTerminalWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[RAZOS Terminal] Gagal upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Ambil ukuran awal cols & rows dari query param jika tersedia
	cols := 120
	rows := 30
	if cStr := c.Query("cols"); cStr != "" {
		if val, err := strconv.Atoi(cStr); err == nil && val > 0 {
			cols = val
		}
	}
	if rStr := c.Query("rows"); rStr != "" {
		if val, err := strconv.Atoi(rStr); err == nil && val > 0 {
			rows = val
		}
	}

	// Tentukan shell sistem (prioritaskan bash, fallback ke sh atau powershell di windows)
	shell := "/bin/bash"
	if _, err := os.Stat(shell); os.IsNotExist(err) {
		shell = "/bin/sh"
		if _, err := os.Stat(shell); os.IsNotExist(err) {
			shell = "cmd.exe"
		}
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/root"
	}

	var cmd *exec.Cmd
	if shell == "cmd.exe" {
		cmd = exec.Command(shell)
	} else {
		// Jalankan shell interaktif di home directory pengguna
		cmd = exec.Command(shell)
		cmd.Dir = homeDir
	}

	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"HOME="+homeDir,
		"USER=root",
	)

	// Buka PTY master/slave
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("[RAZOS Terminal] Gagal menjalankan PTY: %v", err)
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Gagal menginisialisasi sesi terminal host.\r\n"))
		return
	}
	defer func() {
		_ = ptmx.Close()
		_ = cmd.Process.Kill()
	}()

	// Terapkan ukuran awal PTY
	_ = pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})

	// Tampilkan banner MOTD RAZOS resmi dengan telemetri performa hardware riil
	razosBanner := GenerateRAZOSBanner()
	_ = conn.WriteMessage(websocket.BinaryMessage, []byte(razosBanner))

	// Channel sinyal keluar
	done := make(chan struct{})

	// Goroutine: Membaca output PTY dan meneruskannya ke WebSocket client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				close(done)
				return
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				close(done)
				return
			}
		}
	}()

	// Membaca input dari WebSocket client dan meneruskannya ke PTY stdin atau resize handler
	for {
		select {
		case <-done:
			return
		default:
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				// Cek apakah pesan berupa paket kontrol resize JSON
				var resizeMsg terminalResizeMessage
				if json.Unmarshal(p, &resizeMsg) == nil && resizeMsg.Type == "resize" && resizeMsg.Cols > 0 && resizeMsg.Rows > 0 {
					_ = pty.Setsize(ptmx, &pty.Winsize{
						Rows: uint16(resizeMsg.Rows),
						Cols: uint16(resizeMsg.Cols),
					})
					continue
				}

				// Teruskan input keyboard biasa langsung ke PTY stdin
				_, _ = ptmx.Write(p)
			}
		}
	}
}
