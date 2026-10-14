<script setup>
/**
 * RAZOS Desktop Layout Component
 * Updated: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Layout Desktop mode untuk breakpoint > 768px, memadukan desktop icon shortcuts,
 * window manager frame multi-jendela, dan taskbar bawah dengan sistem jendela mandiri tanpa background blur.
 */
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import RAZWindowFrame from '../components/RAZWindowFrame.vue'
import RAZTaskbar from '../components/RAZTaskbar.vue'
import {
  RefreshCw, Folder, Terminal, Cpu, Settings, EyeOff, Layers, Play, Info,
  Trash2, ExternalLink, AlertTriangle, X, Globe, Clock, Activity, HardDrive,
  Network, Search, Grid, Maximize2, Minimize2
} from 'lucide-vue-next'

// Import View Aplikasi
import RAZTaskManagerView from '../views/RAZTaskManagerView.vue'
import RAZTerminalView from '../views/RAZTerminalView.vue'
import RAZAppStoreView from '../views/RAZAppStoreView.vue'
import RAZWebsitesView from '../views/RAZWebsitesView.vue'
import RAZDatabasesView from '../views/RAZDatabasesView.vue'
import RAZSettingsView from '../views/RAZSettingsView.vue'
import RAZFileManagerView from '../views/RAZFileManagerView.vue'
import RAZFileViewerView from '../views/RAZFileViewerView.vue'

// Import View Jendela Khusus / Standalone Child Windows (v2.6.0)
import RAZWebsiteWizardView from '../views/RAZWebsiteWizardView.vue'
import RAZWebsiteEditView from '../views/RAZWebsiteEditView.vue'
import RAZWebsiteGitView from '../views/RAZWebsiteGitView.vue'
import RAZWebsiteDetailView from '../views/RAZWebsiteDetailView.vue'
import RAZFilePropertiesView from '../views/RAZFilePropertiesView.vue'
import RAZDatabaseCreatorView from '../views/RAZDatabaseCreatorView.vue'
import RAZAppStoreManualView from '../views/RAZAppStoreManualView.vue'
import RAZAppStoreInstallView from '../views/RAZAppStoreInstallView.vue'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()

// Mapping ID aplikasi ke komponen View
const viewMap = {
  'file-manager': RAZFileManagerView,
  'task-manager': RAZTaskManagerView,
  'terminal': RAZTerminalView,
  'app-store': RAZAppStoreView,
  'websites': RAZWebsitesView,
  'databases': RAZDatabasesView,
  'settings': RAZSettingsView
}

// Context Menu Klik Kanan pada Area Desktop Kosong
function handleDesktopContextMenu(e) {
  contextMenuStore.show(e, [
    // Sub-opsi Pengaturan Ukuran Ikon Desktop (Windows Style)
    {
      label: 'Ikon Besar (Large)',
      icon: Maximize2,
      checked: desktopIconSize.value === 'large',
      action: () => changeIconSize('large')
    },
    {
      label: 'Ikon Sedang (Medium)',
      icon: Grid,
      checked: desktopIconSize.value === 'medium',
      action: () => changeIconSize('medium')
    },
    {
      label: 'Ikon Kecil (Small)',
      icon: Minimize2,
      checked: desktopIconSize.value === 'small',
      action: () => changeIconSize('small')
    },
    {
      label: 'Ratakan Ikon ke Grid',
      icon: RefreshCw,
      action: () => autoArrangeIcons()
    },
    { separator: true },
    {
      label: 'Muat Ulang Desktop',
      icon: RefreshCw,
      action: () => {
        systemStore.fetchMetrics()
        fetchDesktopConfig()
        systemStore.showNotification('Desktop', 'Tampilan desktop telah disegarkan', 'success')
      }
    },
    { separator: true },
    {
      label: 'Files & Storage',
      icon: Folder,
      action: () => windowStore.openApp('file-manager')
    },
    {
      label: 'Web Terminal',
      icon: Terminal,
      action: () => windowStore.openApp('terminal')
    },
    {
      label: 'Task Manager',
      icon: Cpu,
      action: () => windowStore.openApp('task-manager')
    },
    { separator: true },
    {
      label: 'Tampilkan Desktop',
      icon: EyeOff,
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
    },
    { separator: true },
    {
      label: 'Pengaturan Sistem',
      icon: Settings,
      action: () => windowStore.openApp('settings')
    }
  ], 'Desktop')
}

// State Modal Konfirmasi Uninstall dari Desktop
const showUninstallModal = ref(false)
const appToUninstall = ref(null)
const purgeDataChecked = ref(true)

function openDesktopUninstall(app) {
  appToUninstall.value = app
  purgeDataChecked.value = true
  showUninstallModal.value = true
}

async function confirmDesktopUninstall() {
  if (!appToUninstall.value) return
  const slug = appToUninstall.value.slug
  const success = await windowStore.uninstallStoreApp(slug, purgeDataChecked.value)
  if (success) {
    systemStore.showNotification('Sukses', `Aplikasi ${appToUninstall.value.title} berhasil di-uninstall`, 'success')
    showUninstallModal.value = false
    appToUninstall.value = null
  } else {
    systemStore.showNotification('Gagal', 'Terjadi kesalahan saat menghapus aplikasi', 'error')
  }
}

// Context Menu Klik Kanan pada Ikon Shortcut Desktop
function handleDesktopIconContextMenu(e, app) {
  if (app.isWebsite) {
    // Context Menu untuk Website yang Dibuat di Website Manager
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.launchWebsiteApp(app)
      },
      {
        label: 'Buka di Tab Baru Peramban',
        icon: ExternalLink,
        action: () => windowStore.launchWebsiteApp(app)
      },
      { separator: true },
      {
        label: 'Buka Root Folder',
        icon: Folder,
        action: () => {
          windowStore.openApp('file-manager')
          systemStore.showNotification('Root Folder', `Lokasi: ${app.root_folder}`, 'info')
        }
      },
      {
        label: 'Kelola di Website Manager',
        icon: Globe,
        action: () => windowStore.openApp('websites')
      },
      { separator: true },
      {
        label: 'Info Website',
        icon: Info,
        action: () => {
          systemStore.showNotification('Website Info', `${app.title} (Port: ${app.port} | Domain: ${app.domain})`, 'info')
        }
      }
    ], app.title)
  } else if (app.isStoreApp) {
    // Context Menu untuk Aplikasi Terpasang dari Store / Docker Compose
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.launchStoreApp(app)
      },
      {
        label: 'Buka di Tab Baru Peramban',
        icon: ExternalLink,
        action: () => windowStore.launchStoreApp(app)
      },
      { separator: true },
      {
        label: 'Info Container',
        icon: Info,
        action: () => {
          systemStore.showNotification('Info Container', `${app.title} (Port Host: ${app.port} | Slug: ${app.slug})`, 'info')
        }
      },
      { separator: true },
      {
        label: 'Uninstall Aplikasi',
        icon: Trash2,
        danger: true,
        action: () => openDesktopUninstall(app)
      }
    ], app.title)
  } else {
    // Context Menu untuk Aplikasi Sistem Core
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.openApp(app.id)
      },
      { separator: true },
      {
        label: 'Info Aplikasi',
        icon: Info,
        action: () => {
          systemStore.showNotification('Aplikasi', `${app.title} (Resolusi Default: ${app.width}x${app.height})`, 'info')
        }
      }
    ], app.title)
  }
}

// ==========================================
// DESKTOP WIDGETS & SEARCH LOGIC
// ==========================================
const currentTime = ref('')
const currentDate = ref('')
let clockTimer = null

function updateClock() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  currentDate.value = now.toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})

// Handler Pencarian Desktop (CasaOS Floating Search Bar)
const searchQuery = ref('')

function handleDesktopSearch() {
  if (!searchQuery.value.trim()) return
  const q = encodeURIComponent(searchQuery.value.trim())
  let url = `https://duckduckgo.com/?q=${q}`
  const engine = systemStore.settings?.searchEngine || 'DuckDuckGo'
  if (engine === 'Google') url = `https://www.google.com/search?q=${q}`
  else if (engine === 'Bing') url = `https://www.bing.com/search?q=${q}`
  else if (engine === 'Baidu') url = `https://www.baidu.com/s?wd=${q}`
  window.open(url, '_blank')
  searchQuery.value = ''
}

// Helper Format Ukuran Byte
function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// ==========================================
// WINDOWS DESKTOP GRID & DRAG-AND-DROP SYSTEM
// ==========================================
const DESKTOP_CONFIG_STORAGE_KEY = 'razos_desktop_config'
const DESKTOP_POS_STORAGE_KEY = 'razos_desktop_icon_positions'

// Pilihan Ukuran Ikon Desktop: 'small', 'medium', 'large'
const desktopIconSize = ref('medium')

// Kalkulasi Dimensi Grid Slot & Style berdasarkan Ukuran Ikon
const iconGridMetrics = computed(() => {
  switch (desktopIconSize.value) {
    case 'small':
      return {
        cellW: 80,
        cellH: 84,
        iconBox: 'w-9 h-9',
        fontSize: 'text-[10px]'
      }
    case 'large':
      return {
        cellW: 108,
        cellH: 114,
        iconBox: 'w-16 h-16',
        fontSize: 'text-xs'
      }
    case 'medium':
    default:
      return {
        cellW: 92,
        cellH: 96,
        iconBox: 'w-12 h-12',
        fontSize: 'text-[11px]'
      }
  }
})

// Mapping posisi tiap-tiap icon: { [appId]: { col: number, row: number } }
const iconPositions = ref({})

// Simpan konfigurasi (ukuran & posisi) ke server backend & localStorage
async function persistDesktopConfig(showToast = false) {
  const payload = {
    icon_size: desktopIconSize.value,
    positions: iconPositions.value
  }

  // Simpan ke localStorage untuk loading instan
  try {
    localStorage.setItem(DESKTOP_CONFIG_STORAGE_KEY, JSON.stringify(payload))
    localStorage.setItem(DESKTOP_POS_STORAGE_KEY, JSON.stringify(iconPositions.value))
  } catch (e) {}

  // Push ke server backend agar persisten lintas perangkat (HP, laptop lain, dsb)
  try {
    await fetch('/api/v1/desktop/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (showToast) {
      systemStore.showNotification('Desktop', 'Tata letak desktop disinkronkan', 'success')
    }
  } catch (err) {
    console.warn('Gagal menyimpan konfigurasi desktop ke server:', err)
  }
}

// Muat konfigurasi desktop dari server dan localStorage
async function fetchDesktopConfig() {
  // 1. Baca localStorage dulu agar tampilan langsung muncul tanpa jeda
  try {
    const raw = localStorage.getItem(DESKTOP_CONFIG_STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (parsed.icon_size) desktopIconSize.value = parsed.icon_size
      if (parsed.positions) iconPositions.value = parsed.positions
    } else {
      const oldPos = localStorage.getItem(DESKTOP_POS_STORAGE_KEY)
      if (oldPos) iconPositions.value = JSON.parse(oldPos)
    }
  } catch (e) {}

  // 2. Fetch data mutakhir dari server backend
  try {
    const res = await fetch('/api/v1/desktop/config')
    if (res.ok) {
      const serverCfg = await res.json()
      if (serverCfg.icon_size) {
        desktopIconSize.value = serverCfg.icon_size
      }
      if (serverCfg.positions && Object.keys(serverCfg.positions).length > 0) {
        iconPositions.value = serverCfg.positions
      }
      // Update cache lokal
      try {
        localStorage.setItem(DESKTOP_CONFIG_STORAGE_KEY, JSON.stringify({
          icon_size: desktopIconSize.value,
          positions: iconPositions.value
        }))
        localStorage.setItem(DESKTOP_POS_STORAGE_KEY, JSON.stringify(iconPositions.value))
      } catch (e) {}
    }
  } catch (err) {
    console.warn('Gagal memuat konfigurasi desktop dari server:', err)
  }

  ensureIconPositions()
}

// Ubah ukuran ikon desktop (small, medium, large)
function changeIconSize(size) {
  if (['small', 'medium', 'large'].includes(size)) {
    desktopIconSize.value = size
    persistDesktopConfig()
    systemStore.showNotification('Tampilan Desktop', `Ukuran ikon diubah ke ${size.toUpperCase()}`, 'info')
  }
}

// Ratakan ikon ke grid secara otomatis (Auto-arrange ala Windows)
function autoArrangeIcons() {
  const allApps = windowStore.allDesktopApps || []
  const newPositions = {}
  let curCol = 0
  let curRow = 0
  const maxRowsPerCol = Math.max(4, Math.floor((window.innerHeight - 80) / iconGridMetrics.value.cellH))

  allApps.forEach(app => {
    newPositions[app.id] = { col: curCol, row: curRow }
    curRow++
    if (curRow >= maxRowsPerCol) {
      curRow = 0
      curCol++
    }
  })

  iconPositions.value = newPositions
  persistDesktopConfig(true)
}

// Menghitung koordinat awal grid desktop
function getGridOrigin() {
  return {
    x: 24,
    y: systemStore.settings?.showSearchBar ? 76 : 24
  }
}

// Inisialisasi slot grid dan assign posisi untuk aplikasi baru secara otomatis
function ensureIconPositions() {
  const occupied = new Set()
  const result = { ...iconPositions.value }

  // Catat sel yang sudah terisi
  Object.values(result).forEach(pos => {
    if (pos && typeof pos.col === 'number' && typeof pos.row === 'number') {
      occupied.add(`${pos.col},${pos.row}`)
    }
  })

  // Cari slot kosong berurutan kolom demi kolom (ala Windows desktop)
  let curCol = 0
  let curRow = 0
  const cellH = iconGridMetrics.value.cellH
  const maxRowsPerCol = Math.max(4, Math.floor((window.innerHeight - 80) / cellH))

  const allApps = windowStore.allDesktopApps || []
  let hasNew = false
  allApps.forEach(app => {
    if (!result[app.id]) {
      hasNew = true
      while (occupied.has(`${curCol},${curRow}`)) {
        curRow++
        if (curRow >= maxRowsPerCol) {
          curRow = 0
          curCol++
        }
      }
      result[app.id] = { col: curCol, row: curRow }
      occupied.add(`${curCol},${curRow}`)
      curRow++
      if (curRow >= maxRowsPerCol) {
        curRow = 0
        curCol++
      }
    }
  })

  iconPositions.value = result
  if (hasNew) {
    persistDesktopConfig(false)
  }
}

// Watch perubahan daftar aplikasi agar icon website / store app baru otomatis tertata
watch(() => windowStore.allDesktopApps.length, () => {
  ensureIconPositions()
}, { immediate: true })

onMounted(() => {
  fetchDesktopConfig()
})

// State Drag & Drop
const draggedAppId = ref(null)
const dropPreview = ref(null) // { col, row, x, y }

function handleIconDragStart(e, app) {
  draggedAppId.value = app.id
  e.dataTransfer.effectAllowed = 'move'
  e.dataTransfer.setData('text/plain', app.id)
}

function handleIconDragEnd() {
  draggedAppId.value = null
  dropPreview.value = null
}

function calculateGridCell(clientX, clientY) {
  const origin = getGridOrigin()
  const cellW = iconGridMetrics.value.cellW
  const cellH = iconGridMetrics.value.cellH

  // Batas maksimum kolom agar tidak menimpa panel widget di sebelah kanan
  const maxCol = Math.max(0, Math.floor((window.innerWidth - 320 - origin.x) / cellW))
  const maxRow = Math.max(0, Math.floor((window.innerHeight - 56 - origin.y) / cellH))

  let col = Math.floor((clientX - origin.x) / cellW)
  let row = Math.floor((clientY - origin.y) / cellH)

  col = Math.max(0, Math.min(col, maxCol))
  row = Math.max(0, Math.min(row, maxRow))

  return { col, row }
}

function handleDesktopDragOver(e) {
  if (!draggedAppId.value) return
  const { col, row } = calculateGridCell(e.clientX, e.clientY)
  const origin = getGridOrigin()
  const cellW = iconGridMetrics.value.cellW
  const cellH = iconGridMetrics.value.cellH

  dropPreview.value = {
    col,
    row,
    x: origin.x + col * cellW,
    y: origin.y + row * cellH
  }
}

function handleDesktopDragLeave(e) {
  if (!e.relatedTarget) {
    dropPreview.value = null
  }
}

function handleDesktopDrop(e) {
  if (!draggedAppId.value) return
  const { col, row } = calculateGridCell(e.clientX, e.clientY)
  const currentAppId = draggedAppId.value

  // Cek apakah slot target sudah ditempati icon lain (tukar posisi ala Windows)
  const existingAppEntry = Object.entries(iconPositions.value).find(
    ([id, pos]) => id !== currentAppId && pos.col === col && pos.row === row
  )

  if (existingAppEntry) {
    const [existingId, _] = existingAppEntry
    const oldPos = iconPositions.value[currentAppId] || { col: 0, row: 0 }
    iconPositions.value[existingId] = { ...oldPos }
  }

  // Terapkan posisi baru pada app yang di-drag
  iconPositions.value[currentAppId] = { col, row }
  persistDesktopConfig(true)

  draggedAppId.value = null
  dropPreview.value = null
}
</script>

<template>
  <div
    class="relative w-screen h-screen overflow-hidden bg-cover bg-center select-none"
    :style="{ backgroundImage: `url('${systemStore.settings?.wallpaper || '/assets/images/RAZDesktopBg.webp'}')` }"
    @contextmenu.prevent="handleDesktopContextMenu"
    @dragover.prevent="handleDesktopDragOver"
    @dragleave="handleDesktopDragLeave"
    @drop.prevent="handleDesktopDrop"
  >
    <!-- Subtle Brand Watermark di pojok kanan atas -->
    <div class="absolute top-4 right-6 pointer-events-none opacity-40 z-0">
      <img src="/assets/brand/RAZBrandTransparent.webp" alt="RAZ STUDIO" class="h-8 object-contain" />
    </div>

    <!-- Snap-to-Grid Target Preview (Windows Style) -->
    <div
      v-if="dropPreview"
      :style="{
        position: 'absolute',
        left: `${dropPreview.x}px`,
        top: `${dropPreview.y}px`,
        width: `${iconGridMetrics.cellW}px`,
        height: `${iconGridMetrics.cellH}px`
      }"
      class="rounded-2xl border-2 border-dashed border-sky-400 bg-sky-500/20 pointer-events-none transition-all duration-100 z-10 flex items-center justify-center shadow-lg shadow-sky-500/20 ring-4 ring-sky-400/20"
    >
      <div class="w-8 h-8 rounded-xl bg-sky-400/30 flex items-center justify-center">
        <div class="w-2.5 h-2.5 rounded-full bg-sky-300 animate-ping"></div>
      </div>
    </div>

    <!-- Desktop Layer (z-10, h-[calc(100vh-48px)]) -->
    <div class="relative z-10 w-full h-[calc(100vh-48px)] flex flex-col p-6 pointer-events-none">
      <!-- Top Row: Floating Search Bar (CasaOS Style - Centered, 100% Clickable) -->
      <div
        v-if="systemStore.settings?.showSearchBar"
        class="w-full flex justify-center pb-3 pointer-events-auto"
      >
        <form
          @submit.prevent="handleDesktopSearch"
          class="flex items-center gap-2.5 px-4 py-2 rounded-2xl bg-slate-900/85 backdrop-blur-md border border-slate-700/80 shadow-2xl text-slate-100 hover:border-sky-500/50 transition duration-200 w-full max-w-md"
        >
          <Search class="w-4 h-4 text-sky-400 shrink-0" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="`Cari web dengan ${systemStore.settings?.searchEngine || 'DuckDuckGo'}...`"
            class="bg-transparent border-none outline-none text-xs sm:text-sm text-slate-100 placeholder-slate-400 w-full font-medium select-text cursor-text"
          />
          <button
            type="submit"
            class="text-[10px] font-semibold font-mono px-2 py-0.5 rounded-md bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700/60 uppercase transition active:scale-95 shrink-0"
          >
            {{ systemStore.settings?.searchEngine || 'DuckDuckGo' }}
          </button>
        </form>
      </div>

      <!-- Main Desktop Workspace: Icons Absolute Grid + Right Widgets Column -->
      <div class="flex-1 relative w-full h-full min-h-0 pointer-events-none overflow-hidden">
        <!-- Windows Desktop Grid Icons (Draggable & Absolute Grid Positioned) -->
        <div
          v-for="app in windowStore.allDesktopApps"
          :key="app.id"
          draggable="true"
          @dragstart="handleIconDragStart($event, app)"
          @dragend="handleIconDragEnd"
          @dblclick="app.isWebsite ? windowStore.launchWebsiteApp(app) : (app.isStoreApp ? windowStore.launchStoreApp(app) : windowStore.openApp(app.id))"
          @contextmenu.prevent.stop="handleDesktopIconContextMenu($event, app)"
          :style="{
            position: 'absolute',
            left: `${(iconPositions[app.id]?.col ?? 0) * iconGridMetrics.cellW}px`,
            top: `${(iconPositions[app.id]?.row ?? 0) * iconGridMetrics.cellH}px`,
            width: `${iconGridMetrics.cellW}px`,
            height: `${iconGridMetrics.cellH}px`
          }"
          class="flex flex-col items-center justify-start p-1.5 rounded-2xl hover:bg-slate-800/60 cursor-pointer group transition-all duration-150 text-center select-none pointer-events-auto"
          :class="draggedAppId === app.id ? 'opacity-30 scale-95 ring-2 ring-sky-400' : ''"
        >
          <!-- Ikon Aplikasi -->
          <div
            :class="[iconGridMetrics.iconBox]"
            class="rounded-2xl flex items-center justify-center p-1 group-hover:scale-105 transition-transform duration-150 drop-shadow-md relative shrink-0 pointer-events-none"
          >
            <img :src="app.icon" :alt="app.title" class="w-full h-full object-contain pointer-events-none" />
            <span v-if="app.isStoreApp || app.isWebsite" class="absolute -top-1 -right-1 w-3 h-3 bg-emerald-500 rounded-full border-2 border-slate-900 shadow-sm" :title="app.isWebsite ? 'Website Running' : 'Container Running'"></span>
          </div>
          <!-- Judul / Nama Aplikasi Desktop (Kontras Tinggi & Selalu Terbaca) -->
          <span
            :class="[iconGridMetrics.fontSize]"
            class="mt-1 leading-tight text-slate-100 font-medium group-hover:text-sky-300 drop-shadow-[0_2px_4px_rgba(0,0,0,0.95)] line-clamp-2 break-words w-full px-1 select-none pointer-events-none"
          >
            {{ app.title }}
          </span>
        </div>

        <!-- RIGHT: Desktop Widgets Card (Tampil di sebelah kanan jika minimal 1 widget aktif) -->
        <div
          v-if="systemStore.settings?.widgets && (systemStore.settings.widgets.time || systemStore.settings.widgets.systemStatus || systemStore.settings.widgets.storageStatus || systemStore.settings.widgets.networkStatus)"
          class="absolute top-0 right-0 w-72 rounded-2xl bg-slate-900/80 backdrop-blur-md border border-slate-800/90 shadow-2xl p-4 space-y-4 text-slate-200 pointer-events-auto animate-in fade-in select-none"
        >
          <!-- Widget 1: Time -->
          <div v-if="systemStore.settings.widgets.time" class="pb-3 border-b border-slate-800/80">
            <div class="text-3xl font-extrabold tracking-tight text-slate-100 font-mono">
              {{ currentTime }}
            </div>
            <div class="text-xs text-sky-400 font-medium mt-0.5 capitalize">
              {{ currentDate }}
            </div>
          </div>

          <!-- Widget 2: System Status (CPU, RAM, Temp) -->
          <div v-if="systemStore.settings.widgets.systemStatus" class="space-y-2.5 pb-3 border-b border-slate-800/80">
            <!-- CPU Load & Temperature -->
            <div class="flex items-center justify-between text-xs">
              <span class="flex items-center gap-1.5 font-semibold text-slate-300">
                <Activity class="w-3.5 h-3.5 text-emerald-400" />
                <span>Beban CPU</span>
              </span>
              <span class="font-mono text-emerald-400 font-bold">
                {{ systemStore.metrics.cpu_percent ? systemStore.metrics.cpu_percent.toFixed(1) : 0 }}%
                <span v-if="systemStore.metrics.cpu_temp > 0" class="text-slate-400 font-normal">({{ systemStore.metrics.cpu_temp }}°C)</span>
              </span>
            </div>
            <div class="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
              <div
                class="h-full bg-emerald-500 rounded-full transition-all duration-300"
                :style="{ width: `${Math.min(systemStore.metrics.cpu_percent || 0, 100)}%` }"
              ></div>
            </div>

            <!-- RAM Load -->
            <div class="flex items-center justify-between text-xs pt-1">
              <span class="flex items-center gap-1.5 font-semibold text-slate-300">
                <Cpu class="w-3.5 h-3.5 text-sky-400" />
                <span>Memori RAM</span>
              </span>
              <span class="font-mono text-sky-400 font-bold">
                {{ systemStore.metrics.ram_percent ? systemStore.metrics.ram_percent.toFixed(1) : 0 }}%
              </span>
            </div>
            <div class="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
              <div
                class="h-full bg-sky-500 rounded-full transition-all duration-300"
                :style="{ width: `${Math.min(systemStore.metrics.ram_percent || 0, 100)}%` }"
              ></div>
            </div>
          </div>

          <!-- Widget 3: Storage Status (/DATA) -->
          <div v-if="systemStore.settings.widgets.storageStatus" class="space-y-2 pb-3 border-b border-slate-800/80">
            <div class="flex items-center justify-between text-xs">
              <span class="flex items-center gap-1.5 font-semibold text-slate-300">
                <HardDrive class="w-3.5 h-3.5 text-amber-400" />
                <span>Storage (/DATA)</span>
              </span>
              <span class="font-mono text-amber-400 font-bold text-[11px]">
                {{ formatBytes(systemStore.metrics.disk_used) }} / {{ formatBytes(systemStore.metrics.disk_total) }}
              </span>
            </div>
            <div class="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
              <div
                class="h-full bg-amber-500 rounded-full transition-all duration-300"
                :style="{ width: `${Math.min(systemStore.metrics.disk_percent || 0, 100)}%` }"
              ></div>
            </div>
          </div>

          <!-- Widget 4: Network Status -->
          <div v-if="systemStore.settings.widgets.networkStatus" class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <span class="flex items-center gap-1.5 font-semibold text-slate-300">
                <Network class="w-3.5 h-3.5 text-indigo-400" />
                <span>Kecepatan Jaringan</span>
              </span>
            </div>
            <div class="grid grid-cols-2 gap-2 text-xs font-mono pt-1">
              <div class="p-2 rounded-xl bg-slate-950/60 border border-slate-800">
                <span class="text-[10px] text-slate-400 block">Download</span>
                <span class="text-sky-300 font-bold text-xs truncate block">↓ {{ systemStore.metrics.net_download_speed || '0 KB/s' }}</span>
              </div>
              <div class="p-2 rounded-xl bg-slate-950/60 border border-slate-800">
                <span class="text-[10px] text-slate-400 block">Upload</span>
                <span class="text-indigo-300 font-bold text-xs truncate block">↑ {{ systemStore.metrics.net_upload_speed || '0 KB/s' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Window Management Layer -->
    <div class="absolute inset-0 pointer-events-none z-20">
      <div class="relative w-full h-[calc(100vh-48px)]">
        <template v-for="win in windowStore.openWindows" :key="win.id">
          <RAZWindowFrame :win="win" class="pointer-events-auto">
            <!-- Standalone Universal File Viewer -->
            <RAZFileViewerView v-if="win.id.startsWith('viewer-')" :file="win.fileData" />

            <!-- Standalone File Properties Window -->
            <RAZFilePropertiesView v-else-if="win.id.startsWith('file-prop-')" :win="win" />

            <!-- Standalone Website Wizard Window -->
            <RAZWebsiteWizardView v-else-if="win.id === 'website-create'" :win="win" />

            <!-- Standalone Website Edit Window -->
            <RAZWebsiteEditView v-else-if="win.id.startsWith('website-edit-')" :win="win" />

            <!-- Standalone Website Git Sync Window -->
            <RAZWebsiteGitView v-else-if="win.id.startsWith('website-git-')" :win="win" />

            <!-- Standalone Website Detail Window -->
            <RAZWebsiteDetailView v-else-if="win.id.startsWith('website-detail-')" :win="win" />

            <!-- Standalone Database Creator Window -->
            <RAZDatabaseCreatorView v-else-if="win.id === 'database-create'" :win="win" />

            <!-- Standalone App Store Manual Quick Install Window -->
            <RAZAppStoreManualView v-else-if="win.id === 'appstore-manual'" :win="win" />

            <!-- Standalone App Store Catalog Installer Window -->
            <RAZAppStoreInstallView v-else-if="win.id.startsWith('appstore-install-')" :win="win" />

            <!-- Aplikasi Utama Berdasarkan ID -->
            <component v-else :is="viewMap[win.id]" />
          </RAZWindowFrame>
        </template>
      </div>
    </div>

    <!-- Bottom Desktop Taskbar -->
    <RAZTaskbar />

    <!-- Custom Modal Konfirmasi Uninstall Aplikasi (Floating Window Dialog) -->
    <div
      v-if="showUninstallModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/35 p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md rounded-2xl bg-slate-900 border border-rose-500/40 shadow-2xl shadow-black ring-1 ring-rose-500/20 flex flex-col overflow-hidden text-slate-100">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5 text-rose-400">
            <div class="w-7 h-7 rounded-lg bg-rose-500/20 text-rose-400 flex items-center justify-center border border-rose-500/30">
              <AlertTriangle class="w-4 h-4" />
            </div>
            <h3 class="font-bold text-sm">Konfirmasi Uninstall Aplikasi</h3>
          </div>
          <button @click="showUninstallModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 flex flex-col gap-3.5">
          <div class="flex items-center gap-3 p-3 rounded-xl bg-slate-950/60 border border-slate-800">
            <img :src="appToUninstall?.icon" :alt="appToUninstall?.title" class="w-10 h-10 object-contain rounded-lg" />
            <div class="min-w-0 flex-1">
              <h4 class="font-bold text-sm text-slate-100 truncate">{{ appToUninstall?.title }}</h4>
              <span class="text-xs text-slate-400 font-mono">Port: :{{ appToUninstall?.port }}</span>
            </div>
          </div>

          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda yakin ingin menghentikan dan menghapus container aplikasi ini dari sistem?
          </p>

          <label class="flex items-center gap-2.5 p-3 rounded-xl bg-slate-950 border border-slate-800 cursor-pointer text-xs text-rose-300">
            <input type="checkbox" v-model="purgeDataChecked" class="accent-rose-500 rounded" />
            <span>Hapus juga seluruh direktori data permanen container di storage</span>
          </label>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showUninstallModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="confirmDesktopUninstall"
            class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-lg active:scale-95"
          >
            <Trash2 class="w-4 h-4" />
            <span>Uninstall Container</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
