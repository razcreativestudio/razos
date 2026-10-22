/**
 * RAZOS Multi-User & Role-Based Access Control (RBAC) Store (Pinia)
 * Updated: 2026-10-22 11:35:00
 * Deskripsi: State manager terpusat untuk mengelola banyak pengguna (Multi-User) dan
 * penegakan hierarki izin akses peran (RBAC: Admin, Operator, Member, Guest).
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './RAZAuthStore'

export const useUserStore = defineStore('user', () => {
  const authStore = useAuthStore()

  // Daftar seluruh akun pengguna sistem
  const users = ref([])
  const isLoading = ref(false)
  const userError = ref('')

  // Template izin bawaan untuk masing-masing peran (Default RBAC Matrix)
  const defaultRoleTemplates = {
    admin: {
      role: 'admin',
      label: 'Administrator (Superuser)',
      description: 'Akses penuh tanpa batas ke seluruh modul sistem, terminal root, dan manajemen pengguna.',
      permissions: {
        apps: ['file-manager', 'task-manager', 'terminal', 'app-store', 'websites', 'databases', 'settings', 'cloud-sync', 'users-manager'],
        storage: { read: true, write: true, delete: true },
        can_install_apps: true,
        can_manage_users: true,
        can_mount_cloud: true,
        can_access_root_terminal: true
      }
    },
    operator: {
      role: 'operator',
      label: 'Operator / DevOps Manager',
      description: 'Dapat mengelola Website, Database, App Store, File Storage, dan Cloud Sync. Dibatasi dari Terminal root.',
      permissions: {
        apps: ['file-manager', 'task-manager', 'app-store', 'websites', 'databases', 'cloud-sync'],
        storage: { read: true, write: true, delete: true },
        can_install_apps: true,
        can_manage_users: false,
        can_mount_cloud: true,
        can_access_root_terminal: false
      }
    },
    member: {
      role: 'member',
      label: 'Member / Standard User',
      description: 'Dapat menggunakan File Manager (folder pribadi /DATA/users/<user>), Text Editor, dan memantau status aplikasi.',
      permissions: {
        apps: ['file-manager', 'task-manager'],
        storage: { read: true, write: true, delete: false },
        can_install_apps: false,
        can_manage_users: false,
        can_mount_cloud: false,
        can_access_root_terminal: false
      }
    },
    guest: {
      role: 'guest',
      label: 'Guest / Tamu (Read-Only)',
      description: 'Hanya-baca (Read-Only) pada file publik/shared dan widget metrik. Dilarang mengubah konfigurasi atau menghapus file.',
      permissions: {
        apps: ['file-manager', 'task-manager'],
        storage: { read: true, write: false, delete: false },
        can_install_apps: false,
        can_manage_users: false,
        can_mount_cloud: false,
        can_access_root_terminal: false
      }
    }
  }

  /**
   * Mengambil daftar seluruh pengguna dari backend REST API
   */
  async function fetchUsers() {
    isLoading.value = true
    userError.value = ''
    try {
      const res = await fetch('/api/v1/users', {
        headers: authStore.getAuthHeaders()
      })
      if (!res.ok) throw new Error('Gagal memuat daftar pengguna sistem')
      const data = await res.json()
      users.value = Array.isArray(data.users) ? data.users : (Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('[UserStore] fetchUsers error:', err)
      userError.value = err.message || 'Gagal memuat pengguna'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Membuat pengguna baru dengan konfigurasi peran dan izin spesifik
   * @param {Object} payload { username, password, name, email, role, home_dir, permissions }
   */
  async function createUser(payload) {
    isLoading.value = true
    try {
      const res = await fetch('/api/v1/users', {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(payload)
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal membuat pengguna baru')
      }
      await fetchUsers()
      return { success: true, user: data.user }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Memperbarui informasi atau hak akses pengguna yang ada
   * @param {string} id ID pengguna (contoh: 'usr-admin-1')
   * @param {Object} payload Data perubahan (name, email, role, status, permissions, password)
   */
  async function updateUser(id, payload) {
    isLoading.value = true
    try {
      const res = await fetch(`/api/v1/users/${id}`, {
        method: 'PUT',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(payload)
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal memperbarui pengguna')
      }
      await fetchUsers()
      // Jika user yang diupdate adalah user yang sedang login, update data di authStore
      if (authStore.currentUser?.id === id || authStore.currentUser?.username === payload.username) {
        await authStore.checkAuthStatus()
      }
      return { success: true }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Menghapus pengguna dari sistem
   * @param {string} id ID pengguna
   */
  async function deleteUser(id) {
    isLoading.value = true
    try {
      const res = await fetch(`/api/v1/users/${id}`, {
        method: 'DELETE',
        headers: authStore.getAuthHeaders()
      })
      const data = await res.json()
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Gagal menghapus pengguna')
      }
      await fetchUsers()
      return { success: true }
    } catch (err) {
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Beralih ke akun pengguna lain secara instan (Switch User)
   * @param {string} username Username tujuan
   * @param {string} password Password pengguna
   */
  async function switchUser(username, password) {
    return await authStore.login(username, password)
  }

  return {
    users,
    isLoading,
    userError,
    defaultRoleTemplates,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
    switchUser
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
}
