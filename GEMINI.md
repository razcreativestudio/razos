# GEMINI.md - SYSTEM RULES & ENTERPRISE PROJECT STANDARDS

**PENTING:** Aturan ini adalah **Protokol Mutlak**. Setiap kali AI Agent atau pengembang membuat rencana arsitektur, menulis kode, meninjau struktur, atau membantu pengembangan untuk seluruh ekosistem RAZOS (RAZ Operating System), WAJIB mematuhi seluruh aturan di bawah ini tanpa kecuali.

---

## 1. FILE MANAGEMENT & NAMING CONVENTION (STANDAR RAZ)
* **Prefix Wajib:** Selalu gunakan awalan `RAZ` untuk setiap penamaan file baru yang dibuat di dalam proyek (contoh: `RAZMainServer.go`, `RAZMetricsCollector.go`, `RAZDesktopLayout.vue`, `RAZSettings.css`, `RAZLogo.svg`, `RAZBrandCapsule.png`).
* **Sistem Backup Terpusat:** WAJIB membuat salinan *backup* sebelum mengedit file apa pun. Simpan file tersebut di dalam folder khusus `backup/`. Jika terjadi masalah teknis atau file baru mengalami *corrupt*, segera *restore* menggunakan file *backup* tersebut.
* **Modularisasi File (Pemecahan Beban):** Jika kode dalam satu file menjadi terlalu panjang atau padat sehingga rawan memicu *error* (misal: file mencapai 2000+ baris atau komponen terlalu padat), **wajib pecah file tersebut** berdasarkan fungsi atau *layout*-nya menjadi file-file spesifik yang lebih kecil (misal: `RAZMainLayout.vue` dipecah menjadi `RAZDesktopLayout.vue`, `RAZMobileLayout.vue`, `RAZTaskbar.vue`, dll).
* **Clean Workspace (Hapus File Sampah):** File yang dibuat khusus untuk pengecekan *error*, *testing*, *inject*, *migrate*, atau *hotfix* WAJIB langsung dihapus seluruhnya jika sudah selesai dieksekusi dan tidak lagi dibutuhkan di tahap *production*. Jangan biarkan ada *file* menumpuk! Hindari inject skrip pihak ketiga sembarangan untuk perubahan kode internal.

---

## 2. PROTOKOL UPDATE & DOKUMENTASI (STRICT VERSIONING)
Setiap kali terjadi *update* (perbaikan, penambahan fitur, perubahan kode), wajib jalankan rutinitas dokumentasi berikut:
* **Metadata File:** Selalu perbarui catatan waktu (*timestamp*) dan deskripsi kegiatan *update* di bagian paling atas (header komentar) setiap file kode yang diedit.
* **Update History Harian:** Catat setiap kegiatan perubahan di dalam file *Update History* dengan sangat mendetail. Format file riwayat: `update-history/RAZUpdateHistory_YYYY-MM-DD.md`.
* **Sinkronisasi Dokumen Inti:** Setiap pembaruan kode harus diikuti pembaruan dokumen. Wajib melakukan pembaruan versi (naik versi), penyesuaian di `Projectplan.md`, penyesuaian struktur direktori di `Implementasi.md`, dan sinkronisasi `Agent.md`, `Task_plan.md`, dll.
* **File Implementasi Bertahap:** Selalu mutakhirkan dokumen implementasi agar progres pengerjaan terstruktur, bertahap, dan mudah dilacak titik masalahnya.
* **README.md Sebagai *Single Source of Truth*:** Dokumen `Readme.md` HARUS menjadi pusat informasi paling *up-to-date* untuk arsitektur, fitur, status versi, dan instruksi deployment.

---

## 3. STANDAR KODE BERKUALITAS (EDUKATIF & RAPI)
* **Komentar Edukatif Bahasa Indonesia:** Di setiap blok fungsi, komponen, atau handler logika, berikan keterangan (komentar) yang jelas menggunakan Bahasa Indonesia agar mudah dipelajari dan dipelihara.
* **Zero Dummy Data:** Dilarang menggunakan data palsu (*hardcode*) di frontend. Semua data (daftar aplikasi, status CPU/RAM/Suhu, website, database, proses, container) harus dinamis dari backend/sistem/database SQLite.
* **Micro-Detail Retention:** Tidak boleh ada satu pun peluru poin spesifikasi fitur klien atau prasyarat arsitektur STB yang dilewati demi mempercepat pengerjaan.

---

## 4. UI/UX PRESISI TINGGI & ENTERPRISE LOOK
* **Hanya Custom Modals:** SEMUA bentuk dialog interaksi (konfirmasi instalasi/hapus, peringatan kill task, notifikasi sistem, prompt) WAJIB menggunakan custom modal terintegrasi dengan template UI (Glassmorphism Slate-900). Dilarang keras menggunakan `alert()`, `confirm()`, atau `prompt()` bawaan browser.
* **Desain Premium Tanpa Cacat:** Antarmuka (UI/UX) modern, profesional, presisi tinggi (desktop multi-window & mobile launcher adaptif pada breakpoint 768px). Grid rapi, responsif tanpa overlap margin/padding atau horizontal scroll tak diinginkan.
* **Hide, Don't Fake (No Dead-End):** Jika fungsionalitas fitur belum selesai, sembunyikan tombol/aksinya. Jangan membuat ikon atau link yang tidak menghasilkan aksi apa pun atau kosong saat diklik.

---

## 5. MAKSIMALISASI SKILLS, ASSETS MULTIMEDIA & BRAND INTEGRITY
* **Pemanfaatan Kapabilitas Penuh (Skills Activation):** Agen WAJIB memanfaatkan seluruh kapabilitas dan skills yang tersedia (antislop, generative UI, clean code/architecture, image/video tools, dll) secara proaktif untuk meningkatkan kualitas teknis dan estetika proyek.
* **Penyediaan Aset Mandiri:** Apabila proyek memerlukan gambar latar, ikon aplikasi, logo, mockup UI, ilustrasi status, sound/video splash, atau aset visual lainnya, agen WAJIB memproduksi/men-generate aset tersebut (menggunakan generator gambar, vector SVG, generator visual) secara otonom tanpa menunggu diminta.
* **Integritas Logo Asli & Watermark Wajib:**
  - Folder resmi brand logo: `frontend/public/assets/brand/` (`RAZBrandCapsule.png`, `RAZBrandSquare.jpg`, `RAZBrandTransparent.png`).
  - **Dilarang merusak keaslian bentuk geometris atau tipografi logo RAZ.** Setiap pengeditan hanya boleh dilakukan untuk penyesuaian overlay/komposisi tanpa mendistorsi rasio atau elemen khas logo.
  - Setiap aset visual yang digenerate (wallpaper, banner, splash screen, kartu UI) atau yang membutuhkan identitas resmi WAJIB menyertakan logo resmi RAZ sebagai logomark atau subtle watermark.
* **Standar Aset Ringan & Teroptimasi:** Aset visual harus teroptimasi untuk performa hardware rendah (STB ARM64). Gunakan **SVG murni** untuk seluruh ikon UI dan vektor. Untuk seluruh gambar grafis/bitmap/wallpaper/banner/foto, **WAJIB menggunakan format WebP** dengan kompresi efisien demi kecepatan loading dan penghematan memori RAM/browser.
* **Ekosistem & Ekstensi Skills:** Lakukan instalasi atau penyiapan konfigurasi skills/plugins yang relevan jika diperlukan oleh arsitektur proyek guna mempercepat penyelesaian tugas dengan akurasi tinggi.

---

## 6. REKAYASA SISTEM RAZOS (CORE CONTEXT & HARDWARE LIMITS)
* **Target Hardware:** STB Amlogic S905X (ARM64 aarch64), RAM 1.8 GiB, ZRAM Swap ~900 MiB, eMMC 8 GB (terbatas), MicroSD di `/DATA` (penyimpanan utama).
* **Zero Heavy Runtimes:** Backend wajib Golang murni dengan konsumsi RAM idle <= 35 MB dan ukuran binary <= 25 MB (termasuk static asset frontend yang di-embed via `embed.FS`). Dilarang menggunakan Python/Node.js sebagai runtime backend inti.
* **Operasi File & Docker:** Semua container, volume, database, dan aplikasi diarahkan ke `/DATA/razos/` dan `/DATA/apps/`. Dilarang menulis data dinamis ke root `/` atau `/var/lib/docker/`.
* **Isolasi Proses:** Akses Docker hanya melalui socket Unix lokal `/var/run/docker.sock` dengan Go Docker SDK.
* **Keamanan Terminal:** WebSocket PTY terenkapsulasi token JWT dan batasan command root.

---

## 7. DEFINITION OF DONE (DoD)
Setiap tahapan/tugas dianggap "SELESAI" hanya jika:
1. Logika bisnis berjalan tepat 100% tanpa kompromi fitur dan sesuai limit RAM/resource STB.
2. Tidak ada pesan error atau peringatan merah di browser console maupun backend Go.
3. Seluruh aset yang dibutuhkan (ikon, wallpaper, logo brand, dsb) telah tersedia dan terhubung tanpa link rusak (*broken image*), serta brand mark RAZ diaplikasikan sesuai panduan.
4. Seluruh dokumen telah disinkronkan (`Readme.md`, `Task_plan.md`, `Implementasi.md`, `Update History`), versi dinaikkan, dan backup telah diamankan di `backup/`.
5. Tampilan presisi baik di Desktop mode (>768px) maupun Mobile launcher (<=768px).
