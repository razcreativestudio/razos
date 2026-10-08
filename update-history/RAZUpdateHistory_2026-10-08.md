# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.6.0-file-manager-suite (Penambahan Fitur Advanced Storage & File Manager Ala CasaOS)
**Waktu Eksekusi:** 2026-10-08 23:37 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Backend Storage & File Engine
- Menambahkan modul `backend/internal/filemanager/RAZStorageEngine.go`:
  - **Deteksi Hardware Perangkat:** Membaca status partisi eMMC sistem (`/`), MicroSD penyimpanan utama (`/DATA`), dan drive HDD eksternal / USB Flashdisk (`/media/*` atau `/dev/sd*`) menggunakan `gopsutil/disk`.
  - **Operasi Direktori & Berkas:** List direktori aman dengan ekstensi dan ukuran berkas, pembuatan folder baru, hapus rekursif, ubah nama (rename), serta pembacaan dan penyimpanan konten berkas teks (editor).
- Mendaftarkan rute REST API di `backend/internal/api/RAZApiController.go`:
  - `GET /api/v1/storage/devices`: Status kapasitas disk, mountpoint, filesystem, dan deteksi USB.
  - `GET /api/v1/files/list`: Pembacaan isi folder berdasarkan parameter query `path`.
  - `POST /api/v1/files/folder`: Pembuatan folder baru.
  - `DELETE /api/v1/files/item`: Penghapusan berkas atau direktori.
  - `POST /api/v1/files/rename`: Mengubah nama berkas/folder.
  - `GET /api/v1/files/content` & `POST /api/v1/files/content`: Editor berkas teks langsung di browser.

#### 2. Frontend UI Upgrade (CasaOS-Style Files & Storage Explorer)
- Menambahkan icon vector resmi `frontend/public/assets/icons/RAZIconFileManager.svg`.
- Mendaftarkan aplikasi `file-manager` ("Files & Storage") pada `RAZWindowStore.js`, `RAZDesktopLayout.vue`, dan `RAZMobileLayout.vue`.
- Mengembangkan antarmuka kaya fitur `frontend/src/views/RAZFileManagerView.vue`:
  - **Storage Devices Carousel (CasaOS Look):** Kartu progress bar visual penggunaan ruang partisi eMMC, MicroSD `/DATA`, dan USB external drive.
  - **Hierarchical Breadcrumbs Navigation:** Navigasi penjelajahan direktori tingkat lanjut dengan tombol kembali (*up directory*).
  - **File Type Badges & Icons:** Pembedaan visual otomatis untuk folder, source code (`.go`, `.js`, `.json`, `.vue`), dokumen teks, dan file biner.
  - **In-Browser Text Editor Modal:** Modal gelap untuk membaca, mengedit, dan menyimpan file teks (misal: config, env, json, md, sh) secara langsung tanpa perlu SSH.
  - **Custom Modal Dialogs:** Modal pembuatan folder baru dan pengubahan nama berkas (*rename*).
- Build production terverifikasi sukses dengan `npm run build`.
