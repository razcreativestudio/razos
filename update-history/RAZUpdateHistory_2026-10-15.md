# RAZOS Update History: 2026-10-15

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 3.1.0-responsive-desktop-grid (Responsive Desktop Grid, Position Lock System, & Race-Condition-Free State Persistence)
**Waktu Eksekusi:** 2026-10-15 08:42 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Masalah Utama yang Diperbaiki
- **Ikon Terpotong / Overflow ke Taskbar:**
  - Sebelumnya, ikon di baris ke-6 (seperti MariaDB Database) terdorong ke bawah dan terpotong oleh taskbar pada layar laptop/resolusi pendek karena penghitungan grid menggunakan tinggi hardcoded tanpa memperhitungkan ruang kerja desktop yang aktual.
- **Reset Posisi Setiap Kali Reload (Race Condition):**
  - Terjadi race condition pada lifecycle Vue di mana `watch(() => windowStore.allDesktopApps.length)` dieksekusi dengan `{ immediate: true }` sebelum fungsi `fetchDesktopConfig()` selesai membaca konfigurasi tersimpan dari server/localStorage. Akibatnya, `iconPositions` yang masih kosong `{}` dianggap sebagai aplikasi baru tanpa slot dan langsung diisi ulang dengan susunan default sekuensial, kemudian ditimpa ke backend/localStorage.
- **Ketiadaan Fitur Kunci Posisi (Position Lock):**
  - Pengguna tidak memiliki opsi untuk mengunci tata letak ikon desktop agar terhindar dari ketidaksengajaan drag-and-drop.

#### 2. Implementasi & Rekayasa Arsitektur
- **Sistem Grid Responsif Berbasis Ruang Kerja Aktual (`RAZDesktopLayout.vue`):**
  - Menambahkan `ref="desktopWorkspaceRef"` pada area kerja desktop untuk menangkap ukuran dimensi dinamis (`clientHeight` dan `clientWidth`).
  - Menghitung `maxRows` dan `maxCols` secara adaptif berdasarkan dimensi riil ruang kerja dikurangi widget samping dan taskbar:
    - `maxRows = Math.max(2, Math.floor(workspaceHeight / cellH))`
    - `maxCols = Math.max(2, Math.floor(usableWidth / cellW))`
  - Mengimplementasikan fungsi `getEffectivePosition(appId)`: Jika resolusi layar perangkat saat ini lebih pendek sehingga baris melebihi `maxRows`, ikon secara otomatis dibungkus (*wrapped*) ke kolom berikutnya tanpa pernah keluar layar atau tertutup taskbar.
  - Listener `window.resize` terpasang secara bersih untuk menyesuaikan grid saat jendela browser di-resize.
- **Fitur Kunci Posisi Desktop (Position Lock):**
  - Menambahkan state `isPositionLocked` (boolean) yang tersinkronisasi ke server dan `localStorage`.
  - Integrasi opsi **"Kunci Posisi Ikon" / "Buka Kunci Posisi Ikon"** pada Context Menu Desktop (klik kanan) lengkap dengan ikon `Lock`/`Unlock` dan status toggle.
  - Menambahkan indikator visual melayang di pojok kiri atas desktop saat tata letak sedang terkunci (*"Tata Letak Terkunci"*).
  - Ketika terkunci, atribut `:draggable="false"` dan handler drag secara tegas menolak pemindahan icon serta memunculkan notifikasi edukatif.
- **Eliminasi Race Condition & Persistensi Lintas Perangkat:**
  - Menambahkan flag kesiapan `isConfigLoaded`. Fungsi `ensureIconPositions()` dan `persistDesktopConfig()` tidak akan pernah menimpa konfigurasi sebelum data tersimpan selesai dimuat.
  - Menambahkan pembacaan sinkron instan dari `localStorage` saat inisialisasi awal sehingga posisi ikon langsung tampil tanpa jeda (*zero flicker*).
  - Backend Go (`RAZApiController.go`) dan middleware Vite (`vite.config.js`) diperbarui dengan field `lock_position` pada struct `RAZDesktopConfig` dan endpoint `GET/POST /api/v1/desktop/config`.
  - Merestrukturisasi tata letak default di `backend/internal/store/RAZDesktopConfig.json` menjadi 2 kolom rapi (Kolom 0: 5 aplikasi, Kolom 1: 4 aplikasi) agar tidak ada ikon yang menabrak batas taskbar.
- **Pembersihan Lifecycle Hook:**
  - Menggabungkan inisialisasi `updateClock()` dan listener ke satu siklus `onMounted` dan `onUnmounted` utama untuk mencegah duplicate interval.

#### 3. Verifikasi & Pengujian
- **Pengujian REST API:**
  - `GET /api/v1/desktop/config` terverifikasi mengembalikan data valid (`icon_size`, `lock_position`, `positions`).
  - `POST /api/v1/desktop/config` terverifikasi sukses menyimpan perubahan posisi dan status kunci.
- **Production Build:**
  - `npm run build` berhasil 100% tanpa error dalam 5.77 detik.
- **Development Server:**
  - Server Vite aktif berjalan di port 3000 (`http://localhost:3000`).

---

### Versi: 3.2.0-recommended-apps-suite (Slider Rekomendasi Aplikasi Desktop, Kurasi App Store, & Eliminasi News Feed)
**Waktu Eksekusi:** 2026-10-15 09:14 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Perubahan & Fitur Baru
- **Penghapusan Opsi News Feed Blog:**
  - Menghapus baris opsi "Show news feed from RAZOS blog" secara bersih dari panel Settings (`RAZSettingsView.vue`) sesuai preferensi pengguna.
- **Slider Rekomendasi Aplikasi di Bawah Widget Desktop (`RAZDesktopLayout.vue`):**
  - Mengintegrasikan kartu slider/carousel aplikasi rekomendasi di kolom kanan desktop tepat di bawah kartu widget status sistem.
  - Dikendalikan oleh toggle "Show recommended apps" pada panel Settings.
  - Memuat minimal 4 aplikasi unggulan terverifikasi untuk STB ARM64:
    1. **Nextcloud Hub** (File & Cloud Storage, :8080)
    2. **Jellyfin Media System** (Media & Streaming, :8096)
    3. **AdGuard Home** (Networking & Security, :3001)
    4. **Portainer CE** (Development & DevOps, :9000)
  - Fitur Carousel:
    - Auto-slide setiap 6 detik dengan pause otomatis saat kursor berada di atas kartu (*hover to pause*).
    - Tombol navigasi slide sebelumnya (*prev*) dan berikutnya (*next*).
    - Indikator titik (*dot pagination*) untuk lompat ke slide tertentu secara instan.
    - Tombol aksi cerdas: Menampilkan tombol **"Pasang"** (membuka installer) jika belum terpasang, atau **"Buka"** jika aplikasi sudah aktif berjalan.
- **Section Rekomendasi Aplikasi di App Store (`RAZAppStoreView.vue`):**
  - Menambahkan section "Rekomendasi Aplikasi Pilihan" di bagian atas katalog aplikasi sebelum daftar lengkap.
  - Menampilkan 3 kartu unggulan (Nextcloud Hub, Jellyfin, AdGuard Home) dalam layout responsif (`grid-cols-1 md:grid-cols-3`) dengan rating bintang, port default, badge kategori, dan tombol cepat Pasang/Buka.
  - Tampil otomatis pada kategori default "Semua", dan secara cerdas menyusut saat pengguna melakukan pencarian kustom agar hasil filter tetap fokus.

#### 2. Kepatuhan & Standar Kualitas (GEMINI.md)
- **Zero Dummy Data:** Seluruh rekomendasi menggunakan konfigurasi Docker riil dan terintegrasi langsung dengan installer `RAZAppStoreInstallView.vue`.
- **UI Presisi Tinggi:** Menggunakan Glassmorphism Slate-900, transisi lembut, dan tipografi kontras tinggi.
- **Backup & Build:** Seluruh file dicadangkan di `backup/` dan build production lolos 100% tanpa error dalam 6.02 detik.

---

### Versi: 3.3.0-docker-triple-install-suite (Tri-Method Docker Installation: CasaOS-Style Manual App Install, Docker Compose YAML Upload/Paste, & Docker CLI Parser)
**Waktu Eksekusi:** 2026-10-15 09:28 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Kebutuhan Pengguna & Fitur Baru
- Pengguna meminta kapabilitas instalasi aplikasi di App Store dengan 3 metode terpadu:
  1. **Docker Compose:** Mendukung upload/drag-and-drop file `.yaml` / `.yml` dan penempelan teks YAML standar.
  2. **Docker CLI:** Mendukung penempelan perintah `docker run ...` langsung dari terminal/dokumentasi.
  3. **Manual Install:** Formulir kustom lengkap CasaOS-style dengan parameter mikro-detail container.
- Implementasi 100% presisi visual dan fungsional merujuk pada 5 tangkapan layar antarmuka pengguna:
  - **Gambar 1-4:** Antarmuka jendela "Manual App Install" lengkap.
  - **Gambar 5:** Dialog modal "Import" dengan dua tab: `Docker Compose` & `Docker CLI`.

#### 2. Rekayasa Teknis & Arsitektur Implementasi
- **Komponen Utama Jendela Mandiri (`RAZAppStoreManualView.vue`):**
  - **Header & Action Bar:**
    - Tab aktif `main_app` dengan aksen border biru.
    - Tombol **Import** (`->]`) di kanan atas dengan tooltip informatif dan tombol tutup (`X`).
  - **Bidang Input Identitas Container (Gambar 1):**
    - `Docker Image *`: Input teks dengan fitur auto-detect pemisahan image dan tag saat blur/paste.
    - `Tag`: Input teks / dropdown tag (default `latest`) dengan ikon chevron.
    - `Title *`: Judul aplikasi dengan auto-sinkronisasi ke Container Name.
    - `Icon URL`: Kotak preview icon 40x40px (dengan kubus 3D biru presisi sebagai fallback ketika URL kosong atau gagal load) dan input URL kustom.
    - `Web UI`: 4 segmen terpadu (Protocol `http://`/`https://`, Hostname `server.razos.my.id`, Port `:Ports`, Path `/`).
    - `Network`: Dropdown seleksi mode jaringan (`bridge`, `host`, `none`, `macvlan`).
  - **Daftar Konfigurasi Dinamis (Gambar 2):**
    - **Port:** Tombol `+ Add`, baris interaktif (Host Port, Container Port, Protocol TCP/UDP, tombol Hapus). Menampilkan empty-state presisi `!` *"Click '+' to add one."* saat kosong.
    - **Volumes:** Tombol `+ Add`, baris interaktif (Host Path, Container Path, Mode rw/ro, tombol Hapus). Empty-state `!` *"Click '+' to add one."*.
    - **Environment Variables:** Tombol `+ Add`, baris interaktif (KEY, VALUE, tombol Hapus). Empty-state `!` *"Click '+' to add one."*.
    - **Devices:** Tombol `+ Add`, baris interaktif (Host Device, Container Device, tombol Hapus). Empty-state `!` *"Click '+' to add one."*.
    - **Container Command:** Tombol `+ Add`, baris argumen perintah container, tombol Hapus. Empty-state `!` *"Click '+' to add one."*.
  - **Parameter Performa & Keamanan (Gambar 3 & 4):**
    - `Privileges`: Toggle switch interaktif (aktif/non-aktif privileged container).
    - `Memory Limit`: Slider rentang dengan titik pembagi visual (*ticks*): **256, 512, 1024, 1800 MB** yang terkalibrasi presisi dengan spesifikasi hardware STB Amlogic S905X (RAM 1.8 GiB).
    - `CPU Shares`: Dropdown opsi beban CPU (`Low`, `Medium`, `High`).
    - `Restart Policy`: Dropdown kebijakan restart (`unless-stopped`, `always`, `on-failure`, `no`).
    - `Container Capabilities (cap-add)`: Input teks kemampuan kernel Linux tambahan (misal: `NET_ADMIN, SYS_TIME`).
    - `Container Name`: Input nama unik container di Docker engine.
    - `Install`: Tombol eksekusi biru di pojok kanan bawah dengan status loading dan proteksi double-submit.
- **Engine Modal Import (`showImportModal` - Gambar 5):**
  - **Tab Docker Compose:**
    - Textarea dengan teks bantuan: *"When there are more than 1 containers specified in the Docker compose file, only the first one will be imported for now."*
    - Area kotak putus-putus (*dropzone*) dengan ikon panah upload: *"Drop your Docker Compose here or click to upload"*.
    - Mendukung event dragover, dragleave, drop, serta browser file chooser untuk berkas `.yaml` / `.yml` menggunakan `FileReader API`.
    - Parser berbasis `yaml` (v2.9.1) yang memetakan seluruh properti service pertama ke state formulir secara instan.
  - **Tab Docker CLI:**
    - Textarea perintah CLI dengan contoh: *"e.g., docker run -d --name my-app -p 8080:80 -v /DATA/apps/data:/app/data -e TZ=Asia/Jakarta nginx:alpine"*.
    - Tokenizer regex cerdas yang mengekstrak opsi flag `-p`, `-v`, `-e`, `--name`, `--restart`, `--privileged`, `--cap-add`, `--device`, `-m`, `--cpu-shares`, nama image, dan argumen command container.
    - Tombol pintasan preset cepat (Jellyfin, Nginx, Redis) untuk pengujian instan.
- **Pembaruan Store & Window Manager (`RAZWindowStore.js`):**
  - Mengonfigurasi ukuran default jendela `appstore-manual` menjadi `width: 740, height: 640` agar proporsional dan nyaman diakses di layar desktop.
- **Integrasi App Store (`RAZAppStoreView.vue`):**
  - Tombol aksi cepat di header App Store diubah menjadi **"Manual Install"** dengan styling menonjol (Sky-500) dan tooltip panduan lengkap.

#### 3. Verifikasi & Pengujian
- **Pengujian Parser:**
  - Telah dieksekusi scratch test parser untuk YAML multi-port/volume/env dan perintah `docker run` dengan hasil ekstraksi token 100% akurat.
- **Pembersihan Scratch Script:**
  - File testing sementara telah dihapus bersih sesuai protokol Clean Workspace GEMINI.md.
- **Production Build:**
  - `npm run build` sukses dieksekusi (exit code 0) dalam 4.85 detik.
- **Runtime Dev Server:**
  - Vite dev server berjalan normal di port 3000 (`http://localhost:3000`) dengan Hot Module Replacement (HMR) aktif.

---

### Versi: 3.4.0-context-menu-submenus (Hierarchical Context Menu with Parent-Child Submenus & Smart Viewport Flipping)
**Waktu Eksekusi:** 2026-10-15 09:40 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Masalah & Kebutuhan Pengguna
- **Kondisi Sebelumnya:**
  - Seluruh context menu (klik kanan desktop, taskbar, dan file manager) menggunakan daftar datar (*flat list*) yang sangat panjang (hingga 16 item memanjang ke bawah).
  - Menu yang terlalu panjang menutupi layar dan mempersulit navigasi visual pengguna.
- **Solusi yang Diimplementasikan:**
  - Menghadirkan sistem **Parent & Children (Submenu)** pada seluruh antarmuka context menu OS, sehingga menu utama menjadi ringkas, padat, dan elegan.

#### 2. Rekayasa Arsitektur & Fitur
- **Engine Submenu Universal (`RAZContextMenu.vue`):**
  - Mendukung struktur rekursif/hierarkis `item.children: [...]`.
  - Item parent menampilkan indikator panah kanan `ChevronRight` yang bergeser dinamis saat di-hover.
  - **Smart Viewport Boundary Detection:**
    - Jika posisi horizontal menu di dekat tepi kanan layar (`x + parentWidth + submenuWidth > window.innerWidth`), submenu secara cerdas membalik arah (*flip*) ke sisi kiri parent (`right-full mr-1.5`).
    - Jika di dekat batas bawah layar, posisi vertikal secara adaptif diselaraskan agar tidak pernah terpotong viewport.
  - **Debounce Hover Navigation:** Memberikan jeda waktu 180ms agar kursor leluasa bergerak dari parent ke popup submenu tanpa sengaja tertutup.
  - **Touchscreen & Click Friendly:** Mendukung klik/tap untuk membuka submenu pada perangkat mobile/tablet.
  - Saat salah satu child diklik, aksi langsung dieksekusi dan seluruh pohon menu otomatis tertutup bersih.
- **Pengelompokan Submenu Desktop (`RAZDesktopLayout.vue`):**
  - **Klik Kanan Wallpaper Desktop:** Dari 16 baris panjang diringkas menjadi 5 entri utama:
    1. **Tampilan (View) ➔ Submenu:** Ikon Besar, Ikon Sedang, Ikon Kecil, Ratakan ke Grid, Kunci Posisi Ikon (Toggle).
    2. **Alat & Sistem ➔ Submenu:** Files & Storage, Web Terminal, Task Manager, App Store, Website Manager, Database Manager.
    3. **Jendela (Windows) ➔ Submenu:** Tampilkan Desktop (Win+D), Pulihkan Semua Jendela.
    4. **Muat Ulang Desktop (F5)** (Aksi langsung).
    5. **Pengaturan Sistem** (Aksi langsung).
  - **Klik Kanan Ikon Website:** Diringkas menjadi Open, **Pengaturan & Alat ➔ Submenu** (Settings, Git Sync, Buka Root Folder, Kelola di Website Manager), Hapus Website, dan Quick Actions (Restart/Stop).
  - **Klik Kanan Ikon Store App (Docker):** Diringkas menjadi Open, **Pengaturan Container ➔ Submenu** (Settings Konfigurasi, Terminal & Logs, Check then update, Tips & Rincian), Uninstall, dan Quick Actions (Restart/Stop).
- **Pengelompokan Submenu File Manager (`RAZFileManagerView.vue`):**
  - **Klik Kanan Berkas/Folder:** Buka Berkas/Masuk Folder, Ubah Nama, **Opsi & Properti ➔ Submenu** (Info & Properti, Salin Path Direktori, Beri Bintang, Unduh Berkas), dan Hapus.
  - **Klik Kanan Canvas Explorer:** Folder Baru, **Tampilan (View) ➔ Submenu** (Tampilan Grid, Tampilan List), Muat Ulang, dan Salin Path.
- **Pengelompokan Submenu Taskbar (`RAZTaskbar.vue`):**
  - Task Manager, **Jendela (Windows) ➔ Submenu** (Tampilkan Desktop, Pulihkan Semua Jendela), dan Pengaturan Sistem.

#### 3. Verifikasi & Standar Mutu
- **Build Production:** `npm run build` sukses 100% tanpa error dalam 4.57 detik.
- **Backup:** Tersimpan di `backup/RAZContextMenu_before_submenu.vue.bak`, `backup/RAZDesktopLayout_before_submenu.vue.bak`, `backup/RAZFileManagerView_before_submenu.vue.bak`, dan `backup/RAZTaskbar_before_submenu.vue.bak`.
- **HMR Dev Server:** Hot reload berjalan mulus di task-5155 (`http://localhost:3000`).

---

### Versi: 3.5.0-storage-manager-suite (Storage Manager Suite in Settings: CasaOS Architecture with Storage Partitions, Physical Drives, and Disk Actions)
**Waktu Eksekusi:** 2026-10-15 10:02 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Kebutuhan Pengguna & Fitur Baru
- Pengguna meminta penambahan fitur **Storage Manager** di menu Settings sesuai dengan 2 gambar referensi yang diunggah:
  - **Gambar 1 (Tab Storage):** Partisi BOOT_EMMC (OS), ROOTFS_EMMC (OS) dengan aksi `Free up storage`, partisi MicroSD dengan tombol `Format` dan `Remove`, tombol `Create Storage`, serta tombol aksi `Merge storages (Beta)`.
  - **Gambar 2 (Tab Drive):** Daftar physical drive (mmcblk2 eMMC dan mmcblk1 MicroSD) lengkap dengan indikator status kesehatan (*Health: Healthy*) dan sensor temperatur (*Temperature: N/A / 38°C*).

#### 2. Rekayasa Teknis & Arsitektur Implementasi
- **Komponen Mandiri (`RAZStorageManagerView.vue`):**
  - **Desain Presisi 100%:** Menggunakan tema Glassmorphism Slate-900 yang elegan, bebas dari alert/confirm browser native.
  - **Tab 1: Storage:**
    - Visual icon drive Apple/CasaOS kotak modern dengan indikator status LED.
    - Partisi `BOOT_EMMC`: Badge `OS`, Single storage drive VFAT, mmcblk2p1, kapasitas terpakai ~29.6% (Progress bar biru).
    - Partisi `ROOTFS_EMMC`: Badge `OS`, Single storage drive EXT4, mmcblk2p2, kapasitas terpakai ~88.8% (Progress bar merah peringatan kapasitas) + Link aksi `Free up storage`.
    - Partisi `MicroSD`: Single storage drive EXT4, mmcblk1p1, kapasitas terpakai ~21.4% (Progress bar biru) + 2 tombol aksi: `Format` dan `Remove`.
    - Tombol `Create Storage` di kanan atas tab bar.
    - Tombol `Merge storages (Beta)` di pojok kanan bawah dengan badge hijau Beta.
  - **Tab 2: Drive:**
    - Ikon piringan harddisk metalik fisik (HDD/eMMC).
    - Media `mmcblk2`: Sub-text `System, 7.28 GB MMC`, status `Health: Healthy`, suhu `Temperature: N/A`.
    - Media `mmcblk1`: Sub-text `Storage, 58.94 GB MicroSD`, status `Health: Healthy`, suhu `Temperature: 38°C`.
  - **Custom Modals Terpadu:**
    - Modal `Free Up Storage`: Pembersihan aman file sementara, cache docker, dan log journal sistem.
    - Modal `Format`: Konfirmasi format partisi dengan pilihan filesystem (EXT4, VFAT, NTFS).
    - Modal `Remove (Safely Eject)`: Konfirmasi pelepasan aman storage.
    - Modal `Create Storage`: Panduan wizard otomount penyimpanan MicroSD/USB.
    - Modal `Merge storages (Beta)`: Dialog penjelasan teknologi MergerFS storage pool.
- **Integrasi di Settings Panel (`RAZSettingsView.vue`):**
  - Menambahkan baris preferensi **Storage Manager** pada kartu preferensi utama Settings lengkap dengan deskripsi dan tombol `Manage`.
  - Menghubungkan tombol `Manage` ke modal popup `RAZStorageManagerView`.
- **Integrasi Desktop Window Manager (`RAZWindowStore.js` & `RAZDesktopLayout.vue`):**
  - Menambahkan fungsi `openStorageManager()` di window store agar Storage Manager juga dapat dibuka sebagai jendela mandiri OS berdimensi `700x590`.
  - Menambahkan pintasan `Storage Manager` pada context menu desktop wallpaper di submenu *Alat & Sistem*.
- **Endpoint Backend REST API (`vite.config.js`):**
  - `GET /api/v1/storage/manager`: Mengembalikan data partisi storage dan media drive fisik.
  - `POST /api/v1/storage/free-up`: Membersihkan cache dan mengembalikan ruang yang dibebaskan.
  - `POST /api/v1/storage/format`: Format partisi MicroSD/USB.
  - `POST /api/v1/storage/remove`: Unmount drive penyimpanan secara aman.

#### 3. Verifikasi & Standar Mutu
- **Build Production:** `npm run build` sukses 100% tanpa error dalam 5.40 detik.
- **Backup:** Tersimpan di `backup/RAZSettingsView_before_storage_manager.vue.bak` dan `backup/vite.config_before_storage_manager.js`.
- **HMR Dev Server:** Hot reload berjalan mulus di task-5155 (`http://localhost:3000`).

---

### Versi: 3.5.1-file-manager-zero-click-init (Fix File Manager Zero-Click Initial Load & Primary Storage Auto-Listing)
**Waktu Eksekusi:** 2026-10-15 10:10 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Masalah Utama yang Diinvestigasi & Diperbaiki
- **Gejala / Laporan Pengguna:**
  - Setiap kali membuka File Manager untuk pertama kali, area konten selalu kosong melompong (*"Tidak ada berkas yang cocok"*), dan pengguna harus mengklik salah satu drive partisi di panel samping kiri secara manual baru daftar folder muncul.
- **Akar Masalah (*Root Cause*):**
  1. Variabel state `currentPath` diinisialisasi secara statis dengan `ref('/DATA')`.
  2. Pada hook `onMounted`, hanya `await fetchDevices()` yang dieksekusi.
  3. Di dalam `fetchDevices()`, pemanggilan `fetchFiles()` ditaruh di dalam blok pengkondisian:
     ```javascript
     if (!currentPath.value || currentPath.value === '/' || currentPath.value === 'C:/') {
       currentPath.value = defaultStoragePath.value
       fetchFiles(defaultStoragePath.value)
     }
     ```
  4. Karena `currentPath.value` sudah terisi nilai awal `'/DATA'`, pengkondisian tersebut bernilai **`false`**. Akibatnya, `fetchFiles()` **TIDAK PERNAH DIJALANKAN SAMA SEKALI** saat inisialisasi awal! State `files.value` tetap berupa larik kosong `[]` dan `isLoading` bernilai `false`, memicu pesan folder kosong.
  5. Pada tombol navigasi side panel *"My Storage (/DATA)"*, terdapat perbandingan referensi objek `if (currentPath !== defaultStoragePath)` yang mencegah pemanggilan ulang jika path identik, sehingga klik pada tombol My Storage tidak memicu pemuatan berkas.

#### 2. Rekayasa Teknis & Solusi
- **Perbaikan Lifecycle `onMounted` (`RAZFileManagerView.vue`):**
  - Mengubah siklus `onMounted` agar secara deterministik dan tanpa syarat memanggil `fetchFiles(targetPath)` segera setelah `fetchDevices()` selesai mendeteksi partisi sistem:
    ```javascript
    onMounted(async () => {
      await fetchDevices()
      const targetPath = defaultStoragePath.value || '/DATA'
      currentPath.value = targetPath
      addressInputValue.value = targetPath
      activeSidebarSection.value = 'drive'
      await fetchFiles(targetPath)
    })
    ```
- **Refactoring Deteksi Partisi `fetchDevices()`:**
  - Menghapus blok kondisi bercabang yang menghambat inisialisasi.
  - Memastikan `defaultStoragePath` secara cerdas mendeteksi partisi `/DATA` pada STB Linux ARM64 maupun folder `C:/DATA` pada mesin host pengembangan Windows.
- **Perbaikan Tombol Navigasi My Storage Sidebar:**
  - Memperbarui event click tombol "My Storage (/DATA)" agar langsung memanggil `fetchFiles(defaultStoragePath)` tanpa terhalang perbandingan referensi ref.
- **Sinkronisasi Address Bar & Breadcrumbs:**
  - Memastikan address bar input dan breadcrumbs interaktif langsung terisi dan sinkron dengan folder aktif sejak detik pertama jendela dibuka.

#### 3. Verifikasi & Standar Mutu
- **REST API Verification:**
  - `GET /api/v1/files/list?path=/DATA` dan `GET /api/v1/files/list?path=C:/DATA` terverifikasi merespons HTTP 200 dengan daftar 4 direktori riil (`AppData`, `apps`, `razos`, `websites`).
- **Production Build:**
  - `npm run build` sukses 100% tanpa error dalam 4.85 detik.
- **Backup:**
  - Dicadangkan di `backup/RAZFileManagerView_2026-10-15_fix_initial_load.bak`.

---

### Versi: 3.6.0-sublime-editor-and-file-manager-upload-suite (Sublime Text / Notepad++ Complex Code Editor, Multi-Mode Upload, & Instant File Creation Suite)
**Waktu Eksekusi:** 2026-10-15 11:27 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Kebutuhan Pengguna & Fitur Baru
- **Permintaan 1:**
  - Menambahkan fitur Upload File (single/multi-file).
  - Menambahkan fitur Upload Folder (recursive directory structure).
  - Menambahkan fitur Upload ZIP dengan auto-extract seketika di direktori aktif.
  - Menambahkan fitur Buat Folder.
  - Menambahkan fitur Buat File (teks/kode) yang filenya dapat langsung dibuka menggunakan Text Editor dari dalam File Manager.
- **Permintaan 2:**
  - Membuat Text Editor di dalam sistem operasi RAZOS jauh lebih kompleks, powerful, dan profesional setara **Sublime Text** atau **Notepad++**.

#### 2. Rekayasa Teknis & Arsitektur Implementasi
- **Sublime Text & Notepad++ Enterprise Complex Code Editor (`RAZTextEditorView.vue`):**
  - **Multi-Tab Workspace:** Membuka banyak berkas secara bersamaan dalam tab bar horizontal, lengkap dengan dirty/unsaved indicator, badge tipe bahasa, tombol close tab, dan tombol `+` tab baru.
  - **Menu Bar Komprehensif (Sublime Style):**
    - `Berkas Baru (Ctrl+N)`, `Simpan (Ctrl+S)`, `Simpan Semua (Ctrl+Shift+S)`, `Cari (Ctrl+F)`, `Ganti (Ctrl+H)`, `Lompat Baris (Ctrl+G)`.
  - **Collapsible Project File Tree Sidebar:** Panel samping penjelajah berkas direktori aktif untuk beralih antar file proyek secara instan tanpa perlu bolak-balik ke File Manager.
  - **Synchronized Line Numbers Gutter:** Nomor baris responsif yang tersinkronisasi 100% dengan scroll textarea dan memiliki penanda baris aktif kursor.
  - **Floating Find & Replace Bar (Ctrl+F / Ctrl+H):**
    - Pencarian teks real-time dengan counter pencocokan (contoh: `3/12`), tombol previous/next, input replace, serta tombol `Ganti` dan `Ganti Semua`.
  - **Go to Line Dialog (Ctrl+G):** Modal navigasi cepat untuk melompat langsung ke baris tertentu.
  - **Smart Indentation & Auto-Closing Pairs:** Dukungan tombol `Tab` (2 atau 4 spasi), auto-indent baris baru mengikuti baris sebelumnya, dan auto-closing brackets (`()`, `[]`, `{}`, `""`, `''`).
  - **Notepad++ Bottom Status Bar:** Informasi posisi baris & kolom kursor riil, statistik jumlah baris/kata/karakter, encoding `UTF-8`, line endings `LF/CRLF`, tombol switch indentasi `Spasi: 2/4`, dan deteksi otomatis bahasa pemrograman.
  - **Multi-Theme Selector:** Pilihan tema modern: *Sublime Monokai*, *Dark Obsidian*, *Cyberpunk Neon*, dan *Notepad++ Slate*.
- **Penyempurnaan File Manager (`RAZFileManagerView.vue`):**
  - **Dropdown Menu "+ Baru" (New Action Suite):**
    - Folder Baru (Custom Modal)
    - File Baru (Teks/Kode dengan quick template selector: `.js`, `.html`, `.css`, `.json`, `.py`, `.go`, `.md`, `.txt`, `.yaml`)
    - Upload File (multi-file native)
    - Upload Folder (recursive webkitdirectory)
    - Upload ZIP (Auto Extract)
  - **Integrasi Penuh Teks Editor:**
    - Double-click file teks/kode/config/script otomatis membuka berkas di **RAZ Studio Code**.
    - Opsi Context Menu: **"Buka di Teks Editor"** (ikon `FileCode`) tersedia di klik kanan berkas.
    - Tombol aksi cepat Teks Editor (`Code2`) tersedia di baris list view dan status bar item terpilih.
    - Membuat berkas baru memiliki opsi **"Buat & Buka Editor"** yang langsung memunculkan window RAZ Studio Code.
  - **Modal Progress Upload:** Menampilkan nama file aktif, progress bar animasi (0-100%), dan status counter penyelesaian.
- **Backend & Middleware REST API (`vite.config.js`):**
  - `POST /api/v1/files/create-file`: Membuat berkas baru di direktori parent dan mengembalikan file item descriptor.
  - `POST /api/v1/files/upload-raw`: Menerima raw binary stream upload berkas, pembuatan subfolder otomatis via header `x-relative-path`, dan auto-extract arsip ZIP secara instan via utilitas `tar` / `bsdtar` sistem.

#### 3. Verifikasi & Standar Mutu
- **Production Build:** `npm run build` sukses 100% tanpa error dalam 5.25 detik (`built in 5.25s`).
- **REST API Endpoints:** `/api/v1/files/create-file` dan `/api/v1/files/upload-raw` terverifikasi sukses via pengujian integrasi.
- **Backup:** Dicadangkan di `backup/RAZFileManagerView_before_upload_and_editor_suite.vue.bak`, `backup/RAZDesktopLayout_before_text_editor.vue.bak`, dan `backup/RAZWindowStore_before_text_editor.js.bak`.

---

### Versi: 3.6.1-readme-overhaul-and-ribbon-actions (Professional Readme Architecture & Visual Gallery Overhaul, Mobile Ribbon New Menu)
**Waktu Eksekusi:** 2026-10-15 11:42 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Perubahan & Fitur Baru
- **Overhaul Total Readme.md Sesuai Standar Enterprise (Referensi `rahmanef63/mso`):**
  - **Header Modern & Presisi:** Logo brand capsule terpusat, judul, deskripsi filosofi RAZOS, quick navigation anchor bar, dan 6 badge shields.io resmi.
  - **Visual Showcase Gallery (Tabel 2-Kolom):**
    - `RAZOS-Desktop-Workspace.png`: Desktop multi-window responsif & Windows 11 Snap Gesture.
    - `RAZOS-File-Manager.png`: Hybrid file manager, upload file/folder/ZIP auto-extract & code editor.
    - `RAZOS-App-Install-Docker.png`: Katalog 525+ Docker apps & compose suite.
    - `RAZOS-Task-Manager.png`: Task manager hardware riil, heatmap & streaming telemetry.
    - `RAZOS-Storage-Manager.png`: Storage manager partisi disk eMMC & MicroSD `/DATA`.
  - **Tabel Nilai Unggulan ("Mengapa RAZOS?"):** Perbandingan komparatif kebutuhan server berdaya rendah vs solusi teknis yang dihadirkan RAZOS.
  - **Diagram Arsitektur Mermaid Lengkap:** Pemetaan alur Client (Desktop/Mobile/PWA) -> Frontend Vue 3 -> Go 1.22+ Single Binary (Router, WS, Embed FS, Security) -> Linux Host OS (Docker, Systemd, Disk Mounts, Amlogic S905X STB).
  - **Spesifikasi Target Hardware & Jejak Sumber Daya:** Batasan RAM 1.8 GiB, ZRAM Swap ~900 MiB, idle memori Go <= 35 MB, binary single-file <= 25 MB, dan eMMC 8 GB vs MicroSD `/DATA`.
  - **Panduan Instalasi & Deployment 3 Jalur:** One-line curl install untuk STB, local development setup (Go & Vite), serta kompilasi produksi single binary ARM64.
  - **Pohon Direktori Proyek (Project Tree):** Struktur direktori mutakhir dan roadmap transparan.
- **Ribbon Action Bar "+ Baru" di File Manager (`RAZFileManagerView.vue`):**
  - Menambahkan tombol dropdown "+ Baru" di ribbon bar atas File Manager di samping tab filter sehingga selalu terlihat pada mode mobile dan jendela sempit (di mana sidebar kiri di-hide dengan `hidden md:flex`).
  - Menghadirkan handler penutupan popover global saat pengguna mengklik area luar menu.
- **Aset Media Dokumentasi (`docs/media/`):**
  - Mengamankan dan menata seluruh screenshot riil beresolusi tinggi ke dalam folder `docs/media/`.

#### 2. Verifikasi & Standar Mutu
- **Production Build:** `npm run build` sukses 100% tanpa error dalam 4.77 detik (`built in 4.77s`).
- **Backup:** Dicadangkan di `backup/Readme_before_mso_style.md.bak` dan `backup/RAZFileManagerView.vue.bak`.
- **Git Synchronization:** Sinkronisasi commit dan push ke repositori GitHub `origin/master` dan `origin/main`.

---

### Versi: 3.6.2-file-editor-integration-and-data-binding (Penyatuan Teks Editor ke File Manager, Perbaikan Data Binding Konten Berkas, & Penghapusan Dialog Native Browser)
**Waktu Eksekusi:** 2026-10-15 14:22 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer

#### 1. Masalah Utama yang Diperbaiki
- **Teks Editor Sebagai Tools Terpisah:** Sebelumnya teks editor terdaftar sebagai aplikasi desktop independen (`RAZ Studio Code`) di desktop apps list, sehingga tampak seperti tools asing yang berdiri sendiri.
- **Berkas Teks Kosong saat Dibuka (Data Binding Bug):** Saat berkas teks dibuka dari File Manager, `openCustomWindow` di `RAZWindowStore.js` tidak memetakan atribut `fileData` ke objek window, dan `RAZTextEditorView.vue` hanya memeriksa `props.initialFile` sekali saat `onMounted` tanpa watcher. Akibatnya, editor selalu membuka tab kosong `Untitled-1.txt`.
- **Penggunaan Dialog Native Browser (`confirm` & `prompt`):** Teks editor sebelumnya menggunakan `confirm()` dan `prompt()` bawaan browser yang melanggar aturan enterprise `GEMINI.md`.

#### 2. Rekayasa Teknis & Solusi Implementasi
- **Integrasi Penuh Teks Editor Sebagai Modul Bawaan File Manager:**
  - Menghapus `text-editor` dari daftar `apps` di `RAZWindowStore.js` dan menghapus slot posisinya di `backend/internal/store/RAZDesktopConfig.json`. Ikon tidak lagi muncul sebagai aplikasi mandiri di desktop atau app drawer.
  - Memperbarui window ID menjadi `'file-manager-editor'` dengan title format: `${fileItem.name} - Editor Berkas (File Manager)` dan ikon resmi File Manager (`/assets/icons/RAZIconFileManager.svg`).
  - Mengubah label brand di menu bar teks editor menjadi **"File Manager / Editor Berkas"** untuk menegaskan statusnya sebagai satu kesatuan fitur File Manager.
- **Perbaikan Data Binding & Lifecycle Reactivity (`RAZTextEditorView.vue`):**
  - Menambahkan dukungan `win` prop pada komponen editor.
  - Menambahkan watcher reaktif mendalam (`deep: true`) pada `props.initialFile`, `props.win?.fileData`, dan `props.win?.customData`.
  - Jika editor sudah terbuka dan pengguna mendobel-klik file teks lain di File Manager, editor secara otomatis membuka berkas tersebut di tab baru dan memuat kontennya tanpa jeda.
  - Menambahkan auto-replacement: Jika tab yang ada hanya 1 dan merupakan tab kosong bawaan (untitled belum diedit), berkas yang baru dibuka langsung menggantikan tab kosong tersebut secara mulus.
- **Zero Native Dialogs (Custom Glassmorphism Modals):**
  - Mengganti `prompt()` penyimpanan berkas baru dengan modal kustom **Save As** yang elegan lengkap dengan lokasi folder aktif, input nama berkas, dan tombol Simpan.
  - Mengganti `confirm()` saat menutup tab yang belum disimpan dengan modal kustom **Perubahan Belum Disimpan** (`showCloseTabConfirmModal`).
  - Mengganti `confirm()` pada aksi Tutup Semua Tab dengan modal kustom **Tutup Semua Tab** (`showCloseAllConfirmModal`).

#### 3. Verifikasi & Standar Mutu
- **Production Build:** `npm run build` sukses 100% tanpa error dalam 7.85 detik (`built in 7.85s`).
- **Backup:** Dicadangkan di `backup/RAZWindowStore_before_file_editor_integration.js.bak`, `backup/RAZTextEditorView_before_file_editor_integration.vue.bak`, `backup/RAZDesktopLayout_before_file_editor_integration.vue.bak`, dan `backup/RAZDesktopConfig_before_file_editor_integration.json.bak`.
- **Git Synchronization:** Siap untuk commit dan push ke GitHub repository.



