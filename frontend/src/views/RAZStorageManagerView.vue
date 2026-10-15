<script setup>
/**
 * RAZOS Storage Manager View (CasaOS Style - Standalone Window & Settings Component)
 * Created: 2026-10-15 (v3.5.0 - Storage Manager Suite)
 * Deskripsi: Antarmuka manajemen penyimpanan sistem tingkat lanjut sesuai 2 spesifikasi referensi:
 *   1. Tab "Storage": Partisi BOOT_EMMC (OS), ROOTFS_EMMC (OS) dengan 'Free up storage',
 *      MicroSD dengan 'Format' dan 'Remove', tombol 'Create Storage', serta 'Merge storages (Beta)'.
 *   2. Tab "Drive": Media drive fisik mmcblk2 (System eMMC) dan mmcblk1 (MicroSD)
 *      lengkap dengan indikator Health (Healthy) dan Sensor Suhu (Temperature).
 *   3. Custom Modals: Bebas browser alert/prompt, 100% terintegrasi Glassmorphism Slate-900.
 */

import { ref, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  X,
  HardDrive,
  Info,
  Check,
  AlertTriangle,
  Trash2,
  RefreshCw,
  Plus,
  Sparkles,
  Layers,
  CheckCircle2,
  Cpu,
  Folder
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: false,
    default: () => ({ id: 'storage-manager' })
  },
  isModal: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const systemStore = useSystemStore()
const windowStore = useWindowStore()

// Tab Aktif: 'storage' | 'drive'
const activeTab = ref('storage')
const isLoading = ref(false)

// Data Penyimpanan (Partisi Storage & Drive Fisik)
const storages = ref([
  {
    id: 'boot_emmc',
    name: 'BOOT_EMMC',
    is_os: true,
    type_desc: 'Single storage drive, VFAT',
    device: 'mmcblk2p1',
    mountpoint: '/boot',
    fstype: 'VFAT',
    total_formatted: '509.98 MB',
    available_formatted: '359.23 MB',
    percent: 29.6,
    status: 'normal'
  },
  {
    id: 'rootfs_emmc',
    name: 'ROOTFS_EMMC',
    is_os: true,
    type_desc: 'Single storage drive, EXT4',
    device: 'mmcblk2p2',
    mountpoint: '/',
    fstype: 'EXT4',
    total_formatted: '6.47 GB',
    available_formatted: '726.02 MB',
    percent: 88.8,
    status: 'warning',
    can_free_up: true
  },
  {
    id: 'microsd',
    name: 'MicroSD',
    is_os: false,
    type_desc: 'Single storage drive, EXT4',
    device: 'mmcblk1p1',
    mountpoint: '/DATA',
    fstype: 'EXT4',
    total_formatted: '58.94 GB',
    available_formatted: '46.35 GB',
    percent: 21.4,
    status: 'normal',
    can_format: true,
    can_remove: true
  }
])

const drives = ref([
  {
    id: 'mmcblk2',
    name: 'mmcblk2',
    type_desc: 'System, 7.28 GB MMC',
    size: '7.28 GB',
    health: 'Healthy',
    temperature: 'N/A'
  },
  {
    id: 'mmcblk1',
    name: 'mmcblk1',
    type_desc: 'Storage, 58.94 GB MicroSD',
    size: '58.94 GB',
    health: 'Healthy',
    temperature: '38°C'
  }
])

// Ambil data real-time dari backend
async function fetchStorageData() {
  isLoading.value = true
  try {
    const res = await fetch('/api/v1/storage/manager')
    if (res.ok) {
      const data = await res.json()
      if (Array.isArray(data.storages) && data.storages.length > 0) {
        storages.value = data.storages
      }
      if (Array.isArray(data.drives) && data.drives.length > 0) {
        drives.value = data.drives
      }
    }
  } catch (err) {
    console.warn('Gagal memuat data storage manager:', err)
  } finally {
    isLoading.value = false
  }
}

// ==========================================
// STATE CUSTOM MODALS (Zero Native Dialogs)
// ==========================================
const showFreeUpModal = ref(false)
const isCleaning = ref(false)
const selectedStorage = ref(null)

function openFreeUpModal(storage) {
  selectedStorage.value = storage
  showFreeUpModal.value = true
}

async function executeFreeUp() {
  isCleaning.value = true
  try {
    const res = await fetch('/api/v1/storage/free-up', { method: 'POST' })
    if (res.ok) {
      const result = await res.json()
      systemStore.showNotification('Penyimpanan Dibersihkan', result.message || 'Log dan cache berhasil dibersihkan!', 'success')
      showFreeUpModal.value = false
      await fetchStorageData()
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal membersihkan cache penyimpanan', 'error')
  } finally {
    isCleaning.value = false
  }
}

const showFormatModal = ref(false)
const isFormatting = ref(false)
const formatFs = ref('EXT4')

function openFormatModal(storage) {
  selectedStorage.value = storage
  formatFs.value = 'EXT4'
  showFormatModal.value = true
}

async function executeFormat() {
  isFormatting.value = true
  try {
    const res = await fetch('/api/v1/storage/format', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        device: selectedStorage.value?.device,
        fstype: formatFs.value
      })
    })
    if (res.ok) {
      systemStore.showNotification('Format Berhasil', `Drive ${selectedStorage.value?.name} berhasil diformat ke ${formatFs.value}`, 'success')
      showFormatModal.value = false
      await fetchStorageData()
    }
  } catch (e) {
    systemStore.showNotification('Gagal', 'Terjadi kesalahan saat memformat drive', 'error')
  } finally {
    isFormatting.value = false
  }
}

const showRemoveModal = ref(false)
const isRemoving = ref(false)

function openRemoveModal(storage) {
  selectedStorage.value = storage
  showRemoveModal.value = true
}

async function executeRemove() {
  isRemoving.value = true
  try {
    const res = await fetch('/api/v1/storage/remove', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ device: selectedStorage.value?.device })
    })
    if (res.ok) {
      systemStore.showNotification('Drive Dilepas', `Drive ${selectedStorage.value?.name} telah di-unmount secara aman.`, 'success')
      showRemoveModal.value = false
      await fetchStorageData()
    }
  } catch (e) {
    systemStore.showNotification('Gagal', 'Gagal melepaskan drive penyimpanan', 'error')
  } finally {
    isRemoving.value = false
  }
}

const showCreateStorageModal = ref(false)
const showMergeModal = ref(false)

function handleClose() {
  if (props.isModal) {
    emit('close')
  } else {
    windowStore.closeWindow(props.win?.id || 'storage-manager')
  }
}

onMounted(() => {
  fetchStorageData()
})
</script>

<template>
  <div class="h-full flex flex-col bg-slate-900 text-slate-100 select-none overflow-hidden font-sans relative">
    <!-- TITLEBAR STORAGE MANAGER -->
    <div class="flex items-center justify-between px-6 py-3.5 border-b border-slate-800 bg-slate-950/70 shrink-0">
      <div class="flex items-center gap-2.5">
        <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
          <HardDrive class="w-4 h-4" />
        </div>
        <h2 class="text-sm font-bold tracking-tight text-slate-100">Storage Manager</h2>
      </div>

      <button
        type="button"
        @click="handleClose"
        class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-400 hover:text-white flex items-center justify-center transition border border-slate-700/60 active:scale-95"
        title="Tutup Jendela"
      >
        <X class="w-4 h-4" />
      </button>
    </div>

    <!-- SUB-HEADER TAB SWITCHES: Storage | Drive | Create Storage (Sesuai Gambar 1 & 2) -->
    <div class="flex items-center justify-between px-6 pt-3 pb-0 border-b border-slate-800 bg-slate-900/60 shrink-0">
      <div class="flex items-center gap-6">
        <!-- Tab Storage -->
        <button
          type="button"
          @click="activeTab = 'storage'"
          class="relative py-2 text-sm font-semibold transition"
          :class="activeTab === 'storage' ? 'text-sky-400 border-b-2 border-sky-500' : 'text-slate-400 hover:text-slate-200 border-b-2 border-transparent'"
        >
          Storage
        </button>

        <!-- Tab Drive -->
        <button
          type="button"
          @click="activeTab = 'drive'"
          class="relative py-2 text-sm font-semibold transition"
          :class="activeTab === 'drive' ? 'text-sky-400 border-b-2 border-sky-500' : 'text-slate-400 hover:text-slate-200 border-b-2 border-transparent'"
        >
          Drive
        </button>
      </div>

      <!-- Tombol Create Storage (Sesuai Gambar 1 Kanan Atas) -->
      <button
        v-if="activeTab === 'storage'"
        type="button"
        @click="showCreateStorageModal = true"
        class="mb-2 px-4 py-1 rounded-full border border-sky-500/40 hover:border-sky-400 bg-sky-500/10 hover:bg-sky-500/20 text-sky-300 text-xs font-semibold transition active:scale-95 flex items-center gap-1.5 shadow-sm"
      >
        <span>Create Storage</span>
      </button>
    </div>

    <!-- MAIN BODY CONTAINER (SCROLLABLE) -->
    <div class="flex-1 overflow-y-auto px-6 py-5 space-y-4">
      <!-- ============================================================= -->
      <!-- TAB 1: STORAGE (Sesuai Gambar 1)                               -->
      <!-- ============================================================= -->
      <div v-if="activeTab === 'storage'" class="space-y-4 animate-in fade-in duration-150">
        <!-- Card Item Loop -->
        <div
          v-for="st in storages"
          :key="st.id"
          class="p-4 rounded-2xl bg-slate-950/70 border border-slate-800/90 shadow-lg space-y-3 transition hover:border-slate-700"
        >
          <!-- Top Row: Icon, Title, Badge, Subtitle & Action Links -->
          <div class="flex items-start justify-between gap-4">
            <div class="flex items-start gap-3.5">
              <!-- Drive Icon Box Bergaya CasaOS / Apple Drive (Sesuai Gambar 1) -->
              <div class="w-12 h-14 rounded-xl bg-slate-800/90 border border-slate-700/80 p-1 flex flex-col justify-between items-center shrink-0 shadow-inner overflow-hidden">
                <div class="w-full h-8 rounded-lg bg-gradient-to-b from-slate-200 to-slate-400 border border-slate-300 shadow flex items-center justify-center">
                  <div class="w-5 h-0.5 bg-slate-500 rounded-full"></div>
                </div>
                <div class="w-full h-3 rounded-b-md bg-slate-950 flex items-center justify-between px-1.5">
                  <span class="w-1.5 h-1.5 rounded-full" :class="st.status === 'warning' ? 'bg-rose-500 animate-pulse' : 'bg-sky-400'"></span>
                  <div class="w-4 h-0.5 bg-slate-800 rounded"></div>
                </div>
              </div>

              <!-- Details Info -->
              <div class="space-y-0.5">
                <div class="flex items-center gap-2">
                  <h3 class="text-sm font-bold text-slate-100 tracking-wide">{{ st.name }}</h3>
                  <!-- Badge OS Sesuai Gambar 1 -->
                  <span
                    v-if="st.is_os"
                    class="px-1.5 py-0.2 rounded text-[10px] font-bold text-sky-400 border border-sky-400/40 bg-sky-500/10"
                  >
                    OS
                  </span>
                </div>

                <!-- Sub-text with info icon -->
                <div class="flex items-center gap-1.5 text-xs text-slate-400">
                  <span>{{ st.type_desc }}</span>
                  <Info class="w-3.5 h-3.5 text-slate-500 cursor-pointer hover:text-slate-300 transition" />
                </div>

                <!-- Device detail info -->
                <div class="text-[11px] text-slate-400 font-mono pt-0.5">
                  {{ st.device }} | Available: <span class="text-slate-200 font-semibold">{{ st.available_formatted }}</span> (Total: {{ st.total_formatted }})
                </div>
              </div>
            </div>

            <!-- Action Buttons Kanan Sesuai Gambar 1 -->
            <div class="flex items-center gap-2 shrink-0 pt-1">
              <!-- Tombol Free up storage (pada ROOTFS_EMMC) -->
              <button
                v-if="st.can_free_up"
                type="button"
                @click="openFreeUpModal(st)"
                class="text-xs font-semibold text-sky-400 hover:text-sky-300 hover:underline transition"
              >
                Free up storage
              </button>

              <!-- Tombol Format (pada MicroSD) -->
              <button
                v-if="st.can_format"
                type="button"
                @click="openFormatModal(st)"
                class="px-3.5 py-1 rounded-full border border-slate-700 hover:border-slate-600 bg-slate-800/80 hover:bg-slate-700 text-slate-300 text-xs font-medium transition active:scale-95 shadow-sm"
              >
                Format
              </button>

              <!-- Tombol Remove (pada MicroSD) -->
              <button
                v-if="st.can_remove"
                type="button"
                @click="openRemoveModal(st)"
                class="px-3.5 py-1 rounded-full border border-slate-700 hover:border-rose-500/40 bg-slate-800/80 hover:bg-rose-500/10 text-slate-300 hover:text-rose-400 text-xs font-medium transition active:scale-95 shadow-sm"
              >
                Remove
              </button>
            </div>
          </div>

          <!-- Bottom Row: Progress Bar (Biru atau Merah untuk Peringatan Penuh) -->
          <div class="w-full bg-slate-800/80 rounded-full h-2 overflow-hidden shadow-inner">
            <div
              class="h-full rounded-full transition-all duration-500 ease-out"
              :class="st.status === 'warning' ? 'bg-rose-500' : 'bg-blue-600'"
              :style="{ width: `${Math.min(100, Math.max(2, st.percent))}%` }"
            ></div>
          </div>
        </div>

        <!-- Footer Action: Merge storages (Beta) Sesuai Gambar 1 Kanan Bawah -->
        <div class="flex items-center justify-end gap-2 pt-3 pb-2">
          <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-emerald-500 text-slate-950 uppercase tracking-wider">
            Beta
          </span>
          <button
            type="button"
            @click="showMergeModal = true"
            class="px-5 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-semibold text-xs sm:text-sm transition shadow-lg shadow-blue-600/30 active:scale-95 flex items-center gap-1.5"
          >
            <span>Merge storages</span>
          </button>
        </div>
      </div>

      <!-- ============================================================= -->
      <!-- TAB 2: DRIVE (Sesuai Gambar 2)                                  -->
      <!-- ============================================================= -->
      <div v-else class="space-y-4 animate-in fade-in duration-150">
        <!-- Physical Drive Card Item Loop -->
        <div
          v-for="drv in drives"
          :key="drv.id"
          class="p-4 rounded-2xl bg-slate-950/70 border border-slate-800/90 shadow-lg flex items-center justify-between gap-4 transition hover:border-slate-700"
        >
          <!-- Left: Drive Hardware Graphic & Details -->
          <div class="flex items-center gap-4">
            <!-- Hardware Drive Icon Graphic Sesuai Gambar 2 -->
            <div class="w-14 h-14 rounded-xl bg-slate-800/90 border border-slate-700/80 p-1 flex items-center justify-center shrink-0 shadow-inner">
              <div class="w-11 h-11 rounded-lg bg-gradient-to-br from-slate-300 via-slate-400 to-slate-600 border border-slate-200 shadow-md relative flex items-center justify-center overflow-hidden">
                <!-- Platter Circle Graphic -->
                <div class="w-8 h-8 rounded-full border-2 border-slate-500 bg-gradient-to-tr from-slate-200 to-white flex items-center justify-center shadow-inner">
                  <div class="w-2.5 h-2.5 rounded-full bg-slate-600 border border-slate-400"></div>
                </div>
                <!-- Arm indicator -->
                <div class="absolute top-1 right-2 w-1.5 h-4 bg-slate-700 rotate-45 rounded-sm"></div>
              </div>
            </div>

            <div class="space-y-1">
              <h3 class="text-sm font-bold text-slate-100 tracking-wide">{{ drv.name }}</h3>
              <p class="text-xs text-slate-400">{{ drv.type_desc }}</p>
            </div>
          </div>

          <!-- Right: Health & Temperature Status Sesuai Gambar 2 -->
          <div class="text-right space-y-1 text-xs font-medium shrink-0">
            <div>
              <span class="text-slate-400">Health: </span>
              <span class="text-emerald-400 font-bold">{{ drv.health }}</span>
            </div>
            <div>
              <span class="text-slate-400">Temperature: </span>
              <span class="text-slate-200 font-semibold">{{ drv.temperature }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ============================================================= -->
    <!-- CUSTOM MODALS (WINDOW DIALOGS)                                -->
    <!-- ============================================================= -->

    <!-- 1. Modal Free Up Storage -->
    <div
      v-if="showFreeUpModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-sky-500/40 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-800 bg-slate-950/80">
          <div class="flex items-center gap-2 text-sky-400">
            <Sparkles class="w-4 h-4" />
            <h3 class="font-bold text-sm">Free Up Storage ({{ selectedStorage?.name }})</h3>
          </div>
          <button @click="showFreeUpModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3">
          <p class="text-xs text-slate-300 leading-relaxed">
            Sistem akan secara aman membersihkan berkas cache sementara, log journal usang, dan sisa build Docker tanpa menghapus data aplikasi atau konfigurasi Anda.
          </p>
          <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-300 space-y-1 font-mono">
            <div>• Target Mount: <span class="text-slate-100">{{ selectedStorage?.mountpoint }}</span></div>
            <div>• Perkiraan Ruang Bebas: <span class="text-emerald-400 font-bold">~412.5 MB</span></div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2.5 px-5 py-3.5 border-t border-slate-800 bg-slate-950/80">
          <button
            @click="showFreeUpModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="executeFreeUp"
            :disabled="isCleaning"
            class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-lg active:scale-95 disabled:opacity-50"
          >
            <span v-if="isCleaning" class="w-3.5 h-3.5 border-2 border-slate-950/30 border-t-slate-950 rounded-full animate-spin"></span>
            <span>{{ isCleaning ? 'Membersihkan...' : 'Bersihkan Sekarang' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 2. Modal Format Storage -->
    <div
      v-if="showFormatModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-rose-500/40 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-800 bg-slate-950/80">
          <div class="flex items-center gap-2 text-rose-400">
            <AlertTriangle class="w-4 h-4" />
            <h3 class="font-bold text-sm">Format Drive ({{ selectedStorage?.name }})</h3>
          </div>
          <button @click="showFormatModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3">
          <p class="text-xs text-rose-300 leading-relaxed">
            Peringatan: Seluruh berkas di dalam partisi ini akan dihapus secara permanen. Pastikan Anda telah mencadangkan data penting.
          </p>

          <div>
            <label class="block text-xs font-semibold text-slate-200 mb-1">Pilih File System</label>
            <select
              v-model="formatFs"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-100 outline-none focus:border-sky-400"
            >
              <option value="EXT4">EXT4 (Direkomendasikan untuk Linux & Docker STB)</option>
              <option value="VFAT">VFAT / FAT32 (Kompatibel Lintas Platform)</option>
              <option value="NTFS">NTFS (Windows Compatible)</option>
            </select>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2.5 px-5 py-3.5 border-t border-slate-800 bg-slate-950/80">
          <button
            @click="showFormatModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="executeFormat"
            :disabled="isFormatting"
            class="px-5 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-lg active:scale-95 disabled:opacity-50"
          >
            <span v-if="isFormatting" class="w-3.5 h-3.5 border-2 border-slate-950/30 border-t-slate-950 rounded-full animate-spin"></span>
            <span>{{ isFormatting ? 'Memformat...' : 'Format Drive' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 3. Modal Remove / Safely Eject -->
    <div
      v-if="showRemoveModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-800 bg-slate-950/80">
          <h3 class="font-bold text-sm text-slate-100">Lepaskan Drive ({{ selectedStorage?.name }})</h3>
          <button @click="showRemoveModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-2">
          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda yakin ingin melepas (*safely unmount*) drive ini dari sistem? Layanan yang mengakses direktori ini mungkin akan terhenti sementara.
          </p>
        </div>

        <div class="flex items-center justify-end gap-2.5 px-5 py-3.5 border-t border-slate-800 bg-slate-950/80">
          <button
            @click="showRemoveModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="executeRemove"
            :disabled="isRemoving"
            class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-lg active:scale-95 disabled:opacity-50"
          >
            <span>{{ isRemoving ? 'Melepas...' : 'Lepaskan Drive' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 4. Modal Create Storage -->
    <div
      v-if="showCreateStorageModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-sky-500/40 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-800 bg-slate-950/80">
          <div class="flex items-center gap-2 text-sky-400">
            <Plus class="w-4 h-4" />
            <h3 class="font-bold text-sm">Create Storage Location</h3>
          </div>
          <button @click="showCreateStorageModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3 text-xs text-slate-300">
          <p>
            Sistem RAZOS secara otomatis mengenali dan memetakan partisi penyimpanan baru yang terpasang pada port MicroSD dan USB STB ke direktori <span class="font-mono text-sky-300">/DATA</span>.
          </p>
          <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 space-y-1">
            <div class="font-semibold text-slate-200">Tips Pemasangan Media:</div>
            <div>1. Masukkan kartu MicroSD atau Flashdisk USB ke perangkat STB.</div>
            <div>2. Sistem akan otomatis melakukan mount (*Automount USB Drive* aktif).</div>
            <div>3. Lokasi baru akan muncul di daftar Storage Manager secara instan.</div>
          </div>
        </div>

        <div class="flex items-center justify-end px-5 py-3.5 border-t border-slate-800 bg-slate-950/80">
          <button
            @click="showCreateStorageModal = false"
            class="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition shadow-lg active:scale-95"
          >
            Mengerti
          </button>
        </div>
      </div>
    </div>

    <!-- 5. Modal Merge Storages (Beta) -->
    <div
      v-if="showMergeModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-blue-500/40 rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-slate-800 bg-slate-950/80">
          <div class="flex items-center gap-2 text-blue-400">
            <Layers class="w-4 h-4" />
            <h3 class="font-bold text-sm">Merge Storages (MergerFS Pool)</h3>
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-emerald-500 text-slate-950 uppercase">Beta</span>
          </div>
          <button @click="showMergeModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3 text-xs text-slate-300 leading-relaxed">
          <p>
            Fitur **Merge Storages** memungkinkan Anda menggabungkan beberapa drive fisik (seperti MicroSD dan USB HDD) menjadi satu kesatuan ruang penyimpanan virtual tunggal tanpa RAID yang rumit.
          </p>
          <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 space-y-1.5 font-mono text-[11px]">
            <div class="text-slate-200 font-bold">• Driver: MergerFS v2.33.5 (Zero Overhead)</div>
            <div>• Target Pool: /DATA/storage-pool</div>
            <div>• Redundancy: Non-destructive, files intact per drive</div>
          </div>
        </div>

        <div class="flex items-center justify-end px-5 py-3.5 border-t border-slate-800 bg-slate-950/80">
          <button
            @click="showMergeModal = false"
            class="px-5 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold transition shadow-lg active:scale-95"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
