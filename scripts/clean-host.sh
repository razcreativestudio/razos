#!/bin/bash
# ==============================================================================
# RAZOS Host Cleaner Script
# Updated: 2026-10-08
# Deskripsi: Skrip pembersihan total sisa sistem panel lama (CasaOS, 1Panel, container yatim piatu)
# sebelum instalasi bersih RAZOS pada STB Armbian.
# ==============================================================================

set -e

echo "=========================================================="
echo " [RAZOS] Memulai Pembersihan Total Host STB Armbian..."
echo "=========================================================="

# 1. Hentikan dan copot daemon CasaOS jika ada
if systemctl is-active --quiet casaos 2>/dev/null; then
    echo "[-] Menonaktifkan dan menghapus layanan CasaOS..."
    systemctl stop casaos 2>/dev/null || true
    systemctl disable casaos 2>/dev/null || true
fi

# 2. Hentikan 1Panel jika terpasang
if command -v 1pctl &> /dev/null; then
    echo "[-] Mencopot instalasi 1Panel..."
    1pctl uninstall 2>/dev/null || true
fi

# 3. Pembersihan Docker menyeluruh
if command -v docker &> /dev/null; then
    echo "[-] Menghentikan seluruh container Docker..."
    docker stop $(docker ps -aq) 2>/dev/null || true
    
    echo "[-] Menghapus seluruh container mati..."
    docker rm $(docker ps -aq) 2>/dev/null || true
    
    echo "[-] Menghapus image yang tidak terpakai..."
    docker rmi $(docker images -q) -f 2>/dev/null || true
    
    echo "[-] Pruning volume dan build cache Docker..."
    docker system prune -a --volumes -f 2>/dev/null || true
fi

# 4. Sinkronisasi Kernel & Pembersihan Cache RAM
echo "[-] Membebaskan cache memori kernel (drop_caches)..."
sync
echo 3 > /proc/sys/vm/drop_caches

echo "=========================================================="
echo " [OK] Host bersih sempurna. Kapasitas eMMC telah dipulihkan."
echo " Siap untuk instalasi bersih RAZOS."
echo "=========================================================="
