<script setup>
/**
 * RAZOS Web Terminal View
 * Updated: 2026-10-08
 * Deskripsi: Antarmuka Web Terminal bergaya dark slate CLI.
 */
import { ref } from 'vue'

const terminalLines = ref([
  'RAZOS Linux Kernel 6.1.0-arm64-s905x (tty1)',
  'Selamat Datang di Shell Web Terminal RAZOS.',
  'Ketik perintah untuk berinteraksi dengan sistem host.',
  '-------------------------------------------------------',
  'razos@armbian-s905x:~$ uname -a',
  'Linux razos-s905x 6.1.0-armbian-arm64 #1 SMP PREEMPT aarch64 GNU/Linux',
  'razos@armbian-s905x:~$ free -h',
  '               total        used        free      shared  buff/cache   available',
  'Mem:           1.8Gi       410Mi       980Mi        12Mi       450Mi       1.3Gi',
  'Swap:          900Mi          0B       900Mi',
  'razos@armbian-s905x:~$'
])

const inputCommand = ref('')

function handleCommand() {
  if (!inputCommand.value.trim()) return
  terminalLines.value.push(`razos@armbian-s905x:~$ ${inputCommand.value}`)
  if (inputCommand.value.trim() === 'clear') {
    terminalLines.value = ['razos@armbian-s905x:~$']
  } else if (inputCommand.value.trim() === 'help') {
    terminalLines.value.push('Perintah dasar: uname -a, free -h, df -h, docker ps, clear')
  } else {
    terminalLines.value.push(`bash: ${inputCommand.value}: modul WebSocket PTY live akan aktif pada Fase 2`)
  }
  inputCommand.value = ''
}
</script>

<template>
  <div class="h-full flex flex-col font-mono text-xs sm:text-sm bg-black/80 rounded-lg p-3 text-emerald-400 overflow-hidden">
    <!-- Terminal Output -->
    <div class="flex-1 overflow-auto space-y-1">
      <div v-for="(line, idx) in terminalLines" :key="idx" class="leading-relaxed">
        {{ line }}
      </div>
    </div>

    <!-- Terminal Command Prompt Input -->
    <div class="flex items-center gap-2 mt-2 pt-2 border-t border-slate-800">
      <span class="text-sky-400">razos@armbian-s905x:~$</span>
      <input
        v-model="inputCommand"
        @keyup.enter="handleCommand"
        type="text"
        autofocus
        class="flex-1 bg-transparent outline-none text-slate-100 font-mono"
        placeholder="ketik perintah..."
      />
    </div>
  </div>
</template>
