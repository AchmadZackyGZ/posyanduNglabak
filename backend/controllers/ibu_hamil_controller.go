package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IbuHamilController struct {
	DB *gorm.DB
}

func NewIbuHamilController(db *gorm.DB) *IbuHamilController {
	return &IbuHamilController{DB: db}
}

// Struktur input pendaftaran Ibu Hamil (FR-04)
type RegisterIbuHamilInput struct {
	NIK          string `json:"nik" binding:"required"`
	NamaIbu      string `json:"nama_ibu" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"` // YYYY-MM-DD
	HPL          string `json:"hpl" binding:"required"`           // YYYY-MM-DD (Hari Perkiraan Lahir)
	Alamat       string `json:"alamat" binding:"required"`
	NoHP         string `json:"no_hp" binding:"required"` // Username login mandiri
}

// GetListIbuHamil mengambil seluruh data ibu hamil
func (ic *IbuHamilController) GetListIbuHamil(c *gin.Context) {
	var ibuHamils []models.IbuHamil
	if err := ic.DB.Order("created_at desc").Find(&ibuHamils).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar ibu hamil"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar ibu hamil",
		"data":    ibuHamils,
	})
}

// RegisterIbuHamil mendaftarkan ibu hamil sekaligus membuatkan akun login secara transaksional
func (ic *IbuHamilController) RegisterIbuHamil(c *gin.Context) {
	var input RegisterIbuHamilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tglLahir, err1 := time.Parse("2006-01-02", input.TanggalLahir)
	hpl, err2 := time.Parse("2006-01-02", input.HPL)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	// Memulai blok transaksi database
	tx := ic.DB.Begin()

	// 1. Cari atau buat akun mandiri untuk Ibu Hamil (Role: USER)
	var ibuUser models.User
	if err := tx.Where("username = ?", input.NoHP).First(&ibuUser).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("posyandu123"), bcrypt.DefaultCost)
		ibuUser = models.User{
			Username:     input.NoHP,
			PasswordHash: string(hashedPassword),
			NamaLengkap:  input.NamaIbu,
			Role:         "USER",
			IsActive:     true,
		}
		if err := tx.Create(&ibuUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun otentikasi Ibu Hamil"})
			return
		}
	}

	// 2. Simpan entitas data klinis Ibu Hamil
	ibuHamil := models.IbuHamil{
		UserID:          ibuUser.ID,
		NIK:             input.NIK,
		NamaIbu:         input.NamaIbu,
		TanggalLahir:    tglLahir,
		HPL:             hpl,
		Alamat:          input.Alamat,
		StatusKehamilan: "AKTIF",
	}

	if err := tx.Create(&ibuHamil).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftar. NIK Ibu Hamil mungkin sudah terdaftar."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Ibu Hamil berhasil dicatat",
		"data": gin.H{
			"ibu_hamil_id": ibuHamil.ID,
			"nama":         ibuHamil.NamaIbu,
			"akun_login":   ibuUser.Username,
		},
	})
}

// Struktur input pemeriksaan bulanan kehamilan
type PeriksaIbuHamilInput struct {
	IbuHamilID    string  `json:"ibu_hamil_id" binding:"required"`
	UsiaKehamilan int     `json:"usia_kehamilan" binding:"required"` // Format minggu
	TekananDarah  string  `json:"tekanan_darah" binding:"required"`  // e.g. "120/80"
	BeratBadan    float64 `json:"berat_badan" binding:"required"`
	Catatan       string  `json:"catatan"`
}

// CatatPemeriksaan mengamankan entri kontrol kandungan bulanan oleh Bidan atau Kader
func (ic *IbuHamilController) CatatPemeriksaan(c *gin.Context) {
	var input PeriksaIbuHamilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengambil ID petugas dari memori Context Gin hasil verifikasi JWT
	petugasID, _ := c.Get("userID")

	pemeriksaan := models.PemeriksaanIbuHamil{
		IbuHamilID:     input.IbuHamilID,
		TanggalPeriksa: time.Now(),
		UsiaKehamilan:  input.UsiaKehamilan,
		TekananDarah:   input.TekananDarah,
		BeratBadan:     input.BeratBadan,
		Catatan:        input.Catatan,
		DiperiksaOleh:  petugasID.(string),
	}

	if err := ic.DB.Create(&pemeriksaan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rekam medis Ibu Hamil"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pemeriksaan klinis Ibu Hamil berhasil disimpan",
		"data":    pemeriksaan,
	})
}