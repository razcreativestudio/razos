<script setup>
/**
 * RAZOS Google Drive-Style Storage & Windows File Explorer Hybrid File Manager
 * Updated: 2026-10-26 20:50:00 (v3.7.0 - Windows Multi-Selection, Dedicated Context Menu, and Custom Delete Confirmation Modal)
 * Deskripsi: Antarmuka File Manager tingkat lanjut memadukan arsitektur Windows File Explorer dan Google Drive:
 * - Windows Multi-Selection Engine: Pemilihan multi-berkas dengan Ctrl/Cmd+Click (toggle) dan Shift+Click (range selection).
 * - Keyboard Shortcuts: Ctrl+A (Select All), Escape (Clear Selection), Delete (Hapus Terpilih dengan Konfirmasi).
 * - Dedicated Context Menu: Menu konteks kontekstual untuk berkas tunggal, multi-berkas terpilih (batch actions), dan kanvas area kosong.
 * - Custom Modal Konfirmasi Hapus: Konfirmasi hapus berdesain Glassmorphism Slate-900 (Zero Native Browser Confirm/Alert).
 * - New Action Suite (+ Baru Dropdown): Pembuatan folder baru, file baru, upload file/folder/ZIP auto-extract.
 * - Deep RAZ Studio Code Integration: Membuka berkas kode/teks langsung di Sublime/Notepad++ complex editor.
 * - Fluid Responsive Grid & List View: Auto-fill layout yang tidak pernah gepeng dengan status bar lengkap.
 */
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import { useAuthStore } from '../stores/RAZAuthStore'
import { useCloudStore } from '../stores/RAZCloudStore'
import {
  HardDrive, Folder, File, FileText, FileCode, Plus, Trash2, Edit3,
  ChevronRight, ArrowLeft, ArrowRight, ArrowUp, RefreshCw, X, Save, AlertTriangle, Eye, Usb,
  Grid, List, Search, Star, Download, Info, CheckCircle2, Clock, Upload,
  Share2, MoreVertical, LayoutGrid, Server, Image, Film, Music, FileCheck, Play,
  Copy, Check, BookOpen, FilePlus, FolderPlus, FileUp, FolderUp, Archive, UploadCloud,
  ChevronDown, Sparkles, Loader2, Code2, Cloud, CloudUpload, ShieldAlert, Lock, ArrowUpDown
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()
const contextMenuStore = useContextMenuStore()
const authStore = useAuthStore()
const cloudStore = useCloudStore()

import RAZThumbnail from '../components/RAZThumbnail.vue'

// State Data
const storageDevices = ref([])
const currentPath = ref('/DATA')
const files = ref([])
const isLoading = ref(false)
const selectedItem = ref(null) // Item aktif terakhir
const selectedItems = ref([]) // Array berkas terpilih (multi-select ala Windows)
const lastSelectedIndex = ref(-1) // Index anchor untuk Shift+Click range

// Cek apakah item sedang dipilih
function isItemSelected(item) {
  if (!item || !item.path) return false
  return selectedItems.value.some(it => it.path === item.path)
}

// State Modal Konfirmasi Hapus Kustom (Zero Browser Native Alerts / Confirms)
const showDeleteConfirmModal = ref(false)
const itemsToDelete = ref([])
const isDeletingFiles = ref(false)

// State Sort By & Urutan Tampilan (Windows 11 Sorting Engine)
const sortBy = ref('name') // 'name' | 'date' | 'type' | 'size'
const sortOrder = ref('asc') // 'asc' | 'desc'
const showSortDropdown = ref(false)
const showViewDropdown = ref(false)

function setSort(by) {
  sortBy.value = by
  showSortDropdown.value = false
}

function setOrder(order) {
  sortOrder.value = order
  showSortDropdown.value = false
}

// Tutup seluruh dropdown ribbon bar
function closeRibbonDropdowns() {
  showSortDropdown.value = false
  showViewDropdown.value = false
  showRibbonNewMenu.value = false
}

// State Starred / Favorit Persisten (Tersimpan Permanen di localStorage)
const starredItems = ref([])

function loadStarredItems() {
  try {
    const raw = localStorage.getItem('razos_starred_items')
    if (raw) {
      starredItems.value = JSON.parse(raw)
    }
  } catch (e) {
    starredItems.value = []
  }
}

function saveStarredItems() {
  try {
    localStorage.setItem('razos_starred_items', JSON.stringify(starredItems.value))
  } catch (e) {}
}

function isItemStarred(itemPath) {
  if (!itemPath) return false
  return starredItems.value.some(it => it.path === itemPath)
}

// Inisialisasi awal koleksi berbintang
loadStarredItems()

// State Address Bar & Navigasi Windows Explorer Style
const historyStack = ref([])
const historyIndex = ref(-1)
const isEditingAddress = ref(false)
const addressInputValue = ref('/DATA')
const addressInputRef = ref(null)
const pathCopied = ref(false)

const canGoBack = computed(() => historyIndex.value > 0)
const canGoForward = computed(() => historyIndex.value < historyStack.value.length - 1)

function goBack() {
  if (canGoBack.value) {
    historyIndex.value--
    fetchFiles(historyStack.value[historyIndex.value], true)
  }
}

function goForward() {
  if (canGoForward.value) {
    historyIndex.value++
    fetchFiles(historyStack.value[historyIndex.value], true)
  }
}

function startEditingAddress() {
  isEditingAddress.value = true
  addressInputValue.value = currentPath.value
  nextTick(() => {
    if (addressInputRef.value) {
      addressInputRef.value.focus()
      addressInputRef.value.select()
    }
  })
}

function submitAddress() {
  if (!isEditingAddress.value) return
  isEditingAddress.value = false
  const target = addressInputValue.value.trim()
  if (target) {
    fetchFiles(target)
  }
}

function cancelEditingAddress() {
  isEditingAddress.value = false
  addressInputValue.value = currentPath.value
}

function onAddressInputBlur() {
  // Berikan jeda toleransi klik tombol Go
  setTimeout(() => {
    if (isEditingAddress.value) {
      cancelEditingAddress()
    }
  }, 150)
}

async function copyCurrentPath() {
  try {
    await navigator.clipboard.writeText(currentPath.value)
    pathCopied.value = true
    systemStore.showNotification('Disalin', `Path "${currentPath.value}" disalin ke clipboard`, 'success')
    setTimeout(() => {
      pathCopied.value = false
    }, 2000)
  } catch (err) {
    systemStore.showNotification('Gagal', 'Tidak dapat menyalin path ke clipboard', 'error')
  }
}

const currentFolderShortName = computed(() => {
  const norm = normalizePath(currentPath.value)
  if (norm === '/' || norm === 'C:/') return 'Root'
  const parts = norm.split('/').filter(Boolean)
  return parts.length > 0 ? parts[parts.length - 1] : 'folder'
})

// View Options
const viewMode = ref('grid') // 'grid' | 'list'
const searchQuery = ref('')
const filterType = ref('all') // 'all' | 'folders' | 'code' | 'docs' | 'images' | 'video' | 'audio'
const activeSidebarSection = ref('drive') // 'drive' | 'starred' | 'devices'

// State Dropdown Menu "+ Baru" (Sidebar & Ribbon Top Bar)
const showNewMenu = ref(false)
const showRibbonNewMenu = ref(false)

function closeAllNewMenus() {
  showNewMenu.value = false
  showRibbonNewMenu.value = false
  showSortDropdown.value = false
}

// State Modal Folder & Rename
const showNewFolderModal = ref(false)
const newFolderName = ref('')
const showRenameModal = ref(false)
const renameTargetOldPath = ref('')
const renameTargetNewName = ref('')

// State Modal Buat File Baru (Teks / Kode)
const showNewFileModal = ref(false)
const newFileName = ref('')
const newFileTemplate = ref('txt') // 'txt' | 'js' | 'html' | 'css' | 'json' | 'py' | 'go' | 'md' | 'yaml'

// State Upload File, Folder, dan ZIP Auto-Extract
const fileUploadInputRef = ref(null)
const folderUploadInputRef = ref(null)
const zipUploadInputRef = ref(null)

const isUploading = ref(false)
const uploadProgress = ref(0) // 0 - 100%
const uploadCurrentFileName = ref('')
const uploadTotalCount = ref(0)
const uploadCompletedCount = ref(0)
const showUploadModal = ref(false)
const uploadType = ref('file') // 'file' | 'folder' | 'zip'

// Buka Jendela Mandiri Info & Properti Berkas (v2.6.0)
function openInfoModal(item = selectedItem.value) {
  if (!item) return
  selectedItem.value = item
  windowStore.openFileProperties(item)
}

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

// Path Default untuk My Storage (/DATA)
const defaultStoragePath = ref('/DATA')

// Helper Fungsi Normalisasi Path Universal (Windows & Linux STB)
function normalizePath(p) {
  if (!p) return '/'
  let clean = p.replace(/\\/g, '/').replace(/\/+/g, '/')
  // Jika format drive letter Windows "C:" tanpa slash, tambahkan slash "C:/"
  if (/^[a-zA-Z]:$/.test(clean)) {
    clean += '/'
  }
  // Hilangkan trailing slash jika bukan root drive atau bukan root "/"
  if (clean.length > 3 && clean.endsWith('/')) {
    clean = clean.slice(0, -1)
  }
  return clean
}

// Fetch Devices / Partisi Storage
async function fetchDevices() {
  try {
    const res = await fetch('/api/v1/storage/devices')
    if (res.ok) {
      const data = await res.json()
      storageDevices.value = data
      
      // Deteksi lokasi primary storage (/DATA)
      const dataDev = data.find(d => d.mountpoint === '/DATA' || (d.label && d.label.includes('DATA')))
      if (dataDev) {
        defaultStoragePath.value = dataDev.mountpoint
      } else {
        // Pada host Windows dev, arahkan ke C:/DATA jika ada drive C
        const hasCDrive = data.some(d => (d.device && d.device.startsWith('C:')) || (d.mountpoint && d.mountpoint.startsWith('C:')))
        if (hasCDrive) {
          defaultStoragePath.value = 'C:/DATA'
        } else if (data.length > 0 && data[0].mountpoint) {
          defaultStoragePath.value = data[0].mountpoint
        }
      }
    }
  } catch (err) {
    console.error('Gagal mengambil daftar storage devices:', err)
  }
}

// Partisi Fisik dan Remote Cloud (rclone/WebDAV/Nextcloud)
const physicalDisks = computed(() => {
  return storageDevices.value.filter(dev => !dev.is_cloud)
})

const cloudDisks = computed(() => {
  return storageDevices.value.filter(dev => dev.is_cloud)
})

// Komputasi Total Penyimpanan Riil dari seluruh Disk/Partisi yang terbaca
const totalStorageBytes = computed(() => {
  return storageDevices.value.reduce((acc, dev) => acc + (dev.total || 0), 0)
})

const usedStorageBytes = computed(() => {
  return storageDevices.value.reduce((acc, dev) => acc + (dev.used || 0), 0)
})

const storageUsagePercent = computed(() => {
  if (totalStorageBytes.value === 0) return 0
  return parseFloat(((usedStorageBytes.value / totalStorageBytes.value) * 100).toFixed(1))
})

// Fetch Daftar Berkas dengan Normalisasi Aman & History Tracking
async function fetchFiles(targetPath = currentPath.value, isHistoryNav = false) {
  const cleanPath = normalizePath(targetPath)
  isLoading.value = true
  selectedItem.value = null
  displayLimit.value = 60
  try {
    const res = await fetch(`/api/v1/files/list?path=${encodeURIComponent(cleanPath)}`)
    const data = await res.json()

    // Tangani jika folder dilindungi izin sistem operasi (Access Denied)
    if (data && data.isAccessDenied) {
      systemStore.showNotification('Akses Ditolak', data.error || `Direktori "${cleanPath}" dilindungi sistem operasi.`, 'warning')
      addressInputValue.value = currentPath.value
      isLoading.value = false
      return
    }

    // Tangani jika direktori kosong (server bisa mengembalikan null atau [])
    const fileList = Array.isArray(data) ? data : (data === null ? [] : null)

    if (res.ok && fileList !== null) {
      // Periksa apakah server me-resolve Windows Junction (misal My Documents -> Documents)
      const resolvedHeader = res.headers.get('X-Resolved-Path')
      const actualPath = resolvedHeader ? normalizePath(resolvedHeader) : cleanPath

      files.value = fileList
      currentPath.value = actualPath
      addressInputValue.value = actualPath

      // Rekam jejak ke history stack jika bukan navigasi back/forward
      if (!isHistoryNav) {
        if (historyIndex.value < historyStack.value.length - 1) {
          historyStack.value = historyStack.value.slice(0, historyIndex.value + 1)
        }
        if (historyStack.value.length === 0 || historyStack.value[historyStack.value.length - 1] !== actualPath) {
          historyStack.value.push(actualPath)
          historyIndex.value = historyStack.value.length - 1
        }
      }
    } else {
      // Notifikasi peringatan jika path benar-benar tidak ditemukan (404/500)
      systemStore.showNotification('Peringatan', `Direktori "${cleanPath}" tidak ditemukan`, 'warning')
      addressInputValue.value = currentPath.value
    }
  } catch (err) {
    systemStore.showNotification('Gagal', `Gagal memuat direktori: ${err.message}`, 'error')
    addressInputValue.value = currentPath.value
  }
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

function isMarkdownFile(ext) {
  return ['.md', '.markdown', '.mdown', '.mkd'].includes(ext)
}

function isDocFile(ext) {
  return ['.txt', '.log', '.env', '.conf', '.ini'].includes(ext)
}

// Filtered, Searched & Sorted Files (Windows 11 Sorting Engine)
const filteredFiles = computed(() => {
  // Jika tab Berbintang aktif, gunakan koleksi global starredItems, jika tidak gunakan files folder aktif
  const sourceList = activeSidebarSection.value === 'starred' ? starredItems.value : files.value
  let list = sourceList.filter(item => {
    const matchSearch = item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    if (!matchSearch) return false

    if (filterType.value === 'folders') return item.is_dir
    if (filterType.value === 'images') return !item.is_dir && isImageFile(item.extension)
    if (filterType.value === 'video') return !item.is_dir && isVideoFile(item.extension)
    if (filterType.value === 'audio') return !item.is_dir && isAudioFile(item.extension)
    if (filterType.value === 'docs') return !item.is_dir && (isDocFile(item.extension) || isPdfFile(item.extension) || isMarkdownFile(item.extension))
    if (filterType.value === 'code') return !item.is_dir && isCodeFile(item.extension)

    return true
  })

  // SORTING ENGINE WINDOWS 11 EXPLORER
  list.sort((a, b) => {
    // 1. Folder selalu berada di atas berkas
    if (a.is_dir && !b.is_dir) return -1
    if (!a.is_dir && b.is_dir) return 1

    let comp = 0
    if (sortBy.value === 'name') {
      comp = a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
    } else if (sortBy.value === 'date') {
      const timeA = new Date(a.mod_time).getTime() || 0
      const timeB = new Date(b.mod_time).getTime() || 0
      comp = timeA - timeB
    } else if (sortBy.value === 'type') {
      const extA = (a.extension || (a.is_dir ? 'folder' : '')).toLowerCase()
      const extB = (b.extension || (b.is_dir ? 'folder' : '')).toLowerCase()
      comp = extA.localeCompare(extB)
      if (comp === 0) {
        comp = a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
      }
    } else if (sortBy.value === 'size') {
      const sizeA = Number(a.size) || 0
      const sizeB = Number(b.size) || 0
      comp = sizeA - sizeB
    }

    return sortOrder.value === 'asc' ? comp : -comp
  })

  return list
})

// Progressive Virtual Display Rendering (Anti-Hang saat folder memuat ribuan berkas)
const displayLimit = ref(60)
const visibleFiles = computed(() => {
  return filteredFiles.value.slice(0, displayLimit.value)
})
const hasMoreFiles = computed(() => {
  return visibleFiles.value.length < filteredFiles.value.length
})

function handleExplorerScroll(e) {
  const el = e.target
  if (!el) return
  // Jika mendekati 300px dari dasar scroll container, muat batch 40 berkas berikutnya
  if (el.scrollHeight - (el.scrollTop + el.clientHeight) < 300 && hasMoreFiles.value) {
    displayLimit.value += 40
  }
}

// Reset limit display saat pencarian, filter, atau kategori sidebar berganti
watch([searchQuery, filterType, activeSidebarSection], () => {
  displayLimit.value = 60
})

// Breadcrumbs Navigasi yang Akurat & Terstruktur
const breadcrumbs = computed(() => {
  if (activeSidebarSection.value === 'starred') {
    return [{ name: '⭐ Berbintang (Starred)', path: '' }]
  }
  const norm = normalizePath(currentPath.value)
  const result = []

  // Skenario Windows Path (diawali "C:/" atau "C:/DATA/...")
  const winMatch = norm.match(/^([a-zA-Z]:)(\/.*)?$/)
  if (winMatch) {
    const driveLetter = winMatch[1]
    const rest = winMatch[2] || ''

    // Root drive node
    result.push({ name: `${driveLetter} Disk`, path: `${driveLetter}/` })

    const segments = rest.split('/').filter(s => s.length > 0)
    let accumulated = `${driveLetter}`
    segments.forEach((seg) => {
      accumulated += '/' + seg
      // Beri badge khusus jika ini folder DATA (Cloud Storage)
      const label = (accumulated.toUpperCase() === 'C:/DATA') ? 'My Storage (/DATA)' : seg
      result.push({ name: label, path: accumulated })
    })
    return result
  }

  // Skenario Linux / Unix Path (diawali "/")
  result.push({ name: 'System Root', path: '/' })
  const segments = norm.split('/').filter(s => s.length > 0)
  let accumulated = ''
  segments.forEach((seg) => {
    accumulated += '/' + seg
    const label = (accumulated === '/DATA') ? 'My Storage (/DATA)' : seg
    result.push({ name: label, path: accumulated })
  })
  return result
})

// Cek apakah tombol panah kembali bisa diklik (tidak sedang di root)
const canNavigateUp = computed(() => {
  const norm = normalizePath(currentPath.value)
  if (norm === '/' || /^[a-zA-Z]:\/$/.test(norm)) return false
  return true
})

// Navigasi Naik Satu Tingkat
function navigateUp() {
  const norm = normalizePath(currentPath.value)
  if (!canNavigateUp.value) return

  // Skenario Windows: potong segmen terakhir dengan aman
  const winMatch = norm.match(/^([a-zA-Z]:)(\/.*)?$/)
  if (winMatch) {
    const driveLetter = winMatch[1]
    const rest = winMatch[2] || ''
    const segments = rest.split('/').filter(s => s.length > 0)
    if (segments.length <= 1) {
      // Tinggal 1 segmen (misal "C:/DATA"), kembali ke root "C:/"
      fetchFiles(`${driveLetter}/`)
      return
    }
    segments.pop()
    fetchFiles(`${driveLetter}/${segments.join('/')}`)
    return
  }

  // Skenario Linux: potong dari slash terakhir
  const lastSlash = norm.lastIndexOf('/')
  if (lastSlash <= 0) {
    fetchFiles('/')
    return
  }
  fetchFiles(norm.substring(0, lastSlash))
}

/**
 * Menangani pemilihan berkas ala Windows File Explorer:
 * - Klik biasa: Pilih satu berkas saja
 * - Ctrl + Klik / Cmd + Klik: Toggle status pilihan berkas tanpa mengganggu pilihan lainnya
 * - Shift + Klik: Pemilihan rentang berkas (Range Selection) dari item terakhir ke item yang diklik
 */
function handleItemClick(item, event) {
  if (!item) return

  const currentIndex = visibleFiles.value.findIndex(f => f.path === item.path)

  // 1. Kasus Ctrl / Cmd Key: Toggle Seleksi
  if (event && (event.ctrlKey || event.metaKey)) {
    const existingIndex = selectedItems.value.findIndex(f => f.path === item.path)
    if (existingIndex >= 0) {
      selectedItems.value.splice(existingIndex, 1)
      selectedItem.value = selectedItems.value[selectedItems.value.length - 1] || null
    } else {
      selectedItems.value.push(item)
      selectedItem.value = item
    }
    lastSelectedIndex.value = currentIndex
    return
  }

  // 2. Kasus Shift Key: Range Selection (Rentang Berkas ala Windows Explorer)
  if (event && event.shiftKey && lastSelectedIndex.value >= 0 && currentIndex >= 0) {
    const start = Math.min(lastSelectedIndex.value, currentIndex)
    const end = Math.max(lastSelectedIndex.value, currentIndex)
    const range = visibleFiles.value.slice(start, end + 1)
    selectedItems.value = [...range]
    selectedItem.value = item
    return
  }

  // 3. Kasus Klik Biasa: Pilih hanya berkas ini
  selectedItems.value = [item]
  selectedItem.value = item
  lastSelectedIndex.value = currentIndex
}

/**
 * Pilih semua berkas di direktori aktif (Ctrl + A)
 */
function selectAllFiles() {
  selectedItems.value = [...visibleFiles.value]
  if (selectedItems.value.length > 0) {
    selectedItem.value = selectedItems.value[0]
    lastSelectedIndex.value = 0
  }
}

/**
 * Bersihkan seluruh pilihan berkas (Escape / Klik area kosong)
 */
function clearSelection() {
  selectedItems.value = []
  selectedItem.value = null
  lastSelectedIndex.value = -1
}

/**
 * Handler shortcut keyboard global File Manager
 */
function handleGlobalKeydown(e) {
  // Cegah interupsi jika pengguna sedang mengetik di input form atau dialog
  if (['INPUT', 'TEXTAREA'].includes(e.target.tagName)) return

  // Ctrl + A / Cmd + A: Pilih Semua Berkas
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'a') {
    e.preventDefault()
    selectAllFiles()
    return
  }

  // Escape: Batalkan Seleksi
  if (e.key === 'Escape') {
    clearSelection()
    return
  }

  // Delete: Buka modal konfirmasi hapus untuk item yang dipilih
  if (e.key === 'Delete' && selectedItems.value.length > 0 && authStore.canDeleteStorage()) {
    e.preventDefault()
    openDeleteConfirm(selectedItems.value)
    return
  }
}

// Buka Berkas di Window Mandiri (Standalone Window Viewers & RAZ Studio Code)
function handleItemDoubleClick(item) {
  if (!item) return
  if (item.is_dir) {
    if (activeSidebarSection.value === 'starred') {
      activeSidebarSection.value = 'drive'
    }
    fetchFiles(normalizePath(item.path))
  } else {
    // Jika berkas adalah kode, teks, markdown, script, atau konfigurasi, buka di RAZ Studio Code
    if (isCodeFile(item.extension) || isDocFile(item.extension) || isMarkdownFile(item.extension) || !item.extension) {
      openInTextEditor(item)
    } else if (typeof windowStore.openFileViewer === 'function') {
      windowStore.openFileViewer(item)
    } else {
      openInTextEditor(item)
    }
  }
}

function toggleStar(item, e) {
  e?.stopPropagation()
  if (!item || !item.path) return
  const idx = starredItems.value.findIndex(it => it.path === item.path)
  if (idx >= 0) {
    starredItems.value.splice(idx, 1)
    systemStore.showNotification('Favorit', `Dihapus dari berbintang: ${item.name}`, 'info')
  } else {
    starredItems.value.push({
      name: item.name,
      path: item.path,
      is_dir: !!item.is_dir,
      size: item.size || 0,
      extension: item.extension || '',
      mod_time: item.mod_time || new Date().toISOString()
    })
    systemStore.showNotification('Favorit', `Ditambahkan ke berbintang: ${item.name}`, 'success')
  }
  saveStarredItems()
}

// Buat Folder Baru (RBAC Guarded)
async function confirmCreateFolder() {
  if (!authStore.canWriteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk membuat folder.', 'warning')
    return
  }
  if (!newFolderName.value.trim()) return
  try {
    const res = await fetch('/api/v1/files/folder', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ parent_path: currentPath.value, folder_name: newFolderName.value.trim() })
    })
    if (res.ok) {
      systemStore.showNotification('Sukses', `Folder "${newFolderName.value}" berhasil dibuat`, 'success')
      showNewFolderModal.value = false
      newFolderName.value = ''
      fetchFiles()
    }
  } catch (e) {
    systemStore.showNotification('Error', 'Gagal membuat folder baru', 'error')
  }
}

// Template Konten Awal Berkas Baru
function getInitialTemplateContent(ext) {
  switch (ext) {
    case 'html':
      return `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Dokumen Baru</title>
</head>
<body>
  <h1>Selamat Datang di RAZOS</h1>
</body>
</html>`
    case 'js':
      return `/**
 * RAZOS Script Module
 */
console.log('Hello from RAZ Studio Code!');`
    case 'css':
      return `/* RAZOS Stylesheet */
body {
  margin: 0;
  font-family: system-ui, -apple-system, sans-serif;
}`
    case 'json':
      return `{\n  "name": "razos-project",\n  "version": "1.0.0",\n  "description": ""\n}`
    case 'py':
      return `#!/usr/bin/env python3\n# RAZOS Python Script\n\nprint("Hello from Python!")`
    case 'go':
      return `package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello from Go!")\n}`
    case 'md':
      return `# Judul Berkas Baru\n\nTuliskan catatan atau dokumentasi di sini.`
    case 'yaml':
      return `version: '3.8'\nservices:\n  app:\n    image: nginx:alpine\n    ports:\n      - "80:80"`
    default:
      return ''
  }
}

// Buat Berkas Baru (Teks / Kode) & Langsung Buka di Editor (RBAC Guarded)
async function confirmCreateFile(openInEditor = true) {
  if (!authStore.canWriteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk membuat berkas.', 'warning')
    return
  }
  let name = newFileName.value.trim()
  if (!name) return

  // Jika nama tidak punya ekstensi, tambahkan ekstensi dari template yang dipilih
  if (!name.includes('.')) {
    name = `${name}.${newFileTemplate.value}`
  }

  const ext = name.split('.').pop().toLowerCase()
  const content = getInitialTemplateContent(ext)

  try {
    const res = await fetch('/api/v1/files/create-file', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        parent_path: currentPath.value,
        file_name: name,
        content: content
      })
    })
    const data = await res.json()
    if (res.ok && data.status === 'created') {
      systemStore.showNotification('Sukses', `Berkas "${name}" berhasil dibuat`, 'success')
      showNewFileModal.value = false
      newFileName.value = ''
      await fetchFiles()

      // Buka langsung di RAZ Studio Code (Sublime/Notepad++ Style)
      if (openInEditor && data.item) {
        windowStore.openTextEditor(data.item)
      }
    } else {
      systemStore.showNotification('Gagal', data.error || 'Tidak dapat membuat berkas', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', `Gagal membuat berkas: ${err.message}`, 'error')
  }
}

// Buka Berkas Langsung di RAZ Studio Code (Sublime Text / Notepad++ Style)
function openInTextEditor(item) {
  if (!item || item.is_dir) return
  windowStore.openTextEditor(item)
}

// Trigger Input Upload (RBAC Guarded)
function triggerUpload(type) {
  if (!authStore.canWriteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk mengunggah berkas atau folder.', 'warning')
    return
  }
  uploadType.value = type
  if (type === 'file' && fileUploadInputRef.value) {
    fileUploadInputRef.value.value = ''
    fileUploadInputRef.value.click()
  } else if (type === 'folder' && folderUploadInputRef.value) {
    folderUploadInputRef.value.value = ''
    folderUploadInputRef.value.click()
  } else if (type === 'zip' && zipUploadInputRef.value) {
    zipUploadInputRef.value.value = ''
    zipUploadInputRef.value.click()
  }
}

// Helper Upload Single File ke API raw-stream
async function uploadSingleFileRaw(file, relativePath = '', autoExtract = false) {
  const destDir = currentPath.value
  const res = await fetch('/api/v1/files/upload-raw', {
    method: 'POST',
    headers: {
      'x-dest-dir': encodeURIComponent(destDir),
      'x-filename': encodeURIComponent(file.name),
      'x-relative-path': relativePath ? encodeURIComponent(relativePath) : '',
      'x-auto-extract': autoExtract ? 'true' : 'false'
    },
    body: file
  })
  return await res.json()
}

// Handler Upload Berkas Biasa (Multi-file)
async function handleFileUploadChange(e) {
  const selectedFiles = Array.from(e.target.files || [])
  if (selectedFiles.length === 0) return

  uploadTotalCount.value = selectedFiles.length
  uploadCompletedCount.value = 0
  uploadProgress.value = 0
  isUploading.value = true
  showUploadModal.value = true

  for (let i = 0; i < selectedFiles.length; i++) {
    const file = selectedFiles[i]
    uploadCurrentFileName.value = file.name
    try {
      await uploadSingleFileRaw(file, '', false)
      uploadCompletedCount.value++
      uploadProgress.value = Math.round((uploadCompletedCount.value / uploadTotalCount.value) * 100)
    } catch (err) {
      console.error('Gagal upload berkas:', file.name, err)
    }
  }

  isUploading.value = false
  systemStore.showNotification('Upload Selesai', `${uploadCompletedCount.value} berkas berhasil diunggah`, 'success')
  fetchFiles()
  setTimeout(() => { showUploadModal.value = false }, 1800)
}

// Handler Upload Seluruh Folder (Recursive WebkitDirectory)
async function handleFolderUploadChange(e) {
  const selectedFiles = Array.from(e.target.files || [])
  if (selectedFiles.length === 0) return

  uploadTotalCount.value = selectedFiles.length
  uploadCompletedCount.value = 0
  uploadProgress.value = 0
  isUploading.value = true
  showUploadModal.value = true

  for (let i = 0; i < selectedFiles.length; i++) {
    const file = selectedFiles[i]
    uploadCurrentFileName.value = file.webkitRelativePath || file.name
    try {
      await uploadSingleFileRaw(file, file.webkitRelativePath || '', false)
      uploadCompletedCount.value++
      uploadProgress.value = Math.round((uploadCompletedCount.value / uploadTotalCount.value) * 100)
    } catch (err) {
      console.error('Gagal upload folder item:', file.name, err)
    }
  }

  isUploading.value = false
  systemStore.showNotification('Folder Diunggah', `Folder beserta ${uploadCompletedCount.value} berkas berhasil diunggah`, 'success')
  fetchFiles()
  setTimeout(() => { showUploadModal.value = false }, 1800)
}

// Handler Upload ZIP dengan Auto-Extract Otomatis
async function handleZipUploadChange(e) {
  const selectedFiles = Array.from(e.target.files || [])
  if (selectedFiles.length === 0) return

  const zipFile = selectedFiles[0]
  uploadTotalCount.value = 1
  uploadCompletedCount.value = 0
  uploadProgress.value = 30
  uploadCurrentFileName.value = zipFile.name
  isUploading.value = true
  showUploadModal.value = true

  try {
    const res = await uploadSingleFileRaw(zipFile, '', true)
    uploadProgress.value = 100
    uploadCompletedCount.value = 1
    systemStore.showNotification('Ekstraksi Sukses', res.message || `Arsip "${zipFile.name}" berhasil diekstrak`, 'success')
  } catch (err) {
    systemStore.showNotification('Gagal', `Gagal mengekstrak ZIP: ${err.message}`, 'error')
  } finally {
    isUploading.value = false
    fetchFiles()
    setTimeout(() => { showUploadModal.value = false }, 1800)
  }
}

// Buka Modal Konfirmasi Hapus (Single atau Multi-Item)
function openDeleteConfirm(target) {
  if (!authStore.canDeleteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk menghapus berkas atau folder.', 'warning')
    return
  }
  if (!target) return
  itemsToDelete.value = Array.isArray(target) ? [...target] : [target]
  if (itemsToDelete.value.length === 0) return
  showDeleteConfirmModal.value = true
}

// Kompatibilitas mundur: mengarahkan deleteItem langsung ke modal konfirmasi custom
function deleteItem(item) {
  openDeleteConfirm(item)
}

// Eksekusi Hapus yang Telah Dikonfirmasi Pengguna (RBAC Guarded & Batch Aware)
async function executeDeleteConfirmed() {
  if (!authStore.canDeleteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk menghapus berkas atau folder.', 'warning')
    showDeleteConfirmModal.value = false
    return
  }

  isDeletingFiles.value = true
  let successCount = 0
  let failCount = 0

  for (const item of itemsToDelete.value) {
    try {
      const res = await fetch(`/api/v1/files/item?path=${encodeURIComponent(item.path)}`, { method: 'DELETE' })
      if (res.ok) {
        successCount++
      } else {
        failCount++
      }
    } catch (e) {
      failCount++
    }
  }

  isDeletingFiles.value = false
  showDeleteConfirmModal.value = false
  itemsToDelete.value = []
  clearSelection()
  await fetchFiles(currentPath.value)

  if (successCount > 0 && failCount === 0) {
    systemStore.showNotification('Hapus Sukses', `${successCount} berkas/folder berhasil dihapus secara permanen`, 'success')
  } else if (successCount > 0 && failCount > 0) {
    systemStore.showNotification('Hapus Sebagian', `${successCount} berhasil dihapus, ${failCount} gagal dihapus`, 'warning')
  } else {
    systemStore.showNotification('Gagal', 'Terjadi kendala saat menghapus berkas terpilih', 'error')
  }
}

// Rename Berkas (RBAC Guarded)
function openRenameDialog(item) {
  if (!authStore.canWriteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk mengubah nama berkas/folder.', 'warning')
    return
  }
  renameTargetOldPath.value = item.path
  renameTargetNewName.value = item.name
  showRenameModal.value = true
}

async function confirmRename() {
  if (!authStore.canWriteStorage()) {
    systemStore.showNotification('Akses Dibatasi', 'Peran akun Anda tidak memiliki izin untuk mengubah nama berkas/folder.', 'warning')
    return
  }
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

// Unduh Berkas Langsung
function downloadFile(item) {
  if (!item || item.is_dir) return
  const a = document.createElement('a')
  a.href = `/api/v1/files/raw?path=${encodeURIComponent(item.path)}`
  a.download = item.name
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

// Salin Path Berkas / Folder ke Clipboard
async function copyItemPath(item) {
  if (!item) return
  try {
    await navigator.clipboard.writeText(item.path)
    systemStore.showNotification('Disalin', `Path "${item.path}" disalin ke clipboard`, 'success')
  } catch (e) {
    systemStore.showNotification('Gagal', 'Tidak dapat menyalin path ke clipboard', 'error')
  }
}

// Context Menu Klik Kanan pada Berkas atau Folder (Mendukung Seleksi Tunggal & Multi-Berkas)
function handleFileItemContextMenu(e, item) {
  // Jika item yang diklik kanan belum termasuk di dalam pilihan yang sedang aktif,
  // maka set pilihan hanya menjadi item ini (seperti Windows Explorer)
  if (!isItemSelected(item)) {
    selectedItems.value = [item]
    selectedItem.value = item
    lastSelectedIndex.value = visibleFiles.value.findIndex(f => f.path === item.path)
  }

  // SKENARIO 1: JIKA MULTI-BERKAS DIPILIH (Lebih dari 1 item)
  if (selectedItems.value.length > 1) {
    const multiMenu = [
      {
        label: `Buka ${selectedItems.value.length} Berkas`,
        icon: Eye,
        action: () => {
          selectedItems.value.filter(it => !it.is_dir).slice(0, 5).forEach(it => handleItemDoubleClick(it))
        }
      },
      {
        label: `Salin Seluruh Path (${selectedItems.value.length} Item)`,
        icon: Copy,
        action: async () => {
          const allPaths = selectedItems.value.map(it => it.path).join('\n')
          try {
            await navigator.clipboard.writeText(allPaths)
            systemStore.showNotification('Tersalin', `${selectedItems.value.length} path berkas disalin ke clipboard`, 'success')
          } catch (err) {}
        }
      },
      {
        label: 'Beri Bintang Semua',
        icon: Star,
        action: () => {
          selectedItems.value.forEach(it => {
            if (!isItemStarred(it.path)) toggleStar(it)
          })
          systemStore.showNotification('Favorit', `${selectedItems.value.length} berkas ditambahkan ke favorit`, 'success')
        }
      }
    ]

    // Menu Hapus Multi-Item (RBAC Guarded)
    if (authStore.canDeleteStorage()) {
      multiMenu.push(
        { separator: true },
        {
          label: `Hapus ${selectedItems.value.length} Berkas Terpilih...`,
          icon: Trash2,
          danger: true,
          action: () => openDeleteConfirm(selectedItems.value)
        }
      )
    }

    contextMenuStore.show(e, multiMenu, `${selectedItems.value.length} Item Dipilih`)
    return
  }

  // SKENARIO 2: JIKA BERKAS ATAU FOLDER TUNGGAL
  selectedItem.value = item

  // ==============================================================
  // A. CONTEXT MENU KHUSUS FOLDER (DIREKTORI)
  // ==============================================================
  if (item.is_dir) {
    const folderMenu = [
      {
        label: 'Buka Folder',
        icon: Folder,
        action: () => handleItemDoubleClick(item)
      },
      {
        label: 'Salin Path Direktori',
        icon: Copy,
        action: () => copyItemPath(item)
      },
      {
        label: isItemStarred(item.path) ? 'Hapus dari Bintang' : 'Beri Bintang',
        icon: Star,
        action: () => toggleStar(item)
      }
    ]

    // Ubah nama folder (RBAC Guarded)
    if (authStore.canWriteStorage()) {
      folderMenu.push({
        label: 'Ubah Nama Folder',
        icon: Edit3,
        action: () => openRenameDialog(item)
      })
    }

    folderMenu.push({
      label: 'Info & Properti Folder',
      icon: Info,
      action: () => openInfoModal(item)
    })

    // Hapus folder (RBAC Guarded)
    if (authStore.canDeleteStorage()) {
      folderMenu.push(
        { separator: true },
        {
          label: 'Hapus Folder...',
          icon: Trash2,
          danger: true,
          action: () => openDeleteConfirm(item)
        }
      )
    }

    contextMenuStore.show(e, folderMenu, `Folder: ${item.name}`)
    return
  }

  // ==============================================================
  // B. CONTEXT MENU KHUSUS BERKAS (FILE)
  // ==============================================================
  const fileMenu = [
    {
      label: 'Buka Berkas',
      icon: isVideoFile(item.extension) || isAudioFile(item.extension) ? Play : Eye,
      action: () => handleItemDoubleClick(item)
    }
  ]

  // Jika berupa berkas teks, kode, config, atau markdown, tambahkan opsi buka di Editor Teks
  if (isCodeFile(item.extension) || isDocFile(item.extension) || isMarkdownFile(item.extension) || !item.extension) {
    fileMenu.push({
      label: 'Edit di Editor Teks (Code)',
      icon: Code2,
      action: () => openInTextEditor(item)
    })
  }

  fileMenu.push(
    {
      label: 'Unduh Berkas (Download)',
      icon: Download,
      action: () => downloadFile(item)
    },
    {
      label: 'Salin Path Berkas',
      icon: Copy,
      action: () => copyItemPath(item)
    },
    {
      label: isItemStarred(item.path) ? 'Hapus dari Bintang' : 'Beri Bintang',
      icon: Star,
      action: () => toggleStar(item)
    }
  )

  // Ubah nama berkas (RBAC Guarded)
  if (authStore.canWriteStorage()) {
    fileMenu.push({
      label: 'Ubah Nama Berkas',
      icon: Edit3,
      action: () => openRenameDialog(item)
    })
  }

  fileMenu.push({
    label: 'Info & Properti Berkas',
    icon: Info,
    action: () => openInfoModal(item)
  })

  // Hapus berkas (RBAC Guarded)
  if (authStore.canDeleteStorage()) {
    fileMenu.push(
      { separator: true },
      {
        label: 'Hapus Berkas...',
        icon: Trash2,
        danger: true,
        action: () => openDeleteConfirm(item)
      }
    )
  }

  contextMenuStore.show(e, fileMenu, item.name)
}

// Context Menu Klik Kanan pada Area Kosong Explorer (Hierarki Submenu ala Windows 11)
function handleCanvasContextMenu(e) {
  const canvasMenu = []

  // 1. Menu Sort by (Urutkan Berkas ala Windows 11 Explorer)
  canvasMenu.push(
    {
      label: 'Urutkan (Sort by)',
      icon: ArrowUpDown,
      children: [
        {
          label: 'Nama (Name)',
          checked: sortBy.value === 'name',
          action: () => setSort('name')
        },
        {
          label: 'Tanggal Diubah (Date modified)',
          checked: sortBy.value === 'date',
          action: () => setSort('date')
        },
        {
          label: 'Format (Type)',
          checked: sortBy.value === 'type',
          action: () => setSort('type')
        },
        {
          label: 'Ukuran (Size)',
          checked: sortBy.value === 'size',
          action: () => setSort('size')
        },
        { separator: true },
        {
          label: 'Naik (Ascending)',
          checked: sortOrder.value === 'asc',
          action: () => setOrder('asc')
        },
        {
          label: 'Turun (Descending)',
          checked: sortOrder.value === 'desc',
          action: () => setOrder('desc')
        }
      ]
    },
    {
      label: 'Tampilan (View)',
      icon: LayoutGrid,
      children: [
        {
          label: 'Ikon Besar (Grid)',
          icon: LayoutGrid,
          checked: viewMode.value === 'grid',
          action: () => { viewMode.value = 'grid' }
        },
        {
          label: 'Rincian (List)',
          icon: List,
          checked: viewMode.value === 'list',
          action: () => { viewMode.value = 'list' }
        }
      ]
    },
    { separator: true }
  )

  // Tambahkan opsi Folder, File, dan Upload hanya jika akun memiliki izin tulis
  if (authStore.canWriteStorage()) {
    canvasMenu.push(
      {
        label: 'Folder Baru',
        icon: FolderPlus,
        action: () => { showNewFolderModal.value = true }
      },
      {
        label: 'File Baru (Teks/Kode)',
        icon: FilePlus,
        action: () => { showNewFileModal.value = true }
      },
      { separator: true },
      {
        label: 'Unggah (Upload)',
        icon: UploadCloud,
        children: [
          {
            label: 'Upload File',
            icon: FileUp,
            action: () => triggerUpload('file')
          },
          {
            label: 'Upload Folder Lengkap',
            icon: FolderUp,
            action: () => triggerUpload('folder')
          },
          {
            label: 'Upload & Ekstrak ZIP',
            icon: Archive,
            action: () => triggerUpload('zip')
          }
        ]
      },
      { separator: true }
    )
  }

  // Pintasan Seleksi Berkas
  canvasMenu.push(
    {
      label: 'Pilih Semua Berkas (Ctrl+A)',
      icon: CheckCircle2,
      action: () => selectAllFiles()
    }
  )

  if (selectedItems.value.length > 0) {
    canvasMenu.push(
      {
        label: `Batalkan Pilihan (${selectedItems.value.length})`,
        icon: X,
        action: () => clearSelection()
      }
    )
    if (authStore.canDeleteStorage()) {
      canvasMenu.push(
        {
          label: `Hapus ${selectedItems.value.length} Berkas Terpilih...`,
          icon: Trash2,
          danger: true,
          action: () => openDeleteConfirm(selectedItems.value)
        }
      )
    }
  }

  canvasMenu.push(
    { separator: true },
    {
      label: 'Sinkronisasi Cloud',
      icon: Cloud,
      action: () => windowStore.openCloudSync()
    },
    {
      label: 'Muat Ulang (Refresh)',
      icon: RefreshCw,
      action: () => fetchFiles(currentPath.value)
    },
    {
      label: 'Salin Path Direktori',
      icon: Copy,
      action: () => copyCurrentPath()
    }
  )

  contextMenuStore.show(e, canvasMenu, 'File Manager')
}

onMounted(async () => {
  await fetchDevices()
  // Tentukan path awal berbasis izin RBAC & direktori home pengguna
  const userHome = authStore.getUserHomeDir()
  const targetPath = userHome || defaultStoragePath.value || '/DATA'
  currentPath.value = targetPath
  addressInputValue.value = targetPath
  activeSidebarSection.value = 'drive'
  await fetchFiles(targetPath)
})
</script>

<template>
  <div
    class="h-full flex overflow-hidden text-xs sm:text-sm select-none bg-slate-950 text-slate-100 rounded-xl outline-none"
    @click="closeAllNewMenus"
    tabindex="0"
    @keydown="handleGlobalKeydown"
  >
    <!-- LEFT SIDEBAR (Google Drive Style) -->
    <div class="w-56 border-r border-slate-800/80 bg-slate-900/60 p-3 flex flex-col justify-between hidden md:flex shrink-0">
      <div class="flex flex-col gap-4">
        <!-- New Button with Dropdown (RBAC Guarded) -->
        <div v-if="!authStore.canWriteStorage()" class="w-full h-10 px-3.5 rounded-2xl bg-slate-800/60 border border-slate-700/50 text-slate-400 font-medium flex items-center justify-between text-xs opacity-80 cursor-not-allowed shadow-inner" title="Peran akun ini hanya memiliki izin baca (Read-Only)">
          <div class="flex items-center gap-2">
            <Lock class="w-4 h-4 text-amber-400" />
            <span class="tracking-wide">Hanya Baca</span>
          </div>
          <span class="text-[10px] uppercase font-mono px-1.5 py-0.5 rounded bg-slate-700/60 text-slate-300">RO</span>
        </div>
        <div v-else class="relative">
          <button
            @click.stop="showNewMenu = !showNewMenu"
            class="w-full h-10 px-4 rounded-2xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center justify-between gap-2 shadow-lg transition active:scale-95"
          >
            <div class="flex items-center gap-2">
              <Plus class="w-5 h-5" />
              <span class="text-xs tracking-wide">Baru</span>
            </div>
            <ChevronDown class="w-4 h-4 transition-transform duration-200" :class="showNewMenu ? 'rotate-180' : ''" />
          </button>

          <!-- Dropdown Menu -->
          <div
            v-if="showNewMenu"
            class="absolute top-12 left-0 w-64 bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl z-50 p-1.5 flex flex-col gap-1 text-xs font-sans text-slate-200 backdrop-blur-md animate-in fade-in zoom-in-95"
          >
            <button
              @click="() => { showNewMenu = false; showNewFolderModal = true; }"
              class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
            >
              <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
                <FolderPlus class="w-4 h-4" />
              </div>
              <div>
                <span class="font-medium block text-slate-200">Folder Baru</span>
                <span class="text-[10px] text-slate-400 block">Buat folder di direktori aktif</span>
              </div>
            </button>

            <button
              @click="() => { showNewMenu = false; showNewFileModal = true; }"
              class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
            >
              <div class="w-7 h-7 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0">
                <FilePlus class="w-4 h-4" />
              </div>
              <div>
                <span class="font-medium block text-slate-200">File Baru (Teks/Kode)</span>
                <span class="text-[10px] text-slate-400 block">Buat & langsung buka di Editor</span>
              </div>
            </button>

            <div class="h-px bg-slate-800 my-0.5"></div>

            <button
              @click="() => { showNewMenu = false; triggerUpload('file'); }"
              class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
            >
              <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
                <FileUp class="w-4 h-4" />
              </div>
              <div>
                <span class="font-medium block text-slate-200">Upload File</span>
                <span class="text-[10px] text-slate-400 block">Unggah satu atau banyak berkas</span>
              </div>
            </button>

            <button
              @click="() => { showNewMenu = false; triggerUpload('folder'); }"
              class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
            >
              <div class="w-7 h-7 rounded-lg bg-amber-500/20 text-amber-400 flex items-center justify-center shrink-0">
                <FolderUp class="w-4 h-4" />
              </div>
              <div>
                <span class="font-medium block text-slate-200">Upload Folder</span>
                <span class="text-[10px] text-slate-400 block">Unggah folder beserta seluruh isinya</span>
              </div>
            </button>

            <button
              @click="() => { showNewMenu = false; triggerUpload('zip'); }"
              class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
            >
              <div class="w-7 h-7 rounded-lg bg-purple-500/20 text-purple-400 flex items-center justify-center shrink-0">
                <Archive class="w-4 h-4" />
              </div>
              <div>
                <span class="font-medium block text-purple-300">Upload ZIP (Auto Extract)</span>
                <span class="text-[10px] text-slate-400 block">Unggah & langsung diekstrak</span>
              </div>
            </button>
          </div>
        </div>

        <!-- Navigation Links -->
        <div class="flex flex-col gap-1 text-xs">
          <button
            @click="() => { activeSidebarSection = 'drive'; fetchFiles(defaultStoragePath); }"
            class="px-3 py-2.5 rounded-xl flex items-center gap-3 transition font-medium text-left"
            :class="activeSidebarSection === 'drive' ? 'bg-sky-500/20 text-sky-400 font-bold' : 'text-slate-300 hover:bg-slate-800/60'"
          >
            <Folder class="w-4 h-4 text-sky-400 fill-sky-400/20" />
            <span>My Storage (/DATA)</span>
          </button>

          <button
            @click="activeSidebarSection = 'starred'"
            class="px-3 py-2.5 rounded-xl flex items-center gap-3 transition font-medium text-left"
            :class="activeSidebarSection === 'starred' ? 'bg-sky-500/20 text-sky-400 font-bold' : 'text-slate-300 hover:bg-slate-800/60'"
          >
            <Star class="w-4 h-4 text-amber-400" />
            <div class="flex items-center justify-between flex-1">
              <span>Berbintang (Starred)</span>
              <span v-if="starredItems.length > 0" class="text-[10px] font-mono px-1.5 py-0.5 rounded-full bg-amber-500/20 text-amber-300 font-bold">
                {{ starredItems.length }}
              </span>
            </div>
          </button>
        </div>

        <!-- Devices Section (Physical Disks) -->
        <div class="flex flex-col gap-2 pt-2 border-t border-slate-800/80">
          <span class="text-[10px] font-bold text-slate-500 tracking-wider uppercase px-2">Storage Disks</span>
          <div class="flex flex-col gap-1.5">
            <div
              v-for="(dev, idx) in physicalDisks"
              :key="'phys-' + idx"
              @click="() => { activeSidebarSection = 'devices'; fetchFiles(dev.mountpoint); }"
              class="p-2 rounded-xl bg-slate-950/60 hover:bg-slate-800/60 border border-slate-800/60 cursor-pointer transition flex flex-col gap-1"
              :class="currentPath.startsWith(dev.mountpoint) && activeSidebarSection === 'devices' ? 'ring-1 ring-sky-400/40 bg-sky-950/20' : ''"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5 truncate">
                  <component :is="dev.is_usb ? Usb : HardDrive" class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                  <span class="text-[11px] font-medium text-slate-200 truncate">{{ dev.label }}</span>
                </div>
                <span class="text-[10px] text-slate-400 font-mono">{{ (dev.percent ?? 0).toFixed(0) }}%</span>
              </div>
              <div class="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                <div class="bg-sky-400 h-full" :style="{ width: `${dev.percent ?? 0}%` }"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Cloud Storage Section (Nextcloud / WebDAV) -->
        <div class="flex flex-col gap-2 pt-2 border-t border-slate-800/80">
          <div class="flex items-center justify-between px-2">
            <span class="text-[10px] font-bold text-sky-400/90 tracking-wider uppercase">Cloud Storage</span>
            <button
              @click="windowStore.openCloudSync()"
              class="text-[10px] text-sky-400 hover:text-sky-300 transition hover:underline flex items-center gap-1"
              title="Kelola & Sinkronisasi Cloud"
            >
              <Cloud class="w-3 h-3" />
              <span>Kelola</span>
            </button>
          </div>
          <div class="flex flex-col gap-1.5">
            <div
              v-for="(dev, idx) in cloudDisks"
              :key="'cloud-' + idx"
              @click="() => { activeSidebarSection = 'cloud'; fetchFiles(dev.mountpoint); }"
              class="p-2 rounded-xl bg-sky-950/20 hover:bg-sky-900/30 border border-sky-800/40 cursor-pointer transition flex flex-col gap-1"
              :class="currentPath.startsWith(dev.mountpoint) ? 'ring-1 ring-sky-400/60 bg-sky-950/40' : ''"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5 truncate">
                  <Cloud class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                  <span class="text-[11px] font-medium text-sky-200 truncate">{{ dev.label }}</span>
                </div>
                <span class="text-[10px] text-sky-300 font-mono">{{ (dev.percent ?? 0).toFixed(0) }}%</span>
              </div>
              <div class="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                <div class="bg-gradient-to-r from-sky-500 to-indigo-500 h-full" :style="{ width: `${dev.percent ?? 0}%` }"></div>
              </div>
            </div>

            <!-- Tombol Tambah/Hubungkan Cloud jika belum ada cloud terpasang -->
            <button
              v-if="cloudDisks.length === 0"
              @click="windowStore.openCloudSync()"
              class="p-2 rounded-xl border border-dashed border-slate-700/80 hover:border-sky-500/60 bg-slate-900/40 hover:bg-sky-950/20 text-slate-400 hover:text-sky-300 text-left transition flex items-center gap-2 group"
            >
              <div class="w-6 h-6 rounded-lg bg-sky-500/10 group-hover:bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
                <Plus class="w-3.5 h-3.5" />
              </div>
              <span class="text-[11px]">Hubungkan Cloud...</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Bottom Cloud Storage Meter (100% Real Capacity dari Hardware Host) -->
      <div class="p-3 rounded-2xl bg-slate-950/70 border border-slate-800/80 flex flex-col gap-2">
        <div class="flex items-center justify-between text-[11px]">
          <span class="text-slate-400">Total Digunakan</span>
          <span class="text-sky-400 font-bold font-mono">{{ storageUsagePercent }}%</span>
        </div>
        <div class="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden">
          <div class="bg-sky-400 h-full transition-all duration-500" :style="{ width: `${Math.min(100, Math.max(0, storageUsagePercent))}%` }"></div>
        </div>
        <span class="text-[10px] text-slate-500 font-mono">
          {{ formatBytes(usedStorageBytes) }} dari {{ formatBytes(totalStorageBytes) }} terpakai
        </span>
      </div>
    </div>

    <!-- MAIN DRIVE CONTENT AREA -->
    <div class="flex-1 flex flex-col overflow-hidden bg-slate-950">
      <!-- WINDOWS FILE EXPLORER TOP NAVIGATION & ADDRESS BAR -->
      <div class="h-12 border-b border-slate-800/80 px-3 flex items-center gap-2 bg-slate-900/60 shrink-0">
        <!-- Navigation Trio (Back, Forward, Up) -->
        <div class="flex items-center gap-0.5 shrink-0">
          <button
            @click="goBack"
            :disabled="!canGoBack"
            class="p-1.5 rounded-lg transition"
            :class="canGoBack ? 'hover:bg-slate-800 text-slate-200 cursor-pointer active:scale-95' : 'text-slate-600 opacity-40 cursor-not-allowed'"
            title="Kembali (Back)"
          >
            <ArrowLeft class="w-4 h-4" />
          </button>
          <button
            @click="goForward"
            :disabled="!canGoForward"
            class="p-1.5 rounded-lg transition"
            :class="canGoForward ? 'hover:bg-slate-800 text-slate-200 cursor-pointer active:scale-95' : 'text-slate-600 opacity-40 cursor-not-allowed'"
            title="Maju (Forward)"
          >
            <ArrowRight class="w-4 h-4" />
          </button>
          <button
            @click="navigateUp"
            :disabled="!canNavigateUp"
            class="p-1.5 rounded-lg transition"
            :class="canNavigateUp ? 'hover:bg-slate-800 text-slate-200 cursor-pointer active:scale-95' : 'text-slate-600 opacity-40 cursor-not-allowed'"
            title="Naik Satu Tingkat (Up)"
          >
            <ArrowUp class="w-4 h-4" />
          </button>
        </div>

        <!-- Address Bar (Windows Explorer Style) -->
        <div
          class="flex-1 min-w-0 h-9 bg-slate-950/80 border border-slate-800/90 hover:border-slate-700/80 focus-within:border-sky-500/80 focus-within:ring-1 focus-within:ring-sky-500/20 rounded-xl px-2.5 flex items-center gap-2 transition group text-xs font-mono shadow-inner"
        >
          <!-- Current Location Icon -->
          <Star v-if="activeSidebarSection === 'starred'" class="w-3.5 h-3.5 text-amber-400 shrink-0 fill-amber-400" />
          <HardDrive v-else-if="currentPath.includes('Disk') || /^[a-zA-Z]:/.test(currentPath)" class="w-3.5 h-3.5 text-slate-400 shrink-0" />
          <Folder v-else class="w-3.5 h-3.5 text-sky-400 shrink-0" />

          <!-- Editable Input Mode -->
          <div v-if="isEditingAddress" class="flex-1 h-full flex items-center">
            <input
              ref="addressInputRef"
              v-model="addressInputValue"
              @keydown.enter="submitAddress"
              @keydown.esc="cancelEditingAddress"
              @blur="onAddressInputBlur"
              type="text"
              class="w-full bg-transparent border-0 outline-none text-sky-200 text-xs font-mono h-full px-0 select-all"
              placeholder="Ketik path direktori (misal: /DATA atau C:/DATA)..."
            />
          </div>

          <!-- Breadcrumbs Interactive Mode -->
          <div
            v-else
            @click="startEditingAddress"
            class="flex-1 h-full flex items-center gap-1 overflow-x-auto cursor-text select-none py-0.5"
            style="scrollbar-width: none;"
            title="Klik untuk mengedit address bar"
          >
            <template v-for="(crumb, idx) in breadcrumbs" :key="idx">
              <button
                @click.stop="fetchFiles(crumb.path)"
                class="px-1.5 py-0.5 rounded hover:bg-slate-800/90 text-slate-300 hover:text-sky-300 transition flex items-center gap-1 shrink-0 font-medium"
                :class="idx === breadcrumbs.length - 1 ? 'text-sky-400 font-bold' : ''"
              >
                <span>{{ crumb.name }}</span>
              </button>
              <ChevronRight v-if="idx < breadcrumbs.length - 1" class="w-3 h-3 text-slate-600 shrink-0" />
            </template>
            <!-- Area kosong untuk memudahkan klik edit address -->
            <div class="flex-1 min-w-[24px] h-full cursor-text"></div>
          </div>

          <!-- Right Action Controls in Address Bar -->
          <div class="flex items-center gap-0.5 shrink-0">
            <!-- Go Button when editing -->
            <button
              v-if="isEditingAddress"
              @mousedown.prevent="submitAddress"
              class="p-1 rounded-lg bg-sky-500/20 text-sky-400 hover:bg-sky-500/30 transition"
              title="Pergi ke Path"
            >
              <ArrowRight class="w-3.5 h-3.5" />
            </button>

            <!-- Copy Path Button -->
            <button
              v-if="!isEditingAddress"
              @click.stop="copyCurrentPath"
              class="p-1 rounded-lg hover:bg-slate-800 text-slate-500 hover:text-slate-200 transition"
              :class="pathCopied ? 'text-emerald-400' : ''"
              title="Salin Path Direktori"
            >
              <Check v-if="pathCopied" class="w-3.5 h-3.5 text-emerald-400" />
              <Copy v-else class="w-3.5 h-3.5" />
            </button>

            <!-- Refresh Button -->
            <button
              @click.stop="fetchFiles(currentPath)"
              class="p-1 rounded-lg hover:bg-slate-800 text-slate-500 hover:text-slate-200 transition"
              title="Muat Ulang Folder"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="isLoading ? 'animate-spin text-sky-400' : ''" />
            </button>
          </div>
        </div>

        <!-- Search Box (Windows Explorer Search Style) -->
        <div class="w-40 sm:w-56 h-9 relative shrink-0">
          <Search class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="`Cari di ${currentFolderShortName}...`"
            class="w-full h-full pl-8 pr-7 rounded-xl bg-slate-950/80 border border-slate-800/90 text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-400 focus:ring-1 focus:ring-sky-400/20 transition"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300"
          >
            <X class="w-3 h-3" />
          </button>
        </div>
      </div>

      <!-- FILTER & VIEW MODE RIBBON BAR -->
      <div class="h-10 border-b border-slate-800/60 px-3 sm:px-4 flex items-center justify-between gap-2 sm:gap-3 bg-slate-900/30 text-xs font-mono shrink-0">
        <!-- Quick Action "+ Baru" (Always Visible on Mobile / Small Windows) & Filter Tabs -->
        <div class="flex items-center gap-2 overflow-x-auto py-0.5" style="scrollbar-width: none;">
          <!-- Ribbon "+ Baru" Button with Dropdown (RBAC Guarded) -->
          <div v-if="!authStore.canWriteStorage()" class="flex items-center gap-1.5 px-2.5 py-1 bg-slate-800/60 border border-slate-700/50 text-slate-400 font-medium rounded-lg text-xs opacity-80 cursor-not-allowed shrink-0" title="Akun Anda hanya memiliki izin baca (Read-Only)">
            <Lock class="w-3.5 h-3.5 text-amber-400" />
            <span class="hidden sm:inline">Read Only</span>
          </div>
          <div v-else class="relative shrink-0">
            <button
              @click.stop="showRibbonNewMenu = !showRibbonNewMenu"
              class="flex items-center gap-1.5 px-2.5 py-1 bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold rounded-lg text-xs shadow-sm transition active:scale-95"
              title="Buat folder, berkas teks/kode, atau unggah data"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Baru</span>
              <ChevronDown class="w-3 h-3 ml-0.5 opacity-80" />
            </button>

            <!-- Dropdown Menu -->
            <div
              v-if="showRibbonNewMenu"
              class="absolute top-9 left-0 w-64 bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl z-50 p-1.5 flex flex-col gap-1 text-xs font-sans text-slate-200 backdrop-blur-xl animate-in fade-in zoom-in-95"
            >
              <button
                @click="() => { showRibbonNewMenu = false; showNewFolderModal = true; }"
                class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
              >
                <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
                  <FolderPlus class="w-4 h-4" />
                </div>
                <div>
                  <span class="font-medium block text-slate-200">Folder Baru</span>
                  <span class="text-[10px] text-slate-400 block">Buat folder di direktori aktif</span>
                </div>
              </button>

              <button
                @click="() => { showRibbonNewMenu = false; showNewFileModal = true; }"
                class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
              >
                <div class="w-7 h-7 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center shrink-0">
                  <FilePlus class="w-4 h-4" />
                </div>
                <div>
                  <span class="font-medium block text-slate-200">File Baru (Teks/Kode)</span>
                  <span class="text-[10px] text-slate-400 block">Buat & langsung buka di Editor</span>
                </div>
              </button>

              <div class="h-px bg-slate-800 my-0.5"></div>

              <button
                @click="() => { showRibbonNewMenu = false; triggerUpload('file'); }"
                class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
              >
                <div class="w-7 h-7 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center shrink-0">
                  <FileUp class="w-4 h-4" />
                </div>
                <div>
                  <span class="font-medium block text-slate-200">Upload File</span>
                  <span class="text-[10px] text-slate-400 block">Unggah satu atau banyak berkas</span>
                </div>
              </button>

              <button
                @click="() => { showRibbonNewMenu = false; triggerUpload('folder'); }"
                class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
              >
                <div class="w-7 h-7 rounded-lg bg-amber-500/20 text-amber-400 flex items-center justify-center shrink-0">
                  <FolderUp class="w-4 h-4" />
                </div>
                <div>
                  <span class="font-medium block text-slate-200">Upload Folder</span>
                  <span class="text-[10px] text-slate-400 block">Unggah folder beserta seluruh isinya</span>
                </div>
              </button>

              <button
                @click="() => { showRibbonNewMenu = false; triggerUpload('zip'); }"
                class="w-full px-3 py-2 rounded-xl hover:bg-slate-800 flex items-center gap-3 transition text-left"
              >
                <div class="w-7 h-7 rounded-lg bg-purple-500/20 text-purple-400 flex items-center justify-center shrink-0">
                  <Archive class="w-4 h-4" />
                </div>
                <div>
                  <span class="font-medium block text-purple-300">Upload ZIP (Auto Extract)</span>
                  <span class="text-[10px] text-slate-400 block">Unggah & langsung diekstrak</span>
                </div>
              </button>
            </div>
          </div>

          <div class="h-4 w-px bg-slate-800/80 shrink-0"></div>

          <!-- Filter Tabs Media Viewers -->
          <button
            v-for="f in [{ id: 'all', label: 'Semua' }, { id: 'folders', label: 'Folder' }, { id: 'images', label: 'Gambar' }, { id: 'video', label: 'Video' }, { id: 'audio', label: 'Audio' }, { id: 'docs', label: 'Dokumen' }, { id: 'code', label: 'Kode' }]"
            :key="f.id"
            @click="filterType = f.id"
            class="px-2.5 py-1 rounded-lg transition font-sans text-xs shrink-0"
            :class="filterType === f.id ? 'bg-sky-500/20 text-sky-400 font-bold border border-sky-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'"
          >
            {{ f.label }}
          </button>
        </div>

        <!-- Right Side: Cloud Sync Shortcut, Sort Dropdown & Grid/List View Mode -->
        <div class="flex items-center gap-2 sm:gap-2.5 shrink-0">
          <!-- Tombol Cloud Sync Cepat -->
          <button
            @click="windowStore.openCloudSync()"
            class="px-2 py-1 rounded-lg bg-sky-500/10 hover:bg-sky-500/20 text-sky-400 border border-sky-500/30 flex items-center gap-1.5 transition text-[11px] font-sans font-medium active:scale-95 shadow-sm"
            title="Kelola & Sinkronisasi Cloud Eksternal (WebDAV/Nextcloud)"
          >
            <Cloud class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Cloud Sync</span>
          </button>

          <!-- Tombol Dropdown Sort by (Windows 11 Explorer Standard) -->
          <div class="relative shrink-0">
            <button
              @click.stop="() => { showSortDropdown = !showSortDropdown; }"
              class="px-2 py-1 rounded-lg bg-slate-900/90 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800/80 flex items-center gap-1.5 transition text-[11px] font-sans font-medium active:scale-95 shadow-sm"
              title="Urutkan berkas (Sort by)"
            >
              <ArrowUpDown class="w-3.5 h-3.5 text-sky-400" />
              <span>Sort</span>
              <ChevronDown class="w-3 h-3 text-slate-500" />
            </button>

            <!-- Dropdown Menu Sort by (Solid Dark Windows 11) -->
            <div
              v-if="showSortDropdown"
              class="absolute top-9 right-0 w-48 bg-[#201f27] border border-white/10 rounded-xl shadow-2xl shadow-black/80 z-50 p-1.5 flex flex-col gap-0.5 text-xs font-sans text-slate-200 animate-in fade-in zoom-in-95"
            >
              <button
                @click="setSort('name')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortBy === 'name' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Name</span>
                <div v-if="sortBy === 'name'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>
              <button
                @click="setSort('date')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortBy === 'date' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Date modified</span>
                <div v-if="sortBy === 'date'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>
              <button
                @click="setSort('type')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortBy === 'type' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Type</span>
                <div v-if="sortBy === 'type'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>
              <button
                @click="setSort('size')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortBy === 'size' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Size</span>
                <div v-if="sortBy === 'size'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>

              <div class="h-[1px] bg-white/10 my-1 mx-1"></div>

              <button
                @click="setOrder('asc')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortOrder === 'asc' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Ascending</span>
                <div v-if="sortOrder === 'asc'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>
              <button
                @click="setOrder('desc')"
                class="w-full px-2.5 py-1.5 rounded-lg hover:bg-white/10 flex items-center justify-between transition text-left"
                :class="sortOrder === 'desc' ? 'bg-white/10 text-white font-medium' : 'text-slate-300'"
              >
                <span>Descending</span>
                <div v-if="sortOrder === 'desc'" class="w-2 h-2 rounded-full bg-sky-400"></div>
              </button>
            </div>
          </div>

          <span class="text-[11px] text-slate-500 hidden sm:inline">
            {{ filteredFiles.length }} item
          </span>

          <div class="flex items-center bg-slate-900/90 p-0.5 rounded-xl border border-slate-800/80">
            <button
              @click="viewMode = 'grid'"
              class="p-1 rounded-lg transition"
              :class="viewMode === 'grid' ? 'bg-slate-800 text-sky-400 shadow-sm' : 'text-slate-400 hover:text-slate-200'"
              title="Grid View"
            >
              <LayoutGrid class="w-3.5 h-3.5" />
            </button>
            <button
              @click="viewMode = 'list'"
              class="p-1 rounded-lg transition"
              :class="viewMode === 'list' ? 'bg-slate-800 text-sky-400 shadow-sm' : 'text-slate-400 hover:text-slate-200'"
              title="List View"
            >
              <List class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>

      <!-- MAIN EXPLORER AREA (GRID VS LIST) -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <div class="flex-1 p-4 overflow-y-auto" @scroll="handleExplorerScroll" @contextmenu.prevent="handleCanvasContextMenu">
          <!-- GRID VIEW (Fluid Responsive Auto-Fill Layout) -->
          <div v-if="viewMode === 'grid'" class="grid grid-cols-[repeat(auto-fill,minmax(140px,1fr))] sm:grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-3">
            <div
              v-for="(item, idx) in visibleFiles"
              :key="item.path || idx"
              @click="handleItemClick(item, $event)"
              @dblclick="handleItemDoubleClick(item)"
              @contextmenu.prevent.stop="handleFileItemContextMenu($event, item)"
              class="p-3 rounded-2xl border transition-all duration-200 cursor-pointer flex flex-col justify-between group relative select-none min-h-[140px]"
              :class="isItemSelected(item)
                ? 'bg-sky-500/25 border-sky-400 ring-2 ring-sky-400/50 shadow-lg text-sky-100'
                : 'bg-slate-900/60 border-slate-800/80 hover:border-slate-700 hover:bg-slate-900/90 hover:shadow-md'"
            >
              <!-- Card Top: Icon & Action Controls (Info & Star) -->
              <div class="flex items-start justify-between gap-2 mb-2">
                <!-- Media Thumbnail / Icon Box (Anti-Hang Concurrency Protected) -->
                <div v-if="isImageFile(item.extension)" class="w-11 h-11 rounded-xl overflow-hidden bg-slate-950 border border-slate-800 shrink-0 shadow-inner">
                  <RAZThumbnail
                    :path="item.path"
                    :name="item.name"
                    :size="item.size"
                    :extension="item.extension"
                  />
                </div>
                <div v-else class="w-11 h-11 rounded-xl flex items-center justify-center shrink-0 shadow-inner"
                     :class="item.is_dir ? 'bg-sky-500/15 text-sky-400 border border-sky-500/20' : isVideoFile(item.extension) ? 'bg-purple-500/15 text-purple-400 border border-purple-500/20' : isAudioFile(item.extension) ? 'bg-amber-500/15 text-amber-400 border border-amber-500/20' : isPdfFile(item.extension) ? 'bg-rose-500/15 text-rose-400 border border-rose-500/20' : isMarkdownFile(item.extension) ? 'bg-sky-500/15 text-sky-400 border border-sky-500/30' : isCodeFile(item.extension) ? 'bg-slate-800/80 text-emerald-400 border border-slate-700/60' : 'bg-slate-800/80 text-amber-300 border border-slate-700/60'">
                  <Folder v-if="item.is_dir" class="w-6 h-6 fill-sky-400/20 text-sky-400" />
                  <Film v-else-if="isVideoFile(item.extension)" class="w-5 h-5 text-purple-400" />
                  <Music v-else-if="isAudioFile(item.extension)" class="w-5 h-5 text-amber-400" />
                  <FileText v-else-if="isPdfFile(item.extension)" class="w-5 h-5 text-rose-400" />
                  <BookOpen v-else-if="isMarkdownFile(item.extension)" class="w-5 h-5 text-sky-400" />
                  <FileCode v-else-if="isCodeFile(item.extension)" class="w-5 h-5 text-emerald-400" />
                  <FileText v-else-if="isDocFile(item.extension)" class="w-5 h-5 text-amber-300" />
                  <File v-else class="w-5 h-5 text-slate-400" />
                </div>

                <!-- Right Controls: Info Modal & Star Favorite -->
                <div class="flex items-center gap-0.5">
                  <button
                    @click.stop="openInfoModal(item)"
                    class="p-1 rounded-lg hover:bg-slate-800 text-slate-500 hover:text-sky-400 transition opacity-0 group-hover:opacity-100"
                    title="Detail & Properti Berkas"
                  >
                    <Info class="w-3.5 h-3.5" />
                  </button>
                  <button
                    @click.stop="(e) => toggleStar(item, e)"
                    class="p-1 rounded-lg hover:bg-slate-800 text-slate-500 hover:text-amber-400 transition"
                    :class="isItemStarred(item.path) ? 'text-amber-400 opacity-100' : 'opacity-0 group-hover:opacity-100'"
                    title="Bintang / Favorit"
                  >
                    <Star class="w-3.5 h-3.5" :class="isItemStarred(item.path) ? 'fill-amber-400' : ''" />
                  </button>
                </div>
              </div>

              <!-- Card Bottom: Name & Meta Text -->
              <div class="min-w-0 flex-1 flex flex-col justify-end">
                <h5 class="font-medium text-slate-200 text-xs line-clamp-2 leading-snug group-hover:text-sky-300 transition" :title="item.name">
                  {{ item.name }}
                </h5>
                <div class="flex items-center justify-between text-[10px] text-slate-400 mt-2 font-mono pt-1.5 border-t border-slate-800/60">
                  <span>{{ item.is_dir ? 'Folder' : formatBytes(item.size) }}</span>
                  <span class="uppercase font-semibold text-slate-400 tracking-wider">{{ item.extension ? item.extension.replace('.', '') : '' }}</span>
                </div>
              </div>
            </div>

            <!-- Indikator Loading Batch Virtual Scrolling Grid -->
            <div v-if="hasMoreFiles" class="col-span-full py-4 flex items-center justify-center gap-2 text-xs text-slate-500 font-mono">
              <RefreshCw class="w-3.5 h-3.5 animate-spin text-sky-400" />
              <span>Menampilkan {{ visibleFiles.length }} dari {{ filteredFiles.length }} item (gulir ke bawah untuk memuat lagi)...</span>
            </div>
          </div>

          <!-- LIST VIEW (With Type Badges) -->
          <div v-else class="rounded-2xl border border-slate-800/80 overflow-hidden bg-slate-900/50">
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
                  v-for="(item, idx) in visibleFiles"
                  :key="item.path || idx"
                  @click="handleItemClick(item, $event)"
                  @dblclick="handleItemDoubleClick(item)"
                  @contextmenu.prevent.stop="handleFileItemContextMenu($event, item)"
                  class="hover:bg-slate-800/60 cursor-pointer transition select-none"
                  :class="isItemSelected(item) ? 'bg-sky-500/20 text-sky-200 font-medium' : ''"
                >
                  <td class="py-2.5 px-3 font-medium flex items-center gap-2.5">
                    <button @click.stop="(e) => toggleStar(item, e)" class="text-slate-500 hover:text-amber-400">
                      <Star class="w-3.5 h-3.5" :class="isItemStarred(item.path) ? 'text-amber-400 fill-amber-400' : ''" />
                    </button>
                    <Folder v-if="item.is_dir" class="w-4 h-4 text-sky-400 fill-sky-400/20" />
                    <Image v-else-if="isImageFile(item.extension)" class="w-4 h-4 text-sky-400" />
                    <Film v-else-if="isVideoFile(item.extension)" class="w-4 h-4 text-purple-400" />
                    <Music v-else-if="isAudioFile(item.extension)" class="w-4 h-4 text-amber-400" />
                    <FileText v-else-if="isPdfFile(item.extension)" class="w-4 h-4 text-rose-400" />
                    <BookOpen v-else-if="isMarkdownFile(item.extension)" class="w-4 h-4 text-sky-400" />
                    <FileCode v-else-if="isCodeFile(item.extension)" class="w-4 h-4 text-emerald-400" />
                    <FileText v-else-if="isDocFile(item.extension)" class="w-4 h-4 text-amber-300" />
                    <File v-else class="w-4 h-4 text-slate-400" />
                    <span class="truncate max-w-[280px]">{{ item.name }}</span>
                  </td>
                  <td class="py-2.5 px-3 font-mono text-slate-400">{{ item.is_dir ? '-' : formatBytes(item.size) }}</td>
                  <td class="py-2.5 px-3 uppercase text-[10px] text-slate-400">{{ item.is_dir ? 'Folder' : (item.extension || 'File') }}</td>
                  <td class="py-2.5 px-3 text-slate-400 text-[11px]">{{ formatDate(item.mod_time) }}</td>
                  <td class="py-2.5 px-3 text-right">
                    <div class="flex items-center justify-end gap-1.5" @click.stop>
                      <button @click="openInfoModal(item)" class="p-1 rounded hover:bg-slate-700 text-slate-400 hover:text-sky-400" title="Detail & Properti">
                        <Info class="w-3.5 h-3.5" />
                      </button>
                      <button v-if="!item.is_dir" @click="openInTextEditor(item)" class="p-1 rounded hover:bg-slate-700 text-emerald-400" title="Edit Berkas di Editor Teks">
                        <Code2 class="w-3.5 h-3.5" />
                      </button>
                      <button v-if="!item.is_dir" @click="handleItemDoubleClick(item)" class="p-1 rounded hover:bg-slate-700 text-sky-400" title="Buka Viewer / Player">
                        <Play v-if="isVideoFile(item.extension) || isAudioFile(item.extension)" class="w-3.5 h-3.5" />
                        <Eye v-else class="w-3.5 h-3.5" />
                      </button>
                      <button @click="openRenameDialog(item)" class="p-1 rounded hover:bg-slate-700 text-slate-300" title="Ubah Nama">
                        <Edit3 class="w-3.5 h-3.5" />
                      </button>
                      <button @click="openDeleteConfirm(item)" class="p-1 rounded hover:bg-rose-500/20 text-rose-400" title="Hapus Berkas">
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>

                <!-- Indikator Loading Batch Virtual Scrolling List -->
                <tr v-if="hasMoreFiles">
                  <td colspan="5" class="py-3 text-center text-[11px] text-slate-500 font-mono">
                    <div class="flex items-center justify-center gap-2">
                      <RefreshCw class="w-3 h-3 animate-spin text-sky-400" />
                      <span>Menampilkan {{ visibleFiles.length }} dari {{ filteredFiles.length }} berkas (gulir untuk memuat lagi)...</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Empty Notice -->
          <div v-if="filteredFiles.length === 0 && !isLoading" class="p-12 text-center text-slate-500 flex flex-col items-center gap-2">
            <template v-if="activeSidebarSection === 'starred'">
              <div class="w-16 h-16 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 mb-2">
                <Star class="w-8 h-8 fill-amber-400/30" />
              </div>
              <p class="font-bold text-slate-200 text-sm">Belum Ada Berkas Berbintang</p>
              <span class="text-xs text-slate-400 max-w-sm leading-relaxed">
                Tandai berkas atau folder favorit Anda dari direktori mana pun dengan mengklik ikon bintang (<Star class="w-3 h-3 inline text-amber-400 fill-amber-400" />) agar tersimpan rapi dan mudah diakses di sini.
              </span>
            </template>
            <template v-else>
              <Folder class="w-12 h-12 text-slate-700" />
              <p class="font-medium text-slate-400">Tidak ada berkas yang cocok</p>
              <span class="text-xs text-slate-600">Folder kosong atau kata kunci pencarian tidak ditemukan.</span>
            </template>
          </div>
        </div>

        <!-- BOTTOM STATUS & QUICK ACTION BAR (Windows Explorer / Finder Style) -->
        <div class="h-10 border-t border-slate-800/80 px-4 flex items-center justify-between bg-slate-900/60 text-xs shrink-0 select-none">
          <!-- Kiri: Info Ringkas Item Terpilih -->
          <div class="flex items-center gap-3 truncate text-slate-400">
            <span class="font-mono text-[11px] text-slate-400">{{ filteredFiles.length }} item</span>
            <span class="text-slate-700">|</span>
            <template v-if="selectedItems.length > 1">
              <div class="flex items-center gap-2 truncate">
                <span class="px-2 py-0.5 rounded-lg bg-sky-500/20 text-sky-300 font-bold font-mono text-[11px] border border-sky-500/30">
                  {{ selectedItems.length }} berkas terpilih
                </span>
                <button @click="clearSelection" class="text-[11px] text-slate-400 hover:text-slate-200 underline">
                  Batal
                </button>
              </div>
            </template>
            <template v-else-if="selectedItem">
              <div class="flex items-center gap-2 truncate">
                <span class="text-sky-400 font-semibold truncate max-w-[260px] sm:max-w-md">{{ selectedItem.name }}</span>
                <span class="text-slate-400 font-mono text-[11px]">({{ selectedItem.is_dir ? 'Folder' : formatBytes(selectedItem.size) }})</span>
                <span class="text-slate-400 text-[11px] hidden sm:inline">• Diubah {{ formatDate(selectedItem.mod_time) }}</span>
              </div>
            </template>
            <template v-else>
              <span class="text-[11px] text-slate-400">Klik item untuk melihat info, Ctrl/Shift untuk multi-pilihan</span>
            </template>
          </div>

          <!-- Kanan: Aksi Cepat Item Terpilih -->
          <!-- Multi-Item Action Buttons -->
          <div v-if="selectedItems.length > 1" class="flex items-center gap-1.5 shrink-0">
            <button
              @click="selectAllFiles"
              class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-medium transition"
              title="Pilih Seluruh Berkas (Ctrl+A)"
            >
              Semua (Ctrl+A)
            </button>
            <button
              v-if="authStore.canDeleteStorage()"
              @click="openDeleteConfirm(selectedItems)"
              class="px-2.5 py-1 rounded-lg bg-rose-500/20 hover:bg-rose-500/30 text-rose-300 border border-rose-500/30 text-[11px] font-semibold flex items-center gap-1.5 transition active:scale-95"
              title="Hapus Seluruh Berkas Terpilih"
            >
              <Trash2 class="w-3.5 h-3.5" />
              <span>Hapus ({{ selectedItems.length }})</span>
            </button>
          </div>

          <!-- Single Item Action Buttons -->
          <div v-else-if="selectedItem" class="flex items-center gap-1.5 shrink-0">
            <button
              @click="openInfoModal(selectedItem)"
              class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-[11px] font-medium flex items-center gap-1 transition"
              title="Lihat Detail & Properti Lengkap"
            >
              <Info class="w-3.5 h-3.5 text-sky-400" />
              <span class="hidden sm:inline">Info Berkas</span>
            </button>
            <button
              v-if="!selectedItem.is_dir"
              @click="openInTextEditor(selectedItem)"
              class="px-2.5 py-1 rounded-lg bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-300 text-[11px] font-medium flex items-center gap-1 transition"
              title="Edit Berkas di Editor Teks"
            >
              <Code2 class="w-3.5 h-3.5" />
              <span class="hidden sm:inline">Edit Berkas</span>
            </button>
            <button
              v-if="!selectedItem.is_dir"
              @click="handleItemDoubleClick(selectedItem)"
              class="px-2.5 py-1 rounded-lg bg-sky-500/20 hover:bg-sky-500/30 text-sky-400 text-[11px] font-medium flex items-center gap-1 transition"
              title="Buka Viewer / Player"
            >
              <Eye class="w-3.5 h-3.5" />
              <span class="hidden sm:inline">Buka</span>
            </button>
            <button
              @click="openRenameDialog(selectedItem)"
              class="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
              title="Ubah Nama"
            >
              <Edit3 class="w-3.5 h-3.5" />
            </button>
            <button
              @click="openDeleteConfirm(selectedItem)"
              class="p-1 rounded-lg hover:bg-rose-500/20 text-rose-400 transition"
              title="Hapus"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Buat Folder Baru (Floating Window Dialog) -->
    <div
      v-if="showNewFolderModal"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-4 py-2.5 select-none">
          <div class="flex items-center gap-2">
            <div class="w-6 h-6 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <Folder class="w-3.5 h-3.5" />
            </div>
            <h4 class="font-bold text-slate-100 text-xs">Buat Folder Baru</h4>
          </div>
          <button @click="showNewFolderModal = false" class="w-6 h-6 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="p-4 space-y-3">
          <label class="block text-[11px] font-medium text-slate-300">Nama Folder Baru di Direktori Ini</label>
          <input
            v-model="newFolderName"
            @keyup.enter="confirmCreateFolder"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
            placeholder="Ketik nama folder..."
          />
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-4 py-2.5">
          <button @click="showNewFolderModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition">
            Batal
          </button>
          <button @click="confirmCreateFolder" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow active:scale-95">
            Buat Folder
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Rename File (Floating Window Dialog) -->
    <div
      v-if="showRenameModal"
      class="fixed inset-0 bg-black/35 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-4 py-2.5 select-none">
          <div class="flex items-center gap-2">
            <div class="w-6 h-6 rounded-lg bg-amber-500/20 text-amber-400 flex items-center justify-center border border-amber-500/30">
              <Edit3 class="w-3.5 h-3.5" />
            </div>
            <h4 class="font-bold text-slate-100 text-xs">Ubah Nama Berkas</h4>
          </div>
          <button @click="showRenameModal = false" class="w-6 h-6 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="p-4 space-y-3">
          <label class="block text-[11px] font-medium text-slate-300">Nama Baru Berkas / Folder</label>
          <input
            v-model="renameTargetNewName"
            @keyup.enter="confirmRename"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs outline-none focus:border-sky-400"
          />
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-4 py-2.5">
          <button @click="showRenameModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition">
            Batal
          </button>
          <button @click="confirmRename" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition shadow active:scale-95">
            Simpan Nama
          </button>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Buat File Baru (Teks / Kode) -->
    <div
      v-if="showNewFileModal"
      class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <!-- Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-4 py-2.5 select-none">
          <div class="flex items-center gap-2">
            <div class="w-6 h-6 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center border border-emerald-500/30">
              <FilePlus class="w-3.5 h-3.5" />
            </div>
            <h4 class="font-bold text-slate-100 text-xs">Buat Berkas Baru (Teks / Kode)</h4>
          </div>
          <button @click="showNewFileModal = false" class="w-6 h-6 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="p-4 space-y-3.5">
          <div>
            <label class="block text-[11px] font-medium text-slate-300 mb-1">Nama Berkas & Ekstensi</label>
            <input
              v-model="newFileName"
              @keyup.enter="confirmCreateFile(true)"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 text-xs font-mono outline-none focus:border-emerald-400"
              placeholder="contoh: script.js, index.html, server.go, config.yaml..."
            />
          </div>

          <!-- Quick Template Selector -->
          <div>
            <span class="block text-[10px] font-bold uppercase text-slate-500 tracking-wider mb-1.5">Pilih Format / Ekstensi Cepat</span>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="tpl in [
                  { id: 'txt', label: '.txt (Teks)' },
                  { id: 'js', label: '.js (JavaScript)' },
                  { id: 'html', label: '.html (Web)' },
                  { id: 'css', label: '.css (Style)' },
                  { id: 'json', label: '.json (Data)' },
                  { id: 'py', label: '.py (Python)' },
                  { id: 'go', label: '.go (GoLang)' },
                  { id: 'md', label: '.md (Markdown)' },
                  { id: 'yaml', label: '.yaml (Docker/Config)' }
                ]"
                :key="tpl.id"
                @click="() => {
                  newFileTemplate = tpl.id
                  if (!newFileName.includes('.')) {
                    newFileName = newFileName ? `${newFileName.split('.')[0]}.${tpl.id}` : `untitled.${tpl.id}`
                  } else {
                    newFileName = `${newFileName.split('.')[0]}.${tpl.id}`
                  }
                }"
                class="px-2 py-1 rounded-lg border text-[11px] font-mono transition"
                :class="newFileTemplate === tpl.id ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40 font-bold' : 'bg-slate-950/60 text-slate-400 border-slate-800 hover:text-slate-200'"
              >
                {{ tpl.label }}
              </button>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between bg-slate-950/75 border-t border-slate-800 px-4 py-2.5">
          <button @click="showNewFileModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition">
            Batal
          </button>
          <div class="flex items-center gap-2">
            <button @click="confirmCreateFile(false)" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium transition">
              Buat Saja
            </button>
            <button @click="confirmCreateFile(true)" class="px-3.5 py-1.5 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs transition flex items-center gap-1.5 shadow active:scale-95">
              <Code2 class="w-3.5 h-3.5" />
              <span>Buat & Buka Editor</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Progress Upload File / Folder / ZIP -->
    <div
      v-if="showUploadModal"
      class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black ring-1 ring-white/10 rounded-2xl flex flex-col overflow-hidden">
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-4 py-2.5 select-none">
          <div class="flex items-center gap-2">
            <div class="w-6 h-6 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30">
              <UploadCloud class="w-3.5 h-3.5" />
            </div>
            <h4 class="font-bold text-slate-100 text-xs">
              {{ uploadType === 'zip' ? 'Mengekstrak Arsip ZIP' : 'Mengunggah Berkas' }}
            </h4>
          </div>
          <button v-if="!isUploading" @click="showUploadModal = false" class="w-6 h-6 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="p-4 space-y-3">
          <div class="flex items-center justify-between text-xs font-mono">
            <span class="text-slate-300 truncate max-w-[220px]">{{ uploadCurrentFileName || 'Memproses berkas...' }}</span>
            <span class="text-sky-400 font-bold">{{ uploadProgress }}%</span>
          </div>

          <!-- Progress Bar -->
          <div class="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-sky-500 to-emerald-400 transition-all duration-300"
              :style="{ width: `${uploadProgress}%` }"
            ></div>
          </div>

          <div class="flex items-center justify-between text-[11px] text-slate-400">
            <span>{{ uploadCompletedCount }} dari {{ uploadTotalCount }} item selesai</span>
            <span v-if="isUploading" class="flex items-center gap-1 text-sky-400">
              <RefreshCw class="w-3 h-3 animate-spin" />
              <span>Mengunggah...</span>
            </span>
            <span v-else class="text-emerald-400 flex items-center gap-1 font-semibold">
              <Check class="w-3.5 h-3.5" />
              <span>Selesai</span>
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- CUSTOM MODAL: Konfirmasi Hapus Berkas / Folder (Glassmorphism Slate-900) -->
    <div
      v-if="showDeleteConfirmModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md bg-slate-900 border border-rose-500/40 shadow-2xl shadow-black ring-1 ring-rose-500/20 rounded-2xl flex flex-col overflow-hidden text-slate-100">
        <!-- Modal Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5 text-rose-400">
            <div class="w-7 h-7 rounded-lg bg-rose-500/20 text-rose-400 flex items-center justify-center border border-rose-500/30">
              <Trash2 class="w-4 h-4" />
            </div>
            <h3 class="font-bold text-sm">
              {{ itemsToDelete.length > 1 ? `Hapus ${itemsToDelete.length} Berkas Terpilih` : 'Hapus Berkas / Folder' }}
            </h3>
          </div>
          <button
            @click="showDeleteConfirmModal = false"
            :disabled="isDeletingFiles"
            class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95 disabled:opacity-50"
            title="Batal & Tutup"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Modal Body Content -->
        <div class="p-5 space-y-4">
          <div class="flex items-start gap-3 p-3 rounded-xl bg-rose-500/10 border border-rose-500/25 text-rose-200 text-xs leading-relaxed">
            <AlertTriangle class="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />
            <div>
              <span class="font-bold block mb-0.5 text-rose-300">Peringatan Penghapusan Permanen</span>
              <span>Tindakan ini tidak dapat dibatalkan. Berkas atau direktori berikut akan dihapus secara langsung dari disk host:</span>
            </div>
          </div>

          <!-- Preview Berkas yang Akan Dihapus -->
          <div class="max-h-48 overflow-y-auto rounded-xl bg-slate-950 border border-slate-800/80 divide-y divide-slate-800/50 text-xs font-mono">
            <div
              v-for="target in itemsToDelete"
              :key="target.path"
              class="px-3 py-2 flex items-center justify-between gap-2"
            >
              <div class="flex items-center gap-2 truncate text-slate-200">
                <Folder v-if="target.is_dir" class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                <File v-else class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                <span class="truncate" :title="target.name">{{ target.name }}</span>
              </div>
              <span class="text-[10px] text-slate-500 shrink-0">
                {{ target.is_dir ? 'Folder' : formatBytes(target.size) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Modal Footer Actions -->
        <div class="flex items-center justify-end gap-2 bg-slate-950/80 border-t border-slate-800 px-5 py-3 select-none">
          <button
            @click="showDeleteConfirmModal = false"
            :disabled="isDeletingFiles"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition active:scale-95 disabled:opacity-50"
          >
            Batal
          </button>
          <button
            @click="executeDeleteConfirmed"
            :disabled="isDeletingFiles"
            class="px-5 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs transition flex items-center gap-2 shadow-lg shadow-rose-950/50 active:scale-95 disabled:opacity-50"
          >
            <RefreshCw v-if="isDeletingFiles" class="w-3.5 h-3.5 animate-spin" />
            <Trash2 v-else class="w-3.5 h-3.5" />
            <span>{{ isDeletingFiles ? 'Menghapus...' : 'Ya, Hapus Permanen' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Hidden Native Inputs untuk File, Folder, dan ZIP Upload -->
    <input
      ref="fileUploadInputRef"
      type="file"
      multiple
      @change="handleFileUploadChange"
      class="hidden"
    />
    <input
      ref="folderUploadInputRef"
      type="file"
      webkitdirectory
      directory
      multiple
      @change="handleFolderUploadChange"
      class="hidden"
    />
    <input
      ref="zipUploadInputRef"
      type="file"
      accept=".zip,application/zip"
      @change="handleZipUploadChange"
      class="hidden"
    />
  </div>
</template>
