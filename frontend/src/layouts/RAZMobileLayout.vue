<script setup>
/**
 * RAZOS Mobile Layout Component
 * Updated: 2026-10-08
 * Deskripsi: Layout Mobile Launcher bergaya Android 14 untuk breakpoint <= 768px.
 * Memuat status bar minimalis, widget metrik STB, grid 4-kolom, modal window fullscreen, dan dock bawah.
 */
import { ref } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Maximize2, Cpu, HardDrive, Thermometer, X, LayoutGrid } from 'lucide-vue-next'

// Import View Aplikasi
import RAZTaskManagerView from '../views/RAZTaskManagerView.vue'
import RAZTerminalView from '../views/RAZTerminalView.vue'
import RAZAppStoreView from '../views/RAZAppStoreView.vue'
import RAZWebsitesView from '../views/RAZWebsitesView.vue'
import RAZDatabasesView from '../views/RAZDatabasesView.vue'
import RAZSettingsView from '../views/RAZSettingsView.vue'
import RAZFileManagerView from '../views/RAZFileManagerView.vue'

const windowStore = useWindowStore()
const systemStore = useSystemStore()

const viewMap = {
  'file-manager': RAZFileManagerView,
  'task-manager': RAZTaskManagerView,
  'terminal': RAZTerminalView,
  'app-store': RAZAppStoreView,
  'websites': RAZWebsitesView,
  'databases': RAZDatabasesView,
  'settings': RAZSettingsView
}

// Active mobile modal app
const activeMobileAppId = ref(null)

function openMobileApp(appId) {
  activeMobileAppId.value = appId
}

function closeMobileApp() {
  activeMobileAppId.value = null
}

const currentTime = ref(new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }))
setInterval(() => {
  currentTime.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}, 1000)

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch(() => {})
  } else {
    document.exitFullscreen().catch(() => {})
  }
}
</script>

<template>
  <div
    class="relative w-screen h-screen overflow-hidden bg-cover bg-center flex flex-col justify-between p-4 select-none"
    style="background-image: url('/assets/images/RAZMobileBg.webp')"
  >
    <!-- Top Mobile Status Bar -->
    <div class="flex items-center justify-between text-xs text-slate-300 px-2 pt-1">
      <span class="font-mono font-bold">{{ currentTime }}</span>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1 text-[11px]">
          <Thermometer class="w-3.5 h-3.5 text-amber-400" />
          <span>{{ systemStore.metrics.cpu_temp.toFixed(0) }}°C</span>
        </div>
        <button @click="toggleFullscreen" class="p-1 hover:text-white">
          <Maximize2 class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- Center: Server Widget & App Grid -->
    <div class="flex-1 flex flex-col gap-6 mt-4 overflow-y-auto">
      <!-- Server Health Widget -->
      <div class="p-4 rounded-2xl raz-glass border border-slate-700/60 flex flex-col gap-3 shadow-xl">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <img src="/assets/images/RAZLogo.svg" alt="RAZOS" class="w-5 h-5" />
            <span class="text-xs font-bold text-sky-400 tracking-wider">RAZOS MINI SERVER</span>
          </div>
          <span class="text-[10px] px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 font-medium">Online</span>
        </div>

        <div class="grid grid-cols-2 gap-3 text-xs">
          <div class="flex items-center gap-2 p-2 rounded-xl bg-slate-900/60">
            <Cpu class="w-4 h-4 text-sky-400 shrink-0" />
            <div>
              <div class="text-[10px] text-slate-400">CPU Load</div>
              <div class="font-bold text-slate-200">{{ systemStore.metrics.cpu_percent.toFixed(0) }}%</div>
            </div>
          </div>
          <div class="flex items-center gap-2 p-2 rounded-xl bg-slate-900/60">
            <HardDrive class="w-4 h-4 text-emerald-400 shrink-0" />
            <div>
              <div class="text-[10px] text-slate-400">RAM Load</div>
              <div class="font-bold text-slate-200">{{ systemStore.metrics.ram_percent.toFixed(0) }}%</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Android Home Grid (4 Kolom) -->
      <div class="grid grid-cols-4 gap-4 px-2">
        <div
          v-for="app in windowStore.apps"
          :key="app.id"
          @click="openMobileApp(app.id)"
          class="flex flex-col items-center gap-1 cursor-pointer active:scale-95 transition-transform"
        >
          <div class="w-14 h-14 rounded-2xl raz-glass border border-slate-700/60 flex items-center justify-center p-2.5 shadow-lg">
            <img :src="app.icon" :alt="app.title" class="w-full h-full object-contain" />
          </div>
          <span class="text-[11px] text-slate-200 text-center line-clamp-1 mt-0.5">
            {{ app.title }}
          </span>
        </div>
      </div>
    </div>

    <!-- Bottom Dock Bar -->
    <div class="h-16 rounded-2xl raz-glass border border-slate-700/60 flex items-center justify-around px-4 shadow-2xl mb-1">
      <button @click="openMobileApp('terminal')" class="p-2 flex flex-col items-center">
        <img src="/assets/icons/RAZIconTerminal.svg" alt="Terminal" class="w-6 h-6" />
      </button>
      <button @click="openMobileApp('task-manager')" class="p-2 flex flex-col items-center">
        <img src="/assets/icons/RAZIconTaskManager.svg" alt="Task" class="w-6 h-6" />
      </button>
      <div class="w-10 h-10 rounded-full bg-sky-500/20 border border-sky-400/40 flex items-center justify-center">
        <img src="/assets/images/RAZLogo.svg" alt="Home" class="w-5 h-5" />
      </div>
      <button @click="openMobileApp('websites')" class="p-2 flex flex-col items-center">
        <img src="/assets/icons/RAZIconWebsites.svg" alt="Web" class="w-6 h-6" />
      </button>
      <button @click="openMobileApp('settings')" class="p-2 flex flex-col items-center">
        <img src="/assets/icons/RAZIconSettings.svg" alt="Settings" class="w-6 h-6" />
      </button>
    </div>

    <!-- Mobile Fullscreen App Modal -->
    <div
      v-if="activeMobileAppId"
      class="fixed inset-0 bg-slate-950 z-50 flex flex-col overflow-hidden animate-in fade-in duration-200"
    >
      <!-- Mobile Window Header -->
      <div class="h-12 border-b border-slate-800 bg-slate-900/90 flex items-center justify-between px-4">
        <div class="flex items-center gap-2">
          <img
            :src="windowStore.apps.find(a => a.id === activeMobileAppId)?.icon"
            class="w-5 h-5"
          />
          <span class="font-semibold text-sm text-slate-100">
            {{ windowStore.apps.find(a => a.id === activeMobileAppId)?.title }}
          </span>
        </div>
        <button
          @click="closeMobileApp"
          class="w-8 h-8 rounded-lg bg-slate-800 text-slate-300 flex items-center justify-center hover:bg-rose-600 transition"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Mobile Window Content -->
      <div class="flex-1 overflow-y-auto p-4 bg-slate-950 text-slate-100">
        <component :is="viewMap[activeMobileAppId]" />
      </div>
    </div>
  </div>
</template>
