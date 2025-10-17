package models

import "time"

// Angsuran merepresentasikan jadwal dan pembayaran angsuran untuk pinjaman
// Jika tanggal_bayar NULL maka angsuran belum dibayar
// Denda opsional saat terjadi keterlambatan
type Angsuran struct {
    ID                 uint       `gorm:"primaryKey" json:"id"`
    PinjamanID         uint       `json:"pinjaman_id"`
    Ke                 int        `json:"ke"`
    TanggalJatuhTempo  time.Time  `json:"tanggal_jatuh_tempo"`
    Jumlah             float64    `json:"jumlah"`
    TanggalBayar       *time.Time `json:"tanggal_bayar"`
    Denda              float64    `json:"denda"`
    CreatedAt          time.Time  `json:"created_at"`
}