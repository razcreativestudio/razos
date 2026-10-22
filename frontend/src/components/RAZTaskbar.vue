<script setup>
/**
 * RAZOS Taskbar Component with Windows 11 All Apps Start Menu & Windows 7 Glass Start Orb
 * Updated: 2026-10-22 (v3.7.5 - Windows 11 All-Apps Alphabetical Grouping & Interactive Folder Expansion)
 * Deskripsi: 
 * - Windows 7 Start Orb: Tombol bulat 3D Glass Orb bercahaya dengan logo emas RAZOS & efek gloss refleksi kaca.
 * - Windows 11 Start Menu (All Apps Layout Sesuai Referensi Pengguna):
 *   1. Search capsule bar di bagian atas ("Search for apps, settings, and documents") dengan ikon search cyan.
 *   2. Sub-header bar memuat judul "All" tebal dan tombol navigasi "< Back".
 *   3. Daftar vertikal A-Z rapi memuat aplikasi individual dan folder sistem dengan drop-down chevron interaktif.
 *   4. Bottom footer bar memuat avatar pengguna (Administrator / Active User) di sisi kiri dan tombol daya Power di sisi kanan.
 * - Zero Native Dialogs: Menggunakan custom modal & popover Glassmorphism Slate-900 terintegrasi.
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import { useAuthStore } from '../stores/RAZAuthStore'
import {
  Maximize2, Cpu, HardDrive, Thermometer, Clock, Minus, Square,
  X, Layers, EyeOff, Settings, Search, Power, RefreshCw, LogOut,
  ChevronRight, ChevronLeft, ChevronDown, ArrowLeft, Check, Sparkles, AlertTriangle, User,
  FileText, ShieldCheck, Folder, Terminal, Globe, Database,
  Users, CloudLightning, ArrowRightLeft, Shield
} from 'lucide-vue-next'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()
const authStore = useAuthStore()

// State Start Menu & Navigasi Internal (Default ke tampilan 'All Apps' sesuai referensi)
const isStartMenuOpen = ref(false)
const showAllAppsView = ref(true) // Default langsung tampil All Apps list
const startSearchQuery = ref('')
const showPowerMenu = ref(false)
const showUserMenu = ref(false)
const showConfirmPowerModal = ref(null) // null | 'restart' | 'shutdown' | 'logout'

// State Folder Expand / Collapse di All Apps List
const openFolders = ref({
  'accessibility': false,
  'antigravity': false,
  'anydesk': false
})

function toggleFolder(folderId) {
  openFolders.value[folderId] = !openFolders.value[folderId]
}

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
      showUserMenu.value = false
    }
  }
}

function toggleStartMenu() {
  isStartMenuOpen.value = !isStartMenuOpen.value
  if (isStartMenuOpen.value) {
    // Selalu buka mode 'All Apps' secara default sesuai permintaan
    showAllAppsView.value = true
  } else {
    showPowerMenu.value = false
    showUserMenu.value = false
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
  showUserMenu.value = false
  windowStore.openApp(appId)
}

function handleStartStoreAppClick(app) {
  isStartMenuOpen.value = false
  showPowerMenu.value = false
  showUserMenu.value = false
  windowStore.launchStoreApp(app)
}

// Folder Sistem Khas Windows 11 / RAZOS (Expandable with yellow folder icon)
const systemFolders = [
  {
    id: 'accessibility',
    title: 'Accessibility',
    letter: 'A',
    isFolder: true,
    items: [
      {
        id: 'acc-contrast',
        title: 'High Contrast Mode',
        icon: '/assets/icons/RAZIconSettings.svg',
        action: () => systemStore.showNotification('Aksesibilitas', 'Mode kontras tinggi aktif', 'info')
      },
      {
        id: 'acc-magnifier',
        title: 'Screen Magnifier',
        icon: '/assets/icons/RAZIconSettings.svg',
        action: () => systemStore.showNotification('Aksesibilitas', 'Fitur pembesar layar aktif', 'info')
      },
      {
        id: 'acc-keyboard',
        title: 'On-Screen Virtual Keyboard',
        icon: '/assets/icons/RAZIconSettings.svg',
        action: () => systemStore.showNotification('Aksesibilitas', 'Keyboard virtual siap digunakan', 'info')
      }
    ]
  },
  {
    id: 'antigravity',
    title: 'Antigravity',
    letter: 'A',
    isFolder: true,
    items: [
      {
        id: 'antigravity-code',
        title: 'Antigravity Code Editor',
        icon: '/assets/icons/RAZIconFileManager.svg',
        action: () => handleStartItemClick('file-manager')
      },
      {
        id: 'antigravity-telemetry',
        title: 'Antigravity Telemetry Host',
        icon: '/assets/icons/RAZIconTaskManager.svg',
        action: () => handleStartItemClick('task-manager')
      },
      {
        id: 'antigravity-terminal',
        title: 'Antigravity Root Console',
        icon: '/assets/icons/RAZIconTerminal.svg',
        action: () => handleStartItemClick('terminal')
      }
    ]
  },
  {
    id: 'anydesk',
    title: 'AnyDesk & Remote Tools',
    letter: 'A',
    isFolder: true,
    items: [
      {
        id: 'cloud-sync',
        title: 'Cloud Sync Engine (WebDAV / Nextcloud)',
        icon: '/assets/icons/RAZIconCloudSync.svg',
        action: () => handleStartItemClick('cloud-sync')
      },
      {
        id: 'websites',
        title: 'Website & Virtual Host Manager',
        icon: '/assets/icons/RAZIconWebsites.svg',
        action: () => handleStartItemClick('websites')
      },
      {
        id: 'databases',
        title: 'Database Studio (SQLite & MariaDB)',
        icon: '/assets/icons/RAZIconDatabases.svg',
        action: () => handleStartItemClick('databases')
      }
    ]
  }
]

// 1. ALL APPS DENGAN PENGELOMPOKAN ALFABETIK A-Z (PERSIS SEPERTI SCREENSHOT USER)
const groupedAlphabeticalApps = computed(() => {
  const query = startSearchQuery.value.trim().toLowerCase()

  // Ambil aplikasi bawaan RAZOS yang diizinkan oleh izin RBAC
  const coreApps = windowStore.apps
    .filter(app => authStore.canAccessApp(app.id))
    .map(app => ({
      id: app.id,
      title: app.title,
      icon: app.icon,
      isFolder: false,
      isStore: false,
      letter: (app.title.trim()[0] || '#').toUpperCase(),
      action: () => handleStartItemClick(app.id)
    }))

  // Ambil aplikasi Docker App Store yang terpasang
  const storeApps = windowStore.installedStoreApps
    .filter(s => authStore.canAccessApp('app-store') || authStore.canAccessApp(s.slug))
    .map(s => ({
      id: s.slug || s.id,
      title: s.name,
      icon: s.icon,
      isFolder: false,
      isStore: true,
      letter: (s.name.trim()[0] || '#').toUpperCase(),
      action: () => handleStartStoreAppClick(s)
    }))

  // Ambil folder sistem
  const folders = systemFolders.map(f => ({
    id: f.id,
    title: f.title,
    icon: '/assets/icons/RAZIconWinFolder.svg',
    isFolder: true,
    letter: f.letter || (f.title.trim()[0] || '#').toUpperCase(),
    items: f.items
  }))

  const allItems = [...coreApps, ...storeApps, ...folders]

  // Filter query jika pengguna sedang mengetik di search bar
  let filteredItems = allItems
  if (query) {
    filteredItems = allItems.filter(item => {
      if (item.title.toLowerCase().includes(query)) return true
      if (item.isFolder && item.items) {
        return item.items.some(sub => sub.title.toLowerCase().includes(query))
      }
      return false
    })
  }

  // Urutkan item berdasarkan abjad judul
  filteredItems.sort((a, b) => a.title.localeCompare(b.title))

  // Kelompokkan ke huruf abjad
  const groups = {}
  filteredItems.forEach(item => {
    const char = (item.title.trim()[0] || '#').toUpperCase()
    const letter = /[A-Z]/.test(char) ? char : '#'
    if (!groups[letter]) {
      groups[letter] = []
    }
    groups[letter].push(item)
  })

  // Format ke array terurut berdasarkan abjad
  return Object.keys(groups)
    .sort()
    .map(letter => ({
      letter,
      items: groups[letter]
    }))
})

// 2. PINNED APPS (Windows 11 Pinned Section - Backup jika user klik Back)
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

// 3. RECOMMENDED SECTION
const recommendedItems = computed(() => {
  const items = []
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

    <!-- WINDOWS 11 FLOATING START MENU (PERSIS SEPERTI SCREENSHOT USER) -->
    <div
      v-if="isStartMenuOpen"
      id="raz-win11-start-menu"
      @click.stop
      class="absolute bottom-14 left-2 sm:left-4 z-50 w-[94vw] sm:w-[490px] md:w-[510px] h-[610px] max-h-[85vh] rounded-2xl bg-[#1e1e1e]/95 backdrop-blur-2xl border border-white/10 shadow-[0_20px_50px_rgba(0,0,0,0.9)] p-4 sm:p-5 flex flex-col justify-between overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-200 text-slate-100 font-sans"
    >
      <!-- 1. TOP SEARCH CAPSULE (Windows 11 Search for apps, settings, and documents) -->
      <div class="relative w-full shrink-0 mb-3">
        <Search class="w-4 h-4 text-sky-400 absolute left-4 top-1/2 -translate-y-1/2 pointer-events-none" />
        <input
          v-model="startSearchQuery"
          type="text"
          placeholder="Search for apps, settings, and documents"
          class="w-full h-10 pl-11 pr-10 rounded-full bg-[#2a2a2c]/80 border border-white/5 text-xs text-slate-100 placeholder-slate-400 outline-none focus:border-sky-400/60 focus:ring-1 focus:ring-sky-400/30 transition shadow-inner"
        />
        <button
          v-if="startSearchQuery"
          @click="startSearchQuery = ''"
          class="absolute right-3.5 top-1/2 -translate-y-1/2 p-1 rounded-full hover:bg-white/10 text-slate-400 hover:text-white"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- 2. CONTENT AREA: ALL APPS ALPHABETICAL (DEFAULT) ATAU PINNED -->
      <div class="flex-1 overflow-y-auto pr-1 flex flex-col" style="scrollbar-width: thin;">
        
        <!-- MODE A: ALL APPS VIEW (Sesuai Screenshot User) -->
        <template v-if="showAllAppsView">
          <!-- Sub-Header: 'All' di kiri, '< Back' di kanan -->
          <div class="flex items-center justify-between px-2 pb-2 mb-1 shrink-0">
            <span class="text-sm font-semibold text-slate-100">All</span>
            <button
              @click="showAllAppsView = false"
              class="px-3 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-xs font-medium text-slate-200 hover:text-white flex items-center gap-1.5 transition active:scale-95 border border-white/5"
            >
              <ChevronLeft class="w-3.5 h-3.5" />
              <span>Back</span>
            </button>
          </div>

          <!-- Daftar Aplikasi Terkelompok Berdasarkan Huruf Abjad (A-Z) -->
          <div class="flex flex-col gap-1 pb-2">
            <div
              v-for="group in groupedAlphabeticalApps"
              :key="group.letter"
              class="flex flex-col"
            >
              <!-- Header Huruf (A, B, C, ...) -->
              <div class="text-xs font-semibold text-slate-400 px-3 pt-3 pb-1.5 select-none">
                {{ group.letter }}
              </div>

              <!-- Item List di bawah huruf -->
              <div class="flex flex-col gap-0.5">
                <template v-for="item in group.items" :key="item.id">
                  <!-- Tipe 1: Folder Sistem (Contoh: Accessibility, Antigravity, AnyDesk) -->
                  <div v-if="item.isFolder" class="flex flex-col">
                    <button
                      @click="toggleFolder(item.id)"
                      class="w-full flex items-center justify-between px-3 py-2 rounded-lg hover:bg-white/5 text-left transition group active:scale-[0.99]"
                    >
                      <div class="flex items-center gap-3 min-w-0">
                        <img
                          src="/assets/icons/RAZIconWinFolder.svg"
                          alt="Folder"
                          class="w-6 h-6 object-contain shrink-0 drop-shadow-sm"
                        />
                        <span class="text-xs sm:text-[13px] font-normal text-slate-200 truncate group-hover:text-white">
                          {{ item.title }}
                        </span>
                      </div>
                      <ChevronDown
                        class="w-4 h-4 text-slate-400 transition-transform duration-200 shrink-0"
                        :class="openFolders[item.id] ? 'rotate-180 text-slate-200' : ''"
                      />
                    </button>

                    <!-- Sub-items di dalam Folder jika Di-expand -->
                    <div
                      v-if="openFolders[item.id]"
                      class="pl-9 pr-1 py-1 flex flex-col gap-0.5 animate-in fade-in duration-150"
                    >
                      <button
                        v-for="sub in item.items"
                        :key="sub.id"
                        @click="sub.action"
                        class="w-full flex items-center gap-3 px-3 py-1.5 rounded-lg hover:bg-white/10 text-left transition group active:scale-[0.99]"
                      >
                        <img :src="sub.icon" :alt="sub.title" class="w-5 h-5 object-contain rounded-md shrink-0 opacity-90 group-hover:opacity-100" />
                        <span class="text-xs font-normal text-slate-300 truncate group-hover:text-white">
                          {{ sub.title }}
                        </span>
                      </button>
                    </div>
                  </div>

                  <!-- Tipe 2: Aplikasi Individual Biasa -->
                  <button
                    v-else
                    @click="item.action"
                    class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/5 text-left transition group active:scale-[0.99]"
                  >
                    <img
                      :src="item.icon"
                      :alt="item.title"
                      class="w-6 h-6 sm:w-7 sm:h-7 object-contain rounded-md shrink-0 transition group-hover:scale-105"
                    />
                    <span class="text-xs sm:text-[13px] font-normal text-slate-200 truncate group-hover:text-white">
                      {{ item.title }}
                    </span>
                  </button>
                </template>
              </div>
            </div>

            <!-- Jika tidak ada hasil pencarian -->
            <div
              v-if="groupedAlphabeticalApps.length === 0"
              class="py-12 text-center text-slate-400 text-xs flex flex-col items-center gap-2"
            >
              <Search class="w-8 h-8 opacity-30 text-slate-400" />
              <span>Tidak ada aplikasi yang cocok dengan "{{ startSearchQuery }}"</span>
            </div>
          </div>
        </template>

        <!-- MODE B: PINNED & RECOMMENDED VIEW (Jika User Klik '< Back') -->
        <template v-else>
          <!-- SECTION 1: PINNED APPS -->
          <div class="flex flex-col gap-2.5">
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

          <!-- SECTION 2: RECOMMENDED -->
          <div class="flex flex-col gap-2 pt-3 border-t border-white/10 mt-3">
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

      <!-- 3. BOTTOM FOOTER BAR (Akun Administrator & Tombol Power Khas Windows 11) -->
      <div class="pt-3 border-t border-white/10 flex items-center justify-between shrink-0 relative mt-1">
        <!-- Kiri: Akun Administrator Windows 11 (Avatar bulat + Nama pengguna) -->
        <div class="relative">
          <button
            @click="showUserMenu = !showUserMenu; showPowerMenu = false"
            class="flex items-center gap-3 px-2 py-1.5 rounded-xl hover:bg-white/10 transition group text-left"
            title="Pengaturan Akun Pengguna"
          >
            <img
              :src="authStore.currentUser?.avatar || systemStore.settings?.account?.avatar || '/assets/brand/RAZBrandSquare.webp'"
              alt="Avatar"
              class="w-8 h-8 rounded-full object-cover border border-white/20 shadow-sm transition group-hover:border-white/40"
            />
            <span class="text-xs sm:text-[13px] font-medium text-slate-200 group-hover:text-white transition">
              {{ authStore.currentUser?.name || authStore.currentUser?.username || systemStore.settings?.account?.name || 'Administrator' }}
            </span>
          </button>

          <!-- User Menu Popover -->
          <div
            v-if="showUserMenu"
            class="absolute bottom-12 left-0 w-56 rounded-2xl bg-[#1e1e1e]/98 backdrop-blur-2xl border border-white/15 p-1.5 shadow-2xl flex flex-col gap-0.5 z-50 animate-in fade-in zoom-in-95 duration-150 text-xs"
          >
            <div class="px-3 py-2 border-b border-white/10">
              <div class="font-semibold text-slate-100 truncate">{{ authStore.currentUser?.name || 'Administrator' }}</div>
              <div class="text-[10px] text-slate-400 font-mono capitalize">Peran: {{ authStore.currentUser?.role || 'admin' }}</div>
            </div>

            <button
              @click="handleStartItemClick('settings'); showUserMenu = false"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-slate-200 transition"
            >
              <Settings class="w-3.5 h-3.5 text-sky-400" />
              <span>Pengaturan Akun</span>
            </button>

            <button
              v-if="authStore.canManageUsers()"
              @click="windowStore.openUsersManager(); showUserMenu = false; isStartMenuOpen = false"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-white/10 text-left text-purple-300 transition"
            >
              <Users class="w-3.5 h-3.5 text-purple-400" />
              <span>Kelola Pengguna (RBAC)</span>
            </button>

            <div class="my-0.5 border-t border-white/10"></div>

            <button
              @click="showConfirmPowerModal = 'logout'; showUserMenu = false"
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl hover:bg-rose-500/20 text-left text-rose-300 transition"
            >
              <LogOut class="w-3.5 h-3.5 text-rose-400" />
              <span>Keluar Sesi (Logout)</span>
            </button>
          </div>
        </div>

        <!-- Kanan: Tombol Power Bulat Khas Windows 11 -->
        <div class="relative">
          <button
            @click="showPowerMenu = !showPowerMenu; showUserMenu = false"
            class="w-9 h-9 rounded-full hover:bg-white/10 flex items-center justify-center text-slate-200 hover:text-white transition active:scale-95"
            :class="showPowerMenu ? 'bg-white/15 text-white ring-2 ring-white/20' : ''"
            title="Daya Sistem (Power)"
          >
            <Power class="w-4 h-4" />
          </button>

          <!-- Power Menu Popover -->
          <div
            v-if="showPowerMenu"
            class="absolute bottom-12 right-0 w-48 rounded-2xl bg-[#1e1e1e]/98 backdrop-blur-2xl border border-white/15 p-1.5 shadow-2xl flex flex-col gap-0.5 z-50 animate-in fade-in zoom-in-95 duration-150"
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
