<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import api from '@/lib/axios'
import { useSelectionStore } from '@/stores/selection'
import { useRouter } from 'vue-router'

type Pinjaman = { id: number; nomor_pinjaman?: string; anggota_id: number; status?: string }
type Angsuran = { id: number; pinjaman_id: number; ke: number; tanggal_jatuh_tempo: string; jumlah: number; tanggal_bayar?: string; denda?: number }

type Anggota = { id: number; nama: string; nomor_anggota: string }

const pinjamanList = ref<Pinjaman[]>([])
// Replace local ref with computed synced to global store
const selection = useSelectionStore()
const selectedPinjamanId = computed<number | null>({
  get() { return selection.selectedPinjamanId },
  set(v) { selection.setPinjaman(v) },
})
const router = useRouter()
const angsuranList = ref<Angsuran[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const anggotaList = ref<Anggota[]>([])

const bayarForm = reactive<{ angsuran_id: number | null; jumlah: number; tanggal: string }>(
  { angsuran_id: null, jumlah: 0, tanggal: new Date().toISOString().slice(0, 10) }
)

const belumBayar = computed(() => angsuranList.value.filter(a => !a.tanggal_bayar))

function formatCurrency(n?: number) {
  const v = n ?? 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}

async function fetchAnggota() {
  try {
    const res = await api.get('/api/anggota', { params: { limit: 100 } })
    anggotaList.value = res.data?.data ?? res.data ?? []
  } catch (e: any) {
    anggotaList.value = []
  }
}

async function fetchPinjaman() {
  try {
    const params: Record<string, any> = { status: 'berjalan' }
    if (selection.selectedAnggotaId != null) params.anggota_id = selection.selectedAnggotaId
    const res = await api.get('/api/pinjaman', { params })
    pinjamanList.value = res.data?.data ?? res.data ?? []
    const first = pinjamanList.value[0]
    if (selectedPinjamanId.value == null && first) selectedPinjamanId.value = first.id
  } catch (e: any) {
    pinjamanList.value = []
  }
}

async function fetchAngsuran() {
  loading.value = true
  error.value = null
  try {
    const params: Record<string, any> = {}
    if (selectedPinjamanId.value != null) params.pinjaman_id = selectedPinjamanId.value
    const res = await api.get('/api/angsuran', { params })
    angsuranList.value = res.data?.data ?? res.data ?? []
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Gagal memuat data angsuran'
    angsuranList.value = []
  } finally {
    loading.value = false
  }
}

async function submitBayar() {
  if (bayarForm.angsuran_id == null) return alert('Pilih angsuran yang akan dibayar')
  if (!bayarForm.jumlah || bayarForm.jumlah <= 0) return alert('Jumlah bayar tidak valid')
  try {
    const payload = {
      angsuran_id: bayarForm.angsuran_id,
      jumlah: bayarForm.jumlah,
      tanggal_bayar: new Date(`${bayarForm.tanggal}T00:00:00`).toISOString(),
    }
    const res = await api.post('/api/angsuran/bayar', payload)
    if (res.data) {
      bayarForm.jumlah = 0
      await fetchAngsuran()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Pembayaran angsuran gagal')
  }
}

onMounted(async () => {
  await fetchAnggota()
  await fetchPinjaman()
  await fetchAngsuran()
})
watch(selectedPinjamanId, () => { fetchAngsuran() })
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Angsuran</h1>

    <!-- Filters -->
    <div class="actions-row" style="margin-bottom: 8px; flex-wrap: wrap;">
      <select class="input" v-model.number="selectedPinjamanId" style="min-width: 240px;">
        <option :value="null">Pilih Pinjaman</option>
        <option v-for="p in pinjamanList" :key="p.id" :value="p.id">{{ p.nomor_pinjaman || ('Pinjaman #' + p.id) }}</option>
      </select>
      <button class="btn btn-secondary" @click="fetchAngsuran">Muat Ulang</button>
      <button class="btn btn-light" @click="router.push('/pinjaman')">Lihat Pinjaman</button>
      <span class="muted" v-if="loading">Memuat data…</span>
      <span class="muted" v-if="error">{{ error }}</span>
    </div>

    <!-- Schedule -->
    <div class="card">
      <div class="card-header">Jadwal Angsuran</div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>Ke</th>
              <th>Jatuh Tempo</th>
              <th>Jumlah</th>
              <th>Tanggal Bayar</th>
              <th>Denda</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in angsuranList" :key="a.id">
              <td>{{ a.ke }}</td>
              <td>{{ new Date(a.tanggal_jatuh_tempo).toLocaleDateString('id-ID') }}</td>
              <td>{{ formatCurrency(a.jumlah) }}</td>
              <td>{{ a.tanggal_bayar ? new Date(a.tanggal_bayar).toLocaleDateString('id-ID') : '-' }}</td>
              <td>{{ a.denda != null ? formatCurrency(a.denda) : '-' }}</td>
            </tr>
            <tr v-if="!angsuranList.length">
              <td colspan="5" class="muted">Belum ada jadwal</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Payment Form -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Pembayaran Angsuran</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Pilih Angsuran</label>
          <select class="input" v-model.number="bayarForm.angsuran_id" required>
            <option :value="null">Pilih</option>
            <option v-for="a in belumBayar" :key="a.id" :value="a.id">Angsuran ke-{{ a.ke }} ({{ formatCurrency(a.jumlah) }})</option>
          </select>
        </div>
        <div class="form-row">
          <label class="label">Tanggal Bayar</label>
          <input class="input" v-model="bayarForm.tanggal" type="date" required />
        </div>
        <div class="form-row">
          <label class="label">Jumlah Bayar</label>
          <input class="input" v-model.number="bayarForm.jumlah" type="number" min="1" placeholder="0" required />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" @click="submitBayar" :disabled="!bayarForm.angsuran_id">Bayar</button>
        </div>
      </div>
    </div>
  </div>
</template>