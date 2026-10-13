<script setup>
/**
 * RAZOS Universal Standalone File & Content Viewer
 * Created: 2026-10-13
 * Deskripsi: Komponen view mandiri untuk preview gambar zoomable, video player,
 * audio/music player, PDF document viewer, dan in-browser code/text editor.
 * Beroperasi di window frame terpisah tanpa mengaburkan atau memblokir File Manager.
 */
import { ref, onMounted, computed, watch, nextTick } from 'vue'
import {
  ZoomIn, ZoomOut, RotateCw, RotateCcw, Maximize2, Download, Play, Pause,
  Volume2, Volume1, VolumeX, Save, RefreshCw, CheckCircle2, FileText, Film, Music,
  FileCode, ExternalLink, Info, ShieldAlert, AlertCircle, Repeat, SkipBack, SkipForward, Disc,
  FastForward, Gauge
} from 'lucide-vue-next'
import { useSystemStore } from '../stores/RAZSystemStore'
import RAZMarkdownViewer from '../components/RAZMarkdownViewer.vue'

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
const isMarkdown = computed(() => ['.md', '.markdown', '.mdown', '.mkd'].includes(ext.value))
const isCodeOrText = computed(() => !isImage.value && !isVideo.value && !isAudio.value && !isPdf.value && !isMarkdown.value)

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

// 2. AUDIO VIEWER STATE (Modern Music Player Suite)
const audioRef = ref(null)
const isAudioPlaying = ref(false)
const audioCurrentTime = ref(0)
const audioDuration = ref(0)
const isAudioMuted = ref(false)
const audioVolume = ref(1)
const isAudioLoop = ref(false)
const prevVolume = ref(1)

// Kecepatan Pemutaran (Playback Speed)
const playbackSpeed = ref(1)
const availableSpeeds = [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2]
const showSpeedMenu = ref(false)

// Scrubber Timeline State (Anti-Jitter Seeking & Dragging)
const isSeeking = ref(false)
const seekValue = ref(0)
const isHoveringTimeline = ref(false)
const hoverPercent = ref(0)
const hoverTimeText = ref('00:00')

// Menghitung persentase durasi yang telah terputar (0% - 100%)
const audioProgressPercent = computed(() => {
  const dur = audioDuration.value || 0
  if (dur <= 0) return 0
  const curr = isSeeking.value ? seekValue.value : audioCurrentTime.value
  return Math.min(100, Math.max(0, (curr / dur) * 100))
})

// Waktu aktif yang ditampilkan di label (mengikuti posisi drag saat seeking)
const displayedCurrentTime = computed(() => {
  return isSeeking.value ? seekValue.value : audioCurrentTime.value
})

function toggleAudioPlay() {
  if (!audioRef.value) return
  if (audioRef.value.paused) {
    audioRef.value.play().then(() => {
      isAudioPlaying.value = true
    }).catch(err => {
      console.warn('Audio play gagal dipicu:', err)
      isAudioPlaying.value = false
    })
  } else {
    audioRef.value.pause()
    isAudioPlaying.value = false
  }
}

function seekRelative(secs) {
  if (!audioRef.value) return
  const maxDur = audioDuration.value || 9999
  const target = Math.max(0, Math.min(maxDur, audioRef.value.currentTime + secs))
  audioRef.value.currentTime = target
  audioCurrentTime.value = target
  seekValue.value = target
}

function toggleAudioLoop() {
  isAudioLoop.value = !isAudioLoop.value
  if (audioRef.value) {
    audioRef.value.loop = isAudioLoop.value
  }
}

// Beralih kecepatan pemutaran secara siklik (1x -> 1.25x -> 1.5x -> 2x -> 0.5x -> ...)
function cyclePlaybackSpeed() {
  const currentIndex = availableSpeeds.indexOf(playbackSpeed.value)
  const nextIndex = (currentIndex + 1) % availableSpeeds.length
  setPlaybackSpeed(availableSpeeds[nextIndex])
}

function setPlaybackSpeed(speed) {
  playbackSpeed.value = speed
  if (audioRef.value) {
    audioRef.value.playbackRate = speed
  }
  showSpeedMenu.value = false
}

// Handler Seeking Interaktif: Anti-Jitter Drag & Slide
function onSeekStart(e) {
  isSeeking.value = true
  seekValue.value = Number(e.target.value)
}

function onSeekInput(e) {
  isSeeking.value = true
  seekValue.value = Number(e.target.value)
}

function onSeekEnd(e) {
  const target = Number(e.target.value)
  if (audioRef.value) {
    audioRef.value.currentTime = target
    audioCurrentTime.value = target
  }
  seekValue.value = target
  isSeeking.value = false
}

// Klik langsung pada sembarang titik garis timeline
function onTimelineTrackClick(e) {
  if (!audioRef.value || !audioDuration.value) return
  const rect = e.currentTarget.getBoundingClientRect()
  const clickX = e.clientX - rect.left
  const ratio = Math.max(0, Math.min(1, clickX / rect.width))
  const target = ratio * audioDuration.value
  audioRef.value.currentTime = target
  audioCurrentTime.value = target
  seekValue.value = target
}

// Hover preview tooltip waktu pada timeline
function onTimelineMouseMove(e) {
  if (!audioDuration.value) return
  const rect = e.currentTarget.getBoundingClientRect()
  const mouseX = e.clientX - rect.left
  const ratio = Math.max(0, Math.min(1, mouseX / rect.width))
  hoverPercent.value = ratio * 100
  const previewTime = ratio * audioDuration.value
  hoverTimeText.value = formatAudioTime(previewTime)
  isHoveringTimeline.value = true
}

function onTimelineMouseLeave() {
  isHoveringTimeline.value = false
}

function onAudioTimeUpdate() {
  // Hanya perbarui audioCurrentTime jika pengguna tidak sedang menahan/men-slide slider
  if (audioRef.value && !isSeeking.value) {
    audioCurrentTime.value = audioRef.value.currentTime
  }
}

function onAudioLoaded() {
  if (audioRef.value) {
    audioDuration.value = audioRef.value.duration || 0
    audioRef.value.playbackRate = playbackSpeed.value
  }
}

function onAudioVolumeChange(e) {
  const val = Number(e.target.value)
  audioVolume.value = val
  if (audioRef.value) {
    audioRef.value.volume = val
    audioRef.value.muted = val === 0
    isAudioMuted.value = val === 0
  }
}

function toggleAudioMute() {
  if (!audioRef.value) return
  if (isAudioMuted.value || audioVolume.value === 0) {
    isAudioMuted.value = false
    audioVolume.value = prevVolume.value > 0 ? prevVolume.value : 0.8
    audioRef.value.muted = false
    audioRef.value.volume = audioVolume.value
  } else {
    prevVolume.value = audioVolume.value > 0 ? audioVolume.value : 0.8
    isAudioMuted.value = true
    audioVolume.value = 0
    audioRef.value.muted = true
    audioRef.value.volume = 0
  }
}

function formatAudioTime(secs) {
  if (!secs || isNaN(secs)) return '00:00'
  const m = Math.floor(secs / 60)
  const s = Math.floor(secs % 60)
  return `${m < 10 ? '0' : ''}${m}:${s < 10 ? '0' : ''}${s}`
}

function onAudioEnded() {
  if (!isAudioLoop.value) {
    isAudioPlaying.value = false
    audioCurrentTime.value = 0
  }
}

function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 3. CODE & TEXT EDITOR STATE
const textContent = ref('')
const originalContent = ref('')
const isLoadingText = ref(false)
const isSavingText = ref(false)
const isLockedFile = ref(false)
const lockErrorMessage = ref('')
const lockErrorCode = ref('')
const isDirty = computed(() => textContent.value !== originalContent.value)

async function loadTextContent() {
  if (!isCodeOrText.value) return
  isLoadingText.value = true
  isLockedFile.value = false
  lockErrorMessage.value = ''
  lockErrorCode.value = ''
  try {
    const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(props.file.path)}`)
    const data = await res.json()
    if (res.ok && !data.isLocked && !data.error) {
      textContent.value = data.content || ''
      originalContent.value = data.content || ''
    } else {
      isLockedFile.value = true
      lockErrorMessage.value = data.error || 'Berkas tidak dapat dibuka atau diakses.'
      lockErrorCode.value = data.errorCode || (data.notFound ? 'NOT_FOUND' : 'LOCKED')
    }
  } catch (err) {
    isLockedFile.value = true
    lockErrorMessage.value = `Gagal memuat berkas: ${err.message}`
    lockErrorCode.value = 'ERROR'
  }
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
  } else if (isAudio.value) {
    nextTick(() => {
      if (audioRef.value) {
        audioRef.value.playbackRate = playbackSpeed.value
        audioRef.value.play().then(() => {
          isAudioPlaying.value = true
        }).catch(() => {
          isAudioPlaying.value = false
        })
      }
    })
  }
})

watch(() => props.file.path, () => {
  resetImage()
  if (isCodeOrText.value) {
    loadTextContent()
  } else if (isAudio.value) {
    isAudioPlaying.value = false
    audioCurrentTime.value = 0
    seekValue.value = 0
    nextTick(() => {
      if (audioRef.value) {
        audioRef.value.currentTime = 0
        audioRef.value.playbackRate = playbackSpeed.value
        audioRef.value.play().then(() => {
          isAudioPlaying.value = true
        }).catch(() => {
          isAudioPlaying.value = false
        })
      }
    })
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

    <!-- 3. AUDIO / MUSIC PLAYER (Enterprise Vinyl & Equalizer Suite) -->
    <template v-else-if="isAudio">
      <div class="flex-1 flex flex-col justify-between p-5 bg-gradient-to-b from-slate-900 via-slate-950 to-slate-950 text-slate-100 select-none overflow-hidden h-full">
        <!-- TOP: Turntable Vinyl Disc & Equalizer Visualizer -->
        <div class="flex flex-col items-center justify-center flex-1 min-h-0 py-2 gap-3">
          <!-- Interactive Vinyl Disc -->
          <div
            @click="toggleAudioPlay"
            class="relative w-36 h-36 sm:w-44 sm:h-44 rounded-full border-4 border-amber-500/30 shadow-[0_0_35px_rgba(245,158,11,0.2)] flex items-center justify-center bg-slate-950 cursor-pointer group hover:border-amber-400/60 transition-all select-none"
            :class="isAudioPlaying ? 'animate-[spin_8s_linear_infinite]' : ''"
            :style="isAudioPlaying ? { animationDuration: `${Math.max(1.5, 8 / playbackSpeed)}s` } : {}"
            title="Klik untuk Putar / Jeda"
          >
            <!-- Vinyl Grooves Rings -->
            <div class="absolute inset-2 rounded-full border border-slate-800/80 pointer-events-none"></div>
            <div class="absolute inset-5 rounded-full border border-slate-800/70 pointer-events-none"></div>
            <div class="absolute inset-8 rounded-full border border-slate-800/60 pointer-events-none"></div>

            <!-- Central Vinyl Label -->
            <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-to-tr from-amber-600 to-amber-400 border-2 border-amber-300/80 flex items-center justify-center shadow-inner group-hover:scale-105 transition">
              <Music class="w-7 h-7 text-slate-950 fill-slate-950" />
            </div>

            <!-- Play/Pause Overlay on Hover -->
            <div class="absolute inset-0 rounded-full bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition backdrop-blur-[2px]">
              <Pause v-if="isAudioPlaying" class="w-10 h-10 text-white fill-white drop-shadow-md" />
              <Play v-else class="w-10 h-10 text-white fill-white ml-1 drop-shadow-md" />
            </div>
          </div>

          <!-- Animated Audio Wave / Equalizer Bars -->
          <div class="flex items-end justify-center gap-1 h-6 px-4">
            <span
              v-for="bar in 16"
              :key="bar"
              class="w-1 rounded-full transition-all duration-150"
              :class="isAudioPlaying ? 'bg-gradient-to-t from-amber-500 to-amber-300' : 'bg-slate-800 h-1'"
              :style="isAudioPlaying ? {
                height: `${Math.max(4, Math.sin(bar * 0.8 + audioCurrentTime * 6 * playbackSpeed) * 18 + 12)}px`
              } : { height: '3px' }"
            ></span>
          </div>
        </div>

        <!-- MIDDLE: Track Meta Information -->
        <div class="text-center px-4 shrink-0 my-1">
          <h4 class="font-bold text-slate-100 text-sm sm:text-base truncate leading-snug hover:text-amber-300 transition" :title="file.name">
            {{ file.name }}
          </h4>
          <div class="flex items-center justify-center gap-2 mt-1 font-mono text-[11px] text-slate-400">
            <span class="px-2 py-0.5 rounded-full bg-amber-500/20 text-amber-300 font-bold uppercase tracking-wider text-[10px]">
              {{ ext.replace('.', '') }}
            </span>
            <span>•</span>
            <span>{{ file.size ? formatBytes(file.size) : 'Audio Stream' }}</span>
            <span>•</span>
            <span class="text-emerald-400 font-bold">HQ STEREO</span>
          </div>
        </div>

        <!-- BOTTOM: Scrubber Timeline, Playback Controls, Volume -->
        <div class="w-full flex flex-col gap-3 shrink-0 pt-2 border-t border-slate-800/80">
          <!-- Timeline Scrubber Container (Clickable & Slidable Mulus Anti-Jitter) -->
          <div class="w-full flex flex-col gap-1.5 select-none">
            <div
              class="relative w-full h-3 group flex items-center cursor-pointer"
              @mousemove="onTimelineMouseMove"
              @mouseleave="onTimelineMouseLeave"
              @click="onTimelineTrackClick"
            >
              <!-- Tooltip Preview Waktu Melayang saat Hover -->
              <div
                v-if="isHoveringTimeline && audioDuration > 0"
                class="absolute -top-7 -translate-x-1/2 px-2 py-0.5 rounded-md bg-slate-800/95 border border-slate-700/80 text-amber-300 text-[10px] font-mono shadow-xl pointer-events-none z-30 transition-all duration-75"
                :style="{ left: `${hoverPercent}%` }"
              >
                {{ hoverTimeText }}
              </div>

              <!-- Background Track -->
              <div class="w-full h-1.5 group-hover:h-2 bg-slate-800 rounded-full overflow-hidden transition-all duration-150">
                <!-- Hover Ghost Fill Track -->
                <div
                  v-if="isHoveringTimeline"
                  class="absolute h-full bg-slate-700/60 rounded-full pointer-events-none"
                  :style="{ width: `${hoverPercent}%` }"
                ></div>

                <!-- Filled Progress Track -->
                <div
                  class="h-full bg-gradient-to-r from-amber-500 to-amber-300 rounded-full transition-none"
                  :style="{ width: `${audioProgressPercent}%` }"
                ></div>
              </div>

              <!-- Thumb Knob Indicator (Emas Putih) -->
              <div
                class="absolute top-1/2 -translate-y-1/2 -translate-x-1/2 w-3.5 h-3.5 rounded-full bg-amber-400 border-2 border-white shadow-md shadow-amber-500/50 pointer-events-none transition-transform duration-75"
                :class="isSeeking ? 'scale-125' : 'group-hover:scale-110'"
                :style="{ left: `${audioProgressPercent}%` }"
              ></div>

              <!-- Invisible Range Input di atasnya untuk Dragging Presisi & Touch Mobile -->
              <input
                type="range"
                min="0"
                :max="audioDuration || 100"
                step="0.1"
                :value="isSeeking ? seekValue : audioCurrentTime"
                @mousedown="onSeekStart"
                @touchstart="onSeekStart"
                @input="onSeekInput"
                @change="onSeekEnd"
                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-20"
                title="Geser atau klik untuk mengubah waktu musik"
              />
            </div>

            <!-- Current Time & Total Duration Info -->
            <div class="flex justify-between items-center text-[11px] text-slate-400 font-mono px-0.5">
              <span class="text-amber-400 font-semibold">{{ formatAudioTime(displayedCurrentTime) }}</span>
              <div class="flex items-center gap-1.5">
                <span v-if="playbackSpeed !== 1" class="text-[10px] text-amber-300 font-bold bg-amber-500/20 px-1.5 py-0.2 rounded border border-amber-500/30">
                  {{ playbackSpeed }}x
                </span>
                <span>{{ formatAudioTime(audioDuration) }}</span>
              </div>
            </div>
          </div>

          <!-- Controls Main Bar -->
          <div class="flex items-center justify-between gap-2 px-1">
            <!-- Left: Repeat Mode -->
            <button
              @click="toggleAudioLoop"
              class="p-2 rounded-xl transition relative"
              :class="isAudioLoop ? 'bg-amber-500/20 text-amber-400' : 'text-slate-400 hover:text-white hover:bg-slate-800'"
              :title="isAudioLoop ? 'Ulangi Lagu: Aktif' : 'Ulangi Lagu: Nonaktif'"
            >
              <Repeat class="w-4 h-4" />
              <span v-if="isAudioLoop" class="w-1.5 h-1.5 rounded-full bg-amber-400 absolute bottom-1 right-1"></span>
            </button>

            <!-- Center: Skip -10s, Play/Pause, Skip +10s -->
            <div class="flex items-center gap-2 sm:gap-3">
              <button
                @click="seekRelative(-10)"
                class="p-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition active:scale-95"
                title="Mundur 10 Detik"
              >
                <SkipBack class="w-4 h-4" />
              </button>

              <button
                @click="toggleAudioPlay"
                class="w-12 h-12 rounded-full bg-gradient-to-tr from-amber-500 to-amber-400 hover:from-amber-400 hover:to-amber-300 text-slate-950 flex items-center justify-center shadow-lg shadow-amber-500/30 transition hover:scale-105 active:scale-95"
                :title="isAudioPlaying ? 'Jeda' : 'Putar'"
              >
                <Pause v-if="isAudioPlaying" class="w-5 h-5 fill-slate-950" />
                <Play v-else class="w-5 h-5 fill-slate-950 ml-0.5" />
              </button>

              <button
                @click="seekRelative(10)"
                class="p-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition active:scale-95"
                title="Maju 10 Detik"
              >
                <SkipForward class="w-4 h-4" />
              </button>
            </div>

            <!-- Right: Playback Speed Control Pill & Download -->
            <div class="flex items-center gap-1 relative">
              <!-- Playback Speed Button (Klik untuk percepat / ganti siklis) -->
              <div class="relative">
                <button
                  @click="cyclePlaybackSpeed"
                  @contextmenu.prevent="showSpeedMenu = !showSpeedMenu"
                  class="px-2 py-1.5 rounded-xl text-xs font-mono font-bold transition flex items-center gap-1 active:scale-95"
                  :class="playbackSpeed !== 1
                    ? 'bg-amber-500/25 text-amber-300 border border-amber-500/40 shadow-sm shadow-amber-500/20'
                    : 'text-slate-400 hover:text-white hover:bg-slate-800'"
                  :title="`Kecepatan: ${playbackSpeed}x (Klik untuk ganti, Klik kanan untuk menu)`"
                >
                  <FastForward v-if="playbackSpeed > 1" class="w-3.5 h-3.5 text-amber-400" />
                  <Gauge v-else class="w-3.5 h-3.5 text-slate-400" />
                  <span>{{ playbackSpeed }}x</span>
                </button>

                <!-- Speed Menu Dropdown Popup -->
                <div v-if="showSpeedMenu" class="fixed inset-0 z-30" @click="showSpeedMenu = false"></div>
                <div
                  v-if="showSpeedMenu"
                  class="absolute bottom-full right-0 mb-2 py-1 px-1 bg-slate-900 border border-slate-700/80 rounded-xl shadow-2xl z-40 flex flex-col gap-0.5 min-w-[76px]"
                >
                  <button
                    v-for="s in availableSpeeds"
                    :key="s"
                    @click="setPlaybackSpeed(s)"
                    class="px-2.5 py-1 rounded-lg text-[11px] font-mono text-left transition flex items-center justify-between"
                    :class="playbackSpeed === s ? 'bg-amber-500/20 text-amber-300 font-bold' : 'text-slate-300 hover:bg-slate-800'"
                  >
                    <span>{{ s }}x</span>
                    <span v-if="playbackSpeed === s" class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
                  </button>
                </div>
              </div>

              <!-- Download Button -->
              <button
                @click="triggerDownload"
                class="p-2 rounded-xl text-slate-400 hover:text-sky-400 hover:bg-slate-800 transition"
                title="Unduh Berkas Audio"
              >
                <Download class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Volume Controls Bar -->
          <div class="flex items-center justify-center gap-2 pt-1">
            <button @click="toggleAudioMute" class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition">
              <VolumeX v-if="isAudioMuted || audioVolume === 0" class="w-3.5 h-3.5 text-rose-400" />
              <Volume1 v-else-if="audioVolume < 0.5" class="w-3.5 h-3.5 text-amber-300" />
              <Volume2 v-else class="w-3.5 h-3.5 text-amber-400" />
            </button>
            <input
              type="range"
              min="0"
              max="1"
              step="0.01"
              :value="isAudioMuted ? 0 : audioVolume"
              @input="onAudioVolumeChange"
              class="w-24 sm:w-28 accent-amber-400 h-1 bg-slate-800 rounded-lg cursor-pointer"
            />
            <span class="text-[10px] text-slate-500 font-mono w-8">
              {{ isAudioMuted ? '0%' : Math.round(audioVolume * 100) + '%' }}
            </span>
          </div>
        </div>

        <audio
          ref="audioRef"
          :src="rawFileUrl"
          @timeupdate="onAudioTimeUpdate"
          @loadedmetadata="onAudioLoaded"
          @play="isAudioPlaying = true"
          @pause="isAudioPlaying = false"
          @ended="onAudioEnded"
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

    <!-- 5. DEDICATED MARKDOWN DOCUMENT VIEWER & READER SUITE -->
    <template v-else-if="isMarkdown">
      <RAZMarkdownViewer :file="file" />
    </template>

    <!-- 6. IN-BROWSER CODE & TEXT EDITOR -->
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

      <!-- Tampilan Berkas Terkunci oleh Kernel Sistem / Tidak Dapat Diakses -->
      <div v-if="isLockedFile" class="flex-1 flex flex-col items-center justify-center p-8 text-center bg-slate-950/90 gap-4">
        <div class="w-16 h-16 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 shadow-xl shadow-amber-500/5">
          <ShieldAlert class="w-8 h-8" />
        </div>
        <div class="max-w-md space-y-1.5">
          <h3 class="text-sm font-bold text-slate-100">Berkas Dibatasi oleh Sistem Operasi</h3>
          <p class="text-xs text-slate-400 leading-relaxed">
            {{ lockErrorMessage }}
          </p>
        </div>
        <div class="px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-[11px] font-mono text-slate-400 flex items-center gap-2">
          <span>Status Kode:</span>
          <span class="text-amber-400 font-bold">{{ lockErrorCode }}</span>
        </div>
      </div>

      <!-- Textarea Editor Normal -->
      <div v-else class="flex-1 relative flex">
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
