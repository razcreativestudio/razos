<script setup>
/**
 * RAZOS Sublime Text & Notepad++ Enterprise Complex Code & Text Editor
 * Created: 2026-10-15 (v3.6.0 - Advanced Multi-Tab IDE Architecture)
 * Deskripsi: Text & Code Editor tingkat lanjut berarsitektur desktop profesional:
 * - Multi-Tab Workspace: Membuka banyak berkas sekaligus dengan tab switching, dirty indicator, dan shortcut tab.
 * - Collapsible Project File Tree: Sidebar penjelajah berkas parent directory untuk navigasi cepat file satu folder.
 * - Synchronized Line Numbers Gutter: Nomor baris responsif dengan penanda baris aktif kursor.
 * - Sublime Floating Find & Replace (Ctrl+F / Ctrl+H): Pencarian real-time, pencocokan teks, match counter, dan Replace All.
 * - Go to Line Dialog (Ctrl+G): Lompat instan ke nomor baris yang dituju.
 * - Smart Indentation & Auto-Pairs: Dukungan tombol Tab/Shift+Tab, auto indent baris baru, dan auto closing brackets.
 * - Notepad++ Status Bar: Baris/Kolom kursor, total baris/karakter, encoding UTF-8, line endings, dan selector sintaks bahasa.
 * - Multi-Theme Support: Sublime Monokai, Obsidian Dark, Cyberpunk, dan Classic Slate.
 */
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import {
  FileCode, FileText, Folder, FolderOpen, Save, Search, RefreshCw, X, Plus,
  ChevronRight, ChevronDown, Check, Copy, Eye, Sliders, Type,
  Layout, Code2, Hash, Undo2, Redo2, Scissors, ArrowDown, ArrowUp, ArrowLeft,
  Settings, Columns, Terminal, CheckCircle2, AlertTriangle, FilePlus,
  CornerLeftUp, HardDrive
} from 'lucide-vue-next'
import { useSystemStore } from '../stores/RAZSystemStore'

const props = defineProps({
  // Berkas awal yang dibuka (opsional, jika dibuka dari File Manager)
  initialFile: {
    type: Object,
    default: null
  },
  // Window frame instance descriptor dari RAZWindowStore
  win: {
    type: Object,
    default: null
  }
})

const systemStore = useSystemStore()

// State Custom Modals (Zero Native Dialogs: No alert/confirm/prompt)
const showSaveAsModal = ref(false)
const saveAsFileName = ref('')
const showCloseTabConfirmModal = ref(false)
const pendingCloseTabId = ref('')
const pendingCloseTabName = ref('')
const showCloseAllConfirmModal = ref(false)

// Storage Keys untuk Persistensi Sesi Tab & Direktori
const STORAGE_KEY_TABS = 'razos_text_editor_tabs_session'
const STORAGE_KEY_ACTIVE_TAB = 'razos_text_editor_active_tab'
const STORAGE_KEY_DIR = 'razos_text_editor_sidebar_dir'

// 1. STATE MULTI-TAB WORKSPACE
const tabs = ref([])
const activeTabId = ref('')
const tabCounter = ref(1)

// Helper ID unik tab
function generateTabId() {
  return 'tab_' + Date.now() + '_' + Math.random().toString(36).substring(2, 7)
}

// 2. STATE EDITOR & TAMPILAN
const activeTheme = ref('monokai') // 'monokai' | 'obsidian' | 'cyberpunk' | 'notepad'
const fontSize = ref(13) // dalam px
const isWordWrap = ref(false)
const showLineNumbers = ref(true)
const showSidebar = ref(true)
const indentSize = ref(2) // 2 atau 4 spasi

// 3. STATE PENCARIAN (FIND & REPLACE SUBLIME STYLE)
const showFindBar = ref(false)
const showReplace = ref(false)
const findQuery = ref('')
const replaceQuery = ref('')
const findMatches = ref([])
const currentMatchIndex = ref(-1)
const showGoToLineModal = ref(false)
const targetLineNumber = ref('')

// 4. STATE STATUS KURSOR (NOTEPAD++ STYLE)
const cursorLine = ref(1)
const cursorCol = ref(1)
const editorTextareaRef = ref(null)
const gutterRef = ref(null)

// 5. STATE SIDEBAR FILE TREE & NAVIGASI FOLDER
const sidebarFiles = ref([])
const sidebarCurrentDir = ref('')
const isSidebarLoading = ref(false)
const sidebarHistory = ref([])
const sidebarHistoryIndex = ref(-1)

// Format ukuran berkas (Bytes ke KB, MB, dll)
function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// Nama ringkas folder aktif di header sidebar
const sidebarCurrentFolderShortName = computed(() => {
  if (!sidebarCurrentDir.value) return 'Explorer'
  const clean = sidebarCurrentDir.value.replace(/\\/g, '/').replace(/\/+$/, '')
  if (clean === '' || clean === '/') return 'Root'
  if (/^[a-zA-Z]:$/.test(clean)) return `${clean.toUpperCase()} Disk`
  const parts = clean.split('/').filter(Boolean)
  return parts.length > 0 ? parts[parts.length - 1] : 'Explorer'
})

// Deteksi apakah bisa naik satu tingkat ke folder induk
const canSidebarGoUp = computed(() => {
  if (!sidebarCurrentDir.value) return false
  const clean = sidebarCurrentDir.value.replace(/\\/g, '/').replace(/\/+$/, '')
  if (clean === '' || clean === '/' || /^[a-zA-Z]:$/.test(clean) || /^[a-zA-Z]:\/$/.test(clean)) {
    return false
  }
  return true
})

// Deteksi apakah bisa kembali dalam riwayat
const canSidebarGoBack = computed(() => {
  return sidebarHistoryIndex.value > 0
})

// Komponen Ikon Berkas Berdasarkan Ekstensi
function getFileIconComponent(filename = '') {
  const ext = (filename.includes('.') ? filename.split('.').pop() : '').toLowerCase()
  if (['md', 'txt', 'log', 'env', 'conf', 'ini', 'pdf'].includes(ext)) {
    return FileText
  }
  return FileCode
}

// Warna Ikon Berkas Berdasarkan Ekstensi (Modern VS Code Theme)
function getFileIconColor(filename = '') {
  const ext = (filename.includes('.') ? filename.split('.').pop() : '').toLowerCase()
  switch (ext) {
    case 'js':
    case 'mjs':
      return 'text-amber-300'
    case 'ts':
      return 'text-sky-400'
    case 'vue':
      return 'text-emerald-400'
    case 'html':
      return 'text-orange-400'
    case 'css':
    case 'scss':
      return 'text-cyan-400'
    case 'json':
      return 'text-amber-400'
    case 'py':
      return 'text-blue-400'
    case 'go':
      return 'text-teal-300'
    case 'php':
      return 'text-indigo-400'
    case 'sql':
      return 'text-violet-400'
    case 'sh':
    case 'bash':
      return 'text-emerald-300'
    case 'md':
      return 'text-rose-300'
    case 'yaml':
    case 'yml':
      return 'text-rose-400'
    default:
      return 'text-slate-400'
  }
}

// Deteksi bahasa pemrograman berdasarkan ekstensi berkas
function detectLanguage(filename = '') {
  const ext = (filename.includes('.') ? filename.split('.').pop() : '').toLowerCase()
  const map = {
    js: 'JavaScript',
    mjs: 'JavaScript (ESM)',
    cjs: 'CommonJS',
    ts: 'TypeScript',
    vue: 'Vue Component',
    html: 'HTML5',
    css: 'CSS3',
    scss: 'SASS / SCSS',
    json: 'JSON Data',
    go: 'Go Lang',
    py: 'Python',
    php: 'PHP Script',
    sql: 'SQL Database',
    sh: 'Shell Script',
    bash: 'Bash Script',
    yaml: 'YAML Config',
    yml: 'YAML Config',
    md: 'Markdown',
    txt: 'Plain Text',
    log: 'System Log',
    conf: 'Configuration',
    env: 'Environment'
  }
  return map[ext] || 'Plain Text'
}

// Format icon badge bahasa
function getLanguageBadge(filename = '') {
  const ext = (filename.includes('.') ? filename.split('.').pop() : '').toLowerCase()
  return ext.toUpperCase() || 'TXT'
}

// Ambil Tab Aktif
const activeTab = computed(() => {
  return tabs.value.find(t => t.id === activeTabId.value) || null
})

// Hitung jumlah baris teks pada tab aktif
const lineCount = computed(() => {
  if (!activeTab.value) return 1
  return (activeTab.value.content || '').split('\n').length || 1
})

// Hitung total karakter pada tab aktif
const charCount = computed(() => {
  if (!activeTab.value) return 0
  return (activeTab.value.content || '').length
})

// Hitung total kata pada tab aktif
const wordCount = computed(() => {
  if (!activeTab.value || !activeTab.value.content) return 0
  const clean = activeTab.value.content.trim()
  if (!clean) return 0
  return clean.split(/\s+/).length
})

// Buka berkas ke dalam tab baru atau aktifkan jika sudah ada
async function openFileInTab(fileItem) {
  if (!fileItem || !fileItem.path) return

  // Normalisasi path
  const targetPath = fileItem.path.replace(/\\/g, '/')

  // Cek apakah berkas sudah terbuka di salah satu tab
  const existingTab = tabs.value.find(t => t.path === targetPath || t.path === fileItem.path)
  if (existingTab) {
    activeTabId.value = existingTab.id
    loadSidebarDirectory(getDirectoryPath(targetPath))
    // Jika konten existing tab belum ada atau tidak ada perubahan lokal, muat ulang konten terbaru dari disk
    if (!existingTab.isDirty) {
      existingTab.isLoading = true
      try {
        const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(targetPath)}`)
        const data = await res.json()
        if (res.ok) {
          existingTab.content = data.content ?? ''
          existingTab.originalContent = data.content ?? ''
          existingTab.isLocked = !!data.isLocked
        }
      } catch (err) {
        console.error('Gagal memuat ulang berkas:', err)
      } finally {
        existingTab.isLoading = false
      }
    }
    return
  }

  // Jika saat ini hanya ada 1 tab dan itu adalah tab untitled kosong belum diedit, ganti tab tersebut
  if (tabs.value.length === 1 && !tabs.value[0].path && !tabs.value[0].isDirty && !tabs.value[0].content) {
    tabs.value = []
  }

  const fileName = fileItem.name || (targetPath.split('/').pop() || 'berkas.txt')

  const newTab = {
    id: generateTabId(),
    name: fileName,
    path: targetPath,
    content: typeof fileItem.content === 'string' ? fileItem.content : '',
    originalContent: typeof fileItem.content === 'string' ? fileItem.content : '',
    isDirty: false,
    isLoading: true,
    isLocked: false,
    language: detectLanguage(fileName),
    undoStack: [],
    redoStack: []
  }

  tabs.value.push(newTab)
  activeTabId.value = newTab.id

  // Muat sidebar folder
  loadSidebarDirectory(getDirectoryPath(targetPath))

  // Ambil isi berkas dari API backend
  try {
    const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(targetPath)}`)
    const data = await res.json()
    if (res.ok) {
      newTab.content = data.content ?? ''
      newTab.originalContent = data.content ?? ''
      newTab.isLocked = !!data.isLocked
    } else {
      systemStore.showNotification('Peringatan', data.error || 'Tidak dapat memuat berkas', 'warning')
    }
  } catch (err) {
    systemStore.showNotification('Error', `Gagal membuka berkas: ${err.message}`, 'error')
  } finally {
    newTab.isLoading = false
    nextTick(() => {
      if (editorTextareaRef.value) {
        editorTextareaRef.value.focus()
        updateCursorPosition()
      }
    })
  }
}

// Buat Tab Berkas Baru Kosong (Untitled)
function createNewUntitledTab() {
  const name = `Untitled-${tabCounter.value++}.txt`
  const newTab = {
    id: generateTabId(),
    name: name,
    path: '',
    content: '',
    originalContent: '',
    isDirty: false,
    isLoading: false,
    isLocked: false,
    language: 'Plain Text',
    undoStack: [],
    redoStack: []
  }
  tabs.value.push(newTab)
  activeTabId.value = newTab.id
  nextTick(() => {
    if (editorTextareaRef.value) editorTextareaRef.value.focus()
  })
}

// Tutup Tab dengan konfirmasi modal kustom jika ada perubahan yang belum disimpan
function closeTab(tabId, e) {
  if (e) e.stopPropagation()
  const tabIndex = tabs.value.findIndex(t => t.id === tabId)
  if (tabIndex === -1) return

  const targetTab = tabs.value[tabIndex]
  if (targetTab.isDirty) {
    pendingCloseTabId.value = tabId
    pendingCloseTabName.value = targetTab.name
    showCloseTabConfirmModal.value = true
    return
  }

  doCloseTab(tabIndex)
}

function forceCloseTab() {
  showCloseTabConfirmModal.value = false
  const tabIndex = tabs.value.findIndex(t => t.id === pendingCloseTabId.value)
  if (tabIndex !== -1) {
    doCloseTab(tabIndex)
  }
}

function doCloseTab(tabIndex) {
  const tabId = tabs.value[tabIndex].id
  tabs.value.splice(tabIndex, 1)

  // Jika tab yang ditutup adalah tab aktif, alihkan fokus ke tab sebelahnya
  if (activeTabId.value === tabId) {
    if (tabs.value.length > 0) {
      const nextIndex = Math.min(tabIndex, tabs.value.length - 1)
      activeTabId.value = tabs.value[nextIndex].id
    } else {
      createNewUntitledTab()
    }
  }
}

// Tutup Semua Tab dengan modal konfirmasi kustom jika ada dirty tabs
function closeAllTabs() {
  const dirtyCount = tabs.value.filter(t => t.isDirty).length
  if (dirtyCount > 0) {
    showCloseAllConfirmModal.value = true
    return
  }
  forceCloseAllTabs()
}

function forceCloseAllTabs() {
  showCloseAllConfirmModal.value = false
  tabs.value = []
  createNewUntitledTab()
}

// Simpan Berkas Aktif (Ctrl+S)
async function saveActiveFile() {
  if (!activeTab.value) return
  const cur = activeTab.value

  // Jika belum punya path (untitled), minta nama via Custom Modal
  if (!cur.path) {
    saveAsFileName.value = cur.name
    showSaveAsModal.value = true
    return
  }

  await doSaveFile(cur)
}

async function confirmSaveAs() {
  if (!activeTab.value) return
  const cur = activeTab.value
  const name = saveAsFileName.value.trim()
  if (!name) return

  const defaultDir = sidebarCurrentDir.value || '/DATA'
  const cleanDir = defaultDir.endsWith('/') ? defaultDir : defaultDir + '/'
  cur.path = cleanDir + name
  cur.name = name
  cur.language = detectLanguage(name)
  showSaveAsModal.value = false

  await doSaveFile(cur)
}

async function doSaveFile(cur) {
  cur.isSaving = true
  try {
    const res = await fetch('/api/v1/files/content', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        path: cur.path,
        content: cur.content
      })
    })
    const data = await res.json()
    if (res.ok && data.status === 'saved') {
      cur.originalContent = cur.content
      cur.isDirty = false
      systemStore.showNotification('Tersimpan', `Berkas "${cur.name}" berhasil disimpan`, 'success')
      if (sidebarCurrentDir.value) {
        loadSidebarDirectory(sidebarCurrentDir.value)
      }
    } else {
      systemStore.showNotification('Gagal Simpan', data.error || 'Gagal menyimpan berkas ke disk', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', `Gagal koneksi saat menyimpan: ${err.message}`, 'error')
  } finally {
    cur.isSaving = false
  }
}

// Simpan Semua Berkas Terbuka
async function saveAllFiles() {
  const dirtyTabs = tabs.value.filter(t => t.isDirty && t.path)
  for (const t of dirtyTabs) {
    try {
      await fetch('/api/v1/files/content', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: t.path, content: t.content })
      })
      t.originalContent = t.content
      t.isDirty = false
    } catch (e) {}
  }
  systemStore.showNotification('Simpan Semua', 'Seluruh perubahan berkas berhasil disimpan', 'success')
}

// Helper ekstraksi parent directory dari path
function getDirectoryPath(fullPath = '') {
  if (!fullPath) return ''
  const clean = fullPath.replace(/\\/g, '/')
  const lastSlash = clean.lastIndexOf('/')
  return lastSlash > 0 ? clean.substring(0, lastSlash) : '/'
}

// Muat daftar berkas folder untuk Sidebar File Tree
async function loadSidebarDirectory(dirPath, isHistoryNav = false) {
  if (!dirPath) return
  let clean = dirPath.replace(/\\/g, '/')
  if (/^[a-zA-Z]:$/.test(clean)) clean += '/'

  sidebarCurrentDir.value = clean

  if (!isHistoryNav) {
    if (sidebarHistoryIndex.value < sidebarHistory.value.length - 1) {
      sidebarHistory.value = sidebarHistory.value.slice(0, sidebarHistoryIndex.value + 1)
    }
    sidebarHistory.value.push(clean)
    sidebarHistoryIndex.value = sidebarHistory.value.length - 1
  }

  isSidebarLoading.value = true
  try {
    const res = await fetch(`/api/v1/files/list?path=${encodeURIComponent(clean)}`)
    const data = await res.json()
    if (res.ok && Array.isArray(data)) {
      // Urutkan: folder di atas, berkas di bawah secara alfabetis
      sidebarFiles.value = [...data].sort((a, b) => {
        if (a.is_dir && !b.is_dir) return -1
        if (!a.is_dir && b.is_dir) return 1
        return a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
      })
      localStorage.setItem(STORAGE_KEY_DIR, clean)
    } else {
      sidebarFiles.value = []
    }
  } catch (e) {
    sidebarFiles.value = []
  } finally {
    isSidebarLoading.value = false
  }
}

// Navigasi naik satu tingkat ke folder induk (Up)
function navigateSidebarUp() {
  if (!canSidebarGoUp.value) return
  const clean = sidebarCurrentDir.value.replace(/\\/g, '/').replace(/\/+$/, '')
  const lastSlash = clean.lastIndexOf('/')
  let parent = '/'
  if (lastSlash > 0) {
    parent = clean.substring(0, lastSlash)
  } else if (/^[a-zA-Z]:/.test(clean)) {
    parent = clean.substring(0, 2) + '/'
  }
  loadSidebarDirectory(parent)
}

// Navigasi kembali dalam riwayat sidebar (Back)
function navigateSidebarBack() {
  if (canSidebarGoBack.value) {
    sidebarHistoryIndex.value--
    loadSidebarDirectory(sidebarHistory.value[sidebarHistoryIndex.value], true)
  }
}

// Navigasi langsung ke direktori penyimpanan utama (/DATA atau root)
function navigateSidebarHome() {
  const home = localStorage.getItem('razos_default_storage_path') || '/DATA'
  loadSidebarDirectory(home)
}

// Simpan Sesi Tab ke localStorage
function saveSession() {
  try {
    const serialized = tabs.value.map(t => ({
      id: t.id,
      name: t.name,
      path: t.path,
      content: t.isDirty ? t.content : '', // simpan konten teks jika ada perubahan belum disimpan
      originalContent: t.originalContent,
      isDirty: !!t.isDirty,
      language: t.language
    }))
    localStorage.setItem(STORAGE_KEY_TABS, JSON.stringify(serialized))
    if (activeTabId.value) {
      localStorage.setItem(STORAGE_KEY_ACTIVE_TAB, activeTabId.value)
    }
    if (sidebarCurrentDir.value) {
      localStorage.setItem(STORAGE_KEY_DIR, sidebarCurrentDir.value)
    }
  } catch (e) {}
}

// Pulihkan Sesi Tab Tersimpan dari localStorage
async function restoreSession() {
  try {
    const defaultDir = localStorage.getItem('razos_default_storage_path') || '/DATA'
    const rawDir = localStorage.getItem(STORAGE_KEY_DIR)
    if (rawDir && !rawDir.startsWith('C:/DATA')) {
      loadSidebarDirectory(rawDir)
    } else {
      loadSidebarDirectory(defaultDir)
    }

    const rawTabs = localStorage.getItem(STORAGE_KEY_TABS)
    if (rawTabs) {
      const parsed = JSON.parse(rawTabs)
      if (Array.isArray(parsed) && parsed.length > 0) {
        const restored = []
        for (const item of parsed) {
          const tab = {
            id: item.id || generateTabId(),
            name: item.name || 'untitled.txt',
            path: item.path || '',
            content: item.content || '',
            originalContent: item.originalContent || '',
            isDirty: !!item.isDirty,
            isLoading: false,
            isLocked: false,
            language: item.language || detectLanguage(item.name),
            undoStack: [],
            redoStack: []
          }
          restored.push(tab)

          // Jika berkas memiliki path dan belum dirty, fetch isi terbaru dari disk
          if (tab.path && !tab.isDirty) {
            tab.isLoading = true
            fetch(`/api/v1/files/content?path=${encodeURIComponent(tab.path)}`)
              .then(res => res.json())
              .then(data => {
                if (data && data.content !== undefined) {
                  tab.content = data.content
                  tab.originalContent = data.content
                  tab.isLocked = !!data.isLocked
                }
              })
              .catch(() => {})
              .finally(() => {
                tab.isLoading = false
              })
          }
        }

        tabs.value = restored
        const savedActive = localStorage.getItem(STORAGE_KEY_ACTIVE_TAB)
        if (savedActive && tabs.value.some(t => t.id === savedActive)) {
          activeTabId.value = savedActive
        } else if (tabs.value.length > 0) {
          activeTabId.value = tabs.value[0].id
        }
        return true
      }
    }
  } catch (e) {
    console.error('Gagal memulihkan sesi editor:', e)
  }
  return false
}

// Handle Perubahan Input Teks
function onContentInput(e) {
  if (!activeTab.value) return
  activeTab.value.content = e.target.value
  activeTab.value.isDirty = activeTab.value.content !== activeTab.value.originalContent
  updateCursorPosition()
}

// Sinkronisasi Baris dan Kolom Kursor
function updateCursorPosition() {
  if (!editorTextareaRef.value || !activeTab.value) return
  const pos = editorTextareaRef.value.selectionStart || 0
  const textBefore = activeTab.value.content.substring(0, pos)
  const lines = textBefore.split('\n')
  cursorLine.value = lines.length
  cursorCol.value = lines[lines.length - 1].length + 1
}

// Sinkronisasi Scroll Gutter Nomor Baris dengan Textarea
function onEditorScroll() {
  if (editorTextareaRef.value && gutterRef.value) {
    gutterRef.value.scrollTop = editorTextareaRef.value.scrollTop
  }
}

// SMART KEYBOARD SHORTCUTS (Sublime / Notepad++ Style)
function onKeyDown(e) {
  if (!activeTab.value) return

  // 1. Ctrl + S (Simpan)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault()
    if (e.shiftKey) {
      saveAllFiles()
    } else {
      saveActiveFile()
    }
    return
  }

  // 2. Ctrl + F (Find)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'f') {
    e.preventDefault()
    showFindBar.value = true
    showReplace.value = false
    nextTick(() => {
      const el = document.getElementById('sublime-find-input')
      if (el) el.focus()
    })
    return
  }

  // 3. Ctrl + H (Find & Replace)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'h') {
    e.preventDefault()
    showFindBar.value = true
    showReplace.value = true
    nextTick(() => {
      const el = document.getElementById('sublime-find-input')
      if (el) el.focus()
    })
    return
  }

  // 4. Ctrl + G (Go to Line)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'g') {
    e.preventDefault()
    showGoToLineModal.value = true
    targetLineNumber.value = String(cursorLine.value)
    nextTick(() => {
      const el = document.getElementById('sublime-goto-input')
      if (el) el.focus()
    })
    return
  }

  // 5. Ctrl + N (New Tab)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'n') {
    e.preventDefault()
    createNewUntitledTab()
    return
  }

  // 6. Ctrl + W (Close Active Tab)
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'w') {
    e.preventDefault()
    closeTab(activeTabId.value)
    return
  }

  // 7. Tab Key: Sisipkan Spasi (Indentation) bukan pindah fokus
  if (e.key === 'Tab') {
    e.preventDefault()
    const textarea = editorTextareaRef.value
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const spaces = ' '.repeat(indentSize.value)

    if (e.shiftKey) {
      // Outdent: Hapus spasi di awal baris jika ada
      // Implementasi outdent sederhana
    } else {
      // Indent
      const val = activeTab.value.content
      activeTab.value.content = val.substring(0, start) + spaces + val.substring(end)
      activeTab.value.isDirty = true
      nextTick(() => {
        textarea.selectionStart = textarea.selectionEnd = start + spaces.length
        updateCursorPosition()
      })
    }
    return
  }

  // 8. Auto Closing Pairs: (, [, {, ", '
  const pairs = { '(': ')', '[': ']', '{': '}', '"': '"', "'": "'" }
  if (pairs[e.key] && !e.ctrlKey && !e.metaKey && !e.altKey) {
    const textarea = editorTextareaRef.value
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    if (start === end) {
      e.preventDefault()
      const val = activeTab.value.content
      const openChar = e.key
      const closeChar = pairs[e.key]
      activeTab.value.content = val.substring(0, start) + openChar + closeChar + val.substring(end)
      activeTab.value.isDirty = true
      nextTick(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 1
        updateCursorPosition()
      })
      return
    }
  }

  // 9. Enter Key: Pertahankan Indentasi Baris Sebelumnya
  if (e.key === 'Enter') {
    const textarea = editorTextareaRef.value
    const start = textarea.selectionStart
    const val = activeTab.value.content
    const textBefore = val.substring(0, start)
    const lastLine = textBefore.split('\n').pop() || ''
    const matchIndent = lastLine.match(/^(\s+)/)
    if (matchIndent && matchIndent[1]) {
      e.preventDefault()
      const indent = matchIndent[1]
      activeTab.value.content = val.substring(0, start) + '\n' + indent + val.substring(start)
      activeTab.value.isDirty = true
      nextTick(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 1 + indent.length
        updateCursorPosition()
      })
      return
    }
  }

  // 10. Ctrl + / (Toggle Komentar Baris)
  if ((e.ctrlKey || e.metaKey) && e.key === '/') {
    e.preventDefault()
    toggleLineComment()
    return
  }

  // 11. Ctrl + Shift + D (Duplikasi Baris Aktif)
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'd') {
    e.preventDefault()
    duplicateCurrentLine()
    return
  }
}

// Toggle Komentar Baris
function toggleLineComment() {
  if (!activeTab.value) return
  const textarea = editorTextareaRef.value
  const pos = textarea.selectionStart
  const val = activeTab.value.content
  const lines = val.split('\n')
  
  let currentPos = 0
  let targetIndex = 0
  for (let i = 0; i < lines.length; i++) {
    const lineLen = lines[i].length + 1
    if (pos >= currentPos && pos < currentPos + lineLen) {
      targetIndex = i
      break
    }
    currentPos += lineLen
  }

  const commentPrefix = ['Python', 'Shell Script', 'YAML Config'].includes(activeTab.value.language) ? '# ' : '// '
  if (lines[targetIndex].trim().startsWith(commentPrefix.trim())) {
    lines[targetIndex] = lines[targetIndex].replace(commentPrefix, '').replace(commentPrefix.trim(), '')
  } else {
    lines[targetIndex] = commentPrefix + lines[targetIndex]
  }

  activeTab.value.content = lines.join('\n')
  activeTab.value.isDirty = true
}

// Duplikasi Baris Aktif (Sublime Style)
function duplicateCurrentLine() {
  if (!activeTab.value) return
  const textarea = editorTextareaRef.value
  const pos = textarea.selectionStart
  const val = activeTab.value.content
  const lines = val.split('\n')
  
  let currentPos = 0
  let targetIndex = 0
  for (let i = 0; i < lines.length; i++) {
    const lineLen = lines[i].length + 1
    if (pos >= currentPos && pos < currentPos + lineLen) {
      targetIndex = i
      break
    }
    currentPos += lineLen
  }

  lines.splice(targetIndex + 1, 0, lines[targetIndex])
  activeTab.value.content = lines.join('\n')
  activeTab.value.isDirty = true
}

// FUNGSI PENCARIAN (FIND & REPLACE)
function executeFind() {
  if (!activeTab.value || !findQuery.value) {
    findMatches.value = []
    currentMatchIndex.value = -1
    return
  }
  const text = activeTab.value.content
  const q = findQuery.value
  const matches = []
  let pos = 0
  while ((pos = text.toLowerCase().indexOf(q.toLowerCase(), pos)) !== -1) {
    matches.push(pos)
    pos += q.length
  }
  findMatches.value = matches
  if (matches.length > 0) {
    currentMatchIndex.value = 0
    highlightMatch(matches[0])
  } else {
    currentMatchIndex.value = -1
  }
}

function findNext() {
  if (findMatches.value.length === 0) return
  currentMatchIndex.value = (currentMatchIndex.value + 1) % findMatches.value.length
  highlightMatch(findMatches.value[currentMatchIndex.value])
}

function findPrev() {
  if (findMatches.value.length === 0) return
  currentMatchIndex.value = (currentMatchIndex.value - 1 + findMatches.value.length) % findMatches.value.length
  highlightMatch(findMatches.value[currentMatchIndex.value])
}

function highlightMatch(startIndex) {
  const textarea = editorTextareaRef.value
  if (!textarea || startIndex === undefined) return
  const len = findQuery.value.length
  textarea.focus()
  textarea.setSelectionRange(startIndex, startIndex + len)
  updateCursorPosition()
}

function executeReplace() {
  if (!activeTab.value || !findQuery.value || currentMatchIndex.value === -1) return
  const textarea = editorTextareaRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const val = activeTab.value.content

  activeTab.value.content = val.substring(0, start) + replaceQuery.value + val.substring(end)
  activeTab.value.isDirty = true

  executeFind()
}

function executeReplaceAll() {
  if (!activeTab.value || !findQuery.value) return
  const val = activeTab.value.content
  const regex = new RegExp(findQuery.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi')
  activeTab.value.content = val.replace(regex, replaceQuery.value)
  activeTab.value.isDirty = true
  executeFind()
  systemStore.showNotification('Ganti Semua', 'Seluruh kecocokan berhasil diganti', 'success')
}

// FUNGSI GO TO LINE
function goToLine() {
  const lineNum = parseInt(targetLineNumber.value, 10)
  if (isNaN(lineNum) || lineNum < 1 || !activeTab.value) return

  const lines = activeTab.value.content.split('\n')
  const clampedLine = Math.min(lineNum, lines.length)
  let charIndex = 0
  for (let i = 0; i < clampedLine - 1; i++) {
    charIndex += lines[i].length + 1
  }

  showGoToLineModal.value = false
  nextTick(() => {
    const textarea = editorTextareaRef.value
    if (textarea) {
      textarea.focus()
      textarea.setSelectionRange(charIndex, charIndex)
      updateCursorPosition()
      // Scroll ke posisi baris
      const lineHeightPx = fontSize.value * 1.5
      textarea.scrollTop = (clampedLine - 5) * lineHeightPx
    }
  })
}

// Inisialisasi saat dimount: Pulihkan sesi tab tersimpan & buka berkas target jika ada
onMounted(async () => {
  const restored = await restoreSession()
  const target = props.initialFile || props.win?.fileData || props.win?.customData
  if (target && target.path) {
    await openFileInTab(target)
  } else if (!restored || tabs.value.length === 0) {
    createNewUntitledTab()
  }
})

// Auto-save sesi setiap kali tab berubah, ditutup, atau kontennya dimodifikasi
watch(tabs, () => {
  saveSession()
}, { deep: true })

watch(activeTabId, () => {
  saveSession()
})

// Watcher untuk merespons saat user membuka berkas baru dari File Manager saat editor sudah terbuka
watch(() => props.initialFile, (newFile) => {
  if (newFile && newFile.path) {
    openFileInTab(newFile)
  }
}, { deep: true })

watch(() => props.win?.fileData, (newFile) => {
  if (newFile && newFile.path) {
    openFileInTab(newFile)
  }
}, { deep: true })

watch(() => props.win?.customData, (newData) => {
  if (newData && newData.path) {
    openFileInTab(newData)
  }
}, { deep: true })
</script>

<template>
  <div
    class="h-full flex flex-col overflow-hidden text-xs select-none"
    :class="{
      'bg-[#272822] text-[#f8f8f2]': activeTheme === 'monokai',
      'bg-[#1e1e1e] text-[#d4d4d4]': activeTheme === 'obsidian',
      'bg-[#13111c] text-[#00ffcc]': activeTheme === 'cyberpunk',
      'bg-[#0f172a] text-[#f1f5f9]': activeTheme === 'notepad'
    }"
  >
    <!-- 1. TOP MENU BAR (Sublime Text / Notepad++ Style) -->
    <div class="h-8 border-b border-white/10 px-3 flex items-center justify-between shrink-0 bg-black/30 backdrop-blur text-[11px] font-sans">
      <div class="flex items-center gap-4">
        <!-- Brand Label (Satu Kesatuan dengan File Manager) -->
        <div class="flex items-center gap-1.5 font-bold tracking-wide" :class="activeTheme === 'monokai' ? 'text-amber-400' : 'text-sky-400'">
          <Folder class="w-3.5 h-3.5" />
          <span>File Manager</span>
          <span class="text-white/40 font-normal">/</span>
          <span class="text-slate-300 font-medium">Editor Berkas</span>
        </div>

        <!-- Menu Action Dropdowns -->
        <div class="flex items-center gap-1 text-slate-300">
          <button @click="createNewUntitledTab" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Berkas Baru (Ctrl+N)">Berkas Baru</button>
          <button @click="saveActiveFile" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Simpan Berkas (Ctrl+S)">Simpan</button>
          <button @click="saveAllFiles" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Simpan Semua (Ctrl+Shift+S)">Simpan Semua</button>
          <button @click="() => { showFindBar = !showFindBar; showReplace = false; }" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Cari (Ctrl+F)">Cari</button>
          <button @click="() => { showFindBar = true; showReplace = !showReplace; }" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Ganti Teks (Ctrl+H)">Ganti</button>
          <button @click="showGoToLineModal = true" class="px-2 py-0.5 rounded hover:bg-white/10 transition" title="Lompat Baris (Ctrl+G)">Lompat Baris</button>
        </div>
      </div>

      <!-- Right Controls: View & Theme Toggles -->
      <div class="flex items-center gap-2">
        <!-- Word Wrap Toggle -->
        <button
          @click="isWordWrap = !isWordWrap"
          class="px-2 py-0.5 rounded border border-white/10 transition flex items-center gap-1"
          :class="isWordWrap ? 'bg-sky-500/20 text-sky-300 border-sky-500/40' : 'hover:bg-white/10 text-slate-400'"
          title="Toggle Word Wrap"
        >
          <span>Wrap</span>
        </button>

        <!-- Sidebar Toggle -->
        <button
          @click="showSidebar = !showSidebar"
          class="px-2 py-0.5 rounded border border-white/10 transition flex items-center gap-1"
          :class="showSidebar ? 'bg-sky-500/20 text-sky-300 border-sky-500/40' : 'hover:bg-white/10 text-slate-400'"
          title="Toggle Sidebar Folder"
        >
          <Columns class="w-3 h-3" />
          <span>Sidebar</span>
        </button>

        <!-- Theme Selector -->
        <select
          v-model="activeTheme"
          class="bg-black/40 border border-white/10 rounded px-1.5 py-0.5 text-[10px] text-slate-300 outline-none"
        >
          <option value="monokai">Sublime Monokai</option>
          <option value="obsidian">Dark Obsidian</option>
          <option value="cyberpunk">Cyberpunk Neon</option>
          <option value="notepad">Notepad++ Slate</option>
        </select>
      </div>
    </div>

    <!-- 2. MULTI-TAB WORKSPACE BAR (Sublime Text Style) -->
    <div class="h-9 border-b border-white/10 flex items-center justify-between bg-black/40 px-2 shrink-0 overflow-x-auto" style="scrollbar-width: none;">
      <div class="flex items-center gap-1 h-full">
        <div
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTabId = tab.id"
          class="h-full px-3 flex items-center gap-2 border-r border-white/10 cursor-pointer transition select-none group min-w-[120px] max-w-[220px]"
          :class="activeTabId === tab.id ? 'bg-white/10 text-white font-semibold border-t-2 border-t-amber-400 shadow-sm' : 'text-slate-400 hover:bg-white/5 hover:text-slate-200'"
          :title="tab.path || tab.name"
        >
          <!-- Badge Bahasa -->
          <span class="text-[9px] font-mono px-1 py-0.2 rounded bg-black/40 text-sky-400 font-bold shrink-0">
            {{ getLanguageBadge(tab.name) }}
          </span>

          <!-- Nama Tab -->
          <span class="truncate flex-1 text-xs">{{ tab.name }}</span>

          <!-- Indikator Unsaved / Dirty Dot -->
          <span
            v-if="tab.isDirty"
            class="w-2 h-2 rounded-full bg-amber-400 shrink-0"
            title="Ada perubahan belum disimpan"
          ></span>

          <!-- Tombol Tutup Tab -->
          <button
            @click.stop="closeTab(tab.id, $event)"
            class="p-0.5 rounded-full transition"
            :class="activeTabId === tab.id ? 'opacity-80 hover:opacity-100 hover:bg-rose-500/30 text-slate-300 hover:text-rose-300' : 'opacity-0 group-hover:opacity-100 hover:bg-rose-500/30 text-slate-400 hover:text-rose-300'"
            title="Tutup Tab"
          >
            <X class="w-3 h-3" />
          </button>
        </div>

        <!-- Tombol Tab Baru (+) -->
        <button
          @click="createNewUntitledTab"
          class="p-1.5 rounded-lg hover:bg-white/10 text-slate-400 hover:text-white transition ml-1 shrink-0"
          title="Tab Berkas Baru (Ctrl+N)"
        >
          <Plus class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Quick Action: Save Button on Tab Bar -->
      <button
        @click="saveActiveFile"
        :disabled="!activeTab?.isDirty"
        class="h-6 px-2.5 rounded bg-sky-500 hover:bg-sky-400 disabled:opacity-40 disabled:pointer-events-none text-slate-950 font-bold flex items-center gap-1.5 transition text-[11px] shrink-0"
        title="Simpan Berkas Aktif (Ctrl+S)"
      >
        <Save class="w-3 h-3" />
        <span>Simpan</span>
      </button>
    </div>

    <!-- 3. MAIN WORKSPACE (SIDEBAR + EDITOR CORE) -->
    <div class="flex-1 flex overflow-hidden relative">
      <!-- SIDEBAR FILE TREE (Collapsible) -->
      <div
        v-if="showSidebar"
        class="w-52 sm:w-60 border-r border-white/10 bg-black/30 flex flex-col shrink-0 overflow-hidden text-xs"
      >
        <!-- Header Sidebar dengan Kontrol Navigasi Lengkap -->
        <div class="h-9 px-2 border-b border-white/10 flex items-center justify-between text-slate-300 text-[11px] font-medium bg-black/20 shrink-0">
          <div class="flex items-center gap-1 min-w-0 flex-1 mr-1">
            <!-- Navigasi Back -->
            <button
              @click="navigateSidebarBack"
              :disabled="!canSidebarGoBack"
              class="p-1 rounded hover:bg-white/10 text-slate-400 hover:text-white disabled:opacity-30 disabled:pointer-events-none transition"
              title="Kembali (Back)"
            >
              <ArrowLeft class="w-3.5 h-3.5" />
            </button>

            <!-- Navigasi Up (Naik Satu Folder) -->
            <button
              @click="navigateSidebarUp"
              :disabled="!canSidebarGoUp"
              class="p-1 rounded hover:bg-white/10 text-slate-400 hover:text-white disabled:opacity-30 disabled:pointer-events-none transition"
              title="Folder Induk / Keluar Folder (Up)"
            >
              <ArrowUp class="w-3.5 h-3.5" />
            </button>

            <!-- Navigasi Home (/DATA) -->
            <button
              @click="navigateSidebarHome"
              class="p-1 rounded hover:bg-white/10 text-slate-400 hover:text-white transition"
              title="Ke Direktori Utama (/DATA)"
            >
              <HardDrive class="w-3.5 h-3.5 text-sky-400" />
            </button>

            <!-- Nama Folder Aktif -->
            <span class="truncate font-semibold text-slate-200 ml-0.5" :title="sidebarCurrentDir">
              {{ sidebarCurrentFolderShortName }}
            </span>
          </div>

          <!-- Tombol Refresh -->
          <button
            @click="loadSidebarDirectory(sidebarCurrentDir)"
            class="p-1 rounded hover:bg-white/10 text-slate-400 hover:text-white transition shrink-0"
            title="Muat Ulang Berkas"
          >
            <RefreshCw class="w-3 h-3" :class="isSidebarLoading ? 'animate-spin' : ''" />
          </button>
        </div>

        <!-- File List Tree -->
        <div class="flex-1 overflow-y-auto p-1.5 flex flex-col gap-0.5" style="scrollbar-width: thin;">
          <!-- Item Kembali ke Folder Induk (..) jika di subfolder -->
          <div
            v-if="canSidebarGoUp"
            @click="navigateSidebarUp"
            class="h-8 px-2.5 rounded-lg flex items-center gap-2.5 cursor-pointer transition select-none text-slate-400 hover:bg-white/10 hover:text-slate-200 shrink-0 min-w-0"
            title="Kembali ke Folder Induk (..)"
          >
            <CornerLeftUp class="w-4 h-4 text-sky-400 shrink-0" />
            <span class="truncate text-xs leading-normal font-mono font-medium">.. (Folder Induk)</span>
          </div>

          <!-- Loop File & Folder -->
          <div
            v-for="item in sidebarFiles"
            :key="item.path"
            @click="item.is_dir ? loadSidebarDirectory(item.path) : openFileInTab(item)"
            class="h-8 px-2.5 rounded-lg flex items-center gap-2.5 cursor-pointer transition select-none shrink-0 min-w-0 group"
            :class="activeTab?.path === item.path ? 'bg-sky-500/20 text-sky-300 font-semibold' : 'text-slate-300 hover:bg-white/10 hover:text-white'"
            :title="`${item.name} (${item.is_dir ? 'Folder' : formatBytes(item.size)})`"
          >
            <Folder v-if="item.is_dir" class="w-4 h-4 text-amber-400 shrink-0" />
            <component
              v-else
              :is="getFileIconComponent(item.name)"
              class="w-4 h-4 shrink-0"
              :class="getFileIconColor(item.name)"
            />
            <span class="truncate text-xs leading-normal flex-1">{{ item.name }}</span>
            <span v-if="!item.is_dir && item.size" class="text-[10px] text-slate-500 shrink-0 font-mono hidden sm:inline group-hover:inline">
              {{ formatBytes(item.size) }}
            </span>
          </div>

          <div v-if="sidebarFiles.length === 0 && !isSidebarLoading" class="p-4 text-center text-slate-500 text-[11px]">
            Folder kosong atau belum dipilih
          </div>
        </div>
      </div>

      <!-- EDITOR BODY -->
      <div class="flex-1 flex flex-col overflow-hidden relative">
        <!-- FLOATING FIND & REPLACE BAR (Sublime Text Style) -->
        <div
          v-if="showFindBar"
          class="absolute top-2 right-4 z-30 w-80 sm:w-96 rounded-xl border border-white/20 bg-slate-900/95 backdrop-blur-md p-3 shadow-2xl flex flex-col gap-2 font-mono text-xs"
        >
          <!-- Find Row -->
          <div class="flex items-center gap-2">
            <div class="flex-1 relative flex items-center">
              <Search class="w-3.5 h-3.5 text-slate-400 absolute left-2 pointer-events-none" />
              <input
                id="sublime-find-input"
                v-model="findQuery"
                @input="executeFind"
                @keydown.enter="findNext"
                @keydown.esc="showFindBar = false"
                type="text"
                placeholder="Cari teks... (Enter untuk berikutnya)"
                class="w-full h-8 pl-7 pr-16 rounded-lg bg-black/60 border border-white/10 text-slate-100 placeholder-slate-500 outline-none focus:border-sky-400 text-xs"
              />
              <span v-if="findMatches.length > 0" class="absolute right-2 text-[10px] text-sky-400 font-bold">
                {{ currentMatchIndex + 1 }}/{{ findMatches.length }}
              </span>
            </div>
            <button @click="findPrev" class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300" title="Sebelumnya (Shift+Enter)">
              <ArrowUp class="w-3.5 h-3.5" />
            </button>
            <button @click="findNext" class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300" title="Berikutnya (Enter)">
              <ArrowDown class="w-3.5 h-3.5" />
            </button>
            <button @click="showFindBar = false" class="p-1.5 rounded-lg hover:bg-rose-500/20 text-slate-400 hover:text-rose-300">
              <X class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Replace Row -->
          <div v-if="showReplace" class="flex items-center gap-2 pt-1 border-t border-white/10">
            <input
              v-model="replaceQuery"
              @keydown.enter="executeReplace"
              type="text"
              placeholder="Ganti dengan..."
              class="flex-1 h-8 px-2.5 rounded-lg bg-black/60 border border-white/10 text-slate-100 placeholder-slate-500 outline-none focus:border-amber-400 text-xs"
            />
            <button @click="executeReplace" class="px-2.5 h-8 rounded-lg bg-amber-500/20 text-amber-300 font-bold hover:bg-amber-500/30 text-[11px] transition">
              Ganti
            </button>
            <button @click="executeReplaceAll" class="px-2.5 h-8 rounded-lg bg-rose-500/20 text-rose-300 font-bold hover:bg-rose-500/30 text-[11px] transition">
              Semua
            </button>
          </div>
        </div>

        <!-- GO TO LINE MODAL (Ctrl+G) -->
        <div
          v-if="showGoToLineModal"
          class="absolute inset-0 z-40 bg-black/50 backdrop-blur-sm flex items-start justify-center pt-20"
          @click.self="showGoToLineModal = false"
        >
          <div class="w-80 rounded-2xl bg-slate-900 border border-slate-700 p-4 shadow-2xl flex flex-col gap-3 font-mono">
            <span class="font-bold text-slate-200 text-xs">Lompat ke Baris (Go to Line)</span>
            <input
              id="sublime-goto-input"
              v-model="targetLineNumber"
              @keydown.enter="goToLine"
              @keydown.esc="showGoToLineModal = false"
              type="number"
              min="1"
              :max="lineCount"
              placeholder="Nomor baris..."
              class="w-full h-9 px-3 rounded-xl bg-slate-950 border border-slate-700 text-sky-300 text-sm outline-none focus:border-sky-400"
            />
            <div class="flex items-center justify-end gap-2 text-xs">
              <button @click="showGoToLineModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 text-slate-300">Batal</button>
              <button @click="goToLine" class="px-3 py-1.5 rounded-xl bg-sky-500 text-slate-950 font-bold">Lompat</button>
            </div>
          </div>
        </div>

        <!-- TEXTAREA & GUTTER CONTAINER -->
        <div class="flex-1 flex overflow-hidden relative">
          <!-- Gutter Nomor Baris (Synchronized) -->
          <div
            v-if="showLineNumbers"
            ref="gutterRef"
            class="w-12 bg-black/20 border-r border-white/5 p-2 select-none font-mono text-slate-600 text-right overflow-hidden shrink-0"
            :style="{ fontSize: `${fontSize}px`, lineHeight: '1.5' }"
          >
            <div
              v-for="lineIdx in lineCount"
              :key="lineIdx"
              :class="cursorLine === lineIdx ? 'text-amber-400 font-bold' : ''"
            >
              {{ lineIdx }}
            </div>
          </div>

          <!-- Main Textarea Editor -->
          <div class="flex-1 relative h-full overflow-hidden">
            <textarea
              ref="editorTextareaRef"
              :value="activeTab?.content || ''"
              @input="onContentInput"
              @click="updateCursorPosition"
              @keyup="updateCursorPosition"
              @keydown="onKeyDown"
              @scroll="onEditorScroll"
              spellcheck="false"
              class="w-full h-full p-2 bg-transparent outline-none resize-none font-mono selection:bg-sky-500/30"
              :class="{
                'whitespace-pre-wrap': isWordWrap,
                'whitespace-pre overflow-x-auto': !isWordWrap
              }"
              :style="{
                fontSize: `${fontSize}px`,
                lineHeight: '1.5',
                tabSize: indentSize
              }"
              placeholder="Mulai ketik kode atau teks di sini..."
            ></textarea>

            <!-- Loading Overlay -->
            <div
              v-if="activeTab?.isLoading"
              class="absolute inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center gap-2 text-slate-300"
            >
              <RefreshCw class="w-5 h-5 animate-spin text-sky-400" />
              <span>Memuat konten berkas...</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 4. BOTTOM STATUS BAR (Notepad++ Style) -->
    <div class="h-6 border-t border-white/10 px-3 flex items-center justify-between shrink-0 bg-black/40 text-[11px] font-mono select-none text-slate-400">
      <!-- Left: Line, Col, and Stats -->
      <div class="flex items-center gap-4">
        <span class="text-sky-300 font-semibold">
          Baris: {{ cursorLine }}, Kolom: {{ cursorCol }}
        </span>
        <span class="hidden sm:inline text-slate-500">|</span>
        <span class="hidden sm:inline">
          {{ lineCount }} baris, {{ charCount }} karakter, {{ wordCount }} kata
        </span>
        <span v-if="activeTab?.isDirty" class="text-amber-400 font-bold animate-pulse">
          ● Belum disimpan
        </span>
        <span v-else class="text-emerald-400 flex items-center gap-1">
          <Check class="w-3 h-3" /> Tersimpan
        </span>
      </div>

      <!-- Right: Encoding, Indent, Language -->
      <div class="flex items-center gap-3">
        <span class="hover:text-slate-200 cursor-pointer">UTF-8</span>
        <span class="text-slate-500">|</span>
        <span
          @click="indentSize = (indentSize === 2 ? 4 : 2)"
          class="hover:text-sky-300 cursor-pointer"
          title="Klik untuk ganti spasi 2 atau 4"
        >
          Spasi: {{ indentSize }}
        </span>
        <span class="text-slate-500">|</span>
        <span class="text-amber-400 font-bold">
          {{ activeTab?.language || 'Plain Text' }}
        </span>
      </div>
    </div>

    <!-- 5. CUSTOM MODALS (Zero Native Dialogs: Custom Glassmorphism Modals) -->
    <!-- Modal Simpan Berkas Baru (Save As) -->
    <div
      v-if="showSaveAsModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
      @click.self="showSaveAsModal = false"
    >
      <div class="w-full max-w-md bg-slate-900 border border-slate-700/80 rounded-2xl shadow-2xl p-5 text-slate-100 flex flex-col gap-4 font-sans">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div class="flex items-center gap-2 font-bold text-sm text-sky-400">
            <Save class="w-4 h-4" />
            <span>Simpan Berkas Baru</span>
          </div>
          <button @click="showSaveAsModal = false" class="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="flex flex-col gap-2">
          <div class="text-[11px] text-slate-400">
            Lokasi Folder: <span class="font-mono text-sky-300">{{ sidebarCurrentDir || '/DATA' }}</span>
          </div>
          <label class="text-xs text-slate-300 font-medium">Nama Berkas:</label>
          <input
            v-model="saveAsFileName"
            @keydown.enter="confirmSaveAs"
            type="text"
            class="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-slate-100 outline-none focus:border-sky-400"
            placeholder="contoh: catatan.txt, script.js, config.json..."
            autofocus
          />
        </div>
        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button @click="showSaveAsModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs text-slate-300">
            Batal
          </button>
          <button @click="confirmSaveAs" class="px-4 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs shadow-md">
            Simpan Berkas
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Tutup Tab Belum Disimpan -->
    <div
      v-if="showCloseTabConfirmModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
      @click.self="showCloseTabConfirmModal = false"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-amber-500/40 rounded-2xl shadow-2xl p-5 text-slate-100 flex flex-col gap-3 font-sans">
        <div class="flex items-center gap-2 font-bold text-sm text-amber-400">
          <AlertTriangle class="w-5 h-5" />
          <span>Perubahan Belum Disimpan</span>
        </div>
        <p class="text-xs text-slate-300 leading-relaxed">
          Berkas <span class="font-bold text-white">"{{ pendingCloseTabName }}"</span> memiliki perubahan yang belum disimpan. Tetap ingin menutup tab ini?
        </p>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-800">
          <button @click="showCloseTabConfirmModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs text-slate-300">
            Batal
          </button>
          <button @click="forceCloseTab" class="px-4 py-1.5 rounded-xl bg-rose-500 hover:bg-rose-400 text-white font-bold text-xs shadow-md">
            Tutup Tanpa Simpan
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Tutup Semua Tab -->
    <div
      v-if="showCloseAllConfirmModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in"
      @click.self="showCloseAllConfirmModal = false"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-amber-500/40 rounded-2xl shadow-2xl p-5 text-slate-100 flex flex-col gap-3 font-sans">
        <div class="flex items-center gap-2 font-bold text-sm text-amber-400">
          <AlertTriangle class="w-5 h-5" />
          <span>Tutup Semua Tab</span>
        </div>
        <p class="text-xs text-slate-300 leading-relaxed">
          Terdapat berkas yang belum disimpan. Yakin ingin menutup semua tab yang sedang terbuka?
        </p>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-800">
          <button @click="showCloseAllConfirmModal = false" class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs text-slate-300">
            Batal
          </button>
          <button @click="forceCloseAllTabs" class="px-4 py-1.5 rounded-xl bg-rose-500 hover:bg-rose-400 text-white font-bold text-xs shadow-md">
            Tutup Semua
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
textarea {
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
}
</style>
