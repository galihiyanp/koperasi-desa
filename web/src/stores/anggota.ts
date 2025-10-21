import { defineStore } from 'pinia'

export type AnggotaItem = {
  id: number
  nomor_anggota: string
  nama: string
  nik: string
  alamat?: string
  telp?: string
  status: string
  tanggal_gabung?: string
}

export const useAnggotaStore = defineStore('anggota', {
  state: () => ({
    list: [] as AnggotaItem[],
    page: 1,
    limit: 10,
    hasNext: false,
    q: '' as string,
    status: '' as string,
    loading: false,
    error: null as string | null,
    lastLoadedAt: 0,
  }),
  actions: {
    setList(list: AnggotaItem[]) {
      this.list = list
      this.lastLoadedAt = Date.now()
    },
    setPagination(page: number, limit: number, hasNext: boolean) {
      this.page = page
      this.limit = limit
      this.hasNext = hasNext
    },
    setFilters(q: string, status: string) {
      this.q = q
      this.status = status
    },
    setLoading(val: boolean) { this.loading = val },
    setError(err: string | null) { this.error = err },
  },
})