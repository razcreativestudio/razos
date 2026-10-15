<template>
  <!--
    RAZOS Database Manager (HeidiSQL Client Suite)
    Timestamp: 2026-10-14 13:30:00
    Deskripsi: GUI database management komprehensif layaknya HeidiSQL untuk pengelolaan basis data MySQL/MariaDB.
               Mencakup penjelajah tabel, editor data grid spreadsheet, inspeksi struktur skema kolom,
               konsol eksekusi query SQL bebas (dengan timer eksekusi), serta ekspor dump SQL.
  -->
  <div class="h-full flex flex-col text-slate-100 select-none bg-slate-950 font-sans overflow-hidden">
    
    <!-- 1. TOP TOOLBAR HEIDISQL -->
    <div class="flex flex-wrap items-center justify-between gap-2 px-4 py-2.5 bg-slate-900/90 border-b border-slate-800 shrink-0">
      <!-- Database Selector & Host Info -->
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 px-2.5 py-1 rounded-lg bg-slate-950 border border-slate-800 text-xs">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span class="font-mono text-slate-400">127.0.0.1:3306</span>
          <span class="text-slate-600">|</span>
          <Database class="w-3.5 h-3.5 text-amber-400" />
          
          <!-- Dropdown Pilihan Database -->
          <select 
            v-model="activeDb" 
            @change="handleDbChange"
            class="bg-transparent text-slate-200 font-semibold outline-none cursor-pointer text-xs"
          >
            <option 
              v-for="db in databases" 
              :key="db.name" 
              :value="db.name"
              class="bg-slate-900 text-slate-200"
            >
              {{ db.name }} {{ db.isSystem ? '(System)' : '' }}
            </option>
          </select>
        </div>

        <button 
          @click="fetchAllDatabases" 
          :disabled="isLoadingDbs"
          class="p-1.5 rounded-lg bg-slate-800/80 hover:bg-slate-800 border border-slate-700/60 text-slate-400 hover:text-slate-100 transition"
          title="Muat Ulang Database & Tabel"
        >
          <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoadingDbs }" />
        </button>
      </div>

      <!-- Action Buttons Toolbar -->
      <div class="flex items-center gap-2">
        <!-- Shortcut Buat Tabel Baru -->
        <button
          @click="openCreateTableModal"
          class="px-2.5 py-1.5 rounded-lg bg-sky-500/15 hover:bg-sky-500 text-sky-300 hover:text-slate-950 border border-sky-500/30 font-medium text-xs flex items-center gap-1.5 transition active:scale-95"
          title="Buat Tabel Baru dalam Database Aktif"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tabel Baru</span>
        </button>

        <!-- Shortcut Buat Database Baru -->
        <button
          @click="windowStore.openDatabaseCreator()"
          class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 font-medium text-xs flex items-center gap-1.5 transition active:scale-95"
          title="Deploy Database Baru"
        >
          <Database class="w-3.5 h-3.5 text-amber-400" />
          <span>Database Baru</span>
        </button>

        <!-- Shortcut Web Adminer -->
        <button
          @click="launchAdminer"
          class="px-2.5 py-1.5 rounded-lg bg-amber-500/15 hover:bg-amber-500 text-amber-300 hover:text-slate-950 border border-amber-500/30 font-medium text-xs flex items-center gap-1.5 transition active:scale-95"
          title="Buka Web Adminer GUI"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Adminer</span>
        </button>
      </div>
    </div>

    <!-- 2. MAIN SPLIT BODY (Left Treeview + Right HeidiSQL Workbench) -->
    <div class="flex-1 flex overflow-hidden">
      
      <!-- LEFT PANE: TABLES EXPLORER -->
      <div class="w-64 sm:w-72 border-r border-slate-800 bg-slate-950/70 flex flex-col shrink-0">
        <!-- Filter Search Box -->
        <div class="p-2.5 border-b border-slate-800/80">
          <div class="relative">
            <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-500" />
            <input
              v-model="tableSearchQuery"
              type="text"
              placeholder="Cari nama tabel..."
              class="w-full pl-8 pr-3 py-1.5 rounded-lg bg-slate-900 border border-slate-800 text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500 transition"
            />
          </div>
        </div>

        <!-- Tables List Header -->
        <div class="px-3 py-1.5 bg-slate-900/50 border-b border-slate-800/50 flex items-center justify-between text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
          <span>Daftar Tabel ({{ filteredTables.length }})</span>
          <span class="font-mono text-slate-500 text-[10px]">{{ activeDb }}</span>
        </div>

        <!-- Tables Treeview Scrollable -->
        <div class="flex-1 overflow-y-auto p-1.5 space-y-0.5">
          <div v-if="isLoadingTables" class="p-4 text-center text-xs text-slate-500 flex items-center justify-center gap-2">
            <RotateCw class="w-3.5 h-3.5 animate-spin" />
            <span>Memuat tabel...</span>
          </div>

          <div v-else-if="filteredTables.length === 0" class="p-6 text-center text-xs text-slate-500">
            Tidak ada tabel ditemukan.
          </div>

          <button
            v-else
            v-for="tbl in filteredTables"
            :key="tbl.name"
            @click="selectTable(tbl.name)"
            class="w-full px-2.5 py-1.5 rounded-lg text-left text-xs flex items-center justify-between group transition"
            :class="selectedTable === tbl.name ? 'bg-sky-500/20 text-sky-200 font-semibold border border-sky-500/30' : 'text-slate-300 hover:bg-slate-900 border border-transparent'"
          >
            <div class="flex items-center gap-2 min-w-0">
              <Table class="w-3.5 h-3.5 shrink-0" :class="selectedTable === tbl.name ? 'text-sky-400' : 'text-slate-500 group-hover:text-slate-400'" />
              <span class="truncate">{{ tbl.name }}</span>
            </div>
            <span class="text-[10px] px-1.5 py-0.5 rounded bg-slate-900 text-slate-500 group-hover:text-slate-400 font-mono">
              {{ tbl.rows }}
            </span>
          </button>
        </div>

        <!-- Left Pane Footer: Export Database -->
        <div class="p-2 border-t border-slate-800/80 bg-slate-900/40 flex items-center justify-between text-xs text-slate-400">
          <span class="truncate text-[11px]">Database: {{ activeDb }}</span>
          <button
            v-if="selectedTable"
            @click="exportTableSql(selectedTable)"
            class="text-[11px] text-sky-400 hover:text-sky-300 flex items-center gap-1 font-medium"
            title="Download SQL Dump Tabel"
          >
            <Download class="w-3 h-3" />
            <span>Dump SQL</span>
          </button>
        </div>
      </div>

      <!-- RIGHT PANE: HEIDISQL WORKBENCH -->
      <div class="flex-1 flex flex-col overflow-hidden bg-slate-950">
        
        <!-- WORKBENCH TAB SELECTOR (Data | Struktur | Query | Info) -->
        <div class="flex items-center justify-between border-b border-slate-800 bg-slate-900/60 px-4 pt-1.5 shrink-0">
          <div class="flex items-center gap-1">
            <!-- Tab 1: Data (Isi Tabel) -->
            <button
              @click="activeWorkbenchTab = 'data'"
              class="px-3.5 py-2 text-xs font-semibold rounded-t-lg flex items-center gap-2 border-b-2 transition"
              :class="activeWorkbenchTab === 'data' ? 'border-sky-500 text-sky-400 bg-slate-900' : 'border-transparent text-slate-400 hover:text-slate-200'"
            >
              <Table class="w-3.5 h-3.5" />
              <span>Data Tabel</span>
              <span v-if="selectedTable" class="text-[10px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400 font-mono">{{ tableDataTotal }}</span>
            </button>

            <!-- Tab 2: Struktur Kolom -->
            <button
              @click="activeWorkbenchTab = 'structure'"
              class="px-3.5 py-2 text-xs font-semibold rounded-t-lg flex items-center gap-2 border-b-2 transition"
              :class="activeWorkbenchTab === 'structure' ? 'border-sky-500 text-sky-400 bg-slate-900' : 'border-transparent text-slate-400 hover:text-slate-200'"
            >
              <Columns class="w-3.5 h-3.5" />
              <span>Struktur Skema</span>
            </button>

            <!-- Tab 3: Konsol SQL Query -->
            <button
              @click="activeWorkbenchTab = 'query'"
              class="px-3.5 py-2 text-xs font-semibold rounded-t-lg flex items-center gap-2 border-b-2 transition"
              :class="activeWorkbenchTab === 'query' ? 'border-sky-500 text-sky-400 bg-slate-900' : 'border-transparent text-slate-400 hover:text-slate-200'"
            >
              <Terminal class="w-3.5 h-3.5 text-emerald-400" />
              <span>Konsol SQL (F9)</span>
            </button>

            <!-- Tab 4: Info Tabel -->
            <button
              @click="activeWorkbenchTab = 'info'"
              class="px-3.5 py-2 text-xs font-semibold rounded-t-lg flex items-center gap-2 border-b-2 transition"
              :class="activeWorkbenchTab === 'info' ? 'border-sky-500 text-sky-400 bg-slate-900' : 'border-transparent text-slate-400 hover:text-slate-200'"
            >
              <Info class="w-3.5 h-3.5" />
              <span>Informasi Tabel</span>
            </button>
          </div>

          <!-- Active Table Breadcrumb -->
          <div v-if="selectedTable" class="text-xs text-slate-400 flex items-center gap-1.5 font-mono">
            <span class="text-slate-500">{{ activeDb }}.</span>
            <span class="text-sky-300 font-bold">{{ selectedTable }}</span>
          </div>
        </div>

        <!-- ============================================================== -->
        <!-- TAB CONTENT 1: DATA GRID (ISi TABEL SPREADSHEET HEIDISQL) -->
        <!-- ============================================================== -->
        <div v-if="activeWorkbenchTab === 'data'" class="flex-1 flex flex-col overflow-hidden">
          <div v-if="!selectedTable" class="flex-1 flex flex-col items-center justify-center p-8 text-slate-500 gap-2">
            <Table class="w-10 h-10 text-slate-600" />
            <p class="text-sm font-medium">Pilih tabel di sebelah kiri untuk melihat dan mengelola datanya.</p>
          </div>

          <template v-else>
            <!-- Data Grid Toolbar -->
            <div class="px-4 py-2 border-b border-slate-800/80 bg-slate-900/40 flex flex-wrap items-center justify-between gap-2 text-xs">
              <div class="flex items-center gap-2">
                <!-- Tambah Baris Baru -->
                <button
                  @click="openInsertRowModal"
                  class="px-2.5 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold flex items-center gap-1.5 transition active:scale-95 shadow"
                >
                  <Plus class="w-3.5 h-3.5" />
                  <span>Tambah Baris</span>
                </button>

                <!-- Refresh Data -->
                <button
                  @click="fetchTableData"
                  :disabled="isLoadingData"
                  class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700/60 transition"
                  title="Muat Ulang Data Tabel"
                >
                  <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoadingData }" />
                </button>

                <!-- Search Input In-Table -->
                <div class="relative">
                  <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-500" />
                  <input
                    v-model="dataSearchQuery"
                    @keyup.enter="fetchTableData"
                    type="text"
                    placeholder="Cari dalam isi data..."
                    class="pl-8 pr-3 py-1 rounded-lg bg-slate-950 border border-slate-800 text-xs text-slate-200 placeholder-slate-500 outline-none focus:border-sky-500 transition w-44 sm:w-56"
                  />
                </div>
              </div>

              <!-- Pagination Info & Limits -->
              <div class="flex items-center gap-3 text-slate-400">
                <span>Total: <strong class="text-slate-200">{{ tableDataTotal }}</strong> baris</span>
                <span class="text-slate-600">|</span>
                
                <div class="flex items-center gap-1">
                  <span>Limit:</span>
                  <select 
                    v-model="tableDataLimit" 
                    @change="fetchTableData" 
                    class="bg-slate-900 border border-slate-800 rounded px-1.5 py-0.5 text-slate-200 outline-none"
                  >
                    <option :value="25">25</option>
                    <option :value="50">50</option>
                    <option :value="100">100</option>
                    <option :value="250">250</option>
                  </select>
                </div>

                <!-- Page Nav -->
                <div class="flex items-center gap-1">
                  <button
                    @click="changeDataPage(tableDataPage - 1)"
                    :disabled="tableDataPage <= 1"
                    class="px-2 py-0.5 rounded bg-slate-900 border border-slate-800 disabled:opacity-30 hover:bg-slate-800 text-xs"
                  >
                    &lt;
                  </button>
                  <span class="font-mono text-xs px-1 text-slate-300">{{ tableDataPage }} / {{ tableDataTotalPages }}</span>
                  <button
                    @click="changeDataPage(tableDataPage + 1)"
                    :disabled="tableDataPage >= tableDataTotalPages"
                    class="px-2 py-0.5 rounded bg-slate-900 border border-slate-800 disabled:opacity-30 hover:bg-slate-800 text-xs"
                  >
                    &gt;
                  </button>
                </div>
              </div>
            </div>

            <!-- Spreadsheet Data Table Container -->
            <div class="flex-1 overflow-auto bg-slate-950">
              <div v-if="isLoadingData" class="p-12 text-center text-xs text-slate-500 flex items-center justify-center gap-2">
                <RotateCw class="w-4 h-4 animate-spin text-sky-400" />
                <span>Mengambil data dari MySQL...</span>
              </div>

              <div v-else-if="tableRows.length === 0" class="p-12 text-center text-xs text-slate-500">
                Tabel masih kosong atau tidak ada data yang cocok dengan kriteria pencarian.
              </div>

              <table v-else class="w-full border-collapse text-left text-xs font-mono">
                <!-- Sticky Header Columns -->
                <thead class="sticky top-0 z-10 bg-slate-900 shadow-md">
                  <tr class="border-b border-slate-800 divide-x divide-slate-800">
                    <th class="w-12 px-3 py-2 text-center text-slate-500 font-bold bg-slate-900/95">#</th>
                    <th class="w-16 px-2 py-2 text-center text-slate-400 font-semibold bg-slate-900/95">Aksi</th>
                    <th
                      v-for="col in tableColumns"
                      :key="col"
                      @click="handleSort(col)"
                      class="px-3 py-2 font-semibold text-slate-200 cursor-pointer hover:bg-slate-800/80 transition select-none bg-slate-900/95 whitespace-nowrap"
                    >
                      <div class="flex items-center gap-1.5">
                        <span>{{ col }}</span>
                        <span v-if="sortColumn === col" class="text-sky-400 text-[10px]">
                          {{ sortOrder === 'ASC' ? '▲' : '▼' }}
                        </span>
                      </div>
                    </th>
                  </tr>
                </thead>

                <!-- Body Rows -->
                <tbody class="divide-y divide-slate-800/70">
                  <tr
                    v-for="(row, rIdx) in tableRows"
                    :key="rIdx"
                    class="hover:bg-slate-900/70 divide-x divide-slate-800/50 transition group"
                  >
                    <!-- Row Index Number -->
                    <td class="px-3 py-2 text-center text-slate-600 text-[11px] select-none">
                      {{ (tableDataPage - 1) * tableDataLimit + rIdx + 1 }}
                    </td>

                    <!-- Row Action Buttons (Edit / Delete) -->
                    <td class="px-2 py-1.5 text-center">
                      <div class="flex items-center justify-center gap-1">
                        <button
                          @click="openEditRowModal(row)"
                          class="p-1 rounded text-slate-400 hover:text-sky-300 hover:bg-slate-800 transition"
                          title="Edit Baris"
                        >
                          <Pencil class="w-3 h-3" />
                        </button>
                        <button
                          @click="openDeleteRowModal(row)"
                          class="p-1 rounded text-slate-400 hover:text-rose-400 hover:bg-slate-800 transition"
                          title="Hapus Baris"
                        >
                          <Trash2 class="w-3 h-3" />
                        </button>
                      </div>
                    </td>

                    <!-- Row Cells -->
                    <td
                      v-for="col in tableColumns"
                      :key="col"
                      class="px-3 py-2 max-w-xs truncate text-slate-300"
                      :title="String(row[col])"
                    >
                      <span v-if="row[col] === null" class="text-slate-600 italic font-sans text-[11px]">NULL</span>
                      <span v-else>{{ row[col] }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </template>
        </div>

        <!-- ============================================================== -->
        <!-- TAB CONTENT 2: STRUCTURE (STRUKTUR SKEMA KOLOM) -->
        <!-- ============================================================== -->
        <div v-if="activeWorkbenchTab === 'structure'" class="flex-1 flex flex-col overflow-hidden p-4">
          <div v-if="!selectedTable" class="flex-1 flex flex-col items-center justify-center p-8 text-slate-500 gap-2">
            <Columns class="w-10 h-10 text-slate-600" />
            <p class="text-sm font-medium">Pilih tabel di sebelah kiri untuk melihat struktur skema kolomnya.</p>
          </div>

          <template v-else>
            <div class="flex items-center justify-between mb-3">
              <div>
                <h4 class="font-bold text-slate-200 text-sm">Struktur Kolom Tabel: {{ selectedTable }}</h4>
                <p class="text-xs text-slate-400">Daftar field, tipe data, kunci indeks, dan konfigurasi default.</p>
              </div>

              <button
                @click="fetchTableStructure"
                class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center gap-1.5 border border-slate-700 transition"
              >
                <RotateCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoadingStructure }" />
                <span>Refresh Struktur</span>
              </button>
            </div>

            <!-- Structure Grid -->
            <div class="flex-1 overflow-auto rounded-xl border border-slate-800 bg-slate-900/60">
              <table class="w-full border-collapse text-left text-xs font-mono">
                <thead class="bg-slate-900 border-b border-slate-800 divide-x divide-slate-800">
                  <tr>
                    <th class="px-4 py-2.5 text-slate-300 font-bold">Field / Kolom</th>
                    <th class="px-4 py-2.5 text-slate-300 font-bold">Tipe Data</th>
                    <th class="px-3 py-2.5 text-slate-300 font-bold text-center">Null</th>
                    <th class="px-3 py-2.5 text-slate-300 font-bold text-center">Key</th>
                    <th class="px-4 py-2.5 text-slate-300 font-bold">Default</th>
                    <th class="px-4 py-2.5 text-slate-300 font-bold">Ekstra</th>
                    <th class="px-4 py-2.5 text-slate-300 font-bold">Komentar</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-800/80">
                  <tr v-for="col in tableStructure" :key="col.field" class="hover:bg-slate-800/40 divide-x divide-slate-800/50">
                    <td class="px-4 py-2.5 font-bold text-sky-300 flex items-center gap-1.5">
                      <Key v-if="col.key === 'PRI'" class="w-3.5 h-3.5 text-amber-400 shrink-0" title="Primary Key" />
                      <span>{{ col.field }}</span>
                    </td>
                    <td class="px-4 py-2.5 text-emerald-300">{{ col.type }}</td>
                    <td class="px-3 py-2.5 text-center">
                      <span class="px-1.5 py-0.5 rounded text-[10px]" :class="col.null === 'YES' ? 'bg-slate-800 text-slate-400' : 'bg-rose-500/10 text-rose-300 font-semibold'">
                        {{ col.null }}
                      </span>
                    </td>
                    <td class="px-3 py-2.5 text-center">
                      <span v-if="col.key" class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-amber-500/20 text-amber-300">
                        {{ col.key }}
                      </span>
                      <span v-else class="text-slate-600">-</span>
                    </td>
                    <td class="px-4 py-2.5 text-slate-400">
                      <span v-if="col.default === null" class="italic text-slate-600">NULL</span>
                      <span v-else>{{ col.default }}</span>
                    </td>
                    <td class="px-4 py-2.5 text-purple-300">{{ col.extra || '-' }}</td>
                    <td class="px-4 py-2.5 text-slate-400">{{ col.comment || '-' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </template>
        </div>

        <!-- ============================================================== -->
        <!-- TAB CONTENT 3: SQL QUERY CONSOLE (KONSOL QUERY HEIDISQL F9) -->
        <!-- ============================================================== -->
        <div v-if="activeWorkbenchTab === 'query'" class="flex-1 flex flex-col overflow-hidden p-4 gap-3">
          <!-- Query Controls & Snippets -->
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span class="text-xs font-semibold text-slate-300 flex items-center gap-1">
                <Terminal class="w-3.5 h-3.5 text-emerald-400" />
                <span>Konsol SQL HeidiSQL</span>
              </span>
              
              <!-- Snippets Buttons -->
              <button
                v-if="selectedTable"
                @click="setQuerySnippet(`SELECT * FROM \`${selectedTable}\` LIMIT 50;`)"
                class="px-2 py-0.5 rounded bg-slate-900 hover:bg-slate-800 border border-slate-800 text-[11px] text-slate-400 hover:text-slate-200 transition"
              >
                SELECT * LIMIT 50
              </button>
              <button
                v-if="selectedTable"
                @click="setQuerySnippet(`SELECT COUNT(*) AS total_rows FROM \`${selectedTable}\`;`)"
                class="px-2 py-0.5 rounded bg-slate-900 hover:bg-slate-800 border border-slate-800 text-[11px] text-slate-400 hover:text-slate-200 transition"
              >
                COUNT(*)
              </button>
              <button
                @click="setQuerySnippet('SHOW TABLE STATUS;')"
                class="px-2 py-0.5 rounded bg-slate-900 hover:bg-slate-800 border border-slate-800 text-[11px] text-slate-400 hover:text-slate-200 transition"
              >
                STATUS TABEL
              </button>
            </div>

            <div class="flex items-center gap-2">
              <!-- Execution Time & Rows Badge -->
              <div v-if="queryExecutionStats" class="flex items-center gap-2 text-[11px] font-mono">
                <span class="text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded">
                  ⏱️ {{ queryExecutionStats.timeMs }} ms
                </span>
                <span class="text-sky-400 bg-sky-500/10 border border-sky-500/20 px-2 py-0.5 rounded">
                  📊 {{ queryExecutionStats.rowCount }} baris
                </span>
              </div>

              <!-- Tombol Jalankan SQL (F9) -->
              <button
                @click="executeCustomQuery"
                :disabled="isExecutingQuery"
                class="px-3.5 py-1.5 rounded-lg bg-emerald-500 hover:bg-emerald-400 disabled:opacity-50 text-slate-950 font-bold text-xs flex items-center gap-1.5 transition active:scale-95 shadow-lg shadow-emerald-500/20"
                title="Jalankan Query SQL (Ctrl+Enter)"
              >
                <RotateCw v-if="isExecutingQuery" class="w-3.5 h-3.5 animate-spin" />
                <Play v-else class="w-3.5 h-3.5 fill-current" />
                <span>Jalankan (F9)</span>
              </button>
            </div>
          </div>

          <!-- SQL Editor Textarea -->
          <div class="relative h-40 shrink-0 rounded-xl overflow-hidden border border-slate-800 bg-slate-900">
            <textarea
              v-model="customSqlText"
              @keydown.ctrl.enter.prevent="executeCustomQuery"
              @keydown.f9.prevent="executeCustomQuery"
              class="w-full h-full p-3 font-mono text-xs bg-slate-950 text-emerald-300 outline-none resize-none placeholder-slate-600 leading-relaxed"
              placeholder="Ketik query SQL di sini (misal: SELECT * FROM users LIMIT 10;)... Tekan Ctrl+Enter atau F9 untuk menjalankan."
            ></textarea>
          </div>

          <!-- Error Alert Banner -->
          <div v-if="queryError" class="p-3 rounded-xl bg-rose-500/15 border border-rose-500/30 text-rose-300 text-xs flex items-start gap-2 animate-fadeIn shrink-0">
            <AlertTriangle class="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />
            <div class="flex-1 font-mono leading-relaxed">{{ queryError }}</div>
          </div>

          <!-- Query Result Grid Container -->
          <div class="flex-1 overflow-auto rounded-xl border border-slate-800 bg-slate-900/60">
            <div v-if="!queryResult" class="p-8 text-center text-xs text-slate-500">
              Hasil eksekusi query SQL akan muncul di sini.
            </div>

            <div v-else-if="queryResult.columns.length === 0" class="p-8 text-center text-xs text-emerald-400 font-mono">
              ✓ {{ queryResult.message || 'Query berhasil dieksekusi tanpa hasil baris.' }}
            </div>

            <table v-else class="w-full border-collapse text-left text-xs font-mono">
              <thead class="sticky top-0 bg-slate-900 border-b border-slate-800 divide-x divide-slate-800 shadow">
                <tr>
                  <th class="w-10 px-2.5 py-2 text-center text-slate-500 font-bold">#</th>
                  <th v-for="col in queryResult.columns" :key="col" class="px-3 py-2 font-bold text-slate-200 whitespace-nowrap">
                    {{ col }}
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800/70">
                <tr v-for="(row, idx) in queryResult.rows" :key="idx" class="hover:bg-slate-800/50 divide-x divide-slate-800/50">
                  <td class="px-2.5 py-1.5 text-center text-slate-600 text-[11px]">{{ idx + 1 }}</td>
                  <td v-for="col in queryResult.columns" :key="col" class="px-3 py-1.5 max-w-sm truncate text-slate-300">
                    <span v-if="row[col] === null" class="text-slate-600 italic font-sans text-[11px]">NULL</span>
                    <span v-else>{{ row[col] }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- ============================================================== -->
        <!-- TAB CONTENT 4: TABLE INFO & DUMP -->
        <!-- ============================================================== -->
        <div v-if="activeWorkbenchTab === 'info'" class="flex-1 overflow-y-auto p-6 max-w-2xl mx-auto w-full space-y-6">
          <div v-if="!selectedTable" class="p-12 text-center text-xs text-slate-500">
            Pilih tabel di sebelah kiri untuk melihat rincian informasi dan aksi pemeliharaan.
          </div>

          <template v-else>
            <!-- Header Card -->
            <div class="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 shadow-xl flex items-center justify-between">
              <div>
                <h3 class="text-lg font-bold text-slate-100">{{ selectedTable }}</h3>
                <p class="text-xs text-slate-400">Database: <span class="font-mono text-sky-400">{{ activeDb }}</span></p>
              </div>

              <!-- Quick Export Dump Button -->
              <button
                @click="exportTableSql(selectedTable)"
                class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs flex items-center gap-2 transition active:scale-95 shadow"
              >
                <Download class="w-4 h-4" />
                <span>Unduh SQL Dump</span>
              </button>
            </div>

            <!-- Stats Grid -->
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Storage Engine</span>
                <p class="text-base font-bold text-slate-200 mt-1 font-mono">{{ activeTableStats?.engine || 'InnoDB' }}</p>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Total Baris</span>
                <p class="text-base font-bold text-sky-400 mt-1 font-mono">{{ activeTableStats?.rows || 0 }}</p>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Ukuran Data</span>
                <p class="text-base font-bold text-amber-400 mt-1 font-mono">{{ activeTableStats?.size || '0 KB' }}</p>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Collation</span>
                <p class="text-xs font-bold text-slate-300 mt-1 font-mono truncate">{{ activeTableStats?.collation || 'utf8mb4' }}</p>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Auto Increment</span>
                <p class="text-sm font-bold text-slate-300 mt-1 font-mono">{{ activeTableStats?.autoIncrement || '-' }}</p>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/60 border border-slate-800">
                <span class="text-[11px] text-slate-400">Waktu Dibuat</span>
                <p class="text-xs font-bold text-slate-300 mt-1 font-mono truncate">{{ activeTableStats?.createTime || '-' }}</p>
              </div>
            </div>

            <!-- Danger Zone (Truncate & Drop) -->
            <div class="p-5 rounded-2xl bg-rose-500/5 border border-rose-500/20 space-y-3">
              <h4 class="font-bold text-sm text-rose-300 flex items-center gap-2">
                <AlertTriangle class="w-4 h-4 text-rose-400" />
                <span>Zona Bahaya Pemeliharaan Tabel</span>
              </h4>
              <p class="text-xs text-slate-400 leading-relaxed">
                Tindakan di bawah ini bersifat permanen. Pastikan Anda telah membuat cadangan SQL sebelum mengosongkan atau menghapus tabel.
              </p>

              <div class="flex items-center gap-3 pt-2">
                <button
                  @click="openTruncateModal"
                  class="px-3.5 py-2 rounded-xl bg-amber-500/15 hover:bg-amber-500 text-amber-400 hover:text-slate-950 border border-amber-500/30 text-xs font-semibold transition active:scale-95"
                >
                  Kosongkan Tabel (TRUNCATE)
                </button>
                <button
                  @click="openDropModal"
                  class="px-3.5 py-2 rounded-xl bg-rose-500/15 hover:bg-rose-500 text-rose-400 hover:text-white border border-rose-500/30 text-xs font-semibold transition active:scale-95"
                >
                  Hapus Tabel (DROP TABLE)
                </button>
              </div>
            </div>
          </template>
        </div>

      </div>
    </div>

    <!-- ============================================================== -->
    <!-- MODALS: INSERT/EDIT ROW, CREATE TABLE, TRUNCATE, DROP -->
    <!-- ============================================================== -->

    <!-- Modal 1: Insert / Edit Row -->
    <div v-if="showRowModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col max-h-[85vh] overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3.5 border-b border-slate-800">
          <h3 class="font-bold text-sm flex items-center gap-2">
            <Pencil class="w-4 h-4 text-sky-400" />
            <span>{{ isEditingRow ? 'Edit Baris Data' : 'Tambah Baris Baru' }} ({{ selectedTable }})</span>
          </h3>
          <button @click="showRowModal = false" class="text-slate-400 hover:text-slate-100 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-3.5 overflow-y-auto flex-1">
          <div v-for="col in tableColumns" :key="col" class="space-y-1">
            <div class="flex items-center justify-between text-xs">
              <label class="font-semibold text-slate-300 font-mono">{{ col }}</label>
              <span class="text-[10px] text-slate-500 font-mono">{{ getColumnType(col) }}</span>
            </div>
            <input
              v-model="rowFormData[col]"
              type="text"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 text-xs font-mono outline-none focus:border-sky-400 transition"
              :placeholder="`Masukkan nilai ${col}`"
            />
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showRowModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium transition">
            Batal
          </button>
          <button @click="saveRowData" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition shadow">
            Simpan Baris
          </button>
        </div>
      </div>
    </div>

    <!-- Modal 2: Create Table -->
    <div v-if="showCreateTableModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-xl rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col max-h-[85vh] overflow-hidden text-slate-100">
        <div class="flex items-center justify-between bg-slate-950 px-5 py-3.5 border-b border-slate-800">
          <h3 class="font-bold text-sm flex items-center gap-2">
            <Plus class="w-4 h-4 text-sky-400" />
            <span>Buat Tabel Baru dalam Database: {{ activeDb }}</span>
          </h3>
          <button @click="showCreateTableModal = false" class="text-slate-400 hover:text-slate-100 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-5 space-y-4 overflow-y-auto flex-1 text-xs">
          <div>
            <label class="block font-semibold text-slate-300 mb-1">Nama Tabel</label>
            <input
              v-model.trim="newTableForm.name"
              type="text"
              placeholder="Misal: produk, pesanan, artikel..."
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-700 text-slate-100 outline-none focus:border-sky-400 font-mono"
            />
          </div>

          <!-- Columns Definitions -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="font-semibold text-slate-300">Definisi Kolom</label>
              <button
                type="button"
                @click="addColumnDefinition"
                class="text-xs text-sky-400 hover:text-sky-300 flex items-center gap-1 font-semibold"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Kolom</span>
              </button>
            </div>

            <div class="space-y-2">
              <div 
                v-for="(c, cIdx) in newTableForm.columns" 
                :key="cIdx"
                class="p-2.5 rounded-xl bg-slate-950 border border-slate-800 flex items-center gap-2"
              >
                <input
                  v-model.trim="c.name"
                  type="text"
                  placeholder="Nama Kolom"
                  class="flex-1 px-2.5 py-1.5 rounded-lg bg-slate-900 border border-slate-700 text-xs font-mono outline-none"
                />
                <select
                  v-model="c.type"
                  class="w-28 px-2 py-1.5 rounded-lg bg-slate-900 border border-slate-700 text-xs font-mono outline-none"
                >
                  <option value="INT">INT</option>
                  <option value="BIGINT">BIGINT</option>
                  <option value="VARCHAR(255)">VARCHAR(255)</option>
                  <option value="TEXT">TEXT</option>
                  <option value="DATETIME">DATETIME</option>
                  <option value="TINYINT(1)">TINYINT(1)</option>
                  <option value="DECIMAL(10,2)">DECIMAL(10,2)</option>
                </select>

                <label class="flex items-center gap-1 text-[11px] text-slate-400 cursor-pointer">
                  <input type="checkbox" v-model="c.isPrimary" class="rounded accent-sky-500" />
                  <span>PK</span>
                </label>

                <label class="flex items-center gap-1 text-[11px] text-slate-400 cursor-pointer">
                  <input type="checkbox" v-model="c.autoIncrement" class="rounded accent-sky-500" />
                  <span>AI</span>
                </label>

                <button
                  type="button"
                  @click="removeColumnDefinition(cIdx)"
                  :disabled="newTableForm.columns.length <= 1"
                  class="p-1 rounded text-slate-500 hover:text-rose-400 transition disabled:opacity-20"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="showCreateTableModal = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">Batal</button>
          <button @click="submitCreateTable" class="px-4 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold">Buat Tabel</button>
        </div>
      </div>
    </div>

    <!-- Modal 3: Confirmation Dialog (Custom Modal, Zero Native Alert) -->
    <div v-if="confirmDialog.show" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in">
      <div class="w-full max-w-sm rounded-2xl bg-slate-900 border border-slate-700 shadow-2xl flex flex-col overflow-hidden text-slate-100">
        <div class="p-5 flex flex-col items-center text-center gap-3">
          <div class="w-12 h-12 rounded-full flex items-center justify-center bg-rose-500/20 text-rose-400 border border-rose-500/30">
            <AlertTriangle class="w-6 h-6" />
          </div>
          <div>
            <h4 class="font-bold text-base text-slate-100">{{ confirmDialog.title }}</h4>
            <p class="text-xs text-slate-400 mt-1.5 leading-relaxed">{{ confirmDialog.message }}</p>
          </div>
        </div>
        <div class="flex items-center justify-end gap-2 bg-slate-950 px-5 py-3 border-t border-slate-800">
          <button @click="confirmDialog.show = false" class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium">
            Batal
          </button>
          <button @click="confirmDialog.action" class="px-4 py-2 rounded-xl bg-rose-500 hover:bg-rose-400 text-white text-xs font-bold shadow">
            {{ confirmDialog.confirmText || 'Konfirmasi' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
/**
 * Komponen RAZDatabasesView (HeidiSQL Client Suite)
 * Logika manajemen koneksi database, browser tabel, editor data grid, skema kolom, dan konsol query SQL.
 */
import { ref, reactive, computed, onMounted } from 'vue'
import { useSystemStore } from '../stores/RAZSystemStore'
import { useWindowStore } from '../stores/RAZWindowStore'
import {
  Database,
  Table,
  Columns,
  Terminal,
  Info,
  Plus,
  Trash2,
  Pencil,
  Play,
  Download,
  Search,
  RotateCw,
  ExternalLink,
  Key,
  X,
  AlertTriangle
} from 'lucide-vue-next'

const systemStore = useSystemStore()
const windowStore = useWindowStore()

// State Database & Tables
const databases = ref([])
const activeDb = ref('')
const isLoadingDbs = ref(false)

const tables = ref([])
const selectedTable = ref('')
const isLoadingTables = ref(false)
const tableSearchQuery = ref('')

// Active Workbench Tab: 'data' | 'structure' | 'query' | 'info'
const activeWorkbenchTab = ref('data')

// State Tab 1: Data Grid
const tableColumns = ref([])
const tableRows = ref([])
const tableDataTotal = ref(0)
const tableDataPage = ref(1)
const tableDataLimit = ref(50)
const tableDataTotalPages = ref(1)
const dataSearchQuery = ref('')
const sortColumn = ref('')
const sortOrder = ref('ASC')
const isLoadingData = ref(false)
const tablePrimaryKey = ref('id')

// State Tab 2: Structure
const tableStructure = ref([])
const isLoadingStructure = ref(false)

// State Tab 3: SQL Console
const customSqlText = ref('')
const isExecutingQuery = ref(false)
const queryResult = ref(null)
const queryError = ref('')
const queryExecutionStats = ref(null)

// State Modals
const showRowModal = ref(false)
const isEditingRow = ref(false)
const rowFormData = reactive({})
const originalEditingRow = ref(null)

const showCreateTableModal = ref(false)
const newTableForm = reactive({
  name: '',
  columns: [
    { name: 'id', type: 'INT', isPrimary: true, autoIncrement: true },
    { name: 'name', type: 'VARCHAR(255)', isPrimary: false, autoIncrement: false }
  ]
})

// Custom Confirmation Modal
const confirmDialog = reactive({
  show: false,
  title: '',
  message: '',
  confirmText: '',
  action: null
})

// Computed Filtered Tables
const filteredTables = computed(() => {
  if (!tableSearchQuery.value.trim()) return tables.value
  const q = tableSearchQuery.value.toLowerCase()
  return tables.value.filter(t => t.name.toLowerCase().includes(q))
})

// Computed Active Table Stats
const activeTableStats = computed(() => {
  return tables.value.find(t => t.name === selectedTable.value) || null
})

// Helper Jenis Tipe Data Kolom
const getColumnType = (colName) => {
  const c = tableStructure.value.find(s => s.field === colName)
  return c ? c.type : ''
}

// -----------------------------------------------------------------
// METODE PEMANGGILAN API & MANAJEMEN DATABASE
// -----------------------------------------------------------------

/**
 * Mengambil daftar seluruh database
 */
const fetchAllDatabases = async () => {
  isLoadingDbs.value = true
  try {
    const res = await fetch('/api/v1/databases/list')
    if (res.ok) {
      const data = await res.json()
      databases.value = data
      
      // Pilih database pertama jika belum terpilih
      if (!activeDb.value && data.length > 0) {
        // Cari database non-sistem terlebih dahulu
        const nonSys = data.find(d => !d.isSystem)
        activeDb.value = nonSys ? nonSys.name : data[0].name
        fetchTables()
      }
    }
  } catch (err) {
    console.error('Gagal mengambil daftar database:', err)
  } finally {
    isLoadingDbs.value = false
  }
}

/**
 * Handle pergantian database aktif
 */
const handleDbChange = () => {
  selectedTable.value = ''
  tableRows.value = []
  tableColumns.value = []
  fetchTables()
}

/**
 * Mengambil daftar tabel dalam database aktif
 */
const fetchTables = async () => {
  if (!activeDb.value) return
  isLoadingTables.value = true
  try {
    const res = await fetch(`/api/v1/databases/tables?db=${encodeURIComponent(activeDb.value)}`)
    if (res.ok) {
      const data = await res.json()
      tables.value = data.tables || []
      
      // Jika tabel terpilih sebelumnya ada dalam daftar baru, pertahankan, jika tidak pilih tabel pertama
      if (tables.value.length > 0) {
        if (!tables.value.some(t => t.name === selectedTable.value)) {
          selectTable(tables.value[0].name)
        }
      } else {
        selectedTable.value = ''
      }
    }
  } catch (err) {
    console.error('Gagal mengambil daftar tabel:', err)
  } finally {
    isLoadingTables.value = false
  }
}

/**
 * Memilih tabel untuk ditampilkan di workbench
 */
const selectTable = (tableName) => {
  selectedTable.value = tableName
  tableDataPage.value = 1
  fetchTableData()
  fetchTableStructure()
  customSqlText.value = `SELECT * FROM \`${tableName}\` LIMIT 50;`
}

/**
 * Mengambil data baris tabel (Data Grid Tab)
 */
const fetchTableData = async () => {
  if (!activeDb.value || !selectedTable.value) return
  isLoadingData.value = true
  try {
    const params = new URLSearchParams({
      db: activeDb.value,
      table: selectedTable.value,
      page: tableDataPage.value.toString(),
      limit: tableDataLimit.value.toString()
    })
    if (dataSearchQuery.value.trim()) {
      params.append('search', dataSearchQuery.value.trim())
    }
    if (sortColumn.value) {
      params.append('sort', sortColumn.value)
      params.append('order', sortOrder.value)
    }

    const res = await fetch(`/api/v1/databases/table/data?${params.toString()}`)
    if (res.ok) {
      const data = await res.json()
      tableColumns.value = data.columns || []
      tableRows.value = data.rows || []
      tableDataTotal.value = data.total || 0
      tableDataTotalPages.value = data.totalPages || 1
      tablePrimaryKey.value = data.primaryKey || 'id'
    }
  } catch (err) {
    console.error('Gagal mengambil isi data tabel:', err)
  } finally {
    isLoadingData.value = false
  }
}

/**
 * Mengambil struktur skema kolom tabel
 */
const fetchTableStructure = async () => {
  if (!activeDb.value || !selectedTable.value) return
  isLoadingStructure.value = true
  try {
    const res = await fetch(`/api/v1/databases/table/structure?db=${encodeURIComponent(activeDb.value)}&table=${encodeURIComponent(selectedTable.value)}`)
    if (res.ok) {
      const data = await res.json()
      tableStructure.value = data.columns || []
    }
  } catch (err) {
    console.error('Gagal mengambil struktur tabel:', err)
  } finally {
    isLoadingStructure.value = false
  }
}

/**
 * Sorting kolom data grid
 */
const handleSort = (col) => {
  if (sortColumn.value === col) {
    sortOrder.value = sortOrder.value === 'ASC' ? 'DESC' : 'ASC'
  } else {
    sortColumn.value = col
    sortOrder.value = 'ASC'
  }
  fetchTableData()
}

const changeDataPage = (p) => {
  if (p >= 1 && p <= tableDataTotalPages.value) {
    tableDataPage.value = p
    fetchTableData()
  }
}

// -----------------------------------------------------------------
// KONSOL SQL QUERY (F9)
// -----------------------------------------------------------------

const setQuerySnippet = (sql) => {
  customSqlText.value = sql
}

const executeCustomQuery = async () => {
  const query = customSqlText.value.trim()
  if (!query) return

  isExecutingQuery.value = true
  queryError.value = ''
  try {
    const res = await fetch('/api/v1/databases/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ db: activeDb.value, query })
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      queryError.value = data.error || 'Eksekusi query gagal'
      queryResult.value = null
      queryExecutionStats.value = null
    } else {
      queryResult.value = data
      queryExecutionStats.value = {
        timeMs: data.timeMs || 0,
        rowCount: data.rowCount || (data.rows ? data.rows.length : 0)
      }
      // Jika query DDL/DML, refresh data tabel
      if (!query.toLowerCase().startsWith('select') && !query.toLowerCase().startsWith('show') && !query.toLowerCase().startsWith('describe')) {
        fetchTables()
        fetchTableData()
      }
    }
  } catch (err) {
    queryError.value = err.message || 'Terjadi kesalahan sistem'
  } finally {
    isExecutingQuery.value = false
  }
}

// -----------------------------------------------------------------
// MODAL ROW INSERT & EDIT
// -----------------------------------------------------------------

const openInsertRowModal = () => {
  isEditingRow.value = false
  originalEditingRow.value = null
  Object.keys(rowFormData).forEach(k => delete rowFormData[k])
  tableColumns.value.forEach(c => {
    rowFormData[c] = ''
  })
  showRowModal.value = true
}

const openEditRowModal = (row) => {
  isEditingRow.value = true
  originalEditingRow.value = { ...row }
  Object.keys(rowFormData).forEach(k => delete rowFormData[k])
  tableColumns.value.forEach(c => {
    rowFormData[c] = row[c] === null ? '' : row[c]
  })
  showRowModal.value = true
}

const saveRowData = async () => {
  try {
    let url = '/api/v1/databases/row/insert'
    let payload = {
      db: activeDb.value,
      table: selectedTable.value,
      rowData: { ...rowFormData }
    }

    if (isEditingRow.value && originalEditingRow.value) {
      url = '/api/v1/databases/row/update'
      payload = {
        db: activeDb.value,
        table: selectedTable.value,
        primaryKey: {
          col: tablePrimaryKey.value,
          val: originalEditingRow.value[tablePrimaryKey.value]
        },
        rowData: { ...rowFormData }
      }
    }

    const res = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      systemStore.showNotification('Gagal', data.error || 'Gagal menyimpan baris data', 'error')
    } else {
      systemStore.showNotification('Sukses', data.message || 'Baris data berhasil disimpan', 'success')
      showRowModal.value = false
      fetchTableData()
    }
  } catch (err) {
    systemStore.showNotification('Error', err.message, 'error')
  }
}

const openDeleteRowModal = (row) => {
  confirmDialog.title = 'Hapus Baris Data'
  confirmDialog.message = `Apakah Anda yakin ingin menghapus baris data dengan ${tablePrimaryKey.value} = '${row[tablePrimaryKey.value]}'?`
  confirmDialog.confirmText = 'Hapus Baris'
  confirmDialog.action = async () => {
    confirmDialog.show = false
    try {
      const res = await fetch('/api/v1/databases/row/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          db: activeDb.value,
          table: selectedTable.value,
          primaryKey: {
            col: tablePrimaryKey.value,
            val: row[tablePrimaryKey.value]
          }
        })
      })
      const data = await res.json()
      if (res.ok && data.success) {
        systemStore.showNotification('Sukses', 'Baris data berhasil dihapus', 'success')
        fetchTableData()
      } else {
        systemStore.showNotification('Gagal', data.error || 'Gagal menghapus baris', 'error')
      }
    } catch (err) {
      systemStore.showNotification('Error', err.message, 'error')
    }
  }
  confirmDialog.show = true
}

// -----------------------------------------------------------------
// MODAL CREATE TABLE
// -----------------------------------------------------------------

const openCreateTableModal = () => {
  newTableForm.name = ''
  newTableForm.columns = [
    { name: 'id', type: 'INT', isPrimary: true, autoIncrement: true },
    { name: 'name', type: 'VARCHAR(255)', isPrimary: false, autoIncrement: false }
  ]
  showCreateTableModal.value = true
}

const addColumnDefinition = () => {
  newTableForm.columns.push({
    name: '',
    type: 'VARCHAR(255)',
    isPrimary: false,
    autoIncrement: false
  })
}

const removeColumnDefinition = (index) => {
  if (newTableForm.columns.length > 1) {
    newTableForm.columns.splice(index, 1)
  }
}

const submitCreateTable = async () => {
  if (!newTableForm.name.trim()) {
    systemStore.showNotification('Peringatan', 'Nama tabel harus diisi', 'error')
    return
  }

  try {
    const res = await fetch('/api/v1/databases/table/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        db: activeDb.value,
        tableName: newTableForm.name.trim(),
        columns: newTableForm.columns
      })
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      systemStore.showNotification('Gagal', data.error || 'Gagal membuat tabel baru', 'error')
    } else {
      systemStore.showNotification('Sukses', data.message || 'Tabel berhasil dibuat', 'success')
      showCreateTableModal.value = false
      await fetchTables()
      selectTable(newTableForm.name.trim())
    }
  } catch (err) {
    systemStore.showNotification('Error', err.message, 'error')
  }
}

// -----------------------------------------------------------------
// TRUNCATE & DROP TABLE
// -----------------------------------------------------------------

const openTruncateModal = () => {
  confirmDialog.title = 'Kosongkan Tabel (TRUNCATE)'
  confirmDialog.message = `PERINGATAN: Seluruh data di dalam tabel '${selectedTable.value}' akan dihapus secara permanen. Lanjutkan?`
  confirmDialog.confirmText = 'Kosongkan Sekarang'
  confirmDialog.action = async () => {
    confirmDialog.show = false
    try {
      const res = await fetch('/api/v1/databases/table/truncate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ db: activeDb.value, table: selectedTable.value })
      })
      const data = await res.json()
      if (res.ok && data.success) {
        systemStore.showNotification('Sukses', data.message || 'Tabel berhasil dikosongkan', 'success')
        fetchTableData()
      } else {
        systemStore.showNotification('Gagal', data.error || 'Gagal mengosongkan tabel', 'error')
      }
    } catch (err) {
      systemStore.showNotification('Error', err.message, 'error')
    }
  }
  confirmDialog.show = true
}

const openDropModal = () => {
  confirmDialog.title = 'Hapus Tabel (DROP TABLE)'
  confirmDialog.message = `PERINGATAN: Tabel '${selectedTable.value}' beserta seluruh strukturnya akan dihapus selamanya dari database '${activeDb.value}'. Tindakan ini tidak dapat dibatalkan!`
  confirmDialog.confirmText = 'Hapus Tabel'
  confirmDialog.action = async () => {
    confirmDialog.show = false
    try {
      const res = await fetch('/api/v1/databases/table/drop', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ db: activeDb.value, table: selectedTable.value })
      })
      const data = await res.json()
      if (res.ok && data.success) {
        systemStore.showNotification('Sukses', data.message || 'Tabel berhasil dihapus', 'success')
        selectedTable.value = ''
        fetchTables()
      } else {
        systemStore.showNotification('Gagal', data.error || 'Gagal menghapus tabel', 'error')
      }
    } catch (err) {
      systemStore.showNotification('Error', err.message, 'error')
    }
  }
  confirmDialog.show = true
}

// -----------------------------------------------------------------
// EXPORT SQL DUMP
// -----------------------------------------------------------------

const exportTableSql = (tableName) => {
  if (!activeDb.value || !tableName) return
  window.open(`/api/v1/databases/export?db=${encodeURIComponent(activeDb.value)}&table=${encodeURIComponent(tableName)}`, '_blank')
}

const launchAdminer = async () => {
  try {
    const res = await fetch('/api/v1/databases/adminer/launch', { method: 'POST' })
    if (res.ok) {
      window.open('http://' + window.location.hostname + ':8088', '_blank')
    } else {
      systemStore.showNotification('Gagal', 'Container Adminer sedang tidak aktif', 'error')
    }
  } catch (e) {
    window.open('http://' + window.location.hostname + ':8088', '_blank')
  }
}

onMounted(() => {
  fetchAllDatabases()
})
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
