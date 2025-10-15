export type MenuItem = {
  label: string
  path?: string
  children?: MenuItem[]
}

export const menuItems: MenuItem[] = [
  { label: 'Dashboard', path: '/dashboard' },
  {
    label: 'Keanggotaan',
    path: '/keanggotaan',
  },
  { label: 'Simpanan', path: '/simpanan' },
  { label: 'Penarikan', path: '/penarikan' },
  { label: 'Pinjaman', path: '/pinjaman' },
  { label: 'Angsuran', path: '/angsuran' },
  { label: 'Kas & Jurnal', path: '/kas' },
  { label: 'Laporan', path: '/laporan' },
  { label: 'Pengaturan', path: '/pengaturan' },
  { label: 'Pengguna & Keamanan', path: '/pengguna' },
  { label: 'Bantuan', path: '/bantuan' },
]