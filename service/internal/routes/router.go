package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/controllers"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    uc := controllers.NewUserController(db)

    api := r.Group("/api")
    {
        api.GET("/users", uc.ListUsers)
        api.POST("/users", uc.CreateUser)
        api.GET("/users/:id", uc.GetUser)
        api.PUT("/users/:id", uc.UpdateUser)
        api.DELETE("/users/:id", uc.DeleteUser)
    }
}