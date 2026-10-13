<script setup>
/**
 * RAZOS Markdown Document Reader & Viewer Suite
 * Created: 2026-10-13 (v2.6.1 - Dedicated Markdown Reader Engine)
 * Deskripsi: Komponen pembaca dan editor berkas Markdown (.md, .markdown) dengan:
 * 1. Tiga Mode Tampilan: Pratinjau (Preview), Dua Sisi (Split Live Sync), dan Editor Mentah (Raw).
 * 2. Daftar Isi Otomatis (Table of Contents / Outline): Mendeteksi seluruh heading (#, ##, ###) dan navigasi klik halus.
 * 3. Statistik Dokumen: Jumlah Kata, Jumlah Baris, Karakter, dan Estimasi Waktu Baca.
 * 4. Tipografi Dark Mode Elegan: Blockquote, tabel bergaris rapi, code block, task list, dan link aktif.
 * 5. Aksi Dokumen: Salin Teks (Copy), Simpan Perubahan (Save), Muat Ulang (Reload), dan Unduh Berkas.
 */
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { marked } from 'marked'
import {
  BookOpen, Edit3, Columns, List, Copy, Check, Save, RefreshCw,
  Download, FileText, Clock, AlignLeft, Hash, ExternalLink, Sparkles
} from 'lucide-vue-next'
import { useSystemStore } from '../stores/RAZSystemStore'

const props = defineProps({
  file: {
    type: Object,
    required: true
  }
})

const systemStore = useSystemStore()

// State Dokumen
const content = ref('')
const originalContent = ref('')
const isLoading = ref(false)
const isSaving = ref(false)
const isCopied = ref(false)

// Mode Tampilan: 'preview' | 'split' | 'raw'
const viewMode = ref('preview')

// Drawer / Popover Daftar Isi (Outline)
const showOutline = ref(false)
const previewContainerRef = ref(null)

// Mengambil konten berkas markdown dari backend
async function fetchMarkdown() {
  isLoading.value = true
  try {
    const res = await fetch(`/api/v1/files/content?path=${encodeURIComponent(props.file.path)}`)
    const data = await res.json()
    if (res.ok && !data.error) {
      content.value = data.content || ''
      originalContent.value = data.content || ''
    } else {
      systemStore.showNotification('Gagal', data.error || 'Gagal memuat dokumen markdown', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', `Gagal memuat berkas: ${err.message}`, 'error')
  } finally {
    isLoading.value = false
  }
}

// Menyimpan perubahan dokumen
async function saveMarkdown() {
  if (isSaving.value || !isDirty.value) return
  isSaving.value = true
  try {
    const res = await fetch('/api/v1/files/content', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: props.file.path, content: content.value })
    })
    if (res.ok) {
      originalContent.value = content.value
      systemStore.showNotification('Sukses', `Berkas ${props.file.name} berhasil disimpan`, 'success')
    } else {
      systemStore.showNotification('Gagal', 'Gagal menyimpan perubahan berkas', 'error')
    }
  } catch (err) {
    systemStore.showNotification('Error', 'Terjadi kesalahan sistem saat menyimpan', 'error')
  } finally {
    isSaving.value = false
  }
}

// Indikator apakah dokumen telah dimodifikasi
const isDirty = computed(() => content.value !== originalContent.value)

// Statistik Dokumen
const stats = computed(() => {
  const text = content.value || ''
  const words = text.trim() ? text.trim().split(/\s+/).length : 0
  const chars = text.length
  const lines = text ? text.split('\n').length : 0
  const readingTime = Math.max(1, Math.ceil(words / 200)) // Rata-rata 200 kata/menit
  return { words, chars, lines, readingTime }
})

// Parsing Table of Contents (Outline Heading)
const tableOfContents = computed(() => {
  const lines = (content.value || '').split('\n')
  const headings = []
  let headingIndex = 0

  lines.forEach((line) => {
    const match = line.match(/^(#{1,4})\s+(.+)$/)
    if (match) {
      const level = match[1].length
      const title = match[2].trim().replace(/[*_`]/g, '')
      // Buat anchor ID unik yang aman
      const anchor = `heading-${headingIndex++}-${title.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`
      headings.push({ level, title, anchor })
    }
  })

  return headings
})

// Render HTML Markdown menggunakan marked
const renderedHtml = computed(() => {
  if (!content.value) {
    return '<div class="text-slate-500 italic py-8 text-center">Dokumen kosong. Mulai tulis catatan Anda...</div>'
  }

  try {
    // Kustomisasi renderer marked untuk menyematkan ID anchor pada setiap heading
    let hIdx = 0
    const renderer = new marked.Renderer()
    renderer.heading = ({ text, depth }) => {
      const cleanTitle = text.replace(/<[^>]*>/g, '').replace(/[*_`]/g, '').trim()
      const anchor = `heading-${hIdx++}-${cleanTitle.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`
      return `<h${depth} id="${anchor}" class="scroll-mt-4 group flex items-center gap-2">
        <span>${text}</span>
        <a href="#${anchor}" class="text-slate-600 hover:text-sky-400 opacity-0 group-hover:opacity-100 transition text-xs select-none">#</a>
      </h${depth}>`
    }

    return marked.parse(content.value, { renderer, gfm: true, breaks: true })
  } catch (e) {
    return `<div class="text-rose-400 p-4 font-mono">Gagal me-render markdown: ${e.message}</div>`
  }
})

// Loncat ke heading tertentu saat item TOC diklik
function scrollToHeading(anchorId) {
  nextTick(() => {
    if (!previewContainerRef.value) return
    const el = previewContainerRef.value.querySelector(`#${anchorId}`)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' })
      showOutline.value = false
    }
  })
}

// Salin teks markdown ke clipboard
async function copyMarkdownText() {
  try {
    await navigator.clipboard.writeText(content.value)
    isCopied.value = true
    systemStore.showNotification('Tersalin', 'Seluruh isi berkas Markdown berhasil disalin ke clipboard', 'success')
    setTimeout(() => {
      isCopied.value = false
    }, 2000)
  } catch (err) {
    systemStore.showNotification('Gagal', 'Gagal menyalin teks ke clipboard', 'error')
  }
}

// Unduh berkas fisik
function downloadFile() {
  const blob = new Blob([content.value], { type: 'text/markdown;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = props.file.name || 'document.md'
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(() => {
  fetchMarkdown()
})

watch(() => props.file.path, () => {
  fetchMarkdown()
})
</script>

<template>
  <div class="h-full flex flex-col bg-slate-950 text-slate-100 select-none overflow-hidden text-xs">
    <!-- Top Toolbar Header -->
    <div class="h-11 px-3.5 border-b border-slate-800 bg-slate-900/90 flex items-center justify-between gap-2 shrink-0">
      <!-- Sisi Kiri: Ikon, Nama Berkas & Statistik Cepat -->
      <div class="flex items-center gap-2.5 min-w-0">
        <div class="w-6 h-6 rounded-lg bg-sky-500/20 text-sky-400 flex items-center justify-center border border-sky-500/30 shrink-0">
          <BookOpen class="w-3.5 h-3.5" />
        </div>
        <div class="min-w-0 flex items-center gap-2">
          <span class="font-bold text-slate-100 truncate text-xs sm:text-sm" :title="file.name">
            {{ file.name }}
          </span>
          <span class="text-[10px] px-1.5 py-0.5 rounded bg-sky-500/10 text-sky-300 font-mono font-semibold border border-sky-500/20 uppercase shrink-0">
            MD
          </span>
          <span v-if="isDirty" class="text-[10px] px-1.5 py-0.5 rounded bg-amber-500/20 text-amber-300 font-bold border border-amber-500/30 shrink-0">
            Unsaved
          </span>
        </div>
      </div>

      <!-- Sisi Tengah: Statistik Bacaan (Desktop / Layar Cukup) -->
      <div class="hidden md:flex items-center gap-3 text-[11px] text-slate-400 bg-slate-950/60 px-3 py-1 rounded-xl border border-slate-800">
        <div class="flex items-center gap-1" title="Jumlah Kata">
          <AlignLeft class="w-3 h-3 text-sky-400" />
          <span>{{ stats.words.toLocaleString('id-ID') }} kata</span>
        </div>
        <span class="text-slate-700">•</span>
        <div class="flex items-center gap-1" title="Estimasi Waktu Membaca">
          <Clock class="w-3 h-3 text-amber-400" />
          <span>~{{ stats.readingTime }} mnt baca</span>
        </div>
      </div>

      <!-- Sisi Kanan: Segmented View Switcher & Action Tools -->
      <div class="flex items-center gap-1.5 shrink-0">
        <!-- Segmented Mode Controls -->
        <div class="flex items-center bg-slate-950 p-0.5 rounded-xl border border-slate-800">
          <button
            @click="viewMode = 'preview'"
            class="px-2.5 py-1 rounded-lg font-medium text-[11px] flex items-center gap-1 transition"
            :class="viewMode === 'preview' ? 'bg-sky-500 text-slate-950 font-bold shadow' : 'text-slate-400 hover:text-slate-200'"
            title="Pratinjau Artikel (Preview)"
          >
            <BookOpen class="w-3 h-3" />
            <span class="hidden sm:inline">Pratinjau</span>
          </button>
          <button
            @click="viewMode = 'split'"
            class="px-2.5 py-1 rounded-lg font-medium text-[11px] flex items-center gap-1 transition"
            :class="viewMode === 'split' ? 'bg-sky-500 text-slate-950 font-bold shadow' : 'text-slate-400 hover:text-slate-200'"
            title="Dua Sisi: Editor & Live Preview (Split)"
          >
            <Columns class="w-3 h-3" />
            <span class="hidden sm:inline">Split</span>
          </button>
          <button
            @click="viewMode = 'raw'"
            class="px-2.5 py-1 rounded-lg font-medium text-[11px] flex items-center gap-1 transition"
            :class="viewMode === 'raw' ? 'bg-sky-500 text-slate-950 font-bold shadow' : 'text-slate-400 hover:text-slate-200'"
            title="Editor Kode Mentah (Raw)"
          >
            <Edit3 class="w-3 h-3" />
            <span class="hidden sm:inline">Editor</span>
          </button>
        </div>

        <!-- Tombol Buka/Tutup Outline (Daftar Isi) -->
        <button
          v-if="tableOfContents.length > 0"
          @click="showOutline = !showOutline"
          class="p-1.5 rounded-lg border transition flex items-center gap-1 text-[11px]"
          :class="showOutline
            ? 'bg-sky-500/20 text-sky-300 border-sky-500/40 font-semibold'
            : 'bg-slate-800/80 text-slate-300 border-slate-700/60 hover:text-white hover:bg-slate-700'"
          title="Daftar Isi & Struktur Heading"
        >
          <List class="w-3.5 h-3.5" />
          <span class="hidden lg:inline">Daftar Isi</span>
        </button>

        <!-- Tombol Salin Teks -->
        <button
          @click="copyMarkdownText"
          class="p-1.5 rounded-lg bg-slate-800/80 hover:bg-slate-700 border border-slate-700/60 text-slate-300 hover:text-white transition"
          :title="isCopied ? 'Tersalin!' : 'Salin Seluruh Konten Markdown'"
        >
          <Check v-if="isCopied" class="w-3.5 h-3.5 text-emerald-400 stroke-[2.5]" />
          <Copy v-else class="w-3.5 h-3.5" />
        </button>

        <!-- Tombol Unduh -->
        <button
          @click="downloadFile"
          class="p-1.5 rounded-lg bg-slate-800/80 hover:bg-slate-700 border border-slate-700/60 text-slate-300 hover:text-white transition"
          title="Unduh Berkas .md"
        >
          <Download class="w-3.5 h-3.5" />
        </button>

        <!-- Tombol Simpan (Aktif jika di mode editor/split dan ada perubahan) -->
        <button
          v-if="viewMode !== 'preview' || isDirty"
          @click="saveMarkdown"
          :disabled="isSaving || !isDirty"
          class="px-3 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 disabled:opacity-40 disabled:hover:bg-sky-500 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition shadow active:scale-95 ml-1"
          title="Simpan Perubahan Berkas"
        >
          <Save class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">{{ isSaving ? 'Menyimpan...' : 'Simpan' }}</span>
        </button>
      </div>
    </div>

    <!-- Main Workspace Container -->
    <div class="flex-1 flex overflow-hidden relative">
      <!-- 1. Sisi Kiri / Full: RAW EDITOR -->
      <div
        v-if="viewMode === 'raw' || viewMode === 'split'"
        class="h-full flex flex-col border-r border-slate-800 bg-slate-950/95"
        :class="viewMode === 'split' ? 'w-1/2' : 'w-full'"
      >
        <div class="h-7 px-3 bg-slate-900/60 border-b border-slate-800/80 flex items-center justify-between text-[11px] text-slate-400 shrink-0 font-mono">
          <span>MARKDOWN SOURCE</span>
          <span>{{ stats.lines }} baris • {{ stats.chars }} karakter</span>
        </div>
        <textarea
          v-model="content"
          class="w-full flex-1 p-4 bg-transparent text-slate-200 font-mono text-xs leading-relaxed outline-none resize-none selection:bg-sky-500/30 border-none"
          placeholder="Tulis sintaks markdown di sini..."
          spellcheck="false"
        ></textarea>
      </div>

      <!-- 2. Sisi Kanan / Full: PREVIEW ARTICLE -->
      <div
        v-if="viewMode === 'preview' || viewMode === 'split'"
        ref="previewContainerRef"
        class="h-full flex flex-col overflow-y-auto bg-slate-950/70 p-4 sm:p-6 lg:p-8"
        :class="viewMode === 'split' ? 'w-1/2' : 'w-full'"
      >
        <!-- Artikel Container Berpusat agar Mudah Dibaca -->
        <div class="max-w-3xl w-full mx-auto select-text">
          <div v-html="renderedHtml" class="raz-markdown-reader"></div>
        </div>
      </div>

      <!-- Popover / Sidebar Drawer Table of Contents (Outline) -->
      <div
        v-if="showOutline && tableOfContents.length > 0"
        class="absolute top-2 right-2 bottom-2 w-72 bg-slate-900/95 backdrop-blur-xl border border-slate-800 shadow-2xl rounded-2xl flex flex-col z-30 overflow-hidden animate-in fade-in slide-in-from-right-4"
      >
        <div class="p-3 border-b border-slate-800 flex items-center justify-between bg-slate-950/60">
          <div class="flex items-center gap-2">
            <List class="w-4 h-4 text-sky-400" />
            <span class="font-bold text-xs text-slate-100">Daftar Isi (Outline)</span>
          </div>
          <button
            @click="showOutline = false"
            class="w-5 h-5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-white flex items-center justify-center transition"
          >
            ✕
          </button>
        </div>

        <div class="flex-1 overflow-y-auto p-2 space-y-0.5 text-xs">
          <button
            v-for="item in tableOfContents"
            :key="item.anchor"
            @click="scrollToHeading(item.anchor)"
            class="w-full text-left py-1.5 px-2 rounded-lg hover:bg-slate-800/80 text-slate-300 hover:text-sky-300 transition truncate block"
            :class="{
              'pl-2 font-bold text-slate-100': item.level === 1,
              'pl-4 text-slate-300': item.level === 2,
              'pl-6 text-slate-400 text-[11px]': item.level === 3,
              'pl-8 text-slate-500 text-[11px]': item.level >= 4
            }"
            :title="item.title"
          >
            {{ item.title }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* Tipografi Khusus RAZOS Markdown Reader (Dark Theme) */
.raz-markdown-reader {
  color: #cbd5e1;
  font-size: 0.9375rem;
  line-height: 1.75;
}

.raz-markdown-reader h1 {
  font-size: 1.875rem;
  font-weight: 800;
  color: #f8fafc;
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #334155;
  letter-spacing: -0.025em;
}

.raz-markdown-reader h2 {
  font-size: 1.375rem;
  font-weight: 700;
  color: #f1f5f9;
  margin-top: 1.75rem;
  margin-bottom: 0.75rem;
  padding-bottom: 0.35rem;
  border-bottom: 1px solid rgba(51, 65, 85, 0.4);
}

.raz-markdown-reader h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #e2e8f0;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
}

.raz-markdown-reader h4,
.raz-markdown-reader h5,
.raz-markdown-reader h6 {
  font-size: 0.95rem;
  font-weight: 600;
  color: #cbd5e1;
  margin-top: 1.25rem;
  margin-bottom: 0.5rem;
}

.raz-markdown-reader p {
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
  color: #cbd5e1;
}

.raz-markdown-reader strong {
  color: #f8fafc;
  font-weight: 700;
}

.raz-markdown-reader a {
  color: #38bdf8;
  text-decoration: underline;
  text-underline-offset: 3px;
  transition: color 0.15s ease;
}

.raz-markdown-reader a:hover {
  color: #7dd3fc;
}

.raz-markdown-reader blockquote {
  margin: 1.25rem 0;
  padding: 0.75rem 1.25rem;
  background-color: rgba(15, 23, 42, 0.7);
  border-left: 4px solid #38bdf8;
  border-radius: 0 0.75rem 0.75rem 0;
  color: #94a3b8;
  font-style: italic;
}

.raz-markdown-reader blockquote p {
  margin: 0;
}

.raz-markdown-reader ul {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin: 0.75rem 0;
}

.raz-markdown-reader ol {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin: 0.75rem 0;
}

.raz-markdown-reader li {
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}

.raz-markdown-reader code:not(pre code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.85em;
  background-color: rgba(30, 41, 59, 0.8);
  color: #fbbf24;
  padding: 0.15rem 0.4rem;
  border-radius: 0.375rem;
  border: 1px solid rgba(51, 65, 85, 0.6);
}

.raz-markdown-reader pre {
  margin: 1.25rem 0;
  padding: 1rem 1.25rem;
  background-color: #020617;
  border: 1px solid #1e293b;
  border-radius: 0.875rem;
  overflow-x: auto;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.85rem;
  line-height: 1.6;
  color: #e2e8f0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.raz-markdown-reader pre code {
  background: transparent;
  padding: 0;
  border: none;
  color: inherit;
}

.raz-markdown-reader table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.25rem 0;
  font-size: 0.875rem;
  overflow-x: auto;
  display: block;
}

.raz-markdown-reader th {
  background-color: #0f172a;
  color: #38bdf8;
  font-weight: 700;
  text-align: left;
  padding: 0.625rem 1rem;
  border: 1px solid #334155;
}

.raz-markdown-reader td {
  padding: 0.625rem 1rem;
  border: 1px solid #1e293b;
  color: #cbd5e1;
}

.raz-markdown-reader tr:nth-child(even) {
  background-color: rgba(15, 23, 42, 0.4);
}

.raz-markdown-reader hr {
  border: none;
  border-top: 1px solid #334155;
  margin: 2rem 0;
}

.raz-markdown-reader img {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
  margin: 1rem 0;
  border: 1px solid #334155;
}
</style>
