package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "koperasi-desa/service/internal/models"
)

type UserController struct {
    DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{DB: db}
}

type CreateUserInput struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role"`
}

type UpdateUserInput struct {
    Name     *string `json:"name"`
    Email    *string `json:"email"`
    Password *string `json:"password"`
    Role     *string `json:"role"`
}

func (h *UserController) CreateUser(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hashed),
        Role:     input.Role,
    }
    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

func (h *UserController) ListUsers(c *gin.Context) {
    var users []models.User
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if page < 1 { page = 1 }
    if limit < 1 || limit > 100 { limit = 10 }
    offset := (page - 1) * limit

    if err := h.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": users, "page": page, "limit": limit})
}

func (h *UserController) GetUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := h.DB.First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *UserController) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := h.DB.First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input UpdateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update fields selectively
    if input.Name != nil { user.Name = *input.Name }
    if input.Email != nil { user.Email = *input.Email }
    if input.Role != nil { user.Role = *input.Role }
    if input.Password != nil && *input.Password != "" {
        hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
            return
        }
        user.Password = string(hashed)
    }

    if err := h.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *UserController) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"deleted": true})
}