<script setup>
/**
 * RAZOS Universal Standalone File & Content Viewer
 * Created: 2026-10-13
 * Deskripsi: Komponen view mandiri untuk preview gambar zoomable, video player,
 * audio/music player, PDF document viewer, dan in-browser code/text editor.
 * Beroperasi di window frame terpisah tanpa mengaburkan atau memblokir File Manager.
 */
import { ref, onMounted, computed, watch } from 'vue'
import {
  ZoomIn, ZoomOut, RotateCw, RotateCcw, Maximize2, Download, Play, Pause,
  Volume2, VolumeX, Save, RefreshCw, CheckCircle2, FileText, Film, Music,
  FileCode, ExternalLink, Info
} from 'lucide-vue-next'
import { useSystemStore } from '../stores/RAZSystemStore'

const props = defineProps({
  file: {
    type: Object,
    required: true
  }
})

const systemStore = useSystemStore()

// File format helpers
const ext = computed(() => (props.file.extension || '').toLowerCase())

const isImage = computed(() => ['.webp', '.png', '.jpg', '.jpeg', '.svg', '.gif', '.ico', '.bmp'].includes(ext.value))
const isVideo = computed(() => ['.mp4', '.webm', '.mkv', '.mov', '.avi'].includes(ext.value))
const isAudio = computed(() => ['.mp3', '.wav', '.ogg', '.aac', '.flac', '.m4a'].includes(ext.value))
const isPdf = computed(() => ext.value === '.pdf')
const isCodeOrText = computed(() => !isImage.value && !isVideo.value && !isAudio.value && !isPdf.value)

// Raw stream URL
const rawFileUrl = computed(() => `/api/v1/files/raw?path=${encodeURIComponent(props.file.path)}`)

// 1. IMAGE VIEWER STATE
const imgScale = ref(1)
const imgRotation = ref(0)

function zoomIn() {
  imgScale.value = Math.min(4, imgScale.value + 0.25)
}
function zoomOut() {
  imgScale.value = Math.max(0.25, imgScale.value - 0.25)
}
function resetImage() {
  imgScale.value = 1
  imgRotation.value = 0
}
function rotateImage() {
  imgRotation.value = (imgRotation.value + 90) % 360
}

// 2. AUDIO VIEWER STATE
const audioRef = ref(null)
const isAudioPlaying = ref(false)
const audioCurrentTime = ref(0)
const audioDuration = ref(0)
const isAudioMuted = ref(false)

function toggleAudioPlay() {
  if (!audioRef.value) return
  if (audioRef.value.paused) {
    audioRef.value.play()
    isAudioPlaying.value = true
  } else {
    audioRef.value.pause()
    isAudioPlaying.value = false
  }
}

function onAudioTimeUpdate() {
  if (audioRef.value) {
    audioCurrentTime.value = audioRef.value.currentTime
  }
}

function onAudioLoaded() {
  if (audioRef.value) {
    audioDuration.value = audioRef.value.duration || 0
  }
}

function onAudioSeek(e) {
  if (audioRef.value) {
    audioRef.value.currentTime = Number(e.target.value)
  }
}

function toggleAudioMute() {
  if (audioRef.value) {
    audioRef.value.muted = !audioRef.value.muted
    isAudioMuted.value = audioRef.value.muted
  }
}

function formatAudioTime(secs) {
  if (!secs || isNaN(secs)) return '00:00'
  const m = Math.floor(secs / 60)
  const s = Math.floor(secs % 60)
  return `${m < 10 ? '0' : ''}${m}:${s < 10 ? '0' : ''}${s}`
}

// 3. CODE & TEXT EDITOR STATE
const textContent = ref('')
const originalContent = ref('')
const isLoadingText = ref(false)
const isSavingText = ref(false)
const isDirty = computed(() => textContent.value !== originalContent.value)

async function loadTextContent() {
  if (!isCodeOrText.value) return
  isLoadingText.value = true
  try {
    const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(props.file.path)}`)
    if (res.ok) {
      const data = await res.json()
      textContent.value = data.content || ''
      originalContent.value = data.content || ''
    }
  } catch (err) {}
  isLoadingText.value = false
}

async function saveTextContent() {
  if (isSavingText.value) return
  isSavingText.value = true
  try {
    const res = await fetch('/api/v1/files/content', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: props.file.path, content: textContent.value })
    })
    if (res.ok) {
      originalContent.value = textContent.value
      systemStore.showNotification('Sukses', `Berkas ${props.file.name} berhasil disimpan`, 'success')
    }
  } catch (err) {
    systemStore.showNotification('Gagal', 'Terjadi kesalahan saat menyimpan berkas', 'error')
  }
  isSavingText.value = false
}

function triggerDownload() {
  const link = document.createElement('a')
  link.href = rawFileUrl.value
  link.download = props.file.name
  link.click()
}

onMounted(() => {
  if (isCodeOrText.value) {
    loadTextContent()
  }
})

watch(() => props.file.path, () => {
  resetImage()
  if (isCodeOrText.value) {
    loadTextContent()
  }
})
</script>

<template>
  <div class="h-full flex flex-col bg-slate-950 text-slate-100 select-none overflow-hidden relative">
    <!-- 1. IMAGE VIEWER -->
    <template v-if="isImage">
      <!-- Toolbar Atas Gambar -->
      <div class="h-11 px-3 border-b border-slate-800 bg-slate-900/80 flex items-center justify-between text-xs shrink-0 z-10">
        <div class="flex items-center gap-2">
          <span class="font-mono text-[11px] text-slate-400 font-semibold">{{ Math.round(imgScale * 100) }}%</span>
          <span class="text-slate-600">|</span>
          <span class="text-[11px] text-slate-400 font-mono truncate max-w-[200px]">{{ file.name }}</span>
        </div>
        <div class="flex items-center gap-1.5">
          <button @click="zoomOut" class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition" title="Zoom Out">
            <ZoomOut class="w-3.5 h-3.5" />
          </button>
          <button @click="resetImage" class="px-2 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 font-mono text-[11px] transition" title="Reset Zoom">
            100%
          </button>
          <button @click="zoomIn" class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition" title="Zoom In">
            <ZoomIn class="w-3.5 h-3.5" />
          </button>
          <button @click="rotateImage" class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition" title="Putar 90°">
            <RotateCw class="w-3.5 h-3.5" />
          </button>
          <button @click="triggerDownload" class="p-1.5 rounded-lg bg-sky-500/20 hover:bg-sky-500/30 text-sky-400 transition" title="Download Gambar">
            <Download class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <!-- Canvas Gambar Interaktif -->
      <div class="flex-1 overflow-auto flex items-center justify-center p-4 bg-slate-950/90 relative">
        <div
          class="transition-transform duration-200 ease-out flex items-center justify-center"
          :style="{
            transform: `scale(${imgScale}) rotate(${imgRotation}deg)`
          }"
        >
          <img
            :src="rawFileUrl"
            :alt="file.name"
            class="max-w-full max-h-[72vh] object-contain rounded-lg shadow-2xl pointer-events-none select-none border border-slate-800/80"
          />
        </div>
      </div>
    </template>

    <!-- 2. VIDEO PLAYER -->
    <template v-else-if="isVideo">
      <div class="flex-1 bg-black flex flex-col justify-center items-center overflow-hidden">
        <video
          :src="rawFileUrl"
          controls
          autoplay
          class="w-full h-full max-h-full object-contain"
        >
          Peramban Anda tidak mendukung pemutaran video HTML5.
        </video>
      </div>
    </template>

    <!-- 3. AUDIO / MUSIC PLAYER -->
    <template v-else-if="isAudio">
      <div class="flex-1 flex flex-col items-center justify-center p-6 bg-gradient-to-b from-slate-900 to-slate-950 gap-6">
        <!-- Visual Vinyl Disc -->
        <div
          class="w-36 h-36 rounded-full border-4 border-amber-500/30 shadow-2xl flex items-center justify-center bg-slate-950 relative overflow-hidden"
          :class="isAudioPlaying ? 'animate-[spin_6s_linear_infinite]' : ''"
        >
          <div class="w-20 h-20 rounded-full bg-slate-900 border-2 border-amber-400/50 flex items-center justify-center shadow-inner">
            <Music class="w-8 h-8 text-amber-400" />
          </div>
        </div>

        <div class="text-center max-w-sm px-4">
          <h4 class="font-bold text-slate-100 text-sm truncate leading-snug">{{ file.name }}</h4>
          <span class="text-[11px] text-slate-400 font-mono mt-0.5 block uppercase tracking-wider">{{ ext.replace('.', '') }} AUDIO</span>
        </div>

        <!-- Scrubber Range -->
        <div class="w-full max-w-md flex flex-col gap-1.5">
          <input
            type="range"
            min="0"
            :max="audioDuration || 100"
            :value="audioCurrentTime"
            @input="onAudioSeek"
            class="w-full accent-amber-400 h-1.5 bg-slate-800 rounded-lg cursor-pointer"
          />
          <div class="flex justify-between text-[10px] text-slate-400 font-mono">
            <span>{{ formatAudioTime(audioCurrentTime) }}</span>
            <span>{{ formatAudioTime(audioDuration) }}</span>
          </div>
        </div>

        <!-- Controls -->
        <div class="flex items-center gap-4">
          <button @click="toggleAudioMute" class="p-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition">
            <VolumeX v-if="isAudioMuted" class="w-4 h-4 text-rose-400" />
            <Volume2 v-else class="w-4 h-4 text-slate-300" />
          </button>
          <button
            @click="toggleAudioPlay"
            class="w-12 h-12 rounded-2xl bg-amber-400 hover:bg-amber-300 text-slate-950 flex items-center justify-center shadow-lg transition active:scale-95"
          >
            <Pause v-if="isAudioPlaying" class="w-5 h-5 fill-slate-950" />
            <Play v-else class="w-5 h-5 fill-slate-950 ml-0.5" />
          </button>
          <button @click="triggerDownload" class="p-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition" title="Download Audio">
            <Download class="w-4 h-4" />
          </button>
        </div>

        <audio
          ref="audioRef"
          :src="rawFileUrl"
          @timeupdate="onAudioTimeUpdate"
          @loadedmetadata="onAudioLoaded"
          @play="isAudioPlaying = true"
          @pause="isAudioPlaying = false"
          @ended="isAudioPlaying = false"
          class="hidden"
        ></audio>
      </div>
    </template>

    <!-- 4. PDF VIEWER -->
    <template v-else-if="isPdf">
      <div class="h-10 px-3 border-b border-slate-800 bg-slate-900/80 flex items-center justify-between text-xs shrink-0">
        <span class="text-slate-300 font-mono text-[11px] truncate">{{ file.name }}</span>
        <div class="flex items-center gap-2">
          <a
            :href="rawFileUrl"
            target="_blank"
            class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 flex items-center gap-1.5 transition text-[11px]"
          >
            <ExternalLink class="w-3.5 h-3.5" />
            <span>Tab Baru</span>
          </a>
          <button @click="triggerDownload" class="px-2.5 py-1 rounded-lg bg-rose-500/20 hover:bg-rose-500/30 text-rose-400 flex items-center gap-1.5 transition text-[11px]">
            <Download class="w-3.5 h-3.5" />
            <span>Download</span>
          </button>
        </div>
      </div>
      <iframe :src="rawFileUrl" class="w-full flex-1 border-none bg-slate-950"></iframe>
    </template>

    <!-- 5. IN-BROWSER CODE & TEXT EDITOR -->
    <template v-else>
      <div class="h-11 px-3 border-b border-slate-800 bg-slate-900/80 flex items-center justify-between text-xs shrink-0">
        <div class="flex items-center gap-2">
          <FileCode class="w-4 h-4 text-emerald-400" />
          <span class="font-bold text-slate-200 text-xs truncate max-w-[240px]">{{ file.name }}</span>
          <span v-if="isDirty" class="px-1.5 py-0.5 rounded text-[10px] bg-amber-500/20 text-amber-400 font-mono font-bold">Unsaved</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="loadTextContent"
            :disabled="isLoadingText"
            class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
            title="Muat Ulang Berkas"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="isLoadingText ? 'animate-spin' : ''" />
          </button>
          <button
            @click="saveTextContent"
            :disabled="isSavingText || !isDirty"
            class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 disabled:opacity-40 disabled:hover:bg-sky-500 text-slate-950 font-bold text-xs flex items-center gap-1.5 shadow transition active:scale-95"
          >
            <Save class="w-3.5 h-3.5" />
            <span>{{ isSavingText ? 'Menyimpan...' : 'Simpan' }}</span>
          </button>
        </div>
      </div>

      <div class="flex-1 relative flex">
        <textarea
          v-model="textContent"
          class="w-full h-full p-4 bg-slate-950 text-emerald-400 font-mono text-xs outline-none resize-none leading-relaxed selection:bg-sky-500/30 border-none"
          spellcheck="false"
          placeholder="Berkas kosong..."
        ></textarea>
      </div>
    </template>
  </div>
</template>
