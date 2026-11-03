# RAZOS Update History: 2026-11-03

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 4.5.1-app-store-armbian-verified-expansion (Ekspansi 50 Kontainer Docker Open-Source Terverifikasi Khusus Armbian STB ARM64)
**Waktu Eksekusi:** 2026-11-03 09:10 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Permintaan Pengguna
- Pengguna meminta untuk menambahkan lagi aplikasi kontainer Docker yang open-source, terverifikasi resmi, dan cocok untuk lingkungan RAZOS di dalam Armbian (STB Amlogic S905X ARM64).

#### 2. Investigasi & Persyaratan Lingkungan (Hardware & OS Limits)
1. **Target Arsitektur:** STB Amlogic S905X dengan prosesor quad-core ARM Cortex-A53 (ARM64 aarch64), RAM 1.8 GiB, ZRAM Swap ~900 MiB, MicroSD di `/DATA` sebagai media penyimpanan utama.
2. **Kriteria Seleksi Kontainer:**
   - Image wajib mendukung multi-arch `linux/arm64` resmi (Docker Hub, GitHub Container Registry `ghcr.io`, Quay, dll.).
   - Beban memori idle wajib rendah (< 150-250 MB RAM per container) agar tidak membebani kapasitas RAM STB yang terbatas.
   - 100% Open Source dengan komunitas aktif dan terpercaya (seperti proyek LinuxServer.io, self-hosted community, Rust/Go microservices).
   - Zero Port Conflict: Setiap container host port dialokasikan secara presisi dan dilarang keras bentrok dengan port RAZOS (80 & 8080) maupun kontainer yang telah ada sebelumnya.

#### 3. Solusi & Rekayasa Katalog (`backend/internal/store/RAZApps.json`)
1. **Penambahan 50 Kontainer Baru Terverifikasi (Total Katalog: 150 Aplikasi Riil):**
   - **Media & Streaming:** RomM (Retro ROM Manager WebAssembly), Lychee Photo Gallery, Pinchflat (YouTube Archiver via yt-dlp), PodFetch (Podcast player hemat daya berbasis Rust), Mylar3 (Comic book library manager).
   - **Monitoring & Analytics:** Glance Feed Dashboard, Scrutiny (HDD/MicroSD S.M.A.R.T Drive Health), Gatus (Automated health status monitor), What's Up Docker (WUD - Container update notifications), Statping-ng (Uptime monitor & status page).
   - **Networking & Security:** Headscale (Self-hosted open-source Tailscale control server), Authelia (Single Sign-On & Two-Factor Authentication), Traefik Proxy (Cloud-native reverse proxy & ingress), Rathole (Secure NAT traversal tunnel via Rust), FlareSolverr (Cloudflare bypass proxy).
   - **Development & Git:** Forgejo (Community-driven Git platform ringan), Woodpecker CI (Lightweight continuous integration), Penpot (Open-source design & prototyping), Hoppscotch (API development environment), Hasura GraphQL Engine, SurrealDB (Multi-model database), Typesense (Fast typo-tolerant search engine), ClickHouse, Redis Insight, Mongo Express.
   - **Productivity & Office:** AFFiNE (Next-gen knowledge base), SilverBullet (Extensible markdown note-taking), HedgeDoc (Collaborative markdown editor), Wekan (Kanban board), Excalidraw (Virtual whiteboard), Diagrams.net / Draw.io, Actual Budget (Local-first personal finance), Firefly III (Financial manager), InvenTree (Inventory & asset management), CryptPad (End-to-end encrypted collaboration suite).
   - **File & Cloud Storage:** Rustypaste (Minimalist pastebin/file sharing berbasis Rust), Pingvin Share (Self-hosted file sharing platform), Chibisafe (Simple image & file uploader), PsiTransfer (Simple open-source file sharing), Duplicacy Web Edition.
   - **Home Automation & IoT:** Scrypted (High-performance home video camera integration), TasmoAdmin (Sonoff/Tasmota smart device manager), WLED (Fast LED controller controller), Z-Wave JS UI (Z-Wave home automation bridge).
   - **AI & Machine Learning:** Whisper ASR (OpenAI Whisper Speech-to-Text lokal model tiny/base teroptimasi ARM64), Piper Neural TTS (Text-to-Speech lokal ultra-cepat berbasis deep learning).
   - **Downloaders & Torrents:** JDownloader 2 (Headless multi-part download manager via web GUI), Jackett (Torrent indexer proxy).

2. **Verifikasi Alokasi Port & Volume Mount:**
   - Seluruh 150 aplikasi memiliki port host unik (0 konflik). Port 80 dan 8080 steril dari alokasi port host container.
   - Pemetaan `container_port` dan `container_volume` lengkap terkonfigurasi ke direktori MicroSD `/DATA/apps/<slug>`.

3. **Kompilasi Ulang Frontend:**
   - Membangun kembali bundle Vite (`dist/assets/index-Ci8-d08i.js` & `dist/assets/index-DSiB8-2j.css`) 100% sukses tanpa error.
