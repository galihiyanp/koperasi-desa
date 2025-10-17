<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/lib/axios'

type KasEntry = { id: number; tanggal: string; jenis: 'in' | 'out' | string; keterangan?: string; jumlah: number; ref?: string }

const kasList = ref<KasEntry[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const kasInForm = reactive<{ tanggal: string; keterangan: string; jumlah: number; ref: string }>(
  { tanggal: new Date().toISOString().slice(0, 10), keterangan: '', jumlah: 0, ref: '' }
)
const kasOutForm = reactive<{ tanggal: string; keterangan: string; jumlah: number; ref: string }>(
  { tanggal: new Date().toISOString().slice(0, 10), keterangan: '', jumlah: 0, ref: '' }
)

function formatCurrency(n?: number) {
  const v = n ?? 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}

async function fetchKas() {
  loading.value = true
  error.value = null
  try {
    const res = await api.get('/api/kas')
    kasList.value = res.data?.data ?? res.data ?? []
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Gagal memuat data kas'
    kasList.value = []
  } finally {
    loading.value = false
  }
}

async function submitKasIn() {
  if (!kasInForm.jumlah || kasInForm.jumlah <= 0) return alert('Jumlah penerimaan tidak valid')
  try {
    const payload = {
      tanggal: new Date(`${kasInForm.tanggal}T00:00:00`).toISOString(),
      keterangan: kasInForm.keterangan,
      jumlah: kasInForm.jumlah,
      ref: kasInForm.ref,
    }
    const res = await api.post('/api/kas/in', payload)
    if (res.data) {
      kasInForm.jumlah = 0
      kasInForm.keterangan = ''
      kasInForm.ref = ''
      await fetchKas()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Pencatatan kas-in gagal')
  }
}

async function submitKasOut() {
  if (!kasOutForm.jumlah || kasOutForm.jumlah <= 0) return alert('Jumlah pengeluaran tidak valid')
  try {
    const payload = {
      tanggal: new Date(`${kasOutForm.tanggal}T00:00:00`).toISOString(),
      keterangan: kasOutForm.keterangan,
      jumlah: kasOutForm.jumlah,
      ref: kasOutForm.ref,
    }
    const res = await api.post('/api/kas/out', payload)
    if (res.data) {
      kasOutForm.jumlah = 0
      kasOutForm.keterangan = ''
      kasOutForm.ref = ''
      await fetchKas()
    }
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Pencatatan kas-out gagal')
  }
}

onMounted(async () => {
  await fetchKas()
})
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Kas & Jurnal</h1>

    <div class="actions-row" style="margin-bottom: 8px; flex-wrap: wrap;">
      <button class="btn btn-secondary" @click="fetchKas">Muat Ulang</button>
      <span class="muted" v-if="loading">Memuat data…</span>
      <span class="muted" v-if="error">{{ error }}</span>
    </div>

    <!-- Kas In -->
    <div class="card">
      <div class="card-header">Penerimaan Kas (In)</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Tanggal</label>
          <input class="input" v-model="kasInForm.tanggal" type="date" required />
        </div>
        <div class="form-row">
          <label class="label">Keterangan</label>
          <input class="input" v-model="kasInForm.keterangan" type="text" placeholder="Keterangan" />
        </div>
        <div class="form-row">
          <label class="label">Jumlah</label>
          <input class="input" v-model.number="kasInForm.jumlah" type="number" min="1" placeholder="0" required />
        </div>
        <div class="form-row">
          <label class="label">Ref</label>
          <input class="input" v-model="kasInForm.ref" type="text" placeholder="Ref transaksi" />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" @click="submitKasIn">Catat In</button>
        </div>
      </div>
    </div>

    <!-- Kas Out -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Pengeluaran Kas (Out)</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Tanggal</label>
          <input class="input" v-model="kasOutForm.tanggal" type="date" required />
        </div>
        <div class="form-row">
          <label class="label">Keterangan</label>
          <input class="input" v-model="kasOutForm.keterangan" type="text" placeholder="Keterangan" />
        </div>
        <div class="form-row">
          <label class="label">Jumlah</label>
          <input class="input" v-model.number="kasOutForm.jumlah" type="number" min="1" placeholder="0" required />
        </div>
        <div class="form-row">
          <label class="label">Ref</label>
          <input class="input" v-model="kasOutForm.ref" type="text" placeholder="Ref transaksi" />
        </div>
        <div class="form-actions">
          <button class="btn btn-danger" @click="submitKasOut">Catat Out</button>
        </div>
      </div>
    </div>

    <!-- List Kas -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Buku Kas</div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Jenis</th>
              <th>Keterangan</th>
              <th>Jumlah</th>
              <th>Ref</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="k in kasList" :key="k.id">
              <td>{{ new Date(k.tanggal).toLocaleDateString('id-ID') }}</td>
              <td style="text-transform: capitalize">{{ k.jenis }}</td>
              <td>{{ k.keterangan || '-' }}</td>
              <td>{{ formatCurrency(k.jumlah) }}</td>
              <td>{{ k.ref || '-' }}</td>
            </tr>
            <tr v-if="!kasList.length">
              <td colspan="5" class="muted">Belum ada transaksi kas</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>