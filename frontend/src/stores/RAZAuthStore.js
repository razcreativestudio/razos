/**
 * RAZOS Auth Store (Pinia)
 * Timestamp: 2026-10-14 13:10:00
 * Deskripsi: State management untuk otentikasi pengguna, sesi login, registrasi awal, dan keamanan RAZOS.
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref } from 'vue'

const TOKEN_KEY = 'razos_auth_token'

export const useAuthStore = defineStore('auth', () => {
  // Status Inisialisasi: Apakah sudah ada akun terdaftar sebelumnya?
  const isInitialized = ref(true)
  // Status Otentikasi: Apakah sesi user saat ini aktif & valid?
  const isAuthenticated = ref(false)
  // Data user yang sedang login
  const currentUser = ref(null)
  // Token bearer
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  // Loading status saat pengecekan auth
  const isLoading = ref(true)
  // Pesan error umum
  const authError = ref('')

  /**
   * Helper headers request API dengan Authorization token
   */
  const getAuthHeaders = () => {
    const headers = { 'Content-Type': 'application/json' }
    if (token.value) {
      headers['Authorization'] = `Bearer ${token.value}`
    }
    return headers
  }

  /**
   * Cek status auth sistem dari backend (apakah sudah ada akun, dan apakah token saat ini valid)
   */
  const checkAuthStatus = async () => {
    isLoading.value = true
    authError.value = ''
    try {
      const res = await fetch('/api/v1/auth/status', {
        headers: getAuthHeaders()
      })
      if (!res.ok) throw new Error('Gagal memeriksa status autentikasi')
      const data = await res.json()

      isInitialized.value = !!data.initialized
      isAuthenticated.value = !!data.authenticated
      if (data.authenticated && data.username) {
        currentUser.value = {
          id: data.user_id || data.user?.id || 'usr-current',
          username: data.username,
          name: data.name || data.user?.name || data.username,
          email: data.email || data.user?.email || `${data.username}@razos.local`,
          role: data.role || data.user?.role || 'admin',
          avatar: data.avatar || data.user?.avatar || '/assets/brand/RAZBrandSquare.webp',
          home_dir: data.home_dir || data.user?.home_dir || '/DATA',
          permissions: data.permissions || data.user?.permissions || {
            apps: ['file-manager', 'task-manager', 'terminal', 'app-store', 'websites', 'databases', 'settings', 'cloud-sync', 'users-manager'],
            storage: { read: true, write: true, delete: true },
            can_install_apps: true,
            can_manage_users: true,
            can_mount_cloud: true,
            can_access_root_terminal: true
          }
        }
      } else {
        currentUser.value = null
        if (!data.authenticated && token.value) {
          // Token expired atau tidak valid
          token.value = ''
          localStorage.removeItem(TOKEN_KEY)
        }
      }
    } catch (err) {
      console.error('[AuthStore] Error cek status:', err)
      // Fallback aman: jika gagal network, pertahankan token bila ada atau set unauth
      authError.value = err.message || 'Koneksi ke backend gagal'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Registrasi akun admin pertama kali (Setup Awal)
   */
  const register = async (username, password) => {
    isLoading.value = true
    authError.value = ''
    try {
      const res = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })

      const data = await res.json()
      const isSuccess = data.success || data.status === 'success'
      if (!res.ok || !isSuccess) {
        throw new Error(data.error || 'Registrasi akun gagal')
      }

      // Berhasil registrasi -> otomatis login
      const receivedToken = data.token || (data.data && data.data.token)
      token.value = receivedToken
      if (receivedToken) localStorage.setItem(TOKEN_KEY, receivedToken)
      isAuthenticated.value = true
      isInitialized.value = true
      
      const u = data.user || {}
      currentUser.value = {
        id: u.id || 'usr-admin-1',
        username: u.username || username,
        name: u.name || username,
        email: u.email || `${username}@razos.local`,
        role: u.role || 'admin',
        avatar: u.avatar || '/assets/brand/RAZBrandSquare.webp',
        home_dir: u.home_dir || '/DATA',
        permissions: u.permissions || {
          apps: ['file-manager', 'task-manager', 'terminal', 'app-store', 'websites', 'databases', 'settings', 'cloud-sync', 'users-manager'],
          storage: { read: true, write: true, delete: true },
          can_install_apps: true,
          can_manage_users: true,
          can_mount_cloud: true,
          can_access_root_terminal: true
        }
      }

      return { success: true }
    } catch (err) {
      authError.value = err.message || 'Registrasi gagal'
      return { success: false, error: authError.value }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Login pengguna
   */
  const login = async (username, password) => {
    isLoading.value = true
    authError.value = ''
    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })

      const data = await res.json()
      const isSuccess = data.success || data.status === 'success'
      if (!res.ok || !isSuccess) {
        throw new Error(data.error || 'Username atau password salah')
      }

      // Berhasil login
      const receivedToken = data.token || (data.data && data.data.token)
      token.value = receivedToken
      if (receivedToken) localStorage.setItem(TOKEN_KEY, receivedToken)
      isAuthenticated.value = true
      
      const u = data.user || {}
      currentUser.value = {
        id: u.id || data.user_id || 'usr-current',
        username: u.username || data.username || username,
        name: u.name || data.name || username,
        email: u.email || `${username}@razos.local`,
        role: u.role || data.role || 'admin',
        avatar: u.avatar || '/assets/brand/RAZBrandSquare.webp',
        home_dir: u.home_dir || '/DATA',
        permissions: u.permissions || {
          apps: ['file-manager', 'task-manager', 'terminal', 'app-store', 'websites', 'databases', 'settings', 'cloud-sync', 'users-manager'],
          storage: { read: true, write: true, delete: true },
          can_install_apps: true,
          can_manage_users: true,
          can_mount_cloud: true,
          can_access_root_terminal: true
        }
      }

      return { success: true }
    } catch (err) {
      authError.value = err.message || 'Login gagal'
      return { success: false, error: authError.value }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logout pengguna dan hapus sesi
   */
  const logout = async () => {
    isLoading.value = true
    try {
      if (token.value) {
        await fetch('/api/v1/auth/logout', {
          method: 'POST',
          headers: getAuthHeaders()
        })
      }
    } catch (err) {
      console.warn('[AuthStore] Error saat logout backend:', err)
    } finally {
      token.value = ''
      localStorage.removeItem(TOKEN_KEY)
      isAuthenticated.value = false
      currentUser.value = null
      isLoading.value = false
    }
  }

  /**
   * Ganti password
   */
  const changePassword = async (oldPassword, newPassword) => {
    try {
      const res = await fetch('/api/v1/auth/change-password', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ oldPassword, newPassword })
      })

      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal mengubah password')
      }
      return { success: true }
    } catch (err) {
      return { success: false, error: err.message || 'Gagal mengubah password' }
    }
  }

  // ========================================================
  // RBAC & PERMISSION HELPERS (ROLE-BASED ACCESS CONTROL)
  // ========================================================
  const isAdmin = () => currentUser.value?.role === 'admin'
  const isOperator = () => currentUser.value?.role === 'operator'
  const isMember = () => currentUser.value?.role === 'member'
  const isGuest = () => currentUser.value?.role === 'guest'

  /**
   * Memeriksa apakah user saat ini diizinkan membuka aplikasi tertentu
   */
  const canAccessApp = (appId) => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    const allowedApps = currentUser.value.permissions?.apps || []
    return allowedApps.includes(appId)
  }

  /**
   * Memeriksa izin menulis/mengunggah/membuat file di storage
   */
  const canWriteStorage = () => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    return currentUser.value.permissions?.storage?.write !== false
  }

  /**
   * Memeriksa izin menghapus file di storage
   */
  const canDeleteStorage = () => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    return currentUser.value.permissions?.storage?.delete === true
  }

  /**
   * Memeriksa izin instalasi aplikasi store / container
   */
  const canInstallApps = () => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    return currentUser.value.permissions?.can_install_apps === true
  }

  /**
   * Memeriksa izin manajemen pengguna
   */
  const canManageUsers = () => {
    if (!currentUser.value) return false
    return isAdmin() || currentUser.value.permissions?.can_manage_users === true
  }

  /**
   * Memeriksa izin mount remote cloud
   */
  const canMountCloud = () => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    return currentUser.value.permissions?.can_mount_cloud === true
  }

  /**
   * Memeriksa izin akses terminal root
   */
  const canAccessRootTerminal = () => {
    if (!currentUser.value) return false
    if (isAdmin()) return true
    return currentUser.value.permissions?.can_access_root_terminal === true
  }

  /**
   * Mengambil folder home pengguna saat ini
   */
  const getUserHomeDir = () => {
    return currentUser.value?.home_dir || '/DATA'
  }

  return {
    isInitialized,
    isAuthenticated,
    currentUser,
    token,
    isLoading,
    authError,
    getAuthHeaders,
    checkAuthStatus,
    register,
    login,
    logout,
    changePassword,
    isAdmin,
    isOperator,
    isMember,
    isGuest,
    canAccessApp,
    canWriteStorage,
    canDeleteStorage,
    canInstallApps,
    canManageUsers,
    canMountCloud,
    canAccessRootTerminal,
    getUserHomeDir
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
