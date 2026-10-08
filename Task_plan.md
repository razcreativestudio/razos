# RAZOS Task Plan & Roadmap

## Fase 1: Fondasi Shell & Sistem Backend
- [x] Inisialisasi repositori Git dan konfigurasi GitHub Actions workflow (Kompilasi ARM64).
- [x] Setup struktur direktori backend Go (Gin/Fiber, Gopsutil, Docker SDK).
- [x] Implementasi endpoint metrik sistem (`/api/system/stats`: CPU, RAM, Suhu, ZRAM, Disk).
- [x] Setup frontend SPA (Vue 3, Vite, Tailwind CSS, Pinia, Lucide Icons).
- [x] Implementasi Desktop Window Manager (Draggable, Resizable, Z-Index stacking, Minimize/Maximize).
- [x] Implementasi Mobile Launcher Layout (Android Home Grid, Bottom Dock, App Drawer).
- [x] Integrasi Fullscreen API dan PWA Web Manifest.

## Fase 2: Terminal & Task Manager
- [x] Integrasi backend PTY WebSocket (`creack/pty`) ke bash shell server.
- [x] Pembuatan komponen window terminal menggunakan `xterm.js` dan `xterm-addon-fit`.
- [x] Pembuatan Task Manager: Tabel monitoring proses dan container (`docker stats`).
- [x] Fitur 1-klik Task Manager: Kill Process dan Stop/Restart Container.
- [x] Pembuatan modul Maintenance: RAM Cleaner (`drop_caches`) dan Storage Cleaner (`docker system prune`).

## Fase 3: App Store & Manajemen Container
- [x] Desain skema katalog aplikasi berbasis YAML/JSON.
- [x] Implementasi App Store UI (kategori, deskripsi, form port/volume kustom).
- [x] Engine deploy aplikasi berbasis Docker Compose / Docker Container API ke `/DATA/apps/`.
- [x] Fitur Clean Uninstall: Menghapus container sekaligus opsi memusnahkan folder data terkait.

## Fase 4: Website & Database Management Suite
- [x] Integrasi Reverse Proxy internal (Caddy embedded / Traefik ringan) untuk routing port domain.
- [x] Wizard pembuatan website: Node.js, PHP/HTML statis dengan auto-assign port bebas.
- [x] Deployment container database sentral terisolasi (MariaDB & PostgreSQL).
- [x] Integrasi GUI database web instan (Adminer container pada port internal).
- [x] Form koneksi database instan untuk website yang dibuat.

## Fase 5: OTA Updater & Production Hardening
- [x] Implementasi API updater: Pengecekan rilis GitHub Releases dan penggantian binary mandiri.
- [x] Pembuatan systemd unit file `/etc/systemd/system/razos.service`.
- [x] Skrip pembersihan total sistem lama STB sebelum instalasi bersih.
