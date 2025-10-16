package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/controllers"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    uc := controllers.NewUserController(db)
    ac := controllers.NewAnggotaController(db)

    api := r.Group("/api")
    {
        api.GET("/users", uc.ListUsers)
        api.POST("/users", uc.CreateUser)
        api.GET("/users/:id", uc.GetUser)
        api.PUT("/users/:id", uc.UpdateUser)
        api.DELETE("/users/:id", uc.DeleteUser)

        // Anggota routes sesuai spesifikasi README
        api.GET("/anggota", ac.ListAnggota)
        api.POST("/anggota", ac.CreateAnggota)
        api.GET("/anggota/:id", ac.GetAnggota)
        api.PUT("/anggota/:id", ac.UpdateAnggota)
        api.POST("/anggota/:id/verify", ac.VerifyAnggota)
        api.POST("/anggota/:id/activate", ac.ActivateAnggota)
        api.POST("/anggota/:id/documents", ac.UploadDocuments)
        api.GET("/anggota/:id/documents", ac.ListDocuments)
    }
}