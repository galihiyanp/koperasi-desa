package controllers

import (
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

type SimpananController struct {
    DB *gorm.DB
}

func NewSimpananController(db *gorm.DB) *SimpananController { return &SimpananController{DB: db} }

// Input payloads
type SetoranInput struct {
    AnggotaID uint     `json:"anggota_id" binding:"required"`
    Jenis     string   `json:"jenis" binding:"required"`   // wajib | sukarela | khusus
    Jumlah    float64  `json:"jumlah" binding:"required,gt=0"`
    Tanggal   *time.Time `json:"tanggal"`
}

type PenarikanInput struct {
    AnggotaID uint     `json:"anggota_id" binding:"required"`
    Jenis     string   `json:"jenis" binding:"required"`
    Jumlah    float64  `json:"jumlah" binding:"required,gt=0"`
    Tanggal   *time.Time `json:"tanggal"`
}

// Helper to get latest saldo for anggota+jenis
func (h *SimpananController) latestSaldo(anggotaID uint, jenis string) (float64, error) {
    var last models.Simpanan
    err := h.DB.Where("anggota_id = ? AND jenis = ?", anggotaID, strings.ToLower(jenis)).Order("tanggal DESC, id DESC").First(&last).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0, nil
        }
        return 0, err
    }
    return last.SaldoAkhir, nil
}

// GET /api/simpanan?anggota_id=...&jenis=...&page=...&limit=...
func (h *SimpananController) ListSimpanan(c *gin.Context) {
    var list []models.Simpanan
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if page < 1 { page = 1 }
    if limit < 1 || limit > 100 { limit = 10 }
    offset := (page - 1) * limit

    anggotaIDStr := strings.TrimSpace(c.Query("anggota_id"))
    jenis := strings.ToLower(strings.TrimSpace(c.Query("jenis")))

    tx := h.DB.Model(&models.Simpanan{})
    if anggotaIDStr != "" { tx = tx.Where("anggota_id = ?", anggotaIDStr) }
    if jenis != "" { tx = tx.Where("jenis = ?", jenis) }

    if err := tx.Order("tanggal DESC, id DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Optional saldo rekap per jenis jika anggota_id diberikan
    saldo := gin.H{}
    if anggotaIDStr != "" {
        id64, _ := strconv.ParseUint(anggotaIDStr, 10, 64)
        wajib, _ := h.latestSaldo(uint(id64), "wajib")
        sukarela, _ := h.latestSaldo(uint(id64), "sukarela")
        khusus, _ := h.latestSaldo(uint(id64), "khusus")
        saldo = gin.H{"wajib": wajib, "sukarela": sukarela, "khusus": khusus}
    }

    c.JSON(http.StatusOK, gin.H{"data": list, "page": page, "limit": limit, "saldo": saldo})
}

// POST /api/simpanan/setoran
func (h *SimpananController) Setoran(c *gin.Context) {
    var input SetoranInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    jenis := strings.ToLower(input.Jenis)
    if jenis != "wajib" && jenis != "sukarela" && jenis != "khusus" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "jenis harus wajib/sukarela/khusus"})
        return
    }

    // Ensure anggota exists
    var a models.Anggota
    if err := h.DB.First(&a, input.AnggotaID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "anggota tidak ditemukan"})
        return
    }

    // Transactional insert
    err := h.DB.Transaction(func(tx *gorm.DB) error {
        latest, err := h.latestSaldo(input.AnggotaID, jenis)
        if err != nil { return err }
        tanggal := time.Now()
        if input.Tanggal != nil { tanggal = *input.Tanggal }
        rec := models.Simpanan{
            AnggotaID:  input.AnggotaID,
            Jenis:      jenis,
            Tipe:       "setoran",
            Tanggal:    tanggal,
            Jumlah:     input.Jumlah,
            SaldoAkhir: latest + input.Jumlah,
        }
        return tx.Create(&rec).Error
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"ok": true})
}

// POST /api/simpanan/penarikan
func (h *SimpananController) Penarikan(c *gin.Context) {
    var input PenarikanInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    jenis := strings.ToLower(input.Jenis)
    if jenis != "wajib" && jenis != "sukarela" && jenis != "khusus" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "jenis harus wajib/sukarela/khusus"})
        return
    }
    var a models.Anggota
    if err := h.DB.First(&a, input.AnggotaID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "anggota tidak ditemukan"})
        return
    }

    // Transactional insert with saldo check
    err := h.DB.Transaction(func(tx *gorm.DB) error {
        latest, err := h.latestSaldo(input.AnggotaID, jenis)
        if err != nil { return err }
        if latest < input.Jumlah {
            return gorm.ErrInvalidTransaction
        }
        tanggal := time.Now()
        if input.Tanggal != nil { tanggal = *input.Tanggal }
        rec := models.Simpanan{
            AnggotaID:  input.AnggotaID,
            Jenis:      jenis,
            Tipe:       "penarikan",
            Tanggal:    tanggal,
            Jumlah:     input.Jumlah,
            SaldoAkhir: latest - input.Jumlah,
        }
        return tx.Create(&rec).Error
    })
    if err != nil {
        if err == gorm.ErrInvalidTransaction {
            c.JSON(http.StatusBadRequest, gin.H{"error": "saldo tidak cukup"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusCreated, gin.H{"ok": true})
}