package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LansiaController struct {
	DB *gorm.DB
}

func NewLansiaController(db *gorm.DB) *LansiaController {
	return &LansiaController{DB: db}
}

// Struktur payload untuk registrasi Lansia baru (FR-05)
type RegisterLansiaInput struct {
	NIK            string `json:"nik" binding:"required"`
	NamaLengkap    string `json:"nama_lengkap" binding:"required"`
	TanggalLahir   string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
	JenisKelamin   string `json:"jenis_kelamin" binding:"required"` // L / P
	Alamat         string `json:"alamat" binding:"required"`
	NamaPendamping string `json:"nama_pendamping"`
	KontakDarurat  string `json:"kontak_darurat"`
	Status         string `json:"status"` // Opsional: Untuk mengubah status (misal: AKTIF/PINDAH/MENINGGAL)
}

// Struktur input untuk Update Lansia
type UpdateLansiaInput struct {
	NIK            string `json:"nik" binding:"required"`
	NamaLengkap    string `json:"nama_lengkap" binding:"required"`
	TanggalLahir   string `json:"tanggal_lahir" binding:"required"` // Format: YYYY-MM-DD
	JenisKelamin   string `json:"jenis_kelamin" binding:"required"` // L / P
	Alamat         string `json:"alamat" binding:"required"`
	NamaPendamping string `json:"nama_pendamping"`
	KontakDarurat  string `json:"kontak_darurat"`
	Status         string `json:"status"` // Opsional: Untuk mengubah status (misal: AKTIF/PINDAH/MENINGGAL)
}

// Struktur payload pencatatan rekam medis klinis Lansia
type PeriksaLansiaInput struct {
	LansiaID     string  `json:"lansia_id" binding:"required"`
	TekananDarah string  `json:"tekanan_darah" binding:"required"` // e.g. "130/85"
	GulaDarah    float64 `json:"gula_darah"`                       // mg/dL
	Kolesterol   float64 `json:"kolesterol"`                       // mg/dL
	Catatan      string  `json:"catatan"`
}

// GetListLansia mengambil seluruh data lansia
func (lc *LansiaController) GetListLansia(c *gin.Context) {
	var lansias []models.Lansia
	if err := lc.DB.Order("created_at desc").Find(&lansias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar lansia"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar lansia",
		"data":    lansias,
	})
}

// RegisterLansia memproses pendaftaran data master Lansia di tingkat backend
func (lc *LansiaController) RegisterLansia(c *gin.Context) {
	var input RegisterLansiaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal lahir harus YYYY-MM-DD"})
		return
	}

	lansia := models.Lansia{
		NIK:            input.NIK,
		NamaLengkap:    input.NamaLengkap,
		TanggalLahir:   tglLahir,
		JenisKelamin:   input.JenisKelamin,
		Alamat:         input.Alamat,
		NamaPendamping: input.NamaPendamping,
		KontakDarurat:  input.KontakDarurat,
		Status:         input.Status,
	}

	// Penyimpanan langsung tanpa pembuatan user mandiri (Lansia dikelola penuh oleh Kader)
	if err := lc.DB.Create(&lansia).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftarkan Lansia. NIK mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Lansia berhasil dicatat",
		"data":    lansia,
	})
}

// UpdateLansia mengubah data profil lansia berdasarkan ID
func (lc *LansiaController) UpdateLansia(c *gin.Context) {
	id := c.Param("id")
	var lansia models.Lansia

	// 1. Cek apakah data lansia dengan ID tersebut ada di database
	if err := lc.DB.First(&lansia, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data lansia tidak ditemukan"})
		return
	}

	// 2. Bind data JSON yang dikirim dari Frontend
	var input UpdateLansiaInput
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
	lansia.NIK = input.NIK
	lansia.NamaLengkap = input.NamaLengkap
	lansia.TanggalLahir = tglLahir
	lansia.JenisKelamin = input.JenisKelamin
	lansia.Alamat = input.Alamat
	lansia.NamaPendamping = input.NamaPendamping
	lansia.KontakDarurat = input.KontakDarurat
	
	// Update status jika dikirim dari frontend
	if input.Status != "" {
		lansia.Status = input.Status
	}

	// 5. Simpan perubahan ke database
	if err := lc.DB.Save(&lansia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data lansia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data lansia berhasil diperbarui",
		"data":    lansia,
	})
}

// DeleteLansia menghapus data lansia berdasarkan ID
func (lc *LansiaController) DeleteLansia(c *gin.Context) {
	id := c.Param("id")
	var lansia models.Lansia

	// 1. Pastikan data yang mau dihapus itu ada
	if err := lc.DB.First(&lansia, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data lansia tidak ditemukan"})
		return
	}

	// 2. Eksekusi penghapusan dari database
	if err := lc.DB.Delete(&lansia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data lansia"})
		return
	}

	// 3. Kembalikan response sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Data lansia berhasil dihapus secara permanen",
	})
}


// CatatPemeriksaan mengamankan entri rekam medis bulanan beserta ID pemeriksa
func (lc *LansiaController) CatatPemeriksaan(c *gin.Context) {
	var input PeriksaLansiaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ekstrak klaim identitas dari middleware otentikasi
	petugasID, _ := c.Get("userID")

	pemeriksaan := models.PemeriksaanLansia{
		LansiaID:       input.LansiaID,
		TanggalPeriksa: time.Now(),
		TekananDarah:   input.TekananDarah,
		GulaDarah:      input.GulaDarah,
		Kolesterol:     input.Kolesterol,
		Catatan:        input.Catatan,
		DiperiksaOleh:  petugasID.(string),
	}

	if err := lc.DB.Create(&pemeriksaan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rekam medis Lansia"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pemeriksaan klinis Lansia berhasil disimpan",
		"data":    pemeriksaan,
	})
}