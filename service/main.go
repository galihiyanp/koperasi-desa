package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/gorm"

    dbpkg "koperasi-desa/service/internal/database"
    "koperasi-desa/service/internal/routes"
)

func setupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // CORS untuk frontend Vite
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Register routes
    routes.RegisterRoutes(r, db)
    return r
}

func main() {
    _ = godotenv.Load(".env")

    db := dbpkg.InitDB()
    r := setupRouter(db)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}