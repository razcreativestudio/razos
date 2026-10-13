# RAZOS (RAZ Operating System)

> Sistem operasi berbasis web berkinerja tinggi dan hemat memori, dirancang untuk mengubah STB Armbian berdaya rendah menjadi server pribadi yang modular.

---

## Fitur Utama

- **Antarmuka Desktop Multi-Window:** Pengalaman desktop lengkap dengan windowing, resize, fluid maximize/restore, sistem animasi minimize/restore menuju taskbar 60 FPS (Genie/Windows style), taskbar indikator glow dinamis, sistem Context Menu (klik kanan) desktop yang kaya fitur, serta title bar window anti-wrap dengan text-overflow truncation dan native tooltip.
- **Antarmuka Mobile Android Launcher:** Beralih otomatis ke launcher sentuh lengkap dengan App Drawer dan App Switcher saat dibuka via smartphone.
- **Dukungan Fullscreen & PWA:** Dapat diinstal ke homescreen ponsel atau desktop tanpa address bar peramban.
- **File Manager & Cloud Storage (Google Drive / Windows File Explorer Hybrid):** File explorer tingkat lanjut dengan Windows-style Interactive Address Bar (Breadcrumbs interaktif, editable path text input, history navigation Back/Forward/Up, Copy Path, dan Refresh), inspeksi drive fisik riil (HDD/MicroSD), fluid auto-fill grid & list, arsitektur **Anti-Hang Media Thumbnail** (IntersectionObserver, Concurrency Queue max 4 requests, browser HTTP caching ETag/304, dan perlindungan file HD raksasa), **Progressive Virtual Chunked Rendering** untuk memuat ribuan berkas tanpa lag/hang, context menu berkas/folder/canvas, bintang/favorit, serta standalone multi-window media viewers & players (Image zoom/rotate, Video player, **Turntable Vinyl & Animated Equalizer Waveform Music Player** dengan autoplay, interactive clickable & slidable timeline scrubber anti-jitter, hover floating preview, multi-speed playback 0.5x - 2.0x, kontrol volume & skip 10s, PDF viewer, dan in-browser Code editor).
- **App Store 525+ Docker, Installed Apps Tab & CasaOS-Style Compose YAML Suite:** Katalog 525+ aplikasi Docker siap pakai dengan halaman detail dokumentasi README.md, tab khusus **"Aplikasi Terpasang (Installed Apps)"** untuk manajemen container aktif, peluncuran Web UI & clean uninstall, integrasi **Context Menu** klik kanan shortcut desktop & start menu untuk uninstall instan, serta editor **Docker Compose YAML ala CasaOS** dengan upload/drag-and-drop `.yaml`/`.yml`, syntax inspection real-time, dan deployment multi-container mandiri.
- **Task Manager & System Health:** Monitoring penggunaan sumber daya per-aplikasi secara real-time dengan kill-switch instan.
- **Pembersih Sistem 1-Klik:** RAM Cleaner dan Docker Storage Cleaner terintegrasi untuk menjaga kestabilan perangkat ber-RAM kecil.
- **Website Manager & Container Isolated Suite:** Platform pembuatan dan pengelolaan multi-website terisolasi berbasis Docker (HTML Statis Nginx, PHP 8.2 + Database, Node.js, Next.js, dan Custom Docker stack) dengan penentuan root folder dokumen (`/DATA/apps/websites/<slug>`), auto-assign port bebas host (8100-8250) anti-tabrakan, penamaan container unik (`site-<slug>-web`), pembatasan sumber daya (Memory limits 64M-1024M optimal STB ARM64 & CPU shares), custom command & environment variables, sinkronisasi link icon ke Desktop shortcut & Start Menu, serta tombol cepat pembuka root folder di File Manager.
- **Database Suite & Web Manager GUI:** Server MariaDB dan PostgreSQL terpusat dengan Web GUI Adminer (HeidiSQL / phpMyAdmin web style) siap pakai untuk inspeksi visual tabel dan eksekusi query data, terhubung langsung ke Website Manager.
- **Web Terminal:** Akses SSH/PTY langsung di peramban menggunakan xterm.js.

---

## Kebutuhan Sistem Minimum

- **Perangkat:** STB Amlogic S905X / Single Board Computer ARM64.
- **Sistem Operasi:** Armbian Linux (Debian Bookworm didukung).
- **RAM:** Minimum 1 GB (Direkomendasikan 2 GB).
- **Penyimpanan:** MicroSD minimal 16 GB terpasang di `/DATA`.

---

## Pemasangan Cepat

Jalankan perintah berikut di terminal SSH STB Anda:

```bash
curl -fsSL https://raw.githubusercontent.com/razcreative/razos/main/scripts/install.sh | bash
```

Setelah instalasi selesai, buka peramban Anda di:
👉 `http://IP-STB-ANDA:8080`

---

## Lisensi

Didistribusikan di bawah lisensi MIT. Silakan gunakan dan kembangkan secara bebas.
