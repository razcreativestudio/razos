<script setup>
/**
 * RAZOS App Store View - 525+ Apps Catalog, Installed Apps Management & CasaOS-Style YAML Docker Compose Suite
 * Updated: 2026-10-13 (v2.5.5 - Floating Window Dialog Modals without Background Blur)
 * Deskripsi:
 * - 3 Tab Navigasi Utama:
 *   1. "Katalog Aplikasi (525+)" - Jelajah katalog lengkap dengan filter kategori & live search.
 *   2. "Aplikasi Terpasang" - Kelola semua container yang sedang aktif, inspect port & mountpoint, peluncuran Web UI, serta uninstall bersih.
 *   3. "Instal via YAML (Docker Compose)" - Editor YAML ala CasaOS dengan live parse syntax, drag-and-drop / upload file .yaml/.yml, template preset cepat, dan direct deployment.
 * - Halaman Penuh Detail Aplikasi (Full Page Markdown README.md).
 * - Integrasi mulus dengan RAZWindowStore untuk sinkronisasi ikon Desktop & Start Menu secara instan.
 * - Window Dialog Modals (Tanpa Blur background, visual jendela independen OS).
 */
import { ref, onMounted, computed, watch } from 'vue'
import { marked } from 'marked'
import YAML from 'yaml'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Download, Check, Box, Trash2, Settings2, X, AlertTriangle,
  Search, BookOpen, PlusCircle, ArrowLeft, ExternalLink, HardDrive,
  Cpu, Terminal, Layers, Globe, Shield, Activity, RefreshCw,
  FileCode, Upload, Copy, Play, CheckCircle2, Info, Folder, Sparkles, Server
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()

// State Navigasi Tab: 'catalog' | 'installed' | 'compose'
const activeTab = ref('catalog')

// State Halaman: 'catalog' | 'detail'
const currentView = ref('catalog')
const selectedApp = ref(null)

const apps = ref([])
const isLoading = ref(false)
const selectedCategory = ref('All')
const searchQuery = ref('')

// State Modal Pasang Aplikasi Katalog
const showInstallModal = ref(false)
const activeAppToInstall = ref(null)
const customPort = ref('')
const customVolume = ref('')
const isInstalling = ref(false)

// State Modal Clean Uninstall
const showUninstallModal = ref(false)
const activeAppToUninstall = ref(null)
const purgeDataChecked = ref(true)
const isUninstalling = ref(false)

// State Modal Manual Install Docker (Form Cepat)
const showManualModal = ref(false)
const manualName = ref('')
const manualImage = ref('')
const manualPort = ref('8080')
const manualVolume = ref('/DATA/apps/custom-app')
const manualEnvs = ref('')
const isDeployingManual = ref(false)

// State YAML Compose (CasaOS Style)
const composeYaml = ref(`version: '3.8'
services:
  qbittorrent:
    image: linuxserver/qbittorrent:latest
    container_name: qbittorrent
    restart: unless-stopped
    ports:
      - "8085:8080"
      - "6881:6881"
      - "6881:6881/udp"
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Asia/Jakarta
      - WEBUI_PORT=8080
    volumes:
      - /DATA/apps/qbittorrent/config:/config
      - /DATA/apps/qbittorrent/downloads:/downloads
`)
const yamlError = ref('')
const isDeployingYaml = ref(false)
const dragActive = ref(false)

// Template Preset Docker Compose Siap Pakai
const composePresets = [
  {
    name: 'qBittorrent (Downloader)',
    yaml: `version: '3.8'
services:
  qbittorrent:
    image: linuxserver/qbittorrent:latest
    container_name: qbittorrent
    restart: unless-stopped
    ports:
      - "8085:8080"
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Asia/Jakarta
    volumes:
      - /DATA/apps/qbittorrent/config:/config
      - /DATA/apps/qbittorrent/downloads:/downloads
`
  },
  {
    name: 'Nginx Web Server',
    yaml: `version: '3.8'
services:
  webserver:
    image: nginx:alpine
    container_name: nginx-web
    restart: unless-stopped
    ports:
      - "8080:80"
    volumes:
      - /DATA/apps/nginx/html:/usr/share/nginx/html
      - /DATA/apps/nginx/conf.d:/etc/nginx/conf.d
`
  },
  {
    name: 'Nextcloud Hub',
    yaml: `version: '3.8'
services:
  nextcloud:
    image: nextcloud:fpm-alpine
    container_name: nextcloud
    restart: unless-stopped
    ports:
      - "8081:80"
    volumes:
      - /DATA/apps/nextcloud/data:/var/www/html
    environment:
      - SQLITE_DATABASE=nextcloud
`
  },
  {
    name: 'WordPress CMS',
    yaml: `version: '3.8'
services:
  wordpress:
    image: wordpress:php8.2-apache
    container_name: wordpress
    restart: unless-stopped
    ports:
      - "8082:80"
    environment:
      - WORDPRESS_DB_NAME=wordpress
    volumes:
      - /DATA/apps/wordpress:/var/www/html
`
  }
]

// Kategori Katalog
const categories = [
  'All',
  'Media & Streaming',
  'File & Cloud Storage',
  'Database & Caching',
  'Development & Git',
  'Home Automation & IoT',
  'Networking & Security',
  'Monitoring & Analytics',
  'Productivity & Office',
  'Downloaders & Torrents',
  'AI & Machine Learning'
]

// Ambil daftar aplikasi dari backend
async function fetchApps() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/store/apps')
    if (res.ok) {
      apps.value = await res.json()
    }
  } catch (err) {
    console.error('Gagal mengambil daftar aplikasi:', err)
  }
  isLoading.value = false
}

// Filtered apps untuk Katalog
const filteredApps = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  return apps.value.filter(a => {
    const matchCat = selectedCategory.value === 'All' || a.category === selectedCategory.value
    if (!matchCat) return false
    if (!query) return true
    return (
      a.name.toLowerCase().includes(query) ||
      a.description.toLowerCase().includes(query) ||
      a.image.toLowerCase().includes(query) ||
      a.category.toLowerCase().includes(query)
    )
  })
})

// Daftar Aplikasi yang Sedang Terpasang
const installedApps = computed(() => {
  return apps.value.filter(a => a.installed)
})

// Live Parser Compose YAML
const parsedComposeServices = computed(() => {
  yamlError.value = ''
  if (!composeYaml.value.trim()) return []
  try {
    const doc = YAML.parse(composeYaml.value)
    if (!doc || typeof doc !== 'object' || !doc.services) {
      yamlError.value = 'Dokumen YAML wajib memiliki blok "services:".'
      return []
    }

    const services = []
    for (const [key, svc] of Object.entries(doc.services)) {
      if (!svc || typeof svc !== 'object') continue
      // Ekstrak port host
      let hostPort = '-'
      if (Array.isArray(svc.ports) && svc.ports.length > 0) {
        const p = String(svc.ports[0])
        const parts = p.split(':')
        hostPort = parts[0] || '-'
      }

      // Ekstrak volumes
      const vols = Array.isArray(svc.volumes) ? svc.volumes.map(String) : []

      // Ekstrak envs
      let envCount = 0
      if (Array.isArray(svc.environment)) {
        envCount = svc.environment.length
      } else if (svc.environment && typeof svc.environment === 'object') {
        envCount = Object.keys(svc.environment).length
      }

      services.push({
        id: key,
        name: svc.container_name || key,
        image: svc.image || 'Belum ditentukan',
        hostPort,
        volumes: vols,
        envCount,
        restart: svc.restart || 'no'
      })
    }
    return services
  } catch (e) {
    yamlError.value = `Syntax Error: ${e.message}`
    return []
  }
})

// Buka Halaman Penuh Detail Aplikasi
function openAppDetail(app) {
  selectedApp.value = app
  currentView.value = 'detail'
}

function backToCatalog() {
  currentView.value = 'catalog'
  selectedApp.value = null
}

// Buka Web UI Container (Meluncurkan via RAZWindowStore atau Tab Peramban)
function launchApp(app) {
  windowStore.launchStoreApp(app)
}

// Render Markdown README.md ke HTML Preview
const renderedReadmeHtml = computed(() => {
  if (!selectedApp.value || !selectedApp.value.readme) {
    return '<p class="text-slate-400">Tidak ada dokumentasi README.md yang disertakan.</p>'
  }
  try {
    return marked.parse(selectedApp.value.readme)
  } catch (e) {
    return `<pre class="text-rose-400">${selectedApp.value.readme}</pre>`
  }
})

// Buka Dialog Pasang Aplikasi
function openInstallDialog(app, e) {
  e?.stopPropagation()
  activeAppToInstall.value = app
  customPort.value = app.default_port
  customVolume.value = app.default_volume
  showInstallModal.value = true
}

// Eksekusi Pasang Aplikasi Katalog
async function confirmInstall() {
  if (!activeAppToInstall.value) return
  isInstalling.value = true
  const payload = {
    app_id: activeAppToInstall.value.id,
    port: parseInt(customPort.value) || activeAppToInstall.value.default_port,
    data_path: customVolume.value
  }

  try {
    const res = await fetch('/api/v1/store/install', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container ${activeAppToInstall.value.name} berhasil dipasang!`, 'success')
      showInstallModal.value = false
      if (selectedApp.value && selectedApp.value.id === activeAppToInstall.value.id) {
        selectedApp.value.installed = true
      }
      await fetchApps()
      await windowStore.fetchInstalledApps()
    } else {
      const errData = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Pasang', errData.error || 'Gagal memasang aplikasi', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isInstalling.value = false
  }
}

// Buka Dialog Uninstall Bersih
function openUninstallDialog(app, e) {
  e?.stopPropagation()
  activeAppToUninstall.value = app
  purgeDataChecked.value = true
  showUninstallModal.value = true
}

// Konfirmasi Uninstall
async function confirmUninstall() {
  if (!activeAppToUninstall.value) return
  isUninstalling.value = true
  try {
    const res = await fetch(`/api/v1/store/apps/${activeAppToUninstall.value.slug}?purge_data=${purgeDataChecked.value}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Aplikasi ${activeAppToUninstall.value.name} telah berhasil di-uninstall!`, 'success')
      showUninstallModal.value = false
      if (selectedApp.value && selectedApp.value.id === activeAppToUninstall.value.id) {
        selectedApp.value.installed = false
      }
      await fetchApps()
      await windowStore.fetchInstalledApps()
    } else {
      systemStore.showNotification('Gagal', 'Gagal meng-uninstall container', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isUninstalling.value = false
  }
}

// Submit Manual Custom Docker Container
async function submitManualInstall() {
  if (!manualName.value.trim() || !manualImage.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama aplikasi dan Docker Image wajib diisi!', 'warning')
    return
  }

  isDeployingManual.value = true
  const envList = manualEnvs.value
    ? manualEnvs.value.split('\n').map(s => s.trim()).filter(s => s.length > 0)
    : []

  const payload = {
    name: manualName.value.trim(),
    custom_image: manualImage.value.trim(),
    container_name: manualName.value.toLowerCase().replace(/[^a-z0-9]/g, '-'),
    port: parseInt(manualPort.value) || 8080,
    data_path: manualVolume.value.trim(),
    extra_envs: envList
  }

  try {
    const res = await fetch('/api/v1/store/install', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container kustom ${manualName.value} berhasil di-deploy!`, 'success')
      showManualModal.value = false
      manualName.value = ''
      manualImage.value = ''
      manualPort.value = '8080'
      manualVolume.value = '/DATA/apps/custom-app'
      manualEnvs.value = ''
      await fetchApps()
      await windowStore.fetchInstalledApps()
      activeTab.value = 'installed'
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Pasang', err.error || 'Gagal memasang custom container', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isDeployingManual.value = false
  }
}

// Deploy Docker Compose YAML (CasaOS Style)
async function submitComposeDeployment() {
  if (!composeYaml.value.trim()) {
    systemStore.showNotification('Peringatan', 'Skrip YAML Docker Compose tidak boleh kosong!', 'warning')
    return
  }

  if (yamlError.value) {
    systemStore.showNotification('Syntax Error', 'Perbaiki error pada penulisan YAML sebelum deployment!', 'error')
    return
  }

  isDeployingYaml.value = true
  try {
    const res = await fetch('/api/v1/store/compose', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ yaml: composeYaml.value })
    })

    if (res.ok) {
      const result = await res.json()
      systemStore.showNotification(
        'Deployment Sukses',
        `Berhasil men-deploy ${result.services?.length || 1} container dari skrip YAML Compose!`,
        'success'
      )
      await fetchApps()
      await windowStore.fetchInstalledApps()
      activeTab.value = 'installed'
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Deploy', err.error || 'Gagal memproses skrip Compose YAML', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi backend RAZOS', 'error')
  } finally {
    isDeployingYaml.value = false
  }
}

// Handle File Drag & Drop / Upload untuk YAML
function handleFileUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  readFileContent(file)
}

function handleFileDrop(e) {
  dragActive.value = false
  const file = e.dataTransfer.files?.[0]
  if (!file) return
  if (!file.name.endsWith('.yaml') && !file.name.endsWith('.yml')) {
    systemStore.showNotification('Format Tidak Cocok', 'Harap unggah file berekstensi .yaml atau .yml!', 'warning')
    return
  }
  readFileContent(file)
}

function readFileContent(file) {
  const reader = new FileReader()
  reader.onload = (evt) => {
    composeYaml.value = evt.target.result
    systemStore.showNotification('File Dimuat', `File ${file.name} berhasil dimuat ke editor Compose!`, 'info')
  }
  reader.onerror = () => {
    systemStore.showNotification('Gagal Membaca File', 'Tidak dapat membaca file yang diunggah.', 'error')
  }
  reader.readAsText(file)
}

function applyPreset(preset) {
  composeYaml.value = preset.yaml
  systemStore.showNotification('Template Dimuat', `Template ${preset.name} berhasil dipasang ke editor.`, 'info')
}

onMounted(async () => {
  await fetchApps()
  await windowStore.fetchInstalledApps()
})
</script>

<template>
  <div class="h-full flex flex-col text-xs sm:text-sm select-none bg-slate-950 text-slate-100 overflow-hidden">
    <!-- TOP NAVIGATION BAR (TABS & ACTIONS) -->
    <div class="px-4 pt-3 pb-2 border-b border-slate-800/80 bg-slate-900/60 shrink-0 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
      <!-- Title & Tab Switches -->
      <div class="flex items-center gap-4 overflow-x-auto">
        <div class="flex items-center gap-2 pr-2 border-r border-slate-800">
          <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
            <Box class="w-4 h-4" />
          </div>
          <div>
            <h3 class="font-bold text-slate-100 text-sm leading-tight">App Store RAZOS</h3>
            <span class="text-[10px] text-slate-400 font-mono">Docker Suite STB</span>
          </div>
        </div>

        <!-- Tab Buttons Navigation -->
        <div class="flex items-center gap-1 bg-slate-950/80 p-1 rounded-xl border border-slate-800/80">
          <!-- Tab 1: Katalog -->
          <button
            @click="activeTab = 'catalog'; currentView = 'catalog'"
            class="px-3 py-1.5 rounded-lg font-medium text-xs flex items-center gap-1.5 transition whitespace-nowrap"
            :class="activeTab === 'catalog' && currentView === 'catalog'
              ? 'bg-sky-500 text-slate-950 font-bold shadow'
              : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/50'"
          >
            <Layers class="w-3.5 h-3.5" />
            <span>Katalog ({{ apps.length }})</span>
          </button>

          <!-- Tab 2: Aplikasi Terpasang -->
          <button
            @click="activeTab = 'installed'; currentView = 'catalog'"
            class="px-3 py-1.5 rounded-lg font-medium text-xs flex items-center gap-1.5 transition whitespace-nowrap"
            :class="activeTab === 'installed'
              ? 'bg-emerald-500 text-slate-950 font-bold shadow'
              : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/50'"
          >
            <CheckCircle2 class="w-3.5 h-3.5" />
            <span>Terpasang</span>
            <span
              class="px-1.5 py-0.2 rounded-full text-[10px] font-mono font-bold"
              :class="activeTab === 'installed' ? 'bg-slate-950/30 text-slate-950' : 'bg-emerald-500/20 text-emerald-400'"
            >
              {{ installedApps.length }}
            </span>
          </button>

          <!-- Tab 3: Docker Compose YAML (CasaOS Style) -->
          <button
            @click="activeTab = 'compose'; currentView = 'catalog'"
            class="px-3 py-1.5 rounded-lg font-medium text-xs flex items-center gap-1.5 transition whitespace-nowrap"
            :class="activeTab === 'compose'
              ? 'bg-indigo-500 text-white font-bold shadow'
              : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/50'"
          >
            <FileCode class="w-3.5 h-3.5" />
            <span>Instal via YAML</span>
          </button>
        </div>
      </div>

      <!-- Quick Action Buttons -->
      <div class="flex items-center gap-2 shrink-0">
        <button
          @click="fetchApps"
          class="p-2 rounded-xl bg-slate-800/80 hover:bg-slate-700/80 text-slate-300 hover:text-white transition"
          title="Segarkan Daftar Aplikasi"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
        </button>

        <button
          @click="showManualModal = true"
          class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold flex items-center gap-1.5 text-xs transition border border-slate-700"
        >
          <PlusCircle class="w-3.5 h-3.5 text-sky-400" />
          <span class="hidden md:inline">Instal Cepat</span>
        </button>
      </div>
    </div>

    <!-- MAIN BODY CONTAINER -->
    <div class="flex-1 overflow-hidden relative">

      <!-- ================= TAB 1: KATALOG (525+ APPS) ================= -->
      <div v-if="activeTab === 'catalog' && currentView === 'catalog'" class="h-full flex flex-col gap-3 p-4 overflow-hidden animate-in fade-in">
        <!-- Search Bar & Statistics -->
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5 shrink-0">
          <div class="flex-1 relative">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari dari 525+ aplikasi Docker (Plex, Jellyfin, MySQL, Nextcloud, Ollama, Nginx, qBittorrent...)"
              class="w-full pl-9 pr-3 py-2 rounded-xl bg-slate-900/80 border border-slate-800 text-xs text-slate-100 placeholder-slate-500 outline-none focus:border-sky-400 transition"
            />
          </div>

          <span class="text-[11px] text-slate-500 font-mono hidden md:inline shrink-0">
            Ditemukan: <span class="text-sky-400 font-bold">{{ filteredApps.length }}</span> aplikasi
          </span>
        </div>

        <!-- Horizontal Category Filter Bar -->
        <div class="flex items-center gap-1.5 overflow-x-auto pb-1 shrink-0 scrollbar-none">
          <button
            v-for="cat in categories"
            :key="cat"
            @click="selectedCategory = cat"
            class="px-3 py-1.5 rounded-xl text-xs font-medium transition whitespace-nowrap"
            :class="selectedCategory === cat
              ? 'bg-sky-500/20 text-sky-400 border border-sky-400/40 font-bold shadow'
              : 'text-slate-400 hover:text-slate-200 bg-slate-900/60 border border-slate-800/60'"
          >
            {{ cat }}
          </button>
        </div>

        <!-- App Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5 overflow-y-auto flex-1 pr-1">
          <div
            v-for="app in filteredApps"
            :key="app.id"
            @click="openAppDetail(app)"
            class="p-4 rounded-2xl bg-slate-900/70 border border-slate-800/80 hover:border-sky-500/50 hover:bg-slate-900/90 flex flex-col justify-between gap-3.5 transition-all cursor-pointer group relative shadow-sm hover:shadow-xl hover:-translate-y-0.5"
          >
            <!-- Card Top: Logo Resmi/Identicon & Meta -->
            <div class="flex items-start gap-3.5">
              <div class="w-12 h-12 rounded-2xl bg-slate-950 border border-slate-800/80 p-1 flex items-center justify-center shrink-0 overflow-hidden shadow-inner group-hover:border-sky-400/40 transition">
                <img
                  v-if="app.icon"
                  :src="app.icon"
                  :alt="app.name"
                  class="w-full h-full object-contain rounded-xl"
                  loading="lazy"
                  @error="(e) => e.target.src = `https://api.dicebear.com/7.x/identicon/svg?seed=${app.slug}`"
                />
                <Box v-else class="w-6 h-6 text-sky-400" />
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-1">
                  <h4 class="font-bold text-slate-100 text-sm truncate group-hover:text-sky-300 transition" :title="app.name">
                    {{ app.name }}
                  </h4>
                  <span class="text-[9px] px-2 py-0.5 rounded-full bg-slate-800 text-slate-400 shrink-0 font-medium border border-slate-700/60">
                    {{ app.category }}
                  </span>
                </div>
                <p class="text-xs text-slate-400 mt-1 line-clamp-2 leading-relaxed">
                  {{ app.description }}
                </p>
              </div>
            </div>

            <!-- Card Bottom: Ports & Action Buttons -->
            <div class="flex items-center justify-between border-t border-slate-800/60 pt-2.5 text-xs">
              <div class="text-slate-500 font-mono flex items-center gap-1.5">
                <span>Port:</span>
                <span class="text-sky-400 font-bold">{{ app.default_port }}</span>
              </div>

              <div class="flex items-center gap-2">
                <!-- Jika Terpasang -->
                <template v-if="app.installed">
                  <button
                    @click.stop="launchApp(app)"
                    class="px-2.5 py-1 rounded-xl bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 font-bold flex items-center gap-1 text-[11px] transition"
                    title="Buka Aplikasi Web UI"
                  >
                    <ExternalLink class="w-3.5 h-3.5" />
                    <span>Buka</span>
                  </button>
                  <button
                    @click="(e) => openUninstallDialog(app, e)"
                    class="p-1.5 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 transition"
                    title="Hapus / Uninstall Bersih"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </template>

                <!-- Jika Belum Terpasang -->
                <button
                  v-else
                  @click="(e) => openInstallDialog(app, e)"
                  class="px-3.5 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow active:scale-95"
                >
                  <Download class="w-3.5 h-3.5" />
                  <span>Pasang</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= TAB 2: APLIKASI TERPASANG (INSTALLED APPS) ================= -->
      <div v-else-if="activeTab === 'installed' && currentView === 'catalog'" class="h-full flex flex-col gap-3 p-4 overflow-hidden animate-in fade-in">
        <!-- Header Info Tab Terpasang -->
        <div class="flex items-center justify-between border-b border-slate-800/80 pb-3 shrink-0">
          <div>
            <h4 class="font-bold text-slate-100 text-base">Aplikasi & Container Aktif</h4>
            <p class="text-xs text-slate-400 mt-0.5">
              Daftar container Docker yang saat ini terpasang di STB Anda dan siap diakses.
            </p>
          </div>
          <span class="px-3 py-1 rounded-full bg-emerald-500/20 text-emerald-400 text-xs font-mono font-bold">
            {{ installedApps.length }} Container Berjalan
          </span>
        </div>

        <!-- Empty State Jika Belum Ada Aplikasi Terpasang -->
        <div v-if="installedApps.length === 0" class="flex-1 flex flex-col items-center justify-center p-8 text-center text-slate-400 gap-3">
          <div class="w-16 h-16 rounded-3xl bg-slate-900 border border-slate-800 flex items-center justify-center text-slate-500 mb-2">
            <Box class="w-8 h-8" />
          </div>
          <h4 class="text-base font-bold text-slate-200">Belum Ada Aplikasi yang Terpasang</h4>
          <p class="text-xs text-slate-400 max-w-md">
            Anda belum memasang aplikasi dari katalog maupun dari Docker Compose. Pasang aplikasi untuk memaksimalkan STB Anda!
          </p>
          <div class="flex items-center gap-2 mt-2">
            <button
              @click="activeTab = 'catalog'"
              class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition flex items-center gap-2 shadow"
            >
              <Layers class="w-4 h-4" />
              <span>Buka Katalog (525+)</span>
            </button>
            <button
              @click="activeTab = 'compose'"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold text-xs transition flex items-center gap-2 border border-slate-700"
            >
              <FileCode class="w-4 h-4 text-indigo-400" />
              <span>Impor YAML Compose</span>
            </button>
          </div>
        </div>

        <!-- Grid Container Terpasang -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5 overflow-y-auto flex-1 pr-1">
          <div
            v-for="app in installedApps"
            :key="app.id"
            class="p-4 rounded-2xl bg-slate-900/80 border border-emerald-500/30 hover:border-emerald-400/50 flex flex-col justify-between gap-3.5 transition-all shadow-md"
          >
            <!-- Top Container Info -->
            <div class="flex items-start gap-3.5">
              <div class="w-12 h-12 rounded-2xl bg-slate-950 border border-emerald-500/40 p-1 flex items-center justify-center shrink-0 overflow-hidden shadow-inner">
                <img
                  v-if="app.icon"
                  :src="app.icon"
                  :alt="app.name"
                  class="w-full h-full object-contain rounded-xl"
                  @error="(e) => e.target.src = `https://api.dicebear.com/7.x/identicon/svg?seed=${app.slug}`"
                />
                <Box v-else class="w-6 h-6 text-emerald-400" />
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-1">
                  <h4 class="font-bold text-slate-100 text-sm truncate" :title="app.name">
                    {{ app.name }}
                  </h4>
                  <span class="text-[9px] px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center gap-1 border border-emerald-500/30">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
                    Running
                  </span>
                </div>
                <p class="text-xs text-slate-400 mt-1 line-clamp-1 font-mono">
                  {{ app.image }}
                </p>
                <div class="mt-2 text-[11px] text-slate-400 flex flex-col gap-1 font-mono bg-slate-950/60 p-2 rounded-xl border border-slate-800/80">
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Port Web UI:</span>
                    <span class="text-sky-400 font-bold">{{ app.default_port }}</span>
                  </div>
                  <div class="flex items-center justify-between truncate">
                    <span class="text-slate-500 shrink-0">Path Data:</span>
                    <span class="text-slate-300 truncate ml-2" :title="app.default_volume">{{ app.default_volume }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Bottom Actions -->
            <div class="flex items-center justify-between border-t border-slate-800/60 pt-2.5 text-xs">
              <button
                @click="openAppDetail(app)"
                class="text-xs text-slate-400 hover:text-sky-300 font-medium flex items-center gap-1 transition"
              >
                <BookOpen class="w-3.5 h-3.5" />
                <span>Detail</span>
              </button>

              <div class="flex items-center gap-2">
                <button
                  @click="launchApp(app)"
                  class="px-3 py-1.5 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-400 hover:to-teal-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow active:scale-95"
                >
                  <ExternalLink class="w-3.5 h-3.5" />
                  <span>Buka Web UI</span>
                </button>

                <button
                  @click="(e) => openUninstallDialog(app, e)"
                  class="px-2.5 py-1.5 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 font-bold flex items-center gap-1 text-xs transition border border-rose-500/30"
                  title="Hapus Aplikasi Bersih"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                  <span>Uninstall</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= TAB 3: DOCKER COMPOSE YAML (CASAOS STYLE) ================= -->
      <div v-else-if="activeTab === 'compose' && currentView === 'catalog'" class="h-full flex flex-col lg:flex-row gap-4 p-4 overflow-hidden animate-in fade-in">
        <!-- Panel Kiri: Editor Monospace & Upload/Preset -->
        <div class="flex-1 flex flex-col gap-3 overflow-hidden">
          <!-- Header Bar Editor -->
          <div class="flex items-center justify-between shrink-0">
            <div>
              <h4 class="font-bold text-slate-100 text-base flex items-center gap-2">
                <FileCode class="w-5 h-5 text-indigo-400" />
                <span>Editor Docker Compose (CasaOS Style)</span>
              </h4>
              <p class="text-xs text-slate-400 mt-0.5">
                Ketik atau tempel skrip Docker Compose YAML standar untuk men-deploy multi-container secara instan.
              </p>
            </div>

            <!-- Upload File Trigger -->
            <label class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-bold cursor-pointer transition flex items-center gap-1.5 border border-slate-700">
              <Upload class="w-3.5 h-3.5 text-sky-400" />
              <span>Unggah .yaml / .yml</span>
              <input type="file" accept=".yaml,.yml" class="hidden" @change="handleFileUpload" />
            </label>
          </div>

          <!-- Quick Template Presets -->
          <div class="flex items-center gap-1.5 overflow-x-auto pb-1 shrink-0 scrollbar-none">
            <span class="text-[11px] text-slate-400 font-bold shrink-0">Template Cepat:</span>
            <button
              v-for="preset in composePresets"
              :key="preset.name"
              @click="applyPreset(preset)"
              class="px-2.5 py-1 rounded-lg bg-slate-900 border border-slate-800 hover:border-indigo-400/50 text-[11px] text-slate-300 hover:text-white transition whitespace-nowrap"
            >
              {{ preset.name }}
            </button>
          </div>

          <!-- Drag and Drop Textarea Container -->
          <div
            class="flex-1 relative flex flex-col rounded-2xl border transition-all overflow-hidden bg-slate-950"
            :class="dragActive ? 'border-indigo-400 bg-indigo-950/20' : 'border-slate-800 focus-within:border-indigo-500/70'"
            @dragover.prevent="dragActive = true"
            @dragleave.prevent="dragActive = false"
            @drop.prevent="handleFileDrop"
          >
            <div class="px-3 py-1.5 bg-slate-900/80 border-b border-slate-800/80 flex items-center justify-between text-[11px] text-slate-400">
              <span class="font-mono">docker-compose.yaml</span>
              <span v-if="yamlError" class="text-rose-400 font-mono text-[10px] truncate max-w-sm">{{ yamlError }}</span>
              <span v-else class="text-emerald-400 font-mono text-[10px] flex items-center gap-1">
                <Check class="w-3 h-3" /> Syntax Valid
              </span>
            </div>

            <textarea
              v-model="composeYaml"
              rows="12"
              class="flex-1 w-full p-4 bg-transparent text-slate-100 font-mono text-xs outline-none resize-none leading-relaxed selection:bg-indigo-500/40"
              placeholder="# Tempel konfigurasi docker-compose.yml Anda di sini..."
              spellcheck="false"
            ></textarea>
          </div>
        </div>

        <!-- Panel Kanan: Live Inspector Services & Tombol Deployment -->
        <div class="w-full lg:w-80 flex flex-col justify-between gap-3 bg-slate-900/60 border border-slate-800/80 rounded-2xl p-4 shrink-0 overflow-y-auto">
          <div class="space-y-3">
            <div class="border-b border-slate-800 pb-2">
              <h5 class="font-bold text-slate-100 text-xs uppercase tracking-wider flex items-center gap-1.5">
                <Sparkles class="w-3.5 h-3.5 text-indigo-400" />
                <span>Live Inspector Container</span>
              </h5>
              <p class="text-[11px] text-slate-400 mt-0.5">
                Preview container yang akan dibuat di direktori MicroSD <span class="font-mono text-slate-300">/DATA/apps/</span>.
              </p>
            </div>

            <!-- Jika Error Syntax -->
            <div v-if="yamlError" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-400 text-xs">
              <div class="flex items-center gap-1.5 font-bold mb-1">
                <AlertTriangle class="w-4 h-4 shrink-0" />
                <span>Error Parsing YAML</span>
              </div>
              <p class="text-[11px]">{{ yamlError }}</p>
            </div>

            <!-- List Preview Services -->
            <div v-else class="space-y-2.5">
              <div
                v-for="svc in parsedComposeServices"
                :key="svc.id"
                class="p-3 rounded-xl bg-slate-950 border border-slate-800/90 text-xs space-y-1.5"
              >
                <div class="flex items-center justify-between">
                  <span class="font-bold text-indigo-300 truncate">{{ svc.name }}</span>
                  <span class="text-[9px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 font-mono">Port: {{ svc.hostPort }}</span>
                </div>
                <div class="text-[11px] text-slate-400 font-mono truncate" :title="svc.image">
                  Image: <span class="text-slate-200">{{ svc.image }}</span>
                </div>
                <div class="text-[10px] text-slate-500 flex items-center gap-2 font-mono">
                  <span>Vols: {{ svc.volumes.length }}</span>
                  <span>•</span>
                  <span>Envs: {{ svc.envCount }}</span>
                  <span>•</span>
                  <span>Restart: {{ svc.restart }}</span>
                </div>
              </div>

              <div v-if="parsedComposeServices.length === 0" class="text-center py-6 text-slate-500 text-xs">
                Belum ada container terdeteksi dalam YAML.
              </div>
            </div>
          </div>

          <!-- Tombol Deploy YAML -->
          <div class="pt-3 border-t border-slate-800">
            <button
              @click="submitComposeDeployment"
              :disabled="isDeployingYaml || !!yamlError || parsedComposeServices.length === 0"
              class="w-full py-2.5 rounded-xl bg-gradient-to-r from-indigo-500 to-sky-500 hover:from-indigo-400 hover:to-sky-400 disabled:opacity-50 text-white font-bold text-xs flex items-center justify-center gap-2 shadow-lg transition active:scale-95"
            >
              <Play class="w-4 h-4 fill-white" />
              <span>{{ isDeployingYaml ? 'Mengeksekusi Deployment...' : 'Mulai Deployment Compose' }}</span>
            </button>
            <span class="text-[10px] text-slate-500 block text-center mt-1.5">
              File data otomatis dipetakan ke MicroSD /DATA/apps/
            </span>
          </div>
        </div>
      </div>

      <!-- ================= VIEW 2: HALAMAN PENUH DETAIL APLIKASI (FULL PAGE DETAIL) ================= -->
      <div v-else-if="currentView === 'detail' && selectedApp" class="h-full flex flex-col overflow-hidden animate-in fade-in">
        <!-- Top Navigation & Action Header -->
        <div class="p-4 border-b border-slate-800/80 bg-slate-900/60 flex items-center justify-between gap-3 shrink-0">
          <div class="flex items-center gap-3">
            <button
              @click="backToCatalog"
              class="p-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 transition flex items-center gap-1.5 text-xs font-bold"
              title="Kembali ke Katalog"
            >
              <ArrowLeft class="w-4 h-4" />
              <span class="hidden sm:inline">Kembali</span>
            </button>

            <div class="flex items-center gap-3">
              <div class="w-11 h-11 rounded-2xl bg-slate-950 border border-slate-800 p-1 flex items-center justify-center overflow-hidden">
                <img
                  :src="selectedApp.icon"
                  :alt="selectedApp.name"
                  class="w-full h-full object-contain rounded-xl"
                  @error="(e) => e.target.src = `https://api.dicebear.com/7.x/identicon/svg?seed=${selectedApp.slug}`"
                />
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <h2 class="text-base font-bold text-slate-100">{{ selectedApp.name }}</h2>
                  <span class="text-[10px] px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-400 font-medium">
                    {{ selectedApp.category }}
                  </span>
                </div>
                <span class="text-xs text-slate-400 font-mono">{{ selectedApp.image }}</span>
              </div>
            </div>
          </div>

          <!-- Tombol Aksi Pasang / Status di Header Detail -->
          <div class="flex items-center gap-2">
            <template v-if="selectedApp.installed">
              <button
                @click="launchApp(selectedApp)"
                class="px-3.5 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow"
              >
                <ExternalLink class="w-4 h-4" />
                <span>Buka Web UI</span>
              </button>
              <button
                @click="(e) => openUninstallDialog(selectedApp, e)"
                class="px-3.5 py-2 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 font-bold flex items-center gap-1.5 text-xs transition border border-rose-500/30"
              >
                <Trash2 class="w-4 h-4" />
                <span>Uninstall</span>
              </button>
            </template>

            <button
              v-else
              @click="(e) => openInstallDialog(selectedApp, e)"
              class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 shadow-lg transition active:scale-95"
            >
              <Download class="w-4 h-4" />
              <span>Pasang Aplikasi</span>
            </button>
          </div>
        </div>

        <!-- Main Detail Layout (Sidebar Metadata + Markdown Documentation) -->
        <div class="flex-1 flex overflow-hidden">
          <!-- Main Content: Markdown README.md Visual Preview -->
          <div class="flex-1 p-6 overflow-y-auto space-y-6">
            <!-- Thumbnail Banner Preview -->
            <div v-if="selectedApp.thumbnail" class="w-full h-64 rounded-2xl overflow-hidden border border-slate-800 bg-slate-900 shadow-xl relative">
              <img
                :src="selectedApp.thumbnail"
                alt="App Preview"
                class="w-full h-full object-cover"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-slate-950 via-transparent to-transparent flex items-end p-5">
                <div>
                  <span class="text-xs font-mono text-sky-400 uppercase tracking-widest block font-bold">Katalog Resmi STB</span>
                  <h3 class="text-xl font-bold text-white">{{ selectedApp.name }}</h3>
                </div>
              </div>
            </div>

            <!-- Ringkasan Singkat -->
            <div class="p-4 rounded-2xl bg-slate-900/80 border border-slate-800 text-slate-300 leading-relaxed text-sm">
              <h4 class="font-bold text-slate-100 mb-1">Ikhtisar:</h4>
              <p>{{ selectedApp.description }}</p>
            </div>

            <!-- README.md Rendered Markdown Preview -->
            <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800/80 space-y-4">
              <div class="flex items-center gap-2 border-b border-slate-800 pb-3 text-slate-200 font-bold text-sm">
                <BookOpen class="w-4 h-4 text-sky-400" />
                <span>Dokumentasi Resmi & Panduan (README.md)</span>
              </div>

              <div
                class="prose prose-invert max-w-none text-slate-300 text-xs sm:text-sm leading-relaxed
                       [&_h1]:text-lg [&_h1]:font-bold [&_h1]:text-slate-100 [&_h1]:mb-3
                       [&_h2]:text-base [&_h2]:font-bold [&_h2]:text-slate-200 [&_h2]:mt-4 [&_h2]:mb-2
                       [&_h3]:text-sm [&_h3]:font-bold [&_h3]:text-sky-300 [&_h3]:mt-3 [&_h3]:mb-1.5
                       [&_p]:mb-3 [&_p]:text-slate-300
                       [&_ul]:list-disc [&_ul]:pl-5 [&_ul]:space-y-1.5 [&_ul]:mb-3
                       [&_ol]:list-decimal [&_ol]:pl-5 [&_ol]:space-y-1.5 [&_ol]:mb-3
                       [&_li]:text-slate-300
                       [&_code]:bg-slate-950 [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:rounded [&_code]:text-emerald-400 [&_code]:font-mono [&_code]:text-xs
                       [&_pre]:bg-slate-950 [&_pre]:p-4 [&_pre]:rounded-xl [&_pre]:border [&_pre]:border-slate-800 [&_pre]:overflow-x-auto
                       [&_img]:rounded-xl [&_img]:max-h-72 [&_img]:w-full [&_img]:object-cover [&_img]:my-4 [&_img]:border [&_img]:border-slate-800
                       [&_blockquote]:border-l-4 [&_blockquote]:border-sky-500 [&_blockquote]:pl-3 [&_blockquote]:italic [&_blockquote]:text-slate-400"
                v-html="renderedReadmeHtml"
              ></div>
            </div>
          </div>

          <!-- Right Properties Sidebar -->
          <div class="w-72 border-l border-slate-800/80 bg-slate-900/50 p-5 flex flex-col justify-between hidden lg:flex shrink-0 space-y-4">
            <div class="space-y-4">
              <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Spesifikasi Container</span>

              <div class="space-y-3">
                <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                  <span class="text-[10px] text-slate-500 block">Docker Image:</span>
                  <span class="font-mono text-slate-200 text-xs break-all">{{ selectedApp.image }}</span>
                </div>

                <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                  <span class="text-[10px] text-slate-500 block">Port Akses Bawaan:</span>
                  <span class="font-mono text-sky-400 text-sm font-bold">{{ selectedApp.default_port }}</span>
                </div>

                <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                  <span class="text-[10px] text-slate-500 block">Mountpoint MicroSD:</span>
                  <span class="font-mono text-slate-300 text-xs break-all">{{ selectedApp.default_volume }}</span>
                </div>

                <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                  <span class="text-[10px] text-slate-500 block">Batas Memori RAM:</span>
                  <span class="font-mono text-emerald-400 text-xs font-bold">256 MB (Hemat Daya STB)</span>
                </div>
              </div>
            </div>

            <!-- Tombol Pasang di Bawah Sidebar Detail -->
            <div class="pt-4 border-t border-slate-800">
              <template v-if="!selectedApp.installed">
                <button
                  @click="(e) => openInstallDialog(selectedApp, e)"
                  class="w-full py-2.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center justify-center gap-2 shadow-lg"
                >
                  <Download class="w-4 h-4" />
                  <span>Pasang Aplikasi</span>
                </button>
              </template>
              <template v-else>
                <div class="flex flex-col gap-2">
                  <button
                    @click="launchApp(selectedApp)"
                    class="w-full py-2.5 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs flex items-center justify-center gap-2 shadow-lg"
                  >
                    <ExternalLink class="w-4 h-4" />
                    <span>Buka Web UI</span>
                  </button>
                  <button
                    @click="(e) => openUninstallDialog(selectedApp, e)"
                    class="w-full py-2 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 font-bold text-xs flex items-center justify-center gap-2 border border-rose-500/30"
                  >
                    <Trash2 class="w-4 h-4" />
                    <span>Hapus Aplikasi</span>
                  </button>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ================= CUSTOM MODALS ================= -->

    <!-- CUSTOM MODAL 1: Manual Install Docker Cepat (Floating Window Dialog) -->
    <div
      v-if="showManualModal"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-lg bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <PlusCircle class="w-4 h-4" />
            </div>
            <div>
              <h4 class="font-bold text-slate-100 text-sm">Instal Cepat Docker Container</h4>
              <p class="text-[11px] text-slate-400">Deploy container mandiri dengan parameter cepat.</p>
            </div>
          </div>
          <button @click="showManualModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3.5">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nama Aplikasi / Container</label>
            <input
              v-model="manualName"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: My Custom App"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Docker Image (Docker Hub / Registry)</label>
            <input
              v-model="manualImage"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: nginx:alpine atau redis:7-alpine"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Host Port</label>
              <input
                v-model="manualPort"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
                placeholder="8080"
              />
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Direktori Volume Host</label>
              <input
                v-model="manualVolume"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
                placeholder="/DATA/apps/custom-app"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Environment Variables (Satu per baris)</label>
            <textarea
              v-model="manualEnvs"
              rows="3"
              class="w-full p-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 resize-none"
              placeholder="TZ=Asia/Jakarta&#10;APP_ENV=production"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showManualModal = false"
            class="px-3.5 py-2 rounded-xl bg-slate-800 text-slate-300 font-semibold text-xs hover:bg-slate-700 transition"
          >
            Batal
          </button>
          <button
            @click="submitManualInstall"
            :disabled="isDeployingManual"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow active:scale-95 disabled:opacity-50"
          >
            <span>{{ isDeployingManual ? 'Memasang Container...' : 'Pasang Sekarang' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL 2: Form Konfigurasi Pasang Aplikasi Katalog (Floating Window Dialog) -->
    <div
      v-if="showInstallModal"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <Settings2 class="w-4 h-4" />
            </div>
            <h4 class="font-bold text-slate-100 text-sm">Pasang {{ activeAppToInstall?.name }}</h4>
          </div>
          <button @click="showInstallModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3.5">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Port Akses Host (Container Port)</label>
            <input
              v-model="customPort"
              type="number"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: 8080"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Direktori Penyimpanan Volume</label>
            <input
              v-model="customVolume"
              type="text"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: /DATA/apps/nama-app"
            />
            <span class="text-[10px] text-slate-500 mt-1 block">Data disimpan permanen di MicroSD /DATA/apps</span>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showInstallModal = false"
            class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmInstall"
            :disabled="isInstalling"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow active:scale-95 disabled:opacity-50"
          >
            <span>{{ isInstalling ? 'Sedang Memasang...' : 'Mulai Deployment' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL 3: Konfirmasi Clean Uninstall (Floating Window Dialog) -->
    <div
      v-if="showUninstallModal"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-rose-500/40 shadow-2xl shadow-black ring-1 ring-rose-500/20 rounded-2xl flex flex-col overflow-hidden">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5 text-rose-400">
            <div class="w-7 h-7 rounded-lg bg-rose-500/20 text-rose-400 flex items-center justify-center border border-rose-500/30">
              <AlertTriangle class="w-4 h-4" />
            </div>
            <h4 class="font-bold text-slate-100 text-sm">Konfirmasi Uninstall Aplikasi</h4>
          </div>
          <button @click="showUninstallModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 flex flex-col gap-3.5">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0">
              <AlertTriangle class="w-5 h-5" />
            </div>
            <div>
              <h4 class="font-bold text-slate-100 text-sm">Uninstall {{ activeAppToUninstall?.name }}?</h4>
              <p class="text-xs text-slate-400 mt-0.5">Tindakan ini akan menghentikan dan menghapus container Docker terkait.</p>
            </div>
          </div>

          <label class="flex items-center gap-2.5 p-3 rounded-xl bg-slate-950 border border-slate-800 cursor-pointer">
            <input
              v-model="purgeDataChecked"
              type="checkbox"
              class="w-4 h-4 rounded text-rose-500 focus:ring-0 bg-slate-950 border-slate-700"
            />
            <div class="text-xs">
              <span class="font-medium text-slate-200 block">Musnahkan Folder Data Terkait</span>
              <span class="text-[10px] text-slate-400">Hapus permanen berkas data di /DATA/apps/{{ activeAppToUninstall?.slug }}</span>
            </div>
          </label>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showUninstallModal = false"
            class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmUninstall"
            :disabled="isUninstalling"
            class="px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs transition shadow active:scale-95 disabled:opacity-50"
          >
            <span>{{ isUninstalling ? 'Menghapus...' : 'Ya, Hapus Bersih' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
