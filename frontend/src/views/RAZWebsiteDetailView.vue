<script setup>
/**
 * RAZOS Website Detail Window Component
 * Created: 2026-10-13 (v2.6.0 - Standalone Window Architecture)
 * Deskripsi: Rincian & Kredensial website mandiri.
 * Berjalan sebagai Jendela OS terpisah (tidak memblokir/mengganggu Website Manager).
 */
import { ref, computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Globe, Server, Database, GitBranch, Copy, Check,
  ExternalLink, Pencil, FolderOpen
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: true
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const site = computed(() => props.win.customData || {})
const copiedField = ref(null)

function copyToClipboard(text, fieldName) {
  if (!text) return
  navigator.clipboard.writeText(text).then(() => {
    copiedField.value = fieldName
    systemStore.showNotification('Disalin', `${fieldName} berhasil disalin ke clipboard`, 'info')
    setTimeout(() => {
      copiedField.value = null
    }, 2000)
  })
}

function openEditThisSite() {
  windowStore.openWebsiteEditor(site.value)
}

function openGitThisSite() {
  windowStore.openWebsiteGit(site.value)
}

function openRootFolderInFiles() {
  windowStore.openApp('file-manager')
}
</script>

<template>
  <div class="h-full flex flex-col justify-between text-xs select-none gap-3">
    <!-- Body Content -->
    <div class="flex-1 overflow-y-auto space-y-3 pr-1 text-slate-300">
      <!-- Container Overview Grid -->
      <div class="grid grid-cols-2 gap-2 bg-slate-950 p-3 rounded-xl border border-slate-800">
        <div>
          <span class="text-slate-500 block text-[10px]">CONTAINER NAME</span>
          <span class="font-mono text-slate-200">{{ site.container_name }}</span>
        </div>
        <div>
          <span class="text-slate-500 block text-[10px]">HOST PORT</span>
          <span class="font-mono text-amber-400 font-bold">:{{ site.port }}</span>
        </div>
        <div>
          <span class="text-slate-500 block text-[10px]">RUNTIME STACK</span>
          <span class="uppercase font-semibold text-sky-400">{{ site.runtime }}</span>
        </div>
        <div>
          <span class="text-slate-500 block text-[10px]">STATUS</span>
          <span class="uppercase font-bold text-emerald-400">{{ site.status }}</span>
        </div>
      </div>

      <!-- Root Folder Copy & Open -->
      <div class="bg-slate-950 p-3 rounded-xl border border-slate-800 flex items-center justify-between">
        <div class="min-w-0 pr-2">
          <span class="text-slate-500 block text-[10px]">DOCUMENT ROOT FOLDER</span>
          <span class="font-mono text-[11px] text-slate-200 truncate block">{{ site.root_folder }}</span>
        </div>
        <div class="flex items-center gap-1.5 shrink-0">
          <button
            @click="copyToClipboard(site.root_folder, 'Folder Path')"
            class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
            title="Salin Path Folder"
          >
            <Check v-if="copiedField === 'Folder Path'" class="w-3.5 h-3.5 text-emerald-400" />
            <Copy v-else class="w-3.5 h-3.5" />
          </button>
          <button
            @click="openRootFolderInFiles"
            class="px-2 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] font-medium flex items-center gap-1 transition"
            title="Buka di File Manager"
          >
            <FolderOpen class="w-3.5 h-3.5" />
            <span>Files</span>
          </button>
        </div>
      </div>

      <!-- Database Details jika ada -->
      <div v-if="site.db_name" class="bg-amber-500/10 border border-amber-500/20 p-3 rounded-xl space-y-1.5 text-[11px]">
        <div class="font-bold text-amber-400 flex items-center gap-1.5">
          <Database class="w-3.5 h-3.5" />
          <span>Kredensial Database Terhubung</span>
        </div>
        <div class="grid grid-cols-2 gap-2 text-slate-300 font-mono">
          <div>Host: {{ site.db_host }}:{{ site.db_port }}</div>
          <div>Database: {{ site.db_name }}</div>
          <div>User: {{ site.db_user }}</div>
          <div>Engine: {{ site.db_engine || 'mariadb' }}</div>
        </div>
      </div>

      <!-- Git Sync Details jika ada -->
      <div v-if="site.git_repo" class="bg-sky-500/10 border border-sky-500/20 p-3 rounded-xl space-y-1.5 text-[11px]">
        <div class="font-bold text-sky-400 flex items-center justify-between">
          <span class="flex items-center gap-1.5">
            <GitBranch class="w-3.5 h-3.5" />
            <span>GitHub Repository Terhubung</span>
          </span>
          <button
            @click="openGitThisSite"
            class="text-[10px] text-sky-300 hover:underline flex items-center gap-1"
          >
            <span>Kelola Git</span>
          </button>
        </div>
        <div class="space-y-1 text-slate-300 font-mono text-[10px]">
          <div class="truncate">Repo: {{ site.git_repo }}</div>
          <div>Branch: {{ site.git_branch || 'main' }}</div>
          <div v-if="site.git_last_commit">Commit: {{ site.git_last_commit }}</div>
          <div v-if="site.git_last_pull">Last Pull: {{ site.git_last_pull }}</div>
        </div>
      </div>
    </div>

    <!-- Window Footer Actions -->
    <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
      <button
        @click="windowStore.closeWindow(win.id)"
        class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Tutup
      </button>

      <div class="flex items-center gap-2">
        <button
          @click="openEditThisSite"
          class="px-3.5 py-2 rounded-xl bg-amber-500/15 hover:bg-amber-500/25 border border-amber-500/30 text-amber-300 text-xs font-semibold flex items-center gap-1.5 transition active:scale-95"
          title="Edit Konfigurasi Website"
        >
          <Pencil class="w-3.5 h-3.5" />
          <span>Edit Konfigurasi</span>
        </button>
        <a
          :href="'http://' + (site.domain || 'localhost') + ':' + site.port"
          target="_blank"
          class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow active:scale-95"
        >
          <span>Buka Website</span>
          <ExternalLink class="w-3.5 h-3.5" />
        </a>
      </div>
    </div>
  </div>
</template>
