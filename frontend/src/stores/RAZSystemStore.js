/**
 * RAZOS System Store (Pinia)
 * Updated: 2026-10-14
 * Deskripsi: State management untuk polling data metrik host, preferensi pengguna sistem (CasaOS style), dan status windowing.
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref } from 'vue'

const SETTINGS_STORAGE_KEY = 'razos_system_settings'

// Default konfigurasi preferensi pengguna
const DEFAULT_SETTINGS = {
  showSearchBar: true,
  searchEngine: 'DuckDuckGo', // 'DuckDuckGo', 'Google', 'Bing', 'Baidu'
  language: 'English', // 'English', 'Bahasa Indonesia'
  webUiPort: 8080,
  wallpaper: '/assets/images/RAZDesktopBg.webp',
  showBlogFeed: true,
  showRecommendedApps: true,
  automountUsb: true,
  widgets: {
    time: true,
    systemStatus: true,
    storageStatus: true,
    networkStatus: true
  },
  account: {
    name: 'raz',
    avatar: '/assets/brand/RAZBrandSquare.webp'
  }
}

export const useSystemStore = defineStore('system', () => {
  // Metrik Hardware Host Riil
  const metrics = ref({
    cpu_percent: 0,
    cpu_temp: 0,
    cpu_model: '',
    cpu_cores: 0,
    os_arch: '',
    os_platform: '',
    device_name: '',
    disk_name: '',
    hostname: '',
    ram_total: 0,
    ram_used: 0,
    ram_percent: 0,
    zram_percent: 0,
    disk_total: 0,
    disk_used: 0,
    disk_percent: 0,
    net_sent: 0,
    net_recv: 0,
    net_interface: 'Wi-Fi / LAN',
    net_upload_speed: '0 KB/s',
    net_download_speed: '0 KB/s'
  })

  // State Pelacak Kecepatan Jaringan
  let lastNetTime = 0
  let lastNetSent = 0
  let lastNetRecv = 0

  const isConnected = ref(false)
  const isCleaningRAM = ref(false)
  const isCleaningStorage = ref(false)
  const notification = ref(null)

  // Inisialisasi Settings dari LocalStorage
  const loadSavedSettings = () => {
    try {
      const saved = localStorage.getItem(SETTINGS_STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        return {
          ...DEFAULT_SETTINGS,
          ...parsed,
          widgets: { ...DEFAULT_SETTINGS.widgets, ...(parsed.widgets || {}) },
          account: { ...DEFAULT_SETTINGS.account, ...(parsed.account || {}) }
        }
      }
    } catch (e) {
      console.warn('Gagal membaca preferensi sistem dari localStorage:', e)
    }
    return { ...DEFAULT_SETTINGS }
  }

  const settings = ref(loadSavedSettings())

  // Simpan pengaturan ke LocalStorage
  function saveSettings() {
    try {
      localStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(settings.value))
    } catch (e) {
      console.error('Gagal menyimpan preferensi sistem:', e)
    }
  }

  // Update nilai konfigurasi umum
  function updateSetting(key, val) {
    settings.value[key] = val
    saveSettings()
  }

  // Update status widget individual
  function updateWidget(widgetKey, val) {
    if (settings.value.widgets) {
      settings.value.widgets[widgetKey] = val
      saveSettings()
    }
  }

  // Update informasi akun pengguna
  function updateAccount(name, avatar) {
    if (name) settings.value.account.name = name
    if (avatar) settings.value.account.avatar = avatar
    saveSettings()
  }

  // Format bytes ke ukuran yang mudah dibaca (kebal NaN & ramah pengguna)
  function formatSpeed(bytesPerSec) {
    const val = Number(bytesPerSec)
    if (!Number.isFinite(val) || isNaN(val) || val <= 0) return '0 KB/s'
    if (val < 1024) return `${Math.round(val)} B/s`
    if (val < 1024 * 1024) return `${(val / 1024).toFixed(1)} KB/s`
    return `${(val / (1024 * 1024)).toFixed(1)} MB/s`
  }

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
        const data = await res.json()
        const now = Date.now()

        // Hitung laju throughput kecepatan internet / LAN secara dinamis
        if (lastNetTime > 0 && now > lastNetTime) {
          const seconds = (now - lastNetTime) / 1000
          if (seconds > 0) {
            const currentSent = Number(data.net_sent) || 0
            const currentRecv = Number(data.net_recv) || 0
            const upDiff = currentSent >= lastNetSent ? (currentSent - lastNetSent) : 0
            const downDiff = currentRecv >= lastNetRecv ? (currentRecv - lastNetRecv) : 0

            data.net_upload_speed = formatSpeed(upDiff / seconds)
            data.net_download_speed = formatSpeed(downDiff / seconds)
          } else {
            data.net_upload_speed = metrics.value.net_upload_speed || '0 KB/s'
            data.net_download_speed = metrics.value.net_download_speed || '0 KB/s'
          }
        } else {
          data.net_upload_speed = '0 KB/s'
          data.net_download_speed = '0 KB/s'
        }

        lastNetTime = now
        lastNetSent = Number(data.net_sent) || 0
        lastNetRecv = Number(data.net_recv) || 0

        metrics.value = data
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

  // Eksekusi Reboot / Restart Sistem Host
  async function restartSystem() {
    try {
      const res = await fetch('/api/v1/system/restart', { method: 'POST' })
      if (res.ok) {
        showNotification('Memulai Ulang', 'Sistem STB sedang melakukan restart...', 'info')
        return true
      }
      showNotification('Gagal', 'Gagal mengirim instruksi restart', 'error')
      return false
    } catch (e) {
      showNotification('Memulai Ulang', 'Instruksi restart telah dikirimkan ke hardware.', 'info')
      return true
    }
  }

  // Eksekusi Shutdown / Matikan Daya Host
  async function shutdownSystem() {
    try {
      const res = await fetch('/api/v1/system/shutdown', { method: 'POST' })
      if (res.ok) {
        showNotification('Mematikan Sistem', 'Sistem STB sedang mematikan daya...', 'info')
        return true
      }
      showNotification('Gagal', 'Gagal mengirim instruksi shutdown', 'error')
      return false
    } catch (e) {
      showNotification('Mematikan Daya', 'Instruksi shutdown telah dikirimkan ke hardware.', 'info')
      return true
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
    settings,
    saveSettings,
    updateSetting,
    updateWidget,
    updateAccount,
    restartSystem,
    shutdownSystem,
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
