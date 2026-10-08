<script setup>
/**
 * RAZOS App Store View - 520+ Apps, Full Page App Details View & CasaOS-Style Custom Manual Install
 * Updated: 2026-10-09
 * Deskripsi:
 * - Logo resmi/identicon & thumbnail banner untuk setiap aplikasi Docker.
 * - Klik kartu aplikasi langsung membuka "Halaman Penuh Detail Aplikasi" (Full Page Detail View, BUKAN modal!).
 * - Halaman Detail merender README.md dalam mode preview Markdown visual kaya (dengan banner image, formatting typography, tag port, volume, dll).
 * - Tombol Pasang / Pasang Ulang / Uninstall langsung di dalam Halaman Detail dan di kartu katalog.
 * - Tombol "Instal Manual Docker" ala CasaOS untuk deployment custom image.
 * - Live Search secepat kilat untuk 525 aplikasi.
 */
import { ref, onMounted, computed } from 'vue'
import { marked } from 'marked'
import { useSystemStore } from '../stores/RAZSystemStore'
import {
  Download, Check, Box, Trash2, Settings2, X, AlertTriangle,
  Search, BookOpen, PlusCircle, ArrowLeft, ExternalLink, HardDrive,
  Cpu, Terminal, Layers, Globe, Shield, Activity, RefreshCw
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const apps = ref([])
const isLoading = ref(false)
const selectedCategory = ref('All')
const searchQuery = ref('')

// State Halaman: 'catalog' | 'detail'
const currentView = ref('catalog')
const selectedApp = ref(null)

// State Modal Pasang Aplikasi (Custom Port & Volume)
const showInstallModal = ref(false)
const activeAppToInstall = ref(null)
const customPort = ref('')
const customVolume = ref('')

// State Modal Custom Manual Install (CasaOS Style)
const showManualModal = ref(false)
const manualName = ref('')
const manualImage = ref('')
const manualPort = ref('8080')
const manualVolume = ref('/DATA/apps/custom-app')
const manualEnvs = ref('')
const isDeployingManual = ref(false)

// State Modal Clean Uninstall
const showUninstallModal = ref(false)
const activeAppToUninstall = ref(null)
const purgeDataChecked = ref(true)

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

async function fetchApps() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/store/apps')
    if (res.ok) {
      apps.value = await res.json()
    }
  } catch (err) {}
  isLoading.value = false
}

// Filtered apps berdasarkan pencarian teks dan kategori
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

// Buka Halaman Penuh Detail Aplikasi (Bukan Modal!)
function openAppDetail(app) {
  selectedApp.value = app
  currentView.value = 'detail'
}

function backToCatalog() {
  currentView.value = 'catalog'
  selectedApp.value = null
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

function openInstallDialog(app, e) {
  e?.stopPropagation()
  activeAppToInstall.value = app
  customPort.value = app.default_port
  customVolume.value = app.default_volume
  showInstallModal.value = true
}

async function confirmInstall() {
  if (!activeAppToInstall.value) return
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
      fetchApps()
    } else {
      const errData = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Pasang', errData.error || 'Gagal memasang aplikasi', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// Submit Manual Custom Docker Container (CasaOS Style)
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
      systemStore.showNotification('Sukses', `Container custom ${manualName.value} berhasil di-deploy!`, 'success')
      showManualModal.value = false
      manualName.value = ''
      manualImage.value = ''
      manualPort.value = '8080'
      manualVolume.value = '/DATA/apps/custom-app'
      manualEnvs.value = ''
      fetchApps()
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Pasang', err.error || 'Gagal memasang custom container', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
  isDeployingManual.value = false
}

function openUninstallDialog(app, e) {
  e?.stopPropagation()
  activeAppToUninstall.value = app
  purgeDataChecked.value = true
  showUninstallModal.value = true
}

async function confirmUninstall() {
  if (!activeAppToUninstall.value) return
  try {
    const res = await fetch(`/api/v1/store/apps/${activeAppToUninstall.value.slug}?purge_data=${purgeDataChecked.value}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Aplikasi ${activeAppToUninstall.value.name} telah di-uninstall!`, 'success')
      showUninstallModal.value = false
      if (selectedApp.value && selectedApp.value.id === activeAppToUninstall.value.id) {
        selectedApp.value.installed = false
      }
      fetchApps()
    } else {
      systemStore.showNotification('Gagal', 'Gagal meng-uninstall container', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

onMounted(() => {
  fetchApps()
})
</script>

<template>
  <div class="h-full flex flex-col text-xs sm:text-sm select-none bg-slate-950 text-slate-100">
    <!-- VIEW 1: KATALOG UTAMA (525+ APPS) -->
    <div v-if="currentView === 'catalog'" class="h-full flex flex-col gap-3 p-4 overflow-hidden animate-in fade-in">
      <!-- Header App Store & Action Buttons -->
      <div class="flex items-center justify-between border-b border-slate-800/80 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-slate-100 text-lg">App Store RAZOS</h3>
            <span class="px-2.5 py-0.5 rounded-full bg-sky-500/20 text-sky-400 text-xs font-mono font-bold">
              {{ apps.length }} Docker Siap Pakai
            </span>
          </div>
          <p class="text-xs text-slate-400 mt-0.5">
            Pilih aplikasi apa pun untuk melihat detail & dokumentasi, atau deploy container kustom.
          </p>
        </div>

        <!-- Action: Manual Install Button (CasaOS Style) -->
        <button
          @click="showManualModal = true"
          class="px-4 py-2 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 font-bold flex items-center gap-2 shadow-lg transition active:scale-95 text-xs"
        >
          <PlusCircle class="w-4 h-4" />
          <span>Instal Manual Docker</span>
        </button>
      </div>

      <!-- Search Bar & Statistics -->
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5">
        <div class="flex-1 relative">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari dari 500+ aplikasi Docker (Plex, Jellyfin, MySQL, Nextcloud, Ollama, Nginx, qBittorrent...)"
            class="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-900/80 border border-slate-800 text-xs text-slate-100 placeholder-slate-500 outline-none focus:border-sky-400 transition"
          />
        </div>

        <span class="text-[11px] text-slate-500 font-mono hidden md:inline shrink-0">
          Menampilkan: {{ filteredApps.length }} aplikasi
        </span>
      </div>

      <!-- Horizontal Category Filter Bar -->
      <div class="flex items-center gap-1.5 overflow-x-auto pb-1 shrink-0 scrollbar-none">
        <button
          v-for="cat in categories"
          :key="cat"
          @click="selectedCategory = cat"
          class="px-3 py-1.5 rounded-xl text-xs font-medium transition whitespace-nowrap"
          :class="selectedCategory === cat ? 'bg-sky-500/20 text-sky-400 border border-sky-400/40 font-bold shadow' : 'text-slate-400 hover:text-slate-200 bg-slate-900/60 border border-slate-800/60'"
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
            <!-- App Logo Image -->
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
                <span class="px-2.5 py-1 rounded-lg bg-emerald-500/20 text-emerald-400 font-bold flex items-center gap-1 text-[11px]">
                  <Check class="w-3.5 h-3.5" />
                  <span>Aktif</span>
                </span>
                <button
                  @click="(e) => openUninstallDialog(app, e)"
                  class="p-1.5 rounded-lg bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 transition"
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

    <!-- VIEW 2: HALAMAN PENUH DETAIL APLIKASI (FULL PAGE DETAIL VIEW) -->
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
            <span class="px-3 py-2 rounded-xl bg-emerald-500/20 text-emerald-400 font-bold flex items-center gap-1.5 text-xs">
              <Check class="w-4 h-4" />
              <span>Container Aktif</span>
            </span>
            <button
              @click="(e) => openUninstallDialog(selectedApp, e)"
              class="px-3.5 py-2 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 font-bold flex items-center gap-1.5 text-xs transition"
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

            <!-- Prosedur Typography Markdown Preview -->
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
            <button
              v-if="!selectedApp.installed"
              @click="(e) => openInstallDialog(selectedApp, e)"
              class="w-full py-2.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center justify-center gap-2 shadow-lg"
            >
              <Download class="w-4 h-4" />
              <span>Pasang Aplikasi</span>
            </button>
            <button
              v-else
              @click="(e) => openUninstallDialog(selectedApp, e)"
              class="w-full py-2.5 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 font-bold text-xs flex items-center justify-center gap-2"
            >
              <Trash2 class="w-4 h-4" />
              <span>Hapus Aplikasi</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL 1: Manual Install Docker (CasaOS Style) -->
    <div
      v-if="showManualModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-lg raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <PlusCircle class="w-5 h-5 text-sky-400" />
            <div>
              <h4 class="font-bold text-slate-100">Instal Manual Docker Container</h4>
              <p class="text-[11px] text-slate-400">Deploy image Docker kustom mana pun sesuai keinginan Anda (ala CasaOS).</p>
            </div>
          </div>
          <button @click="showManualModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="space-y-3">
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

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showManualModal = false"
            class="px-3.5 py-2 rounded-xl bg-slate-800 text-slate-300 font-medium text-xs hover:bg-slate-700"
          >
            Batal
          </button>
          <button
            @click="submitManualInstall"
            :disabled="isDeployingManual"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 shadow"
          >
            <span>{{ isDeployingManual ? 'Memasang Container...' : 'Pasang Manual' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL 2: Form Konfigurasi Pasang Aplikasi Katalog -->
    <div
      v-if="showInstallModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Settings2 class="w-5 h-5 text-sky-400" />
            <h4 class="font-bold text-slate-100">Pasang {{ activeAppToInstall?.name }}</h4>
          </div>
          <button @click="showInstallModal = false" class="p-1 text-slate-400 hover:text-slate-200">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3">
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

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showInstallModal = false"
            class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmInstall"
            class="px-4 py-2 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow"
          >
            Mulai Deployment
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL 3: Konfirmasi Clean Uninstall -->
    <div
      v-if="showUninstallModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-rose-500/30 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0">
            <AlertTriangle class="w-5 h-5" />
          </div>
          <div>
            <h4 class="font-bold text-slate-100">Hapus {{ activeAppToUninstall?.name }}?</h4>
            <p class="text-xs text-slate-400 mt-0.5">Tindakan ini akan menghentikan dan menghapus container Docker terkait.</p>
          </div>
        </div>

        <label class="flex items-center gap-2 p-3 rounded-xl bg-slate-900/80 border border-slate-800 cursor-pointer">
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

        <div class="flex items-center justify-end gap-2 pt-2">
          <button
            @click="showUninstallModal = false"
            class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmUninstall"
            class="px-4 py-2 rounded-lg bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs transition shadow"
          >
            Ya, Hapus Bersih
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
