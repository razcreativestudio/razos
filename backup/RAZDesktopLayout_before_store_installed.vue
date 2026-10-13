<script setup>
/**
 * RAZOS Desktop Layout Component
 * Updated: 2026-10-08
 * Deskripsi: Layout Desktop mode untuk breakpoint > 768px, memadukan desktop icon shortcuts,
 * window manager frame, dan taskbar bawah dengan wallpaper WebP teroptimasi.
 */
import { useWindowStore } from '../stores/RAZWindowStore'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useContextMenuStore } from '../stores/RAZContextMenuStore'
import RAZWindowFrame from '../components/RAZWindowFrame.vue'
import RAZTaskbar from '../components/RAZTaskbar.vue'
import {
  RefreshCw, Folder, Terminal, Cpu, Settings, EyeOff, Layers, Play, Info
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

// Context Menu Klik Kanan pada Ikon Shortcut Desktop
function handleDesktopIconContextMenu(e, app) {
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

    <!-- Desktop Icons Grid -->
    <div class="relative z-10 p-6 grid grid-flow-col grid-rows-6 gap-6 w-max">
      <div
        v-for="app in windowStore.apps"
        :key="app.id"
        @dblclick="windowStore.openApp(app.id)"
        @contextmenu.prevent.stop="handleDesktopIconContextMenu($event, app)"
        class="w-20 flex flex-col items-center gap-1.5 p-2 rounded-xl hover:bg-slate-800/50 cursor-pointer group transition text-center"
      >
        <div class="w-12 h-12 rounded-2xl flex items-center justify-center p-1 group-hover:scale-105 transition-transform duration-150 drop-shadow-md">
          <img :src="app.icon" :alt="app.title" class="w-full h-full object-contain" />
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
            <RAZFileViewerView v-if="win.id.startsWith('viewer-')" :file="win.fileData" />
            <component v-else :is="viewMap[win.id]" />
          </RAZWindowFrame>
        </template>
      </div>
    </div>

    <!-- Bottom Desktop Taskbar -->
    <RAZTaskbar />
  </div>
</template>
