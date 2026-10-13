# RAZOS (RAZ Operating System)

> Sistem operasi berbasis web berkinerja tinggi dan hemat memori, dirancang untuk mengubah STB Armbian berdaya rendah menjadi server pribadi yang modular.

---

## Fitur Utama

- **Antarmuka Desktop Multi-Window:** Pengalaman desktop lengkap dengan windowing, resize, fluid maximize/restore, sistem animasi minimize/restore menuju taskbar 60 FPS (Genie/Windows style), taskbar indikator glow dinamis, sistem Context Menu (klik kanan) desktop yang kaya fitur, serta title bar window anti-wrap dengan text-overflow truncation dan native tooltip.
- **Antarmuka Mobile Android Launcher:** Beralih otomatis ke launcher sentuh lengkap dengan App Drawer dan App Switcher saat dibuka via smartphone.
- **Dukungan Fullscreen & PWA:** Dapat diinstal ke homescreen ponsel atau desktop tanpa address bar peramban.
- **File Manager & Cloud Storage (Google Drive / Windows File Explorer Hybrid):** File explorer tingkat lanjut dengan Windows-style Interactive Address Bar (Breadcrumbs interaktif, editable path text input, history navigation Back/Forward/Up, Copy Path, dan Refresh), inspeksi drive fisik riil (HDD/MicroSD), fluid auto-fill grid & list, arsitektur **Anti-Hang Media Thumbnail** (IntersectionObserver, Concurrency Queue max 4 requests, browser HTTP caching ETag/304, dan perlindungan file HD raksasa), **Progressive Virtual Chunked Rendering** untuk memuat ribuan berkas tanpa lag/hang, context menu berkas/folder/canvas, bintang/favorit, serta standalone multi-window media viewers & players (Image zoom/rotate, Video player, **Turntable Vinyl & Animated Equalizer Waveform Music Player** dengan autoplay, interactive clickable & slidable timeline scrubber anti-jitter, hover floating preview, multi-speed playback 0.5x - 2.0x, kontrol volume & skip 10s, PDF viewer, **Dedicated Markdown Document Reader & Editor Suite** dengan mode Preview, Split live-sync, Raw, Table of Contents outline heading otomatis, statistik kata & estimasi waktu baca, serta in-browser Code editor).
- **App Store 525+ Docker, Installed Apps Tab & CasaOS-Style Compose YAML Suite:** Katalog 525+ aplikasi Docker siap pakai dengan halaman detail dokumentasi README.md, tab khusus **"Aplikasi Terpasang (Installed Apps)"** untuk manajemen container aktif, peluncuran Web UI & clean uninstall, integrasi **Context Menu** klik kanan shortcut desktop & start menu untuk uninstall instan, serta editor **Docker Compose YAML ala CasaOS** dengan upload/drag-and-drop `.yaml`/`.yml`, syntax inspection real-time, dan deployment multi-container mandiri.
- **Task Manager & System Health:** Monitoring penggunaan sumber daya per-aplikasi secara real-time dengan kill-switch instan.
- **Pembersih Sistem 1-Klik:** RAM Cleaner dan Docker Storage Cleaner terintegrasi untuk menjaga kestabilan perangkat ber-RAM kecil.
- **Website Manager & STB-Optimized Multi-Hosting Suite:** Platform pembuatan dan pengelolaan multi-website terisolasi berbasis Docker yang dioptimasi khusus untuk keterbatasan hardware STB ARM64 (RAM 1.8 GiB) dengan mengeliminasi runtime berat (seperti Next.js SSR monolitik) dan menghadirkan **4 Kategori Hosting Ideal**:
  1. *Website Statis & Landing Page* (HTML, CSS, JS, Tailwind, Astro/Hugo/Vite static build via Nginx Alpine, RAM 15–40 MB).
  2. *Micro-Framework Web & REST API* (Go Micro-service binary native RAM 20–40 MB, PHP Native 40–80 MB, Node.js Express/Fastify 70–120 MB, Python FastAPI/Flask 80–150 MB) didukung **SQLite Berbasis Berkas (0 MB Service RAM)** sebagai default utama atau MariaDB Alpine.
  3. *CMS Ringan* (Flat-File CMS Grav/Kirby Style tanpa MySQL RAM 35–50 MB, WordPress Optimized).
  4. *Web Dashboard & Internal Tools* (Linktree Pribadi / Bookmark Portal < 40 MB RAM, File Sharing Ringan).
  Dilengkapi alokasi port host anti-tabrakan (8100-8250), isolasi folder `/DATA/apps/websites/<slug>`, resource limit (Memory 32M-512M STB presets & CPU shares), sinkronisasi icon ke desktop, **Action Button Edit** dengan modal konfigurasi 5-tab, serta **Integrasi GitHub Repository & Personal Access Token (PAT)** untuk auto-pull kode terbaru ke root folder secara instan.
- **Database Suite & Web Manager GUI:** Server MariaDB, PostgreSQL, dan SQLite terpusat dengan Web GUI Adminer (HeidiSQL / phpMyAdmin web style) siap pakai untuk inspeksi visual tabel dan eksekusi query data, terhubung langsung ke Website Manager.
- **Web Terminal:** Akses SSH/PTY langsung di peramban menggunakan xterm.js.
- **Detached Standalone Windows Architecture (Seperti File Manager - Zero In-App Blocking):** Seluruh dialog konfigurasi dan form wizard (Wizard Pembuatan Website, Edit Konfigurasi 5-Tab, GitHub PAT & Git Pull Sync, Rincian Website, Properti Berkas File Manager, dan Buat Database) kini bertransformasi menjadi **Jendela Desktop Mandiri Terpisah (`RAZWindowStore`)** persis seperti sistem pembuka berkas File Manager. Jendela-jendela ini memiliki frame mandiri, titlebar desktop, draggable, resizable, minimizeable ke taskbar dengan ikonnya sendiri, dan **sama sekali tidak mengunci/mengganggu aplikasi utama**. Pengguna dapat leluasa mengklik, menggulir, dan menjelajahi Website Manager, File Manager, atau App Store secara simultan tanpa halangan overlay modal.

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
