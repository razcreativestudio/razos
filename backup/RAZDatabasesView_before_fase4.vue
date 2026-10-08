<script setup>
/**
 * RAZOS Database Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka status database engine MariaDB & PostgreSQL dan shortcut Adminer Web GUI.
 * Zero dummy data: Menampilkan daftar riil database dari backend API.
 */
import { ref, onMounted } from 'vue'
import { Database, ExternalLink, HardDrive, Inbox } from 'lucide-vue-next'

const databases = ref([])
const isLoading = ref(false)

async function fetchDatabases() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/databases')
    if (res.ok) {
      databases.value = await res.json()
    } else {
      databases.value = []
    }
  } catch (err) {
    databases.value = []
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchDatabases()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Database Engine Suite</h3>
        <p class="text-xs text-slate-400">Database server terisolasi dengan akses antarmuka web GUI Adminer.</p>
      </div>

      <a
        href="http://localhost:8088"
        target="_blank"
        class="px-3 py-1.5 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition"
      >
        <ExternalLink class="w-3.5 h-3.5" />
        <span>Buka Web Adminer</span>
      </a>
    </div>

    <!-- Database Cards (Jika ada database terdaftar) -->
    <div v-if="databases.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div
        v-for="db in databases"
        :key="db.id"
        class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-col gap-2"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Database class="w-4 h-4 text-amber-400" />
            <span class="font-semibold text-slate-200">{{ db.name }}</span>
          </div>
          <span class="text-[10px] px-2 py-0.5 rounded bg-slate-800 text-amber-300 font-mono">{{ db.engine }}</span>
        </div>

        <div class="grid grid-cols-2 gap-2 text-xs text-slate-400 mt-1">
          <div>User: <span class="text-slate-200 font-mono">{{ db.user }}</span></div>
          <div>Port: <span class="text-slate-200 font-mono">{{ db.port }}</span></div>
          <div>Ukuran: <span class="text-slate-200 font-mono">{{ db.size }}</span></div>
          <div>Host: <span class="text-slate-200 font-mono">127.0.0.1</span></div>
        </div>
      </div>
    </div>

    <!-- Empty State Elegan (Zero Dummy Data) -->
    <div v-else class="flex-1 rounded-xl bg-slate-900/40 border border-slate-800/60 flex flex-col items-center justify-center text-slate-500 gap-3 p-8">
      <div class="w-12 h-12 rounded-2xl bg-slate-800/50 flex items-center justify-center text-slate-400">
        <Inbox class="w-6 h-6" />
      </div>
      <div class="text-center">
        <p class="font-medium text-slate-300">Belum ada basis data yang dibuat</p>
        <p class="text-xs text-slate-500 mt-1">Gunakan wizard pembuatan website untuk men-generate database MariaDB atau PostgreSQL otomatis.</p>
      </div>
    </div>
  </div>
</template>
