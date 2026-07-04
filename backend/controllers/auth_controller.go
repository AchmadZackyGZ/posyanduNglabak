package controllers

import (
	"net/http"
	"os"
	"time"

	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 1. Tambahkan Struct Input untuk Registrasi Publik
type RegisterInput struct {
	Username    string `json:"username" binding:"required"` // Biasanya diisi Nomor HP oleh warga
	Password    string `json:"password" binding:"required,min=6"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format username atau password tidak valid"})
		return
	}

	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda sedang dinonaktifkan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Terbitkan Token JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menerbitkan token otentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"nama_lengkap": user.NamaLengkap,
			"role":         user.Role,
		},
	})
}

// 2. Tambahkan Fungsi Register untuk Warga (Role "USER")
func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua kolom wajib diisi (Password min 6 karakter)"})
		return
	}

	// Cek apakah username (nomor HP) sudah pernah dipakai
	var existingUser models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username / Nomor HP tersebut sudah terdaftar"})
		return
	}

	// Lakukan hashing pada password demi keamanan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses password"})
		return
	}

	// Bentuk entitas User baru
	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		NamaLengkap:  input.NamaLengkap,
		Role:         "USER", // HARCODED: Kunci hak akses pendaftar hanya sebagai warga biasa
		IsActive:     true,   // Langsung aktif agar bisa langsung login
	}

	// Eksekusi penyimpanan ke database
	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan akun ke sistem"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil! Silakan login menggunakan akun Anda.",
		"data": gin.H{
			"username":     user.Username,
			"nama_lengkap": user.NamaLengkap,
			"role":         user.Role,
		},
	})
}