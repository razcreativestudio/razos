<script setup>
/**
 * RAZOS Database Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka status database engine MariaDB & PostgreSQL dan shortcut Adminer Web GUI.
 */
import { ref } from 'vue'
import { Database, ExternalLink, HardDrive } from 'lucide-vue-next'

const dbs = ref([
  {
    id: 1,
    name: 'razos_core_db',
    engine: 'MariaDB 10.11',
    size: '14.2 MB',
    user: 'raz_admin',
    port: 3306
  },
  {
    id: 2,
    name: 'analytics_db',
    engine: 'PostgreSQL 16',
    size: '48.5 MB',
    user: 'raz_pg',
    port: 5432
  }
])
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Database Engine Suite</h3>
        <p class="text-xs text-slate-400">Database server MariaDB & PostgreSQL dengan akses web Adminer.</p>
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

    <!-- Database Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div
        v-for="db in dbs"
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
  </div>
</template>
