# RAZOS Update History: 2026-10-14

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 2.7.0-casaos-settings-suite (Integrasi Menu Settings Lengkap CasaOS-Style, Manajemen Akun, & Desktop Widgets)
**Waktu Eksekusi:** 2026-10-14 09:45 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Implementasi 3 Tab Navigasi Header Settings (CasaOS Standard)
- **Tab 1: Account (Manajemen Profil & Kredensial Pengguna)**
  - Mengintegrasikan header profil dengan tombol pill `Logout` beraksen rose.
  - Avatar lingkaran profil pengguna beresolusi tinggi dengan indikator badge aktif.
  - Field kredensial dinamis:
    - **Name:** Menampilkan username saat ini dengan tombol aksi edit (ikon pensil) yang memicu dialog modal kustom.
    - **Password:** Menampilkan masked password (`••••••••`) dengan tombol aksi edit (ikon pensil) untuk penggantian kata sandi baru secara aman.
- **Tab 2: Settings (Preferensi Sistem & OTA Updater)**
  - **Show search bar:** Toggle switch untuk mengaktifkan/menonaktifkan floating web search bar di desktop host.
  - **Search Engine:** Dropdown pemilihan mesin telusur default (DuckDuckGo, Google, Bing, Baidu).
  - **Language:** Dropdown konfigurasi bahasa antarmuka (English, Bahasa Indonesia).
  - **WebUI Port:** Menampilkan port WebUI saat ini (`8080`) dengan tombol aksi "Change" yang memicu modal kustom validasi port.
  - **Wallpaper:** Thumbnail preview wallpaper aktif dengan tombol aksi "Change" yang memicu modal galeri wallpaper preset (Cosmic, Gradient, Cyber, Minimalist) atau custom URL.
  - **Show news feed from RAZOS blog:** Toggle switch preferensi feed komunitas.
  - **Show recommended apps:** Toggle switch untuk rekomendasi katalog aplikasi.
  - **Automount USB drive:** Toggle switch automount media penyimpanan eksternal.
  - **OTA Updater Card:** Kartu status versi sistem (`v2.6.0`), indikator "Currently, at the latest version ✓", tombol "Check Update", dan "Update Now".
  - **Footer Actions:** Tombol daya **Restart** (reboot aman) dan **Shutdown** (poweroff aman) dengan dialog modal konfirmasi.
- **Tab 3: Widget settings (Konfigurasi Widget Desktop Real-Time)**
  - **Time:** Toggle switch widget jam digital presisi dan tanggal bahasa Indonesia.
  - **System status:** Toggle switch widget status real-time beban CPU (persentase & suhu SoC STB) serta beban memori RAM.
  - **Storage status:** Toggle switch widget ruang disk penyimpanan utama (/DATA MicroSD).
  - **Network status:** Toggle switch widget throughput kecepatan jaringan (Download & Upload real-time).

#### 2. Reaktivitas Frontend & Desktop Layout Integration
- **Persistensi State (`RAZSystemStore.js`):**
  - Seluruh konfigurasi pengguna (toggles, search engine, wallpaper, port, widgets, info akun) tersimpan persisten ke `localStorage` (`razos_system_settings`).
  - Menambahkan kalkulasi kecepatan throughput jaringan (`net_upload_speed` dan `net_download_speed`) berbasis selisih byte terhadap waktu polling.
- **Desktop Layout Enhancement (`RAZDesktopLayout.vue`):**
  - Background desktop terikat secara reaktif ke `systemStore.settings.wallpaper`.
  - Floating Search Bar di bagian atas desktop yang langsung membuka kueri pencarian ke tab baru sesuai search engine pilihan pengguna saat di-submit.
  - Floating Widget Card di sebelah kiri desktop yang menampilkan Time, System CPU/RAM, Storage /DATA, dan Network I/O sesuai preferensi toggle yang aktif di Settings.

#### 3. Backend System Endpoints
- **Modul Pemeliharaan Host (`RAZMaintenance.go`):**
  - Menambahkan fungsi `RestartHost()` yang mengeksekusi reboot aman pada Linux STB.
  - Menambahkan fungsi `ShutdownHost()` yang mengeksekusi poweroff aman pada Linux STB.
- **API Controller (`RAZApiController.go`):**
  - Mendaftarkan endpoint `POST /api/v1/system/restart` dan `POST /api/v1/system/shutdown`.
- **Metrics Collector (`RAZMetricsCollector.go`):**
  - Menambahkan penangkapan metrik I/O jaringan `NetSent` dan `NetRecv` dari gopsutil.

#### 4. Kepatuhan Aturan & Standar Proyek (GEMINI.md)
- **Zero Browser Alerts:** Seluruh dialog interaksi (Edit Name, Edit Password, Change Port, Change Wallpaper, Confirm Restart, Confirm Shutdown, Confirm Logout) menggunakan Custom Modal Glassmorphism Slate-900.
- **Zero Dummy Data:** Terhubung dinamis ke `RAZSystemStore` dan backend API.
- **Backup Aman:** Seluruh file sebelum modifikasi telah dicadangkan di folder `backup/`.
- **Verifikasi Build:** `npm run build` berhasil 100% tanpa error kompilasi.
