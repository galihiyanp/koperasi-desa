import { defineStore } from 'pinia'

export const useSelectionStore = defineStore('selection', {
  state: (): { selectedAnggotaId: number | null; selectedPinjamanId: number | null } => ({
    selectedAnggotaId: null,
    selectedPinjamanId: null,
  }),
  actions: {
    setAnggota(id: number | null) {
      this.selectedAnggotaId = id
    },
    setPinjaman(id: number | null) {
      this.selectedPinjamanId = id
    },
    reset() {
      this.selectedAnggotaId = null
      this.selectedPinjamanId = null
    },
  },
})