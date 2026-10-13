<script setup>
/**
 * RAZOS Ultra-Fast Anti-Hang Thumbnail Component
 * Created: 2026-10-13 (v2.4.6)
 * Deskripsi:
 * - Mengatasi masalah UI freeze / "Not Responding" saat membuka direktori dengan ratusan/ribuan berkas media.
 * - IntersectionObserver: Hanya memicu fetch saat kartu berada di dalam viewport atau mendekati viewport (rootMargin: 150px).
 * - Global Concurrency Queue: Membatasi maksimal 4 permintaan fetch bersamaan untuk menjaga UI thread dan I/O tetap responsif.
 * - Auto Queue Cancellation: Membatalkan permintaan item yang di-scroll cepat melewati viewport sebelum sempat dimuat.
 * - Memory & HTTP Caching: Menyimpan cache status item yang sudah dimuat agar tidak terjadi re-fetch berulang saat scroll.
 * - Safe Size Protection: File gambar berukuran raksasa (> 20 MB) dilindungi dengan badge HD/4K dan fallback ikon agar tidak meledakkan memori browser.
 * - Fallback Graceful: Placeholder shimmer dan transisi fade-in halus tanpa tata letak bergeser (zero layout shift).
 */
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Image as ImageIcon, AlertCircle } from 'lucide-vue-next'

const props = defineProps({
  path: {
    type: String,
    required: true
  },
  name: {
    type: String,
    default: ''
  },
  size: {
    type: Number,
    default: 0
  },
  extension: {
    type: String,
    default: ''
  },
  customClass: {
    type: String,
    default: 'w-full h-full object-cover'
  }
})

// Batas aman ukuran gambar untuk preview otomatis (20 MB)
const MAX_SAFE_PREVIEW_SIZE = 20 * 1024 * 1024

// State Komponen
const containerRef = ref(null)
const isVisible = ref(false)
const isLoading = ref(true)
const isLoaded = ref(false)
const isError = ref(false)
const isTooLarge = ref(false)
const imageSrc = ref('')

// Module-level Singleton Concurrency Queue & Cache
let observer = null
let queueTask = null

// Modul Cache Global
if (!window.__RAZ_THUMBNAIL_CACHE__) {
  window.__RAZ_THUMBNAIL_CACHE__ = new Set()
}
if (!window.__RAZ_QUEUE_STATE__) {
  window.__RAZ_QUEUE_STATE__ = {
    runningCount: 0,
    maxConcurrent: 4,
    queue: []
  }
}

const queueState = window.__RAZ_QUEUE_STATE__
const thumbCache = window.__RAZ_THUMBNAIL_CACHE__

function processNextQueue() {
  if (queueState.runningCount >= queueState.maxConcurrent || queueState.queue.length === 0) {
    return
  }

  // Prioritaskan item terbaru (LIFO) yang saat ini terlihat di viewport
  const task = queueState.queue.pop()
  if (!task || task.cancelled) {
    processNextQueue()
    return
  }

  queueState.runningCount++
  task.execute(() => {
    queueState.runningCount--
    processNextQueue()
  })
}

function enqueueThumbnailTask(executeFn) {
  const task = {
    cancelled: false,
    execute: (done) => {
      executeFn(done)
    }
  }
  queueState.queue.push(task)
  processNextQueue()
  return task
}

function loadThumbnail() {
  if (isLoaded.value || isError.value) return

  // Cek apakah ukuran file melebihi batas aman
  if (props.size > MAX_SAFE_PREVIEW_SIZE) {
    isTooLarge.value = true
    isLoading.value = false
    return
  }

  const encodedPath = encodeURIComponent(props.path)
  const targetUrl = `/api/v1/files/thumbnail?path=${encodedPath}&t=${encodeURIComponent(props.extension)}`

  // Jika sudah ada di cache global peramban
  if (thumbCache.has(props.path)) {
    imageSrc.value = targetUrl
    isLoaded.value = true
    isLoading.value = false
    return
  }

  // Masukkan ke antrean pembatasan konkurensi (Concurrency Queue)
  queueTask = enqueueThumbnailTask((done) => {
    if (!isVisible.value) {
      done()
      return
    }

    const img = new Image()
    img.decoding = 'async'
    img.loading = 'lazy'

    img.onload = () => {
      thumbCache.add(props.path)
      imageSrc.value = targetUrl
      isLoaded.value = true
      isLoading.value = false
      done()
    }

    img.onerror = () => {
      isError.value = true
      isLoading.value = false
      done()
    }

    img.src = targetUrl
  })
}

function handleImageError() {
  isError.value = true
  isLoading.value = false
}

onMounted(() => {
  // Intersection Observer untuk memantau kapan kartu memasuki viewport
  if ('IntersectionObserver' in window) {
    observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (entry.isIntersecting) {
          isVisible.value = true
          loadThumbnail()
        } else {
          isVisible.value = false
          // Batalkan tugas antrean jika di-scroll lewat sebelum sempat dieksekusi
          if (queueTask && !isLoaded.value) {
            queueTask.cancelled = true
            queueTask = null
          }
        }
      },
      {
        root: null,
        rootMargin: '120px 0px 120px 0px', // Buffer preload 120px sebelum masuk viewport
        threshold: 0.01
      }
    )

    if (containerRef.value) {
      observer.observe(containerRef.value)
    }
  } else {
    // Fallback jika browser lawas tidak mendukung IntersectionObserver
    isVisible.value = true
    loadThumbnail()
  }
})

onBeforeUnmount(() => {
  if (observer && containerRef.value) {
    observer.unobserve(containerRef.value)
    observer.disconnect()
  }
  if (queueTask) {
    queueTask.cancelled = true
    queueTask = null
  }
})

// Pantau perubahan path jika komponen digunakan kembali (reused)
watch(
  () => props.path,
  () => {
    isLoaded.value = false
    isError.value = false
    isTooLarge.value = false
    isLoading.value = true
    imageSrc.value = ''
    if (isVisible.value) {
      loadThumbnail()
    }
  }
)
</script>

<template>
  <div
    ref="containerRef"
    class="w-full h-full relative flex items-center justify-center overflow-hidden bg-slate-950 select-none"
  >
    <!-- Placeholder Skeleton Shimmer saat sedang memuat -->
    <div
      v-if="isLoading && !isTooLarge"
      class="absolute inset-0 bg-slate-900 animate-pulse flex items-center justify-center"
    >
      <ImageIcon class="w-4 h-4 text-slate-700 animate-spin" />
    </div>

    <!-- Tampilan Berkas Gambar Berhasil Dimuat -->
    <img
      v-if="isLoaded && imageSrc"
      :src="imageSrc"
      :alt="name"
      :class="[customClass, 'transition-opacity duration-300 opacity-100']"
      loading="lazy"
      decoding="async"
      @error="handleImageError"
    />

    <!-- Fallback Jika Gambar Sangat Besar (> 20 MB) Demi Mencegah Hang -->
    <div
      v-else-if="isTooLarge"
      class="w-full h-full flex flex-col items-center justify-center bg-slate-900/90 text-amber-400 p-1 text-center"
      title="File gambar terlalu besar untuk preview langsung (klik dua kali untuk membuka)"
    >
      <ImageIcon class="w-4 h-4 text-amber-400 mb-0.5" />
      <span class="text-[8px] font-mono font-bold bg-amber-500/20 px-1 rounded text-amber-300">HD FILE</span>
    </div>

    <!-- Fallback Jika Terjadi Error / Format Tidak Didukung Browser -->
    <div
      v-else-if="isError"
      class="w-full h-full flex items-center justify-center bg-slate-900/60 text-slate-500"
      title="Pratinjau tidak tersedia"
    >
      <ImageIcon class="w-5 h-5 text-slate-600" />
    </div>
  </div>
</template>
