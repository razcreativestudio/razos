<script setup>
/**
 * RAZOS Website Git Sync Window Component
 * Created: 2026-10-13 (v2.6.0 - Standalone Window Architecture)
 * Deskripsi: Komponen pengelolaan GitHub & Git Pull mandiri.
 * Berjalan sebagai Jendela OS terpisah (tidak memblokir/mengganggu Website Manager).
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  GitBranch, GitPullRequest, Key, Eye, EyeOff, RotateCw,
  FolderOpen, Terminal, CheckCircle2
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ customData: {} })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const site = computed(() => props.win.customData || {})

const gitModalRepo = ref('')
const gitModalBranch = ref('main')
const gitModalToken = ref('')
const gitModalShowToken = ref(false)
const gitModalIsPulling = ref(false)
const gitModalLogs = ref('')

function initGitData() {
  const s = site.value
  if (!s || !s.id) return
  gitModalRepo.value = s.git_repo || ''
  gitModalBranch.value = s.git_branch || 'main'
  gitModalToken.value = s.git_token || ''
  gitModalLogs.value = s.git_last_commit
    ? `[INFO] Commit Terakhir: ${s.git_last_commit}\n[INFO] Waktu Pull: ${s.git_last_pull || '-'}`
    : 'Belum ada log tarikan kode (pull) sebelumnya. Masukkan repo dan klik Tarik Kode Sekarang.'
}

// Buka folder root di File Manager
function openWebsiteFolder() {
  if (!site.value || !site.value.root_folder) return
  windowStore.openApp('file-manager')
}

// Simpan konfigurasi Git
async function saveGitModalConfig() {
  const s = site.value
  if (!s || !s.id) return
  try {
    const res = await fetch(`/api/v1/websites/${s.id}/git`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        git_repo: gitModalRepo.value.trim(),
        git_branch: gitModalBranch.value.trim() || 'main',
        git_token: gitModalToken.value.trim()
      })
    })
    const data = await res.json()
    if (res.ok) {
      s.git_repo = data.website.git_repo
      s.git_branch = data.website.git_branch
      s.git_token = data.website.git_token
      systemStore.showNotification('Sukses', 'Konfigurasi GitHub & PAT berhasil disimpan', 'success')
      systemStore.triggerWebsitesRefresh()
    } else {
      systemStore.showNotification('Gagal', data.error || 'Gagal menyimpan konfigurasi Git', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// Eksekusi Git Pull ke Root Folder Website
async function executeGitPull() {
  const s = site.value
  if (!s || !s.id) return
  if (!gitModalRepo.value.trim()) {
    systemStore.showNotification('Peringatan', 'Isi alamat GitHub Repository terlebih dahulu', 'warning')
    return
  }

  gitModalIsPulling.value = true
  gitModalLogs.value = `[${new Date().toLocaleTimeString()}] Menghubungkan ke GitHub (${gitModalBranch.value || 'main'})...\n`
  gitModalLogs.value += `Target Folder: ${s.root_folder}\n`

  try {
    // Simpan dulu config terbaru
    await saveGitModalConfig()

    const res = await fetch(`/api/v1/websites/${s.id}/git/pull`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    if (res.ok) {
      s.git_last_pull = data.website?.git_last_pull || new Date().toISOString()
      s.git_last_commit = data.website?.git_last_commit || ''
      gitModalLogs.value += (data.logs || 'Git pull selesai tanpa output tambahan.') + '\n'
      gitModalLogs.value += `[${new Date().toLocaleTimeString()}] SUKSES: Kode terbaru berhasil ditarik ke root folder!\n`
      systemStore.showNotification(
        'Git Pull Berhasil',
        `Pembaruan kode GitHub terbaru berhasil ditarik ke "${s.name}"`,
        'success'
      )
      systemStore.triggerWebsitesRefresh()
    } else {
      gitModalLogs.value += '\n[ERROR] ' + (data.error || 'Gagal melakukan pull dari GitHub') + '\n'
      if (data.logs) {
        gitModalLogs.value += data.logs + '\n'
      }
      systemStore.showNotification('Git Pull Gagal', data.error || 'Gagal menarik kode dari GitHub', 'error')
    }
  } catch (e) {
    gitModalLogs.value += `\n[FATAL] Gagal menghubungi backend API: ${e.message}\n`
    systemStore.showNotification('Error', 'Gagal terhubung ke service backend', 'error')
  } finally {
    gitModalIsPulling.value = false
  }
}

onMounted(() => {
  initGitData()
})
</script>

<template>
  <form @submit.prevent="saveGitConfig" class="h-full flex flex-col justify-between text-xs select-none gap-3">
    <!-- Body Form & Log -->
    <div class="flex-1 overflow-y-auto space-y-3.5 pr-1">
      <!-- Folder Sasaran -->
      <div class="p-2.5 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between">
        <div class="min-w-0 pr-2">
          <span class="text-[10px] text-slate-500 block uppercase font-semibold">Root Folder Sasaran</span>
          <span class="font-mono text-slate-300 truncate block">{{ site.root_folder }}</span>
        </div>
        <button
          @click="openWebsiteFolder"
          class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-medium flex items-center gap-1 transition"
        >
          <FolderOpen class="w-3 h-3" />
          <span>Buka di Files</span>
        </button>
      </div>

      <!-- Input Repository & Branch -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="block font-medium text-slate-300 mb-1">
            GitHub Repository URL <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="gitModalRepo"
            type="text"
            placeholder="https://github.com/username/project.git"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
          />
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Branch</label>
          <input
            v-model="gitModalBranch"
            type="text"
            placeholder="main"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
          />
        </div>
      </div>

      <!-- Input PAT Token -->
      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="font-medium text-slate-300 flex items-center gap-1.5">
            <Key class="w-3.5 h-3.5 text-amber-400" />
            <span>GitHub Personal Access Token (PAT)</span>
          </label>
          <span class="text-[10px] text-slate-500">Wajib jika private repo</span>
        </div>
        <div class="relative">
          <input
            v-model="gitModalToken"
            :type="gitModalShowToken ? 'text' : 'password'"
            placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
            class="w-full pl-3 pr-10 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
          />
          <button
            type="button"
            @click="gitModalShowToken = !gitModalShowToken"
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200"
          >
            <EyeOff v-if="gitModalShowToken" class="w-4 h-4" />
            <Eye v-else class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Terminal Logs Git Pull -->
      <div>
        <div class="flex items-center justify-between text-slate-400 text-xs mb-1">
          <span class="font-mono flex items-center gap-1 text-[11px]">
            <Terminal class="w-3 h-3 text-sky-400" />
            <span>Status Log Tarikan Kode (Pull Output)</span>
          </span>
          <button @click="gitModalLogs = ''" class="hover:text-slate-200 text-[10px]">Bersihkan</button>
        </div>
        <pre class="bg-black/90 p-3 rounded-xl border border-slate-800 text-[10px] font-mono text-emerald-400 whitespace-pre-wrap max-h-40 overflow-y-auto leading-relaxed">{{ gitModalLogs }}</pre>
      </div>
    </div>

    <!-- Window Footer Actions -->
    <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
      <button
        @click="windowStore.closeWindow(win.id)"
        class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Tutup
      </button>

      <div class="flex items-center gap-2">
        <button
          @click="saveGitModalConfig"
          class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-sky-400 border border-sky-500/30 text-xs font-bold transition active:scale-95"
        >
          Simpan Pengaturan
        </button>
        <button
          @click="executeGitPull"
          :disabled="gitModalIsPulling || !gitModalRepo"
          class="px-4 py-2 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow-lg shadow-sky-500/20 active:scale-95 disabled:opacity-50"
        >
          <RotateCw v-if="gitModalIsPulling" class="w-3.5 h-3.5 animate-spin" />
          <GitPullRequest v-else class="w-3.5 h-3.5 stroke-[2.5]" />
          <span>{{ gitModalIsPulling ? 'Menarik Kode...' : 'Tarik Kode Sekarang (Pull)' }}</span>
        </button>
      </div>
    </div>
  </form>
</template>
