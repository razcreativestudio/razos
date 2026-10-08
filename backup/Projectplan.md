# RAZOS Project Plan

## 1. Ringkasan Eksekutif
RAZOS (RAZ Operating System) adalah sistem manajemen server berbasis web yang memadukan pengalaman desktop multi-window (mode PC) dan Android launcher (mode ponsel), dirancang khusus untuk mengubah perangkat keras berdaya rendah seperti STB Armbian menjadi mini cloud server berkinerja tinggi.

## 2. Metrik Keberhasilan
- Penggunaan RAM idle panel: <= 35 MB.
- Ukuran binary akhir: <= 25 MB (termasuk static asset frontend yang di-embed).
- Waktu booting service: < 2 detik.
- Waktu deploy container aplikasi baru: < 30 detik.

## 3. Matriks Alokasi Sumber Daya
- **Penyimpanan Primer:** MicroSD `/DATA` untuk Docker root, direktori web, dan database.
- **Penyimpanan Sistem:** eMMC 8 GB hanya untuk base OS Armbian dan binary RAZOS.
- **Limit Memori:** Setiap container tambahan dibatasi memory limit (`mem_limit`) agar sistem tidak memicu Out-of-Memory (OOM) Killer.
