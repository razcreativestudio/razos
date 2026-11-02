# RAZOS Update History: 2026-11-02

Dokumen riwayat pembaruan harian ekosistem RAZ Operating System (RAZOS).

---

### Versi: 4.5.0-app-store-purge-fake-and-verified-containers (Pembersihan Kontainer Palsu di App Store & Integrasi 100 Kontainer Docker Terverifikasi Siap Pakai)
**Waktu Eksekusi:** 2026-11-02 09:15 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Temuan Pengguna
- Pengguna melaporkan bahwa di dalam App Store RAZOS sebagian besar aplikasi yang tampil adalah aplikasi palsu (fake/dummy) yang tidak dapat digunakan atau ditarik dari Docker Hub (seperti `Downloaders Bridge Suite #8`, `AI Hub Suite #9`, `Media Service Suite #10`, dll. dengan image `razos/app-...:stable` yang tidak nyata).
- Pengguna meminta agar seluruh aplikasi palsu tersebut dibersihkan dan diganti/diisi dengan kontainer Docker riil yang benar-benar dapat dipasang dan digunakan di perangkat STB ARM64.

#### 2. Investigasi & Akar Masalah (Root Cause Analysis)
1. **Procedural Dummy Suites di `RAZApps.json`:**
   - Ditemukan sebanyak 454 entri dummy buatan berbasis template `Suite #1` hingga `Suite #454` dengan image fiktif `razos/app-file-gateway-1:stable`, dll.
   - Karena image tersebut tidak pernah dipublikasikan di Docker Hub atau registry publik mana pun, proses instalasi selalu gagal pada tahap Docker pull.
2. **Konflik Port Default & Port Tabrakan:**
   - Sejumlah aplikasi riil sebelumnya (seperti Nextcloud, Zigbee2MQTT, Dozzle, Trilium, Open WebUI) menggunakan `default_port: 8080` yang bentrok langsung dengan port WebUI RAZOS (8080) dan bentrok antar satu sama lain.
   - Beberapa aplikasi (Pi-hole, Caddy, Vaultwarden, BookStack) menggunakan port 80 yang bentrok dengan port HTTP utama RAZOS (80).
   - Akibatnya, saat dijalankan, Docker mengembalikan error `bind: address already in use`.
3. **Pemisahan Port Host vs Port Internal Kontainer:**
   - Struktur manifes sebelumnya hanya memiliki satu field `default_port` yang dipetakan langsung ke port internal container, sehingga aplikasi web server yang di dalamnya listening di port 80 tidak dapat dibuka jika dipetakan ke port host lain.

#### 3. Solusi & Rekayasa Kode
1. **Pembaruan Backend Engine (`backend/internal/store/RAZAppStoreEngine.go`):**
   - Menambahkan field `ContainerPort` dan `ContainerVolume` pada struktur data `AppManifest`.
   - Pada fungsi `DeployApp`, port internal container (`natPort`) diarahkan secara presisi ke `targetApp.ContainerPort` (misal port 80 internal), sedangkan port host STB (`targetPort`) dialokasikan ke port unik yang bebas konflik.
   - Menyesuaikan jalur volume mount target kontainer (`ContainerVolume`) secara spesifik sesuai kebutuhan aplikasi (misal `/config` untuk image linuxserver, `/var/lib/mysql` untuk database, `/var/www/html` untuk Nextcloud/web server).
2. **Pembersihan Total & Kurasi 100 Kontainer Riil (`backend/internal/store/RAZApps.json`):**
   - Menghapus 100% (seluruh 454 entri) aplikasi dummy `Suite #` fiktif.
   - Mempertahankan dan memvalidasi 71 aplikasi riil populer, serta menambahkan 29 kontainer self-hosted teruji baru (seperti Homarr, Homepage, Dashy, Flame, IT-Tools, CyberChef, Ntfy, Gotify, Speedtest Tracker, Beszel, Glances, Memos, Linkding, Wallabag, Shlink, Mealie, Grocy, NocoDB, PocketBase, Directus, Jellyseerr, MeTube, PyLoad-ng, Mosquitto, Calibre Desktop, RustDesk Server, Watchtower, Tailscale, Cloudflared).
   - Memastikan seluruh 100 kontainer memiliki image Docker Hub/GHCR multi-arch dengan dukungan resmi arsitektur ARM64 (`linux/arm64`) yang kompatibel dengan STB Amlogic S905X.
   - Memverifikasi 0 bentrok port host (`port conflicts: 0`), dan memastikan tidak ada satu pun aplikasi yang menggunakan port host 80 atau 8080.
3. **Penyelarasan Tampilan Frontend (`frontend/src/views/RAZAppStoreView.vue`):**
   - Menghapus teks statis `525+` dan mengubah badge katalog menjadi dinamis mengikuti jumlah aplikasi riil `Katalog ({{ apps.length }})`.
   - Memperbarui placeholder pencarian menjadi `"Cari aplikasi Docker (Plex, Jellyfin, MySQL, Nextcloud, Ollama, Nginx, WireGuard...)"`.
   - Mengubah teks tombol empty state menjadi `"Buka Katalog Aplikasi"`.
4. **Pembaruan Mock Vite (`frontend/vite.config.js`):**
   - Menyelaraskan komentar dan sinkronisasi data katalog resmi.
5. **Kompilasi Ulang Frontend:**
   - Menjalankan build produksi Vite (`dist/assets/index-Ci8-d08i.js` & `dist/assets/index-DSiB8-2j.css`) 100% sukses tanpa error dalam 4.96 detik.
