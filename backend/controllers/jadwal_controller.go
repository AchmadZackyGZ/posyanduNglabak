package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JadwalController struct {
	DB *gorm.DB
}

func NewJadwalController(db *gorm.DB) *JadwalController {
	return &JadwalController{DB: db}
}

// Struktur payload untuk membuat jadwal baru
type CreateJadwalInput struct {
	NamaKegiatan string `json:"nama_kegiatan" binding:"required"`
	Tanggal      string `json:"tanggal" binding:"required"` // Format: YYYY-MM-DD
	JamMulai     string `json:"jam_mulai" binding:"required"`
	JamSelesai   string `json:"jam_selesai" binding:"required"`
	Lokasi       string `json:"lokasi" binding:"required"`
}

// CreateJadwal menyimpan agenda Posyandu baru ke database
func (jc *JadwalController) CreateJadwal(c *gin.Context) {
	var input CreateJadwalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsing format tanggal dari string ke time.Time
	tglKegiatan, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	// Ambil ID petugas yang sedang login dari konteks JWT
	petugasID, _ := c.Get("userID")

	jadwal := models.JadwalKegiatan{
		NamaKegiatan: input.NamaKegiatan,
		Tanggal:      tglKegiatan,
		JamMulai:     input.JamMulai,
		JamSelesai:   input.JamSelesai,
		Lokasi:       input.Lokasi,
		Status:       "Akan Datang", // Default status saat baru dibuat
		CreatedByID:  petugasID.(string),
	}

	if err := jc.DB.Create(&jadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jadwal kegiatan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jadwal kegiatan berhasil ditambahkan",
		"data":    jadwal,
	})
}

// GetListJadwal mengambil seluruh agenda Posyandu
func (jc *JadwalController) GetListJadwal(c *gin.Context) {
	var jadwalList []models.JadwalKegiatan
	
	// Kita urutkan berdasarkan tanggal terdekat (Ascending)
	if err := jc.DB.Order("tanggal asc").Find(&jadwalList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar jadwal",
		"data":    jadwalList,
	})
}