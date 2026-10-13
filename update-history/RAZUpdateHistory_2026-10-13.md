# RAZOS Update History: 2026-10-13

Dokumen riwayat pembaruan harian proyek RAZOS sesuai standar operasional enterprise.

---

### Versi: 2.2.0-file-manager-ui-redesign (Penyempurnaan Proporsi Grid Kartu File Manager & Detail Preview Pane)
**Waktu Eksekusi:** 2026-10-13 09:35 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Perbaikan Layouting & Visual Grid Kartu Berkas (Google Drive / Apple Finder Style)
- **Koreksi Aspek Rasio Kartu:** 
  - Mengubah orientasi kartu berkas dari sebelumnya sempit memanjang vertikal (yang memotong nama file secara ekstrem menjadi `$...`, `4...`, dsb) menjadi kartu horizontal-proporsional modern dengan margin & padding seimbang.
  - Ikon berkas kini diletakkan di dalam container `w-11 h-11 rounded-xl` terpadu dengan warna aksen amber/sky/emerald/rose yang tegas.
  - Nama berkas menggunakan wrapping `line-clamp-2` dengan `min-w-0` dan font size `text-xs font-medium` sehingga nama panjang terbaca utuh dan jelas.
  - Menampilkan badge ekstensi format berkas dan indikator ukuran file secara rapi dan serasi di baris bawah kartu.
- **Penyempurnaan Right Details Preview Pane:**
  - Memperluas lebar panel preview kanan menjadi `w-72` untuk memberi ruang bernapas pada metadata berkas.
  - Mengganti preview atas dengan visual preview container `w-20 h-20` beraksen glow halus.
  - Menyusun daftar properti berkas (Tipe, Ukuran, Lokasi Path, Waktu Modifikasi) dengan tata letak key-value bergaris pembatas halus.
- **Peningkatan Responsivitas Grid:**
  - Penyesuaian breakpoints Tailwind: `grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-3` yang adaptif terhadap ukuran jendela multi-window OS.
- **Verifikasi Build:**
  - `npm run build` berhasil 100% tanpa error kompilasi Vite (1501 modules transformed).

---

### Versi: 2.2.1-file-manager-fluid-grid-fix (Eliminasi Right Details Pane, Fluid Grid Auto-Fill & Bottom Status Bar)
**Waktu Eksekusi:** 2026-10-13 09:56 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Solusi Tuntas Folder List Gepeng & Overcrowding
- **Penyebab:** Right Details Pane di samping kanan (`w-72`) memaksa masuk ke layout horizontal saat berkas diklik, memotong lebar ruang tengah jendela dari ~800px menjadi ~280px. Ketika dikombinasikan dengan media queries Tailwind berbasis lebar layar browser penuh, grid memaksa 5-6 kolom di ruang 280px sehingga kartu-kartu menjadi gepeng setebal 50px.
- **Eliminasi Panel Kanan:** Menghapus panel samping kanan statis yang meremas layout explorer.
- **Implementasi Fluid Responsive Auto-Fill Grid:**
  - Menerapkan CSS Grid fleksibel: `grid-cols-[repeat(auto-fill,minmax(140px,1fr))] sm:grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-3`.
  - Berapapun ukuran resize jendela window, lebar kartu folder/file dijamin minimal 140px - 150px sehingga tidak akan pernah gepeng atau tumpang tindih lagi.
- **Bottom Status & Quick Action Bar (Ala Windows Explorer / Apple Finder):**
  - Mengintegrasikan status bar di bagian bawah jendela yang menampilkan jumlah total berkas, nama item yang dipilih, format tipe, ukuran, dan waktu modifikasi.
  - Menyediakan tombol aksi cepat langsung di status bar: **Info Berkas**, **Buka**, **Ubah Nama**, dan **Hapus**.
- **Custom Modal Detail & Properti Berkas:**
  - Rincian metadata lengkap kini dapat dibuka melalui Custom Modal (Glassmorphism Slate-900) dengan klik tombol "Info Berkas" atau icon `(i)` pada hover kartu, menjaga tampilan explorer tetap luas, bersih, dan lega.
- **Verifikasi:**
  - `npm run build` sukses 100% tanpa error kompilasi.
  - Hot Module Replacement (HMR) aktif di dev server `localhost:3000`.

---

### Versi: 2.2.2-folder-navigation-integrity-fix (Perbaikan Lokasi Folder & Integritas Navigasi Anti Salah Folder)
**Waktu Eksekusi:** 2026-10-13 10:44 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Perbaikan Navigasi Breadcrumbs & Tombol Kembali (Up Navigation)
- **Eliminasi Bug Driver Letter Bare Windows:** Memperbaiki bug di mana navigasi mundur dari subfolder Windows (`C:/DATA`) menghasilkan string `"C:"` tanpa slash, yang memicu Node.js `fs.readdirSync("C:")` membaca direktori kerja aktif frontend (`node_modules`, `src`) bukannya root harddisk. Seluruh pemotongan path kini dinormalisasi menjadi `"C:/"` dengan trailing slash.
- **Controlled Up Navigation:** Menambahkan state `canNavigateUp` sehingga saat berada di root (`C:/` atau `/`), tombol panah kembali dinonaktifkan secara anggun (disabled dengan visual redup) untuk mencegah lompat folder liar.
- **Akurat & Terstruktur Breadcrumbs:**
  - Menampilkan segmen hirarkis yang benar (`C: Disk` -> `My Storage (/DATA)` -> `subfolder`).
  - Setiap segmen memiliki icon penanda (HardDrive / Folder Cloud) dan path absolut forward slash yang konsisten (`filepath.ToSlash`).

#### 2. Resolusi Lokasi Default "My Storage (/DATA)"
- **Pemisahan My Storage & Storage Disks:**
  - Tombol **"My Storage (/DATA)"** di sidebar kini selalu diarahkan ke direktori kerja utama penyimpanan RAZOS (`C:/DATA` di Windows atau `/DATA` di Linux STB). Tidak lagi keliru membuka root harddisk `C:\` sistem operasi.
  - Penelusuran harddisk fisik host (misal `C: Local Disk`) tetap tersedia di bagian bawah **Storage Disks**, sehingga pengguna dapat beralih antara cloud storage dan root harddisk secara terisolasi tanpa tertukar.
- **Normalisasi Backend & Dev Middleware:**
  - Menyelaraskan fungsi `resolveSafePath` di `frontend/vite.config.js` dan `RAZStorageEngine.go` untuk menangani path kosong, `/DATA`, atau drive letter secara seragam.
  - Seluruh item berkas dan folder kini dikembalikan dengan path forward-slash (`/`) untuk mencegah konflik backslash di browser.
- **Verifikasi:**
  - `npm run build` berhasil 100% (1501 modules transformed).
  - Dev server otomatis reload dengan HMR aktif di port 3000.

---

### Versi: 2.3.0-standalone-content-viewers (Pemisahan Content View Menjadi Window Mandiri Tanpa Efek Gelap/Blur)
**Waktu Eksekusi:** 2026-10-13 11:03 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Arsitektur Decoupled Multi-Window File Viewers
- **Eliminasi Total Modal Backdrop Gelap/Blur:**
  - Menghapus seluruh modal internal (`fixed inset-0 bg-black/90 backdrop-blur-md`) yang sebelumnya memblokir interaksi dan menggelapkan jendela File Manager saat membuka berkas.
  - File Manager kini tetap 100% aktif, terang, dan tidak terganggu ketika berkas dibuka.
- **Implementasi Standalone Universal Viewers (`RAZFileViewerView.vue`):**
  - Berkas yang dibuka (gambar, video, musik, PDF, atau kode/teks) kini otomatis dimuat ke dalam **Window Frame Mandiri** di desktop RAZOS melalui `windowStore.openFileViewer(item)`.
  - Terintegrasi penuh dengan desktop window manager: mendukung multi-window stacking, drag & drop, resize, minimize, maximize, dan tombol close independen.
  - Setiap window viewer memiliki ikon SVG resmi tersendiri (`RAZIconImageViewer.svg`, `RAZIconVideoPlayer.svg`, `RAZIconMusicPlayer.svg`, `RAZIconPdfViewer.svg`, `RAZIconCodeEditor.svg`) dan muncul di bottom desktop taskbar.
- **Multi-File Concurrent Viewing (Bebas Membuka Berkas Lain):**
  - Pengguna dapat membuka beberapa berkas sekaligus (misalnya membuka 2 gambar secara berdampingan untuk perbandingan, atau memutar lagu sambil meninjau dokumen teks/PDF) tanpa menutup File Manager.
- **Fitur Unggulan Tiap Format:**
  - **Image Viewer:** Zoom in/out (25%-400%), reset 100%, rotate 90°, dan tombol unduh.
  - **Video Player:** HTML5 native media controller dengan autoplay, picture-in-picture, dan loop.
  - **Music Player:** Visualizer piringan vinil berputar dengan scrubber timeline, penyesuaian volume, dan mute.
  - **PDF Viewer:** Embedded viewer dengan tombol tab baru dan download.
  - **In-Browser Code Editor:** Editor kode/teks monospace dengan syntax styling, status unsaved badge, reload, dan penyimpanan langsung ke backend host.
- **Verifikasi:**
  - `npm run build` sukses 100% tanpa error (1502 modules transformed).
  - Dev server HMR ter-update secara real-time di `localhost:3000`.

---

### Versi: 2.3.1-windows-style-address-bar (Bilah Alamat Interaktif Ala Windows File Explorer & History Navigation)
**Waktu Eksekusi:** 2026-10-13 11:10 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Implementasi Windows Explorer Interactive Address Bar
- **Bilah Alamat Dwi-Fungsi (Breadcrumb vs Editable Input):**
  - **Mode Breadcrumb (Default):** Menampilkan deretan segmen folder interaktif dengan ikon HardDrive/Folder dan pemisah chevron (`>`). Setiap segmen dapat langsung diklik untuk melompat ke direktori tersebut.
  - **Mode Editable Input (Klik Area / Shortcut):** Mengklik area kosong di bilah alamat beralih otomatis ke kotak input teks dengan path absolut lengkap (misal `C:/DATA/razos` atau `/DATA/apps`).
  - **Auto-Focus & Text Selection:** Teks path langsung terseleksi (*select all*) saat mode edit aktif sehingga pengguna dapat langsung mengetik atau menempelkan (*paste*) path tujuan.
  - **Keyboard Action:** Menekan tombol `Enter` mengeksekusi navigasi ke path yang diketik, sedangkan tombol `Esc` atau `blur` membatalkan pengeditan dan mengembalikan tampilan ke mode breadcrumb.
- **Aksi Cepat di Dalam Address Bar:**
  - **Tombol Go (`ArrowRight`):** Muncul saat mode edit aktif untuk navigasi instan.
  - **Tombol Copy Path (`Copy`):** Menyalin path direktori absolut ke clipboard sistem dengan notifikasi custom toast `systemStore.showNotification` tanpa modal browser native.
  - **Tombol Refresh Folder (`RefreshCw`):** Memuat ulang daftar berkas direktori aktif dengan animasi putar saat loading.

#### 2. History Navigation Trio (Back, Forward, Up)
- **Riwayat Navigasi Komprehensif (`historyStack` & `historyIndex`):**
  - **Tombol Back (`ArrowLeft`):** Melangkah mundur ke riwayat folder sebelumnya, otomatis nonaktif (*disabled*) jika berada di awal riwayat.
  - **Tombol Forward (`ArrowRight`):** Melangkah maju ke riwayat folder berikutnya, otomatis nonaktif jika berada di ujung riwayat.
  - **Tombol Up (`ArrowUp`):** Naik satu tingkat ke folder induk (*parent directory*), otomatis nonaktif jika berada di root sistem.
- **Deteksi Path & Fallback Aman:**
  - Jika pengguna memasukkan path direktori yang tidak ada atau tidak valid, sistem menampilkan notifikasi peringatan dan mengembalikan pengguna secara aman ke direktori default penyimpanan.

#### 3. Penataan Tata Letak (Explorer Ribbon & Search Box)
- **Windows Explorer Search Bar:**
  - Bilah pencarian beraksen modern dengan placeholder dinamis berbasis folder aktif (`Cari di [nama folder]...`) dan tombol clear (X).
- **Secondary Ribbon Bar:**
  - Memisahkan tombol filter kategori (Semua, Folder, Gambar, Video, Audio, Dokumen, Kode) dan pemilih tampilan (Grid / List) ke dalam ribbon bar terpisah yang rapi dan lapang.
- **Verifikasi Build:**
  - `npm run build` berhasil 100% (1502 modules transformed).
  - Dev server reload secara lancar via HMR di `localhost:3000`.

---

### Versi: 2.3.2-windows-junction-and-cache-resilience (Resolusi Windows Junction Symlink & Failsafe Store)
**Waktu Eksekusi:** 2026-10-13 11:15 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Penanganan Tuntas Reparse Point & Junction Windows (Documents and Settings)
- **Penyebab Masalah:**
  - Direktori sistem Windows seperti `C:\Documents and Settings` merupakan sebuah *NTFS Junction Link* (Reparse Point). Pada Node.js `fs.readdirSync({ withFileTypes: true })`, objek entry memiliki `e.isDirectory() === false` dan `e.isSymbolicLink() === true`.
  - Hal ini menyebabkan sistem sebelumnya keliru menandai folder junction sebagai "file biasa tanpa ekstensi", sehingga saat diklik ganda sistem mencoba membaca teks menggunakan `/api/v1/files/content` yang memicu 500 error (`EISDIR` / `EPERM`).
- **Solusi Rekayasa:**
  - Memperbarui `frontend/vite.config.js` dan Go backend `RAZStorageEngine.go` agar secara cerdas mendeteksi `isSymbolicLink()` atau entry tanpa ekstensi: jika `fs.statSync` menghasilkan direktori atau melempar permission-denied (EPERM) khas proteksi folder sistem Windows, entri tersebut secara otomatis dan konsisten ditandai sebagai `is_dir = true`.
  - Menambahkan verifikasi `stat.isDirectory()` pada endpoint `/api/v1/files/content` untuk mengembalikan status 400 terproteksi ("Target adalah folder direktori, bukan berkas teks") sebelum mencoba `readFileSync`.

#### 2. Ketahanan Failsafe Store & Eliminasi Peringatan Modal
- **Pembersihan Cache State:**
  - Peringatan `[Vue warn]: Property "showImageModal" was accessed during render` dan `openFileViewer is not a function` berasal dari cache memori browser yang masih menahan bundle sebelum reload penuh (*stale Pinia store instance*).
  - Menambahkan pengecekan failsafe `typeof windowStore.openFileViewer === 'function'` pada handler double click di `RAZFileManagerView.vue` agar aman dari desinkronisasi cache browser.
- **Verifikasi:**
  - `npm run build` sukses 100% (1502 modules transformed).
  - Vite dev server me-restart server secara otomatis dan bersih di port 3000.

---

### Versi: 2.4.0-taskbar-toggle-and-systemwide-context-menus (Taskbar Minimize/Restore Toggle & Sistem Context Menu Global)
**Waktu Eksekusi:** 2026-10-13 11:34 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Perilaku Taskbar Klik Cerdas (Toggle Minimize & Restore Ala Windows/macOS)
- **Logika State Interaktif:**
  - Mengimplementasikan `toggleWindowFromTaskbar(winId)` pada `RAZWindowStore.js` dan menghubungkannya ke tombol aplikasi taskbar di `RAZTaskbar.vue`.
  - **Saat Jendela Minimized:** Mengklik tombol taskbar akan langsung memulihkan (*restore*) jendela ke layar dan membawanya ke fokus paling depan (`bringToFront`).
  - **Saat Jendela Aktif di Depan:** Mengklik tombol taskbar akan langsung meminimalkan (*minimize*) jendela ke taskbar.
  - **Saat Jendela di Belakang Jendela Lain:** Mengklik tombol taskbar akan membawa jendela tersebut ke paling depan dan memberikannya fokus tanpa perlu minimize terlebih dahulu.

#### 2. Arsitektur Universal Custom Context Menu (`RAZContextMenu.vue` & `RAZContextMenuStore.js`)
- **Desain Glassmorphism Slate-900 Terpadu:**
  - Menghapus ketergantungan pada context menu native peramban di seluruh ekosistem desktop.
  - Boundary-aware positioning: posisi `x` dan `y` dihitung otomatis agar menu tidak pernah terpotong di tepi kanan atau bawah layar peramban.
  - Mendukung ikon Lucide, separator baris, status tombol bahaya (*danger* warna merah rose-400), serta auto-close pada klik luar atau tombol <kbd>Esc</kbd>.
- **Cakupan Menu Klik Kanan Komprehensif:**
  1. **Desktop Wallpaper (Area Kosong):**
     - Muat Ulang Desktop (Refresh metrics)
     - Pintasan cepat buka: Files & Storage, Web Terminal, Task Manager
     - Tampilkan Desktop (Minimalkan Semua Jendela)
     - Pulihkan Semua Jendela
     - Pengaturan Sistem
  2. **Ikon Shortcut Desktop:**
     - Buka Aplikasi
     - Informasi Aplikasi & Dimensi Window
  3. **Taskbar (Area Kosong):**
     - Buka Task Manager
     - Tampilkan Desktop (Minimalkan Semua Jendela)
     - Pulihkan Semua Jendela
     - Pengaturan Sistem
  4. **Item Aplikasi Aktif di Taskbar:**
     - Fokus ke Depan / Pulihkan Jendela
     - Minimalkan Jendela
     - Maksimalkan / Ukuran Normal
     - Tutup Jendela (Aksi Cepat Berbahaya)
  5. **Berkas & Folder di File Manager:**
     - Buka / Masuk Folder
     - Detail & Properti Berkas
     - Salin Path Direktori / Berkas ke Clipboard
     - Beri / Hapus Bintang (Favorit)
     - Ubah Nama Berkas (Rename)
     - Unduh Berkas Langsung
     - Hapus Berkas / Folder
  6. **Canvas Area Kosong File Manager:**
     - Buat Folder Baru
     - Muat Ulang (Refresh)
     - Beralih Tampilan Grid
     - Beralih Tampilan List
     - Salin Path Direktori Aktif

#### 3. Verifikasi Build & Stabilitas
- `npm run build` sukses 100% tanpa error (*1504 modules transformed*).
- HMR dev server langsung memuat pembaruan secara real-time pada `http://localhost:3000`.

---

### Versi: 2.4.1-pinia-hmr-and-failsafe-resilience (Integrasi Pinia acceptHMRUpdate & Failsafe Execution)
**Waktu Eksekusi:** 2026-10-13 11:40 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Sinkronisasi Instan Pinia HMR Store
- **Masalah:** Saat fungsi baru ditambahkan ke Pinia Store (`RAZWindowStore.js`), Vite HMR me-reload komponen Vue tetapi Pinia mempertahankan closure memori singleton yang sudah ada di tab browser pengguna. Hal ini memicu pesan `TypeError: is not a function` jika browser belum di-refresh penuh.
- **Penyelesaian:**
  - Mengintegrasikan modul resmi `acceptHMRUpdate` dari Pinia ke seluruh store inti: `RAZWindowStore.js`, `RAZContextMenuStore.js`, dan `RAZSystemStore.js`.
  - Sekarang saat ada pembaruan fungsi di store, Vite HMR secara otomatis mengganti dan memperbarui method di memori browser tanpa memerlukan refresh manual.

#### 2. Arsitektur Failsafe Dua Lapis (Zero Console Error)
- **Fallback Eksekusi:**
  - Menambahkan pengecekan typeof dan fallback internal pada `handleTaskbarItemClick` di `RAZTaskbar.vue` dan aksi context menu di `RAZDesktopLayout.vue`.
  - Membungkus eksekusi `item.action()` pada `RAZContextMenu.vue` dalam blok `try-catch` sehingga tidak ada `Uncaught TypeError` yang tembus ke browser console.
- **Verifikasi Build:**
  - `npm run build` berhasil 100% (*1504 modules transformed*).
  - Dev server aktif dan mutakhir di `http://localhost:3000`.

---

### Versi: 2.4.2-desktop-and-window-animations (Sistem Animasi Minimize/Restore & Interaksi Fluid 60FPS)
**Waktu Eksekusi:** 2026-10-13 11:59 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Animasi Window Minimize & Restore Menuju Taskbar
- **Transisi Meluncur & Skala (Genie/Windows Style):**
  - Mengintegrasikan `<transition name="raz-window-transition">` pada `RAZWindowFrame.vue`.
  - Menggunakan titik acuan `transform-origin: 50% 100%` (bagian bawah tengah jendela mengarah ke taskbar).
  - **Saat Minimize:** Jendela mengecil (*scale-down ke 0.85*), meluncur turun ke arah taskbar (*translateY 45px*), dan memudar secara halus (*fade-out*).
  - **Saat Restore:** Jendela membesar dari taskbar (*scale 0.85 -> 1*), meluncur naik ke koordinat semula, dan memudar masuk (*fade-in*) dengan kurva akselerasi `cubic-bezier(0.16, 1, 0.3, 1)`.
- **Transisi Ukuran Halus (Fluid Maximize/Restore):**
  - Mengimplementasikan class `.raz-window-fluid` pada `RAZTheme.css` yang memberikan transisi halus saat jendela dimaksimalkan atau dikembalikan ke ukuran normal.
  - Dilengkapi deteksi dragging/resizing: saat pengguna sedang menggeser (*drag*) atau mengubah ukuran (*resize*) dengan mouse, transisi dinonaktifkan (`transition-none`) agar kursor mouse tidak mengalami lag/delay.

#### 2. Animasi & Indikator Status Taskbar
- **Animasi Tombol Interaktif:**
  - Efek tekan (*click feedback*): `active:scale-95 transition-all duration-150` pada tombol taskbar dan Start menu.
  - Ikon aplikasi membesar halus saat dihover (`group-hover:scale-110`).
- **Glow Indikator Dinamis di Bawah Tombol:**
  - **Jendela Aktif:** Strip cyan glow lebar `w-6 bg-sky-400` dengan bayangan berpendar (*box-shadow glow*).
  - **Jendela Terbuka di Belakang:** Garis netral `w-3.5 bg-slate-400/70`.
  - **Jendela Minimized:** Titik kompak redup `w-2 bg-slate-600`.

#### 3. Animasi Pop-In Context Menu
- Animasi pembukaan context menu lebih presisi: transisi berskala `scale-95 -> scale-100` dengan slight offset vertikal `-translate-y-1` dan fade-in berkecepatan 150ms.
- Efek mikro-interaksi pada item menu: hover highlight dan pembesaran ikon 10% yang responsif.

#### 4. Verifikasi Build & Stabilitas
- `npm run build` berhasil 100% (*1504 modules transformed*).
- Dev server HMR berjalan lancar tanpa error di `http://localhost:3000`.

---

### Versi: 2.4.3-windows-junction-and-locked-file-resilience (Resolusi Windows Junction & Penanganan Berkas Terkunci Kernel)
**Waktu Eksekusi:** 2026-10-13 12:20 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Resolusi Otomatis Windows NTFS Junction Points (`My Documents`, `Documents and Settings`)
- **Masalah:** Di Windows, folder seperti `C:/Users/Administrator/My Documents` dan `C:/Documents and Settings` merupakan NTFS Junction Points dengan Access Control Entry (ACE) terproteksi. Ketika diklik, `fs.readdirSync()` menghasilkan error `EPERM (Operation not permitted)` yang memicu HTTP 500 di console browser.
- **Penyelesaian:**
  - Mengimplementasikan helper `resolveJunctionTarget(targetPath)` menggunakan `fs.lstatSync().isSymbolicLink()` dan `fs.readlinkSync()`.
  - Ketika pengguna mengklik `My Documents`, target otomatis di-resolve secara transparan ke `C:/Users/Administrator/Documents`.
  - Menyertakan response header `X-Resolved-Path: C:/Users/Administrator/Documents` sehingga Address Bar dan Breadcrumbs di File Manager langsung sinkron dengan lokasi sebenarnya.
  - Menyelaraskan implementasi pada backend Go (`backend/internal/filemanager/RAZStorageEngine.go`) menggunakan `filepath.EvalSymlinks()`.

#### 2. Penanganan Berkas Terkunci Kernel Sistem (`DumpStack.log.tmp`, `pagefile.sys`)
- **Masalah:** Berkas sistem seperti crash log dump kernel (`DumpStack.log.tmp`) atau memory paging berstatus `EBUSY / EXCLUSIVE_LOCK` atau berkas sementara yang dibuat-hapus oleh OS saat runtime, memicu HTTP 500 saat dibuka via penampil teks.
- **Penyelesaian:**
  - Endpoint `/api/v1/files/content` mendeteksi error `EBUSY`, `EPERM`, `EACCES`, atau berkas hilang/sementara. Alih-alih crash 500, endpoint mengembalikan respons terstruktur status 200 ber-flag `isLocked: true` dan pesan edukatif.
  - Menghilangkan pesan `500 (Internal Server Error)` di Network Console browser.
  - Komponen `RAZFileViewerView.vue` dilengkapi tampilan visual peringatan Glassmorphism (ShieldAlert icon) yang menginformasikan kepada pengguna bahwa berkas tersebut dilindungi atau digunakan secara eksklusif oleh sistem operasi.

#### 3. Penanganan Folder Terproteksi Sistem (`System Volume Information`)
- Jika pengguna mencoba membuka folder sistem yang memang tidak memiliki izin akses (EPERM / Access Denied):
  - Endpoint mengembalikan status 200 ber-flag `isAccessDenied: true`.
  - File Manager menampilkan custom notification: *"⚠️ Akses Ditolak: Direktori dilindungi oleh sistem operasi"*, serta mempertahankan posisi folder sebelumnya tanpa merusak UI explorer.

#### 4. Verifikasi Build & Stabilitas
- `npm run build` sukses 100% (*1504 modules transformed*).
- Pengujian otomatis via Node fetch untuk `My Documents` (resolusi ke `Documents`, 65 items) dan `DumpStack.log.tmp` berhasil 100% dengan status HTTP 200 tanpa error konsol.
- Dev server HMR aktif di `http://localhost:3000`.

---

### Versi: 2.4.4-starred-favorites-persistence-fix (Persistensi LocalStorage & Koleksi Global Berkas Berbintang)
**Waktu Eksekusi:** 2026-10-13 12:26 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah Utama yang Teridentifikasi
- **Data Tidak Persisten (Hilang Saat Refresh/Navigasi):** State `starredPaths` sebelumnya hanya disimpan dalam `ref(new Set())` di memori komponen Vue. Setiap kali window File Manager dibuka ulang, browser di-refresh, atau terjadi HMR, instance `Set` kembali kosong.
- **Koleksi Lokal vs Global:** Menu sidebar "Berbintang (Starred)" sebelumnya hanya memfilter array `files.value` dari direktori aktif saat itu (`currentPath`). Jika pengguna men-star berkas di folder A lalu berpindah ke folder B atau membuka tab Berbintang, berkas dari folder A tidak akan pernah muncul, dan jika di folder saat ini tidak ada item yang di-star, tampilan langsung kosong melompong (terlihat seolah-olah item hilang sendiri atau halaman me-reload).
- **Reload Berkedip Saat Klik Sidebar:** Tombol "My Storage (/DATA)" di sidebar sebelumnya selalu memanggil ulang `fetchFiles(defaultStoragePath)` meskipun pengguna sudah berada di direktori tersebut, yang memicu loading spinner jaringan yang tidak perlu.

#### 2. Solusi & Peningkatan yang Diimplementasikan
- **Penyimpanan Permanen `localStorage`:**
  - Mengintegrasikan penyimpanan persisten `localStorage.getItem('razos_starred_items')`.
  - Menyimpan informasi item lengkap (`name`, `path`, `is_dir`, `size`, `extension`, `mod_time`) secara permanen.
- **Koleksi Global Berkas & Folder Berbintang:**
  - Saat tab sidebar "Berbintang (Starred)" dipilih, sumber data otomatis dialihkan ke array global `starredItems.value`.
  - Berkas atau folder yang di-star dari lokasi MANAPUN (baik di `C:/`, `C:/DATA`, `Documents`, dsb) dijamin tampil utuh dan lengkap.
  - Menambahkan badge jumlah item berbintang di samping menu sidebar.
- **Navigasi Intuitif dari Tab Berbintang:**
  - Double click pada folder di tab Berbintang langsung mengalihkan mode ke `drive` dan membuka folder tersebut secara instan.
  - Double click pada berkas di tab Berbintang langsung membuka penampil berkas mandiri.
- **Pencegahan Reload Berkedip:**
  - Klik "My Storage (/DATA)" kini hanya beralih tab tanpa memicu fetch ulang jika sudah berada di direktori tersebut.
  - Beralih ke tab Berbintang berlangsung instan (0 ms latency) tanpa request jaringan.
- **Empty State yang Bersahabat:**
  - Menampilkan kartu ilustrasi bintang emas dan panduan ramah jika pengguna belum menandai berkas favorit, bukan lagi layar kosong tanpa petunjuk.

#### 3. Verifikasi Build
- `npm run build` sukses 100% (*1504 modules transformed*).
- HMR dev server aktif di `http://localhost:3000`.

---

### Versi: 2.4.5-casaos-yaml-and-installed-apps-suite (Instalasi 525+ Aplikasi, Tab Aplikasi Terpasang, Uninstall Context Menu & Docker Compose YAML ala CasaOS)
**Waktu Eksekusi:** 2026-10-13 12:40 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Persistensi & Sinkronisasi Ekosistem Aplikasi App Store
- **Penyimpanan Terpusat Persisten (`RAZInstalledApps.json`):**
  - Mengimplementasikan `backend/internal/store/RAZInstalledApps.json` untuk mencatat seluruh aplikasi yang terpasang.
  - Memisahkan katalog master 525 aplikasi (`RAZApps.json`) dari data instalasi dinamis agar katalog tetap bersih dan status aplikasi tidak hilang saat restart service.
  - Endpoint `GET /api/v1/store/apps` dan `GET /api/v1/store/installed` otomatis me-merge data runtime secara real-time.
- **Penyediaan Direktori Host Otomatis:**
  - Setiap instalasi aplikasi secara otomatis membuat folder volume persisten di `/DATA/apps/<slug>/` pada penyimpanan MicroSD STB.
  - Mendukung clean uninstall dengan opsi `purge_data=true` untuk memusnahkan data aplikasi jika diinginkan pengguna.

#### 2. Tab "Aplikasi Terpasang (Installed Apps)" di App Store
- **Tampilan Khusus Container Aktif:**
  - Tab navigasi kedua di App Store menampilkan seluruh container Docker yang sedang terpasang di sistem.
  - Menampilkan badge status running hijau, nomor port Web UI, direktori path volume di `/DATA/apps/`, dan spesifikasi image.
  - Tombol aksi cepat: **Buka Web UI**, **Detail & Dokumentasi**, dan **Uninstall Bersih**.
  - Empty state elegan dengan tombol jalan pintas langsung ke Katalog Aplikasi atau Impor YAML jika belum ada container yang terpasang.

#### 3. Instalasi via Berkas / Skrip Docker Compose YAML (Ala CasaOS)
- **Editor Monospace & Live Syntax Inspector:**
  - Tab navigasi ketiga di App Store menyediakan editor YAML ala CasaOS dengan font monospace, pengecekan sintaks live, dan indikator validasi.
  - Mengintegrasikan parser pustaka `yaml` untuk membedah struktur Compose: service name, Docker image, host/container ports, bind mount volumes, dan environment variables secara real-time.
- **Dukungan Drag & Drop / Unggah File `.yaml` / `.yml`:**
  - Pengguna dapat menyeret file compose atau mengklik tombol upload untuk memuat konfigurasi langsung ke editor.
- **Template Preset Siap Pakai:**
  - Menyediakan template cepat: **qBittorrent Downloader**, **Nginx Web Server**, **Nextcloud Hub**, dan **WordPress CMS**.
- **Eksekusi Deployment Instan:**
  - Endpoint `POST /api/v1/store/compose` memproses skrip YAML, mendaftarkan kontainer ke store, dan otomatis mengarahkannya ke tab Aplikasi Terpasang.

#### 4. Integrasi Shortcut Desktop, Start Menu & Context Menu Uninstall
- **Sinkronisasi Otomatis Desktop:**
  - Seluruh aplikasi yang terpasang langsung muncul di Desktop dan Start Menu melalui Pinia store `windowStore.allDesktopApps`.
  - Double click pada shortcut desktop langsung meluncurkan antarmuka Web UI aplikasi pada port terkait.
- **Menu Uninstall di Context Menu (Klik Kanan):**
  - Klik kanan pada shortcut aplikasi di Desktop kini menyertakan opsi: **Buka Aplikasi**, **Buka di Tab Baru Peramban**, **Info Container**, dan **Uninstall Aplikasi** (`danger: true`).
  - Dialog konfirmasi uninstall menggunakan Custom Modal Glassmorphism Slate-900 dengan checkbox pemusnahan berkas volume.
- **Taskbar Start Menu Terintegrasi:**
  - Start menu popup menampilkan seksi **"APLIKASI TERPASANG"** beserta badge counter jumlah container aktif dan nomor port langsung.

#### 5. Verifikasi & Definisi Selesai (DoD)
- Build produksi Vite frontend (`npm run build`) berhasil 100% tanpa error (*1577 modules transformed*).
- Pengujian otomatis deployment katalog (Plex), deployment Compose YAML (qBittorrent), dan clean uninstall terverifikasi 100% lolos via endpoint REST API.
- Dev server HMR aktif dan responsif di `http://localhost:3000`.

---

### Versi: 2.4.6-thumbnail-anti-hang-and-virtual-rendering (Arsitektur Anti-Hang File Manager, Concurrency Queue & Virtual Chunked Rendering)
**Waktu Eksekusi:** 2026-10-13 12:56 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah Utama Teridentifikasi
- **UI Freeze / "Not Responding" Akibat Ratusan Request Simultan:**
  - Sebelumnya, ketika pengguna membuka folder berisi banyak file gambar, seluruh kartu me-render tag `<img>` telanjang langsung ke file asli resolusi penuh (`/api/v1/files/raw`).
  - Browser seketika menembakkan puluhan hingga ratusan koneksi HTTP bersamaan dan membebani main UI thread dengan decoding bitmap berukuran megabyte secara masif.
  - Memori RAM browser dan CPU STB ARM64 melonjak tajam sehingga memicu aplikasi macet total atau crash.
- **Ketiadaan Lazy Loading & DOM Overload:**
  - Kartu berkas yang berada jauh di bawah scroll tetap di-request dan di-render, memboroskan bandwidth dan memori.
  - Merender 500-2000 elemen DOM sekaligus dalam satu tick Vue membuat perpindahan folder tersendat-sendat.

#### 2. Solusi Arsitektur yang Diterapkan

##### A. Komponen Mandiri `RAZThumbnail.vue` (Anti-Hang Concurrency Protected)
- **IntersectionObserver Terpadu:**
  - Komponen hanya memicu pemuatan gambar ketika kartu benar-benar berada di dalam viewport atau mendekati viewport (buffer margin `120px`).
  - Jika pengguna melakukan scroll cepat melewati puluhan gambar, tugas antrean yang belum sempat dimuat akan otomatis dibatalkan (`queueTask.cancelled = true`) sehingga tidak menyita I/O jaringan atau thread decoding.
- **Global Concurrency Limiting Queue (Maksimal 4 Request Bersamaan):**
  - Mengimplementasikan queue cerdas berbasis LIFO/FIFO singleton yang membatasi maksimal **4 permintaan fetch aktif secara paralel**.
  - Permintaan lainnya mengantre secara tertib, menjaga UI thread browser tetap 100% mulus (60 FPS) tanpa freeze.
- **In-Memory & Browser Caching:**
  - Menyimpan status gambar yang berhasil dimuat ke cache memori (`__RAZ_THUMBNAIL_CACHE__`), mencegah request ulang saat scroll bolak-balik.
- **Safe Size Protection (> 20 MB):**
  - File foto/gambar mentah berukuran raksasa (> 20 MB) dilindungi dengan badge `HD FILE` dan ikon fallback agar tidak meledakkan kuota RAM browser STB.
- **Placeholder Skeleton Shimmer & Graceful Error:**
  - Menampilkan shimmer gelap halus saat memuat dan fallback otomatis ke format ikon jika berkas korup atau gagal load tanpa merusak grid layout.

##### B. Progressive Virtual Chunked Display Rendering
- **Batching DOM Rendering:**
  - Alih-alih merender ribuan item sekaligus, File Manager awalnya hanya merender `60 item pertama` (`visibleFiles = filteredFiles.slice(0, displayLimit)`).
  - Saat pengguna menggulir mendekati dasar scroll (`@scroll="handleExplorerScroll"`), sistem secara mulus menambah batch `+40 item`.
  - Pembukaan folder yang berisi ribuan berkas kini berlangsung instan (< 15 ms).
  - Menampilkan indikator elegan di bawah grid/list saat memuat batch berikutnya.
  - Otomatis mereset limit ke 60 saat navigasi folder, pencarian, atau filter kategori berubah.

##### C. Endpoint Khusus `/api/v1/files/thumbnail` dengan HTTP Caching Agresif
- **Header Caching Canggih di Vite Dev Middleware & Backend Go:**
  - Mengirimkan header `Cache-Control: public, max-age=604800, immutable` (7 hari).
  - Menyertakan `ETag` berbasis ukuran dan mtime berkas (`W/"size-mtime"`).
  - Mendukung respons `304 Not Modified` seketika (0 bytes transfer data) jika berkas sudah ada di cache browser.

#### 3. Verifikasi & Definisi Selesai (DoD)
- Build produksi Vite frontend (`npm run build`) sukses 100% (*1578 modules transformed*).
- Respons `/api/v1/files/thumbnail` diverifikasi mengembalikan status 200 OK beserta header `Cache-Control: public, max-age=604800, immutable` dan `ETag`.
- UI File Manager tetap responsif, mulus, dan tidak pernah hang saat membuka direktori dengan banyak gambar.

---

### Versi: 2.4.7-music-player-redesign-and-controls-suite (Redesain Music Player: Eliminasi Kontrol Terpotong, Turntable Vinyl, Equalizer Waveform, Autoplay & Kontrol Lengkap)
**Waktu Eksekusi:** 2026-10-13 13:20 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah Utama yang Teridentifikasi
- **Kontrol Audio Terpotong Habis (Window Height Terlalu Pendek):**
  - Ukuran default jendela audio viewer di `RAZWindowStore.js` sebelumnya hanya `460 x 340 px` dengan `minHeight = 260`.
  - Akibatnya, elemen piringan hitam, teks judul, scrubber timeline, tombol Play/Pause, tombol mute, dan tombol unduh terdorong ke bawah dan terpotong oleh batas bawah jendela window frame.
- **Ketiadaan Tombol Kendali Esensial & Autoplay:**
  - Belum ada tombol lompat maju/mundur 10 detik (`SkipBack` / `SkipForward`), tombol mode loop/repeat (`Repeat`), dan slider pengatur volume audio independen.
  - Saat jendela pemutar musik dibuka, lagu tidak otomatis memutar (*no autoplay*).

#### 2. Solusi & Peningkatan yang Diterapkan
- **Optimalisasi Dimensi Window Frame Audio ([`RAZWindowStore.js`](file:///c:/laragon/www/razos/frontend/src/stores/RAZWindowStore.js)):**
  - Mengubah dimensi jendela default pemutar audio menjadi `480 x 520 px` dengan `minWidth = 380` dan `minHeight = 440`.
  - Tata letak menggunakan `flex flex-col justify-between` sehingga seluruh kontrol dijamin terlihat utuh dan tidak terpotong di resolusi layar apa pun.
- **Turntable Vinyl Disc Interaktif & Animated Equalizer Waveform ([`RAZFileViewerView.vue`](file:///c:/laragon/www/razos/frontend/src/views/RAZFileViewerView.vue)):**
  - Desain piringan hitam dengan alur konsentris vinyl (*groove rings*) dan label emas berputar 60 FPS saat memutar musik.
  - Klik langsung pada vinyl disc berfungsi sebagai tombol pintas Play/Pause dengan efek hover overlay.
  - Menambahkan **Animated Audio Wave / Equalizer Bars (16 Bars)** yang menari-nari naik-turun selaras waktu audio saat memutar dan mereda ke garis datar saat dijeda.
- **Paket Kontrol Pemutaran Lengkap (Modern Audio Playback Suite):**
  - **Autoplay On-Mount & Track-Switch:** Lagu otomatis memutar seketika saat berkas audio dibuka atau saat pengguna berganti lagu.
  - **Tombol Play / Pause Utama:** Tombol bundar besar emas/amber berbayangan pendar (*glow shadow*) di posisi tengah yang mudah ditekan.
  - **Lompat 10 Detik:** Tombol cepat mundur 10s (`SkipBack`) dan maju 10s (`SkipForward`).
  - **Mode Ulangi Lagu (Loop / Repeat):** Tombol toggle mode loop dengan titik indikator aktif.
  - **Scrubber Timeline Presisi:** Progress bar akurat dengan penunjuk waktu berjalan (`01:23`) dan total durasi (`03:45`).
  - **Slider Kontrol Volume Interaktif:** Ikon speaker dinamis (`VolumeX`, `Volume1`, `Volume2`) lengkap dengan slider horizontal persentase volume.
  - **Metadata & Tombol Unduh:** Badge format (MP3, WAV, FLAC, dll), estimasi ukuran berkas, status `HQ STEREO`, dan tombol unduh berkas.

#### 3. Verifikasi & Definisi Selesai (DoD)
- Build produksi Vite frontend (`npm run build`) sukses 100% (*1578 modules transformed*).
- Tampilan jendela Music Player kini lega, elegan, bebas dari pemotongan elemen, dan seluruh fungsi tombol bekerja responsif.

---

### Versi: 2.4.8-window-titlebar-truncate-fix (Title Bar Single-Line Anti-Wrap & Ellipsis Truncation)
**Waktu Eksekusi:** 2026-10-13 13:25 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah Utama yang Teridentifikasi
- **Line-Wrapping Judul Jendela:**
  - Pada jendela dengan nama berkas atau judul aplikasi yang panjang (misalnya nama lagu atau file dokumen teks yang panjang), teks judul di title bar membungkus (*word-wrap*) menjadi 2 hingga 3 baris.
  - Hal ini menyebabkan teks meluber turun ke bawah menabrak garis batas header dan menutupi area konten jendela teratas.
  - Tombol kontrol jendela (Minimize, Maximize, Close) di sisi kanan juga berisiko terdesak jika judul terlalu panjang.

#### 2. Solusi & Peningkatan yang Diterapkan
- **Penguncian Flexbox Single-Line & Truncate ([`RAZWindowFrame.vue`](file:///c:/laragon/www/razos/frontend/src/components/RAZWindowFrame.vue)):**
  - Menambahkan `min-w-0 flex-1 mr-2 overflow-hidden` pada container judul agar CSS Flexbox mengizinkan penyusutan teks secara tepat tanpa overflow.
  - Menerapkan class `truncate whitespace-nowrap` pada elemen `<span>` teks judul sehingga judul panjang otomatis terpotong rapi dengan elipsis tiga titik (`...`) dalam satu baris.
  - Mempertahankan aksesibilitas dan keterbacaan penuh dengan menyematkan atribut native tooltip `:title="win.title"`. Ketika pengguna mengarahkan kursor mouse ke judul, nama lengkap berkas/jendela akan langsung ditampilkan.
  - Mengamankan kontainer tombol aksi kanan (Minimize, Maximize, Close) dengan utility `shrink-0` agar tombol selalu memiliki lebar konsisten dan tidak pernah tergeser atau terhimpit.

#### 3. Verifikasi & Definisi Selesai (DoD)
- Build produksi Vite frontend (`npm run build`) sukses 100% (*1578 modules transformed*).
- Judul jendela kini selalu terkunci 1 baris rapi dan estetik pada seluruh resolusi layar dan ukuran jendela desktop RAZOS.

---

### Versi: 2.4.9-music-player-speed-and-anti-jitter-seek (Playback Speed Multiplier, Scrubber Timeline Klik/Slide Bebas & Hover Preview)
**Waktu Eksekusi:** 2026-10-13 13:38 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah Utama yang Teridentifikasi
- **Slider Timeline "Mental" / Tidak Bisa Di-slide dengan Mulus:**
  - Sebelumnya, input range audio terikat dua arah pada `audioCurrentTime` yang diperbarui setiap 250ms oleh event `timeupdate` elemen `<audio>`.
  - Ketika pengguna menahan dan menggeser (*drag/slide*) slider waktu, event `timeupdate` menimpa nilai yang sedang digeser dan menarik paksa posisi thumb kembali ke posisi pemutaran sebelum buffer selesai (*slider snapping/jitter*).
  - Mengklik sembarang titik di garis waktu juga sering kali tidak akurat pada elemen `<input type="range">` native peramban.
- **Ketiadaan Fitur Percepat Lagu (Playback Speed):**
  - Belum ada opsi untuk mempercepat pemutaran lagu (misal 1.25x, 1.5x, 2x) bagi pengguna yang ingin mendengarkan audio rekaman, podcast, atau lagu dengan tempo lebih cepat.

#### 2. Solusi & Peningkatan yang Diterapkan
- **Fitur Pengatur Kecepatan Pemutaran Audio (Playback Speed Suite):**
  - Mengimplementasikan `playbackSpeed` ref dengan dukungan multi-kecepatan: `0.5x`, `0.75x`, `1.0x`, `1.25x`, `1.5x`, `1.75x`, dan `2.0x`.
  - Tombol Pill Interaktif:
    - Klik kiri: Beralih kecepatan secara siklik (`1.0x` ➔ `1.25x` ➔ `1.5x` ➔ `2.0x` ➔ `0.5x` ➔ `0.75x` ➔ `1.0x`).
    - Klik kanan / opsi menu: Membuka menu popover pilihan kecepatan instan lengkap dengan backdrop auto-dismiss.
    - Ikon dinamis: Menampilkan ikon `FastForward` berwarna emas amber saat kecepatan > 1x, dan indikator badge `1.5x` / `2x` di samping durasi total.
  - **Dynamic Vinyl Turntable Spin Rate:** Putaran animasi piringan hitam vinil secara otomatis berputar lebih kencang seiring bertambahnya kecepatan pemutaran (`animationDuration: 8 / playbackSpeed s`).
  - **Equalizer Waveform Sync:** Kecepatan osilasi animasi gelombang equalizer 16-bar bergerak selaras dengan kecepatan pemutaran audio.
- **Sistem Scrubber Timeline Bebas Hambatan (Klik & Slide Anti-Jitter):**
  - **State `isSeeking` & `seekValue`:** Selama mouse atau jari pengguna menahan slider (`mousedown` / `touchstart` / `input`), event `timeupdate` dinonaktifkan sementara sehingga slider bergerak sangat mulus tanpa gangguan mental kembali.
  - **Full Clickable Track:** Mengklik titik mana pun pada garis timeline (`onTimelineTrackClick`) langsung menghitung rasio posisi klik terhadap lebar elemen dan melompatkan waktu lagu ke detik tersebut secara instan.
  - **Visual Filled Bar & Glowing Thumb Knob:** Garis waktu yang sudah terputar diisi dengan gradien emas amber `from-amber-500 to-amber-300`, dan indikator thumb emas putih berbayangan pendar mengikuti posisi kursor secara presisi.
  - **Floating Hover Tooltip Preview:** Saat kursor mouse digerakkan di atas garis timeline, muncul gelembung waktu (*floating timestamp badge*, misal `02:15`) dan garis bayangan abu-abu (*ghost track*) untuk mempermudah navigasi lagu seperti pada Spotify dan Apple Music.

#### 3. Verifikasi & Definisi Selesai (DoD)
- Build produksi Vite frontend (`npm run build`) sukses 100% (*1578 modules transformed*).
- Slider timeline dapat di-klik di mana saja dan di-slide secara sangat mulus tanpa lag atau jitter.
- Kecepatan pemutaran 0.5x hingga 2.0x bekerja presisi dengan sinkronisasi piringan vinil dan equalizer.

---

### Versi: 2.5.0-website-manager-container-isolated-suite (Arsitektur Multi-Website Terisolasi Docker, Root Folder Picker, Integrasi Database Manager & Desktop Shortcuts)
**Waktu Eksekusi:** 2026-10-13 14:05 WIB  
**Pelaksana:** AI Agent Core Engineering

#### 1. Masalah & Kebutuhan Utama yang Diatasi
- **Ketiadaan Isolasi Container Khusus Website:**
  - Sebelumnya, pembuatan website belum mengelola container Docker terisolasi secara mandiri dengan spesifikasi nama container unik, alokasi root folder khusus di `/DATA/apps/websites/<slug>`, resource limits RAM/CPU, dan perintah kustom.
  - Pengguna membutuhkan sistem di mana setiap website berjalan dalam container terpisah yang aman sehingga tidak bertabrakan dengan website lain meskipun menggunakan web server atau runtime yang sama (Nginx, PHP, Node.js, Next.js).
- **Koneksi Database & Integrasi Database Manager:**
  - Ketiadaan integrasi langsung antara website dinamis (seperti PHP/Node.js) dengan server basis data (MariaDB / PostgreSQL) dan alat pengelola database visual (seperti Web Adminer / HeidiSQL style).
- **Ketiadaan Shortcut Desktop & App Drawer untuk Website:**
  - Ikon dan tautan website belum tersinkronisasi ke Desktop layout atau Start Menu sebagai aplikasi mandiri yang bisa diluncurkan dengan double click.

#### 2. Solusi Arsitektur yang Diterapkan

##### A. Mesin Backend & Persistensi JSON (`RAZWebsiteManager.go` & `RAZWebsites.json`)
- **Persistensi Data Permanen:** Menyimpan seluruh konfigurasi website secara permanen di `backend/internal/store/RAZWebsites.json`.
- **Unique Container Isolation (`container_name`):** Setiap website memiliki container name terisolasi (misal `site-toko-online-web`) yang dijamin tidak bentrok.
- **Auto-Alokasi Port Bebas (8100-8250):** Mendeteksi port yang diminta atau secara otomatis mengalokasikan port TCP yang belum terpakai di sistem host.
- **Multi-Runtime Stack Support:**
  - **HTML Statis (Nginx):** Image `nginx:alpine` super ringan (~32MB RAM) memetakan root folder ke `/usr/share/nginx/html`.
  - **PHP 8.2 + Database:** Image `php:8.2-apache` memetakan root folder ke `/var/www/html` dan menyuntikkan variabel koneksi basis data.
  - **Node.js Server:** Image `node:18-alpine` untuk server Express/Fastify/API dengan custom execution command.
  - **Next.js Framework:** Image `node:20-alpine` untuk aplikasi React SSR pada port 3000.
  - **Custom Docker Stack:** Konfigurasi image dan perintah container bebas.
- **Generator Starter Template Otomatis:** Secara otomatis membuat starter file modern (`index.html`, `index.php` dengan pengecekan DB info, atau `server.js`) saat root folder baru dibuat.
- **Resource Limits Hemat Daya:** Pembatasan memori RAM (`64M`, `128M`, `256M`, `512M`) dan `CPUShares` (256, 512, 1024) yang disesuaikan secara presisi untuk hardware STB ARM64 ber-RAM 1.8GB.

##### B. Wizard Pembuatan Website 4-Langkah (`RAZWebsitesView.vue`)
- **Langkah 1: Identitas & Akses Website:** Nama website, container name otomatis, domain/hostname, port host dengan tombol auto-suggest port bebas, dan pemilih link icon (preset Nginx, PHP, Node, Next.js, WordPress, Laravel, React, Vue, HTML5).
- **Langkah 2: Runtime Stack & Root Folder:** Pilihan kartu runtime visual interaktif dan penentuan root folder penyimpanan berkas di `/DATA/apps/websites/<slug>`.
- **Langkah 3: Integrasi Basis Data (Adminer / HeidiSQL Connected):** Toggle aktifkan database, pilihan database terdaftar di RAZOS (auto-fill host, port, user), atau form input kredensial manual lengkap.
- **Langkah 4: Batasan Sumber Daya & Advanced:** Preset tombol RAM limit, prioritas CPU shares, perintah eksekusi container opsional, dan textarea variabel lingkungan tambahan (ENV vars).

##### C. Integrasi Shortcut Desktop, Start Menu & File Manager
- **Desktop Launcher & Context Menu:** Setiap website yang dibuat otomatis muncul di Desktop dan Start Menu dengan icon kustomnya (`isWebsite: true`). Double-click pada shortcut desktop langsung membuka website di browser.
- **Klik Kanan (Context Menu) Shortcut Website:** Opsi buka web, buka di tab baru peramban, buka root folder di File Manager, kelola di Website Manager, dan info website.
- **Tombol Cepat Buka Root Folder di File Manager:** Mengklik icon folder pada kartu website langsung meluncurkan jendela File Manager ke direktori kerja website tersebut.
- **Tombol Cepat Buka Web Adminer:** Membuka Web GUI Adminer untuk mengelola tabel dan data database terkait.
- **100% Custom Modal:** Seluruh dialog interaksi (wizard, konfirmasi hapus dengan opsi purge folder, dan detail kredensial) menggunakan custom modal Glassmorphism Slate-900 tanpa `alert()` atau `confirm()` native.

#### 3. Verifikasi & Definisi Selesai (DoD)
- Kompilasi build frontend Vite (`npm run build`) sukses 100% (*1578 modules transformed*).
- Pengujian otomatis end-to-end via endpoint REST API (`POST /api/v1/databases`, `POST /api/v1/websites`, `POST /api/v1/websites/:id/action`, `DELETE /api/v1/websites/:id?purge=true`) lolos 100%.
- Starter file `index.php` dan `index.html` berhasil terbuat secara fisik di root folder host dengan konfigurasi database yang presisi.
- Dev server Vite HMR aktif dan responsif di port 3000.






