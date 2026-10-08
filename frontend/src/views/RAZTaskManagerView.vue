<script setup>
/**
 * RAZOS Task Manager App View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka pemantauan metrik sistem, konsumsi RAM, Suhu SOC STB,
 * serta tombol 1-klik Clean RAM (drop_caches) & Clean Storage (Docker Prune).
 */
import { useSystemStore } from '../stores/RAZSystemStore'
import { Cpu, HardDrive, Thermometer, Database, Sparkles, Trash2, RefreshCw } from 'lucide-vue-next'

const systemStore = useSystemStore()

function formatBytes(bytes) {
  if (!bytes) return '0 MB'
  const mb = bytes / (1024 * 1024)
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(0) + ' MB'
}
</script>

<template>
  <div class="h-full flex flex-col gap-4 text-xs sm:text-sm">
    <!-- Header Summary Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <!-- CPU Metric Card -->
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-1">
        <div class="flex items-center justify-between text-slate-400">
          <span class="font-medium">Beban CPU</span>
          <Cpu class="w-4 h-4 text-sky-400" />
        </div>
        <div class="text-xl font-bold text-sky-400">
          {{ systemStore.metrics.cpu_percent.toFixed(1) }}%
        </div>
        <div class="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mt-1">
          <div class="bg-sky-400 h-full transition-all duration-300" :style="{ width: `${systemStore.metrics.cpu_percent}%` }"></div>
        </div>
      </div>

      <!-- RAM Metric Card -->
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-1">
        <div class="flex items-center justify-between text-slate-400">
          <span class="font-medium">RAM Fisik</span>
          <HardDrive class="w-4 h-4 text-emerald-400" />
        </div>
        <div class="text-xl font-bold text-emerald-400">
          {{ systemStore.metrics.ram_percent.toFixed(1) }}%
        </div>
        <div class="text-[11px] text-slate-400">
          {{ formatBytes(systemStore.metrics.ram_used) }} / {{ formatBytes(systemStore.metrics.ram_total) }}
        </div>
        <div class="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mt-1">
          <div class="bg-emerald-400 h-full transition-all duration-300" :style="{ width: `${systemStore.metrics.ram_percent}%` }"></div>
        </div>
      </div>

      <!-- SOC Temp Card -->
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-1">
        <div class="flex items-center justify-between text-slate-400">
          <span class="font-medium">Suhu STB SOC</span>
          <Thermometer class="w-4 h-4 text-amber-400" />
        </div>
        <div class="text-xl font-bold text-amber-400">
          {{ systemStore.metrics.cpu_temp.toFixed(1) }}°C
        </div>
        <div class="text-[11px] text-slate-400">
          Status: {{ systemStore.metrics.cpu_temp < 65 ? 'Stabil / Dingin' : 'Perlu Kipas Tambahan' }}
        </div>
      </div>

      <!-- Storage Card -->
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-1">
        <div class="flex items-center justify-between text-slate-400">
          <span class="font-medium">Penyimpanan /DATA</span>
          <Database class="w-4 h-4 text-purple-400" />
        </div>
        <div class="text-xl font-bold text-purple-400">
          {{ systemStore.metrics.disk_percent.toFixed(1) }}%
        </div>
        <div class="text-[11px] text-slate-400">
          {{ formatBytes(systemStore.metrics.disk_used) }} / {{ formatBytes(systemStore.metrics.disk_total) }}
        </div>
      </div>
    </div>

    <!-- Maintenance Action Section -->
    <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-3">
      <div>
        <h4 class="font-semibold text-slate-200">Pemeliharaan Cepat STB (1-Klik)</h4>
        <p class="text-xs text-slate-400">Bebaskan cache RAM Linux atau bersihkan sampah image Docker tanpa restart server.</p>
      </div>

      <div class="flex items-center gap-2 w-full sm:w-auto">
        <!-- Clean RAM -->
        <button
          @click="systemStore.cleanRAM"
          :disabled="systemStore.isCleaningRAM"
          class="flex-1 sm:flex-initial px-3 py-2 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 font-medium flex items-center justify-center gap-2 transition disabled:opacity-50"
        >
          <Sparkles class="w-4 h-4" />
          <span>{{ systemStore.isCleaningRAM ? 'Membersihkan...' : 'Clean RAM' }}</span>
        </button>

        <!-- Clean Storage -->
        <button
          @click="systemStore.cleanStorage"
          :disabled="systemStore.isCleaningStorage"
          class="flex-1 sm:flex-initial px-3 py-2 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/30 font-medium flex items-center justify-center gap-2 transition disabled:opacity-50"
        >
          <Trash2 class="w-4 h-4" />
          <span>{{ systemStore.isCleaningStorage ? 'Pruning...' : 'Clean Storage' }}</span>
        </button>

        <!-- Refresh -->
        <button
          @click="systemStore.fetchMetrics"
          class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium flex items-center justify-center gap-2 transition"
          title="Segarkan Metrik"
        >
          <RefreshCw class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Container / Process Monitoring Placeholder Table -->
    <div class="flex-1 rounded-xl bg-slate-900/60 border border-slate-800/80 p-3 overflow-auto flex flex-col">
      <div class="flex items-center justify-between mb-2">
        <span class="font-semibold text-slate-300">Daftar Service & Container Aktif</span>
        <span class="text-xs text-slate-400">Total: 4 Container</span>
      </div>
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 text-slate-400">
            <th class="py-2 px-2">Nama Container</th>
            <th class="py-2 px-2">Port Terbuka</th>
            <th class="py-2 px-2">CPU</th>
            <th class="py-2 px-2">Memori</th>
            <th class="py-2 px-2">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-slate-300">
          <tr>
            <td class="py-2.5 px-2 font-mono text-sky-400">razos-core</td>
            <td class="py-2.5 px-2">8080 (Web UI)</td>
            <td class="py-2.5 px-2">0.4%</td>
            <td class="py-2.5 px-2">24 MB</td>
            <td class="py-2.5 px-2"><span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">Running</span></td>
          </tr>
          <tr>
            <td class="py-2.5 px-2 font-mono text-sky-400">razos-mariadb</td>
            <td class="py-2.5 px-2">3306 (Internal)</td>
            <td class="py-2.5 px-2">0.2%</td>
            <td class="py-2.5 px-2">128 MB</td>
            <td class="py-2.5 px-2"><span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">Running</span></td>
          </tr>
          <tr>
            <td class="py-2.5 px-2 font-mono text-sky-400">razos-postgres</td>
            <td class="py-2.5 px-2">5432 (Internal)</td>
            <td class="py-2.5 px-2">0.1%</td>
            <td class="py-2.5 px-2">92 MB</td>
            <td class="py-2.5 px-2"><span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">Running</span></td>
          </tr>
          <tr>
            <td class="py-2.5 px-2 font-mono text-sky-400">razos-adminer</td>
            <td class="py-2.5 px-2">8088 (Web DB)</td>
            <td class="py-2.5 px-2">0.0%</td>
            <td class="py-2.5 px-2">18 MB</td>
            <td class="py-2.5 px-2"><span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">Running</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
