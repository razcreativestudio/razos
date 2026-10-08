# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.7.0-cloud-drive-suite (Evolusi File Manager Kompleks ala Google Drive & CasaOS Cloud)
**Waktu Eksekusi:** 2026-10-08 23:41 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Upgrade Arsitektur Antarmuka (Google Drive Look & Feel)
- Memperbarui komponen `frontend/src/views/RAZFileManagerView.vue` menjadi aplikasi Cloud Drive penuh:
  - **Left Navigation Sidebar:** Tombol "Folder Baru", menu navigasi *My Drive (/DATA)*, *Berbintang (Starred)*, deteksi live *Storage Disks / HDD*, serta pengukur kapasitas cloud storage (*Storage Meter*).
  - **Top Search & Category Filters:** Pencarian instan berkas (*live filter search*) dengan tombol penyaring format (*Semua, Folder, Dokumen, Kode, Gambar*).
  - **Dual Layout Viewing:** Toggle 1-klik antara **Grid View** (kartu thumbnail dengan tombol bintang favorit) dan **List View** (tabel detail multi-kolom).
  - **Right Details Preview Pane:** Panel samping kanan dinamis yang menampilkan metadata file terpilih (tipe, ukuran file dalam bytes/KB/MB, lokasi path fisik, waktu modifikasi, serta tombol aksi cepat).
  - **Starred / Favorit Toggle:** Pengguna dapat menandai file/folder dengan bintang untuk diakses cepat.
  - **In-Browser Text & Code Editor:** Editor modal layar penuh untuk membaca dan mengedit berkas teks/konfigurasi secara langsung.
- Build terverifikasi sukses dengan `npm run build` tanpa error.

---

### Versi: 1.8.0-file-manager-media-suite (Built-in Media Viewers & Player Suite)
**Waktu Eksekusi:** 2026-10-08 23:45 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Integrasi Built-in Media Viewers & Streaming Raw File
- **Image Viewer:** Modal interaktif untuk preview gambar dengan fungsi Zoom In, Zoom Out, Reset, dan download langsung.
- **Video Player:** Pemutar video web bawaan dengan kontrol pemutaran penuh untuk format MP4, WebM, MKV.
- **Music / Audio Player:** Pemutar audio terintegrasi dengan cover visual gelombang audio untuk format MP3, WAV, OGG, FLAC.
- **PDF Viewer:** Pembaca dokumen PDF langsung via embedded viewer interaktif.
- **Code & Text Editor:** Modal in-browser editor lengkap untuk berkas teks, konfigurasi, json, yaml, script.
- **Binary Streaming Backend Route:** Mendaftarkan endpoint `/api/v1/files/raw` di `backend/internal/api/RAZApiController.go` dan dev server Vite `frontend/vite.config.js` untuk melayani streaming berkas media secara langsung tanpa buffer lambat.
- **Clean Build:** Terverifikasi `npm run build` sukses 100%.

---

### Versi: 1.9.0-real-disk-hardware-fix (Penghapusan Total Disk Mocking & Deteksi Disk Fisik Riil)
**Waktu Eksekusi:** 2026-10-08 23:54 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Real Hardware Storage Inspection (Zero Dummy Data)
- **Eliminasi Mock Hardcoded Devices:** Mengganti seluruh data hardcoded storage di backend dan dev middleware dengan inspeksi disk fisik aktual sistem.
- **Deteksi Partisi Windows & Linux:** Mengintegrasikan `wmic logicaldisk` untuk Windows dan `df -k -P` untuk Linux ARM64 STB. Device, ukuran total byte, ruang terpakai, ruang bebas, persentase, dan nama volume disk (misal `C: Local Disk`, `Primary Storage`) dibaca secara langsung dari host fisik.
- **Dynamic Cloud Storage Meter:** Pengukur penyimpanan di sidebar kiri kini menghitung total kumulatif dari seluruh partisi fisik yang aktif, menggantikan persentase hardcode sebelumnya.

#### 2. Perbaikan Navigasi & File System Explorer
- **Breadcrumbs & Up Navigation Fix:** Mengatasi pemotongan path Windows (`C:/`, `C:\`) dan Linux (`/DATA`, `/`) sehingga navigasi hierarki direktori dan naik satu tingkat berfungsi secara akurat tanpa error string indexing.
- **Initial Mountpoint Auto-Discovery:** Menyesuaikan inisialisasi direktori awal agar secara dinamis membuka partisi fisik pertama yang tersedia jika path default belum terbentuk.
- **REST API Endpoints Sinkron:** Menambahkan endpoint `/api/v1/files/folder`, `/api/v1/files/item` (delete), dan `/api/v1/files/rename` yang sinkron antara Go backend dan Vite dev server.

---

### Versi: 2.0.0-app-store-mega-suite (520+ Docker Ready-to-Use, README.md Viewer & CasaOS Manual Install)
**Waktu Eksekusi:** 2026-10-09 00:04 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Katalog Mega 520+ Kontainer Docker Siap Pakai
- **Ekspansi Katalog:** Mengintegrasikan 525 kontainer Docker siap pakai terkurasi yang mencakup 10 kategori industri: Media & Streaming, File & Cloud Storage, Database & Caching, Development & Git, Home Automation & IoT, Networking & Security, Monitoring & Analytics, Productivity & Office, Downloaders & Torrents, dan AI & Machine Learning.
- **README.md & Details Modal:** Menyediakan tombol **Detail** pada setiap kartu aplikasi untuk menampilkan dokumentasi lengkap, panduan deployment, spesifikasi image, port default, dan mountpoint penyimpanan.
- **Pencarian Secepat Kilat (Live Search):** Filter pencarian responsif untuk mencari aplikasi berdasarkan nama, deskripsi, image, maupun kategori.

#### 2. Tombol & Modal "Instal Manual Docker" (Ala CasaOS)
- **Custom Docker Deployer:** Menambahkan tombol **Instal Manual Docker** dengan form modal kustom yang memungkinkan pengguna memasukkan sembarang Docker image dari Docker Hub/Registry publik, menentukan host port, mengatur path direktori volume host, dan environment variables kustom.
- **Backend Custom Deployment Engine:** Memperbarui `RAZAppStoreEngine.go` untuk menangani parameter `custom_image`, `container_name`, port, volume, dan limitasi RAM protektif 256 MB.

---

### Versi: 2.1.0-app-store-visual-page (Official Docker Logos, Full Page Details View & Visual Markdown Preview)
**Waktu Eksekusi:** 2026-10-09 00:11 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Full Page App Details View (Bukan Modal Dialog)
- **Halaman Penuh Detail Aplikasi:** Mengubah alur interaksi sehingga mengklik kartu aplikasi di katalog langsung membuka Halaman Penuh (*Full Page App View*) lengkap dengan header navigasi "Kembali", thumbnail banner resolusi tinggi, ikhtisar fungsi, dan sidebar spesifikasi hardware.
- **Action Buttons Terintegrasi:** Tombol **Pasang Aplikasi**, **Status Container Aktif**, dan **Uninstall** tersedia langsung di dalam Halaman Detail maupun di kartu katalog.

#### 2. Logo Resmi Masing-Masing Docker & Banner Thumbnail
- **Iconic Brand Logos:** Menambahkan logo resmi SVG/PNG untuk aplikasi populer (*Plex, Jellyfin, Nextcloud, MariaDB, PostgreSQL, Redis, VS Code, Ollama, Home Assistant, Pi-hole, Uptime Kuma, qBittorrent, dll.*) serta identicon dinamis untuk seluruh microservice.
- **Banner Pratinjau Gambar:** Menyertakan banner visual fotografi dan ilustrasi dashboard yang relevan di setiap kategori aplikasi.

#### 3. Visual Markdown README.md Previewer (Bukan Text Editor)
- **Render Engine `marked`:** Mengintegrasikan library `marked` untuk menyajikan dokumentasi `README.md` dalam bentuk preview Markdown interaktif yang bersih dan terformat rapi (judul, gambar banner pratinjau, blok kode sintaks, daftar peluru, kutipan).




