# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.3.0-phase3-complete (Penyelesaian Fase 3: App Store Engine, Custom Port/Volume, & Clean Uninstall)
**Waktu Eksekusi:** 2026-10-08 23:23 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Backend App Store Engine & Proteksi RAM STB
- Menambahkan skema katalog JSON resmi di `backend/internal/store/RAZApps.json` memuat metadata Caddy, MariaDB, PostgreSQL, Adminer, Uptime Kuma, dan FileBrowser.
- Mengembangkan engine deploy container di `backend/internal/store/RAZAppStoreEngine.go`:
  - **Validasi Hardware STB:** Menolak instalasi container jika kapasitas RAM fisik host tersisa kurang dari 150 MB untuk mencegah OOM-Killer.
  - **Memory Limits per Container:** Setiap container tambahan otomatis dibatasi alokasi memori maksimal 256 MB.
  - **Isolasi MicroSD `/DATA/apps`:** Volume data container otomatis diarahkan ke partisi MicroSD eksternal.
  - **Clean Uninstall:** Fitur hapus container sekaligus opsi menghapus seluruh folder data terkait secara tuntas.
- Mendaftarkan endpoint REST API di `backend/internal/api/RAZApiController.go`:
  - `GET /api/v1/store/apps`: Mengambil status dan katalog aplikasi.
  - `POST /api/v1/store/install`: Payload instalasi custom port dan data volume.
  - `DELETE /api/v1/store/apps/:slug?purge_data=true`: Clean uninstall container dan penghapusan data.

#### 2. Frontend App Store UI & Custom Modals
- Memperbarui `RAZAppStoreView.vue`:
  - Tab kategori: *All, Network, Database, Monitoring, Storage*.
  - Modal kustom Glassmorphism untuk form konfigurasi port host dan path volume direktori MicroSD.
  - Modal konfirmasi kustom untuk **Clean Uninstall** dengan checkbox interaktif pemusnahan folder data fisik.
  - Tidak menggunakan alert/confirm native browser (sesuai aturan GEMINI.md).
- Memperbarui roadmap pada [Task_plan.md](file:///c:/laragon/www/razos/Task_plan.md) (Fase 3 Selesai 100%).
