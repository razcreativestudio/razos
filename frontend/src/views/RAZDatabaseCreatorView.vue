<script setup>
/**
 * RAZOS Database Creator View (Standalone Desktop Window)
 * Updated: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Form pembuatan database baru MariaDB / PostgreSQL yang berjalan di jendela desktop mandiri.
 * Tidak mengunci jendela aplikasi Database Manager utama.
 */
import { ref } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { Database, ShieldCheck, Check, X, Server } from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ id: 'database-create' })
  }
})

const dbName = ref('')
const dbEngine = ref('mariadb')
const dbUser = ref('raz_user')
const dbPassword = ref('')
const isSubmitting = ref(false)

async function handleSubmit() {
  if (!dbName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama database wajib diisi', 'error')
    return
  }

  isSubmitting.value = true
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
      // Picu refresh reaktif global ke semua window
      systemStore.triggerDatabasesRefresh()
      // Tutup window pembuatan DB
      windowStore.closeWindow(props.win?.id || 'database-create')
    } else {
      systemStore.showNotification('Gagal', 'Gagal membuat database baru', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isSubmitting.value = false
  }
}

function handleCancel() {
  windowStore.closeWindow(props.win?.id || 'database-create')
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="h-full flex flex-col justify-between p-4 bg-slate-900/90 text-xs sm:text-sm select-none">
    <div class="space-y-4 overflow-y-auto pr-1">
      <!-- Info Header -->
      <div class="flex items-center gap-3 p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-300 text-xs">
        <Server class="w-5 h-5 flex-shrink-0 text-amber-400" />
        <div>
          <p class="font-semibold text-slate-100">Buat Database Baru</p>
          <p class="text-[11px] text-slate-400">Database akan otomatis di-binding ke port engine lokal STB.</p>
        </div>
      </div>

      <!-- Nama Database -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Nama Database</label>
        <input
          v-model="dbName"
          type="text"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 transition"
          placeholder="Contoh: toko_online_db"
          autofocus
        />
        <p class="text-[10px] text-slate-500 mt-1">Gunakan huruf kecil tanpa spasi (a-z, 0-9, underscore).</p>
      </div>

      <!-- Pilihan Engine -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Pilihan Engine Database</label>
        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            @click="dbEngine = 'mariadb'"
            class="py-2.5 px-3 rounded-xl border text-xs font-medium flex items-center justify-center gap-2 transition"
            :class="dbEngine === 'mariadb' ? 'bg-amber-500/20 text-amber-300 border-amber-400/60 shadow-sm' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
          >
            <Database class="w-3.5 h-3.5" />
            <span>MariaDB (Port 3306)</span>
          </button>
          <button
            type="button"
            @click="dbEngine = 'postgres'"
            class="py-2.5 px-3 rounded-xl border text-xs font-medium flex items-center justify-center gap-2 transition"
            :class="dbEngine === 'postgres' ? 'bg-amber-500/20 text-amber-300 border-amber-400/60 shadow-sm' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
          >
            <Database class="w-3.5 h-3.5" />
            <span>PostgreSQL (Port 5432)</span>
          </button>
        </div>
      </div>

      <!-- User Pengguna -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Username Pengguna</label>
        <input
          v-model="dbUser"
          type="text"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 transition"
          placeholder="Contoh: raz_user"
        />
      </div>

      <!-- Password -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">Kata Sandi Database</label>
        <input
          v-model="dbPassword"
          type="password"
          class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 transition"
          placeholder="Kata sandi..."
        />
        <p class="text-[10px] text-slate-500 mt-1">Kosongkan jika ingin menggunakan password default root.</p>
      </div>
    </div>

    <!-- Action Buttons Footer -->
    <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-800/80 mt-3">
      <button
        @click="handleCancel"
        type="button"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs transition"
      >
        Batal
      </button>
      <button
        @click="handleSubmit"
        :disabled="isSubmitting"
        type="button"
        class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow active:scale-95 disabled:opacity-50"
      >
        <Check class="w-3.5 h-3.5" />
        <span>{{ isSubmitting ? 'Memproses...' : 'Simpan & Aktifkan DB' }}</span>
      </button>
    </div>
  </form>
</template>
