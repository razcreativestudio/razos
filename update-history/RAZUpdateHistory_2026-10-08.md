# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.2.0-phase2-complete (Penyelesaian Fase 2: Terminal PTY WebSocket & Task Manager Live Suite)
**Waktu Eksekusi:** 2026-10-08 23:18 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Backend PTY WebSocket Shell & Docker Manager
- Menambahkan modul `backend/internal/terminal/RAZTerminalHandler.go` menggunakan `creack/pty` dan `gorilla/websocket` untuk membuka live interactive bash shell (`/bin/bash` atau fallback shell).
- Menambahkan modul `backend/internal/docker/RAZDockerManager.go` menggunakan Docker SDK Go dan `gopsutil/process` untuk membaca container list, stop/restart container, membaca top proses sistem berdasarkan PID, serta fungsi kill process.
- Mendaftarkan endpoint baru di `backend/internal/api/RAZApiController.go`:
  - `GET /ws/terminal`: WebSocket PTY duplex stream.
  - `GET /api/v1/containers`: Daftar container aktif.
  - `POST /api/v1/containers/:id/stop` & `restart`: Kontrol kontainer 1-klik.
  - `GET /api/v1/processes`: Daftar proses host aktif.
  - `POST /api/v1/processes/:pid/kill`: Terminasi proses PID host seketika.

#### 2. Frontend UI Upgrade (xterm.js & Tabbed Task Manager)
- Memasang `@xterm/xterm` dan `@xterm/addon-fit`.
- Memperbarui `RAZTerminalView.vue`:
  - Emulasi terminal dark slate xterm-256color yang auto-fit terhadap ukuran window frame saat di-drag atau di-resize.
  - Integrasi live WebSocket PTY ke host, lengkap dengan interaktif command fallback shell.
- Memperbarui `RAZTaskManagerView.vue`:
  - Tab navigasi dinamis: **Proses Host (Top PID)** dan **Container Docker (Live)**.
  - Aksi 1-klik interaktif: Kill Process per PID dengan notifikasi custom modal/toast, serta Restart/Stop Container.
  - Tombol pemeliharaan host 1-klik: Clean RAM (`drop_caches`) & Clean Storage (`docker system prune`).
- Memperbarui roadmap pada [Task_plan.md](file:///c:/laragon/www/razos/Task_plan.md) (Fase 2 Selesai 100%).
