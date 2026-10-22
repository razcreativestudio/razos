/**
 * RAZOS External Cloud Synchronization Store (Pinia)
 * Updated: 2026-10-22 11:36:00
 * Deskripsi: State manager terpusat untuk integrasi sinkronisasi cloud eksternal
 * (Nextcloud, ownCloud, WebDAV, rclone remote) dengan File Manager & Local Storage.
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './RAZAuthStore'

export const useCloudStore = defineStore('cloud', () => {
  const authStore = useAuthStore()

  // Daftar akun remote cloud yang terdaftar
  const remotes = ref([])
  // Riwayat log aktivitas sinkronisasi berkas
  const syncLogs = ref([])
  const isLoading = ref(false)
  const isSyncing = ref(false)
  const syncingRemoteId = ref(null)
  const syncProgress = ref({ percent: 0, message: '' })
  const cloudError = ref('')

  // Daftar penyedia (providers) yang didukung
  const supportedProviders = [
    {
      id: 'nextcloud',
      name: 'Nextcloud Hub',
      icon: '/assets/icons/RAZIconCloudNextcloud.svg',
      defaultPort: 443,
      placeholderUrl: 'https://cloud.domain.com/remote.php/dav/files/username/',
      description: 'Platform kolaborasi cloud pribadi lengkap dengan kalender, kontak, dan file sharing terenkripsi.'
    },
    {
      id: 'webdav',
      name: 'WebDAV Server',
      icon: '/assets/icons/RAZIconCloudWebdav.svg',
      defaultPort: 443,
      placeholderUrl: 'https://dav.box.com/dav/ atau http://nas.local:5005/',
      description: 'Protokol standar industri RFC 4918 untuk mengakses dan menyinkronkan file server remote (Synology, QNAP, Box).'
    },
    {
      id: 'googledrive',
      name: 'Google Drive (rclone)',
      icon: '/assets/icons/RAZIconCloudGdrive.svg',
      defaultPort: 443,
      placeholderUrl: 'remote: atau gdrive:/FolderBackup',
      description: 'Penyimpanan cloud Google Drive menggunakan rclone token interaktif.'
    },
    {
      id: 'onedrive',
      name: 'Microsoft OneDrive (rclone)',
      icon: '/assets/icons/RAZIconCloudOnedrive.svg',
      defaultPort: 443,
      placeholderUrl: 'onedrive:/RAZOS_Backups',
      description: 'Sinkronisasi berkas ke Microsoft OneDrive / SharePoint.'
    }
  ]

  /**
   * Mengambil daftar remote cloud dari backend
   */
  async function fetchRemotes() {
    isLoading.value = true
    cloudError.value = ''
    try {
      const res = await fetch('/api/v1/cloud/remotes', {
        headers: authStore.getAuthHeaders()
      })
      if (!res.ok) throw new Error('Gagal memuat akun cloud eksternal')
      const data = await res.json()
      remotes.value = Array.isArray(data.remotes) ? data.remotes : (Array.isArray(data) ? data : [])
      if (Array.isArray(data.sync_logs)) {
        syncLogs.value = data.sync_logs
      }
    } catch (err) {
      console.error('[CloudStore] fetchRemotes error:', err)
      cloudError.value = err.message || 'Gagal memuat remote cloud'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Mengambil riwayat log sinkronisasi
   */
  async function fetchSyncLogs() {
    try {
      const res = await fetch('/api/v1/cloud/logs', {
        headers: authStore.getAuthHeaders()
      })
      if (!res.ok) return
      const data = await res.json()
      syncLogs.value = Array.isArray(data.logs) ? data.logs : []
    } catch (e) {}
  }

  /**
   * Menguji koneksi WebDAV / Nextcloud ke server remote (Test Connection)
   * @param {Object} connectionData { provider, url, username, password }
   */
  async function testConnection(connectionData) {
    isLoading.value = true
    try {
      const res = await fetch('/api/v1/cloud/test-connection', {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(connectionData)
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Uji koneksi gagal. Periksa URL dan kredensial Anda.')
      }
      return { success: true, message: data.message || 'Koneksi ke server remote berhasil diverifikasi!' }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Menambahkan akun remote cloud baru
   * @param {Object} remoteData Konfigurasi cloud remote
   */
  async function addRemote(remoteData) {
    isLoading.value = true
    try {
      const res = await fetch('/api/v1/cloud/remotes', {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(remoteData)
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal menyimpan konfigurasi cloud')
      }
      await fetchRemotes()
      return { success: true, remote: data.remote }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Memperbarui konfigurasi cloud remote
   */
  async function updateRemote(id, remoteData) {
    isLoading.value = true
    try {
      const res = await fetch(`/api/v1/cloud/remotes/${id}`, {
        method: 'PUT',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(remoteData)
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal memperbarui konfigurasi cloud')
      }
      await fetchRemotes()
      return { success: true }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Menghapus akun remote cloud dan melepaskan direktori mount lokal
   */
  async function deleteRemote(id) {
    isLoading.value = true
    try {
      const res = await fetch(`/api/v1/cloud/remotes/${id}`, {
        method: 'DELETE',
        headers: authStore.getAuthHeaders()
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal menghapus remote cloud')
      }
      await fetchRemotes()
      return { success: true }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Memasang (Mount) atau Melepas (Unmount) remote cloud ke sistem berkas lokal
   */
  async function toggleMount(id) {
    const remote = remotes.value.find(r => r.id === id)
    if (!remote) return { success: false, error: 'Remote tidak ditemukan' }
    const action = remote.is_mounted ? 'unmount' : 'mount'
    isLoading.value = true
    try {
      const res = await fetch('/api/v1/cloud/mount', {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify({ id, action })
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || `Gagal ${action} remote cloud`)
      }
      await fetchRemotes()
      return { success: true, is_mounted: data.is_mounted }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Menjalankan sinkronisasi berkas manual atau terjadwal
   * @param {string} id ID remote
   * @param {string} mode 'two-way' | 'upload-only' | 'download-only'
   */
  async function syncRemote(id, mode) {
    isSyncing.value = true
    syncingRemoteId.value = id
    syncProgress.value = { percent: 15, message: 'Menginisialisasi handshake ke server cloud...' }
    try {
      setTimeout(() => {
        if (isSyncing.value) syncProgress.value = { percent: 55, message: 'Memindai berkas baru & membandingkan timestamp...' }
      }, 700)

      setTimeout(() => {
        if (isSyncing.value) syncProgress.value = { percent: 85, message: 'Mentransfer berkas & memperbarui indeks storage...' }
      }, 1500)

      const res = await fetch('/api/v1/cloud/sync', {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify({ id, mode })
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Sinkronisasi gagal')
      }
      syncProgress.value = { percent: 100, message: 'Sinkronisasi selesai 100%!' }
      await fetchRemotes()
      await fetchSyncLogs()
      return { success: true, message: data.message }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      setTimeout(() => {
        isSyncing.value = false
        syncingRemoteId.value = null
        syncProgress.value = { percent: 0, message: '' }
      }, 800)
    }
  }

  return {
    remotes,
    syncLogs,
    isLoading,
    isSyncing,
    syncingRemoteId,
    syncProgress,
    cloudError,
    supportedProviders,
    fetchRemotes,
    fetchSyncLogs,
    testConnection,
    addRemote,
    updateRemote,
    deleteRemote,
    toggleMount,
    syncRemote
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCloudStore, import.meta.hot))
}
