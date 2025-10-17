package models

import "time"

// Pinjaman merepresentasikan entitas pinjaman
// Status: pengajuan | disetujui | berjalan | lunas
type Pinjaman struct {
    ID                uint       `gorm:"primaryKey" json:"id"`
    AnggotaID         uint       `json:"anggota_id"`
    NomorPinjaman     string     `gorm:"size:64;uniqueIndex" json:"nomor_pinjaman"`
    TanggalPengajuan  time.Time  `json:"tanggal_pengajuan"`
    TanggalDisetujui  *time.Time `json:"tanggal_disetujui"`
    TanggalPencairan  *time.Time `json:"tanggal_pencairan"`
    Nominal           float64    `json:"nominal"`
    TenorBulan        int        `json:"tenor_bulan"`
    BungaPersen       float64    `json:"bunga_persen"`
    Status            string     `gorm:"size:32" json:"status"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`

    Angsuran          []Angsuran `json:"-"`
}