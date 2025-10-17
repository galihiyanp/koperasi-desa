<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/lib/axios'

// Keys yang digunakan
const KEY_PROFILE = 'settings.profile'
const KEY_FINANCIAL = 'settings.financial'
const KEY_CATEGORIES = 'settings.categories'
const KEY_FORMAT = 'settings.format'
const KEY_INTEGRATIONS = 'settings.integrations'

// State
const loading = ref(false)
const error = ref<string | null>(null)

// Profil koperasi
const profile = reactive<{ name: string; address: string; phone: string; email: string }>({
  name: '', address: '', phone: '', email: ''
})

// Parameter finansial
const financial = reactive<{ suku_bunga_default: number; biaya_admin_default: number }>({
  suku_bunga_default: 1.5,
  biaya_admin_default: 0,
})

// Kategori simpanan & pinjaman
const categories = reactive<{ simpanan: string[]; pinjaman: string[] }>({
  simpanan: ['wajib', 'sukarela', 'khusus'],
  pinjaman: ['produktif', 'konsumtif'],
})
// Tambah input refs untuk kategori baru
const newSimpanan = ref('')
const newPinjaman = ref('')

// Format nomor transaksi
const format = reactive<{ nomor_anggota: string; nomor_pinjaman: string }>({
  nomor_anggota: 'AG-{YYYY}-{seq}',
  nomor_pinjaman: 'PJ-{YYYY}-{seq}',
})

// Integrasi
const integrations = reactive<{ whatsapp_enabled: boolean; email_enabled: boolean; whatsapp_number?: string; email_smtp?: string }>({
  whatsapp_enabled: false,
  email_enabled: false,
  whatsapp_number: '',
  email_smtp: '',
})

async function getSetting(key: string) {
  try {
    const res = await api.get(`/api/settings/${encodeURIComponent(key)}`)
    return res.data?.value as string | undefined
  } catch (e) {
    return undefined
  }
}
async function putSetting(key: string, value: any) {
  await api.put(`/api/settings/${encodeURIComponent(key)}`, { value: typeof value === 'string' ? value : JSON.stringify(value) })
}

async function loadAll() {
  loading.value = true
  error.value = null
  try {
    // Load tiap key
    const p = await getSetting(KEY_PROFILE)
    if (p) {
      try { Object.assign(profile, JSON.parse(p)) } catch { /* string fallback */ }
    }
    const f = await getSetting(KEY_FINANCIAL)
    if (f) {
      try { Object.assign(financial, JSON.parse(f)) } catch { /* ignore */ }
    }
    const c = await getSetting(KEY_CATEGORIES)
    if (c) {
      try { Object.assign(categories, JSON.parse(c)) } catch { /* ignore */ }
    }
    const fm = await getSetting(KEY_FORMAT)
    if (fm) {
      try { Object.assign(format, JSON.parse(fm)) } catch { /* ignore */ }
    }
    const ig = await getSetting(KEY_INTEGRATIONS)
    if (ig) {
      try { Object.assign(integrations, JSON.parse(ig)) } catch { /* ignore */ }
    }
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal memuat pengaturan'
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

// Simpan fungsi per bagian
async function saveProfile() { await putSetting(KEY_PROFILE, profile) }
async function saveFinancial() { await putSetting(KEY_FINANCIAL, financial) }
async function saveCategories() { await putSetting(KEY_CATEGORIES, categories) }
async function saveFormat() { await putSetting(KEY_FORMAT, format) }
async function saveIntegrations() { await putSetting(KEY_INTEGRATIONS, integrations) }

function addItem(list: string[], v: string) {
  const s = v.trim(); if (!s) return; if (!list.includes(s)) list.push(s)
}
function removeItem(list: string[], i: number) { list.splice(i, 1) }
</script>

<template>
  <div class="content">
    <h1 class="page-title">Pengaturan</h1>
    <p class="muted">Kelola profil koperasi, parameter finansial, kategori, format nomor, dan integrasi.</p>

    <div v-if="loading" class="card"><p>Memuat pengaturan...</p></div>
    <div v-if="error" class="card error"><p>{{ error }}</p></div>

    <!-- Profil Koperasi -->
    <div class="card">
      <h2>Profil Koperasi</h2>
      <div class="grid">
        <label>Nama
          <input v-model="profile.name" placeholder="Nama Koperasi" />
        </label>
        <label>Alamat
          <input v-model="profile.address" placeholder="Alamat" />
        </label>
        <label>Telepon
          <input v-model="profile.phone" placeholder="08xxxxxxxx" />
        </label>
        <label>Email
          <input v-model="profile.email" placeholder="email@domain.com" />
        </label>
      </div>
      <button class="btn" @click="saveProfile">Simpan Profil</button>
    </div>

    <!-- Parameter Finansial -->
    <div class="card">
      <h2>Parameter Finansial</h2>
      <div class="grid">
        <label>Suku Bunga Default (%)
          <input type="number" step="0.01" v-model.number="financial.suku_bunga_default" />
        </label>
        <label>Biaya Administrasi Default (Rp)
          <input type="number" step="0.01" v-model.number="financial.biaya_admin_default" />
        </label>
      </div>
      <button class="btn" @click="saveFinancial">Simpan Parameter</button>
    </div>

    <!-- Kategori Simpanan & Pinjaman -->
    <div class="card">
      <h2>Kategori Simpanan & Pinjaman</h2>
      <div class="grid two">
        <div>
          <h3>Simpanan</h3>
          <ul class="pill-list">
            <li v-for="(j,i) in categories.simpanan" :key="i">
              {{ j }} <button class="link danger" @click="removeItem(categories.simpanan, i)">hapus</button>
            </li>
          </ul>
          <div class="inline">
            <input id="newSimpanan" placeholder="Tambah kategori" v-model="newSimpanan" />
            <button class="btn" @click="() => { addItem(categories.simpanan, newSimpanan); newSimpanan=''; }">Tambah</button>
          </div>
        </div>
        <div>
          <h3>Pinjaman</h3>
          <ul class="pill-list">
            <li v-for="(j,i) in categories.pinjaman" :key="i">
              {{ j }} <button class="link danger" @click="removeItem(categories.pinjaman, i)">hapus</button>
            </li>
          </ul>
          <div class="inline">
            <input id="newPinjaman" placeholder="Tambah kategori" v-model="newPinjaman" />
            <button class="btn" @click="() => { addItem(categories.pinjaman, newPinjaman); newPinjaman=''; }">Tambah</button>
          </div>
        </div>
      </div>
      <button class="btn" @click="saveCategories">Simpan Kategori</button>
    </div>

    <!-- Format Nomor Transaksi -->
    <div class="card">
      <h2>Format Nomor Transaksi</h2>
      <p class="muted">Gunakan token seperti {YYYY}, {MM}, {DD}, {seq}. Contoh: PJ-{YYYY}-{seq}</p>
      <div class="grid">
        <label>Format Nomor Anggota
          <input v-model="format.nomor_anggota" />
        </label>
        <label>Format Nomor Pinjaman
          <input v-model="format.nomor_pinjaman" />
        </label>
      </div>
      <button class="btn" @click="saveFormat">Simpan Format</button>
    </div>

    <!-- Integrasi -->
    <div class="card">
      <h2>Integrasi</h2>
      <div class="grid">
        <label>
          <input type="checkbox" v-model="integrations.whatsapp_enabled" /> Aktifkan WhatsApp Gateway
        </label>
        <label>
          Nomor WhatsApp
          <input v-model="integrations.whatsapp_number" placeholder="62xxxxxxxxxx" />
        </label>
        <label>
          <input type="checkbox" v-model="integrations.email_enabled" /> Aktifkan Email Gateway
        </label>
        <label>
          Konfigurasi SMTP
          <input v-model="integrations.email_smtp" placeholder="smtp://user:pass@host:port" />
        </label>
      </div>
      <button class="btn" @click="saveIntegrations">Simpan Integrasi</button>
    </div>
  </div>
</template>

<style scoped>
.content { padding: 16px; }
.page-title { font-size: 20px; font-weight: 600; }
.muted { color: #666; }
.card { background: #fff; border: 1px solid #eee; border-radius: 8px; padding: 16px; margin-bottom: 16px; }
.card.error { border-color: #e33; color: #b00; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.grid.two { grid-template-columns: 1fr 1fr; }
label { display: flex; flex-direction: column; gap: 6px; font-size: 14px; }
input, select { border: 1px solid #ccc; border-radius: 6px; padding: 8px; font-size: 14px; }
.btn { margin-top: 8px; padding: 8px 12px; border-radius: 6px; border: 1px solid #ccc; background: #f8f8f8; cursor: pointer; }
.inline { display: flex; gap: 8px; align-items: center; }
.pill-list { display: flex; flex-wrap: wrap; gap: 8px; list-style: none; padding: 0; }
.pill-list li { background: #f3f3f3; border: 1px solid #ddd; border-radius: 999px; padding: 6px 10px; }
.link { background: none; border: none; color: #06c; cursor: pointer; text-decoration: underline; }
.link.danger { color: #c00; }
</style>