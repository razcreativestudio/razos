# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.1.0-phase1-complete (Penyelesaian Fase 1: Fondasi Shell, Backend Go, & SPA Frontend)
**Waktu Eksekusi:** 2026-10-08 22:57 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Inisialisasi Repositori & CI/CD Workflow
- Inisialisasi repositori Git dan penambahan file [.gitignore](file:///c:/laragon/www/razos/.gitignore).
- Konfigurasi workflow GitHub Actions [.github/workflows/build.yml](file:///c:/laragon/www/razos/.github/workflows/build.yml) untuk otomatisasi build frontend Vite dan cross-compile binary static ARM64 (`aarch64` target STB) serta AMD64.

#### 2. Arsitektur Backend Go
- Pembuatan modul struktur folder backend Go di bawah `backend/`:
  - `backend/cmd/server/RAZMainServer.go`: Web server utama menggunakan Gin Release Mode dengan CORS dan recovery middleware.
  - `backend/internal/api/RAZApiController.go`: REST API Controller `/api/v1` (metrik, pemeliharaan sistem, health check).
  - `backend/internal/system/RAZMetricsCollector.go`: Engine pengumpul metrik CPU, RAM fisik, Suhu SoC Amlogic (`/sys/class/thermal`), ZRAM swap, dan partisi `/DATA`.
  - `backend/internal/system/RAZMaintenance.go`: Eksekusi kernel Linux pembersihan cache RAM (`drop_caches`) dan penyimpanan (`docker system prune`).
  - `backend/go.mod`: Definisi dependensi backend ringan.

#### 3. Implementasi Frontend SPA (Vue 3, Vite, Tailwind CSS, Pinia)
- Setup build Vite + Tailwind CSS dengan skema tema glassmorphism Slate-900.
- State Management Pinia:
  - `RAZSystemStore.js`: Polling metrik hardware setiap 3 detik, custom notification toasts, dan aksi Clean RAM/Storage.
  - `RAZWindowStore.js`: Window management desktop multi-tasking (draggable, resizable, minimize, maximize, stacking z-index).
- Komponen Modular:
  - `RAZWindowFrame.vue`: Container window multi-tasking dengan resize handles dan titlebar controls.
  - `RAZTaskbar.vue`: Panel taskbar bawah desktop lengkap dengan start menu, task items aktif, serta system tray (CPU, RAM, Suhu STB, jam digital, tombol fullscreen).
  - `RAZDesktopLayout.vue`: Tampilan desktop mode (>768px) dengan desktop app icon grid, watermark resmi RAZ STUDIO, dan wallpaper WebP teroptimasi.
  - `RAZMobileLayout.vue`: Tampilan mobile launcher (<=768px) dengan Android home grid 4-kolom, widget server STB, dock bawah, dan modal app window fullscreen.
- Views Aplikasi:
  - `RAZTaskManagerView.vue`: Metrik real-time & tombol 1-klik pemeliharaan STB.
  - `RAZTerminalView.vue`: Shell terminal interface.
  - `RAZAppStoreView.vue`: Katalog aplikasi Docker ke `/DATA/apps`.
  - `RAZWebsitesView.vue`: Manajemen multi-website & port reverse proxy.
  - `RAZDatabasesView.vue`: Status MariaDB & PostgreSQL serta link Adminer Web.
  - `RAZSettingsView.vue`: Pengaturan identitas host STB.
- Berhasil diverifikasi dengan `npm run build` tanpa error (output bersih di `dist/`).
- Memperbarui roadmap pada [Task_plan.md](file:///c:/laragon/www/razos/Task_plan.md) (Fase 1 Selesai 100%).
