<script setup>
/**
 * RAZOS Taskbar Component with Windows 11 Start Menu & Windows 7 Glass Start Orb
 * Updated: 2026-10-15 (v3.6.5 - Windows 11 Enterprise Start Menu & Windows 7 Orb Architecture)
 * Deskripsi: 
 * - Windows 7 Start Orb: Tombol bulat 3D Glass Orb bercahaya dengan logo emas RAZOS & efek gloss refleksi kaca.
 * - Windows 11 Start Menu: Menu mengambang rounded-3xl glassmorphism berisikan:
 *   1. Search capsule bar di bagian atas ("Search for apps, settings, and documents")
 *   2. Section "Pinned" aplikasi dengan grid 6 kolom & tombol "All >" untuk melihat seluruh aplikasi
 *   3. Section "Recommended" aplikasi & berkas terbaru dengan grid 2 kolom
 *   4. Bottom footer bar memuat avatar pengguna (Administrator), kredit RAZ Creative Studio, dan tombol daya Power.
 * - Zero Native Dialogs: Menggunakan custom popover & glassmorphism modal untuk aksi restart, shutdown, dan logout.
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import { useAuthStore } from '../stores/RAZAuthStore'
import {
  Maximize2, Cpu, HardDrive, Thermometer, Clock, Minus, Square,
  X, Layers, EyeOff, Settings, Search, Power, RefreshCw, LogOut,
  ChevronRight, ArrowLeft, Check, Sparkles, AlertTriangle, User,
  FileText, ShieldCheck, Folder, Terminal, Globe, Database,
  Users, CloudLightning, ArrowRightLeft
} from 'lucide-vue-next'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()
const authStore = useAuthStore()

// State Start Menu & Navigasi Internal
const isStartMenuOpen = ref(false)
const showAllAppsView = ref(false)
const startSearchQuery = ref('')
const showPowerMenu = ref(false)
const showConfirmPowerModal = ref(null) // null | 'restart' | 'shutdown' | 'logout'

// Format Waktu Real-Time Taskbar
const currentTime = ref(new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }))
let timeInterval = null

onMounted(() => {
  timeInterval = setInterval(() => {
    currentTime.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
  }, 1000)
  window.addEventListener('click', handleGlobalClick)
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
  window.removeEventListener('click', handleGlobalClick)
})

// Tutup Start Menu jika klik di luar
function handleGlobalClick(e) {
  const startMenuEl = document.getElementById('raz-win11-start-menu')
  const startOrbEl = document.getElementById('raz-win7-start-orb')
  if (isStartMenuOpen.value && startMenuEl && startOrbEl) {
    if (!startMenuEl.contains(e.target) && !startOrbEl.contains(e.target)) {
      isStartMenuOpen.value = false
      showPowerMenu.value = false
    }
  }
}

function toggleStartMenu() {
  isStartMenuOpen.value = !isStartMenuOpen.value
  if (!isStartMenuOpen.value) {
    showPowerMenu.value = false
    startSearchQuery.value = ''
  }
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch(() => {})
  } else {
    document.exitFullscreen().catch(() => {})
  }
}

function handleStartItemClick(appId) {
  isStartMenuOpen.value = false
  showPowerMenu.value = false
  windowStore.openApp(appId)
}

function handleStartStoreAppClick(app) {
  isStartMenuOpen.value = false
  showPowerMenu.value = false
  windowStore.launchStoreApp(app)
}

// 1. PINNED APPS (Windows 11 Pinned Section - Filtered by RBAC)
const pinnedApps = computed(() => {
  const allowed = windowStore.apps.filter(app => authStore.canAccessApp(app.id))
  const query = startSearchQuery.value.trim().toLowerCase()
  if (!query) {
    return allowed
  }
  return allowed.filter(app => 
    (app.title && app.title.toLowerCase().includes(query)) ||
    (app.id && app.id.toLowerCase().includes(query))
  )
})

// 2. ALL APPS (Alphabetical View A-Z - Filtered by RBAC)
const allAppsAlphabetical = computed(() => {
  const allowedCore = windowStore.apps.filter(app => authStore.canAccessApp(app.id))
  const combined = [
    ...allowedCore.map(a => ({ ...a, isStore: false })),
    ...windowStore.installedStoreApps.filter(s => authStore.canAccessApp('app-store') || authStore.canAccessApp(s.slug)).map(s => ({
      id: s.slug || s.id,
      title: s.name,
      icon: s.icon,
      isStore: true,
      raw: s
    }))
  ]

  const query = startSearchQuery.value.trim().toLowerCase()
  const filtered = query 
    ? combined.filter(item => item.title.toLowerCase().includes(query))
    : combined

  return [...filtered].sort((a, b) => a.title.localeCompare(b.title))
})

// 3. RECOMMENDED SECTION (Windows 11 Style)
const recommendedItems = computed(() => {
  const items = []

  // Tambahkan aplikasi baru dari store jika ada
  if (windowStore.installedStoreApps.length > 0) {
    const latest = windowStore.installedStoreApps.slice(0, 2)
    latest.forEach(app => {
      items.push({
        id: app.slug || app.id,
        title: app.name,
        subtitle: 'Baru ditambahkan',
        icon: app.icon,
        action: () => handleStartStoreAppClick(app)
      })
    })
  }

  // Rekomendasi bawaan yang esensial
  const recApps = [
    {
      id: 'file-manager',
      title: 'File Manager & Explorer',
      subtitle: 'Baru saja dibuka',
      icon: '/assets/icons/RAZIconFileManager.svg'
    },
    {
      id: 'settings',
      title: 'Pengaturan Sistem',
      subtitle: 'Host & Jaringan',
      icon: '/assets/icons/RAZIconSettings.svg'
    },
    {
      id: 'terminal',
      title: 'Terminal Konsol',
      subtitle: 'Sesi Root Aktif',
      icon: '/assets/icons/RAZIconTerminal.svg'
    },
    {
      id: 'task-manager',
      title: 'Task Manager Host',
      subtitle: 'CPU & RAM Monitor',
      icon: '/assets/icons/RAZIconTaskManager.svg'
    }
  ]

  for (const app of recApps) {
    if (items.length >= 6) break
    items.push({
      id: app.id,
      title: app.title,
      subtitle: app.subtitle,
      icon: app.icon,
      action: () => handleStartItemClick(app.id)
    })
  }

  const query = startSearchQuery.value.trim().toLowerCase()
  if (!query) return items
  return items.filter(i => i.title.toLowerCase().includes(query))
})

// Eksekusi Daya (Restart / Shutdown / Logout)
async function executePowerAction() {
  const action = showConfirmPowerModal.value
  showConfirmPowerModal.value = null
  showPowerMenu.value = false
  isStartMenuOpen.value = false

  if (action === 'restart') {
    await systemStore.restartSystem()
  } else if (action === 'shutdown') {
    await systemStore.shutdownSystem()
  } else if (action === 'logout') {
    await authStore.logout()
    systemStore.showNotification('Sesi Berakhir', 'Anda telah berhasil keluar dari akun RAZOS.', 'info')
  }
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
    <!-- Left: Windows 7 Start Orb Button -->
    <div class="relative flex items-center gap-3">
      <!-- WINDOWS 7 3D GLASS START ORB -->
      <button
        id="raz-win7-start-orb"
        @click.stop="toggleStartMenu"
        class="w-10 h-10 sm:w-11 sm:h-11 rounded-full flex items-center justify-center relative cursor-pointer select-none transition-all duration-300 active:scale-90 group focus:outline-none"
        :class="isStartMenuOpen
          ? 'scale-95 shadow-[0_0_20px_rgba(251,191,36,0.85)] border-2 border-amber-300 ring-2 ring-amber-400/50'
          : 'hover:scale-105 hover:shadow-[0_0_18px_rgba(245,158,11,0.7)] border-2 border-amber-400/60 shadow-lg shadow-black/80'"
        title="Start Menu RAZOS"
      >
        <!-- Background Radial Orb Gradient -->
        <div class="absolute inset-0 rounded-full bg-gradient-to-b from-slate-900 via-amber-950/40 to-black overflow-hidden">
          <!-- Upper Half Glass Highlight (Pantulan Kaca Khas Windows 7) -->
          <div class="absolute top-0 inset-x-1.5 h-4.5 rounded-t-full bg-gradient-to-b from-white/45 via-white/15 to-transparent pointer-events-none"></div>
          <!-- Inner Ambient Glow Saat Hover -->
          <div class="absolute inset-0 rounded-full bg-amber-400/0 group-hover:bg-amber-400/20 transition duration-300"></div>
        </div>

        <!-- Golden RAZOS Cloud Knot Logo di Tengah Orb -->
        <img
          src="/assets/images/RAZLogo.svg"
          alt="RAZOS Start"
          class="relative w-6 h-6 object-contain drop-shadow-[0_2px_5px_rgba(0,0,0,0.9)] transition-transform duration-300 group-hover:scale-110"
        />
      </button>

      <!-- Active Window Task Items -->
      <div class="flex items-center gap-1.5 overflow-x-auto max-w-lg">
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

    <!-- WINDOWS 11 FLOATING START MENU -->
    <div
      v-if="isStartMenuOpen"
      id="raz-win11-start-menu"
      @click.stop
      class="absolute bottom-14 left-2 sm:left-4 z-50 w-[94vw] sm:w-[580px] md:w-[620px] h-[590px] max-h-[85vh] rounded-3xl bg-slate-900/95 backdrop-blur-2xl border border-white/10 shadow-2xl shadow-black/90 p-5 sm:p-6 flex flex-col justify-between overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-200 text-slate-100 font-sans"
    >
      <!-- 1. TOP SEARCH CAPSULE (Windows 11 Search for apps, settings, and documents) -->
      <div class="relative w-full shrink-0 mb-4">
        <Search class="w-4 h-4 text-sky-400 absolute left-4 top-1/2 -translate-y-1/2 pointer-events-none" />
        <input
          v-model="startSearchQuery"
          type="text"
          placeholder="Search for apps, settings, and documents"
          class="w-full h-11 pl-11 pr-10 rounded-full bg-slate-950/80 border border-white/10 text-xs text-slate-100 placeholder-slate-400 outline-none focus:border-amber-400/60 focus:ring-2 focus:ring-amber-500/20 transition shadow-inner"
        />
        <button
          v-if="startSearchQuery"
          @click="startSearchQuery = ''"
          class="absolute right-3.5 top-1/2 -translate-y-1/2 p-1 rounded-full hover:bg-white/10 text-slate-400 hover:text-white"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- 2. CONTENT AREA: PINNED & RECOMMENDED atau ALL APPS LIST -->
      <div class="flex-1 overflow-y-auto pr-1 flex flex-col gap-5" style="scrollbar-width: thin;">
        
        <!-- MODE A: TAMPILAN ALL APPS (Saat klik 'All >') -->
        <template v-if="showAllAppsView">
          <div class="flex items-center justify-between pb-1 border-b border-white/10">
            <button
              @click="showAllAppsView = false"
              class="flex items-center gap-1.5 text-xs font-semibold text-sky-400 hover:text-sky-300 transition"
            >
              <ArrowLeft class="w-3.5 h-3.5" />
              <span>Kembali ke Pinned</span>
            </button>
            <span class="text-xs font-bold text-slate-300">Semua Aplikasi ({{ allAppsAlphabetical.length }})</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-1.5">
            <button
              v-for="item in allAppsAlphabetical"
              :key="item.id"
              @click="item.isStore ? handleStartStoreAppClick(item.raw) : handleStartItemClick(item.id)"
              class="flex items-center gap-3 p-2 rounded-xl hover:bg-white/10 text-left transition group active:scale-98"
            >
              <img :src="item.icon" :alt="item.title" class="w-7 h-7 object-contain rounded-lg shrink-0 group-hover:scale-110 transition" />
              <div class="min-w-0 flex-1">
                <div class="text-xs font-medium text-slate-200 truncate">{{ item.title }}</div>
                <div class="text-[10px] text-slate-400 font-mono">{{ item.isStore ? 'Store App' : 'System App' }}</div>
              </div>
            </button>
          </div>
        </template>

        <!-- MODE B: STANDARD WINDOWS 11 (PINNED + RECOMMENDED) -->
        <template v-else>
          <!-- SECTION 1: PINNED APPS -->
          <div class="flex flex-col gap-2.5">
            <!-- Pinned Header -->
            <div class="flex items-center justify-between px-1">
              <span class="text-xs font-bold text-slate-200 tracking-wide">Pinned</span>
              <button
                @click="showAllAppsView = true"
                class="px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-[11px] font-semibold text-slate-300 hover:text-white flex items-center gap-1 transition"
              >
                <span>All</span>
                <ChevronRight class="w-3.5 h-3.5" />
              </button>
            </div>

            <!-- Pinned Grid 6 Kolom Sesuai Screenshot Windows 11 -->
            <div class="grid grid-cols-3 sm:grid-cols-6 gap-2">
              <button
                v-for="app in pinnedApps"
                :key="app.id"
                @click="handleStartItemClick(app.id)"
                class="p-2 rounded-2xl hover:bg-white/10 flex flex-col items-center justify-center gap-2 transition duration-150 active:scale-95 group text-center"
                :title="app.title"
              >
                <div class="w-10 h-10 flex items-center justify-center group-hover:scale-110 transition duration-200">
                  <img :src="app.icon" :alt="app.title" class="w-9 h-9 object-contain drop-shadow" />
                </div>
                <span class="text-[11px] font-medium text-slate-200 truncate w-full group-hover:text-amber-300 transition">
                  {{ app.title }}
                </span>
              </button>

              <!-- Jika ada installed store apps, tambahkan ke pinned grid -->
              <button
                v-for="app in windowStore.installedStoreApps.slice(0, 3)"
                :key="app.slug || app.id"
                @click="handleStartStoreAppClick(app)"
                class="p-2 rounded-2xl hover:bg-white/10 flex flex-col items-center justify-center gap-2 transition duration-150 active:scale-95 group text-center"
                :title="app.name"
              >
                <div class="w-10 h-10 flex items-center justify-center group-hover:scale-110 transition duration-200">
                  <img :src="app.icon" :alt="app.name" class="w-9 h-9 object-contain rounded-xl drop-shadow" />
                </div>
                <span class="text-[11px] font-medium text-slate-200 truncate w-full group-hover:text-sky-300 transition">
                  {{ app.name }}
                </span>
              </button>
            </div>
          </div>

          <!-- SECTION 2: RECOMMENDED (Sesuai Screenshot Windows 11) -->
          <div class="flex flex-col gap-2 pt-2 border-t border-white/10">
            <!-- Recommended Header -->
            <div class="flex items-center justify-between px-1">
              <span class="text-xs font-bold text-slate-200 tracking-wide">Recommended</span>
              <button
                @click="showAllAppsView = true"
                class="px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-[11px] font-semibold text-slate-300 hover:text-white flex items-center gap-1 transition"
              >
                <span>More</span>
                <ChevronRight class="w-3.5 h-3.5" />
              </button>
            </div>

            <!-- Recommended 2 Kolom Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
              <button
                v-for="item in recommendedItems"
                :key="item.id"
                @click="item.action"
                class="px-2.5 py-2 rounded-2xl hover:bg-white/10 flex items-center gap-3 transition active:scale-98 text-left group"
                :title="item.title"
              >
                <div class="w-9 h-9 rounded-xl bg-slate-800/80 p-1.5 flex items-center justify-center shrink-0 border border-white/5 group-hover:scale-105 transition shadow-sm">
                  <img :src="item.icon" :alt="item.title" class="w-6 h-6 object-contain" />
                </div>
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-semibold text-slate-200 truncate group-hover:text-amber-300 transition">
                    {{ item.title }}
                  </div>
                  <div class="text-[10px] text-slate-400 truncate">
                    {{ item.subtitle }}
                  </div>
                </div>
              </button>
            </div>
          </div>
        </template>
      </div>

      <!-- 3. BOTTOM FOOTER BAR (Akun Administrator, RAZ Brand, dan Tombol Power) -->
      <div class="pt-4 border-t border-white/10 flex items-center justify-between shrink-0 relative">
        <!-- Kiri: Akun Administrator Windows 11 -->
        <button
          @click="handleStartItemClick('settings')"
          class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-xl hover:bg-white/10 transition group"
          title="Pengaturan Akun Pengguna"
        >
          <img
            :src="systemStore.settings?.account?.avatar || '/assets/brand/RAZBrandSquare.webp'"
            alt="Avatar"
            class="w-8 h-8 rounded-full object-cover border border-amber-500/40 shadow-sm group-hover:border-amber-400 transition"
          />
          <div class="text-left">
            <span class="text-xs font-semibold text-slate-100 block group-hover:text-amber-300 transition">
              {{ authStore.currentUser?.username || systemStore.settings?.account?.name || 'Administrator' }}
            </span>
            <span class="text-[9px] text-emerald-400 font-mono block leading-none">● Online</span>
          </div>
        </button>

        <!-- Tengah: Brand RAZ Creative Studio Tetap Ada -->
        <div class="hidden sm:flex items-center gap-1.5 opacity-60 hover:opacity-100 transition px-2" title="Created by RAZ Creative Studio">
          <span class="text-[9px] text-slate-400">by</span>
          <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ Creative Studio" class="h-3.5 object-contain" />
        </div>

        <!-- Kanan: Tombol Power Windows 11 -->
        <div class="relative">
          <button
            @click="showPowerMenu = !showPowerMenu"
            class="w-9 h-9 rounded-full hover:bg-white/10 flex items-center justify-center text-slate-300 hover:text-amber-300 transition active:scale-95"
            :class="showPowerMenu ? 'bg-white/10 text-amber-300 ring-2 ring-amber-400/30' : ''"
            title="Daya Sistem (Power)"
          >
            <Power class="w-4 h-4" />
          </button>

          <!-- Power Menu Popover -->
          <div
            v-if="showPowerMenu"
            class="absolute bottom-12 right-0 w-48 rounded-2xl bg-slate-950/95 backdrop-blur-xl border border-white/15 p-1.5 shadow-2xl flex flex-col gap-0.5 z-50 animate-in fade-in zoom-in-95 duration-150"
          >
            <button
              v-if="authStore.canManageUsers()"
              @click="windowStore.openUsersManager(); showPowerMenu = false; isStartMenuOpen = false"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-xs text-purple-300 transition"
            >
              <Users class="w-3.5 h-3.5 text-purple-400" />
              <span>Kelola Pengguna & RBAC</span>
            </button>
            <button
              v-if="authStore.canMountCloud()"
              @click="windowStore.openCloudSync(); showPowerMenu = false; isStartMenuOpen = false"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-xs text-sky-300 transition"
            >
              <CloudLightning class="w-3.5 h-3.5 text-sky-400" />
              <span>Sinkronisasi Cloud</span>
            </button>
            <div v-if="authStore.canManageUsers() || authStore.canMountCloud()" class="my-1 border-t border-white/10"></div>
            <button
              @click="showConfirmPowerModal = 'restart'"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-xs text-slate-200 transition"
            >
              <RefreshCw class="w-3.5 h-3.5 text-amber-400" />
              <span>Mulai Ulang (Restart)</span>
            </button>
            <button
              @click="showConfirmPowerModal = 'shutdown'"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-rose-500/20 text-left text-xs text-rose-300 transition"
            >
              <Power class="w-3.5 h-3.5 text-rose-400" />
              <span>Matikan (Shutdown)</span>
            </button>
            <div class="my-1 border-t border-white/10"></div>
            <button
              @click="showConfirmPowerModal = 'logout'"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-xs text-slate-300 transition"
            >
              <LogOut class="w-3.5 h-3.5 text-slate-400" />
              <span>Keluar Sesi (Logout)</span>
            </button>
          </div>
        </div>
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
      <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-slate-900/80 font-mono text-slate-200 border border-slate-700/60 shadow-sm">
        <Clock class="w-3.5 h-3.5 text-sky-400" />
        <span>{{ currentTime }}</span>
      </div>
    </div>

    <!-- CUSTOM POWER CONFIRMATION MODAL (Zero Native Dialogs) -->
    <div
      v-if="showConfirmPowerModal"
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
      @click.self="showConfirmPowerModal = null"
    >
      <div class="w-full max-w-sm rounded-3xl bg-slate-900 border border-slate-700 shadow-2xl p-6 flex flex-col gap-4 font-sans text-slate-100">
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-2xl flex items-center justify-center shadow-lg"
            :class="showConfirmPowerModal === 'shutdown' ? 'bg-rose-500/20 text-rose-400 border border-rose-500/40' : 'bg-amber-500/20 text-amber-400 border border-amber-500/40'"
          >
            <AlertTriangle class="w-5 h-5" />
          </div>
          <div>
            <h4 class="font-bold text-sm text-white">
              {{ showConfirmPowerModal === 'restart' ? 'Mulai Ulang Sistem (Restart)' : showConfirmPowerModal === 'shutdown' ? 'Matikan Sistem (Shutdown)' : 'Keluar Sesi (Logout)' }}
            </h4>
            <p class="text-[11px] text-slate-400">Konfirmasi tindakan daya RAZOS</p>
          </div>
        </div>

        <p class="text-xs text-slate-300 leading-relaxed">
          {{ showConfirmPowerModal === 'restart'
            ? 'Apakah Anda yakin ingin me-restart sistem host RAZOS sekarang?'
            : showConfirmPowerModal === 'shutdown'
            ? 'Apakah Anda yakin ingin mematikan daya sistem host RAZOS?'
            : 'Apakah Anda ingin mengakhiri sesi login RAZOS saat ini?'
          }}
        </p>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showConfirmPowerModal = null"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-300 transition"
          >
            Batal
          </button>
          <button
            @click="executePowerAction"
            class="px-4 py-2 rounded-xl text-xs font-bold transition shadow-md"
            :class="showConfirmPowerModal === 'shutdown' ? 'bg-rose-500 hover:bg-rose-400 text-white' : 'bg-amber-500 hover:bg-amber-400 text-slate-950'"
          >
            Konfirmasi
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Glassmorphism taskbar styling */
.raz-glass {
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}
</style>
