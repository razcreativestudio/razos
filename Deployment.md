# Panduan Deployment & Pembersihan Host STB

## 1. Pembersihan Total Docker Lama
Jalankan skrip ini via SSH untuk menghapus seluruh sisa container CasaOS, 1Panel, dan database lama:

```bash
#!/bin/bash
set -e
echo "Memulai pembersihan total Docker..."
docker stop $(docker ps -aq) 2>/dev/null || true
docker rm $(docker ps -aq) 2>/dev/null || true
docker rmi $(docker images -q) -f 2>/dev/null || true
docker system prune -a --volumes -f
systemctl stop casaos 2>/dev/null || true
systemctl disable casaos 2>/dev/null || true
1pctl uninstall 2>/dev/null || true
echo "Host bersih. Sisa penyimpanan eMMC telah dipulihkan."
```

## 2. Pemasangan RAZOS via Skrip Otomatis

```bash
# Buat direktori kerja di MicroSD
mkdir -p /DATA/razos/bin
mkdir -p /DATA/razos/data
mkdir -p /DATA/apps

# Unduh binary rilis terbaru ARM64
curl -fsSL https://github.com/username/razos/releases/latest/download/razos-linux-arm64 -o /DATA/razos/bin/razos
chmod +x /DATA/razos/bin/razos

# Pasang Systemd Service
cat <<EOF> /etc/systemd/system/razos.service
[Unit]
Description=RAZOS - Web-based Cloud OS
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/DATA/razos
ExecStart=/DATA/razos/bin/razos
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now razos
```

## 3. Integrasi Cloudflare Tunnel

Arahkan hostname publik Anda (misalnya `os.razos.my.id`) di dashboard Cloudflare Zero Trust ke:

* **Service:** `HTTP`
* **URL:** `localhost:8080`
