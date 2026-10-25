# RAZOS Update History: 2026-10-25

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 3.8.0-smart-installer-self-healing (Self-Healing Binary Download & Automatic Host Build Fallback)
**Waktu Eksekusi:** 2026-10-25 12:02 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pengguna melaporkan error saat menjalankan skrip 1-Click Installer di STB Armbian:
  `curl: (22) The requested URL returned error: 404`
- Hal ini terjadi karena URL mengarah ke nama organisasi lama dan asset rilis binary ARM64 belum dipublikasikan di menu GitHub Releases, menyebabkan file executable kosong dan service `razos.service` gagal berjalan.

#### 2. Solusi & Rekayasa Skrip (`scripts/install.sh`)
- **Pembaruan URL Repository:** Menyesuaikan endpoint resmi ke `https://github.com/razcreativestudio/razos`.
- **Validasi Integritas Binary:** Menambahkan pengecekan ukuran berkas (> 100 KB) untuk memastikan berkas yang terunduh bukan berkas kosong atau halaman error 404.
- **Mode Fallback Otomatis (Self-Healing):** Jika asset rilis belum tersedia di GitHub Releases, installer secara otomatis menginstal dependensi ringan (`git`, `golang-go`), meng-clone repository, dan melakukan kompilasi Go native langsung di STB ARM64.
- **Verifikasi Status Service:** Memastikan `systemctl is-active razos` sebelum menampilkan status instalasi sukses.


---

### Versi: 3.8.1-go-module-tidy-resolution (Go Module Resolution & Direct Native Compile)
**Waktu Eksekusi:** 2026-10-25 12:18 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pada saat kompilasi mandiri di STB (`root@raz-server`), Go compiler mengeluarkan galat `missing go.sum entry for module providing package...` dan binary `/DATA/razos/bin/razos` tidak terbentuk sehingga systemd gagal menjalankan service (exit-code 203/EXEC).
- Penyebab: Perintah kompilasi dijalankan langsung di subfolder `backend/cmd/server/` tanpa terlebih dahulu mengunduh checksum dependensi (`go mod tidy`) di direktori root modul (`backend/`). Selain itu, versi Go pada `backend/go.mod` diturunkan ke `go 1.19` agar kompatibel dengan paket default Debian Bookworm (`golang-1.19-go`).

#### 2. Solusi & Perubahan
- **Sinkronisasi Modul:** Menjalankan `go mod tidy` pada root modul `/DATA/razos/src/backend` sebelum kompilasi.
- **Kompilasi Standalone:** Menjalankan `go build -ldflags="-s -w" -o /DATA/razos/bin/razos ./cmd/server/RAZMainServer.go` dari root backend.
- **Pembaruan `scripts/install.sh`:** Memperbarui alur fallback agar otomatis mengeksekusi `go mod tidy` di root modul backend sebelum memanggil `go build`.
- **Kompatibilitas Go:** Menyesuaikan `backend/go.mod` ke `go 1.19`.


---

### Versi: 3.8.3-opentelemetry-dependency-pinning (OpenTelemetry & Docker SDK Compatibility Pinning)
**Waktu Eksekusi:** 2026-10-25 12:38 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pada saat kompilasi dengan Go 1.22 di STB ARM64, dependensi transitif Docker SDK (`otelhttp`) secara otomatis mencoba menarik versi rilis bleeding-edge `v0.72.0` yang menuntut `go >= 1.26.0`, memicu error `toolchain upgrade needed` dan kegagalan kompilasi.

#### 2. Solusi & Perubahan
- Mengunci (*pinning*) versi OpenTelemetry (`otelhttp v0.49.0`, `otel v1.24.0`), Docker utilities (`go-connections`, `go-units`, `image-spec`), serta submodul `golang.org/x/net v0.22.0` dan `golang.org/x/sys v0.18.0` di dalam berkas `backend/go.mod`.
- Memastikan seluruh rantai dependensi 100% cocok dan terverifikasi di atas Go 1.22 tanpa memerlukan upgrade toolchain tambahan.


---

### Versi: 3.8.4-go-syntax-and-docker-sdk-cleanup (Go Syntax & Docker SDK Cleanup)
**Waktu Eksekusi:** 2026-10-25 13:26 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pada saat kompilasi Go native di STB ARM64, Go compiler mendeteksi beberapa galat kode internal:
  1. `RAZWebsiteManager.go`: Karakter backtick ter-escape (`\`\``) di dalam raw string literal Go merusak parsing compiler.
  2. `RAZAppStoreEngine.go` & `RAZDockerManager.go`: Pemanggilan `cli.ContainerStop` dan `cli.ContainerRestart` menggunakan `nil` alih-alih `container.StopOptions{}` pada Docker SDK v25.
  3. `RAZStorageEngine.go` & `RAZDatabaseSuite.go` & `RAZDockerManager.go`: Import tidak terpakai (`io`, `fmt`, `os/exec`).
  4. Modul `otlptracehttp` mencoba menarik versi `v1.47.0` yang meminta `go >= 1.26.0`.

#### 2. Solusi & Perubahan
- Merefaktor template string Node.js pada `RAZWebsiteManager.go` menjadi string concatenation JavaScript standar tanpa karakter backtick.
- Mengganti seluruh parameter penghentian/restart kontainer Docker menjadi `container.StopOptions{}`.
- Membersihkan seluruh unused import dan mengunci `otlptracehttp v1.24.0` dan `sdk v1.24.0` pada `backend/go.mod`.


---

### Versi: 3.8.5-dual-port-listening-and-route-completion (Port 80 Direct Access & Full API Routes)
**Waktu Eksekusi:** 2026-10-25 14:19 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pengguna ingin mengakses dashboard RAZOS langsung via IP tanpa harus mengetik port `:8080` (menggunakan default HTTP port 80).
- Frontend mengalami 404 pada `/api/v1/store/installed` dan `/api/v1/cloud/remotes`.
- Terjadi galat 500 pada pembacaan direktori `C:/DATA` di Linux ARM64 dan pembacaan katalog `RAZApps.json`.
- Terdapat potensi `TypeError: Cannot read properties of undefined (reading 'toFixed')` saat metrik hardware/storage sedang dimuat.

#### 2. Solusi & Perubahan
- **Dual Port Listening (`RAZMainServer.go`):** Server Go sekarang secara bersamaan mendengarkan (*listen*) di port utama `80` (akses langsung IP `http://192.168.100.117`) dan port sekunder `8080` (kompatibilitas backward).
- **Normalisasi Path Multi-OS (`RAZStorageEngine.go`):** Permintaan path diawali `C:/DATA` atau `C:` secara otomatis dinormalkan ke `/DATA` dan `/` pada sistem operasi Linux.
- **Penyelesaian Rute API (`RAZApiController.go`):** Menambahkan handler lengkap untuk `/api/v1/store/installed` serta modul Cloud Storage Suite (`/api/v1/cloud/remotes`, `/logs`, `/test-connection`, `/sync`).
- **Penyematan Katalog Mandiri (`RAZAppStoreEngine.go`):** Menggunakan `//go:embed RAZApps.json` sehingga data katalog selalu tersedia secara atomik tanpa ketergantungan path disk.
- **Null-Safety Guard Frontend:** Mengamankan seluruh pemanggilan `.toFixed()` pada `RAZTaskbar.vue`, `RAZFileManagerView.vue`, dan `RAZMobileLayout.vue`.


---

### Versi: 3.8.6-appstore-pull-and-website-ip-routing (App Store Docker Pull & IP-Based Website Routing)
**Waktu Eksekusi:** 2026-10-25 14:52 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Instalasi App Store (`/api/v1/store/install`) gagal dengan error 500 karena Docker Engine mencoba membuat kontainer sebelum image ditarik (`ImagePull`).
- Uninstall aplikasi yang belum terpasang (`DELETE /api/v1/store/apps/:slug`) mengembalikan error 500 jika kontainer belum ada di Docker.
- Pembuatan website di Website Manager tidak bisa diakses dari browser karena tautan diarahkan ke `http://localhost:<port>` alih-alih `http://<IP_STB>:<port>`.

#### 2. Solusi & Perubahan
- **Automated Image Pull (`RAZAppStoreEngine.go` & `RAZWebsiteManager.go`):** Memeriksa apakah image Docker sudah ada di cache lokal, dan jika belum ada otomatis melakukan `cli.ImagePull` dengan penanganan stream lengkap sebelum `ContainerCreate`.
- **Fault-Tolerant Uninstaller (`RAZAppStoreEngine.go`):** Menangani error "No such container" secara elegan sehingga proses uninstall atau pembersihan kontainer yang gagal terpasang tidak lagi memicu error 500.
- **Dynamic Hostname Binding (`RAZWebsitesView.vue` & `RAZWebsiteDetailView.vue`):** Mengganti referensi `localhost` dengan `window.location.hostname` dinamis sehingga tautan website langsung mengarah ke IP STB (`http://192.168.100.117:<port>`).
- **Pembaruan Bundle JS:** Melakukan rebuild frontend Vite (`index-CXrjxw0e.js`) dan memastikan seluruh bundle terdaftar di Git.

---

### Versi: 3.8.7-filemanager-text-editor-sync-and-crud (Text File Loading, Path Normalization & Complete Storage CRUD)
**Waktu Eksekusi:** 2026-10-25 15:15 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah yang Diperbaiki
- Pengguna melaporkan: "buka file text tidak load pakai file manager" (membuka file teks dari File Manager tidak memuat isi file atau tetap kosong).
- Penyebab utama:
  1. `RAZStorageEngine.go` pada `ReadFileContent`, `SaveFileContent`, `CreateFolder`, `DeleteFileOrFolder`, `RenameFileOrFolder`, `handleGetRawFile`, dan `handleGetThumbnail` tidak melakukan normalisasi path lintas sistem operasi sehingga path Windows (`C:/DATA/...`) yang tersimpan di localStorage atau dikirim frontend langsung menyebabkan error `500 Internal Server Error (stat C:/DATA/...: no such file or directory)` di Linux STB ARM64.
  2. Pada `RAZTextEditorView.vue`, saat `openFileInTab(fileItem)` dijalankan dan file sudah tercatat di `tabs.value` (misal dari sesi tab tersimpan yang gagal termuat atau kosong), fungsi langsung keluar (`return`) tanpa memuat ulang (`fetch`) isi file terbaru dari disk.
  3. `navigateSidebarHome` dan `restoreSession` di editor teks memiliki fallback hardcode `'C:/DATA'` yang memicu galat 500 saat dijalankan di lingkungan Linux.
  4. Endpoint API pembuatan file baru (`/api/v1/files/create-file`) dan unggah berkas mentah (`/api/v1/files/upload-raw`) belum terdaftar di `RAZApiController.go`.
  5. Layout mobile (`RAZMobileLayout.vue`) belum mendaftarkan view `file-manager-editor` / `text-editor`.

#### 2. Solusi & Perubahan
- **Universal Storage Path Normalizer (`RAZStorageEngine.go`):** Mengimplementasikan fungsi tunggal `NormalizeStoragePath(targetPath)` yang secara otonom mendeteksi sistem operasi (Linux vs Windows), mengonversi backslash, memetakan `C:/DATA` ke `/DATA` dan `C:` ke `/` di Linux, serta memetakan `/DATA` ke `C:/DATA` di Windows.
- **Penyelesaian Endpoint Storage (`RAZApiController.go`):**
  - Mendaftarkan rute `POST /api/v1/files/create-file` (`handleCreateFile`) untuk membuat file baru dengan konten template instan.
  - Mendaftarkan rute `POST /api/v1/files/upload-raw` (`handleUploadRaw`) untuk streaming upload berkas, struktur folder bertingkat, dan auto-extract file ZIP.
  - Menerapkan `NormalizeStoragePath` pada seluruh handler berkas (`handleGetFileContent`, `handleSaveFileContent`, `handleGetRawFile`, `handleGetThumbnail`).
- **Pemuatan Ulang Cerdas Tab Teks (`RAZTextEditorView.vue`):**
  - Memperbarui `openFileInTab`: jika tab berkas sudah ada namun belum dimodifikasi (`!existingTab.isDirty`), editor secara proaktif memanggil `/api/v1/files/content` untuk memuat isi berkas terbaru dari disk.
  - Memperbaiki fallback direktori home dan pemulihan sesi tab ke `/DATA` dinamis tanpa hardcode `C:/DATA`.
- **Integrasi Editor Mobile (`RAZMobileLayout.vue`):** Menambahkan `RAZTextEditorView` pada pemetaan view mobile.
- **Rebuild Bundle Frontend:** Melakukan kompilasi produksi Vite (`index-6xxbGSKe.js` & `index-DXBa-kC6.css`).
