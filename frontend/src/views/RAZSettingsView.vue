<script setup>
/**
 * RAZOS Settings View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka konfigurasi sistem, identitas host, dan rilis versi.
 */
import { ref } from 'vue'

const hostname = ref('razos-s905x')
const timezone = ref('Asia/Jakarta (WIB)')
const storageTarget = ref('/DATA (MicroSD 32GB)')
</script>

<template>
  <div class="h-full flex flex-col gap-4 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Pengaturan Sistem</h3>
        <p class="text-xs text-slate-400">Konfigurasi umum host dan preferensi sistem RAZOS.</p>
      </div>
      <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ" class="h-5 object-contain" />
    </div>

    <!-- Settings Form Cards -->
    <div class="space-y-3 overflow-auto flex-1">
      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Hostname Perangkat</label>
        <input
          v-model="hostname"
          type="text"
          class="px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 outline-none focus:border-sky-500 font-mono"
        />
      </div>

      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Zona Waktu Sistem</label>
        <input
          v-model="timezone"
          type="text"
          disabled
          class="px-3 py-2 rounded-lg bg-slate-950/60 border border-slate-800 text-slate-400 font-mono"
        />
      </div>

      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2">
        <label class="font-medium text-slate-300">Direktori Penyimpanan Utama</label>
        <input
          v-model="storageTarget"
          type="text"
          disabled
          class="px-3 py-2 rounded-lg bg-slate-950/60 border border-slate-800 text-slate-400 font-mono"
        />
      </div>

      <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <div class="font-medium text-slate-200">Tentang RAZOS</div>
          <div class="text-xs text-slate-400">Versi 1.0.0-rc1 | RAZ CREATIVE STUDIO</div>
        </div>
        <img src="/assets/images/RAZLogo.svg" alt="RAZOS" class="w-8 h-8" />
      </div>
    </div>
  </div>
</template>
