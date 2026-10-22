<script setup>
/**
 * RAZOS External Cloud Synchronization View (rclone & WebDAV / Nextcloud Suite)
 * Updated: 2026-10-22 11:42:00
 * Deskripsi: Komponen antarmuka enterprise untuk menghubungkan server cloud eksternal
 * (Nextcloud, ownCloud, WebDAV Server, Google Drive, OneDrive) secara langsung dengan
 * filesystem RAZOS dan File Manager.
 */

import { ref, computed, onMounted } from 'vue'
import {
  Cloud, CloudLightning, RefreshCw, Plus, Trash2, Edit, Check, X,
  FolderCheck, HardDrive, CheckCircle2, AlertCircle, Play, Pause,
  ExternalLink, ArrowDownUp, ShieldCheck, Activity, Globe, Database,
  Sliders, ArrowUpRight, Folder, FileText
} from 'lucide-vue-next'
import { useCloudStore } from '../stores/RAZCloudStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useAuthStore } from '../stores/RAZAuthStore'

const cloudStore = useCloudStore()
const windowStore = useWindowStore()
const systemStore = useSystemStore()
const authStore = useAuthStore()

// Tab aktif: 'remotes' | 'logs'
const activeTab = ref('remotes')

// State Modal Tambah Remote
const showAddModal = ref(false)
const testResult = ref(null)
const isTesting = ref(false)

const newRemoteForm = ref({
  name: '',
  provider: 'nextcloud',
  url: '',
  username: '',
  password: '',
  remote_path: '/',
  local_mount_path: '/DATA/cloud/',
  sync_mode: 'two-way',
  sync_interval: '1h',
  auto_mount: true,
  auto_sync: true
})

// State Modal Edit Remote
const showEditModal = ref(false)
const editingRemote = ref(null)

// State Modal Konfirmasi Hapus
const showDeleteModal = ref(false)
const remoteToDelete = ref(null)

onMounted(async () => {
  await cloudStore.fetchRemotes()
})

// Format byte ke representasi human-readable
function formatBytes(bytes, decimals = 1) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

// Buka modal tambah remote
function openAddRemote() {
  testResult.value = null
  newRemoteForm.value = {
    name: '',
    provider: 'nextcloud',
    url: 'https://cloud.domain.com/remote.php/dav/files/user/',
    username: '',
    password: '',
    remote_path: '/',
    local_mount_path: '/DATA/cloud/Nextcloud',
    sync_mode: 'two-way',
    sync_interval: '1h',
    auto_mount: true,
    auto_sync: true
  }
  showAddModal.value = true
}

// Ubah template URL dan mount path otomatis saat provider diganti
function handleProviderChange() {
  const p = newRemoteForm.value.provider
  if (p === 'nextcloud') {
    newRemoteForm.value.url = 'https://cloud.domain.com/remote.php/dav/files/user/'
    newRemoteForm.value.local_mount_path = '/DATA/cloud/Nextcloud'
    if (!newRemoteForm.value.name) newRemoteForm.value.name = 'Nextcloud Hub'
  } else if (p === 'webdav') {
    newRemoteForm.value.url = 'https://dav.box.com/dav/ atau http://nas.local:5005/'
    newRemoteForm.value.local_mount_path = '/DATA/cloud/WebDAV'
    if (!newRemoteForm.value.name) newRemoteForm.value.name = 'WebDAV Server'
  } else if (p === 'googledrive') {
    newRemoteForm.value.url = 'gdrive:/Backups'
    newRemoteForm.value.local_mount_path = '/DATA/cloud/GoogleDrive'
    if (!newRemoteForm.value.name) newRemoteForm.value.name = 'Google Drive'
  } else if (p === 'onedrive') {
    newRemoteForm.value.url = 'onedrive:/Documents'
    newRemoteForm.value.local_mount_path = '/DATA/cloud/OneDrive'
    if (!newRemoteForm.value.name) newRemoteForm.value.name = 'OneDrive'
  }
}

// Jalankan Uji Koneksi (Test Connection)
async function runTestConnection() {
  isTesting.value = true
  testResult.value = null
  const res = await cloudStore.testConnection({
    provider: newRemoteForm.value.provider,
    url: newRemoteForm.value.url,
    username: newRemoteForm.value.username,
    password: newRemoteForm.value.password
  })
  isTesting.value = false
  testResult.value = res
}

// Simpan Remote Cloud Baru
async function submitAddRemote() {
  const f = newRemoteForm.value
  if (!f.name || !f.url) {
    systemStore.showNotification('Validasi', 'Nama dan URL server wajib diisi!', 'warning')
    return
  }

  const res = await cloudStore.addRemote(f)
  if (res.success) {
    systemStore.showNotification('Sukses', `Koneksi cloud ${f.name} berhasil ditambahkan!`, 'success')
    showAddModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Terjadi kesalahan saat menyimpan', 'error')
  }
}

// Buka Modal Edit Remote
function openEditRemote(remote) {
  editingRemote.value = JSON.parse(JSON.stringify(remote))
  showEditModal.value = true
}

async function submitEditRemote() {
  if (!editingRemote.value) return
  const r = editingRemote.value
  const res = await cloudStore.updateRemote(r.id, r)
  if (res.success) {
    systemStore.showNotification('Sukses', `Konfigurasi ${r.name} berhasil diperbarui!`, 'success')
    showEditModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal menyimpan pembaruan', 'error')
  }
}

// Buka Modal Hapus Remote
function openDeleteRemote(remote) {
  remoteToDelete.value = remote
  showDeleteModal.value = true
}

async function submitDeleteRemote() {
  if (!remoteToDelete.value) return
  const res = await cloudStore.deleteRemote(remoteToDelete.value.id)
  if (res.success) {
    systemStore.showNotification('Dihapus', `Koneksi cloud ${remoteToDelete.value.name} telah dilepas.`, 'info')
    showDeleteModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal menghapus remote', 'error')
  }
}

// Jalankan Sinkronisasi Manual
async function handleSyncNow(remote) {
  systemStore.showNotification('Sinkronisasi', `Memulai sinkronisasi ${remote.name}...`, 'info')
  const res = await cloudStore.syncRemote(remote.id, remote.sync_mode)
  if (res.success) {
    systemStore.showNotification('Sukses', `Sinkronisasi ${remote.name} berhasil!`, 'success')
  } else {
    systemStore.showNotification('Gagal Sync', res.error || 'Sinkronisasi mengalami kegagalan', 'error')
  }
}

// Pasang / Lepas Mount Drive
async function handleToggleMount(remote) {
  const actionText = remote.is_mounted ? 'Melepas mount' : 'Menghubungkan mount'
  systemStore.showNotification('Mount', `${actionText} ${remote.name}...`, 'info')
  const res = await cloudStore.toggleMount(remote.id)
  if (res.success) {
    systemStore.showNotification('Sukses', `Mount ${remote.name} berhasil diubah!`, 'success')
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal mengubah mount', 'error')
  }
}

// Buka Langsung Lokasi Mount di File Manager
function openInFileManager(remote) {
  windowStore.openApp('file-manager')
  systemStore.showNotification('File Manager', `Membuka penyimpanan cloud: ${remote.local_mount_path}`, 'info')
}
</script>

<template>
  <div class="w-full h-full flex flex-col bg-slate-950/90 text-slate-100 select-none overflow-hidden font-sans">
    <!-- Header Toolbar -->
    <div class="h-14 border-b border-slate-800/80 px-5 flex items-center justify-between shrink-0 bg-slate-900/60 backdrop-blur-md">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-sky-500/20 border border-sky-500/30 flex items-center justify-center text-sky-400">
          <CloudLightning class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-sm font-bold text-slate-100 flex items-center gap-2">
            Sinkronisasi Cloud Eksternal (Nextcloud & WebDAV Suite)
            <span class="text-[10px] px-2 py-0.5 rounded-full bg-slate-800 text-sky-400 border border-slate-700/60 font-mono">
              rclone RFC 4918
            </span>
          </h2>
          <p class="text-[11px] text-slate-400">Sinkronisasi berkas dua arah, auto-mount ke File Manager, dan integrasi storage cloud</p>
        </div>
      </div>

      <div class="flex items-center gap-2.5">
        <!-- Tab Switcher -->
        <div class="flex items-center p-1 bg-slate-900/80 rounded-xl border border-slate-800/80 text-xs">
          <button
            type="button"
            @click="activeTab = 'remotes'"
            class="px-3 py-1 rounded-lg font-medium transition"
            :class="activeTab === 'remotes' ? 'bg-sky-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200'"
          >
            Akun Terhubung ({{ cloudStore.remotes.length }})
          </button>
          <button
            type="button"
            @click="activeTab = 'logs'"
            class="px-3 py-1 rounded-lg font-medium transition"
            :class="activeTab === 'logs' ? 'bg-sky-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200'"
          >
            Log & Riwayat
          </button>
        </div>

        <button
          type="button"
          @click="cloudStore.fetchRemotes"
          class="p-2 rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700/60 transition active:scale-95"
          title="Muat Ulang Data"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': cloudStore.isLoading }" />
        </button>

        <button
          v-if="authStore.canMountCloud()"
          type="button"
          @click="openAddRemote"
          class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-sky-600 hover:bg-sky-500 text-white font-semibold text-xs transition active:scale-95 shadow-lg shadow-sky-600/20 border border-sky-400/30"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Akun Cloud</span>
        </button>
      </div>
    </div>

    <!-- Active Global Sync Notification Bar (Jika sedang ada remote yang sinkronisasi) -->
    <div
      v-if="cloudStore.isSyncing"
      class="px-5 py-2 bg-sky-950/80 border-b border-sky-500/40 flex items-center justify-between text-xs text-sky-200 animate-in fade-in"
    >
      <div class="flex items-center gap-2.5">
        <RefreshCw class="w-4 h-4 text-sky-400 animate-spin" />
        <span class="font-semibold">{{ cloudStore.syncProgress.message }}</span>
      </div>
      <div class="w-36 h-2 rounded-full bg-slate-900 border border-sky-500/30 overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-sky-400 to-emerald-400 transition-all duration-300"
          :style="{ width: `${cloudStore.syncProgress.percent}%` }"
        ></div>
      </div>
    </div>

    <!-- TAB 1: DAFTAR AKUN REMOTE CLOUD -->
    <div v-if="activeTab === 'remotes'" class="flex-1 overflow-y-auto p-5 scrollbar-thin">
      <div v-if="cloudStore.remotes.length === 0" class="h-64 flex flex-col items-center justify-center text-center text-slate-400">
        <Cloud class="w-12 h-12 text-slate-600 mb-2" />
        <p class="text-sm font-medium">Belum ada akun cloud eksternal yang terhubung</p>
        <p class="text-xs text-slate-500 mt-1">Klik "Tambah Akun Cloud" untuk menghubungkan Nextcloud, WebDAV, atau Google Drive</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="remote in cloudStore.remotes"
          :key="remote.id"
          class="rounded-2xl bg-slate-900/70 border border-slate-800/90 hover:border-sky-500/40 p-4 transition duration-200 flex flex-col justify-between shadow-xl"
        >
          <div>
            <!-- Top Header & Provider Info -->
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-3">
                <img
                  :src="remote.provider === 'nextcloud' ? '/assets/icons/RAZIconCloudNextcloud.svg' : '/assets/icons/RAZIconCloudWebdav.svg'"
                  alt="Provider"
                  class="w-10 h-10 rounded-xl object-contain shadow-md bg-slate-800 p-1 border border-slate-700/80"
                />
                <div>
                  <h3 class="text-sm font-bold text-slate-100 flex items-center gap-2">
                    {{ remote.name }}
                    <span
                      class="text-[9px] font-semibold uppercase px-1.5 py-0.5 rounded border tracking-wider"
                      :class="remote.status === 'connected' ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30' : 'bg-amber-500/20 text-amber-300 border-amber-500/30'"
                    >
                      {{ remote.status }}
                    </span>
                  </h3>
                  <p class="text-[11px] text-slate-400 font-mono truncate max-w-xs">{{ remote.url }}</p>
                </div>
              </div>

              <!-- Mount Status Badge -->
              <div
                class="flex items-center gap-1.5 px-2.5 py-1 rounded-xl text-xs font-semibold border transition"
                :class="remote.is_mounted ? 'bg-sky-500/20 text-sky-300 border-sky-500/40' : 'bg-slate-800 text-slate-400 border-slate-700'"
              >
                <FolderCheck class="w-3.5 h-3.5" />
                <span>{{ remote.is_mounted ? 'Ter-Mount' : 'Unmounted' }}</span>
              </div>
            </div>

            <!-- Storage Capacity Bar (CasaOS Style) -->
            <div class="p-3 rounded-xl bg-slate-950/60 border border-slate-800/80 space-y-2 mb-3">
              <div class="flex items-center justify-between text-xs">
                <span class="text-slate-400">Penyimpanan Terpakai:</span>
                <span class="font-mono text-slate-200 font-bold">
                  {{ formatBytes(remote.stats?.used_bytes) }} / {{ formatBytes(remote.stats?.total_bytes) }}
                </span>
              </div>
              <div class="w-full h-2 rounded-full bg-slate-800 overflow-hidden">
                <div
                  class="h-full bg-gradient-to-r from-sky-400 to-blue-500 transition-all duration-300"
                  :style="{ width: `${Math.min(100, Math.max(5, ((remote.stats?.used_bytes || 0) / (remote.stats?.total_bytes || 1)) * 100))}%` }"
                ></div>
              </div>
              <div class="flex items-center justify-between text-[10px] text-slate-500">
                <span>{{ remote.stats?.files_count || 0 }} Berkas Cloud</span>
                <span>Sinkronisasi Terakhir: {{ remote.last_sync ? new Date(remote.last_sync).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) : 'Belum' }}</span>
              </div>
            </div>

            <!-- Detail Mount & Sync Policy -->
            <div class="grid grid-cols-2 gap-2 text-[11px] mb-3">
              <div class="p-2 rounded-lg bg-slate-950/40 border border-slate-800/60">
                <span class="block text-[10px] text-slate-500 uppercase">Mount Point Lokal:</span>
                <span class="font-mono text-slate-300 truncate block">{{ remote.local_mount_path }}</span>
              </div>
              <div class="p-2 rounded-lg bg-slate-950/40 border border-slate-800/60">
                <span class="block text-[10px] text-slate-500 uppercase">Mode & Interval:</span>
                <span class="text-slate-300 capitalize truncate block">{{ remote.sync_mode }} ({{ remote.sync_interval }})</span>
              </div>
            </div>
          </div>

          <!-- Bottom Action Buttons -->
          <div class="pt-3 border-t border-slate-800/80 flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-1.5">
              <!-- Tombol Buka di File Manager -->
              <button
                type="button"
                @click="openInFileManager(remote)"
                class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-sky-400 text-xs font-semibold border border-slate-700/60 transition active:scale-95"
                title="Buka direktori ini di File Manager"
              >
                <Folder class="w-3.5 h-3.5" />
                <span>Buka di Files</span>
              </button>

              <!-- Tombol Sync Now -->
              <button
                type="button"
                @click="handleSyncNow(remote)"
                :disabled="cloudStore.isSyncing && cloudStore.syncingRemoteId === remote.id"
                class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl bg-sky-600/90 hover:bg-sky-500 text-white text-xs font-semibold shadow-md shadow-sky-600/20 transition active:scale-95"
              >
                <RefreshCw
                  class="w-3.5 h-3.5"
                  :class="{ 'animate-spin': cloudStore.isSyncing && cloudStore.syncingRemoteId === remote.id }"
                />
                <span>Sync Now</span>
              </button>
            </div>

            <div class="flex items-center gap-1">
              <!-- Toggle Mount/Unmount -->
              <button
                type="button"
                @click="handleToggleMount(remote)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                :title="remote.is_mounted ? 'Lepas Mount' : 'Hubungkan Mount'"
              >
                <HardDrive class="w-4 h-4" :class="remote.is_mounted ? 'text-sky-400' : 'text-slate-400'" />
              </button>

              <!-- Edit Config -->
              <button
                type="button"
                @click="openEditRemote(remote)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                title="Edit Konfigurasi"
              >
                <Edit class="w-4 h-4" />
              </button>

              <!-- Hapus Remote -->
              <button
                type="button"
                @click="openDeleteRemote(remote)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 transition"
                title="Hapus Remote"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: RIWAYAT & LOG SINKRONISASI -->
    <div v-else-if="activeTab === 'logs'" class="flex-1 overflow-y-auto p-5 scrollbar-thin">
      <div v-if="cloudStore.syncLogs.length === 0" class="h-64 flex flex-col items-center justify-center text-center text-slate-400">
        <FileText class="w-12 h-12 text-slate-600 mb-2" />
        <p class="text-sm font-medium">Belum ada riwayat aktivitas sinkronisasi</p>
      </div>

      <div v-else class="space-y-2 max-w-4xl mx-auto">
        <div
          v-for="log in cloudStore.syncLogs"
          :key="log.id"
          class="p-3 rounded-2xl bg-slate-900/80 border border-slate-800/80 flex items-start justify-between gap-3 text-xs"
        >
          <div class="flex items-start gap-3">
            <div
              class="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 mt-0.5"
              :class="log.status === 'success' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-rose-500/20 text-rose-400'"
            >
              <CheckCircle2 v-if="log.status === 'success'" class="w-4 h-4" />
              <AlertCircle v-else class="w-4 h-4" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-100">{{ log.remote_name }}</span>
                <span class="text-[10px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400 uppercase font-mono">{{ log.type }}</span>
              </div>
              <p class="text-slate-300 mt-0.5">{{ log.message }}</p>
            </div>
          </div>

          <div class="text-right shrink-0">
            <span class="text-[11px] font-mono text-slate-400">{{ new Date(log.timestamp).toLocaleTimeString('id-ID') }}</span>
            <span class="block text-[10px] text-slate-500">{{ log.duration_sec ? `${log.duration_sec}s` : '' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL: TAMBAH REMOTE CLOUD BARU                      -->
    <!-- ==================================================== -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-6 text-slate-100 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <CloudLightning class="w-5 h-5 text-sky-400" />
            <h3 class="text-sm font-bold text-slate-100">Hubungkan Akun Cloud Eksternal</h3>
          </div>
          <button @click="showAddModal = false" class="text-slate-400 hover:text-white p-1 rounded-lg">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitAddRemote" class="space-y-4 text-xs">
          <!-- Pilihan Provider -->
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1.5">Penyedia Layanan (Cloud Provider)</label>
            <div class="grid grid-cols-2 gap-2">
              <label
                v-for="prov in cloudStore.supportedProviders"
                :key="prov.id"
                class="p-2.5 rounded-xl border cursor-pointer transition flex items-center gap-2.5"
                :class="newRemoteForm.provider === prov.id ? 'bg-sky-950/40 border-sky-500 text-sky-200' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
              >
                <input
                  type="radio"
                  name="provRadio"
                  :value="prov.id"
                  v-model="newRemoteForm.provider"
                  @change="handleProviderChange"
                  class="accent-sky-500"
                />
                <div>
                  <span class="font-bold text-xs block text-slate-100">{{ prov.name }}</span>
                </div>
              </label>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Nama Tampilan Remote *</label>
              <input
                v-model="newRemoteForm.name"
                type="text"
                required
                placeholder="misal: Nextcloud Kantor"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Mount Path Lokal *</label>
              <input
                v-model="newRemoteForm.local_mount_path"
                type="text"
                required
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 font-mono text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">URL Endpoint Server WebDAV *</label>
            <input
              v-model="newRemoteForm.url"
              type="text"
              required
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 font-mono text-slate-100 focus:outline-none focus:border-sky-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Username / Akun</label>
              <input
                v-model="newRemoteForm.username"
                type="text"
                placeholder="username"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Password / App Token</label>
              <input
                v-model="newRemoteForm.password"
                type="password"
                placeholder="Password atau App Token"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Mode Sinkronisasi</label>
              <select
                v-model="newRemoteForm.sync_mode"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              >
                <option value="two-way">Dua Arah (Two-Way Sync)</option>
                <option value="upload-only">Upload Saja (Backup Lokal ke Cloud)</option>
                <option value="download-only">Download Saja (Cloud ke Lokal)</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Interval Auto-Sync</label>
              <select
                v-model="newRemoteForm.sync_interval"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              >
                <option value="15m">Setiap 15 Menit</option>
                <option value="1h">Setiap 1 Jam</option>
                <option value="6h">Setiap 6 Jam</option>
                <option value="24h">Harian (24 Jam)</option>
                <option value="manual">Manual Saja</option>
              </select>
            </div>
          </div>

          <!-- Uji Koneksi Status Box -->
          <div v-if="testResult" class="p-3 rounded-xl border text-xs" :class="testResult.success ? 'bg-emerald-950/40 border-emerald-500/50 text-emerald-200' : 'bg-rose-950/40 border-rose-500/50 text-rose-200'">
            <div class="flex items-center gap-2">
              <CheckCircle2 v-if="testResult.success" class="w-4 h-4 text-emerald-400 shrink-0" />
              <AlertCircle v-else class="w-4 h-4 text-rose-400 shrink-0" />
              <span>{{ testResult.message || testResult.error }}</span>
            </div>
          </div>

          <div class="flex items-center justify-between pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="runTestConnection"
              :disabled="isTesting"
              class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-sky-400 font-semibold flex items-center gap-1.5 transition active:scale-95 border border-slate-700/60"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isTesting }" />
              <span>{{ isTesting ? 'Menguji...' : 'Uji Koneksi Server' }}</span>
            </button>

            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="showAddModal = false"
                class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="cloudStore.isLoading"
                class="px-4 py-2 rounded-xl bg-sky-600 hover:bg-sky-500 text-white font-semibold flex items-center gap-1.5 shadow-lg shadow-sky-600/30"
              >
                <Check class="w-4 h-4" />
                <span>Simpan Remote</span>
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL: EDIT KONFIGURASI REMOTE                       -->
    <!-- ==================================================== -->
    <div
      v-if="showEditModal && editingRemote"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-6 text-slate-100 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <Edit class="w-5 h-5 text-sky-400" />
            <h3 class="text-sm font-bold text-slate-100">Edit Konfigurasi: {{ editingRemote.name }}</h3>
          </div>
          <button @click="showEditModal = false" class="text-slate-400 hover:text-white p-1 rounded-lg">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitEditRemote" class="space-y-4 text-xs">
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">Nama Tampilan</label>
            <input
              v-model="editingRemote.name"
              type="text"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
            />
          </div>

          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">URL Endpoint WebDAV</label>
            <input
              v-model="editingRemote.url"
              type="text"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 font-mono text-slate-100 focus:outline-none focus:border-sky-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Mode Sinkronisasi</label>
              <select
                v-model="editingRemote.sync_mode"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              >
                <option value="two-way">Dua Arah (Two-Way)</option>
                <option value="upload-only">Upload Saja (Backup)</option>
                <option value="download-only">Download Saja</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Interval Auto-Sync</label>
              <select
                v-model="editingRemote.sync_interval"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              >
                <option value="15m">Setiap 15 Menit</option>
                <option value="1h">Setiap 1 Jam</option>
                <option value="6h">Setiap 6 Jam</option>
                <option value="24h">Harian (24 Jam)</option>
                <option value="manual">Manual Saja</option>
              </select>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="showEditModal = false"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-sky-600 hover:bg-sky-500 text-white font-semibold"
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL: HAPUS REMOTE                                  -->
    <!-- ==================================================== -->
    <div
      v-if="showDeleteModal && remoteToDelete"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-rose-500/30 shadow-2xl p-5 text-slate-100 space-y-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-rose-500/20 border border-rose-500/30 flex items-center justify-center text-rose-400 shrink-0">
            <Cloud class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-slate-100">Hapus Remote Cloud?</h3>
            <p class="text-xs text-slate-400">Koneksi ke {{ remoteToDelete.name }} akan dihapus dan unmount secara aman.</p>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            type="button"
            @click="showDeleteModal = false"
            class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs"
          >
            Batal
          </button>
          <button
            type="button"
            @click="submitDeleteRemote"
            class="px-3.5 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs shadow-lg shadow-rose-600/30"
          >
            Hapus Remote
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
