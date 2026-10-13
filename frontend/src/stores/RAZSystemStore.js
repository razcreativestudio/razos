/**
 * RAZOS System Store (Pinia)
 * Updated: 2026-10-08
 * Deskripsi: State management untuk polling data metrik host dan status windowing.
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref } from 'vue'

export const useSystemStore = defineStore('system', () => {
  // Metrik Hardware
  const metrics = ref({
    cpu_percent: 0,
    cpu_temp: 0,
    ram_total: 0,
    ram_used: 0,
    ram_percent: 0,
    zram_percent: 0,
    disk_total: 0,
    disk_used: 0,
    disk_percent: 0
  })

  const isConnected = ref(false)
  const isCleaningRAM = ref(false)
  const isCleaningStorage = ref(false)
  const notification = ref(null)

  // Tampilkan custom notification modal/toast
  function showNotification(title, message, type = 'info') {
    notification.value = { title, message, type }
    setTimeout(() => {
      notification.value = null
    }, 4000)
  }

  // Fetch data metrik dari backend Go API
  async function fetchMetrics() {
    try {
      const res = await fetch('/api/v1/system/metrics')
      if (res.ok) {
        metrics.value = await res.json()
        isConnected.value = true
      } else {
        isConnected.value = false
      }
    } catch (err) {
      isConnected.value = false
    }
  }

  // Eksekusi Clean RAM 1-klik
  async function cleanRAM() {
    isCleaningRAM.value = true
    try {
      const res = await fetch('/api/v1/system/clean-ram', { method: 'POST' })
      if (res.ok) {
        showNotification('Sukses', 'RAM cache berhasil dibersihkan dari kernel!', 'success')
        await fetchMetrics()
      } else {
        showNotification('Gagal', 'Gagal membersihkan cache RAM', 'error')
      }
    } catch (e) {
      showNotification('Error', 'Koneksi ke backend gagal', 'error')
    } finally {
      isCleaningRAM.value = false
    }
  }

  // Eksekusi Clean Storage (Docker Prune) 1-klik
  async function cleanStorage() {
    isCleaningStorage.value = true
    try {
      const res = await fetch('/api/v1/system/clean-storage', { method: 'POST' })
      if (res.ok) {
        showNotification('Sukses', 'Sisa container & storage berhasil dirapikan!', 'success')
        await fetchMetrics()
      } else {
        showNotification('Gagal', 'Gagal mengeksekusi storage cleaner', 'error')
      }
    } catch (e) {
      showNotification('Error', 'Koneksi ke backend gagal', 'error')
    } finally {
      isCleaningStorage.value = false
    }
  }

  // Event Triggers untuk sinkronisasi antar-jendela OS
  const websitesRefreshTrigger = ref(0)
  const databasesRefreshTrigger = ref(0)
  const appsRefreshTrigger = ref(0)

  function triggerWebsitesRefresh() {
    websitesRefreshTrigger.value++
  }

  function triggerDatabasesRefresh() {
    databasesRefreshTrigger.value++
  }

  function triggerAppsRefresh() {
    appsRefreshTrigger.value++
  }

  return {
    metrics,
    isConnected,
    isCleaningRAM,
    isCleaningStorage,
    notification,
    websitesRefreshTrigger,
    databasesRefreshTrigger,
    appsRefreshTrigger,
    triggerWebsitesRefresh,
    triggerDatabasesRefresh,
    triggerAppsRefresh,
    showNotification,
    fetchMetrics,
    cleanRAM,
    cleanStorage
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSystemStore, import.meta.hot))
}
