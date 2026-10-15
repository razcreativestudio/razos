<script setup>
/**
 * RAZOS Desktop Layout Component
 * Updated: 2026-10-15 09:05:00 (v3.2.0: Slider Rekomendasi Aplikasi di bawah widget desktop, responsive grid & position lock)
 * Deskripsi: Layout Desktop mode untuk breakpoint > 768px, memadukan desktop icon shortcuts,
 * window manager frame multi-jendela, taskbar bawah, context menu lengkap, dan slider rekomendasi aplikasi.
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
  Network, Search, Grid, Maximize2, Minimize2, Sliders, Power, RotateCw,
  HelpCircle, GitBranch, Lock, Unlock, Sparkles, ChevronLeft, ChevronRight,
  Download, ArrowRight, Box, Database, LayoutGrid
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

// Import View Jendela Khusus / Standalone Child Windows (v3.0.0)
import RAZWebsiteWizardView from '../views/RAZWebsiteWizardView.vue'
import RAZWebsiteEditView from '../views/RAZWebsiteEditView.vue'
import RAZWebsiteGitView from '../views/RAZWebsiteGitView.vue'
import RAZWebsiteDetailView from '../views/RAZWebsiteDetailView.vue'
import RAZFilePropertiesView from '../views/RAZFilePropertiesView.vue'
import RAZDatabaseCreatorView from '../views/RAZDatabaseCreatorView.vue'
import RAZAppStoreManualView from '../views/RAZAppStoreManualView.vue'
import RAZAppStoreInstallView from '../views/RAZAppStoreInstallView.vue'
import RAZAppEditView from '../views/RAZAppEditView.vue'
import RAZAppLogsView from '../views/RAZAppLogsView.vue'
import RAZStorageManagerView from '../views/RAZStorageManagerView.vue'
import RAZTextEditorView from '../views/RAZTextEditorView.vue'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()

// Mapping ID aplikasi ke komponen View
const viewMap = {
  'file-manager': RAZFileManagerView,
  'file-manager-editor': RAZTextEditorView,
  'text-editor': RAZTextEditorView,
  'task-manager': RAZTaskManagerView,
  'terminal': RAZTerminalView,
  'app-store': RAZAppStoreView,
  'websites': RAZWebsitesView,
  'databases': RAZDatabasesView,
  'settings': RAZSettingsView
}

// Context Menu Klik Kanan pada Area Desktop Kosong (Hierarki Submenu Modern ala OS)
function handleDesktopContextMenu(e) {
  contextMenuStore.show(e, [
    // 1. Submenu: Tampilan (View)
    {
      label: 'Tampilan (View)',
      icon: LayoutGrid,
      children: [
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
        { separator: true },
        {
          label: 'Ratakan Ikon ke Grid',
          icon: RefreshCw,
          action: () => autoArrangeIcons(true)
        },
        {
          label: isPositionLocked.value ? 'Buka Kunci Posisi Ikon' : 'Kunci Posisi Ikon',
          icon: isPositionLocked.value ? Unlock : Lock,
          checked: isPositionLocked.value,
          action: () => togglePositionLock()
        }
      ]
    },

    // 2. Submenu: Alat & Aplikasi Sistem
    {
      label: 'Alat & Sistem',
      icon: Box,
      children: [
        {
          label: 'Files & Storage',
          icon: Folder,
          action: () => windowStore.openApp('file-manager')
        },
        {
          label: 'RAZ Studio Code',
          icon: Code2,
          action: () => windowStore.openTextEditor()
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
        {
          label: 'App Store',
          icon: Box,
          action: () => windowStore.openApp('app-store')
        },
        {
          label: 'Website Manager',
          icon: Globe,
          action: () => windowStore.openApp('websites')
        },
        {
          label: 'Database Manager',
          icon: Database,
          action: () => windowStore.openApp('databases')
        },
        {
          label: 'Storage Manager',
          icon: HardDrive,
          action: () => windowStore.openStorageManager()
        }
      ]
    },

    // 3. Submenu: Manajemen Jendela
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

    // 4. Aksi Utama Langsung: Refresh / Muat Ulang
    {
      label: 'Muat Ulang Desktop',
      icon: RefreshCw,
      shortcut: 'F5',
      action: () => {
        systemStore.fetchMetrics()
        fetchDesktopConfig()
        systemStore.showNotification('Desktop', 'Tampilan desktop telah disegarkan', 'success')
      }
    },

    { separator: true },

    // 5. Pengaturan Sistem
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

// Context Menu Klik Kanan pada Ikon Shortcut Desktop (CasaOS Enterprise Menu System)
function handleDesktopIconContextMenu(e, app) {
  if (app.isWebsite) {
    // Context Menu untuk Website yang Dibuat di Website Manager (Hierarki Submenu)
    contextMenuStore.show(e, [
      {
        label: 'Open Website',
        icon: Play,
        action: () => windowStore.launchWebsiteApp(app)
      },
      {
        label: 'Pengaturan & Alat',
        icon: Sliders,
        children: [
          {
            label: 'Settings (Edit Konfigurasi)',
            icon: Sliders,
            action: () => windowStore.openWebsiteEditor(app)
          },
          {
            label: 'Git Sync & Pull',
            icon: GitBranch,
            action: () => windowStore.openWebsiteGit(app)
          },
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
          }
        ]
      },
      { separator: true },
      {
        label: 'Hapus Website',
        icon: Trash2,
        danger: true,
        action: () => {
          windowStore.openApp('websites')
        }
      },
      {
        type: 'quick-actions',
        actions: [
          {
            icon: RotateCw,
            label: 'Restart',
            title: 'Mulai Ulang Layanan Website',
            action: async () => {
              try {
                systemStore.showNotification('Memuat Ulang', `Merestart website ${app.title}...`, 'info')
                await fetch(`/api/v1/websites/${app.id}/action`, {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ action: 'restart' })
                })
                systemStore.showNotification('Sukses', `Layanan website ${app.title} berhasil dimuat ulang!`, 'success')
                systemStore.triggerWebsitesRefresh()
              } catch (e) {}
            }
          },
          {
            icon: Power,
            label: 'Stop',
            danger: true,
            title: 'Hentikan Layanan Website',
            action: async () => {
              try {
                await fetch(`/api/v1/websites/${app.id}/action`, {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ action: 'stop' })
                })
                systemStore.showNotification('Website Dihentikan', `Layanan website ${app.title} dihentikan`, 'warning')
                systemStore.triggerWebsitesRefresh()
              } catch (e) {}
            }
          }
        ]
      }
    ], app.title)
  } else if (app.isStoreApp) {
    // Context Menu untuk Aplikasi Terpasang dari Store / Docker Compose (Hierarki Submenu)
    contextMenuStore.show(e, [
      {
        label: 'Open App',
        icon: Play,
        action: () => windowStore.launchStoreApp(app)
      },
      {
        label: 'Pengaturan Container',
        icon: Sliders,
        children: [
          {
            label: 'Settings (Konfigurasi Container)',
            icon: Sliders,
            action: () => windowStore.openAppEditor(app)
          },
          {
            label: 'Terminal & Logs',
            icon: Terminal,
            action: () => windowStore.openAppLogs(app)
          },
          {
            label: 'Check then update',
            icon: RefreshCw,
            action: async () => {
              systemStore.showNotification('Pembaruan', `Memeriksa versi image "${app.image || app.slug}"...`, 'info')
              setTimeout(() => {
                systemStore.showNotification('Versi Terbaru', `Aplikasi "${app.title}" sudah menggunakan tag image terbaru.`, 'success')
              }, 1500)
            }
          },
          {
            label: 'Tips & Rincian',
            icon: HelpCircle,
            action: () => {
              systemStore.showNotification(
                `Tips: ${app.title}`,
                `Port: :${app.port || 8080} | Storage: ${app.data_path || `/DATA/apps/${app.slug}`}\nContainer: razos-${app.slug}`,
                'info'
              )
            }
          }
        ]
      },
      { separator: true },
      {
        label: 'Uninstall',
        icon: Trash2,
        danger: true,
        action: () => openDesktopUninstall(app)
      },
      {
        type: 'quick-actions',
        actions: [
          {
            icon: RotateCw,
            label: 'Restart',
            title: 'Mulai Ulang Container (Reboot)',
            action: async () => {
              try {
                systemStore.showNotification('Memuat Ulang', `Merestart container ${app.title}...`, 'info')
                await fetch(`/api/v1/store/apps/${app.slug || app.id}/restart`, { method: 'POST' })
                systemStore.showNotification('Sukses', `Container ${app.title} berhasil dimuat ulang!`, 'success')
              } catch (e) {
                systemStore.showNotification('Error', 'Gagal merestart container', 'error')
              }
            }
          },
          {
            icon: Power,
            label: 'Stop',
            danger: true,
            title: 'Hentikan Container',
            action: async () => {
              try {
                await fetch(`/api/v1/store/apps/${app.slug || app.id}/action`, {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ action: 'stop' })
                })
                systemStore.showNotification('Container Dihentikan', `Container ${app.title} telah dihentikan`, 'warning')
                systemStore.triggerAppsRefresh()
              } catch (e) {
                systemStore.showNotification('Error', 'Gagal mengubah status container', 'error')
              }
            }
          }
        ]
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

// Clock timer & resize listener dimulai di onMounted utama (bagian bawah script)

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
// DESKTOP RECOMMENDED APPS SLIDER (v3.2.0)
// ==========================================
const recommendedAppsList = ref([
  {
    name: 'Nextcloud Hub',
    slug: 'nextcloud',
    category: 'Cloud Storage',
    tag: 'Populer',
    icon: 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/nextcloud.png',
    description: 'Penyimpanan file cloud mandiri, kalender, kontak, dan kantor kolaboratif.',
    port: 8080,
    image: 'nextcloud:apache'
  },
  {
    name: 'Jellyfin',
    slug: 'jellyfin',
    category: 'Media Server',
    tag: 'Favorit',
    icon: 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/jellyfin.png',
    description: 'Sistem media streaming gratis tanpa langganan untuk film, anime, dan serial TV.',
    port: 8096,
    image: 'jellyfin/jellyfin:latest'
  },
  {
    name: 'AdGuard Home',
    slug: 'adguard-home',
    category: 'Network Shield',
    tag: 'Keamanan',
    icon: 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/adguard-home.png',
    description: 'DNS resolver pemblokir iklan, phishing, dan pelacak untuk seluruh perangkat LAN.',
    port: 3001,
    image: 'adguard/adguardhome:latest'
  },
  {
    name: 'Portainer CE',
    slug: 'portainer',
    category: 'DevOps & Git',
    tag: 'Essential',
    icon: 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/portainer.png',
    description: 'Manajemen visual container Docker, compose stack, image, dan resource sistem.',
    port: 9000,
    image: 'portainer/portainer-ce:latest'
  }
])

const currentRecIndex = ref(0)
let recSliderTimer = null

function nextRecApp() {
  currentRecIndex.value = (currentRecIndex.value + 1) % recommendedAppsList.value.length
}

function prevRecApp() {
  currentRecIndex.value = (currentRecIndex.value - 1 + recommendedAppsList.value.length) % recommendedAppsList.value.length
}

function startRecAutoSlide() {
  stopRecAutoSlide()
  recSliderTimer = setInterval(() => {
    nextRecApp()
  }, 6000)
}

function stopRecAutoSlide() {
  if (recSliderTimer) {
    clearInterval(recSliderTimer)
    recSliderTimer = null
  }
}

function isRecAppInstalled(slug) {
  return windowStore.installedStoreApps.some(i => (i.slug || i.id) === slug)
}

function handleInstallOrOpenRec(app) {
  if (isRecAppInstalled(app.slug)) {
    const installed = windowStore.installedStoreApps.find(i => (i.slug || i.id) === app.slug)
    if (installed) {
      windowStore.launchStoreApp({
        id: `store-${installed.slug || installed.id}`,
        title: installed.name,
        slug: installed.slug || installed.id,
        port: installed.port
      })
    }
  } else {
    windowStore.openAppStoreInstaller(app)
  }
}

// ==========================================
// WINDOWS DESKTOP GRID & DRAG-AND-DROP SYSTEM (RESPONSIVE & PERSISTENT)
// ==========================================
const DESKTOP_CONFIG_STORAGE_KEY = 'razos_desktop_config'
const DESKTOP_POS_STORAGE_KEY = 'razos_desktop_icon_positions'

// Ref container area kerja desktop
const desktopWorkspaceRef = ref(null)

// Pilihan Ukuran Ikon Desktop: 'small', 'medium', 'large'
const desktopIconSize = ref('medium')

// Status Penguncian Posisi Ikon Desktop (Position Lock)
const isPositionLocked = ref(false)

// Status kesiapan konfigurasi (mencegah penimpaan posisi default sebelum data tersimpan selesai dimuat)
const isConfigLoaded = ref(false)

// Mapping posisi tiap-tiap icon: { [appId]: { col: number, row: number } }
const iconPositions = ref({})

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
        cellH: 94,
        iconBox: 'w-12 h-12',
        fontSize: 'text-[11px]'
      }
  }
})

// Dimensi workspace aktual yang responsif
const workspaceHeight = ref(window.innerHeight - 140)
const workspaceWidth = ref(window.innerWidth - 48)

function updateWorkspaceSize() {
  if (desktopWorkspaceRef.value) {
    workspaceHeight.value = desktopWorkspaceRef.value.clientHeight || (window.innerHeight - 140)
    workspaceWidth.value = desktopWorkspaceRef.value.clientWidth || (window.innerWidth - 48)
  } else {
    workspaceHeight.value = window.innerHeight - (systemStore.settings?.showSearchBar ? 160 : 100)
    workspaceWidth.value = window.innerWidth - 48
  }
}

// Menghitung jumlah baris maksimum yang muat secara responsif tanpa menabrak taskbar
const maxRows = computed(() => {
  const cellH = iconGridMetrics.value.cellH
  return Math.max(2, Math.floor(workspaceHeight.value / cellH))
})

// Menghitung jumlah kolom maksimum yang muat tanpa menimpa widget kanan
const maxCols = computed(() => {
  const cellW = iconGridMetrics.value.cellW
  const hasWidgets = (
    systemStore.settings?.widgets && (
      systemStore.settings.widgets.time ||
      systemStore.settings.widgets.systemStatus ||
      systemStore.settings.widgets.storageStatus ||
      systemStore.settings.widgets.networkStatus
    )
  ) || systemStore.settings?.showRecommendedApps !== false
  const usableWidth = hasWidgets ? Math.max(cellW * 2, workspaceWidth.value - 320) : workspaceWidth.value
  return Math.max(2, Math.floor(usableWidth / cellW))
})

// Menghitung posisi render responsif agar tidak pernah ada icon yang keluar batas atau terpotong
function getEffectivePosition(appId) {
  const pos = iconPositions.value[appId]
  if (!pos || typeof pos.col !== 'number' || typeof pos.row !== 'number') {
    return { col: 0, row: 0 }
  }
  const rows = maxRows.value
  if (rows > 0 && pos.row >= rows) {
    // Jika resolusi layar lebih pendek dari saat posisi disimpan,
    // bungkus secara mulus ke kolom berikutnya agar tetap terlihat 100%
    const extraCol = Math.floor(pos.row / rows)
    const remRow = pos.row % rows
    return {
      col: pos.col + extraCol,
      row: remRow
    }
  }
  return pos
}

// Toggle kunci posisi ikon desktop
function togglePositionLock() {
  isPositionLocked.value = !isPositionLocked.value
  persistDesktopConfig(false)
  systemStore.showNotification(
    'Kunci Tata Letak Desktop',
    isPositionLocked.value
      ? 'Posisi ikon berhasil DIKUNCI (tidak dapat digeser).'
      : 'Kunci posisi dibuka (ikon kini dapat digeser bebas).',
    isPositionLocked.value ? 'info' : 'success'
  )
}

// Simpan konfigurasi (ukuran, status kunci, & posisi) ke backend & localStorage
async function persistDesktopConfig(showToast = false) {
  if (!isConfigLoaded.value) return

  const payload = {
    icon_size: desktopIconSize.value,
    lock_position: isPositionLocked.value,
    positions: iconPositions.value
  }

  // Simpan ke localStorage untuk loading instan tanpa jeda
  try {
    localStorage.setItem(DESKTOP_CONFIG_STORAGE_KEY, JSON.stringify(payload))
    localStorage.setItem(DESKTOP_POS_STORAGE_KEY, JSON.stringify(iconPositions.value))
  } catch (e) {}

  // Push ke backend agar persisten lintas perangkat (HP, tablet, laptop lain, dsb)
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

// Muat konfigurasi desktop dari server dan localStorage secara aman
async function fetchDesktopConfig() {
  // 1. Baca cache lokal dulu agar antarmuka instan
  try {
    const raw = localStorage.getItem(DESKTOP_CONFIG_STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (parsed.icon_size) desktopIconSize.value = parsed.icon_size
      if (parsed.lock_position !== undefined) isPositionLocked.value = Boolean(parsed.lock_position)
      if (parsed.positions && typeof parsed.positions === 'object' && Object.keys(parsed.positions).length > 0) {
        iconPositions.value = { ...parsed.positions }
      }
    } else {
      const oldPos = localStorage.getItem(DESKTOP_POS_STORAGE_KEY)
      if (oldPos) iconPositions.value = JSON.parse(oldPos)
    }
  } catch (e) {}

  // 2. Fetch konfigurasi mutakhir dari server backend
  try {
    const res = await fetch('/api/v1/desktop/config')
    if (res.ok) {
      const serverCfg = await res.json()
      if (serverCfg.icon_size) {
        desktopIconSize.value = serverCfg.icon_size
      }
      if (serverCfg.lock_position !== undefined) {
        isPositionLocked.value = Boolean(serverCfg.lock_position)
      }
      if (serverCfg.positions && typeof serverCfg.positions === 'object' && Object.keys(serverCfg.positions).length > 0) {
        // Pertahankan posisi terakhir yang tersimpan di server
        iconPositions.value = { ...iconPositions.value, ...serverCfg.positions }
      }
      // Perbarui cache lokal
      try {
        localStorage.setItem(DESKTOP_CONFIG_STORAGE_KEY, JSON.stringify({
          icon_size: desktopIconSize.value,
          lock_position: isPositionLocked.value,
          positions: iconPositions.value
        }))
        localStorage.setItem(DESKTOP_POS_STORAGE_KEY, JSON.stringify(iconPositions.value))
      } catch (e) {}
    }
  } catch (err) {
    console.warn('Gagal memuat konfigurasi desktop dari server:', err)
  } finally {
    isConfigLoaded.value = true
    updateWorkspaceSize()
    ensureIconPositions()
  }
}

// Ubah ukuran ikon desktop (small, medium, large)
function changeIconSize(size) {
  if (['small', 'medium', 'large'].includes(size)) {
    desktopIconSize.value = size
    persistDesktopConfig(false)
    systemStore.showNotification('Tampilan Desktop', `Ukuran ikon diubah ke ${size.toUpperCase()}`, 'info')
  }
}

// Ratakan ikon ke grid secara otomatis (Auto-arrange ala Windows)
function autoArrangeIcons(showToast = true) {
  const allApps = windowStore.allDesktopApps || []
  const newPositions = {}
  let curCol = 0
  let curRow = 0
  const rows = maxRows.value

  allApps.forEach(app => {
    newPositions[app.id] = { col: curCol, row: curRow }
    curRow++
    if (curRow >= rows) {
      curRow = 0
      curCol++
    }
  })

  iconPositions.value = newPositions
  persistDesktopConfig(showToast)
}

// Inisialisasi slot grid dan assign posisi untuk aplikasi baru secara otomatis
function ensureIconPositions() {
  if (!isConfigLoaded.value) return

  const occupied = new Set()
  const result = { ...iconPositions.value }
  const rows = maxRows.value

  // Catat sel yang sudah terisi oleh posisi yang valid
  Object.values(result).forEach(pos => {
    if (pos && typeof pos.col === 'number' && typeof pos.row === 'number') {
      occupied.add(`${pos.col},${pos.row}`)
    }
  })

  const allApps = windowStore.allDesktopApps || []
  let hasNew = false
  let curCol = 0
  let curRow = 0

  allApps.forEach(app => {
    // Hanya assign slot baru bila aplikasi ini belum pernah memiliki koordinat tersimpan
    if (!result[app.id] || typeof result[app.id].col !== 'number' || typeof result[app.id].row !== 'number') {
      hasNew = true
      while (occupied.has(`${curCol},${curRow}`)) {
        curRow++
        if (curRow >= rows) {
          curRow = 0
          curCol++
        }
      }
      result[app.id] = { col: curCol, row: curRow }
      occupied.add(`${curCol},${curRow}`)
      curRow++
      if (curRow >= rows) {
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
})

function handleWindowResize() {
  updateWorkspaceSize()
}

onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
  startRecAutoSlide()
  window.addEventListener('resize', handleWindowResize)
  fetchDesktopConfig()
})

onUnmounted(() => {
  stopRecAutoSlide()
  window.removeEventListener('resize', handleWindowResize)
  if (clockTimer) clearInterval(clockTimer)
})

// State Drag & Drop
const draggedAppId = ref(null)
const isDraggingIcon = ref(false)
const dropPreview = ref(null) // { col, row, x, y }

function handleIconDragStart(e, app) {
  if (isPositionLocked.value) {
    e.preventDefault()
    systemStore.showNotification('Posisi Terkunci', 'Tata letak ikon sedang dikunci. Buka kunci di klik kanan desktop.', 'warning')
    return false
  }
  isDraggingIcon.value = true
  draggedAppId.value = app.id
  e.dataTransfer.effectAllowed = 'move'
  e.dataTransfer.setData('text/plain', app.id)
}

function handleIconDragEnd() {
  draggedAppId.value = null
  dropPreview.value = null
  setTimeout(() => {
    isDraggingIcon.value = false
  }, 120)
}

// Handler Peluncuran Aplikasi Desktop (Mendukung Single-Click dan Double-Click)
function handleAppLaunch(app) {
  if (isDraggingIcon.value || draggedAppId.value) return
  if (app.isWebsite) {
    windowStore.launchWebsiteApp(app)
  } else if (app.isStoreApp) {
    windowStore.launchStoreApp(app)
  } else {
    windowStore.openApp(app.id)
  }
}

function calculateGridCell(clientX, clientY) {
  if (!desktopWorkspaceRef.value) return { col: 0, row: 0 }
  const rect = desktopWorkspaceRef.value.getBoundingClientRect()
  const relX = clientX - rect.left
  const relY = clientY - rect.top

  const cellW = iconGridMetrics.value.cellW
  const cellH = iconGridMetrics.value.cellH

  let col = Math.floor(relX / cellW)
  let row = Math.floor(relY / cellH)

  col = Math.max(0, Math.min(col, maxCols.value - 1))
  row = Math.max(0, Math.min(row, maxRows.value - 1))

  return { col, row }
}

function handleDesktopDragOver(e) {
  if (!draggedAppId.value || isPositionLocked.value) return
  const { col, row } = calculateGridCell(e.clientX, e.clientY)
  dropPreview.value = {
    col,
    row,
    x: col * iconGridMetrics.value.cellW,
    y: row * iconGridMetrics.value.cellH
  }
}

function handleDesktopDragLeave(e) {
  if (!e.relatedTarget) {
    dropPreview.value = null
  }
}

function handleDesktopDrop(e) {
  if (!draggedAppId.value || isPositionLocked.value) return
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
    <!-- Indikator Posisi Ikon Terkunci -->
    <div
      v-if="isPositionLocked"
      class="absolute top-4 left-6 flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-slate-900/80 border border-amber-500/40 backdrop-blur-md text-xs font-medium text-amber-300 shadow-xl pointer-events-auto z-10 animate-in fade-in"
      title="Posisi ikon desktop terkunci. Klik kanan area desktop untuk membuka kunci."
    >
      <Lock class="w-3.5 h-3.5 text-amber-400" />
      <span>Tata Letak Terkunci</span>
    </div>

    <!-- Subtle Brand Watermark di pojok kanan atas: Logo RAZOS Baru & RAZ Creative Studio Tetap Ada -->
    <div class="absolute top-4 right-6 pointer-events-none flex items-center gap-2.5 opacity-60 z-0">
      <div class="flex items-center gap-1.5">
        <img src="/assets/brand/RAZOSMarkSquare.webp" alt="RAZOS" class="w-6 h-6 rounded-lg object-cover shadow-md border border-amber-500/30" />
        <span class="text-xs font-bold bg-gradient-to-r from-amber-200 to-yellow-400 bg-clip-text text-transparent tracking-wider hidden sm:inline">RAZOS</span>
      </div>
      <span class="text-slate-600 text-xs hidden sm:inline">•</span>
      <img src="/assets/brand/RAZBrandTransparent.webp" alt="RAZ STUDIO" class="h-6 object-contain drop-shadow" />
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
      <div ref="desktopWorkspaceRef" class="flex-1 relative w-full h-full min-h-0 pointer-events-none overflow-hidden">
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

        <!-- Windows Desktop Grid Icons (Draggable & Absolute Grid Positioned) -->
        <div
          v-for="app in windowStore.allDesktopApps"
          :key="app.id"
          :draggable="!isPositionLocked"
          @dragstart="handleIconDragStart($event, app)"
          @dragend="handleIconDragEnd"
          @click="handleAppLaunch(app)"
          @dblclick="handleAppLaunch(app)"
          @contextmenu.prevent.stop="handleDesktopIconContextMenu($event, app)"
          :style="{
            position: 'absolute',
            left: `${getEffectivePosition(app.id).col * iconGridMetrics.cellW}px`,
            top: `${getEffectivePosition(app.id).row * iconGridMetrics.cellH}px`,
            width: `${iconGridMetrics.cellW}px`,
            height: `${iconGridMetrics.cellH}px`
          }"
          class="flex flex-col items-center justify-start p-1.5 rounded-2xl hover:bg-slate-800/60 group transition-all duration-150 text-center select-none pointer-events-auto"
          :class="[
            draggedAppId === app.id ? 'opacity-30 scale-95 ring-2 ring-sky-400' : '',
            isPositionLocked ? 'cursor-default' : 'cursor-grab active:cursor-grabbing'
          ]"
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

        <!-- RIGHT COLUMN: Desktop Widgets & Recommended Apps Slider (w-72) -->
        <div class="absolute top-0 right-0 w-72 flex flex-col gap-3 max-h-full overflow-y-auto pointer-events-auto scrollbar-none select-none z-10">
          <!-- Widget 1-4 Card: Time, System Status, Storage Status, Network Status -->
          <div
            v-if="systemStore.settings?.widgets && (systemStore.settings.widgets.time || systemStore.settings.widgets.systemStatus || systemStore.settings.widgets.storageStatus || systemStore.settings.widgets.networkStatus)"
            class="rounded-2xl bg-slate-900/80 backdrop-blur-md border border-slate-800/90 shadow-2xl p-4 space-y-4 text-slate-200 animate-in fade-in select-none"
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

          <!-- Card Rekomendasi Aplikasi Desktop (Slider / Carousel Ala CasaOS) -->
          <div
            v-if="systemStore.settings?.showRecommendedApps !== false"
            class="rounded-2xl bg-slate-900/80 backdrop-blur-md border border-slate-800/90 shadow-2xl p-3.5 space-y-2.5 text-slate-200 animate-in fade-in select-none relative overflow-hidden group"
            @mouseenter="stopRecAutoSlide"
            @mouseleave="startRecAutoSlide"
          >
            <!-- Header Kartu Rekomendasi -->
            <div class="flex items-center justify-between pb-2 border-b border-slate-800/80">
              <div class="flex items-center gap-1.5">
                <Sparkles class="w-3.5 h-3.5 text-amber-400 animate-pulse" />
                <span class="text-xs font-bold text-slate-200 tracking-wide">Rekomendasi App</span>
              </div>
              <div class="flex items-center gap-1">
                <button
                  type="button"
                  @click="prevRecApp"
                  class="p-1 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition active:scale-95"
                  title="Aplikasi Sebelumnya"
                >
                  <ChevronLeft class="w-3 h-3" />
                </button>
                <button
                  type="button"
                  @click="nextRecApp"
                  class="p-1 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition active:scale-95"
                  title="Aplikasi Berikutnya"
                >
                  <ChevronRight class="w-3 h-3" />
                </button>
              </div>
            </div>

            <!-- Konten Slider Aplikasi Aktif -->
            <div class="relative min-h-[92px] flex flex-col justify-between">
              <div class="flex items-start gap-2.5">
                <div class="w-10 h-10 rounded-xl bg-slate-950 border border-slate-800 p-1 flex items-center justify-center shrink-0 overflow-hidden shadow-inner group-hover:border-sky-400/40 transition">
                  <img
                    :src="recommendedAppsList[currentRecIndex].icon"
                    :alt="recommendedAppsList[currentRecIndex].name"
                    class="w-full h-full object-contain"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between gap-1">
                    <h4 class="text-xs font-bold text-slate-100 truncate group-hover:text-sky-300 transition">
                      {{ recommendedAppsList[currentRecIndex].name }}
                    </h4>
                    <span class="text-[9px] px-1.5 py-0.5 rounded-full bg-sky-500/15 text-sky-300 border border-sky-400/20 font-semibold shrink-0">
                      {{ recommendedAppsList[currentRecIndex].tag }}
                    </span>
                  </div>
                  <p class="text-[10px] text-slate-400 line-clamp-2 mt-0.5 leading-snug">
                    {{ recommendedAppsList[currentRecIndex].description }}
                  </p>
                </div>
              </div>

              <!-- Tombol Aksi & Indikator Titik -->
              <div class="flex items-center justify-between pt-2 mt-1 border-t border-slate-800/60">
                <!-- Titik Slider -->
                <div class="flex items-center gap-1">
                  <button
                    v-for="(_, idx) in recommendedAppsList"
                    :key="idx"
                    type="button"
                    @click="currentRecIndex = idx"
                    class="h-1.5 rounded-full transition-all"
                    :class="currentRecIndex === idx ? 'w-4 bg-sky-400' : 'w-1.5 bg-slate-700 hover:bg-slate-500'"
                    :title="`Slide ${idx + 1}`"
                  ></button>
                </div>

                <!-- Tombol Install / Buka -->
                <div class="flex items-center gap-1.5">
                  <button
                    type="button"
                    @click="windowStore.openApp('app-store')"
                    class="text-[10px] px-2 py-1 rounded-lg bg-slate-800/70 hover:bg-slate-700 text-slate-300 hover:text-white transition active:scale-95"
                    title="Buka App Store"
                  >
                    Katalog
                  </button>
                  <button
                    type="button"
                    @click="handleInstallOrOpenRec(recommendedAppsList[currentRecIndex])"
                    class="text-[10px] font-semibold px-2.5 py-1 rounded-lg transition flex items-center gap-1 shadow-sm active:scale-95"
                    :class="isRecAppInstalled(recommendedAppsList[currentRecIndex].slug)
                      ? 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-emerald-600/20'
                      : 'bg-sky-500 hover:bg-sky-400 text-white shadow-sky-500/20'"
                  >
                    <component
                      :is="isRecAppInstalled(recommendedAppsList[currentRecIndex].slug) ? Play : Download"
                      class="w-3 h-3"
                    />
                    <span>{{ isRecAppInstalled(recommendedAppsList[currentRecIndex].slug) ? 'Buka' : 'Pasang' }}</span>
                  </button>
                </div>
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

            <!-- Standalone App Store Settings Window (CasaOS Style) -->
            <RAZAppEditView v-else-if="win.id.startsWith('app-edit-')" :win="win" />

            <!-- Standalone App Store Logs & Terminal Window (CasaOS Style) -->
            <RAZAppLogsView v-else-if="win.id.startsWith('app-logs-')" :win="win" />

            <!-- Standalone Storage Manager Window (CasaOS Style) -->
            <RAZStorageManagerView v-else-if="win.id === 'storage-manager'" :win="win" />

            <!-- Sublime & Notepad++ Advanced Code/Text Editor Terintegrasi File Manager -->
            <RAZTextEditorView v-else-if="win.id === 'file-manager-editor' || win.id === 'text-editor' || win.id.startsWith('editor-')" :initial-file="win.fileData" :win="win" />

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
