<script setup>
/**
 * RAZOS App Store View - 500+ Apps & CasaOS-Style Custom Manual Install
 * Updated: 2026-10-09
 * Deskripsi: Antarmuka App Store lengkap dengan katalog 520+ kontainer Docker siap pakai:
 * - Details & Readme.md Viewer untuk setiap aplikasi.
 * - Tombol Pasang 1-klik dengan custom port & volume bindings.
 * - Tombol "Instal Manual" (Custom Docker App) ala CasaOS (Image, Port, Web UI, Volume, Env).
 * - Kategori komprehensif (Media, Cloud, Database, Dev, IoT, Networking, Monitoring, dll).
 * - Live Search secepat kilat untuk 500+ aplikasi.
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import {
  Download, Check, Box, Trash2, Settings2, X, AlertTriangle,
  Search, BookOpen, PlusCircle, ExternalLink, HardDrive, Cpu, Terminal
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const apps = ref([])
const isLoading = ref(false)
const selectedCategory = ref('All')
const searchQuery = ref('')

// State Modal Pasang Aplikasi (Custom Port & Volume)
const showInstallModal = ref(false)
const activeAppToInstall = ref(null)
const customPort = ref('')
const customVolume = ref('')

// State Modal Details & Readme.md
const showDetailModal = ref(false)
const activeAppDetail = ref(null)

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

function openDetailDialog(app) {
  activeAppDetail.value = app
  showDetailModal.value = true
}

function openInstallDialog(app) {
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
      if (showDetailModal.value) showDetailModal.value = false
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

function openUninstallDialog(app) {
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
      if (showDetailModal.value) showDetailModal.value = false
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
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm select-none">
    <!-- Header App Store & Action Buttons -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2.5">
      <div>
        <div class="flex items-center gap-2">
          <h3 class="font-bold text-slate-100 text-base">App Store RAZOS</h3>
          <span class="px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-400 text-[11px] font-mono font-bold">
            {{ apps.length }} Siap Pakai
          </span>
        </div>
        <p class="text-xs text-slate-400 mt-0.5">Katalog Docker siap pakai 1-klik terisolasi di MicroSD (/DATA/apps) & Custom Install.</p>
      </div>

      <!-- Action: Manual Install Button (CasaOS Style) -->
      <button
        @click="showManualModal = true"
        class="px-3.5 py-2 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 font-bold flex items-center gap-2 shadow-lg transition active:scale-95 text-xs"
      >
        <PlusCircle class="w-4 h-4" />
        <span>Instal Manual Docker</span>
      </button>
    </div>

    <!-- Search Bar & Category Navigation -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5">
      <!-- Search Input -->
      <div class="flex-1 relative">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari dari 500+ aplikasi Docker (Plex, Jellyfin, MySQL, Nextcloud, Ollama, Nginx...)"
          class="w-full pl-9 pr-3 py-2 rounded-xl bg-slate-900/80 border border-slate-800 text-xs text-slate-100 placeholder-slate-500 outline-none focus:border-sky-400 transition"
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
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 overflow-auto flex-1 pr-1">
      <div
        v-for="app in filteredApps"
        :key="app.id"
        class="p-3.5 rounded-2xl bg-slate-900/70 border border-slate-800/80 hover:border-slate-700/80 flex flex-col justify-between gap-3 transition group relative"
      >
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-xl bg-sky-500/10 border border-sky-500/30 flex items-center justify-center text-sky-400 shrink-0">
            <Box class="w-5 h-5" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-1">
              <h4 class="font-bold text-slate-200 text-xs truncate" :title="app.name">{{ app.name }}</h4>
              <span class="text-[9px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 shrink-0 font-medium">
                {{ app.category }}
              </span>
            </div>
            <p class="text-[11px] text-slate-400 mt-1 line-clamp-2 leading-relaxed">{{ app.description }}</p>
          </div>
        </div>

        <div class="flex items-center justify-between border-t border-slate-800/60 pt-2 text-[11px]">
          <div class="text-slate-500 font-mono flex items-center gap-1.5">
            <span>Port:</span>
            <span class="text-sky-400 font-bold">{{ app.default_port }}</span>
          </div>

          <div class="flex items-center gap-1.5">
            <!-- Tombol Details / Readme.md -->
            <button
              @click="openDetailDialog(app)"
              class="px-2.5 py-1.5 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 font-medium flex items-center gap-1 transition"
              title="Lihat Details / README.md"
            >
              <BookOpen class="w-3.5 h-3.5 text-slate-400" />
              <span>Detail</span>
            </button>

            <!-- Jika Terpasang -->
            <template v-if="app.installed">
              <span class="px-2 py-1 rounded-md bg-emerald-500/20 text-emerald-400 font-bold flex items-center gap-1">
                <Check class="w-3.5 h-3.5" />
                <span>Aktif</span>
              </span>
              <button
                @click="openUninstallDialog(app)"
                class="p-1.5 rounded-md bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 transition"
                title="Uninstall Bersih"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </template>

            <!-- Jika Belum Terpasang -->
            <button
              v-else
              @click="openInstallDialog(app)"
              class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center gap-1 transition shadow active:scale-95"
            >
              <Download class="w-3.5 h-3.5" />
              <span>Pasang</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 1. CUSTOM MODAL: Details / README.md Viewer -->
    <div
      v-if="showDetailModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-2xl max-h-[85vh] raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-sky-500/20 text-sky-400 flex items-center justify-center font-bold">
              <Box class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-bold text-slate-100 text-base">{{ activeAppDetail?.name }}</h3>
              <span class="text-xs text-sky-400 font-mono">{{ activeAppDetail?.image }}</span>
            </div>
          </div>
          <button @click="showDetailModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Konten Markdown README.md Terformat -->
        <div class="flex-1 overflow-y-auto p-4 rounded-xl bg-slate-950/80 border border-slate-800 text-slate-300 text-xs leading-relaxed space-y-3 font-sans">
          <div class="p-3 bg-sky-500/10 border border-sky-500/30 rounded-xl text-sky-200">
            <h5 class="font-bold mb-1">Deskripsi Singkat:</h5>
            <p>{{ activeAppDetail?.description }}</p>
          </div>

          <div class="space-y-2">
            <h5 class="font-bold text-slate-200 text-sm border-b border-slate-800 pb-1">Dokumentasi & Spesifikasi README.md:</h5>
            <pre class="font-mono text-emerald-400 whitespace-pre-wrap bg-black/60 p-3 rounded-lg border border-slate-800/80">{{ activeAppDetail?.readme || 'Tidak ada file README.md spesifik.' }}</pre>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-2">
            <div class="p-2.5 rounded-lg bg-slate-900 border border-slate-800">
              <span class="text-[10px] text-slate-500 block">Port Akses:</span>
              <span class="font-mono text-slate-200 font-bold">{{ activeAppDetail?.default_port }}</span>
            </div>
            <div class="p-2.5 rounded-lg bg-slate-900 border border-slate-800">
              <span class="text-[10px] text-slate-500 block">Kategori Aplikasi:</span>
              <span class="text-slate-200 font-medium">{{ activeAppDetail?.category }}</span>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-slate-800">
          <button
            @click="showDetailModal = false"
            class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs"
          >
            Tutup
          </button>
          <button
            v-if="!activeAppDetail?.installed"
            @click="() => { showDetailModal = false; openInstallDialog(activeAppDetail); }"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 shadow"
          >
            <Download class="w-4 h-4" />
            <span>Pasang Sekarang</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 2. CUSTOM MODAL: Manual Install Docker (CasaOS Style) -->
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

    <!-- 3. CUSTOM MODAL: Form Konfigurasi Pasang Aplikasi Katalog -->
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

    <!-- 4. CUSTOM MODAL: Konfirmasi Clean Uninstall -->
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
