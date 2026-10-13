<script setup>
/**
 * RAZOS App Store Catalog Installer View (Standalone Desktop Window)
 * Created: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Form Konfigurasi Pasang Aplikasi Katalog Docker yang berjalan di Jendela OS mandiri.
 * Tidak mengunci atau memblokir jendela utama App Store.
 */
import { ref, computed, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { Download, Check, X, Box, HardDrive, Server } from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ customData: {} })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const app = computed(() => props.win?.customData || {})

const customPort = ref('')
const customVolume = ref('')
const isInstalling = ref(false)

onMounted(() => {
  if (app.value) {
    customPort.value = app.value.default_port || '8080'
    customVolume.value = app.value.default_volume || `/DATA/apps/${app.value.slug || 'app'}`
  }
})

async function handleInstall() {
  if (!app.value || !app.value.id) return
  isInstalling.value = true
  const payload = {
    app_id: app.value.id,
    port: parseInt(customPort.value) || app.value.default_port,
    data_path: customVolume.value
  }

  try {
    const res = await fetch('/api/v1/store/install', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Container ${app.value.name} berhasil dipasang!`, 'success')
      // Picu refresh global ke App Store dan Desktop
      systemStore.triggerAppsRefresh()
      await windowStore.fetchInstalledApps()
      // Tutup jendela instalasi mandiri ini
      windowStore.closeWindow(props.win?.id || `appstore-install-${app.value.slug}`)
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

function handleCancel() {
  windowStore.closeWindow(props.win?.id || `appstore-install-${app.value?.slug}`)
}
</script>

<template>
  <form @submit.prevent="handleInstall" class="h-full flex flex-col justify-between p-4 bg-slate-900/90 text-xs sm:text-sm select-none gap-3">
    <div class="space-y-4 overflow-y-auto pr-1">
      <!-- App Preview Banner -->
      <div class="flex items-center gap-3.5 p-3.5 bg-slate-950/70 rounded-2xl border border-slate-800">
        <img
          :src="app.icon"
          :alt="app.name"
          class="w-12 h-12 object-contain rounded-xl p-1.5 bg-slate-900 border border-slate-800 shrink-0"
        />
        <div class="min-w-0 flex-1">
          <h4 class="font-bold text-slate-100 text-sm truncate">{{ app.name }}</h4>
          <span class="text-[11px] text-slate-400 block truncate">{{ app.tagline || app.category }}</span>
        </div>
      </div>

      <!-- Konfigurasi Port Host -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Port Akses Host</label>
        <div class="relative">
          <input
            v-model="customPort"
            type="number"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
            placeholder="8080"
          />
        </div>
        <p class="text-[10px] text-slate-500 mt-1">Port host lokal STB untuk mengakses Web UI container.</p>
      </div>

      <!-- Konfigurasi Mountpoint Storage -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Direktori Penyimpanan Data Permanen</label>
        <input
          v-model="customVolume"
          type="text"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
          placeholder="/DATA/apps/nama-aplikasi"
        />
        <p class="text-[10px] text-slate-500 mt-1">Direktori penyimpanan persisten di MicroSD/Storage eksternal.</p>
      </div>
    </div>

    <!-- Action Buttons Footer -->
    <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-800/80 shrink-0">
      <button
        @click="handleCancel"
        type="button"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs transition"
      >
        Batal
      </button>

      <button
        type="submit"
        :disabled="isInstalling"
        class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow-lg active:scale-95 disabled:opacity-50"
      >
        <Download class="w-3.5 h-3.5" />
        <span>{{ isInstalling ? 'Sedang Memasang...' : 'Pasang Sekarang' }}</span>
      </button>
    </div>
  </form>
</template>
