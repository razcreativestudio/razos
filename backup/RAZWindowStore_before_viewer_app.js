/**
 * RAZOS Window Manager Store (Pinia)
 * Updated: 2026-10-08
 * Deskripsi: State manager untuk window multi-tasking (buka, tutup, minimize, maximize, fokus Z-Index).
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWindowStore = defineStore('window', () => {
  // Daftar aplikasi yang terdaftar
  const apps = ref([
    {
      id: 'file-manager',
      title: 'Files & Storage',
      icon: '/assets/icons/RAZIconFileManager.svg',
      width: 820,
      height: 520,
      minWidth: 480,
      minHeight: 340,
    },
    {
      id: 'task-manager',
      title: 'Task Manager',
      icon: '/assets/icons/RAZIconTaskManager.svg',
      width: 720,
      height: 480,
      minWidth: 400,
      minHeight: 300,
    },
    {
      id: 'terminal',
      title: 'Web Terminal',
      icon: '/assets/icons/RAZIconTerminal.svg',
      width: 680,
      height: 420,
      minWidth: 400,
      minHeight: 280,
    },
    {
      id: 'app-store',
      title: 'App Store',
      icon: '/assets/icons/RAZIconAppStore.svg',
      width: 800,
      height: 520,
      minWidth: 450,
      minHeight: 350,
    },
    {
      id: 'websites',
      title: 'Website Manager',
      icon: '/assets/icons/RAZIconWebsites.svg',
      width: 750,
      height: 480,
      minWidth: 400,
      minHeight: 320,
    },
    {
      id: 'databases',
      title: 'Database Manager',
      icon: '/assets/icons/RAZIconDatabase.svg',
      width: 700,
      height: 450,
      minWidth: 400,
      minHeight: 300,
    },
    {
      id: 'settings',
      title: 'Settings',
      icon: '/assets/icons/RAZIconSettings.svg',
      width: 650,
      height: 450,
      minWidth: 400,
      minHeight: 300,
    }
  ])

  // Window aktif
  const openWindows = ref([])
  const activeWindowId = ref(null)
  const zIndexCounter = ref(10)

  // Buka window aplikasi
  function openApp(appId) {
    const existing = openWindows.value.find(w => w.id === appId)
    if (existing) {
      if (existing.isMinimized) {
        existing.isMinimized = false
      }
      bringToFront(appId)
      return
    }

    const appDef = apps.value.find(a => a.id === appId)
    if (!appDef) return

    zIndexCounter.value++
    // Posisi acak teratur agar tidak menumpuk tepat di satu titik
    const offset = openWindows.value.length * 28
    const win = {
      ...appDef,
      x: 60 + (offset % 200),
      y: 40 + (offset % 150),
      zIndex: zIndexCounter.value,
      isMinimized: false,
      isMaximized: false
    }

    openWindows.value.push(win)
    activeWindowId.value = appId
  }

  // Tutup window
  function closeWindow(appId) {
    openWindows.value = openWindows.value.filter(w => w.id !== appId)
    if (activeWindowId.value === appId) {
      const remaining = openWindows.value.filter(w => !w.isMinimized)
      activeWindowId.value = remaining.length > 0 ? remaining[remaining.length - 1].id : null
    }
  }

  // Minimize window
  function minimizeWindow(appId) {
    const win = openWindows.value.find(w => w.id === appId)
    if (win) {
      win.isMinimized = true
      if (activeWindowId.value === appId) {
        const remaining = openWindows.value.filter(w => !w.isMinimized)
        activeWindowId.value = remaining.length > 0 ? remaining[remaining.length - 1].id : null
      }
    }
  }

  // Toggle maximize
  function toggleMaximize(appId) {
    const win = openWindows.value.find(w => w.id === appId)
    if (win) {
      win.isMaximized = !win.isMaximized
      bringToFront(appId)
    }
  }

  // Bawa window ke posisi paling depan (fokus Z-Index)
  function bringToFront(appId) {
    const win = openWindows.value.find(w => w.id === appId)
    if (win) {
      zIndexCounter.value++
      win.zIndex = zIndexCounter.value
      activeWindowId.value = appId
    }
  }

  return {
    apps,
    openWindows,
    activeWindowId,
    openApp,
    closeWindow,
    minimizeWindow,
    toggleMaximize,
    bringToFront
  }
})
