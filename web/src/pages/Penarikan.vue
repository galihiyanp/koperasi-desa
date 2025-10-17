<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import api from '@/lib/axios'

type Anggota = { id: number; nama: string; nomor_anggota: string }
type Transaksi = {
  id: number
  anggota_id: number
  jenis: 'wajib' | 'sukarela' | 'khusus' | string
  tipe: 'setoran' | 'penarikan' | string
  tanggal: string
  jumlah: number
  saldo_akhir?: number
  created_at?: string
}
type Saldo = { wajib: number; sukarela: number; khusus: number }

const anggotaList = ref<Anggota[]>([])
const selectedAnggotaId = ref<number | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const transaksi = ref<Transaksi[]>([])
const saldo = ref<Saldo>({ wajib: 0, sukarela: 0, khusus: 0 })

const penarikanForm = reactive<{ jenis: 'wajib' | 'sukarela' | 'khusus'; jumlah: number; tanggal: string }>({
  jenis: 'sukarela',
  jumlah: 0,
  tanggal: new Date().toISOString().slice(0, 10),
})

const filteredPenarikan = computed(() => {
  const sid = selectedAnggotaId.value
  return transaksi.value.filter((t) => (sid != null ? t.anggota_id === sid : true)).filter((t) => t.tipe === 'penarikan')
})

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
  } catch (e: any) {
    anggotaList.value = []
  }
}

async function fetchTransaksi() {
  loading.value = true
  error.value = null
  try {
    const params: Record<string, any> = {}
    if (selectedAnggotaId.value != null) params.anggota_id = selectedAnggotaId.value
    const res = await api.get('/api/simpanan', { params })
    transaksi.value = res.data?.data ?? res.data ?? []
    saldo.value = res.data?.saldo ?? { wajib: 0, sukarela: 0, khusus: 0 }
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Gagal memuat data penarikan'
    transaksi.value = []
    saldo.value = { wajib: 0, sukarela: 0, khusus: 0 }
  } finally {
    loading.value = false
  }
}

async function submitPenarikan() {
  const sid = selectedAnggotaId.value
  if (sid == null) return alert('Pilih anggota terlebih dahulu')
  if (!penarikanForm.jumlah || penarikanForm.jumlah <= 0) return alert('Jumlah penarikan tidak valid')
  const available = penarikanForm.jenis === 'wajib'
    ? saldo.value.wajib
    : penarikanForm.jenis === 'sukarela'
      ? saldo.value.sukarela
      : saldo.value.khusus
  if (penarikanForm.jumlah > available) return alert('Saldo tidak cukup')
  try {
    const payload = {
      anggota_id: sid,
      jenis: penarikanForm.jenis,
      jumlah: penarikanForm.jumlah,
      tanggal: new Date(`${penarikanForm.tanggal}T00:00:00`).toISOString(),
    }
    const res = await api.post('/api/simpanan/penarikan', payload)
    if (res.data?.ok) {
      penarikanForm.jumlah = 0
      await fetchTransaksi()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Penarikan gagal')
  }
}

watch(selectedAnggotaId, () => {
  fetchTransaksi()
})

onMounted(async () => {
  await fetchAnggota()
  await fetchTransaksi()
})
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Penarikan</h1>

    <!-- Filters -->
    <div class="actions-row" style="margin-bottom: 8px; flex-wrap: wrap;">
      <select class="input" v-model.number="selectedAnggotaId" style="min-width: 240px;">
        <option :value="null">Pilih Anggota</option>
        <option v-for="a in anggotaList" :key="a.id" :value="a.id">{{ a.nama }} ({{ a.nomor_anggota }})</option>
      </select>
      <button class="btn btn-secondary" @click="fetchTransaksi">Muat Ulang</button>
      <span class="muted" v-if="loading">Memuat data…</span>
      <span class="muted" v-if="error">{{ error }}</span>
    </div>

    <!-- Penarikan Form -->
    <div class="card">
      <div class="card-header">Penarikan Simpanan</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Jenis (wajib diisi)</label>
          <select class="input" v-model="penarikanForm.jenis" required>
            <option value="wajib">Wajib</option>
            <option value="sukarela">Sukarela</option>
            <option value="khusus">Khusus</option>
          </select>
        </div>
        <div class="form-row">
          <label class="label">Tanggal</label>
          <input class="input" v-model="penarikanForm.tanggal" type="date" required />
        </div>
        <div class="form-row">
          <label class="label">Jumlah</label>
          <input class="input" v-model.number="penarikanForm.jumlah" type="number" min="1" :max="penarikanForm.jenis === 'wajib' ? saldo.wajib : penarikanForm.jenis === 'sukarela' ? saldo.sukarela : saldo.khusus" placeholder="0" required />
        </div>
        <div class="muted" style="margin-bottom: 8px;">
          Saldo tersedia ({{ penarikanForm.jenis }}):
          <strong>{{ formatCurrency(penarikanForm.jenis === 'wajib' ? saldo.wajib : penarikanForm.jenis === 'sukarela' ? saldo.sukarela : saldo.khusus) }}</strong>
        </div>
        <div class="form-actions">
          <button class="btn btn-danger" @click="submitPenarikan" :disabled="!selectedAnggotaId">Tarik</button>
        </div>
      </div>
    </div>

    <!-- Recap Available -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Rekap Saldo Tersedia</div>
      <div class="card-content">
        <div style="display: grid; grid-template-columns: 1fr; gap: 8px;">
          <div style="display:flex; justify-content: space-between;">
            <span>Wajib</span>
            <strong>{{ formatCurrency(saldo.wajib) }}</strong>
          </div>
          <div style="display:flex; justify-content: space-between;">
            <span>Sukarela</span>
            <strong>{{ formatCurrency(saldo.sukarela) }}</strong>
          </div>
          <div style="display:flex; justify-content: space-between;">
            <span>Khusus</span>
            <strong>{{ formatCurrency(saldo.khusus) }}</strong>
          </div>
        </div>
      </div>
    </div>

    <!-- History: Penarikan only -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Riwayat Penarikan</div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Jenis</th>
              <th>Jumlah</th>
              <th>Saldo Akhir</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in filteredPenarikan" :key="t.id">
              <td>{{ new Date(t.tanggal).toLocaleDateString('id-ID') }}</td>
              <td style="text-transform: capitalize">{{ t.jenis }}</td>
              <td>{{ formatCurrency(t.jumlah) }}</td>
              <td>{{ t.saldo_akhir != null ? formatCurrency(t.saldo_akhir) : '-' }}</td>
            </tr>
            <tr v-if="!filteredPenarikan.length">
              <td colspan="4" class="muted">Belum ada penarikan</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>