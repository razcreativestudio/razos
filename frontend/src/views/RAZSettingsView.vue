<script setup>
/**
 * RAZOS Settings View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka konfigurasi sistem, identitas host, dan OTA Updater rilis GitHub.
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { ArrowUpCircle, CheckCircle2, RefreshCw, Sparkles, ShieldCheck } from 'lucide-vue-next'

const systemStore = useSystemStore()
const hostname = ref('razos-s905x')
const timezone = ref('Asia/Jakarta (WIB)')
const storageTarget = ref('/DATA (MicroSD 32GB)')

// State OTA Updater
const isCheckingUpdate = ref(false)
const updateStatus = ref({
  current_version: '1.5.0',
  latest_version: '1.5.0',
  has_update: false,
  release_notes: 'Sistem dalam kondisi paling mutakhir.',
  download_url: ''
})

async function checkUpdate() {
  isCheckingUpdate.value = true
  try {
    const res = await fetch('/api/v1/updater/status')
    if (res.ok) {
      updateStatus.value = await res.json()
      if (updateStatus.value.has_update) {
        systemStore.showNotification('Pembaruan Tersedia', `Versi ${updateStatus.value.latest_version} siap diinstal!`, 'info')
      } else {
        systemStore.showNotification('Sistem Up to Date', 'RAZOS menggunakan versi terbaru', 'success')
      }
    }
  } catch (err) {}
  isCheckingUpdate.value = false
}

async function applyUpdate() {
  if (!updateStatus.value.download_url) return
  try {
    const res = await fetch('/api/v1/updater/apply', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ download_url: updateStatus.value.download_url })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', 'Pembaruan OTA berhasil dipasang!', 'success')
      checkUpdate()
    } else {
      systemStore.showNotification('Gagal', 'Gagal memasang pembaruan otomatis', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

onMounted(() => {
  checkUpdate()
})
</script>

<template>
  <div class="h-full flex flex-col gap-4 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Pengaturan Sistem & OTA Updater</h3>
        <p class="text-xs text-slate-400">Konfigurasi umum host dan pembaruan rilis sistem RAZOS.</p>
      </div>
      <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ" class="h-5 object-contain" />
    </div>

    <!-- Settings Content Cards -->
    <div class="space-y-3 overflow-auto flex-1 pr-1">
      <!-- OTA Updater Card -->
      <div class="p-4 rounded-xl bg-slate-900/90 border border-sky-500/30 flex flex-col gap-3 shadow-lg">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <ArrowUpCircle class="w-5 h-5 text-sky-400" />
            <span class="font-bold text-slate-100 text-sm">Pembaruan Sistem (OTA Updater)</span>
          </div>
          <span class="text-[11px] px-2 py-0.5 rounded font-mono bg-sky-500/20 text-sky-300">
            v{{ updateStatus.current_version }}
          </span>
        </div>

        <p class="text-xs text-slate-400 leading-relaxed">
          {{ updateStatus.release_notes }}
        </p>

        <div class="flex items-center justify-between pt-2 border-t border-slate-800">
          <div class="text-[11px] text-slate-400">
            Status: <span class="text-emerald-400 font-medium">{{ updateStatus.has_update ? 'Versi Baru Tersedia' : 'Paling Mutakhir' }}</span>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="checkUpdate"
              :disabled="isCheckingUpdate"
              class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs flex items-center gap-1.5 transition"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="isCheckingUpdate ? 'animate-spin' : ''" />
              <span>Periksa Update</span>
            </button>

            <button
              v-if="updateStatus.has_update"
              @click="applyUpdate"
              class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow"
            >
              <Sparkles class="w-3.5 h-3.5" />
              <span>Perbarui Sekarang</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Host Settings -->
      <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Hostname Perangkat</label>
        <input
          v-model="hostname"
          type="text"
          class="px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 outline-none focus:border-sky-500 font-mono text-xs"
        />
      </div>

      <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Zona Waktu Sistem</label>
        <input
          v-model="timezone"
          type="text"
          disabled
          class="px-3 py-2 rounded-lg bg-slate-950/60 border border-slate-800 text-slate-400 font-mono text-xs"
        />
      </div>

      <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Direktori Penyimpanan Utama</label>
        <input
          v-model="storageTarget"
          type="text"
          disabled
          class="px-3 py-2 rounded-lg bg-slate-950/60 border border-slate-800 text-slate-400 font-mono text-xs"
        />
      </div>

      <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <div class="font-medium text-slate-200">Tentang RAZOS</div>
          <div class="text-xs text-slate-400">Versi 1.5.0-production | RAZ CREATIVE STUDIO</div>
        </div>
        <img src="/assets/images/RAZLogo.svg" alt="RAZOS" class="w-8 h-8" />
      </div>
    </div>
  </div>
</template>
