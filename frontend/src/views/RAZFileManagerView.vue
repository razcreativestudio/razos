<script setup>
/**
 * RAZOS Advanced Storage & File Manager View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka File Manager canggih ala CasaOS yang mendeteksi HDD/USB fisik,
 * partisi MicroSD (/DATA), navigasi direktori hierarkis (breadcrumbs),
 * preview & edit berkas teks (editor), buat folder baru, rename, dan hapus berkas.
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import {
  HardDrive, Folder, File, FileText, FileCode, Plus, Trash2, Edit3,
  ChevronRight, ArrowLeft, RefreshCw, X, Save, AlertTriangle, Eye, Usb
} from 'lucide-vue-next'

const systemStore = useSystemStore()

// State Devices & Files
const storageDevices = ref([])
const currentPath = ref('/DATA')
const files = ref([])
const isLoading = ref(false)
const selectedItem = ref(null)

// State Modal Text Editor
const showEditorModal = ref(false)
const editingFileName = ref('')
const editingFilePath = ref('')
const editingContent = ref('')
const isSavingFile = ref(false)

// State Modal Buat Folder Baru
const showNewFolderModal = ref(false)
const newFolderName = ref('')

// State Modal Rename
const showRenameModal = ref(false)
const renameTargetOldPath = ref('')
const renameTargetNewName = ref('')

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// Fetch Devices / Partisi Storage
async function fetchDevices() {
  try {
    const res = await fetch('/api/v1/storage/devices')
    if (res.ok) {
      storageDevices.value = await res.json()
    }
  } catch (err) {}
}

// Fetch Daftar Berkas
async function fetchFiles(targetPath = currentPath.value) {
  isLoading.value = true
  selectedItem.value = null
  try {
    const res = await fetch(`/api/v1/files/list?path=${encodeURIComponent(targetPath)}`)
    if (res.ok) {
      files.value = await res.json()
      currentPath.value = targetPath
    } else {
      // Fallback ke path aman jika path tidak ada
      if (targetPath !== '.') {
        fetchFiles('.')
      }
    }
  } catch (err) {}
  isLoading.value = false
}

// Navigasi Breadcrumbs
const breadcrumbs = computed(() => {
  const parts = currentPath.value.split(/[\\/]/).filter(p => p.length > 0)
  const result = [{ name: 'Root', path: currentPath.value.startsWith('/') ? '/' : '.' }]
  let acc = currentPath.value.startsWith('/') ? '' : ''
  parts.forEach(part => {
    acc += (currentPath.value.startsWith('/') ? '/' : '\\') + part
    result.push({ name: part, path: acc })
  })
  return result
})

function navigateUp() {
  const parts = currentPath.value.split(/[\\/]/).filter(p => p.length > 0)
  if (parts.length <= 1) {
    fetchFiles('/')
    return
  }
  parts.pop()
  const up = (currentPath.value.startsWith('/') ? '/' : '') + parts.join(currentPath.value.startsWith('/') ? '/' : '\\')
  fetchFiles(up)
}

function handleItemClick(item) {
  selectedItem.value = item
}

function handleItemDoubleClick(item) {
  if (item.is_dir) {
    fetchFiles(item.path)
  } else {
    openFileEditor(item)
  }
}

// Buka Editor File Teks
async function openFileEditor(item) {
  editingFileName.value = item.name
  editingFilePath.value = item.path
  try {
    const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(item.path)}`)
    if (res.ok) {
      const data = await res.json()
      editingContent.value = data.content
      showEditorModal.value = true
    } else {
      systemStore.showNotification('Gagal', 'Tidak dapat membuka berkas sebagai teks', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Gagal membaca berkas', 'error')
  }
}

async function saveFileEditor() {
  isSavingFile.value = true
  try {
    const res = await fetch('/api/v1/files/content', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: editingFilePath.value, content: editingContent.value })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', 'Perubahan berkas berhasil disimpan', 'success')
      showEditorModal.value = false
      fetchFiles()
    }
  } catch (e) {}
  isSavingFile.value = false
}

// Buat Folder Baru
async function confirmCreateFolder() {
  if (!newFolderName.value.trim()) return
  try {
    const res = await fetch('/api/v1/files/folder', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ parent_path: currentPath.value, folder_name: newFolderName.value.trim() })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Folder ${newFolderName.value} berhasil dibuat`, 'success')
      showNewFolderModal.value = false
      newFolderName.value = ''
      fetchFiles()
    }
  } catch (e) {}
}

// Hapus Berkas/Folder
async function deleteItem(item) {
  try {
    const res = await fetch(`/api/v1/files/item?path=${encodeURIComponent(item.path)}`, { method: 'DELETE' })
    if (res.ok) {
      systemStore.showNotification('Sukses', `${item.name} berhasil dihapus`, 'success')
      fetchFiles()
    }
  } catch (e) {}
}

// Rename Berkas
function openRenameDialog(item) {
  renameTargetOldPath.value = item.path
  renameTargetNewName.value = item.name
  showRenameModal.value = true
}

async function confirmRename() {
  if (!renameTargetNewName.value.trim()) return
  const dir = currentPath.value
  const newPath = dir.endsWith('/') || dir.endsWith('\\') ? dir + renameTargetNewName.value : dir + '/' + renameTargetNewName.value
  try {
    const res = await fetch('/api/v1/files/rename', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ old_path: renameTargetOldPath.value, new_path: newPath })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', 'Nama berkas berhasil diubah', 'success')
      showRenameModal.value = false
      fetchFiles()
    }
  } catch (e) {}
}

onMounted(() => {
  fetchDevices()
  fetchFiles()
})
</script>

<template>
  <div class="h-full flex flex-col gap-3 text-xs sm:text-sm select-none">
    <!-- Top Header: Storage Devices Overview (CasaOS Style) -->
    <div class="flex items-center gap-2 overflow-x-auto pb-1 border-b border-slate-800">
      <div
        v-for="(dev, idx) in storageDevices"
        :key="idx"
        @click="fetchFiles(dev.mountpoint)"
        class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 hover:border-sky-400/50 cursor-pointer transition flex items-center gap-3 shrink-0"
        :class="currentPath.startsWith(dev.mountpoint) ? 'ring-1 ring-sky-400/50 bg-sky-950/20' : ''"
      >
        <div class="w-8 h-8 rounded-lg bg-sky-500/10 text-sky-400 flex items-center justify-center shrink-0">
          <component :is="dev.is_usb ? Usb : HardDrive" class="w-4 h-4" />
        </div>
        <div class="text-left">
          <div class="font-bold text-slate-200 text-xs">{{ dev.label }}</div>
          <div class="text-[10px] text-slate-400">
            {{ formatBytes(dev.used) }} / {{ formatBytes(dev.total) }} ({{ dev.percent.toFixed(0) }}%)
          </div>
          <div class="w-24 bg-slate-800 h-1 rounded-full overflow-hidden mt-1">
            <div class="bg-sky-400 h-full" :style="{ width: `${dev.percent}%` }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Navigation & Action Bar -->
    <div class="flex items-center justify-between gap-2 bg-slate-900/70 p-2 rounded-xl border border-slate-800">
      <!-- Back & Breadcrumb Path -->
      <div class="flex items-center gap-1.5 overflow-x-auto flex-1 font-mono text-xs">
        <button
          @click="navigateUp"
          class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
          title="Naik Satu Folder"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
        </button>

        <div class="flex items-center gap-1 text-slate-300">
          <template v-for="(crumb, idx) in breadcrumbs" :key="idx">
            <span
              @click="fetchFiles(crumb.path)"
              class="px-1.5 py-0.5 rounded hover:bg-slate-800 cursor-pointer text-slate-300 hover:text-sky-400"
            >
              {{ crumb.name }}
            </span>
            <ChevronRight v-if="idx < breadcrumbs.length - 1" class="w-3 h-3 text-slate-600" />
          </template>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-1.5 shrink-0">
        <button
          @click="showNewFolderModal = true"
          class="px-2.5 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1 transition shadow"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Folder</span>
        </button>

        <button
          @click="fetchFiles(currentPath)"
          class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
          title="Refresh"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="isLoading ? 'animate-spin' : ''" />
        </button>
      </div>
    </div>

    <!-- File Explorer List Table -->
    <div class="flex-1 rounded-xl bg-slate-900/50 border border-slate-800/80 p-2 overflow-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 text-slate-400">
            <th class="py-2 px-2">Nama</th>
            <th class="py-2 px-2">Ukuran</th>
            <th class="py-2 px-2">Format</th>
            <th class="py-2 px-2 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/50 text-slate-300">
          <tr
            v-for="(item, idx) in files"
            :key="idx"
            @click="handleItemClick(item)"
            @dblclick="handleItemDoubleClick(item)"
            class="hover:bg-slate-800/60 cursor-pointer transition rounded-lg"
            :class="selectedItem?.path === item.path ? 'bg-sky-500/10 text-sky-300' : ''"
          >
            <!-- Name + Icon -->
            <td class="py-2 px-2 font-medium flex items-center gap-2">
              <Folder v-if="item.is_dir" class="w-4 h-4 text-sky-400 fill-sky-400/20" />
              <FileCode v-else-if="['.js', '.json', '.go', '.vue', '.html', '.css', '.sh', '.yaml'].includes(item.extension)" class="w-4 h-4 text-emerald-400" />
              <FileText v-else-if="['.txt', '.md', '.log', '.env'].includes(item.extension)" class="w-4 h-4 text-amber-400" />
              <File v-else class="w-4 h-4 text-slate-400" />
              <span class="truncate max-w-[280px]">{{ item.name }}</span>
            </td>

            <!-- Size -->
            <td class="py-2 px-2 text-slate-400 font-mono">
              {{ item.is_dir ? '-' : formatBytes(item.size) }}
            </td>

            <!-- Type -->
            <td class="py-2 px-2 text-slate-400 uppercase text-[10px]">
              {{ item.is_dir ? 'Folder' : (item.extension ? item.extension.replace('.', '') : 'File') }}
            </td>

            <!-- Action Buttons -->
            <td class="py-2 px-2 text-right">
              <div class="flex items-center justify-end gap-1.5" @click.stop>
                <button
                  v-if="!item.is_dir"
                  @click="openFileEditor(item)"
                  class="p-1 rounded hover:bg-slate-700 text-slate-300"
                  title="Buka / Edit"
                >
                  <Eye class="w-3.5 h-3.5" />
                </button>
                <button
                  @click="openRenameDialog(item)"
                  class="p-1 rounded hover:bg-slate-700 text-slate-300"
                  title="Rename"
                >
                  <Edit3 class="w-3.5 h-3.5" />
                </button>
                <button
                  @click="deleteItem(item)"
                  class="p-1 rounded hover:bg-rose-500/20 text-rose-400"
                  title="Hapus"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Empty Folder Notice -->
      <div v-if="files.length === 0 && !isLoading" class="p-8 text-center text-slate-500">
        Folder ini kosong.
      </div>
    </div>

    <!-- CUSTOM MODAL: Text Editor File -->
    <div
      v-if="showEditorModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-2xl h-[80vh] raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <FileCode class="w-4 h-4 text-sky-400" />
            <span class="font-bold text-slate-100">{{ editingFileName }}</span>
          </div>
          <button @click="showEditorModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <textarea
          v-model="editingContent"
          class="flex-1 w-full p-3 bg-black/80 text-emerald-400 font-mono text-xs rounded-xl border border-slate-800 outline-none resize-none leading-relaxed"
        ></textarea>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            @click="showEditorModal = false"
            class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs"
          >
            Tutup
          </button>
          <button
            @click="saveFileEditor"
            :disabled="isSavingFile"
            class="px-4 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-1.5 shadow"
          >
            <Save class="w-3.5 h-3.5" />
            <span>Simpan Berkas</span>
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Buat Folder Baru -->
    <div
      v-if="showNewFolderModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-sm raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <h4 class="font-bold text-slate-100">Buat Folder Baru</h4>
        <input
          v-model="newFolderName"
          @keyup.enter="confirmCreateFolder"
          type="text"
          class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
          placeholder="Nama folder..."
        />
        <div class="flex items-center justify-end gap-2 pt-2">
          <button @click="showNewFolderModal = false" class="px-3 py-1.5 rounded-lg bg-slate-800 text-slate-300 text-xs">
            Batal
          </button>
          <button @click="confirmCreateFolder" class="px-4 py-1.5 rounded-lg bg-sky-500 text-slate-950 font-bold text-xs">
            Buat
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Rename File -->
    <div
      v-if="showRenameModal"
      class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-sm raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <h4 class="font-bold text-slate-100">Ubah Nama Berkas</h4>
        <input
          v-model="renameTargetNewName"
          @keyup.enter="confirmRename"
          type="text"
          class="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
        />
        <div class="flex items-center justify-end gap-2 pt-2">
          <button @click="showRenameModal = false" class="px-3 py-1.5 rounded-lg bg-slate-800 text-slate-300 text-xs">
            Batal
          </button>
          <button @click="confirmRename" class="px-4 py-1.5 rounded-lg bg-sky-500 text-slate-950 font-bold text-xs">
            Simpan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
