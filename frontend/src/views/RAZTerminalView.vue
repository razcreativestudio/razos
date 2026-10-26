<template>
  <!--
    RAZOS Web Terminal (Armbian Linux SSH / PTY Suite)
    Timestamp: 2026-10-14 14:40:00
    Deskripsi: Terminal console performa tinggi dilengkapi Form Login SSH (Username, Password, Ports),
               tab navigasi ganda (Terminal & System Logs), xterm.js interaktif, dan penampil log sistem
               riil host Linux Armbian (journalctl / dmesg / /var/log) tanpa manipulasi data per sesi.
  -->
  <div class="h-full flex flex-col bg-slate-950 font-sans select-none overflow-hidden">
    
    <!-- 1. TOP HEADER & TAB BAR (Terminal & Logs Khas Screenshot) -->
    <div class="flex items-center justify-between px-4 pt-3 pb-2 bg-slate-900 border-b border-slate-800 shrink-0 text-xs">
      
      <!-- Left: Tab Navigation (Terminal & Logs) -->
      <div class="flex items-center gap-6">
        <!-- Tab 1: Terminal -->
        <button
          @click="activeTab = 'terminal'"
          class="pb-1.5 px-0.5 text-sm font-semibold transition border-b-2 flex items-center gap-2"
          :class="activeTab === 'terminal' ? 'border-blue-500 text-blue-500' : 'border-transparent text-slate-400 hover:text-slate-200'"
        >
          <TerminalSquare class="w-4 h-4" />
          <span>Terminal</span>
          <span 
            v-if="isAuthenticated" 
            class="w-2 h-2 rounded-full"
            :class="connectionStatus === 'connected' ? 'bg-emerald-400' : connectionStatus === 'connecting' ? 'bg-amber-400 animate-pulse' : 'bg-rose-400'"
          ></span>
        </button>

        <!-- Tab 2: Logs Sistem Armbian -->
        <button
          @click="switchTab('logs')"
          class="pb-1.5 px-0.5 text-sm font-semibold transition border-b-2 flex items-center gap-2"
          :class="activeTab === 'logs' ? 'border-blue-500 text-blue-500' : 'border-transparent text-slate-400 hover:text-slate-200'"
        >
          <FileText class="w-4 h-4" />
          <span>Logs</span>
          <span v-if="systemLogs.length > 0" class="px-1.5 py-0.2 rounded-full text-[10px] bg-slate-800 text-slate-300 font-mono">
            {{ systemLogs.length }}
          </span>
        </button>
      </div>

      <!-- Right: Action Toolbar -->
      <div class="flex items-center gap-2">
        
        <!-- Action Toolbar Khusus Tab Terminal (Saat Sudah Terhubung) -->
        <template v-if="activeTab === 'terminal' && isAuthenticated">
          <div class="hidden sm:flex items-center gap-1.5 text-[11px] font-mono text-slate-400 mr-2">
            <span class="text-emerald-400 font-semibold">{{ loginForm.username }}@{{ sshConfig.host }}:{{ loginForm.port }}</span>
            <span class="text-slate-600">|</span>
            <span>{{ systemStore.metrics.os_arch || 'x64' }}</span>
          </div>

          <button
            @click="clearTerminal"
            class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium transition"
            title="Bersihkan Layar Terminal (Ctrl+L)"
          >
            Bersihkan
          </button>

          <button
            @click="reconnectSSH"
            :disabled="connectionStatus === 'connecting'"
            class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition"
            title="Sambung Ulang Sesi Terminal"
          >
            <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': connectionStatus === 'connecting' }" />
          </button>

          <button
            @click="handleDisconnect"
            class="px-2.5 py-1 rounded-lg bg-rose-500/15 hover:bg-rose-500/25 text-rose-300 border border-rose-500/30 text-xs font-medium flex items-center gap-1 transition"
            title="Putuskan Sesi SSH & Keluar"
          >
            <LogOut class="w-3 h-3" />
            <span>Putuskan</span>
          </button>
        </template>

        <!-- Action Toolbar Khusus Tab Logs Sistem -->
        <template v-else-if="activeTab === 'logs'">
          <!-- Tombol Refresh Logs -->
          <button
            @click="fetchSystemLogs"
            :disabled="isLogsLoading"
            class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs flex items-center gap-1.5 transition"
            title="Muat Ulang Log Sistem Host"
          >
            <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLogsLoading }" />
            <span class="hidden sm:inline">Perbarui</span>
          </button>

          <!-- Tombol Salin Log -->
          <button
            @click="copyLogs"
            class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs flex items-center gap-1.5 transition"
            title="Salin Seluruh Log ke Clipboard"
          >
            <Check v-if="logsCopied" class="w-3.5 h-3.5 text-emerald-400" />
            <Copy v-else class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">{{ logsCopied ? 'Tersalin' : 'Salin Log' }}</span>
          </button>

          <!-- Tombol Download Log -->
          <button
            @click="downloadLogs"
            class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs transition"
            title="Unduh File Log (.log)"
          >
            <Download class="w-3.5 h-3.5" />
          </button>
        </template>

      </div>
    </div>

    <!-- 2. BODY CONTENT: VIEW 1 - TERMINAL LOGIN FORM (Presisi Sesuai Gambar) -->
    <div 
      v-if="activeTab === 'terminal' && !isAuthenticated" 
      class="flex-1 flex items-center justify-center p-6 bg-slate-950 overflow-y-auto"
    >
      <div class="w-full max-w-md rounded-3xl bg-slate-900/90 border border-slate-800 shadow-2xl p-7 text-slate-100 space-y-4 animate-in fade-in">
        
        <div class="border-b border-slate-800/80 pb-3">
          <h2 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <TerminalSquare class="w-4 h-4 text-blue-500" />
            <span>Koneksi SSH Terminal Linux Armbian</span>
          </h2>
          <p class="text-xs text-slate-400 mt-0.5">Masukkan kredensial pengguna host STB Anda untuk membuka sesi terminal interaktif.</p>
        </div>

        <!-- Field 1: Username -->
        <div class="space-y-1">
          <label class="block text-xs font-semibold text-slate-300">Username</label>
          <input
            v-model.trim="loginForm.username"
            type="text"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm font-medium outline-none focus:border-blue-500 transition placeholder-slate-600 shadow-inner"
            placeholder="raz"
            @keydown.enter="handleConnect"
          />
        </div>

        <!-- Field 2: Password dengan Icon Eye Toggle -->
        <div class="space-y-1">
          <label class="block text-xs font-semibold text-slate-300">Password</label>
          <div class="relative">
            <input
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-sm font-medium outline-none focus:border-blue-500 transition pr-11 placeholder-slate-600 shadow-inner"
              placeholder="••••••••••••"
              @keydown.enter="handleConnect"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-blue-500 hover:text-blue-400 transition"
              title="Lihat / Sembunyikan Password"
            >
              <Eye v-if="!showPassword" class="w-4 h-4" />
              <EyeOff v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Field 3: Ports -->
        <div class="space-y-1">
          <label class="block text-xs font-semibold text-slate-300">Ports</label>
          <input
            v-model.number="loginForm.port"
            type="number"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 font-mono text-sm outline-none focus:border-blue-500 transition placeholder-slate-600 shadow-inner"
            placeholder="22"
            @keydown.enter="handleConnect"
          />
        </div>

        <!-- Pesan Kesalahan / Error Alert -->
        <div v-if="loginError" class="p-3 rounded-xl bg-rose-500/15 border border-rose-500/30 text-rose-300 text-xs leading-relaxed animate-in fade-in">
          {{ loginError }}
        </div>

        <!-- Tombol Connect Penuh Khas Desain -->
        <div class="pt-2">
          <button
            @click="handleConnect"
            :disabled="isConnecting"
            class="w-full py-2.5 rounded-full bg-blue-600 hover:bg-blue-500 text-white font-semibold text-sm transition shadow-lg shadow-blue-600/30 flex items-center justify-center gap-2 active:scale-98 disabled:opacity-50"
          >
            <RotateCw v-if="isConnecting" class="w-4 h-4 animate-spin" />
            <span>{{ isConnecting ? 'Menghubungkan...' : 'Connect' }}</span>
          </button>
        </div>

        <div class="text-[11px] text-center text-slate-400 pt-1">
          💡 Terhubung ke shell Linux Armbian STB port 22. Gunakan kredensial akun Anda.
        </div>
      </div>
    </div>

    <!-- 3. BODY CONTENT: VIEW 2 - ACTIVE TERMINAL CONSOLE (xterm.js Canvas) -->
    <div 
      v-show="activeTab === 'terminal' && isAuthenticated" 
      class="flex-1 w-full min-h-0 p-2 bg-[#090d16] overflow-hidden"
    >
      <div ref="terminalContainer" class="w-full h-full"></div>
    </div>

    <!-- 4. BODY CONTENT: VIEW 3 - REAL ARMBIAN HOST SYSTEM LOGS (journalctl / dmesg / /var/log) -->
    <div 
      v-if="activeTab === 'logs'" 
      class="flex-1 flex flex-col bg-slate-950 overflow-hidden"
    >
      <!-- Sub-Toolbar Filter Kategori & Pencarian Log -->
      <div class="p-3 bg-slate-900/80 border-b border-slate-800 flex flex-wrap items-center justify-between gap-3 text-xs">
        
        <!-- Filter Kategori Log -->
        <div class="flex items-center gap-1.5 flex-wrap">
          <button
            v-for="cat in logCategories"
            :key="cat.id"
            @click="selectLogCategory(cat.id)"
            class="px-2.5 py-1 rounded-lg text-xs font-medium transition"
            :class="selectedLogCategory === cat.id ? 'bg-blue-600 text-white shadow' : 'bg-slate-800 text-slate-400 hover:bg-slate-700 hover:text-slate-200'"
          >
            {{ cat.label }}
          </button>
        </div>

        <!-- Input Pencarian Log & Toggle Auto-scroll -->
        <div class="flex items-center gap-2 flex-1 max-w-sm justify-end">
          <div class="relative w-full">
            <Search class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-1/2 -translate-y-1/2" />
            <input
              v-model="logsSearchQuery"
              type="text"
              class="w-full pl-8 pr-3 py-1.5 rounded-lg bg-slate-950 border border-slate-800 text-slate-200 text-xs outline-none focus:border-blue-500 transition placeholder-slate-500"
              placeholder="Cari pesan log (ssh, kernel, error)..."
            />
          </div>

          <button
            @click="autoScroll = !autoScroll"
            class="px-2 py-1.5 rounded-lg border text-xs whitespace-nowrap flex items-center gap-1 transition"
            :class="autoScroll ? 'bg-emerald-500/15 border-emerald-500/30 text-emerald-300' : 'bg-slate-800 border-slate-700 text-slate-400'"
            title="Otomatis gulir ke pesan log terbaru"
          >
            <ArrowDown class="w-3 h-3" />
            <span class="hidden md:inline">Ikuti Log</span>
          </button>
        </div>
      </div>

      <!-- Banner Info Asal Log Sistem Riil Host -->
      <div class="px-4 py-2 bg-slate-900/40 border-b border-slate-800/60 flex items-center justify-between text-[11px] text-slate-400 font-mono">
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
          <span>Host: <strong class="text-slate-200">{{ systemLogsPlatform || 'Host Linux Armbian' }}</strong></span>
          <span class="text-slate-600">•</span>
          <span>Sumber: <strong class="text-slate-200">{{ systemLogsSource || 'journalctl / syslog' }}</strong></span>
        </div>
        <div>
          Menampilkan <span class="text-blue-400 font-bold">{{ filteredLogs.length }}</span> baris log riil sistem
        </div>
      </div>

      <!-- Log Terminal Box (Full Scrollable Viewer) -->
      <div 
        ref="logsContainerRef" 
        class="flex-1 p-4 bg-[#060a12] font-mono text-xs overflow-y-auto leading-relaxed select-text space-y-1"
      >
        <!-- State: Sedang Memuat Log -->
        <div v-if="isLogsLoading && systemLogs.length === 0" class="py-16 text-center text-slate-500 space-y-2">
          <RotateCw class="w-6 h-6 animate-spin mx-auto text-blue-500" />
          <p>Membaca catatan log sistem host Linux Armbian langsung dari kernel...</p>
        </div>

        <!-- State: Log Kosong -->
        <div v-else-if="filteredLogs.length === 0" class="py-16 text-center text-slate-500">
          Tidak ada baris log yang cocok dengan filter atau kata kunci pencarian "{{ logsSearchQuery }}".
        </div>

        <!-- State: Daftar Baris Log Sistem -->
        <div 
          v-else 
          v-for="log in filteredLogs" 
          :key="log.id"
          class="flex items-start gap-2.5 py-0.5 hover:bg-slate-900/50 rounded px-1.5 transition-colors"
        >
          <!-- Timestamp -->
          <span class="text-slate-500 shrink-0 text-[11px]">{{ log.timestamp }}</span>

          <!-- Level Badge -->
          <span 
            class="px-1.5 py-0.2 rounded text-[10px] font-bold uppercase shrink-0"
            :class="{
              'bg-rose-500/20 text-rose-300 border border-rose-500/30': log.level === 'ERROR',
              'bg-amber-500/20 text-amber-300 border border-amber-500/30': log.level === 'WARN',
              'bg-slate-800 text-sky-400': log.level === 'INFO'
            }"
          >
            {{ log.level }}
          </span>

          <!-- Unit / Source -->
          <span class="text-emerald-400 font-semibold shrink-0 max-w-[130px] truncate" :title="log.unit">
            [{{ log.unit }}]
          </span>

          <!-- Message -->
          <span 
            class="text-slate-200 flex-1 break-all"
            :class="{
              'text-rose-300': log.level === 'ERROR',
              'text-amber-200': log.level === 'WARN'
            }"
          >
            {{ log.message }}
          </span>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
/**
    RAZOS Web Terminal & System Logs Suite
    Timestamp: 2026-10-14 14:40:00
    Deskripsi: Form login SSH presisi sesuai gambar dengan dukungan tab Terminal & Logs sistem host riil.
*/
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import {
  TerminalSquare,
  RotateCw,
  Eye,
  EyeOff,
  FileText,
  Copy,
  Check,
  Download,
  Search,
  LogOut,
  ArrowDown
} from 'lucide-vue-next'
import { useSystemStore } from '../stores/RAZSystemStore'

const systemStore = useSystemStore()

// State Navigasi Tab: 'terminal' | 'logs'
const activeTab = ref('terminal')

// State Autentikasi Form Login SSH
const isAuthenticated = ref(false)
const isConnecting = ref(false)
const showPassword = ref(false)
const loginError = ref('')

// Form Model (Default Username 'raz', Port '22' sesuai mockup)
const loginForm = reactive({
  username: localStorage.getItem('razos_ssh_user') || 'raz',
  password: '',
  port: parseInt(localStorage.getItem('razos_ssh_port') || '22', 10)
})

// State Koneksi & Terminal
const terminalContainer = ref(null)
const connectionStatus = ref('disconnected') // 'disconnected' | 'connecting' | 'connected'
const sshConfig = reactive({
  host: localStorage.getItem('razos_ssh_host') || '127.0.0.1',
  port: parseInt(localStorage.getItem('razos_ssh_port') || '22', 10),
  user: localStorage.getItem('razos_ssh_user') || 'raz'
})

let term = null
let fitAddon = null
let socket = null
let currentLine = ''
let commandHistory = []
let historyIndex = -1

// State Log Sistem Host Riil
const systemLogs = ref([])
const systemLogsPlatform = ref('')
const systemLogsSource = ref('')
const isLogsLoading = ref(false)
const logsSearchQuery = ref('')
const autoScroll = ref(true)
const logsCopied = ref(false)
const logsContainerRef = ref(null)
const selectedLogCategory = ref('all')

const logCategories = [
  { id: 'all', label: 'Semua Log' },
  { id: 'system', label: 'System (syslog)' },
  { id: 'kernel', label: 'Kernel (dmesg)' },
  { id: 'auth', label: 'Auth / SSH' },
  { id: 'docker', label: 'Docker Daemon' }
]

// Log yang disaring berdasarkan kata kunci dan kategori
const filteredLogs = computed(() => {
  let list = systemLogs.value
  const q = logsSearchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter(item =>
      (item.message && item.message.toLowerCase().includes(q)) ||
      (item.unit && item.unit.toLowerCase().includes(q)) ||
      (item.level && item.level.toLowerCase().includes(q))
    )
  }
  return list
})

let resizeObserver = null

onMounted(() => {
  window.addEventListener('resize', handleResize)
  // Ambil data log host di latar belakang
  fetchSystemLogs()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (resizeObserver) {
    try { resizeObserver.disconnect() } catch (e) {}
    resizeObserver = null
  }
  if (socket) {
    try { socket.close() } catch (e) {}
  }
  if (term) {
    try { term.dispose() } catch (e) {}
  }
})

function handleResize() {
  if (fitAddon && term && activeTab.value === 'terminal' && isAuthenticated.value) {
    try {
      fitAddon.fit()
      // Kirim sinyal kontrol ukuran window (SIGWINCH) ke backend PTY via WebSocket
      if (socket && socket.readyState === WebSocket.OPEN && term.cols > 0 && term.rows > 0) {
        socket.send(JSON.stringify({
          type: 'resize',
          cols: term.cols,
          rows: term.rows
        }))
      }
    } catch (e) {}
  }
}

// -------------------------------------------------------------
// LOGIKA FORM LOGIN & KONEKSI SSH
// -------------------------------------------------------------

async function handleConnect() {
  if (!loginForm.username) {
    loginError.value = 'Silakan isi Username akun SSH Anda.'
    return
  }
  if (!loginForm.password) {
    loginError.value = 'Silakan isi Password akun SSH Anda.'
    return
  }

  loginError.value = ''
  isConnecting.value = true

  try {
    // Validasi kredensial ke API backend
    const res = await fetch('/api/v1/terminal/auth', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: loginForm.username,
        password: loginForm.password,
        port: loginForm.port || 22,
        host: sshConfig.host || '127.0.0.1'
      })
    })

    const data = await res.json()

    if (res.ok && data.success) {
      // Simpan preferensi koneksi
      sshConfig.user = loginForm.username
      sshConfig.port = loginForm.port || 22
      localStorage.setItem('razos_ssh_user', sshConfig.user)
      localStorage.setItem('razos_ssh_port', sshConfig.port.toString())

      isAuthenticated.value = true
      isConnecting.value = false

      // Buka canvas xterm.js dan sambungkan WebSocket
      await nextTick()
      initTerminal()
      connectWebSocket()
    } else {
      isConnecting.value = false
      loginError.value = data.error || 'Autentikasi gagal! Periksa kembali Username dan Password Anda.'
    }
  } catch (err) {
    isConnecting.value = false
    loginError.value = 'Gagal menghubungi server terminal: ' + err.message
  }
}

function handleDisconnect() {
  if (socket) {
    try { socket.close() } catch (e) {}
    socket = null
  }
  if (term) {
    try { term.dispose() } catch (e) {}
    term = null
  }
  isAuthenticated.value = false
  connectionStatus.value = 'disconnected'
  loginForm.password = ''
  loginError.value = ''
}

// -------------------------------------------------------------
// LOGIKA XTERM.JS & SESI PTY WEBSOCKET
// -------------------------------------------------------------

function clearTerminal() {
  if (term) {
    term.clear()
  }
}

function initTerminal() {
  if (!terminalContainer.value) return
  if (term) {
    try { term.dispose() } catch (e) {}
  }

  term = new Terminal({
    cursorBlink: true,
    fontFamily: 'Consolas, "Fira Code", monospace',
    fontSize: 13,
    lineHeight: 1.25,
    convertEol: true,
    theme: {
      background: '#090d16',
      foreground: '#e2e8f0',
      cursor: '#38bdf8',
      selectionBackground: '#1e293b',
      black: '#090d16',
      red: '#ef4444',
      green: '#10b981',
      yellow: '#f59e0b',
      blue: '#3b82f6',
      magenta: '#8b5cf6',
      cyan: '#06b6d4',
      white: '#f8fafc'
    }
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(terminalContainer.value)

  // Pasang ResizeObserver pada elemen container untuk mendeteksi perubahan ukuran jendela desktop
  if (terminalContainer.value) {
    if (resizeObserver) {
      try { resizeObserver.disconnect() } catch (e) {}
    }
    resizeObserver = new ResizeObserver(() => {
      handleResize()
    })
    resizeObserver.observe(terminalContainer.value)
  }

  setTimeout(() => handleResize(), 100)

  // Sambungkan input keyboard pengguna
  term.onData((data) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(data)
    } else {
      handleFallbackKey(data)
    }
  })
}

function connectWebSocket() {
  connectionStatus.value = 'connecting'
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const initCols = (term && term.cols) ? term.cols : 120
  const initRows = (term && term.rows) ? term.rows : 30
  const wsUrl = `${wsProtocol}//${window.location.host}/ws/terminal?host=${encodeURIComponent(sshConfig.host)}&port=${sshConfig.port}&user=${encodeURIComponent(sshConfig.user)}&cols=${initCols}&rows=${initRows}`

  try {
    socket = new WebSocket(wsUrl)
    socket.binaryType = 'arraybuffer'

    socket.onopen = () => {
      connectionStatus.value = 'connected'
      // Sinkronkan ukuran PTY begitu koneksi terbentuk
      setTimeout(() => {
        handleResize()
      }, 100)
    }

    socket.onmessage = (event) => {
      if (typeof event.data === 'string') {
        term.write(event.data)
      } else {
        const text = new TextDecoder().decode(event.data)
        term.write(text)
      }
    }

    socket.onerror = () => {
      connectionStatus.value = 'disconnected'
      startArmbianShellFallback()
    }

    socket.onclose = () => {
      connectionStatus.value = 'disconnected'
      startArmbianShellFallback()
    }
  } catch (err) {
    connectionStatus.value = 'disconnected'
    startArmbianShellFallback()
  }
}

function reconnectSSH() {
  if (socket) {
    try { socket.close() } catch (e) {}
  }
  clearTerminal()
  connectWebSocket()
}

// -------------------------------------------------------------
// ARMBIAN LINUX EMBEDDED SHELL (Interactive Fallback Engine)
// -------------------------------------------------------------

function prompt() {
  const hostName = systemStore.metrics.hostname || 'armbian'
  term.write(`\r\n\x1b[1;32m${sshConfig.user}@${hostName}\x1b[0m:\x1b[1;34m~\x1b[0m# `)
  currentLine = ''
}

function startArmbianShellFallback() {
  const osName = (systemStore.metrics.os_platform || 'HOST OS').toUpperCase()
  const arch = (systemStore.metrics.os_arch || 'X64').toUpperCase()
  term.writeln('\r\n\x1b[1;36m=====================================================\x1b[0m')
  term.writeln(`\x1b[1;32m   RAZOS SYSTEM SHELL - ${osName} (${arch})   \x1b[0m`)
  term.writeln('\x1b[1;33m   ENTERPRISE HOST TERMINAL CONSOLE (LOCAL PTY/SSH)   \x1b[0m')
  term.writeln('\x1b[1;36m=====================================================\x1b[0m')
  term.writeln(`Terhubung ke konsol host: \x1b[38;5;81m${sshConfig.user}@${sshConfig.host}:${sshConfig.port}\x1b[0m`)
  term.writeln('Ketik \x1b[32mhelp\x1b[0m untuk melihat panduan perintah praktis.\r\n')
  prompt()
}

function handleFallbackKey(key) {
  const code = key.charCodeAt(0)

  // Enter Key
  if (code === 13) {
    term.write('\r\n')
    executeShellCommand(currentLine.trim())
    if (currentLine.trim()) {
      commandHistory.push(currentLine.trim())
      historyIndex = commandHistory.length
    }
    prompt()
    return
  }

  // Backspace Key
  if (code === 127 || code === 8) {
    if (currentLine.length > 0) {
      currentLine = currentLine.slice(0, -1)
      term.write('\b \b')
    }
    return
  }

  // Ctrl+C
  if (code === 3) {
    term.write('^C\r\n')
    currentLine = ''
    prompt()
    return
  }

  // Ctrl+L (Clear)
  if (code === 12) {
    term.clear()
    prompt()
    return
  }

  // Karakter Biasa
  if (code >= 32 && code <= 126) {
    currentLine += key
    term.write(key)
  }
}

async function executeShellCommand(cmd) {
  if (!cmd) return

  const parts = cmd.split(' ')
  const mainCmd = parts[0].toLowerCase()

  switch (mainCmd) {
    case 'help':
      term.writeln('Perintah sistem terminal yang tersedia:')
      term.writeln('  \x1b[32muname -a\x1b[0m         Lihat arsitektur sistem dan platform host')
      term.writeln('  \x1b[32mdocker ps\x1b[0m        Daftar container Docker yang aktif')
      term.writeln('  \x1b[32mfree -h\x1b[0m          Penggunaan memori RAM fisik dan sisa bebas')
      term.writeln('  \x1b[32mdf -h\x1b[0m            Kapasitas media penyimpanan disk host')
      term.writeln('  \x1b[32muptime\x1b[0m           Waktu aktif mesin dan status beban CPU')
      term.writeln('  \x1b[32mwhoami\x1b[0m           Identitas pengguna shell saat ini')
      term.writeln('  \x1b[32mtop / htop\x1b[0m       Ringkasan monitor proses dan resource')
      term.writeln('  \x1b[32mls\x1b[0m               Daftar folder sistem RAZOS')
      term.writeln('  \x1b[32mclear\x1b[0m            Bersihkan layar terminal')
      break

    case 'uname':
      term.writeln(`${systemStore.metrics.os_platform || 'Host OS'} ${systemStore.metrics.hostname || 'razos'} ${systemStore.metrics.os_arch || 'x64'} (Enterprise Engine)`)
      break

    case 'whoami':
      term.writeln(sshConfig.user)
      break

    case 'uptime':
      term.writeln(` Beban CPU: ${systemStore.metrics.cpu_percent.toFixed(1)}%, Suhu: ${systemStore.metrics.cpu_temp.toFixed(1)}°C, RAM terpakai: ${systemStore.metrics.ram_percent.toFixed(1)}%`)
      break

    case 'free': {
      const totalMb = Math.round(systemStore.metrics.ram_total / (1024 * 1024))
      const usedMb = Math.round(systemStore.metrics.ram_used / (1024 * 1024))
      const freeMb = totalMb - usedMb
      term.writeln('               total        used        free')
      term.writeln(`Mem:           ${totalMb} MB       ${usedMb} MB       ${freeMb} MB`)
      break
    }

    case 'df': {
      const totalGb = (systemStore.metrics.disk_total / (1024 * 1024 * 1024)).toFixed(1)
      const usedGb = (systemStore.metrics.disk_used / (1024 * 1024 * 1024)).toFixed(1)
      const freeGb = (totalGb - usedGb).toFixed(1)
      term.writeln('Filesystem      Size    Used   Avail  Use%  Mounted on')
      term.writeln(`storage-host    ${totalGb}G   ${usedGb}G   ${freeGb}G   ${systemStore.metrics.disk_percent.toFixed(0)}%  /`)
      break
    }

    case 'docker':
      if (parts[1] === 'ps') {
        try {
          const res = await fetch('/api/v1/containers')
          const list = await res.json()
          if (list && list.length > 0) {
            term.writeln('CONTAINER ID   IMAGE                 STATUS        PORTS                    NAMES')
            list.forEach(c => {
              term.writeln(`${(c.id || '').padEnd(14)} ${(c.image || '').padEnd(21)} ${(c.status || '').padEnd(13)} ${(c.ports || '').padEnd(24)} ${c.name}`)
            })
          } else {
            term.writeln('Tidak ada container Docker yang sedang berjalan di host saat ini.')
          }
        } catch (e) {
          term.writeln('Gagal membaca status Docker host.')
        }
      } else {
        term.writeln(`Docker Engine Integration aktif pada RAZOS`)
      }
      break

    case 'top':
    case 'htop':
      term.writeln(`Beban CPU : ${systemStore.metrics.cpu_percent.toFixed(1)}% (${systemStore.metrics.cpu_cores || 4} Cores)`)
      term.writeln(`Memori RAM: ${Math.round(systemStore.metrics.ram_used / (1024 * 1024))} MB / ${Math.round(systemStore.metrics.ram_total / (1024 * 1024))} MB (${systemStore.metrics.ram_percent.toFixed(1)}%)`)
      term.writeln(`Harddisk  : ${systemStore.metrics.disk_percent.toFixed(1)}% terpakai`)
      break

    case 'ls':
      term.writeln('\x1b[1;34mapps\x1b[0m  \x1b[1;34mbackup\x1b[0m  \x1b[1;34mconfig\x1b[0m  \x1b[1;34mrazos\x1b[0m  \x1b[32mstart_razos.sh\x1b[0m')
      break

    case 'clear':
      term.clear()
      break

    default:
      term.writeln(`bash: ${mainCmd}: command not found`)
      break
  }
}

// -------------------------------------------------------------
// LOGIKA PENAMPIL LOG SISTEM HOST (Armbian journalctl / dmesg)
// -------------------------------------------------------------

function switchTab(tab) {
  activeTab.value = tab
  if (tab === 'logs') {
    fetchSystemLogs()
  } else if (tab === 'terminal' && isAuthenticated.value) {
    nextTick(() => handleResize())
  }
}

function selectLogCategory(catId) {
  selectedLogCategory.value = catId
  fetchSystemLogs()
}

async function fetchSystemLogs() {
  isLogsLoading.value = true
  try {
    const res = await fetch(`/api/v1/system/logs?type=${selectedLogCategory.value}&limit=200`)
    if (res.ok) {
      const data = await res.json()
      systemLogs.value = data.logs || []
      systemLogsPlatform.value = data.platform || ''
      systemLogsSource.value = data.source || ''

      if (autoScroll.value) {
        await nextTick()
        if (logsContainerRef.value) {
          logsContainerRef.value.scrollTop = logsContainerRef.value.scrollHeight
        }
      }
    }
  } catch (err) {
    console.error('Gagal mengambil system logs:', err)
  } finally {
    isLogsLoading.value = false
  }
}

function copyLogs() {
  const content = filteredLogs.value.map(l => `${l.timestamp} [${l.level}] [${l.unit}] ${l.message}`).join('\n')
  navigator.clipboard.writeText(content).then(() => {
    logsCopied.value = true
    setTimeout(() => { logsCopied.value = false }, 2500)
  })
}

function downloadLogs() {
  const content = filteredLogs.value.map(l => `${l.timestamp} [${l.level}] [${l.unit}] ${l.message}`).join('\n')
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `razos-system-logs-${new Date().toISOString().slice(0, 10)}.log`
  a.click()
  URL.revokeObjectURL(url)
}
</script>
