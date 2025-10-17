package controllers

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

type PinjamanController struct { DB *gorm.DB }
func NewPinjamanController(db *gorm.DB) *PinjamanController { return &PinjamanController{DB: db} }

// GET /api/pinjaman?anggota_id=...&status=...&page=...&limit=...
func (h *PinjamanController) ListPinjaman(c *gin.Context) {
    var list []models.Pinjaman
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if page < 1 { page = 1 }
    if limit < 1 || limit > 100 { limit = 10 }
    offset := (page - 1) * limit

    anggotaID := strings.TrimSpace(c.Query("anggota_id"))
    status := strings.TrimSpace(c.Query("status"))

    tx := h.DB.Model(&models.Pinjaman{})
    if anggotaID != "" { tx = tx.Where("anggota_id = ?", anggotaID) }
    if status != "" { tx = tx.Where("status = ?", strings.ToLower(status)) }

    if err := tx.Order("created_at DESC, id DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": list, "page": page, "limit": limit})
}

// POST /api/pinjaman/pengajuan
// { anggota_id, nominal, tenor_bulan, bunga_persen, tanggal_pengajuan? }
type PinjamanPengajuanInput struct {
    AnggotaID    uint       `json:"anggota_id"`
    Nominal      float64    `json:"nominal"`
    TenorBulan   int        `json:"tenor_bulan"`
    BungaPersen  float64    `json:"bunga_persen"`
    Tanggal      *time.Time `json:"tanggal_pengajuan"`
}

func (h *PinjamanController) Pengajuan(c *gin.Context) {
    var in PinjamanPengajuanInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if in.AnggotaID == 0 || in.Nominal <= 0 || in.TenorBulan <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
        return
    }
    // ensure anggota exists
    var a models.Anggota
    if err := h.DB.First(&a, in.AnggotaID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "anggota tidak ditemukan"})
        return
    }

    tanggal := time.Now()
    if in.Tanggal != nil { tanggal = *in.Tanggal }

    p := models.Pinjaman{
        AnggotaID:        in.AnggotaID,
        TanggalPengajuan: tanggal,
        Nominal:          in.Nominal,
        TenorBulan:       in.TenorBulan,
        BungaPersen:      in.BungaPersen,
        Status:           "pengajuan",
    }
    if err := h.DB.Create(&p).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // generate nomor pinjaman setelah ID tersedia
    p.NomorPinjaman = fmt.Sprintf("PJ-%d-%06d", time.Now().Year(), p.ID)
    _ = h.DB.Save(&p).Error

    c.JSON(http.StatusCreated, p)
}

// POST /api/pinjaman/verifikasi { pinjaman_id }
type PinjamanActionInput struct { PinjamanID uint `json:"pinjaman_id"` }
func (h *PinjamanController) Verifikasi(c *gin.Context) {
    var in PinjamanActionInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    var p models.Pinjaman
    if err := h.DB.First(&p, in.PinjamanID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "pinjaman tidak ditemukan"})
        return
    }
    now := time.Now()
    p.TanggalDisetujui = &now
    p.Status = "disetujui"
    if err := h.DB.Save(&p).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/pinjaman/pencairan { pinjaman_id }
func (h *PinjamanController) Pencairan(c *gin.Context) {
    var in PinjamanActionInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := h.DB.Transaction(func(tx *gorm.DB) error {
        var p models.Pinjaman
        if err := tx.First(&p, in.PinjamanID).Error; err != nil { return err }
        now := time.Now()
        p.TanggalPencairan = &now
        p.Status = "berjalan"
        if err := tx.Save(&p).Error; err != nil { return err }

        // generate schedule angsuran sederhana: (nominal + bunga_flat) / tenor
        total := p.Nominal + (p.Nominal * p.BungaPersen / 100.0)
        angsur := total / float64(p.TenorBulan)
        baseDate := now
        var batch []models.Angsuran
        for i := 1; i <= p.TenorBulan; i++ {
            due := baseDate.AddDate(0, i, 0)
            batch = append(batch, models.Angsuran{
                PinjamanID:        p.ID,
                Ke:                i,
                TanggalJatuhTempo: due,
                Jumlah:            angsur,
                Denda:             0,
            })
        }
        if err := tx.Create(&batch).Error; err != nil { return err }
        return nil
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}