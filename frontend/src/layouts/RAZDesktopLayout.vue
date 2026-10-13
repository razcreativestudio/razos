<script setup>
/**
 * RAZOS Desktop Layout Component
 * Updated: 2026-10-13 (v2.6.0 - Detached Standalone Windows Architecture)
 * Deskripsi: Layout Desktop mode untuk breakpoint > 768px, memadukan desktop icon shortcuts,
 * window manager frame multi-jendela, dan taskbar bawah dengan sistem jendela mandiri tanpa background blur.
 */
import { ref } from 'vue'
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import RAZWindowFrame from '../components/RAZWindowFrame.vue'
import RAZTaskbar from '../components/RAZTaskbar.vue'
import {
  RefreshCw, Folder, Terminal, Cpu, Settings, EyeOff, Layers, Play, Info,
  Trash2, ExternalLink, AlertTriangle, X, Globe
} from 'lucide-vue-next'

// Import View Aplikasi
import RAZTaskManagerView from '../views/RAZTaskManagerView.vue'
import RAZTerminalView from '../views/RAZTerminalView.vue'
import RAZAppStoreView from '../views/RAZAppStoreView.vue'
import RAZWebsitesView from '../views/RAZWebsitesView.vue'
import RAZDatabasesView from '../views/RAZDatabasesView.vue'
import RAZSettingsView from '../views/RAZSettingsView.vue'
import RAZFileManagerView from '../views/RAZFileManagerView.vue'
import RAZFileViewerView from '../views/RAZFileViewerView.vue'

// Import View Jendela Khusus / Standalone Child Windows (v2.6.0)
import RAZWebsiteWizardView from '../views/RAZWebsiteWizardView.vue'
import RAZWebsiteEditView from '../views/RAZWebsiteEditView.vue'
import RAZWebsiteGitView from '../views/RAZWebsiteGitView.vue'
import RAZWebsiteDetailView from '../views/RAZWebsiteDetailView.vue'
import RAZFilePropertiesView from '../views/RAZFilePropertiesView.vue'
import RAZDatabaseCreatorView from '../views/RAZDatabaseCreatorView.vue'
import RAZAppStoreManualView from '../views/RAZAppStoreManualView.vue'
import RAZAppStoreInstallView from '../views/RAZAppStoreInstallView.vue'

const windowStore = useWindowStore()
const systemStore = useSystemStore()
const contextMenuStore = useContextMenuStore()

// Mapping ID aplikasi ke komponen View
const viewMap = {
  'file-manager': RAZFileManagerView,
  'task-manager': RAZTaskManagerView,
  'terminal': RAZTerminalView,
  'app-store': RAZAppStoreView,
  'websites': RAZWebsitesView,
  'databases': RAZDatabasesView,
  'settings': RAZSettingsView
}

// Context Menu Klik Kanan pada Area Desktop Kosong
function handleDesktopContextMenu(e) {
  contextMenuStore.show(e, [
    {
      label: 'Muat Ulang Desktop',
      icon: RefreshCw,
      action: () => {
        systemStore.fetchMetrics()
        systemStore.showNotification('Desktop', 'Tampilan desktop telah disegarkan', 'success')
      }
    },
    { separator: true },
    {
      label: 'Files & Storage',
      icon: Folder,
      action: () => windowStore.openApp('file-manager')
    },
    {
      label: 'Web Terminal',
      icon: Terminal,
      action: () => windowStore.openApp('terminal')
    },
    {
      label: 'Task Manager',
      icon: Cpu,
      action: () => windowStore.openApp('task-manager')
    },
    { separator: true },
    {
      label: 'Tampilkan Desktop',
      icon: EyeOff,
      action: () => {
        if (typeof windowStore.minimizeAllWindows === 'function') {
          windowStore.minimizeAllWindows()
        } else {
          windowStore.openWindows.forEach(w => { w.isMinimized = true })
          windowStore.activeWindowId = null
        }
      }
    },
    {
      label: 'Pulihkan Semua Jendela',
      icon: Layers,
      action: () => {
        if (typeof windowStore.restoreAllWindows === 'function') {
          windowStore.restoreAllWindows()
        } else {
          windowStore.openWindows.forEach(w => { w.isMinimized = false })
          if (windowStore.openWindows.length > 0) {
            windowStore.bringToFront(windowStore.openWindows[windowStore.openWindows.length - 1].id)
          }
        }
      }
    },
    { separator: true },
    {
      label: 'Pengaturan Sistem',
      icon: Settings,
      action: () => windowStore.openApp('settings')
    }
  ], 'Desktop')
}

// State Modal Konfirmasi Uninstall dari Desktop
const showUninstallModal = ref(false)
const appToUninstall = ref(null)
const purgeDataChecked = ref(true)

function openDesktopUninstall(app) {
  appToUninstall.value = app
  purgeDataChecked.value = true
  showUninstallModal.value = true
}

async function confirmDesktopUninstall() {
  if (!appToUninstall.value) return
  const slug = appToUninstall.value.slug
  const success = await windowStore.uninstallStoreApp(slug, purgeDataChecked.value)
  if (success) {
    systemStore.showNotification('Sukses', `Aplikasi ${appToUninstall.value.title} berhasil di-uninstall`, 'success')
    showUninstallModal.value = false
    appToUninstall.value = null
  } else {
    systemStore.showNotification('Gagal', 'Terjadi kesalahan saat menghapus aplikasi', 'error')
  }
}

// Context Menu Klik Kanan pada Ikon Shortcut Desktop
function handleDesktopIconContextMenu(e, app) {
  if (app.isWebsite) {
    // Context Menu untuk Website yang Dibuat di Website Manager
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.launchWebsiteApp(app)
      },
      {
        label: 'Buka di Tab Baru Peramban',
        icon: ExternalLink,
        action: () => windowStore.launchWebsiteApp(app)
      },
      { separator: true },
      {
        label: 'Buka Root Folder',
        icon: Folder,
        action: () => {
          windowStore.openApp('file-manager')
          systemStore.showNotification('Root Folder', `Lokasi: ${app.root_folder}`, 'info')
        }
      },
      {
        label: 'Kelola di Website Manager',
        icon: Globe,
        action: () => windowStore.openApp('websites')
      },
      { separator: true },
      {
        label: 'Info Website',
        icon: Info,
        action: () => {
          systemStore.showNotification('Website Info', `${app.title} (Port: ${app.port} | Domain: ${app.domain})`, 'info')
        }
      }
    ], app.title)
  } else if (app.isStoreApp) {
    // Context Menu untuk Aplikasi Terpasang dari Store / Docker Compose
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.launchStoreApp(app)
      },
      {
        label: 'Buka di Tab Baru Peramban',
        icon: ExternalLink,
        action: () => windowStore.launchStoreApp(app)
      },
      { separator: true },
      {
        label: 'Info Container',
        icon: Info,
        action: () => {
          systemStore.showNotification('Info Container', `${app.title} (Port Host: ${app.port} | Slug: ${app.slug})`, 'info')
        }
      },
      { separator: true },
      {
        label: 'Uninstall Aplikasi',
        icon: Trash2,
        danger: true,
        action: () => openDesktopUninstall(app)
      }
    ], app.title)
  } else {
    // Context Menu untuk Aplikasi Sistem Core
    contextMenuStore.show(e, [
      {
        label: `Buka ${app.title}`,
        icon: Play,
        action: () => windowStore.openApp(app.id)
      },
      { separator: true },
      {
        label: 'Info Aplikasi',
        icon: Info,
        action: () => {
          systemStore.showNotification('Aplikasi', `${app.title} (Resolusi Default: ${app.width}x${app.height})`, 'info')
        }
      }
    ], app.title)
  }
}
</script>

<template>
  <div
    class="relative w-screen h-screen overflow-hidden bg-cover bg-center select-none"
    style="background-image: url('/assets/images/RAZDesktopBg.webp')"
    @contextmenu.prevent="handleDesktopContextMenu"
  >
    <!-- Subtle Brand Watermark di pojok kanan atas -->
    <div class="absolute top-4 right-6 pointer-events-none opacity-40 z-0">
      <img src="/assets/brand/RAZBrandTransparent.webp" alt="RAZ STUDIO" class="h-8 object-contain" />
    </div>

    <!-- Desktop Icons Grid (Aplikasi Sistem + Aplikasi Terpasang dari Store + Website) -->
    <div class="relative z-10 p-6 grid grid-flow-col grid-rows-6 gap-6 w-max">
      <div
        v-for="app in windowStore.allDesktopApps"
        :key="app.id"
        @dblclick="app.isWebsite ? windowStore.launchWebsiteApp(app) : (app.isStoreApp ? windowStore.launchStoreApp(app) : windowStore.openApp(app.id))"
        @contextmenu.prevent.stop="handleDesktopIconContextMenu($event, app)"
        class="w-20 flex flex-col items-center gap-1.5 p-2 rounded-xl hover:bg-slate-800/50 cursor-pointer group transition text-center relative"
      >
        <div class="w-12 h-12 rounded-2xl flex items-center justify-center p-1 group-hover:scale-105 transition-transform duration-150 drop-shadow-md relative">
          <img :src="app.icon" :alt="app.title" class="w-full h-full object-contain" />
          <span v-if="app.isStoreApp || app.isWebsite" class="absolute -top-1 -right-1 w-3 h-3 bg-emerald-500 rounded-full border-2 border-slate-900" :title="app.isWebsite ? 'Website Running' : 'Container Running'"></span>
        </div>
        <span class="text-xs text-slate-100 font-medium group-hover:text-sky-300 drop-shadow line-clamp-2">
          {{ app.title }}
        </span>
      </div>
    </div>

    <!-- Window Management Layer -->
    <div class="absolute inset-0 pointer-events-none z-20">
      <div class="relative w-full h-[calc(100vh-48px)]">
        <template v-for="win in windowStore.openWindows" :key="win.id">
          <RAZWindowFrame :win="win" class="pointer-events-auto">
            <!-- Standalone Universal File Viewer -->
            <RAZFileViewerView v-if="win.id.startsWith('viewer-')" :file="win.fileData" />

            <!-- Standalone File Properties Window -->
            <RAZFilePropertiesView v-else-if="win.id.startsWith('file-prop-')" :win="win" />

            <!-- Standalone Website Wizard Window -->
            <RAZWebsiteWizardView v-else-if="win.id === 'website-create'" :win="win" />

            <!-- Standalone Website Edit Window -->
            <RAZWebsiteEditView v-else-if="win.id.startsWith('website-edit-')" :win="win" />

            <!-- Standalone Website Git Sync Window -->
            <RAZWebsiteGitView v-else-if="win.id.startsWith('website-git-')" :win="win" />

            <!-- Standalone Website Detail Window -->
            <RAZWebsiteDetailView v-else-if="win.id.startsWith('website-detail-')" :win="win" />

            <!-- Standalone Database Creator Window -->
            <RAZDatabaseCreatorView v-else-if="win.id === 'database-create'" :win="win" />

            <!-- Standalone App Store Manual Quick Install Window -->
            <RAZAppStoreManualView v-else-if="win.id === 'appstore-manual'" :win="win" />

            <!-- Standalone App Store Catalog Installer Window -->
            <RAZAppStoreInstallView v-else-if="win.id.startsWith('appstore-install-')" :win="win" />

            <!-- Aplikasi Utama Berdasarkan ID -->
            <component v-else :is="viewMap[win.id]" />
          </RAZWindowFrame>
        </template>
      </div>
    </div>

    <!-- Bottom Desktop Taskbar -->
    <RAZTaskbar />

    <!-- Custom Modal Konfirmasi Uninstall Aplikasi (Floating Window Dialog) -->
    <div
      v-if="showUninstallModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/35 p-4 animate-in fade-in"
    >
      <div class="w-full max-w-md rounded-2xl bg-slate-900 border border-rose-500/40 shadow-2xl shadow-black ring-1 ring-rose-500/20 flex flex-col overflow-hidden text-slate-100">
        <!-- Window Titlebar -->
        <div class="flex items-center justify-between bg-slate-950/85 border-b border-slate-800 px-5 py-3 select-none">
          <div class="flex items-center gap-2.5 text-rose-400">
            <div class="w-7 h-7 rounded-lg bg-rose-500/20 text-rose-400 flex items-center justify-center border border-rose-500/30">
              <AlertTriangle class="w-4 h-4" />
            </div>
            <h3 class="font-bold text-sm">Konfirmasi Uninstall Aplikasi</h3>
          </div>
          <button @click="showUninstallModal = false" class="w-7 h-7 rounded-lg bg-slate-800/80 hover:bg-rose-500 text-slate-400 hover:text-white transition flex items-center justify-center border border-slate-700/60 active:scale-95" title="Tutup Jendela">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 flex flex-col gap-3.5">
          <div class="flex items-center gap-3 p-3 rounded-xl bg-slate-950/60 border border-slate-800">
            <img :src="appToUninstall?.icon" :alt="appToUninstall?.title" class="w-10 h-10 object-contain rounded-lg" />
            <div class="min-w-0 flex-1">
              <h4 class="font-bold text-sm text-slate-100 truncate">{{ appToUninstall?.title }}</h4>
              <span class="text-xs text-slate-400 font-mono">Port: :{{ appToUninstall?.port }}</span>
            </div>
          </div>

          <p class="text-xs text-slate-300 leading-relaxed">
            Apakah Anda yakin ingin menghentikan dan menghapus container aplikasi ini dari sistem?
          </p>

          <label class="flex items-center gap-2.5 p-3 rounded-xl bg-slate-950 border border-slate-800 cursor-pointer text-xs text-rose-300">
            <input type="checkbox" v-model="purgeDataChecked" class="accent-rose-500 rounded" />
            <span>Hapus juga seluruh direktori data permanen container di storage</span>
          </label>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950/75 border-t border-slate-800 px-5 py-3">
          <button
            @click="showUninstallModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            Batal
          </button>
          <button
            @click="confirmDesktopUninstall"
            class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-lg active:scale-95"
          >
            <Trash2 class="w-4 h-4" />
            <span>Uninstall Container</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
