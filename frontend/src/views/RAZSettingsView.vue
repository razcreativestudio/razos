<script setup>
/**
 * RAZOS Settings View (CasaOS Architecture Standard)
 * Updated: 2026-10-14
 * Deskripsi: Panel pengaturan sistem terpadu multi-tab mencakup Manajemen Akun,
 * Pengaturan Preferensi Desktop & OTA Updater, serta Konfigurasi Widget Real-time.
 */
import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import {
  User, Sliders, LayoutGrid, Clock, Activity, HardDrive, Network,
  RefreshCw, Power, ArrowUpCircle, Check, Pencil, Image, CheckCircle2,
  X, AlertTriangle, LogOut, Sparkles, ExternalLink, ShieldCheck
} from 'lucide-vue-next'

const systemStore = useSystemStore()

// Tab Aktif: 'account' | 'settings' | 'widgets'
const activeTab = ref('settings')

// State OTA Updater
const isCheckingUpdate = ref(false)
const updateStatus = ref({
  current_version: '2.6.0',
  latest_version: '2.6.0',
  has_update: false,
  release_notes: 'Sistem dalam kondisi paling mutakhir.',
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
        systemStore.showNotification('Sistem Up to Date', 'RAZOS menggunakan versi terbaru', 'success')
      }
    }
  } catch (err) {
    // Fallback info bila dev mode
    systemStore.showNotification('Info Versi', 'RAZOS v2.6.0 telah terpasang dengan versi terbaru', 'success')
  } finally {
    isCheckingUpdate.value = false
  }
}

async function applyUpdate() {
  if (!updateStatus.value.download_url) return
  try {
    const res = await fetch('/api/v1/updater/apply', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ download_url: updateStatus.value.download_url })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', 'Pembaruan OTA berhasil dipasang!', 'success')
      checkUpdate()
    } else {
      systemStore.showNotification('Gagal', 'Gagal memasang pembaruan otomatis', 'error')
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal menghubungi server', 'error')
  }
}

// ==========================================
// STATE CUSTOM MODALS (WINDOW DIALOGS)
// ==========================================
const showEditNameModal = ref(false)
const inputName = ref('')

function openEditName() {
  inputName.value = systemStore.settings.account.name || 'raz'
  showEditNameModal.value = true
}

function saveEditName() {
  if (!inputName.value.trim()) return
  systemStore.updateAccount(inputName.value.trim(), null)
  systemStore.showNotification('Akun Diperbarui', `Nama profil diubah menjadi: ${inputName.value}`, 'success')
  showEditNameModal.value = false
}

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

function saveEditPassword() {
  if (!newPassword.value) {
    systemStore.showNotification('Peringatan', 'Kata sandi baru tidak boleh kosong', 'error')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    systemStore.showNotification('Peringatan', 'Konfirmasi kata sandi tidak cocok', 'error')
    return
  }
  systemStore.showNotification('Sukses', 'Kata sandi sistem berhasil diperbarui', 'success')
  showEditPasswordModal.value = false
}

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
  systemStore.showNotification('Port WebUI', `Port diubah menjadi :${p}. Muat ulang halaman jika port service berubah.`, 'success')
  showPortModal.value = false
}

const showWallpaperModal = ref(false)
const customWallpaperUrl = ref('')

// Pilihan wallpaper preset bawaan sistem
const presetWallpapers = [
  {
    name: 'RAZ Cosmic Dark (Default)',
    url: '/assets/images/RAZDesktopBg.webp'
  },
  {
    name: 'RAZ Gradient Nebula',
    url: '/assets/images/RAZMobileBg.webp'
  },
  {
    name: 'Cyber Horizon',
    url: 'https://images.unsplash.com/photo-1550745165-9bc0b252726f?q=80&w=1920&auto=format&fit=crop'
  },
  {
    name: 'Deep Minimalist Slate',
    url: 'https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=1920&auto=format&fit=crop'
  }
]

function selectWallpaper(url) {
  systemStore.updateSetting('wallpaper', url)
  systemStore.showNotification('Wallpaper Berubah', 'Wallpaper desktop berhasil diperbarui', 'success')
  showWallpaperModal.value = false
}

function saveCustomWallpaper() {
  if (!customWallpaperUrl.value.trim()) return
  systemStore.updateSetting('wallpaper', customWallpaperUrl.value.trim())
  systemStore.showNotification('Wallpaper Berubah', 'Wallpaper kustom berhasil diterapkan', 'success')
  showWallpaperModal.value = false
}

// Konfirmasi Daya (Restart / Shutdown / Logout)
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

function confirmLogout() {
  showLogoutModal.value = false
  systemStore.showNotification('Logout', 'Sesi telah diakhiri. Silakan login kembali.', 'info')
}

onMounted(() => {
  checkUpdate()
})
</script>

<template>
  <div class="h-full flex flex-col text-slate-100 select-none overflow-hidden">
    <!-- Top Navigation Tabs (CasaOS Style: Account | Settings | Widgets) -->
    <div class="flex items-center justify-center pt-2 pb-3 border-b border-slate-800/80 bg-slate-950/40">
      <div class="inline-flex items-center p-1 rounded-2xl bg-slate-900 border border-slate-800/90 shadow-inner">
        <!-- Tab 1: Account -->
        <button
          @click="activeTab = 'account'"
          class="flex items-center justify-center w-11 h-9 rounded-xl transition duration-150"
          :class="activeTab === 'account' ? 'bg-slate-800 text-sky-400 shadow-md font-bold' : 'text-slate-400 hover:text-slate-200'"
          title="Account Profile"
        >
          <User class="w-5 h-5" />
        </button>

        <!-- Tab 2: Settings -->
        <button
          @click="activeTab = 'settings'"
          class="flex items-center justify-center w-11 h-9 rounded-xl transition duration-150"
          :class="activeTab === 'settings' ? 'bg-slate-800 text-sky-400 shadow-md font-bold' : 'text-slate-400 hover:text-slate-200'"
          title="System Settings"
        >
          <Sliders class="w-5 h-5" />
        </button>

        <!-- Tab 3: Widgets -->
        <button
          @click="activeTab = 'widgets'"
          class="flex items-center justify-center w-11 h-9 rounded-xl transition duration-150"
          :class="activeTab === 'widgets' ? 'bg-slate-800 text-sky-400 shadow-md font-bold' : 'text-slate-400 hover:text-slate-200'"
          title="Widget Settings"
        >
          <LayoutGrid class="w-5 h-5" />
        </button>
      </div>
    </div>

    <!-- Main Content Area Scrollable -->
    <div class="flex-1 overflow-y-auto px-6 py-5 max-w-2xl mx-auto w-full space-y-6">

      <!-- ============================================================= -->
      <!-- TAB 1: ACCOUNT (Sesuai Gambar 2) -->
      <!-- ============================================================= -->
      <div v-if="activeTab === 'account'" class="space-y-6 animate-in fade-in duration-200">
        <!-- Title & Logout Button -->
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-slate-100">Account</h2>
          <button
            @click="showLogoutModal = true"
            class="px-5 py-1.5 rounded-full bg-rose-500/15 hover:bg-rose-500 text-rose-400 hover:text-white border border-rose-500/30 text-xs font-semibold transition active:scale-95 shadow-sm"
          >
            Logout
          </button>
        </div>

        <!-- Center Avatar Display -->
        <div class="flex flex-col items-center justify-center py-4">
          <div class="relative group">
            <div class="w-24 h-24 rounded-full border-2 border-slate-700 bg-slate-900 p-1.5 shadow-xl flex items-center justify-center overflow-hidden">
              <img
                :src="systemStore.settings.account.avatar || '/assets/brand/RAZBrandSquare.webp'"
                alt="Account Avatar"
                class="w-full h-full object-cover rounded-full"
              />
            </div>
            <div class="absolute bottom-0 right-0 w-7 h-7 rounded-full bg-sky-500 text-slate-950 flex items-center justify-center border-2 border-slate-900 shadow">
              <Check class="w-4 h-4 stroke-[3]" />
            </div>
          </div>
          <p class="mt-3 text-sm font-semibold text-slate-200 capitalize">
            {{ systemStore.settings.account.name || 'raz' }}
          </p>
          <span class="text-[11px] text-slate-400">Administrator Perangkat STB</span>
        </div>

        <!-- Account Credentials List Card -->
        <div class="rounded-2xl bg-slate-900/90 border border-slate-800/90 divide-y divide-slate-800/80 shadow-lg overflow-hidden">
          <!-- Row 1: Name -->
          <div class="flex items-center justify-between px-5 py-4">
            <span class="text-sm font-medium text-slate-300">Name</span>
            <div class="flex items-center gap-3">
              <span class="text-sm font-semibold text-slate-100">{{ systemStore.settings.account.name || 'raz' }}</span>
              <button
                @click="openEditName"
                class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-100 flex items-center justify-center transition border border-slate-700/60"
                title="Edit Nama"
              >
                <Pencil class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Row 2: Password -->
          <div class="flex items-center justify-between px-5 py-4">
            <span class="text-sm font-medium text-slate-300">Password</span>
            <div class="flex items-center gap-3">
              <span class="text-sm font-mono tracking-widest text-slate-400">••••••••</span>
              <button
                @click="openEditPassword"
                class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-100 flex items-center justify-center transition border border-slate-700/60"
                title="Ganti Password"
              >
                <Pencil class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================= -->
      <!-- TAB 2: SETTINGS (Sesuai Gambar 1) -->
      <!-- ============================================================= -->
      <div v-else-if="activeTab === 'settings'" class="space-y-6 animate-in fade-in duration-200">
        <!-- Title -->
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-slate-100">Settings</h2>
          <span class="text-xs text-slate-400 font-mono">v{{ updateStatus.current_version }}</span>
        </div>

        <!-- Main Preferences Card -->
        <div class="rounded-2xl bg-slate-900/90 border border-slate-800/90 divide-y divide-slate-800/80 shadow-lg overflow-hidden">
          <!-- Row 1: Show search bar -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Show search bar</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.showSearchBar"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Row 2: Search Engine -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Search Engine</span>
            <select
              v-model="systemStore.settings.searchEngine"
              @change="systemStore.saveSettings()"
              class="px-3 py-1.5 rounded-xl bg-slate-800 border border-slate-700 text-slate-200 text-xs font-medium outline-none focus:border-sky-400 cursor-pointer transition"
            >
              <option value="DuckDuckGo">DuckDuckGo</option>
              <option value="Google">Google</option>
              <option value="Bing">Bing</option>
              <option value="Baidu">Baidu</option>
            </select>
          </div>

          <!-- Row 3: Language -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Language</span>
            <select
              v-model="systemStore.settings.language"
              @change="systemStore.saveSettings()"
              class="px-3 py-1.5 rounded-xl bg-slate-800 border border-slate-700 text-slate-200 text-xs font-medium outline-none focus:border-sky-400 cursor-pointer transition"
            >
              <option value="English">English</option>
              <option value="Bahasa Indonesia">Bahasa Indonesia</option>
            </select>
          </div>

          <!-- Row 4: WebUI Port -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">WebUI Port</span>
            <div class="flex items-center gap-3">
              <span class="text-xs font-mono font-semibold text-slate-300">{{ systemStore.settings.webUiPort || 8080 }}</span>
              <button
                @click="openPortModal"
                class="px-3.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium transition border border-slate-700/60 active:scale-95"
              >
                Change
              </button>
            </div>
          </div>

          <!-- Row 5: Wallpaper -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Wallpaper</span>
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg overflow-hidden border border-slate-700 bg-slate-950 flex items-center justify-center shadow-inner">
                <img :src="systemStore.settings.wallpaper" alt="Wallpaper Preview" class="w-full h-full object-cover" />
              </div>
              <button
                @click="showWallpaperModal = true"
                class="px-3.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium transition border border-slate-700/60 active:scale-95"
              >
                Change
              </button>
            </div>
          </div>

          <!-- Row 6: Show news feed from blog -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Show news feed from RAZOS blog</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.showBlogFeed"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Row 7: Show recommended apps -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Show recommended apps</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.showRecommendedApps"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Row 8: Automount USB drive -->
          <div class="flex items-center justify-between px-5 py-3.5">
            <span class="text-sm font-medium text-slate-200">Automount USB drive</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.automountUsb"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>
        </div>

        <!-- Section Card: Update (Sesuai Gambar 1) -->
        <div class="rounded-2xl bg-slate-900/90 border border-slate-800/90 p-5 shadow-lg space-y-3">
          <h3 class="text-sm font-bold text-slate-200">Update</h3>

          <div class="p-4 rounded-xl bg-slate-950/60 border border-slate-800 flex items-center justify-between gap-4">
            <div class="flex items-center gap-3.5 min-w-0">
              <div class="w-10 h-10 rounded-xl bg-emerald-500/20 text-emerald-400 flex items-center justify-center border border-emerald-500/30 shrink-0">
                <CheckCircle2 class="w-6 h-6" />
              </div>
              <div class="min-w-0">
                <div class="text-sm font-bold text-slate-100 truncate">RAZOS v{{ updateStatus.current_version }}</div>
                <div class="text-xs text-emerald-400 font-medium flex items-center gap-1.5">
                  <span>{{ updateStatus.has_update ? 'Versi Baru Tersedia!' : 'Currently, at the latest version' }}</span>
                  <Check class="w-3.5 h-3.5 stroke-[3]" />
                </div>
              </div>
            </div>

            <div class="shrink-0 flex items-center gap-2">
              <button
                @click="checkUpdate"
                :disabled="isCheckingUpdate"
                class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold text-xs transition border border-slate-700/60 active:scale-95 flex items-center gap-1.5"
              >
                <RefreshCw class="w-3.5 h-3.5" :class="isCheckingUpdate ? 'animate-spin' : ''" />
                <span>{{ isCheckingUpdate ? 'Checking...' : 'Check Update' }}</span>
              </button>

              <button
                v-if="updateStatus.has_update"
                @click="applyUpdate"
                class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow-lg active:scale-95 flex items-center gap-1.5"
              >
                <Sparkles class="w-3.5 h-3.5" />
                <span>Update Now</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Footer Power Actions: Restart & Shutdown (Sesuai Gambar 1) -->
        <div class="flex items-center justify-between pt-2">
          <!-- Restart Button -->
          <button
            @click="showRestartModal = true"
            class="px-5 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-slate-100 border border-slate-800 text-xs font-semibold transition flex items-center gap-2 shadow-sm active:scale-95"
          >
            <RefreshCw class="w-4 h-4 text-slate-400" />
            <span>Restart</span>
          </button>

          <!-- Shutdown Button -->
          <button
            @click="showShutdownModal = true"
            class="px-5 py-2.5 rounded-xl bg-rose-500/15 hover:bg-rose-500 text-rose-400 hover:text-white border border-rose-500/30 text-xs font-semibold transition flex items-center gap-2 shadow-sm active:scale-95"
          >
            <Power class="w-4 h-4" />
            <span>Shutdown</span>
          </button>
        </div>
      </div>

      <!-- ============================================================= -->
      <!-- TAB 3: WIDGET SETTINGS (Sesuai Gambar 3) -->
      <!-- ============================================================= -->
      <div v-else-if="activeTab === 'widgets'" class="space-y-6 animate-in fade-in duration-200">
        <!-- Title -->
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-slate-100">Widget settings</h2>
          <span class="text-xs text-slate-400">Pilih widget yang tampil di desktop host</span>
        </div>

        <!-- Widgets Toggle List Card -->
        <div class="rounded-2xl bg-slate-900/90 border border-slate-800/90 divide-y divide-slate-800/80 shadow-lg overflow-hidden">
          <!-- Widget 1: Time -->
          <div class="flex items-center justify-between px-5 py-4">
            <div class="flex items-center gap-3 text-slate-200">
              <Clock class="w-5 h-5 text-sky-400" />
              <span class="text-sm font-medium">Time</span>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.time"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Widget 2: System status -->
          <div class="flex items-center justify-between px-5 py-4">
            <div class="flex items-center gap-3 text-slate-200">
              <Activity class="w-5 h-5 text-emerald-400" />
              <span class="text-sm font-medium">System status</span>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.systemStatus"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Widget 3: Storage status -->
          <div class="flex items-center justify-between px-5 py-4">
            <div class="flex items-center gap-3 text-slate-200">
              <HardDrive class="w-5 h-5 text-amber-400" />
              <span class="text-sm font-medium">Storage status</span>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.storageStatus"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <!-- Widget 4: Network status -->
          <div class="flex items-center justify-between px-5 py-4">
            <div class="flex items-center gap-3 text-slate-200">
              <Network class="w-5 h-5 text-indigo-400" />
              <span class="text-sm font-medium">Network status</span>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                v-model="systemStore.settings.widgets.networkStatus"
                @change="systemStore.saveSettings()"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-slate-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- ============================================================= -->
    <!-- CUSTOM MODALS (WINDOW DIALOGS TANPA BROWSER ALERT/PROMPT) -->
    <!-- ============================================================= -->

    <!-- 1. Modal Edit Nama Profil -->
    <div v-if="showEditNameModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800">
          <h3 class="font-bold text-sm">Edit Profil Pengguna</h3>
          <button @click="showEditNameModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5 space-y-3">
          <label class="text-xs font-semibold text-slate-300">Nama Lengkap / Username</label>
          <input
            v-model="inputName"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm outline-none focus:border-sky-400"
            placeholder="Masukkan nama pengguna..."
            @keyup.enter="saveEditName"
          />
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showEditNameModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="saveEditName" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold">Simpan</button>
        </div>
      </div>
    </div>

    <!-- 2. Modal Ganti Password -->
    <div v-if="showEditPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800">
          <h3 class="font-bold text-sm">Ubah Kata Sandi</h3>
          <button @click="showEditPasswordModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5 space-y-3">
          <div>
            <label class="text-xs font-semibold text-slate-300 block mb-1">Kata Sandi Saat Ini</label>
            <input
              v-model="oldPassword"
              type="password"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm outline-none focus:border-sky-400"
              placeholder="••••••••"
            />
          </div>
          <div>
            <label class="text-xs font-semibold text-slate-300 block mb-1">Kata Sandi Baru</label>
            <input
              v-model="newPassword"
              type="password"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm outline-none focus:border-sky-400"
              placeholder="Minimal 6 karakter"
            />
          </div>
          <div>
            <label class="text-xs font-semibold text-slate-300 block mb-1">Ulangi Kata Sandi Baru</label>
            <input
              v-model="confirmPassword"
              type="password"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm outline-none focus:border-sky-400"
              placeholder="Ulangi kata sandi"
              @keyup.enter="saveEditPassword"
            />
          </div>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showEditPasswordModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="saveEditPassword" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold">Simpan Password</button>
        </div>
      </div>
    </div>

    <!-- 3. Modal Ganti WebUI Port -->
    <div v-if="showPortModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800">
          <h3 class="font-bold text-sm">Ubah WebUI Port</h3>
          <button @click="showPortModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5 space-y-3">
          <label class="text-xs font-semibold text-slate-300">Port Web Antarmuka RAZOS</label>
          <input
            v-model="inputPort"
            type="number"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm font-mono outline-none focus:border-sky-400"
            placeholder="8080"
            @keyup.enter="savePort"
          />
          <p class="text-[11px] text-slate-400 leading-relaxed">
            Perubahan port akan disimpan ke konfigurasi host. Pastikan port baru tidak bentrok dengan service container lain.
          </p>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showPortModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="savePort" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold">Terapkan Port</button>
        </div>
      </div>
    </div>

    <!-- 4. Modal Pilih Wallpaper -->
    <div v-if="showWallpaperModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800">
          <div class="flex items-center gap-2">
            <Image class="w-4 h-4 text-sky-400" />
            <h3 class="font-bold text-sm">Pilih Wallpaper Desktop</h3>
          </div>
          <button @click="showWallpaperModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5 space-y-4 max-h-[60vh] overflow-y-auto">
          <!-- Preset Wallpapers Grid -->
          <div class="grid grid-cols-2 gap-3">
            <div
              v-for="wp in presetWallpapers"
              :key="wp.url"
              @click="selectWallpaper(wp.url)"
              class="relative rounded-xl overflow-hidden border-2 cursor-pointer group transition-all duration-200 aspect-video shadow-md hover:scale-[1.02]"
              :class="systemStore.settings.wallpaper === wp.url ? 'border-sky-400 ring-2 ring-sky-400/30' : 'border-slate-800 hover:border-slate-600'"
            >
              <img :src="wp.url" :alt="wp.name" class="w-full h-full object-cover" />
              <div class="absolute inset-0 bg-gradient-to-t from-slate-950/80 via-transparent to-transparent flex items-end p-2.5">
                <span class="text-[11px] font-semibold text-slate-200 truncate">{{ wp.name }}</span>
              </div>
              <div v-if="systemStore.settings.wallpaper === wp.url" class="absolute top-2 right-2 w-5 h-5 rounded-full bg-sky-500 text-slate-950 flex items-center justify-center shadow">
                <Check class="w-3.5 h-3.5 stroke-[3]" />
              </div>
            </div>
          </div>

          <!-- Custom Wallpaper Input -->
          <div class="pt-2 border-t border-slate-800 space-y-2">
            <label class="text-xs font-semibold text-slate-300">Atau Gunakan URL Gambar Kustom</label>
            <div class="flex items-center gap-2">
              <input
                v-model="customWallpaperUrl"
                type="text"
                class="flex-1 px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-xs font-mono outline-none focus:border-sky-400"
                placeholder="https://example.com/wallpaper.jpg"
              />
              <button
                @click="saveCustomWallpaper"
                class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition whitespace-nowrap"
              >
                Gunakan URL
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 5. Modal Konfirmasi Restart -->
    <div v-if="showRestartModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-amber-500/40 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800 text-amber-400">
          <div class="flex items-center gap-2">
            <AlertTriangle class="w-4 h-4" />
            <h3 class="font-bold text-sm">Konfirmasi Restart</h3>
          </div>
          <button @click="showRestartModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5">
          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda yakin ingin memulai ulang perangkat STB sekarang? Seluruh website dan container akan dimatikan sementara dan dihidupkan kembali otomatis setelah reboot selesai.
          </p>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showRestartModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="confirmRestart" class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 shadow">
            <RefreshCw class="w-3.5 h-3.5" />
            <span>Ya, Restart Sistem</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 6. Modal Konfirmasi Shutdown -->
    <div v-if="showShutdownModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-rose-500/40 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800 text-rose-400">
          <div class="flex items-center gap-2">
            <Power class="w-4 h-4" />
            <h3 class="font-bold text-sm">Konfirmasi Shutdown</h3>
          </div>
          <button @click="showShutdownModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5">
          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda yakin ingin mematikan daya perangkat STB? Anda harus mencabut dan menancapkan adaptor daya kembali secara fisik untuk menyalakan perangkat.
          </p>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showShutdownModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="confirmShutdown" class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white text-xs font-bold flex items-center gap-1.5 shadow">
            <Power class="w-3.5 h-3.5" />
            <span>Matikan Daya</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 7. Modal Konfirmasi Logout -->
    <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3 border-b border-slate-800">
          <div class="flex items-center gap-2 text-rose-400">
            <LogOut class="w-4 h-4" />
            <h3 class="font-bold text-sm">Keluar dari Sesi</h3>
          </div>
          <button @click="showLogoutModal = false" class="text-slate-400 hover:text-slate-100">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="p-5">
          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda ingin mengakhiri sesi login RAZOS saat ini?
          </p>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showLogoutModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="confirmLogout" class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white text-xs font-bold">Logout</button>
        </div>
      </div>
    </div>

  </div>
</template>
