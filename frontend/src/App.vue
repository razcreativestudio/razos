<script setup>
/**
 * RAZOS Main App Component
 * Updated: 2026-10-14 13:14:00
 * Deskripsi: Root component dengan deteksi autentikasi pengguna (Login / Setup Akun),
 * deteksi breakpoint responsif (Desktop > 768px vs Mobile <= 768px),
 * inisialisasi polling metrik backend, dan custom notification toast.
 */
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useSystemStore } from './stores/RAZSystemStore'
import { useAuthStore } from './stores/RAZAuthStore'
import RAZDesktopLayout from './layouts/RAZDesktopLayout.vue'
import RAZMobileLayout from './layouts/RAZMobileLayout.vue'
import RAZAuthView from './views/RAZAuthView.vue'
import RAZContextMenu from './components/RAZContextMenu.vue'
import { CheckCircle2, AlertCircle, Info, Loader2 } from 'lucide-vue-next'

const systemStore = useSystemStore()
const authStore = useAuthStore()
const isMobile = ref(window.innerWidth <= 768)

function handleResize() {
  isMobile.value = window.innerWidth <= 768
}

let metricsInterval = null

const startMetricsPolling = () => {
  if (metricsInterval) clearInterval(metricsInterval)
  systemStore.fetchMetrics()
  metricsInterval = setInterval(() => {
    systemStore.fetchMetrics()
  }, 3000)
}

const stopMetricsPolling = () => {
  if (metricsInterval) {
    clearInterval(metricsInterval)
    metricsInterval = null
  }
}

// Pantau perubahan status login pengguna
watch(() => authStore.isAuthenticated, (newVal) => {
  if (newVal) {
    startMetricsPolling()
  } else {
    stopMetricsPolling()
  }
})

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  // Periksa status auth awal
  await authStore.checkAuthStatus()
  if (authStore.isAuthenticated) {
    startMetricsPolling()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  stopMetricsPolling()
})
</script>

<template>
  <div class="w-screen h-screen overflow-hidden bg-slate-950 font-sans text-slate-100 antialiased select-none">
    
    <!-- State Loading Pengecekan Autentikasi Sistem -->
    <div 
      v-if="authStore.isLoading" 
      class="w-full h-full flex flex-col items-center justify-center bg-slate-950 text-slate-400 gap-3"
    >
      <Loader2 class="w-8 h-8 animate-spin text-blue-500" />
      <span class="text-xs font-medium tracking-wide">Memuat Sistem RAZOS...</span>
    </div>

    <!-- Tampilan Form Login / Registrasi Akun Awal jika belum login -->
    <RAZAuthView v-else-if="!authStore.isAuthenticated" />

    <!-- Tampilan Desktop / Mobile jika sudah login -->
    <template v-else>
      <!-- Render Layout Sesuai Breakpoint 768px -->
      <RAZMobileLayout v-if="isMobile" />
      <RAZDesktopLayout v-else />

      <!-- Global Enterprise Context Menu -->
      <RAZContextMenu />
    </template>

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
  </div>
</template>
