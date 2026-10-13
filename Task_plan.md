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
- [x] Desain skema katalog aplikasi berbasis YAML/JSON (525+ aplikasi Docker siap pakai).
- [x] Implementasi App Store UI (katalog 525+ apps, filter kategori, pencarian live cepat, form port/volume kustom).
- [x] Full Page App Details View dengan render dokumentasi README.md Markdown visual.
- [x] Tab "Aplikasi Terpasang (Installed Apps)" terdedikasi untuk inspeksi, peluncuran port Web UI, dan manajemen status aktif.
- [x] Editor Docker Compose YAML ala CasaOS dengan dukungan upload/drag-and-drop `.yaml`/`.yml`, live parsing container, preset template cepat, dan eksekusi deployment mandiri.
- [x] Integrasi sinkronisasi Desktop & Start Menu otomatis (`windowStore.allDesktopApps`).
- [x] Fitur Clean Uninstall terintegrasi di Context Menu klik kanan desktop icon, start menu, dan kartu App Store (opsi `purge_data=true`).

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
