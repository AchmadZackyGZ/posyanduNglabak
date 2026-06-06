package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models" // Sesuaikan dengan nama modul Anda

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BalitaController struct {
	DB *gorm.DB
}

func NewBalitaController(db *gorm.DB) *BalitaController {
	return &BalitaController{DB: db}
}

// Struktur input pencatatan timbang (FR-03)
type TimbangBalitaInput struct {
	BalitaID      string  `json:"balita_id" binding:"required"`
	BeratBadan    float64 `json:"berat_badan" binding:"required"`
	TinggiBadan   float64 `json:"tinggi_badan" binding:"required"`
	LingkarKepala float64 `json:"lingkar_kepala"`
	StatusGizi    string  `json:"status_gizi" binding:"required"` // NORMAL / STUNTING / dll
	Catatan       string  `json:"catatan"`
}

// Struktur input untuk pendaftaran Balita (FR-02)
type RegisterBalitaInput struct {
	NIK          string `json:"nik" binding:"required"`
	NamaBalita   string `json:"nama_balita" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
	JenisKelamin string `json:"jenis_kelamin" binding:"required"` // L / P
	NamaOrangTua string `json:"nama_orang_tua" binding:"required"`
	Alamat       string `json:"alamat" binding:"required"`
	NoHP         string `json:"no_hp" binding:"required"` // Digunakan sebagai username ortu
}

// Struktur input untuk Update Balita (Tanpa NoHP karena NoHP terikat di tabel User)
type UpdateBalitaInput struct {
	NIK          string `json:"nik" binding:"required"`
	NamaBalita   string `json:"nama_balita" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required"`
	NamaOrangTua string `json:"nama_orang_tua" binding:"required"`
	Alamat       string `json:"alamat" binding:"required"`
}

// GetListBalita mengambil seluruh data balita secara menurun (terbaru di atas)
func (bc *BalitaController) GetListBalita(c *gin.Context) {
	var balitas []models.Balita
	if err := bc.DB.Order("created_at desc").Find(&balitas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data balita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar balita",
		"data":    balitas,
	})
}

// RegisterBalita mendaftarkan pasien sekaligus membuatkan akun USER untuk orang tua
func (bc *BalitaController) RegisterBalita(c *gin.Context) {
	var input RegisterBalitaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal lahir harus YYYY-MM-DD"})
		return
	}

	// Gunakan transaksi DB agar jika salah satu gagal, semuanya di-rollback
	tx := bc.DB.Begin()

	// 1. Cari atau buat akun Orang Tua (Role: USER) menggunakan NoHP sebagai username
	var ortu models.User
	if err := tx.Where("username = ?", input.NoHP).First(&ortu).Error; err != nil {
		// Jika belum ada, buat akun baru dengan password default "posyandu123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("posyandu123"), bcrypt.DefaultCost)
		ortu = models.User{
			Username:     input.NoHP,
			PasswordHash: string(hashedPassword),
			NamaLengkap:  input.NamaOrangTua,
			Role:         "USER",
			IsActive:     true,
		}
		if err := tx.Create(&ortu).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun orang tua"})
			return
		}
	}

	// 2. Simpan Data Balita dengan mengaitkannya ke UserID orang tua
	balita := models.Balita{
		UserID:       ortu.ID,
		NIK:          input.NIK,
		NamaBalita:   input.NamaBalita,
		TanggalLahir: tglLahir,
		JenisKelamin: input.JenisKelamin,
		NamaOrangTua: input.NamaOrangTua,
		Alamat:       input.Alamat,
	}

	if err := tx.Create(&balita).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftarkan balita. NIK mungkin sudah terdaftar."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Balita berhasil",
		"data": gin.H{
			"balita_id": balita.ID,
			"nama":      balita.NamaBalita,
			"ortu_akun": ortu.Username,
		},
	})
}

func (bc *BalitaController) UpdateBalita(c *gin.Context){
	id := c.Param("id")
	var balita models.Balita

	// 1. Cek apakah balita dengan ID tersebut ada di database
	if err := bc.DB.First(&balita, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data balita tidak ditemukan"})
		return
	}

	// 2. Tangkap data JSON yang dikirim dari Frontend
	var input UpdateBalitaInput 
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Konversi format tanggal lahir
	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal lahir harus YYYY-MM-DD"})
		return
	}

	// 4. Timpa data lama dengan data baru
	balita.NIK = input.NIK
	balita.NamaBalita = input.NamaBalita
	balita.TanggalLahir = tglLahir
	balita.JenisKelamin = input.JenisKelamin
	balita.NamaOrangTua = input.NamaOrangTua
	balita.Alamat = input.Alamat

	// 5. Simpan perubahan ke database
	if err := bc.DB.Save(&balita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data balita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data balita berhasil diperbarui",
		"data":    balita,
	})
}

func (bc *BalitaController) DeleteBalita(c *gin.Context){
	id := c.Param("id")
	var balita models.Balita

	// 1. Pastikan data yang mau dihapus itu ada
	if err := bc.DB.First(&balita, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data balita tidak ditemukan"})
		return
	}

	// 2. Eksekusi penghapusan dari database
	if err := bc.DB.Delete(&balita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data balita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data balita berhasil dihapus secara permanen",
	})
}


// CatatPemeriksaan menyimpan hasil ukur bulanan oleh Kader/Bidan
func (bc *BalitaController) CatatPemeriksaan(c *gin.Context) {
	var input TimbangBalitaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID petugas (Kader/Bidan) yang sedang login dari Context Middleware
	kaderID, _ := c.Get("userID")

	pemeriksaan := models.PemeriksaanBalita{
		BalitaID:       input.BalitaID,
		TanggalPeriksa: time.Now(),
		BeratBadan:     input.BeratBadan,
		TinggiBadan:    input.TinggiBadan,
		LingkarKepala:  input.LingkarKepala,
		StatusGizi:     input.StatusGizi,
		Catatan:        input.Catatan,
		DiperiksaOleh:  kaderID.(string),
	}

	if err := bc.DB.Create(&pemeriksaan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rekam medis balita"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hasil pemeriksaan berhasil dicatat",
		"data":    pemeriksaan,
	})
}