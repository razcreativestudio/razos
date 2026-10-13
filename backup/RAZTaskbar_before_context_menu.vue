<script setup>
/**
 * RAZOS Taskbar Component
 * Updated: 2026-10-08
 * Deskripsi: Taskbar desktop bawah layar memuat Start Menu, active tasks, dan system tray.
 */
import { ref, computed } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Maximize2, Cpu, HardDrive, Thermometer, Clock } from 'lucide-vue-next'

const windowStore = useWindowStore()
const systemStore = useSystemStore()

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
</script>

<template>
  <div class="fixed bottom-0 left-0 w-full h-12 raz-glass border-t border-slate-700/60 z-50 flex items-center justify-between px-3">
    <!-- Left: Start Button -->
    <div class="relative flex items-center gap-2">
      <button
        @click="isStartMenuOpen = !isStartMenuOpen"
        class="h-9 px-3 flex items-center gap-2 rounded-lg bg-sky-500/10 hover:bg-sky-500/20 text-sky-400 font-semibold border border-sky-400/30 transition shadow-sm"
      >
        <img src="/assets/images/RAZLogo.svg" alt="RAZOS" class="w-5 h-5" />
        <span class="text-xs uppercase tracking-wider hidden sm:inline">RAZOS</span>
      </button>

      <!-- Start Menu Popup -->
      <div
        v-if="isStartMenuOpen"
        @mouseleave="isStartMenuOpen = false"
        class="absolute bottom-12 left-0 w-64 raz-glass rounded-xl p-2 shadow-2xl border border-slate-700/80 mb-1 z-50 flex flex-col gap-1"
      >
        <div class="px-3 py-2 border-b border-slate-800 flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-400 tracking-wider">APLIKASI SISTEM</span>
          <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ" class="h-4 object-contain" />
        </div>
        <button
          v-for="app in windowStore.apps"
          :key="app.id"
          @click="handleStartItemClick(app.id)"
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-slate-800/80 text-left text-sm text-slate-200 transition"
        >
          <img :src="app.icon" :alt="app.title" class="w-5 h-5" />
          <span>{{ app.title }}</span>
        </button>
      </div>

      <!-- Active Window Task Items -->
      <div class="flex items-center gap-1.5 ml-2 overflow-x-auto max-w-lg">
        <button
          v-for="win in windowStore.openWindows"
          :key="win.id"
          @click="win.isMinimized ? (win.isMinimized = false, windowStore.bringToFront(win.id)) : windowStore.bringToFront(win.id)"
          class="h-8 px-2.5 flex items-center gap-2 rounded-lg text-xs font-medium transition border"
          :class="windowStore.activeWindowId === win.id && !win.isMinimized
            ? 'bg-slate-800 text-sky-400 border-sky-400/40 shadow'
            : 'bg-slate-900/50 text-slate-400 hover:text-slate-200 border-slate-700/40'"
        >
          <img :src="win.icon" :alt="win.title" class="w-3.5 h-3.5" />
          <span class="truncate max-w-[100px]">{{ win.title }}</span>
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
