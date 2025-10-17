<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import api from '@/lib/axios'

// Ambil kontak dukungan dari pengaturan profile dan integrasi
const KEY_PROFILE = 'settings.profile'
const KEY_INTEGRATIONS = 'settings.integrations'

const loading = ref(false)
const error = ref<string | null>(null)

const profile = reactive<{ name?: string; address?: string; phone?: string; email?: string }>({})
const integrations = reactive<{ whatsapp_enabled?: boolean; email_enabled?: boolean; whatsapp_number?: string; email_smtp?: string }>({})

async function getSetting(key: string) {
  try {
    const res = await api.get(`/api/settings/${encodeURIComponent(key)}`)
    return res.data?.value as string | undefined
  } catch (e) {
    return undefined
  }
}
async function loadSettings() {
  loading.value = true
  error.value = null
  try {
    const p = await getSetting(KEY_PROFILE)
    if (p) { try { Object.assign(profile, JSON.parse(p)) } catch { /* noop */ } }
    const ig = await getSetting(KEY_INTEGRATIONS)
    if (ig) { try { Object.assign(integrations, JSON.parse(ig)) } catch { /* noop */ } }
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal memuat pengaturan bantuan'
  } finally {
    loading.value = false
  }
}

onMounted(loadSettings)

// FAQ & Panduan (ringkas, sesuai README Menu & Fitur)
const query = ref('')
const faq = ref<Array<{ q: string; a: string }>>([
  { q: 'Bagaimana cara membuat pengajuan pinjaman?', a: 'Masuk ke menu Pinjaman, klik Pengajuan, isi data (anggota, nominal, tenor, bunga), lalu kirim. Admin verifikasi sebelum pencairan.' },
  { q: 'Bagaimana cara melakukan pembayaran angsuran?', a: 'Buka menu Angsuran, pilih pinjaman terkait, klik Bayar pada angsuran jatuh tempo, masukkan jumlah, simpan.' },
  { q: 'Bagaimana cara setoran simpanan?', a: 'Masuk ke menu Simpanan, pilih Setoran, pilih jenis (wajib/sukarela/khusus), input nominal, dan Simpan.' },
  { q: 'Di mana melihat laporan transaksi?', a: 'Gunakan menu Laporan untuk melihat data Anggota, Simpanan, Pinjaman, dan Angsuran. Tersedia ekspor CSV/PDF.' },
  { q: 'Bagaimana mengatur profil koperasi?', a: 'Masuk ke menu Pengaturan, isi nama/alamat/kontak, parameter finansial, kategori, dan integrasi. Simpan per bagian.' },
])
const filteredFaq = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return faq.value
  return faq.value.filter(f => f.q.toLowerCase().includes(q) || f.a.toLowerCase().includes(q))
})

// Helper link
const whatsappLink = computed(() => {
  if (!integrations.whatsapp_enabled || !integrations.whatsapp_number) return ''
  const msg = encodeURIComponent('Halo, saya butuh bantuan terkait aplikasi Koperasi Desa.')
  const num = integrations.whatsapp_number.replace(/[^0-9]/g, '')
  return `https://wa.me/${num}?text=${msg}`
})
</script>

<template>
  <div class="content">
    <h1 class="page-title">Bantuan</h1>
    <p class="muted">FAQ & Panduan penggunaan, serta kontak dukungan.</p>

    <div v-if="loading" class="card"><p>Memuat informasi...</p></div>
    <div v-if="error" class="card error"><p>{{ error }}</p></div>

    <!-- Pencarian FAQ -->
    <div class="card">
      <h2>FAQ</h2>
      <input class="search" v-model="query" placeholder="Cari topik (mis. pinjaman, angsuran)" />
      <ul class="faq">
        <li v-for="(f, i) in filteredFaq" :key="i">
          <strong>{{ f.q }}</strong>
          <p>{{ f.a }}</p>
        </li>
      </ul>
    </div>

    <!-- Panduan Ringkas -->
    <div class="card">
      <h2>Panduan Ringkas</h2>
      <ul class="guide">
        <li>Keanggotaan: Daftarkan anggota, verifikasi, dan aktivasi. Unggah dokumen bila perlu.</li>
        <li>Simpanan: Atur kategori di Pengaturan. Melakukan setoran/penarikan via menu terkait.</li>
        <li>Pinjaman: Pengajuan, verifikasi, pencairan. Jadwal angsuran dibuat otomatis saat pencairan.</li>
        <li>Angsuran: Lihat daftar, bayar yang jatuh tempo, pantau denda/tunggakan.</li>
        <li>Laporan: Gunakan tab untuk Anggota/Simpanan/Pinjaman/Angsuran. Ekspor CSV.</li>
        <li>Pengaturan: Profil koperasi, parameter finansial, kategori, format nomor, integrasi.</li>
      </ul>
    </div>

    <!-- Kontak Dukungan -->
    <div class="card">
      <h2>Kontak Dukungan</h2>
      <div class="grid">
        <div>
          <p><strong>Koperasi:</strong> {{ profile.name || '-' }}</p>
          <p><strong>Alamat:</strong> {{ profile.address || '-' }}</p>
          <p><strong>Telepon:</strong> {{ profile.phone || '-' }}</p>
          <p><strong>Email:</strong> {{ profile.email || '-' }}</p>
        </div>
        <div>
          <div v-if="integrations.whatsapp_enabled && integrations.whatsapp_number">
            <a :href="whatsappLink" target="_blank" class="btn primary">Hubungi via WhatsApp</a>
          </div>
          <div v-else class="muted">WhatsApp Gateway belum diaktifkan di Pengaturan.</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.content { padding: 16px; }
.page-title { font-size: 20px; font-weight: 600; }
.muted { color: #666; }
.card { background: #fff; border: 1px solid #eee; border-radius: 8px; padding: 16px; margin-bottom: 16px; }
.card.error { border-color: #e33; color: #b00; }
.search { width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 6px; margin-bottom: 12px; }
.faq { list-style: none; padding: 0; display: flex; flex-direction: column; gap: 10px; }
.faq li strong { display: block; margin-bottom: 4px; }
.guide { list-style: disc; padding-left: 18px; display: flex; flex-direction: column; gap: 6px; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.btn { display: inline-block; padding: 8px 12px; border-radius: 6px; border: 1px solid #ccc; background: #f8f8f8; cursor: pointer; text-decoration: none; text-align: center; }
.btn.primary { background: #0b77ff; color: #fff; border-color: #0b77ff; }
</style>