<script setup lang="ts">
import { onMounted } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'

const dashboard = useDashboardStore()

onMounted(() => {
  dashboard.fetchOverview()
})
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Dashboard</h1>

    <!-- Metrics cards -->
    <div class="metrics">
      <div class="card">
        <div class="card-header">Jumlah Anggota</div>
        <div class="card-content text-3xl font-bold">
          {{ dashboard.overview.anggotaCount ?? '-' }}
        </div>
        <div class="px-4 pb-4 text-sm text-gray-500">Total anggota aktif</div>
      </div>

      <div class="card">
        <div class="card-header">Saldo Simpanan</div>
        <div class="card-content text-3xl font-bold">
          {{ formatCurrency(dashboard.overview.simpananTotal) }}
        </div>
        <div class="px-4 pb-4 text-sm text-gray-500">Akumulasi simpanan</div>
      </div>

      <div class="card">
        <div class="card-header">Pinjaman Aktif</div>
        <div class="card-content text-3xl font-bold">
          {{ dashboard.overview.pinjamanAktif ?? '-' }}
        </div>
        <div class="px-4 pb-4 text-sm text-gray-500">Jumlah pinjaman berjalan</div>
      </div>

      <div class="card">
        <div class="card-header">Kas Harian</div>
        <div class="card-content text-3xl font-bold">
          {{ formatCurrency(dashboard.overview.kasHarian) }}
        </div>
        <div class="px-4 pb-4 text-sm text-gray-500">Net kas hari ini</div>
      </div>
    </div>

    <!-- Tasks / Notifications -->
    <div class="card" style="margin-top: 24px;">
      <div class="card-header">Notifikasi & Tugas</div>
      <div class="card-content">
        <ul>
          <li v-for="(t, i) in dashboard.overview.tasks" :key="i" style="display:flex; align-items:flex-start; justify-content:space-between; margin-bottom:12px;">
            <div>
              <p style="font-weight:600;">{{ t.title }}</p>
              <p class="muted" style="font-size: 0.9rem;">{{ t.description }}</p>
            </div>
            <span class="badge">{{ t.badge }}</span>
          </li>
        </ul>
        <div v-if="!dashboard.overview.tasks?.length" class="muted">Tidak ada notifikasi</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  methods: {
    formatCurrency(value?: number) {
      if (value == null) return '-'
      return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
    },
  },
}
</script>