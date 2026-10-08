# Spesifikasi REST API & WebSocket RAZOS

Seluruh endpoint diawali dengan prefix `/api/v1`. Autentikasi menggunakan header `Authorization: Bearer <TOKEN>`.

## 1. Autentikasi
- `POST /auth/login`
  - Body: `{"username": "admin", "password": "..."}`
  - Response: `{"token": "JWT_TOKEN"}`
- `GET /auth/me`
  - Response: `{"username": "admin"}`

## 2. Metrik & Pemeliharaan Sistem
- `GET /system/metrics`
  - Response: CPU percent, RAM used/total, Suhu, ZRAM percent, Disk percent.
- `POST /system/clean-ram`
  - Menjalankan pembersihan cache RAM.
  - Response: `{"status": "success", "freed_ram_mb": 120}`
- `POST /system/clean-storage`
  - Menjalankan pemangkasan container/image/volume mati.
  - Response: `{"status": "success", "space_reclaimed": "1.2GB"}`

## 3. Manajemen Container & Task Manager
- `GET /containers`
  - Mendapatkan daftar container aktif beserta penggunaan CPU & RAM.
- `POST /containers/:id/stop`
  - Menghentikan container target.
- `POST /containers/:id/restart`
  - Memulai ulang container target.
- `DELETE /containers/:id`
  - Query: `?purge_data=true|false`
  - Menghapus container dan data terkait.

## 4. Manajemen Website
- `GET /websites`
  - Daftar situs web yang terdaftar.
- `POST /websites`
  - Body: `{"name": "Blog", "domain": "blog.razos.my.id", "runtime": "php", "db_engine": "mariadb"}`
  - Response: `{"status": "created", "port": 8101}`

## 5. WebSocket Terminal
- `WS /ws/terminal`
  - Protocol: WebSocket
  - Handshake: Header `Sec-WebSocket-Protocol: token.<JWT_TOKEN>`
  - Data Channel: Meneruskan raw input/output PTY binary stream.
