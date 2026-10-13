<script setup>
/**
 * RAZOS Database Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka status database engine MariaDB & PostgreSQL dan shortcut Adminer Web GUI.
 * Dilengkapi dialog modal kustom pembuatan database baru serta peluncuran container Adminer Web.
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { Database, ExternalLink, Plus, Inbox, X, HardDrive, ShieldCheck } from 'lucide-vue-next'

const systemStore = useSystemStore()
const databases = ref([])
const isLoading = ref(false)

// State Modal Pembuatan Database Baru
const showCreateDbModal = ref(false)
const dbName = ref('')
const dbEngine = ref('mariadb')
const dbUser = ref('raz_user')
const dbPassword = ref('')

async function fetchDatabases() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/databases')
    if (res.ok) {
      databases.value = await res.json()
    }
  } catch (err) {}
  isLoading.value = false
}

function openCreateDbDialog() {
  dbName.value = ''
  dbEngine.value = 'mariadb'
  dbUser.value = 'raz_user'
  dbPassword.value = ''
  showCreateDbModal.value = true
}

async function confirmCreateDatabase() {
  if (!dbName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama database wajib diisi', 'error')
    return
  }

  const payload = {
    name: dbName.value.trim(),
    engine: dbEngine.value,
    user: dbUser.value.trim(),
    password: dbPassword.value
  }

  try {
    const res = await fetch('/api/v1/databases', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      const created = await res.json()
      systemStore.showNotification('Sukses', `Database ${created.name} (${created.engine.toUpperCase()}) berhasil dibuat!`, 'success')
      showCreateDbModal.value = false
      fetchDatabases()
    } else {
      systemStore.showNotification('Gagal', 'Gagal membuat database baru', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

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
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h3 class="font-bold text-slate-100 text-base">Database Engine Suite</h3>
        <p class="text-xs text-slate-400">Database server MariaDB & PostgreSQL dengan akses web Adminer instan.</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="openCreateDbDialog"
          class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Buat Database</span>
        </button>

        <button
          @click="launchAdminer"
          class="px-3 py-1.5 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-semibold flex items-center gap-1.5 text-xs transition shadow"
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
    </div>

    <!-- CUSTOM MODAL: Form Buat Database Baru -->
    <div
      v-if="showCreateDbModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Database class="w-5 h-5 text-amber-400" />
            <h4 class="font-bold text-slate-100">Buat Database Baru</h4>
          </div>
          <button @click="showCreateDbModal = false" class="p-1 text-slate-400 hover:text-slate-200">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nama Database</label>
            <input
              v-model="dbName"
              type="text"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
              placeholder="Contoh: my_app_db"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Pilihan Engine Database</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="dbEngine = 'mariadb'"
                class="py-2 rounded-lg border text-xs font-medium transition"
                :class="dbEngine === 'mariadb' ? 'bg-amber-500/20 text-amber-300 border-amber-400/50' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                MariaDB (Port 3306)
              </button>
              <button
                type="button"
                @click="dbEngine = 'postgres'"
                class="py-2 rounded-lg border text-xs font-medium transition"
                :class="dbEngine === 'postgres' ? 'bg-amber-500/20 text-amber-300 border-amber-400/50' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                PostgreSQL (Port 5432)
              </button>
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Username Pengguna</label>
            <input
              v-model="dbUser"
              type="text"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
              placeholder="Contoh: raz_user"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Kata Sandi Database</label>
            <input
              v-model="dbPassword"
              type="password"
              class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
              placeholder="Kata sandi..."
            />
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showCreateDbModal = false"
            class="px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs transition"
          >
            Batal
          </button>
          <button
            @click="confirmCreateDatabase"
            class="px-4 py-2 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs transition shadow"
          >
            Simpan & Aktifkan DB
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
