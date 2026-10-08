<script setup>
/**
 * RAZOS Task Manager App View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka pemantauan metrik sistem, konsumsi RAM, Suhu SOC STB,
 * manajemen proses riil (Kill PID), status container Docker (Stop/Restart),
 * serta tombol 1-klik Clean RAM (drop_caches) & Clean Storage (Docker Prune).
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Cpu, HardDrive, Thermometer, Database, Sparkles, Trash2, RefreshCw, XOctagon, RotateCw } from 'lucide-vue-next'

const systemStore = useSystemStore()
const activeTab = ref('processes') // 'processes' | 'containers'
const processes = ref([])
const containers = ref([])
const isLoading = ref(false)

function formatBytes(bytes) {
  if (!bytes) return '0 MB'
  const mb = bytes / (1024 * 1024)
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(0) + ' MB'
}

async function fetchProcesses() {
  try {
    const res = await fetch('/api/v1/processes')
    if (res.ok) {
      processes.value = await res.json()
    }
  } catch (err) {}
}

async function fetchContainers() {
  try {
    const res = await fetch('/api/v1/containers')
    if (res.ok) {
      containers.value = await res.json()
    }
  } catch (err) {}
}

async function refreshAll() {
  isLoading.value = true
  await Promise.all([
    systemStore.fetchMetrics(),
    fetchProcesses(),
    fetchContainers()
  ])
  isLoading.value = false
}

// 1-Klik Kill Process
async function killProcess(pid) {
  try {
    const res = await fetch(`/api/v1/processes/${pid}/kill`, { method: 'POST' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Proses PID ${pid} berhasil dihentikan`, 'success')
      fetchProcesses()
    } else {
      systemStore.showNotification('Gagal', `Gagal menghentikan PID ${pid}`, 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// 1-Klik Stop Container
async function stopContainer(id) {
  try {
    const res = await fetch(`/api/v1/containers/${id}/stop`, { method: 'POST' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container ${id} dihentikan`, 'success')
      fetchContainers()
    }
  } catch (err) {}
}

// 1-Klik Restart Container
async function restartContainer(id) {
  try {
    const res = await fetch(`/api/v1/containers/${id}/restart`, { method: 'POST' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container ${id} dimulai ulang`, 'success')
      fetchContainers()
    }
  } catch (err) {}
}

onMounted(() => {
  refreshAll()
})
</script>

<template>
  <div class="h-full flex flex-col gap-4 text-xs sm:text-sm">
    <!-- Header Summary Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <!-- CPU Metric Card -->
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-1">
        <div class="flex items-center justify-between text-slate-400">
          <span class="font-medium">Beban CPU Host</span>
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
          <span class="font-medium">RAM Fisik Host</span>
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
          <span class="font-medium">Suhu SOC STB</span>
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
        <p class="text-xs text-slate-400">Bebaskan cache RAM Linux atau pangkas sampah container Docker.</p>
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
          @click="refreshAll"
          :disabled="isLoading"
          class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium flex items-center justify-center gap-2 transition"
          title="Segarkan Metrik"
        >
          <RefreshCw class="w-4 h-4" :class="isLoading ? 'animate-spin' : ''" />
        </button>
      </div>
    </div>

    <!-- Tab Selection: Processes vs Containers -->
    <div class="flex items-center gap-2 border-b border-slate-800 pb-2">
      <button
        @click="activeTab = 'processes'"
        class="px-3 py-1.5 rounded-lg text-xs font-semibold transition"
        :class="activeTab === 'processes' ? 'bg-sky-500/20 text-sky-400 border border-sky-400/40' : 'text-slate-400 hover:text-slate-200'"
      >
        Proses Host (Top PID)
      </button>
      <button
        @click="activeTab = 'containers'"
        class="px-3 py-1.5 rounded-lg text-xs font-semibold transition"
        :class="activeTab === 'containers' ? 'bg-sky-500/20 text-sky-400 border border-sky-400/40' : 'text-slate-400 hover:text-slate-200'"
      >
        Container Docker (Live)
      </button>
    </div>

    <!-- Table Section -->
    <div class="flex-1 rounded-xl bg-slate-900/60 border border-slate-800/80 p-3 overflow-auto flex flex-col">
      <!-- Tabel Proses Host -->
      <table v-if="activeTab === 'processes' && processes.length > 0" class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 text-slate-400">
            <th class="py-2 px-2">PID</th>
            <th class="py-2 px-2">Nama Program</th>
            <th class="py-2 px-2">CPU</th>
            <th class="py-2 px-2">Memori</th>
            <th class="py-2 px-2 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-slate-300">
          <tr v-for="proc in processes" :key="proc.pid">
            <td class="py-2 px-2 font-mono text-sky-400">{{ proc.pid }}</td>
            <td class="py-2 px-2 font-medium">{{ proc.name }}</td>
            <td class="py-2 px-2">{{ proc.cpu.toFixed(1) }}%</td>
            <td class="py-2 px-2">{{ proc.memory.toFixed(1) }}%</td>
            <td class="py-2 px-2 text-right">
              <button
                @click="killProcess(proc.pid)"
                class="px-2 py-1 rounded bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 transition inline-flex items-center gap-1 text-[11px]"
                title="Hentikan Proses (Kill)"
              >
                <XOctagon class="w-3 h-3" />
                <span>Kill</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Tabel Container Docker -->
      <table v-else-if="activeTab === 'containers' && containers.length > 0" class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 text-slate-400">
            <th class="py-2 px-2">Container ID</th>
            <th class="py-2 px-2">Nama</th>
            <th class="py-2 px-2">Port</th>
            <th class="py-2 px-2">Status</th>
            <th class="py-2 px-2 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-slate-300">
          <tr v-for="c in containers" :key="c.id">
            <td class="py-2 px-2 font-mono text-sky-400">{{ c.id }}</td>
            <td class="py-2 px-2 font-medium">{{ c.name }}</td>
            <td class="py-2 px-2">{{ c.ports || '-' }}</td>
            <td class="py-2 px-2">
              <span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">{{ c.status }}</span>
            </td>
            <td class="py-2 px-2 text-right flex items-center justify-end gap-1.5">
              <button
                @click="restartContainer(c.id)"
                class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 transition text-[11px]"
                title="Restart"
              >
                <RotateCw class="w-3 h-3" />
              </button>
              <button
                @click="stopContainer(c.id)"
                class="px-2 py-1 rounded bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 transition text-[11px]"
                title="Stop"
              >
                <XOctagon class="w-3 h-3" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Empty State jika belum ada item -->
      <div v-else class="flex-1 flex flex-col items-center justify-center text-slate-500 gap-2 p-6">
        <p class="text-xs">Tidak ada data {{ activeTab === 'processes' ? 'proses aktif' : 'container Docker' }}.</p>
      </div>
    </div>
  </div>
</template>
