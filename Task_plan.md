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
- [ ] Integrasi backend PTY WebSocket (`creack/pty`) ke bash shell server.
- [ ] Pembuatan komponen window terminal menggunakan `xterm.js` dan `xterm-addon-fit`.
- [ ] Pembuatan Task Manager: Tabel monitoring proses dan container (`docker stats`).
- [ ] Fitur 1-klik Task Manager: Kill Process dan Stop/Restart Container.
- [ ] Pembuatan modul Maintenance: RAM Cleaner (`drop_caches`) dan Storage Cleaner (`docker system prune`).

## Fase 3: App Store & Manajemen Container
- [ ] Desain skema katalog aplikasi berbasis YAML/JSON.
- [ ] Implementasi App Store UI (kategori, deskripsi, form port/volume kustom).
- [ ] Engine deploy aplikasi berbasis Docker Compose / Docker Container API ke `/DATA/apps/`.
- [ ] Fitur Clean Uninstall: Menghapus container sekaligus opsi memusnahkan folder data terkait.

## Fase 4: Website & Database Management Suite
- [ ] Integrasi Reverse Proxy internal (Caddy embedded / Traefik ringan) untuk routing port domain.
- [ ] Wizard pembuatan website: Node.js, PHP/HTML statis dengan auto-assign port bebas.
- [ ] Deployment container database sentral terisolasi (MariaDB & PostgreSQL).
- [ ] Integrasi GUI database web instan (Adminer container pada port internal).
- [ ] Form koneksi database instan untuk website yang dibuat.

## Fase 5: OTA Updater & Production Hardening
- [ ] Implementasi API updater: Pengecekan rilis GitHub Releases dan penggantian binary mandiri.
- [ ] Pembuatan systemd unit file `/etc/systemd/system/razos.service`.
- [ ] Skrip pembersihan total sistem lama STB sebelum instalasi bersih.
