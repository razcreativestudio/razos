<script setup>
/**
 * RAZOS Window Frame Component
 * Updated: 2026-10-08
 * Deskripsi: Komponen container window dengan fitur draggable, resizable, stacking z-index,
 * minimize, maximize, dan styling glassmorphism Slate-900.
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { Minus, Square, Copy, X } from 'lucide-vue-next'

const props = defineProps({
  win: {
    type: Object,
    required: true
  }
})

const windowStore = useWindowStore()

// State interaksi Dragging
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0, winX: 0, winY: 0 })

// State interaksi Resizing
const isResizing = ref(false)
const resizeDirection = ref(null)
const resizeStart = ref({ mouseX: 0, mouseY: 0, width: 0, height: 0, winX: 0, winY: 0 })

// State Animasi Windows 11 Snap-to-Top Maximize Preview
const isSnapTopPreview = ref(false)

// Simpan ukuran dan koordinat jendela sebelum maximize agar dapat dipulihkan dengan tepat
const unmaximizedBounds = ref({
  width: props.win.width || 700,
  height: props.win.height || 480,
  x: props.win.x || 60,
  y: props.win.y || 40
})

// Dimensi dan Koordinat Dinamis yang Aman (Mencegah NaN atau Nilai Rusak)
const winLeft = computed(() => {
  if (props.win.isMaximized) return '0px'
  const val = Number(props.win.x)
  return `${Number.isFinite(val) ? val : 60}px`
})

const winTop = computed(() => {
  if (props.win.isMaximized) return '0px'
  const val = Number(props.win.y)
  return `${Number.isFinite(val) ? val : 40}px`
})

const winWidth = computed(() => {
  if (props.win.isMaximized) return '100vw'
  const val = Number(props.win.width)
  return `${Number.isFinite(val) && val >= 200 ? val : 700}px`
})

const winHeight = computed(() => {
  if (props.win.isMaximized) return 'calc(100vh - 48px)'
  const val = Number(props.win.height)
  return `${Number.isFinite(val) && val >= 150 ? val : 480}px`
})

const winZIndex = computed(() => {
  const val = Number(props.win.zIndex)
  return Number.isFinite(val) && val > 0 ? val : 10
})

function startDrag(e) {
  // Hanya drag jika klik kiri mouse
  if (e.button !== undefined && e.button !== 0) return

  windowStore.bringToFront(props.win.id)
  
  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY

  const curX = Number.isFinite(Number(props.win.x)) ? Number(props.win.x) : 60
  const curY = Number.isFinite(Number(props.win.y)) ? Number(props.win.y) : 40
  const curW = Number.isFinite(Number(props.win.width)) && Number(props.win.width) >= 200 ? Number(props.win.width) : 700
  const curH = Number.isFinite(Number(props.win.height)) && Number(props.win.height) >= 150 ? Number(props.win.height) : 480

  props.win.x = curX
  props.win.y = curY
  props.win.width = curW
  props.win.height = curH

  // Jika window sedang maximized dan di-drag ke bawah, unmaximize seketika ala Windows 11
  if (props.win.isMaximized) {
    props.win.isMaximized = false
    const restoredW = unmaximizedBounds.value.width || curW || 700
    const restoredH = unmaximizedBounds.value.height || curH || 480
    props.win.width = restoredW
    props.win.height = restoredH
    // Posisikan window di bawah kursor mouse secara presisi
    props.win.x = Math.max(10, clientX - Math.round(restoredW / 2))
    props.win.y = Math.max(10, clientY - 18)
  } else {
    unmaximizedBounds.value = {
      width: curW,
      height: curH,
      x: curX,
      y: curY
    }
  }

  isDragging.value = true

  dragStart.value = {
    x: clientX,
    y: clientY,
    winX: props.win.x,
    winY: props.win.y
  }

  window.addEventListener('mousemove', onDrag)
  window.addEventListener('mouseup', stopDrag)
  window.addEventListener('touchmove', onDrag)
  window.addEventListener('touchend', stopDrag)
}

function onDrag(e) {
  if (!isDragging.value) return
  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY

  const dx = clientX - dragStart.value.x
  const dy = clientY - dragStart.value.y

  props.win.x = Math.max(0, dragStart.value.winX + dx)
  props.win.y = Math.max(0, dragStart.value.winY + dy)

  // Cek apakah kursor mentok ke bagian paling atas layar (Windows 11 Snap Gesture <= 12px)
  if (clientY <= 12) {
    isSnapTopPreview.value = true
  } else {
    isSnapTopPreview.value = false
  }
}

function handleToggleMaximize() {
  if (!props.win.isMaximized) {
    unmaximizedBounds.value = {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y
    }
  } else {
    props.win.width = unmaximizedBounds.value.width
    props.win.height = unmaximizedBounds.value.height
    props.win.x = unmaximizedBounds.value.x
    props.win.y = unmaximizedBounds.value.y
  }
  windowStore.toggleMaximize(props.win.id)
}

function stopDrag() {
  isDragging.value = false
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDrag)
  window.removeEventListener('touchmove', onDrag)
  window.removeEventListener('touchend', stopDrag)

  // Jika dilepas di bagian atas layar, eksekusi maximize otomatis (Windows 11 Snap Gesture)
  if (isSnapTopPreview.value) {
    isSnapTopPreview.value = false
    props.win.isMaximized = true
    windowStore.bringToFront(props.win.id)
    windowStore.saveWindowBounds(props.win.id, {
      width: unmaximizedBounds.value.width,
      height: unmaximizedBounds.value.height,
      x: unmaximizedBounds.value.x,
      y: unmaximizedBounds.value.y,
      isMaximized: true
    })
  } else if (!props.win.isMaximized) {
    unmaximizedBounds.value = {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y
    }
    // Simpan ukuran dan posisi terakhir ke localStorage
    windowStore.saveWindowBounds(props.win.id, {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y,
      isMaximized: false
    })
  }
}

function startResize(e, dir) {
  if (props.win.isMaximized) return
  e.stopPropagation()
  isResizing.value = true
  resizeDirection.value = dir
  windowStore.bringToFront(props.win.id)

  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY

  resizeStart.value = {
    mouseX: clientX,
    mouseY: clientY,
    width: props.win.width,
    height: props.win.height,
    winX: props.win.x,
    winY: props.win.y
  }

  window.addEventListener('mousemove', onResize)
  window.addEventListener('mouseup', stopResize)
  window.addEventListener('touchmove', onResize)
  window.addEventListener('touchend', stopResize)
}

function onResize(e) {
  if (!isResizing.value) return
  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY

  const dx = clientX - resizeStart.value.mouseX
  const dy = clientY - resizeStart.value.mouseY

  if (resizeDirection.value.includes('e')) {
    props.win.width = Math.max(props.win.minWidth || 300, resizeStart.value.width + dx)
  }
  if (resizeDirection.value.includes('s')) {
    props.win.height = Math.max(props.win.minHeight || 200, resizeStart.value.height + dy)
  }
}

function stopResize() {
  isResizing.value = false
  resizeDirection.value = null
  window.removeEventListener('mousemove', onResize)
  window.removeEventListener('mouseup', stopResize)
  window.removeEventListener('touchmove', onResize)
  window.removeEventListener('touchend', stopResize)

  if (!props.win.isMaximized) {
    unmaximizedBounds.value = {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y
    }
    // Simpan ukuran baru hasil resize ke localStorage
    windowStore.saveWindowBounds(props.win.id, {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y,
      isMaximized: false
    })
  }
}
</script>

<template>
  <transition name="raz-window-transition">
    <div
      v-show="!win.isMinimized"
      @mousedown="windowStore.bringToFront(win.id)"
      :style="{
        position: 'absolute',
        left: winLeft,
        top: winTop,
        width: winWidth,
        height: winHeight,
        zIndex: winZIndex
      }"
      class="flex flex-col rounded-xl overflow-hidden raz-glass raz-window select-none pointer-events-auto"
      :class="[
        windowStore.activeWindowId === win.id ? 'ring-1 ring-sky-400/40 shadow-2xl' : 'opacity-95 shadow-lg',
        isDragging || isResizing ? 'transition-none' : 'raz-window-fluid'
      ]"
    >
      <!-- Windows 11 Snap-to-Top Maximize Preview Halo Animation -->
      <Teleport to="body">
        <transition name="fade">
          <div
            v-if="isSnapTopPreview"
            class="fixed inset-x-2 top-2 h-[calc(100vh-56px)] rounded-2xl bg-sky-500/10 border-2 border-sky-400/80 backdrop-blur-[2px] pointer-events-none z-[9999] shadow-[0_0_40px_rgba(56,189,248,0.4)] animate-pulse"
          >
            <div class="absolute inset-0 bg-gradient-to-b from-sky-400/20 via-transparent to-sky-400/5 rounded-2xl"></div>
          </div>
        </transition>
      </Teleport>

      <!-- Header Window / Titlebar (Double Click to Maximize/Restore) -->
      <div
        @mousedown="startDrag"
        @touchstart.passive="startDrag"
        @dblclick="handleToggleMaximize"
        class="h-10 bg-slate-900/90 border-b border-slate-700/50 flex items-center justify-between px-3 cursor-move"
      >
        <!-- Title & Icon (Truncate Single Line Anti-Wrap) -->
        <div class="flex items-center gap-2 text-xs sm:text-sm font-medium text-slate-200 min-w-0 flex-1 mr-2 overflow-hidden pointer-events-none">
          <img :src="win.icon" :alt="win.title" class="w-4 h-4 shrink-0 pointer-events-none" />
          <span class="truncate whitespace-nowrap font-medium pointer-events-none" :title="win.title">{{ win.title }}</span>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-1.5 shrink-0 pointer-events-auto" @mousedown.stop @touchstart.stop>
          <!-- Minimize -->
          <button
            @click.stop="windowStore.minimizeWindow(win.id)"
            class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition active:scale-95"
            title="Minimize"
          >
            <Minus class="w-3.5 h-3.5" />
          </button>

          <!-- Maximize / Restore -->
          <button
            @click.stop="handleToggleMaximize"
            class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition active:scale-95"
            :title="win.isMaximized ? 'Restore' : 'Maximize'"
          >
            <component :is="win.isMaximized ? Copy : Square" class="w-3.5 h-3.5" />
          </button>

          <!-- Close -->
          <button
            @click.stop="windowStore.closeWindow(win.id)"
            class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-rose-600/80 text-slate-400 hover:text-white transition active:scale-95"
            title="Close"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <!-- Window Body Content -->
      <div class="flex-1 overflow-auto bg-slate-950/70 p-4 text-slate-100 select-auto pointer-events-auto">
        <slot></slot>
      </div>

      <!-- Resize Handles (Hanya jika tidak maximized) -->
      <template v-if="!win.isMaximized">
        <div @mousedown.stop="(e) => startResize(e, 'se')" class="absolute bottom-0 right-0 w-3.5 h-3.5 cursor-se-resize pointer-events-auto z-10"></div>
        <div @mousedown.stop="(e) => startResize(e, 'e')" class="absolute top-10 right-0 w-1.5 h-full cursor-e-resize pointer-events-auto z-10"></div>
        <div @mousedown.stop="(e) => startResize(e, 's')" class="absolute bottom-0 left-0 w-full h-1.5 cursor-s-resize pointer-events-auto z-10"></div>
      </template>
    </div>
  </transition>
</template>
