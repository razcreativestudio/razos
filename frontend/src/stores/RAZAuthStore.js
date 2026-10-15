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
        currentUser.value = { username: data.username }
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
      currentUser.value = { username: data.username || (data.user && data.user.username) || username }

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
      currentUser.value = { username: data.username || (data.user && data.user.username) || username }

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
    changePassword
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
