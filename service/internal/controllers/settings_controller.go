package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

// SettingsController mengelola konfigurasi sistem berbasis key-value
// Nilai value disepakati bertipe string; untuk data kompleks gunakan JSON string.
type SettingsController struct { DB *gorm.DB }
func NewSettingsController(db *gorm.DB) *SettingsController { return &SettingsController{DB: db} }

// GET /api/settings
// Mengembalikan daftar semua settings
func (h *SettingsController) ListSettings(c *gin.Context) {
    var list []models.Setting
    if err := h.DB.Order("`key` ASC").Find(&list).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": list})
}

// GET /api/settings/:key
func (h *SettingsController) GetSetting(c *gin.Context) {
    key := c.Param("key")
    var s models.Setting
    if err := h.DB.First(&s, "`key` = ?", key).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "setting tidak ditemukan"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, s)
}

// PUT /api/settings/:key
// Body: { value: string }
// Upsert nilai setting berdasarkan key
func (h *SettingsController) PutSetting(c *gin.Context) {
    key := c.Param("key")
    var body struct { Value string `json:"value"` }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if key == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "key wajib diisi"})
        return
    }
    s := models.Setting{Key: key, Value: body.Value}
    // upsert berdasarkan primary key (Name)
    if err := h.DB.Save(&s).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}