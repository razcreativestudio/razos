<template>
  <!--
    RAZOS Authentication View (Login & Initial Setup)
    Timestamp: 2026-10-14 13:12:00
    Deskripsi: Tampilan login dan setup akun awal sistem RAZOS dengan estetika Glassmorphism Slate-900,
               fitur toggle lihat password, validasi realtime, dan custom notification.
  -->
  <div class="relative w-full h-screen overflow-hidden select-none bg-slate-950 flex items-center justify-center font-sans">
    <!-- Background Wallpaper dengan Overlay Halus -->
    <div 
      class="absolute inset-0 bg-cover bg-center transition-all duration-700 scale-105"
      :style="{ backgroundImage: `url('/assets/images/RAZDesktopBg.webp')` }"
    >
      <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-md"></div>
    </div>

    <!-- Vignette & Dekorasi Glow Cahaya -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-blue-600/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-indigo-600/20 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Main Card Container -->
    <div class="relative z-10 w-full max-w-md mx-4">
      <div class="bg-slate-900/85 backdrop-blur-2xl border border-slate-700/60 shadow-2xl shadow-black/60 rounded-3xl p-7 sm:p-9 text-slate-100 transition-all duration-300">
        
        <!-- Header: Logo Brand & Identitas Sistem -->
        <div class="flex flex-col items-center text-center mb-7">
          <div class="relative mb-3 group">
            <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-indigo-500 rounded-2xl blur opacity-30 group-hover:opacity-60 transition duration-500"></div>
            <img 
              src="/assets/brand/RAZBrandTransparent.webp" 
              alt="RAZOS Logo" 
              class="relative w-16 h-16 object-contain drop-shadow-md"
            />
          </div>
          
          <h1 class="text-2xl font-bold tracking-tight bg-gradient-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent">
            {{ isSetupMode ? 'Setup Akun Administrator' : 'Selamat Datang di RAZOS' }}
          </h1>
          <p class="text-xs text-slate-400 mt-1 max-w-xs leading-relaxed">
            {{ isSetupMode 
              ? 'Sistem belum memiliki akun. Buat username & password utama untuk mengamankan perangkat Anda.' 
              : 'Silakan masuk dengan kredensial sistem Anda untuk mengakses desktop.' 
            }}
          </p>
        </div>

        <!-- Banner Pesan Error / Informasi -->
        <div 
          v-if="errorMessage" 
          class="mb-5 p-3.5 rounded-xl bg-rose-500/15 border border-rose-500/30 text-rose-300 text-xs flex items-start gap-2.5 animate-fadeIn"
        >
          <AlertCircle class="w-4 h-4 shrink-0 mt-0.5 text-rose-400" />
          <div class="flex-1 leading-snug">{{ errorMessage }}</div>
        </div>

        <!-- FORM SETUP (REGISTRASI AWAL) -->
        <form v-if="isSetupMode" @submit.prevent="handleRegister" class="space-y-4">
          <!-- Input Username -->
          <div class="space-y-1.5">
            <label class="block text-xs font-medium text-slate-300">Username Admin</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <User class="w-4 h-4" />
              </div>
              <input
                v-model.trim="form.username"
                type="text"
                required
                autocomplete="username"
                placeholder="Misal: admin atau raz"
                class="w-full pl-10 pr-4 py-2.5 bg-slate-800/80 hover:bg-slate-800 border border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 rounded-xl text-sm text-slate-100 placeholder-slate-500 outline-none transition"
              />
            </div>
          </div>

          <!-- Input Password -->
          <div class="space-y-1.5">
            <label class="block text-xs font-medium text-slate-300">Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="new-password"
                placeholder="Minimal 4 karakter"
                class="w-full pl-10 pr-11 py-2.5 bg-slate-800/80 hover:bg-slate-800 border border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 rounded-xl text-sm text-slate-100 placeholder-slate-500 outline-none transition"
              />
              <!-- Tombol Lihat / Sembunyikan Password -->
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-200 transition"
                tabindex="-1"
                :title="showPassword ? 'Sembunyikan Password' : 'Lihat Password'"
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Input Konfirmasi Password -->
          <div class="space-y-1.5">
            <label class="block text-xs font-medium text-slate-300">Konfirmasi Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <ShieldCheck class="w-4 h-4" />
              </div>
              <input
                v-model="form.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                autocomplete="new-password"
                placeholder="Ketik ulang password"
                class="w-full pl-10 pr-11 py-2.5 bg-slate-800/80 hover:bg-slate-800 border border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 rounded-xl text-sm text-slate-100 placeholder-slate-500 outline-none transition"
              />
              <!-- Tombol Lihat / Sembunyikan Konfirmasi Password -->
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-200 transition"
                tabindex="-1"
                :title="showConfirmPassword ? 'Sembunyikan Password' : 'Lihat Password'"
              >
                <EyeOff v-if="showConfirmPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Tombol Submit Registrasi -->
          <div class="pt-2">
            <button
              type="submit"
              :disabled="authStore.isLoading"
              class="w-full py-2.5 px-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 active:scale-[0.99] disabled:opacity-50 text-white font-medium text-sm rounded-xl shadow-lg shadow-blue-500/20 transition flex items-center justify-center gap-2"
            >
              <Loader2 v-if="authStore.isLoading" class="w-4 h-4 animate-spin" />
              <UserPlus v-else class="w-4 h-4" />
              <span>{{ authStore.isLoading ? 'Membuat Akun...' : 'Daftar & Masuk ke RAZOS' }}</span>
            </button>
          </div>
        </form>

        <!-- FORM LOGIN (JIKA SUDAH ADA AKUN) -->
        <form v-else @submit.prevent="handleLogin" class="space-y-4">
          <!-- Input Username -->
          <div class="space-y-1.5">
            <label class="block text-xs font-medium text-slate-300">Username</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <User class="w-4 h-4" />
              </div>
              <input
                v-model.trim="form.username"
                type="text"
                required
                autocomplete="username"
                placeholder="Masukkan username"
                class="w-full pl-10 pr-4 py-2.5 bg-slate-800/80 hover:bg-slate-800 border border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 rounded-xl text-sm text-slate-100 placeholder-slate-500 outline-none transition"
              />
            </div>
          </div>

          <!-- Input Password -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <label class="block text-xs font-medium text-slate-300">Password</label>
            </div>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="current-password"
                placeholder="Masukkan password"
                class="w-full pl-10 pr-11 py-2.5 bg-slate-800/80 hover:bg-slate-800 border border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 rounded-xl text-sm text-slate-100 placeholder-slate-500 outline-none transition"
              />
              <!-- Tombol Lihat / Sembunyikan Password -->
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-200 transition"
                tabindex="-1"
                :title="showPassword ? 'Sembunyikan Password' : 'Lihat Password'"
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Tombol Submit Login -->
          <div class="pt-2">
            <button
              type="submit"
              :disabled="authStore.isLoading"
              class="w-full py-2.5 px-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 active:scale-[0.99] disabled:opacity-50 text-white font-medium text-sm rounded-xl shadow-lg shadow-blue-500/20 transition flex items-center justify-center gap-2"
            >
              <Loader2 v-if="authStore.isLoading" class="w-4 h-4 animate-spin" />
              <LogIn v-else class="w-4 h-4" />
              <span>{{ authStore.isLoading ? 'Memverifikasi...' : 'Masuk ke Sistem' }}</span>
            </button>
          </div>
        </form>

        <!-- Footer Info -->
        <div class="mt-6 pt-5 border-t border-slate-800/80 text-center text-xs text-slate-500 flex items-center justify-center gap-2">
          <span>RAZOS Enterprise v1.1</span>
          <span>•</span>
          <span>Protected System</span>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * Komponen RAZAuthView
 * Menangani logika alur registrasi awal akun admin dan login harian.
 */
import { ref, reactive, computed } from 'vue'
import { useAuthStore } from '../stores/RAZAuthStore'
import {
  User,
  Lock,
  ShieldCheck,
  Eye,
  EyeOff,
  UserPlus,
  LogIn,
  AlertCircle,
  Loader2
} from 'lucide-vue-next'

const authStore = useAuthStore()

// State Form
const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

// State Toggle Lihat Password
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const errorMessage = ref('')

// Tentukan apakah dalam mode setup (jika sistem belum ada akun sama sekali)
const isSetupMode = computed(() => !authStore.isInitialized)

/**
 * Handle Registrasi Akun Awal
 */
const handleRegister = async () => {
  errorMessage.value = ''

  if (!form.username) {
    errorMessage.value = 'Silakan isi username.'
    return
  }

  if (form.password.length < 4) {
    errorMessage.value = 'Password minimal harus 4 karakter.'
    return
  }

  if (form.password !== form.confirmPassword) {
    errorMessage.value = 'Konfirmasi password tidak cocok dengan password yang dimasukkan!'
    return
  }

  const res = await authStore.register(form.username, form.password)
  if (!res.success) {
    errorMessage.value = res.error || 'Gagal mendaftarkan akun.'
  }
}

/**
 * Handle Login Pengguna
 */
const handleLogin = async () => {
  errorMessage.value = ''

  if (!form.username || !form.password) {
    errorMessage.value = 'Silakan isi username dan password.'
    return
  }

  const res = await authStore.login(form.username, form.password)
  if (!res.success) {
    errorMessage.value = res.error || 'Username atau password salah!'
  }
}
</script>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fadeIn {
  animation: fadeIn 0.2s ease-out forwards;
}
</style>
