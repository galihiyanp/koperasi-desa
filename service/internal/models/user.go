package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:128" json:"name"`
    Email     string    `gorm:"size:128;uniqueIndex" json:"email"`
    Password  string    `gorm:"size:255" json:"-"`
    Role      string    `gorm:"size:64" json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}