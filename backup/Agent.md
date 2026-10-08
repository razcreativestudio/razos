# Agent Guidelines & Context: RAZOS Core Engineering

Dokumen ini menjadi acuan kerja bagi AI Agent dan pengembang dalam memelihara serta membangun RAZOS (RAZ Operating System).

> **STANDAR OPERASIONAL MUTLAK:** Seluruh aktivitas rekayasa kode, modifikasi file, dan pembaruan arsitektur WAJIB merujuk dan tunduk pada aturan di [GEMINI.md](file:///c:/laragon/www/razos/GEMINI.md).

---

## 1. Konteks Lingkungan Eksekusi & Batasan Hardware
- **Arsitektur Target:** Linux ARM64 (aarch64).
- **Target Hardware:** STB Amlogic S905X, RAM 1.8 GiB, ZRAM Swap ~900 MiB.
- **Penyimpanan:** eMMC internal 8 GB (terbatas) dan MicroSD eksternal di `/DATA` (penyimpanan utama).
- **Prinsip Utama:** Nol dependensi runtime berat (hindari runtime Python/Node.js di backend produksi). Wajib menggunakan binary statis Golang dengan konsumsi RAM di bawah 30-35 MB.

---

## 2. Aturan Rekayasa (Engineering Rules)
- **Konvensi Penamaan (Prefix RAZ):** Seluruh file baru buatan pengembang wajib memakai prefix `RAZ` (misal: `RAZMain.go`, `RAZTaskbar.vue`).
- **Sistem Backup Terpusat:** Sebelum mengedit file apa pun, wajib membuat salinan cadangan di folder `backup/`.
- **Operasi File & Docker:** Semua volume container, file database, dan artefak website wajib diarahkan ke direktori `/DATA/razos/` dan `/DATA/apps/`. Dilarang menulis data dinamis ke root `/` atau `/var/lib/docker/` tanpa pembersihan berkala.
- **Isolasi Proses:** Interaksi Docker Engine wajib melalui Unix Domain Socket resmi (`/var/run/docker.sock`) menggunakan Docker SDK Go resmi.
- **Keamanan PTY/Terminal:** WebSocket PTY wajib memiliki validasi token sesi JWT dan pembatasan izin command.
- **Responsivitas Antarmuka:** Frontend harus membedakan breakpoint desktop (>768px) dan mobile (<=768px) dengan arsitektur layout terpisah secara modular.
- **Komentar Edukatif:** Komentar kode wajib menggunakan Bahasa Indonesia yang mudah dipahami.
- **Zero Dummy Data & No Dead-End UI:** Semua data harus dinamis dan tombol tanpa fungsi tidak boleh ditampilkan ke pengguna.

---

## 3. Protokol Dokumentasi & Update History
Setiap pembaruan wajib dicatat pada `update-history/` harian, melakukan penyesuaian versi di dokumen inti (`Readme.md`, `Task_plan.md`, `Implementasi.md`, `Projectplan.md`), dan membersihkan file sementara sebelum dinyatakan selesai.
