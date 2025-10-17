<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/lib/axios'

type Anggota = {
  id: number
  nomor_anggota: string
  nama: string
  nik: string
  alamat?: string
  telp?: string
  status: 'pending' | 'verified' | 'active' | string
  tanggal_gabung?: string
  created_at?: string
  updated_at?: string
}

type Activity = { id: number; action: string; note?: string; created_at?: string }
type Document = { id: number; filename: string; url: string; uploaded_at?: string }

const anggota = ref<Anggota[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const page = ref(1)
const limit = ref(10)
const hasNext = ref(false)

const q = ref('')
const status = ref<string>('')

const showForm = ref(false)
const mode = ref<'create' | 'edit'>('create')
const form = reactive<{ id?: number; nomor_anggota: string; nama: string; nik: string; alamat: string; telp: string; tanggal_gabung?: string; files: File[] }>({
  id: undefined,
  nomor_anggota: '',
  nama: '',
  nik: '',
  alamat: '',
  telp: '',
  tanggal_gabung: undefined,
  files: [],
})

const selected = ref<Anggota | null>(null)
const activities = ref<Activity[]>([])
const documents = ref<Document[]>([])

// Pastikan tautan dokumen mengarah ke server backend yang menyajikan folder uploads
const apiBase = import.meta.env.VITE_API_BASE_URL || ''
function docUrl(d: Document) {
  const url = d.url || ''
  if (/^https?:\/\//.test(url)) return url
  const base = apiBase.endsWith('/') ? apiBase : apiBase + '/'
  const path = url.startsWith('/') ? url.slice(1) : url
  return base + path
}

async function fetchAnggota() {
  loading.value = true
  error.value = null
  try {
    const res = await api.get('/api/anggota', { params: { page: page.value, limit: limit.value, q: q.value || undefined, status: status.value || undefined } })
    const data = res.data?.data ?? []
    anggota.value = data
    hasNext.value = (data.length ?? 0) >= limit.value
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal memuat data anggota'
  } finally {
    loading.value = false
  }
}

function openCreate() {
  mode.value = 'create'
  showForm.value = true
  form.id = undefined
  form.nomor_anggota = ''
  form.nama = ''
  form.nik = ''
  form.alamat = ''
  form.telp = ''
  form.tanggal_gabung = undefined
  form.files = []
}

function openEdit(a: Anggota) {
  mode.value = 'edit'
  showForm.value = true
  form.id = a.id
  form.nomor_anggota = a.nomor_anggota
  form.nama = a.nama
  form.nik = a.nik
  form.alamat = a.alamat || ''
  form.telp = a.telp || ''
  form.tanggal_gabung = a.tanggal_gabung
  form.files = []
}

function cancelForm() { showForm.value = false }

async function submitForm() {
  try {
    if (mode.value === 'create') {
      const payload: Record<string, any> = {
        nomor_anggota: form.nomor_anggota,
        nama: form.nama,
        nik: form.nik,
        alamat: form.alamat,
        telp: form.telp,
        tanggal_gabung: form.tanggal_gabung,
      }
      const res = await api.post('/api/anggota', payload)
      const created = res.data
      if (form.files.length) {
        const fd = new FormData()
        form.files.forEach((f) => fd.append('files', f))
        await api.post(`/api/anggota/${created.id}/documents`, fd, { headers: { 'Content-Type': 'multipart/form-data' } })
      }
    } else if (mode.value === 'edit' && form.id != null) {
      const payload: Record<string, any> = {
        nomor_anggota: form.nomor_anggota,
        nama: form.nama,
        nik: form.nik,
        alamat: form.alamat,
        telp: form.telp,
        tanggal_gabung: form.tanggal_gabung,
      }
      await api.put(`/api/anggota/${form.id}`, payload)
      if (form.files.length) {
        const fd = new FormData()
        form.files.forEach((f) => fd.append('files', f))
        await api.post(`/api/anggota/${form.id}/documents`, fd, { headers: { 'Content-Type': 'multipart/form-data' } })
      }
    }
    await fetchAnggota()
    showForm.value = false
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Operasi gagal')
  }
}

async function viewProfile(a: Anggota) {
  selected.value = a
  try {
    const res = await api.get(`/api/anggota/${a.id}`)
    const data = res.data?.data ?? a
    selected.value = data
    activities.value = res.data?.activities ?? []
    const docsRes = await api.get(`/api/anggota/${a.id}/documents`)
    documents.value = docsRes.data?.data ?? []
  } catch (e: any) {
    activities.value = []
    documents.value = []
  }
}

async function verifySelected() {
  if (!selected.value) return
  try {
    await api.post(`/api/anggota/${selected.value.id}/verify`)
    await viewProfile(selected.value)
    await fetchAnggota()
  } catch (e: any) {
    alert('Verifikasi gagal')
  }
}

async function activateSelected() {
  if (!selected.value) return
  try {
    await api.post(`/api/anggota/${selected.value.id}/activate`)
    await viewProfile(selected.value)
    await fetchAnggota()
  } catch (e: any) {
    alert('Aktivasi gagal')
  }
}

function onFilesChange(e: Event) {
  const input = e.target as HTMLInputElement
  form.files = Array.from(input.files || [])
}

function prevPage() { if (page.value > 1) { page.value -= 1; fetchAnggota() } }
function nextPage() { if (hasNext.value) { page.value += 1; fetchAnggota() } }

onMounted(() => { fetchAnggota() })
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Keanggotaan</h1>

    <div class="actions-row" style="margin-bottom: 8px;">
      <button class="btn btn-primary" @click="openCreate">Pendaftaran Anggota</button>
      <input class="input" v-model="q" type="text" placeholder="Cari nama/nik/nomor anggota" style="max-width: 280px;" />
      <select class="input" v-model="status" style="max-width: 180px;">
        <option value="">Semua Status</option>
        <option value="pending">Pending</option>
        <option value="verified">Verified</option>
        <option value="active">Active</option>
      </select>
      <button class="btn btn-secondary" @click="fetchAnggota">Terapkan</button>
    </div>

    <!-- Form Pendaftaran / Ubah Anggota -->
    <div v-if="showForm" class="card" style="margin-top: 8px;">
      <div class="card-header">{{ mode === 'create' ? 'Pendaftaran Anggota' : 'Ubah Data Anggota' }}</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Nomor Anggota</label>
          <input class="input" v-model="form.nomor_anggota" type="text" placeholder="Contoh: AG-2025-0001" />
        </div>
        <div class="form-row">
          <label class="label">Nama Lengkap</label>
          <input class="input" v-model="form.nama" type="text" placeholder="Nama lengkap" />
        </div>
        <div class="form-row">
          <label class="label">NIK</label>
          <input class="input" v-model="form.nik" type="text" placeholder="Nomor Induk Kependudukan" />
        </div>
        <div class="form-row">
          <label class="label">Alamat</label>
          <input class="input" v-model="form.alamat" type="text" placeholder="Alamat domisili" />
        </div>
        <div class="form-row">
          <label class="label">Telepon</label>
          <input class="input" v-model="form.telp" type="text" placeholder="08xxxxxxxxxx" />
        </div>
        <div class="form-row">
          <label class="label">Tanggal Gabung</label>
          <input class="input" v-model="form.tanggal_gabung" type="date" />
        </div>
        <div class="form-row">
          <label class="label">Upload Dokumen</label>
          <input class="input" type="file" multiple @change="onFilesChange" />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" @click="submitForm">Simpan</button>
          <button class="btn btn-secondary" @click="cancelForm">Batal</button>
        </div>
      </div>
    </div>

    <!-- Daftar Anggota -->
    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Verifikasi & Aktivasi Anggota</div>
      <div class="card-content">
        <div v-if="error" class="muted" style="margin-bottom: 12px;">{{ error }}</div>
        <div v-if="loading" class="muted">Memuat...</div>
        <table v-else class="table">
          <thead>
            <tr>
              <th>Nomor</th>
              <th>Nama</th>
              <th>NIK</th>
              <th>Status</th>
              <th>Gabung</th>
              <th style="width: 280px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in anggota" :key="a.id">
              <td>{{ a.nomor_anggota }}</td>
              <td>{{ a.nama }}</td>
              <td>{{ a.nik }}</td>
              <td>
                <span class="status-badge" :class="a.status">{{ a.status }}</span>
              </td>
              <td>{{ a.tanggal_gabung || '-' }}</td>
              <td>
                <div class="table-actions">
                  <button class="btn btn-secondary" @click="openEdit(a)">Ubah</button>
                  <button class="btn btn-secondary" @click="viewProfile(a)">Profil</button>
                  <button class="btn btn-primary" :disabled="a.status === 'verified' || a.status === 'active'" @click="selected = a; verifySelected()">Verifikasi</button>
                  <button class="btn btn-primary" :disabled="a.status === 'active'" @click="selected = a; activateSelected()">Aktivasi</button>
                </div>
              </td>
            </tr>
            <tr v-if="!anggota.length">
              <td colspan="6" class="muted">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
        <div class="pagination">
          <button class="btn btn-secondary" :disabled="page <= 1" @click="prevPage">Sebelumnya</button>
          <span class="muted">Halaman {{ page }}</span>
          <button class="btn btn-secondary" :disabled="!hasNext" @click="nextPage">Berikutnya</button>
        </div>
      </div>
    </div>

    <!-- Profil & Riwayat Aktivitas -->
    <div v-if="selected" class="card" style="margin-top: 16px;">
      <div class="card-header">Profil Anggota: {{ selected?.nama }} ({{ selected?.nomor_anggota }})</div>
      <div class="card-content">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
          <div>
            <p><strong>NIK:</strong> {{ selected?.nik }}</p>
            <p><strong>Alamat:</strong> {{ selected?.alamat || '-' }}</p>
            <p><strong>Telepon:</strong> {{ selected?.telp || '-' }}</p>
            <p><strong>Status:</strong> <span class="status-badge" :class="selected?.status">{{ selected?.status }}</span></p>
            <p><strong>Gabung:</strong> {{ selected?.tanggal_gabung || '-' }}</p>
          </div>
          <div>
            <p class="label" style="margin-bottom: 6px;">Dokumen</p>
            <ul>
              <li v-for="d in documents" :key="d.id">
                <a :href="docUrl(d)" target="_blank" rel="noopener noreferrer">{{ d.filename }}</a>
              </li>
              <li v-if="!documents.length" class="muted">Tidak ada dokumen</li>
            </ul>
          </div>
        </div>
        <div style="margin-top: 16px;">
          <p class="label" style="margin-bottom: 6px;">Riwayat Aktivitas</p>
          <ul>
            <li v-for="act in activities" :key="act.id">
              {{ act.created_at }} - <strong>{{ act.action }}</strong> <span v-if="act.note">({{ act.note }})</span>
            </li>
            <li v-if="!activities.length" class="muted">Belum ada aktivitas</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>