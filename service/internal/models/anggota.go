package models

import "time"

// Anggota merepresentasikan member koperasi
type Anggota struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    NomorAnggota  string    `gorm:"size:64;uniqueIndex" json:"nomor_anggota"`
    Nama          string    `gorm:"size:128" json:"nama"`
    NIK           string    `gorm:"size:32" json:"nik"`
    Alamat        string    `gorm:"size:255" json:"alamat"`
    Telp          string    `gorm:"size:32" json:"telp"`
    Status        string    `gorm:"size:32" json:"status"`
    TanggalGabung time.Time `json:"tanggal_gabung"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    Documents     []AnggotaDocument `json:"documents"`
}

// AnggotaDocument menyimpan metadata dokumen anggota (KTP, KK, dll)
type AnggotaDocument struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    AnggotaID uint      `json:"anggota_id"`
    Filename  string    `gorm:"size:255" json:"filename"`
    URL       string    `gorm:"size:255" json:"url"`
    UploadedAt time.Time `json:"uploaded_at"`
}

// AnggotaActivity mencatat aktivitas penting (pendaftaran, verifikasi, aktivasi)
type AnggotaActivity struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    AnggotaID uint      `json:"anggota_id"`
    Action    string    `gorm:"size:64" json:"action"`
    Note      string    `gorm:"size:255" json:"note"`
    CreatedAt time.Time `json:"created_at"`
}