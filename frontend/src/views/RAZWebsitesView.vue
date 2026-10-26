<script setup>
/**
 * RAZOS Website Manager Suite
 * Updated: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Platform pengelolaan multi-website terisolasi berbasis Docker Container.
 * Fitur:
 * 1. Root folder terisolasi di /DATA/apps/websites/<slug>
 * 2. Port host independen (8100-8250) anti-tabrakan
 * 3. Dukungan runtime stack: HTML Statis (Nginx), PHP 8.2 + DB, Node.js, Next.js, Custom
 * 4. Integrasi Database Manager (Adminer / HeidiSQL web-style) dengan auto-fill koneksi
 * 5. Link icon kustom untuk launcher Desktop dan Start Menu
 * 6. Resource limits: Memory limit (RAM STB-optimized) & CPU shares
 * 7. Custom container command & Environment variables
 * 8. Standalone Desktop Windows: Seluruh form dialog (Wizard, Edit, Git Sync, Rincian)
 *    berjalan di jendela desktop mandiri tanpa mengunci/memblokir Website Manager.
 */
import { ref, onMounted, computed, watch } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Globe, Plus, ExternalLink, Inbox, Trash2, X, Server, Layers,
  Play, Square, RotateCw, FolderOpen, Database, Cpu, HardDrive,
  Copy, Check, ShieldCheck, Search, Filter, Terminal, Sparkles,
  Info, AlertTriangle, ArrowRight, CheckCircle2,
  GitBranch, GitPullRequest, Pencil
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const websites = ref([])
const isLoading = ref(false)
const searchQuery = ref('')
const selectedRuntimeFilter = ref('all')
const pullingSiteId = ref(null)

// State Modal Konfirmasi Hapus Website
const showDeleteModal = ref(false)
const siteToDelete = ref(null)
const purgeRootFolder = ref(false)

// Ambil daftar website dari backend
async function fetchWebsites() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/websites')
    if (res.ok) {
      websites.value = await res.json()
    }
  } catch (err) {
    console.error('Gagal mengambil daftar website:', err)
  }
  isLoading.value = false
}

// Watch trigger sinkronisasi global dari jendela lain (misal setelah buat/edit website)
watch(() => systemStore.websitesRefreshTrigger, () => {
  fetchWebsites()
})

// Toggle Status Container Website (Start / Stop / Restart)
async function toggleWebsiteStatus(site, action) {
  try {
    const res = await fetch(`/api/v1/websites/${site.id}/action`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action })
    })
    if (res.ok) {
      const updated = await res.json()
      site.status = updated.status
      systemStore.showNotification(
        'Status Diperbarui',
        `Website "${site.name}" sekarang berstatus ${site.status.toUpperCase()}`,
        'success'
      )
    }
  } catch (e) {
    systemStore.showNotification('Gagal', `Gagal mengubah status website: ${e.message}`, 'error')
  }
}

// Quick Git Pull langsung dari kartu
async function executeGitPull(site) {
  if (!site) return
  if (!site.git_repo) {
    // Jika belum konfigurasi git repo, buka jendela pengaturan git
    windowStore.openWebsiteGit(site)
    return
  }

  pullingSiteId.value = site.id
  try {
    const res = await fetch(`/api/v1/websites/${site.id}/git/pull`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    if (res.ok) {
      site.git_last_pull = data.website?.git_last_pull || new Date().toISOString()
      site.git_last_commit = data.website?.git_last_commit || ''
      systemStore.showNotification(
        'Git Pull Berhasil',
        `Pembaruan kode GitHub terbaru berhasil ditarik ke "${site.name}"`,
        'success'
      )
      fetchWebsites()
    } else {
      systemStore.showNotification(
        'Git Pull Gagal',
        data.error || 'Terjadi kesalahan saat menarik kode dari GitHub',
        'error'
      )
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi backend: ' + e.message, 'error')
  } finally {
    pullingSiteId.value = null
  }
}

// Buka folder website di File Manager
function openWebsiteFolder(site) {
  windowStore.openApp('file-manager')
  systemStore.showNotification('Root Folder', `Lokasi folder: ${site.root_folder}`, 'info')
}

// Buka Web Adminer Database GUI
async function openDatabaseManager(site) {
  try {
    await fetch('/api/v1/databases/adminer/launch', { method: 'POST' })
    const host = window.location.hostname || 'localhost'
    const adminerUrl = `http://${host}:8088/?server=${encodeURIComponent(site.db_host || '127.0.0.1')}&username=${encodeURIComponent(site.db_user || 'root')}&db=${encodeURIComponent(site.db_name || '')}`
    window.open(adminerUrl, '_blank')
  } catch (e) {
    window.open(`http://${window.location.hostname || 'localhost'}:8088`, '_blank')
  }
}

// Dialog konfirmasi hapus website
const currentHostname = computed(() => window.location.hostname || 'localhost')
function promptDeleteWebsite(site) {
  siteToDelete.value = site
  purgeRootFolder.value = false
  showDeleteModal.value = true
}

async function executeDeleteWebsite() {
  if (!siteToDelete.value) return
  const id = siteToDelete.value.id
  const name = siteToDelete.value.name
  try {
    const res = await fetch(`/api/v1/websites/${id}?purge=${purgeRootFolder.value}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Website "${name}" telah dihapus`, 'success')
      showDeleteModal.value = false
      siteToDelete.value = null
      fetchWebsites()
    } else {
      systemStore.showNotification('Gagal', 'Gagal menghapus website', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// Filter pencarian dan runtime stack
const filteredWebsites = computed(() => {
  return websites.value.filter(site => {
    const q = searchQuery.value.toLowerCase().trim()
    const matchesSearch = !q ||
      (site.name && site.name.toLowerCase().includes(q)) ||
      (site.domain && site.domain.toLowerCase().includes(q)) ||
      (site.port && site.port.toString().includes(q)) ||
      (site.container_name && site.container_name.toLowerCase().includes(q)) ||
      (site.runtime && site.runtime.toLowerCase().includes(q))

    const matchesRuntime = selectedRuntimeFilter.value === 'all' ||
      site.runtime === selectedRuntimeFilter.value

    return matchesSearch && matchesRuntime
  })
})

function getRuntimeBadge(runtime) {
  switch (runtime) {
    case 'static': return { label: 'HTML Statis (Nginx)', color: 'text-sky-300 bg-sky-500/10 border-sky-500/20' }
    case 'php': return { label: 'PHP Native + DB', color: 'text-indigo-300 bg-indigo-500/10 border-indigo-500/20' }
    case 'nodejs': return { label: 'Node.js Express', color: 'text-emerald-300 bg-emerald-500/10 border-emerald-500/20' }
    case 'golang': return { label: 'Go API Micro', color: 'text-cyan-300 bg-cyan-500/10 border-cyan-500/20' }
    case 'python': return { label: 'Python FastAPI', color: 'text-yellow-300 bg-yellow-500/10 border-yellow-500/20' }
    case 'flatcms': return { label: 'Flat CMS (Grav)', color: 'text-teal-300 bg-teal-500/10 border-teal-500/20' }
    case 'wordpress': return { label: 'WordPress Light', color: 'text-blue-300 bg-blue-500/10 border-blue-500/20' }
    case 'dashboard': return { label: 'Dashboard Portal', color: 'text-purple-300 bg-purple-500/10 border-purple-500/20' }
    case 'nextjs': return { label: 'Next.js (Static)', color: 'text-slate-200 bg-slate-500/15 border-slate-500/25' }
    default: return { label: (runtime || 'Custom').toUpperCase(), color: 'text-slate-300 bg-slate-500/10 border-slate-500/20' }
  }
}

onMounted(() => {
  fetchWebsites()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm select-none">
    <!-- Top Action Toolbar -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2.5">
      <div>
        <div class="flex items-center gap-2">
          <h3 class="font-bold text-slate-100 text-base">Website Manager</h3>
          <span class="text-[10px] px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-300 border border-sky-500/30 font-semibold font-mono">
            {{ websites.length }} Web Terpasang
          </span>
        </div>
        <p class="text-xs text-slate-400 mt-0.5">
          Deploy dan kelola website container mandiri dengan isolasi root folder & port host bebas tabrakan.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="fetchWebsites"
          :disabled="isLoading"
          class="p-2 rounded-xl bg-slate-900 hover:bg-slate-800 border border-slate-800 text-slate-300 hover:text-white transition"
          title="Muat Ulang Daftar Website"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
        </button>

        <button
          @click="windowStore.openWebsiteCreator()"
          class="px-3.5 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow-lg shadow-sky-500/20 active:scale-95"
        >
          <Plus class="w-4 h-4 stroke-[2.5]" />
          <span>Buat Website Baru</span>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-2.5 bg-slate-900/60 p-2.5 rounded-xl border border-slate-800/80">
      <div class="relative flex-1 min-w-[220px]">
        <Search class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama website, port, domain, container..."
          class="w-full pl-8 pr-3 py-1.5 bg-slate-950/80 border border-slate-800 rounded-lg text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500 transition"
        />
      </div>

      <!-- Filter Kategori Runtime STB -->
      <div class="flex items-center gap-1.5 overflow-x-auto text-[11px] pb-0.5 lg:pb-0">
        <button
          v-for="f in [
            { id: 'all', label: 'Semua' },
            { id: 'static', label: 'HTML Statis' },
            { id: 'golang', label: 'Go API' },
            { id: 'php', label: 'PHP Native' },
            { id: 'nodejs', label: 'Node.js' },
            { id: 'python', label: 'Python' },
            { id: 'flatcms', label: 'Flat CMS' },
            { id: 'wordpress', label: 'WordPress' },
            { id: 'dashboard', label: 'Dashboard' }
          ]"
          :key="f.id"
          @click="selectedRuntimeFilter = f.id"
          class="px-2.5 py-1 rounded-lg font-medium transition shrink-0 whitespace-nowrap"
          :class="selectedRuntimeFilter === f.id
            ? 'bg-sky-500/20 text-sky-300 border border-sky-500/40 font-bold shadow-sm'
            : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'"
        >
          {{ f.label }}
        </button>
      </div>
    </div>

    <!-- Website Cards Grid -->
    <div
      v-if="filteredWebsites.length > 0"
      class="flex-1 overflow-y-auto grid gap-3.5 pr-0.5"
      style="grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));"
    >
      <div
        v-for="site in filteredWebsites"
        :key="site.id"
        class="raz-glass rounded-2xl border border-slate-800/90 p-4 sm:p-5 flex flex-col justify-between hover:border-sky-500/40 hover:shadow-2xl hover:shadow-sky-500/5 transition-all duration-200 group"
      >
        <!-- Card Top Section -->
        <div>
          <!-- Header Card: Icon, Name, Domain, Status Badge -->
          <div class="flex items-start gap-3.5 mb-3">
            <div class="w-12 h-12 rounded-xl bg-slate-900/90 border border-slate-700/80 p-2.5 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform shadow">
              <img
                v-if="site.icon"
                :src="site.icon"
                :alt="site.name"
                class="w-full h-full object-contain"
                @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'"
              />
              <Globe v-else class="w-6 h-6 text-sky-400" />
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-2">
                <h4 class="font-bold text-slate-100 text-sm sm:text-base truncate group-hover:text-sky-300 transition-colors leading-tight" :title="site.name">
                  {{ site.name }}
                </h4>

                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider flex items-center gap-1.5 border shrink-0"
                  :class="site.status === 'running'
                    ? 'bg-emerald-500/15 text-emerald-300 border-emerald-500/30'
                    : 'bg-rose-500/15 text-rose-300 border-rose-500/30'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="site.status === 'running' ? 'bg-emerald-400 animate-pulse' : 'bg-rose-400'"></span>
                  {{ site.status }}
                </span>
              </div>

              <!-- Host Link & Runtime Tag -->
              <div class="flex items-center gap-2 mt-1">
                <a
                  :href="'http://' + currentHostname + ':' + site.port"
                  target="_blank"
                  class="inline-flex items-center gap-1 font-mono text-[11px] text-sky-400 hover:text-sky-300 hover:underline truncate max-w-[210px]"
                  :title="'Buka http://' + currentHostname + ':' + site.port"
                >
                  <span class="truncate">{{ currentHostname }}:{{ site.port }}</span>
                  <ExternalLink class="w-3 h-3 shrink-0" />
                </a>

                <span
                  v-if="site.domain"
                  class="text-[10px] px-1.5 py-0.5 rounded bg-slate-800/80 text-slate-400 border border-slate-700/60 font-mono truncate max-w-[110px]"
                  :title="'Domain reverse proxy: ' + site.domain"
                >
                  {{ site.domain }}
                </span>

                <span
                  class="text-[10px] px-2 py-0.5 rounded-md border font-medium truncate"
                  :class="getRuntimeBadge(site.runtime).color"
                >
                  {{ getRuntimeBadge(site.runtime).label }}
                </span>
              </div>
            </div>
          </div>

          <!-- Specs Info Box -->
          <div class="p-3 rounded-xl bg-slate-950/70 border border-slate-900/90 text-xs space-y-1.5">
            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-400 shrink-0">Container Name:</span>
              <span class="font-mono text-slate-200 font-semibold truncate text-right flex-1" :title="site.container_name">
                {{ site.container_name }}
              </span>
            </div>

            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-400 shrink-0">Alokasi Port:</span>
              <span class="font-mono text-sky-400 font-bold truncate text-right flex-1">
                :{{ site.port }} <span class="text-[10px] text-slate-500 font-normal">(Host Binding)</span>
              </span>
            </div>

            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-400 shrink-0">Root Folder:</span>
              <span class="font-mono text-slate-300 text-[11px] truncate text-right flex-1" :title="site.root_folder">
                {{ site.root_folder }}
              </span>
            </div>

            <div class="flex items-center justify-between gap-2.5 pt-1 border-t border-slate-900">
              <span class="text-slate-400 shrink-0">Limit RAM:</span>
              <span class="font-mono text-emerald-400 font-semibold truncate text-right flex-1">
                {{ site.memory_limit || '128M' }} ({{ site.cpu_shares || 512 }} CPU)
              </span>
            </div>

            <!-- Database Info (Jika Ada) -->
            <div v-if="site.db_name" class="flex items-center justify-between gap-2.5 pt-1 border-t border-slate-900">
              <span class="text-amber-400/90 flex items-center gap-1.5 shrink-0 whitespace-nowrap font-medium">
                <Database class="w-3.5 h-3.5 text-amber-400" />
                <span>Basis Data:</span>
              </span>
              <span class="font-mono text-amber-300 font-bold truncate text-right flex-1" :title="site.db_name">
                {{ site.db_name }}
              </span>
            </div>

            <!-- Git Info (Jika Ada) -->
            <div v-if="site.git_repo" class="flex items-center justify-between gap-2.5 pt-1 border-t border-slate-900">
              <span class="text-sky-400/90 flex items-center gap-1.5 shrink-0 whitespace-nowrap font-medium">
                <GitBranch class="w-3.5 h-3.5 text-sky-400" />
                <span>GitHub:</span>
              </span>
              <span class="font-mono text-slate-300 text-[10px] truncate text-right flex-1" :title="site.git_last_commit || site.git_repo">
                {{ site.git_branch || 'main' }} • {{ site.git_last_commit ? site.git_last_commit.substring(0, 7) : 'Terkoneksi' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Card Bottom Action Bar -->
        <div class="flex items-center justify-between gap-2 pt-3 mt-3 border-t border-slate-800/80">
          <!-- Left Actions Toolbar: Start/Stop, Folder, Edit, Git, Database, Info -->
          <div class="flex items-center gap-1 bg-slate-950/60 p-1 rounded-xl border border-slate-900">
            <!-- Toggle Start / Stop -->
            <button
              v-if="site.status === 'running'"
              @click="toggleWebsiteStatus(site, 'stop')"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-amber-400 transition"
              title="Hentikan Container Website"
            >
              <Square class="w-3.5 h-3.5 fill-amber-400/20" />
            </button>
            <button
              v-else
              @click="toggleWebsiteStatus(site, 'start')"
              class="p-1.5 rounded-lg bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 transition"
              title="Jalankan Container Website"
            >
              <Play class="w-3.5 h-3.5 fill-emerald-400" />
            </button>

            <!-- Buka Root Folder di File Manager -->
            <button
              @click="openWebsiteFolder(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white transition"
              title="Buka Root Folder di File Manager RAZOS"
            >
              <FolderOpen class="w-3.5 h-3.5" />
            </button>

            <!-- Edit Konfigurasi Website (Buka Jendela Mandiri) -->
            <button
              @click="windowStore.openWebsiteEditor(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-amber-300 transition"
              title="Edit Konfigurasi Website (Jendela Mandiri)"
            >
              <Pencil class="w-3.5 h-3.5" />
            </button>

            <!-- Quick Git Pull -->
            <button
              @click="executeGitPull(site)"
              :disabled="pullingSiteId === site.id"
              class="p-1.5 rounded-lg transition"
              :class="site.git_repo
                ? 'bg-sky-500/15 hover:bg-sky-500/25 text-sky-300'
                : 'bg-slate-900 hover:bg-slate-800 text-slate-400 hover:text-slate-200'"
              :title="site.git_repo ? 'Tarik Pembaruan GitHub Terbaru (Git Pull)' : 'Hubungkan Repository GitHub'"
            >
              <RotateCw v-if="pullingSiteId === site.id" class="w-3.5 h-3.5 animate-spin text-sky-400" />
              <GitPullRequest v-else class="w-3.5 h-3.5" />
            </button>

            <!-- Pengaturan Git & PAT Token (Buka Jendela Mandiri) -->
            <button
              @click="windowStore.openWebsiteGit(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-sky-300 transition"
              title="Pengaturan GitHub & PAT Token (Jendela Mandiri)"
            >
              <GitBranch class="w-3.5 h-3.5" />
            </button>

            <!-- Buka Web Database GUI (Adminer) jika ada database -->
            <button
              v-if="site.db_name"
              @click="openDatabaseManager(site)"
              class="p-1.5 rounded-lg bg-amber-500/20 hover:bg-amber-500/30 text-amber-300 transition"
              title="Kelola Database di Web Adminer"
            >
              <Database class="w-3.5 h-3.5" />
            </button>

            <!-- Detail Konfigurasi (Buka Jendela Mandiri) -->
            <button
              @click="windowStore.openWebsiteDetail(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 transition"
              title="Lihat Detail & Kredensial (Jendela Mandiri)"
            >
              <Info class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Right Actions: Launch Web Button & Delete -->
          <div class="flex items-center gap-1.5">
            <a
              :href="'http://' + currentHostname + ':' + site.port"
              target="_blank"
              class="px-3.5 py-1.5 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow-md shadow-sky-500/20 active:scale-95"
              :title="'Buka website di http://' + currentHostname + ':' + site.port"
            >
              <span>Buka</span>
              <ExternalLink class="w-3.5 h-3.5 stroke-[2.5]" />
            </a>

            <button
              @click="promptDeleteWebsite(site)"
              class="p-1.5 rounded-xl hover:bg-rose-500/20 text-rose-400 hover:text-rose-300 transition"
              title="Hapus Website"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else
      class="flex-1 rounded-2xl bg-slate-900/40 border border-slate-800/60 flex flex-col items-center justify-center text-slate-500 gap-3 p-8 text-center"
    >
      <div class="w-14 h-14 rounded-2xl bg-sky-500/10 border border-sky-500/20 flex items-center justify-center text-sky-400 shadow-xl shadow-sky-500/5">
        <Globe class="w-7 h-7" />
      </div>
      <div class="max-w-md">
        <h4 class="font-bold text-slate-200 text-sm sm:text-base">Belum Ada Website yang Dikonfigurasi</h4>
        <p class="text-xs text-slate-400 mt-1 leading-relaxed">
          Gunakan Website Manager untuk membuat container website terisolasi (HTML Statis, PHP dengan Database, Node.js, atau Next.js) dengan root folder dan alokasi port bebas.
        </p>
      </div>
      <button
        @click="windowStore.openWebsiteCreator()"
        class="mt-2 px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 transition shadow-lg shadow-sky-500/20 active:scale-95"
      >
        <Plus class="w-4 h-4 stroke-[2.5]" />
        <span>Buat Website Pertama Anda</span>
      </button>
    </div>

    <!-- CUSTOM MODAL: Konfirmasi Hapus Website (Floating Window Dialog Ringan) -->
    <div
      v-if="showDeleteModal && siteToDelete"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-rose-500/40 shadow-2xl shadow-black ring-1 ring-rose-500/20 rounded-2xl flex flex-col overflow-hidden">
        <!-- Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2 text-rose-400">
            <AlertTriangle class="w-4 h-4" />
            <h4 class="font-bold text-slate-100 text-sm">Konfirmasi Hapus Website</h4>
          </div>
          <button @click="showDeleteModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3.5">
          <p class="text-xs text-slate-300">
            Apakah Anda yakin ingin menghapus container website berikut dari sistem?
          </p>

          <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 space-y-1 text-xs">
            <div class="text-slate-300 font-semibold">{{ siteToDelete.name }}</div>
            <div class="font-mono text-slate-400 text-[11px]">Container: {{ siteToDelete.container_name }}</div>
            <div class="font-mono text-amber-400 text-[11px]">Port: :{{ siteToDelete.port }}</div>
          </div>

          <label class="flex items-center gap-2.5 p-2 rounded-lg bg-rose-500/10 border border-rose-500/20 text-xs text-rose-300 cursor-pointer">
            <input type="checkbox" v-model="purgeRootFolder" class="w-4 h-4 accent-rose-500 rounded" />
            <span>Hapus juga seluruh berkas di root folder (`{{ siteToDelete.root_folder }}`)</span>
          </label>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showDeleteModal = false"
            class="px-3.5 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="executeDeleteWebsite"
            class="px-4 py-1.5 rounded-xl bg-rose-500 hover:bg-rose-400 text-white text-xs font-bold transition shadow active:scale-95 flex items-center gap-1.5"
          >
            <Trash2 class="w-3.5 h-3.5" />
            <span>Hapus Permanen</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
