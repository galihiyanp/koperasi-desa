<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { useRouter } from 'vue-router'
import { useNotificationsStore } from '@/stores/notifications'

const router = useRouter()
const dashboard = useDashboardStore()

const loading = ref(false)
const error = ref<string | null>(null)

const notifications = useNotificationsStore()

async function refreshOverview() {
  loading.value = true
  error.value = null
  notifications.notify({
    type: 'info',
    title: 'Memuat ringkasan',
    message: 'Mengambil data terbaru dashboard...',
    detail: 'Harap tunggu sebentar',
    timeout: 3000,
  })
  try {
    await dashboard.fetchOverview()
    notifications.notify({
      type: 'success',
      title: 'Ringkasan diperbarui',
      message: 'Data dashboard berhasil dimuat',
      timeout: 3500,
    })
  } catch (e: any) {
    error.value = e?.message || 'Terjadi kesalahan'
    notifications.notify({
      type: 'error',
      title: 'Gagal memuat',
      message: error.value ?? 'Terjadi kesalahan',
      detail: 'Coba muat ulang atau periksa koneksi',
      timeout: 6000,
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshOverview()
})

function formatCurrency(value?: number) {
  if (value == null) return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

const kasProgress = computed(() => {
  const v = dashboard.overview.kasHarian ?? 0
  const target = 5_000_000
  const pct = Math.max(0, Math.min(100, Math.round((v / target) * 100)))
  return pct
})
</script>

<template>
  <div class="dashboard">
    <!-- Header & Quick Actions -->
    <div class="header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="muted">Ringkasan koperasi dan tindakan cepat</p>
      </div>
      <div class="header-actions">
        <button class="btn" @click="refreshOverview" :disabled="loading">{{ loading ? 'Memuat…' : 'Refresh' }}</button>
        <div class="quick-actions">
          <button class="btn" @click="router.push('/keanggotaan')">Keanggotaan</button>
          <button class="btn" @click="router.push('/pinjaman')">Pinjaman</button>
          <button class="btn" @click="router.push('/angsuran')">Angsuran</button>
          <button class="btn" @click="router.push('/laporan')">Laporan</button>
        </div>
      </div>
    </div>

    <div v-if="error" class="card error"><div class="card-content">{{ error }}</div></div>

    <!-- Metrics Cards -->
    <div class="metrics-grid" v-if="!loading">
      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-title">Jumlah Anggota</span>
          <span class="metric-icon users"></span>
        </div>
        <div class="metric-value">{{ dashboard.overview.anggotaCount ?? '-' }}</div>
        <div class="metric-sub">Total anggota aktif</div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-title">Saldo Simpanan</span>
          <span class="metric-icon savings"></span>
        </div>
        <div class="metric-value">{{ formatCurrency(dashboard.overview.simpananTotal) }}</div>
        <div class="metric-sub">Akumulasi simpanan</div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-title">Pinjaman Aktif</span>
          <span class="metric-icon loans"></span>
        </div>
        <div class="metric-value">{{ dashboard.overview.pinjamanAktif ?? '-' }}</div>
        <div class="metric-sub">Jumlah pinjaman berjalan</div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <span class="metric-title">Kas Harian</span>
          <span class="metric-icon cash"></span>
        </div>
        <div class="metric-value">{{ formatCurrency(dashboard.overview.kasHarian) }}</div>
        <div class="metric-sub">Net kas hari ini</div>
      </div>
    </div>

    <!-- Metrics Skeleton -->
    <div class="metrics-grid" v-else>
      <div class="metric-card skeleton" v-for="i in 4" :key="i">
        <div class="skeleton-line" style="width: 60%;"></div>
        <div class="skeleton-line" style="width: 40%; height: 24px; margin-top: 8px;"></div>
        <div class="skeleton-line" style="width: 50%;"></div>
      </div>
    </div>

    <!-- Quick Links -->
    <div class="links-grid">
      <div class="link-card" @click="router.push('/keanggotaan')">
        <div class="link-icon users"></div>
        <div class="link-title">Kelola Anggota</div>
        <div class="link-sub">Daftar, verifikasi, aktivasi</div>
      </div>
      <div class="link-card" @click="router.push('/pinjaman')">
        <div class="link-icon loans"></div>
        <div class="link-title">Pinjaman</div>
        <div class="link-sub">Pengajuan, pencairan, berjalan</div>
      </div>
      <div class="link-card" @click="router.push('/angsuran')">
        <div class="link-icon cash"></div>
        <div class="link-title">Angsuran</div>
        <div class="link-sub">Bayar jatuh tempo</div>
      </div>
      <div class="link-card" @click="router.push('/laporan')">
        <div class="link-icon report"></div>
        <div class="link-title">Laporan</div>
        <div class="link-sub">Anggota, simpanan, pinjaman, kas</div>
      </div>
    </div>

    <!-- Grid: Kas Summary + Notifications -->
    <div class="grid-two">
      <div class="card">
        <div class="card-header">Ringkasan Kas</div>
        <div class="card-content">
          <div class="kas-row">
            <span>Net hari ini</span>
            <strong>{{ formatCurrency(dashboard.overview.kasHarian) }}</strong>
          </div>
          <div class="progress">
            <div class="bar" :style="{ width: kasProgress + '%' }"></div>
          </div>
          <p class="muted">Estimasi pencapaian terhadap target harian</p>
        </div>
      </div>

      <div class="card">
        <div class="card-header">Notifikasi & Tugas</div>
        <div class="card-content">
          <ul class="tasks">
            <li v-for="(t, i) in dashboard.overview.tasks" :key="i">
              <div class="task-left">
                <p class="task-title">{{ t.title }}</p>
                <p class="task-desc">{{ t.description }}</p>
              </div>
              <span class="badge">{{ t.badge }}</span>
            </li>
          </ul>
          <div v-if="!dashboard.overview.tasks?.length" class="muted">Tidak ada notifikasi</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-title { margin: 0; font-size: 22px; font-weight: 700; }
.muted { color: #6b7280; }

.header { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; margin-bottom: 16px; position: sticky; top: 0; z-index: 20; background: var(--color-surface); border-bottom: 1px solid var(--color-border); box-shadow: var(--shadow-sm); }
.header-actions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.quick-actions { display: flex; gap: 8px; flex-wrap: wrap; }

.btn { border: 1px solid var(--color-border); padding: 8px 12px; border-radius: var(--radius-sm); background: var(--color-surface); cursor: pointer; }
.btn:hover { background: #f8fafc; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }

.metrics-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.metric-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); box-shadow: var(--shadow-sm); padding: 12px; }
.metric-header { display: flex; align-items: center; justify-content: space-between; }
.metric-title { font-weight: 600; }
.metric-icon { width: 28px; height: 28px; border-radius: 50%; background: #eef2ff; border: 1px solid #dbeafe; }
.metric-icon.users { background-image: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%); }
.metric-icon.savings { background-image: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%); }
.metric-icon.loans { background-image: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%); }
.metric-icon.cash { background-image: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%); }
.metric-value { font-size: 26px; font-weight: 700; margin-top: 8px; }
.metric-sub { font-size: 12px; color: #6b7280; margin-top: 2px; }

.skeleton { position: relative; overflow: hidden; }
.skeleton::after { content: ''; position: absolute; inset: 0; background: linear-gradient(90deg, rgba(0,0,0,0.03) 0%, rgba(0,0,0,0.06) 50%, rgba(0,0,0,0.03) 100%); animation: shimmer 1.2s infinite; }
.skeleton-line { height: 10px; background: #eef2ff; border-radius: 6px; }
@keyframes shimmer { 0% { transform: translateX(-100%);} 100% { transform: translateX(100%);} }

.links-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-top: 16px; }
.link-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); box-shadow: var(--shadow-sm); padding: 12px; cursor: pointer; }
.link-card:hover { background: #f8fafc; }
.link-icon { width: 36px; height: 36px; border-radius: 8px; border: 1px solid var(--color-border); background: #eef2ff; margin-bottom: 6px; }
.link-icon.users { background-image: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%); }
.link-icon.loans { background-image: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%); }
.link-icon.cash { background-image: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%); }
.link-icon.report { background-image: linear-gradient(135deg, #fde68a 0%, #fcd34d 100%); }
.link-title { font-weight: 600; }
.link-sub { font-size: 12px; color: #6b7280; }

.grid-two { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 16px; }
.card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); box-shadow: var(--shadow-sm); }
.card-header { padding: 12px; border-bottom: 1px solid var(--color-border); font-weight: 600; }
.card-content { padding: 12px; }
.card.error { border-color: #e33; color: #b00; }

.kas-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.progress { height: 8px; background: #eef2ff; border-radius: 999px; overflow: hidden; border: 1px solid #e5e7eb; }
.bar { height: 100%; background: #0ea5e9; width: 0%; transition: width 0.3s ease; }

.tasks { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; }
.tasks li { display: flex; align-items: flex-start; justify-content: space-between; gap: 10px; }
.task-title { font-weight: 600; }
.task-desc { font-size: 13px; color: #6b7280; }
.badge { display: inline-block; padding: 4px 8px; border-radius: 999px; border: 1px solid var(--color-border); font-size: 12px; }

@media (max-width: 1024px) {
  .metrics-grid { grid-template-columns: repeat(2, 1fr); }
  .links-grid { grid-template-columns: repeat(2, 1fr); }
  .grid-two { grid-template-columns: 1fr; }
}

@media (max-width: 640px) {
  .metrics-grid, .links-grid { grid-template-columns: 1fr; }
}
</style>