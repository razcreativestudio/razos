#!/bin/bash
# ==============================================================================
# RAZOS 1-Click Smart Installer Script
# Updated: 2026-10-25 (v3.8.1 - Fix Go Module Dependencies & Automatic go mod tidy)
# Deskripsi: Skrip pemasangan mandiri RAZOS pada STB Linux ARM64 (aarch64).
# Mengonfigurasi direktori MicroSD (/DATA), mengunduh binary rilis resmi,
# atau melakukan kompilasi otomatis langsung di host jika release binary belum tersedia.
# ==============================================================================

set -e

echo "=========================================================="
echo "          RAZOS - RAZ Operating System Installer         "
echo "                 Khusus STB Armbian ARM64                 "
echo "=========================================================="

ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ] && [ "$ARCH" != "arm64" ]; then
    echo "[!] Peringatan: Arsitektur terdeteksi $ARCH. Binary dioptimalkan untuk ARM64 (aarch64)."
fi

# 1. Siapkan direktori penyimpanan kerja di MicroSD
echo "[+] Menyiapkan direktori penyimpanan kerja di /DATA..."
mkdir -p /DATA/razos/bin
mkdir -p /DATA/razos/data
mkdir -p /DATA/razos/src
mkdir -p /DATA/apps

# 2. Unduh binary rilis ARM64 atau Fallback Kompilasi
echo "[+] Mengunduh binary statis RAZOS terbaru..."
RELEASE_URL="https://github.com/razcreativestudio/razos/releases/latest/download/razos-linux-arm64"

DOWNLOAD_SUCCESS=false
if curl -fsSL --connect-timeout 10 "$RELEASE_URL" -o /DATA/razos/bin/razos 2>/dev/null; then
    # Validasi apakah file bukan file kosong / 404 html
    if [ -s /DATA/razos/bin/razos ] && [ $(wc -c < /DATA/razos/bin/razos) -gt 100000 ]; then
        DOWNLOAD_SUCCESS=true
        chmod +x /DATA/razos/bin/razos
        echo "[✓] Binary resmi RAZOS ARM64 berhasil diunduh."
    fi
fi

# Fallback: Jika binary release belum tersedia di GitHub Releases (HTTP 404), build mandiri
if [ "$DOWNLOAD_SUCCESS" = false ]; then
    echo "[!] Binary pra-rilis GitHub belum tersedia. Memulai mode kompilasi otomatis di STB..."
    
    # Pastikan git dan golang terpasang
    if ! command -v git &> /dev/null || ! command -v go &> /dev/null; then
        echo "[+] Memasang dependensi kompilasi ringan (git, golang-go)..."
        apt update -y
        apt install -y git golang-go
    fi

    echo "[+] Mengambil source code terbaru dari GitHub..."
    rm -rf /DATA/razos/src
    git clone --depth 1 https://github.com/razcreativestudio/razos.git /DATA/razos/src

    echo "[+] Menyiapkan dependensi modul Go (go mod tidy)..."
    cd /DATA/razos/src/backend
    go mod tidy

    echo "[+] Melakukan kompilasi binary server RAZOS ARM64..."
    go build -ldflags="-s -w" -o /DATA/razos/bin/razos ./cmd/server/RAZMainServer.go
    chmod +x /DATA/razos/bin/razos
    echo "[✓] Kompilasi binary RAZOS ARM64 berhasil diselesaikan."
fi

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
echo "[+] Memulai dan mengaktifkan service RAZOS..."
systemctl daemon-reload
systemctl enable --now razos

# 5. Verifikasi Status Layanan
sleep 2
if systemctl is-active --quiet razos; then
    IP_HOST=$(hostname -I | awk '{print $1}')
    echo "=========================================================="
    echo " [OK] Instalasi & Service RAZOS Berjalan Sempurna!"
    echo " Dashboard Web siap diakses di: http://${IP_HOST}:8080"
    echo "=========================================================="
else
    echo "=========================================================="
    echo " [!] Peringatan: Service razos belum aktif. Cek log:"
    journalctl -u razos -n 15 --no-pager
    echo "=========================================================="
fi
