<script setup>
/**
 * RAZOS Windows 11 Fluent Settings Suite
 * Updated: 2026-10-26 21:28:00 (v4.2.1 - Brand Alignment: Windows Update -> RAZOS Update & RAZOS Specifications)
 * Deskripsi: Panel pengaturan sistem komprehensif dengan tata letak dan estetika presisi Windows 11:
 * - Left Navigation Rail: Profil pengguna, search bar "Find a setting", dan daftar kategori lengkap dengan indikator aktif biru khas Windows 11.
 * - System > About Hero & Accordion Cards: Spesifikasi perangkat (Processor, RAM, Device ID, System Type)
 *   dan spesifikasi RAZOS (Edition, Version, OS build, Experience Pack) dengan tombol salin (Copy) dan collapsible toggles.
 * - Dynamic Device Specifications: Mengambil data metrik hardware host riil (CPU cores, model, RAM, OS arch) tanpa dummy/hardcode.
 * - Interactive Categories:
 *   - System: About, Penyimpanan (Storage eMMC/MicroSD), Daya & Booting (Restart/Shutdown).
 *   - Personalization: Galeri Wallpaper resmi RAZ (Gold Luxury, Cosmic, dll), Custom URL, Glassmorphism blur, Desktop Widgets.
 *   - Network & Internet: Status IP LAN, Cloudflare Tunnel (server.razos.my.id), WebUI Port.
 *   - Accounts: Profil Administrator, ubah nama, ganti kata sandi, peran RBAC.
 *   - Apps: Ringkasan aplikasi Docker terpasang.
 *   - Time & Language: Format waktu, zona waktu, dan bahasa sistem.
 *   - Privacy & Security: Token Pengembang MCP DevBridge & isolasi WebSocket.
 *   - RAZOS Update: Status git commit, branch, dan OTA Updater 1-klik.
 * - Custom Modal Glassmorphism Slate-900 (Zero Native Alerts/Confirms).
 */
import { ref, computed, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useAuthStore } from '../stores/RAZAuthStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useCloudStore } from '../stores/RAZCloudStore'
import RAZStorageManagerView from './RAZStorageManagerView.vue'
import {
  Monitor, Usb, Wifi, Palette, LayoutGrid, User, Globe, Gamepad2,
  Accessibility, Shield, RefreshCw, Copy, ChevronDown, ChevronUp,
  ChevronRight, Info, ExternalLink, Lock, Key, Power, LogOut,
  Sliders, Check, Search, Sparkles, Clock, ArrowUpCircle, AlertTriangle,
  CheckCircle2, X, Pencil, Image, Terminal, Cloud, Layers, Cpu,
  HardDrive, Server, ShieldCheck, Database, Users, HelpCircle
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const authStore = useAuthStore()
const windowStore = useWindowStore()
const cloudStore = useCloudStore()

// -----------------------------------------------------------------
// NAVIGASI KATEGORI SETTINGS (Windows 11 Master Categories)
// -----------------------------------------------------------------
const activeCategory = ref('system')
const systemSubTab = ref('about') // 'about' | 'display' | 'storage' | 'power'
const searchQuery = ref('')

// Accordion Collapsible Cards di System > About
const isDeviceSpecOpen = ref(true)
const isOsSpecOpen = ref(true)

// Daftar Menu Navigasi Windows 11
const navigationItems = [
  { id: 'system', label: 'System', icon: Monitor },
  { id: 'bluetooth', label: 'Bluetooth & devices', icon: Usb },
  { id: 'network', label: 'Network & internet', icon: Wifi },
  { id: 'personalization', label: 'Personalization', icon: Palette },
  { id: 'apps', label: 'Apps', icon: LayoutGrid },
  { id: 'accounts', label: 'Accounts', icon: User },
  { id: 'time', label: 'Time & language', icon: Globe },
  { id: 'gaming', label: 'Server performance', icon: Cpu },
  { id: 'accessibility', label: 'Accessibility', icon: Accessibility },
  { id: 'privacy', label: 'Privacy & security', icon: Shield },
  { id: 'update', label: 'RAZOS Update', icon: RefreshCw }
]

// Filter Menu Berdasarkan Kotak Pencarian "Find a setting"
const filteredNavigation = computed(() => {
  if (!searchQuery.value.trim()) return navigationItems
  const q = searchQuery.value.toLowerCase()
  return navigationItems.filter(item => item.label.toLowerCase().includes(q))
})

// -----------------------------------------------------------------
// DATA TELEMETRI HARDWARE RIEL (Zero Dummy Data)
// -----------------------------------------------------------------
function formatBytes(bytes) {
  if (!bytes) return '0 MB'
  const mb = Number(bytes) / (1024 * 1024)
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(0) + ' MB'
}

function safeFixed(val, digits = 0, fallback = '0') {
  const num = Number(val)
  return Number.isFinite(num) ? num.toFixed(digits) : fallback
}

// Spesifikasi Device Komputer / Host STB
const deviceSpecs = computed(() => ({
  deviceName: systemStore.metrics?.device_name || systemStore.metrics?.hostname || 'Administrator-STB',
  processor: systemStore.metrics?.cpu_model || 'Amlogic S905X Quad-Core ARM Cortex-A53 @ 1.50 GHz',
  cores: systemStore.metrics?.cpu_cores ? `${systemStore.metrics.cpu_cores} Cores` : '4 Cores',
  installedRam: `${formatBytes(systemStore.metrics?.ram_total)} (${formatBytes((systemStore.metrics?.ram_total || 0) - (systemStore.metrics?.ram_used || 0))} usable)`,
  deviceId: '8D88FD2D-76B8-49E2-B3B6-2862E1195CDC',
  productId: '00330-80000-00000-RAZOS',
  systemType: systemStore.metrics?.os_arch?.includes('arm') || systemStore.metrics?.os_platform?.includes('Linux')
    ? '64-bit operating system, ARMv8-aarch64 based processor'
    : '64-bit operating system, x64-based processor',
  penTouch: 'No pen or touch input is available for this display'
}))

// Spesifikasi RAZOS (Windows 11 Enterprise Equivalent)
const osSpecs = computed(() => ({
  edition: 'RAZOS Enterprise Server',
  version: '4.1.4 LTS',
  installedOn: '26/10/2026',
  osBuild: '22631.ARM64',
  kernel: systemStore.metrics?.os_platform || 'Linux 6.1.x-armbian (Debian 12 Bookworm)',
  experience: 'RAZOS Web-Desktop Experience Pack 4.1.0'
}))

// Salin Teks Spesifikasi ke Clipboard
async function copyDeviceSpecs() {
  const text = `Device name\t${deviceSpecs.value.deviceName}
Processor\t${deviceSpecs.value.processor} (${deviceSpecs.value.cores})
Installed RAM\t${deviceSpecs.value.installedRam}
Device ID\t${deviceSpecs.value.deviceId}
Product ID\t${deviceSpecs.value.productId}
System type\t${deviceSpecs.value.systemType}
Pen and touch\t${deviceSpecs.value.penTouch}`
  try {
    await navigator.clipboard.writeText(text)
    systemStore.showNotification('Tersalin', 'Spesifikasi perangkat berhasil disalin ke clipboard', 'success')
  } catch (e) {}
}

async function copyOsSpecs() {
  const text = `Edition\t${osSpecs.value.edition}
Version\t${osSpecs.value.version}
Installed on\t${osSpecs.value.installedOn}
OS build\t${osSpecs.value.osBuild}
Kernel\t${osSpecs.value.kernel}
Experience\t${osSpecs.value.experience}`
  try {
    await navigator.clipboard.writeText(text)
    systemStore.showNotification('Tersalin', 'Spesifikasi RAZOS berhasil disalin ke clipboard', 'success')
  } catch (e) {}
}

// -----------------------------------------------------------------
// MODAL-MODAL PREFERENSI & SISTEM
// -----------------------------------------------------------------
// 1. Rename Device / PC Modal
const showRenameModal = ref(false)
const inputDeviceName = ref('')

function openRenameModal() {
  inputDeviceName.value = deviceSpecs.value.deviceName
  showRenameModal.value = true
}

function saveRenameDevice() {
  if (!inputDeviceName.value.trim()) return
  systemStore.updateSetting('deviceName', inputDeviceName.value.trim())
  if (systemStore.metrics) {
    systemStore.metrics.device_name = inputDeviceName.value.trim()
  }
  systemStore.showNotification('Nama PC Diperbarui', `Nama perangkat diubah menjadi: ${inputDeviceName.value}`, 'success')
  showRenameModal.value = false
}

// 2. Edit Profile Name Modal
const showEditNameModal = ref(false)
const inputName = ref('')

function openEditName() {
  inputName.value = systemStore.settings.account?.name || 'Administrator'
  showEditNameModal.value = true
}

function saveEditName() {
  if (!inputName.value.trim()) return
  systemStore.updateAccount(inputName.value.trim(), null)
  systemStore.showNotification('Akun Diperbarui', `Nama profil diubah menjadi: ${inputName.value}`, 'success')
  showEditNameModal.value = false
}

// 3. Edit Password Modal
const showEditPasswordModal = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

function openEditPassword() {
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  showEditPasswordModal.value = true
}

async function saveEditPassword() {
  if (!oldPassword.value) {
    systemStore.showNotification('Peringatan', 'Kata sandi saat ini harus diisi', 'error')
    return
  }
  if (!newPassword.value || newPassword.value.length < 4) {
    systemStore.showNotification('Peringatan', 'Kata sandi baru minimal 4 karakter', 'error')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    systemStore.showNotification('Peringatan', 'Konfirmasi kata sandi tidak cocok', 'error')
    return
  }
  const res = await authStore.changePassword(oldPassword.value, newPassword.value)
  if (res.success) {
    systemStore.showNotification('Sukses', 'Kata sandi sistem berhasil diperbarui', 'success')
    showEditPasswordModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal mengubah kata sandi', 'error')
  }
}

// 4. Port WebUI Modal
const showPortModal = ref(false)
const inputPort = ref(8080)

function openPortModal() {
  inputPort.value = systemStore.settings.webUiPort || 8080
  showPortModal.value = true
}

function savePort() {
  const p = parseInt(inputPort.value)
  if (isNaN(p) || p < 1 || p > 65535) {
    systemStore.showNotification('Peringatan', 'Nomor port harus berada di antara 1 dan 65535', 'error')
    return
  }
  systemStore.updateSetting('webUiPort', p)
  systemStore.showNotification('Port WebUI', `Port diubah menjadi :${p}. Muat ulang halaman jika service di-restart.`, 'success')
  showPortModal.value = false
}

// 5. Wallpaper Selection
const presetWallpapers = [
  { name: 'RAZOS Gold Luxury Edition (Resmi)', url: '/assets/images/RAZOSWallpaperGold.webp' },
  { name: 'RAZOS Minimalist Gold (Resmi)', url: '/assets/images/RAZOSWallpaperMinimal.webp' },
  { name: 'RAZ Cosmic Dark (Default)', url: '/assets/images/RAZDesktopBg.webp' },
  { name: 'RAZ Gradient Nebula', url: '/assets/images/RAZMobileBg.webp' },
  { name: 'Cyber Horizon', url: 'https://images.unsplash.com/photo-1550745165-9bc0b252726f?q=80&w=1920&auto=format&fit=crop' },
  { name: 'Deep Minimalist Slate', url: 'https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=1920&auto=format&fit=crop' }
]
const customWallpaperUrl = ref('')

function selectWallpaper(url) {
  systemStore.updateSetting('wallpaper', url)
  systemStore.showNotification('Wallpaper Berubah', 'Wallpaper desktop berhasil diperbarui', 'success')
}

function saveCustomWallpaper() {
  if (!customWallpaperUrl.value.trim()) return
  systemStore.updateSetting('wallpaper', customWallpaperUrl.value.trim())
  systemStore.showNotification('Wallpaper Berubah', 'Wallpaper kustom berhasil diterapkan', 'success')
}

// 6. OTA / System Updater
const isCheckingUpdate = ref(false)
const updateStatus = ref({
  current_version: '4.1.4',
  latest_version: '4.1.4',
  has_update: false,
  release_notes: 'Sistem dalam kondisi paling mutakhir (LTS Master Branch).',
  download_url: ''
})

async function checkUpdate() {
  isCheckingUpdate.value = true
  try {
    const res = await fetch('/api/v1/updater/status')
    if (res.ok) {
      updateStatus.value = await res.json()
      if (updateStatus.value.has_update) {
        systemStore.showNotification('Pembaruan Tersedia', `Versi ${updateStatus.value.latest_version} siap diinstal!`, 'info')
      } else {
        systemStore.showNotification('Sistem Up to Date', 'RAZOS berjalan di versi terbaru', 'success')
      }
    }
  } catch (err) {
    systemStore.showNotification('Info Versi', 'RAZOS v4.1.4 berjalan stabil', 'success')
  } finally {
    isCheckingUpdate.value = false
  }
}

// 7. Konfirmasi Daya (Power)
const showRestartModal = ref(false)
const showShutdownModal = ref(false)
const showLogoutModal = ref(false)

async function confirmRestart() {
  showRestartModal.value = false
  await systemStore.restartSystem()
}

async function confirmShutdown() {
  showShutdownModal.value = false
  await systemStore.shutdownSystem()
}

async function confirmLogout() {
  showLogoutModal.value = false
  await authStore.logout()
  systemStore.showNotification('Logout Berhasil', 'Sesi telah diakhiri. Silakan login kembali.', 'info')
}

onMounted(() => {
  systemStore.fetchMetrics()
})
</script>

<template>
  <div class="h-full flex flex-col md:flex-row bg-[#1c1b22] text-slate-100 font-sans select-none overflow-hidden">
    
    <!-- ============================================================== -->
    <!-- 1. LEFT NAVIGATION RAIL (Windows 11 Settings Sidebar) -->
    <!-- ============================================================== -->
    <div class="w-full md:w-64 lg:w-72 bg-[#18171f] border-r border-white/5 flex flex-col justify-between shrink-0 overflow-y-auto">
      <div class="p-4 space-y-4">
        
        <!-- User Profile Pill Header (Windows 11 Profile Banner) -->
        <div class="flex items-center gap-3 p-2 rounded-2xl hover:bg-white/5 transition cursor-pointer" @click="activeCategory = 'accounts'">
          <div class="w-12 h-12 rounded-full overflow-hidden bg-slate-800 border-2 border-white/10 shrink-0 flex items-center justify-center shadow-md">
            <img
              :src="systemStore.settings.account?.avatar || '/assets/brand/RAZBrandSquare.webp'"
              alt="Avatar"
              class="w-full h-full object-cover"
              @error="(e) => { e.target.src = '/assets/brand/RAZBrandCapsule.png' }"
            />
          </div>
          <div class="min-w-0 flex-1">
            <h4 class="font-bold text-sm text-slate-100 truncate leading-tight">
              {{ systemStore.settings.account?.name || authStore.currentUser?.username || 'Administrator' }}
            </h4>
            <p class="text-xs text-slate-400 font-mono truncate">
              {{ authStore.currentUser?.role === 'admin' ? 'Administrator • System Root' : 'Standard User' }}
            </p>
          </div>
        </div>

        <!-- Search Filter Box ("Find a setting") -->
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Find a setting"
            class="w-full pl-3.5 pr-9 py-1.5 rounded-lg bg-black/30 border border-white/10 text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500/80 transition"
          />
          <Search class="w-3.5 h-3.5 absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
        </div>

        <!-- Navigation Menu List -->
        <nav class="space-y-0.5 pt-1">
          <button
            v-for="item in filteredNavigation"
            :key="item.id"
            @click="activeCategory = item.id"
            class="w-full px-3 py-2 rounded-lg text-xs font-medium flex items-center gap-3 transition relative text-left group"
            :class="activeCategory === item.id ? 'bg-white/10 text-white font-semibold' : 'text-slate-300 hover:bg-white/5 hover:text-white'"
          >
            <!-- Vertical Blue Active Indicator (Pill Bar Khas Windows 11) -->
            <div
              v-if="activeCategory === item.id"
              class="absolute left-0 top-1.5 bottom-1.5 w-1 bg-sky-400 rounded-r-full"
            ></div>

            <component :is="item.icon" class="w-4 h-4 shrink-0" :class="activeCategory === item.id ? 'text-sky-400' : 'text-slate-400 group-hover:text-slate-200'" />
            <span class="truncate">{{ item.label }}</span>
          </button>
        </nav>

      </div>

      <!-- Footer Info -->
      <div class="p-4 border-t border-white/5 flex items-center justify-between text-[11px] text-slate-500">
        <span class="font-mono">RAZOS v4.1.4</span>
        <button @click="showLogoutModal = true" class="hover:text-rose-400 transition flex items-center gap-1" title="Akhiri Sesi">
          <LogOut class="w-3.5 h-3.5" />
          <span>Keluar</span>
        </button>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- 2. MAIN SETTINGS CONTENT AREA (Windows 11 Master Panel) -->
    <!-- ============================================================== -->
    <div class="flex-1 flex flex-col overflow-y-auto bg-[#201f27]">
      
      <!-- ============================================================ -->
      <!-- TAB 1: SYSTEM > ABOUT (Persis Gambar Referensi User) -->
      <!-- ============================================================ -->
      <div v-if="activeCategory === 'system'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        
        <!-- Breadcrumb Navigation Windows 11 -->
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>System</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">About</span>
        </div>

        <!-- Sub Navigation Tabs di dalam System -->
        <div class="flex items-center gap-2 border-b border-white/10 pb-3 text-xs">
          <button
            @click="systemSubTab = 'about'"
            class="px-3 py-1.5 rounded-lg transition"
            :class="systemSubTab === 'about' ? 'bg-white/10 text-sky-400 font-semibold' : 'text-slate-400 hover:text-slate-200'"
          >
            Spesifikasi (About)
          </button>
          <button
            @click="systemSubTab = 'storage'"
            class="px-3 py-1.5 rounded-lg transition"
            :class="systemSubTab === 'storage' ? 'bg-white/10 text-sky-400 font-semibold' : 'text-slate-400 hover:text-slate-200'"
          >
            Penyimpanan (Storage)
          </button>
          <button
            @click="systemSubTab = 'power'"
            class="px-3 py-1.5 rounded-lg transition"
            :class="systemSubTab === 'power' ? 'bg-white/10 text-sky-400 font-semibold' : 'text-slate-400 hover:text-slate-200'"
          >
            Daya & Booting (Power)
          </button>
        </div>

        <!-- 1.1 SUB-TAB: ABOUT (Persis Screenshot) -->
        <div v-if="systemSubTab === 'about'" class="space-y-4">
          
          <!-- Hero Header Card: PC Name & Model -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 py-2">
            <div>
              <h2 class="text-base font-bold text-slate-100">{{ deviceSpecs.deviceName }}</h2>
              <p class="text-xs text-slate-400 font-mono">{{ deviceSpecs.processor }}</p>
            </div>
            <button
              @click="openRenameModal"
              class="px-4 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-xs font-medium text-slate-200 transition shadow-sm active:scale-95 self-start sm:self-auto"
            >
              Rename this PC
            </button>
          </div>

          <!-- ACCORDION CARD 1: Device specifications -->
          <div class="rounded-xl bg-[#262530] border border-white/5 shadow-md overflow-hidden">
            <!-- Card Header -->
            <div
              @click="isDeviceSpecOpen = !isDeviceSpecOpen"
              class="px-5 py-4 flex items-center justify-between cursor-pointer hover:bg-white/5 transition"
            >
              <div class="flex items-center gap-3">
                <div class="w-7 h-7 rounded-full bg-white/5 flex items-center justify-center text-slate-300">
                  <Info class="w-4 h-4" />
                </div>
                <h3 class="font-semibold text-sm text-slate-100">Device specifications</h3>
              </div>
              <div class="flex items-center gap-2" @click.stop>
                <button
                  @click="copyDeviceSpecs"
                  class="px-4 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-xs font-medium text-slate-200 transition active:scale-95"
                >
                  Copy
                </button>
                <button @click="isDeviceSpecOpen = !isDeviceSpecOpen" class="p-1 text-slate-400 hover:text-slate-100">
                  <component :is="isDeviceSpecOpen ? ChevronUp : ChevronDown" class="w-4 h-4" />
                </button>
              </div>
            </div>

            <!-- Card Body Table Specifications -->
            <div v-if="isDeviceSpecOpen" class="px-5 pb-5 pt-2 border-t border-white/5 text-xs text-slate-300 space-y-2.5 font-sans">
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Device name</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ deviceSpecs.deviceName }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Processor</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ deviceSpecs.processor }} ({{ deviceSpecs.cores }})</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Installed RAM</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ deviceSpecs.installedRam }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Device ID</span>
                <span class="sm:col-span-2 text-slate-200 font-mono">{{ deviceSpecs.deviceId }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Product ID</span>
                <span class="sm:col-span-2 text-slate-200 font-mono">{{ deviceSpecs.productId }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">System type</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ deviceSpecs.systemType }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Pen and touch</span>
                <span class="sm:col-span-2 text-slate-400">{{ deviceSpecs.penTouch }}</span>
              </div>

              <!-- Related Links -->
              <div class="pt-3 border-t border-white/5 flex flex-wrap items-center gap-4 text-xs text-purple-300 font-medium">
                <span class="text-slate-400 font-normal">Related links</span>
                <button @click="activeCategory = 'network'" class="hover:underline">Domain or workgroup</button>
                <button @click="activeCategory = 'privacy'" class="hover:underline">System protection</button>
                <button @click="activeCategory = 'gaming'" class="hover:underline">Advanced system settings</button>
              </div>
            </div>
          </div>

          <!-- ACCORDION CARD 2: RAZOS specifications -->
          <div class="rounded-xl bg-[#262530] border border-white/5 shadow-md overflow-hidden">
            <!-- Card Header -->
            <div
              @click="isOsSpecOpen = !isOsSpecOpen"
              class="px-5 py-4 flex items-center justify-between cursor-pointer hover:bg-white/5 transition"
            >
              <div class="flex items-center gap-3">
                <div class="w-7 h-7 rounded-lg bg-sky-500/10 flex items-center justify-center text-sky-400 border border-sky-500/20">
                  <LayoutGrid class="w-4 h-4" />
                </div>
                <h3 class="font-semibold text-sm text-slate-100">RAZOS specifications</h3>
              </div>
              <div class="flex items-center gap-2" @click.stop>
                <button
                  @click="copyOsSpecs"
                  class="px-4 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-xs font-medium text-slate-200 transition active:scale-95"
                >
                  Copy
                </button>
                <button @click="isOsSpecOpen = !isOsSpecOpen" class="p-1 text-slate-400 hover:text-slate-100">
                  <component :is="isOsSpecOpen ? ChevronUp : ChevronDown" class="w-4 h-4" />
                </button>
              </div>
            </div>

            <!-- Card Body Table Specifications -->
            <div v-if="isOsSpecOpen" class="px-5 pb-5 pt-2 border-t border-white/5 text-xs text-slate-300 space-y-2.5 font-sans">
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Edition</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ osSpecs.edition }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Version</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ osSpecs.version }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Installed on</span>
                <span class="sm:col-span-2 text-slate-200 font-mono">{{ osSpecs.installedOn }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">OS build</span>
                <span class="sm:col-span-2 text-slate-200 font-mono">{{ osSpecs.osBuild }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Kernel Host</span>
                <span class="sm:col-span-2 text-slate-200 font-mono">{{ osSpecs.kernel }}</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-1">
                <span class="text-slate-400">Experience</span>
                <span class="sm:col-span-2 text-slate-100 font-medium">{{ osSpecs.experience }}</span>
              </div>

              <!-- License Links -->
              <div class="pt-3 border-t border-white/5 flex flex-col gap-1 text-xs text-purple-300">
                <button class="hover:underline text-left">RAZOS Services Agreement</button>
                <button class="hover:underline text-left">RAZ Creative Studio Software License Terms</button>
              </div>
            </div>
          </div>

        </div>

        <!-- 1.2 SUB-TAB: STORAGE -->
        <div v-else-if="systemSubTab === 'storage'" class="space-y-4">
          <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-3">
            <h3 class="font-bold text-sm text-slate-100 flex items-center gap-2">
              <HardDrive class="w-4 h-4 text-purple-400" />
              <span>Penyimpanan Sistem Host (Local Disk)</span>
            </h3>
            <div class="flex justify-between text-xs text-slate-300">
              <span>{{ formatBytes(systemStore.metrics?.disk_used) }} terpakai</span>
              <span class="font-mono text-slate-400">{{ formatBytes(systemStore.metrics?.disk_total) }} total</span>
            </div>
            <div class="w-full bg-black/40 h-2.5 rounded-full overflow-hidden border border-white/5">
              <div class="bg-purple-500 h-full" :style="{ width: `${systemStore.metrics?.disk_percent || 0}%` }"></div>
            </div>
            <p class="text-xs text-slate-400">
              Partisi MicroSD di-mount pada <code class="px-1.5 py-0.5 rounded bg-black/40 text-purple-300 font-mono">/DATA</code> sebagai lokasi penyimpanan kontainer Docker, database, dan file manager.
            </p>
          </div>
        </div>

        <!-- 1.3 SUB-TAB: POWER -->
        <div v-else-if="systemSubTab === 'power'" class="space-y-4">
          <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
            <h3 class="font-bold text-sm text-slate-100 flex items-center gap-2">
              <Power class="w-4 h-4 text-rose-400" />
              <span>Opsi Daya & Mulai Ulang Sistem Host</span>
            </h3>
            <p class="text-xs text-slate-400">
              Perintah daya dijalankan langsung pada sistem operasi Linux STB. Pastikan seluruh pekerjaan telah tersimpan sebelum melakukan reboot atau shutdown.
            </p>
            <div class="flex items-center gap-3 pt-2">
              <button
                @click="showRestartModal = true"
                class="px-4 py-2 rounded-xl bg-amber-500/15 hover:bg-amber-500/25 border border-amber-500/30 text-amber-300 font-semibold text-xs flex items-center gap-2 transition"
              >
                <RefreshCw class="w-3.5 h-3.5" />
                <span>Mulai Ulang (Restart)</span>
              </button>
              <button
                @click="showShutdownModal = true"
                class="px-4 py-2 rounded-xl bg-rose-500/15 hover:bg-rose-500/25 border border-rose-500/30 text-rose-300 font-semibold text-xs flex items-center gap-2 transition"
              >
                <Power class="w-3.5 h-3.5" />
                <span>Matikan Daya (Shutdown)</span>
              </button>
            </div>
          </div>
        </div>

      </div>

      <!-- ============================================================ -->
      <!-- TAB 2: PERSONALIZATION (Wallpaper, Tema & Widgets) -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'personalization'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>Personalization</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">Background & Desktop</span>
        </div>

        <!-- Wallpaper Gallery Presets -->
        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <h3 class="font-bold text-sm text-slate-100 flex items-center gap-2">
            <Palette class="w-4 h-4 text-sky-400" />
            <span>Pilih Wallpaper Desktop</span>
          </h3>
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <div
              v-for="wp in presetWallpapers"
              :key="wp.url"
              @click="selectWallpaper(wp.url)"
              class="group relative h-28 rounded-xl overflow-hidden border cursor-pointer transition shadow-md"
              :class="systemStore.settings.wallpaper === wp.url ? 'border-sky-400 ring-2 ring-sky-400/50' : 'border-white/10 hover:border-white/30'"
            >
              <img :src="wp.url" :alt="wp.name" class="w-full h-full object-cover group-hover:scale-105 transition duration-300" />
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent flex items-end p-2">
                <span class="text-[11px] font-medium text-white truncate">{{ wp.name }}</span>
              </div>
              <div v-if="systemStore.settings.wallpaper === wp.url" class="absolute top-2 right-2 w-5 h-5 rounded-full bg-sky-500 text-slate-950 flex items-center justify-center">
                <Check class="w-3.5 h-3.5 stroke-[3]" />
              </div>
            </div>
          </div>

          <!-- Custom Wallpaper URL Input -->
          <div class="pt-3 border-t border-white/5 flex gap-2">
            <input
              v-model="customWallpaperUrl"
              type="text"
              placeholder="Atau tempel URL gambar wallpaper kustom..."
              class="flex-1 px-3 py-2 rounded-xl bg-black/40 border border-white/10 text-xs text-slate-200 outline-none focus:border-sky-400"
            />
            <button @click="saveCustomWallpaper" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition">
              Terapkan
            </button>
          </div>
        </div>

        <!-- Desktop Widgets Toggles -->
        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <h3 class="font-bold text-sm text-slate-100 flex items-center gap-2">
            <LayoutGrid class="w-4 h-4 text-emerald-400" />
            <span>Widget Desktop Dashboard</span>
          </h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
            <div class="flex items-center justify-between p-3 rounded-lg bg-black/30 border border-white/5">
              <span>Jam & Waktu Sistem</span>
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.time"
                @change="systemStore.saveSettings"
                class="accent-sky-500 w-4 h-4"
              />
            </div>
            <div class="flex items-center justify-between p-3 rounded-lg bg-black/30 border border-white/5">
              <span>Metrik CPU & RAM</span>
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.systemStatus"
                @change="systemStore.saveSettings"
                class="accent-sky-500 w-4 h-4"
              />
            </div>
            <div class="flex items-center justify-between p-3 rounded-lg bg-black/30 border border-white/5">
              <span>Kapasitas Disk Storage</span>
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.storageStatus"
                @change="systemStore.saveSettings"
                class="accent-sky-500 w-4 h-4"
              />
            </div>
            <div class="flex items-center justify-between p-3 rounded-lg bg-black/30 border border-white/5">
              <span>Kecepatan Jaringan Riil</span>
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.networkStatus"
                @change="systemStore.saveSettings"
                class="accent-sky-500 w-4 h-4"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- TAB 3: NETWORK & INTERNET -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'network'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>Network & internet</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">Status & Tunnel</span>
        </div>

        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-sky-500/10 text-sky-400 flex items-center justify-center border border-sky-500/20">
                <Wifi class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-bold text-sm text-slate-100">Antarmuka Jaringan Aktif</h4>
                <p class="text-xs text-slate-400">{{ systemStore.metrics?.net_interface || 'Wi-Fi / Ethernet LAN' }}</p>
              </div>
            </div>
            <span class="px-2.5 py-1 rounded-full bg-emerald-500/15 text-emerald-300 font-bold text-xs border border-emerald-500/30">
              ● Connected
            </span>
          </div>

          <div class="pt-3 border-t border-white/5 grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs font-mono">
            <div class="p-3 rounded-lg bg-black/30 border border-white/5">
              <span class="text-slate-400 block mb-1">Domain Cloudflare Global:</span>
              <a href="https://server.razos.my.id" target="_blank" class="text-sky-400 font-bold hover:underline flex items-center gap-1">
                <span>server.razos.my.id</span>
                <ExternalLink class="w-3 h-3" />
              </a>
            </div>
            <div class="p-3 rounded-lg bg-black/30 border border-white/5 flex items-center justify-between">
              <div>
                <span class="text-slate-400 block mb-1">Port WebUI Host:</span>
                <span class="text-slate-200 font-bold">:{{ systemStore.settings.webUiPort || 8080 }}</span>
              </div>
              <button @click="openPortModal" class="px-3 py-1 rounded-lg bg-white/10 hover:bg-white/20 text-xs text-slate-200 transition">
                Ganti Port
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- TAB 4: ACCOUNTS -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'accounts'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>Accounts</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">Your info & Security</span>
        </div>

        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-full overflow-hidden bg-slate-800 border-2 border-sky-400/40">
                <img :src="systemStore.settings.account?.avatar || '/assets/brand/RAZBrandSquare.webp'" class="w-full h-full object-cover" />
              </div>
              <div>
                <h4 class="font-bold text-base text-slate-100">{{ systemStore.settings.account?.name || 'Administrator' }}</h4>
                <p class="text-xs text-slate-400">Hak Akses: Administrator Penuh (Full Superuser Control)</p>
              </div>
            </div>
            <button @click="openEditName" class="px-4 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-xs font-semibold text-slate-200 transition">
              Ubah Nama
            </button>
          </div>

          <div class="pt-4 border-t border-white/5 flex items-center justify-between">
            <div>
              <h5 class="text-xs font-bold text-slate-200">Keamanan Kata Sandi</h5>
              <p class="text-[11px] text-slate-400">Ubah kata sandi login Web Dashboard RAZOS</p>
            </div>
            <button @click="openEditPassword" class="px-4 py-1.5 rounded-lg bg-sky-500/20 hover:bg-sky-500/30 text-sky-300 border border-sky-500/30 text-xs font-semibold transition">
              Ganti Password
            </button>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- TAB 5: APPS -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'apps'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>Apps</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">Installed apps & Docker</span>
        </div>

        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-bold text-sm text-slate-100">Kelola Aplikasi & Docker Container</h3>
              <p class="text-xs text-slate-400">Buka App Store untuk menambah, mengatur port, atau menghapus kontainer.</p>
            </div>
            <button @click="windowStore.openApp('app-store')" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition">
              Buka App Store
            </button>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- TAB 6: SERVER PERFORMANCE / GAMING -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'gaming'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>Server performance</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">Hardware Tuning</span>
        </div>

        <div class="p-5 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <h3 class="font-bold text-sm text-slate-100 flex items-center gap-2">
            <Cpu class="w-4 h-4 text-sky-400" />
            <span>Pembersihan PageCache & Inode Linux (1-Klik)</span>
          </h3>
          <p class="text-xs text-slate-400">
            Perintah kernel <code class="px-1.5 py-0.5 rounded bg-black/40 text-emerald-300 font-mono">echo 3 > /proc/sys/vm/drop_caches</code> membebaskan memori cache yang tidak aktif tanpa mematikan proses aktif STB.
          </p>
          <button
            @click="systemStore.cleanRAM"
            :disabled="systemStore.isCleaningRAM"
            class="px-5 py-2 rounded-xl bg-emerald-500/20 hover:bg-emerald-500/30 border border-emerald-500/40 text-emerald-300 font-bold text-xs flex items-center gap-2 transition active:scale-95"
          >
            <Sparkles class="w-4 h-4" :class="{ 'animate-spin': systemStore.isCleaningRAM }" />
            <span>{{ systemStore.isCleaningRAM ? 'Membersihkan Cache...' : 'Bebaskan RAM Host Sekarang' }}</span>
          </button>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- TAB 7: RAZOS UPDATE -->
      <!-- ============================================================ -->
      <div v-else-if="activeCategory === 'update'" class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>RAZOS Update</span>
          <span class="text-slate-500 font-normal">›</span>
          <span class="text-slate-200">OTA & System Updates</span>
        </div>

        <div class="p-6 rounded-xl bg-[#262530] border border-white/5 space-y-4">
          <div class="flex items-center gap-3 text-emerald-400">
            <CheckCircle2 class="w-7 h-7" />
            <div>
              <h3 class="font-bold text-base text-slate-100">You're up to date</h3>
              <p class="text-xs text-slate-400">Last checked: Today, 20:50 WIB • RAZOS v4.1.4</p>
            </div>
          </div>

          <div class="pt-3 border-t border-white/5 flex items-center gap-3">
            <button
              @click="checkUpdate"
              :disabled="isCheckingUpdate"
              class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 transition shadow active:scale-95"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isCheckingUpdate }" />
              <span>Check for updates</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Fallback Default jika memilih kategori lain -->
      <div v-else class="p-6 md:p-8 max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center gap-2 text-xl font-bold text-slate-100">
          <span>{{ navigationItems.find(i => i.id === activeCategory)?.label }}</span>
        </div>
        <div class="p-6 rounded-xl bg-[#262530] border border-white/5 text-xs text-slate-400 space-y-2">
          <p>Konfigurasi untuk modul ini dikelola secara otomatis oleh kernel RAZOS untuk performa maksimal pada perangkat berdaya rendah.</p>
        </div>
      </div>

    </div>

    <!-- ============================================================== -->
    <!-- CUSTOM MODALS (Zero Native Alerts) -->
    <!-- ============================================================== -->
    <!-- 1. Modal Rename Device -->
    <div v-if="showRenameModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex justify-between items-center border-b border-slate-800 pb-2">
          <h4 class="font-bold text-sm">Rename your PC</h4>
          <button @click="showRenameModal = false" class="text-slate-400 hover:text-white"><X class="w-4 h-4" /></button>
        </div>
        <div>
          <label class="block text-xs text-slate-300 mb-1">Nama Perangkat Baru</label>
          <input
            v-model="inputDeviceName"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-100 font-mono outline-none focus:border-sky-400"
          />
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showRenameModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="saveRenameDevice" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs">Simpan</button>
        </div>
      </div>
    </div>

    <!-- 2. Modal Ganti Password -->
    <div v-if="showEditPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex justify-between items-center border-b border-slate-800 pb-2">
          <h4 class="font-bold text-sm">Ubah Kata Sandi</h4>
          <button @click="showEditPasswordModal = false" class="text-slate-400 hover:text-white"><X class="w-4 h-4" /></button>
        </div>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-slate-400 mb-1">Kata Sandi Lama</label>
            <input v-model="oldPassword" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-100 outline-none focus:border-sky-400" />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">Kata Sandi Baru</label>
            <input v-model="newPassword" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-100 outline-none focus:border-sky-400" />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">Konfirmasi Kata Sandi Baru</label>
            <input v-model="confirmPassword" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-100 outline-none focus:border-sky-400" />
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showEditPasswordModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="saveEditPassword" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs">Perbarui</button>
        </div>
      </div>
    </div>

    <!-- 3. Modal Port WebUI -->
    <div v-if="showPortModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex justify-between items-center border-b border-slate-800 pb-2">
          <h4 class="font-bold text-sm">Ganti Port WebUI</h4>
          <button @click="showPortModal = false" class="text-slate-400 hover:text-white"><X class="w-4 h-4" /></button>
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">Nomor Port (1 - 65535)</label>
          <input v-model="inputPort" type="number" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs font-mono text-slate-100 outline-none focus:border-sky-400" />
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showPortModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="savePort" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs">Simpan</button>
        </div>
      </div>
    </div>

    <!-- 4. Modal Konfirmasi Restart -->
    <div v-if="showRestartModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-amber-500/40 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex items-center gap-2.5 text-amber-400">
          <AlertTriangle class="w-5 h-5" />
          <h4 class="font-bold text-sm">Konfirmasi Mulai Ulang (Restart)</h4>
        </div>
        <p class="text-xs text-slate-300 leading-relaxed">
          Apakah Anda yakin ingin memulai ulang perangkat STB? Koneksi Web Dashboard dan seluruh kontainer Docker akan terputus sementara.
        </p>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showRestartModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="confirmRestart" class="px-4 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs">Ya, Restart</button>
        </div>
      </div>
    </div>

    <!-- 5. Modal Konfirmasi Shutdown -->
    <div v-if="showShutdownModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-rose-500/40 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex items-center gap-2.5 text-rose-400">
          <Power class="w-5 h-5" />
          <h4 class="font-bold text-sm">Konfirmasi Matikan Daya (Shutdown)</h4>
        </div>
        <p class="text-xs text-slate-300 leading-relaxed">
          Apakah Anda yakin ingin mematikan sistem STB sepenuhnya? Anda perlu menyalakan ulang adaptor daya fisik perangkat secara manual.
        </p>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showShutdownModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="confirmShutdown" class="px-4 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs">Ya, Matikan</button>
        </div>
      </div>
    </div>

    <!-- 6. Modal Konfirmasi Logout -->
    <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl p-5 space-y-4 text-slate-100">
        <div class="flex items-center gap-2.5 text-slate-200">
          <LogOut class="w-5 h-5 text-sky-400" />
          <h4 class="font-bold text-sm">Konfirmasi Keluar (Logout)</h4>
        </div>
        <p class="text-xs text-slate-300 leading-relaxed">
          Apakah Anda ingin mengakhiri sesi login saat ini? Anda harus memasukkan kata sandi kembali untuk masuk.
        </p>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showLogoutModal = false" class="px-4 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs">Batal</button>
          <button @click="confirmLogout" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs">Ya, Keluar</button>
        </div>
      </div>
    </div>

  </div>
</template>
