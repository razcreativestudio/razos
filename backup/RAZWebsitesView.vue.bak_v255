<script setup>
/**
 * RAZOS Website Manager Suite
 * Updated: 2026-10-13
 * Deskripsi: Platform pengelolaan multi-website terisolasi berbasis Docker Container.
 * Fitur:
 * 1. Root folder terisolasi di /DATA/apps/websites/<slug>
 * 2. Port host independen (8100-8250) anti-tabrakan
 * 3. Dukungan runtime stack: HTML Statis (Nginx), PHP 8.2 + DB, Node.js, Next.js, Custom
 * 4. Integrasi Database Manager (Adminer / HeidiSQL web-style) dengan auto-fill koneksi
 * 5. Link icon kustom untuk launcher Desktop dan Start Menu
 * 6. Resource limits: Memory limit (RAM STB-optimized) & CPU shares
 * 7. Custom container command & Environment variables
 * 8. Zero native alert/confirm - 100% Custom Modals Glassmorphism Slate-900
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Globe, Plus, ExternalLink, Inbox, Trash2, X, Server, Layers,
  Play, Square, RotateCw, FolderOpen, Database, Cpu, HardDrive,
  Copy, Check, ShieldCheck, Search, Filter, Terminal, Sparkles,
  Info, AlertTriangle, ArrowRight, CheckCircle2,
  GitBranch, GitPullRequest, Key, Eye, EyeOff, Download, RefreshCw,
  Pencil
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const websites = ref([])
const databases = ref([])
const isLoading = ref(false)
const searchQuery = ref('')
const selectedRuntimeFilter = ref('all')

// State Modal Wizard Pembuatan Website
const showCreateModal = ref(false)
const wizardStep = ref(1) // 1: Identitas & Port, 2: Runtime & Folder, 3: Database, 4: Resource & Advanced
const isSubmitting = ref(false)

// Form Wizard State
const formSiteName = ref('')
const formContainerName = ref('')
const formDomain = ref('')
const formPort = ref(8101)
const formIcon = ref('')
const formRuntime = ref('static') // 'static' | 'php' | 'nodejs' | 'nextjs' | 'custom'
const formRootFolder = ref('')
const formCommand = ref('')
const formMemLimit = ref('128M')
const formCpuShares = ref(512)
const formEnvVars = ref('')

// Form Database State (SQLite Direkomendasikan untuk STB ARM64 - 0 MB Service Overhead)
const formEnableDb = ref(false)
const formSelectedDbId = ref('')
const formDbEngine = ref('sqlite')
const formDbHost = ref('127.0.0.1')
const formDbPort = ref(3306)
const formDbName = ref('database.sqlite')
const formDbUser = ref('raz_user')
const formDbPass = ref('')

// Form Git State di Wizard
const formGitRepo = ref('')
const formGitBranch = ref('main')
const formGitToken = ref('')
const formShowGitToken = ref(false)

// State Modal Pengaturan Git & Git Pull
const showGitModal = ref(false)
const gitModalSite = ref(null)
const gitModalRepo = ref('')
const gitModalBranch = ref('main')
const gitModalToken = ref('')
const gitModalShowToken = ref(false)
const gitModalIsPulling = ref(false)
const gitModalLogs = ref('')
const pullingSiteId = ref(null)

// State Modal Edit Website
const showEditModal = ref(false)
const editSiteId = ref(null)
const editSiteName = ref('')
const editDomain = ref('')
const editPort = ref(8101)
const editIcon = ref('')
const editRuntime = ref('static')
const editRootFolder = ref('')
const editCommand = ref('')
const editMemLimit = ref('128M')
const editCpuShares = ref(512)
const editEnvVars = ref('')

const editEnableDb = ref(false)
const editSelectedDbId = ref('')
const editDbEngine = ref('sqlite')
const editDbHost = ref('127.0.0.1')
const editDbPort = ref(3306)
const editDbName = ref('database.sqlite')
const editDbUser = ref('raz_user')
const editDbPass = ref('')

const editGitRepo = ref('')
const editGitBranch = ref('main')
const editGitToken = ref('')
const editShowGitToken = ref(false)

const editActiveTab = ref('general')
const isSavingEdit = ref(false)

// State Modal Hapus Website
const showDeleteModal = ref(false)
const siteToDelete = ref(null)
const purgeRootFolder = ref(false)

// State Modal Detail Website
const showDetailModal = ref(false)
const selectedSiteDetail = ref(null)
const copiedField = ref(null)

// Preset Icon Pilihan Cepat (Lengkap Sesuai Ekosistem STB Hosting)
const presetIcons = [
  { name: 'Nginx / Statis', url: '/assets/icons/RAZIconWebsites.svg' },
  { name: 'Go (Golang)', url: 'https://cdn.simpleicons.org/go/00ADD8' },
  { name: 'PHP', url: 'https://cdn.simpleicons.org/php/777BB4' },
  { name: 'Node.js', url: 'https://cdn.simpleicons.org/nodedotjs/5FA04E' },
  { name: 'Python', url: 'https://cdn.simpleicons.org/python/3776AB' },
  { name: 'Flat-File CMS', url: 'https://cdn.simpleicons.org/markdown/white' },
  { name: 'WordPress', url: 'https://cdn.simpleicons.org/wordpress/21759B' },
  { name: 'Linktree / Portal', url: 'https://cdn.simpleicons.org/linktree/43E660' },
  { name: 'Filebrowser', url: 'https://cdn.simpleicons.org/files/white' },
  { name: 'SQLite', url: 'https://cdn.simpleicons.org/sqlite/003B57' },
  { name: 'TailwindCSS', url: 'https://cdn.simpleicons.org/tailwindcss/06B6D4' }
]

// Preset Memory Limits Direkomendasikan untuk STB ARM64 (RAM 1.8 GiB)
const memoryPresets = [
  { label: '32M (Ultra Ringan - Statis/Portal)', val: '32M' },
  { label: '64M (Ringan - Go/Flat-CMS)', val: '64M' },
  { label: '128M (Rekomendasi STB - PHP/API)', val: '128M' },
  { label: '256M (Batas Aman STB)', val: '256M' },
  { label: '512M (Batas Puncak)', val: '512M' }
]

// 4 Kategori Hosting Ideal untuk Server STB Amlogic S905X ARM64
const stbRuntimeCategories = [
  {
    title: '1. Website Statis & Landing Page',
    badge: '🔥 Paling Sempurna (15 – 40 MB RAM)',
    badgeColor: 'text-emerald-400 bg-emerald-500/10 border-emerald-500/30',
    desc: 'Bebas database. Server STB Anda akan sangat dingin, cepat, dan responsif.',
    runtimes: [
      {
        id: 'static',
        name: 'HTML & CSS Statis (Nginx Alpine)',
        stack: 'HTML5, CSS, JS, Tailwind, Astro/Hugo/Vite static build',
        ram: '15 – 40 MB',
        icon: '/assets/icons/RAZIconWebsites.svg',
        desc: 'Landing page studio/agensi, web portofolio, brosur online, company profile.'
      }
    ]
  },
  {
    title: '2. Micro-Framework Web & REST API',
    badge: '⚡ Ringan – Sedang (20 – 120 MB RAM)',
    badgeColor: 'text-sky-400 bg-sky-500/10 border-sky-500/30',
    desc: 'Untuk web dinamis, form input, REST API, dashboard internal, atau bot bridge.',
    runtimes: [
      {
        id: 'golang',
        name: 'Go (Golang) Micro-Service',
        stack: 'Go Native Binary (Goroutines, net/http)',
        ram: '20 – 40 MB',
        icon: 'https://cdn.simpleicons.org/go/00ADD8',
        desc: 'Paling direkomendasikan untuk STB! Native kecepatan tinggi, memori sangat hemat.'
      },
      {
        id: 'php',
        name: 'PHP Native + SQLite / MariaDB',
        stack: 'PHP 8.2 FPM / Apache Alpine',
        ram: '40 – 80 MB',
        icon: 'https://cdn.simpleicons.org/php/777BB4',
        desc: 'Web dinamis PHP standar tanpa framework monolitik raksasa yang boros.'
      },
      {
        id: 'nodejs',
        name: 'Node.js Express / Fastify Light',
        stack: 'Node 18/20 Alpine (Express / Fastify)',
        ram: '70 – 120 MB',
        icon: 'https://cdn.simpleicons.org/nodedotjs/5FA04E',
        desc: 'REST API atau backend JavaScript ringan terisolasi dalam container.'
      },
      {
        id: 'python',
        name: 'Python FastAPI / Flask Micro-service',
        stack: 'Python 3.11 Alpine (FastAPI / Flask)',
        ram: '80 – 150 MB',
        icon: 'https://cdn.simpleicons.org/python/3776AB',
        desc: 'Micro-service Python teroptimasi untuk API data dan otomatisasi tugas.'
      }
    ]
  },
  {
    title: '3. CMS Ringan (Content Management System)',
    badge: '📦 Ramah STB (35 – 250 MB RAM)',
    badgeColor: 'text-indigo-400 bg-indigo-500/10 border-indigo-500/30',
    desc: 'Untuk website yang punya dashboard admin untuk artikel atau publikasi konten.',
    runtimes: [
      {
        id: 'flatcms',
        name: 'Flat-File CMS (Grav / Kirby Style)',
        stack: 'Markdown Based (Tanpa Database SQL)',
        ram: '35 – 50 MB',
        icon: 'https://cdn.simpleicons.org/markdown/white',
        desc: 'Sangat cepat dan ramah RAM STB karena tidak membutuhkan service MariaDB/MySQL.'
      },
      {
        id: 'wordpress',
        name: 'WordPress Optimized (Light)',
        stack: 'WordPress Core + MariaDB Alpine Caching',
        ram: '150 – 250 MB',
        icon: 'https://cdn.simpleicons.org/wordpress/21759B',
        desc: 'WordPress dengan MariaDB Alpine (disarankan tanpa page-builder berat seperti Elementor).'
      }
    ]
  },
  {
    title: '4. Web Dashboard & Internal Tools (Self-Hosted)',
    badge: '💡 Self-Hosted Utility ( < 40 MB RAM)',
    badgeColor: 'text-amber-400 bg-amber-500/10 border-amber-500/30',
    desc: 'Website utility tim, portal bookmark, monitoring, dan sharing file ringan.',
    runtimes: [
      {
        id: 'dashboard',
        name: 'Linktree Pribadi / Bookmark Portal',
        stack: 'Homepage / Flame / Dashy Style Web',
        ram: '< 40 MB',
        icon: 'https://cdn.simpleicons.org/linktree/43E660',
        desc: 'Portal kumpulan tautan pribadi, status server, atau bookmark portal responsif.'
      },
      {
        id: 'filebrowser',
        name: 'File Sharing Ringan (Filebrowser)',
        stack: 'Go Native Filebrowser Service',
        ram: '20 – 35 MB',
        icon: 'https://cdn.simpleicons.org/files/white',
        desc: 'Alternatif Google Drive ultra ringan untuk manajemen dan unduh berkas di STB.'
      }
    ]
  }
]

async function fetchWebsites() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/websites')
    if (res.ok) {
      websites.value = await res.json()
      // Sinkronkan juga ke desktop store jika fungsi tersedia
      if (typeof windowStore.fetchWebsitesList === 'function') {
        windowStore.fetchWebsitesList()
      }
    }
  } catch (err) {
    console.error('Gagal mengambil daftar website:', err)
  }
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

// Otomatis update port bebas saat membuka wizard
async function suggestFreePort() {
  const usedPorts = new Set(websites.value.map(w => w.port))
  for (let p = 8100; p <= 8250; p++) {
    if (!usedPorts.has(p)) {
      formPort.value = p
      break
    }
  }
}

// Buka Dialog Wizard Pembuatan Website Baru
function openCreateDialog() {
  wizardStep.value = 1
  formSiteName.value = ''
  formContainerName.value = ''
  formDomain.value = ''
  formIcon.value = '/assets/icons/RAZIconWebsites.svg'
  formRuntime.value = 'static'
  formRootFolder.value = ''
  formCommand.value = ''
  formMemLimit.value = '128M'
  formCpuShares.value = 512
  formEnvVars.value = ''

  formEnableDb.value = false
  formSelectedDbId.value = ''
  formDbEngine.value = 'mariadb'
  formDbHost.value = '127.0.0.1'
  formDbPort.value = 3306
  formDbName.value = ''
  formDbUser.value = 'raz_user'
  formDbPass.value = ''

  formGitRepo.value = ''
  formGitBranch.value = 'main'
  formGitToken.value = ''
  formShowGitToken.value = false

  suggestFreePort()
  fetchDatabases()
  showCreateModal.value = true
}

// Auto-generate container name & root folder saat user mengetik nama website
function onSiteNameChange() {
  const slug = formSiteName.value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '')

  if (!formContainerName.value || formContainerName.value.startsWith('site-')) {
    formContainerName.value = slug ? `site-${slug}-web` : ''
  }
  if (!formDomain.value || formDomain.value.endsWith('.local')) {
    formDomain.value = slug ? `${slug}.razos.local` : ''
  }
  if (!formRootFolder.value || formRootFolder.value.includes('/websites/')) {
    formRootFolder.value = slug ? `/DATA/apps/websites/${slug}` : ''
  }
  if (formEnableDb.value && (!formDbName.value || formDbName.value.startsWith('db_'))) {
    formDbName.value = slug ? `db_${slug.replace(/-/g, '_')}` : ''
  }
}

// Reaksi perubahan pilihan runtime stack (Wizard)
function onRuntimeChange(rt) {
  formRuntime.value = rt
  switch (rt) {
    case 'golang':
      formIcon.value = 'https://cdn.simpleicons.org/go/00ADD8'
      formMemLimit.value = '64M'
      formCommand.value = './main'
      formEnableDb.value = false
      break
    case 'php':
      formIcon.value = 'https://cdn.simpleicons.org/php/777BB4'
      formMemLimit.value = '128M'
      formCommand.value = ''
      formEnableDb.value = true
      formDbEngine.value = 'sqlite'
      formDbName.value = 'database.sqlite'
      break
    case 'nodejs':
      formIcon.value = 'https://cdn.simpleicons.org/nodedotjs/5FA04E'
      formMemLimit.value = '128M'
      formCommand.value = 'node server.js'
      formEnableDb.value = false
      break
    case 'python':
      formIcon.value = 'https://cdn.simpleicons.org/python/3776AB'
      formMemLimit.value = '128M'
      formCommand.value = 'python app.py'
      formEnableDb.value = false
      break
    case 'flatcms':
      formIcon.value = 'https://cdn.simpleicons.org/markdown/white'
      formMemLimit.value = '64M'
      formCommand.value = ''
      formEnableDb.value = false
      break
    case 'wordpress':
      formIcon.value = 'https://cdn.simpleicons.org/wordpress/21759B'
      formMemLimit.value = '256M'
      formCommand.value = ''
      formEnableDb.value = true
      formDbEngine.value = 'mariadb'
      formDbPort.value = 3306
      break
    case 'dashboard':
      formIcon.value = 'https://cdn.simpleicons.org/linktree/43E660'
      formMemLimit.value = '32M'
      formCommand.value = ''
      formEnableDb.value = false
      break
    case 'filebrowser':
      formIcon.value = 'https://cdn.simpleicons.org/files/white'
      formMemLimit.value = '32M'
      formCommand.value = ''
      formEnableDb.value = false
      break
    case 'static':
    default:
      formIcon.value = '/assets/icons/RAZIconWebsites.svg'
      formMemLimit.value = '32M'
      formCommand.value = ''
      formEnableDb.value = false
      break
  }
}

// Reaksi perubahan pilihan runtime stack (Edit Modal)
function onEditRuntimeChange(rt) {
  editRuntime.value = rt
  switch (rt) {
    case 'golang':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/go/00ADD8'
      if (!editCommand.value) editCommand.value = './main'
      break
    case 'php':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/php/777BB4'
      break
    case 'nodejs':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/nodedotjs/5FA04E'
      if (!editCommand.value) editCommand.value = 'node server.js'
      break
    case 'python':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/python/3776AB'
      if (!editCommand.value) editCommand.value = 'python app.py'
      break
    case 'flatcms':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/markdown/white'
      break
    case 'wordpress':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/wordpress/21759B'
      break
    case 'dashboard':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/linktree/43E660'
      break
    case 'filebrowser':
      if (!editIcon.value || editIcon.value.includes('RAZIconWebsites')) editIcon.value = 'https://cdn.simpleicons.org/files/white'
      break
    case 'static':
    default:
      if (!editIcon.value || editIcon.value.includes('simpleicons')) editIcon.value = '/assets/icons/RAZIconWebsites.svg'
      break
  }
}

// Auto-fill dari dropdown database yang sudah terdaftar di RAZOS
function onRegisteredDbSelect() {
  if (!formSelectedDbId.value) return
  const db = databases.value.find(d => d.id === parseInt(formSelectedDbId.value))
  if (db) {
    formDbEngine.value = db.engine || 'mariadb'
    formDbHost.value = db.host || '127.0.0.1'
    formDbPort.value = db.port || 3306
    formDbName.value = db.name
    formDbUser.value = db.user || 'raz_user'
  }
}

// Eksekusi Pembuatan Website
async function confirmCreateWebsite() {
  if (!formSiteName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama website wajib diisi', 'error')
    wizardStep.value = 1
    return
  }

  isSubmitting.value = true
  const payload = {
    name: formSiteName.value.trim(),
    container_name: formContainerName.value.trim(),
    domain: formDomain.value.trim(),
    port: parseInt(formPort.value) || 8101,
    icon: formIcon.value.trim(),
    runtime: formRuntime.value,
    root_folder: formRootFolder.value.trim(),
    command: formCommand.value.trim(),
    memory_limit: formMemLimit.value,
    cpu_shares: parseInt(formCpuShares.value) || 512,
    env_vars: formEnvVars.value.trim(),
    database_id: formEnableDb.value && formSelectedDbId.value ? parseInt(formSelectedDbId.value) : null,
    db_engine: formEnableDb.value ? formDbEngine.value : '',
    db_host: formEnableDb.value ? formDbHost.value : '',
    db_port: formEnableDb.value ? parseInt(formDbPort.value) : 3306,
    db_name: formEnableDb.value ? formDbName.value : '',
    db_user: formEnableDb.value ? formDbUser.value : '',
    db_password: formEnableDb.value ? formDbPass.value : '',
    git_repo: formGitRepo.value.trim(),
    git_branch: formGitBranch.value.trim() || 'main',
    git_token: formGitToken.value.trim()
  }

  try {
    const res = await fetch('/api/v1/websites', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      const created = await res.json()
      systemStore.showNotification(
        'Sukses',
        `Website "${created.name}" berhasil dibuat pada port ${created.port}!`,
        'success'
      )
      showCreateModal.value = false
      fetchWebsites()
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal', err.error || 'Gagal membuat website', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi backend', 'error')
  }
  isSubmitting.value = false
}

// Buka Modal Pengaturan Git & PAT
function openGitModal(site) {
  gitModalSite.value = site
  gitModalRepo.value = site.git_repo || ''
  gitModalBranch.value = site.git_branch || 'main'
  gitModalToken.value = site.git_token || ''
  gitModalShowToken.value = false
  gitModalLogs.value = ''
  showGitModal.value = true
}

// Simpan Konfigurasi Git untuk Website
async function saveGitModalConfig() {
  if (!gitModalSite.value) return
  const siteId = gitModalSite.value.id
  try {
    const res = await fetch(`/api/v1/websites/${siteId}/git/config`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        git_repo: gitModalRepo.value.trim(),
        git_branch: gitModalBranch.value.trim() || 'main',
        git_token: gitModalToken.value.trim()
      })
    })
    if (res.ok) {
      const data = await res.json()
      gitModalSite.value.git_repo = data.website.git_repo
      gitModalSite.value.git_branch = data.website.git_branch
      gitModalSite.value.git_token = data.website.git_token
      systemStore.showNotification('Sukses', 'Konfigurasi GitHub & PAT berhasil disimpan', 'success')
      fetchWebsites()
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal', err.error || 'Gagal menyimpan konfigurasi Git', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// Eksekusi Git Pull ke Root Folder Website
async function executeGitPull(site) {
  const target = site || gitModalSite.value
  if (!target) return
  if (!target.git_repo) {
    // Jika belum ada konfigurasi repo, buka dialog pengaturan
    openGitModal(target)
    return
  }

  pullingSiteId.value = target.id
  gitModalIsPulling.value = true
  gitModalLogs.value = `[${new Date().toLocaleTimeString()}] Menghubungkan ke GitHub (${target.git_branch || 'main'})...\n`
  gitModalLogs.value += `Target Folder: ${target.root_folder}\n`

  try {
    const res = await fetch(`/api/v1/websites/${target.id}/git/pull`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    if (res.ok) {
      target.git_last_pull = data.website?.git_last_pull || new Date().toISOString()
      target.git_last_commit = data.website?.git_last_commit || ''
      gitModalLogs.value += (data.logs || 'Git pull selesai tanpa output tambahan.') + '\n'
      gitModalLogs.value += `[${new Date().toLocaleTimeString()}] SUKSES: Kode terbaru berhasil ditarik ke root folder!\n`
      systemStore.showNotification(
        'Git Pull Berhasil',
        `Pembaruan kode GitHub terbaru berhasil ditarik ke "${target.name}"`,
        'success'
      )
      fetchWebsites()
    } else {
      gitModalLogs.value += '\n[ERROR] ' + (data.error || 'Gagal melakukan pull dari GitHub') + '\n'
      if (data.logs) {
        gitModalLogs.value += data.logs + '\n'
      }
      systemStore.showNotification(
        'Git Pull Gagal',
        data.error || 'Terjadi kesalahan saat menarik kode dari GitHub',
        'error'
      )
    }
  } catch (e) {
    gitModalLogs.value += '\n[FATAL] Gagal menghubungi backend: ' + e.message + '\n'
    systemStore.showNotification('Error', 'Gagal menghubungi backend', 'error')
  } finally {
    pullingSiteId.value = null
    gitModalIsPulling.value = false
  }
}

// Toggle Status Container Website (Start / Stop / Restart)
async function toggleWebsiteStatus(site, action) {
  try {
    const res = await fetch(`/api/v1/websites/${site.id}/action`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action })
    })
    if (res.ok) {
      const updated = await res.json()
      site.status = updated.status
      systemStore.showNotification(
        'Status Diperbarui',
        `Website "${site.name}" sekarang berstatus ${site.status.toUpperCase()}`,
        'success'
      )
    }
  } catch (e) {
    systemStore.showNotification('Gagal', `Gagal mengubah status website: ${e.message}`, 'error')
  }
}

// Buka Modal Edit Konfigurasi Website
function openEditDialog(site) {
  editSiteId.value = site.id
  editSiteName.value = site.name || ''
  editDomain.value = site.domain || ''
  editPort.value = site.port || 8101
  editIcon.value = site.icon || '/assets/icons/RAZIconWebsites.svg'
  editRuntime.value = site.runtime || 'static'
  editRootFolder.value = site.root_folder || ''
  editCommand.value = site.command || ''
  editMemLimit.value = site.memory_limit || '128M'
  editCpuShares.value = site.cpu_shares || 512
  editEnvVars.value = site.env_vars || ''

  editEnableDb.value = !!site.db_name
  editSelectedDbId.value = site.database_id ? String(site.database_id) : ''
  editDbEngine.value = site.db_engine || 'mariadb'
  editDbHost.value = site.db_host || '127.0.0.1'
  editDbPort.value = site.db_port || 3306
  editDbName.value = site.db_name || ''
  editDbUser.value = site.db_user || 'raz_user'
  editDbPass.value = site.db_password || ''

  editGitRepo.value = site.git_repo || ''
  editGitBranch.value = site.git_branch || 'main'
  editGitToken.value = site.git_token || ''
  editShowGitToken.value = false

  editActiveTab.value = 'general'
  fetchDatabases()
  showEditModal.value = true
}

// Auto-fill dari database terdaftar saat form edit
function onEditRegisteredDbSelect() {
  if (!editSelectedDbId.value) return
  const db = databases.value.find(d => d.id === parseInt(editSelectedDbId.value))
  if (db) {
    editDbEngine.value = db.engine || 'mariadb'
    editDbHost.value = db.host || '127.0.0.1'
    editDbPort.value = db.port || 3306
    editDbName.value = db.name
    editDbUser.value = db.user || 'raz_user'
  }
}

// Simpan Perubahan Edit Konfigurasi Website
async function confirmUpdateWebsite() {
  if (!editSiteName.value.trim()) {
    systemStore.showNotification('Peringatan', 'Nama website wajib diisi', 'error')
    editActiveTab.value = 'general'
    return
  }

  isSavingEdit.value = true
  const payload = {
    name: editSiteName.value.trim(),
    domain: editDomain.value.trim(),
    port: parseInt(editPort.value) || 8101,
    icon: editIcon.value.trim(),
    runtime: editRuntime.value,
    root_folder: editRootFolder.value.trim(),
    command: editCommand.value.trim(),
    memory_limit: editMemLimit.value,
    cpu_shares: parseInt(editCpuShares.value) || 512,
    env_vars: editEnvVars.value.trim(),
    database_id: editEnableDb.value && editSelectedDbId.value ? parseInt(editSelectedDbId.value) : null,
    db_engine: editEnableDb.value ? editDbEngine.value : '',
    db_host: editEnableDb.value ? editDbHost.value : '',
    db_port: editEnableDb.value ? parseInt(editDbPort.value) : 3306,
    db_name: editEnableDb.value ? editDbName.value : '',
    db_user: editEnableDb.value ? editDbUser.value : '',
    db_password: editEnableDb.value ? editDbPass.value : '',
    git_repo: editGitRepo.value.trim(),
    git_branch: editGitBranch.value.trim() || 'main',
    git_token: editGitToken.value.trim()
  }

  try {
    const res = await fetch(`/api/v1/websites/${editSiteId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (res.ok) {
      const updated = await res.json()
      systemStore.showNotification(
        'Sukses',
        `Konfigurasi website "${updated.name}" berhasil diperbarui!`,
        'success'
      )
      showEditModal.value = false
      fetchWebsites()
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal', err.error || 'Gagal memperbarui website', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi backend', 'error')
  }
  isSavingEdit.value = false
}

// Buka dialog konfirmasi hapus website
function promptDeleteWebsite(site) {
  siteToDelete.value = site
  purgeRootFolder.value = false
  showDeleteModal.value = true
}

// Konfirmasi hapus website
async function executeDeleteWebsite() {
  if (!siteToDelete.value) return
  const id = siteToDelete.value.id
  const name = siteToDelete.value.name
  try {
    const res = await fetch(`/api/v1/websites/${id}?purge=${purgeRootFolder.value}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Website "${name}" telah dihapus`, 'success')
      showDeleteModal.value = false
      fetchWebsites()
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghapus website', 'error')
  }
}

// Buka Root Folder langsung di File Manager RAZOS
function openWebsiteFolder(site) {
  windowStore.openApp('file-manager')
  systemStore.showNotification(
    'Membuka Berkas',
    `Direktori website: ${site.root_folder}`,
    'info'
  )
}

// Buka Web Database GUI (Adminer / HeidiSQL Web-style)
function openDatabaseManager(site) {
  windowStore.openApp('databases')
  const host = window.location.hostname
  window.open(`http://${host}:8088`, '_blank')
}

// Tampilkan Detail Modal
function openDetailDialog(site) {
  selectedSiteDetail.value = site
  showDetailModal.value = true
}

// Salin teks ke clipboard
function copyToClipboard(text, fieldName) {
  navigator.clipboard.writeText(text)
  copiedField.value = fieldName
  setTimeout(() => {
    copiedField.value = null
  }, 2000)
}

// Filter Website Berdasarkan Pencarian & Runtime
const filteredWebsites = computed(() => {
  return websites.value.filter(site => {
    const matchesSearch = !searchQuery.value ||
      site.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      site.container_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      site.domain.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      String(site.port).includes(searchQuery.value)

    const matchesRuntime = selectedRuntimeFilter.value === 'all' ||
      site.runtime === selectedRuntimeFilter.value

    return matchesSearch && matchesRuntime
  })
})

onMounted(() => {
  fetchWebsites()
  fetchDatabases()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm select-none">
    <!-- Top Action Bar -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-800 pb-3">
      <div>
        <div class="flex items-center gap-2">
          <Globe class="w-5 h-5 text-sky-400" />
          <h3 class="font-bold text-slate-100 text-base">Website Manager Suite</h3>
          <span class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-sky-500/20 text-sky-300 font-bold">
            {{ websites.length }} Sites
          </span>
        </div>
        <p class="text-xs text-slate-400 mt-0.5">
          Platform hosting multi-website terisolasi berbasis Docker dioptimasi untuk STB ARM64 (Statis, Go, PHP, Node.js, Python, Flat CMS).
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <!-- Shortcut Buka Database Manager -->
        <button
          @click="openDatabaseManager(null)"
          class="px-3 py-1.5 rounded-xl bg-amber-500/10 hover:bg-amber-500/20 text-amber-300 border border-amber-500/30 font-medium flex items-center gap-1.5 text-xs transition active:scale-95"
          title="Buka Web Database GUI (Adminer / HeidiSQL style)"
        >
          <Database class="w-3.5 h-3.5" />
          <span>Kelola DB</span>
        </button>

        <!-- Tombol Wizard Tambah Website -->
        <button
          @click="openCreateDialog"
          class="px-3.5 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow-lg shadow-sky-500/20 active:scale-95"
        >
          <Plus class="w-4 h-4 stroke-[2.5]" />
          <span>Buat Website Baru</span>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-2.5 bg-slate-900/60 p-2.5 rounded-xl border border-slate-800/80">
      <div class="relative flex-1 min-w-[220px]">
        <Search class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama website, port, domain, container..."
          class="w-full pl-8 pr-3 py-1.5 bg-slate-950/80 border border-slate-800 rounded-lg text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500 transition"
        />
      </div>

      <!-- Filter Kategori Runtime STB -->
      <div class="flex items-center gap-1.5 overflow-x-auto text-[11px] pb-0.5 lg:pb-0">
        <button
          v-for="f in [
            { id: 'all', label: 'Semua' },
            { id: 'static', label: 'HTML Statis' },
            { id: 'golang', label: 'Go API' },
            { id: 'php', label: 'PHP Native' },
            { id: 'nodejs', label: 'Node.js' },
            { id: 'python', label: 'Python' },
            { id: 'flatcms', label: 'Flat CMS' },
            { id: 'wordpress', label: 'WordPress' },
            { id: 'dashboard', label: 'Dashboard' }
          ]"
          :key="f.id"
          @click="selectedRuntimeFilter = f.id"
          class="px-2.5 py-1 rounded-lg font-medium transition shrink-0 whitespace-nowrap"
          :class="selectedRuntimeFilter === f.id
            ? 'bg-sky-500/20 text-sky-300 border border-sky-500/40 font-bold shadow-sm'
            : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'"
        >
          {{ f.label }}
        </button>
      </div>
    </div>

    <!-- Website Cards Grid: Auto-fill layout yang lega dan responsif -->
    <div
      v-if="filteredWebsites.length > 0"
      class="flex-1 overflow-y-auto grid gap-3.5 pr-0.5"
      style="grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));"
    >
      <div
        v-for="site in filteredWebsites"
        :key="site.id"
        class="raz-glass rounded-2xl border border-slate-800/90 p-4 sm:p-5 flex flex-col justify-between hover:border-sky-500/40 hover:shadow-2xl hover:shadow-sky-500/5 transition-all duration-200 group"
      >
        <!-- Card Top Section -->
        <div>
          <!-- Header Card: Icon, Name, Domain, Status Badge -->
          <div class="flex items-start gap-3.5 mb-3">
            <!-- App Icon (Proporsional & Clean) -->
            <div class="w-12 h-12 rounded-xl bg-slate-900/90 border border-slate-700/80 p-2.5 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform shadow">
              <img
                v-if="site.icon"
                :src="site.icon"
                :alt="site.name"
                class="w-full h-full object-contain"
                @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'"
              />
              <Globe v-else class="w-6 h-6 text-sky-400" />
            </div>

            <!-- Nama Website & Domain -->
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-2">
                <h4 class="font-bold text-slate-100 text-sm sm:text-base truncate group-hover:text-sky-300 transition-colors leading-tight" :title="site.name">
                  {{ site.name }}
                </h4>

                <!-- Status Indicator Badge -->
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider flex items-center gap-1.5 border shrink-0"
                  :class="site.status === 'running'
                    ? 'bg-emerald-500/15 text-emerald-300 border-emerald-500/30'
                    : 'bg-rose-500/15 text-rose-300 border-rose-500/30'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="site.status === 'running' ? 'bg-emerald-400 animate-pulse' : 'bg-rose-400'"></span>
                  {{ site.status }}
                </span>
              </div>

              <!-- Host Link & Runtime Tag -->
              <div class="flex items-center gap-2 mt-1">
                <a
                  :href="'http://' + (site.domain || 'localhost') + ':' + site.port"
                  target="_blank"
                  class="inline-flex items-center gap-1 font-mono text-[11px] text-sky-400 hover:text-sky-300 hover:underline truncate max-w-[210px]"
                  :title="'http://' + (site.domain || 'localhost') + ':' + site.port"
                >
                  <span class="truncate">{{ site.domain || 'localhost' }}</span>
                  <span class="text-slate-600">:</span>
                  <span class="text-amber-400 font-bold shrink-0">{{ site.port }}</span>
                </a>
                <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-slate-800 text-slate-300 uppercase shrink-0 font-semibold border border-slate-700/60">
                  {{ site.runtime }}
                </span>
              </div>
            </div>
          </div>

          <!-- Metadata Box (Lega, Rapi, & Tidak Terpotong) -->
          <div class="bg-slate-950/80 rounded-xl p-3 border border-slate-900/90 space-y-2 text-[11px]">
            <!-- Container Name -->
            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-500 flex items-center gap-1.5 shrink-0 whitespace-nowrap">
                <Server class="w-3.5 h-3.5 text-slate-400" />
                <span class="font-medium">Container:</span>
              </span>
              <span class="font-mono text-slate-200 font-semibold truncate text-right flex-1" :title="site.container_name">
                {{ site.container_name }}
              </span>
            </div>

            <!-- Root Folder -->
            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-500 flex items-center gap-1.5 shrink-0 whitespace-nowrap">
                <FolderOpen class="w-3.5 h-3.5 text-sky-400" />
                <span class="font-medium">Root Folder:</span>
              </span>
              <button
                @click="openWebsiteFolder(site)"
                class="font-mono text-slate-300 hover:text-sky-300 truncate text-right flex-1 hover:underline transition-colors text-[10.5px]"
                :title="site.root_folder"
              >
                {{ site.root_folder }}
              </button>
            </div>

            <!-- Resources -->
            <div class="flex items-center justify-between gap-2.5">
              <span class="text-slate-500 flex items-center gap-1.5 shrink-0 whitespace-nowrap">
                <Cpu class="w-3.5 h-3.5 text-emerald-400" />
                <span class="font-medium">Resources:</span>
              </span>
              <span class="font-mono text-slate-300 text-[10.5px] truncate text-right">
                RAM <strong class="text-emerald-400 font-semibold">{{ site.memory_limit || '128M' }}</strong> • CPU {{ site.cpu_shares || 512 }}
              </span>
            </div>

            <!-- Database Info (Jika Ada) -->
            <div v-if="site.db_name" class="flex items-center justify-between gap-2.5 pt-1 border-t border-slate-900">
              <span class="text-amber-400/90 flex items-center gap-1.5 shrink-0 whitespace-nowrap font-medium">
                <Database class="w-3.5 h-3.5 text-amber-400" />
                <span>Basis Data:</span>
              </span>
              <span class="font-mono text-amber-300 font-bold truncate text-right flex-1" :title="site.db_name">
                {{ site.db_name }}
              </span>
            </div>

            <!-- Git Info (Jika Ada) -->
            <div v-if="site.git_repo" class="flex items-center justify-between gap-2.5 pt-1 border-t border-slate-900">
              <span class="text-sky-400/90 flex items-center gap-1.5 shrink-0 whitespace-nowrap font-medium">
                <GitBranch class="w-3.5 h-3.5 text-sky-400" />
                <span>GitHub:</span>
              </span>
              <span class="font-mono text-slate-300 text-[10px] truncate text-right flex-1" :title="site.git_last_commit || site.git_repo">
                {{ site.git_branch || 'main' }} • {{ site.git_last_commit ? site.git_last_commit.substring(0, 7) : 'Terkoneksi' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Card Bottom Action Bar: Tertata Bersih & Proporsional -->
        <div class="flex items-center justify-between gap-2 pt-3 mt-3 border-t border-slate-800/80">
          <!-- Left Actions Toolbar: Start/Stop, Folder, Edit, Git, Database, Info -->
          <div class="flex items-center gap-1 bg-slate-950/60 p-1 rounded-xl border border-slate-900">
            <!-- Toggle Start / Stop -->
            <button
              v-if="site.status === 'running'"
              @click="toggleWebsiteStatus(site, 'stop')"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-amber-400 transition"
              title="Hentikan Container Website"
            >
              <Square class="w-3.5 h-3.5 fill-amber-400/20" />
            </button>
            <button
              v-else
              @click="toggleWebsiteStatus(site, 'start')"
              class="p-1.5 rounded-lg bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 transition"
              title="Jalankan Container Website"
            >
              <Play class="w-3.5 h-3.5 fill-emerald-400" />
            </button>

            <!-- Buka Root Folder di File Manager -->
            <button
              @click="openWebsiteFolder(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white transition"
              title="Buka Root Folder di File Manager RAZOS"
            >
              <FolderOpen class="w-3.5 h-3.5" />
            </button>

            <!-- Edit Konfigurasi Website -->
            <button
              @click="openEditDialog(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-amber-300 transition"
              title="Edit Konfigurasi Website"
            >
              <Pencil class="w-3.5 h-3.5" />
            </button>

            <!-- Quick Git Pull -->
            <button
              @click="executeGitPull(site)"
              :disabled="pullingSiteId === site.id"
              class="p-1.5 rounded-lg transition"
              :class="site.git_repo
                ? 'bg-sky-500/15 hover:bg-sky-500/25 text-sky-300'
                : 'bg-slate-900 hover:bg-slate-800 text-slate-400 hover:text-slate-200'"
              :title="site.git_repo ? 'Tarik Pembaruan GitHub Terbaru (Git Pull)' : 'Hubungkan Repository GitHub'"
            >
              <RotateCw v-if="pullingSiteId === site.id" class="w-3.5 h-3.5 animate-spin text-sky-400" />
              <GitPullRequest v-else class="w-3.5 h-3.5" />
            </button>

            <!-- Pengaturan Git & PAT Token -->
            <button
              @click="openGitModal(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-sky-300 transition"
              title="Pengaturan GitHub & PAT Token"
            >
              <GitBranch class="w-3.5 h-3.5" />
            </button>

            <!-- Buka Web Database GUI (Adminer) jika ada database -->
            <button
              v-if="site.db_name"
              @click="openDatabaseManager(site)"
              class="p-1.5 rounded-lg bg-amber-500/20 hover:bg-amber-500/30 text-amber-300 transition"
              title="Kelola Database di Web Adminer"
            >
              <Database class="w-3.5 h-3.5" />
            </button>

            <!-- Detail Konfigurasi -->
            <button
              @click="openDetailDialog(site)"
              class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 transition"
              title="Lihat Detail & Kredensial"
            >
              <Info class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Right Actions: Launch Web Button & Delete -->
          <div class="flex items-center gap-1.5">
            <a
              :href="'http://' + (site.domain || 'localhost') + ':' + site.port"
              target="_blank"
              class="px-3.5 py-1.5 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 font-bold flex items-center gap-1.5 text-xs transition shadow-md shadow-sky-500/20 active:scale-95"
            >
              <span>Buka</span>
              <ExternalLink class="w-3.5 h-3.5 stroke-[2.5]" />
            </a>

            <button
              @click="promptDeleteWebsite(site)"
              class="p-1.5 rounded-xl hover:bg-rose-500/20 text-rose-400 hover:text-rose-300 transition"
              title="Hapus Website"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else
      class="flex-1 rounded-2xl bg-slate-900/40 border border-slate-800/60 flex flex-col items-center justify-center text-slate-500 gap-3 p-8 text-center"
    >
      <div class="w-14 h-14 rounded-2xl bg-sky-500/10 border border-sky-500/20 flex items-center justify-center text-sky-400 shadow-xl shadow-sky-500/5">
        <Globe class="w-7 h-7" />
      </div>
      <div class="max-w-md">
        <h4 class="font-bold text-slate-200 text-sm sm:text-base">Belum Ada Website yang Dikonfigurasi</h4>
        <p class="text-xs text-slate-400 mt-1 leading-relaxed">
          Gunakan Website Manager untuk membuat container website terisolasi (HTML Statis, PHP dengan Database, Node.js, atau Next.js) dengan root folder dan alokasi port bebas.
        </p>
      </div>
      <button
        @click="openCreateDialog"
        class="mt-2 px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 transition shadow-lg shadow-sky-500/20 active:scale-95"
      >
        <Plus class="w-4 h-4 stroke-[2.5]" />
        <span>Buat Website Pertama Anda</span>
      </button>
    </div>

    <!-- CUSTOM MODAL: Wizard Pembuatan Website Terisolasi -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-3 sm:p-4 animate-in fade-in"
    >
      <div class="w-full max-w-2xl raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4 max-h-[92vh] overflow-hidden">
        <!-- Modal Header -->
        <div class="flex items-center justify-between border-b border-slate-800 pb-3 shrink-0">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <Layers class="w-4 h-4" />
            </div>
            <div>
              <h4 class="font-bold text-slate-100 text-sm sm:text-base">Wizard Pembuatan Website Terisolasi</h4>
              <p class="text-[11px] text-slate-400">Atur root folder, port, runtime stack, database link, dan batasan sumber daya.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Wizard Steps Indicator -->
        <div class="grid grid-cols-4 gap-1.5 bg-slate-950/80 p-1.5 rounded-xl border border-slate-800/80 text-center text-[11px] shrink-0">
          <button
            v-for="(st, idx) in [
              { num: 1, label: '1. Identitas & Port' },
              { num: 2, label: '2. Runtime & Folder' },
              { num: 3, label: '3. Koneksi Database' },
              { num: 4, label: '4. Resource Limit' }
            ]"
            :key="st.num"
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

        <!-- Wizard Body (Scrollable Form) -->
        <div class="flex-1 overflow-y-auto space-y-4 pr-1 text-xs">
          <!-- STEP 1: Identitas, Container Name, Domain & Port -->
          <div v-if="wizardStep === 1" class="space-y-3.5">
            <div>
              <label class="block font-medium text-slate-300 mb-1">Nama Website / Project <span class="text-rose-400">*</span></label>
              <input
                v-model="formSiteName"
                @input="onSiteNameChange"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400 transition"
                placeholder="Contoh: Toko Online UMKM"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-medium text-slate-300 mb-1">Nama Container (Terisolasi Anti-Tabrakan)</label>
                <input
                  v-model="formContainerName"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                  placeholder="Contoh: site-toko-umkm-web"
                />
                <span class="text-[10px] text-slate-500 mt-1 block">Nama unik Docker container agar tidak bentrok.</span>
              </div>

              <div>
                <label class="block font-medium text-slate-300 mb-1">Port Launching Host</label>
                <div class="flex items-center gap-2">
                  <input
                    v-model.number="formPort"
                    type="number"
                    min="80"
                    max="65535"
                    class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-amber-300 font-mono text-xs outline-none focus:border-sky-400 transition"
                  />
                  <button
                    type="button"
                    @click="suggestFreePort"
                    class="px-2.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono text-[10px] shrink-0"
                    title="Cari port bebas baru"
                  >
                    Auto Port
                  </button>
                </div>
                <span class="text-[10px] text-slate-500 mt-1 block">Akses langsung di peramban: :{{ formPort }}</span>
              </div>
            </div>

            <div>
              <label class="block font-medium text-slate-300 mb-1">Domain Hostname (Routing Proxy Internal)</label>
              <input
                v-model="formDomain"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                placeholder="Contoh: toko.razos.local"
              />
            </div>

            <!-- Icon Selector & Link Icon -->
            <div>
              <label class="block font-medium text-slate-300 mb-1">Link Icon (Untuk Desktop & App Drawer)</label>
              <div class="flex items-center gap-2">
                <input
                  v-model="formIcon"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400 transition"
                  placeholder="URL Ikon (misal: https://cdn.simpleicons.org/...)"
                />
                <div class="w-9 h-9 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center p-1.5 shrink-0">
                  <img :src="formIcon" class="w-full h-full object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
                </div>
              </div>
              <!-- Preset Icon Badges -->
              <div class="flex flex-wrap items-center gap-1.5 mt-2">
                <button
                  v-for="p in presetIcons"
                  :key="p.name"
                  type="button"
                  @click="formIcon = p.url"
                  class="px-2 py-0.5 rounded-lg bg-slate-900 border border-slate-800 hover:border-slate-700 text-[10px] text-slate-300 flex items-center gap-1 transition"
                >
                  <img :src="p.url" class="w-3 h-3 object-contain" />
                  <span>{{ p.name }}</span>
                </button>
              </div>
            </div>
          </div>

          <!-- STEP 2: Runtime Stack & Root Folder (STB Optimized) -->
          <div v-else-if="wizardStep === 2" class="space-y-3.5">
            <!-- Banner Edukasi STB S905X ARM64 -->
            <div class="p-2.5 rounded-xl bg-gradient-to-r from-sky-500/15 via-indigo-500/10 to-transparent border border-sky-500/25 flex items-start gap-2.5">
              <Sparkles class="w-4 h-4 text-sky-400 shrink-0 mt-0.5" />
              <div class="text-[11px] leading-relaxed">
                <span class="font-bold text-sky-300">Preset Khusus STB ARM64 (RAM 1.8 GiB):</span>
                <p class="text-slate-300 text-[10px] mt-0.5">
                  Telah difilter khusus agar sistem STB tidak lag atau kehabisan RAM. Runtime berat (Next.js SSR monolitik & MySQL 8 standar) telah ditiadakan demi stabilitas server.
                </p>
              </div>
            </div>

            <!-- Loop 4 Kategori Hosting STB -->
            <div class="space-y-3 max-h-[46vh] overflow-y-auto pr-1">
              <div
                v-for="cat in stbRuntimeCategories"
                :key="cat.title"
                class="rounded-xl border border-slate-800/80 bg-slate-950/60 p-3 space-y-2"
              >
                <!-- Category Header -->
                <div class="flex items-center justify-between gap-2 border-b border-slate-800/60 pb-1.5">
                  <div>
                    <h5 class="font-bold text-slate-200 text-xs flex items-center gap-1.5">
                      {{ cat.title }}
                    </h5>
                    <p class="text-[10px] text-slate-400">{{ cat.desc }}</p>
                  </div>
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold border shrink-0 whitespace-nowrap"
                    :class="cat.badgeColor"
                  >
                    {{ cat.badge }}
                  </span>
                </div>

                <!-- Runtimes Grid -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-0.5">
                  <div
                    v-for="rt in cat.runtimes"
                    :key="rt.id"
                    @click="onRuntimeChange(rt.id)"
                    class="p-2.5 rounded-xl border cursor-pointer transition flex items-start gap-2.5"
                    :class="formRuntime === rt.id
                      ? 'bg-sky-500/20 border-sky-400 text-sky-200 shadow-sm'
                      : 'bg-slate-900/60 border-slate-800 text-slate-400 hover:border-slate-700 hover:bg-slate-900'"
                  >
                    <div class="w-7 h-7 rounded-lg bg-slate-800 border border-slate-700/80 flex items-center justify-center shrink-0 p-1 mt-0.5">
                      <img :src="rt.icon" class="w-full h-full object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
                    </div>
                    <div class="min-w-0 flex-1">
                      <div class="flex items-center justify-between gap-1">
                        <h6 class="font-bold text-slate-200 text-[11px] truncate">{{ rt.name }}</h6>
                        <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-slate-800 text-slate-300 font-bold shrink-0">
                          {{ rt.ram }}
                        </span>
                      </div>
                      <p class="text-[10px] text-slate-400 mt-0.5 leading-snug">{{ rt.desc }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Root Folder (Document Root) -->
            <div>
              <label class="block font-medium text-slate-300 mb-1">Root Folder Penyimpanan Berkas Website</label>
              <div class="flex items-center gap-2">
                <input
                  v-model="formRootFolder"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                  placeholder="/DATA/apps/websites/nama-site"
                />
              </div>
              <p class="text-[10px] text-slate-500 mt-1">
                Folder ini akan otomatis dibuat di penyimpanan host dan dapat diedit langsung melalui File Manager RAZOS.
              </p>
            </div>

            <!-- Integrasi GitHub & PAT (Opsional) -->
            <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800/90 space-y-2.5">
              <div class="flex items-center gap-2">
                <GitBranch class="w-4 h-4 text-sky-400" />
                <div>
                  <h5 class="font-bold text-slate-200 text-xs">Integrasi GitHub & Auto-Pull (Opsional)</h5>
                  <p class="text-[10px] text-slate-400">Clone atau pull kode terbaru dari repository GitHub langsung ke root folder.</p>
                </div>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 pt-1">
                <div class="sm:col-span-2">
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">GitHub Repository URL</label>
                  <input
                    v-model="formGitRepo"
                    type="text"
                    class="w-full px-3 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="https://github.com/username/repo-web"
                  />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">Branch</label>
                  <input
                    v-model="formGitBranch"
                    type="text"
                    class="w-full px-3 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="main"
                  />
                </div>
              </div>

              <div>
                <label class="block text-[11px] font-medium text-slate-300 mb-1">GitHub Personal Access Token (PAT)</label>
                <div class="relative">
                  <input
                    :type="formShowGitToken ? 'text' : 'password'"
                    v-model="formGitToken"
                    class="w-full pl-8 pr-9 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="ghp_xxxxxxxxxxxxxxxxxxxx (untuk repo private/autentikasi)"
                  />
                  <Key class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-1/2 -translate-y-1/2" />
                  <button
                    type="button"
                    @click="formShowGitToken = !formShowGitToken"
                    class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200 transition"
                  >
                    <EyeOff v-if="formShowGitToken" class="w-3.5 h-3.5" />
                    <Eye v-else class="w-3.5 h-3.5" />
                  </button>
                </div>
                <p class="text-[10px] text-slate-500 mt-1">
                  PAT disimpan aman di lokal STB untuk auto-pull repositori private tanpa perlu memasukkan password berulang kali.
                </p>
              </div>
            </div>
          </div>

          <!-- STEP 3: Integrasi Database (HeidiSQL / Adminer Connected) -->
          <div v-if="wizardStep === 3" class="space-y-4">
            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-950 border border-slate-800">
              <div class="flex items-center gap-2.5">
                <Database class="w-5 h-5 text-amber-400" />
                <div>
                  <h5 class="font-bold text-slate-200">Koneksikan Basis Data untuk Website</h5>
                  <p class="text-[11px] text-slate-400">Hubungkan website dengan MariaDB atau PostgreSQL server.</p>
                </div>
              </div>
              <input
                type="checkbox"
                v-model="formEnableDb"
                class="w-5 h-5 accent-amber-400 rounded cursor-pointer"
              />
            </div>

            <div v-if="formEnableDb" class="space-y-3.5 animate-in fade-in">
              <!-- Pilih Database yang Sudah Ada di RAZOS -->
              <div v-if="databases.length > 0">
                <label class="block font-medium text-slate-300 mb-1">Gunakan Database Terdaftar di RAZOS (Pintasan)</label>
                <select
                  v-model="formSelectedDbId"
                  @change="onRegisteredDbSelect"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-amber-400"
                >
                  <option value="">-- Input Konfigurasi Kustom / Database Baru --</option>
                  <option v-for="db in databases" :key="db.id" :value="db.id">
                    {{ db.name }} ({{ db.engine.toUpperCase() }} - Port {{ db.port }}) - User: {{ db.user }}
                  </option>
                </select>
              </div>

              <!-- Detail Kredensial Database -->
              <!-- Pilihan Database Engine -->
              <div>
                <label class="block font-medium text-slate-300 mb-1">Database Engine</label>
                <select
                  v-model="formDbEngine"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-amber-400"
                >
                  <option value="sqlite">⭐ SQLite (Paling Ideal STB - Berbasis File - 0 MB Service RAM)</option>
                  <option value="mariadb">MariaDB / MySQL (Varian Alpine ~120-180 MB RAM)</option>
                  <option value="postgres">PostgreSQL Server</option>
                </select>
              </div>

              <!-- JIKA SQLITE DIPILIH (Paling Hemat RAM & Direkomendasikan) -->
              <div v-if="formDbEngine === 'sqlite'" class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 space-y-2.5">
                <div class="flex items-center gap-2 text-emerald-300 text-xs font-bold">
                  <Database class="w-4 h-4 text-emerald-400" />
                  <span>SQLite Embedded Database (Tanpa Daemon Service)</span>
                </div>
                <p class="text-[11px] text-slate-300 leading-relaxed">
                  Database berbasis berkas langsung di dalam root folder website. <strong>0 MB konsumsi RAM service terpisah</strong>, sangat cepat dan aman dari memory leak STB.
                </p>
                <div>
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">Nama File Database SQLite</label>
                  <input
                    v-model="formDbName"
                    type="text"
                    placeholder="database.sqlite"
                    class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-emerald-400"
                  />
                  <p class="text-[10px] text-slate-500 mt-1">
                    Path berkas otomatis: <span class="font-mono text-slate-400">{{ (formRootFolder || '/DATA/apps/websites/site') + '/' + (formDbName || 'database.sqlite') }}</span>
                  </p>
                </div>
              </div>

              <!-- JIKA MARIADB ATAU POSTGRES DIPILIH -->
              <div v-else class="space-y-3">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Host Database</label>
                    <input
                      v-model="formDbHost"
                      type="text"
                      placeholder="127.0.0.1"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Port Database</label>
                    <input
                      v-model.number="formDbPort"
                      type="number"
                      placeholder="3306"
                      class="w-full px-2 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400 text-center"
                    />
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Nama Database</label>
                    <input
                      v-model="formDbName"
                      type="text"
                      placeholder="db_tokoonline"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>

                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Username DB</label>
                    <input
                      v-model="formDbUser"
                      type="text"
                      placeholder="raz_user"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>

                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Password DB</label>
                    <input
                      v-model="formDbPass"
                      type="password"
                      placeholder="Kata sandi..."
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                </div>

                <!-- Catatan Terintegrasi Adminer -->
                <div class="p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-300 text-[11px] flex items-center gap-2">
                  <Database class="w-4 h-4 shrink-0" />
                  <span>
                    Variabel koneksi ini otomatis di-inject ke environment container (`DB_HOST`, `DB_NAME`, dll) dan dapat dikelola langsung via <strong>Web Adminer (HeidiSQL-style)</strong>.
                  </span>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-6 text-slate-500">
              <p>Website ini akan berjalan secara terisolasi tanpa koneksi database.</p>
            </div>
          </div>

          <!-- STEP 4: Resource Limit (RAM/CPU) & Advanced -->
          <div v-if="wizardStep === 4" class="space-y-4">
            <div>
              <label class="block font-medium text-slate-300 mb-2">Batas Penggunaan Memori RAM (Optimal untuk STB 1.8GB)</label>
              <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
                <button
                  v-for="mem in memoryPresets"
                  :key="mem.val"
                  type="button"
                  @click="formMemLimit = mem.val"
                  class="p-2.5 rounded-xl border text-left text-xs font-medium transition"
                  :class="formMemLimit === mem.val
                    ? 'bg-sky-500/20 border-sky-400 text-sky-300 font-bold'
                    : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
                >
                  {{ mem.label }}
                </button>
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-medium text-slate-300 mb-1">CPU Shares (Prioritas Penjadwalan CPU)</label>
                <select
                  v-model.number="formCpuShares"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-sky-400"
                >
                  <option :value="256">256 (Prioritas Ringan / Hemat)</option>
                  <option :value="512">512 (Standar Normal)</option>
                  <option :value="1024">1024 (Prioritas Maksimal)</option>
                </select>
              </div>

              <div>
                <label class="block font-medium text-slate-300 mb-1">Perintah Eksekusi Container (Cmd Opsional)</label>
                <input
                  v-model="formCommand"
                  type="text"
                  placeholder="Contoh: npm run start"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-sky-400"
                />
              </div>
            </div>

            <div>
              <label class="block font-medium text-slate-300 mb-1">Variabel Lingkungan Tambahan (ENV Vars)</label>
              <textarea
                v-model="formEnvVars"
                rows="2"
                placeholder="NODE_ENV=production&#10;CUSTOM_KEY=VALUE"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-emerald-400 font-mono text-xs outline-none focus:border-sky-400 resize-none leading-relaxed"
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
      </div>
    </div>

    <!-- CUSTOM MODAL: Konfirmasi Hapus Website -->
    <div
      v-if="showDeleteModal && siteToDelete"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-rose-500/30 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-rose-500/20 text-rose-400 flex items-center justify-center shrink-0 border border-rose-500/30">
            <AlertTriangle class="w-5 h-5" />
          </div>
          <div>
            <h4 class="font-bold text-slate-100 text-sm sm:text-base">Hapus Website Terisolasi?</h4>
            <p class="text-xs text-slate-400">Tindakan ini akan menghentikan dan mencabut container dari sistem.</p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 space-y-1 text-xs">
          <div class="text-slate-300 font-semibold">{{ siteToDelete.name }}</div>
          <div class="font-mono text-slate-400 text-[11px]">Container: {{ siteToDelete.container_name }}</div>
          <div class="font-mono text-amber-400 text-[11px]">Port: :{{ siteToDelete.port }}</div>
        </div>

        <label class="flex items-center gap-2.5 p-2 rounded-lg bg-rose-500/10 border border-rose-500/20 text-xs text-rose-300 cursor-pointer">
          <input type="checkbox" v-model="purgeRootFolder" class="w-4 h-4 accent-rose-500 rounded" />
          <span>Hapus juga seluruh berkas di root folder (`{{ siteToDelete.root_folder }}`)</span>
        </label>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showDeleteModal = false"
            class="px-3 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="executeDeleteWebsite"
            class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white font-bold text-xs transition shadow-lg shadow-rose-500/20 active:scale-95"
          >
            Hapus Website Sekarang
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Rincian & Kredensial Website -->
    <div
      v-if="showDetailModal && selectedSiteDetail"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-lg raz-glass rounded-2xl border border-slate-700/80 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div class="flex items-center gap-2.5">
            <img :src="selectedSiteDetail.icon" class="w-7 h-7 object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
            <h4 class="font-bold text-slate-100 text-sm sm:text-base">{{ selectedSiteDetail.name }}</h4>
          </div>
          <button @click="showDetailModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-2.5 text-xs text-slate-300">
          <div class="grid grid-cols-2 gap-2 bg-slate-950 p-3 rounded-xl border border-slate-800">
            <div>
              <span class="text-slate-500 block text-[10px]">CONTAINER NAME</span>
              <span class="font-mono text-slate-200">{{ selectedSiteDetail.container_name }}</span>
            </div>
            <div>
              <span class="text-slate-500 block text-[10px]">HOST PORT</span>
              <span class="font-mono text-amber-400 font-bold">:{{ selectedSiteDetail.port }}</span>
            </div>
            <div>
              <span class="text-slate-500 block text-[10px]">RUNTIME STACK</span>
              <span class="uppercase font-semibold text-sky-400">{{ selectedSiteDetail.runtime }}</span>
            </div>
            <div>
              <span class="text-slate-500 block text-[10px]">STATUS</span>
              <span class="uppercase font-bold text-emerald-400">{{ selectedSiteDetail.status }}</span>
            </div>
          </div>

          <!-- Root Folder Copy -->
          <div class="bg-slate-950 p-2.5 rounded-xl border border-slate-800 flex items-center justify-between">
            <div class="min-w-0 pr-2">
              <span class="text-slate-500 block text-[10px]">DOCUMENT ROOT FOLDER</span>
              <span class="font-mono text-[11px] text-slate-200 truncate block">{{ selectedSiteDetail.root_folder }}</span>
            </div>
            <button
              @click="copyToClipboard(selectedSiteDetail.root_folder, 'folder')"
              class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
              title="Salin Path Folder"
            >
              <Check v-if="copiedField === 'folder'" class="w-3.5 h-3.5 text-emerald-400" />
              <Copy v-else class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Database Details jika ada -->
          <div v-if="selectedSiteDetail.db_name" class="bg-amber-500/10 border border-amber-500/20 p-3 rounded-xl space-y-1.5 text-[11px]">
            <div class="font-bold text-amber-400 flex items-center gap-1.5">
              <Database class="w-3.5 h-3.5" />
              <span>Kredensial Database Terhubung</span>
            </div>
            <div class="grid grid-cols-2 gap-2 text-slate-300 font-mono">
              <div>Host: {{ selectedSiteDetail.db_host }}:{{ selectedSiteDetail.db_port }}</div>
              <div>Database: {{ selectedSiteDetail.db_name }}</div>
              <div>User: {{ selectedSiteDetail.db_user }}</div>
              <div>Engine: {{ selectedSiteDetail.db_engine || 'mariadb' }}</div>
            </div>
          </div>

          <!-- Git Sync Details jika ada -->
          <div v-if="selectedSiteDetail.git_repo" class="bg-sky-500/10 border border-sky-500/20 p-3 rounded-xl space-y-1.5 text-[11px]">
            <div class="font-bold text-sky-400 flex items-center justify-between">
              <span class="flex items-center gap-1.5">
                <GitBranch class="w-3.5 h-3.5" />
                <span>GitHub Repository Terhubung</span>
              </span>
              <button
                @click="showDetailModal = false; openGitModal(selectedSiteDetail)"
                class="text-[10px] text-sky-300 hover:underline flex items-center gap-1"
              >
                <span>Kelola Git</span>
                <ArrowRight class="w-3 h-3" />
              </button>
            </div>
            <div class="space-y-1 text-slate-300 font-mono text-[10px]">
              <div class="truncate">Repo: {{ selectedSiteDetail.git_repo }}</div>
              <div>Branch: {{ selectedSiteDetail.git_branch || 'main' }}</div>
              <div v-if="selectedSiteDetail.git_last_commit">Commit: {{ selectedSiteDetail.git_last_commit }}</div>
              <div v-if="selectedSiteDetail.git_last_pull">Last Pull: {{ selectedSiteDetail.git_last_pull }}</div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showDetailModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Tutup
          </button>
          <button
            @click="showDetailModal = false; openEditDialog(selectedSiteDetail)"
            class="px-3.5 py-2 rounded-xl bg-amber-500/15 hover:bg-amber-500/25 border border-amber-500/30 text-amber-300 text-xs font-semibold flex items-center gap-1.5 transition active:scale-95"
            title="Edit Konfigurasi Website"
          >
            <Pencil class="w-3.5 h-3.5" />
            <span>Edit Konfigurasi</span>
          </button>
          <a
            :href="'http://' + (selectedSiteDetail.domain || 'localhost') + ':' + selectedSiteDetail.port"
            target="_blank"
            class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow active:scale-95"
          >
            <span>Buka Website</span>
            <ExternalLink class="w-3.5 h-3.5" />
          </a>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Pengaturan GitHub & Git Sync -->
    <div
      v-if="showGitModal && gitModalSite"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-3 sm:p-4 animate-in fade-in"
    >
      <div class="w-full max-w-xl raz-glass rounded-2xl border border-sky-500/40 p-5 shadow-2xl flex flex-col gap-4 max-h-[90vh] overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-800 pb-3 shrink-0">
          <div class="flex items-center gap-2.5">
            <div class="w-9 h-9 rounded-xl bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <GitBranch class="w-5 h-5" />
            </div>
            <div>
              <h4 class="font-bold text-slate-100 text-sm sm:text-base flex items-center gap-2">
                <span>Pengaturan GitHub & Git Pull</span>
                <span class="text-xs px-2 py-0.5 rounded-full bg-slate-800 text-sky-300 font-normal">
                  {{ gitModalSite.name }}
                </span>
              </h4>
              <p class="text-[11px] text-slate-400">Sinkronisasi repositori kode langsung ke root folder website.</p>
            </div>
          </div>
          <button @click="showGitModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Body Form & Log -->
        <div class="flex-1 overflow-y-auto space-y-3.5 pr-1 text-xs">
          <!-- Folder Sasaran -->
          <div class="p-2.5 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between">
            <div class="min-w-0 pr-2">
              <span class="text-[10px] text-slate-500 block uppercase font-semibold">Root Folder Sasaran</span>
              <span class="font-mono text-slate-300 truncate block">{{ gitModalSite.root_folder }}</span>
            </div>
            <button
              @click="openWebsiteFolder(gitModalSite)"
              class="px-2 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-medium flex items-center gap-1 transition"
            >
              <FolderOpen class="w-3 h-3" />
              <span>Buka</span>
            </button>
          </div>

          <!-- Input Repository & Branch -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="sm:col-span-2">
              <label class="block font-medium text-slate-300 mb-1">GitHub Repository URL <span class="text-rose-400">*</span></label>
              <input
                v-model="gitModalRepo"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                placeholder="https://github.com/username/repository.git"
              />
            </div>
            <div>
              <label class="block font-medium text-slate-300 mb-1">Branch</label>
              <input
                v-model="gitModalBranch"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                placeholder="main"
              />
            </div>
          </div>

          <!-- Input Personal Access Token (PAT) -->
          <div>
            <label class="block font-medium text-slate-300 mb-1">GitHub Personal Access Token (PAT)</label>
            <div class="relative">
              <input
                :type="gitModalShowToken ? 'text' : 'password'"
                v-model="gitModalToken"
                class="w-full pl-8 pr-9 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                placeholder="ghp_xxxxxxxxxxxxxxxxxxxx (opsional untuk public, wajib untuk private)"
              />
              <Key class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-1/2 -translate-y-1/2" />
              <button
                type="button"
                @click="gitModalShowToken = !gitModalShowToken"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200 transition"
              >
                <EyeOff v-if="gitModalShowToken" class="w-3.5 h-3.5" />
                <Eye v-else class="w-3.5 h-3.5" />
              </button>
            </div>
            <p class="text-[10px] text-slate-500 mt-1">
              Token PAT disimpan aman di konfigurasi lokal RAZOS untuk otentikasi Git.
            </p>
          </div>

          <!-- Info Status Git -->
          <div v-if="gitModalSite.git_last_pull || gitModalSite.git_last_commit" class="p-2.5 rounded-xl bg-slate-950 border border-slate-800 space-y-1 text-[11px]">
            <div v-if="gitModalSite.git_last_commit" class="flex items-center justify-between">
              <span class="text-slate-500">Commit Terakhir:</span>
              <span class="font-mono text-sky-400 font-semibold">{{ gitModalSite.git_last_commit }}</span>
            </div>
            <div v-if="gitModalSite.git_last_pull" class="flex items-center justify-between">
              <span class="text-slate-500">Terakhir Ditarik:</span>
              <span class="font-mono text-slate-300">{{ gitModalSite.git_last_pull }}</span>
            </div>
          </div>

          <!-- Terminal Output Log -->
          <div v-if="gitModalLogs" class="space-y-1">
            <div class="flex items-center justify-between text-[11px] text-slate-400">
              <span class="flex items-center gap-1 font-mono text-sky-400">
                <Terminal class="w-3 h-3" />
                <span>Console Git Log:</span>
              </span>
              <button @click="gitModalLogs = ''" class="hover:text-slate-200 text-[10px]">Bersihkan</button>
            </div>
            <pre class="bg-black/90 p-3 rounded-xl border border-slate-800 text-[10px] font-mono text-emerald-400 whitespace-pre-wrap max-h-40 overflow-y-auto leading-relaxed">{{ gitModalLogs }}</pre>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
          <button
            @click="showGitModal = false"
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
              @click="executeGitPull(gitModalSite)"
              :disabled="gitModalIsPulling || !gitModalRepo"
              class="px-4 py-2 rounded-xl bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow-lg shadow-sky-500/20 active:scale-95 disabled:opacity-50"
            >
              <RotateCw v-if="gitModalIsPulling" class="w-3.5 h-3.5 animate-spin" />
              <GitPullRequest v-else class="w-3.5 h-3.5 stroke-[2.5]" />
              <span>{{ gitModalIsPulling ? 'Menarik Kode...' : 'Tarik Kode Sekarang (Pull)' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Edit Konfigurasi Website -->
    <div
      v-if="showEditModal && editSiteId"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-3 sm:p-4 animate-in fade-in"
    >
      <div class="w-full max-w-2xl raz-glass rounded-2xl border border-amber-500/40 p-5 shadow-2xl flex flex-col gap-4 max-h-[92vh] overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-800 pb-3 shrink-0">
          <div class="flex items-center gap-2.5">
            <div class="w-9 h-9 rounded-xl bg-amber-500/20 text-amber-400 flex items-center justify-center border border-amber-500/30">
              <Pencil class="w-4 h-4" />
            </div>
            <div>
              <h4 class="font-bold text-slate-100 text-sm sm:text-base flex items-center gap-2">
                <span>Edit Konfigurasi Website</span>
                <span class="text-xs px-2 py-0.5 rounded-full bg-slate-800 text-amber-300 font-mono">
                  #{{ editSiteId }}
                </span>
              </h4>
              <p class="text-[11px] text-slate-400">Perbarui identitas, domain, port, runtime, root folder, database, atau Git.</p>
            </div>
          </div>
          <button @click="showEditModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

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
            @click="editActiveTab = t.id"
            class="px-3 py-1.5 rounded-lg font-medium transition shrink-0"
            :class="editActiveTab === t.id
              ? 'bg-amber-500 text-slate-950 font-bold shadow'
              : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'"
          >
            {{ t.label }}
          </button>
        </div>

        <!-- Body Form (Scrollable) -->
        <div class="flex-1 overflow-y-auto space-y-3.5 pr-1 text-xs">
          <!-- TAB 1: General (Identitas, Domain, Port, Icon) -->
          <div v-if="editActiveTab === 'general'" class="space-y-3.5 animate-in fade-in">
            <div>
              <label class="block font-medium text-slate-300 mb-1">Nama Website / Project <span class="text-rose-400">*</span></label>
              <input
                v-model="editSiteName"
                type="text"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-amber-400 transition"
                placeholder="Contoh: Toko Online UMKM"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-medium text-slate-300 mb-1">Domain / Hostname Lokal</label>
                <input
                  v-model="editDomain"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 transition"
                  placeholder="toko.razos.local"
                />
              </div>

              <div>
                <label class="block font-medium text-slate-300 mb-1">Port Launching Host</label>
                <input
                  v-model.number="editPort"
                  type="number"
                  min="80"
                  max="65535"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-amber-300 font-mono text-xs outline-none focus:border-amber-400 transition"
                />
              </div>
            </div>

            <!-- Icon URL & Presets -->
            <div>
              <label class="block font-medium text-slate-300 mb-1">Link Ikon Website</label>
              <div class="flex items-center gap-2">
                <input
                  v-model="editIcon"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-amber-400 transition"
                  placeholder="URL Ikon (misal: https://cdn.simpleicons.org/...)"
                />
                <div class="w-9 h-9 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center p-1.5 shrink-0">
                  <img :src="editIcon" class="w-full h-full object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
                </div>
              </div>
              <div class="flex flex-wrap items-center gap-1.5 mt-2">
                <button
                  v-for="p in presetIcons"
                  :key="p.name"
                  type="button"
                  @click="editIcon = p.url"
                  class="px-2 py-0.5 rounded-lg bg-slate-900 border border-slate-800 hover:border-slate-700 text-[10px] text-slate-300 flex items-center gap-1 transition"
                >
                  <img :src="p.url" class="w-3 h-3 object-contain" />
                  <span>{{ p.name }}</span>
                </button>
              </div>
            </div>
          </div>

          <!-- TAB 2: Runtime & Root Folder (STB Optimized) -->
          <div v-else-if="editActiveTab === 'runtime'" class="space-y-3.5 animate-in fade-in">
            <!-- Banner Edukasi STB S905X ARM64 -->
            <div class="p-2.5 rounded-xl bg-gradient-to-r from-amber-500/15 via-sky-500/10 to-transparent border border-amber-500/25 flex items-start gap-2.5">
              <Sparkles class="w-4 h-4 text-amber-400 shrink-0 mt-0.5" />
              <div class="text-[11px] leading-relaxed">
                <span class="font-bold text-amber-300">Pilihan Runtime Ramah Hardware STB (RAM 1.8 GiB):</span>
                <p class="text-slate-300 text-[10px] mt-0.5">
                  Ubah arsitektur runtime sesuai kebutuhan. Runtime Next.js SSR telah digantikan dengan opsi micro-service cepat dan hemat RAM.
                </p>
              </div>
            </div>

            <!-- Loop 4 Kategori STB di Edit Modal -->
            <div class="space-y-3 max-h-[46vh] overflow-y-auto pr-1">
              <div
                v-for="cat in stbRuntimeCategories"
                :key="cat.title"
                class="rounded-xl border border-slate-800/80 bg-slate-950/60 p-3 space-y-2"
              >
                <!-- Category Header -->
                <div class="flex items-center justify-between gap-2 border-b border-slate-800/60 pb-1.5">
                  <div>
                    <h5 class="font-bold text-slate-200 text-xs flex items-center gap-1.5">
                      {{ cat.title }}
                    </h5>
                    <p class="text-[10px] text-slate-400">{{ cat.desc }}</p>
                  </div>
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold border shrink-0 whitespace-nowrap"
                    :class="cat.badgeColor"
                  >
                    {{ cat.badge }}
                  </span>
                </div>

                <!-- Runtimes Grid -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-0.5">
                  <div
                    v-for="rt in cat.runtimes"
                    :key="rt.id"
                    @click="onEditRuntimeChange(rt.id)"
                    class="p-2.5 rounded-xl border cursor-pointer transition flex items-start gap-2.5"
                    :class="editRuntime === rt.id
                      ? 'bg-amber-500/20 border-amber-400 text-amber-200 shadow-sm'
                      : 'bg-slate-900/60 border-slate-800 text-slate-400 hover:border-slate-700 hover:bg-slate-900'"
                  >
                    <div class="w-7 h-7 rounded-lg bg-slate-800 border border-slate-700/80 flex items-center justify-center shrink-0 p-1 mt-0.5">
                      <img :src="rt.icon" class="w-full h-full object-contain" @error="$event.target.src = '/assets/icons/RAZIconWebsites.svg'" />
                    </div>
                    <div class="min-w-0 flex-1">
                      <div class="flex items-center justify-between gap-1">
                        <h6 class="font-bold text-slate-200 text-[11px] truncate">{{ rt.name }}</h6>
                        <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-slate-800 text-slate-300 font-bold shrink-0">
                          {{ rt.ram }}
                        </span>
                      </div>
                      <p class="text-[10px] text-slate-400 mt-0.5 leading-snug">{{ rt.desc }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <label class="block font-medium text-slate-300 mb-1">Root Folder Penyimpanan Berkas Website</label>
              <div class="flex items-center gap-2">
                <input
                  v-model="editRootFolder"
                  type="text"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 transition"
                />
                <button
                  type="button"
                  @click="windowStore.openApp('file-manager')"
                  class="px-2.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 flex items-center gap-1 shrink-0"
                  title="Lihat di File Manager"
                >
                  <FolderOpen class="w-3.5 h-3.5" />
                  <span>Buka</span>
                </button>
              </div>
              <p class="text-[10px] text-slate-500 mt-1">
                Lokasi folder fisik berkas website di hosting atau storage RAZOS.
              </p>
            </div>
          </div>

          <!-- TAB 3: Basis Data -->
          <div v-else-if="editActiveTab === 'database'" class="space-y-3.5 animate-in fade-in">
            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-950 border border-slate-800">
              <div class="flex items-center gap-2.5">
                <Database class="w-5 h-5 text-amber-400" />
                <div>
                  <h5 class="font-bold text-slate-200">Koneksikan Basis Data untuk Website</h5>
                  <p class="text-[11px] text-slate-400">Hubungkan website dengan MariaDB atau PostgreSQL server.</p>
                </div>
              </div>
              <input
                type="checkbox"
                v-model="editEnableDb"
                class="w-5 h-5 accent-amber-400 rounded cursor-pointer"
              />
            </div>

            <div v-if="editEnableDb" class="space-y-3.5">
              <div v-if="databases.length > 0">
                <label class="block font-medium text-slate-300 mb-1">Gunakan Database Terdaftar di RAZOS (Pintasan)</label>
                <select
                  v-model="editSelectedDbId"
                  @change="onEditRegisteredDbSelect"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-amber-400"
                >
                  <option value="">-- Input Konfigurasi Kustom / Database Baru --</option>
                  <option v-for="db in databases" :key="db.id" :value="String(db.id)">
                    {{ db.name }} ({{ db.engine.toUpperCase() }} - Port {{ db.port }}) - User: {{ db.user }}
                  </option>
                </select>
              </div>

              <!-- Pilihan Database Engine di Edit Modal -->
              <div>
                <label class="block font-medium text-slate-300 mb-1">Database Engine</label>
                <select
                  v-model="editDbEngine"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-amber-400"
                >
                  <option value="sqlite">⭐ SQLite (Paling Ideal STB - Berbasis File - 0 MB Service RAM)</option>
                  <option value="mariadb">MariaDB / MySQL (Varian Alpine ~120-180 MB RAM)</option>
                  <option value="postgres">PostgreSQL Server</option>
                </select>
              </div>

              <!-- JIKA SQLITE DIPILIH -->
              <div v-if="editDbEngine === 'sqlite'" class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 space-y-2.5">
                <div class="flex items-center gap-2 text-emerald-300 text-xs font-bold">
                  <Database class="w-4 h-4 text-emerald-400" />
                  <span>SQLite Embedded Database (Tanpa Daemon Service)</span>
                </div>
                <p class="text-[11px] text-slate-300 leading-relaxed">
                  Database berbasis berkas langsung di dalam root folder website. <strong>0 MB konsumsi RAM service terpisah</strong>, sangat cepat dan aman dari memory leak STB.
                </p>
                <div>
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">Nama File Database SQLite</label>
                  <input
                    v-model="editDbName"
                    type="text"
                    placeholder="database.sqlite"
                    class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-emerald-400"
                  />
                  <p class="text-[10px] text-slate-500 mt-1">
                    Path berkas otomatis: <span class="font-mono text-slate-400">{{ (editRootFolder || '/DATA/apps/websites/site') + '/' + (editDbName || 'database.sqlite') }}</span>
                  </p>
                </div>
              </div>

              <!-- JIKA MARIADB ATAU POSTGRES DIPILIH -->
              <div v-else class="space-y-3">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Host Database</label>
                    <input
                      v-model="editDbHost"
                      type="text"
                      placeholder="127.0.0.1"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Port Database</label>
                    <input
                      v-model.number="editDbPort"
                      type="number"
                      placeholder="3306"
                      class="w-full px-2 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400 text-center"
                    />
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Nama Database</label>
                    <input
                      v-model="editDbName"
                      type="text"
                      placeholder="db_nama_website"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Username</label>
                    <input
                      v-model="editDbUser"
                      type="text"
                      placeholder="raz_user"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                  <div>
                    <label class="block font-medium text-slate-300 mb-1">Password</label>
                    <input
                      v-model="editDbPass"
                      type="password"
                      placeholder="••••••••"
                      class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono text-xs outline-none focus:border-amber-400"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- TAB 4: GitHub & PAT -->
          <div v-else-if="editActiveTab === 'git'" class="space-y-3.5 animate-in fade-in">
            <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 space-y-2.5">
              <div class="flex items-center gap-2">
                <GitBranch class="w-4 h-4 text-sky-400" />
                <div>
                  <h5 class="font-bold text-slate-200 text-xs">Konfigurasi GitHub & Git Pull</h5>
                  <p class="text-[10px] text-slate-400">Atur repositori Git untuk auto-pull kode ke root folder.</p>
                </div>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 pt-1">
                <div class="sm:col-span-2">
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">GitHub Repository URL</label>
                  <input
                    v-model="editGitRepo"
                    type="text"
                    class="w-full px-3 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="https://github.com/username/repo-web"
                  />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-slate-300 mb-1">Branch</label>
                  <input
                    v-model="editGitBranch"
                    type="text"
                    class="w-full px-3 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="main"
                  />
                </div>
              </div>

              <div>
                <label class="block text-[11px] font-medium text-slate-300 mb-1">GitHub Personal Access Token (PAT)</label>
                <div class="relative">
                  <input
                    :type="editShowGitToken ? 'text' : 'password'"
                    v-model="editGitToken"
                    class="w-full pl-8 pr-9 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
                    placeholder="ghp_xxxxxxxxxxxxxxxxxxxx (opsional untuk public, wajib untuk private)"
                  />
                  <Key class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-1/2 -translate-y-1/2" />
                  <button
                    type="button"
                    @click="editShowGitToken = !editShowGitToken"
                    class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200 transition"
                  >
                    <EyeOff v-if="editShowGitToken" class="w-3.5 h-3.5" />
                    <Eye v-else class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- TAB 5: Resource Limit & ENV -->
          <div v-else-if="editActiveTab === 'advanced'" class="space-y-3.5 animate-in fade-in">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-medium text-slate-300 mb-1">Memory Limit (RAM)</label>
                <input
                  v-model="editMemLimit"
                  type="text"
                  placeholder="128M"
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
                <label class="block font-medium text-slate-300 mb-1">CPU Shares (Prioritas CPU)</label>
                <input
                  v-model.number="editCpuShares"
                  type="number"
                  placeholder="512"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
                />
              </div>
            </div>

            <div>
              <label class="block font-medium text-slate-300 mb-1">Perintah Eksekusi Container (Command)</label>
              <input
                v-model="editCommand"
                type="text"
                placeholder="Contoh: node server.js atau python app.py"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400"
              />
            </div>

            <div>
              <label class="block font-medium text-slate-300 mb-1">Variabel Lingkungan (ENV Vars)</label>
              <textarea
                v-model="editEnvVars"
                rows="3"
                placeholder="KEY=VALUE&#10;NODE_ENV=production"
                class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-amber-400 resize-none"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
          <button
            @click="showEditModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>

          <button
            @click="confirmUpdateWebsite"
            :disabled="isSavingEdit"
            class="px-5 py-2 rounded-xl bg-gradient-to-r from-amber-500 to-amber-400 hover:from-amber-400 hover:to-amber-300 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow-lg shadow-amber-500/20 active:scale-95 disabled:opacity-50"
          >
            <RotateCw v-if="isSavingEdit" class="w-3.5 h-3.5 animate-spin" />
            <Check v-else class="w-3.5 h-3.5 stroke-[2.5]" />
            <span>{{ isSavingEdit ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
