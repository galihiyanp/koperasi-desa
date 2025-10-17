package models

import "time"

// Simpanan merepresentasikan transaksi simpanan per anggota
// Mendukung setoran dan penarikan dengan saldo_akhir per jenis.
type Simpanan struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    AnggotaID  uint      `json:"anggota_id"`
    Jenis      string    `gorm:"size:32" json:"jenis"`      // wajib | sukarela | khusus
    Tipe       string    `gorm:"size:16" json:"tipe"`       // setoran | penarikan
    Tanggal    time.Time `json:"tanggal"`
    Jumlah     float64   `json:"jumlah"`
    SaldoAkhir float64   `json:"saldo_akhir"`
    CreatedAt  time.Time `json:"created_at"`
}