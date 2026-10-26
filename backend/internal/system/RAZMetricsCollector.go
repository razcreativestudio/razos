/**
 * RAZOS System Metrics Collector
 * Updated: 2026-10-08
 * Deskripsi: Modul pengambil data metrik hardware host (CPU, RAM, Suhu, ZRAM, Disk)
 * khusus dirancang untuk lingkungan STB ARM64 dengan konsumsi CPU/memori ultra-rendah.
 */

package system

import (
	"os"
	"strconv"
	"strings"

	"github.com/razos/razos/internal/filemanager"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// DiskMetric mendeskripsikan ringkasan metrik setiap storage fisik aktif
type DiskMetric struct {
	Name       string  `json:"name"`       // Nama aslinya (misal: /DATA, /, C:, /media/usb)
	Mountpoint string  `json:"mountpoint"` // Jalur mountpoint
	Device     string  `json:"device"`     // Jalur device partisi
	Fstype     string  `json:"fstype"`     // Format filesystem (ext4, ntfs, dll)
	Total      uint64  `json:"total"`      // Kapasitas total dalam bytes
	Used       uint64  `json:"used"`       // Kapasitas terpakai dalam bytes
	Free       uint64  `json:"free"`       // Kapasitas kosong dalam bytes
	Percent    float64 `json:"percent"`    // Persentase penggunaan disk
}

// Metrics merepresentasikan struktur data performa sistem secara komprehensif
type Metrics struct {
	CPUPercent  float64      `json:"cpu_percent"`  // Persentase penggunaan CPU keseluruhan
	CPUTemp     float64      `json:"cpu_temp"`     // Suhu SOC STB dalam derajat Celcius
	RAMTotal    uint64       `json:"ram_total"`    // Total RAM dalam bytes
	RAMUsed     uint64       `json:"ram_used"`     // RAM terpakai dalam bytes
	RAMPercent  float64      `json:"ram_percent"`  // Persentase RAM fisik terpakai
	ZRAMPercent float64      `json:"zram_percent"` // Persentase ZRAM / Swap terpakai
	DiskTotal   uint64       `json:"disk_total"`   // Total penyimpanan data (/DATA)
	DiskUsed    uint64       `json:"disk_used"`    // Penyimpanan terpakai (/DATA)
	DiskPercent float64      `json:"disk_percent"` // Persentase disk terpakai
	DiskName    string       `json:"disk_name"`    // Nama storage utama
	Disks       []DiskMetric `json:"disks"`        // Daftar seluruh storage disk fisik aktif
	NetSent     uint64       `json:"net_sent"`     // Total byte dikirim jaringan
	NetRecv     uint64       `json:"net_recv"`     // Total byte diterima jaringan
}

// GetMetrics mengumpulkan metrik sistem secara real-time dari kernel Linux host
func GetMetrics() (*Metrics, error) {
	// Ambil data virtual memory fisik
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Ambil data swap / ZRAM memory
	sMem, _ := mem.SwapMemory()

	// Ambil persentase beban CPU
	cpuStats, _ := cpu.Percent(0, false)
	cpuVal := 0.0
	if len(cpuStats) > 0 {
		cpuVal = cpuStats[0]
	}

	// Kumpulkan seluruh storage fisik aktif melalui filemanager engine
	devices, _ := filemanager.GetStorageDevices()
	var diskList []DiskMetric
	diskTotal := uint64(0)
	diskUsed := uint64(0)
	diskPercent := 0.0
	diskName := "/DATA"

	for _, dev := range devices {
		diskList = append(diskList, DiskMetric{
			Name:       dev.Name,
			Mountpoint: dev.Mountpoint,
			Device:     dev.Device,
			Fstype:     dev.Fstype,
			Total:      dev.Total,
			Used:       dev.Used,
			Free:       dev.Free,
			Percent:    dev.Percent,
		})
	}

	// Tentukan primary disk: prioritaskan /DATA, fallback ke disk pertama
	primaryFound := false
	for _, d := range diskList {
		if d.Mountpoint == "/DATA" || d.Mountpoint == "C:\\DATA" {
			diskTotal = d.Total
			diskUsed = d.Used
			diskPercent = d.Percent
			diskName = d.Name
			primaryFound = true
			break
		}
	}
	if !primaryFound && len(diskList) > 0 {
		diskTotal = diskList[0].Total
		diskUsed = diskList[0].Used
		diskPercent = diskList[0].Percent
		diskName = diskList[0].Name
	}

	// Fallback jika tidak ada device terdeteksi dari GetStorageDevices
	if len(diskList) == 0 {
		diskPath := "/DATA"
		if _, err := os.Stat(diskPath); os.IsNotExist(err) {
			diskPath = "/"
		}
		if diskStats, _ := disk.Usage(diskPath); diskStats != nil {
			diskTotal = diskStats.Total
			diskUsed = diskStats.Used
			diskPercent = diskStats.UsedPercent
			diskName = diskPath
			diskList = append(diskList, DiskMetric{
				Name:       diskPath,
				Mountpoint: diskPath,
				Total:      diskTotal,
				Used:       diskUsed,
				Percent:    diskPercent,
			})
		}
	}

	swapPercent := 0.0
	if sMem != nil {
		swapPercent = sMem.UsedPercent
	}

	// Ambil statistik I/O jaringan keseluruhan
	netStats, _ := net.IOCounters(false)
	netSent := uint64(0)
	netRecv := uint64(0)
	if len(netStats) > 0 {
		netSent = netStats[0].BytesSent
		netRecv = netStats[0].BytesRecv
	}

	return &Metrics{
		CPUPercent:  cpuVal,
		CPUTemp:     readCPUTemp(),
		RAMTotal:    vMem.Total,
		RAMUsed:     vMem.Used,
		RAMPercent:  vMem.UsedPercent,
		ZRAMPercent: swapPercent,
		DiskTotal:   diskTotal,
		DiskUsed:    diskUsed,
		DiskPercent: diskPercent,
		DiskName:    diskName,
		Disks:       diskList,
		NetSent:     netSent,
		NetRecv:     netRecv,
	}, nil
}

// readCPUTemp membaca thermal zone SoC Amlogic S905X STB dari sysfs
func readCPUTemp() float64 {
	// Jalur thermal sensor standar kernel Armbian/Linux
	paths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			raw := strings.TrimSpace(string(data))
			if val, err := strconv.ParseFloat(raw, 64); err == nil {
				// Nilai suhu di sysfs umumnya dalam miliderajat Celcius
				if val > 1000 {
					return val / 1000.0
				}
				return val
			}
		}
	}
	return 0.0
}
