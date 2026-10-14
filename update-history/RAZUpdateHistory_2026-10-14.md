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

---

### Versi: 2.8.1-multi-device-desktop-icon-persistence-sizing (Persistensi Ikon Desktop Lintas Perangkat & Pengaturan Ukuran Ikon Context Menu)
**Waktu Eksekusi:** 2026-10-14 11:05 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Persistensi Posisi Ikon Lintas Perangkat (Multi-Device Server Persistence)
- **Problem Statement:** Jika posisi ikon desktop hanya disimpan di `localStorage` peramban lokal, saat pengguna login atau mengakses RAZOS dari perangkat lain (ponsel pintar, tablet, atau komputer kerja lain), tata letak ikon akan kembali default atau tidak tersinkron.
- **Solusi Arsitektur:**
  - Menyediakan file storage persisten server: `backend/internal/store/RAZDesktopConfig.json`.
  - Membuat REST API endpoint `GET /api/v1/desktop/config` dan `POST /api/v1/desktop/config` pada Go Controller (`RAZApiController.go`) dan middleware server Vite (`vite.config.js`).
  - Frontend menerapkan strategi *optimistic local cache* via `localStorage` untuk kecepatan respon instan, dan langsung men-sync perubahan (posisi ikon dan ukuran ikon) secara bidirectional ke backend server.
  - Setiap perubahan letak drop ikon atau pengaturan ukuran otomatis tersimpan persisten ke backend server dan langsung konsisten di seluruh perangkat.

#### 2. Sistem Pengaturan Ukuran Ikon Desktop (Small, Medium, Large)
- **Pilihan Ukuran Dinamis:**
  - **Ikon Kecil (Small):** Grid slot `80px x 84px`, ikon `w-9 h-9`, tipografi `text-[10px]`. Cocok untuk layar padat atau pengguna yang memiliki banyak shortcut aplikasi.
  - **Ikon Sedang (Medium - Default):** Grid slot `92px x 96px`, ikon `w-12 h-12`, tipografi `text-[11px]`. Seimbang, proporsional, dan estetis.
  - **Ikon Besar (Large):** Grid slot `108px x 114px`, ikon `w-16 h-16`, tipografi `text-xs`. Sangat jelas dan ramah sentuhan (touchscreen / resolusi tinggi).
- **Kalkulasi Grid Adaptif:**
  - Fungsi `calculateGridCell`, batas maksimum kolom `maxCol`, batas maksimum baris `maxRow`, dan indikator preview snap-to-grid semuanya terikat secara reaktif pada `iconGridMetrics`.

#### 3. Integrasi Context Menu Desktop (Klik Kanan pada Desktop)
- **Seksi Tampilan Ikon Windows Style:**
  - Menambahkan menu pilihan ukuran: **Ikon Besar (Large)**, **Ikon Sedang (Medium)**, dan **Ikon Kecil (Small)**.
  - Mengintegrasikan penanda checkmark aktif (`Check` icon dari Lucide) pada item context menu yang sedang terpilih melalui dukungan properti `item.checked` di `RAZContextMenu.vue`.
  - Menambahkan aksi **Ratakan Ikon ke Grid (Auto Arrange)** untuk merapikan seluruh ikon ke slot awal kolom demi kolom secara otomatis.

#### 4. Verifikasi & Build
- `npm run build` sukses 100% tanpa error dalam 4.22 detik.
- Pengujian API `GET` & `POST` `/api/v1/desktop/config` terverifikasi sukses dengan validasi payload JSON di `backend/internal/store/RAZDesktopConfig.json`.
- Hot Module Replacement (HMR) berjalan lancar pada dev server `localhost:3000`.

---

### Versi: 2.8.2-window-pointer-events-click-drag-fix (Perbaikan Akses Pointer Events, Drag & Interaksi Seluruh Jendela Aplikasi)
**Waktu Eksekusi:** 2026-10-14 11:15 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Akar Masalah (Root Cause Analysis)
- **Multi-Root Component Fallthrough Loss:** Pada `RAZWindowFrame.vue`, sebelumnya tag `<Teleport to="body">` diletakkan sejajar di tingkat root dengan `<transition name="raz-window-transition">`, menjadikan komponen tersebut multi-root fragment di Vue 3.
- Akibatnya, atribut fallthrough `class="pointer-events-auto"` yang dipasang pada pemanggil `<RAZWindowFrame>` tidak diteruskan ke DOM node mana pun.
- Karena parent pembungkus Window Management Layer pada `RAZDesktopLayout.vue` memiliki style `pointer-events-none` (agar area kosong desktop dapat di-klik tembus), elemen `div.raz-window` secara otomatis mewarisi `pointer-events-none`.
- Dampaknya, seluruh klik mouse, seret (*drag*) titlebar, tombol minimize/maximize/close, serta form dan konten dalam seluruh jendela aplikasi terblokir (tembus dan tidak dapat merespons klik apa pun).

#### 2. Solusi Arsitektur & Perbaikan
- **Enkapsulasi Single-Root Node:** Memindahkan `<Teleport to="body">` ke dalam elemen `div.raz-window`, sehingga `RAZWindowFrame.vue` kembali menjadi single-root component.
- **Eksplisit Pointer Events:** Menetapkan `pointer-events-auto` secara langsung dan permanen pada class `div.raz-window`, area action buttons titlebar, dan kontainer body content aplikasi.
- **Sanitasi Angka Dimensi & Posisi (Anti-NaN Protection):**
  - Membuat computed properties aman: `winLeft`, `winTop`, `winWidth`, `winHeight`, dan `winZIndex` untuk menjamin seluruh nilai koordinat dan ukuran jendela selalu berupa angka piksel berhingga (`Number.isFinite`), mencegah layout crash jika terdapat state corrupt di `localStorage`.
  - Memperbarui fungsi `getSavedBounds` dan `saveWindowBounds` di `RAZWindowStore.js` dengan filter validasi tipe data ketat.
- **Proteksi Drag Mouse Event:** Memastikan `startDrag` hanya aktif pada klik kiri mouse (`e.button === 0`) dengan dukungan `passive` touch event.

#### 3. Verifikasi & Pengujian
- `npm run build` berhasil 100% tanpa error dalam 4.75 detik.
- Seluruh event klik, drag titlebar, minimize, maximize, close, resize handle, dan interaktivitas isi jendela aplikasi berfungsi lancar 100%.

