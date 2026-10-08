<script setup>
/**
 * RAZOS App Store View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka katalog aplikasi Docker 1-klik untuk di-deploy ke direktori MicroSD /DATA/apps/.
 * Dilengkapi dialog modal kustom untuk form custom port/volume dan konfirmasi Clean Uninstall.
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Download, Check, Box, Trash2, Settings2, X, AlertTriangle } from 'lucide-vue-next'

const systemStore = useSystemStore()
const apps = ref([])
const isLoading = ref(false)
const selectedCategory = ref('All')

// State Modal Pasang Aplikasi (Custom Port & Volume)
const showInstallModal = ref(false)
const activeAppToInstall = ref(null)
const customPort = ref('')
const customVolume = ref('')

// State Modal Uninstall Bersih (Clean Uninstall)
const showUninstallModal = ref(false)
const activeAppToUninstall = ref(null)
const purgeDataChecked = ref(true)

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
      fetchApps()
    } else {
      const errData = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Pasang', errData.error || 'Gagal memasang aplikasi', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
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
      systemStore.showNotification('Sukses', `Aplikasi ${activeAppToUninstall.value.name} telah di-uninstall bersih!`, 'success')
      showUninstallModal.value = false
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
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <!-- Header App Store -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">App Store RAZOS</h3>
        <p class="text-xs text-slate-400">Instalasi aplikasi terisolasi ke MicroSD (/DATA/apps) dengan limit RAM protektif.</p>
      </div>
      <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ" class="h-5 object-contain" />
    </div>

    <!-- Category Filter Bar -->
    <div class="flex items-center gap-1.5 overflow-x-auto pb-1">
      <button
        v-for="cat in ['All', 'Network', 'Database', 'Monitoring', 'Storage']"
        :key="cat"
        @click="selectedCategory = cat"
        class="px-3 py-1 rounded-lg text-xs font-medium transition"
        :class="selectedCategory === cat ? 'bg-sky-500/20 text-sky-400 border border-sky-400/40' : 'text-slate-400 hover:text-slate-200 bg-slate-900/60'"
      >
        {{ cat }}
      </button>
    </div>

    <!-- App Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3 overflow-auto flex-1 pr-1">
      <div
        v-for="app in apps.filter(a => selectedCategory === 'All' || a.category === selectedCategory)"
        :key="app.id"
        class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col justify-between gap-3 hover:border-slate-700 transition"
      >
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-xl bg-sky-500/10 border border-sky-500/30 flex items-center justify-center text-sky-400 shrink-0">
            <Box class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-semibold text-slate-200 text-sm">{{ app.name }}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-400">{{ app.category }}</span>
            </div>
            <p class="text-xs text-slate-400 mt-1 leading-relaxed">{{ app.description }}</p>
          </div>
        </div>

        <div class="flex items-center justify-between border-t border-slate-800/60 pt-2 text-xs">
          <div class="text-slate-400">
            Port bawaan: <span class="font-mono text-sky-400">{{ app.default_port }}</span>
          </div>

          <div class="flex items-center gap-2">
            <!-- Jika Terpasang -->
            <template v-if="app.installed">
              <span class="px-2 py-1 rounded-md bg-emerald-500/20 text-emerald-400 text-[11px] font-medium flex items-center gap-1">
                <Check class="w-3.5 h-3.5" />
                <span>Aktif</span>
              </span>
              <button
                @click="openUninstallDialog(app)"
                class="px-2.5 py-1 rounded-md bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 text-[11px] font-medium flex items-center gap-1 transition"
                title="Uninstall Bersih"
              >
                <Trash2 class="w-3 h-3" />
                <span>Hapus</span>
              </button>
            </template>

            <!-- Jika Belum Terpasang -->
            <button
              v-else
              @click="openInstallDialog(app)"
              class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow"
            >
              <Download class="w-3.5 h-3.5" />
              <span>Pasang</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Form Konfigurasi Pasang Aplikasi -->
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
            <span class="text-[10px] text-slate-500 mt-1 block">Data akan disimpan permanen di MicroSD /DATA/apps</span>
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

    <!-- CUSTOM MODAL: Konfirmasi Clean Uninstall -->
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

        <!-- Opsi Bersihkan Folder Data -->
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
