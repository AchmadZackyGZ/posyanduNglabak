package controllers

import (
	"net/http"
	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// --- DTO Inputs ---
type CreateUserInput struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Role        string `json:"role" binding:"required"` // Pilihan: ADMIN, BIDAN, KADER
}

type UpdateUserInput struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Role        string `json:"role" binding:"required"`
	IsActive    *bool  `json:"is_active" binding:"required"` // Pakai pointer agar bisa menerima nilai 'false' murni dari JSON
}

type ResetPasswordInput struct {
	NewPassword string `json:"new_password" binding:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	NamaLengkap string `json:"nama_lengkap"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
}

// 1. Ambil Daftar Pengguna (Kecuali ADMIN)
func (uc *UserController) GetListUsers(c *gin.Context) {
	var users []UserResponse
	
	// Gunakan Model(&models.User{}) agar GORM tahu tabel mana yang diakses,
	// lalu petakan hasilnya langsung ke struct UserResponse
	if err := uc.DB.Model(&models.User{}).
		Select("id, username, nama_lengkap, role, is_active").
		Where("role IN ?", []string{"KADER", "BIDAN"}). // SECARA SPESIFIK HANYA MENARIK STAFF
		Order("created_at desc").
		Find(&users).Error; err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil", 
		"data": users,
	})
}

// 1b. Ambil Daftar Warga (Khusus Role USER)
func (uc *UserController) GetPublicUsers(c *gin.Context) {
	var users []UserResponse
	
	// Hanya tarik data yang role-nya adalah 'USER'
	if err := uc.DB.Model(&models.User{}).
		Select("id, username, nama_lengkap, role, is_active").
		Where("role = ?", "USER").
		Order("created_at desc").
		Find(&users).Error; err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data warga"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil", 
		"data": users,
	})
}

// 2. Buat Pengguna Baru
func (uc *UserController) CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah username sudah dipakai
	var existingUser models.User
	if err := uc.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah terdaftar di sistem"})
		return
	}

	// Enkripsi kata sandi
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses kata sandi"})
		return
	}

	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		NamaLengkap:  input.NamaLengkap,
		Role:         input.Role,
		IsActive:     true,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengguna berhasil ditambahkan"})
}

// 3. Update Profil/Status Pengguna
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := uc.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	user.NamaLengkap = input.NamaLengkap
	user.Role = input.Role
	user.IsActive = *input.IsActive

	uc.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Data pengguna berhasil diperbarui"})
}

// 4. Reset Kata Sandi (Hak Istimewa Admin)
func (uc *UserController) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := uc.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)

	uc.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Kata sandi pengguna berhasil direset"})
}