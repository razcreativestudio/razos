# Product Requirements Document (PRD) - RAZOS

## 1. Deskripsi Produk
RAZOS adalah portal sistem operasi web mandiri untuk server rumahan dan mini-VPS, menggantikan dashboard konvensional dengan antarmuka berbasis desktop windowing dan mobile launcher adaptif.

## 2. Persona Pengguna
Administrator server rumahan (self-hoster) yang menggunakan perangkat keras single-board/STB dan menginginkan kendali komprehensif atas aplikasi, website, dan database tanpa dependensi yang menguras RAM.

## 3. Persyaratan Fungsional
- **Multi-Window Management:** Pengguna dapat membuka lebih dari satu aplikasi sekaligus di layar desktop.
- **Adaptive Mode Switching:** Transisi otomatis ke launcher bergaya Android pada viewport layar sentuh.
- **Task Management & Kill Switch:** Visualisasi konsumsi resource per-container dengan opsi terminasi seketika.
- **Website Provisioning:** Pembuatan container web mandiri dengan isolasi port dan koneksi database instan.
- **Database Suite:** Akses langsung ke database server dan antarmuka manajemen Adminer tanpa konfigurasi manual.
- **Zero-Config Web Terminal:** Akses baris perintah Linux langsung dari peramban secara aman.
- **System Maintenance:** Pembersih RAM dan penyimpanan otomatis dengan satu tombol.

## 4. Persyaratan Non-Fungsional
- **Efisiensi:** Konsumsi RAM backend statis <= 35 MB.
- **Ketersediaan:** Auto-restart via systemd saat terjadi kegagalan (*crash*).
- **Keandalan:** Menolak instalasi container jika kapasitas RAM bebas < 150 MB.
