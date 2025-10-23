import { createApp } from 'vue'
import App from './App.vue'
import './index.css'

import { createPinia } from 'pinia'
import router from './router'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Set judul halaman secara konsisten berdasarkan meta rute
const appName = import.meta.env.VITE_APP_NAME || 'Sistem Koperasi Kantor Desa Binangun'
document.title = appName
router.afterEach((to) => {
  const metaTitle = (to.meta && (to.meta as any).title) as string | undefined
  document.title = metaTitle ? `${metaTitle} — ${appName}` : appName
})

app.mount('#app')
