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
	StatusKehamilan string `json:"status_kehamilan" binding:"required"` // Normal / Risiko Tinggi / dll
	NoHP         string `json:"no_hp" binding:"required"` // Username login mandiri
}

// Struktur input pemeriksaan bulanan kehamilan
type PeriksaIbuHamilInput struct {
	IbuHamilID    string  `json:"ibu_hamil_id" binding:"required"`
	UsiaKehamilan int     `json:"usia_kehamilan" binding:"required"` // Format minggu
	TekananDarah  string  `json:"tekanan_darah" binding:"required"`  // e.g. "120/80"
	BeratBadan    float64 `json:"berat_badan" binding:"required"`
	Catatan       string  `json:"catatan"`
}

// Struktur input untuk Update Ibu Hamil (Tanpa NoHP karena terikat akun User)
type UpdateIbuHamilInput struct {
	NIK             string `json:"nik" binding:"required"`
	NamaIbu         string `json:"nama_ibu" binding:"required"`
	HPL             string `json:"hpl" binding:"required"` // Format: YYYY-MM-DD
	StatusKehamilan string `json:"status_kehamilan" binding:"required"` // Normal / Risiko Tinggi / dll
	Alamat          string `json:"alamat" binding:"required"`
}

// GetListIbuHamil mengambil seluruh data ibu hamil
func (ic *IbuHamilController) GetListIbuHamil(c *gin.Context) {
	 // 1. Buat struct kustom (DTO) untuk menggabungkan data IbuHamil murni dengan tambahan NoHP
	type IbuHamilResponse struct {
		models.IbuHamil
		NoHP string `json:"no_hp"` // Menangkap username dari tabel users	
	}

	var results []IbuHamilResponse

	// 2. Lakukan Query JOIN ke tabel users
	// Kita menarik 'username' dari tabel users dan memberinya nama alias 'no_hp'
	err := ic.DB.Table("ibu_hamils").
		Select("ibu_hamils.*, users.username as no_hp").
		Joins("left join users on users.id = ibu_hamils.user_id").
		Order("ibu_hamils.created_at desc").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ibu hamil"})
		return
	}

	// 3. Kirimkan JSON yang sudah lengkap dengan NoHP ke SvelteKit
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar ibu hamil",
		"data":    results,
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
		StatusKehamilan: input.StatusKehamilan,
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

// GetRiwayatPeriksa mengambil seluruh data rekam medis kontrol kehamilan
func (ic *IbuHamilController) GetRiwayatPeriksa(c *gin.Context) {
	var riwayat []models.PemeriksaanIbuHamil
	
	// Preload("Pemeriksa") untuk menarik data Bidan/Kader dari tabel User
	query := ic.DB.Preload("Pemeriksa").Order("tanggal_periksa desc")

	// Fitur Filter: Jika ada query ?ibu_hamil_id=xxx, filter hanya untuk ibu hamil tersebut
	if ibuHamilID := c.Query("ibu_hamil_id"); ibuHamilID != "" {
		query = query.Where("ibu_hamil_id = ?", ibuHamilID)
	}

	if err := query.Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat pemeriksaan Ibu Hamil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat pemeriksaan",
		"data":    riwayat,
	})
}

// UpdateIbuHamil mengubah data profil ibu hamil berdasarkan ID
func (ic *IbuHamilController) UpdateIbuHamil(c *gin.Context) {
	id := c.Param("id")
	var ibuHamil models.IbuHamil

	// 1. Cek keberadaan data
	if err := ic.DB.First(&ibuHamil, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Ibu Hamil tidak ditemukan"})
		return
	}

	// 2. Bind JSON dari Frontend
	var input UpdateIbuHamilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Konversi format tanggal HPL
	hplParsed, err := time.Parse("2006-01-02", input.HPL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format HPL harus YYYY-MM-DD"})
		return
	}

	// 4. Timpa data lama
	ibuHamil.NIK = input.NIK
	ibuHamil.NamaIbu = input.NamaIbu
	ibuHamil.HPL = hplParsed
	ibuHamil.StatusKehamilan = input.StatusKehamilan
	ibuHamil.Alamat = input.Alamat

	// 5. Simpan ke database
	if err := ic.DB.Save(&ibuHamil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data ibu hamil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data Ibu Hamil berhasil diperbarui",
		"data":    ibuHamil,
	})
}

// DeleteIbuHamil menghapus data ibu hamil berdasarkan ID
func (ic *IbuHamilController) DeleteIbuHamil(c *gin.Context) {
	id := c.Param("id")
	var ibuHamil models.IbuHamil

	// 1. Pastikan data ada
	if err := ic.DB.First(&ibuHamil, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Ibu Hamil tidak ditemukan"})
		return
	}

	// 2. Eksekusi penghapusan
	if err := ic.DB.Delete(&ibuHamil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data ibu hamil"})
		return
	}

	// 3. Kembalikan response sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Ibu Hamil berhasil dihapus secara permanen",
	})
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