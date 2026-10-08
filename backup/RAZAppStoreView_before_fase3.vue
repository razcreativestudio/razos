<script setup>
/**
 * RAZOS App Store View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka katalog aplikasi Docker 1-klik untuk di-deploy ke direktori MicroSD /DATA/apps.
 */
import { ref } from 'vue'
import { Download, Check, Box } from 'lucide-vue-next'

const catalog = ref([
  {
    id: 'caddy',
    name: 'Caddy Server',
    desc: 'Web server & Reverse Proxy otomatis dengan sertifikat SSL gratis.',
    category: 'Network',
    installed: true,
  },
  {
    id: 'mariadb',
    name: 'MariaDB Enterprise',
    desc: 'Relational Database Server SQL cepat dan hemat memori.',
    category: 'Database',
    installed: true,
  },
  {
    id: 'postgres',
    name: 'PostgreSQL Server',
    desc: 'Object-relational SQL database canggih untuk aplikasi enterprise.',
    category: 'Database',
    installed: true,
  },
  {
    id: 'adminer',
    name: 'Adminer Web GUI',
    desc: 'Tool manajemen database berbasis web satu file yang super ringan.',
    category: 'Database',
    installed: true,
  },
  {
    id: 'portainer',
    name: 'Portainer CE',
    desc: 'Dashboard visual manajemen container Docker.',
    category: 'Tools',
    installed: false,
  },
  {
    id: 'uptime-kuma',
    name: 'Uptime Kuma',
    desc: 'Monitoring status server dan website self-hosted.',
    category: 'Monitoring',
    installed: false,
  }
])
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">App Store RAZOS</h3>
        <p class="text-xs text-slate-400">Instal aplikasi Docker yang dioptimalkan untuk penyimpanan MicroSD (/DATA/apps).</p>
      </div>
      <img src="/assets/brand/RAZBrandCapsule.webp" alt="RAZ" class="h-5 object-contain" />
    </div>

    <!-- App Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3 overflow-auto flex-1 pr-1">
      <div
        v-for="app in catalog"
        :key="app.id"
        class="p-3 rounded-xl bg-slate-900/80 border border-slate-800 flex items-start justify-between gap-3 hover:border-slate-700 transition"
      >
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-xl bg-sky-500/10 border border-sky-500/30 flex items-center justify-center text-sky-400 shrink-0">
            <Box class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-semibold text-slate-200">{{ app.name }}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-400">{{ app.category }}</span>
            </div>
            <p class="text-xs text-slate-400 mt-1 leading-relaxed">{{ app.desc }}</p>
          </div>
        </div>

        <div>
          <button
            v-if="app.installed"
            class="px-2.5 py-1.5 rounded-lg bg-slate-800 text-emerald-400 flex items-center gap-1.5 text-xs font-medium cursor-default"
          >
            <Check class="w-3.5 h-3.5" />
            <span>Terpasang</span>
          </button>
          <button
            v-else
            class="px-2.5 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow"
          >
            <Download class="w-3.5 h-3.5" />
            <span>Pasang</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
