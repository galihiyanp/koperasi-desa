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
  { name: 'dashboard', path: '/dashboard', component: Dashboard, meta: { title: 'Dashboard' } },
  { name: 'keanggotaan', path: '/keanggotaan', component: Keanggotaan, meta: { title: 'Keanggotaan' } },
  // Rute halaman profil anggota menggunakan komponen yang sama
  { name: 'keanggotaan-profil', path: '/keanggotaan/profil/:id', component: Keanggotaan, meta: { title: 'Profil Anggota' } },
  { name: 'simpanan', path: '/simpanan', component: Simpanan, meta: { title: 'Simpanan' } },
  { name: 'penarikan', path: '/penarikan', component: Penarikan, meta: { title: 'Penarikan' } },
  { name: 'pinjaman', path: '/pinjaman', component: Pinjaman, meta: { title: 'Pinjaman' } },
  { name: 'angsuran', path: '/angsuran', component: Angsuran, meta: { title: 'Angsuran' } },
  { name: 'kas', path: '/kas', component: Kas, meta: { title: 'Kas & Jurnal' } },
  { name: 'laporan', path: '/laporan', component: Laporan, meta: { title: 'Laporan' } },
  { name: 'pengaturan', path: '/pengaturan', component: Pengaturan, meta: { title: 'Pengaturan' } },
  { name: 'pengguna', path: '/pengguna', component: UsersSecurity, meta: { title: 'Pengguna & Keamanan' } },
  { name: 'bantuan', path: '/bantuan', component: Bantuan, meta: { title: 'Bantuan' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router