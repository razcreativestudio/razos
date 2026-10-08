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

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Metrics merepresentasikan struktur data performa sistem secara komprehensif
type Metrics struct {
	CPUPercent  float64 `json:"cpu_percent"`  // Persentase penggunaan CPU keseluruhan
	CPUTemp     float64 `json:"cpu_temp"`     // Suhu SOC STB dalam derajat Celcius
	RAMTotal    uint64  `json:"ram_total"`    // Total RAM dalam bytes
	RAMUsed     uint64  `json:"ram_used"`     // RAM terpakai dalam bytes
	RAMPercent  float64 `json:"ram_percent"`  // Persentase RAM fisik terpakai
	ZRAMPercent float64 `json:"zram_percent"` // Persentase ZRAM / Swap terpakai
	DiskTotal   uint64  `json:"disk_total"`   // Total penyimpanan data (/DATA)
	DiskUsed    uint64  `json:"disk_used"`    // Penyimpanan terpakai (/DATA)
	DiskPercent float64 `json:"disk_percent"` // Persentase disk terpakai
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

	// Prioritaskan mount point /DATA (MicroSD STB), fallback ke root / jika /DATA belum dimount
	diskPath := "/DATA"
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		diskPath = "/"
	}
	diskStats, _ := disk.Usage(diskPath)

	diskTotal := uint64(0)
	diskUsed := uint64(0)
	diskPercent := 0.0
	if diskStats != nil {
		diskTotal = diskStats.Total
		diskUsed = diskStats.Used
		diskPercent = diskStats.UsedPercent
	}

	swapPercent := 0.0
	if sMem != nil {
		swapPercent = sMem.UsedPercent
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
