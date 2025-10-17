package models

import "time"

// Setting menyimpan konfigurasi sistem berbasis key-value.
// Nilai disimpan sebagai string (dapat berupa JSON string untuk struktur kompleks).
// Contoh key:
// - settings.profile
// - settings.financial
// - settings.categories
// - settings.format
// - settings.integrations
//
// Lihat README bagian Pengaturan untuk spesifikasi.
type Setting struct {
    Key       string    `gorm:"primaryKey;size:128;column:key" json:"key"`
    Value     string    `gorm:"type:text" json:"value"`
    UpdatedAt time.Time `json:"updated_at"`
}