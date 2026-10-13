# RAZOS (RAZ Operating System)

> Sistem operasi berbasis web berkinerja tinggi dan hemat memori, dirancang untuk mengubah STB Armbian berdaya rendah menjadi server pribadi yang modular.

---

## Fitur Utama

- **Antarmuka Desktop Multi-Window:** Pengalaman desktop lengkap dengan windowing, resize, minimize, taskbar, dan system tray.
- **Antarmuka Mobile Android Launcher:** Beralih otomatis ke launcher sentuh lengkap dengan App Drawer dan App Switcher saat dibuka via smartphone.
- **Dukungan Fullscreen & PWA:** Dapat diinstal ke homescreen ponsel atau desktop tanpa address bar peramban.
- **File Manager & Cloud Storage (Google Drive / Windows File Explorer Hybrid):** File explorer tingkat lanjut dengan Windows-style Interactive Address Bar (Breadcrumbs interaktif, editable path text input, history navigation Back/Forward/Up, Copy Path, dan Refresh), inspeksi drive fisik riil (HDD/MicroSD), fluid auto-fill grid & list, bintang/favorit, serta standalone multi-window media viewers & players (Image zoom/rotate, Video player, Music visualizer player, PDF viewer, dan in-browser Code editor).
- **App Store 520+ Docker & Manual Custom Container:** Katalog 520+ aplikasi Docker siap pakai dengan halaman detail dokumentasi README.md, serta fitur Manual Install Docker Compose/Run ala CasaOS.
- **Task Manager & System Health:** Monitoring penggunaan sumber daya per-aplikasi secara real-time dengan kill-switch instan.
- **Pembersih Sistem 1-Klik:** RAM Cleaner dan Docker Storage Cleaner terintegrasi untuk menjaga kestabilan perangkat ber-RAM kecil.
- **Website Suite:** Deploy banyak website secara instan dengan alokasi port otomatis dan reverse proxy internal.
- **Database Suite:** MariaDB, PostgreSQL, dan GUI Adminer bawaan yang siap dikoneksikan ke website.
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
