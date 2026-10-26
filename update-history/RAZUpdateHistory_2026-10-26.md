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

---

### Versi: 3.8.9-website-manager-port-mapping-and-ip-url-fix (Website Manager Docker Port Mapping & Direct IP Access Fix)
**Waktu Eksekusi:** 2026-10-26 09:12 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Website yang dibuat melalui Website Manager (`webraz`, port 8100, Static HTML Nginx) tidak dapat dibuka:
  1. Tombol "Buka" di Website Manager membuka tab `about:blank`.
  2. Mengakses manual melalui `http://192.168.100.117:8100` menghasilkan error browser "Situs ini tidak dapat dijangkau" (Connection Refused).

#### 2. Akar Masalah (Root Cause Analysis)
1. **URL Buka Mengarah ke Domain Lokal yang Belum Ada di DNS Client:**
   - Template tautan pada `RAZWebsitesView.vue` dan `RAZWebsiteDetailView.vue` menggunakan `:href="'http://' + (site.domain || currentHostname) + ':' + site.port"`.
   - Wizard Website Manager secara otomatis mengisi domain dengan `slug.local` (contoh: `webraz.local`).
   - Di komputer klien pengguna, domain `.local` tidak terdaftar pada DNS LAN lokal maupun `/etc/hosts`, sehingga browser membatalkan resolusi nama host dan membuka halaman kosong (`about:blank` / `ERR_NAME_NOT_RESOLVED`).
2. **Port Docker Tidak Di-expose pada `container.Config` Backend Go:**
   - Di `backend/internal/website/RAZWebsiteManager.go`, pemanggilan `cli.ContainerCreate` hanya menyertakan `PortBindings` pada `HostConfig`, tetapi luput mendefinisikan `ExposedPorts` di `container.Config`.
   - Pada Docker Go SDK, tanpa `ExposedPorts: nat.PortSet{natPort: struct{}{}}`, Docker Daemon tidak membuka proxy listening pada host `0.0.0.0:8100`.
3. **Sinkronisasi Status & Self-Healing Container:**
   - Jika container mati atau gagal start saat dibuat, status di UI tetap menampilkan data statis JSON tanpa verifikasi `ContainerInspect`.

#### 3. Solusi & Rekayasa Kode
1. **Pembaruan Docker Container Config (`RAZWebsiteManager.go`):**
   - Menambahkan `ExposedPorts: nat.PortSet{natPort: struct{}{}}` dan `WorkingDir: mountDest` pada `container.Config`.
   - Menambahkan default command untuk runtime `golang` (`go run main.go`), `nodejs` (`node server.js`), dan `python` (`python app.py`) jika command kosong agar container tidak langsung exit.
   - Memastikan starter template file dan direktori dokumen selalu dibuat sebelum container dijalankan.
   - Menambahkan mekanisme self-healing pada `ToggleWebsiteStatus`: jika kontainer tidak ditemukan atau start gagal saat ditekan tombol Start/Restart, sistem secara otomatis melakukan re-deploy container baru dengan port exposure yang valid.
   - Mengintegrasikan pengecekan `cli.ContainerInspect` secara real-time pada `GetWebsites()` agar status UI selalu sinkron dengan kenyataan status container Docker.
2. **Pembaruan UI Link & Tombol Buka (`RAZWebsitesView.vue` & `RAZWebsiteDetailView.vue`):**
   - Mengubah rute tautan tombol "Buka" dan link hostname menjadi `http://${currentHostname}:${site.port}` (`http://192.168.100.117:8100`) sehingga langsung mengarah ke IP host STB yang valid di jaringan LAN.
   - Menampilkan informasi `site.domain` sebagai badge terpisah untuk konfigurasi reverse proxy.
3. **Penyelarasan Wizard (`RAZWebsiteWizardView.vue`):**
   - Mengubah label field menjadi "Domain Reverse Proxy (Opsional)" agar pengguna memahami bahwa website selalu dapat diakses langsung via IP dan port.
4. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-CS45O_7f.js` & `dist/assets/index-BaanchvU.css`).

---

### Versi: 3.9.0-app-store-docker-pull-and-clean-uninstall-fix (App Store Docker Pull Engine & Deep Clean Error Fix)
**Waktu Eksekusi:** 2026-10-26 09:21 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Setiap kali pengguna memasang aplikasi dari App Store (misal Immich, dll.), proses instalasi selalu gagal dengan notifikasi modal:
  `Error response from daemon: No such image: ghcr.io/immich-app/immich-server:release`
- Ketika pengguna mengklik tombol "Deep Clean & Hapus Sampah", endpoint DELETE mengembalikan error 500 karena kontainer belum sempat terdaftar di Docker daemon.

#### 2. Akar Masalah (Root Cause Analysis)
1. **Silenced Error pada `cli.ImagePull` Docker Go SDK:**
   - Pada `backend/internal/store/RAZAppStoreEngine.go`, saat image belum ada di lokal host, backend memanggil `cli.ImagePull(...)` dan membuang output stream-nya menggunakan `io.Copy(io.Discard, out)`.
   - Di Docker Engine SDK, bila terjadi kegagalan unduh (seperti network timeout, DNS failure, 401 unauth, atau rate limit), error dikirimkan di dalam baris-baris JSON stream (`{"errorDetail":{"message":"..."},"error":"..."}`).
   - Karena stream dibuang begitu saja tanpa di-decode dan tanpa verifikasi ulang `ImageInspectWithRaw`, backend langsung melanjutkan eksekusi ke `ContainerCreate`. Docker daemon pun merespons error `No such image`.
2. **Missing `ExposedPorts` di `container.Config`:**
   - Binding port host ke container membutuhkan deklarasi `ExposedPorts` eksplisit di `container.Config` agar network layer host memetakan port publik secara presisi.
3. **Rollback Deep Clean Error 500:**
   - Pemanggilan `CleanUninstall` melemparkan error jika kontainer yang hendak dihapus memang belum pernah berhasil dibuat (`no such container`).

#### 3. Solusi & Rekayasa Kode
1. **Penguatan Engine Pull Image (`RAZAppStoreEngine.go`):**
   - Mengimplementasikan fungsi `pullDockerImage` berbasis multi-layer:
     - **Prioritas 1 (Docker CLI):** Menjalankan `/usr/bin/docker pull` dengan konteks timeout 10 menit (sangat andal menangani layer besar, retry, multi-arch ARM64, dan DNS host).
     - **Fallback (Go SDK):** Menggunakan `cli.ImagePull` dengan JSON message streaming decoder untuk menangkap detail pesan error asli dari daemon secara presisi.
     - **Verifikasi Akhir:** Melakukan `ImageInspectWithRaw` untuk memastikan image 100% tersimpan sebelum `ContainerCreate` dipanggil.
   - Menambahkan `ExposedPorts: nat.PortSet{natPort: struct{}{}}` pada `container.Config`.
   - Menyelaraskan target mount volume kontainer agar menunjuk ke `/data` jika default volume adalah path host `/DATA/...`.
2. **Self-Healing pada `CleanUninstall`:**
   - Mengabaikan pesan error `no such container` / `not found` sehingga deep clean selalu berhasil (status 200 OK) dan folder `/DATA/apps/<slug>` langsung terhapus bersih tanpa error 500.
3. **Penyelarasan UI Instalasi (`RAZAppStoreInstallView.vue`):**
   - Memperbarui status teks pada Tahap 5 agar transparan menginformasikan proses download layer image dan inisialisasi kontainer.
4. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-sIXCvAIv.js`).

---

### Versi: 3.9.1-terminal-resize-sigwinch-and-motd-banner (Web Terminal PTY Resize Sync & Armbian MOTD Banner)
**Waktu Eksekusi:** 2026-10-26 09:27 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Saat jendela Web Terminal di RAZOS dimaksimalkan / diperbesar, tampilan terminal tidak melebar memenuhi layar (terkunci di kotak sempit sebelah kiri dengan ruang hitam kosong di kanan), dan baris teks terpotong atau wrap aneh dengan tanda `>`.
- Terminal tidak menampilkan banner sambutan resmi sistem Armbian (logo ASCII merah `Armbian OS`, statistik performa, suhu CPU, penggunaan memori/disk, dsb) seperti halnya saat login via PuTTY SSH.

#### 2. Akar Masalah (Root Cause Analysis)
1. **Tidak Ada Handler Resize Window PTY (Missing SIGWINCH):**
   - Di `backend/internal/terminal/RAZTerminalHandler.go`, PTY diinisialisasi tanpa ukuran dinamis dan WebSocket reader hanya meneruskan raw string ke `ptmx.Write` tanpa handler `pty.Setsize`.
   - Di frontend (`RAZTerminalView.vue`), `handleResize` hanya memanggil `fitAddon.fit()`, namun **tidak pernah mengirim pesan cols/rows ke WebSocket backend**.
   - Selain itu, event listener hanya terpasang di `window.onresize` peramban (bukan pada elemen kontainer jendela desktop OS), sehingga saat jendela desktop dimaksimalkan, event tidak memicu penyesuaian ukuran.
   - Akibatnya, driver PTY Linux kernel tetap menganggap ukuran terminal adalah 80x24 kolom, sehingga bash membungkus (*line-wrap*) teks di kolom ke-80 dan menyisakan ruang kosong besar di layar yang dimaksimalkan.
2. **Bash Dijalankan Sebagai Non-Login Shell:**
   - Backend memanggil `exec.Command(shell)` tanpa opsi `--login` (`-l`), dan tidak mengirimkan berkas MOTD dinamis Armbian (`/run/motd.dynamic` / `/etc/motd`).
   - Akibatnya, bash tidak memuat profile login Armbian dan tidak mencetak logo/banner telemetri host seperti halnya sesi SSH di PuTTY.

#### 3. Solusi & Rekayasa Kode
1. **Backend PTY Terminal Handler (`RAZTerminalHandler.go`):**
   - Menjalankan shell bash dalam mode login (`/bin/bash -l`) di direktori home pengguna (`/root`).
   - Membaca dan langsung menyalurkan berkas MOTD resmi Armbian (`/run/motd.dynamic` atau `/etc/motd`) sesaat setelah sesi WebSocket tersambung sehingga logo ASCII merah Armbian OS dan ringkasan performa langsung tampil persis seperti PuTTY.
   - Menambahkan decoder pesan kontrol resize berformat JSON (`{"type":"resize","cols":...,"rows":...}`).
   - Memanggil `pty.Setsize(ptmx, &pty.Winsize{Rows: rows, Cols: cols})` saat menerima pesan resize, yang secara otomatis memicu sinyal `SIGWINCH` kernel ke bash. Bash kini otomatis melebarkan lebar kolom mengikuti ukuran jendela penuh tanpa wrap terpotong.
2. **Frontend Terminal Suite (`RAZTerminalView.vue`):**
   - Memasang `ResizeObserver` pada `terminalContainer` untuk mendeteksi perubahan ukuran jendela desktop RAZOS secara instan (resize, snap, maximize).
   - Mengirimkan paket resize `{ type: 'resize', cols, rows }` ke WebSocket backend setiap kali ukuran berubah.
   - Mengirimkan query parameter `cols` & `rows` awal saat pembukaan WebSocket URL agar PTY langsung terbuka sesuai dimensi jendela sejak milidetik pertama.
   - Mengaktifkan `convertEol: true` pada konfigurasi xterm untuk mencegah efek *stair-stepping* pada baris teks panjang.
3. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-Cp0feTnm.js`).

---

### Versi: 3.9.2-razos-terminal-motd-banner-and-realtime-telemetry (Custom RAZOS Terminal Banner & Real-Time Performance Telemetry)
**Waktu Eksekusi:** 2026-10-26 12:08 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna meminta agar tampilan banner selamat datang di Web Terminal disesuaikan menjadi banner resmi **RAZOS** (bukan logo Armbian standar) lengkap dengan data telemetri performa sistem riil (CPU Load, Suhu, RAM, ZRAM, Storage `/` & `/DATA`, Uptime, IP LAN, dan Container).

#### 2. Solusi & Rekayasa Kode
1. **Pembuatan Generator Banner RAZOS (`RAZTerminalBanner.go`):**
   - Merancang logo ASCII Art presisi untuk brand **RAZOS** dengan palet warna gradien Cyan/Sky-Blue (`\x1b[38;5;39m` ke `\x1b[38;5;51m`).
   - Mengumpulkan data telemetri riil dari kernel Linux host menggunakan library `gopsutil` dan sensor thermal `/sys/class/thermal/`:
     - **CPU Load & Cores**: Utilisasi CPU real-time dan jumlah core SOC Amlogic S905X.
     - **Uptime**: Format hari, jam, menit dari `/proc/uptime`.
     - **Memory Usage & ZRAM**: RAM fisik aktif dalam MiB/GiB dan persentase ZRAM Swap terpakai.
     - **CPU Temp**: Suhu thermal zone SOC dengan pewarnaan dinamis (Hijau < 65°C, Kuning < 78°C, Merah >= 78°C).
     - **Disk Storage**: Kapasitas riil partisi eMMC root `/` dan MicroSD `/DATA`.
     - **Active Containers**: Jumlah kontainer Docker aktif via Docker API negotiation.
     - **Network & IP**: Alamat IPv4 LAN aktif dan tautan Web Dashboard.
   - Mengembangkan fungsi perataan visual `padRight` dan `stripANSI` agar kolom telemetri sejajar rapi secara matematis terlepas dari kode warna ANSI.
2. **Penyambungan ke PTY WebSocket Handler (`RAZTerminalHandler.go`):**
   - Mengganti pencetakan banner statis Armbian dengan pemanggilan `GenerateRAZOSBanner()`.
   - Mengarahkan direktori awal shell interaktif ke home user (`/root`).

---

### Versi: 3.9.3-website-hosting-verification-and-port-binding (Website Hosting Container Port Re-Creation & Live Verification)
**Waktu Eksekusi:** 2026-10-26 12:26 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Temuan
- Website `webraz` yang telah dibuat sebelumnya pada port 8100 masih tidak dapat dijangkau (`TcpTestSucceeded: False` / Connection Refused).
- Hal ini terjadi karena kontainer `raz-web-webraz` dibuat sebelum perbaikan `ExposedPorts` & `PortBindings` diterapkan ke backend, dan Docker tidak dapat mengubah port binding pada kontainer yang sudah terbentuk hanya dengan restart biasa.

#### 2. Solusi & Verifikasi
1. **Trigger Self-Healing Redeploy:**
   - Melakukan re-deploy kontainer secara bersih melalui endpoint `POST /api/v1/websites/1/action` (`restart`).
   - Backend Go secara otomatis menghapus kontainer lama yang cacat port, mengunduh image `nginx:alpine`, lalu membangun kontainer baru dengan mapping port `0.0.0.0:8100 -> 80/tcp`.
2. **Verifikasi Status Live (Verified Result):**
   - `Test-NetConnection -ComputerName 192.168.100.117 -Port 8100`: Berhasil (`TcpTestSucceeded : True`).
   - `Invoke-WebRequest -Uri "http://192.168.100.117:8100"`: Berhasil mengembalikan `HTTP/1.1 200 OK` dengan landing page RAZ Creative Studio (`/DATA/www/webraz/index.html`).
   - Website kini 100% online dan dapat diakses dari browser jaringan lokal.

---

### Versi: 3.9.4-app-store-live-terminal-box-and-scanner (App Store Real-Time Terminal Execution Box & Docker Pipe Scanner)
**Waktu Eksekusi:** 2026-10-26 12:38 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Saat pengguna memasang aplikasi di App Store (misal Immich), progres di Tahap 5 tampak tertahan pada teks statis `$ docker pull ...` tanpa informasi detail apa yang sedang dikerjakan terminal di belakang layar (apakah mengunduh layer, mengekstrak, menyiapkan storage, atau membuat kontainer).
- Pengguna meminta agar dibuatkan teks boks terminal di jendela modal instalasi untuk menampilkan aktivitas terminal secara riil dan mengetahui proses instalasi sudah sampai mana.

#### 2. Solusi & Rekayasa Kode
1. **Pembaruan Backend Store Engine (`RAZAppStoreEngine.go`):**
   - Mengembangkan sistem logging instalasi thread-safe (`LogItem`, `AddInstallLog`, `GetInstallLogs`, `ClearInstallLogs`).
   - Mengintegrasikan pembacaan output stream riil dari perintah `docker pull` menggunakan `cmd.StdoutPipe()` & `cmd.StderrPipe()` via `bufio.Scanner` sehingga setiap layer download (`Pulling fs layer`, `Downloading [===> ]`, `Extracting`, `Pull complete`) langsung tercatat baris demi baris.
   - Mencatat setiap fase instalasi: Validasi RAM STB (`[SYS]`), Inisialisasi Storage `/DATA/apps/<slug>` (`[STORAGE]`), Pemetaan Port Host (`[NET]`), Eksekusi Docker Pull (`[PULL]`), Pembuatan Kontainer (`[CREATE]`), dan Container Start (`[ONLINE]`).
2. **Pembaruan REST API Controller (`RAZApiController.go`):**
   - Menambahkan endpoint REST `GET /api/v1/store/install/logs/:slug` (`handleGetStoreInstallLogs`) agar UI frontend dapat menarik log aktivitas terminal secara real-time.
3. **Pembaruan UI Terminal Box (`RAZAppStoreInstallView.vue`):**
   - Mengganti teks statis 1 baris dengan komponen **Live Terminal Output Console Box**:
     - Header jendela terminal dengan 3 dot macOS/Linux style, teks status `installer@razos:~# docker pull`, indikator pulsing `LIVE`, counter jumlah baris log, dan tombol cepat "Salin" log ke clipboard.
     - Wadah log monospaced hitam transparan (`bg-slate-950 border border-slate-800/90`) dengan autoscroll otomatis ke baris paling baru (`nextTick` scrollTop).
     - Pewarnaan tag dinamis: Sky-blue untuk `[INFO]`/`[SYS]`, Amber untuk `[DOCKER]`, Indigo untuk `[PULL]`, Emerald untuk `[ONLINE]`/`[SUCCESS]`, dan Rose untuk `[ERR]`.
     - Kursor prompt berkedip (`animate-pulse`) yang menandakan proses terminal aktif.
     - Menampilkan terminal log riil pada state error bila proses terhenti agar penyebab kegagalan langsung transparan.
4. **Penyelarasan Dimensi Jendela (`RAZWindowStore.js`):**
   - Menyesuaikan ukuran default jendela instalasi menjadi `width: 560, height: 520` (min: `420x380`) agar checklist dan boks terminal muat proporsional tanpa desak-desakan.
5. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-9Nk6Q-7G.js` & `dist/assets/index-8as4WQE_.css`).

---

### Versi: 4.0.0-razos-native-devapi-mcp-bridge (RAZOS Native DevAPI MCP Bridge & AI Hot-Reload Automation)
**Waktu Eksekusi:** 2026-10-26 12:55 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna meminta agar dibuatkan MCP Server yang menghubungkan langsung AI Antigravity ke Armbian STB RAZOS, sehingga instruksi perbaikan dapat langsung dieksekusi di STB tanpa pengguna harus membuka PuTTY secara manual, sambil memastikan repository lokal laptop dan GitHub tetap tersinkronisasi 100%.

#### 2. Solusi & Rekayasa Arsitektur
1. **Pembangunan Modul Backend Native DevBridge (`backend/internal/dev/RAZDevBridge.go`):**
   - Modul internal Go terlindungi token rahasia `X-RAZOS-Dev-Token` (`razos-super-secret-dev-bridge-2026`).
   - Menyediakan 5 endpoint utama pengembang:
     - `GET /api/v1/dev/status`: Telemetri real-time STB, commit git aktif, branch, dan sisa RAM.
     - `POST /api/v1/dev/exec`: Eksekusi perintah bash di STB dengan output stdout/stderr & exit code.
     - `POST /api/v1/dev/hot-reload`: Otomatisasi 3 langkah (Git fetch/reset, Go build, dan restart systemd service secara asinkron).
     - `POST /api/v1/dev/sync-file`: Penulisan langsung berkas kode dari laptop ke `/DATA/razos/src/` tanpa harus menunggu git push.
     - `GET /api/v1/dev/logs`: Pembacaan log live `journalctl` atau `docker logs`.
2. **Integrasi Router API (`backend/internal/api/RAZApiController.go`):**
   - Mendaftarkan endpoint group `/api/v1/dev` yang otomatis aktif saat backend RAZOS dijalankan.
3. **Pembangunan MCP Server Antigravity (`RAZDevBridgeMCP.mjs`):**
   - Ditempatkan pada `C:\Users\Administrator\.gemini\antigravity\mcp\razos-bridge\RAZDevBridgeMCP.mjs`.
   - Menggunakan Node.js native ESM tanpa ketergantungan pihak ketiga (zero external npm dependencies).
   - Menyediakan 5 tools MCP: `razos_status`, `razos_hot_reload`, `razos_exec`, `razos_sync_file`, dan `razos_logs`.
   - Mendaftarkan server ke `mcp_config.json` Antigravity secara permanen.

---

### Versi: 4.1.0-app-store-uninstaller-window-with-live-terminal (Standalone Uninstaller Window & Real-Time Terminal Log Execution)
**Waktu Eksekusi:** 2026-10-26 13:28 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna meminta agar ketika tombol "Uninstall Container" diklik, sistem tidak hanya sekadar menutup dialog atau memutar spinner tak tentu, melainkan membuka jendela baru seperti pada proses instalasi yang memiliki teks boks terminal untuk menampilkan aktivitas dan tahapan proses uninstall secara transparan (stop container, remove container & volume, storage purge, dsb).

#### 2. Solusi & Rekayasa Kode
1. **Pembaruan Backend Store Engine (`RAZAppStoreEngine.go`):**
   - Memutakhirkan `CleanUninstall` agar menyalurkan log aktivitas instalasi/uninstalasi secara terstruktur ke `AddInstallLog`:
     - `[STOP]`: docker stop -t 5 razos-<slug>
     - `[REMOVE]`: docker rm -v -f razos-<slug>
     - `[STORAGE]`: rm -rf /DATA/apps/<slug> (jika opsi purge data diaktifkan)
     - `[DONE]`: Pelepasan memori dan resource host selesai 100%.
2. **Pembuatan Komponen Mandiri (`RAZAppStoreUninstallView.vue`):**
   - Jendela pelacak proses uninstall interaktif:
     - Header jendela dengan ikon aplikasi, progres persentase, dan badge tahapan aktif.
     - Progress bar dengan gradien Rose ke Red (`from-rose-500 via-red-500 to-amber-400`).
     - Checklist 5 tahap: Validasi Container, Menghentikan Container, Menghapus Container & Volume, Pembersihan Direktori Storage, dan Pelepasan Resource Host.
     - **Live Terminal Execution Box**: Wadah terminal monospaced dengan autoscroll, indikator `LIVE UNINSTALL` berkedip, counter baris log, dan tombol salin log ke clipboard.
     - State sukses & gagal yang informatif dengan tombol penutup yang menyinkronkan daftar aplikasi terpasang di desktop.
3. **Penyelarasan Store Jendela (`RAZWindowStore.js`):**
   - Menambahkan fungsi `openAppStoreUninstaller(app, purgeData)` dengan dimensi proporsional `560x520 px`.
4. **Integrasi UI Desktop & App Store (`RAZDesktopLayout.vue` & `RAZAppStoreView.vue`):**
   - Mengalihkan aksi konfirmasi tombol "Uninstall Container" untuk memanggil `windowStore.openAppStoreUninstaller` sehingga proses uninstall langsung divisualisasikan di jendela terminal mandiri.
5. **Kompilasi Ulang Frontend:**
   - Melakukan build produksi Vite (`dist/assets/index-zAuGkjJo.js` & `dist/assets/index-C_VRwv9X.css`).

---

### Versi: 4.1.1-app-store-edit-port-reactivity-and-conflict-fix (Perbaikan Port Reaktif Pengaturan Container & Pencegahan Bentrok Port 8080)
**Waktu Eksekusi:** 2026-10-26 13:38 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna melaporkan bahwa saat membuka pengaturan container (seperti Immich) melalui jendela edit aplikasi (`RAZAppEditView.vue`), nilai port yang muncul pada form host dan container berubah menjadi `8080 : 8080` padahal di App Store katalog Immich memiliki port `2283`, sehingga aplikasi tidak dapat dibuka dan memicu kebingungan konfigurasi.

#### 2. Investigasi & Akar Masalah (Root Cause Analysis)
1. **Fallback Hardcode 8080 Tanpa Watcher Reaktif:**
   - Pada `frontend/src/views/RAZAppEditView.vue`, `populateForm()` hanya dijalankan satu kali saat `onMounted`. Karena data props jendela (`props.win?.customData`) bisa dimuat secara asinkron atau membawa struktur pembungkus `storeData`, `app.value.port` bernilai `undefined`, yang kemudian jatuh ke default fallback statis `8080`.
2. **Kehilangan Properti Port pada Objek Desktop Icon:**
   - Pada `frontend/src/stores/RAZWindowStore.js`, saat menyusun `allDesktopApps`, objek instansi hanya menyalin sebagian field dan `port` diisi dari `inst.port`. Backend endpoint `/api/v1/store/installed` sebelumnya mengembalikan `default_port` dan bukan `port`, sehingga nilai `port` menjadi kosong dan jatuh ke `undefined`.
3. **Pemberian Fallback Port pada Peluncur Aplikasi (`launchStoreApp`):**
   - Fungsi `launchStoreApp(app)` di `RAZWindowStore.js` langsung mengarahkan ke port `8080` bila `app.port` kosong, tanpa mengecek `app.default_port` atau katalog spesifik aplikasi (seperti Immich port `2283`).

#### 3. Solusi & Rekayasa Kode
1. **Pembaruan Backend Store Engine (`backend/internal/store/RAZAppStoreEngine.go`):**
   - Menambahkan field `Port int` pada struct `AppManifest` (`json:"port,omitempty"`).
   - Memutakhirkan `GetCatalog()` agar saat mendeteksi container terpasang (`/razos-<slug>`), Docker inspect membaca `PublicPort` riil dari daftar port container yang aktif. Jika belum terpetakan, mengisi otomatis dengan `DefaultPort` dari manifes.
   - Menambahkan fallback `req.CustomImage = req.Image` pada `DeployApp()` untuk memastikan kestabilan saat menerima payload dari form edit.
2. **Pembaruan Jendela Pengaturan Container (`frontend/src/views/RAZAppEditView.vue`):**
   - Mengimpor `watch` dari Vue dan menambahkan watcher mendalam (`immediate: true, deep: true`) pada `props.win?.customData` agar form langsung terisi otomatis saat data aplikasi masuk.
   - Mengganti fallback hardcode `8080` menjadi port dinamis sesuai katalog container (`a.default_port || (a.slug === 'immich' ? 2283 : 8090)`).
   - Menambahkan unpacking pintar pada computed `app` (`d.storeData ? { ...d.storeData, ...d } : d`).
   - Menambahkan validasi keamanan port: Menolak penyimpanan jika pengguna memasukkan port `80` atau `8080` (port antarmuka inti RAZOS) dengan modal notifikasi peringatan.
3. **Penyelarasan Store Jendela (`frontend/src/stores/RAZWindowStore.js`):**
   - Pada `allDesktopApps`, meratakan seluruh properti `inst` (`...inst`) serta menyertakan `port: livePort` dan `default_port: inst.default_port`.
   - Pada `openAppEditor(app)`, melakukan unpacking `storeData` sebelum membuka jendela kustom.
   - Pada `launchStoreApp(app)`, menyertakan resolusi port dinamis dengan fallback `default_port` dan spesifik slug Immich (`2283`).
4. **Kompilasi Ulang Frontend:**
   - Membangun bundle produksi Vite terbaru (`dist/assets/index-C2_ptWDS.js` & `dist/assets/index-C_VRwv9X.css`).

---

### Versi: 4.1.2-immich-companion-microservices-stack (Deployment Service Pendamping PostgreSQL 16 pgvector & Redis untuk Immich)
**Waktu Eksekusi:** 2026-10-26 14:15 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Immich tidak dapat diakses setelah instalasi karena container Immich resmi (`immich-server:release`) mewajibkan adanya service pendamping PostgreSQL (dengan ekstensi `pgvector`) dan Redis. Tanpa kedua service tersebut, Node.js di dalam container mengalami `getaddrinfo ENOTFOUND database` dan mati dalam hitungan detik (crash-loop).

#### 2. Investigasi & Akar Masalah
- Immich v3.x membutuhkan ekstensi `vector` (pgvector) pada database PostgreSQL untuk pencarian embedding foto dan metadata.
- Komponen worker microservices Immich juga memerlukan alokasi heap Node.js yang memadai agar tidak terhenti oleh `ERR_WORKER_OUT_OF_MEMORY`.

#### 3. Solusi & Rekayasa Sistem di STB
1. **Penyediaan Network Terisolasi (`razos-apps-net`):**
   - Membuat docker network bridge khusus agar antar container dapat saling berkomunikasi via DNS alias (`database` dan `redis`).
2. **Pemasangan Service Redis (`razos-redis`):**
   - Menggunakan image ringan `redis:alpine` dengan limit memori hemat STB (64 MB) dan alias `redis`.
3. **Pemasangan Service Database PostgreSQL 16 dengan pgvector (`razos-postgres`):**
   - Menggunakan image resmi `pgvector/pgvector:pg16` arsitektur ARM64 dengan limit memori 256 MB dan volume persisten di `/DATA/apps/immich/database`.
   - Mengaktifkan ekstensi `vector` v0.8.7 pada database `immich`.
4. **Rekonfigurasi & Booting Immich Server (`razos-immich`):**
   - Menghubungkan container ke `razos-apps-net` dengan port binding `2283:2283`.
   - Mengonfigurasi variabel lingkungan:
     - `DB_HOSTNAME=database`
     - `DB_USERNAME=postgres`
     - `DB_PASSWORD=postgres`
     - `DB_DATABASE_NAME=immich`
     - `REDIS_HOSTNAME=redis`
     - `IMMICH_MACHINE_LEARNING_ENABLED=false` (menghemat CPU/RAM STB)
     - `NODE_OPTIONS="--max-old-space-size=512"`
   - Volume upload foto diarahkan persisten ke `/DATA/apps/immich/upload`.
5. **Verifikasi Operasional:**
   - Container status: `Up (healthy)`.
   - HTTP Health Check: `HTTP/1.1 200 OK` pada `http://127.0.0.1:2283`.
   - API Server Version: `{"major":3,"minor":3,"patch":5,"prerelease":null}`.
   - Sisa RAM Bebas STB: `1.40 GB` (sangat lega).

---

### Versi: 4.1.3-remote-mcp-access-via-cloudflare-tunnel (Dukungan Akses Global Remote MCP via Cloudflare Tunnel)
**Waktu Eksekusi:** 2026-10-26 14:22 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna bertanya apakah saat berada di luar rumah (tidak satu jaringan Wi-Fi lokal dengan STB), pembaruan sistem atau perbaikan kerusakan masih bisa dilakukan melalui MCP Antigravity. Pengguna menyediakan domain aktif `server.razos.my.id`.

#### 2. Investigasi & Verifikasi
- Domain `server.razos.my.id` sudah terhubung aktif ke STB RAZOS melalui Cloudflare Zero Trust Tunnel (`cloudflared`).
- Uji endpoint: `curl -H 'X-RAZOS-Dev-Token: ...' https://server.razos.my.id/api/v1/dev/status` merespons sukses `HTTP 200 OK` dengan payload telemetri real-time.

#### 3. Rekayasa MCP Server (`RAZDevBridgeMCP.mjs`)
- Memutakhirkan `RAZDevBridgeMCP.mjs` dengan mekanisme **Smart Dual-Endpoint Failover**:
  1. **Primary Endpoint:** `https://server.razos.my.id/api/v1/dev` (akses HTTPS aman dari mana saja di seluruh dunia).
  2. **Fallback Endpoint:** `http://192.168.100.117:8080/api/v1/dev` & `http://192.168.100.117/api/v1/dev` (akses LAN lokal super cepat saat satu Wi-Fi).
- Jika laptop berada di luar rumah (kafe, kantor, tethering HP), MCP otomatis menggunakan jalur Cloudflare Tunnel tanpa perlu konfigurasi manual.
- Menyimpan cache endpoint aktif yang berhasil merespons agar request berikutnya tetap cepat.

---

### Versi: 4.1.4-task-manager-fix-and-file-manager-windows-suite (Perbaikan Jendela Kosong Task Manager & Fitur Lengkap File Manager: Multi-Selection, Dedicated Context Menu, Custom Delete Modal)
**Waktu Eksekusi:** 2026-10-26 20:50 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
1. **Task Manager Blank Window:** Jendela Task Manager terbuka namun layarnya kosong/hitam total tanpa menampilkan tab, metrik, maupun daftar proses.
2. **Context Menu Khusus File Manager:** Membutuhkan context menu kontekstual khusus untuk berkas tunggal, multi-berkas (batch actions), dan kanvas area kosong.
3. **Multi-Selection File ala Windows Explorer:** Pemilihan banyak berkas sekaligus menggunakan tombol `Ctrl` (toggle pilihan) dan `Shift` (pemilihan rentang berkas / range selection), serta pintasan keyboard `Ctrl+A` (Pilih Semua) dan `Escape` (Batalkan Seleksi).
4. **Modal Konfirmasi Hapus Berkas Kustom:** Modal konfirmasi sebelum menghapus berkas/folder dengan desain Glassmorphism Slate-900 terintegrasi tema RAZOS tanpa dialog `alert()` / `confirm()` bawaan browser.

#### 2. Investigasi & Akar Masalah (Root Cause Analysis)
1. **Penyebab Blank Window Task Manager (`RAZTaskManagerView.vue`):**
   - Pada baris template daftar proses, komponen memanggil `proc.memoryMb.toFixed(1)`. Sementara endpoint Go backend `/api/v1/processes` mengembalikan field JSON `memory` (float64 percent), bukan `memoryMb`. Hal ini memicu `TypeError: Cannot read properties of undefined (reading 'toFixed')` di dalam render loop Vue 3 yang mengakibatkan seluruh komponen gagal me-render dan menjadi layar kosong.
   - Panggilan langsung `.toFixed()` pada `metrics.cpu_percent`, `metrics.ram_percent`, dan `metrics.cpu_temp` rentan melempar error saat metrik belum selesai dimuat atau bernilai null.
2. **Keterbatasan Seleksi & Hapus File Manager (`RAZFileManagerView.vue`):**
   - Sebelumnya `selectedItem` hanya menyimpan objek tunggal sehingga tidak mendukung multi-selection.
   - Fungsi `deleteItem` sebelumnya langsung mengirim request `DELETE` ke API backend tanpa menampilkan konfirmasi peringatan kepada pengguna.

#### 3. Solusi & Rekayasa Kode
1. **Pembaruan Task Manager (`frontend/src/views/RAZTaskManagerView.vue`):**
   - Mengembangkan helper `safeFixed(val, digits, fallback)` yang kebal crash terhadap nilai `undefined`, `null`, maupun `string` tidak valid.
   - Memutakhirkan `filteredProcesses` agar memetakan field `cpu`, `memory`, dan `memoryMb` dengan fallback aman.
   - Mengganti seluruh pemanggilan `.toFixed` di seluruh tab (Processes, Performance, Containers, Services) dengan `safeFixed`.
2. **Pembaruan File Manager Multi-Selection (`frontend/src/views/RAZFileManagerView.vue`):**
   - Menambahkan state `selectedItems` (array objek berkas) dan `lastSelectedIndex` (anchor index untuk rentang Shift-Click).
   - Mengembangkan logika `handleItemClick(item, event)`:
     - **Ctrl / Cmd + Click**: Toggle status pilihan item tanpa membatalkan pilihan item lainnya.
     - **Shift + Click**: Pemilihan rentang berkas (Range Selection) dari berkas yang terakhir dipilih hingga berkas saat ini.
     - **Klik Biasa**: Memilih satu berkas saja dan mereset pilihan lainnya.
   - Menambahkan fungsi `selectAllFiles()` (`Ctrl+A`) dan `clearSelection()` (`Escape`).
   - Menambahkan event listener keyboard global `@keydown="handleGlobalKeydown"` pada root explorer.
3. **Dedicated Context Menu File Manager:**
   - **Mode Multi-Berkas (>1 item terpilih):** Menampilkan menu aksi massal: `Buka X Berkas`, `Salin Seluruh Path (X Item)`, `Beri Bintang Semua`, dan `Hapus X Berkas Terpilih...` (Danger).
   - **Mode Berkas Tunggal:** Menampilkan menu lengkap: Buka, Edit Berkas (Editor Teks), Ubah Nama, Unduh Berkas, Salin Path, Bintang, Info & Properti, serta Hapus.
   - **Mode Area Kanvas Kosong:** Menampilkan menu: Folder Baru, File Baru (Teks/Kode), Menu Upload (File, Folder, ZIP), Pilih Semua (Ctrl+A), Bersihkan Pilihan, Segarkan, dan Salin Path Direktori.
4. **Bottom Status Bar Dinamis:**
   - Menampilkan badge jumlah berkas terpilih (`X berkas terpilih`) beserta tombol batal seleksi, shortcut Pilih Semua, dan tombol cepat `Hapus (X Berkas)`.
5. **Custom Modal Konfirmasi Hapus Berkas (`showDeleteConfirmModal`):**
   - Didesain mengikuti standar UI Glassmorphism Slate-900 dengan ring merah `border-rose-500/40`.
   - Menampilkan daftar nama berkas/folder yang akan dihapus beserta ukuran dan tipe ikonnya.
   - Memiliki tombol "Batal" dan "Ya, Hapus Permanen" dengan indikator loading status saat proses penghapusan massal berjalan.
   - 100% patuh aturan GEMINI.md: Zero browser native alerts/confirms.
6. **Kompilasi Ulang Frontend:**
   - Menjalankan Vite build produksi (`dist/assets/index-BmNHPoeP.js` & `dist/assets/index-D-vNd4Fk.css`) terverifikasi sukses tanpa error.

---

### Versi: 4.2.0-windows-11-settings-suite-overhaul (Transformasi Total Tampilan Settings ala Windows 11 Fluent Design)
**Waktu Eksekusi:** 2026-10-26 21:00 WIB  
**Pelaksana:** Antigravity AI Senior Systems Engineer  

#### 1. Masalah / Kebutuhan Pengguna
- Pengguna meminta agar tampilan Settings dirombak total menyerupai antarmuka resmi **Windows 11 Settings** persis seperti screenshot referensi pengguna:
  - Layout dua kolom: Left Navigation Rail (sidebar) dengan profil user, input pencarian "Find a setting", dan daftar kategori lengkap dengan indikator aktif biru khas Windows 11.
  - Halaman `System > About` dengan Hero section (nama PC & tombol "Rename this PC") dan Accordion Collapsible Cards:
    - **Device specifications:** Ikon bulat `(i)`, tombol "Copy", dan tabel spesifikasi hardware host riil (Device name, Processor CPU & Cores, Installed RAM usable, Device ID, Product ID, System Type 64-bit ARM/x64).
    - **Windows specifications / RAZOS specifications:** Ikon 4-kotak logo, tombol "Copy", dan tabel detail OS (Edition, Version, Installed on, OS build, Kernel host, Experience pack).
  - Kategori pengaturan fungsional lainnya: Personalization (pilihan wallpaper RAZOS resmi, tema, widget), Network & Internet (IP LAN, Cloudflare Tunnel, port WebUI), Accounts (profil Administrator, ganti password), Apps, Server Performance, Privacy & Security, dan Windows/RAZOS Update.

#### 2. Solusi & Rekayasa Kode
1. **Pembaruan Dimensi Jendela (`frontend/src/stores/RAZWindowStore.js`):**
   - Menyesuaikan dimensi default jendela `settings` menjadi `width: 860, height: 560, minWidth: 480, minHeight: 360` agar tata letak dua kolom Windows 11 tampil proporsional tanpa terpotong.
2. **Rekayasa Ulang Komponen Settings (`frontend/src/views/RAZSettingsView.vue`):**
   - Mengembangkan Left Navigation Rail dengan avatar profil Administrator, kotak pencarian "Find a setting" reaktif, dan 11 kategori menu Windows 11 dengan pill bar biru vertikal aktif.
   - Mengembangkan tampilan `System > About` dengan styling `#1c1b22` / `#201f27` / `#262530` persis Windows 11 Dark Mode.
   - Mengintegrasikan Accordion Cards (`Device specifications` & `RAZOS specifications`) dengan tombol `Copy` fungsional yang menyalin detail spesifikasi ke clipboard dalam format tabulasi rapi.
   - Mengintegrasikan tombol `Rename this PC` dengan custom modal dialog.
   - Menghubungkan seluruh data hardware riil dari `systemStore.metrics` (CPU model, cores, RAM total/used, OS platform).
   - Mempertahankan seluruh fitur kontrol sistem penting: Ubah password, ganti port WebUI, galeri wallpaper resmi RAZOS, pembersihan RAM 1-klik (`drop_caches`), restart/shutdown STB, dan OTA update.
3. **Kompilasi Ulang Frontend:**
   - Menjalankan build produksi Vite (`dist/assets/index-Bjj8GvcN.js` & `dist/assets/index-pQlwp5ri.css`) sukses tanpa error.






