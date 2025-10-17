<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import api from '@/lib/axios'

// Tabs: anggota, simpanan, pinjaman, angsuran, kas
const tab = ref<'anggota' | 'simpanan' | 'pinjaman' | 'angsuran' | 'kas'>('anggota')

// Common utils
function formatCurrency(n?: number) {
  const v = n ?? 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(v)
}
function formatDate(s?: string | Date | null) {
  if (!s) return '-'
  const d = typeof s === 'string' ? new Date(s) : s
  if (!d || isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('id-ID')
}

// Anggota laporan
 type Anggota = { id: number; nomor_anggota: string; nama: string; nik?: string; status: string; tanggal_gabung?: string }
 const anggota = ref<Anggota[]>([])
 const anggotaStatus = ref<string>('') // '', pending, verified, active
 const anggotaLoading = ref(false)
 const anggotaError = ref<string | null>(null)
 async function fetchAnggota() {
   anggotaLoading.value = true
   anggotaError.value = null
   try {
     const res = await api.get('/api/anggota', { params: { limit: 200, status: anggotaStatus.value || undefined } })
     anggota.value = res.data?.data ?? []
   } catch (e: any) {
     anggotaError.value = e?.message ?? 'Gagal memuat laporan anggota'
   } finally {
     anggotaLoading.value = false
   }
 }

 // Simpanan laporan
 type Simpanan = { id: number; anggota_id: number; jenis: 'wajib'|'sukarela'|'khusus'|string; tipe: 'setoran'|'penarikan'|string; tanggal: string; jumlah: number; saldo_akhir?: number }
 const simpanan = ref<Simpanan[]>([])
 const simpananJenis = ref<string>('') // '', wajib, sukarela, khusus
 const simpananPeriode = ref<'harian'|'mingguan'|'bulanan'>('harian')
 const simpananStart = ref<string>('')
 const simpananEnd = ref<string>('')
 const simpananLoading = ref(false)
 const simpananError = ref<string | null>(null)
 function computeDefaultRange() {
   const now = new Date()
   if (simpananPeriode.value === 'harian') {
     simpananStart.value = new Date(now.getFullYear(), now.getMonth(), now.getDate()).toISOString().slice(0,10)
     simpananEnd.value = simpananStart.value
   } else if (simpananPeriode.value === 'mingguan') {
     const day = now.getDay() // 0-6
     const diffStart = day === 0 ? 6 : day - 1
     const start = new Date(now)
     start.setDate(now.getDate() - diffStart)
     const end = new Date(start)
     end.setDate(start.getDate() + 6)
     simpananStart.value = start.toISOString().slice(0,10)
     simpananEnd.value = end.toISOString().slice(0,10)
   } else {
     const start = new Date(now.getFullYear(), now.getMonth(), 1)
     const end = new Date(now.getFullYear(), now.getMonth() + 1, 0)
     simpananStart.value = start.toISOString().slice(0,10)
     simpananEnd.value = end.toISOString().slice(0,10)
   }
 }
 async function fetchSimpanan() {
   simpananLoading.value = true
   simpananError.value = null
   try {
     const res = await api.get('/api/simpanan', { params: { limit: 500, jenis: simpananJenis.value || undefined } })
     const all: Simpanan[] = res.data?.data ?? []
     // Client-side filter by tanggal range
     if (simpananStart.value && simpananEnd.value) {
       const start = new Date(simpananStart.value + 'T00:00:00')
       const end = new Date(simpananEnd.value + 'T23:59:59')
       simpanan.value = all.filter(x => {
         const t = new Date(x.tanggal)
         return t >= start && t <= end
       })
     } else {
       simpanan.value = all
     }
   } catch (e: any) {
     simpananError.value = e?.message ?? 'Gagal memuat laporan simpanan'
   } finally {
     simpananLoading.value = false
   }
 }
 watch(simpananPeriode, () => computeDefaultRange())

 // Pinjaman laporan
 type Pinjaman = { id: number; anggota_id: number; nomor_pinjaman?: string; tanggal_pengajuan?: string; nominal: number; tenor_bulan: number; bunga_persen: number; status: string }
 const pinjaman = ref<Pinjaman[]>([])
 const pinjamanStatus = ref<string>('') // pengajuan|disetujui|berjalan|lunas
 const pinjamanLoading = ref(false)
 const pinjamanError = ref<string | null>(null)
 async function fetchPinjaman() {
   pinjamanLoading.value = true
   pinjamanError.value = null
   try {
     const res = await api.get('/api/pinjaman', { params: { limit: 200, status: pinjamanStatus.value || undefined } })
     pinjaman.value = res.data?.data ?? []
   } catch (e: any) {
     pinjamanError.value = e?.message ?? 'Gagal memuat laporan pinjaman'
   } finally {
     pinjamanLoading.value = false
   }
 }

 // Angsuran laporan
 type Angsuran = { id: number; pinjaman_id: number; ke: number; tanggal_jatuh_tempo: string; jumlah: number; tanggal_bayar?: string|null; denda?: number }
 const angsuran = ref<Angsuran[]>([])
 const angsuranMode = ref<'semua'|'bayar'|'tunggakan'|'denda'>('semua')
 const angsuranLoading = ref(false)
 const angsuranError = ref<string | null>(null)
 const angsuranFiltered = computed(() => {
   const now = new Date()
   if (angsuranMode.value === 'bayar') return angsuran.value.filter(a => a.tanggal_bayar)
   if (angsuranMode.value === 'tunggakan') return angsuran.value.filter(a => !a.tanggal_bayar && new Date(a.tanggal_jatuh_tempo) < now)
   if (angsuranMode.value === 'denda') return angsuran.value.filter(a => (a.denda ?? 0) > 0)
   return angsuran.value
 })
 async function fetchAngsuran() {
   angsuranLoading.value = true
   angsuranError.value = null
   try {
     const res = await api.get('/api/angsuran', { params: { /* fetch all; backend supports pinjaman_id filter only */ } })
     angsuran.value = res.data?.data ?? []
   } catch (e: any) {
     angsuranError.value = e?.message ?? 'Gagal memuat laporan angsuran'
   } finally {
     angsuranLoading.value = false
   }
 }

 // Kas laporan (placeholder integrasi; backend belum menyediakan /api/kas di router)
 type Kas = { id: number; tanggal: string; jenis: 'in'|'out'|string; keterangan?: string; jumlah: number; ref?: string }
 const kas = ref<Kas[]>([])
 const kasPeriode = ref<'harian'|'mingguan'|'bulanan'>('harian')
 const kasStart = ref('')
 const kasEnd = ref('')
 const kasLoading = ref(false)
 const kasError = ref<string | null>(null)
 async function fetchKas() {
   kasLoading.value = true
   kasError.value = null
   try {
     // Attempt: if backend supports /api/kas, otherwise show message
     const res = await api.get('/api/kas', { params: { periode: kasPeriode.value } })
     kas.value = res.data?.data ?? []
   } catch (e: any) {
     kasError.value = 'Endpoint /api/kas belum tersedia di backend'
     kas.value = []
   } finally {
     kasLoading.value = false
   }
 }
 function computeKasRange() {
   const now = new Date()
   if (kasPeriode.value === 'harian') {
     kasStart.value = new Date(now.getFullYear(), now.getMonth(), now.getDate()).toISOString().slice(0,10)
     kasEnd.value = kasStart.value
   } else if (kasPeriode.value === 'mingguan') {
     const day = now.getDay()
     const diffStart = day === 0 ? 6 : day - 1
     const start = new Date(now)
     start.setDate(now.getDate() - diffStart)
     const end = new Date(start)
     end.setDate(start.getDate() + 6)
     kasStart.value = start.toISOString().slice(0,10)
     kasEnd.value = end.toISOString().slice(0,10)
   } else {
     const start = new Date(now.getFullYear(), now.getMonth(), 1)
     const end = new Date(now.getFullYear(), now.getMonth() + 1, 0)
     kasStart.value = start.toISOString().slice(0,10)
     kasEnd.value = end.toISOString().slice(0,10)
   }
 }
 watch(kasPeriode, () => computeKasRange())

 // Export helpers
 function toCSV(rows: any[], headers: { key: string; label: string }[]) {
   const headerLine = headers.map(h => '"' + h.label + '"').join(',')
   const lines = rows.map(r => headers.map(h => {
     const v = (r as any)[h.key]
     const s = typeof v === 'number' ? v.toString() : (v ?? '').toString()
     return '"' + s.replace(/"/g, '""') + '"'
   }).join(','))
   return [headerLine, ...lines].join('\n')
 }
 function downloadFile(filename: string, content: string, mime = 'text/csv') {
   const blob = new Blob([content], { type: mime })
   const url = URL.createObjectURL(blob)
   const a = document.createElement('a')
   a.href = url
   a.download = filename
   a.click()
   URL.revokeObjectURL(url)
 }
 function exportAnggotaCSV() {
   const headers = [
     { key: 'nomor_anggota', label: 'Nomor Anggota' },
     { key: 'nama', label: 'Nama' },
     { key: 'nik', label: 'NIK' },
     { key: 'status', label: 'Status' },
     { key: 'tanggal_gabung', label: 'Tanggal Gabung' },
   ]
   const csv = toCSV(anggota.value, headers)
   downloadFile('laporan_anggota.csv', csv)
 }
 function exportSimpananCSV() {
   const headers = [
     { key: 'anggota_id', label: 'Anggota ID' },
     { key: 'jenis', label: 'Jenis' },
     { key: 'tipe', label: 'Tipe' },
     { key: 'tanggal', label: 'Tanggal' },
     { key: 'jumlah', label: 'Jumlah' },
     { key: 'saldo_akhir', label: 'Saldo Akhir' },
   ]
   const csv = toCSV(simpanan.value, headers)
   downloadFile('laporan_simpanan.csv', csv)
 }
 function exportPinjamanCSV() {
   const headers = [
     { key: 'nomor_pinjaman', label: 'Nomor Pinjaman' },
     { key: 'anggota_id', label: 'Anggota ID' },
     { key: 'tanggal_pengajuan', label: 'Tanggal Pengajuan' },
     { key: 'nominal', label: 'Nominal' },
     { key: 'tenor_bulan', label: 'Tenor (bulan)' },
     { key: 'bunga_persen', label: 'Bunga (%)' },
     { key: 'status', label: 'Status' },
   ]
   const csv = toCSV(pinjaman.value, headers)
   downloadFile('laporan_pinjaman.csv', csv)
 }
 function exportAngsuranCSV() {
   const headers = [
     { key: 'pinjaman_id', label: 'Pinjaman ID' },
     { key: 'ke', label: 'Ke' },
     { key: 'tanggal_jatuh_tempo', label: 'Jatuh Tempo' },
     { key: 'jumlah', label: 'Jumlah' },
     { key: 'tanggal_bayar', label: 'Tanggal Bayar' },
     { key: 'denda', label: 'Denda' },
   ]
   const csv = toCSV(angsuranFiltered.value, headers)
   downloadFile('laporan_angsuran.csv', csv)
 }
 function exportKasCSV() {
   const headers = [
     { key: 'tanggal', label: 'Tanggal' },
     { key: 'jenis', label: 'Jenis' },
     { key: 'keterangan', label: 'Keterangan' },
     { key: 'jumlah', label: 'Jumlah' },
     { key: 'ref', label: 'Ref' },
   ]
   const csv = toCSV(kas.value, headers)
   downloadFile('laporan_kas.csv', csv)
 }
 function printPDF() { window.print() }

 onMounted(() => {
   computeDefaultRange()
   fetchAnggota()
   fetchSimpanan()
   fetchPinjaman()
   fetchAngsuran()
   computeKasRange()
   fetchKas()
 })
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Laporan</h1>

    <!-- Tabs -->
    <div class="actions-row" style="margin-bottom: 12px; gap: 8px; flex-wrap: wrap;">
      <button class="btn" :class="{ 'btn-primary': tab==='anggota' }" @click="tab='anggota'">Anggota</button>
      <button class="btn" :class="{ 'btn-primary': tab==='simpanan' }" @click="tab='simpanan'">Simpanan</button>
      <button class="btn" :class="{ 'btn-primary': tab==='pinjaman' }" @click="tab='pinjaman'">Pinjaman</button>
      <button class="btn" :class="{ 'btn-primary': tab==='angsuran' }" @click="tab='angsuran'">Angsuran</button>
      <button class="btn" :class="{ 'btn-primary': tab==='kas' }" @click="tab='kas'">Kas</button>
      <span class="spacer"></span>
      <button class="btn btn-secondary" @click="printPDF">Cetak PDF</button>
    </div>

    <!-- Anggota -->
    <div v-if="tab==='anggota'" class="card">
      <div class="card-header">Laporan Anggota</div>
      <div class="card-content">
        <div class="actions-row" style="margin-bottom: 8px;">
          <select class="input" v-model="anggotaStatus" style="max-width: 200px;">
            <option value="">Semua Status</option>
            <option value="pending">Pending</option>
            <option value="verified">Verified</option>
            <option value="active">Active</option>
          </select>
          <button class="btn btn-secondary" @click="fetchAnggota">Terapkan</button>
          <button class="btn" @click="exportAnggotaCSV">Export CSV</button>
          <span class="muted" v-if="anggotaLoading">Memuat…</span>
          <span class="muted" v-if="anggotaError">{{ anggotaError }}</span>
        </div>
        <table class="table">
          <thead>
            <tr>
              <th>Nomor</th>
              <th>Nama</th>
              <th>NIK</th>
              <th>Status</th>
              <th>Gabung</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in anggota" :key="a.id">
              <td>{{ a.nomor_anggota }}</td>
              <td>{{ a.nama }}</td>
              <td>{{ a.nik || '-' }}</td>
              <td><span class="status-badge" :class="a.status">{{ a.status }}</span></td>
              <td>{{ formatDate(a.tanggal_gabung) }}</td>
            </tr>
            <tr v-if="!anggota.length">
              <td colspan="5" class="muted">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Simpanan -->
    <div v-if="tab==='simpanan'" class="card">
      <div class="card-header">Laporan Simpanan</div>
      <div class="card-content">
        <div class="actions-row" style="margin-bottom: 8px; gap: 8px; flex-wrap: wrap;">
          <select class="input" v-model="simpananPeriode" style="max-width: 180px;">
            <option value="harian">Harian</option>
            <option value="mingguan">Mingguan</option>
            <option value="bulanan">Bulanan</option>
          </select>
          <input class="input" v-model="simpananStart" type="date" />
          <input class="input" v-model="simpananEnd" type="date" />
          <select class="input" v-model="simpananJenis" style="max-width: 180px;">
            <option value="">Semua Jenis</option>
            <option value="wajib">Wajib</option>
            <option value="sukarela">Sukarela</option>
            <option value="khusus">Khusus</option>
          </select>
          <button class="btn btn-secondary" @click="fetchSimpanan">Terapkan</button>
          <button class="btn" @click="exportSimpananCSV">Export CSV</button>
          <span class="muted" v-if="simpananLoading">Memuat…</span>
          <span class="muted" v-if="simpananError">{{ simpananError }}</span>
        </div>
        <table class="table">
          <thead>
            <tr>
              <th>Anggota ID</th>
              <th>Jenis</th>
              <th>Tipe</th>
              <th>Tanggal</th>
              <th>Jumlah</th>
              <th>Saldo Akhir</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in simpanan" :key="s.id">
              <td>{{ s.anggota_id }}</td>
              <td>{{ s.jenis }}</td>
              <td>{{ s.tipe }}</td>
              <td>{{ formatDate(s.tanggal) }}</td>
              <td>{{ formatCurrency(s.jumlah) }}</td>
              <td>{{ typeof s.saldo_akhir === 'number' ? formatCurrency(s.saldo_akhir) : '-' }}</td>
            </tr>
            <tr v-if="!simpanan.length">
              <td colspan="6" class="muted">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pinjaman -->
    <div v-if="tab==='pinjaman'" class="card">
      <div class="card-header">Laporan Pinjaman</div>
      <div class="card-content">
        <div class="actions-row" style="margin-bottom: 8px;">
          <select class="input" v-model="pinjamanStatus" style="max-width: 200px;">
            <option value="">Semua Status</option>
            <option value="pengajuan">Pengajuan</option>
            <option value="disetujui">Disetujui</option>
            <option value="berjalan">Berjalan</option>
            <option value="lunas">Lunas</option>
          </select>
          <button class="btn btn-secondary" @click="fetchPinjaman">Terapkan</button>
          <button class="btn" @click="exportPinjamanCSV">Export CSV</button>
          <span class="muted" v-if="pinjamanLoading">Memuat…</span>
          <span class="muted" v-if="pinjamanError">{{ pinjamanError }}</span>
        </div>
        <table class="table">
          <thead>
            <tr>
              <th>Nomor Pinjaman</th>
              <th>Anggota ID</th>
              <th>Tanggal Pengajuan</th>
              <th>Nominal</th>
              <th>Tenor</th>
              <th>Bunga (%)</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in pinjaman" :key="p.id">
              <td>{{ p.nomor_pinjaman }}</td>
              <td>{{ p.anggota_id }}</td>
              <td>{{ formatDate(p.tanggal_pengajuan) }}</td>
              <td>{{ formatCurrency(p.nominal) }}</td>
              <td>{{ p.tenor_bulan }}</td>
              <td>{{ p.bunga_persen }}</td>
              <td><span class="status-badge" :class="p.status">{{ p.status }}</span></td>
            </tr>
            <tr v-if="!pinjaman.length">
              <td colspan="7" class="muted">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Angsuran -->
    <div v-if="tab==='angsuran'" class="card">
      <div class="card-header">Laporan Angsuran</div>
      <div class="card-content">
        <div class="actions-row" style="margin-bottom: 8px;">
          <select class="input" v-model="angsuranMode" style="max-width: 220px;">
            <option value="semua">Semua</option>
            <option value="bayar">Sudah Bayar</option>
            <option value="tunggakan">Tunggakan (Belum Bayar & Jatuh Tempo Lewat)</option>
            <option value="denda">Kena Denda</option>
          </select>
          <button class="btn btn-secondary" @click="fetchAngsuran">Muat Ulang</button>
          <button class="btn" @click="exportAngsuranCSV">Export CSV</button>
          <span class="muted" v-if="angsuranLoading">Memuat…</span>
          <span class="muted" v-if="angsuranError">{{ angsuranError }}</span>
        </div>
        <table class="table">
          <thead>
            <tr>
              <th>Pinjaman ID</th>
              <th>Ke</th>
              <th>Jatuh Tempo</th>
              <th>Jumlah</th>
              <th>Tanggal Bayar</th>
              <th>Denda</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in angsuranFiltered" :key="a.id">
              <td>{{ a.pinjaman_id }}</td>
              <td>{{ a.ke }}</td>
              <td>{{ formatDate(a.tanggal_jatuh_tempo) }}</td>
              <td>{{ formatCurrency(a.jumlah) }}</td>
              <td>{{ a.tanggal_bayar ? formatDate(a.tanggal_bayar) : '-' }}</td>
              <td>{{ formatCurrency(a.denda || 0) }}</td>
            </tr>
            <tr v-if="!angsuranFiltered.length">
              <td colspan="6" class="muted">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Kas -->
    <div v-if="tab==='kas'" class="card">
      <div class="card-header">Laporan Kas (Penerimaan/Pengeluaran)</div>
      <div class="card-content">
        <div class="actions-row" style="margin-bottom: 8px; gap: 8px; flex-wrap: wrap;">
          <select class="input" v-model="kasPeriode" style="max-width: 180px;">
            <option value="harian">Harian</option>
            <option value="mingguan">Mingguan</option>
            <option value="bulanan">Bulanan</option>
          </select>
          <input class="input" v-model="kasStart" type="date" />
          <input class="input" v-model="kasEnd" type="date" />
          <button class="btn btn-secondary" @click="fetchKas">Terapkan</button>
          <button class="btn" @click="exportKasCSV" :disabled="!!kasError">Export CSV</button>
          <span class="muted" v-if="kasLoading">Memuat…</span>
          <span class="muted" v-if="kasError">{{ kasError }}</span>
        </div>
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
            <tr v-for="k in kas" :key="(k as any).id">
              <td>{{ formatDate(k.tanggal) }}</td>
              <td>{{ k.jenis }}</td>
              <td>{{ k.keterangan || '-' }}</td>
              <td>{{ formatCurrency(k.jumlah) }}</td>
              <td>{{ k.ref || '-' }}</td>
            </tr>
            <tr v-if="!kas.length">
              <td colspan="5" class="muted">{{ kasError ? 'Endpoint tidak tersedia' : 'Tidak ada data' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.spacer { flex: 1; }
.btn-primary { background: var(--color-primary); color: white; }
.btn { border: 1px solid var(--color-border); padding: 8px 12px; border-radius: var(--radius-sm); background: var(--color-surface); cursor: pointer; }
.btn:hover { background: #f8fafc; }
.input { border: 1px solid var(--color-border); border-radius: var(--radius-sm); padding: 8px; }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { border: 1px solid var(--color-border); padding: 8px; text-align: left; }
.card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); box-shadow: var(--shadow-sm); margin-top: 12px; }
.card-header { padding: 12px; border-bottom: 1px solid var(--color-border); font-weight: 600; }
.card-content { padding: 12px; }
.actions-row { display: flex; align-items: center; gap: 10px; }
.page-title { margin: 0 0 8px 0; }
.status-badge { display: inline-block; padding: 4px 8px; border-radius: 999px; border: 1px solid var(--color-border); font-size: 12px; }
</style>