<script setup>
/**
 * RAZOS Website Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka pengelolaan website & reverse proxy internal.
 * Dilengkapi Wizard Pembuatan Website (Node.js, PHP, Static) dengan auto-assign port 8100-8200
 * dan koneksi database instan.
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Globe, Plus, ExternalLink, Inbox, Trash2, X, Server, Layers } from 'lucide-vue-next'

const systemStore = useSystemStore()
const websites = ref([])
const databases = ref([])
const isLoading = ref(false)

// State Modal Wizard Buat Website
const showCreateModal = ref(false)
const siteName = ref('')
const siteDomain = ref('')
const siteRuntime = ref('php')
const selectedDbId = ref('')

async function fetchWebsites() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/websites')
    if (res.ok) {
      websites.value = await res.json()
    }
  } catch (err) {}
  isLoading.value = false
}

async function fetchDatabases() {
  try {
    const res = await fetch('/api/v1/databases')
    if (res.ok) {
      databases.value = await res.json()
    }
  } catch (err) {}
}

function openCreateDialog() {
  siteName.value = ''
  siteDomain.value = ''
  siteRuntime.value = 'php'
  selectedDbId.value = ''
  fetchDatabases()
  showCreateModal.value = true
}

async function confirmCreateWebsite() {
  if (!siteName.value || !siteDomain.value) {
    systemStore.showNotification('Peringatan', 'Nama dan Domain website wajib diisi', 'error')
    return
  }

  const payload = {
    name: siteName.value.trim(),
    domain: siteDomain.value.trim(),
    runtime: siteRuntime.value,
    database_id: selectedDbId.value ? parseInt(selectedDbId.value) : null
  }

  try {
    const res = await fetch('/api/v1/websites', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      const created = await res.json()
      systemStore.showNotification('Sukses', `Website ${created.name} berhasil dibuat pada port ${created.port}!`, 'success')
      showCreateModal.value = false
      fetchWebsites()
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal', err.error || 'Gagal membuat website', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi backend', 'error')
  }
}

async function deleteWebsite(site) {
  try {
    const res = await fetch(`/api/v1/websites/${site.id}`, { method: 'DELETE' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Website ${site.name} telah dihapus`, 'success')
      fetchWebsites()
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghapus website', 'error')
  }
}

onMounted(() => {
  fetchWebsites()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Website Suite & Reverse Proxy</h3>
        <p class="text-xs text-slate-400">Pengelolaan multi-website terisolasi pada alokasi port 8100-8200.</p>
      </div>

      <button
        @click="openCreateDialog"
        class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Wizard Buat Website</span>
      </button>
    </div>

    <!-- Website List Table -->
    <div v-if="websites.length > 0" class="flex-1 rounded-xl bg-slate-900/60 border border-slate-800/80 p-3 overflow-auto">
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
            <td class="py-2.5 px-2 uppercase font-semibold text-[11px] text-slate-400">{{ site.runtime }}</td>
            <td class="py-2.5 px-2">
              <span class="px-2 py-0.5 rounded text-[10px] bg-emerald-500/20 text-emerald-400 font-medium">{{ site.status }}</span>
            </td>
            <td class="py-2.5 px-2 text-right flex items-center justify-end gap-2">
              <a
                :href="'http://' + site.domain"
                target="_blank"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 transition text-[11px]"
              >
                <ExternalLink class="w-3 h-3" />
                <span>Buka</span>
              </a>
              <button
                @click="deleteWebsite(site)"
                class="p-1 rounded hover:bg-rose-500/20 text-rose-400 transition"
                title="Hapus Website"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-else class="flex-1 rounded-xl bg-slate-900/40 border border-slate-800/60 flex flex-col items-center justify-center text-slate-500 gap-3 p-8">
      <div class="w-12 h-12 rounded-2xl bg-slate-800/50 flex items-center justify-center text-slate-400">
        <Inbox class="w-6 h-6" />
      </div>
      <div class="text-center">
        <p class="font-medium text-slate-300">Belum ada website yang dikonfigurasi</p>
        <p class="text-xs text-slate-500 mt-1">Gunakan tombol "Wizard Buat Website" di atas untuk membuat situs PHP, Node.js, atau HTML statis baru.</p>
      </div>
    </div>

    <!-- CUSTOM MODAL: Wizard Pembuatan Website -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Layers class="w-5 h-5 text-sky-400" />
            <h4 class="font-bold text-slate-100">Wizard Pembuatan Website</h4>
          </div>
          <button @click="showCreateModal = false" class="p-1 text-slate-400 hover:text-slate-200">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nama Website / Project</label>
            <input
              v-model="siteName"
              type="text"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: Blog Perusahaan"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Domain Hostname (Routing Proxy)</label>
            <input
              v-model="siteDomain"
              type="text"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="Contoh: blog.razos.my.id"
            />
            <span class="text-[10px] text-slate-500 mt-1 block">Akan di-route otomatis oleh proxy internal ke port bebas (8100-8200).</span>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Pilihan Lingkungan Runtime</label>
            <div class="grid grid-cols-3 gap-2">
              <button
                type="button"
                @click="siteRuntime = 'php'"
                class="py-2 rounded-lg border text-xs font-medium transition"
                :class="siteRuntime === 'php' ? 'bg-sky-500/20 text-sky-400 border-sky-400/50' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                PHP 8.2
              </button>
              <button
                type="button"
                @click="siteRuntime = 'nodejs'"
                class="py-2 rounded-lg border text-xs font-medium transition"
                :class="siteRuntime === 'nodejs' ? 'bg-sky-500/20 text-sky-400 border-sky-400/50' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                Node.js
              </button>
              <button
                type="button"
                @click="siteRuntime = 'static'"
                class="py-2 rounded-lg border text-xs font-medium transition"
                :class="siteRuntime === 'static' ? 'bg-sky-500/20 text-sky-400 border-sky-400/50' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                HTML Statis
              </button>
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Koneksi Database Sentral (Opsional)</label>
            <select
              v-model="selectedDbId"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-sky-400"
            >
              <option value="">Tanpa Database</option>
              <option v-for="db in databases" :key="db.id" :value="db.id">
                {{ db.name }} ({{ db.engine.toUpperCase() }} - Port {{ db.port }})
              </option>
            </select>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showCreateModal = false"
            class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmCreateWebsite"
            class="px-4 py-2 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow"
          >
            Generate & Pasang Website
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
