<script setup>
/**
 * RAZOS File Properties Window Component
 * Created: 2026-10-13 (v2.6.0 - Standalone Window Architecture)
 * Deskripsi: Rincian & Properti Berkas mandiri ala Windows Explorer Properties.
 * Berjalan sebagai Jendela OS terpisah (tidak memblokir/mengganggu File Manager).
 */
import { computed } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import RAZThumbnail from '../components/RAZThumbnail.vue'
import {
  Info, Folder, Film, Music, FileText, FileCode, File, Eye, Copy, Check
} from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: true
  }
})

const systemStore = useSystemStore()
const windowStore = useWindowStore()

const file = computed(() => props.win.customData || {})

function isImageFile(ext) {
  return ['.webp', '.png', '.jpg', '.jpeg', '.svg', '.gif', '.ico', '.bmp'].includes((ext || '').toLowerCase())
}
function isVideoFile(ext) {
  return ['.mp4', '.webm', '.mkv', '.mov', '.avi'].includes((ext || '').toLowerCase())
}
function isAudioFile(ext) {
  return ['.mp3', '.wav', '.ogg', '.aac', '.flac', '.m4a'].includes((ext || '').toLowerCase())
}
function isPdfFile(ext) {
  return (ext || '').toLowerCase() === '.pdf'
}
function isCodeFile(ext) {
  return ['.js', '.vue', '.json', '.html', '.css', '.md', '.go', '.py', '.sh', '.yaml', '.yml', '.env'].includes((ext || '').toLowerCase())
}
function isDocFile(ext) {
  return ['.txt', '.log', '.csv', '.rtf'].includes((ext || '').toLowerCase())
}

function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return isNaN(d.getTime()) ? dateStr : d.toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

function openThisFile() {
  if (file.value && !file.value.is_dir) {
    windowStore.openFileViewer(file.value)
  }
}
</script>

<template>
  <div class="h-full flex flex-col justify-between text-xs select-none gap-3">
    <!-- Body Content -->
    <div class="flex-1 overflow-y-auto space-y-3.5 pr-1 text-slate-300">
      <!-- Preview Ikon & Nama Berkas -->
      <div class="flex items-center gap-4 p-3.5 bg-slate-950/70 rounded-2xl border border-slate-800/80">
        <div v-if="isImageFile(file.extension)" class="w-16 h-16 rounded-xl overflow-hidden bg-slate-900 border border-slate-800 shrink-0 shadow">
          <RAZThumbnail
            :path="file.path"
            :name="file.name"
            :size="file.size"
            :extension="file.extension"
          />
        </div>
        <div v-else class="w-16 h-16 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-center shrink-0 shadow">
          <Folder v-if="file.is_dir" class="w-8 h-8 text-sky-400 fill-sky-400/20" />
          <Film v-else-if="isVideoFile(file.extension)" class="w-7 h-7 text-purple-400" />
          <Music v-else-if="isAudioFile(file.extension)" class="w-7 h-7 text-amber-400" />
          <FileText v-else-if="isPdfFile(file.extension)" class="w-7 h-7 text-rose-400" />
          <FileCode v-else-if="isCodeFile(file.extension)" class="w-7 h-7 text-emerald-400" />
          <FileText v-else-if="isDocFile(file.extension)" class="w-7 h-7 text-sky-400" />
          <File v-else class="w-7 h-7 text-slate-400" />
        </div>
        <div class="min-w-0 flex-1">
          <h5 class="font-bold text-slate-100 text-sm break-all leading-snug">{{ file.name }}</h5>
          <span class="text-[11px] text-slate-400 font-mono block mt-0.5">
            {{ file.is_dir ? 'Folder Direktori' : (file.extension?.toUpperCase() || 'Berkas') }}
          </span>
        </div>
      </div>

      <!-- Metadata Properties Grid -->
      <div class="space-y-2 bg-slate-950/50 p-3.5 rounded-2xl border border-slate-800/80 font-mono text-[11px]">
        <div class="flex justify-between py-1 border-b border-slate-800/50">
          <span class="text-slate-400 font-sans">Ukuran:</span>
          <span class="text-slate-200 font-medium">
            {{ file.is_dir ? '-' : formatBytes(file.size) + ' (' + (file.size || 0).toLocaleString('id-ID') + ' bytes)' }}
          </span>
        </div>
        <div class="flex justify-between py-1 border-b border-slate-800/50">
          <span class="text-slate-400 font-sans">Format Ekstensi:</span>
          <span class="text-slate-200 uppercase">{{ file.is_dir ? 'Folder' : (file.extension || 'Tanpa Ekstensi') }}</span>
        </div>
        <div class="flex flex-col gap-1 py-1 border-b border-slate-800/50 font-sans">
          <span class="text-slate-400">Lokasi Path Fisik:</span>
          <span class="text-slate-300 font-mono text-[11px] bg-slate-900/80 p-2 rounded-lg break-all border border-slate-800 select-text">
            {{ file.path }}
          </span>
        </div>
        <div class="flex justify-between py-1 font-sans">
          <span class="text-slate-400">Terakhir Dimodifikasi:</span>
          <span class="text-slate-300 font-mono text-[11px]">{{ formatDate(file.mod_time) }}</span>
        </div>
      </div>
    </div>

    <!-- Window Footer Actions -->
    <div class="flex items-center justify-between border-t border-slate-800 pt-3 shrink-0">
      <button
        @click="windowStore.closeWindow(win.id)"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
      >
        Tutup
      </button>

      <button
        v-if="!file.is_dir"
        @click="openThisFile"
        class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold flex items-center gap-1.5 transition shadow active:scale-95"
      >
        <Eye class="w-3.5 h-3.5" />
        <span>Buka Berkas</span>
      </button>
    </div>
  </div>
</template>
