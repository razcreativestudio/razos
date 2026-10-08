# Arsitektur Sistem RAZOS

## 1. Diagram Tingkat Tinggi

```
+-----------------------------------------------------------------------+
|                         Browser Client (PWA)                         |
|  [Breakpoint > 768px: Desktop Mode] | [Breakpoint <= 768px: Mobile]   |
+------------------------------------+----------------------------------+
| (HTTP REST / WebSocket)
v
+-----------------------------------------------------------------------+
|                    RAZOS Core Engine (Golang Binary)                  |
|                                                                       |
|  +--------------------+  +--------------------+  +------------------+ |
|  | Reverse Proxy Core |  | WebSocket Terminal |  | REST Controllers | |
|  +--------------------+  +--------------------+  +------------------+ |
|            |                       |                       |          |
|  +--------------------+  +--------------------+  +------------------+ |
|  | Caddy Routing Hub  |  |   PTY Local Bash   |  | Services Domain  | |
|  +--------------------+  +--------------------+  +------------------+ |
+------------------------------------------------------------+----------+
|
+-----------------------------------------+
v                                         v
+-------------------------------------+   +-----------------------------+
|    Docker Engine (Local Daemon)     |   |      Linux Kernel Space     |
|  /var/run/docker.sock               |   |  /proc/stat, /sys/thermal   |
|                                     |   |  /proc/sys/vm/drop_caches   |
|  [Central MariaDB] [Central PG]     |   +-----------------------------+
|  [Adminer Web GUI] [User Websites]  |
|  [App Store Containers]             |
+-------------------------------------+
```

## 2. Aliran Data Komponen
1. **Frontend-Backend Binding:** Single Page Application (Vue 3) di-build menjadi asset statis dan diintegrasikan ke dalam binary Golang menggunakan `embed.FS` standar. Server berjalan di satu port tunggal (`:8080`).
2. **Reverse Proxy & Website Isolation:** Setiap website yang dibuat dialokasikan ke port unik (`8100-8200`). Core proxy internal memetakan domain/subdomain langsung ke port container target tanpa perlu restart service web server utama.
