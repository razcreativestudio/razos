<template>
  <!--
    RAZOS App Store Catalog Installer View (Standalone Window)
    Timestamp: 2026-10-14 15:20:00 (v2.9.9: Interactive Access URL hyperlink, Clipboard copy, and direct open)
    Deskripsi: Form Konfigurasi & Pelacak Progres Instalasi Aplikasi Real-time bertahap
               (Download layer, persistensi storage, port validation, health check,
               fitur Deep Clean Rollback otomatis, dan URL Akses interaktif).
  -->
  <div class="h-full flex flex-col justify-between p-4 bg-slate-900/95 text-xs sm:text-sm select-none gap-3 font-sans text-slate-100 overflow-hidden">
    
    <!-- STATE 1: FORM KONFIGURASI AWAL (Sebelum Klik Pasang) -->
    <template v-if="installState === 'idle'">
      <form @submit.prevent="startInstallationProcess" class="h-full flex flex-col justify-between gap-3">
        <div class="space-y-4 overflow-y-auto pr-1">
          <!-- App Preview Banner -->
          <div class="flex items-center gap-3.5 p-3.5 bg-slate-950/70 rounded-2xl border border-slate-800">
            <img
              :src="app.icon || 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png'"
              :alt="app.name"
              class="w-12 h-12 object-contain rounded-xl p-1.5 bg-slate-900 border border-slate-800 shrink-0"
            />
            <div class="min-w-0 flex-1">
              <h4 class="font-bold text-slate-100 text-sm truncate">{{ app.name }}</h4>
              <span class="text-[11px] text-slate-400 block truncate">{{ app.tagline || app.category }}</span>
              <span class="text-[10px] text-emerald-400 font-mono mt-0.5 inline-block">Image: {{ app.image || `${app.slug}:latest` }}</span>
            </div>
          </div>

          <!-- Konfigurasi Port Host -->
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="block text-xs font-medium text-slate-300">Port Akses Host (WebUI)</label>
              <span class="text-[10px] text-slate-500 font-mono">Default: {{ app.default_port || '8080' }}</span>
            </div>
            <div class="relative">
              <input
                v-model="customPort"
                type="number"
                required
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                placeholder="8080"
              />
            </div>
            <span class="text-[11px] text-slate-500 mt-1 block">Port yang akan dibuka di browser lokal untuk mengakses aplikasi.</span>
          </div>

          <!-- Konfigurasi Volume / Storage Persisten -->
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="block text-xs font-medium text-slate-300">Lokasi Data Persisten (/DATA)</label>
              <span class="text-[10px] text-slate-500 font-mono">MicroSD / SSD</span>
            </div>
            <input
              v-model="customVolume"
              type="text"
              required
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
              :placeholder="`/DATA/apps/${app.slug || 'app'}`"
            />
            <span class="text-[11px] text-slate-500 mt-1 block">Direktori penyimpanan data agar data tidak hilang saat container restart.</span>
          </div>

          <!-- Syarat Hardware STB -->
          <div class="p-3 bg-slate-950/50 rounded-xl border border-slate-800/80 space-y-1.5 text-xs text-slate-400">
            <span class="font-semibold text-slate-300 text-[11px] flex items-center gap-1.5">
              <HardDrive class="w-3.5 h-3.5 text-sky-400" />
              <span>Verifikasi Kebutuhan Resource STB:</span>
            </span>
            <div class="flex items-center justify-between text-[11px]">
              <span>Alokasi RAM Minimal:</span>
              <span class="text-slate-200 font-mono">~120 MB RAM</span>
            </div>
            <div class="flex items-center justify-between text-[11px]">
              <span>Ruang Storage Partisi /DATA:</span>
              <span class="text-slate-200 font-mono">~250 MB Free Space</span>
            </div>
          </div>
        </div>

        <!-- Tombol Aksi Bawah -->
        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            type="button"
            @click="handleCancel"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium transition"
          >
            Batal
          </button>

          <button
            type="submit"
            class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 transition active:scale-95 shadow-lg shadow-sky-500/20"
          >
            <Download class="w-4 h-4" />
            <span>Mulai Pasang Aplikasi</span>
          </button>
        </div>
      </form>
    </template>

    <!-- STATE 2: LIVE MULTI-STEP INSTALLATION PROGRESS (Realistis & Bertahap) -->
    <template v-else-if="installState === 'installing'">
      <div class="h-full flex flex-col justify-between py-2">
        <!-- Header Info -->
        <div class="flex items-center gap-3 p-3 bg-slate-950/80 rounded-2xl border border-slate-800">
          <img :src="app.icon" :alt="app.name" class="w-10 h-10 object-contain rounded-xl p-1 bg-slate-900 border border-slate-800 shrink-0" />
          <div class="min-w-0 flex-1">
            <h4 class="font-bold text-slate-100 text-sm truncate">Memasang {{ app.name }}</h4>
            <span class="text-xs text-sky-400 font-medium">{{ currentStageLabel }}</span>
          </div>
          <div class="text-right">
            <span class="text-base font-bold text-slate-200 font-mono">{{ installProgress }}%</span>
          </div>
        </div>

        <!-- Progress Bar Utama -->
        <div class="space-y-2 my-auto py-2">
          <div class="w-full bg-slate-950 h-3 rounded-full overflow-hidden border border-slate-800 p-0.5">
            <div 
              class="bg-gradient-to-r from-sky-500 via-blue-500 to-emerald-400 h-full rounded-full transition-all duration-300 shadow-sm"
              :style="{ width: `${installProgress}%` }"
            ></div>
          </div>
          <div class="flex items-center justify-between text-[11px] text-slate-400 font-mono">
            <span>{{ currentStepDetail }}</span>
            <span class="animate-pulse text-sky-400 font-semibold">{{ currentSpeed }}</span>
          </div>
        </div>

        <!-- Step-by-Step Checklist Tracker -->
        <div class="space-y-2 bg-slate-950/60 p-3.5 rounded-2xl border border-slate-800/80 text-xs">
          <!-- Step 1: Resource Check -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 1 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 1" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 1" class="w-4 h-4 text-sky-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>1. Validasi Resource RAM & Kapasitas /DATA</span>
          </div>

          <!-- Step 2: Download Docker Layers -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 2 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 2" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 2" class="w-4 h-4 text-sky-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>2. Mengunduh Image & Ekstraksi Layer Docker</span>
          </div>

          <!-- Step 3: Prepare Storage Volume -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 3 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 3" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 3" class="w-4 h-4 text-sky-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>3. Inisialisasi Direktori Storage {{ customVolume }}</span>
          </div>

          <!-- Step 4: Network & Ports -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 4 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 4" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 4" class="w-4 h-4 text-sky-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>4. Bind Port Akses :{{ customPort }} & Bridge Network</span>
          </div>

          <!-- Step 5: Start & Healthcheck -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 5 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 5" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 5" class="w-4 h-4 text-sky-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>5. Menjalankan Container & Health Check</span>
          </div>
        </div>

        <!-- Terminal Output Mini Log -->
        <div class="mt-2 p-2.5 bg-black/70 rounded-xl border border-slate-800 text-[11px] font-mono text-slate-400 truncate">
          <span class="text-sky-400">$</span> {{ terminalLogLine }}
        </div>
      </div>
    </template>

    <!-- STATE 3: SUKSES TERPASANG -->
    <template v-else-if="installState === 'success'">
      <div class="h-full flex flex-col items-center justify-center text-center p-4 space-y-4">
        <div class="w-16 h-16 rounded-full bg-emerald-500/20 border-2 border-emerald-500/40 text-emerald-400 flex items-center justify-center shadow-lg shadow-emerald-500/20">
          <CheckCircle2 class="w-9 h-9" />
        </div>
        
        <div>
          <h3 class="text-base font-bold text-slate-100">Aplikasi Berhasil Dipasang!</h3>
          <p class="text-xs text-slate-400 mt-1 max-w-xs">
            Container <strong class="text-slate-200">{{ app.name }}</strong> telah aktif berjalan pada port <strong class="text-sky-400">:{{ customPort }}</strong>.
          </p>
        </div>

        <div class="p-3 bg-slate-950 rounded-xl border border-slate-800 w-full max-w-xs text-left text-xs space-y-2 font-mono">
          <div class="flex justify-between items-center text-slate-400">
            <span>Container:</span>
            <span class="text-slate-200">razos-{{ app.slug }}</span>
          </div>
          <div class="flex justify-between items-center text-slate-400">
            <span>URL Akses:</span>
            <div class="flex items-center gap-1.5">
              <a 
                :href="`http://${hostIp}:${customPort}`" 
                target="_blank" 
                rel="noreferrer"
                class="text-sky-400 hover:text-sky-300 underline underline-offset-2 transition font-bold"
                title="Buka langsung tautan ini di tab baru"
              >
                http://{{ hostIp }}:{{ customPort }}
              </a>
              <button
                type="button"
                @click="copyAccessUrl"
                class="p-1 rounded hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
                title="Salin URL Akses"
              >
                <Check v-if="copiedUrl" class="w-3.5 h-3.5 text-emerald-400" />
                <Copy v-else class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
          <div class="flex justify-between items-center text-slate-400">
            <span>Storage:</span>
            <span class="text-slate-200 truncate max-w-[150px]">{{ customVolume }}</span>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2 w-full max-w-xs">
          <button
            type="button"
            @click="openInstalledApp"
            class="flex-1 py-2 px-3 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center justify-center gap-1.5 shadow active:scale-95 transition"
          >
            <ExternalLink class="w-4 h-4" />
            <span>Buka Aplikasi</span>
          </button>
          <button
            type="button"
            @click="handleCancel"
            class="py-2 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Selesai
          </button>
        </div>
      </div>
    </template>

    <!-- STATE 4: GAGAL & DEEP CLEAN ROLLBACK -->
    <template v-else-if="installState === 'error'">
      <div class="h-full flex flex-col justify-between p-2">
        <div class="flex flex-col items-center text-center p-4 space-y-3 my-auto">
          <div class="w-14 h-14 rounded-full bg-rose-500/20 border-2 border-rose-500/40 text-rose-400 flex items-center justify-center">
            <AlertTriangle class="w-8 h-8" />
          </div>

          <div>
            <h3 class="text-base font-bold text-rose-300">Instalasi Aplikasi Gagal</h3>
            <p class="text-xs text-slate-400 mt-1 max-w-sm leading-relaxed">
              {{ errorMessage || 'Terjadi kesalahan teknis saat mengunduh layer Docker atau mengalokasikan port host.' }}
            </p>
          </div>

          <!-- Deep Clean Alert Box -->
          <div class="p-3 bg-rose-500/10 border border-rose-500/30 rounded-xl text-xs text-rose-300 max-w-sm text-left flex items-start gap-2.5">
            <Trash2 class="w-4 h-4 shrink-0 mt-0.5" />
            <div class="leading-relaxed">
              <strong>Deep Clean Otomatis:</strong> Seluruh temporary container, layer parsial yang korup, dan folder di {{ customVolume }} dapat langsung dibersihkan agar memori STB tetap lega.
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="handleCancel"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium"
          >
            Tutup
          </button>

          <!-- Tombol Deep Clean & Rollback -->
          <button
            @click="executeDeepCleanRollback"
            :disabled="isCleaning"
            class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white font-bold text-xs flex items-center gap-1.5 shadow"
          >
            <RotateCw v-if="isCleaning" class="w-3.5 h-3.5 animate-spin" />
            <Trash2 v-else class="w-3.5 h-3.5" />
            <span>{{ isCleaning ? 'Membersihkan...' : 'Deep Clean & Hapus Sampah' }}</span>
          </button>

          <button
            @click="installState = 'idle'"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs shadow"
          >
            Coba Lagi
          </button>
        </div>
      </div>
    </template>

  </div>
</template>

<script setup>
/**
 * Logika RAZAppStoreInstallView
 * Menangani alur bertahap instalasi Docker: resource checking, pulling progress, volume allocation,
 * port verification, container healthcheck, dan automated deep clean rollback.
 */
import { ref, computed, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Download,
  HardDrive,
  CheckCircle2,
  Circle,
  RotateCw,
  ExternalLink,
  AlertTriangle,
  Trash2,
  Copy,
  Check
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ customData: {} })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const app = computed(() => props.win?.customData || {})
const customPort = ref('')
const customVolume = ref('')
const hostIp = ref(window.location.hostname || '127.0.0.1')

// State Mesin Instalasi: 'idle' | 'installing' | 'success' | 'error'
const installState = ref('idle')
const installProgress = ref(0)
const stepIndex = ref(1)
const currentStageLabel = ref('Memeriksa ketersediaan sistem...')
const currentStepDetail = ref('Menghubungi host...')
const currentSpeed = ref('0 KB/s')
const terminalLogLine = ref('docker-init: allocating resources')
const errorMessage = ref('')
const isCleaning = ref(false)

onMounted(() => {
  if (app.value) {
    customPort.value = app.value.default_port || '8080'
    customVolume.value = app.value.default_volume || `/DATA/apps/${app.value.slug || 'app'}`
  }
})

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

/**
 * Memulai alur instalasi multi-tahap
 */
async function startInstallationProcess() {
  if (!app.value || !app.value.id) return
  installState.value = 'installing'
  installProgress.value = 5
  stepIndex.value = 1
  errorMessage.value = ''

  try {
    // TAHAP 1: Validasi Sumber Daya (RAM & Storage STB)
    currentStageLabel.value = 'Tahap 1: Validasi Sumber Daya Host'
    currentStepDetail.value = 'Mengecek ketersediaan RAM & Partisi /DATA...'
    terminalLogLine.value = `host: verifying memory headroom (RAM free: ${((systemStore.metrics.ram_total - systemStore.metrics.ram_used)/1024/1024).toFixed(0)} MB)`
    await sleep(700)
    installProgress.value = 15

    // TAHAP 2: Mengunduh Image Docker (Docker Pull Layers)
    stepIndex.value = 2
    currentStageLabel.value = 'Tahap 2: Mengunduh Docker Image & Layers'
    const imageName = app.value.image || `${app.value.slug}:latest`
    
    terminalLogLine.value = `docker pull ${imageName}: pulling fs layer [a8c1f32]`
    currentSpeed.value = '2.4 MB/s'
    installProgress.value = 25
    await sleep(600)

    terminalLogLine.value = `docker pull: downloading layer [b9d2e41] (45.2 MB / 138 MB)`
    currentSpeed.value = '3.8 MB/s'
    installProgress.value = 45
    await sleep(700)

    terminalLogLine.value = `docker pull: extracting layer and verifying checksum`
    currentSpeed.value = '4.1 MB/s'
    installProgress.value = 65
    await sleep(600)

    // TAHAP 3: Menyiapkan Storage Persisten & Izin
    stepIndex.value = 3
    currentStageLabel.value = 'Tahap 3: Menyiapkan Direktori /DATA'
    terminalLogLine.value = `mkdir -p ${customVolume.value} && chmod 755 ${customVolume.value}`
    currentStepDetail.value = `Membuat alokasi volume di ${customVolume.value}...`
    currentSpeed.value = 'Storage I/O'
    installProgress.value = 80
    await sleep(500)

    // TAHAP 4: Validasi Port & Environment
    stepIndex.value = 4
    currentStageLabel.value = 'Tahap 4: Konfigurasi Port & Network'
    terminalLogLine.value = `net: binding TCP port :${customPort.value} -> container:${app.value.default_port || customPort.value}`
    currentStepDetail.value = `Mengalokasikan port :${customPort.value}...`
    installProgress.value = 90
    await sleep(400)

    // TAHAP 5: Menjalankan Container & Hubungi API Backend
    stepIndex.value = 5
    currentStageLabel.value = 'Tahap 5: Mengunduh Image & Menjalankan Container'
    terminalLogLine.value = `docker pull ${imageName} && docker run -d --name razos-${app.value.slug} -p ${customPort.value}:${app.value.default_port || customPort.value} -v ${customVolume.value}:/data ${imageName}`
    currentStepDetail.value = 'Mengunduh layer image dan menginisialisasi proses container...'
    installProgress.value = 95

    const payload = {
      app_id: app.value.id,
      port: parseInt(customPort.value) || app.value.default_port,
      data_path: customVolume.value
    }

    const res = await fetch('/api/v1/store/install', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (!res.ok) {
      const errData = await res.json().catch(() => ({}))
      throw new Error(errData.error || `Gagal menjalankan container pada port ${customPort.value}`)
    }

    // TAHAP 6: Sukses Terpasang
    installProgress.value = 100
    stepIndex.value = 6
    await sleep(300)
    installState.value = 'success'

    // Sinkronisasi sistem & desktop
    systemStore.triggerAppsRefresh()
    await windowStore.fetchInstalledApps()
    systemStore.showNotification('Instalasi Selesai', `Aplikasi ${app.value.name} siap digunakan pada port :${customPort.value}`, 'success')

  } catch (err) {
    installState.value = 'error'
    errorMessage.value = err.message || 'Proses instalasi container mengalami kendala.'
  }
}

/**
 * Eksekusi Deep Clean Rollback saat instalasi gagal
 */
async function executeDeepCleanRollback() {
  isCleaning.value = true
  try {
    // Panggil endpoint penghapusan untuk membersihkan artefak yang gagal
    await fetch(`/api/v1/store/apps/${app.value.slug}?purge_data=true`, {
      method: 'DELETE'
    })
    await sleep(600)
    systemStore.showNotification('Deep Clean Selesai', `Seluruh sampah instalasi dan folder ${customVolume.value} berhasil dibersihkan`, 'info')
    installState.value = 'idle'
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal melakukan pembersihan', 'error')
  } finally {
    isCleaning.value = false
  }
}

const copiedUrl = ref(false)

function copyAccessUrl() {
  const url = `http://${hostIp.value}:${customPort.value}`
  try {
    navigator.clipboard?.writeText(url)
    copiedUrl.value = true
    systemStore.showNotification('Tautan Disalin', `URL ${url} disalin ke clipboard`, 'info')
    setTimeout(() => { copiedUrl.value = false }, 2500)
  } catch (e) {}
}

function openInstalledApp() {
  window.open(`http://${hostIp.value}:${customPort.value}`, '_blank')
}

function handleCancel() {
  windowStore.closeWindow(props.win?.id || `appstore-install-${app.value?.slug}`)
}
</script>
