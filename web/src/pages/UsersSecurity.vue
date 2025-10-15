<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/lib/axios'

type User = {
  id: number
  name: string
  email: string
  role?: string
  created_at?: string
  updated_at?: string
}

const users = ref<User[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const page = ref(1)
const limit = ref(10)
const hasNext = ref(false)

const showForm = ref(false)
const mode = ref<'create' | 'edit'>('create')
const form = reactive<{ id?: number; name: string; email: string; password: string; role: string }>({
  id: undefined,
  name: '',
  email: '',
  password: '',
  role: 'user',
})

async function fetchUsers() {
  loading.value = true
  error.value = null
  try {
    const res = await api.get('/api/users', { params: { page: page.value, limit: limit.value } })
    users.value = res.data?.data ?? []
    hasNext.value = (users.value?.length ?? 0) >= limit.value
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal memuat data pengguna'
  } finally {
    loading.value = false
  }
}

function openCreate() {
  mode.value = 'create'
  showForm.value = true
  form.id = undefined
  form.name = ''
  form.email = ''
  form.password = ''
  form.role = 'user'
}

function openEdit(u: User) {
  mode.value = 'edit'
  showForm.value = true
  form.id = u.id
  form.name = u.name
  form.email = u.email
  form.password = ''
  form.role = u.role || 'user'
}

function cancelForm() {
  showForm.value = false
}

async function submitForm() {
  try {
    if (mode.value === 'create') {
      const payload: Record<string, any> = { name: form.name, email: form.email, role: form.role, password: form.password }
      await api.post('/api/users', payload)
    } else if (mode.value === 'edit' && form.id != null) {
      const payload: Record<string, any> = { name: form.name, email: form.email, role: form.role }
      if (form.password && form.password.length >= 6) payload.password = form.password
      await api.put(`/api/users/${form.id}`, payload)
    }
    await fetchUsers()
    showForm.value = false
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Operasi gagal')
  }
}

async function deleteUser(u: User) {
  if (!confirm(`Hapus pengguna ${u.email}?`)) return
  try {
    await api.delete(`/api/users/${u.id}`)
    await fetchUsers()
  } catch (e: any) {
    alert(e?.response?.data?.error || e?.message || 'Gagal menghapus pengguna')
  }
}

function prevPage() {
  if (page.value > 1) {
    page.value -= 1
    fetchUsers()
  }
}

function nextPage() {
  if (hasNext.value) {
    page.value += 1
    fetchUsers()
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="dashboard">
    <h1 class="page-title">Pengguna & Keamanan</h1>

    <div class="actions-row">
      <button class="btn btn-primary" @click="openCreate">Tambah Pengguna</button>
    </div>

    <div v-if="showForm" class="card" style="margin-top: 16px;">
      <div class="card-header">{{ mode === 'create' ? 'Tambah Pengguna' : 'Ubah Pengguna' }}</div>
      <div class="card-content">
        <div class="form-row">
          <label class="label">Nama</label>
          <input class="input" v-model="form.name" type="text" placeholder="Nama lengkap" />
        </div>
        <div class="form-row">
          <label class="label">Email</label>
          <input class="input" v-model="form.email" type="email" placeholder="email@example.com" />
        </div>
        <div class="form-row">
          <label class="label">Role</label>
          <select class="input" v-model="form.role">
            <option value="user">User</option>
            <option value="admin">Admin</option>
            <option value="superadmin">Super Admin</option>
          </select>
        </div>
        <div class="form-row">
          <label class="label">Password</label>
          <input class="input" v-model="form.password" type="password" :placeholder="mode === 'create' ? 'Minimal 6 karakter' : 'Isi untuk mengganti password'" />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" @click="submitForm">Simpan</button>
          <button class="btn btn-secondary" @click="cancelForm">Batal</button>
        </div>
      </div>
    </div>

    <div class="card" style="margin-top: 16px;">
      <div class="card-header">Daftar Pengguna</div>
      <div class="card-content">
        <div v-if="error" class="muted" style="margin-bottom: 12px;">{{ error }}</div>
        <div v-if="loading" class="muted">Memuat...</div>
        <table v-else class="table">
          <thead>
            <tr>
              <th>Nama</th>
              <th>Email</th>
              <th>Role</th>
              <th>Dibuat</th>
              <th style="width: 160px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td>{{ u.name }}</td>
              <td>{{ u.email }}</td>
              <td>{{ u.role || '-' }}</td>
              <td>{{ u.created_at || '-' }}</td>
              <td>
                <div class="table-actions">
                  <button class="btn btn-secondary" @click="openEdit(u)">Ubah</button>
                  <button class="btn btn-danger" @click="deleteUser(u)">Hapus</button>
                </div>
              </td>
            </tr>
            <tr v-if="!users.length">
              <td colspan="5" class="muted">Tidak ada data</td>
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
  </div>
</template>

<style scoped>
.page-title { font-size: 22px; font-weight: 600; margin-bottom: 12px; }
.actions-row { display: flex; gap: 8px; }
.form-row { display: grid; grid-template-columns: 160px 1fr; gap: 8px; align-items: center; margin-bottom: 10px; }
.label { font-size: 14px; color: #374151; }
.input { padding: 8px 10px; border-radius: 8px; border: 1px solid #e5e7eb; font-size: 14px; }
.form-actions { display: flex; gap: 8px; margin-top: 6px; }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { text-align: left; padding: 10px 12px; border-bottom: 1px solid #f3f4f6; }
.table thead th { font-weight: 600; color: #374151; background: #f9fafb; }
.table-actions { display: flex; gap: 8px; }
.pagination { display: flex; gap: 12px; align-items: center; margin-top: 12px; }
.btn { padding: 6px 12px; border-radius: 6px; font-size: 13px; border: 1px solid #e5e7eb; cursor: pointer; }
.btn-primary { background: #1d4ed8; color: #fff; border-color: #1d4ed8; }
.btn-primary:hover { background: #1e40af; border-color: #1e40af; }
.btn-secondary { background: #f3f4f6; }
.btn-secondary:hover { background: #e5e7eb; }
.btn-danger { background: #ef4444; color: #fff; border-color: #ef4444; }
.btn-danger:hover { background: #dc2626; border-color: #dc2626; }
</style>