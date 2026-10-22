<script setup>
/**
 * RAZOS Users & Role-Based Access Control (RBAC) Manager View
 * Updated: 2026-10-22 11:40:00
 * Deskripsi: Komponen antarmuka enterprise untuk mengelola akun pengguna multi-user,
 * penugasan peran (Admin, Operator, Member, Guest), dan konfigurasi izin akses granular per-aplikasi dan per-storage.
 */

import { ref, computed, onMounted } from 'vue'
import {
  Users, UserPlus, Shield, ShieldCheck, ShieldAlert, Key, Edit, Trash2,
  Check, X, Search, Lock, Unlock, Folder, Terminal, Box, Globe, Database,
  Settings, RefreshCw, UserCheck, AlertTriangle, ArrowRightLeft, Sparkles
} from 'lucide-vue-next'
import { useUserStore } from '../stores/RAZUserStore'
import { useAuthStore } from '../stores/RAZAuthStore'
import { useSystemStore } from '../stores/RAZSystemStore'

const userStore = useUserStore()
const authStore = useAuthStore()
const systemStore = useSystemStore()

const searchQuery = ref('')
const filterRole = ref('all') // 'all' | 'admin' | 'operator' | 'member' | 'guest'

// State Modal Tambah Pengguna
const showAddUserModal = ref(false)
const newUserForm = ref({
  username: '',
  name: '',
  email: '',
  password: '',
  role: 'member',
  home_dir: '/DATA/users/',
  permissions: {
    apps: ['file-manager', 'task-manager'],
    storage: { read: true, write: true, delete: false },
    can_install_apps: false,
    can_manage_users: false,
    can_mount_cloud: false,
    can_access_root_terminal: false
  }
})

// State Modal Edit Pengguna
const showEditUserModal = ref(false)
const editingUser = ref(null)

// State Modal Reset Password
const showPasswordModal = ref(false)
const passwordTargetUser = ref(null)
const newPasswordInput = ref('')

// State Modal Konfirmasi Hapus
const showDeleteModal = ref(false)
const userToDelete = ref(null)

// State Modal Switch User
const showSwitchModal = ref(false)
const switchTargetUsername = ref('')
const switchPassword = ref('')
const switchError = ref('')

onMounted(async () => {
  await userStore.fetchUsers()
})

// Filter daftar user berdasarkan query pencarian dan peran
const filteredUsers = computed(() => {
  let list = userStore.users || []
  if (filterRole.value !== 'all') {
    list = list.filter(u => u.role === filterRole.value)
  }
  const q = searchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter(u =>
      (u.username && u.username.toLowerCase().includes(q)) ||
      (u.name && u.name.toLowerCase().includes(q)) ||
      (u.email && u.email.toLowerCase().includes(q)) ||
      (u.role && u.role.toLowerCase().includes(q))
    )
  }
  return list
})

// Helper warna & badge peran
function getRoleBadge(role) {
  switch (role) {
    case 'admin':
      return {
        label: 'Administrator',
        bg: 'bg-purple-500/20 text-purple-300 border-purple-500/40',
        icon: ShieldCheck
      }
    case 'operator':
      return {
        label: 'Operator / DevOps',
        bg: 'bg-sky-500/20 text-sky-300 border-sky-500/40',
        icon: Shield
      }
    case 'member':
      return {
        label: 'Member',
        bg: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40',
        icon: UserCheck
      }
    case 'guest':
    default:
      return {
        label: 'Guest (Read-Only)',
        bg: 'bg-amber-500/20 text-amber-300 border-amber-500/40',
        icon: ShieldAlert
      }
  }
}

// Buka Modal Tambah User dengan Preset Peran Otomatis
function openAddUser() {
  newUserForm.value = {
    username: '',
    name: '',
    email: '',
    password: '',
    role: 'member',
    home_dir: '/DATA/users/',
    permissions: JSON.parse(JSON.stringify(userStore.defaultRoleTemplates.member.permissions))
  }
  showAddUserModal.value = true
}

// Update form izin otomatis saat pilihan dropdown role berganti
function handleRolePresetChange(isEdit = false) {
  const target = isEdit ? editingUser.value : newUserForm.value
  const template = userStore.defaultRoleTemplates[target.role]
  if (template) {
    target.permissions = JSON.parse(JSON.stringify(template.permissions))
    if (!isEdit && target.username) {
      target.home_dir = target.role === 'admin' ? '/DATA' : `/DATA/users/${target.username}`
    }
  }
}

// Simpan Akun Pengguna Baru
async function submitCreateUser() {
  const f = newUserForm.value
  if (!f.username || f.username.length < 3) {
    systemStore.showNotification('Validasi', 'Username minimal 3 karakter!', 'warning')
    return
  }
  if (!f.password || f.password.length < 4) {
    systemStore.showNotification('Validasi', 'Password minimal 4 karakter!', 'warning')
    return
  }

  const payload = {
    username: f.username.trim().toLowerCase(),
    password: f.password.trim(),
    name: f.name.trim() || f.username.trim(),
    email: f.email.trim() || `${f.username}@razos.local`,
    role: f.role,
    home_dir: f.role === 'admin' ? '/DATA' : `/DATA/users/${f.username.trim().toLowerCase()}`,
    permissions: f.permissions
  }

  const res = await userStore.createUser(payload)
  if (res.success) {
    systemStore.showNotification('Sukses', `Pengguna ${payload.name} berhasil dibuat!`, 'success')
    showAddUserModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Terjadi kesalahan', 'error')
  }
}

// Buka Modal Edit User
function openEditUser(user) {
  editingUser.value = JSON.parse(JSON.stringify(user))
  showEditUserModal.value = true
}

// Simpan Perubahan Pengguna
async function submitEditUser() {
  if (!editingUser.value) return
  const u = editingUser.value
  const payload = {
    name: u.name,
    email: u.email,
    role: u.role,
    status: u.status,
    home_dir: u.home_dir,
    permissions: u.permissions
  }

  const res = await userStore.updateUser(u.id, payload)
  if (res.success) {
    systemStore.showNotification('Sukses', `Pengaturan akun ${u.name} berhasil diperbarui!`, 'success')
    showEditUserModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal menyimpan perubahan', 'error')
  }
}

// Buka Modal Reset Password
function openResetPassword(user) {
  passwordTargetUser.value = user
  newPasswordInput.value = ''
  showPasswordModal.value = true
}

async function submitResetPassword() {
  if (!passwordTargetUser.value || !newPasswordInput.value) return
  if (newPasswordInput.value.length < 4) {
    systemStore.showNotification('Validasi', 'Password minimal 4 karakter!', 'warning')
    return
  }

  const res = await userStore.updateUser(passwordTargetUser.value.id, {
    password: newPasswordInput.value.trim()
  })
  if (res.success) {
    systemStore.showNotification('Sukses', `Password untuk ${passwordTargetUser.value.username} berhasil diubah!`, 'success')
    showPasswordModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal mereset password', 'error')
  }
}

// Konfirmasi Hapus User
function openDeleteUser(user) {
  userToDelete.value = user
  showDeleteModal.value = true
}

async function submitDeleteUser() {
  if (!userToDelete.value) return
  const res = await userStore.deleteUser(userToDelete.value.id)
  if (res.success) {
    systemStore.showNotification('Dihapus', `Pengguna ${userToDelete.value.username} telah dihapus dari sistem.`, 'info')
    showDeleteModal.value = false
  } else {
    systemStore.showNotification('Gagal', res.error || 'Gagal menghapus pengguna', 'error')
  }
}

// Buka Dialog Cepat Switch User
function openSwitchUser(targetUsername = '') {
  switchTargetUsername.value = targetUsername
  switchPassword.value = ''
  switchError.value = ''
  showSwitchModal.value = true
}

async function submitSwitchUser() {
  if (!switchTargetUsername.value || !switchPassword.value) {
    switchError.value = 'Silakan isi username dan password!'
    return
  }
  const res = await userStore.switchUser(switchTargetUsername.value, switchPassword.value)
  if (res.success) {
    systemStore.showNotification('Sesi Berubah', `Berhasil masuk sebagai ${switchTargetUsername.value}`, 'success')
    showSwitchModal.value = false
    await userStore.fetchUsers()
  } else {
    switchError.value = res.error || 'Gagal beralih pengguna'
  }
}
</script>

<template>
  <div class="w-full h-full flex flex-col bg-slate-950/90 text-slate-100 select-none overflow-hidden font-sans">
    <!-- Header Toolbar -->
    <div class="h-14 border-b border-slate-800/80 px-5 flex items-center justify-between shrink-0 bg-slate-900/60 backdrop-blur-md">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-purple-500/20 border border-purple-500/30 flex items-center justify-center text-purple-400">
          <Users class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-sm font-bold text-slate-100 flex items-center gap-2">
            Manajemen Pengguna & Izin Akses Peran (RBAC)
            <span class="text-[10px] px-2 py-0.5 rounded-full bg-slate-800 text-slate-400 border border-slate-700/60 font-mono">
              {{ filteredUsers.length }} User
            </span>
          </h2>
          <p class="text-[11px] text-slate-400">Kelola hierarki multi-user, pembatasan direktori, dan kontrol akses per-aplikasi</p>
        </div>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="userStore.fetchUsers"
          class="p-2 rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700/60 transition active:scale-95"
          title="Muat Ulang Data"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': userStore.isLoading }" />
        </button>

        <button
          v-if="authStore.canManageUsers()"
          type="button"
          @click="openAddUser"
          class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-semibold text-xs transition active:scale-95 shadow-lg shadow-purple-600/20 border border-purple-400/30"
        >
          <UserPlus class="w-4 h-4" />
          <span>Tambah Pengguna</span>
        </button>
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="px-5 py-3 border-b border-slate-800/60 flex flex-wrap items-center justify-between gap-3 bg-slate-900/30">
      <!-- Search Box -->
      <div class="relative flex-1 min-w-[200px] max-w-sm">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama, username, atau email..."
          class="w-full bg-slate-900/90 border border-slate-800 rounded-xl pl-9 pr-3 py-1.5 text-xs text-slate-100 placeholder-slate-500 focus:outline-none focus:border-purple-500 transition"
        />
      </div>

      <!-- Filter Buttons -->
      <div class="flex items-center gap-1.5 p-1 bg-slate-900/80 rounded-xl border border-slate-800/80 text-xs">
        <button
          v-for="roleKey in ['all', 'admin', 'operator', 'member', 'guest']"
          :key="roleKey"
          type="button"
          @click="filterRole = roleKey"
          class="px-2.5 py-1 rounded-lg font-medium capitalize transition"
          :class="filterRole === roleKey ? 'bg-purple-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200'"
        >
          {{ roleKey }}
        </button>
      </div>
    </div>

    <!-- Main Content: Users List Table / Cards -->
    <div class="flex-1 overflow-y-auto p-5 scrollbar-thin">
      <div v-if="filteredUsers.length === 0" class="h-64 flex flex-col items-center justify-center text-center text-slate-400">
        <Users class="w-12 h-12 text-slate-600 mb-2" />
        <p class="text-sm font-medium">Tidak ada pengguna yang sesuai dengan filter pencarian</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="user in filteredUsers"
          :key="user.id"
          class="rounded-2xl bg-slate-900/70 border border-slate-800/90 hover:border-purple-500/40 p-4 transition duration-200 flex flex-col justify-between shadow-xl relative group"
        >
          <!-- Current User Indicator Glow -->
          <div
            v-if="authStore.currentUser?.username === user.username"
            class="absolute top-3 right-3 flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/20 border border-emerald-500/40 text-[10px] font-bold text-emerald-300"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>Akun Anda</span>
          </div>

          <div>
            <!-- Top Avatar & Name -->
            <div class="flex items-center gap-3 mb-3">
              <img
                :src="user.avatar || '/assets/brand/RAZBrandSquare.webp'"
                alt="Avatar"
                class="w-11 h-11 rounded-xl object-cover border border-slate-700/80 shadow-md bg-slate-800"
              />
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-1.5">
                  <h3 class="text-sm font-bold text-slate-100 truncate">{{ user.name || user.username }}</h3>
                </div>
                <p class="text-[11px] text-slate-400 font-mono truncate">@{{ user.username }}</p>
                <p class="text-[10px] text-slate-500 truncate">{{ user.email }}</p>
              </div>
            </div>

            <!-- Role Badge & Home Dir -->
            <div class="flex flex-wrap items-center gap-1.5 mb-3">
              <span
                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-lg border text-[10px] font-bold tracking-wide"
                :class="getRoleBadge(user.role).bg"
              >
                <component :is="getRoleBadge(user.role).icon" class="w-3 h-3" />
                <span>{{ getRoleBadge(user.role).label }}</span>
              </span>

              <span
                v-if="user.status === 'disabled'"
                class="px-2 py-0.5 rounded-lg border border-rose-500/40 bg-rose-500/20 text-rose-300 text-[10px] font-bold"
              >
                Dinonaktifkan
              </span>
            </div>

            <!-- Directory Home & Quota Info -->
            <div class="p-2.5 rounded-xl bg-slate-950/60 border border-slate-800/80 space-y-1.5 text-[11px] mb-3">
              <div class="flex items-center justify-between text-slate-400">
                <span class="flex items-center gap-1.5">
                  <Folder class="w-3.5 h-3.5 text-sky-400" />
                  <span>Home Dir:</span>
                </span>
                <span class="font-mono text-slate-200 truncate max-w-[150px]">{{ user.home_dir || '/DATA' }}</span>
              </div>
              <div class="flex items-center justify-between text-slate-400">
                <span class="flex items-center gap-1.5">
                  <Shield class="w-3.5 h-3.5 text-amber-400" />
                  <span>Izin Storage:</span>
                </span>
                <span class="text-slate-300 font-medium">
                  {{ user.permissions?.storage?.write ? (user.permissions?.storage?.delete ? 'Full R/W/D' : 'Read/Write') : 'Read-Only' }}
                </span>
              </div>
            </div>

            <!-- Allowed Apps Badges -->
            <div class="space-y-1 mb-2">
              <span class="text-[10px] text-slate-400 uppercase font-semibold tracking-wider">Aplikasi Diizinkan:</span>
              <div class="flex flex-wrap gap-1">
                <span
                  v-for="appId in (user.permissions?.apps || [])"
                  :key="appId"
                  class="px-1.5 py-0.5 rounded-md bg-slate-800/80 border border-slate-700/60 text-[10px] text-slate-300 font-medium capitalize"
                >
                  {{ appId.replace('-', ' ') }}
                </span>
              </div>
            </div>
          </div>

          <!-- Bottom Action Buttons -->
          <div class="pt-3 border-t border-slate-800/80 flex items-center justify-between gap-1.5">
            <!-- Switch to this user button -->
            <button
              v-if="authStore.currentUser?.username !== user.username"
              type="button"
              @click="openSwitchUser(user.username)"
              class="flex items-center gap-1 px-2 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-medium transition active:scale-95"
              title="Beralih ke akun ini"
            >
              <ArrowRightLeft class="w-3 h-3 text-sky-400" />
              <span>Switch</span>
            </button>
            <div v-else class="text-[11px] text-emerald-400 font-medium">Aktif Saat Ini</div>

            <!-- Edit, Key, Delete (Hanya jika admin atau punya hak manage_users) -->
            <div v-if="authStore.canManageUsers()" class="flex items-center gap-1">
              <button
                type="button"
                @click="openResetPassword(user)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-amber-300 transition"
                title="Ganti Password"
              >
                <Key class="w-3.5 h-3.5" />
              </button>

              <button
                type="button"
                @click="openEditUser(user)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-sky-300 transition"
                title="Edit Hak Akses / Peran"
              >
                <Edit class="w-3.5 h-3.5" />
              </button>

              <button
                v-if="user.username !== 'raz' && authStore.currentUser?.username !== user.username"
                type="button"
                @click="openDeleteUser(user)"
                class="p-1.5 rounded-lg bg-slate-800 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 transition"
                title="Hapus Pengguna"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL 1: TAMBAH PENGGUNA BARU (CUSTOM GLASSMORPHISM) -->
    <!-- ==================================================== -->
    <div
      v-if="showAddUserModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-6 text-slate-100 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <UserPlus class="w-5 h-5 text-purple-400" />
            <h3 class="text-sm font-bold text-slate-100">Tambah Akun Pengguna Baru</h3>
          </div>
          <button @click="showAddUserModal = false" class="text-slate-400 hover:text-white p-1 rounded-lg">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitCreateUser" class="space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Username *</label>
              <input
                v-model="newUserForm.username"
                type="text"
                required
                placeholder="misal: devops_john"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 placeholder-slate-600 focus:outline-none focus:border-purple-500"
              />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Nama Lengkap</label>
              <input
                v-model="newUserForm.name"
                type="text"
                placeholder="misal: John Doe"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 placeholder-slate-600 focus:outline-none focus:border-purple-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Email</label>
              <input
                v-model="newUserForm.email"
                type="email"
                placeholder="john@example.com"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 placeholder-slate-600 focus:outline-none focus:border-purple-500"
              />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Password Awal *</label>
              <input
                v-model="newUserForm.password"
                type="password"
                required
                placeholder="Minimal 4 karakter"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 placeholder-slate-600 focus:outline-none focus:border-purple-500"
              />
            </div>
          </div>

          <!-- Pilihan Peran (Role RBAC) -->
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1.5">Peran Sistem (Role)</label>
            <div class="grid grid-cols-2 gap-2">
              <label
                v-for="(tpl, roleKey) in userStore.defaultRoleTemplates"
                :key="roleKey"
                class="p-2.5 rounded-xl border cursor-pointer transition flex flex-col gap-1"
                :class="newUserForm.role === roleKey ? 'bg-purple-950/40 border-purple-500 text-purple-200' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
              >
                <div class="flex items-center justify-between">
                  <span class="font-bold text-xs capitalize">{{ roleKey }}</span>
                  <input
                    type="radio"
                    name="newRole"
                    :value="roleKey"
                    v-model="newUserForm.role"
                    @change="handleRolePresetChange(false)"
                    class="accent-purple-500"
                  />
                </div>
                <span class="text-[10px] leading-tight opacity-80">{{ tpl.label }}</span>
              </label>
            </div>
          </div>

          <!-- Konfigurasi Izin Granular Per-Aplikasi -->
          <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 space-y-2.5">
            <span class="block text-[11px] font-semibold text-purple-300">Izin Modul Aplikasi:</span>
            <div class="grid grid-cols-3 gap-2">
              <label
                v-for="app in [
                  { id: 'file-manager', label: 'File Manager' },
                  { id: 'task-manager', label: 'Task Manager' },
                  { id: 'terminal', label: 'Web Terminal' },
                  { id: 'app-store', label: 'App Store' },
                  { id: 'websites', label: 'Websites' },
                  { id: 'databases', label: 'Databases' },
                  { id: 'settings', label: 'Settings' },
                  { id: 'cloud-sync', label: 'Cloud Sync' }
                ]"
                :key="app.id"
                class="flex items-center gap-2 cursor-pointer text-[11px] text-slate-300"
              >
                <input
                  type="checkbox"
                  :value="app.id"
                  v-model="newUserForm.permissions.apps"
                  class="rounded accent-purple-500"
                />
                <span>{{ app.label }}</span>
              </label>
            </div>

            <!-- Izin Operasi Storage -->
            <div class="pt-2 border-t border-slate-800/80 flex items-center justify-between text-[11px]">
              <span class="text-slate-400">Hak Akses Berkas:</span>
              <div class="flex items-center gap-3">
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="newUserForm.permissions.storage.read" class="accent-purple-500" />
                  <span>Read</span>
                </label>
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="newUserForm.permissions.storage.write" class="accent-purple-500" />
                  <span>Write</span>
                </label>
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="newUserForm.permissions.storage.delete" class="accent-purple-500" />
                  <span>Delete</span>
                </label>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="showAddUserModal = false"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="userStore.isLoading"
              class="px-4 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-semibold flex items-center gap-1.5 shadow-lg shadow-purple-600/30"
            >
              <Check class="w-4 h-4" />
              <span>Simpan Pengguna</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL 2: EDIT PENGGUNA & HAK AKSES                   -->
    <!-- ==================================================== -->
    <div
      v-if="showEditUserModal && editingUser"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-6 text-slate-100 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <Edit class="w-5 h-5 text-sky-400" />
            <h3 class="text-sm font-bold text-slate-100">Edit Pengguna: @{{ editingUser.username }}</h3>
          </div>
          <button @click="showEditUserModal = false" class="text-slate-400 hover:text-white p-1 rounded-lg">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitEditUser" class="space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Nama Lengkap</label>
              <input
                v-model="editingUser.name"
                type="text"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Email</label>
              <input
                v-model="editingUser.email"
                type="email"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Peran (Role)</label>
              <select
                v-model="editingUser.role"
                @change="handleRolePresetChange(true)"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500 capitalize"
              >
                <option value="admin">Admin (Superuser)</option>
                <option value="operator">Operator / DevOps</option>
                <option value="member">Member</option>
                <option value="guest">Guest (Read-Only)</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-slate-300 mb-1">Status Akun</label>
              <select
                v-model="editingUser.status"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
              >
                <option value="active">Aktif</option>
                <option value="disabled">Dinonaktifkan</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">Direktori Home (Root Path Scoping)</label>
            <input
              v-model="editingUser.home_dir"
              type="text"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 font-mono text-slate-100 focus:outline-none focus:border-sky-500"
            />
          </div>

          <!-- Konfigurasi Izin Granular Per-Aplikasi -->
          <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 space-y-2.5">
            <span class="block text-[11px] font-semibold text-sky-300">Izin Modul Aplikasi:</span>
            <div class="grid grid-cols-3 gap-2">
              <label
                v-for="app in [
                  { id: 'file-manager', label: 'File Manager' },
                  { id: 'task-manager', label: 'Task Manager' },
                  { id: 'terminal', label: 'Web Terminal' },
                  { id: 'app-store', label: 'App Store' },
                  { id: 'websites', label: 'Websites' },
                  { id: 'databases', label: 'Databases' },
                  { id: 'settings', label: 'Settings' },
                  { id: 'cloud-sync', label: 'Cloud Sync' }
                ]"
                :key="app.id"
                class="flex items-center gap-2 cursor-pointer text-[11px] text-slate-300"
              >
                <input
                  type="checkbox"
                  :value="app.id"
                  v-model="editingUser.permissions.apps"
                  class="rounded accent-sky-500"
                />
                <span>{{ app.label }}</span>
              </label>
            </div>

            <!-- Izin Operasi Storage -->
            <div class="pt-2 border-t border-slate-800/80 flex items-center justify-between text-[11px]">
              <span class="text-slate-400">Hak Akses Berkas:</span>
              <div class="flex items-center gap-3">
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="editingUser.permissions.storage.read" class="accent-sky-500" />
                  <span>Read</span>
                </label>
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="editingUser.permissions.storage.write" class="accent-sky-500" />
                  <span>Write</span>
                </label>
                <label class="flex items-center gap-1.5 cursor-pointer">
                  <input type="checkbox" v-model="editingUser.permissions.storage.delete" class="accent-sky-500" />
                  <span>Delete</span>
                </label>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="showEditUserModal = false"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="userStore.isLoading"
              class="px-4 py-2 rounded-xl bg-sky-600 hover:bg-sky-500 text-white font-semibold flex items-center gap-1.5 shadow-lg shadow-sky-600/30"
            >
              <Check class="w-4 h-4" />
              <span>Simpan Perubahan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL 3: GANTI PASSWORD PENGGUNA                     -->
    <!-- ==================================================== -->
    <div
      v-if="showPasswordModal && passwordTargetUser"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-5 text-slate-100 space-y-4">
        <div class="flex items-center justify-between pb-2 border-b border-slate-800">
          <div class="flex items-center gap-2">
            <Key class="w-4 h-4 text-amber-400" />
            <h3 class="text-xs font-bold text-slate-100">Ganti Password: @{{ passwordTargetUser.username }}</h3>
          </div>
          <button @click="showPasswordModal = false" class="text-slate-400 hover:text-white p-1">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitResetPassword" class="space-y-3 text-xs">
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">Password Baru</label>
            <input
              v-model="newPasswordInput"
              type="password"
              required
              placeholder="Minimal 4 karakter"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-amber-500"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="showPasswordModal = false"
              class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-3.5 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs"
            >
              Perbarui Password
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL 4: KONFIRMASI HAPUS PENGGUNA                   -->
    <!-- ==================================================== -->
    <div
      v-if="showDeleteModal && userToDelete"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-rose-500/30 shadow-2xl p-5 text-slate-100 space-y-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-rose-500/20 border border-rose-500/30 flex items-center justify-center text-rose-400 shrink-0">
            <AlertTriangle class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-slate-100">Hapus Akun Pengguna?</h3>
            <p class="text-xs text-slate-400">Akun @{{ userToDelete.username }} ({{ userToDelete.name }}) akan dihapus secara permanen dari sistem.</p>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
          <button
            type="button"
            @click="showDeleteModal = false"
            class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs"
          >
            Batal
          </button>
          <button
            type="button"
            @click="submitDeleteUser"
            class="px-3.5 py-1.5 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs shadow-lg shadow-rose-600/30"
          >
            Hapus Pengguna
          </button>
        </div>
      </div>
    </div>

    <!-- ==================================================== -->
    <!-- MODAL 5: SWITCH USER (BERALIH PENGGUNA)              -->
    <!-- ==================================================== -->
    <div
      v-if="showSwitchModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm"
    >
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-800 shadow-2xl p-5 text-slate-100 space-y-4">
        <div class="flex items-center justify-between pb-2 border-b border-slate-800">
          <div class="flex items-center gap-2">
            <ArrowRightLeft class="w-4 h-4 text-sky-400" />
            <h3 class="text-xs font-bold text-slate-100">Beralih Akun (Switch User)</h3>
          </div>
          <button @click="showSwitchModal = false" class="text-slate-400 hover:text-white p-1">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div v-if="switchError" class="p-2 rounded-xl bg-rose-500/20 border border-rose-500/40 text-rose-300 text-xs">
          {{ switchError }}
        </div>

        <form @submit.prevent="submitSwitchUser" class="space-y-3 text-xs">
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">Username</label>
            <input
              v-model="switchTargetUsername"
              type="text"
              required
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-[11px] font-semibold text-slate-300 mb-1">Password</label>
            <input
              v-model="switchPassword"
              type="password"
              required
              placeholder="Masukkan password akun"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-slate-100 focus:outline-none focus:border-sky-500"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              @click="showSwitchModal = false"
              class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-3.5 py-1.5 rounded-xl bg-sky-600 hover:bg-sky-500 text-white font-bold text-xs shadow-lg shadow-sky-600/30"
            >
              Masuk Akun
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
