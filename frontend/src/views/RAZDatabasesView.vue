<script setup>
/**
 * RAZOS Database Manager View
 * Updated: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Antarmuka status database engine MariaDB & PostgreSQL dan shortcut Adminer Web GUI.
 * Form pembuatan database baru terbuka di Jendela OS mandiri tanpa memblokir Database Manager.
 */
import { ref, onMounted, watch } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { Database, ExternalLink, Plus, Inbox, RotateCw } from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()
const databases = ref([])
const isLoading = ref(false)

async function fetchDatabases() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/databases')
    if (res.ok) {
      databases.value = await res.json()
    }
  } catch (err) {
    console.error('Gagal mengambil daftar database:', err)
  }
  isLoading.value = false
}

// Watch trigger reaktif dari jendela Database Creator mandiri
watch(() => systemStore.databasesRefreshTrigger, () => {
  fetchDatabases()
})

async function launchAdminer() {
  try {
    const res = await fetch('/api/v1/databases/adminer/launch', { method: 'POST' })
    if (res.ok) {
      window.open('http://' + window.location.hostname + ':8088', '_blank')
    } else {
      systemStore.showNotification('Gagal', 'Container Adminer sedang tidak aktif', 'error')
    }
  } catch (e) {
    // Fallback jika mode dev
    window.open('http://' + window.location.hostname + ':8088', '_blank')
  }
}

onMounted(() => {
  fetchDatabases()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm select-none">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Database Engine Suite</h3>
        <p class="text-xs text-slate-400">Database server MariaDB & PostgreSQL dengan akses web Adminer instan.</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="fetchDatabases"
          :disabled="isLoading"
          class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 border border-slate-800 text-slate-300 hover:text-white transition"
          title="Muat Ulang Status Database"
        >
          <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
        </button>

        <!-- Buka Form Buat Database di Jendela OS Mandiri -->
        <button
          @click="windowStore.openDatabaseCreator()"
          class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow active:scale-95"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Buat Database</span>
        </button>

        <button
          @click="launchAdminer"
          class="px-3 py-1.5 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow active:scale-95"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Buka Web Adminer</span>
        </button>
      </div>
    </div>

    <!-- Database Cards Grid -->
    <div v-if="databases.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-3 overflow-auto flex-1">
      <div
        v-for="db in databases"
        :key="db.id"
        class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2.5 hover:border-slate-700 transition"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Database class="w-4 h-4 text-amber-400" />
            <span class="font-semibold text-slate-200 text-sm">{{ db.name }}</span>
          </div>
          <span class="text-[10px] px-2 py-0.5 rounded bg-slate-800 text-amber-300 font-mono uppercase font-bold">{{ db.engine }}</span>
        </div>

        <div class="grid grid-cols-2 gap-2 text-xs text-slate-400 mt-1 bg-slate-950/60 p-2.5 rounded-lg border border-slate-900">
          <div>Pengguna: <span class="text-slate-200 font-mono">{{ db.user }}</span></div>
          <div>Port: <span class="text-slate-200 font-mono">{{ db.port }}</span></div>
          <div>Ukuran Data: <span class="text-slate-200 font-mono">{{ db.size }}</span></div>
          <div>Host Akses: <span class="text-slate-200 font-mono">{{ db.host }}</span></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="flex-1 rounded-xl bg-slate-900/40 border border-slate-800/60 flex flex-col items-center justify-center text-slate-500 gap-3 p-8">
      <div class="w-12 h-12 rounded-2xl bg-slate-800/50 flex items-center justify-center text-slate-400">
        <Inbox class="w-6 h-6" />
      </div>
      <div class="text-center">
        <p class="font-medium text-slate-300">Belum ada basis data yang dibuat</p>
        <p class="text-xs text-slate-500 mt-1">Gunakan tombol "Buat Database" di atas untuk men-deploy basis data MariaDB atau PostgreSQL baru.</p>
      </div>
      <button
        @click="windowStore.openDatabaseCreator()"
        class="mt-2 px-3.5 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold text-xs flex items-center gap-1.5 transition shadow"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Buat Database Sekarang</span>
      </button>
    </div>
  </div>
</template>
