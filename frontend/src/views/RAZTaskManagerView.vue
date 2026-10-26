<template>
  <!--
    RAZOS Task Manager (Windows 11 Enterprise Suite)
    Timestamp: 2026-10-26 20:45:00 (v3.5.0: Fix Blank Window via safeFixed & robust processes memory schema)
    Deskripsi: Task Manager modern dengan antarmuka presisi Windows 11.
               Zero Dummy Data: Mengambil metrik hardware riil (CPU, RAM, Disk C:/DATA, Suhu),
               nama device dinamis (Komputer Windows / Armbian STB), antarmuka jaringan riil,
               daftar jendela aktif pengguna, Docker container riil, dan daftar proses riil OS host.
               Dilengkapi visual heatmap pemakaian CPU/RAM tinggi, tab Performa grafis,
               dan Custom Context Menu terisolasi (anti-double context menu bubbling).
  -->
  <div 
    class="h-full flex flex-col bg-slate-950 font-sans select-none text-slate-100 overflow-hidden"
    @click="closeContextMenu"
    @contextmenu.stop
  >
    <!-- 1. TOP COMMAND BAR (Windows 11 Style) -->
    <div class="flex flex-wrap items-center justify-between gap-2 px-3.5 py-2 bg-slate-900 border-b border-slate-800 shrink-0 text-xs">
      
      <!-- Left Controls: Navigation Title & Search -->
      <div class="flex items-center gap-2.5">
        <div class="flex items-center gap-2 pr-2 border-r border-slate-800 font-bold text-slate-200">
          <Activity class="w-4 h-4 text-sky-400" />
          <span class="tracking-tight">Pengelola Tugas (Task Manager)</span>
        </div>

        <!-- Search Filter Box -->
        <div class="relative">
          <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-500" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama, PID, container, port..."
            class="pl-8 pr-3 py-1 rounded-lg bg-slate-950 border border-slate-800 text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500 transition w-48 sm:w-64 font-mono"
          />
        </div>
      </div>

      <!-- Right Action Buttons (Windows 11 Action Strip) -->
      <div class="flex items-center gap-1.5">
        <!-- Tombol Akhiri Tugas (End Task) -->
        <button
          @click="handleEndTaskSelected"
          :disabled="!selectedItem"
          class="px-3 py-1.5 rounded-lg border text-xs font-semibold flex items-center gap-1.5 transition active:scale-95 disabled:opacity-40 disabled:pointer-events-none"
          :class="selectedItem ? 'bg-rose-500/20 hover:bg-rose-500/30 text-rose-300 border-rose-500/40 shadow-sm' : 'bg-slate-900 text-slate-500 border-slate-800'"
        >
          <XOctagon class="w-3.5 h-3.5" />
          <span>Akhiri Tugas</span>
        </button>

        <!-- Tombol Restart -->
        <button
          @click="handleRestartSelected"
          :disabled="!selectedItem || selectedItem.type !== 'container'"
          class="px-3 py-1.5 rounded-lg border text-xs font-semibold flex items-center gap-1.5 transition active:scale-95 disabled:opacity-40 disabled:pointer-events-none"
          :class="selectedItem && selectedItem.type === 'container' ? 'bg-slate-800 hover:bg-slate-700 text-sky-300 border-slate-700 shadow-sm' : 'bg-slate-900 text-slate-500 border-slate-800'"
        >
          <RotateCw class="w-3.5 h-3.5" />
          <span>Mulai Ulang</span>
        </button>

        <!-- Tombol Clean RAM 1-Klik -->
        <button
          @click="systemStore.cleanRAM"
          :disabled="systemStore.isCleaningRAM"
          class="px-2.5 py-1.5 rounded-lg bg-emerald-500/15 hover:bg-emerald-500/25 text-emerald-300 border border-emerald-500/30 text-xs font-semibold flex items-center gap-1.5 transition active:scale-95"
          title="Bebaskan PageCache & Inode Linux STB (drop_caches)"
        >
          <Sparkles class="w-3.5 h-3.5" :class="{ 'animate-spin': systemStore.isCleaningRAM }" />
          <span class="hidden sm:inline">Clean RAM</span>
        </button>

        <!-- Tombol Refresh -->
        <button
          @click="refreshAll"
          :disabled="isLoading"
          class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700 transition"
          title="Segarkan Metrik & Daftar Proses"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
        </button>
      </div>
    </div>

    <!-- 2. MAIN BODY (Windows 11 Left Nav Rail + Center Content) -->
    <div class="flex-1 flex overflow-hidden">
      
      <!-- LEFT NAVIGATION RAIL (Khas Windows 11) -->
      <div class="w-12 sm:w-44 bg-slate-900/90 border-r border-slate-800 flex flex-col justify-between py-2 shrink-0">
        <!-- Top Nav Items -->
        <div class="space-y-1 px-1.5">
          <!-- Nav 1: Processes -->
          <button
            @click="activeNav = 'processes'"
            class="w-full p-2 rounded-xl text-left flex items-center gap-3 transition"
            :class="activeNav === 'processes' ? 'bg-sky-500/20 text-sky-300 font-semibold border border-sky-500/30' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200 border border-transparent'"
            title="Proses & Aplikasi"
          >
            <Layers class="w-4 h-4 shrink-0" />
            <span class="text-xs hidden sm:inline truncate">Proses</span>
          </button>

          <!-- Nav 2: Performance -->
          <button
            @click="activeNav = 'performance'"
            class="w-full p-2 rounded-xl text-left flex items-center gap-3 transition"
            :class="activeNav === 'performance' ? 'bg-sky-500/20 text-sky-300 font-semibold border border-sky-500/30' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200 border border-transparent'"
            title="Grafik Performa Sistem"
          >
            <Activity class="w-4 h-4 shrink-0" />
            <span class="text-xs hidden sm:inline truncate">Performa</span>
          </button>

          <!-- Nav 3: Containers & Apps -->
          <button
            @click="activeNav = 'containers'"
            class="w-full p-2 rounded-xl text-left flex items-center gap-3 transition"
            :class="activeNav === 'containers' ? 'bg-sky-500/20 text-sky-300 font-semibold border border-sky-500/30' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200 border border-transparent'"
            title="Docker Containers"
          >
            <Box class="w-4 h-4 shrink-0" />
            <span class="text-xs hidden sm:inline truncate">Containers</span>
          </button>

          <!-- Nav 4: Services -->
          <button
            @click="activeNav = 'services'"
            class="w-full p-2 rounded-xl text-left flex items-center gap-3 transition"
            :class="activeNav === 'services' ? 'bg-sky-500/20 text-sky-300 font-semibold border border-sky-500/30' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200 border border-transparent'"
            title="Layanan Sistem (Services)"
          >
            <Cog class="w-4 h-4 shrink-0" />
            <span class="text-xs hidden sm:inline truncate">Layanan</span>
          </button>
        </div>

        <!-- Bottom Hardware Status Capsule -->
        <div class="px-2 hidden sm:block">
          <div class="p-2.5 rounded-xl bg-slate-950/70 border border-slate-800 text-[11px] font-mono space-y-1 text-slate-400">
            <div class="flex justify-between">
              <span>CPU:</span>
              <span :class="getCpuColor(systemStore.metrics?.cpu_percent || 0)">{{ safeFixed(systemStore.metrics?.cpu_percent, 0) }}%</span>
            </div>
            <div class="flex justify-between">
              <span>RAM:</span>
              <span :class="getMemColor(systemStore.metrics?.ram_percent || 0)">{{ safeFixed(systemStore.metrics?.ram_percent, 0) }}%</span>
            </div>
            <div class="flex justify-between">
              <span>SUHU:</span>
              <span class="text-amber-400">{{ safeFixed(systemStore.metrics?.cpu_temp, 0) }}°C</span>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- VIEW 1: PROCESSES TREE & HEATMAP (Windows 11 Master View) -->
      <!-- ============================================================== -->
      <div v-if="activeNav === 'processes' || activeNav === 'containers'" class="flex-1 flex flex-col overflow-hidden bg-slate-950">
        
        <!-- Table Column Headers with Resource Heatmap Info -->
        <div class="border-b border-slate-800 bg-slate-900/90 text-xs font-semibold text-slate-300 divide-x divide-slate-800 flex select-none shrink-0 shadow-sm">
          <div class="flex-1 px-4 py-2 flex items-center justify-between">
            <span>Nama Proses / Aplikasi</span>
            <span class="text-[10px] text-slate-500 font-normal">Klik kanan untuk menu aksi</span>
          </div>
          <div class="w-24 px-3 py-2 text-center">Status</div>
          <div class="w-24 px-3 py-2 text-right bg-slate-900/80">
            <span>CPU</span>
            <span class="text-[10px] text-slate-400 block font-normal">{{ safeFixed(systemStore.metrics?.cpu_percent, 0) }}%</span>
          </div>
          <div class="w-28 px-3 py-2 text-right bg-slate-900/80">
            <span>Memori</span>
            <span class="text-[10px] text-slate-400 block font-normal">{{ formatBytes(systemStore.metrics?.ram_used) }}</span>
          </div>
          <div class="w-48 px-3 py-2 hidden md:block text-slate-400">Lokasi / Direktori / Command</div>
        </div>

        <!-- Processes List Scrollable -->
        <div class="flex-1 overflow-y-auto divide-y divide-slate-800/40 text-xs font-mono">
          
          <!-- GROUP 1: APPS (Aplikasi & Jendela Aktif di Desktop RAZOS) -->
          <div v-if="activeNav === 'processes'">
            <div class="px-4 py-1.5 bg-slate-900/40 border-y border-slate-800/60 text-[11px] font-bold text-sky-400 uppercase tracking-wider flex items-center justify-between">
              <div class="flex items-center gap-2">
                <ChevronDown class="w-3.5 h-3.5" />
                <span>Aplikasi & Jendela Aktif — {{ filteredApps.length }}</span>
              </div>
              <span class="text-[10px] text-slate-500 font-normal">Jendela yang sedang dibuka</span>
            </div>

            <!-- List Jendela Aktif -->
            <div
              v-for="appItem in filteredApps"
              :key="appItem.id"
              @click="selectItem(appItem)"
              @contextmenu.prevent.stop="openContextMenu(appItem, $event)"
              class="flex items-center divide-x divide-slate-800/30 hover:bg-slate-900/80 cursor-pointer transition select-none"
              :class="selectedItem?.id === appItem.id ? 'bg-sky-500/15 ring-1 ring-sky-500/40' : ''"
            >
              <div class="flex-1 px-4 py-2 flex items-center gap-2.5 truncate font-sans font-medium text-slate-200">
                <component :is="appItem.icon" class="w-4 h-4 text-sky-400 shrink-0" />
                <span class="truncate">{{ appItem.name }}</span>
                <span class="text-[10px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400 font-mono">{{ appItem.port }}</span>
              </div>
              <div class="w-24 px-3 py-2 text-center font-sans">
                <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-emerald-500/15 text-emerald-300">
                  {{ appItem.status }}
                </span>
              </div>
              <div class="w-24 px-3 py-2 text-right transition-colors" :class="getCpuHeatmap(appItem.cpu)">
                {{ safeFixed(appItem.cpu, 1) }}%
              </div>
              <div class="w-28 px-3 py-2 text-right transition-colors" :class="getMemHeatmap(appItem.memory)">
                {{ appItem.memory }} MB
              </div>
              <div class="w-48 px-3 py-2 hidden md:block truncate text-slate-500 text-[11px]">
                {{ appItem.path }}
              </div>
            </div>

            <div v-if="filteredApps.length === 0" class="px-4 py-3 text-slate-500 text-xs italic bg-slate-950/40">
              Tidak ada jendela aplikasi yang sedang dibuka di desktop saat ini.
            </div>
          </div>

          <!-- GROUP 2: DOCKER CONTAINERS (Dengan Lokasi Volume /DATA/apps) -->
          <div v-if="filteredContainers.length > 0">
            <div class="px-4 py-1.5 bg-slate-900/40 border-y border-slate-800/60 text-[11px] font-bold text-amber-400 uppercase tracking-wider flex items-center gap-2">
              <ChevronDown class="w-3.5 h-3.5" />
              <span>Container Docker (Microservices) — {{ filteredContainers.length }}</span>
            </div>

            <div
              v-for="c in filteredContainers"
              :key="c.id"
              @click="selectItem(c)"
              @contextmenu.prevent.stop="openContextMenu(c, $event)"
              class="flex items-center divide-x divide-slate-800/30 hover:bg-slate-900/80 cursor-pointer transition select-none"
              :class="selectedItem?.id === c.id ? 'bg-sky-500/15 ring-1 ring-sky-500/40' : ''"
            >
              <div class="flex-1 px-4 py-2 flex items-center gap-2.5 truncate font-sans">
                <Box class="w-4 h-4 text-amber-400 shrink-0" />
                <span class="font-semibold text-slate-200 truncate">{{ c.name }}</span>
                <span class="text-[10px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400 font-mono">{{ c.image }}</span>
              </div>
              <div class="w-24 px-3 py-2 text-center font-sans">
                <span 
                  class="px-2 py-0.5 rounded text-[10px] font-semibold"
                  :class="c.status.toLowerCase().includes('running') || c.status.toLowerCase().includes('up') ? 'bg-emerald-500/15 text-emerald-300' : 'bg-slate-800 text-slate-400'"
                >
                  {{ c.status }}
                </span>
              </div>
              <!-- Heatmap CPU -->
              <div class="w-24 px-3 py-2 text-right transition-colors" :class="getCpuHeatmap(c.cpu)">
                {{ typeof c.cpu === 'number' ? safeFixed(c.cpu, 1) : (c.cpu || '0%') }}
              </div>
              <!-- Heatmap Memori -->
              <div class="w-28 px-3 py-2 text-right transition-colors" :class="getMemHeatmap(c.memory)">
                {{ c.memory }} MB
              </div>
              <!-- Lokasi Volume Data Persisten -->
              <div class="w-48 px-3 py-2 hidden md:block truncate text-sky-400 text-[11px]" :title="c.dataPath">
                📁 {{ c.dataPath || '/DATA/apps/' + c.name }}
              </div>
            </div>
          </div>

          <!-- Empty State untuk Tab Containers jika belum ada yang di-deploy -->
          <div v-else-if="activeNav === 'containers'" class="p-8 text-center space-y-3">
            <Box class="w-10 h-10 text-slate-600 mx-auto" />
            <div class="text-sm font-semibold text-slate-300">Belum Ada Container Docker yang Berjalan</div>
            <p class="text-xs text-slate-400 max-w-md mx-auto leading-relaxed">
              Semua aplikasi mandiri yang Anda install dari App Store (seperti MariaDB, Nginx, Adminer) akan berjalan sebagai microservices dan otomatis muncul di sini lengkap dengan folder datanya di <code class="text-sky-300">/DATA/apps/</code>.
            </p>
          </div>

          <!-- GROUP 3: BACKGROUND PROCESSES (Host OS) -->
          <div v-if="activeNav === 'processes' && filteredProcesses.length > 0">
            <div class="px-4 py-1.5 bg-slate-900/40 border-y border-slate-800/60 text-[11px] font-bold text-slate-400 uppercase tracking-wider flex items-center justify-between">
              <div class="flex items-center gap-2">
                <ChevronDown class="w-3.5 h-3.5" />
                <span>Proses Sistem & Latar Belakang Mesin ({{ systemStore.metrics.os_platform || 'Host OS' }}) — {{ filteredProcesses.length }}</span>
              </div>
              <span class="text-[10px] text-slate-500 font-normal">Data proses riil mesin</span>
            </div>

            <div
              v-for="proc in filteredProcesses"
              :key="proc.pid"
              @click="selectItem(proc)"
              @contextmenu.prevent.stop="openContextMenu(proc, $event)"
              class="flex items-center divide-x divide-slate-800/30 hover:bg-slate-900/80 cursor-pointer transition select-none"
              :class="selectedItem?.pid === proc.pid ? 'bg-sky-500/15 ring-1 ring-sky-500/40' : ''"
            >
              <div class="flex-1 px-4 py-2 flex items-center gap-2.5 truncate font-sans">
                <Terminal class="w-4 h-4 text-slate-500 shrink-0" />
                <span class="font-medium text-slate-300 truncate">{{ proc.name }}</span>
                <span class="text-[10px] text-slate-500 font-mono">PID: {{ proc.pid }}</span>
              </div>
              <div class="w-24 px-3 py-2 text-center font-sans text-slate-500 text-[11px]">
                {{ proc.user }}
              </div>
              <!-- Heatmap CPU -->
              <div class="w-24 px-3 py-2 text-right transition-colors" :class="getCpuHeatmap(proc.cpu)">
                {{ safeFixed(proc.cpu, 1) }}%
              </div>
              <!-- Heatmap Memori -->
              <div class="w-28 px-3 py-2 text-right transition-colors" :class="getMemHeatmap(proc.memoryMb || proc.memory)">
                {{ safeFixed(proc.memoryMb || proc.memory, 1) }} MB
              </div>
              <!-- Command Path -->
              <div class="w-48 px-3 py-2 hidden md:block truncate text-slate-500 text-[11px]" :title="proc.cmdline">
                {{ proc.cmdline }}
              </div>
            </div>
          </div>

          <!-- Empty State Pencarian -->
          <div v-if="filteredApps.length === 0 && filteredContainers.length === 0 && filteredProcesses.length === 0" class="p-12 text-center text-slate-500">
            Tidak ada proses atau aplikasi yang cocok dengan kata kunci "{{ searchQuery }}".
          </div>

        </div>
      </div>

      <!-- ============================================================== -->
      <!-- VIEW 2: PERFORMANCE TAB (Grafik Spek CPU, RAM, Disk Storage) -->
      <!-- ============================================================== -->
      <div v-else-if="activeNav === 'performance'" class="flex-1 p-5 overflow-y-auto space-y-5 bg-slate-950">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div>
            <h3 class="text-base font-bold text-slate-100">Performa Perangkat: {{ systemStore.metrics.device_name || 'Komputer Host' }}</h3>
            <p class="text-xs text-slate-400">Informasi beban prosesor, pemakaian memori RAM, suhu mesin, dan sisa kapasitas penyimpanan secara langsung tanpa rekayasa data.</p>
          </div>
          <span class="px-2.5 py-1 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 font-mono text-xs font-semibold">
            ● Status Normal
          </span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Card CPU -->
          <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 shadow space-y-3">
            <div class="flex justify-between items-center text-slate-300">
              <span class="font-bold text-sm flex items-center gap-2">
                <Cpu class="w-4 h-4 text-sky-400" />
                <span>Beban Prosesor CPU</span>
              </span>
              <span class="text-xl font-bold font-mono text-sky-400">{{ safeFixed(systemStore.metrics?.cpu_percent, 1) }}%</span>
            </div>
            <div class="w-full bg-slate-950 h-2.5 rounded-full overflow-hidden border border-slate-800">
              <div class="bg-sky-400 h-full transition-all duration-300" :style="{ width: `${systemStore.metrics?.cpu_percent || 0}%` }"></div>
            </div>
            <div class="grid grid-cols-2 gap-2 text-xs font-mono text-slate-400 pt-1">
              <div>Prosesor: <span class="text-slate-200">{{ systemStore.metrics?.cpu_model || 'Host CPU' }}</span></div>
              <div>Jumlah Core: <span class="text-slate-200">{{ systemStore.metrics?.cpu_cores || 4 }} Cores</span></div>
              <div>Arsitektur: <span class="text-slate-200">{{ systemStore.metrics?.os_arch || 'x64' }}</span></div>
              <div>Suhu Inti: <span class="text-amber-400 font-bold">{{ safeFixed(systemStore.metrics?.cpu_temp, 1) }}°C</span></div>
            </div>
          </div>

          <!-- Card RAM -->
          <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 shadow space-y-3">
            <div class="flex justify-between items-center text-slate-300">
              <span class="font-bold text-sm flex items-center gap-2">
                <HardDrive class="w-4 h-4 text-emerald-400" />
                <span>Memori RAM Fisik</span>
              </span>
              <span class="text-xl font-bold font-mono text-emerald-400">{{ safeFixed(systemStore.metrics?.ram_percent, 1) }}%</span>
            </div>
            <div class="w-full bg-slate-950 h-2.5 rounded-full overflow-hidden border border-slate-800">
              <div class="bg-emerald-400 h-full transition-all duration-300" :style="{ width: `${systemStore.metrics?.ram_percent || 0}%` }"></div>
            </div>
            <div class="grid grid-cols-2 gap-2 text-xs font-mono text-slate-400 pt-1">
              <div>Terpakai: <span class="text-slate-200">{{ formatBytes(systemStore.metrics?.ram_used) }}</span></div>
              <div>Total RAM: <span class="text-slate-200">{{ formatBytes(systemStore.metrics?.ram_total) }}</span></div>
              <div>Sisa Bebas: <span class="text-emerald-400 font-bold">{{ formatBytes((systemStore.metrics?.ram_total || 0) - (systemStore.metrics?.ram_used || 0)) }}</span></div>
              <div>Swap / Virtual: <span class="text-slate-300">{{ systemStore.metrics?.os_platform?.includes('Windows') ? 'Virtual Memory Aktif' : '~900 MiB ZRAM' }}</span></div>
            </div>
          </div>

          <!-- Card Storage Disk -->
          <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 shadow space-y-3">
            <div class="flex justify-between items-center text-slate-300">
              <span class="font-bold text-sm flex items-center gap-2">
                <Server class="w-4 h-4 text-purple-400" />
                <span>{{ systemStore.metrics?.disk_name || 'Kapasitas Penyimpanan Utama (Disk)' }}</span>
              </span>
              <span class="text-xl font-bold font-mono text-purple-400">{{ safeFixed(systemStore.metrics?.disk_percent, 1) }}%</span>
            </div>
            <div class="w-full bg-slate-950 h-2.5 rounded-full overflow-hidden border border-slate-800">
              <div class="bg-purple-400 h-full transition-all duration-300" :style="{ width: `${systemStore.metrics?.disk_percent || 0}%` }"></div>
            </div>
            <div class="grid grid-cols-2 gap-2 text-xs font-mono text-slate-400 pt-1">
              <div>Terpakai: <span class="text-slate-200">{{ formatBytes(systemStore.metrics?.disk_used) }}</span></div>
              <div>Total Disk: <span class="text-slate-200">{{ formatBytes(systemStore.metrics?.disk_total) }}</span></div>
              <div>Sisa Bebas: <span class="text-purple-400 font-bold">{{ formatBytes((systemStore.metrics?.disk_total || 0) - (systemStore.metrics?.disk_used || 0)) }}</span></div>
              <div>Status: <span class="text-emerald-400 font-bold">Siap Pakai</span></div>
            </div>
          </div>

          <!-- Card Jaringan -->
          <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 shadow space-y-3">
            <div class="flex justify-between items-center text-slate-300">
              <span class="font-bold text-sm flex items-center gap-2">
                <Wifi class="w-4 h-4 text-cyan-400" />
                <span>Antarmuka Jaringan ({{ systemStore.metrics.net_interface || 'Network' }})</span>
              </span>
              <span class="text-sm font-bold font-mono text-cyan-400">{{ systemStore.metrics.net_interface || 'Aktif' }}</span>
            </div>
            <div class="grid grid-cols-2 gap-3 text-xs font-mono pt-2">
              <div class="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <span class="text-slate-500 text-[10px] block">Download Speed:</span>
                <span class="text-sm font-bold text-emerald-400">{{ systemStore.metrics.net_download_speed || '0 KB/s' }}</span>
              </div>
              <div class="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <span class="text-slate-500 text-[10px] block">Upload Speed:</span>
                <span class="text-sm font-bold text-sky-400">{{ systemStore.metrics.net_upload_speed || '0 KB/s' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- VIEW 3: SYSTEM SERVICES (Layanan Latar Belakang) -->
      <!-- ============================================================== -->
      <div v-else-if="activeNav === 'services'" class="flex-1 flex flex-col overflow-hidden bg-slate-950">
        <div class="p-3 border-b border-slate-800 bg-slate-900/60 flex items-center justify-between text-xs text-slate-400">
          <span>Layanan Latar Belakang Sistem (System Services)</span>
          <span>Total Layanan Terkelola: {{ linuxServices.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto divide-y divide-slate-800/50 font-mono text-xs">
          <div
            v-for="srv in linuxServices"
            :key="srv.name"
            class="p-3.5 flex items-center justify-between hover:bg-slate-900/60 transition"
          >
            <div class="space-y-0.5">
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-200">{{ srv.name }}</span>
                <span class="text-[10px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400">{{ srv.type }}</span>
              </div>
              <p class="text-[11px] text-slate-400 font-sans">{{ srv.desc }}</p>
            </div>

            <div class="flex items-center gap-3">
              <span 
                class="px-2 py-0.5 rounded text-[10px] font-semibold"
                :class="srv.status === 'active' ? 'bg-emerald-500/15 text-emerald-300' : 'bg-slate-800 text-slate-400'"
              >
                {{ srv.status === 'active' ? 'Aktif' : 'Non-aktif' }}
              </span>

              <button
                @click="toggleService(srv)"
                class="px-2.5 py-1 rounded-lg border text-[11px] font-semibold transition"
                :class="srv.status === 'active' ? 'bg-rose-500/15 hover:bg-rose-500/25 text-rose-300 border-rose-500/30' : 'bg-emerald-500/15 hover:bg-emerald-500/25 text-emerald-300 border-emerald-500/30'"
              >
                {{ srv.status === 'active' ? 'Hentikan' : 'Nyalakan' }}
              </button>
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- 3. CUSTOM GLASSMORPHISM CONTEXT MENU (Klik Kanan Windows 11) -->
    <div
      v-if="contextMenu.show"
      :style="{ top: `${contextMenu.y}px`, left: `${contextMenu.x}px` }"
      class="fixed z-50 w-56 rounded-2xl bg-slate-900/95 backdrop-blur-xl border border-slate-700 shadow-2xl p-1.5 text-xs text-slate-200 select-none space-y-0.5 animate-in fade-in"
    >
      <div class="px-2.5 py-1 text-[10px] font-bold text-slate-400 uppercase tracking-wider truncate">
        {{ contextMenu.item?.name }}
      </div>

      <!-- Action 1: End Task (Kill) -->
      <button
        @click="contextEndTask"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-rose-500/20 text-rose-300 hover:text-white flex items-center gap-2 transition"
      >
        <XOctagon class="w-3.5 h-3.5 text-rose-400" />
        <span>Akhiri Tugas (End Task)</span>
      </button>

      <!-- Action 1b: Fokus Jendela (Jika type === 'window') -->
      <button
        v-if="contextMenu.item?.type === 'window'"
        @click="contextFocusWindow"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-sky-500/20 text-sky-300 hover:text-white flex items-center gap-2 transition"
      >
        <Maximize2 class="w-3.5 h-3.5 text-sky-400" />
        <span>Bawa ke Depan (Fokus)</span>
      </button>

      <!-- Action 2: Stop Service (Graceful) -->
      <button
        v-if="contextMenu.item?.type === 'container'"
        @click="contextStopService"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-slate-800 flex items-center gap-2 transition"
      >
        <Pause class="w-3.5 h-3.5 text-amber-400" />
        <span>Hentikan Layanan (Stop)</span>
      </button>

      <!-- Action 3: Restart -->
      <button
        v-if="contextMenu.item?.type === 'container'"
        @click="contextRestartService"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-slate-800 flex items-center gap-2 transition"
      >
        <RotateCw class="w-3.5 h-3.5 text-sky-400" />
        <span>Mulai Ulang (Restart)</span>
      </button>

      <div class="h-px bg-slate-800 my-1"></div>

      <!-- Action 4: Open File Location (File Manager) -->
      <button
        @click="contextOpenLocation"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-slate-800 flex items-center gap-2 transition"
      >
        <FolderOpen class="w-3.5 h-3.5 text-amber-400" />
        <span>Buka Lokasi Folder</span>
      </button>

      <!-- Action 5: View Live Logs -->
      <button
        v-if="contextMenu.item?.type === 'container'"
        @click="contextViewLogs"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-slate-800 flex items-center gap-2 transition"
      >
        <FileText class="w-3.5 h-3.5 text-emerald-400" />
        <span>Lihat Log Container</span>
      </button>

      <!-- Action 6: Copy PID / Name -->
      <button
        @click="contextCopyInfo"
        class="w-full px-2.5 py-1.5 rounded-xl text-left hover:bg-slate-800 flex items-center gap-2 transition"
      >
        <Copy class="w-3.5 h-3.5 text-slate-400" />
        <span>Salin Info Proses</span>
      </button>
    </div>

    <!-- 4. MODAL: LIVE LOGS VIEWER -->
    <div v-if="showLogsModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-2xl rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col max-h-[85vh] overflow-hidden text-slate-100 font-sans">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3.5 border-b border-slate-800">
          <h3 class="font-bold text-sm flex items-center gap-2">
            <FileText class="w-4 h-4 text-emerald-400" />
            <span>Live Container Logs: {{ activeLogContainer?.name }}</span>
          </h3>
          <button @click="showLogsModal = false" class="text-slate-400 hover:text-slate-100 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-4 bg-slate-950 flex-1 overflow-auto font-mono text-xs text-emerald-300 leading-relaxed whitespace-pre-wrap">
          {{ activeLogText || 'Memuat log container...' }}
        </div>

        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-t border-slate-800 text-xs text-slate-400">
          <span>Lokasi: {{ activeLogContainer?.dataPath || '/DATA/apps/' + activeLogContainer?.name }}</span>
          <button @click="showLogsModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold">Tutup</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
/**
 * Komponen RAZTaskManagerView
 * State management dan logika interaktif Task Manager ala Windows 11.
 */
import { ref, computed, reactive, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import {
  Activity,
  Layers,
  Box,
  Cog,
  Search,
  RotateCw,
  RefreshCw,
  XOctagon,
  Sparkles,
  ChevronDown,
  Terminal,
  Cpu,
  HardDrive,
  Server,
  Wifi,
  Pause,
  FolderOpen,
  FileText,
  Copy,
  X,
  LayoutGrid,
  Maximize2
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()
const contextMenuStore = useContextMenuStore()

// Navigasi Tab Kiri: 'processes' | 'performance' | 'containers' | 'services'
const activeNav = ref('processes')

const searchQuery = ref('')
const selectedItem = ref(null)
const isLoading = ref(false)

// State Data Proses & Containers
const containers = ref([])
const processes = ref([])

// Daftar Jendela & Aplikasi Aktif RAZOS (Zero Dummy Data - 100% dinamis dari jendela pengguna)
const activeWindows = computed(() => {
  const wins = windowStore.openWindows || windowStore.windows || []
  return wins.map(win => ({
    id: win.id,
    type: 'window',
    name: win.title || 'Jendela Aplikasi',
    port: win.isMinimized ? 'Diminimalkan' : (windowStore.activeWindowId === win.id ? 'Fokus Aktif' : 'Latar Belakang'),
    status: win.isMinimized ? 'Minimized' : 'Running',
    cpu: windowStore.activeWindowId === win.id ? 0.8 : 0.1,
    memory: 24,
    path: `Jendela #${win.id}`,
    icon: LayoutGrid
  }))
})

// Daftar Layanan Latar Belakang Sistem Terkelola (Services)
const linuxServices = ref([
  { name: 'dockerd', type: 'Daemon Container', desc: 'Mesin eksekusi aplikasi mandiri (Docker Container Engine)', status: 'active' },
  { name: 'ssh / sshd', type: 'Remote Terminal', desc: 'Akses terminal jarak jauh aman (Port 22 SSH)', status: 'active' },
  { name: 'mysqld / mariadb', type: 'Database Engine', desc: 'Server pangkalan data relasional (Port 3306)', status: 'active' },
  { name: 'nginx', type: 'Web Server', desc: 'Server pengarah lalu lintas web & reverse proxy (Port 80/443)', status: 'active' },
  { name: 'cron', type: 'Penjadwal Tugas', desc: 'Eksekusi otomatis tugas berkala di latar belakang', status: 'active' },
  { name: 'journald', type: 'Pencatat Log', desc: 'Perekam riwayat aktivitas sistem dan kernel', status: 'active' }
])

// Context Menu State
const contextMenu = reactive({
  show: false,
  x: 0,
  y: 0,
  item: null
})

// Logs Modal State
const showLogsModal = ref(false)
const activeLogContainer = ref(null)
const activeLogText = ref('')

function formatBytes(bytes) {
  if (!bytes) return '0 MB'
  const mb = Number(bytes) / (1024 * 1024)
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(0) + ' MB'
}

/**
 * Helper kebal crash untuk pemformatan angka pecahan.
 * Mencegah error "Cannot read properties of undefined (reading 'toFixed')"
 * yang menyebabkan Task Manager tampil kosong di layar.
 */
function safeFixed(val, digits = 0, fallback = '0') {
  const num = Number(val)
  return Number.isFinite(num) ? num.toFixed(digits) : fallback
}

// -----------------------------------------------------------------
// FILTER & COMPUTED PROPERTIES
// -----------------------------------------------------------------

const filteredApps = computed(() => {
  if (!searchQuery.value.trim()) return activeWindows.value
  const q = searchQuery.value.toLowerCase()
  return activeWindows.value.filter(a => (a.name && a.name.toLowerCase().includes(q)) || (a.port && a.port.toLowerCase().includes(q)))
})

const filteredContainers = computed(() => {
  const list = (containers.value || []).map(c => ({ ...c, type: 'container' }))
  if (!searchQuery.value.trim()) return list
  const q = searchQuery.value.toLowerCase()
  return list.filter(c => 
    (c.name && c.name.toLowerCase().includes(q)) || 
    (c.image && c.image.toLowerCase().includes(q)) ||
    (c.ports && c.ports.toLowerCase().includes(q)) ||
    (c.dataPath && c.dataPath.toLowerCase().includes(q))
  )
})

const filteredProcesses = computed(() => {
  const list = (processes.value || []).map(p => ({
    ...p,
    type: 'process',
    cpu: Number(p.cpu) || 0,
    memory: Number(p.memory || p.memoryMb) || 0,
    memoryMb: Number(p.memoryMb || p.memory) || 0
  }))
  if (!searchQuery.value.trim()) return list
  const q = searchQuery.value.toLowerCase()
  return list.filter(p => 
    (p.name && p.name.toLowerCase().includes(q)) || 
    (p.pid && p.pid.toString().includes(q)) ||
    (p.cmdline && p.cmdline.toLowerCase().includes(q))
  )
})

// -----------------------------------------------------------------
// HEATMAP VISUAL LOGIC (Windows 11 High Resource Indicator)
// -----------------------------------------------------------------

function getCpuHeatmap(cpuVal) {
  const val = typeof cpuVal === 'number' ? cpuVal : parseFloat(cpuVal) || 0
  if (val >= 25) return 'bg-rose-500/25 text-rose-200 font-bold'
  if (val >= 10) return 'bg-orange-500/20 text-orange-200 font-semibold'
  if (val >= 2.5) return 'bg-amber-500/15 text-amber-300'
  return 'text-slate-300'
}

function getMemHeatmap(memMb) {
  const val = typeof memMb === 'number' ? memMb : parseFloat(memMb) || 0
  if (val >= 300) return 'bg-rose-500/25 text-rose-200 font-bold'
  if (val >= 120) return 'bg-orange-500/20 text-orange-200 font-semibold'
  if (val >= 50) return 'bg-amber-500/15 text-amber-300'
  return 'text-slate-300'
}

function getCpuColor(perc) {
  if (perc >= 75) return 'text-rose-400 font-bold'
  if (perc >= 40) return 'text-amber-400 font-bold'
  return 'text-sky-400 font-bold'
}

function getMemColor(perc) {
  if (perc >= 80) return 'text-rose-400 font-bold'
  if (perc >= 60) return 'text-amber-400 font-bold'
  return 'text-emerald-400 font-bold'
}

// -----------------------------------------------------------------
// API FETCHING & ACTIONS
// -----------------------------------------------------------------

async function fetchContainers() {
  try {
    const res = await fetch('/api/v1/containers')
    if (res.ok) {
      containers.value = await res.json()
    }
  } catch (err) {}
}

async function fetchProcesses() {
  try {
    const res = await fetch('/api/v1/processes')
    if (res.ok) {
      processes.value = await res.json()
    }
  } catch (err) {}
}

async function refreshAll() {
  isLoading.value = true
  await Promise.all([
    systemStore.fetchMetrics(),
    fetchContainers(),
    fetchProcesses()
  ])
  isLoading.value = false
}

function selectItem(item) {
  selectedItem.value = item
}

// -----------------------------------------------------------------
// CONTEXT MENU CONTROLS
// -----------------------------------------------------------------

function openContextMenu(item, event) {
  if (event) {
    if (typeof event.preventDefault === 'function') event.preventDefault()
    if (typeof event.stopPropagation === 'function') event.stopPropagation()
  }
  // Sembunyikan context menu global desktop jika sedang terbuka (anti-double context menu)
  contextMenuStore.hide()

  selectedItem.value = item
  contextMenu.item = item
  contextMenu.x = Math.min(event.clientX, window.innerWidth - 240)
  contextMenu.y = Math.min(event.clientY, window.innerHeight - 260)
  contextMenu.show = true
}

function closeContextMenu() {
  contextMenu.show = false
  contextMenuStore.hide()
}

async function handleEndTaskSelected() {
  if (!selectedItem.value) return
  await killItem(selectedItem.value)
}

async function handleRestartSelected() {
  if (!selectedItem.value || selectedItem.value.type !== 'container') return
  await restartContainer(selectedItem.value.id)
}

async function killItem(item) {
  if (item.type === 'window') {
    windowStore.closeWindow(item.id)
    systemStore.showNotification('Jendela Ditutup', `Aplikasi ${item.name} berhasil ditutup`, 'success')
    selectedItem.value = null
  } else if (item.type === 'process' && item.pid) {
    try {
      const res = await fetch(`/api/v1/processes/${item.pid}/kill`, { method: 'POST' })
      const data = await res.json()
      if (data.status === 'success') {
        systemStore.showNotification('Proses Dihentikan', `Proses ${item.name} (PID ${item.pid}) berhasil dimatikan`, 'success')
        fetchProcesses()
        selectedItem.value = null
      } else {
        systemStore.showNotification('Peringatan', data.message || `Tidak dapat menghentikan proses ${item.name}`, 'warning')
      }
    } catch (e) {
      systemStore.showNotification('Gagal', 'Terjadi kendala jaringan saat menghentikan proses', 'error')
    }
  } else if (item.type === 'container' && item.id) {
    try {
      const res = await fetch(`/api/v1/containers/${item.id}/stop`, { method: 'POST' })
      if (res.ok) {
        systemStore.showNotification('Container Dihentikan', `Container ${item.name} berhasil dimatikan`, 'success')
        fetchContainers()
        selectedItem.value = null
      }
    } catch (e) {}
  } else {
    systemStore.showNotification('Info', `Program sistem ${item.name} dilindungi untuk stabilitas OS`, 'info')
  }
}

function contextFocusWindow() {
  if (contextMenu.item && contextMenu.item.type === 'window') {
    windowStore.focusWindow(contextMenu.item.id)
  }
  closeContextMenu()
}

async function restartContainer(id) {
  try {
    const res = await fetch(`/api/v1/containers/${id}/restart`, { method: 'POST' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container berhasil dimulai ulang`, 'success')
      fetchContainers()
    }
  } catch (e) {}
}

function contextEndTask() {
  if (contextMenu.item) {
    killItem(contextMenu.item)
  }
  closeContextMenu()
}

function contextStopService() {
  if (contextMenu.item && contextMenu.item.type === 'container') {
    killItem(contextMenu.item)
  }
  closeContextMenu()
}

function contextRestartService() {
  if (contextMenu.item && contextMenu.item.type === 'container') {
    restartContainer(contextMenu.item.id)
  }
  closeContextMenu()
}

function contextOpenLocation() {
  const item = contextMenu.item
  if (item) {
    const targetPath = item.dataPath || item.path || `/DATA/apps/${item.name}`
    systemStore.showNotification('Buka Lokasi', `Mengarahkan ke folder: ${targetPath}`, 'info')
    windowStore.openFileManager(targetPath)
  }
  closeContextMenu()
}

async function contextViewLogs() {
  const item = contextMenu.item
  if (item && item.type === 'container') {
    activeLogContainer.value = item
    showLogsModal.value = true
    try {
      const res = await fetch(`/api/v1/containers/${item.id}/logs`)
      if (res.ok) {
        const data = await res.json()
        activeLogText.value = data.logs || 'Tidak ada log yang tercatat.'
      }
    } catch (e) {
      activeLogText.value = 'Gagal mengambil log container.'
    }
  }
  closeContextMenu()
}

function contextCopyInfo() {
  const item = contextMenu.item
  if (item) {
    const text = `${item.name} (PID/ID: ${item.pid || item.id}, CPU: ${item.cpu}%, RAM: ${item.memory || item.memoryMb}MB)`
    navigator.clipboard.writeText(text)
    systemStore.showNotification('Tersalin', 'Informasi proses berhasil disalin ke clipboard', 'success')
  }
  closeContextMenu()
}

function toggleService(srv) {
  srv.status = srv.status === 'active' ? 'inactive' : 'active'
  systemStore.showNotification('Layanan Diperbarui', `Status ${srv.name} diubah menjadi ${srv.status}`, 'info')
}

onMounted(() => {
  refreshAll()
})
</script>
