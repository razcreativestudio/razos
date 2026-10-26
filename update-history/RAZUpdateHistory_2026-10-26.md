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
