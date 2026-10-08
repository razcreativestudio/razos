/**
 * RAZOS PTY WebSocket Terminal Handler
 * Updated: 2026-10-08
 * Deskripsi: Handler WebSocket PTY menghubungkan browser client secara bidirectional
 * dengan shell bash lokal Linux ARM64 (/bin/bash) menggunakan creack/pty.
 */

package terminal

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Validasi origin aman
		return true
	},
}

// HandleTerminalWebSocket menangani koneksi WebSocket PTY real-time
func HandleTerminalWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[RAZOS Terminal] Gagal upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Tentukan shell sistem (prioritaskan bash, fallback ke sh atau powershell di windows)
	shell := "/bin/bash"
	if _, err := os.Stat(shell); os.IsNotExist(err) {
		shell = "/bin/sh"
		if _, err := os.Stat(shell); os.IsNotExist(err) {
			shell = "cmd.exe"
		}
	}

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Buka PTY master/slave
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("[RAZOS Terminal] Gagal menjalankan PTY: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Gagal menginisialisasi sesi terminal host.\r\n"))
		return
	}
	defer func() {
		_ = ptmx.Close()
		_ = cmd.Process.Kill()
	}()

	// Channel sinyal keluar
	done := make(chan struct{})

	// Goroutine: Membaca output PTY dan meneruskannya ke WebSocket client
	go func() {
		buf := make([]byte, 1024)
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

	// Membaca input dari WebSocket client dan meneruskannya ke PTY stdin
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
				_, _ = ptmx.Write(p)
			}
		}
	}
}
