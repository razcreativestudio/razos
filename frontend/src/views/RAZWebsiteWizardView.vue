<script setup>
/**
 * RAZOS Website Wizard Window Component
 * Created: 2026-10-13 (v2.6.0 - Standalone Window Architecture)
 * Updated: 2026-10-14 15:10:00 (v2.9.8 Hotfix: Cegah auto-submit form pada wizard step buttons dengan type="button")
 * Deskripsi: Form Wizard mandiri untuk pembuatan website terisolasi berbasis Docker.
 * Berjalan sebagai Jendela OS terpisah (tidak memblokir/mengganggu Website Manager).
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Layers, CheckCircle2, ArrowRight, ShieldCheck, Database,
  Terminal, Server, Globe, Sparkles, FolderOpen, AlertTriangle
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ id: 'website-create' })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

// State Wizard
const wizardStep = ref(1)
const isSubmitting = ref(false)

// Form State
const formSiteName = ref('')
const formContainerName = ref('')
const formDomain = ref('')
const formPort = ref(8100)
const formRuntime = ref('static')
const formRootFolder = ref('/DATA/apps/websites/')
const formCommand = ref('')
const formMemLimit = ref('128M')
const formCpuShares = ref(512)
const formEnvVars = ref('')

// Database State
const formEnableDb = ref(false)
const formDbEngine = ref('sqlite')
const formDbHost = ref('127.0.0.1')
const formDbPort = ref(3306)
const formDbName = ref('database.sqlite')
const formDbUser = ref('raz_user')
const formDbPass = ref('')

// Git State
const formGitRepo = ref('')
const formGitBranch = ref('main')
const formGitToken = ref('')

// Presets RAM STB
const memoryPresets = [
  { label: 'Sangat Hemat (32MB)', val: '32M' },
  { label: 'Hemat (64MB)', val: '64M' },
  { label: 'Normal STB (128MB)', val: '128M' },
  { label: 'Sedang (256MB)', val: '256M' },
  { label: 'Maksimal STB (512MB)', val: '512M' }
]

// 4 Kategori Hosting STB-Friendly
const runtimeCategories = [
  {
    id: 'static',
    title: '1. Website Statis & Landing Page',
    badge: '15 - 40 MB RAM (Sempurna & Super Dingin)',
    desc: 'Untuk landing page, profil perusahaan, portofolio, brosur online.',
    tech: 'HTML, CSS, JS, Tailwind, Astro, Hugo, Vite',
    server: 'Nginx (Alpine) / Caddy',
    defaultCmd: '',
    defaultMem: '64M',
    ramBadgeClass: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30'
  },
  {
    id: 'go-micro',
    title: '2. Micro-Framework Web & REST API',
    badge: '20 - 90 MB RAM (Sangat Ringan & Efisien)',
    desc: 'Untuk API backend, webhook handler, bot bridge, web dinamis.',
    tech: 'Golang (Binary murni), Rust, PHP 8.2 Native, Node.js (Express)',
    server: 'Golang / PHP-FPM / Express',
    defaultCmd: '',
    defaultMem: '128M',
    ramBadgeClass: 'bg-sky-500/20 text-sky-300 border-sky-500/30'
  },
  {
    id: 'flat-cms',
    title: '3. CMS Ringan (Content Management)',
    badge: '35 - 90 MB RAM (Blogging & Artikel Tanpa MySQL)',
    desc: 'Blogging & website artikel tanpa service database yang memakan RAM.',
    tech: 'Grav CMS, Kirby, Bludit (Flat-File / Markdown)',
    server: 'PHP-FPM + Nginx Alpine',
    defaultCmd: '',
    defaultMem: '128M',
    ramBadgeClass: 'bg-amber-500/20 text-amber-300 border-amber-400/30'
  },
  {
    id: 'internal-tools',
    title: '4. Web Dashboard & Internal Tools',
    badge: '25 - 60 MB RAM (Portal Tim & Monitoring)',
    desc: 'Untuk dashboard utilitas tim, bookmark portal pribadi, status checker.',
    tech: 'Homepage, Flame, Dashy, Filebrowser, Uptime Kuma',
    server: 'Single Container Dedicated',
    defaultCmd: '',
    defaultMem: '128M',
    ramBadgeClass: 'bg-purple-500/20 text-purple-300 border-purple-400/30'
  }
]

// Sinkronisasi Nama Project ke Container & Folder Slug
function onSiteNameChange() {
  const slug = formSiteName.value
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9_-]/g, '-')
    .replace(/-+/g, '-')

  if (slug) {
    formContainerName.value = 'raz-web-' + slug
    formDomain.value = slug + '.local'
    formRootFolder.value = '/DATA/apps/websites/' + slug
  }
}

// Cari Port Kosong Otomatis
async function detectNextFreePort() {
  try {
    const res = await fetch('/api/v1/websites')
    if (res.ok) {
      const data = await res.json()
      const sites = data.websites || []
      const usedPorts = sites.map(s => s.port)
      let nextPort = 8100
      while (usedPorts.includes(nextPort)) {
        nextPort++
      }
      formPort.value = nextPort
    }
  } catch (e) {
    formPort.value = 8100
  }
}

function selectRuntimeCategory(cat) {
  formRuntime.value = cat.id
  formMemLimit.value = cat.defaultMem
  if (cat.defaultCmd) {
    formCommand.value = cat.defaultCmd
  }
}

async function confirmCreateWebsite() {
  if (!formSiteName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama website wajib diisi', 'warning')
    wizardStep.value = 1
    return
  }

  isSubmitting.value = true
  try {
    const payload = {
      name: formSiteName.value.trim(),
      domain: formDomain.value.trim(),
      port: Number(formPort.value),
      runtime: formRuntime.value,
      root_folder: formRootFolder.value.trim(),
      container_name: formContainerName.value.trim(),
      command: formCommand.value.trim(),
      mem_limit: formMemLimit.value.trim(),
      cpu_shares: Number(formCpuShares.value) || 512,
      env_vars: formEnvVars.value.trim(),
      db_engine: formEnableDb.value ? formDbEngine.value : '',
      db_host: formEnableDb.value ? formDbHost.value : '',
      db_port: formEnableDb.value ? Number(formDbPort.value) : 0,
      db_name: formEnableDb.value ? formDbName.value : '',
      db_user: formEnableDb.value ? formDbUser.value : '',
      db_pass: formEnableDb.value ? formDbPass.value : '',
      git_repo: formGitRepo.value.trim(),
      git_branch: formGitBranch.value.trim() || 'main',
      git_token: formGitToken.value.trim()
    }

    const res = await fetch('/api/v1/websites', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    const data = await res.json()
    if (res.ok) {
      systemStore.showNotification(
        'Berhasil',
        `Website "${payload.name}" berhasil dibuat di port :${payload.port}!`,
        'success'
      )
      // Memicu Website Manager memuat ulang daftar secara reaktif
      systemStore.triggerWebsitesRefresh()
      // Tutup jendela wizard mandiri ini
      windowStore.closeWindow(props.win?.id || 'website-create')
    } else {
      systemStore.showNotification('Gagal', data.error || 'Gagal membuat website', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Terjadi kesalahan sistem', 'error')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  detectNextFreePort()
})
</script>

<template>
  <form @submit.prevent="confirmCreateWebsite" class="h-full flex flex-col justify-between text-xs sm:text-sm select-none gap-3">
    <!-- Wizard Stepper Indicator -->
    <div class="grid grid-cols-4 gap-1.5 bg-slate-950/80 p-1.5 rounded-xl border border-slate-800/80 text-center text-[11px] shrink-0">
      <button
        v-for="(st, idx) in [
          { num: 1, label: '1. Identitas & Port' },
          { num: 2, label: '2. Runtime & Folder' },
          { num: 3, label: '3. Koneksi Database' },
          { num: 4, label: '4. Resource Limit' }
        ]"
        :key="st.num"
        type="button"
        @click="wizardStep = st.num"
        class="py-1.5 rounded-lg font-medium transition flex items-center justify-center gap-1"
        :class="wizardStep === st.num
          ? 'bg-sky-500 text-slate-950 font-bold shadow'
          : wizardStep > st.num
            ? 'text-sky-400 bg-sky-500/10'
            : 'text-slate-500 hover:text-slate-300'"
      >
        <span>{{ st.label }}</span>
      </button>
    </div>

    <!-- Wizard Scrollable Body -->
    <div class="flex-1 overflow-y-auto space-y-4 pr-1 text-xs">
      <!-- STEP 1: Identitas & Port -->
      <div v-if="wizardStep === 1" class="space-y-3.5">
        <div>
          <label class="block font-medium text-slate-300 mb-1">
            Nama Website / Project <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="formSiteName"
            @input="onSiteNameChange"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
            placeholder="Contoh: Landing Page Klien / Web Portofolio"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-slate-300 mb-1">Nama Container Docker</label>
            <input
              v-model="formContainerName"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="raz-web-myproject"
            />
          </div>

          <div>
            <label class="block font-medium text-slate-300 mb-1">Domain Reverse Proxy (Opsional)</label>
            <input
              v-model="formDomain"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="myproject.local / domain.com"
            />
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">
            Host Port Akses (Otomatis / Anti Tabrakan) <span class="text-rose-400">*</span>
          </label>
          <div class="flex items-center gap-2">
            <input
              v-model.number="formPort"
              type="number"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-amber-400 font-mono text-xs font-bold outline-none focus:border-sky-400"
              placeholder="8100"
            />
            <button
              type="button"
              @click="detectNextFreePort"
              class="px-3 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium shrink-0 transition"
            >
              Cari Port Bebas
            </button>
          </div>
          <span class="text-[10px] text-slate-400 mt-1 block">
            Port terisolasi untuk tiap project (rentang 8100 - 8250).
          </span>
        </div>
      </div>

      <!-- STEP 2: Runtime Stack & Root Folder -->
      <div v-if="wizardStep === 2" class="space-y-3.5">
        <div>
          <label class="block font-medium text-slate-300 mb-1.5">
            Pilih Kategori Runtime Stack (Dioptimasi untuk STB ARM64)
          </label>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
            <div
              v-for="cat in runtimeCategories"
              :key="cat.id"
              @click="selectRuntimeCategory(cat)"
              class="p-3 rounded-xl border cursor-pointer transition flex flex-col justify-between gap-1.5"
              :class="formRuntime === cat.id
                ? 'bg-sky-500/15 border-sky-400/60 shadow-lg'
                : 'bg-slate-950/70 border-slate-800/80 hover:border-slate-700'"
            >
              <div>
                <div class="flex items-center justify-between mb-1">
                  <h5 class="font-bold text-slate-100 text-xs">{{ cat.title }}</h5>
                  <span class="w-3.5 h-3.5 rounded-full border flex items-center justify-center"
                    :class="formRuntime === cat.id ? 'border-sky-400 bg-sky-400' : 'border-slate-600'">
                    <span v-if="formRuntime === cat.id" class="w-1.5 h-1.5 rounded-full bg-slate-950"></span>
                  </span>
                </div>
                <p class="text-[11px] text-slate-400 leading-snug">{{ cat.desc }}</p>
                <div class="mt-1 text-[10px] text-slate-300 font-mono">
                  Stack: {{ cat.tech }}
                </div>
              </div>
              <span class="text-[10px] font-mono px-2 py-0.5 rounded-md border w-fit" :class="cat.ramBadgeClass">
                {{ cat.badge }}
              </span>
            </div>
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">
            Path Root Folder Dokumen Web <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="formRootFolder"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
            placeholder="/DATA/apps/websites/nama-project"
          />
          <span class="text-[10px] text-slate-400 mt-1 block">
            Direktori fisik berkas website di MicroSD /DATA. File Manager terintegrasi penuh.
          </span>
        </div>
      </div>

      <!-- STEP 3: Koneksi Database -->
      <div v-if="wizardStep === 3" class="space-y-3.5">
        <label class="flex items-center gap-2.5 p-3 rounded-xl bg-slate-950 border border-slate-800 cursor-pointer">
          <input type="checkbox" v-model="formEnableDb" class="w-4 h-4 accent-sky-500 rounded" />
          <div>
            <span class="font-bold text-slate-200 block">Hubungkan Database untuk Website Ini</span>
            <span class="text-[11px] text-slate-400">Pilih SQLite (Rekomendasi STB) atau MariaDB / PostgreSQL.</span>
          </div>
        </label>

        <div v-if="formEnableDb" class="space-y-3 p-3.5 rounded-xl bg-slate-950/60 border border-slate-800">
          <div>
            <label class="block font-medium text-slate-300 mb-1.5">Pilihan Engine Database</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="formDbEngine = 'sqlite'"
                class="p-2.5 rounded-xl border text-left transition flex flex-col justify-between"
                :class="formDbEngine === 'sqlite' ? 'bg-emerald-500/15 border-emerald-400 text-emerald-300' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                <span class="font-bold text-xs">SQLite (Berbasis Berkas)</span>
                <span class="text-[10px] opacity-80 mt-1">0 MB Service RAM (Paling Ringan & Aman)</span>
              </button>

              <button
                type="button"
                @click="formDbEngine = 'mariadb'"
                class="p-2.5 rounded-xl border text-left transition flex flex-col justify-between"
                :class="formDbEngine === 'mariadb' ? 'bg-amber-500/15 border-amber-400 text-amber-300' : 'bg-slate-950 border-slate-800 text-slate-400'"
              >
                <span class="font-bold text-xs">MariaDB / MySQL</span>
                <span class="text-[10px] opacity-80 mt-1">Port :3306 (Adminer GUI)</span>
              </button>
            </div>
          </div>

          <div v-if="formDbEngine !== 'sqlite'" class="grid grid-cols-2 gap-2.5">
            <div>
              <label class="block text-[11px] font-medium text-slate-400 mb-1">Nama Database</label>
              <input v-model="formDbName" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-800 text-xs text-slate-100 font-mono" placeholder="db_myproject" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-slate-400 mb-1">Username DB</label>
              <input v-model="formDbUser" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-800 text-xs text-slate-100 font-mono" placeholder="raz_user" />
            </div>
          </div>

          <div v-else>
            <label class="block text-[11px] font-medium text-slate-400 mb-1">Nama Berkas Database SQLite</label>
            <input v-model="formDbName" type="text" class="w-full px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-800 text-xs text-slate-100 font-mono" placeholder="database.sqlite" />
          </div>
        </div>
      </div>

      <!-- STEP 4: Resource Limits & Git Pull -->
      <div v-if="wizardStep === 4" class="space-y-3.5">
        <div>
          <label class="block font-medium text-slate-300 mb-1.5">Batas Memori RAM Container</label>
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
            <button
              v-for="p in memoryPresets"
              :key="p.val"
              type="button"
              @click="formMemLimit = p.val"
              class="p-2 rounded-xl border text-center transition"
              :class="formMemLimit === p.val ? 'bg-sky-500/20 text-sky-300 border-sky-400 font-bold' : 'bg-slate-950 border-slate-800 text-slate-400'"
            >
              <span class="block text-xs">{{ p.val }}</span>
              <span class="text-[9px] opacity-75 block truncate">{{ p.label }}</span>
            </button>
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">GitHub Repository (Opsional)</label>
          <input
            v-model="formGitRepo"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
            placeholder="https://github.com/username/repository.git"
          />
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Environment Variables (ENV Vars)</label>
          <textarea
            v-model="formEnvVars"
            rows="2"
            placeholder="NODE_ENV=production&#10;APP_KEY=secret"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-emerald-400 font-mono text-xs outline-none focus:border-sky-400 resize-none"
          ></textarea>
        </div>
      </div>
    </div>

    <!-- Wizard Navigation Footer -->
    <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
      <button
        v-if="wizardStep > 1"
        type="button"
        @click="wizardStep--"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Kembali
      </button>
      <div v-else></div>

      <div class="flex items-center gap-2">
        <button
          v-if="wizardStep < 4"
          type="button"
          @click="wizardStep++"
          class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow"
        >
          <span>Lanjut</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>

        <button
          v-else
          type="button"
          :disabled="isSubmitting"
          @click="confirmCreateWebsite"
          class="px-5 py-2 rounded-xl bg-gradient-to-r from-sky-500 to-emerald-500 hover:from-sky-400 hover:to-emerald-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow-lg active:scale-95 disabled:opacity-50"
        >
          <CheckCircle2 class="w-4 h-4" />
          <span>{{ isSubmitting ? 'Memproses Container...' : 'Deploy & Pasang Website' }}</span>
        </button>
      </div>
    </div>
  </form>
</template>
