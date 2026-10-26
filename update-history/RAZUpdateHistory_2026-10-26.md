# RAZOS Update History: 2026-10-26

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 3.8.8-multi-drive-storage-widget (Multi-Drive Desktop Widget & Real Mountpoint Naming)
**Waktu Eksekusi:** 2026-10-26 08:58 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna meminta agar widget status sistem di desktop dapat menampilkan seluruh pembaca drive (storage/partisi fisik) yang aktif di STB (bukan hanya `/DATA`), dan nama storage ditampilkan menggunakan nama aslinya saja (misal `/DATA`, `/`, atau jalur mount asli perangkat eksternal USB/HDD) tanpa label artifisial.

#### 2. Solusi & Rekayasa Kode
1. **Pembaruan Backend Metrics Collector (`RAZMetricsCollector.go`):**
   - Menambahkan struktur data `DiskMetric` yang mencakup `name`, `mountpoint`, `device`, `fstype`, `total`, `used`, `free`, dan `percent`.
   - Mengintegrasikan pemanggilan `filemanager.GetStorageDevices()` ke dalam `GetMetrics()` sehingga field `disks` pada `/api/v1/system/metrics` mengembalikan seluruh partisi fisik aktif secara real-time.
   - Mengatur nama asli mountpoint (`/DATA`, `/`, dll) sebagai nilai nama identitas disk.
2. **Pembersihan Nama Storage (`RAZStorageEngine.go`):**
   - Mengubah penetapan `Label` dan `Name` pada `GetStorageDevices` agar menggunakan nama asli mountpoint partisi tanpa embel-embel teks panjang (`Primary Storage...`, `System Storage...`).
3. **Penyelarasan Store Sistem (`RAZSystemStore.js`):**
   - Mendaftarkan array `disks: []` pada state bawaan `metrics`.
   - Menambahkan fungsi `fetchStorageDevices()` sebagai fallback multi-layer jika data array `disks` belum terisi dari pooling metrik.
4. **Pembaruan Widget Desktop (`RAZDesktopLayout.vue`):**
   - Membuat computed property `storageList` yang secara dinamis menyusun daftar seluruh drive fisik aktif.
   - Mengganti tampilan widget disk tunggal yang sebelumnya ter-hardcode `Storage (/DATA)` menjadi list multi-drive yang merender setiap storage beserta progress bar, kapasitas terpakai, dan nama aslinya dalam font monospaced yang presisi.
5. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`index-CdT-ZkUD.js` & `index-DXBa-kC6.css`).

---

### Versi: 3.8.9-website-manager-port-mapping-and-ip-url-fix (Website Manager Docker Port Mapping & Direct IP Access Fix)
**Waktu Eksekusi:** 2026-10-26 09:12 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Website yang dibuat melalui Website Manager (`webraz`, port 8100, Static HTML Nginx) tidak dapat dibuka:
  1. Tombol "Buka" di Website Manager membuka tab `about:blank`.
  2. Mengakses manual melalui `http://192.168.100.117:8100` menghasilkan error browser "Situs ini tidak dapat dijangkau" (Connection Refused).

#### 2. Akar Masalah (Root Cause Analysis)
1. **URL Buka Mengarah ke Domain Lokal yang Belum Ada di DNS Client:**
   - Template tautan pada `RAZWebsitesView.vue` dan `RAZWebsiteDetailView.vue` menggunakan `:href="'http://' + (site.domain || currentHostname) + ':' + site.port"`.
   - Wizard Website Manager secara otomatis mengisi domain dengan `slug.local` (contoh: `webraz.local`).
   - Di komputer klien pengguna, domain `.local` tidak terdaftar pada DNS LAN lokal maupun `/etc/hosts`, sehingga browser membatalkan resolusi nama host dan membuka halaman kosong (`about:blank` / `ERR_NAME_NOT_RESOLVED`).
2. **Port Docker Tidak Di-expose pada `container.Config` Backend Go:**
   - Di `backend/internal/website/RAZWebsiteManager.go`, pemanggilan `cli.ContainerCreate` hanya menyertakan `PortBindings` pada `HostConfig`, tetapi luput mendefinisikan `ExposedPorts` di `container.Config`.
   - Pada Docker Go SDK, tanpa `ExposedPorts: nat.PortSet{natPort: struct{}{}}`, Docker Daemon tidak membuka proxy listening pada host `0.0.0.0:8100`.
3. **Sinkronisasi Status & Self-Healing Container:**
   - Jika container mati atau gagal start saat dibuat, status di UI tetap menampilkan data statis JSON tanpa verifikasi `ContainerInspect`.

#### 3. Solusi & Rekayasa Kode
1. **Pembaruan Docker Container Config (`RAZWebsiteManager.go`):**
   - Menambahkan `ExposedPorts: nat.PortSet{natPort: struct{}{}}` dan `WorkingDir: mountDest` pada `container.Config`.
   - Menambahkan default command untuk runtime `golang` (`go run main.go`), `nodejs` (`node server.js`), dan `python` (`python app.py`) jika command kosong agar container tidak langsung exit.
   - Memastikan starter template file dan direktori dokumen selalu dibuat sebelum container dijalankan.
   - Menambahkan mekanisme self-healing pada `ToggleWebsiteStatus`: jika kontainer tidak ditemukan atau start gagal saat ditekan tombol Start/Restart, sistem secara otomatis melakukan re-deploy container baru dengan port exposure yang valid.
   - Mengintegrasikan pengecekan `cli.ContainerInspect` secara real-time pada `GetWebsites()` agar status UI selalu sinkron dengan kenyataan status container Docker.
2. **Pembaruan UI Link & Tombol Buka (`RAZWebsitesView.vue` & `RAZWebsiteDetailView.vue`):**
   - Mengubah rute tautan tombol "Buka" dan link hostname menjadi `http://${currentHostname}:${site.port}` (`http://192.168.100.117:8100`) sehingga langsung mengarah ke IP host STB yang valid di jaringan LAN.
   - Menampilkan informasi `site.domain` sebagai badge terpisah untuk konfigurasi reverse proxy.
3. **Penyelarasan Wizard (`RAZWebsiteWizardView.vue`):**
   - Mengubah label field menjadi "Domain Reverse Proxy (Opsional)" agar pengguna memahami bahwa website selalu dapat diakses langsung via IP dan port.
4. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-CS45O_7f.js` & `dist/assets/index-BaanchvU.css`).

---

### Versi: 3.9.0-app-store-docker-pull-and-clean-uninstall-fix (App Store Docker Pull Engine & Deep Clean Error Fix)
**Waktu Eksekusi:** 2026-10-26 09:21 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Setiap kali pengguna memasang aplikasi dari App Store (misal Immich, dll.), proses instalasi selalu gagal dengan notifikasi modal:
  `Error response from daemon: No such image: ghcr.io/immich-app/immich-server:release`
- Ketika pengguna mengklik tombol "Deep Clean & Hapus Sampah", endpoint DELETE mengembalikan error 500 karena kontainer belum sempat terdaftar di Docker daemon.

#### 2. Akar Masalah (Root Cause Analysis)
1. **Silenced Error pada `cli.ImagePull` Docker Go SDK:**
   - Pada `backend/internal/store/RAZAppStoreEngine.go`, saat image belum ada di lokal host, backend memanggil `cli.ImagePull(...)` dan membuang output stream-nya menggunakan `io.Copy(io.Discard, out)`.
   - Di Docker Engine SDK, bila terjadi kegagalan unduh (seperti network timeout, DNS failure, 401 unauth, atau rate limit), error dikirimkan di dalam baris-baris JSON stream (`{"errorDetail":{"message":"..."},"error":"..."}`).
   - Karena stream dibuang begitu saja tanpa di-decode dan tanpa verifikasi ulang `ImageInspectWithRaw`, backend langsung melanjutkan eksekusi ke `ContainerCreate`. Docker daemon pun merespons error `No such image`.
2. **Missing `ExposedPorts` di `container.Config`:**
   - Binding port host ke container membutuhkan deklarasi `ExposedPorts` eksplisit di `container.Config` agar network layer host memetakan port publik secara presisi.
3. **Rollback Deep Clean Error 500:**
   - Pemanggilan `CleanUninstall` melemparkan error jika kontainer yang hendak dihapus memang belum pernah berhasil dibuat (`no such container`).

#### 3. Solusi & Rekayasa Kode
1. **Penguatan Engine Pull Image (`RAZAppStoreEngine.go`):**
   - Mengimplementasikan fungsi `pullDockerImage` berbasis multi-layer:
     - **Prioritas 1 (Docker CLI):** Menjalankan `/usr/bin/docker pull` dengan konteks timeout 10 menit (sangat andal menangani layer besar, retry, multi-arch ARM64, dan DNS host).
     - **Fallback (Go SDK):** Menggunakan `cli.ImagePull` dengan JSON message streaming decoder untuk menangkap detail pesan error asli dari daemon secara presisi.
     - **Verifikasi Akhir:** Melakukan `ImageInspectWithRaw` untuk memastikan image 100% tersimpan sebelum `ContainerCreate` dipanggil.
   - Menambahkan `ExposedPorts: nat.PortSet{natPort: struct{}{}}` pada `container.Config`.
   - Menyelaraskan target mount volume kontainer agar menunjuk ke `/data` jika default volume adalah path host `/DATA/...`.
2. **Self-Healing pada `CleanUninstall`:**
   - Mengabaikan pesan error `no such container` / `not found` sehingga deep clean selalu berhasil (status 200 OK) dan folder `/DATA/apps/<slug>` langsung terhapus bersih tanpa error 500.
3. **Penyelarasan UI Instalasi (`RAZAppStoreInstallView.vue`):**
   - Memperbarui status teks pada Tahap 5 agar transparan menginformasikan proses download layer image dan inisialisasi kontainer.
4. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-sIXCvAIv.js`).
