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

---

### Versi: 2.7.1-desktop-layout-right-widgets-fix (Penataan Widget Sebelah Kanan & Perbaikan Interaktivitas Search Bar)
**Waktu Eksekusi:** 2026-10-14 10:25 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Penataan Ulang Posisi Widget ke Sebelah Kanan Layar
- Memindahkan kontainer kartu **Desktop Status Widgets** (Jam & Tanggal, Status CPU & Suhu/RAM, Storage /DATA, dan Kecepatan Network) ke sisi kanan layar (`ml-auto` / `justify-between`).
- Sisi kiri layar kini sepenuhnya didedikasikan untuk **Grid Ikon Aplikasi Desktop** yang bersih, rapi, dan leluasa.
- Mengatur `pointer-events-none` pada area kosong antarmuka dan `pointer-events-auto` pada elemen interaktif, sehingga klik kanan pada wallpaper kosong tetap memicu Context Menu Desktop tanpa terhalang.

#### 2. Perbaikan Interaktivitas Floating Search Bar
- **Penyebab:** Sebelumnya search bar berada di tumpukan z-index yang tertutup oleh kontainer flex desktop icons sehingga event klik dan input keyboard terhalang (unclickable).
- **Solusi:** 
  - Menyusun search bar di header teratas alur DOM (`w-full flex justify-center pb-3 pointer-events-auto`) tanpa ada elemen lain yang menumpuk.
  - Menambahkan styling `select-text` dan `cursor-text` pada input pencarian agar pengguna dapat dengan mudah mengklik, mengetik kueri, dan menekan `Enter` atau tombol mesin telusur untuk langsung membuka hasil pencarian ke tab baru.
- **Verifikasi:**
  - `npm run build` sukses 100% dalam 4.12 detik.
  - Hot Module Replacement (HMR) aktif di dev server `localhost:3000`.

---

### Versi: 2.7.2-desktop-icon-titles-fix (Perbaikan Visibilitas & Spasi Label Judul Ikon Desktop)
**Waktu Eksekusi:** 2026-10-14 10:48 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Masalah Terpotongnya Label Judul Ikon Desktop
- **Penyebab:** Pada layar dengan tinggi viewport tertentu, grid CSS `grid-rows-6` (`repeat(6, minmax(0, 1fr))`) membagi sisa tinggi ruang desktop menjadi 6 baris tipis (~60px). Ketinggian tersebut hanya cukup untuk kotak ikon saja, sehingga teks judul (`app.title`) tertekan keluar dan terpotong (hilang/tidak terlihat).
- **Solusi:**
  - Mengubah aturan grid row menjadi `grid-rows-[repeat(6,minmax(84px,max-content))]` dengan `gap-x-4 gap-y-3`.
  - Memberikan batas minimum ketinggian kartu item `min-h-[84px]` sehingga setiap slot ikon selalu memiliki ruang vertikal yang cukup untuk ikon (`w-12 h-12`) dan 2 baris teks nama aplikasi.
  - Menambahkan styling kontras tinggi `drop-shadow-[0_2px_4px_rgba(0,0,0,0.95)]` pada `span` judul aplikasi agar selalu terbaca tajam dan jelas di atas wallpaper apa pun.
- **Verifikasi:**
  - `npm run build` berhasil 100% dalam 4.20 detik.
  - Hot Module Replacement (HMR) langsung ter-update di `localhost:3000`.

---

### Versi: 2.8.0-windows-grid-dnd-snap-maximize (Windows Desktop Grid & Drag-Drop, Windows 11 Snap Maximize & Window Bounds Persistence)
**Waktu Eksekusi:** 2026-10-14 10:55 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Windows Desktop Grid System & Drag-and-Drop dengan Snap Target
- **Sistem Grid Koordinat Desktop:**
  - Menerapkan sistem penataan grid slot koordinat `(col, row)` ala Windows Desktop (`CELL_W = 96px`, `CELL_H = 100px`).
  - Setiap ikon aplikasi diposisikan secara mutlak (`absolute`) berdasarkan koordinat slot grid.
- **Drag-and-Drop Interaktif:**
  - Ikon desktop dapat diklik dan diseret (`drag-and-drop`) bebas ke slot grid mana pun di desktop.
  - Menampilkan animasi indikator kotak target snap semi-transparan berbingkai putus-putus (`border-dashed border-sky-400 bg-sky-500/20`) saat menyeret ikon untuk memperlihatkan posisi landing secara presisi.
  - Jika slot tujuan sudah terisi, sistem secara otomatis menukar posisi (*swap positions*) dengan ikon yang ada persis seperti Windows desktop.
- **Persistensi Posisi Ikon Desktop:**
  - Seluruh posisi ikon tersimpan secara permanen ke `localStorage` (`razos_desktop_icon_positions`).
  - Saat reload peramban atau login kembali, tata letak ikon desktop tetap berada pada posisi terakhirnya.

#### 2. Windows 11 Snap-to-Top Maximize Gesture & Animasi Halo Preview
- **Gesture Drag Mentok ke Atas:**
  - Saat title bar jendela aplikasi diseret dan mentok ke bagian paling atas layar (`clientY <= 12px`), muncul animasi halo preview fullscreen semi-transparan berpendar biru (`bg-sky-500/10 border-sky-400/80 shadow-[0_0_40px_rgba(56,189,248,0.4)] animate-pulse`) yang di-teleport langsung ke body (Windows 11 Snap style).
  - Ketika mouse dilepaskan, jendela otomatis beranimasi menjadi **Maximize**.
  - Menyeret title bar jendela yang sedang maximized ke bawah otomatis mengembalikan jendela ke ukuran semula (*unmaximize on pull down*) dengan kursor berada tepat di tengah titlebar.

#### 3. Persistensi Ukuran & Posisi Jendela per Aplikasi
- **Penyimpanan Ukuran Terakhir:**
  - Ukuran jendela (`width`, `height`) serta koordinat posisi (`x`, `y`) disimpan secara otomatis ke `localStorage` (`razos_window_bounds`) setiap kali jendela selesai dipindahkan, di-resize, atau ditutup (`closeWindow`).
  - Saat aplikasi dibuka kembali di kemudian hari, jendela akan langsung muncul dengan ukuran dan posisi yang sama persis seperti saat terakhir digunakan.

#### 4. Verifikasi & Build
- `npm run build` sukses 100% tanpa error dalam 4.23 detik.
- Hot Module Replacement (HMR) aktif di dev server `localhost:3000`.
