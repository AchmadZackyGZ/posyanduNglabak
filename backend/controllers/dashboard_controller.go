package controllers

import (
	"net/http"
	"posyandu-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB *gorm.DB
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{DB: db}
}

// GetSummary mengambil total agregasi metrik untuk beranda dashboard
func (dc *DashboardController) GetSummary(c *gin.Context) {
	var totalBalita, balitaDitimbang, ibuHamilAktif, totalLansia int64
	now := time.Now()

	// 1. Hitung Total Balita
	dc.DB.Model(&models.Balita{}).Count(&totalBalita)

	// 2. Hitung Balita Ditimbang Bulan Ini (Menggunakan ekstrak tanggal native PostgreSQL)
	dc.DB.Model(&models.PemeriksaanBalita{}).
		Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", int(now.Month()), now.Year()).
		Count(&balitaDitimbang)

	// 3. Hitung Ibu Hamil Aktif
	dc.DB.Model(&models.IbuHamil{}).Where("status_kehamilan = ?", "AKTIF").Count(&ibuHamilAktif)

	// 4. Hitung Lansia Aktif
	dc.DB.Model(&models.Lansia{}).Where("status = ?", "AKTIF").Count(&totalLansia)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil ringkasan dashboard",
		"data": gin.H{
			"total_balita":     totalBalita,
			"balita_ditimbang": balitaDitimbang,
			"ibu_hamil_aktif":  ibuHamilAktif,
			"total_lansia":     totalLansia,
		},
	})
}