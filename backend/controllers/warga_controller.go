package controllers

import (
	"net/http"
	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WargaController struct {
	DB *gorm.DB
}

func NewWargaController(db *gorm.DB) *WargaController {
	return &WargaController{DB: db}
}

// GetKMSAnak mengambil data balita dan riwayat timbang khusus untuk orang tua yang login
func (wc *WargaController) GetKMSAnak(c *gin.Context) {
	// Ambil ID User Warga dari token JWT yang sedang aktif
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}

	var balitas []models.Balita

	// Preload "Pemeriksaans" akan otomatis menempelkan riwayat KMS ke data balita
	// Kita urutkan tanggal_periksa ASC (dari terlama ke terbaru) agar grafiknya nanti maju ke depan
	if err := wc.DB.Preload("Pemeriksaans", func(db *gorm.DB) *gorm.DB {
		return db.Order("tanggal_periksa ASC")
	}).Where("user_id = ?", userID).Find(&balitas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data KMS anak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memuat data KMS",
		"data":    balitas,
	})
}