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
