package controllers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

type AnggotaController struct {
    DB *gorm.DB
}

func NewAnggotaController(db *gorm.DB) *AnggotaController {
    return &AnggotaController{DB: db}
}

type CreateAnggotaInput struct {
    NomorAnggota  string    `json:"nomor_anggota" binding:"required"`
    Nama          string    `json:"nama" binding:"required"`
    NIK           string    `json:"nik" binding:"required"`
    Alamat        string    `json:"alamat"`
    Telp          string    `json:"telp"`
    TanggalGabung time.Time `json:"tanggal_gabung"`
}

type UpdateAnggotaInput struct {
    NomorAnggota  *string    `json:"nomor_anggota"`
    Nama          *string    `json:"nama"`
    NIK           *string    `json:"nik"`
    Alamat        *string    `json:"alamat"`
    Telp          *string    `json:"telp"`
    Status        *string    `json:"status"`
    TanggalGabung *time.Time `json:"tanggal_gabung"`
}

func (h *AnggotaController) CreateAnggota(c *gin.Context) {
    var input CreateAnggotaInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if input.TanggalGabung.IsZero() { input.TanggalGabung = time.Now() }
    anggota := models.Anggota{
        NomorAnggota:  input.NomorAnggota,
        Nama:          input.Nama,
        NIK:           input.NIK,
        Alamat:        input.Alamat,
        Telp:          input.Telp,
        Status:        "pending",
        TanggalGabung: input.TanggalGabung,
    }
    if err := h.DB.Create(&anggota).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Log activity
    h.DB.Create(&models.AnggotaActivity{AnggotaID: anggota.ID, Action: "registered", Note: "Pendaftaran anggota", CreatedAt: time.Now()})
    c.JSON(http.StatusCreated, anggota)
}

func (h *AnggotaController) ListAnggota(c *gin.Context) {
    var list []models.Anggota
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if page < 1 { page = 1 }
    if limit < 1 || limit > 100 { limit = 10 }
    offset := (page - 1) * limit

    status := strings.TrimSpace(c.Query("status"))
    q := strings.TrimSpace(c.Query("q"))

    tx := h.DB.Model(&models.Anggota{})
    if status != "" { tx = tx.Where("status = ?", status) }
    if q != "" { tx = tx.Where("nama LIKE ? OR nik LIKE ? OR nomor_anggota LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%") }

    if err := tx.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": list, "page": page, "limit": limit})
}

func (h *AnggotaController) GetAnggota(c *gin.Context) {
    id := c.Param("id")
    var anggota models.Anggota
    if err := h.DB.Preload("Documents").First(&anggota, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "anggota not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    var acts []models.AnggotaActivity
    _ = h.DB.Where("anggota_id = ?", anggota.ID).Order("created_at DESC").Find(&acts).Error
    c.JSON(http.StatusOK, gin.H{"data": anggota, "activities": acts})
}

func (h *AnggotaController) UpdateAnggota(c *gin.Context) {
    id := c.Param("id")
    var anggota models.Anggota
    if err := h.DB.First(&anggota, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "anggota not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input UpdateAnggotaInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.NomorAnggota != nil { anggota.NomorAnggota = *input.NomorAnggota }
    if input.Nama != nil { anggota.Nama = *input.Nama }
    if input.NIK != nil { anggota.NIK = *input.NIK }
    if input.Alamat != nil { anggota.Alamat = *input.Alamat }
    if input.Telp != nil { anggota.Telp = *input.Telp }
    if input.TanggalGabung != nil { anggota.TanggalGabung = *input.TanggalGabung }
    if input.Status != nil {
        s := strings.ToLower(*input.Status)
        if s == "pending" || s == "verified" || s == "active" {
            anggota.Status = s
        }
    }

    if err := h.DB.Save(&anggota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, anggota)
}

func (h *AnggotaController) VerifyAnggota(c *gin.Context) {
    id := c.Param("id")
    var anggota models.Anggota
    if err := h.DB.First(&anggota, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "anggota not found"})
        return
    }
    anggota.Status = "verified"
    if err := h.DB.Save(&anggota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    h.DB.Create(&models.AnggotaActivity{AnggotaID: anggota.ID, Action: "verified", Note: "Verifikasi anggota", CreatedAt: time.Now()})
    c.JSON(http.StatusOK, gin.H{"verified": true})
}

func (h *AnggotaController) ActivateAnggota(c *gin.Context) {
    id := c.Param("id")
    var anggota models.Anggota
    if err := h.DB.First(&anggota, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "anggota not found"})
        return
    }
    anggota.Status = "active"
    if err := h.DB.Save(&anggota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    h.DB.Create(&models.AnggotaActivity{AnggotaID: anggota.ID, Action: "activated", Note: "Aktivasi anggota", CreatedAt: time.Now()})
    c.JSON(http.StatusOK, gin.H{"activated": true})
}

func (h *AnggotaController) UploadDocuments(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.Atoi(idStr)
    var anggota models.Anggota
    if err := h.DB.First(&anggota, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "anggota not found"})
        return
    }

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
        return
    }
    files := form.File["files"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no files"})
        return
    }

    baseDir := filepath.Join("uploads", "anggota", fmt.Sprintf("%d", anggota.ID))
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
        return
    }
    var saved []models.AnggotaDocument
    for _, f := range files {
        filename := filepath.Base(f.Filename)
        path := filepath.Join(baseDir, filename)
        if err := c.SaveUploadedFile(f, path); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
            return
        }
        url := "/uploads/anggota/" + fmt.Sprintf("%d", anggota.ID) + "/" + filename
        doc := models.AnggotaDocument{AnggotaID: uint(id), Filename: filename, URL: url, UploadedAt: time.Now()}
        if err := h.DB.Create(&doc).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record document"})
            return
        }
        saved = append(saved, doc)
    }
    c.JSON(http.StatusOK, gin.H{"uploaded": saved})
}

func (h *AnggotaController) ListDocuments(c *gin.Context) {
    id := c.Param("id")
    var docs []models.AnggotaDocument
    if err := h.DB.Where("anggota_id = ?", id).Order("uploaded_at DESC").Find(&docs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": docs})
}