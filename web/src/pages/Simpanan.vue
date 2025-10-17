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

type Recap = { wajib: number; sukarela: number; khusus: number; total: number }

const anggotaList = ref<Anggota[]>([])
const selectedAnggotaId = ref<number | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const transaksi = ref<Transaksi[]>([])

const setoranForm = reactive<{ jenis: 'wajib' | 'sukarela' | 'khusus'; jumlah: number; tanggal: string }>({
  jenis: 'wajib',
  jumlah: 0,
  tanggal: new Date().toISOString().slice(0, 10),
})


const filteredTransaksi = computed(() => {
  const sid = selectedAnggotaId.value
  return transaksi.value.filter((t) => (sid != null ? t.anggota_id === sid : true))
})

const recap = computed<Recap>(() => {
  const sums: Recap = { wajib: 0, sukarela: 0, khusus: 0, total: 0 }
  for (const t of filteredTransaksi.value) {
    const sign = t.tipe === 'penarikan' ? -1 : 1
    if (t.jenis === 'wajib') sums.wajib += sign * (t.jumlah || 0)
    else if (t.jenis === 'sukarela') sums.sukarela += sign * (t.jumlah || 0)
    else if (t.jenis === 'khusus') sums.khusus += sign * (t.jumlah || 0)
  }
  sums.total = (sums.wajib || 0) + (sums.sukarela || 0) + (sums.khusus || 0)
  return sums
})

function formatCurrency(n?: number) {
  const v = n ?? 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}

async function fetchAnggota() {
  try {
    const res = await api.get('/api/anggota', { params: { limit: 100 } })
    anggotaList.value = res.data?.data ?? res.data ?? []
    if (selectedAnggotaId.value == null && anggotaList.value.length > 0) {
      const first = anggotaList.value[0]
      if (first) selectedAnggotaId.value = first.id
    }
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
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Gagal memuat data simpanan'
    transaksi.value = []
  } finally {
    loading.value = false
  }
}

async function submitSetoran() {
  const sid = selectedAnggotaId.value
  if (sid == null) return alert('Pilih anggota terlebih dahulu')
  if (!setoranForm.jumlah || setoranForm.jumlah <= 0) return alert('Jumlah setoran tidak valid')
  try {
    const payload = {
      anggota_id: sid,
      jenis: setoranForm.jenis,
      jumlah: setoranForm.jumlah,
      tanggal: new Date(`${setoranForm.tanggal}T00:00:00`).toISOString(),
    }
    const res = await api.post('/api/simpanan/setoran', payload)
    if (res.data?.ok) {
      setoranForm.jumlah = 0
      await fetchTransaksi()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Setoran gagal')
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
    <h1 class="page-title">Simpanan</h1>

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

    <!-- Forms -->
    <div style="display: grid; grid-template-columns: 1fr; gap: 16px;">
      <div class="card">
        <div class="card-header">Setoran Simpanan</div>
        <div class="card-content">
          <div class="form-row">
            <label class="label">Jenis</label>
            <select class="input" v-model="setoranForm.jenis">
              <option value="wajib">Wajib</option>
              <option value="sukarela">Sukarela</option>
              <option value="khusus">Khusus</option>
            </select>
          </div>
          <div class="form-row">
            <label class="label">Tanggal</label>
            <input class="input" v-model="setoranForm.tanggal" type="date" />
          </div>
          <div class="form-row">
            <label class="label">Jumlah</label>
            <input class="input" v-model.number="setoranForm.jumlah" type="number" min="0" placeholder="0" />
          </div>
          <div class="form-actions">
            <button class="btn btn-primary" @click="submitSetoran" :disabled="!selectedAnggotaId">Setor</button>
          </div>
        </div>
      </div>

    </div>

    <!-- Recap -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Rekap Saldo Simpanan</div>
      <div class="card-content">
        <div style="display: grid; grid-template-columns: 1fr; gap: 8px;">
          <div style="display:flex; justify-content: space-between;">
            <span>Wajib</span>
            <strong>{{ formatCurrency(recap.wajib) }}</strong>
          </div>
          <div style="display:flex; justify-content: space-between;">
            <span>Sukarela</span>
            <strong>{{ formatCurrency(recap.sukarela) }}</strong>
          </div>
          <div style="display:flex; justify-content: space-between;">
            <span>Khusus</span>
            <strong>{{ formatCurrency(recap.khusus) }}</strong>
          </div>
          <hr style="border:0; border-top:1px solid #f3f4f6; margin: 6px 0;" />
          <div style="display:flex; justify-content: space-between;">
            <span>Total</span>
            <strong>{{ formatCurrency(recap.total) }}</strong>
          </div>
        </div>
      </div>
    </div>

    <!-- History -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Riwayat Transaksi Simpanan</div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Jenis</th>
              <th>Tipe</th>
              <th>Jumlah</th>
              <th>Saldo Akhir</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in filteredTransaksi" :key="t.id">
              <td>{{ new Date(t.tanggal).toLocaleDateString('id-ID') }}</td>
              <td style="text-transform: capitalize">{{ t.jenis }}</td>
              <td style="text-transform: capitalize">
                <span class="status-badge" :class="t.tipe === 'setoran' ? 'verified' : 'pending'">{{ t.tipe }}</span>
              </td>
              <td>{{ formatCurrency(t.jumlah) }}</td>
              <td>{{ t.saldo_akhir != null ? formatCurrency(t.saldo_akhir) : '-' }}</td>
            </tr>
            <tr v-if="!filteredTransaksi.length">
              <td colspan="5" class="muted">Belum ada transaksi</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>