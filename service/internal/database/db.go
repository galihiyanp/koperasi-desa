package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "koperasi-desa/service/internal/models"
)

func InitDB() *gorm.DB {
    host := getenv("DB_HOST", "127.0.0.1")
    port := getenv("DB_PORT", "3306")
    name := getenv("DB_NAME", "koperasi_desa")
    user := getenv("DB_USER", "root")
    pass := getenv("DB_PASS", "")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4", user, pass, host, port, name)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
    if err != nil {
        log.Printf("failed to connect MySQL: %v", err)
        // Fallback dev-only ke SQLite agar server tetap berjalan
        db, err = gorm.Open(sqlite.Open("file:koperasi_dev.db?cache=shared&_fk=1"), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
        if err != nil {
            log.Fatalf("failed to open fallback SQLite: %v", err)
        }
        log.Printf("using SQLite fallback database 'koperasi_dev.db'")
    }

    // Migrations
    if err := db.AutoMigrate(
        &models.User{},
        &models.Anggota{},
        &models.AnggotaDocument{},
        &models.AnggotaActivity{},
        &models.Simpanan{},
        &models.Pinjaman{},
        &models.Angsuran{},
        &models.Setting{},
    ); err != nil {
        log.Fatalf("failed to migrate: %v", err)
    }

    return db
}

func getenv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}