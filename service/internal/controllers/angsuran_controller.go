package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

type AngsuranController struct { DB *gorm.DB }
func NewAngsuranController(db *gorm.DB) *AngsuranController { return &AngsuranController{DB: db} }

// GET /api/angsuran?pinjaman_id=...
func (h *AngsuranController) ListAngsuran(c *gin.Context) {
    var list []models.Angsuran
    pinjamanID := c.Query("pinjaman_id")

    tx := h.DB.Model(&models.Angsuran{})
    if pinjamanID != "" { tx = tx.Where("pinjaman_id = ?", pinjamanID) }

    if err := tx.Order("ke ASC").Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": list})
}

// POST /api/angsuran/bayar
// { angsuran_id, jumlah, tanggal_bayar }
type BayarAngsuranInput struct {
    AngsuranID   uint       `json:"angsuran_id"`
    Jumlah       float64    `json:"jumlah"`
    TanggalBayar *time.Time `json:"tanggal_bayar"`
}

func (h *AngsuranController) Bayar(c *gin.Context) {
    var in BayarAngsuranInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if in.AngsuranID == 0 || in.Jumlah <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
        return
    }

    err := h.DB.Transaction(func(tx *gorm.DB) error {
        var a models.Angsuran
        if err := tx.First(&a, in.AngsuranID).Error; err != nil { return err }
        if a.TanggalBayar != nil { return gorm.ErrInvalidTransaction }

        // validasi jumlah (minimal sama dengan jumlah angsuran)
        if in.Jumlah < a.Jumlah { return gorm.ErrInvalidTransaction }

        // set tanggal bayar, hitung denda sederhana jika terlambat
        tanggal := time.Now()
        if in.TanggalBayar != nil { tanggal = *in.TanggalBayar }
        a.TanggalBayar = &tanggal
        if tanggal.After(a.TanggalJatuhTempo) {
            // denda flat 1% dari jumlah
            a.Denda = a.Jumlah * 0.01
        }
        if err := tx.Save(&a).Error; err != nil { return err }

        // Jika semua angsuran sudah dibayar, set status pinjaman ke 'lunas'
        var remaining int64
        if err := tx.Model(&models.Angsuran{}).Where("pinjaman_id = ? AND tanggal_bayar IS NULL", a.PinjamanID).Count(&remaining).Error; err != nil { return err }
        if remaining == 0 {
            _ = tx.Model(&models.Pinjaman{}).Where("id = ?", a.PinjamanID).Update("status", "lunas").Error
        }
        return nil
    })

    if err != nil {
        if err == gorm.ErrInvalidTransaction {
            c.JSON(http.StatusBadRequest, gin.H{"error": "transaksi tidak valid"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}