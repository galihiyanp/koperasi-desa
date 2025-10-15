import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const Dashboard = () => import('@/pages/Dashboard.vue')
const Placeholder = () => import('@/pages/Placeholder.vue')
const UsersSecurity = () => import('@/pages/UsersSecurity.vue')

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/dashboard' },
  { name: 'dashboard', path: '/dashboard', component: Dashboard },
  { name: 'keanggotaan', path: '/keanggotaan', component: Placeholder },
  { name: 'simpanan', path: '/simpanan', component: Placeholder },
  { name: 'penarikan', path: '/penarikan', component: Placeholder },
  { name: 'pinjaman', path: '/pinjaman', component: Placeholder },
  { name: 'angsuran', path: '/angsuran', component: Placeholder },
  { name: 'kas', path: '/kas', component: Placeholder },
  { name: 'laporan', path: '/laporan', component: Placeholder },
  { name: 'pengaturan', path: '/pengaturan', component: Placeholder },
  { name: 'pengguna', path: '/pengguna', component: UsersSecurity },
  { name: 'bantuan', path: '/bantuan', component: Placeholder },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router