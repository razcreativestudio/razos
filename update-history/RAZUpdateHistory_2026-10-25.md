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
