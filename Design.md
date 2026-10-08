# Spesifikasi Desain UI/UX RAZOS

## 1. Filosofi Desain
Mengadopsi kombinasi kesederhanaan desktop modern (Fluent/macOS) pada desktop dan Material You (Android 14) pada mobile launcher.

## 2. Tema & Palet Warna
- **Latar Belakang:** `#0f172a` (Slate 900) dengan efek glassmorphism (`backdrop-blur-md bg-slate-900/80`).
- **Aksen Utama:** `#3b82f6` (Blue 500) untuk aksi primer dan status aktif.
- **Status Bahaya / Kill:** `#ef4444` (Red 500) untuk RAM tinggi dan tombol terminasi.
- **Teks:** `#f8fafc` (Slate 50) dengan subteks `#94a3b8` (Slate 400).

## 3. Perilaku Komponen

### Desktop Mode (> 768px)
- **Taskbar:** Menempel di bawah layar, tinggi 48px.
  - Sisi Kiri: Tombol Start Menu / App Drawer.
  - Sisi Tengah: Task items aplikasi aktif yang sedang dibuka.
  - Sisi Kanan: System tray (Beban CPU, RAM, Suhu STB, Jam, Tombol Fullscreen).
- **Window Frame:**
  - Header window memuat judul, tombol Minimize, Maximize, dan Close.
  - Mendukung interaksi seret (*drag*) dan ubah ukuran (*resize*) di 4 sudut.

### Mobile Mode (<= 768px)
- **Top Bar:** Status bar minimalis dengan indikator waktu, suhu, dan tombol Fullscreen.
- **Home Screen:** Grid 4 kolom untuk ikon aplikasi dengan widget ringkasan server di atas.
- **Dock:** 4 aplikasi utama + 1 tombol Drawer di bagian bawah layar.
- **App Drawer:** Modal geser dari bawah ke atas menampilkan daftar lengkap aplikasi dengan filter pencarian.
- **App Switcher (Multitasking):** Tampilan kartu horizontal aplikasi yang berjalan dengan tombol "Tutup Semua".
