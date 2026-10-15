<script setup>
/**
 * RAZOS Website Edit Window Component
 * Created: 2026-10-13 (v2.6.0 - Standalone Window Architecture)
 * Updated: 2026-10-14 15:10:00 (v2.9.8 Hotfix: Cegah auto-submit form saat navigasi tab dengan type="button")
 * Deskripsi: Form Edit Konfigurasi website mandiri 5-tab.
 * Berjalan sebagai Jendela OS terpisah (tidak memblokir/mengganggu Website Manager).
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Pencil, Check, RotateCw, Globe, Server, Database, GitBranch,
  Sliders, Eye, EyeOff, FolderOpen, AlertTriangle
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

// State Tab & Form
const editActiveTab = ref('general')
const isSavingEdit = ref(false)

const editSiteName = ref('')
const editDomain = ref('')
const editPort = ref(8100)
const editIcon = ref('')
const editRuntime = ref('static')
const editRootFolder = ref('')
const editDbEngine = ref('')
const editDbHost = ref('127.0.0.1')
const editDbPort = ref(3306)
const editDbName = ref('')
const editDbUser = ref('')
const editDbPass = ref('')
const editGitRepo = ref('')
const editGitBranch = ref('main')
const editGitToken = ref('')
const editShowGitToken = ref(false)
const editMemLimit = ref('128M')
const editCpuShares = ref(512)
const editCommand = ref('')
const editEnvVars = ref('')

const memoryPresets = [
  { label: '32MB', val: '32M' },
  { label: '64MB', val: '64M' },
  { label: '128MB', val: '128M' },
  { label: '256MB', val: '256M' },
  { label: '512MB', val: '512M' }
]

function populateForm() {
  const s = site.value
  if (!s || !s.id) return
  editSiteName.value = s.name || ''
  editDomain.value = s.domain || ''
  editPort.value = s.port || 8100
  editIcon.value = s.icon || ''
  editRuntime.value = s.runtime || 'static'
  editRootFolder.value = s.root_folder || ''
  editDbEngine.value = s.db_engine || ''
  editDbHost.value = s.db_host || '127.0.0.1'
  editDbPort.value = s.db_port || 3306
  editDbName.value = s.db_name || ''
  editDbUser.value = s.db_user || ''
  editDbPass.value = s.db_pass || ''
  editGitRepo.value = s.git_repo || ''
  editGitBranch.value = s.git_branch || 'main'
  editGitToken.value = s.git_token || ''
  editMemLimit.value = s.mem_limit || '128M'
  editCpuShares.value = s.cpu_shares || 512
  editCommand.value = s.command || ''
  editEnvVars.value = s.env_vars || ''
}

async function confirmUpdateWebsite() {
  if (!editSiteName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama website wajib diisi', 'warning')
    return
  }

  isSavingEdit.value = true
  try {
    const payload = {
      name: editSiteName.value.trim(),
      domain: editDomain.value.trim(),
      port: Number(editPort.value),
      icon: editIcon.value.trim(),
      runtime: editRuntime.value,
      root_folder: editRootFolder.value.trim(),
      db_engine: editDbEngine.value,
      db_host: editDbHost.value,
      db_port: Number(editDbPort.value) || 0,
      db_name: editDbName.value.trim(),
      db_user: editDbUser.value.trim(),
      db_pass: editDbPass.value.trim(),
      git_repo: editGitRepo.value.trim(),
      git_branch: editGitBranch.value.trim() || 'main',
      git_token: editGitToken.value.trim(),
      mem_limit: editMemLimit.value.trim(),
      cpu_shares: Number(editCpuShares.value) || 512,
      command: editCommand.value.trim(),
      env_vars: editEnvVars.value.trim()
    }

    const res = await fetch(`/api/v1/websites/${site.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    const data = await res.json()
    if (res.ok) {
      systemStore.showNotification(
        'Sukses',
        `Konfigurasi "${payload.name}" berhasil diperbarui!`,
        'success'
      )
      // Memicu sinkronisasi data ke Website Manager & Desktop Icons
      systemStore.triggerWebsitesRefresh()
      windowStore.fetchWebsitesList()
      // Tutup jendela editor mandiri ini
      windowStore.closeWindow(props.win.id)
    } else {
      systemStore.showNotification('Gagal', data.error || 'Gagal memperbarui konfigurasi', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Terjadi kesalahan sistem saat menyimpan', 'error')
  } finally {
    isSavingEdit.value = false
  }
}

onMounted(() => {
  populateForm()
})
</script>

<template>
  <form @submit.prevent="confirmUpdateWebsite" class="h-full flex flex-col justify-between text-xs sm:text-sm select-none gap-3">
    <!-- Navigation Tabs -->
    <div class="flex items-center gap-1 overflow-x-auto bg-slate-950/80 p-1 rounded-xl border border-slate-800 text-[11px] shrink-0">
      <button
        v-for="t in [
          { id: 'general', label: '1. Identitas & Port' },
          { id: 'runtime', label: '2. Runtime & Folder' },
          { id: 'database', label: '3. Basis Data' },
          { id: 'git', label: '4. GitHub & PAT' },
          { id: 'advanced', label: '5. Resource & ENV' }
        ]"
        :key="t.id"
        type="button"
        @click="editActiveTab = t.id"
        class="px-3 py-1.5 rounded-lg font-medium transition shrink-0"
        :class="editActiveTab === t.id
          ? 'bg-amber-500 text-slate-950 font-bold shadow'
          : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- Scrollable Tab Content Body -->
    <div class="flex-1 overflow-y-auto space-y-3.5 pr-1 text-xs">
      <!-- TAB 1: General (Identitas, Domain, Port, Icon) -->
      <div v-if="editActiveTab === 'general'" class="space-y-3.5 animate-in fade-in">
        <div>
          <label class="block font-medium text-slate-300 mb-1">
            Nama Website / Project <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="editSiteName"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-amber-400"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-slate-300 mb-1">Domain Lokal / Hostname</label>
            <input
              v-model="editDomain"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
            />
          </div>

          <div>
            <label class="block font-medium text-slate-300 mb-1">
              Host Port Akses <span class="text-rose-400">*</span>
            </label>
            <input
              v-model.number="editPort"
              type="number"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-amber-400 font-mono text-xs font-bold outline-none focus:border-amber-400"
            />
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">URL / Path Ikon Aplikasi (Launcher)</label>
          <div class="flex items-center gap-2">
            <input
              v-model="editIcon"
              type="text"
              placeholder="/assets/icons/RAZIconWebsites.svg"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
            />
            <div class="w-8 h-8 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-center shrink-0">
              <img :src="editIcon || '/assets/icons/RAZIconWebsites.svg'" class="w-5 h-5 object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: Runtime Stack & Root Folder -->
      <div v-if="editActiveTab === 'runtime'" class="space-y-3.5 animate-in fade-in">
        <div>
          <label class="block font-medium text-slate-300 mb-1.5">Runtime Stack</label>
          <select
            v-model="editRuntime"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-amber-400"
          >
            <option value="static">1. Website Statis & Landing Page (HTML/Tailwind/Vite Nginx)</option>
            <option value="go-micro">2. Micro-Framework Web & REST API (Golang / PHP 8.2 Native)</option>
            <option value="flat-cms">3. CMS Ringan (Grav CMS / Kirby / Bludit)</option>
            <option value="internal-tools">4. Web Dashboard & Internal Tools (Homepage / Linktree)</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">
            Path Root Folder Dokumen Web <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="editRootFolder"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
          />
        </div>
      </div>

      <!-- TAB 3: Database -->
      <div v-if="editActiveTab === 'database'" class="space-y-3.5 animate-in fade-in">
        <div>
          <label class="block font-medium text-slate-300 mb-1.5">Engine Database Terhubung</label>
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              @click="editDbEngine = ''"
              class="py-2 rounded-xl border text-xs font-medium transition"
              :class="!editDbEngine ? 'bg-slate-800 border-slate-600 text-white' : 'bg-slate-950 border-slate-800 text-slate-400'"
            >
              Tanpa Database
            </button>
            <button
              type="button"
              @click="editDbEngine = 'sqlite'"
              class="py-2 rounded-xl border text-xs font-medium transition"
              :class="editDbEngine === 'sqlite' ? 'bg-emerald-500/20 border-emerald-400 text-emerald-300 font-bold' : 'bg-slate-950 border-slate-800 text-slate-400'"
            >
              SQLite (Zero RAM)
            </button>
            <button
              type="button"
              @click="editDbEngine = 'mariadb'"
              class="py-2 rounded-xl border text-xs font-medium transition"
              :class="editDbEngine === 'mariadb' ? 'bg-amber-500/20 border-amber-400 text-amber-300 font-bold' : 'bg-slate-950 border-slate-800 text-slate-400'"
            >
              MariaDB / MySQL
            </button>
          </div>
        </div>

        <div v-if="editDbEngine === 'sqlite'" class="p-3 rounded-xl bg-slate-950 border border-slate-800">
          <label class="block text-[11px] text-slate-400 mb-1">Nama File SQLite</label>
          <input v-model="editDbName" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" placeholder="database.sqlite" />
        </div>

        <div v-else-if="editDbEngine === 'mariadb'" class="space-y-3 p-3.5 rounded-xl bg-slate-950 border border-slate-800">
          <div class="grid grid-cols-2 gap-2.5">
            <div>
              <label class="block text-[11px] text-slate-400 mb-1">Host Database</label>
              <input v-model="editDbHost" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" />
            </div>
            <div>
              <label class="block text-[11px] text-slate-400 mb-1">Port</label>
              <input v-model.number="editDbPort" type="number" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2.5">
            <div>
              <label class="block text-[11px] text-slate-400 mb-1">Nama Database</label>
              <input v-model="editDbName" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" />
            </div>
            <div>
              <label class="block text-[11px] text-slate-400 mb-1">User DB</label>
              <input v-model="editDbUser" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" />
            </div>
          </div>
          <div>
            <label class="block text-[11px] text-slate-400 mb-1">Password DB</label>
            <input v-model="editDbPass" type="password" class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono" placeholder="Biarkan jika tidak diubah" />
          </div>
        </div>
      </div>

      <!-- TAB 4: Git & PAT -->
      <div v-if="editActiveTab === 'git'" class="space-y-3.5 animate-in fade-in">
        <div>
          <label class="block font-medium text-slate-300 mb-1">GitHub Repository URL</label>
          <input
            v-model="editGitRepo"
            type="text"
            placeholder="https://github.com/username/repo.git"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
          />
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Branch</label>
          <input
            v-model="editGitBranch"
            type="text"
            placeholder="main"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
          />
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Personal Access Token (PAT)</label>
          <div class="relative">
            <input
              v-model="editGitToken"
              :type="editShowGitToken ? 'text' : 'password'"
              placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
              class="w-full pl-3 pr-10 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
            />
            <button
              type="button"
              @click="editShowGitToken = !editShowGitToken"
              class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200"
            >
              <EyeOff v-if="editShowGitToken" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- TAB 5: Advanced (Resource Limits & Command) -->
      <div v-if="editActiveTab === 'advanced'" class="space-y-3.5 animate-in fade-in">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-slate-300 mb-1">Batas Memori RAM</label>
            <input
              v-model="editMemLimit"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
            />
            <div class="flex flex-wrap gap-1 mt-1.5">
              <button
                v-for="m in memoryPresets"
                :key="m.val"
                type="button"
                @click="editMemLimit = m.val"
                class="px-2 py-0.5 rounded-lg text-[10px] bg-slate-900 border border-slate-800 hover:border-slate-700 text-slate-300"
              >
                {{ m.val }}
              </button>
            </div>
          </div>

          <div>
            <label class="block font-medium text-slate-300 mb-1">CPU Shares</label>
            <input
              v-model.number="editCpuShares"
              type="number"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
            />
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Perintah Eksekusi Container (Command)</label>
          <input
            v-model="editCommand"
            type="text"
            placeholder="Contoh: node server.js"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
          />
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Environment Variables (ENV Vars)</label>
          <textarea
            v-model="editEnvVars"
            rows="3"
            placeholder="KEY=VALUE&#10;NODE_ENV=production"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 resize-none"
          ></textarea>
        </div>
      </div>
    </div>

    <!-- Footer Actions -->
    <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
      <button
        type="button"
        @click="windowStore.closeWindow(win.id)"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Batal
      </button>

      <button
        type="submit"
        :disabled="isSavingEdit"
        class="px-5 py-2 rounded-xl bg-gradient-to-r from-amber-500 to-amber-400 hover:from-amber-400 hover:to-amber-300 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow-lg shadow-amber-500/20 active:scale-95 disabled:opacity-50"
      >
        <RotateCw v-if="isSavingEdit" class="w-3.5 h-3.5 animate-spin" />
        <Check v-else class="w-3.5 h-3.5 stroke-[2.5]" />
        <span>{{ isSavingEdit ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
      </button>
    </div>
  </form>
</template>
