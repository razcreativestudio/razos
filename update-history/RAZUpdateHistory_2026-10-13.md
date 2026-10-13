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
