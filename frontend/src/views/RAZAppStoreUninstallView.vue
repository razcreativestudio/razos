<template>
  <!--
    RAZOS App Store Uninstaller View (Standalone Window)
    File: RAZAppStoreUninstallView.vue
    Timestamp: 2026-10-26 (v4.1.0)
    Deskripsi: Jendela Pelacak Progres Uninstall Aplikasi & Container Docker secara Real-Time
               dengan Live Terminal Output Console Box, autoscroll, dan verifikasi pelepasan resource.
  -->
  <div class="h-full flex flex-col justify-between p-4 bg-slate-900/95 text-xs sm:text-sm select-none gap-3 font-sans text-slate-100 overflow-hidden">
    
    <!-- STATE 1: SEDANG UNINSTALL BERJALAN -->
    <template v-if="uninstallState === 'running'">
      <div class="h-full flex flex-col justify-between py-1 gap-2">
        <!-- Header Info -->
        <div class="flex items-center gap-3 p-3 bg-slate-950/80 rounded-2xl border border-slate-800 shrink-0">
          <img 
            :src="app.icon || '/assets/icons/RAZIconAppStore.svg'" 
            :alt="app.name || app.title" 
            class="w-10 h-10 object-contain rounded-xl p-1 bg-slate-900 border border-slate-800 shrink-0" 
          />
          <div class="min-w-0 flex-1">
            <h4 class="font-bold text-slate-100 text-sm truncate">Uninstall {{ app.name || app.title }}</h4>
            <span class="text-xs text-rose-400 font-medium">{{ currentStageLabel }}</span>
          </div>
          <div class="text-right">
            <span class="text-base font-bold text-rose-300 font-mono">{{ uninstallProgress }}%</span>
          </div>
        </div>

        <!-- Progress Bar Utama -->
        <div class="space-y-1.5 my-1 shrink-0">
          <div class="w-full bg-slate-950 h-2.5 rounded-full overflow-hidden border border-slate-800 p-0.5">
            <div 
              class="bg-gradient-to-r from-rose-500 via-red-500 to-amber-400 h-full rounded-full transition-all duration-300 shadow-sm"
              :style="{ width: `${uninstallProgress}%` }"
            ></div>
          </div>
          <div class="flex items-center justify-between text-[11px] text-slate-400 font-mono">
            <span class="truncate max-w-[280px]">{{ currentStepDetail }}</span>
            <span class="animate-pulse text-rose-400 font-semibold shrink-0">{{ currentSpeed }}</span>
          </div>
        </div>

        <!-- Step-by-Step Checklist Tracker -->
        <div class="space-y-1.5 bg-slate-950/60 p-3 rounded-2xl border border-slate-800/80 text-xs shrink-0">
          <!-- Step 1: Validasi Container -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 1 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 1" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 1" class="w-4 h-4 text-rose-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>1. Validasi Container razos-{{ appSlug }}</span>
          </div>

          <!-- Step 2: Stop Container -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 2 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 2" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 2" class="w-4 h-4 text-rose-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>2. Menghentikan Container (docker stop)</span>
          </div>

          <!-- Step 3: Remove Container & Volumes -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 3 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 3" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 3" class="w-4 h-4 text-rose-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>3. Menghapus Container & Volume Docker (docker rm -v)</span>
          </div>

          <!-- Step 4: Storage Directory Purge -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 4 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 4" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 4" class="w-4 h-4 text-rose-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>4. {{ purgeData ? `Pembersihan Direktori /DATA/apps/${appSlug}` : 'Mempertahankan Data Storage' }}</span>
          </div>

          <!-- Step 5: Final Memory Release -->
          <div class="flex items-center gap-2.5" :class="stepIndex >= 5 ? 'text-slate-200' : 'text-slate-500'">
            <CheckCircle2 v-if="stepIndex > 5" class="w-4 h-4 text-emerald-400 shrink-0" />
            <RotateCw v-else-if="stepIndex === 5" class="w-4 h-4 text-rose-400 animate-spin shrink-0" />
            <Circle v-else class="w-4 h-4 text-slate-600 shrink-0" />
            <span>5. Pelepasan Resource & Sinkronisasi Sistem</span>
          </div>
        </div>

        <!-- Terminal Output Live Console Box -->
        <div class="mt-1 flex-1 min-h-[140px] flex flex-col rounded-xl bg-slate-950 border border-slate-800/90 overflow-hidden shadow-2xl">
          <!-- Terminal Title Bar -->
          <div class="flex items-center justify-between px-3 py-1.5 bg-slate-900/90 border-b border-slate-800/80 select-none shrink-0">
            <div class="flex items-center gap-2">
              <div class="flex items-center gap-1.5">
                <span class="w-2.5 h-2.5 rounded-full bg-rose-500/80"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-amber-500/80"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-500/80"></span>
              </div>
              <span class="text-[10px] font-mono text-slate-400 font-semibold tracking-wider flex items-center gap-1.5 ml-1">
                <span>installer@razos:~# docker rm -f</span>
                <span class="inline-block w-1.5 h-1.5 rounded-full bg-rose-400 animate-pulse"></span>
              </span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-[9px] font-mono px-1.5 py-0.5 rounded bg-slate-800 text-rose-300 font-semibold">
                {{ terminalLogs.length }} lines
              </span>
              <button 
                type="button"
                @click="copyTerminalLogs"
                class="text-slate-400 hover:text-slate-200 transition px-1.5 py-0.5 rounded hover:bg-slate-800 flex items-center gap-1 text-[10px] font-mono"
                title="Salin Seluruh Log Terminal"
              >
                <Copy class="w-3 h-3" />
                <span>Salin</span>
              </button>
            </div>
          </div>

          <!-- Terminal Logs Body (Auto-scrolled) -->
          <div 
            ref="terminalScrollRef" 
            class="p-2.5 font-mono text-[10.5px] leading-relaxed overflow-y-auto flex-1 space-y-1 select-text scrollbar-thin scrollbar-thumb-slate-800"
          >
            <div v-for="(log, idx) in terminalLogs" :key="idx" class="break-all flex items-start gap-1.5">
              <span class="text-slate-500 text-[10px] shrink-0 font-mono select-none">[{{ log.time }}]</span>
              <span :class="getLogLevelClass(log.type)" class="shrink-0 font-bold font-mono text-[10px] select-none">
                {{ log.tag || '$' }}
              </span>
              <span :class="getLogTextClass(log.type)" class="flex-1 font-mono">
                {{ log.text }}
              </span>
            </div>
            <!-- Live Blinking Prompt -->
            <div v-if="uninstallState === 'running'" class="flex items-center gap-1 text-slate-500 text-[10px] pt-0.5">
              <span class="text-rose-400 font-bold">$</span>
              <span class="inline-block w-2 h-3.5 bg-rose-400/90 animate-pulse"></span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- STATE 2: SUKSES UNINSTALL -->
    <template v-else-if="uninstallState === 'success'">
      <div class="h-full flex flex-col items-center justify-center text-center p-4 space-y-4">
        <div class="w-16 h-16 rounded-full bg-emerald-500/20 border-2 border-emerald-500/40 text-emerald-400 flex items-center justify-center shadow-lg shadow-emerald-500/20">
          <CheckCircle2 class="w-9 h-9" />
        </div>
        
        <div>
          <h3 class="text-base font-bold text-slate-100">Aplikasi Berhasil Dihapus!</h3>
          <p class="text-xs text-slate-400 mt-1 max-w-xs">
            Container <strong class="text-slate-200">razos-{{ appSlug }}</strong> dan volumenya telah berhasil dibersihkan dari sistem STB.
          </p>
        </div>

        <div class="p-3 bg-slate-950 rounded-xl border border-slate-800 w-full max-w-xs text-left text-xs space-y-2 font-mono">
          <div class="flex justify-between items-center text-slate-400">
            <span>Status:</span>
            <span class="text-emerald-400 font-bold">Uninstalled</span>
          </div>
          <div class="flex justify-between items-center text-slate-400">
            <span>Container:</span>
            <span class="text-slate-200">razos-{{ appSlug }}</span>
          </div>
          <div class="flex justify-between items-center text-slate-400">
            <span>Storage /DATA:</span>
            <span class="text-slate-200">{{ purgeData ? 'Dibersihkan' : 'Dipertahankan' }}</span>
          </div>
        </div>

        <button
          type="button"
          @click="closeWindow"
          class="w-full max-w-xs py-2 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold transition active:scale-95 shadow"
        >
          Selesai & Tutup
        </button>
      </div>
    </template>

    <!-- STATE 3: GAGAL UNINSTALL -->
    <template v-else-if="uninstallState === 'error'">
      <div class="h-full flex flex-col justify-between p-2">
        <div class="flex flex-col items-center text-center p-4 space-y-3 my-auto">
          <div class="w-14 h-14 rounded-full bg-rose-500/20 border-2 border-rose-500/40 text-rose-400 flex items-center justify-center">
            <AlertTriangle class="w-8 h-8" />
          </div>

          <div>
            <h3 class="text-base font-bold text-rose-300">Uninstall Mengalami Kendala</h3>
            <p class="text-xs text-slate-400 mt-1 max-w-sm leading-relaxed">
              {{ errorMessage || 'Terjadi kesalahan teknis saat menghentikan container atau membersihkan storage.' }}
            </p>
          </div>

          <!-- Terminal Logs on Error -->
          <div v-if="terminalLogs.length > 0" class="w-full max-w-sm rounded-xl bg-slate-950 border border-rose-500/30 overflow-hidden text-left shadow-lg">
            <div class="flex items-center justify-between px-3 py-1 bg-slate-900 border-b border-slate-800 text-[10px] font-mono text-slate-400">
              <span class="text-rose-400 font-semibold flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-rose-500 animate-pulse"></span>
                <span>Terminal Execution Log</span>
              </span>
              <button @click="copyTerminalLogs" class="hover:text-white transition flex items-center gap-1">
                <Copy class="w-3 h-3" />
                <span>Salin</span>
              </button>
            </div>
            <div class="p-2.5 font-mono text-[10px] leading-relaxed max-h-28 overflow-y-auto space-y-1">
              <div v-for="(log, idx) in terminalLogs.slice(-10)" :key="idx" class="break-all flex items-start gap-1.5">
                <span class="text-slate-500 text-[9px] shrink-0 font-mono">[{{ log.time }}]</span>
                <span :class="getLogLevelClass(log.type)" class="shrink-0 font-bold text-[9px]">{{ log.tag }}</span>
                <span :class="getLogTextClass(log.type)">{{ log.text }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="closeWindow"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium"
          >
            Tutup
          </button>
          <button
            @click="startUninstallProcess"
            class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white font-bold text-xs shadow"
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
 * Logika RAZAppStoreUninstallView
 * Menangani alur bertahap penghapusan container Docker dengan Live Terminal Output Logger
 */
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  CheckCircle2,
  Circle,
  RotateCw,
  AlertTriangle,
  Copy
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
const appSlug = computed(() => app.value.slug || app.value.id || 'app')
const purgeData = computed(() => app.value.purgeData === true)

// State Mesin Uninstall: 'running' | 'success' | 'error'
const uninstallState = ref('running')
const uninstallProgress = ref(0)
const stepIndex = ref(1)
const currentStageLabel = ref('Memulai proses uninstall...')
const currentStepDetail = ref('Menghubungi host...')
const currentSpeed = ref('Docker Daemon')
const errorMessage = ref('')

// State Live Terminal Logs
const terminalLogs = ref([])
const terminalScrollRef = ref(null)
let pollTimer = null

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

function appendLocalLog(type, tag, text) {
  const now = new Date()
  const timeStr = now.toTimeString().split(' ')[0]
  terminalLogs.value.push({ time: timeStr, type, tag, text })
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => {
    if (terminalScrollRef.value) {
      terminalScrollRef.value.scrollTop = terminalScrollRef.value.scrollHeight
    }
  })
}

async function fetchBackendLogs() {
  const slug = appSlug.value
  if (!slug) return
  try {
    const res = await fetch(`/api/v1/store/install/logs/${slug}`)
    if (res.ok) {
      const data = await res.json()
      if (Array.isArray(data.logs) && data.logs.length > 0) {
        terminalLogs.value = data.logs
        scrollToBottom()

        const lastLog = data.logs[data.logs.length - 1]
        if (lastLog && lastLog.text) {
          currentStepDetail.value = lastLog.text
        }
      }
    }
  } catch (e) {}
}

function getLogLevelClass(type) {
  switch (type) {
    case 'info': return 'text-sky-400'
    case 'docker': return 'text-amber-400'
    case 'success': return 'text-emerald-400'
    case 'error': return 'text-rose-400'
    default: return 'text-slate-400'
  }
}

function getLogTextClass(type) {
  switch (type) {
    case 'error': return 'text-rose-300'
    case 'success': return 'text-emerald-300 font-semibold'
    case 'docker': return 'text-amber-200'
    default: return 'text-slate-400'
  }
}

function copyTerminalLogs() {
  if (terminalLogs.value.length === 0) return
  const text = terminalLogs.value.map(l => `[${l.time}] ${l.tag || '$'} ${l.text}`).join('\n')
  try {
    navigator.clipboard?.writeText(text)
    systemStore.showNotification('Log Disalin', 'Seluruh log terminal uninstall berhasil disalin', 'info')
  } catch (e) {}
}

/**
 * Eksekusi proses uninstall bertahap dengan terminal log
 */
async function startUninstallProcess() {
  const slug = appSlug.value
  if (!slug) return

  uninstallState.value = 'running'
  uninstallProgress.value = 10
  stepIndex.value = 1
  errorMessage.value = ''
  terminalLogs.value = []

  appendLocalLog('info', '[INIT]', `Memulai proses uninstall aplikasi razos-${slug}...`)

  try {
    // TAHAP 1: Validasi Container
    currentStageLabel.value = 'Tahap 1: Validasi Container & Resource Host'
    currentStepDetail.value = `Memeriksa container razos-${slug}...`
    appendLocalLog('info', '[SYS]', `Cek status container razos-${slug} di Docker Engine...`)
    await sleep(400)
    uninstallProgress.value = 25

    // TAHAP 2: Stop Container
    stepIndex.value = 2
    currentStageLabel.value = 'Tahap 2: Menghentikan Container'
    currentStepDetail.value = `Mengirim perintah docker stop razos-${slug}...`
    appendLocalLog('docker', '[STOP]', `docker stop -t 5 razos-${slug}`)
    uninstallProgress.value = 45
    await sleep(400)

    // TAHAP 3 & 4: Panggil Endpoint DELETE ke Backend
    stepIndex.value = 3
    currentStageLabel.value = 'Tahap 3: Menghapus Container & Volume'
    currentStepDetail.value = `Menghapus container razos-${slug}...`
    uninstallProgress.value = 65

    // Mulai polling log server aktif dari backend
    pollTimer = setInterval(fetchBackendLogs, 400)

    const res = await fetch(`/api/v1/store/apps/${slug}?purge_data=${purgeData.value}`, {
      method: 'DELETE'
    })

    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    await fetchBackendLogs()

    if (!res.ok) {
      const errData = await res.json().catch(() => ({}))
      const msg = errData.error || `Gagal melakukan uninstall aplikasi ${slug}`
      appendLocalLog('error', '[ERR]', msg)
      throw new Error(msg)
    }

    // TAHAP 4: Pembersihan Data
    stepIndex.value = 4
    if (purgeData.value) {
      currentStageLabel.value = 'Tahap 4: Pembersihan Direktori Storage'
      currentStepDetail.value = `Membersihkan direktori /DATA/apps/${slug}...`
      appendLocalLog('info', '[STORAGE]', `rm -rf /DATA/apps/${slug}`)
    } else {
      currentStageLabel.value = 'Tahap 4: Direktori Data Dipertahankan'
      currentStepDetail.value = 'Opsi simpan data aktif.'
      appendLocalLog('info', '[STORAGE]', `Folder /DATA/apps/${slug} dipertahankan sesuai permintaan.`)
    }
    uninstallProgress.value = 85
    await sleep(300)

    // TAHAP 5: Finalisasi & Pelepasan Resource
    stepIndex.value = 5
    currentStageLabel.value = 'Tahap 5: Pelepasan Resource Host'
    currentStepDetail.value = 'Sinkronisasi status aplikasi sistem...'
    appendLocalLog('success', '[CLEAN]', 'Resource RAM & Network port telah berhasil dilepaskan.')
    uninstallProgress.value = 100
    stepIndex.value = 6
    await sleep(400)

    uninstallState.value = 'success'

    // Sinkronisasi sistem & taskbar
    systemStore.triggerAppsRefresh()
    await windowStore.fetchInstalledApps()
    systemStore.showNotification('Uninstall Selesai', `Aplikasi ${app.value.name || app.value.title || slug} berhasil dihapus dari sistem`, 'success')

  } catch (err) {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    uninstallState.value = 'error'
    errorMessage.value = err.message || 'Proses uninstall container mengalami kendala.'
    appendLocalLog('error', '[FATAL]', errorMessage.value)
  }
}

function closeWindow() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  windowStore.closeWindow(props.win?.id || `appstore-uninstall-${appSlug.value}`)
}

onMounted(() => {
  startUninstallProcess()
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
})
</script>
