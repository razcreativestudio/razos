<script setup>
/**
 * RAZOS Website Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka pengelolaan website & port reverse proxy.
 */
import { ref } from 'vue'
import { Globe, Plus, ExternalLink } from 'lucide-vue-next'

const websites = ref([
  {
    id: 1,
    name: 'Portal Perusahaan',
    domain: 'portal.razos.my.id',
    port: 8101,
    runtime: 'Node.js',
    status: 'Active'
  },
  {
    id: 2,
    name: 'Blog Utama',
    domain: 'blog.razos.my.id',
    port: 8102,
    runtime: 'PHP 8.2',
    status: 'Active'
  }
])
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Website Suite & Reverse Proxy</h3>
        <p class="text-xs text-slate-400">Pengelolaan multi-website terisolasi pada alokasi port 8100-8200.</p>
      </div>

      <button class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition">
        <Plus class="w-3.5 h-3.5" />
        <span>Tambah Website</span>
      </button>
    </div>

    <!-- Website List Table -->
    <div class="flex-1 rounded-xl bg-slate-900/60 border border-slate-800/80 p-3 overflow-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 text-slate-400">
            <th class="py-2 px-2">Nama Website</th>
            <th class="py-2 px-2">Domain Hostname</th>
            <th class="py-2 px-2">Port Proxy</th>
            <th class="py-2 px-2">Runtime</th>
            <th class="py-2 px-2">Status</th>
            <th class="py-2 px-2 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-slate-300">
          <tr v-for="site in websites" :key="site.id">
            <td class="py-2.5 px-2 font-medium flex items-center gap-2">
              <Globe class="w-3.5 h-3.5 text-sky-400" />
              <span>{{ site.name }}</span>
            </td>
            <td class="py-2.5 px-2 font-mono text-sky-400">{{ site.domain }}</td>
            <td class="py-2.5 px-2 font-mono">{{ site.port }}</td>
            <td class="py-2.5 px-2">{{ site.runtime }}</td>
            <td class="py-2.5 px-2">
              <span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400">{{ site.status }}</span>
            </td>
            <td class="py-2.5 px-2 text-right">
              <a
                :href="'http://' + site.domain"
                target="_blank"
                class="inline-flex items-center gap-1 px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
              >
                <ExternalLink class="w-3 h-3" />
                <span>Buka</span>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
