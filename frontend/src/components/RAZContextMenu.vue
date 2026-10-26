<script setup>
/**
 * RAZOS Universal Context Menu Component
 * Updated: 2026-10-26 21:20:00 (v4.3.0 - Windows 11 Dark Solid Styling & Identical Parent-Child Submenu Theme)
 * Deskripsi: Komponen visual Custom Context Menu berdesain Windows 11 Dark (#201f27).
 * - Solid Non-Transparent Surface: Warna latar pekat tidak tembus pandang berlebihan sesuai standar Windows 11.
 * - Identical Child Architecture: Submenu anak memiliki dimensi, radius, border, dan warna yang 100% konsisten dengan menu induk.
 * - Windows 11 Radio Dot & Check Indicators: Mendukung status terpilih untuk grup Sort by dan View mode.
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Check, ChevronRight } from 'lucide-vue-next'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'

const contextMenuStore = useContextMenuStore()

// State submenu yang sedang aktif terbuka
const activeSubmenuIndex = ref(null)
let submenuCloseTimer = null

// Deteksi apakah submenu sebaiknya membuka ke kiri karena keterbatasan ruang layar kanan
const shouldOpenSubmenuLeft = computed(() => {
  if (typeof window === 'undefined') return false
  const estimatedParentWidth = 220
  const estimatedSubmenuWidth = 210
  return (contextMenuStore.x + estimatedParentWidth + estimatedSubmenuWidth) > window.innerWidth
})

// Buka submenu saat mouse hover ke parent item
function handleItemMouseEnter(item, idx) {
  if (submenuCloseTimer) {
    clearTimeout(submenuCloseTimer)
    submenuCloseTimer = null
  }
  if (item.children && Array.isArray(item.children) && item.children.length > 0) {
    activeSubmenuIndex.value = idx
  } else {
    activeSubmenuIndex.value = null
  }
}

// Tutup submenu dengan debounce halus agar navigasi kursor nyaman
function handleItemMouseLeave() {
  submenuCloseTimer = setTimeout(() => {
    activeSubmenuIndex.value = null
  }, 180)
}

// Batalkan penutupan submenu saat kursor masuk ke area submenu
function keepSubmenuOpen() {
  if (submenuCloseTimer) {
    clearTimeout(submenuCloseTimer)
    submenuCloseTimer = null
  }
}

// Eksekusi aksi item menu (baik parent tanpa children maupun child item)
function handleItemClick(item) {
  if (item.disabled) return
  // Jika item memiliki children, klik akan membuka/toggle submenu (ramah layar sentuh)
  if (item.children && Array.isArray(item.children) && item.children.length > 0) {
    activeSubmenuIndex.value = activeSubmenuIndex.value === item ? null : item
    return
  }

  // Tutup seluruh menu setelah eksekusi
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
  // Cegah penutupan mendadak jika event berasal dari klik kanan (button 2) atau race-condition saat menu baru saja dibuka
  if (e && e.button === 2) return
  if (contextMenuStore.openedAt && (Date.now() - contextMenuStore.openedAt < 120)) return

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
  if (submenuCloseTimer) {
    clearTimeout(submenuCloseTimer)
  }
})
</script>

<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="transform scale-95 -translate-y-1 opacity-0"
      enter-to-class="transform scale-100 translate-y-0 opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="transform scale-100 translate-y-0 opacity-100"
      leave-to-class="transform scale-95 -translate-y-1 opacity-0"
    >
      <div
        id="raz-custom-context-menu"
        v-if="contextMenuStore.isOpen"
        class="fixed z-[99999] min-w-[210px] max-w-xs rounded-xl bg-[#201f27] border border-white/10 p-1.5 shadow-2xl shadow-black/80 text-xs select-none font-sans text-slate-100"
        :style="{ left: `${contextMenuStore.x}px`, top: `${contextMenuStore.y}px` }"
        @click.stop
        @contextmenu.prevent.stop
      >
        <!-- Header Judul Menu (Opsional) -->
        <div
          v-if="contextMenuStore.title"
          class="px-3 py-1.5 border-b border-white/10 mb-1 text-[11px] font-semibold text-slate-400 tracking-wider uppercase truncate"
        >
          {{ contextMenuStore.title }}
        </div>

        <!-- Daftar Menu Items Utama -->
        <div class="flex flex-col gap-0.5">
          <template v-for="(item, idx) in contextMenuStore.items" :key="idx">
            <!-- Separator Garis Pembatas -->
            <div
              v-if="item.separator"
              class="h-[1px] bg-white/10 my-1 mx-1.5"
            ></div>

            <!-- Quick Action Bar (CasaOS Style: 2 Column Bottom Buttons for Restart & Power) -->
            <div
              v-else-if="item.type === 'quick-actions'"
              class="grid grid-cols-2 gap-1.5 pt-1 mt-1 border-t border-slate-800/80"
            >
              <button
                v-for="(act, aIdx) in item.actions"
                :key="aIdx"
                type="button"
                @click="handleItemClick(act)"
                :disabled="act.disabled"
                class="py-2 px-2.5 rounded-xl flex items-center justify-center gap-1.5 font-semibold text-xs transition active:scale-95 border border-slate-800/60"
                :class="[
                  act.danger
                    ? 'text-rose-400 hover:bg-rose-500/15 hover:border-rose-500/30'
                    : act.success
                    ? 'text-emerald-400 hover:bg-emerald-500/15 hover:border-emerald-500/30'
                    : 'text-slate-300 hover:bg-slate-800 hover:text-white hover:border-slate-700'
                ]"
                :title="act.title || act.label"
              >
                <component :is="act.icon" class="w-4 h-4 shrink-0" />
                <span v-if="act.label" class="text-[11px] truncate">{{ act.label }}</span>
              </button>
            </div>

            <!-- Item Menu Reguler / Item dengan Submenu (Children) -->
            <div
              v-else
              class="relative"
              @mouseenter="handleItemMouseEnter(item, idx)"
              @mouseleave="handleItemMouseLeave"
            >
              <button
                type="button"
                @click="handleItemClick(item)"
                :disabled="item.disabled"
                class="w-full px-2.5 py-1.5 rounded-xl flex items-center justify-between gap-3 text-left transition-all duration-150 group active:scale-[0.98]"
                :class="[
                  item.disabled
                    ? 'opacity-40 cursor-not-allowed text-slate-500'
                    : item.danger
                    ? 'text-rose-400 hover:bg-rose-500/15 hover:text-rose-300'
                    : activeSubmenuIndex === idx
                    ? 'bg-sky-500/20 text-sky-300'
                    : 'text-slate-200 hover:bg-sky-500/15 hover:text-sky-300'
                ]"
              >
                <div class="flex items-center gap-2.5 min-w-0">
                  <!-- Ikon Komponen Lucide -->
                  <component
                    :is="item.icon"
                    v-if="item.icon"
                    class="w-4 h-4 shrink-0 transition-transform duration-150 group-hover:scale-110"
                    :class="[
                      item.danger
                        ? 'text-rose-400 group-hover:text-rose-300'
                        : 'text-slate-400 group-hover:text-sky-400'
                    ]"
                  />
                  <span class="truncate font-medium">{{ item.label }}</span>
                </div>

                <!-- Indikator Kanan: Submenu Arrow, Shortcut Hint, atau Status Checkmark -->
                <div class="flex items-center gap-1.5 shrink-0">
                  <!-- Submenu Arrow (Chevron Right) jika memiliki children -->
                  <ChevronRight
                    v-if="item.children && item.children.length > 0"
                    class="w-3.5 h-3.5 text-slate-500 group-hover:text-sky-300 group-hover:translate-x-0.5 transition"
                  />

                  <!-- Shortcut Hint Text -->
                  <span
                    v-else-if="item.shortcut"
                    class="text-[10px] font-mono text-slate-500 tracking-tighter"
                  >
                    {{ item.shortcut }}
                  </span>

                  <!-- Status Checkmark -->
                  <Check
                    v-else-if="item.checked"
                    class="w-3.5 h-3.5 text-sky-400 shrink-0"
                  />
                </div>
              </button>

              <!-- ============================================== -->
              <!-- SUBMENU POPUP (CHILDREN MENU)                  -->
              <!-- Identik 100% dengan parent, solid non-transparent Windows 11 Dark -->
              <!-- ============================================== -->
              <transition
                enter-active-class="transition duration-120 ease-out"
                enter-from-class="transform scale-95 opacity-0"
                enter-to-class="transform scale-100 opacity-100"
                leave-active-class="transition duration-100 ease-in"
                leave-from-class="transform scale-100 opacity-100"
                leave-to-class="transform scale-95 opacity-0"
              >
                <div
                  v-if="item.children && item.children.length > 0 && activeSubmenuIndex === idx"
                  @mouseenter="keepSubmenuOpen"
                  class="absolute top-0 z-[100000] min-w-[210px] max-w-xs rounded-xl bg-[#201f27] border border-white/10 p-1.5 shadow-2xl shadow-black/80 text-xs text-slate-100 font-sans"
                  :class="[
                    shouldOpenSubmenuLeft
                      ? 'right-full mr-1.5'
                      : 'left-full ml-1.5'
                  ]"
                >
                  <div class="flex flex-col gap-0.5">
                    <template v-for="(child, cIdx) in item.children" :key="cIdx">
                      <!-- Separator di dalam Submenu -->
                      <div
                        v-if="child.separator"
                        class="h-[1px] bg-white/10 my-1 mx-1.5"
                      ></div>

                      <!-- Tombol Child Item -->
                      <button
                        v-else
                        type="button"
                        @click="handleItemClick(child)"
                        :disabled="child.disabled"
                        class="w-full px-2.5 py-1.5 rounded-lg flex items-center justify-between gap-3 text-left transition-all duration-150 group active:scale-[0.98]"
                        :class="[
                          child.disabled
                            ? 'opacity-40 cursor-not-allowed text-slate-500'
                            : child.danger
                            ? 'text-rose-400 hover:bg-rose-500/15 hover:text-rose-300'
                            : child.checked
                            ? 'bg-white/10 text-white font-medium'
                            : 'text-slate-200 hover:bg-white/10 hover:text-white'
                        ]"
                      >
                        <div class="flex items-center gap-2.5 min-w-0">
                          <!-- Radio Bullet Dot atau Ikon Child -->
                          <div
                            v-if="child.checked && !child.icon"
                            class="w-2 h-2 rounded-full bg-sky-400 shrink-0"
                          ></div>
                          <component
                            :is="child.icon"
                            v-else-if="child.icon"
                            class="w-4 h-4 shrink-0 transition-transform duration-150 group-hover:scale-110"
                            :class="[
                              child.danger
                                ? 'text-rose-400 group-hover:text-rose-300'
                                : 'text-slate-400 group-hover:text-sky-400'
                            ]"
                          />
                          <span class="truncate font-medium">{{ child.label }}</span>
                        </div>

                        <!-- Checkmark atau Shortcut Child -->
                        <div class="flex items-center gap-1.5 shrink-0">
                          <span
                            v-if="child.shortcut"
                            class="text-[10px] font-mono text-slate-500 tracking-tighter"
                          >
                            {{ child.shortcut }}
                          </span>
                          <Check
                            v-if="child.checked && child.icon"
                            class="w-3.5 h-3.5 text-sky-400 shrink-0"
                          />
                        </div>
                      </button>
                    </template>
                  </div>
                </div>
              </transition>
            </div>
          </template>
        </div>
      </div>
    </transition>
  </Teleport>
</template>
