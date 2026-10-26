<script setup>
/**
 * RAZOS App Store Settings & Container Reconfiguration View (CasaOS Enterprise Style)
 * Created: 2026-10-14 (v3.0.0 - Full App Lifecycle Management)
 * Updated: 2026-10-26 (v3.1.2 - Reaktif Form Port Binding & Pencegahan Konflik Port 8080)
 * Deskripsi: Form pengaturan detail aplikasi Docker terpasang (Image & Tag, Title, Icon,
 * Web UI Port, Network, Volume Mounts, Environment Variables, RAM/CPU Limits, dan Restart Policy).
 * Menyimpan konfigurasi dan memicu reboot container otomatis dengan parameter terbaru.
 */
import { ref, computed, onMounted, watch } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Sliders,
  Check,
  RotateCw,
  Terminal,
  ExternalLink,
  X,
  HardDrive,
  Cpu,
  Layers,
  Globe,
  Plus,
  Trash2,
  Box
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

// Data aplikasi yang sedang diedit (dukung unpacking storeData)
const app = computed(() => {
  const d = props.win?.customData || {}
  return d.storeData ? { ...d.storeData, ...d } : d
})
const hostIp = ref(window.location.hostname || '127.0.0.1')

// State Form Pengaturan
const editTitle = ref('')
const editImage = ref('')
const editTag = ref('latest')
const editIcon = ref('')
const editProtocol = ref('http://')
const editPort = ref(0)
const editContainerPort = ref(0)
const editWebPath = ref('/')
const editNetwork = ref('bridge')
const editVolumeHost = ref('')
const editVolumeContainer = ref('/data')
const editMemLimit = ref('256M')
const editCpuShares = ref(512)
const editRestartPolicy = ref('always')
const editEnvs = ref('')
const isSaving = ref(false)

const memoryPresets = [
  { label: '32MB', val: '32M' },
  { label: '64MB', val: '64M' },
  { label: '128MB', val: '128M' },
  { label: '256MB', val: '256M' },
  { label: '512MB', val: '512M' },
  { label: '1GB', val: '1024M' },
  { label: '2GB', val: '2048M' }
]

function populateForm() {
  const a = app.value
  if (!a || (!a.slug && !a.id && !a.name && !a.title)) return

  editTitle.value = a.name || a.title || 'Aplikasi Terpasang'
  
  // Ekstraksi Image & Tag
  const rawImage = a.image || a.custom_image || `${a.slug || 'app'}:latest`
  if (rawImage.includes(':')) {
    const parts = rawImage.split(':')
    editImage.value = parts[0]
    editTag.value = parts[1] || 'latest'
  } else {
    editImage.value = rawImage
    editTag.value = a.tag || 'latest'
  }

  editIcon.value = a.icon || 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png'

  // Tentukan default port dinamis sesuai katalog container (hindari fallback 8080)
  const defaultAppPort = a.default_port || (a.slug === 'immich' ? 2283 : 8090)
  editPort.value = a.port || a.default_port || defaultAppPort
  editContainerPort.value = a.default_port || a.container_port || a.port || defaultAppPort

  editVolumeHost.value = a.data_path || a.default_volume || a.installed_volume || `/DATA/apps/${a.slug || 'app'}`
  editNetwork.value = a.network || 'bridge'
  editMemLimit.value = a.mem_limit || '256M'
  editCpuShares.value = a.cpu_shares || 512
  editRestartPolicy.value = a.restart_policy || 'always'

  // ENV vars
  if (Array.isArray(a.extra_envs) && a.extra_envs.length > 0) {
    editEnvs.value = a.extra_envs.join('\n')
  } else if (typeof a.env_vars === 'string') {
    editEnvs.value = a.env_vars
  } else if (Array.isArray(a.env)) {
    editEnvs.value = a.env.join('\n')
  }
}

// Pantau perubahan props.win.customData secara reaktif agar port dan data langsung terisi
watch(
  () => props.win?.customData,
  (val) => {
    if (val && Object.keys(val).length > 0) {
      populateForm()
    }
  },
  { immediate: true, deep: true }
)

async function handleSaveSettings() {
  if (!editTitle.value.trim() || !editImage.value.trim()) {
    systemStore.showNotification('Peringatan', 'Judul dan Docker Image wajib diisi!', 'warning')
    return
  }

  const portNum = Number(editPort.value)
  if (portNum === 80 || portNum === 8080) {
    systemStore.showNotification(
      'Port Konflik',
      'Port 80 dan 8080 digunakan oleh antarmuka sistem RAZOS. Harap pilih port lain (contoh: 2283 atau 8090)!',
      'warning'
    )
    return
  }

  isSaving.value = true
  const slug = app.value.slug || app.value.id
  const defaultAppPort = app.value.default_port || (slug === 'immich' ? 2283 : 8090)

  const payload = {
    app_id: slug,
    title: editTitle.value.trim(),
    name: editTitle.value.trim(),
    image: `${editImage.value.trim()}:${editTag.value.trim() || 'latest'}`,
    custom_image: `${editImage.value.trim()}:${editTag.value.trim() || 'latest'}`,
    tag: editTag.value.trim() || 'latest',
    icon: editIcon.value.trim(),
    port: portNum || defaultAppPort,
    container_port: Number(editContainerPort.value) || defaultAppPort,
    data_path: editVolumeHost.value.trim(),
    network: editNetwork.value,
    mem_limit: editMemLimit.value,
    cpu_shares: Number(editCpuShares.value) || 512,
    restart_policy: editRestartPolicy.value,
    extra_envs: editEnvs.value
      ? editEnvs.value.split('\n').map(s => s.trim()).filter(s => s.length > 0)
      : []
  }

  try {
    const res = await fetch(`/api/v1/store/apps/${slug}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      systemStore.showNotification(
        'Pengaturan Disimpan',
        `Container "${payload.title}" berhasil diperbarui & dimuat ulang (reboot)!`,
        'success'
      )
      // Sinkronisasi data ke desktop layout dan store apps
      systemStore.triggerAppsRefresh()
      await windowStore.fetchInstalledApps()
      // Tutup jendela setting mandiri
      windowStore.closeWindow(props.win?.id || `app-edit-${slug}`)
    } else {
      const err = await res.json().catch(() => ({}))
      systemStore.showNotification('Gagal Menyimpan', err.error || 'Terjadi kesalahan pada backend', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  } finally {
    isSaving.value = false
  }
}

function openLogs() {
  windowStore.openAppLogs(app.value)
}

function openWebUI() {
  windowStore.launchStoreApp({
    ...app.value,
    port: editPort.value
  })
}

function handleCancel() {
  const slug = app.value.slug || app.value.id
  windowStore.closeWindow(props.win?.id || `app-edit-${slug}`)
}

onMounted(() => {
  populateForm()
})
</script>

<template>
  <form
    @submit.prevent="handleSaveSettings"
    class="h-full flex flex-col justify-between bg-slate-950 font-sans text-slate-100 select-none overflow-hidden"
  >
    <!-- 1. TOP TITLEBAR SUB-HEADER (CasaOS Style with Quick Action Icons) -->
    <div class="flex items-center justify-between px-5 py-3 border-b border-slate-800/90 bg-slate-900/80 shrink-0">
      <div class="flex items-center gap-3 min-w-0">
        <div class="w-8 h-8 rounded-xl bg-slate-800 p-1 flex items-center justify-center border border-slate-700/60 shrink-0">
          <img
            :src="editIcon"
            class="w-full h-full object-contain"
            @error="$event.target.src = 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png'"
          />
        </div>
        <div class="min-w-0">
          <h3 class="font-bold text-sm text-slate-100 truncate tracking-tight">
            {{ editTitle || 'App Settings' }}
          </h3>
          <span class="text-[11px] text-sky-400 font-mono block truncate">
            razos-{{ app.slug || 'container' }}
          </span>
        </div>
      </div>

      <!-- Quick Action Icons in Header -->
      <div class="flex items-center gap-1.5 shrink-0">
        <button
          type="button"
          @click="openLogs"
          class="p-2 rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-sky-300 transition border border-slate-700/60"
          title="Buka Terminal & Logs Container"
        >
          <Terminal class="w-4 h-4" />
        </button>
        <button
          type="button"
          @click="openWebUI"
          class="p-2 rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-emerald-300 transition border border-slate-700/60"
          title="Buka Web UI di Tab Baru"
        >
          <ExternalLink class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- 2. SCROLLABLE SETTINGS FORM BODY (CasaOS Exact Form Layout) -->
    <div class="flex-1 overflow-y-auto px-6 py-4 space-y-4 text-xs">
      <!-- Breadcrumb / Active App Tab -->
      <div class="border-b border-slate-800 pb-2">
        <span class="text-sky-400 font-semibold border-b-2 border-sky-400 pb-2 inline-block px-1">
          {{ app.slug || 'app-configuration' }}
        </span>
      </div>

      <!-- FIELD 1: Docker Image & Tag -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="block font-medium text-slate-300 mb-1">
            Docker Image <span class="text-rose-400">*</span>
          </label>
          <div class="relative">
            <input
              v-model="editImage"
              type="text"
              required
              class="w-full pl-3 pr-8 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
              placeholder="Contoh: immich-app/immich-server"
            />
            <Check class="w-4 h-4 text-emerald-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Tag</label>
          <input
            v-model="editTag"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
            placeholder="latest"
          />
        </div>
      </div>

      <!-- FIELD 2: Title / App Name -->
      <div>
        <label class="block font-medium text-slate-300 mb-1">
          Title / Nama Aplikasi <span class="text-rose-400">*</span>
        </label>
        <div class="relative">
          <input
            v-model="editTitle"
            type="text"
            required
            class="w-full pl-3 pr-8 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400 transition"
            placeholder="Nama tampilan aplikasi di desktop"
          />
          <Check class="w-4 h-4 text-emerald-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>
      </div>

      <!-- FIELD 3: Icon URL & Live Preview -->
      <div>
        <label class="block font-medium text-slate-300 mb-1">Icon URL</label>
        <div class="flex items-center gap-2">
          <div class="w-9 h-9 rounded-xl bg-slate-900 border border-slate-800 p-1 flex items-center justify-center shrink-0">
            <img
              :src="editIcon"
              class="w-full h-full object-contain"
              @error="$event.target.src = 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png'"
            />
          </div>
          <input
            v-model="editIcon"
            type="text"
            class="flex-1 px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400 transition"
            placeholder="https://... atau /assets/icons/..."
          />
        </div>
      </div>

      <!-- FIELD 4: Web UI Access URL -->
      <div>
        <label class="block font-medium text-slate-300 mb-1">Web UI Akses Peramban</label>
        <div class="grid grid-cols-12 gap-2">
          <div class="col-span-3 sm:col-span-2">
            <select
              v-model="editProtocol"
              class="w-full px-2 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
            >
              <option value="http://">http://</option>
              <option value="https://">https://</option>
            </select>
          </div>
          <div class="col-span-5 sm:col-span-6">
            <input
              :value="hostIp"
              disabled
              class="w-full px-3 py-2 rounded-xl bg-slate-900/50 border border-slate-800/80 text-slate-400 font-mono text-xs cursor-not-allowed"
            />
          </div>
          <div class="col-span-3 sm:col-span-3">
            <div class="relative">
              <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-500 font-mono">:</span>
              <input
                v-model.number="editPort"
                type="number"
                class="w-full pl-5 pr-2 py-2 rounded-xl bg-slate-900 border border-slate-800 text-amber-400 font-bold font-mono text-xs outline-none focus:border-sky-400"
                placeholder="Port"
              />
            </div>
          </div>
          <div class="col-span-1 sm:col-span-1">
            <input
              v-model="editWebPath"
              type="text"
              class="w-full px-2 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 font-mono text-xs text-center outline-none focus:border-sky-400"
            />
          </div>
        </div>
      </div>

      <!-- FIELD 5: Network Mode & Port Mapping -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
        <div>
          <label class="block font-medium text-slate-300 mb-1">Network Mode</label>
          <select
            v-model="editNetwork"
            class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
          >
            <option value="bridge">Bridge (Default Isolasi Aman)</option>
            <option value="host">Host Mode (Port Langsung Host)</option>
            <option value="macvlan">Macvlan (Dedicated IP LAN)</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Port Mapping (Host : Container)</label>
          <div class="flex items-center gap-2">
            <input
              v-model.number="editPort"
              type="number"
              class="w-1/2 px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-amber-400 font-mono font-bold text-xs outline-none focus:border-sky-400"
              placeholder="Host Port"
            />
            <span class="text-slate-500 font-mono">:</span>
            <input
              v-model.number="editContainerPort"
              type="number"
              class="w-1/2 px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="Container Port"
            />
          </div>
        </div>
      </div>

      <!-- FIELD 6: Storage Volume Mount (/DATA Persisten) -->
      <div>
        <label class="block font-medium text-slate-300 mb-1">
          Storage Volume Path (/DATA MicroSD / SSD)
        </label>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <div>
            <span class="text-[10px] text-slate-400 block mb-0.5">Path Host STB:</span>
            <input
              v-model="editVolumeHost"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="/DATA/apps/nama-app"
            />
          </div>
          <div>
            <span class="text-[10px] text-slate-400 block mb-0.5">Mount Container:</span>
            <input
              v-model="editVolumeContainer"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 font-mono text-xs outline-none focus:border-sky-400"
              placeholder="/data"
            />
          </div>
        </div>
      </div>

      <!-- FIELD 7: Environment Variables -->
      <div>
        <label class="block font-medium text-slate-300 mb-1">Environment Variables (ENV Vars)</label>
        <textarea
          v-model="editEnvs"
          rows="3"
          placeholder="KEY=VALUE&#10;TZ=Asia/Jakarta&#10;DB_PASSWORD=secret"
          class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-emerald-400 font-mono text-xs outline-none focus:border-sky-400 resize-none leading-relaxed"
        ></textarea>
      </div>

      <!-- FIELD 8: Hardware Resource Limits (STB Friendly) -->
      <div class="space-y-2 pt-1 border-t border-slate-800/80">
        <label class="block font-medium text-slate-300">Batas Memori RAM Container</label>
        <div class="flex items-center gap-2">
          <input
            v-model="editMemLimit"
            type="text"
            class="w-32 px-3 py-1.5 rounded-xl bg-slate-900 border border-slate-800 text-sky-400 font-mono font-bold text-xs outline-none focus:border-sky-400"
          />
          <div class="flex flex-wrap gap-1">
            <button
              v-for="p in memoryPresets"
              :key="p.val"
              type="button"
              @click="editMemLimit = p.val"
              class="px-2 py-1 rounded-lg text-[10px] bg-slate-900 border border-slate-800 hover:border-slate-700 text-slate-300 transition"
              :class="editMemLimit === p.val ? 'bg-sky-500/20 text-sky-300 border-sky-400 font-bold' : ''"
            >
              {{ p.label }}
            </button>
          </div>
        </div>
      </div>

      <!-- Restart Policy -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label class="block font-medium text-slate-300 mb-1">Restart Policy</label>
          <select
            v-model="editRestartPolicy"
            class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
          >
            <option value="always">Always (Selalu Nyala Otomatis)</option>
            <option value="unless-stopped">Unless Stopped (Kecuali Dimatikan Manual)</option>
            <option value="on-failure">On Failure (Nyalakan Ulang Jika Crash)</option>
            <option value="no">No (Jangan Nyalakan Otomatis)</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">CPU Shares</label>
          <input
            v-model.number="editCpuShares"
            type="number"
            class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-100 font-mono text-xs outline-none focus:border-sky-400"
            placeholder="512"
          />
        </div>
      </div>
    </div>

    <!-- 3. FOOTER ACTIONS (CasaOS Blue Save Button) -->
    <div class="flex items-center justify-between px-6 py-3 border-t border-slate-800 bg-slate-900/90 shrink-0">
      <button
        type="button"
        @click="handleCancel"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Batal
      </button>

      <button
        type="submit"
        :disabled="isSaving"
        class="px-6 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs flex items-center gap-1.5 transition shadow-lg shadow-blue-600/25 active:scale-95 disabled:opacity-50"
      >
        <RotateCw v-if="isSaving" class="w-3.5 h-3.5 animate-spin" />
        <Check v-else class="w-3.5 h-3.5 stroke-[2.5]" />
        <span>{{ isSaving ? 'Menyimpan & Reboot...' : 'Save' }}</span>
      </button>
    </div>
  </form>
</template>
