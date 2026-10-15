<script setup>
/**
 * RAZOS App Store Manual Install View (CasaOS Style - Standalone Desktop Window)
 * Updated: 2026-10-15 (v2.7.0 - Full Feature Parity & Import Docker Compose/CLI)
 * Deskripsi: Antarmuka instalasi aplikasi mandiri tingkat lanjut sesuai 5 spesifikasi referensi:
 *   1. Form Manual App Install lengkap (Image, Tag, Title, Icon, Web UI, Network, dynamic Ports,
 *      Volumes, Envs, Devices, Command, Privileges switch, Memory Limit slider 256-1800, CPU Shares,
 *      Restart Policy, Cap-Add, Container Name).
 *   2. Modal Import terintegrasi: Docker Compose (upload file YAML / drag-and-drop & paste)
 *      dan Docker CLI (tokenizer perintah 'docker run').
 *   3. Desain Glassmorphism Slate-900 / Slate-950 modern, responsif, dan zero browser alert.
 */
import { ref, reactive, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import YAML from 'yaml'
import {
  LogIn,
  X,
  Plus,
  Trash2,
  Upload,
  Info,
  Check,
  AlertCircle,
  HelpCircle,
  Box,
  Layers,
  Terminal,
  FileCode,
  Globe,
  HardDrive,
  Cpu,
  ChevronDown
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ id: 'appstore-manual' })
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

// State Formulir Instalasi Manual Aplikasi
const form = reactive({
  image: '',
  tag: 'latest',
  title: '',
  icon: '',
  webUiProtocol: 'http://',
  webUiHost: typeof window !== 'undefined' ? (window.location.hostname || 'server.razos.my.id') : 'server.razos.my.id',
  webUiPort: '',
  webUiPath: '/',
  network: 'bridge',
  ports: [],
  volumes: [],
  envs: [],
  devices: [],
  commands: [],
  privileged: false,
  memoryLimit: 512,
  cpuShares: 'Low',
  restartPolicy: 'unless-stopped',
  capAdd: '',
  containerName: ''
})

// Status Proses Deploy
const isDeploying = ref(false)

// State Modal Import (Docker Compose / Docker CLI)
const showImportModal = ref(false)
const importTab = ref('compose') // 'compose' | 'cli'
const composeInput = ref('')
const cliInput = ref('')
const isDragActive = ref(false)
const importError = ref('')
const fileInputRef = ref(null)

// Validasi Icon URL & Status Gambar Rusak
const iconLoadError = ref(false)
function onIconError() {
  iconLoadError.value = true
}
function onIconInput() {
  iconLoadError.value = false
}

// Auto-sync jika user mengetik nama image dengan format "repo/image:tag"
function onImageBlur() {
  if (form.image.includes(':')) {
    const parts = form.image.split(':')
    form.image = parts[0].trim()
    if (parts[1]) {
      form.tag = parts[1].trim()
    }
  }
  // Auto-generate Title dan Container Name jika masih kosong
  if (!form.title.trim() && form.image.trim()) {
    const rawName = form.image.split('/').pop().replace(/[^a-zA-Z0-9_-]/g, ' ')
    form.title = rawName.charAt(0).toUpperCase() + rawName.slice(1)
  }
  if (!form.containerName.trim() && form.title.trim()) {
    form.containerName = form.title.toLowerCase().replace(/[^a-z0-9_-]/g, '-')
  }
}

// Auto-sync Container Name saat Title berubah
function onTitleInput() {
  if (!form.containerName || form.containerName === form.title.toLowerCase().slice(0, -1).replace(/[^a-z0-9_-]/g, '-')) {
    form.containerName = form.title.toLowerCase().replace(/[^a-z0-9_-]/g, '-')
  }
}

// Sinkronisasi otomatis port Web UI dari daftar port pertama jika port Web UI masih kosong
function syncWebUiPort() {
  if (!form.webUiPort && form.ports.length > 0 && form.ports[0].host) {
    form.webUiPort = form.ports[0].host
  }
}

// ==========================================
// OPERASI DAFTAR DINAMIS (PORTS, VOLS, ENVS, DLL)
// ==========================================
function addPort() {
  form.ports.push({
    host: '',
    container: '',
    protocol: 'TCP'
  })
}
function removePort(index) {
  form.ports.splice(index, 1)
}

function addVolume() {
  form.volumes.push({
    host: `/DATA/apps/${form.containerName || 'custom-app'}/data`,
    container: '/data',
    mode: 'rw'
  })
}
function removeVolume(index) {
  form.volumes.splice(index, 1)
}

function addEnv() {
  form.envs.push({
    key: '',
    value: ''
  })
}
function removeEnv(index) {
  form.envs.splice(index, 1)
}

function addDevice() {
  form.devices.push({
    host: '/dev/dri',
    container: '/dev/dri'
  })
}
function removeDevice(index) {
  form.devices.splice(index, 1)
}

function addCommand() {
  form.commands.push({
    value: ''
  })
}
function removeCommand(index) {
  form.commands.splice(index, 1)
}

// ==========================================
// IMPORT ENGINE (DOCKER COMPOSE & DOCKER CLI)
// ==========================================
function openImportDialog() {
  importError.value = ''
  showImportModal.value = true
}

function closeImportDialog() {
  showImportModal.value = false
  importError.value = ''
}

// Buka dialog file explorer untuk upload file YAML
function triggerFileInput() {
  if (fileInputRef.value) {
    fileInputRef.value.click()
  }
}

// Handle pemilihan file YAML via browser file dialog
function handleFileSelect(e) {
  const file = e.target.files?.[0]
  if (!file) return
  readFileContent(file)
}

// Handle Drag and Drop file YAML
function handleDrop(e) {
  isDragActive.value = false
  const file = e.dataTransfer?.files?.[0]
  if (!file) return
  readFileContent(file)
}

// Baca konten berkas menggunakan FileReader
function readFileContent(file) {
  const reader = new FileReader()
  reader.onload = (evt) => {
    composeInput.value = evt.target?.result || ''
    importError.value = ''
  }
  reader.onerror = () => {
    importError.value = 'Gagal membaca berkas file'
  }
  reader.readAsText(file)
}

// Parsing Docker Compose YAML ke state form
function applyComposeYaml(yamlText) {
  if (!yamlText.trim()) {
    throw new Error('Skrip Docker Compose tidak boleh kosong!')
  }

  const doc = YAML.parse(yamlText)
  if (!doc) throw new Error('Format YAML tidak dapat di-parse.')

  const services = doc.services || doc
  const serviceKeys = Object.keys(services)
  if (serviceKeys.length === 0) {
    throw new Error('Tidak ditemukan blok "services" dalam skrip Compose.')
  }

  // Gunakan service pertama sesuai catatan di antarmuka CasaOS
  const serviceKey = serviceKeys[0]
  const svc = services[serviceKey]
  if (!svc || typeof svc !== 'object') {
    throw new Error(`Service "${serviceKey}" tidak memiliki konfigurasi valid.`)
  }

  // 1. Ekstrak Image & Tag
  let imgStr = svc.image || ''
  if (imgStr.includes(':')) {
    const idx = imgStr.lastIndexOf(':')
    form.image = imgStr.substring(0, idx)
    form.tag = imgStr.substring(idx + 1)
  } else {
    form.image = imgStr
    form.tag = 'latest'
  }

  // 2. Ekstrak Title & Container Name
  const rawTitle = svc.container_name || serviceKey || form.image.split('/').pop() || 'Custom App'
  form.title = rawTitle.charAt(0).toUpperCase() + rawTitle.slice(1)
  form.containerName = svc.container_name || rawTitle.toLowerCase().replace(/[^a-z0-9_-]/g, '-')

  // 3. Ekstrak Network
  if (svc.network_mode) {
    form.network = svc.network_mode
  }

  // 4. Ekstrak Ports
  form.ports = []
  if (Array.isArray(svc.ports)) {
    svc.ports.forEach(p => {
      if (typeof p === 'object' && p !== null) {
        form.ports.push({
          host: String(p.published || ''),
          container: String(p.target || ''),
          protocol: String(p.protocol || 'TCP').toUpperCase()
        })
      } else {
        const str = String(p).trim()
        let proto = 'TCP'
        let cleanStr = str
        if (str.toLowerCase().endsWith('/udp')) {
          proto = 'UDP'
          cleanStr = str.replace(/\/udp$/i, '')
        } else if (str.toLowerCase().endsWith('/tcp')) {
          cleanStr = str.replace(/\/tcp$/i, '')
        }
        const parts = cleanStr.split(':')
        if (parts.length >= 2) {
          form.ports.push({
            host: parts[0],
            container: parts[1],
            protocol: proto
          })
        } else if (parts.length === 1) {
          form.ports.push({
            host: parts[0],
            container: parts[0],
            protocol: proto
          })
        }
      }
    })
  }

  // Auto-set Web UI Port jika ada
  if (form.ports.length > 0 && form.ports[0].host) {
    form.webUiPort = form.ports[0].host
  }

  // 5. Ekstrak Volumes
  form.volumes = []
  if (Array.isArray(svc.volumes)) {
    svc.volumes.forEach(v => {
      if (typeof v === 'object' && v !== null) {
        form.volumes.push({
          host: String(v.source || ''),
          container: String(v.target || ''),
          mode: v.read_only ? 'ro' : 'rw'
        })
      } else {
        const parts = String(v).split(':')
        if (parts.length >= 2) {
          form.volumes.push({
            host: parts[0],
            container: parts[1],
            mode: parts[2] || 'rw'
          })
        }
      }
    })
  }

  // 6. Ekstrak Environment Variables
  form.envs = []
  if (Array.isArray(svc.environment)) {
    svc.environment.forEach(e => {
      const str = String(e)
      const eqIdx = str.indexOf('=')
      if (eqIdx !== -1) {
        form.envs.push({
          key: str.substring(0, eqIdx).trim(),
          value: str.substring(eqIdx + 1).trim()
        })
      } else {
        form.envs.push({ key: str.trim(), value: '' })
      }
    })
  } else if (svc.environment && typeof svc.environment === 'object') {
    Object.entries(svc.environment).forEach(([k, v]) => {
      form.envs.push({ key: String(k), value: String(v ?? '') })
    })
  }

  // 7. Ekstrak Devices
  form.devices = []
  if (Array.isArray(svc.devices)) {
    svc.devices.forEach(d => {
      const parts = String(d).split(':')
      form.devices.push({
        host: parts[0] || '',
        container: parts[1] || parts[0] || ''
      })
    })
  }

  // 8. Ekstrak Commands
  form.commands = []
  if (Array.isArray(svc.command)) {
    svc.command.forEach(c => form.commands.push({ value: String(c) }))
  } else if (typeof svc.command === 'string' && svc.command.trim()) {
    form.commands.push({ value: svc.command.trim() })
  }

  // 9. Privileged
  form.privileged = !!svc.privileged

  // 10. Restart Policy
  if (svc.restart) {
    form.restartPolicy = String(svc.restart).toLowerCase()
  }

  // 11. Cap-Add
  if (Array.isArray(svc.cap_add)) {
    form.capAdd = svc.cap_add.join(', ')
  } else if (typeof svc.cap_add === 'string') {
    form.capAdd = svc.cap_add
  }

  // 12. Memory Limit
  if (svc.deploy?.resources?.limits?.memory) {
    const memStr = String(svc.deploy.resources.limits.memory).toUpperCase()
    if (memStr.endsWith('M') || memStr.endsWith('MB')) {
      form.memoryLimit = Math.min(1800, Math.max(256, parseInt(memStr) || 512))
    } else if (memStr.endsWith('G') || memStr.endsWith('GB')) {
      form.memoryLimit = Math.min(1800, Math.max(256, Math.round(parseFloat(memStr) * 1024) || 1024))
    }
  } else if (svc.mem_limit) {
    const memStr = String(svc.mem_limit).toUpperCase()
    if (memStr.endsWith('M') || memStr.endsWith('MB')) {
      form.memoryLimit = Math.min(1800, Math.max(256, parseInt(memStr) || 512))
    } else if (memStr.endsWith('G') || memStr.endsWith('GB')) {
      form.memoryLimit = Math.min(1800, Math.max(256, Math.round(parseFloat(memStr) * 1024) || 1024))
    }
  }
}

// Tokenizer perintah Docker CLI ('docker run ...') ke state form
function applyDockerCLI(cliString) {
  let str = cliString.trim()
  if (!str) {
    throw new Error('Perintah Docker CLI tidak boleh kosong!')
  }

  if (str.startsWith('docker run')) {
    str = str.replace(/^docker\s+run\s+/, '')
  }

  // Tokenize regex dengan mempertimbangkan tanda petik kutip tunggal dan ganda
  const regex = /[^\s"']+|"([^"]*)"|'([^']*)'/g
  const tokens = []
  let match
  while ((match = regex.exec(str)) !== null) {
    tokens.push(match[1] !== undefined ? match[1] : (match[2] !== undefined ? match[2] : match[0]))
  }

  let foundImage = ''
  let foundTag = 'latest'
  let foundName = ''
  const parsedPorts = []
  const parsedVols = []
  const parsedEnvs = []
  const parsedDevices = []
  const parsedCommands = []
  let parsedPrivileged = false
  let parsedRestart = 'unless-stopped'
  let parsedNetwork = 'bridge'
  let parsedCap = []
  let parsedMem = 512

  let i = 0
  while (i < tokens.length) {
    const token = tokens[i]

    if (token === '-d' || token === '--detach' || token === '-it' || token === '-i' || token === '-t' || token === '--rm') {
      i++
      continue
    }

    if (token === '--privileged') {
      parsedPrivileged = true
      i++
      continue
    }

    if (token.startsWith('--name=')) {
      foundName = token.split('=')[1]
      i++
      continue
    } else if (token === '--name' && i + 1 < tokens.length) {
      foundName = tokens[i + 1]
      i += 2
      continue
    }

    if (token.startsWith('-p=') || token.startsWith('--publish=')) {
      parsePortToken(token.split('=')[1], parsedPorts)
      i++
      continue
    } else if ((token === '-p' || token === '--publish') && i + 1 < tokens.length) {
      parsePortToken(tokens[i + 1], parsedPorts)
      i += 2
      continue
    }

    if (token.startsWith('-v=') || token.startsWith('--volume=')) {
      parseVolumeToken(token.split('=')[1], parsedVols)
      i++
      continue
    } else if ((token === '-v' || token === '--volume') && i + 1 < tokens.length) {
      parseVolumeToken(tokens[i + 1], parsedVols)
      i += 2
      continue
    }

    if (token.startsWith('-e=') || token.startsWith('--env=')) {
      parseEnvToken(token.substring(token.indexOf('=') + 1), parsedEnvs)
      i++
      continue
    } else if ((token === '-e' || token === '--env') && i + 1 < tokens.length) {
      parseEnvToken(tokens[i + 1], parsedEnvs)
      i += 2
      continue
    }

    if (token.startsWith('--restart=')) {
      parsedRestart = token.split('=')[1]
      i++
      continue
    } else if (token === '--restart' && i + 1 < tokens.length) {
      parsedRestart = tokens[i + 1]
      i += 2
      continue
    }

    if (token.startsWith('--net=') || token.startsWith('--network=')) {
      parsedNetwork = token.split('=')[1]
      i++
      continue
    } else if ((token === '--net' || token === '--network') && i + 1 < tokens.length) {
      parsedNetwork = tokens[i + 1]
      i += 2
      continue
    }

    if (token.startsWith('--device=')) {
      parseDeviceToken(token.split('=')[1], parsedDevices)
      i++
      continue
    } else if (token === '--device' && i + 1 < tokens.length) {
      parseDeviceToken(tokens[i + 1], parsedDevices)
      i += 2
      continue
    }

    if (token.startsWith('--cap-add=')) {
      parsedCap.push(token.split('=')[1])
      i++
      continue
    } else if (token === '--cap-add' && i + 1 < tokens.length) {
      parsedCap.push(tokens[i + 1])
      i += 2
      continue
    }

    if (token.startsWith('-m=') || token.startsWith('--memory=')) {
      parsedMem = parseMemoryString(token.split('=')[1])
      i++
      continue
    } else if ((token === '-m' || token === '--memory') && i + 1 < tokens.length) {
      parsedMem = parseMemoryString(tokens[i + 1])
      i += 2
      continue
    }

    // Token non-opsi pertama dianggap sebagai Image
    if (!foundImage && !token.startsWith('-')) {
      if (token.includes(':')) {
        const idx = token.lastIndexOf(':')
        foundImage = token.substring(0, idx)
        foundTag = token.substring(idx + 1)
      } else {
        foundImage = token
        foundTag = 'latest'
      }
      i++
      continue
    }

    // Token setelah image dianggap argumen container command
    if (foundImage) {
      parsedCommands.push({ value: token })
      i++
      continue
    }

    i++
  }

  if (!foundImage) {
    throw new Error('Tidak dapat menemukan Docker Image dalam perintah CLI.')
  }

  // Terapkan hasil parse ke form
  form.image = foundImage
  form.tag = foundTag
  form.containerName = foundName || foundImage.split('/').pop().replace(/[^a-z0-9_-]/g, '-')
  form.title = (foundName || form.containerName).replace(/[-_]/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
  form.privileged = parsedPrivileged
  form.restartPolicy = parsedRestart
  form.network = parsedNetwork
  form.memoryLimit = parsedMem
  form.capAdd = parsedCap.join(', ')
  if (parsedPorts.length > 0) form.ports = parsedPorts
  if (parsedVols.length > 0) form.volumes = parsedVols
  if (parsedEnvs.length > 0) form.envs = parsedEnvs
  if (parsedDevices.length > 0) form.devices = parsedDevices
  if (parsedCommands.length > 0) form.commands = parsedCommands

  if (form.ports.length > 0 && form.ports[0].host) {
    form.webUiPort = form.ports[0].host
  }
}

// Helper parser token port
function parsePortToken(val, targetArr) {
  if (!val) return
  let proto = 'TCP'
  let clean = val
  if (clean.toLowerCase().endsWith('/udp')) {
    proto = 'UDP'
    clean = clean.replace(/\/udp$/i, '')
  } else if (clean.toLowerCase().endsWith('/tcp')) {
    clean = clean.replace(/\/tcp$/i, '')
  }
  const parts = clean.split(':')
  if (parts.length >= 2) {
    targetArr.push({ host: parts[0], container: parts[1], protocol: proto })
  } else {
    targetArr.push({ host: parts[0], container: parts[0], protocol: proto })
  }
}

// Helper parser token volume
function parseVolumeToken(val, targetArr) {
  if (!val) return
  const parts = val.split(':')
  if (parts.length >= 2) {
    targetArr.push({ host: parts[0], container: parts[1], mode: parts[2] || 'rw' })
  }
}

// Helper parser token env
function parseEnvToken(val, targetArr) {
  if (!val) return
  const eqIdx = val.indexOf('=')
  if (eqIdx !== -1) {
    targetArr.push({ key: val.substring(0, eqIdx).trim(), value: val.substring(eqIdx + 1).trim() })
  } else {
    targetArr.push({ key: val.trim(), value: '' })
  }
}

// Helper parser token device
function parseDeviceToken(val, targetArr) {
  if (!val) return
  const parts = val.split(':')
  targetArr.push({ host: parts[0] || '', container: parts[1] || parts[0] || '' })
}

// Helper parser memory string (misal: "512m", "1g")
function parseMemoryString(val) {
  if (!val) return 512
  const s = String(val).toUpperCase().trim()
  if (s.endsWith('M') || s.endsWith('MB')) {
    return Math.min(1800, Math.max(256, parseInt(s) || 512))
  } else if (s.endsWith('G') || s.endsWith('GB')) {
    return Math.min(1800, Math.max(256, Math.round(parseFloat(s) * 1024) || 1024))
  }
  const num = parseInt(s)
  return isNaN(num) ? 512 : Math.min(1800, Math.max(256, num))
}

// Eksekusi Submit Import Modal
function handleImportSubmit() {
  importError.value = ''
  try {
    if (importTab.value === 'compose') {
      applyComposeYaml(composeInput.value)
      systemStore.showNotification('Sukses Import', 'Konfigurasi Docker Compose berhasil dimuat ke formulir!', 'success')
    } else {
      applyDockerCLI(cliInput.value)
      systemStore.showNotification('Sukses Import', 'Perintah Docker CLI berhasil dimuat ke formulir!', 'success')
    }
    showImportModal.value = false
  } catch (err) {
    importError.value = err.message || 'Gagal memproses konfigurasi'
  }
}

// ==========================================
// FORM SUBMIT (INSTALL CONTAINER)
// ==========================================
async function handleInstallSubmit() {
  if (!form.image.trim()) {
    systemStore.showNotification('Peringatan', 'Docker Image wajib diisi!', 'warning')
    return
  }
  if (!form.title.trim()) {
    systemStore.showNotification('Peringatan', 'Title aplikasi wajib diisi!', 'warning')
    return
  }

  isDeploying.value = true

  // Susun payload terstruktur kompatibel dengan backend Docker engine
  const fullImage = form.tag ? `${form.image.trim()}:${form.tag.trim()}` : form.image.trim()
  const primaryPort = parseInt(form.webUiPort || form.ports[0]?.host || '8080', 10) || 8080
  const envList = form.envs.filter(e => e.key.trim().length > 0).map(e => `${e.key.trim()}=${e.value.trim()}`)

  const payload = {
    name: form.title.trim(),
    custom_image: fullImage,
    image: form.image.trim(),
    tag: form.tag.trim() || 'latest',
    title: form.title.trim(),
    container_name: form.containerName.trim() || form.title.toLowerCase().replace(/[^a-z0-9_-]/g, '-'),
    icon: form.icon.trim() || 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png',
    web_ui: {
      protocol: form.webUiProtocol,
      host: form.webUiHost,
      port: primaryPort,
      path: form.webUiPath
    },
    port: primaryPort,
    network: form.network,
    ports: form.ports,
    volumes: form.volumes,
    extra_envs: envList,
    envs: form.envs,
    devices: form.devices,
    commands: form.commands.map(c => c.value),
    privileged: form.privileged,
    memory_limit: form.memoryLimit,
    cpu_shares: form.cpuShares,
    restart_policy: form.restartPolicy,
    cap_add: form.capAdd,
    data_path: form.volumes[0]?.host || `/DATA/apps/${form.containerName.trim() || 'custom-app'}`
  }

  try {
    const res = await fetch('/api/v1/store/install', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      systemStore.showNotification('Sukses Dipasang', `Aplikasi ${form.title} berhasil dipasang ke sistem!`, 'success')
      // Refresh global store & desktop
      systemStore.triggerAppsRefresh()
      await windowStore.fetchInstalledApps()
      // Tutup window
      windowStore.closeWindow(props.win?.id || 'appstore-manual')
    } else {
      const errData = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Memasang', errData.error || 'Terjadi kesalahan saat deploy container', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Gagal menghubungi server lokal', 'error')
  } finally {
    isDeploying.value = false
  }
}

// Tutup jendela
function handleCloseWindow() {
  windowStore.closeWindow(props.win?.id || 'appstore-manual')
}
</script>

<template>
  <div class="h-full flex flex-col bg-slate-900 text-slate-100 select-none relative overflow-hidden font-sans">
    <!-- SUB-HEADER TAB BAR & ACTIONS (Sesuai Gambar 1) -->
    <div class="flex items-center justify-between px-6 pt-3 pb-0 border-b border-slate-800 bg-slate-900/90 shrink-0">
      <!-- Tab Aktif main_app -->
      <div class="flex items-center gap-6">
        <div class="relative py-2 text-sm font-semibold text-sky-400 border-b-2 border-sky-500 cursor-pointer">
          main_app
        </div>
      </div>

      <!-- Action Buttons Kanan Atas: Import & Close -->
      <div class="flex items-center gap-2 pb-2">
        <button
          type="button"
          @click="openImportDialog"
          class="p-1.5 rounded-lg text-slate-400 hover:text-sky-300 hover:bg-slate-800 transition flex items-center gap-1.5"
          title="Import Docker Compose / Docker CLI"
        >
          <!-- Ikon Import ->] -->
          <LogIn class="w-4 h-4 text-sky-400" />
          <span class="text-xs font-medium text-slate-300 hidden sm:inline">Import</span>
        </button>

        <div class="w-[1px] h-4 bg-slate-700/80 mx-1"></div>

        <button
          type="button"
          @click="handleCloseWindow"
          class="p-1.5 rounded-lg text-slate-400 hover:text-rose-400 hover:bg-slate-800 transition"
          title="Tutup Jendela"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- MAIN SCROLLABLE FORM BODY (Sesuai Gambar 1 - 4) -->
    <form @submit.prevent="handleInstallSubmit" class="flex-1 overflow-y-auto px-6 py-5 space-y-5 text-xs sm:text-sm">
      <!-- ROW 1: DOCKER IMAGE & TAG -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <!-- Docker Image * (2 Kolom) -->
        <div class="sm:col-span-2">
          <label class="block text-xs font-semibold text-slate-200 mb-1.5">
            Docker Image <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="form.image"
            @blur="onImageBlur"
            type="text"
            required
            placeholder="e.g. hello-world:latest"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition"
          />
        </div>

        <!-- Tag (1 Kolom) -->
        <div>
          <label class="block text-xs font-semibold text-slate-200 mb-1.5">
            Tag
          </label>
          <div class="relative">
            <input
              v-model="form.tag"
              type="text"
              placeholder="latest"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition pr-8"
            />
            <ChevronDown class="w-4 h-4 text-slate-500 absolute right-3 top-3 pointer-events-none" />
          </div>
        </div>
      </div>

      <!-- ROW 2: TITLE * -->
      <div>
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Title <span class="text-rose-400">*</span>
        </label>
        <input
          v-model="form.title"
          @input="onTitleInput"
          type="text"
          required
          placeholder="e.g., Your App Name"
          class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition"
        />
      </div>

      <!-- ROW 3: ICON URL & PREVIEW BOX -->
      <div>
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Icon URL
        </label>
        <div class="flex items-center gap-2.5">
          <!-- Preview Box Icon Sesuai Gambar 1 -->
          <div class="w-11 h-11 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-center shrink-0 overflow-hidden shadow-inner p-1">
            <img
              v-if="form.icon && !iconLoadError"
              :src="form.icon"
              @error="onIconError"
              alt="Icon Preview"
              class="w-full h-full object-contain"
            />
            <!-- Fallback Icon Cube 3D Biru Sesuai Gambar 1 -->
            <div v-else class="w-7 h-7 rounded-lg bg-sky-500/20 border border-sky-400/40 flex items-center justify-center text-sky-400">
              <Box class="w-4 h-4" />
            </div>
          </div>

          <input
            v-model="form.icon"
            @input="onIconInput"
            type="text"
            placeholder="Your custom icon URL"
            class="flex-1 px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition"
          />
        </div>
      </div>

      <!-- ROW 4: WEB UI (SEGMENTED CONTROL) -->
      <div>
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Web UI
        </label>
        <div class="grid grid-cols-12 rounded-xl bg-slate-950 border border-slate-800 overflow-hidden divide-y sm:divide-y-0 sm:divide-x divide-slate-800">
          <!-- Protocol Select -->
          <div class="col-span-12 sm:col-span-3 relative flex items-center">
            <select
              v-model="form.webUiProtocol"
              class="w-full bg-transparent px-3 py-2.5 text-xs sm:text-sm text-slate-100 outline-none appearance-none cursor-pointer"
            >
              <option value="http://" class="bg-slate-900 text-slate-100">http://</option>
              <option value="https://" class="bg-slate-900 text-slate-100">https://</option>
            </select>
            <ChevronDown class="w-3.5 h-3.5 text-slate-500 absolute right-3 pointer-events-none" />
          </div>

          <!-- Hostname -->
          <div class="col-span-12 sm:col-span-4">
            <input
              v-model="form.webUiHost"
              type="text"
              placeholder="server.razos.my.id"
              class="w-full bg-transparent px-3 py-2.5 text-xs sm:text-sm text-slate-100 placeholder-slate-500 outline-none"
            />
          </div>

          <!-- Port -->
          <div class="col-span-12 sm:col-span-3 flex items-center">
            <span class="text-slate-500 pl-3 select-none">:</span>
            <input
              v-model="form.webUiPort"
              type="text"
              placeholder="Ports"
              class="w-full bg-transparent px-1.5 py-2.5 text-xs sm:text-sm text-slate-100 placeholder-slate-500 outline-none font-mono"
            />
          </div>

          <!-- Path -->
          <div class="col-span-12 sm:col-span-2">
            <input
              v-model="form.webUiPath"
              type="text"
              placeholder="/"
              class="w-full bg-transparent px-3 py-2.5 text-xs sm:text-sm text-slate-100 placeholder-slate-500 outline-none font-mono"
            />
          </div>
        </div>
      </div>

      <!-- ROW 5: NETWORK MODE -->
      <div>
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Network
        </label>
        <div class="relative">
          <select
            v-model="form.network"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm outline-none appearance-none cursor-pointer focus:border-sky-400 transition"
          >
            <option value="bridge" class="bg-slate-900 text-slate-100">bridge</option>
            <option value="host" class="bg-slate-900 text-slate-100">host</option>
            <option value="none" class="bg-slate-900 text-slate-100">none</option>
            <option value="macvlan" class="bg-slate-900 text-slate-100">macvlan</option>
          </select>
          <ChevronDown class="w-4 h-4 text-slate-500 absolute right-3.5 top-3.5 pointer-events-none" />
        </div>
      </div>

      <!-- SECTION 6: PORTS (DINAMIS DENGAN + ADD & EMPTY STATE) -->
      <div class="space-y-2 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-200">Port</label>
          <button
            type="button"
            @click="addPort"
            class="px-3 py-1 rounded-full bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium flex items-center gap-1 border border-slate-700 transition active:scale-95"
          >
            <Plus class="w-3.5 h-3.5 text-sky-400" />
            <span>Add</span>
          </button>
        </div>

        <!-- Empty State Sesuai Gambar 2 -->
        <div
          v-if="form.ports.length === 0"
          class="flex items-center gap-2 p-2 text-slate-400 text-xs"
        >
          <div class="w-4 h-4 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-[10px]">
            !
          </div>
          <span>Click &ldquo;+&rdquo; to add one.</span>
        </div>

        <!-- List Baris Port -->
        <div v-else class="space-y-2">
          <div
            v-for="(p, idx) in form.ports"
            :key="idx"
            class="grid grid-cols-12 gap-2 items-center bg-slate-950/60 p-2 rounded-xl border border-slate-800/80"
          >
            <div class="col-span-5">
              <input
                v-model="p.host"
                @input="syncWebUiPort"
                type="number"
                placeholder="Host Port (e.g. 8080)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-4">
              <input
                v-model="p.container"
                type="number"
                placeholder="Container (80)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-2">
              <select
                v-model="p.protocol"
                class="w-full px-2 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none"
              >
                <option value="TCP">TCP</option>
                <option value="UDP">UDP</option>
              </select>
            </div>
            <div class="col-span-1 flex justify-center">
              <button
                type="button"
                @click="removePort(idx)"
                class="p-1 rounded-lg hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 transition"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 7: VOLUMES (DINAMIS DENGAN + ADD & EMPTY STATE) -->
      <div class="space-y-2 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-200">Volumes</label>
          <button
            type="button"
            @click="addVolume"
            class="px-3 py-1 rounded-full bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium flex items-center gap-1 border border-slate-700 transition active:scale-95"
          >
            <Plus class="w-3.5 h-3.5 text-sky-400" />
            <span>Add</span>
          </button>
        </div>

        <!-- Empty State Sesuai Gambar 2 -->
        <div
          v-if="form.volumes.length === 0"
          class="flex items-center gap-2 p-2 text-slate-400 text-xs"
        >
          <div class="w-4 h-4 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-[10px]">
            !
          </div>
          <span>Click &ldquo;+&rdquo; to add one.</span>
        </div>

        <!-- List Baris Volume -->
        <div v-else class="space-y-2">
          <div
            v-for="(v, idx) in form.volumes"
            :key="idx"
            class="grid grid-cols-12 gap-2 items-center bg-slate-950/60 p-2 rounded-xl border border-slate-800/80"
          >
            <div class="col-span-6">
              <input
                v-model="v.host"
                type="text"
                placeholder="Host Path (/DATA/apps/...)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-4">
              <input
                v-model="v.container"
                type="text"
                placeholder="Container Path (/data)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-1">
              <select
                v-model="v.mode"
                class="w-full px-1.5 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-[11px] text-slate-100 font-mono outline-none"
              >
                <option value="rw">rw</option>
                <option value="ro">ro</option>
              </select>
            </div>
            <div class="col-span-1 flex justify-center">
              <button
                type="button"
                @click="removeVolume(idx)"
                class="p-1 rounded-lg hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 transition"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 8: ENVIRONMENT VARIABLES -->
      <div class="space-y-2 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-200">Environment Variables</label>
          <button
            type="button"
            @click="addEnv"
            class="px-3 py-1 rounded-full bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium flex items-center gap-1 border border-slate-700 transition active:scale-95"
          >
            <Plus class="w-3.5 h-3.5 text-sky-400" />
            <span>Add</span>
          </button>
        </div>

        <!-- Empty State Sesuai Gambar 2 -->
        <div
          v-if="form.envs.length === 0"
          class="flex items-center gap-2 p-2 text-slate-400 text-xs"
        >
          <div class="w-4 h-4 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-[10px]">
            !
          </div>
          <span>Click &ldquo;+&rdquo; to add one.</span>
        </div>

        <!-- List Baris Envs -->
        <div v-else class="space-y-2">
          <div
            v-for="(e, idx) in form.envs"
            :key="idx"
            class="grid grid-cols-12 gap-2 items-center bg-slate-950/60 p-2 rounded-xl border border-slate-800/80"
          >
            <div class="col-span-5">
              <input
                v-model="e.key"
                type="text"
                placeholder="KEY (e.g. TZ)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-6">
              <input
                v-model="e.value"
                type="text"
                placeholder="VALUE (e.g. Asia/Jakarta)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-1 flex justify-center">
              <button
                type="button"
                @click="removeEnv(idx)"
                class="p-1 rounded-lg hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 transition"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 9: DEVICES -->
      <div class="space-y-2 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-200">Devices</label>
          <button
            type="button"
            @click="addDevice"
            class="px-3 py-1 rounded-full bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium flex items-center gap-1 border border-slate-700 transition active:scale-95"
          >
            <Plus class="w-3.5 h-3.5 text-sky-400" />
            <span>Add</span>
          </button>
        </div>

        <!-- Empty State Sesuai Gambar 2 -->
        <div
          v-if="form.devices.length === 0"
          class="flex items-center gap-2 p-2 text-slate-400 text-xs"
        >
          <div class="w-4 h-4 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-[10px]">
            !
          </div>
          <span>Click &ldquo;+&rdquo; to add one.</span>
        </div>

        <!-- List Baris Devices -->
        <div v-else class="space-y-2">
          <div
            v-for="(d, idx) in form.devices"
            :key="idx"
            class="grid grid-cols-12 gap-2 items-center bg-slate-950/60 p-2 rounded-xl border border-slate-800/80"
          >
            <div class="col-span-6">
              <input
                v-model="d.host"
                type="text"
                placeholder="Host (e.g. /dev/dri)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-5">
              <input
                v-model="d.container"
                type="text"
                placeholder="Container (/dev/dri)"
                class="w-full px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
              />
            </div>
            <div class="col-span-1 flex justify-center">
              <button
                type="button"
                @click="removeDevice(idx)"
                class="p-1 rounded-lg hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 transition"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 10: CONTAINER COMMAND -->
      <div class="space-y-2 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-200">Container Command</label>
          <button
            type="button"
            @click="addCommand"
            class="px-3 py-1 rounded-full bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium flex items-center gap-1 border border-slate-700 transition active:scale-95"
          >
            <Plus class="w-3.5 h-3.5 text-sky-400" />
            <span>Add</span>
          </button>
        </div>

        <!-- Empty State Sesuai Gambar 2 & 3 -->
        <div
          v-if="form.commands.length === 0"
          class="flex items-center gap-2 p-2 text-slate-400 text-xs"
        >
          <div class="w-4 h-4 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-[10px]">
            !
          </div>
          <span>Click &ldquo;+&rdquo; to add one.</span>
        </div>

        <!-- List Baris Command -->
        <div v-else class="space-y-2">
          <div
            v-for="(c, idx) in form.commands"
            :key="idx"
            class="flex items-center gap-2 bg-slate-950/60 p-2 rounded-xl border border-slate-800/80"
          >
            <input
              v-model="c.value"
              type="text"
              placeholder="e.g. --bind 0.0.0.0 or npm start"
              class="flex-1 px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
            />
            <button
              type="button"
              @click="removeCommand(idx)"
              class="p-1 rounded-lg hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 transition"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>

      <!-- SECTION 11: PRIVILEGES TOGGLE SWITCH (Sesuai Gambar 2 & 3) -->
      <div class="pt-2 border-t border-slate-800/80">
        <label class="block text-xs font-semibold text-slate-200 mb-2">
          Privileges
        </label>
        <button
          type="button"
          @click="form.privileged = !form.privileged"
          class="w-12 h-6 rounded-full transition-colors duration-200 ease-in-out relative p-0.5 focus:outline-none"
          :class="form.privileged ? 'bg-sky-500' : 'bg-slate-700'"
        >
          <span
            class="w-5 h-5 rounded-full bg-white shadow-md transform transition-transform duration-200 ease-in-out block"
            :class="form.privileged ? 'translate-x-6' : 'translate-x-0'"
          ></span>
        </button>
      </div>

      <!-- SECTION 12: MEMORY LIMIT SLIDER (Sesuai Gambar 3: 256, 512, 1024, 1800) -->
      <div class="pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between mb-2">
          <label class="text-xs font-semibold text-slate-200">Memory Limit</label>
          <span class="text-xs font-mono font-bold text-sky-400 bg-sky-500/10 px-2 py-0.5 rounded-md border border-sky-500/20">
            {{ form.memoryLimit }} MB
          </span>
        </div>

        <!-- Slider Bar -->
        <input
          v-model.number="form.memoryLimit"
          type="range"
          min="256"
          max="1800"
          step="64"
          class="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-sky-400"
        />

        <!-- Ticks Label Sesuai Gambar 3: 256, 512, 1024, 1800 -->
        <div class="flex justify-between text-[11px] text-slate-500 font-mono mt-1 px-1">
          <span>256</span>
          <span>512</span>
          <span>1024</span>
          <span>1800</span>
        </div>
      </div>

      <!-- SECTION 13: CPU SHARES (Sesuai Gambar 3 & 4) -->
      <div class="pt-2 border-t border-slate-800/80">
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          CPU Shares
        </label>
        <div class="relative">
          <select
            v-model="form.cpuShares"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm outline-none appearance-none cursor-pointer focus:border-sky-400 transition"
          >
            <option value="Low" class="bg-slate-900 text-slate-100">Low</option>
            <option value="Medium" class="bg-slate-900 text-slate-100">Medium</option>
            <option value="High" class="bg-slate-900 text-slate-100">High</option>
          </select>
          <ChevronDown class="w-4 h-4 text-slate-500 absolute right-3.5 top-3.5 pointer-events-none" />
        </div>
      </div>

      <!-- SECTION 14: RESTART POLICY (Sesuai Gambar 3 & 4) -->
      <div class="pt-2 border-t border-slate-800/80">
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Restart Policy
        </label>
        <div class="relative">
          <select
            v-model="form.restartPolicy"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm outline-none appearance-none cursor-pointer focus:border-sky-400 transition"
          >
            <option value="unless-stopped" class="bg-slate-900 text-slate-100">unless-stopped</option>
            <option value="always" class="bg-slate-900 text-slate-100">always</option>
            <option value="on-failure" class="bg-slate-900 text-slate-100">on-failure</option>
            <option value="no" class="bg-slate-900 text-slate-100">no</option>
          </select>
          <ChevronDown class="w-4 h-4 text-slate-500 absolute right-3.5 top-3.5 pointer-events-none" />
        </div>
      </div>

      <!-- SECTION 15: CONTAINER CAPABILITIES (CAP-ADD) (Sesuai Gambar 3 & 4) -->
      <div class="pt-2 border-t border-slate-800/80">
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Container Capabilities (cap-add)
        </label>
        <input
          v-model="form.capAdd"
          type="text"
          placeholder="e.g. NET_ADMIN, SYS_TIME"
          class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition font-mono"
        />
      </div>

      <!-- SECTION 16: CONTAINER NAME (Sesuai Gambar 4) -->
      <div class="pt-2 border-t border-slate-800/80">
        <label class="block text-xs font-semibold text-slate-200 mb-1.5">
          Container Name
        </label>
        <input
          v-model="form.containerName"
          type="text"
          placeholder="Name of app container"
          class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs sm:text-sm placeholder-slate-500 outline-none focus:border-sky-400 transition font-mono"
        />
      </div>

      <!-- FOOTER ACTION: INSTALL BUTTON (Sesuai Gambar 4) -->
      <div class="flex items-center justify-end pt-4 pb-2">
        <button
          type="submit"
          :disabled="isDeploying"
          class="px-6 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-semibold text-xs sm:text-sm flex items-center gap-2 shadow-lg shadow-blue-600/30 active:scale-95 transition disabled:opacity-50"
        >
          <span v-if="isDeploying" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span>{{ isDeploying ? 'Installing...' : 'Install' }}</span>
        </button>
      </div>
    </form>

    <!-- ============================================================== -->
    <!-- MODAL IMPORT (DOCKER COMPOSE & DOCKER CLI) (Sesuai Gambar 5) -->
    <!-- ============================================================== -->
    <div
      v-if="showImportModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-2xl bg-slate-900 border border-slate-700/80 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <!-- Header Modal Import -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800 bg-slate-950/80">
          <h3 class="text-base font-bold text-slate-100">Import</h3>
          <button
            type="button"
            @click="closeImportDialog"
            class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Tab Switches: Docker Compose & Docker CLI (Sesuai Gambar 5) -->
        <div class="flex items-center gap-6 px-6 pt-3 border-b border-slate-800 bg-slate-900/50">
          <button
            type="button"
            @click="importTab = 'compose'"
            class="py-2 text-sm font-semibold transition border-b-2"
            :class="importTab === 'compose' ? 'text-sky-400 border-sky-500' : 'text-slate-400 border-transparent hover:text-slate-200'"
          >
            Docker Compose
          </button>
          <button
            type="button"
            @click="importTab = 'cli'"
            class="py-2 text-sm font-semibold transition border-b-2"
            :class="importTab === 'cli' ? 'text-sky-400 border-sky-500' : 'text-slate-400 border-transparent hover:text-slate-200'"
          >
            Docker CLI
          </button>
        </div>

        <!-- Body Modal Konten Tab -->
        <div class="p-6 space-y-4">
          <!-- TAB 1: DOCKER COMPOSE YAML -->
          <template v-if="importTab === 'compose'">
            <!-- Textarea YAML -->
            <textarea
              v-model="composeInput"
              rows="7"
              placeholder="When there are more than 1 containers specified in the Docker compose file, only the first one will be imported for now."
              class="w-full px-4 py-3 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition resize-none placeholder-slate-500"
            ></textarea>

            <!-- Dropzone Upload File YAML Sesuai Gambar 5 -->
            <div
              @dragover.prevent="isDragActive = true"
              @dragleave.prevent="isDragActive = false"
              @drop.prevent="handleDrop"
              @click="triggerFileInput"
              class="border-2 border-dashed rounded-xl p-5 text-center cursor-pointer transition flex flex-col items-center justify-center gap-2 group"
              :class="isDragActive ? 'border-sky-400 bg-sky-500/10' : 'border-slate-800 hover:border-slate-700 bg-slate-950/40'"
            >
              <input
                ref="fileInputRef"
                type="file"
                accept=".yaml,.yml,.txt"
                @change="handleFileSelect"
                class="hidden"
              />
              <div class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-slate-300 group-hover:text-sky-400 group-hover:bg-slate-700 transition">
                <Upload class="w-5 h-5" />
              </div>
              <p class="text-xs text-slate-400 group-hover:text-slate-200 transition">
                Drop your Docker Compose here or click to upload
              </p>
            </div>
          </template>

          <!-- TAB 2: DOCKER CLI ('docker run ...') -->
          <template v-else>
            <!-- Textarea CLI Command -->
            <textarea
              v-model="cliInput"
              rows="9"
              placeholder="e.g., docker run -d --name my-app -p 8080:80 -v /DATA/apps/data:/app/data -e TZ=Asia/Jakarta nginx:alpine"
              class="w-full px-4 py-3 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition resize-none placeholder-slate-500"
            ></textarea>

            <!-- Bantuan Preset CLI Cepat -->
            <div class="flex items-center gap-2 overflow-x-auto text-xs text-slate-400 pt-1">
              <span>Contoh Cepat:</span>
              <button
                type="button"
                @click="cliInput = 'docker run -d --name jellyfin -p 8096:8096 -v /DATA/media:/media -e TZ=Asia/Jakarta --restart=unless-stopped jellyfin/jellyfin:latest'"
                class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-mono shrink-0"
              >
                Jellyfin
              </button>
              <button
                type="button"
                @click="cliInput = 'docker run -d --name nginx-web -p 8080:80 -v /DATA/apps/nginx/html:/usr/share/nginx/html nginx:alpine'"
                class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-mono shrink-0"
              >
                Nginx
              </button>
              <button
                type="button"
                @click="cliInput = 'docker run -d --name redis-cache -p 6379:6379 -v /DATA/apps/redis:/data redis:7-alpine redis-server --appendonly yes'"
                class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-mono shrink-0"
              >
                Redis
              </button>
            </div>
          </template>

          <!-- Pesan Error Parse jika ada -->
          <div
            v-if="importError"
            class="flex items-center gap-2 p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-400 text-xs"
          >
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ importError }}</span>
          </div>
        </div>

        <!-- Footer Modal Import (Cancel & Submit Sesuai Gambar 5) -->
        <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-slate-800 bg-slate-950/80">
          <button
            type="button"
            @click="closeImportDialog"
            class="px-5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs sm:text-sm transition"
          >
            Cancel
          </button>
          <button
            type="button"
            @click="handleImportSubmit"
            class="px-6 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-medium text-xs sm:text-sm transition shadow-lg shadow-blue-600/30 active:scale-95"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
