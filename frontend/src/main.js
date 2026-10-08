/**
 * RAZOS Frontend Application Entrypoint
 * Updated: 2026-10-08
 * Deskripsi: Inisialisasi Vue 3, Pinia state store, dan stylesheet tema RAZOS.
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/RAZTheme.css'

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
