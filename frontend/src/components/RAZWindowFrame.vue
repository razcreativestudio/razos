<script setup>
/**
 * RAZOS Window Frame Component
 * Updated: 2026-10-08
 * Deskripsi: Komponen container window dengan fitur draggable, resizable, stacking z-index,
 * minimize, maximize, dan styling glassmorphism Slate-900.
 */
import { ref, onMounted, onUnmounted } from 'vue'
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
  width: props.win.width,
  height: props.win.height,
  x: props.win.x,
  y: props.win.y
})

function startDrag(e) {
  windowStore.bringToFront(props.win.id)
  
  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY

  // Jika window sedang maximized dan di-drag ke bawah, unmaximize seketika ala Windows 11
  if (props.win.isMaximized) {
    props.win.isMaximized = false
    props.win.width = unmaximizedBounds.value.width || 700
    props.win.height = unmaximizedBounds.value.height || 480
    // Posisikan window di bawah kursor mouse secara presisi
    props.win.x = Math.max(10, clientX - Math.round(props.win.width / 2))
    props.win.y = Math.max(10, clientY - 18)
  } else {
    unmaximizedBounds.value = {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y
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

function stopDrag() {
  isDragging.value = false
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDrag)
  window.removeEventListener('touchmove', onDrag)
  window.removeEventListener('touchend', stopDrag)

  // Jika dilepas di bagian atas layar, eksekusi maximize otomatis
  if (isSnapTopPreview.value) {
    isSnapTopPreview.value = false
    props.win.isMaximized = true
    windowStore.bringToFront(props.win.id)
  } else if (!props.win.isMaximized) {
    unmaximizedBounds.value = {
      width: props.win.width,
      height: props.win.height,
      x: props.win.x,
      y: props.win.y
    }
    // Simpan ukuran dan posisi terakhir ke localStorage
    windowStore.saveWindowBounds(props.win.id, unmaximizedBounds.value)
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
    windowStore.saveWindowBounds(props.win.id, unmaximizedBounds.value)
  }
}
</script>

<template>
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
  <transition name="raz-window-transition">
    <div
      v-show="!win.isMinimized"
      @mousedown="windowStore.bringToFront(win.id)"
      :style="win.isMaximized ? {
        position: 'absolute',
        left: '0px',
        top: '0px',
        width: '100vw',
        height: 'calc(100vh - 48px)',
        zIndex: win.zIndex
      } : {
        position: 'absolute',
        left: `${win.x}px`,
        top: `${win.y}px`,
        width: `${win.width}px`,
        height: `${win.height}px`,
        zIndex: win.zIndex
      }"
      class="flex flex-col rounded-xl overflow-hidden raz-glass raz-window select-none"
      :class="[
        windowStore.activeWindowId === win.id ? 'ring-1 ring-sky-400/40 shadow-2xl' : 'opacity-95 shadow-lg',
        isDragging || isResizing ? 'transition-none' : 'raz-window-fluid'
      ]"
    >
    <!-- Header Window / Titlebar -->
    <div
      @mousedown="startDrag"
      class="h-10 bg-slate-900/90 border-b border-slate-700/50 flex items-center justify-between px-3 cursor-move"
    >
      <!-- Title & Icon (Truncate Single Line Anti-Wrap) -->
      <div class="flex items-center gap-2 text-xs sm:text-sm font-medium text-slate-200 min-w-0 flex-1 mr-2 overflow-hidden">
        <img :src="win.icon" :alt="win.title" class="w-4 h-4 shrink-0" />
        <span class="truncate whitespace-nowrap font-medium" :title="win.title">{{ win.title }}</span>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-1.5 shrink-0" @mousedown.stop>
        <!-- Minimize -->
        <button
          @click="windowStore.minimizeWindow(win.id)"
          class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
          title="Minimize"
        >
          <Minus class="w-3.5 h-3.5" />
        </button>

        <!-- Maximize / Restore -->
        <button
          @click="windowStore.toggleMaximize(win.id)"
          class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
          :title="win.isMaximized ? 'Restore' : 'Maximize'"
        >
          <component :is="win.isMaximized ? Copy : Square" class="w-3.5 h-3.5" />
        </button>

        <!-- Close -->
        <button
          @click="windowStore.closeWindow(win.id)"
          class="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-rose-600/80 text-slate-400 hover:text-white transition"
          title="Close"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- Window Body Content -->
    <div class="flex-1 overflow-auto bg-slate-950/70 p-4 text-slate-100">
      <slot></slot>
    </div>

    <!-- Resize Handles (Hanya jika tidak maximized) -->
    <template v-if="!win.isMaximized">
      <div @mousedown="(e) => startResize(e, 'se')" class="absolute bottom-0 right-0 w-3.5 h-3.5 cursor-se-resize"></div>
      <div @mousedown="(e) => startResize(e, 'e')" class="absolute top-10 right-0 w-1.5 h-full cursor-e-resize"></div>
      <div @mousedown="(e) => startResize(e, 's')" class="absolute bottom-0 left-0 w-full h-1.5 cursor-s-resize"></div>
    </template>
  </div>
  </transition>
</template>
