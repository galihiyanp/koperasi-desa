import { defineStore } from 'pinia'
import api from '@/lib/axios'

type Task = { title: string; description: string; badge?: string }

type Overview = {
  anggotaCount?: number
  simpananTotal?: number
  pinjamanAktif?: number
  kasHarian?: number
  tasks: Task[]
}

export const useDashboardStore = defineStore('dashboard', {
  state: (): { overview: Overview } => ({
    overview: { tasks: [] },
  }),
  actions: {
    async fetchOverview() {
      try {
        const [anggotaRes, pinjamanRes, kasRes] = await Promise.allSettled([
          api.get('/api/anggota'),
          api.get('/api/pinjaman'),
          api.get('/api/kas?periode=today'),
        ])

        const anggotaCount = anggotaRes.status === 'fulfilled' ? (anggotaRes.value.data?.length ?? 0) : 128
        const pinjamanAktif = pinjamanRes.status === 'fulfilled' ? (pinjamanRes.value.data?.filter((p: any) => p.status === 'berjalan').length ?? 0) : 24
        const kasHarian = kasRes.status === 'fulfilled' ? (kasRes.value.data?.net ?? 0) : 1500000

        // Simpanan total contoh (mock). Endpoint bisa ditambahkan kemudian.
        const simpananTotal = 125_000_000

        this.overview = {
          anggotaCount,
          simpananTotal,
          pinjamanAktif,
          kasHarian,
          tasks: [
            { title: 'Pengajuan pinjaman menunggu persetujuan', description: '5 pengajuan baru', badge: 'Prioritas' },
            { title: 'Angsuran jatuh tempo minggu ini', description: '12 anggota perlu diingatkan', badge: 'Pengingat' },
          ],
        }
      } catch (e) {
        // Fallback jika service belum tersedia
        this.overview = {
          anggotaCount: 128,
          simpananTotal: 125_000_000,
          pinjamanAktif: 24,
          kasHarian: 1_500_000,
          tasks: [
            { title: 'Pengajuan pinjaman menunggu persetujuan', description: '5 pengajuan baru', badge: 'Prioritas' },
            { title: 'Angsuran jatuh tempo minggu ini', description: '12 anggota perlu diingatkan', badge: 'Pengingat' },
          ],
        }
      }
    },
  },
})