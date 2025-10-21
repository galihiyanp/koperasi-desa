<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import api from '@/lib/axios'
import { useSelectionStore } from '@/stores/selection'
import { useRouter, useRoute } from 'vue-router'
import { useAnggotaStore } from '@/stores/anggota'

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


type Document = { id: number; filename: string; url: string; uploaded_at?: string }

type Pinjaman = {
  id: number
  anggota_id: number
  nomor_pinjaman?: string
  tanggal_pengajuan?: string
  tanggal_disetujui?: string
  tanggal_pencairan?: string
  nominal: number
  tenor_bulan: number
  bunga_persen: number
  status: string
}

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
// Integrasi store global untuk sinkronisasi pilihan anggota
const selection = useSelectionStore()
const router = useRouter()
const route = useRoute()
const anggotaStore = useAnggotaStore()
const isProfilePage = computed(() => route.name === 'keanggotaan-profil')

const documents = ref<Document[]>([])
const pinjaman = ref<Pinjaman[]>([])
const profileLoading = ref(false)
const profileError = ref<string | null>(null)
const pinjamanError = ref<string | null>(null)

// Pastikan tautan dokumen mengarah ke server backend yang menyajikan folder uploads
const apiBase = import.meta.env.DEV ? '/' : (import.meta.env.VITE_API_BASE_URL || '')
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
  anggotaStore.setLoading(true)
  anggotaStore.setError(null)
  try {
    const res = await api.get('/api/anggota', { params: { page: page.value, limit: limit.value, q: q.value || undefined, status: status.value || undefined } })
    const data = res.data?.data ?? []
    anggota.value = data
    hasNext.value = (data.length ?? 0) >= limit.value
    // Sinkronisasi ke store agar state tetap bertahan saat navigasi
    anggotaStore.setList(data)
    anggotaStore.setPagination(page.value, limit.value, hasNext.value)
    anggotaStore.setFilters(q.value, status.value)
    // Preselect berdasarkan store jika ada
    const pre = selection.selectedAnggotaId
    if (pre != null) {
      const found = anggota.value.find(x => x.id === pre) || null
      if (found) selected.value = found
    }
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal memuat data anggota'
    anggotaStore.setError(error.value)
  } finally {
    loading.value = false
    anggotaStore.setLoading(false)
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



async function verifySelected() {
  if (!selected.value) return
  try {
    await api.post(`/api/anggota/${selected.value.id}/verify`)
    await fetchAnggota()
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Verifikasi gagal')
  }
}

async function activateSelected() {
  if (!selected.value) return
  try {
    await api.post(`/api/anggota/${selected.value.id}/activate`)
    await fetchAnggota()
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Aktivasi gagal')
  }
}

async function loadProfileById(id: number) {
  profileLoading.value = true
  profileError.value = null
  pinjamanError.value = null
  selected.value = null
  documents.value = []
  pinjaman.value = []
  try {
    const res = await api.get(`/api/anggota/${id}`)
    selected.value = res.data?.data ?? null
    const docsRes = await api.get(`/api/anggota/${id}/documents`)
    documents.value = docsRes.data?.data ?? []
  } catch (e: any) {
    profileError.value = e?.response?.data?.error || e?.message || 'Gagal memuat profil anggota'
  }
  try {
    const pjRes = await api.get('/api/pinjaman', { params: { anggota_id: id, limit: 100 } })
    pinjaman.value = pjRes.data?.data ?? []
  } catch (e: any) {
    pinjamanError.value = e?.response?.data?.error || e?.message || 'Gagal memuat data pinjaman'
  } finally {
    profileLoading.value = false
  }
}

function openProfilePage(id: number) {
  const href = router.resolve({ name: 'keanggotaan-profil', params: { id } }).href
  window.open(href, '_blank')
}

function onFilesChange(e: Event) {
  const input = e.target as HTMLInputElement
  form.files = Array.from(input.files || [])
}

function prevPage() { if (page.value > 1) { page.value -= 1; fetchAnggota() } }
function nextPage() { if (hasNext.value) { page.value += 1; fetchAnggota() } }

onMounted(() => {
  if (isProfilePage.value) {
    const id = Number(route.params.id)
    if (id) loadProfileById(id)
  } else {
    // Hydrate dari store jika tersedia, untuk mempertahankan state
    if (anggotaStore.list.length > 0) {
      anggota.value = anggotaStore.list
      page.value = anggotaStore.page
      limit.value = anggotaStore.limit
      hasNext.value = anggotaStore.hasNext
      q.value = anggotaStore.q
      status.value = anggotaStore.status
      loading.value = anggotaStore.loading
      error.value = anggotaStore.error
    } else {
      fetchAnggota()
    }
  }
})

watch(() => isProfilePage.value, (isProfile) => {
  if (isProfile) {
    const id = Number(route.params.id)
    if (id) loadProfileById(id)
  } else {
    // Saat kembali ke halaman daftar, gunakan data dari store
    if (anggotaStore.list.length > 0) {
      anggota.value = anggotaStore.list
      page.value = anggotaStore.page
      limit.value = anggotaStore.limit
      hasNext.value = anggotaStore.hasNext
      q.value = anggotaStore.q
      status.value = anggotaStore.status
      loading.value = anggotaStore.loading
      error.value = anggotaStore.error
    } else {
      // Fallback: jika store kosong (mis. tab baru), lakukan fetch
      fetchAnggota()
    }
  }
})
</script>

<template>
  <div class="dashboard">
    <template v-if="isProfilePage">
      <h1 class="page-title">Profil Anggota</h1>
      <div v-if="profileLoading" class="muted">Memuat profil...</div>
      <div v-if="profileError" class="card error" style="margin-bottom:12px;">{{ profileError }}</div>
      <div v-if="!selected && !profileLoading && !profileError" class="card"><p>Data anggota tidak ditemukan.</p></div>

      <div v-if="selected" class="card" style="margin-top: 8px;">
        <div class="card-header">{{ selected?.nama }} ({{ selected?.nomor_anggota }})</div>
        <div class="card-content profile-grid">
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
      </div>

      <div class="card" style="margin-top: 16px;">
        <div class="card-header">Data Peminjaman</div>
        <div class="card-content">
          <div v-if="pinjamanError" class="muted" style="margin-bottom:12px;">{{ pinjamanError }}</div>
          <table class="table">
            <thead>
              <tr>
                <th>Nomor</th>
                <th>Tanggal Pengajuan</th>
                <th>Nominal</th>
                <th>Tenor</th>
                <th>Bunga</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in pinjaman" :key="p.id">
                <td>{{ p.nomor_pinjaman || '-' }}</td>
                <td>{{ p.tanggal_pengajuan || '-' }}</td>
                <td>{{ new Intl.NumberFormat('id-ID',{style:'currency',currency:'IDR',maximumFractionDigits:0}).format(p.nominal || 0) }}</td>
                <td>{{ p.tenor_bulan }} bln</td>
                <td>{{ p.bunga_persen }}%</td>
                <td><span class="status-badge" :class="p.status">{{ p.status }}</span></td>
              </tr>
              <tr v-if="!pinjaman.length">
                <td colspan="6" class="muted">Tidak ada data pinjaman</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div style="margin-top:12px;">
        <button class="btn btn-secondary" @click="router.push('/keanggotaan')">Kembali ke Keanggotaan</button>
      </div>
    </template>

    <template v-else>
      <h1 class="page-title">Keanggotaan</h1>
      
      <!-- Baris Aksi -->
      <div class="actions-row">
        <button class="btn btn-primary" @click="openCreate">
          <i class="fas fa-plus"></i> Tambah Anggota
        </button>
        
        <!-- Filter & Search -->
        <div class="search-filters">
          <input v-model="q" type="text" placeholder="Cari nama, NIK, atau no. anggota..." class="input search-input">
          <select v-model="status" class="input status-filter">
            <option value="">Semua Status</option>
            <option value="pending">Pending</option>
            <option value="verified">Verified</option>
            <option value="active">Active</option>
          </select>
          <button class="btn btn-secondary" @click="fetchAnggota">Terapkan Filter</button>
        </div>
      </div>

      <!-- Form Pendaftaran/Edit -->
      <div v-if="showForm" class="card form-card">
        <h3>{{ mode === 'create' ? 'Pendaftaran Anggota Baru' : 'Edit Anggota' }}</h3>
        <form @submit.prevent="submitForm" class="form-grid">
          <div class="form-group">
            <label class="label">Nomor Anggota</label>
            <input v-model="form.nomor_anggota" type="text" placeholder="Contoh: AG-2025-0001" class="input">
          </div>
          <div class="form-group">
            <label class="label">Nama Lengkap *</label>
            <input v-model="form.nama" type="text" placeholder="Nama lengkap" class="input" required>
          </div>
          <div class="form-group">
            <label class="label">NIK *</label>
            <input v-model="form.nik" type="text" placeholder="Nomor Induk Kependudukan" class="input" required>
          </div>
          <div class="form-group">
            <label class="label">Alamat</label>
            <input v-model="form.alamat" type="text" placeholder="Alamat domisili" class="input">
          </div>
          <div class="form-group">
            <label class="label">Telepon</label>
            <input v-model="form.telp" type="tel" placeholder="08xxxxxxxxxx" class="input">
          </div>
          <div class="form-group">
            <label class="label">Tanggal Gabung</label>
            <input v-model="form.tanggal_gabung" type="date" class="input">
          </div>
          <div class="form-group">
            <label class="label">Upload Dokumen</label>
            <input type="file" multiple @change="onFilesChange" class="input">
          </div>
          <div class="form-actions">
            <button type="submit" class="btn btn-primary">
              {{ mode === 'create' ? 'Daftar' : 'Update' }}
            </button>
            <button type="button" @click="cancelForm" class="btn btn-secondary">
              Batal
            </button>
          </div>
        </form>
      </div>

      <!-- Tabel Data Anggota -->
      <div class="card">
        <div class="card-header">Verifikasi & Aktivasi Anggota</div>
        <div class="card-content">
          <div v-if="error" class="muted" style="margin-bottom: 12px;">{{ error }}</div>
          <div v-if="loading" class="muted">Memuat...</div>
          <div v-else class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>Nomor</th>
                  <th>Nama</th>
                  <th>NIK</th>
                  <th>Status</th>
                  <th>Gabung</th>
                  <th style="width: 360px;">Aksi</th>
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
                    <div class="table-actions" style="display:flex; gap:6px; flex-wrap:wrap;">
                      <button class="btn btn-secondary" @click="openEdit(a)">Ubah</button>
                      <button class="btn btn-secondary" @click="openProfilePage(a.id)">Profil</button>
                      <button class="btn btn-primary" :disabled="a.status === 'verified' || a.status === 'active'" @click="selected = a; verifySelected()">Verifikasi</button>
                      <button class="btn btn-primary" :disabled="a.status === 'active'" @click="selected = a; activateSelected()">Aktivasi</button>
                      <button class="btn btn-light" @click="selection.setAnggota(a.id); router.push('/pinjaman')">Pinjaman</button>
                      <button class="btn btn-light" @click="selection.setAnggota(a.id); router.push('/angsuran')">Angsuran</button>
                    </div>
                  </td>
                </tr>
                <tr v-if="!anggota.length">
                  <td colspan="6" class="muted">Tidak ada data</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="pagination">
            <button class="btn btn-secondary" :disabled="page <= 1" @click="prevPage">Sebelumnya</button>
            <span class="muted">Halaman {{ page }}</span>
            <button class="btn btn-secondary" :disabled="!hasNext" @click="nextPage">Berikutnya</button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* Actions Row */
.actions-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.search-filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  min-width: 250px;
}

.status-filter {
  min-width: 150px;
}

/* Form Card */
.form-card {
  margin-bottom: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  gap: 12px;
  justify-content: flex-start;
}

/* Table Container */
.table-container {
  overflow-x: auto;
}

/* Profile Grid (untuk halaman profil) */
.profile-grid { 
  display: grid; 
  grid-template-columns: 1fr; 
  gap: 16px; 
}

@media (min-width: 768px) { 
  .profile-grid { 
    grid-template-columns: 1fr 1fr; 
  } 
}

/* Status Badge */
.status-badge { 
  padding: 2px 6px; 
  border-radius: 12px; 
  font-size: 12px; 
}

/* Error Card */
.card.error { 
  border: 1px solid #e33; 
  color: #b00; 
}

/* Responsive Design */
@media (max-width: 768px) {
  .actions-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-filters {
    flex-direction: column;
  }
  
  .search-input,
  .status-filter {
    min-width: auto;
    width: 100%;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
}
</style>