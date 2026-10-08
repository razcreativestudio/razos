# RAZOS Update History: 2026-10-08

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 1.5.0-production-ready (Penyelesaian Fase 5: OTA Updater, Systemd Hardening & Pembersihan Host)
**Waktu Eksekusi:** 2026-10-08 23:32 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Backend Over-The-Air (OTA) Updater Engine
- Mengembangkan modul `backend/internal/updater/RAZUpdaterEngine.go`:
  - Pengecekan rilis versi baru secara aman ke GitHub Releases API (`razcreative/razos`).
  - Fitur `PerformSelfUpdate`: Unduhan binary otomatis yang menggantikan binary aktif di `/DATA/razos/bin/razos` secara atomik.
- Mendaftarkan endpoint updater di `backend/internal/api/RAZApiController.go`:
  - `GET /api/v1/updater/status`: Membaca status versi saat ini dan rilis terbaru.
  - `POST /api/v1/updater/apply`: Menjalankan update mandiri seketika.

#### 2. Skrip Produksi & Systemd Hardening
- Menambahkan skrip bash eksekusi mandiri di folder `scripts/`:
  - `scripts/razos.service`: Unit file systemd resmi dengan proteksi batas memori STB (`MemoryHigh=1500M`, `MemoryMax=1700M`), auto-restart pada crash, dan eksekusi direktori MicroSD `/DATA/razos`.
  - `scripts/clean-host.sh`: Skrip pembersihan total sisa container CasaOS, 1Panel, dan docker system prune sebelum fresh install.
  - `scripts/install.sh`: Skrip instalasi 1-baris cepat dengan auto-detect arsitektur Linux ARM64.

#### 3. Frontend OTA Updater UI & Penyelesaian Seluruh Roadmap
- Memperbarui `RAZSettingsView.vue`:
  - Kartu OTA Updater interaktif dengan status rilis, changelog, dan tombol "Perbarui Sekarang".
- Memperbarui dokumen [Readme.md](file:///c:/laragon/www/razos/Readme.md) dengan instruksi instalasi resmi.
- Memperbarui [Task_plan.md](file:///c:/laragon/www/razos/Task_plan.md): Seluruh target Fase 1, 2, 3, 4, dan 5 telah terselesaikan 100% (*All Tasks Done*).
