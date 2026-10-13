<script setup>
/**
 * RAZOS Universal Context Menu Component
 * Updated: 2026-10-13 (Enterprise Context Menu System)
 * Deskripsi: Komponen visual Custom Context Menu berdesain Glassmorphism Slate-900.
 * Menampilkan item aksi klik kanan yang responsif, rapi, dengan dukungan ikon Lucide,
 * keyboard shortcut hint, pemisah seksi (divider), dan status danger (merah).
 */

import { onMounted, onUnmounted } from 'vue'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'

const contextMenuStore = useContextMenuStore()

function handleItemClick(item) {
  if (item.disabled) return
  contextMenuStore.hide()
  if (typeof item.action === 'function') {
    try {
      item.action()
    } catch (err) {
      console.warn('Gagal mengeksekusi aksi context menu:', err)
    }
  }
}

// Tutup menu saat klik sembarang di luar atau tekan tombol Escape
function handleGlobalClick(e) {
  if (contextMenuStore.isOpen) {
    contextMenuStore.hide()
  }
}

function handleGlobalKeydown(e) {
  if (e.key === 'Escape' && contextMenuStore.isOpen) {
    contextMenuStore.hide()
  }
}

onMounted(() => {
  window.addEventListener('click', handleGlobalClick)
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  window.removeEventListener('click', handleGlobalClick)
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="contextMenuStore.isOpen"
        class="fixed z-[99999] min-w-[210px] max-w-xs rounded-2xl bg-slate-900/95 backdrop-blur-xl border border-slate-700/80 p-1.5 shadow-2xl text-xs select-none font-sans"
        :style="{ left: `${contextMenuStore.x}px`, top: `${contextMenuStore.y}px` }"
        @click.stop
        @contextmenu.prevent
      >
        <!-- Header Judul Menu (Opsional) -->
        <div
          v-if="contextMenuStore.title"
          class="px-3 py-1.5 border-b border-slate-800/80 mb-1 text-[11px] font-semibold text-slate-400 tracking-wider uppercase truncate"
        >
          {{ contextMenuStore.title }}
        </div>

        <!-- Daftar Menu Items -->
        <div class="flex flex-col gap-0.5">
          <template v-for="(item, idx) in contextMenuStore.items" :key="idx">
            <!-- Separator Garis Pembatas -->
            <div
              v-if="item.separator"
              class="h-[1px] bg-slate-800/80 my-1 mx-1.5"
            ></div>

            <!-- Tombol Aksi Menu -->
            <button
              v-else
              @click="handleItemClick(item)"
              :disabled="item.disabled"
              class="w-full px-2.5 py-1.5 rounded-xl flex items-center justify-between gap-3 text-left transition group"
              :class="[
                item.disabled
                  ? 'opacity-40 cursor-not-allowed text-slate-500'
                  : item.danger
                  ? 'text-rose-400 hover:bg-rose-500/15 hover:text-rose-300'
                  : 'text-slate-200 hover:bg-sky-500/15 hover:text-sky-300'
              ]"
            >
              <div class="flex items-center gap-2.5 min-w-0">
                <!-- Ikon Komponen Lucide -->
                <component
                  :is="item.icon"
                  v-if="item.icon"
                  class="w-4 h-4 shrink-0 transition"
                  :class="[
                    item.danger
                      ? 'text-rose-400 group-hover:text-rose-300'
                      : 'text-slate-400 group-hover:text-sky-400'
                  ]"
                />
                <span class="truncate font-medium">{{ item.label }}</span>
              </div>

              <!-- Shortcut Hint Text (Opsional) -->
              <span
                v-if="item.shortcut"
                class="text-[10px] font-mono text-slate-500 tracking-tighter shrink-0"
              >
                {{ item.shortcut }}
              </span>
            </button>
          </template>
        </div>
      </div>
    </transition>
  </Teleport>
</template>
