import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const Dashboard = () => import('@/pages/Dashboard.vue')
const Placeholder = () => import('@/pages/Placeholder.vue')
const UsersSecurity = () => import('@/pages/UsersSecurity.vue')
const Keanggotaan = () => import('@/pages/Keanggotaan.vue')
const Simpanan = () => import('@/pages/Simpanan.vue')
const Penarikan = () => import('@/pages/Penarikan.vue')
const Pinjaman = () => import('@/pages/Pinjaman.vue')
const Angsuran = () => import('@/pages/Angsuran.vue')
const Kas = () => import('@/pages/Kas.vue')
const Laporan = () => import('@/pages/Laporan.vue')
const Pengaturan = () => import('@/pages/Pengaturan.vue')
const Bantuan = () => import('@/pages/Bantuan.vue')

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/dashboard' },
  { name: 'dashboard', path: '/dashboard', component: Dashboard },
  { name: 'keanggotaan', path: '/keanggotaan', component: Keanggotaan },
  // Rute halaman profil anggota menggunakan komponen yang sama
  { name: 'keanggotaan-profil', path: '/keanggotaan/profil/:id', component: Keanggotaan },
  { name: 'simpanan', path: '/simpanan', component: Simpanan },
  { name: 'penarikan', path: '/penarikan', component: Penarikan },
  { name: 'pinjaman', path: '/pinjaman', component: Pinjaman },
  { name: 'angsuran', path: '/angsuran', component: Angsuran },
  { name: 'kas', path: '/kas', component: Kas },
  { name: 'laporan', path: '/laporan', component: Laporan },
  { name: 'pengaturan', path: '/pengaturan', component: Pengaturan },
  { name: 'pengguna', path: '/pengguna', component: UsersSecurity },
  { name: 'bantuan', path: '/bantuan', component: Bantuan },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router