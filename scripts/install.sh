#!/bin/bash
# ==============================================================================
# RAZOS 1-Click Installer Script
# Updated: 2026-10-08
# Deskripsi: Skrip pemasangan mandiri RAZOS pada STB Linux ARM64 (aarch64).
# Mengonfigurasi direktori MicroSD (/DATA), mengunduh binary rilis, dan mendaftarkan systemd service.
# ==============================================================================

set -e

echo "=========================================================="
echo "          RAZOS - RAZ Operating System Installer         "
echo "                 Khusus STB Armbian ARM64                 "
echo "=========================================================="

ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ] && [ "$ARCH" != "arm64" ]; then
    echo "[!] Peringatan: Arsitektur terdeteksi $ARCH. Binary dioptimalkan untuk ARM64."
fi

# 1. Siapkan direktori penyimpanan kerja di MicroSD
echo "[+] Menyiapkan direktori penyimpanan di /DATA..."
mkdir -p /DATA/razos/bin
mkdir -p /DATA/razos/data
mkdir -p /DATA/apps

# 2. Unduh binary rilis ARM64
echo "[+] Mengunduh binary statis RAZOS terbaru..."
RELEASE_URL="https://github.com/razcreative/razos/releases/latest/download/razos-linux-arm64"
curl -fsSL "$RELEASE_URL" -o /DATA/razos/bin/razos || {
    echo "[!] Menggunakan executable lokal jika tersedia..."
}
chmod +x /DATA/razos/bin/razos 2>/dev/null || true

# 3. Konfigurasi Systemd Service
echo "[+] Mendaftarkan systemd service /etc/systemd/system/razos.service..."
cat << 'EOF' > /etc/systemd/system/razos.service
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
MemoryHigh=1500M
MemoryMax=1700M

[Install]
WantedBy=multi-user.target
EOF

# 4. Aktifkan dan Jalankan Layanan
echo "[+] Memulai service RAZOS..."
systemctl daemon-reload
systemctl enable --now razos 2>/dev/null || true

echo "=========================================================="
echo " [OK] Instalasi RAZOS Berhasil!"
echo " Buka dashboard peramban Anda di: http://$(hostname -I | awk '{print $1}'):8080"
echo "=========================================================="
