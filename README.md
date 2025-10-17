# Koperasi Desa – Sistem Koperasi Kantor Desa

Dokumentasi lengkap pengembangan Sistem Koperasi Kantor Desa yang mencakup arsitektur, fitur, menu, alur bisnis, setup lingkungan, dan panduan operasional. Sistem ini dirancang modular agar mudah dikembangkan dan di-deploy.

## Ringkasan
- Sistem mendukung dua peran utama: `Anggota` dan `Pegawai/Admin`.
- Fitur inti: keanggotaan, simpanan, penarikan, pinjaman, angsuran, kas & jurnal sederhana, laporan, dan pengaturan.
- Frontend berbasis `Vue 3 + Vite + Pinia + shadcn-vue`.
- Backend berbasis `Go 1.22 + Gin + GORM + MySQL`.

## Arsitektur & Tech Stack
- Frontend: `Vue 3`, `Vite`, `Pinia`, `shadcn-vue`, `Axios`
- Backend: `Go 1.22`, `Gin`, `GORM`
- Database: `MySQL`
- Autentikasi: `JWT` (Bearer Token)
- Build & Dev: `Node.js >= 18`, `Go >= 1.22`
- Opsional: `Docker Compose` untuk orkestrasi layanan

## Struktur Proyek
```
.
├── README.md            # Dokumentasi proyek
├── web/                 # Frontend (Vue 3 + Vite)
└── service/             # Backend (Go + Gin + GORM)
```

## Peran Pengguna
- Anggota
  - Mengelola profil
  - Melihat saldo dan riwayat simpanan
  - Mengajukan pinjaman dan melihat statusnya
  - Melihat jadwal angsuran dan melakukan pembayaran
- Pegawai/Admin
  - Manajemen anggota (pendaftaran, verifikasi, aktivasi)
  - Pengelolaan simpanan, penarikan, pinjaman, angsuran
  - Pencatatan kas (penerimaan/pengeluaran) & jurnal sederhana
  - Manajemen pengguna, pengaturan sistem, dan pembuatan laporan

## Menu & Fitur Lengkap
- Dashboard
  - Ringkasan jumlah anggota, simpanan, pinjaman aktif, kas harian
  - Notifikasi tugas (pengajuan menunggu persetujuan, jatuh tempo angsuran)
- Keanggotaan
  - Pendaftaran Anggota (input data, upload dokumen)
  - Verifikasi & Aktivasi Anggota
  - Profil Anggota & Riwayat Aktivitas
- Simpanan
  - Setoran (Wajib, Sukarela, Khusus)
  - Penarikan
  - Rekap Saldo Simpanan per Anggota
  - Riwayat Transaksi Simpanan
- Angsuran
  - Daftar Angsuran (per anggota atau semua pinjaman)
  - Jadwal Angsuran (lihat rincian pokok, bunga, denda, tanggal jatuh tempo)
  - Pembayaran Angsuran (input bayar, otomatis hitung denda jika terlambat)
  - Pelunasan Dipercepat (hitung pelunasan total, update status pinjaman)
  - Riwayat Bayar (semua cicilan yang sudah dibayar)
  - Tunggakan & Denda (daftar yang belum bayar, nominal tunggakan, denda)
  - Notifikasi Jatuh Tempo (opsional, via WhatsApp/Email)
- Pinjaman
  - Pengajuan Pinjaman (nominal, tenor, bunga)
  - Analisis & Persetujuan/Verifikasi
  - Pencairan Pinjaman
  - Jadwal Angsuran & Pembayaran Angsuran
  - Denda & Tunggakan
- Kas & Jurnal Sederhana
  - Penerimaan Kas
  - Pengeluaran Kas
  - Buku Kas Umum
  - Jurnal Transaksi (mapping sederhana, bukan modul akuntansi penuh)
- Laporan
  - Laporan Anggota (aktif/non-aktif)
  - Laporan Simpanan (periode harian/mingguan/bulanan)
  - Laporan Pinjaman (pengajuan, disetujui, berjalan, lunas)
  - Laporan Angsuran (bayar, tunggakan, denda)
  - Laporan Kas (cash-in/cash-out)
  - Ekspor CSV/PDF
- Pengaturan
  - Profil Koperasi (nama, alamat, kontak)
  - Parameter Finansial (suku bunga, biaya administrasi)
  - Kategori Simpanan & Pinjaman
  - Format Nomor Transaksi
  - Integrasi (Email/WhatsApp Gateway, opsional)
- Pengguna & Keamanan
  - Manajemen Pengguna (admin, petugas)
  - Role-Based Access Control (RBAC)
  - Audit Log (opsional)
- Bantuan
  - FAQ & Panduan Penggunaan
  - Kontak Dukungan

## Alur Bisnis Utama
- Pendaftaran Anggota → Verifikasi → Aktivasi
- Simpanan: Setoran → Update Saldo → Riwayat
- Penarikan: Permohonan → Persetujuan → Proses Penarikan
- Pinjaman: Pengajuan → Analisis → Persetujuan → Pencairan → Pembayaran Angsuran → Pelunasan
- Kas: Catat transaksi penerimaan/pengeluaran → Rekap buku kas

## Skema Basis Data (Ringkas)
- `anggota(id, nomor_anggota, nama, nik, alamat, telp, status, tanggal_gabung)`
- `pegawai(id, nama, email, role, status)`
- `simpanans(id, anggota_id, jenis, tanggal, jumlah, saldo_akhir)`
- `penarikans(id, anggota_id, jenis, tanggal, jumlah)`
- `pinjamans(id, anggota_id, nomor_pinjaman, tanggal_pengajuan, nominal, tenor_bulan, bunga_persen, status)`
- `angsuran(id, pinjaman_id, ke, tanggal_jatuh_tempo, jumlah, tanggal_bayar, denda)`
- `kas(id, tanggal, jenis, keterangan, jumlah, ref)`
- `users(id, email, password_hash, role, status)`
- `audit_log(id, user_id, action, entity, entity_id, timestamp, note)`
- `settings(key, value)`

Catatan: Nama tabel dapat disesuaikan. Backend menggunakan GORM sehingga migrasi dapat dilakukan otomatis berdasarkan model.

## API (Contoh Endpoint)
- Autentikasi
  - `POST /api/auth/login` → JWT
  - `POST /api/auth/logout`
- Anggota
  - `GET /api/anggota`
  - `POST /api/anggota`
  - `GET /api/anggota/:id` / `PUT /api/anggota/:id`
- Simpanan & Penarikan
  - `GET /api/simpanan?anggota_id=...`
  - `POST /api/simpanan/setoran`
  - `POST /api/simpanan/penarikan`
- Pinjaman & Angsuran
  - `GET /api/pinjaman`
  - `POST /api/pinjaman/pengajuan`
  - `POST /api/pinjaman/verifikasi`
  - `POST /api/pinjaman/pencairan`
  - `GET /api/angsuran?pinjaman_id=...`
  - `POST /api/angsuran/bayar`
- Kas & Jurnal
  - `GET /api/kas`
  - `POST /api/kas/in` / `POST /api/kas/out`
- Laporan
  - `GET /api/laporan/simpanan?periode=...`
  - `GET /api/laporan/pinjaman?status=...`
  - `GET /api/laporan/kas?periode=...`

## HTTP Client (Axios)
Untuk komunikasi HTTP di frontend, proyek ini menggunakan `Axios` sebagai client.

### Instalasi
- Jalankan perintah berikut di folder `web`:

```bash
cd web
npm install axios
```

### Konfigurasi Instance
Buat instance Axios agar konsisten menggunakan `VITE_API_BASE_URL` dan otomatis menambahkan JWT (jika tersedia).

```ts
// web/src/lib/axios.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Tambahkan Authorization header jika token tersimpan
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

### Contoh Penggunaan
- Login dan simpan token:

```ts
// web/src/services/auth.ts
import api from '@/lib/axios';

export async function login(email: string, password: string) {
  const { data } = await api.post('/api/auth/login', { email, password });
  localStorage.setItem('token', data.token);
  return data;
}
```

- Memuat daftar anggota:

```ts
// web/src/services/anggota.ts
import api from '@/lib/axios';

export async function getAnggota() {
  const { data } = await api.get('/api/anggota');
  return data;
}
```

## Setup Lingkungan – Development
### Prasyarat
- `Node.js >= 18`, `npm >= 9`
- `Go >= 1.22`
- `MySQL >= 8`

### Konfigurasi Environment
Frontend (`web/.env`):
- `VITE_APP_NAME` → "Sistem Koperasi Kantor Desa"
- `VITE_API_BASE_URL` → contoh: `http://localhost:8080`

Backend (`service/.env`):
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_AUDIENCE`

### Menjalankan Aplikasi
- Frontend
  - `cd web`
  - `npm install`
  - `npm run dev` → buka `http://localhost:5173` (default Vite)
- Backend
  - `cd service`
  - `go mod tidy`
  - `go run ./cmd/main.go` atau `make run` (jika tersedia)

### Database & Migrasi
- Pastikan `MySQL` berjalan dan kredensial sesuai di `service/.env`.
- Jalankan migrasi otomatis via GORM (model akan dimigrasikan saat service start).
  - Driver: `gorm.io/driver/mysql`
  - Contoh DSN: `user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`

## Deployment (Opsional: Docker Compose)
Contoh ringkas layanan backend + database (MySQL):
```yaml
version: '3.9'
services:
  db:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: koperasi
      MYSQL_USER: koperasi_user
      MYSQL_PASSWORD: strong_password
      MYSQL_ROOT_PASSWORD: root_password
    ports:
      - "3306:3306"
    command: ["--default-authentication-plugin=mysql_native_password"]
    volumes:
      - mysql_data:/var/lib/mysql

  service:
    build: ./service
    environment:
      DB_HOST: db
      DB_PORT: 3306
      DB_USER: koperasi_user
      DB_PASS: strong_password
      DB_NAME: koperasi
      JWT_SECRET: super_secret_key
    ports:
      - "8080:8080"
    depends_on:
      - db

volumes:
  mysql_data:
```

## Keamanan
- Gunakan `HTTPS` di production.
- Simpan `JWT_SECRET` secara aman (env/secret manager).
- Terapkan `RBAC` untuk membatasi akses per peran.
- Validasi input dan sanitasi data.
- Audit log untuk aksi penting (opsional).

## Pengujian
- Backend: unit test untuk service, repository, dan handler.
- Frontend: e2e test untuk alur anggota & admin.
- Postman: koleksi request untuk uji integrasi API.

## Roadmap
- Notifikasi WhatsApp/Email (pengingat angsuran, pengajuan disetujui)
- Laporan grafik interaktif
- Multi-tenant (jika dibutuhkan)
- Aplikasi mobile (PWA/Hybrid)

## Kontribusi
- Buat branch fitur, ajukan PR, sertakan deskripsi perubahan & bukti pengujian.

## Lisensi
- Internal/Proprietary (sesuaikan kebutuhan).
