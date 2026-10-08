# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.4.0-phase4-complete (Penyelesaian Fase 4: Reverse Proxy Dinamis, Website Wizard, & Database Suite)
**Waktu Eksekusi:** 2026-10-08 23:28 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Backend Reverse Proxy, Website Wizard & Database Central Suite
- Mengembangkan modul `backend/internal/proxy/RAZProxyRouter.go`: Reverse proxy dinamis dalam memori yang memetakan domain hostname langsung ke port lokal container target tanpa memerlukan reload proses Go.
- Mengembangkan modul `backend/internal/website/RAZWebsiteManager.go`:
  - **Auto-Assign Port Bebas:** Pemindaian TCP port otomatis pada rentang 8100-8200.
  - **Penyimpanan Terisolasi:** Pembuatan direktori web otomatis di `/DATA/apps/websites/<nama>`.
  - Pilihan runtime: PHP 8.2, Node.js, dan HTML statis.
- Mengembangkan modul `backend/internal/database/RAZDatabaseSuite.go`:
  - Registrasi database terpusat (MariaDB port 3306 & PostgreSQL port 5432).
  - Integrasi container Adminer Web GUI pada port internal 8088 dengan konsumsi RAM rendah (< 64 MB).
- Mendaftarkan endpoint REST API di `backend/internal/api/RAZApiController.go`:
  - `GET /api/v1/websites` & `POST /api/v1/websites` (Wizard deploy website).
  - `DELETE /api/v1/websites/:id`: Hapus website dan unbind rute proxy.
  - `GET /api/v1/databases` & `POST /api/v1/databases`: Pembuatan database baru.
  - `POST /api/v1/databases/adminer/launch`: Peluncuran instan Adminer GUI.

#### 2. Frontend UI Upgrade (Wizard Website & Manajemen DB)
- Memperbarui `RAZWebsitesView.vue`:
  - Custom Modal Wizard Pembuatan Website (Nama, Domain Hostname, Pilihan Runtime, dan Koneksi Database instan).
  - Tautan langsung pembukaan URL domain serta tombol hapus website.
- Memperbarui `RAZDatabasesView.vue`:
  - Custom Modal Pembuatan Database Baru (Nama DB, Engine MariaDB/PostgreSQL, Username, Password).
  - Tombol peluncuran 1-klik ke antarmuka Web Adminer GUI (`http://localhost:8088`).
- Memperbarui roadmap pada [Task_plan.md](file:///c:/laragon/www/razos/Task_plan.md) (Fase 4 Selesai 100%).
