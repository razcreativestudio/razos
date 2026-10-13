/**
 * RAZOS Global Context Menu Store (Pinia)
 * Updated: 2026-10-13 (Enterprise Context Menu System)
 * Deskripsi: State manager terpusat untuk menampilkan Custom Context Menu (klik kanan) di seluruh
 * desktop OS (Desktop wallpaper, Desktop Icons, Taskbar, Taskbar App Items, File Manager items, dan Explorer Canvas).
 * Dilengkapi kalkulasi posisi otomatis (boundary-aware) agar menu tidak keluar dari layar/viewport.
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useContextMenuStore = defineStore('contextMenu', () => {
  const isOpen = ref(false)
  const x = ref(0)
  const y = ref(0)
  const items = ref([])
  const title = ref('')

  /**
   * Tampilkan context menu pada koordinat event klik kanan
   * @param {MouseEvent} event Event mouse native
   * @param {Array} menuItems Daftar menu { label, icon, action, danger, separator, shortcut }
   * @param {string} menuTitle Judul header menu opsional
   */
  function show(event, menuItems, menuTitle = '') {
    if (event && event.preventDefault) {
      event.preventDefault()
      event.stopPropagation()
    }

    title.value = menuTitle
    items.value = menuItems || []

    const clickX = event?.clientX ?? 0
    const clickY = event?.clientY ?? 0

    // Perkiraan dimensi menu untuk boundary check
    const menuWidth = 220
    const estimatedHeight = ((menuItems?.length || 4) * 36) + (menuTitle ? 40 : 16)

    let posX = clickX
    let posY = clickY

    // Cegah menu terpotong di tepi kanan
    if (clickX + menuWidth > window.innerWidth) {
      posX = window.innerWidth - menuWidth - 10
    }

    // Cegah menu terpotong di tepi bawah
    if (clickY + estimatedHeight > window.innerHeight) {
      posY = window.innerHeight - estimatedHeight - 10
    }

    x.value = Math.max(8, posX)
    y.value = Math.max(8, posY)
    isOpen.value = true
  }

  /**
   * Tutup context menu
   */
  function hide() {
    isOpen.value = false
    items.value = []
    title.value = ''
  }

  return {
    isOpen,
    x,
    y,
    items,
    title,
    show,
    hide
  }
})
