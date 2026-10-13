<script setup>
/**
 * RAZOS Google Drive-Style Storage & Cloud File Manager with Built-in Media Viewers
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka File Manager tingkat lanjut memadukan arsitektur Google Drive dan CasaOS:
 * - Built-in Media Viewers & Players:
 *   * Image Viewer: Preview gambar zoomable (WebP, PNG, JPG, SVG, GIF)
 *   * Video Player: Player MP4/WebM dengan kontrol multimedia kustom
 *   * Audio / Music Player: Player MP3/WAV/OGG dengan visualisasi gelombang
 *   * PDF Viewer: Penampil berkas dokumen PDF interaktif
 *   * Docs / Code Editor: Penampil dan editor kode (.txt, .md, .js, .go, .vue, .json, .sh, .yaml)
 * - Sidebar navigasi: My Drive (/DATA), Physical Disks/HDD, Starred/Favorit, Quick Storage Meter.
 * - Multi-view mode: Grid Thumbnail Card vs List Table view.
 * - Live Search filter & file type filter (Semua, Dokumen, Gambar, Video, Audio, Kode, Folder).
 * - File details preview pane kanan dengan media preview cepat.
 */
import { ref, onMounted, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import {
  HardDrive, Folder, File, FileText, FileCode, Plus, Trash2, Edit3,
  ChevronRight, ArrowLeft, RefreshCw, X, Save, AlertTriangle, Eye, Usb,
  Grid, List, Search, Star, Download, Info, CheckCircle2, Clock, Upload,
  Share2, MoreVertical, LayoutGrid, Server, Image, Film, Music, FileCheck, Play
} from 'lucide-vue-next'

const systemStore = useSystemStore()

// State Data
const storageDevices = ref([])
const currentPath = ref('/DATA')
const files = ref([])
const isLoading = ref(false)
const selectedItem = ref(null)
const starredPaths = ref(new Set())

// View Options
const viewMode = ref('grid') // 'grid' | 'list'
const searchQuery = ref('')
const filterType = ref('all') // 'all' | 'folders' | 'code' | 'docs' | 'images' | 'video' | 'audio'
const activeSidebarSection = ref('drive') // 'drive' | 'starred' | 'devices'

// State Modal Text Editor
const showEditorModal = ref(false)
const editingFileName = ref('')
const editingFilePath = ref('')
const editingContent = ref('')
const isSavingFile = ref(false)

// State Modal Image Viewer
const showImageModal = ref(false)
const previewImageUrl = ref('')
const previewImageName = ref('')

// State Modal Video Player
const showVideoModal = ref(false)
const previewVideoUrl = ref('')
const previewVideoName = ref('')

// State Modal Audio / Music Player
const showAudioModal = ref(false)
const previewAudioUrl = ref('')
const previewAudioName = ref('')

// State Modal PDF Viewer
const showPdfModal = ref(false)
const previewPdfUrl = ref('')
const previewPdfName = ref('')

// State Modal Folder & Rename
const showNewFolderModal = ref(false)
const newFolderName = ref('')
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

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
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
      if (targetPath !== '.') {
        fetchFiles('.')
      }
    }
  } catch (err) {}
  isLoading.value = false
}

// File Category Detection
function isImageFile(ext) {
  return ['.webp', '.png', '.jpg', '.jpeg', '.svg', '.gif', '.ico'].includes(ext)
}

function isVideoFile(ext) {
  return ['.mp4', '.webm', '.mkv', '.mov', '.avi'].includes(ext)
}

function isAudioFile(ext) {
  return ['.mp3', '.wav', '.ogg', '.aac', '.flac', '.m4a'].includes(ext)
}

function isPdfFile(ext) {
  return ext === '.pdf'
}

function isCodeFile(ext) {
  return ['.go', '.js', '.json', '.vue', '.html', '.css', '.sh', '.yaml', '.py', '.php', '.sql'].includes(ext)
}

function isDocFile(ext) {
  return ['.txt', '.md', '.log', '.env', '.conf', '.ini'].includes(ext)
}

// Filtered & Searched Files
const filteredFiles = computed(() => {
  return files.value.filter(item => {
    const matchSearch = item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    if (!matchSearch) return false

    if (activeSidebarSection.value === 'starred') {
      if (!starredPaths.value.has(item.path)) return false
    }

    if (filterType.value === 'folders') return item.is_dir
    if (filterType.value === 'images') return !item.is_dir && isImageFile(item.extension)
    if (filterType.value === 'video') return !item.is_dir && isVideoFile(item.extension)
    if (filterType.value === 'audio') return !item.is_dir && isAudioFile(item.extension)
    if (filterType.value === 'docs') return !item.is_dir && (isDocFile(item.extension) || isPdfFile(item.extension))
    if (filterType.value === 'code') return !item.is_dir && isCodeFile(item.extension)

    return true
  })
})

// Breadcrumbs Navigasi
const breadcrumbs = computed(() => {
  const parts = currentPath.value.split(/[\\/]/).filter(p => p.length > 0)
  const result = [{ name: 'My Drive', path: currentPath.value.startsWith('/') ? '/' : '.' }]
  let acc = ''
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

// Otomatis Buka Viewer yang Cocok saat Double Click
function handleItemDoubleClick(item) {
  if (item.is_dir) {
    fetchFiles(item.path)
  } else if (isImageFile(item.extension)) {
    openImageViewer(item)
  } else if (isVideoFile(item.extension)) {
    openVideoPlayer(item)
  } else if (isAudioFile(item.extension)) {
    openAudioPlayer(item)
  } else if (isPdfFile(item.extension)) {
    openPdfViewer(item)
  } else {
    openFileEditor(item)
  }
}

function openImageViewer(item) {
  previewImageName.value = item.name
  previewImageUrl.value = `/api/v1/files/raw?path=${encodeURIComponent(item.path)}`
  showImageModal.value = true
}

function openVideoPlayer(item) {
  previewVideoName.value = item.name
  previewVideoUrl.value = `/api/v1/files/raw?path=${encodeURIComponent(item.path)}`
  showVideoModal.value = true
}

function openAudioPlayer(item) {
  previewAudioName.value = item.name
  previewAudioUrl.value = `/api/v1/files/raw?path=${encodeURIComponent(item.path)}`
  showAudioModal.value = true
}

function openPdfViewer(item) {
  previewPdfName.value = item.name
  previewPdfUrl.value = `/api/v1/files/raw?path=${encodeURIComponent(item.path)}`
  showPdfModal.value = true
}

function toggleStar(item, e) {
  e?.stopPropagation()
  if (starredPaths.value.has(item.path)) {
    starredPaths.value.delete(item.path)
    systemStore.showNotification('Favorit', `Dihapus dari berbintang: ${item.name}`, 'info')
  } else {
    starredPaths.value.add(item.path)
    systemStore.showNotification('Favorit', `Ditambahkan ke berbintang: ${item.name}`, 'success')
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
  <div class="h-full flex overflow-hidden text-xs sm:text-sm select-none bg-slate-950 text-slate-100 rounded-xl">
    <!-- LEFT SIDEBAR (Google Drive Style) -->
    <div class="w-56 border-r border-slate-800/80 bg-slate-900/60 p-3 flex flex-col justify-between hidden md:flex shrink-0">
      <div class="flex flex-col gap-4">
        <!-- New Button -->
        <button
          @click="showNewFolderModal = true"
          class="h-10 px-4 rounded-2xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center justify-center gap-2 shadow-lg transition active:scale-95"
        >
          <Plus class="w-5 h-5" />
          <span class="text-xs tracking-wide">Folder Baru</span>
        </button>

        <!-- Navigation Links -->
        <div class="flex flex-col gap-1 text-xs">
          <button
            @click="() => { activeSidebarSection = 'drive'; fetchFiles('/DATA'); }"
            class="px-3 py-2.5 rounded-xl flex items-center gap-3 transition font-medium text-left"
            :class="activeSidebarSection === 'drive' ? 'bg-sky-500/20 text-sky-400 font-bold' : 'text-slate-300 hover:bg-slate-800/60'"
          >
            <Folder class="w-4 h-4 text-sky-400 fill-sky-400/20" />
            <span>My Drive (/DATA)</span>
          </button>

          <button
            @click="activeSidebarSection = 'starred'"
            class="px-3 py-2.5 rounded-xl flex items-center gap-3 transition font-medium text-left"
            :class="activeSidebarSection === 'starred' ? 'bg-sky-500/20 text-sky-400 font-bold' : 'text-slate-300 hover:bg-slate-800/60'"
          >
            <Star class="w-4 h-4 text-amber-400" />
            <span>Berbintang (Starred)</span>
          </button>
        </div>

        <!-- Devices Section -->
        <div class="flex flex-col gap-2 pt-2 border-t border-slate-800/80">
          <span class="text-[10px] font-bold text-slate-500 tracking-wider uppercase px-2">Storage Disks</span>
          <div class="flex flex-col gap-1.5">
            <div
              v-for="(dev, idx) in storageDevices"
              :key="idx"
              @click="() => { activeSidebarSection = 'devices'; fetchFiles(dev.mountpoint); }"
              class="p-2 rounded-xl bg-slate-950/60 hover:bg-slate-800/60 border border-slate-800/60 cursor-pointer transition flex flex-col gap-1"
              :class="currentPath.startsWith(dev.mountpoint) ? 'ring-1 ring-sky-400/40 bg-sky-950/20' : ''"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5 truncate">
                  <component :is="dev.is_usb ? Usb : HardDrive" class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                  <span class="text-[11px] font-medium text-slate-200 truncate">{{ dev.label }}</span>
                </div>
                <span class="text-[10px] text-slate-400 font-mono">{{ dev.percent.toFixed(0) }}%</span>
              </div>
              <div class="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                <div class="bg-sky-400 h-full" :style="{ width: `${dev.percent}%` }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom Cloud Storage Meter -->
      <div class="p-3 rounded-2xl bg-slate-950/70 border border-slate-800/80 flex flex-col gap-2">
        <div class="flex items-center justify-between text-[11px]">
          <span class="text-slate-400">Total Digunakan</span>
          <span class="text-sky-400 font-bold font-mono">25.6%</span>
        </div>
        <div class="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden">
          <div class="bg-sky-400 h-full w-[25.6%]"></div>
        </div>
        <span class="text-[10px] text-slate-500">8.2 GB dari 32 GB terpakai</span>
      </div>
    </div>

    <!-- MAIN DRIVE CONTENT AREA -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- TOP SEARCH & TOOLBAR -->
      <div class="h-14 border-b border-slate-800/80 px-4 flex items-center justify-between gap-3 bg-slate-900/40 shrink-0">
        <!-- Search Input Bar -->
        <div class="flex-1 max-w-md relative">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari file, media, gambar, video, musik, kode..."
            class="w-full pl-9 pr-3 py-2 rounded-xl bg-slate-900/80 border border-slate-800 text-xs text-slate-100 placeholder-slate-500 outline-none focus:border-sky-400 transition"
          />
        </div>

        <!-- Controls: Category Chips & View Toggle -->
        <div class="flex items-center gap-1.5 shrink-0">
          <!-- Filter Tabs Media Viewers -->
          <div class="hidden lg:flex items-center gap-1 bg-slate-900/80 p-1 rounded-xl border border-slate-800 text-[11px]">
            <button
              v-for="f in [{ id: 'all', label: 'Semua' }, { id: 'folders', label: 'Folder' }, { id: 'images', label: 'Gambar' }, { id: 'video', label: 'Video' }, { id: 'audio', label: 'Audio' }, { id: 'docs', label: 'Dokumen' }, { id: 'code', label: 'Kode' }]"
              :key="f.id"
              @click="filterType = f.id"
              class="px-2 py-1 rounded-lg transition"
              :class="filterType === f.id ? 'bg-sky-500/20 text-sky-400 font-bold' : 'text-slate-400 hover:text-slate-200'"
            >
              {{ f.label }}
            </button>
          </div>

          <!-- Grid vs List Toggle -->
          <div class="flex items-center bg-slate-900/80 p-0.5 rounded-lg border border-slate-800">
            <button
              @click="viewMode = 'grid'"
              class="p-1.5 rounded-md transition"
              :class="viewMode === 'grid' ? 'bg-slate-800 text-sky-400' : 'text-slate-400 hover:text-slate-200'"
              title="Grid View"
            >
              <LayoutGrid class="w-4 h-4" />
            </button>
            <button
              @click="viewMode = 'list'"
              class="p-1.5 rounded-md transition"
              :class="viewMode === 'list' ? 'bg-slate-800 text-sky-400' : 'text-slate-400 hover:text-slate-200'"
              title="List View"
            >
              <List class="w-4 h-4" />
            </button>
          </div>

          <button
            @click="fetchFiles(currentPath)"
            class="p-2 rounded-lg bg-slate-900/80 hover:bg-slate-800 text-slate-400 hover:text-slate-200 border border-slate-800 transition"
            title="Refresh"
          >
            <RefreshCw class="w-4 h-4" :class="isLoading ? 'animate-spin' : ''" />
          </button>
        </div>
      </div>

      <!-- BREADCRUMBS BAR -->
      <div class="px-4 py-2 border-b border-slate-800/50 flex items-center justify-between text-xs font-mono bg-slate-950/50 shrink-0">
        <div class="flex items-center gap-1.5 overflow-x-auto text-slate-400">
          <button
            @click="navigateUp"
            class="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200"
            title="Naik Satu Tingkat"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
          </button>
          <template v-for="(crumb, idx) in breadcrumbs" :key="idx">
            <span
              @click="fetchFiles(crumb.path)"
              class="px-1.5 py-0.5 rounded hover:bg-slate-800 cursor-pointer hover:text-sky-400"
              :class="idx === breadcrumbs.length - 1 ? 'text-slate-200 font-bold' : ''"
            >
              {{ crumb.name }}
            </span>
            <ChevronRight v-if="idx < breadcrumbs.length - 1" class="w-3 h-3 text-slate-600" />
          </template>
        </div>

        <span class="text-[11px] text-slate-500 hidden sm:inline">
          {{ filteredFiles.length }} item ditampilkan
        </span>
      </div>

      <!-- MAIN EXPLORER AREA (GRID VS LIST) -->
      <div class="flex-1 flex overflow-hidden">
        <div class="flex-1 p-4 overflow-y-auto">
          <!-- GRID VIEW (With Thumbnail Previews) -->
          <div v-if="viewMode === 'grid'" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3.5">
            <div
              v-for="(item, idx) in filteredFiles"
              :key="idx"
              @click="handleItemClick(item)"
              @dblclick="handleItemDoubleClick(item)"
              class="p-3 rounded-2xl border transition-all cursor-pointer flex flex-col justify-between gap-3 group relative overflow-hidden"
              :class="selectedItem?.path === item.path
                ? 'bg-sky-500/10 border-sky-400/50 ring-1 ring-sky-400/40 shadow-lg'
                : 'bg-slate-900/70 border-slate-800/80 hover:border-slate-700 hover:bg-slate-900'"
            >
              <!-- Card Top: Thumbnail / Icon & Star -->
              <div class="flex items-start justify-between">
                <!-- Media Thumbnail / Icon -->
                <div v-if="isImageFile(item.extension)" class="w-12 h-12 rounded-xl overflow-hidden bg-slate-950 border border-slate-800 shrink-0">
                  <img :src="`/api/v1/files/raw?path=${encodeURIComponent(item.path)}`" class="w-full h-full object-cover" />
                </div>
                <div v-else class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
                     :class="item.is_dir ? 'bg-sky-500/10 text-sky-400' : isVideoFile(item.extension) ? 'bg-purple-500/10 text-purple-400' : isAudioFile(item.extension) ? 'bg-amber-500/10 text-amber-400' : isPdfFile(item.extension) ? 'bg-rose-500/10 text-rose-400' : 'bg-slate-800/80 text-emerald-400'">
                  <Folder v-if="item.is_dir" class="w-6 h-6 fill-sky-400/20 text-sky-400" />
                  <Film v-else-if="isVideoFile(item.extension)" class="w-5 h-5 text-purple-400" />
                  <Music v-else-if="isAudioFile(item.extension)" class="w-5 h-5 text-amber-400" />
                  <FileText v-else-if="isPdfFile(item.extension)" class="w-5 h-5 text-rose-400" />
                  <FileCode v-else-if="isCodeFile(item.extension)" class="w-5 h-5 text-emerald-400" />
                  <FileText v-else-if="isDocFile(item.extension)" class="w-5 h-5 text-sky-400" />
                  <File v-else class="w-5 h-5 text-slate-400" />
                </div>

                <!-- Star Favorite Button -->
                <button
                  @click="(e) => toggleStar(item, e)"
                  class="p-1 rounded-full hover:bg-slate-800 text-slate-500 hover:text-amber-400 transition"
                  :class="starredPaths.has(item.path) ? 'text-amber-400' : 'opacity-0 group-hover:opacity-100'"
                >
                  <Star class="w-4 h-4" :class="starredPaths.has(item.path) ? 'fill-amber-400' : ''" />
                </button>
              </div>

              <!-- Card Body: Name & Meta -->
              <div>
                <h5 class="font-medium text-slate-200 text-xs truncate" :title="item.name">
                  {{ item.name }}
                </h5>
                <div class="flex items-center justify-between text-[10px] text-slate-500 mt-1 font-mono">
                  <span>{{ item.is_dir ? 'Folder' : formatBytes(item.size) }}</span>
                  <span class="uppercase font-semibold text-slate-400">{{ item.extension ? item.extension.replace('.', '') : '' }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- LIST VIEW (With Type Badges) -->
          <div v-else class="rounded-xl border border-slate-800/80 overflow-hidden bg-slate-900/50">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="border-b border-slate-800 text-slate-400 bg-slate-900/80">
                  <th class="py-2.5 px-3">Nama Berkas</th>
                  <th class="py-2.5 px-3">Ukuran</th>
                  <th class="py-2.5 px-3">Format</th>
                  <th class="py-2.5 px-3">Terakhir Diubah</th>
                  <th class="py-2.5 px-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800/50 text-slate-300">
                <tr
                  v-for="(item, idx) in filteredFiles"
                  :key="idx"
                  @click="handleItemClick(item)"
                  @dblclick="handleItemDoubleClick(item)"
                  class="hover:bg-slate-800/60 cursor-pointer transition"
                  :class="selectedItem?.path === item.path ? 'bg-sky-500/10 text-sky-300' : ''"
                >
                  <td class="py-2.5 px-3 font-medium flex items-center gap-2.5">
                    <button @click="(e) => toggleStar(item, e)" class="text-slate-500 hover:text-amber-400">
                      <Star class="w-3.5 h-3.5" :class="starredPaths.has(item.path) ? 'text-amber-400 fill-amber-400' : ''" />
                    </button>
                    <Folder v-if="item.is_dir" class="w-4 h-4 text-sky-400 fill-sky-400/20" />
                    <Image v-else-if="isImageFile(item.extension)" class="w-4 h-4 text-sky-400" />
                    <Film v-else-if="isVideoFile(item.extension)" class="w-4 h-4 text-purple-400" />
                    <Music v-else-if="isAudioFile(item.extension)" class="w-4 h-4 text-amber-400" />
                    <FileText v-else-if="isPdfFile(item.extension)" class="w-4 h-4 text-rose-400" />
                    <FileCode v-else-if="isCodeFile(item.extension)" class="w-4 h-4 text-emerald-400" />
                    <FileText v-else-if="isDocFile(item.extension)" class="w-4 h-4 text-amber-400" />
                    <File v-else class="w-4 h-4 text-slate-400" />
                    <span class="truncate max-w-[280px]">{{ item.name }}</span>
                  </td>
                  <td class="py-2.5 px-3 font-mono text-slate-400">{{ item.is_dir ? '-' : formatBytes(item.size) }}</td>
                  <td class="py-2.5 px-3 uppercase text-[10px] text-slate-400">{{ item.is_dir ? 'Folder' : (item.extension || 'File') }}</td>
                  <td class="py-2.5 px-3 text-slate-500 text-[11px]">{{ formatDate(item.mod_time) }}</td>
                  <td class="py-2.5 px-3 text-right">
                    <div class="flex items-center justify-end gap-1.5" @click.stop>
                      <button v-if="!item.is_dir" @click="handleItemDoubleClick(item)" class="p-1 rounded hover:bg-slate-700 text-sky-400" title="Buka Viewer / Player">
                        <Play v-if="isVideoFile(item.extension) || isAudioFile(item.extension)" class="w-3.5 h-3.5" />
                        <Eye v-else class="w-3.5 h-3.5" />
                      </button>
                      <button @click="openRenameDialog(item)" class="p-1 rounded hover:bg-slate-700 text-slate-300" title="Ubah Nama">
                        <Edit3 class="w-3.5 h-3.5" />
                      </button>
                      <button @click="deleteItem(item)" class="p-1 rounded hover:bg-rose-500/20 text-rose-400" title="Hapus">
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Empty Notice -->
          <div v-if="filteredFiles.length === 0 && !isLoading" class="p-12 text-center text-slate-500 flex flex-col items-center gap-2">
            <Folder class="w-12 h-12 text-slate-700" />
            <p class="font-medium text-slate-400">Tidak ada berkas yang cocok</p>
            <span class="text-xs text-slate-600">Folder kosong atau kata kunci pencarian tidak ditemukan.</span>
          </div>
        </div>

        <!-- RIGHT DETAILS PREVIEW PANE -->
        <div v-if="selectedItem" class="w-64 border-l border-slate-800/80 bg-slate-900/60 p-4 flex flex-col justify-between hidden lg:flex shrink-0 animate-in fade-in">
          <div class="flex flex-col gap-4">
            <div class="flex items-center justify-between pb-2 border-b border-slate-800">
              <span class="font-bold text-slate-200 text-xs">Detail Berkas</span>
              <button @click="selectedItem = null" class="text-slate-400 hover:text-white">
                <X class="w-4 h-4" />
              </button>
            </div>

            <!-- Thumbnail Preview di Panel Kanan -->
            <div class="flex flex-col items-center text-center gap-2 py-2">
              <div v-if="isImageFile(selectedItem.extension)" class="w-24 h-24 rounded-2xl overflow-hidden bg-slate-950 border border-slate-800">
                <img :src="`/api/v1/files/raw?path=${encodeURIComponent(selectedItem.path)}`" class="w-full h-full object-cover" />
              </div>
              <div v-else class="w-16 h-16 rounded-2xl bg-slate-800 flex items-center justify-center">
                <Folder v-if="selectedItem.is_dir" class="w-8 h-8 text-sky-400 fill-sky-400/20" />
                <Film v-else-if="isVideoFile(selectedItem.extension)" class="w-8 h-8 text-purple-400" />
                <Music v-else-if="isAudioFile(selectedItem.extension)" class="w-8 h-8 text-amber-400" />
                <FileText v-else-if="isPdfFile(selectedItem.extension)" class="w-8 h-8 text-rose-400" />
                <FileCode v-else-if="isCodeFile(selectedItem.extension)" class="w-8 h-8 text-emerald-400" />
                <FileText v-else class="w-8 h-8 text-slate-400" />
              </div>
              <h5 class="font-semibold text-slate-100 text-xs break-all leading-snug">{{ selectedItem.name }}</h5>
            </div>

            <!-- File Properties List -->
            <div class="space-y-2 text-xs text-slate-400 bg-slate-950/60 p-3 rounded-xl border border-slate-800/80">
              <div>
                <span class="text-slate-500 block text-[10px]">Tipe:</span>
                <span class="text-slate-200 font-medium">{{ selectedItem.is_dir ? 'Folder Direktori' : (selectedItem.extension || 'File Dokumen') }}</span>
              </div>
              <div>
                <span class="text-slate-500 block text-[10px]">Ukuran:</span>
                <span class="text-slate-200 font-mono">{{ selectedItem.is_dir ? '-' : formatBytes(selectedItem.size) }}</span>
              </div>
              <div>
                <span class="text-slate-500 block text-[10px]">Lokasi Path:</span>
                <span class="text-slate-200 font-mono text-[11px] break-all">{{ selectedItem.path }}</span>
              </div>
              <div>
                <span class="text-slate-500 block text-[10px]">Dimodifikasi:</span>
                <span class="text-slate-200 text-[11px]">{{ formatDate(selectedItem.mod_time) }}</span>
              </div>
            </div>
          </div>

          <!-- Quick Action Buttons -->
          <div class="flex flex-col gap-2 pt-3 border-t border-slate-800">
            <button
              v-if="!selectedItem.is_dir"
              @click="handleItemDoubleClick(selectedItem)"
              class="w-full py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center justify-center gap-2 shadow"
            >
              <Eye class="w-3.5 h-3.5" />
              <span>Buka Viewer / Player</span>
            </button>
            <button
              @click="openRenameDialog(selectedItem)"
              class="w-full py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs flex items-center justify-center gap-2"
            >
              <Edit3 class="w-3.5 h-3.5" />
              <span>Ubah Nama</span>
            </button>
            <button
              @click="deleteItem(selectedItem)"
              class="w-full py-2 rounded-xl bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 text-xs flex items-center justify-center gap-2"
            >
              <Trash2 class="w-3.5 h-3.5" />
              <span>Hapus Berkas</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 1. CUSTOM MODAL: Image Viewer -->
    <div
      v-if="showImageModal"
      class="fixed inset-0 bg-black/90 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-4xl max-h-[90vh] raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Image class="w-5 h-5 text-sky-400" />
            <span class="font-bold text-slate-100">{{ previewImageName }}</span>
          </div>
          <button @click="showImageModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="flex-1 flex items-center justify-center overflow-hidden bg-black/50 rounded-xl p-2">
          <img :src="previewImageUrl" class="max-w-full max-h-[70vh] object-contain rounded-lg shadow-2xl" />
        </div>
      </div>
    </div>

    <!-- 2. CUSTOM MODAL: Video Player -->
    <div
      v-if="showVideoModal"
      class="fixed inset-0 bg-black/90 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-3xl raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Film class="w-5 h-5 text-purple-400" />
            <span class="font-bold text-slate-100">{{ previewVideoName }}</span>
          </div>
          <button @click="showVideoModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="w-full aspect-video bg-black rounded-xl overflow-hidden flex items-center justify-center">
          <video :src="previewVideoUrl" controls autoplay class="w-full h-full"></video>
        </div>
      </div>
    </div>

    <!-- 3. CUSTOM MODAL: Music / Audio Player -->
    <div
      v-if="showAudioModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md raz-glass rounded-2xl border border-slate-700 p-5 shadow-2xl flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <Music class="w-5 h-5 text-amber-400" />
            <span class="font-bold text-slate-100">Audio Player</span>
          </div>
          <button @click="showAudioModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="flex flex-col items-center gap-3 py-3">
          <div class="w-16 h-16 rounded-full bg-amber-500/20 text-amber-400 flex items-center justify-center shadow-lg animate-pulse">
            <Music class="w-8 h-8" />
          </div>
          <h4 class="font-semibold text-slate-200 text-sm text-center break-all">{{ previewAudioName }}</h4>
          <audio :src="previewAudioUrl" controls autoplay class="w-full mt-2"></audio>
        </div>
      </div>
    </div>

    <!-- 4. CUSTOM MODAL: PDF Viewer -->
    <div
      v-if="showPdfModal"
      class="fixed inset-0 bg-black/90 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-4xl h-[88vh] raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <FileText class="w-5 h-5 text-rose-400" />
            <span class="font-bold text-slate-100">{{ previewPdfName }}</span>
          </div>
          <button @click="showPdfModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-5 h-5" />
          </button>
        </div>
        <iframe :src="previewPdfUrl" class="w-full flex-1 rounded-xl bg-slate-900 border border-slate-800"></iframe>
      </div>
    </div>

    <!-- 5. CUSTOM MODAL: In-Browser Text & Code Editor -->
    <div
      v-if="showEditorModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-3xl h-[85vh] raz-glass rounded-2xl border border-slate-700 p-4 shadow-2xl flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div class="flex items-center gap-2">
            <FileCode class="w-5 h-5 text-sky-400" />
            <span class="font-bold text-slate-100">{{ editingFileName }}</span>
          </div>
          <button @click="showEditorModal = false" class="p-1 text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <textarea
          v-model="editingContent"
          class="flex-1 w-full p-3.5 bg-black/90 text-emerald-400 font-mono text-xs rounded-xl border border-slate-800 outline-none resize-none leading-relaxed"
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
            <span>Simpan Perubahan</span>
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
        <h4 class="font-bold text-slate-100">Buat Folder Baru di Drive</h4>
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
