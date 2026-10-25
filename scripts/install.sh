#!/bin/bash
# ==============================================================================
# RAZOS 1-Click Smart Installer Script
# Updated: 2026-10-25 (v3.8.2 - Automated Go 1.22 Toolchain Setup for ARM64)
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
    
    # Pastikan git dan curl terpasang
    if ! command -v git &> /dev/null || ! command -v curl &> /dev/null; then
        echo "[+] Memasang dependensi ringan (git, curl)..."
        apt update -y
        apt install -y git curl
    fi

    # Pastikan Go 1.22+ terpasang untuk mendukung modul modern (cmp, slices, log/slog)
    GO_NEED_INSTALL=false
    if ! command -v go &> /dev/null; then
        GO_NEED_INSTALL=true
    else
        GO_CUR_VER=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//' | cut -d'.' -f1,2)
        if [ "$GO_CUR_VER" = "1.18" ] || [ "$GO_CUR_VER" = "1.19" ] || [ "$GO_CUR_VER" = "1.20" ]; then
            GO_NEED_INSTALL=true
        fi
    fi

    if [ "$GO_NEED_INSTALL" = true ]; then
        echo "[+] Memasang Go 1.22 ARM64 resmi dari go.dev..."
        curl -fsSL https://go.dev/dl/go1.22.7.linux-arm64.tar.gz | tar -C /usr/local -xz
        ln -sf /usr/local/go/bin/go /usr/bin/go
        ln -sf /usr/local/go/bin/gofmt /usr/bin/gofmt
    fi

    echo "[+] Mengambil source code terbaru dari GitHub..."
    rm -rf /DATA/razos/src
    git clone --depth 1 https://github.com/razcreativestudio/razos.git /DATA/razos/src

    echo "[+] Menyiapkan dependensi modul Go (GOTOOLCHAIN=local, GOPATH di MicroSD)..."
    export GOTOOLCHAIN=local
    export GOPATH=/DATA/razos/gopath
    export GOCACHE=/DATA/razos/gocache
    mkdir -p "$GOPATH" "$GOCACHE"
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
