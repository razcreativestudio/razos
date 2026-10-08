# Kebijakan Keamanan & Privasi RAZOS

## 1. Manajemen Autentikasi
- Kata sandi dienkripsi menggunakan hashing **Argon2id** atau **Bcrypt** (cost 12).
- Autentikasi API menggunakan stateless **JWT (JSON Web Token)** yang disimpan dalam cookie HttpOnly, SameSite Strict.
- Sesi kedaluwarsa otomatis dalam 24 jam tidak aktif.

## 2. Isolasi Docker Engine
- Endpoint Docker Engine tidak boleh diekspos ke TCP publik (`0.0.0.0:2375`).
- Interaksi backend ke Docker terbatas hanya melalui socket lokal Unix `/var/run/docker.sock` dengan sanitasi ketat pada argumen pembuatan container.

## 3. Pembatasan Shell Web Terminal
- Koneksi terminal WebSocket diwajibkan menyertakan token JWT yang valid pada handshake awal.
- Shell dijalankan langsung di bawah sesi root host lokal secara terisolasi dan langsung diputus jika koneksi peramban terputus.

## 4. Privasi & Telemetri
- RAZOS bebas dari analitik pihak ketiga atau pelacak telemetri.
- Pengecekan versi hanya menghubungi endpoint resmi API GitHub Release tanpa mengirimkan metadata hardware pengguna.
