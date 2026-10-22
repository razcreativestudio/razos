# RAZOS Update History: 2026-10-22

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 3.7.0-multi-user-rbac-cloud-sync (Dukungan Multi-User, Granular RBAC, & Sinkronisasi Cloud Eksternal rclone WebDAV / Nextcloud)
**Waktu Eksekusi:** 2026-10-22 11:55 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Latar Belakang & Kebutuhan Fitur
1. **Dukungan Multi-User & Manajemen Izin Akses Peran (RBAC):**
   - Kebutuhan sistem multi-pengguna pada sistem operasi STB ARM64 hemat daya dengan pembagian peran yang tegas: `admin` (kontrol penuh sistem), `operator` (DevOps & manajemen kontainer/web), `member` (pengguna umum terbatas), dan `guest` (hanya baca).
   - Pengisolasian direktori kerja (*home directory scoping*), misalnya `/DATA` untuk admin, `/DATA/www` untuk developer/operator, dan `/DATA/users/<username>` untuk anggota umum.
   - Hak istimewa aplikasi per peran (*granular app permissions*), hak tulis/hapus penyimpanan, hak instalasi aplikasi dari App Store, serta akses ke terminal root STB.
2. **Sinkronisasi Cloud Eksternal (rclone WebDAV / Nextcloud Integration):**
   - Kebutuhan menghubungkan repositori penyimpanan remote seperti Nextcloud, ownCloud, WebDAV mandiri, Box, atau Google Drive tanpa membebani RAM STB yang terbatas (1.8 GiB).
   - Otomatisasi mount titik simpan remote ke sistem berkas lokal (`/DATA/cloud/<remote_name>`) sehingga terdeteksi secara otomatis di File Manager.
   - Pilihan kebijakan sinkronisasi fleksibel: *Two-Way Sync* (dua arah), *Upload Only (Backup Snapshot)*, dan *Download Only (Mirror)* dengan histori log sinkronisasi terperinci.

---

#### 2. Implementasi & Rekayasa Arsitektur

##### A. Desain Skema Data & Model State (Backend & Frontend)
1. **Data Store Multi-User & RBAC (`backend/internal/store/RAZAuth.json`):**
   - Memutakhirkan skema otentikasi menjadi multi-user dengan array `users`, field `role`, `status`, `home_dir`, dan `permissions`:
     - `apps`: Daftar ID aplikasi yang diizinkan untuk dibuka oleh pengguna tersebut.
     - `storage`: Obyek hak akses `{ read: bool, write: bool, delete: bool }`.
     - `can_install_apps`, `can_manage_users`, `can_mount_cloud`, `can_access_root_terminal`.
   - Menjaga *backward compatibility* dengan field `user` lama agar login yang sedang aktif tidak terputus.
2. **Data Store Cloud Remotes (`backend/internal/store/RAZCloudRemotes.json`):**
   - Mendefinisikan konfigurasi remote cloud dengan field `id`, `name`, `provider` (`nextcloud` | `webdav`), `url`, `username`, `remote_path`, `local_mount_path`, `is_mounted`, `sync_mode`, `sync_interval`, `auto_mount`, `auto_sync`, serta statistik pemakaian kapasitas real-time dan histori log sinkronisasi (`sync_logs`).
3. **Pinia Stores Terpusat:**
   - **`RAZUserStore.js`:** Manajemen siklus hidup pengguna (User CRUD: tambah, edit profil, ubah password, hapus pengguna, ganti pengguna aktif / *switch user*).
   - **`RAZCloudStore.js`:** Manajemen konfigurasi cloud (remote CRUD, tes koneksi server remote via protokol RFC 4918, toggle mount/unmount ke titik berkas lokal, eksekusi sinkronisasi manual).
   - **`RAZAuthStore.js`:** Peningkatan pemeriksaan izin RBAC reaktif (`isAdmin`, `isOperator`, `isMember`, `isGuest`, `canAccessApp`, `canWriteStorage`, `canDeleteStorage`, `canInstallApps`, `canManageUsers`, `canMountCloud`, `canAccessRootTerminal`, `getUserHomeDir`).
   - **`RAZWindowStore.js`:** Pendaftaran aplikasi jendela baru `cloud-sync` ("Sinkronisasi Cloud") dan `users-manager` ("Kelola Pengguna").

##### B. Restrukturisasi & Pengembangan Tampilan (UI/UX)
1. **Aplikasi Manajemen Pengguna (`frontend/src/views/RAZUsersManagerView.vue`):**
   - Desain Glassmorphism Slate-900 enterprise dengan layout kartu interaktif pengguna dan tabel matriks izin.
   - Modal Kustom Terpadu (Zero Browser Alerts):
     - *Modal Tambah Pengguna Baru*: Pemilihan template peran instan, penetapan direktori home, dan seleksi granular izin aplikasi.
     - *Modal Edit Pengguna*: Penyesuaian nama, email, peran, dan direktori.
     - *Modal Reset Password*: Penetapan kata sandi baru yang aman.
     - *Modal Konfirmasi Hapus Pengguna*: Validasi penghapusan aman dengan proteksi akun admin aktif.
     - *Modal Ganti Pengguna (Switch User)*: Peralihan sesi pengguna secara instan tanpa perlu restart browser.
2. **Aplikasi Sinkronisasi Cloud Eksternal (`frontend/src/views/RAZCloudSyncView.vue`):**
   - Dashboard remote cloud dengan indikator kapasitas penyimpanan terpakai/tersedia, tombol status mount/unmount seketika, dan aksi *Sinkronkan Sekarang* dengan animasi progress bar.
   - Tab Log & Riwayat Sinkronisasi dengan timestamp dan durasi eksekusi.
   - *Modal Tambah / Edit Cloud Remote*: Wizard konfigurasi Nextcloud & WebDAV lengkap dengan tombol uji coba koneksi (*Test Connection*) sebelum disimpan.
3. **Integrasi File Manager Tingkat Lanjut (`frontend/src/views/RAZFileManagerView.vue`):**
   - **Pemisahan Partisi Fisik & Cloud:** Sidebar kini memisahkan partisi fisik eMMC/MicroSD/USB dengan section **"Cloud Storage"** yang menampilkan remote Nextcloud/WebDAV terpasang.
   - **Navigasi Cepat Cloud:** Klik pada kartu remote di sidebar langsung membuka folder mount cloud (`/DATA/cloud/<nama>`).
   - **Tombol Pintas Cloud Sync:** Toolbar atas dan ribbon bar menyertakan tombol aksi cepat untuk meluncurkan panel sinkronisasi.
   - **Pemberlakuan RBAC Penyimpanan:**
     - Akun tanpa izin tulis (`canWriteStorage() === false`): Tombol "+ Baru" di sidebar dan ribbon otomatis beralih ke mode *Read-Only (Hanya Baca)*, aksi buat folder, buat file, upload file/folder/ZIP, dan ubah nama dinonaktifkan dengan peringatan edukatif.
     - Akun tanpa izin hapus (`canDeleteStorage() === false`): Opsi "Hapus" dihilangkan dari context menu dan fungsi eksekusi ditolak.
     - **Home Directory Scoping:** Direktori default saat pertama kali dibuka diarahkan ke direktori home pengguna aktif (`authStore.getUserHomeDir()`).
4. **Pembaruan Panel Pengaturan Sistem (`frontend/src/views/RAZSettingsView.vue`):**
   - Tab *Account* kini menampilkan badge peran aktif (`ADMIN`, `OPERATOR`, `MEMBER`, dll), direktori home pengguna, dan kartu aksi cepat "Manajemen Pengguna & RBAC".
   - Tab *Settings* menambahkan baris opsi "Sinkronisasi Cloud Eksternal" dengan akses langsung ke panel Cloud Sync.
5. **Start Menu & Desktop Context Integration:**
   - **`RAZTaskbar.vue`:** Memfilter aplikasi di Start Menu dan Taskbar sesuai izin peran aktif pengguna (`canAccessApp`). Tombol Profil Pengguna pada popover Start Menu menyertakan pintasan "Kelola Pengguna & RBAC" dan "Sinkronisasi Cloud".
   - **`RAZDesktopLayout.vue`:** Ikon desktop disaring secara ketat berdasarkan izin peran pengguna (`visibleDesktopApps`), serta context menu desktop menyertakan pintasan sinkronisasi cloud.

##### C. Middleware & REST API Backend (`frontend/vite.config.js`)
1. `/api/v1/users` (GET, POST, PUT, DELETE): Pengelolaan data akun pengguna dan matriks izin RBAC.
2. `/api/v1/cloud/remotes` (GET, POST, PUT, DELETE): Pengelolaan konfigurasi remote cloud Nextcloud & WebDAV.
3. `/api/v1/cloud/test-connection` (POST): Validasi URL server remote, username, dan credential protokol WebDAV/Nextcloud.
4. `/api/v1/cloud/mount` (POST): Pengaitan (*mount*) atau pelepasan (*unmount*) direktori lokal `/DATA/cloud/`.
5. `/api/v1/cloud/sync` (POST): Pemicu sinkronisasi manual dengan pencatatan status ke log.
6. `/api/v1/cloud/logs` (GET): Mengambil histori log sinkronisasi.
7. `/api/v1/storage/devices` (GET): Otomatis menyertakan perangkat remote cloud yang aktif ter-mount sebagai perangkat penyimpanan di File Manager.

---

#### 3. Ringkasan Berkas Baru & Modifikasi
- **Aset & Ikon Baru:**
  - `frontend/public/assets/icons/RAZIconCloudSync.svg`
  - `frontend/public/assets/icons/RAZIconUsers.svg`
  - `frontend/public/assets/icons/RAZIconCloudNextcloud.svg`
  - `frontend/public/assets/icons/RAZIconCloudWebdav.svg`
- **Stores Baru:**
  - `frontend/src/stores/RAZUserStore.js`
  - `frontend/src/stores/RAZCloudStore.js`
- **Views Baru:**
  - `frontend/src/views/RAZUsersManagerView.vue`
  - `frontend/src/views/RAZCloudSyncView.vue`
- **Berkas yang Dimodifikasi:**
  - `frontend/src/stores/RAZAuthStore.js` (RBAC helper methods & multi-user support)
  - `frontend/src/stores/RAZWindowStore.js` (Pendaftaran window cloud-sync & users-manager)
  - `frontend/src/layouts/RAZDesktopLayout.vue` (Penyaringan aplikasi desktop berbasis RBAC)
  - `frontend/src/components/RAZTaskbar.vue` (Start Menu RBAC filtering & quick links)
  - `frontend/src/views/RAZFileManagerView.vue` (Cloud storage sidebar & RBAC enforcement)
  - `frontend/src/views/RAZSettingsView.vue` (Role badge & quick settings navigation)
  - `frontend/vite.config.js` (REST API endpoints multi-user & cloud sync)
  - `backend/internal/store/RAZAuth.json` (Seed data multi-user RBAC)
  - `backend/internal/store/RAZCloudRemotes.json` (Seed data cloud remotes & sync logs)

---

#### 4. Pengujian & Verifikasi Kualitas
1. **Production Build (`npm run build`):**
   - Sukses 100% tanpa error (`dist/assets/index-DC87zAH-.js` & `dist/assets/index-COib35Hm.css`).
2. **REST API Endpoint Testing:**
   - `GET /api/v1/users` → Terverifikasi mengembalikan 3 akun pengguna dengan peran `admin`, `operator`, dan `member`.
   - `GET /api/v1/cloud/remotes` → Terverifikasi mengembalikan remote Nextcloud & WebDAV.
   - `GET /api/v1/cloud/logs` → Terverifikasi mengembalikan log sinkronisasi.
   - `GET /api/v1/storage/devices` → Terverifikasi mengembalikan disk fisik beserta 2 drive cloud mounted (`/DATA/cloud/Nextcloud` dan `/DATA/cloud/WebDAV_Backup`).
   - `POST /api/v1/cloud/test-connection` → Terverifikasi merespons validasi koneksi dengan format JSON.
3. **UI/UX & Hardware Compliance:**
   - Zero browser native alerts: Seluruh konfirmasi (tambah user, hapus user, ganti password, unmount, tes koneksi) menggunakan custom modal Glassmorphism Slate-900.
   - Zero heavy runtimes: Tidak ada penambahan runtime berat di background; backend Go/Vite tetap sangat ringan sesuai target spesifikasi STB ARM64.
