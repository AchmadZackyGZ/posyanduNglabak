package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models"

	"github.com/gin-gonic/gin"
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

// --- DTO Update ---
type UpdatePemeriksaanIbuHamilInput struct {
	UsiaKehamilan int     `json:"usia_kehamilan"`
	TekananDarah  string  `json:"tekanan_darah"`
	BeratBadan    float64 `json:"berat_badan"`
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
	var results []models.IbuHamil

	// FIX: Jauh lebih bersih! Tidak perlu JOIN lagi karena NoHP sudah ada di tabel IbuHamil
	if err := ic.DB.Order("created_at desc").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ibu hamil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar ibu hamil",
		"data":    results,
	})
}

// RegisterIbuHamil mendaftarkan rekam medis ibu hamil MURNI tanpa membuat akun user
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

	// FIX: Simpan entitas data klinis Ibu Hamil SECARA LANGSUNG
	// Tidak ada lagi logika pembuatan akun (User) siluman di sini!
	ibuHamil := models.IbuHamil{
		NIK:             input.NIK,
		NamaIbu:         input.NamaIbu,
		NoHP:            input.NoHP, // Disimpan langsung ke kolom baru
		TanggalLahir:    tglLahir,
		HPL:             hpl,
		Alamat:          input.Alamat,
		StatusKehamilan: input.StatusKehamilan,
	}

	if err := ic.DB.Create(&ibuHamil).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftar. NIK Ibu Hamil mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Ibu Hamil berhasil dicatat",
		"data": gin.H{
			"ibu_hamil_id": ibuHamil.ID,
			"nama":         ibuHamil.NamaIbu,
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

func (ihc *IbuHamilController) UpdatePemeriksaan(c *gin.Context) {
	id := c.Param("id")
	var input UpdatePemeriksaanIbuHamilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	var pemeriksaan models.PemeriksaanIbuHamil
	if err := ihc.DB.Where("id = ?", id).First(&pemeriksaan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data rekam medis tidak ditemukan"})
		return
	}

	pemeriksaan.UsiaKehamilan = input.UsiaKehamilan
	pemeriksaan.TekananDarah = input.TekananDarah
	pemeriksaan.BeratBadan = input.BeratBadan
	pemeriksaan.Catatan = input.Catatan

	ihc.DB.Save(&pemeriksaan)
	c.JSON(http.StatusOK, gin.H{"message": "Rekam medis ibu hamil berhasil diperbarui"})
}

func (ihc *IbuHamilController) DeletePemeriksaan(c *gin.Context) {
	id := c.Param("id")
	if err := ihc.DB.Where("id = ?", id).Delete(&models.PemeriksaanIbuHamil{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus rekam medis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rekam medis ibu hamil berhasil dihapus"})
}