<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSelectionStore } from '@/stores/selection'
import api from '@/lib/axios'

type Anggota = { id: number; nama: string; nomor_anggota: string }
type Pinjaman = {
  id: number
  anggota_id: number
  nomor_pinjaman?: string
  tanggal_pengajuan?: string
  nominal: number
  tenor_bulan: number
  bunga_persen: number
  status: 'pengajuan' | 'disetujui' | 'berjalan' | 'lunas' | string
}

const anggotaList = ref<Anggota[]>([])
const selection = useSelectionStore()
const selectedAnggotaId = computed<number | null>({
  get() { return selection.selectedAnggotaId },
  set(v) { selection.setAnggota(v) },
})
const statusFilter = ref<string>('')
const loading = ref(false)
const error = ref<string | null>(null)

const pinjamanList = ref<Pinjaman[]>([])
const router = useRouter()

const pengajuanForm = reactive<{ anggota_id: number | null; nominal: number; tenor_bulan: number; bunga_persen: number; tanggal: string }>(
  {
    anggota_id: null,
    nominal: 0,
    tenor_bulan: 12,
    bunga_persen: 12,
    tanggal: new Date().toISOString().slice(0, 10),
  }
)

function formatCurrency(n?: number) {
  const v = n ?? 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}

async function fetchAnggota() {
  try {
    const res = await api.get('/api/anggota', { params: { limit: 100 } })
    anggotaList.value = res.data?.data ?? res.data ?? []
    const first = anggotaList.value[0]
    if (selectedAnggotaId.value == null && first) selectedAnggotaId.value = first.id
    if (pengajuanForm.anggota_id == null && first) pengajuanForm.anggota_id = first.id
  } catch (e: any) {
    anggotaList.value = []
  }
}

async function fetchPinjaman() {
  loading.value = true
  error.value = null
  try {
    const params: Record<string, any> = {}
    if (selectedAnggotaId.value != null) params.anggota_id = selectedAnggotaId.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/pinjaman', { params })
    pinjamanList.value = res.data?.data ?? res.data ?? []
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Gagal memuat data pinjaman'
    pinjamanList.value = []
  } finally {
    loading.value = false
  }
}

async function submitPengajuan() {
  if (pengajuanForm.anggota_id == null) return alert('Pilih anggota untuk pengajuan')
  if (!pengajuanForm.nominal || pengajuanForm.nominal <= 0) return alert('Nominal pinjaman tidak valid')
  if (!pengajuanForm.tenor_bulan || pengajuanForm.tenor_bulan <= 0) return alert('Tenor tidak valid')
  try {
    const payload = {
      anggota_id: pengajuanForm.anggota_id,
      nominal: pengajuanForm.nominal,
      tenor_bulan: pengajuanForm.tenor_bulan,
      bunga_persen: pengajuanForm.bunga_persen,
      tanggal_pengajuan: new Date(`${pengajuanForm.tanggal}T00:00:00`).toISOString(),
    }
    const res = await api.post('/api/pinjaman/pengajuan', payload)
    if (res.data) {
      pengajuanForm.nominal = 0
      await fetchPinjaman()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Pengajuan pinjaman gagal')
  }
}

async function verifikasiPinjaman(id: number) {
  try {
    const res = await api.post('/api/pinjaman/verifikasi', { pinjaman_id: id })
    if (res.data) await fetchPinjaman()
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Verifikasi gagal')
  }
}

async function pencairanPinjaman(id: number) {
  try {
    const res = await api.post('/api/pinjaman/pencairan', { pinjaman_id: id })
    if (res.data) await fetchPinjaman()
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Pencairan gagal')
  }
}

function gotoAngsuran(p: Pinjaman) {
  selection.setPinjaman(p.id)
  router.push('/angsuran')
}

onMounted(async () => {
  await fetchAnggota()
  await fetchPinjaman()
})
watch(selectedAnggotaId, () => { fetchPinjaman() })
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Pinjaman</h1>

    <!-- Filters -->
    <div class="actions-row" style="margin-bottom: 8px; flex-wrap: wrap;">
      <select class="input" v-model.number="selectedAnggotaId" style="min-width: 240px;">
        <option :value="null">Pilih Anggota</option>
        <option v-for="a in anggotaList" :key="a.id" :value="a.id">{{ a.nama }} ({{ a.nomor_anggota }})</option>
      </select>
      <select class="input" v-model="statusFilter" style="min-width: 200px;">
        <option value="">Semua Status</option>
        <option value="pengajuan">Pengajuan</option>
        <option value="disetujui">Disetujui</option>
        <option value="berjalan">Berjalan</option>
        <option value="lunas">Lunas</option>
      </select>
      <button class="btn btn-secondary" @click="fetchPinjaman">Muat Ulang</button>
      <span class="muted" v-if="loading">Memuat data…</span>
      <span class="muted" v-if="error">{{ error }}</span>
    </div>

    <!-- Pengajuan Form -->
    <div class="card">
      <div class="card-header">Pengajuan Pinjaman</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Anggota</label>
          <select class="input" v-model.number="pengajuanForm.anggota_id" required>
            <option v-for="a in anggotaList" :key="a.id" :value="a.id">{{ a.nama }} ({{ a.nomor_anggota }})</option>
          </select>
        </div>
        <div class="form-row">
          <label class="label">Tanggal Pengajuan</label>
          <input class="input" v-model="pengajuanForm.tanggal" type="date" required />
        </div>
        <div class="form-row">
          <label class="label">Nominal</label>
          <input class="input" v-model.number="pengajuanForm.nominal" type="number" min="1" placeholder="0" required />
        </div>
        <div class="form-row" style="display:grid; grid-template-columns: 1fr 1fr; gap: 8px;">
          <div>
            <label class="label">Tenor (bulan)</label>
            <input class="input" v-model.number="pengajuanForm.tenor_bulan" type="number" min="1" placeholder="12" required />
          </div>
          <div>
            <label class="label">Bunga (%)</label>
            <input class="input" v-model.number="pengajuanForm.bunga_persen" type="number" min="0" step="0.1" placeholder="12" />
          </div>
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" @click="submitPengajuan">Ajukan</button>
        </div>
      </div>
    </div>

    <!-- List Pinjaman -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Daftar Pinjaman</div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>Nomor</th>
              <th>Anggota</th>
              <th>Tanggal</th>
              <th>Nominal</th>
              <th>Tenor</th>
              <th>Bunga</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in pinjamanList" :key="p.id">
              <td>{{ p.nomor_pinjaman || '-' }}</td>
              <td>{{ p.anggota_id }}</td>
              <td>{{ p.tanggal_pengajuan ? new Date(p.tanggal_pengajuan).toLocaleDateString('id-ID') : '-' }}</td>
              <td>{{ formatCurrency(p.nominal) }}</td>
              <td>{{ p.tenor_bulan }} bln</td>
              <td>{{ p.bunga_persen }}%</td>
              <td style="text-transform: capitalize">{{ p.status }}</td>
              <td>
                <div style="display:flex; gap:6px;">
                  <button class="btn btn-secondary" @click="verifikasiPinjaman(p.id)" :disabled="p.status !== 'pengajuan'">Verifikasi</button>
                  <button class="btn btn-primary" @click="pencairanPinjaman(p.id)" :disabled="p.status !== 'disetujui'">Cairkan</button>
                  <button class="btn btn-light" @click="gotoAngsuran(p)">Angsuran</button>
                </div>
              </td>
            </tr>
            <tr v-if="!pinjamanList.length">
              <td colspan="8" class="muted">Belum ada data pinjaman</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>