<script setup>
/**
 * RAZOS Web Terminal View (xterm.js + FitAddon)
 * Updated: 2026-10-08
 * Deskripsi: Komponen terminal console performa tinggi berbasis xterm.js dan FitAddon.
 * Terhubung ke backend PTY shell atau mode simulasi interaktif responsif.
 */
import { ref, onMounted, onUnmounted } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

const terminalContainer = ref(null)
let term = null
let fitAddon = null
let socket = null

onMounted(() => {
  if (!terminalContainer.value) return

  // Inisialisasi Terminal
  term = new Terminal({
    cursorBlink: true,
    fontFamily: 'Consolas, "Fira Code", monospace',
    fontSize: 13,
    lineHeight: 1.2,
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
  fitAddon.fit()

  // Sambungkan ke WebSocket PTY
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${wsProtocol}//${window.location.host}/ws/terminal`

  try {
    socket = new WebSocket(wsUrl)
    socket.binaryType = 'arraybuffer'

    socket.onopen = () => {
      term.writeln('\x1b[1;36m[RAZOS Terminal]\x1b[0m Terhubung ke PTY Host Linux ARM64...')
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
      startInteractiveFallback()
    }

    socket.onclose = () => {
      term.writeln('\r\n\x1b[33m[RAZOS Terminal] Sesi PTY terputus. Beralih ke shell lokal.\x1b[0m')
      startInteractiveFallback()
    }

    term.onData((data) => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(data)
      }
    })
  } catch (err) {
    startInteractiveFallback()
  }

  window.addEventListener('resize', handleResize)
})

function handleResize() {
  if (fitAddon) {
    fitAddon.fit()
  }
}

// Fallback Shell Interaktif jika server WebSocket backend sedang offline
function startInteractiveFallback() {
  term.writeln('\r\n\x1b[1;32mRAZOS STB Shell (Kernel 6.1 aarch64)\x1b[0m')
  term.writeln('Ketik \x1b[33mhelp\x1b[0m untuk bantuan perintah sistem.\r\n')
  term.write('\x1b[1;36mrazos@armbian-s905x\x1b[0m:\x1b[1;34m~\x1b[0m$ ')

  let lineBuffer = ''
  term.onData((e) => {
    switch (e) {
      case '\r': // Enter
        term.writeln('')
        executeFallbackCommand(lineBuffer.trim())
        lineBuffer = ''
        term.write('\x1b[1;36mrazos@armbian-s905x\x1b[0m:\x1b[1;34m~\x1b[0m$ ')
        break
      case '\u007F': // Backspace
        if (lineBuffer.length > 0) {
          lineBuffer = lineBuffer.slice(0, -1)
          term.write('\b \b')
        }
        break
      default:
        if (e >= ' ' && e <= '~') {
          lineBuffer += e
          term.write(e)
        }
    }
  })
}

function executeFallbackCommand(cmd) {
  if (!cmd) return
  if (cmd === 'clear') {
    term.clear()
  } else if (cmd === 'help') {
    term.writeln('Perintah sistem tersedia: uname -a, free -h, df -h, docker ps, clear')
  } else if (cmd === 'uname -a') {
    term.writeln('Linux razos-s905x 6.1.0-armbian-arm64 #1 SMP PREEMPT aarch64 GNU/Linux')
  } else if (cmd === 'free -h') {
    term.writeln('               total        used        free      shared  buff/cache   available')
    term.writeln('Mem:           1.8Gi       420Mi       970Mi        12Mi       450Mi       1.3Gi')
    term.writeln('Swap:          900Mi          0B       900Mi')
  } else if (cmd === 'df -h') {
    term.writeln('Filesystem      Size  Used Avail Use% Mounted on')
    term.writeln('/dev/mmcblk0p1  7.2G  2.1G  4.8G  31% /')
    term.writeln('/dev/mmcblk1p1   29G  4.2G   24G  15% /DATA')
  } else if (cmd === 'docker ps') {
    term.writeln('CONTAINER ID   IMAGE          COMMAND       STATUS         PORTS')
    term.writeln('a1b2c3d4e5f6   razos-core     "./razos"     Up 2 hours     0.0.0.0:8080->8080/tcp')
  } else {
    term.writeln(`bash: ${cmd}: command not found`)
  }
}

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (socket) socket.close()
  if (term) term.dispose()
})
</script>

<template>
  <div class="h-full w-full bg-[#090d16] rounded-lg p-2 overflow-hidden flex flex-col">
    <div ref="terminalContainer" class="flex-1 w-full h-full overflow-hidden"></div>
  </div>
</template>
