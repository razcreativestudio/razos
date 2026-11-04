# RAZOS Update History: 2026-11-04

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 4.5.3-file-manager-chunked-upload-engine (Implementasi Chunked Slice Upload 4 MB & Real-time Progress Bar di File Manager)
**Waktu Eksekusi:** 2026-11-04 09:55 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Temuan Pengguna
- Pengguna melaporkan bahwa saat mengunggah berkas besar (contoh berkas arsip `eMart v6.7.rar`), proses unggah macet / *stuck* di 0% dan tidak kunjung berjalan (*"upload filenya gak jalan jalan"*).

#### 2. Investigasi & Akar Masalah (Root Cause Analysis)
1. **Ketiadaan Progres Transmisi pada `fetch()` Bawaan:**
   - Implementasi sebelumnya mengirim seluruh berkas dalam satu *request* `fetch()` raksasa tanpa *event listener* progres byte. Nilai `uploadProgress` hanya dihitung setelah seluruh berkas tuntas diterima server.
   - Akibatnya, pada berkas berukuran besar (puluhan hingga ratusan MB), persentase progress terkunci di angka 0% selama bermenit-menit sehingga tampak membeku (*freeze*).
2. **Batasan Body Size Reverse Proxy / Cloudflare Tunnel (100 MB Limit):**
   - Domain publik STB (`server.razos.my.id`) dialirkan melalui Cloudflare Tunnel. Pada paket gratis, Cloudflare membatasi payload HTTP POST maksimal sebesar 100 MB. Berkas lebih besar dari 100 MB otomatis ditolak (HTTP 413 Payload Too Large) atau diputus koneksinya di tengah jalan.
3. **CORS Headers:**
   - Header kustom pengunggahan berkas (`x-chunk-index`, `x-total-chunks`, `x-dest-dir`, dsb.) belum didaftarkan lengkap pada whitelist CORS server Go.

#### 3. Solusi & Rekayasa Arsitektur
1. **Mesin Unggah Berbasis Chunked Slicing 4 MB (`frontend/src/views/RAZFileManagerView.vue`):**
   - Berkas dipotong secara asinkron menjadi segmen-segmen kecil berukuran 4 MB menggunakan `file.slice()`.
   - Ukuran 4 MB sangat aman dan bebas dari batasan proxy Cloudflare (100 MB), serta menjaga alokasi RAM STB ARM64 tetap berada di bawah 1 MB per transfer.
   - Dilengkapi *automatic retry* hingga 3 kali per-chunk jika terjadi *network drop* atau *timeout*.
   - Input berkas di-reset nilainya agar pengguna dapat mengunggah kembali berkas yang sama tanpa perlu me-refresh halaman.
2. **Backend Streaming & Append Handler (`backend/internal/api/RAZApiController.go`):**
   - Endpoint `/api/v1/files/upload-raw` membaca header `x-chunk-index` dan `x-total-chunks`.
   - Chunk pertama (0) membuat berkas baru dengan `os.O_CREATE | os.O_TRUNC`, sedangkan chunk lanjutan (1, 2, ...) ditulis langsung menggunakan mode `os.O_APPEND` dan `io.Copy(out, c.Request.Body)` dengan penyangga memori hemat daya.
   - Kompatibilitas mundur 100% terjaga jika request tidak menyertakan chunk index (single file upload biasa / curl).
3. **Penyelarasan Whitelist CORS (`backend/cmd/server/RAZMainServer.go`):**
   - Menambahkan header `x-dest-dir`, `x-filename`, `x-relative-path`, `x-auto-extract`, `x-chunk-index`, `x-total-chunks`, `x-file-size` pada `cors.Config.AllowHeaders`.
4. **UI/UX Modal Unggah Modern:**
   - Menampilkan persentase real-time 0-100% yang bergerak per-chunk secara halus.
   - Menampilkan metrik byte terunggah vs total (contoh: `45.2 MB / 210.4 MB`).
   - Menampilkan indikator kecepatan transmisi aktual (contoh: `2.4 MB/s`).
   - Menyediakan tombol "Batal" (*Cancel Upload*) yang aman dan menghentikan pengiriman chunk secara instan.
   - Banner notifikasi error jika terjadi penolakan server.
5. **Kompilasi Ulang Frontend:**
   - Build produksi Vite (`dist/assets/index-BQp6uWBU.js`) sukses 100% tanpa error dalam 5.74 detik.

---

### Versi: 4.5.4-file-manager-gdrive-floating-background-upload (Google Drive Style Floating Minimizable Upload Manager & Anti-Refresh Guard)
**Waktu Eksekusi:** 2026-11-04 10:25 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Permintaan Pengguna
- Pengguna meminta agar proses unggah berkas dapat berjalan di latar belakang (*background process*), dapat diminimalkan ke pojok taskbar (seperti pada Google Drive) sehingga pengguna tetap bisa bebas beraktivitas, serta adanya proteksi konfirmasi/peringatan jika tab browser tidak sengaja di-refresh saat proses unggah masih berlangsung.

#### 2. Investigasi & Solusi Desain UI/UX
1. **Penghapusan Modal Backdrop Penuh (Zero Blocking):**
   - Mengganti dialog modal tengah layar yang memblokir layar (`fixed inset-0 bg-black/40`) dengan **Floating Dock Widget** menggunakan `<Teleport to="body">` yang melayang di pojok kanan bawah desktop (`fixed bottom-14 right-6 z-[9999]`).
   - Kontainer induk bersifat `pointer-events-none` sementara kartu widget bersifat `pointer-events-auto`, sehingga pengguna bebas membuka folder lain, membaca berkas, atau menjalankan aplikasi lain di File Manager / Desktop sembari proses unggah tetap berjalan di background.
2. **Dua Mode Tampilan Ala Google Drive:**
   - **Mode Penuh (Expanded Card):** Menampilkan nama berkas aktif, progress bar beranimasi, byte stats, kecepatan transfer, ringkasan jumlah item, tombol minimize, dan tombol batal.
   - **Mode Minim (Taskbar Floating Pill):** Dikompresi menjadi kapsul anggun di atas taskbar (`Mengunggah 22% (3/6) • 933 KB/s`) lengkap dengan indikator animasi dan tombol perbesar/tutup.
3. **Penguncian Path Tujuan (*Locked Destination Directory*):**
   - Jalur direktori penyimpanan dikunci saat unggahan dimulai (`lockedUploadDestDir`). Pengguna dapat berpindah-pindah folder tanpa khawatir berkas terkirim ke folder yang salah.
4. **Proteksi Anti-Refresh & Anti-Close (`beforeunload` Guard):**
   - Mendaftarkan event listener `beforeunload` saat proses unggah aktif. Jika pengguna menekan F5, Ctrl+R, atau menutup tab browser, browser secara otomatis memunculkan dialog peringatan untuk mencegah hilangnya progres unggah.
5. **Custom Modal Konfirmasi Pembatalan (Glassmorphism Slate-900):**
   - Menekan tombol Batal atau Tutup saat unggahan berlangsung tidak akan langsung memutus koneksi, melainkan memunculkan dialog konfirmasi kustom: *"Batalkan Proses Unggah?"* dengan tombol *"Lanjutkan Unggah"* dan *"Ya, Batalkan"*.
6. **Kompilasi Ulang Bundle Produksi:**
   - Build Vite (`dist/assets/index-DBZpIQut.js` & `dist/assets/index-dhKDcG6e.css`) sukses 100% tanpa error dalam 5.57 detik.

