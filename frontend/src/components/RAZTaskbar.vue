<script setup>
/**
 * RAZOS Taskbar Component
 * Updated: 2026-10-08
 * Deskripsi: Taskbar desktop bawah layar memuat Start Menu, active tasks, dan system tray.
 */
import { ref, computed } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import {
  Maximize2, Cpu, HardDrive, Thermometer, Clock, Minus, Square,
  X, Layers, EyeOff, Settings
} from 'lucide-vue-next'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()

const isStartMenuOpen = ref(false)

// Format Waktu Real-Time
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

function handleStartItemClick(appId) {
  isStartMenuOpen.value = false
  windowStore.openApp(appId)
}

function handleStartStoreAppClick(app) {
  isStartMenuOpen.value = false
  windowStore.launchStoreApp(app)
}

// Context menu klik kanan pada item taskbar aplikasi aktif
function handleTaskContextMenu(e, win) {
  contextMenuStore.show(e, [
    {
      label: win.isMinimized ? 'Pulihkan Jendela' : 'Fokus ke Depan',
      icon: Maximize2,
      action: () => {
        win.isMinimized = false
        windowStore.bringToFront(win.id)
      }
    },
    {
      label: 'Minimalkan',
      icon: Minus,
      disabled: win.isMinimized,
      action: () => windowStore.minimizeWindow(win.id)
    },
    {
      label: win.isMaximized ? 'Ukuran Normal' : 'Maksimalkan',
      icon: Square,
      action: () => windowStore.toggleMaximize(win.id)
    },
    { separator: true },
    {
      label: 'Tutup Jendela',
      icon: X,
      danger: true,
      action: () => windowStore.closeWindow(win.id)
    }
  ], win.title)
}

// Handler klik item taskbar dengan failsafe jika instance store lama
function handleTaskbarItemClick(winId) {
  if (typeof windowStore.toggleWindowFromTaskbar === 'function') {
    windowStore.toggleWindowFromTaskbar(winId)
  } else {
    const win = windowStore.openWindows.find(w => w.id === winId)
    if (win) {
      if (win.isMinimized) {
        win.isMinimized = false
        windowStore.bringToFront(winId)
      } else if (windowStore.activeWindowId === winId) {
        windowStore.minimizeWindow(winId)
      } else {
        windowStore.bringToFront(winId)
      }
    }
  }
}

// Context menu klik kanan pada area kosong taskbar (Hierarki Submenu)
function handleTaskbarContextMenu(e) {
  contextMenuStore.show(e, [
    {
      label: 'Task Manager',
      icon: Cpu,
      action: () => windowStore.openApp('task-manager')
    },
    {
      label: 'Jendela (Windows)',
      icon: Layers,
      children: [
        {
          label: 'Tampilkan Desktop',
          icon: EyeOff,
          shortcut: 'Win+D',
          action: () => {
            if (typeof windowStore.minimizeAllWindows === 'function') {
              windowStore.minimizeAllWindows()
            } else {
              windowStore.openWindows.forEach(w => { w.isMinimized = true })
              windowStore.activeWindowId = null
            }
          }
        },
        {
          label: 'Pulihkan Semua Jendela',
          icon: Layers,
          action: () => {
            if (typeof windowStore.restoreAllWindows === 'function') {
              windowStore.restoreAllWindows()
            } else {
              windowStore.openWindows.forEach(w => { w.isMinimized = false })
              if (windowStore.openWindows.length > 0) {
                windowStore.bringToFront(windowStore.openWindows[windowStore.openWindows.length - 1].id)
              }
            }
          }
        }
      ]
    },
    { separator: true },
    {
      label: 'Pengaturan Sistem',
      icon: Settings,
      action: () => windowStore.openApp('settings')
    }
  ], 'Menu Taskbar')
}
</script>

<template>
  <div
    class="fixed bottom-0 left-0 w-full h-12 raz-glass border-t border-slate-700/60 z-50 flex items-center justify-between px-3 select-none"
    @contextmenu.prevent="handleTaskbarContextMenu"
  >
    <!-- Left: Start Button -->
    <div class="relative flex items-center gap-2">
      <button
        @click="isStartMenuOpen = !isStartMenuOpen"
        class="h-9 px-3 flex items-center gap-2 rounded-lg bg-amber-500/10 hover:bg-amber-500/20 text-amber-300 hover:text-amber-200 font-semibold border border-amber-500/30 transition shadow-sm active:scale-95"
        title="Buka Menu Utama RAZOS"
      >
        <img src="/assets/images/RAZLogo.svg" alt="RAZOS" class="w-5 h-5 object-contain drop-shadow" />
        <span class="text-xs uppercase tracking-wider hidden sm:inline font-bold bg-gradient-to-r from-amber-200 to-yellow-400 bg-clip-text text-transparent">RAZOS</span>
      </button>

      <!-- Start Menu Popup -->
      <div
        v-if="isStartMenuOpen"
        @mouseleave="isStartMenuOpen = false"
        class="absolute bottom-12 left-0 w-64 raz-glass rounded-xl p-2 shadow-2xl border border-slate-700/80 mb-1 z-50 flex flex-col gap-1"
      >
        <div class="px-3 py-2 border-b border-slate-800 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <img src="/assets/brand/RAZOSMarkSquare.webp" alt="RAZOS" class="w-4 h-4 rounded object-cover shadow-sm" />
            <span class="text-xs font-bold text-amber-300 tracking-wider">APLIKASI SISTEM</span>
          </div>
          <div class="flex items-center gap-1 opacity-75 hover:opacity-100 transition" title="Created by RAZ Creative Studio">
            <span class="text-[9px] text-slate-400">by</span>
            <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ Creative Studio" class="h-3.5 object-contain" />
          </div>
        </div>
        <div class="flex flex-col gap-0.5">
          <button
            v-for="app in windowStore.apps"
            :key="app.id"
            @click="handleStartItemClick(app.id)"
            class="flex items-center gap-3 px-3 py-1.5 rounded-lg hover:bg-slate-800/80 text-left text-xs text-slate-200 transition"
          >
            <img :src="app.icon" :alt="app.title" class="w-4 h-4 object-contain" />
            <span>{{ app.title }}</span>
          </button>
        </div>

        <!-- Bagian Aplikasi Terpasang dari Store -->
        <template v-if="windowStore.installedStoreApps.length > 0">
          <div class="px-3 py-1.5 border-b border-t border-slate-800 flex items-center justify-between text-[11px] font-semibold text-slate-400 tracking-wider mt-1">
            <span>APLIKASI TERPASANG</span>
            <span class="text-[10px] font-mono text-emerald-400 font-bold">{{ windowStore.installedStoreApps.length }}</span>
          </div>
          <div class="max-h-40 overflow-y-auto flex flex-col gap-0.5 pr-1" style="scrollbar-width: thin;">
            <button
              v-for="app in windowStore.installedStoreApps"
              :key="app.slug || app.id"
              @click="handleStartStoreAppClick(app)"
              class="flex items-center justify-between px-3 py-1.5 rounded-lg hover:bg-slate-800/80 text-left text-xs text-slate-200 transition group"
            >
              <div class="flex items-center gap-2.5 truncate">
                <img :src="app.icon" :alt="app.name" class="w-4 h-4 object-contain rounded" />
                <span class="truncate font-medium">{{ app.name }}</span>
              </div>
              <span class="text-[10px] font-mono text-slate-400 group-hover:text-sky-300">:{{ app.port }}</span>
            </button>
          </div>
        </template>
      </div>

      <!-- Active Window Task Items -->
      <div class="flex items-center gap-1.5 ml-2 overflow-x-auto max-w-lg">
        <button
          v-for="win in windowStore.openWindows"
          :key="win.id"
          @click="handleTaskbarItemClick(win.id)"
          @contextmenu.prevent.stop="handleTaskContextMenu($event, win)"
          class="h-8.5 px-2.5 py-1 flex flex-col justify-center items-center rounded-lg text-xs font-medium transition-all duration-200 border relative group active:scale-95"
          :class="windowStore.activeWindowId === win.id && !win.isMinimized
            ? 'bg-slate-800/90 text-sky-300 border-sky-400/40 shadow-sm'
            : win.isMinimized
            ? 'bg-slate-900/40 text-slate-500 hover:text-slate-300 border-slate-800/40'
            : 'bg-slate-900/60 text-slate-300 hover:text-white border-slate-700/50'"
          :title="win.isMinimized ? `${win.title} (Diminimalkan - Klik untuk munculkan)` : `${win.title} (Klik untuk minimalkan)`"
        >
          <div class="flex items-center gap-1.5 pb-0.5">
            <img :src="win.icon" :alt="win.title" class="w-3.5 h-3.5 transition-transform duration-200 group-hover:scale-110" />
            <span class="truncate max-w-[100px]">{{ win.title }}</span>
          </div>

          <!-- Indikator Status Window (Glow Aktif / Terbuka / Minimized) -->
          <div
            class="h-0.5 rounded-full transition-all duration-300"
            :class="windowStore.activeWindowId === win.id && !win.isMinimized
              ? 'w-6 bg-sky-400 shadow-[0_0_8px_rgba(56,189,248,0.9)]'
              : win.isMinimized
              ? 'w-2 bg-slate-600'
              : 'w-3.5 bg-slate-400/70'"
          ></div>
        </button>
      </div>
    </div>

    <!-- Right: System Tray (CPU, RAM, Temp, Clock, Fullscreen) -->
    <div class="flex items-center gap-3 text-xs text-slate-300">
      <!-- CPU Usage -->
      <div class="flex items-center gap-1 hidden md:flex" title="Penggunaan CPU">
        <Cpu class="w-3.5 h-3.5 text-sky-400" />
        <span>{{ systemStore.metrics.cpu_percent.toFixed(0) }}%</span>
      </div>

      <!-- RAM Usage -->
      <div class="flex items-center gap-1 hidden md:flex" title="Penggunaan RAM">
        <HardDrive class="w-3.5 h-3.5 text-emerald-400" />
        <span>{{ systemStore.metrics.ram_percent.toFixed(0) }}%</span>
      </div>

      <!-- SOC Temp -->
      <div class="flex items-center gap-1 hidden sm:flex" title="Suhu SoC STB">
        <Thermometer class="w-3.5 h-3.5 text-amber-400" />
        <span>{{ systemStore.metrics.cpu_temp.toFixed(0) }}°C</span>
      </div>

      <!-- Fullscreen Toggle -->
      <button
        @click="toggleFullscreen"
        class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
        title="Toggle Fullscreen"
      >
        <Maximize2 class="w-3.5 h-3.5" />
      </button>

      <!-- Digital Clock -->
      <div class="flex items-center gap-1.5 px-2 py-1 rounded bg-slate-900/60 font-mono text-slate-200 border border-slate-800">
        <Clock class="w-3.5 h-3.5 text-sky-400" />
        <span>{{ currentTime }}</span>
      </div>
    </div>
  </div>
</template>
