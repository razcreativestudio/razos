<script setup>
/**
 * RAZOS App Store Manual Deploy View (Standalone Desktop Window)
 * Created: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Form Deploy Cepat Docker Container kustom yang berjalan sebagai Jendela OS mandiri.
 * Tidak mengunci atau menutupi jendela utama App Store.
 */
import { ref } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { PlusCircle, Check, X, Box, HardDrive, Terminal } from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ id: 'appstore-manual' })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const manualName = ref('')
const manualImage = ref('')
const manualPort = ref('8080')
const manualVolume = ref('/DATA/apps/custom-app')
const manualEnvs = ref('')
const isDeployingManual = ref(false)

async function handleSubmit() {
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
      // Picu refresh global ke App Store dan Desktop
      systemStore.triggerAppsRefresh()
      await windowStore.fetchInstalledApps()
      // Tutup jendela instal cepat mandiri
      windowStore.closeWindow(props.win?.id || 'appstore-manual')
    } else {
      const errData = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Deploy', errData.error || 'Gagal men-deploy container', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isDeployingManual.value = false
  }
}

function handleCancel() {
  windowStore.closeWindow(props.win?.id || 'appstore-manual')
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="h-full flex flex-col justify-between p-4 bg-slate-900/90 text-xs sm:text-sm select-none gap-3">
    <div class="space-y-3.5 overflow-y-auto pr-1">
      <!-- Info Header -->
      <div class="flex items-center gap-3 p-3 rounded-xl bg-sky-500/10 border border-sky-500/20 text-sky-300 text-xs">
        <Box class="w-5 h-5 flex-shrink-0 text-sky-400" />
        <div>
          <p class="font-semibold text-slate-100">Instal Cepat Docker Container</p>
          <p class="text-[11px] text-slate-400">Deploy container mandiri dengan parameter cepat tanpa compose YAML.</p>
        </div>
      </div>

      <!-- Nama Aplikasi -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">
          Nama Aplikasi / Container <span class="text-rose-400">*</span>
        </label>
        <input
          v-model="manualName"
          type="text"
          placeholder="Contoh: My Custom App"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400 transition"
          autofocus
        />
      </div>

      <!-- Docker Image -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">
          Docker Image (Docker Hub / Registry) <span class="text-rose-400">*</span>
        </label>
        <input
          v-model="manualImage"
          type="text"
          placeholder="Contoh: nginx:alpine atau redis:7-alpine"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
        />
      </div>

      <!-- Host Port & Volume -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label class="block text-xs font-medium text-slate-300 mb-1">Host Port</label>
          <input
            v-model="manualPort"
            type="number"
            placeholder="8080"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
          />
        </div>

        <div>
          <label class="block text-xs font-medium text-slate-300 mb-1">Direktori Volume Host</label>
          <input
            v-model="manualVolume"
            type="text"
            placeholder="/DATA/apps/custom-app"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
          />
        </div>
      </div>

      <!-- Environment Variables -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Environment Variables (Satu per baris)</label>
        <textarea
          v-model="manualEnvs"
          rows="3"
          placeholder="TZ=Asia/Jakarta&#10;APP_ENV=production"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition resize-none"
        ></textarea>
      </div>
    </div>

    <!-- Actions Footer -->
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
        :disabled="isDeployingManual"
        class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow-lg active:scale-95 disabled:opacity-50"
      >
        <Check class="w-3.5 h-3.5" />
        <span>{{ isDeployingManual ? 'Memasang Container...' : 'Pasang Sekarang' }}</span>
      </button>
    </div>
  </form>
</template>
