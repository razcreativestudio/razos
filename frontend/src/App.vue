<script setup>
/**
 * RAZOS Main App Component
 * Updated: 2026-10-08
 * Deskripsi: Root component dengan deteksi breakpoint responsif (Desktop > 768px vs Mobile <= 768px),
 * inisialisasi polling metrik backend, dan custom notification toast.
 */
import { ref, onMounted, onUnmounted } from 'vue'
import { useSystemStore } from './stores/RAZSystemStore'
import RAZDesktopLayout from './layouts/RAZDesktopLayout.vue'
import RAZMobileLayout from './layouts/RAZMobileLayout.vue'
import RAZContextMenu from './components/RAZContextMenu.vue'
import { CheckCircle2, AlertCircle, Info } from 'lucide-vue-next'

const systemStore = useSystemStore()
const isMobile = ref(window.innerWidth <= 768)

function handleResize() {
  isMobile.value = window.innerWidth <= 768
}

let metricsInterval = null

onMounted(() => {
  window.addEventListener('resize', handleResize)
  // Polling metrik setiap 3 detik
  systemStore.fetchMetrics()
  metricsInterval = setInterval(() => {
    systemStore.fetchMetrics()
  }, 3000)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (metricsInterval) clearInterval(metricsInterval)
})
</script>

<template>
  <div class="w-screen h-screen overflow-hidden bg-slate-950 font-sans text-slate-100 antialiased">
    <!-- Render Layout Sesuai Breakpoint 768px -->
    <RAZMobileLayout v-if="isMobile" />
    <RAZDesktopLayout v-else />

    <!-- Custom Notification Toast (Pengganti alert bawaan browser) -->
    <transition
      enter-active-class="transform transition ease-out duration-300"
      enter-from-class="translate-y-4 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="systemStore.notification"
        class="fixed top-5 right-5 z-[9999] flex items-center gap-3 p-3.5 rounded-xl raz-glass border shadow-2xl max-w-sm"
        :class="systemStore.notification.type === 'success'
          ? 'border-emerald-500/40 text-emerald-300'
          : systemStore.notification.type === 'error'
          ? 'border-rose-500/40 text-rose-300'
          : 'border-sky-500/40 text-sky-300'"
      >
        <component
          :is="systemStore.notification.type === 'success' ? CheckCircle2 : systemStore.notification.type === 'error' ? AlertCircle : Info"
          class="w-5 h-5 shrink-0"
        />
        <div>
          <h5 class="text-xs font-bold leading-tight">{{ systemStore.notification.title }}</h5>
          <p class="text-xs text-slate-300 mt-0.5 leading-snug">{{ systemStore.notification.message }}</p>
        </div>
      </div>
    </transition>

    <!-- Global Enterprise Context Menu -->
    <RAZContextMenu />
  </div>
</template>
