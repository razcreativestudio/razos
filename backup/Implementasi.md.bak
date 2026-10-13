# Panduan Implementasi Teknis RAZOS

## 1. Struktur Direktori Proyek

```
razos/
├── .github/workflows/build.yml
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── api/
│   │   ├── docker/
│   │   ├── system/
│   │   ├── terminal/
│   │   └── updater/
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/manifest.json
│   ├── src/
│   │   ├── components/
│   │   ├── layouts/
│   │   │   ├── DesktopLayout.vue
│   │   │   └── MobileLayout.vue
│   │   ├── stores/
│   │   ├── App.vue
│   │   └── main.js
│   ├── package.json
│   └── vite.config.js
├── scripts/
│   ├── install.sh
│   └── clean-host.sh
└── Makefile
```

## 2. Cuplikan Kode Backend: Collector Metrik Sistem (Go)
```go
package system

import (
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	CPUPercent  float64 `json:"cpu_percent"`
	CPUTemp     float64 `json:"cpu_temp"`
	RAMTotal    uint64  `json:"ram_total"`
	RAMUsed     uint64  `json:"ram_used"`
	RAMPercent  float64 `json:"ram_percent"`
	ZRAMPercent float64 `json:"zram_percent"`
	DiskPercent float64 `json:"disk_percent"`
}

func GetMetrics() (*Metrics, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	s, _ := mem.SwapMemory()
	c, _ := cpu.Percent(0, false)
	d, _ := disk.Usage("/DATA")

	cpuVal := 0.0
	if len(c) > 0 {
		cpuVal = c[0]
	}

	return &Metrics{
		CPUPercent:  cpuVal,
		CPUTemp:     readCPUTemp(),
		RAMTotal:    v.Total,
		RAMUsed:     v.Used,
		RAMPercent:  v.UsedPercent,
		ZRAMPercent: s.UsedPercent,
		DiskPercent: d.UsedPercent,
	}, nil
}

func readCPUTemp() float64 {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0
	}
	raw := strings.TrimSpace(string(data))
	val, _ := strconv.ParseFloat(raw, 64)
	return val / 1000.0
}
```

## 3. Cuplikan Kode Backend: RAM Cleaner API

```go
package system

import (
	"os/exec"
)

func CleanRAM() error {
	cmd := exec.Command("sh", "-c", "sync && echo 3 > /proc/sys/vm/drop_caches")
	return cmd.Run()
}
```
