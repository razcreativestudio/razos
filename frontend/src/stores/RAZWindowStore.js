/**
 * RAZOS Window Manager Store (Pinia)
 * Updated: 2026-10-13 (Standalone Universal File & Content Viewer Integration)
 * Deskripsi: State manager untuk window multi-tasking desktop OS (buka, tutup, minimize, maximize,
 * fokus Z-Index, dan dynamic standalone file viewers untuk gambar, video, musik, pdf, teks/kode).
 */

import { defineStore, acceptHMRUpdate } from 'pinia'
import { ref, computed } from 'vue'

export const useWindowStore = defineStore('window', () => {
  // Aplikasi bawaan sistem operasi
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

  // Kunci penyimpanan preferensi ukuran & posisi jendela di localStorage
  const WINDOW_BOUNDS_STORAGE_KEY = 'razos_window_bounds'

  // Ambil data ukuran, posisi, dan status maximized tersimpan dari localStorage dengan sanitasi angka
  function getSavedBounds(id) {
    try {
      const raw = localStorage.getItem(WINDOW_BOUNDS_STORAGE_KEY)
      if (raw) {
        const parsed = JSON.parse(raw)
        const item = parsed[id]
        if (item && typeof item === 'object') {
          const w = Number(item.width)
          const h = Number(item.height)
          const x = Number(item.x)
          const y = Number(item.y)
          return {
            width: Number.isFinite(w) && w >= 200 ? w : null,
            height: Number.isFinite(h) && h >= 150 ? h : null,
            x: Number.isFinite(x) ? x : null,
            y: Number.isFinite(y) ? y : null,
            isMaximized: item.isMaximized === true
          }
        }
      }
    } catch (e) {}
    return null
  }

  // Simpan data ukuran, koordinat posisi, dan status maximized jendela ke localStorage dengan sanitasi angka
  function saveWindowBounds(id, bounds) {
    if (!id || !bounds) return
    try {
      const raw = localStorage.getItem(WINDOW_BOUNDS_STORAGE_KEY)
      const parsed = raw ? JSON.parse(raw) : {}
      const existing = parsed[id] || {}

      const w = Number(bounds.width)
      const h = Number(bounds.height)
      const x = Number(bounds.x)
      const y = Number(bounds.y)
      const isMax = bounds.isMaximized !== undefined ? !!bounds.isMaximized : (existing.isMaximized === true)

      // Pertahankan ukuran dan koordinat unmaximized jika data saat ini valid
      const validW = Number.isFinite(w) && w >= 200 ? Math.round(w) : (existing.width || 700)
      const validH = Number.isFinite(h) && h >= 150 ? Math.round(h) : (existing.height || 480)
      const validX = Number.isFinite(x) ? Math.round(x) : (existing.x ?? 60)
      const validY = Number.isFinite(y) ? Math.round(y) : (existing.y ?? 40)

      parsed[id] = {
        width: validW,
        height: validH,
        x: validX,
        y: validY,
        isMaximized: isMax
      }
      localStorage.setItem(WINDOW_BOUNDS_STORAGE_KEY, JSON.stringify(parsed))
    } catch (e) {}
  }

  // Buka window aplikasi dengan memulihkan ukuran, posisi, dan status maximize terakhir
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
    const offset = openWindows.value.length * 28
    const saved = getSavedBounds(appId)

    const winWidth = saved?.width || appDef.width
    const winHeight = saved?.height || appDef.height
    const maxX = Math.max(20, (window.innerWidth || 1200) - winWidth - 20)
    const maxY = Math.max(20, (window.innerHeight || 800) - winHeight - 60)

    const winX = (saved && saved.x !== null && saved.x !== undefined) ? Math.min(Math.max(20, saved.x), maxX) : (60 + (offset % 200))
    const winY = (saved && saved.y !== null && saved.y !== undefined) ? Math.min(Math.max(20, saved.y), maxY) : (40 + (offset % 150))
    const isMaximized = saved?.isMaximized === true

    const win = {
      ...appDef,
      width: winWidth,
      height: winHeight,
      x: winX,
      y: winY,
      zIndex: zIndexCounter.value,
      isMinimized: false,
      isMaximized: isMaximized
    }

    openWindows.value.push(win)
    activeWindowId.value = appId
  }

  // Tutup window dan simpan ukuran, posisi, serta status maximize terakhirnya
  function closeWindow(appId) {
    const win = openWindows.value.find(w => w.id === appId)
    if (win) {
      saveWindowBounds(appId, {
        width: win.width,
        height: win.height,
        x: win.x,
        y: win.y,
        isMaximized: win.isMaximized === true
      })
    }
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

  // Toggle maximize dan simpan status terbaru
  function toggleMaximize(appId) {
    const win = openWindows.value.find(w => w.id === appId)
    if (win) {
      win.isMaximized = !win.isMaximized
      saveWindowBounds(appId, {
        width: win.width,
        height: win.height,
        x: win.x,
        y: win.y,
        isMaximized: win.isMaximized === true
      })
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

  // Buka window universal file viewer (gambar, video, audio, pdf, code/text)
  function openFileViewer(fileItem) {
    if (!fileItem || !fileItem.path) return
    const winId = 'viewer-' + btoa(unescape(encodeURIComponent(fileItem.path))).replace(/=/g, '')
    
    // Cek apakah window untuk file ini sudah terbuka
    const existing = openWindows.value.find(w => w.id === winId)
    if (existing) {
      existing.isMinimized = false
      bringToFront(winId)
      return
    }

    const ext = (fileItem.extension || '').toLowerCase()
    let icon = '/assets/icons/RAZIconCodeEditor.svg'
    let width = 740
    let height = 520
    let minWidth = 400
    let minHeight = 300

    if (['.webp', '.png', '.jpg', '.jpeg', '.svg', '.gif', '.ico', '.bmp'].includes(ext)) {
      icon = '/assets/icons/RAZIconImageViewer.svg'
      width = 720
      height = 540
      minWidth = 360
      minHeight = 280
    } else if (['.mp4', '.webm', '.mkv', '.mov', '.avi'].includes(ext)) {
      icon = '/assets/icons/RAZIconVideoPlayer.svg'
      width = 760
      height = 500
      minWidth = 420
      minHeight = 300
    } else if (['.mp3', '.wav', '.ogg', '.aac', '.flac', '.m4a'].includes(ext)) {
      icon = '/assets/icons/RAZIconMusicPlayer.svg'
      width = 480
      height = 520
      minWidth = 380
      minHeight = 440
    } else if (ext === '.pdf') {
      icon = '/assets/icons/RAZIconPdfViewer.svg'
      width = 800
      height = 580
      minWidth = 450
      minHeight = 350
    } else if (['.md', '.markdown', '.mdown', '.mkd'].includes(ext)) {
      icon = '/assets/icons/RAZIconCodeEditor.svg'
      width = 840
      height = 620
      minWidth = 460
      minHeight = 360
    }

    zIndexCounter.value++
    const offset = (openWindows.value.length * 32) % 200
    const saved = getSavedBounds(winId)

    const winWidth = saved?.width || width
    const winHeight = saved?.height || height
    const isMaximized = saved?.isMaximized === true

    const win = {
      id: winId,
      title: fileItem.name,
      icon: icon,
      fileData: fileItem,
      width: winWidth,
      height: winHeight,
      minWidth: minWidth,
      minHeight: minHeight,
      x: (saved && saved.x !== null && saved.x !== undefined) ? saved.x : (120 + offset),
      y: (saved && saved.y !== null && saved.y !== undefined) ? saved.y : (60 + offset),
      zIndex: zIndexCounter.value,
      isMinimized: false,
      isMaximized: isMaximized
    }

    openWindows.value.push(win)
    activeWindowId.value = winId
  }

  // Generic helper untuk membuka jendela kustom / child window independen
  function openCustomWindow(config) {
    if (!config || !config.id) return
    const existing = openWindows.value.find(w => w.id === config.id)
    if (existing) {
      existing.isMinimized = false
      if (config.data) {
        existing.customData = config.data
      }
      bringToFront(config.id)
      return
    }

    zIndexCounter.value++
    const offset = (openWindows.value.length * 32) % 200
    const saved = getSavedBounds(config.id)

    const winWidth = saved?.width || config.width || 600
    const winHeight = saved?.height || config.height || 480
    const maxX = Math.max(20, (window.innerWidth || 1200) - winWidth - 20)
    const maxY = Math.max(20, (window.innerHeight || 800) - winHeight - 60)

    const winX = (saved && saved.x !== null && saved.x !== undefined) ? Math.min(Math.max(20, saved.x), maxX) : Math.max(40, 100 + offset)
    const winY = (saved && saved.y !== null && saved.y !== undefined) ? Math.min(Math.max(20, saved.y), maxY) : Math.max(30, 60 + offset)
    const isMaximized = saved?.isMaximized === true

    const win = {
      id: config.id,
      title: config.title || 'Jendela Dialog',
      icon: config.icon || '/assets/icons/RAZBrandCapsule.png',
      width: winWidth,
      height: winHeight,
      minWidth: config.minWidth || 380,
      minHeight: config.minHeight || 300,
      x: winX,
      y: winY,
      zIndex: zIndexCounter.value,
      isMinimized: false,
      isMaximized: isMaximized,
      customData: config.data || null,
      customType: config.type || config.id
    }

    openWindows.value.push(win)
    activeWindowId.value = config.id
  }

  // Standalone Window: Wizard Pembuatan Website Baru
  function openWebsiteCreator() {
    openCustomWindow({
      id: 'website-create',
      title: 'Wizard Pembuatan Website Terisolasi',
      icon: '/assets/icons/RAZIconWebsites.svg',
      width: 700,
      height: 580,
      minWidth: 460,
      minHeight: 400,
      type: 'website-create'
    })
  }

  // Standalone Window: Edit Konfigurasi Website
  function openWebsiteEditor(site) {
    if (!site || !site.id) return
    openCustomWindow({
      id: `website-edit-${site.id}`,
      title: `Edit: ${site.name}`,
      icon: site.icon || '/assets/icons/RAZIconWebsites.svg',
      width: 720,
      height: 600,
      minWidth: 460,
      minHeight: 420,
      type: 'website-edit',
      data: site
    })
  }

  // Standalone Window: Pengaturan GitHub & Git Pull
  function openWebsiteGit(site) {
    if (!site || !site.id) return
    openCustomWindow({
      id: `website-git-${site.id}`,
      title: `Git Sync: ${site.name}`,
      icon: site.icon || '/assets/icons/RAZIconWebsites.svg',
      width: 620,
      height: 540,
      minWidth: 440,
      minHeight: 380,
      type: 'website-git',
      data: site
    })
  }

  // Standalone Window: Rincian & Kredensial Website
  function openWebsiteDetail(site) {
    if (!site || !site.id) return
    openCustomWindow({
      id: `website-detail-${site.id}`,
      title: `Rincian: ${site.name}`,
      icon: site.icon || '/assets/icons/RAZIconWebsites.svg',
      width: 560,
      height: 500,
      minWidth: 400,
      minHeight: 340,
      type: 'website-detail',
      data: site
    })
  }

  // Standalone Window: Manual Install Docker Container Cepat
  function openAppStoreManual() {
    openCustomWindow({
      id: 'appstore-manual',
      title: 'Instal Cepat Docker Container',
      icon: '/assets/icons/RAZIconAppStore.svg',
      width: 540,
      height: 520,
      minWidth: 400,
      minHeight: 360,
      type: 'appstore-manual'
    })
  }

  // Standalone Window: Form Konfigurasi Pasang Aplikasi Katalog
  function openAppStoreInstaller(app) {
    if (!app || !app.slug) return
    openCustomWindow({
      id: `appstore-install-${app.slug}`,
      title: `Pasang ${app.name}`,
      icon: app.icon || '/assets/icons/RAZIconAppStore.svg',
      width: 500,
      height: 420,
      minWidth: 380,
      minHeight: 300,
      type: 'appstore-install',
      data: app
    })
  }

  // Standalone Window: Form Buat Database Baru
  function openDatabaseCreator() {
    openCustomWindow({
      id: 'database-create',
      title: 'Buat Basis Data Baru',
      icon: '/assets/icons/RAZIconDatabase.svg',
      width: 500,
      height: 460,
      minWidth: 380,
      minHeight: 320,
      type: 'database-create'
    })
  }

  // Standalone Window: Informasi & Properti Berkas
  function openFileProperties(fileItem) {
    if (!fileItem || !fileItem.path) return
    const id = 'file-prop-' + btoa(unescape(encodeURIComponent(fileItem.path))).replace(/=/g, '')
    openCustomWindow({
      id: id,
      title: `Properti: ${fileItem.name}`,
      icon: '/assets/icons/RAZIconFileManager.svg',
      width: 480,
      height: 460,
      minWidth: 360,
      minHeight: 320,
      type: 'file-properties',
      data: fileItem
    })
  }

  // Toggle window dari taskbar (minimize jika aktif di depan, atau restore/bringToFront jika minimize/di belakang)
  function toggleWindowFromTaskbar(winId) {
    const win = openWindows.value.find(w => w.id === winId)
    if (!win) return

    if (win.isMinimized) {
      // Jika sedang diminimalkan, munculkan kembali dan fokuskan
      win.isMinimized = false
      bringToFront(winId)
    } else {
      // Jika jendela sedang aktif di paling depan, minimalkan
      if (activeWindowId.value === winId) {
        minimizeWindow(winId)
      } else {
        // Jika jendela terbuka tapi tertutup jendela lain, bawa ke paling depan
        bringToFront(winId)
      }
    }
  }

  // Minimalkan semua jendela yang terbuka (Show Desktop)
  function minimizeAllWindows() {
    openWindows.value.forEach(w => {
      w.isMinimized = true
    })
    activeWindowId.value = null
  }

  // Pulihkan semua jendela yang diminimalkan
  function restoreAllWindows() {
    openWindows.value.forEach(w => {
      w.isMinimized = false
    })
    if (openWindows.value.length > 0) {
      bringToFront(openWindows.value[openWindows.value.length - 1].id)
    }
  }

  // State Aplikasi yang Terinstall dari App Store (Dinamis dari Container/Catalog/YAML)
  const installedStoreApps = ref([])

  async function fetchInstalledApps() {
    try {
      const res = await fetch('/api/v1/store/installed')
      if (res.ok) {
        installedStoreApps.value = await res.json()
      }
    } catch (e) {}
  }

  // State Website yang Dikelola oleh Website Manager
  const managedWebsites = ref([])

  async function fetchWebsitesList() {
    try {
      const res = await fetch('/api/v1/websites')
      if (res.ok) {
        managedWebsites.value = await res.json()
      }
    } catch (e) {}
  }

  // Luncurkan Web UI Aplikasi Store
  function launchStoreApp(app) {
    if (!app) return
    const port = app.port || app.storeData?.port || 8080
    const url = `http://${window.location.hostname}:${port}`
    window.open(url, '_blank')
  }

  // Luncurkan Website yang Dikelola Website Manager
  function launchWebsiteApp(site) {
    if (!site) return
    const port = site.port || site.websiteData?.port || 80
    const url = `http://${window.location.hostname}:${port}`
    window.open(url, '_blank')
  }

  // Uninstall Aplikasi Store dengan pembersihan data
  async function uninstallStoreApp(slug, purgeData = true) {
    try {
      const res = await fetch(`/api/v1/store/apps/${slug}?purge_data=${purgeData}`, { method: 'DELETE' })
      if (res.ok) {
        await fetchInstalledApps()
        return true
      }
    } catch (e) {}
    return false
  }

  // Seluruh Aplikasi Desktop (Aplikasi Bawaan Sistem + Aplikasi Terpasang dari Store + Website)
  const allDesktopApps = computed(() => {
    const list = [...apps.value]

    // Tambahkan Container Terpasang dari App Store
    installedStoreApps.value.forEach(inst => {
      list.push({
        id: `store-${inst.slug || inst.id}`,
        title: inst.name,
        icon: inst.icon || 'https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/docker.png',
        isStoreApp: true,
        storeData: inst,
        port: inst.port,
        slug: inst.slug || inst.id
      })
    })

    // Tambahkan Website yang Dibuat di Website Manager
    managedWebsites.value.forEach(site => {
      list.push({
        id: `website-${site.id}`,
        title: site.name,
        icon: site.icon || '/assets/icons/RAZIconWebsites.svg',
        isWebsite: true,
        websiteData: site,
        port: site.port,
        domain: site.domain,
        root_folder: site.root_folder,
        status: site.status
      })
    })

    return list
  })

  // Inisialisasi awal daftar aplikasi dan website
  fetchInstalledApps()
  fetchWebsitesList()

  return {
    apps,
    installedStoreApps,
    managedWebsites,
    allDesktopApps,
    fetchInstalledApps,
    fetchWebsitesList,
    launchStoreApp,
    launchWebsiteApp,
    uninstallStoreApp,
    openWindows,
    activeWindowId,
    openApp,
    openCustomWindow,
    openFileViewer,
    openWebsiteCreator,
    openWebsiteEditor,
    openWebsiteGit,
    openWebsiteDetail,
    openAppStoreManual,
    openAppStoreInstaller,
    openDatabaseCreator,
    openFileProperties,
    closeWindow,
    saveWindowBounds,
    minimizeWindow,
    toggleMaximize,
    bringToFront,
    toggleWindowFromTaskbar,
    minimizeAllWindows,
    restoreAllWindows
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useWindowStore, import.meta.hot))
}
