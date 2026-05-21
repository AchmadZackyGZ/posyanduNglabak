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

type GrafikPertumbuhan struct {
	Bulan      string  `json:"bulan"`
	RataTinggi float64 `json:"rata_tinggi"`
	RataBerat  float64 `json:"rata_berat"`
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

	// 5. Ambil data grafik 6 bulan terakhir
	var grafik []GrafikPertumbuhan

	// Menghitung tanggal awal dari 5 bulan ke belakang (jadi total 6 bulan terhitung dengan bulan ini)
	// Misalnya sekarang Mei tanggal 21 -> Ditarik mundur ke tanggal 1 Desember.
	enamBulanLalu := now.AddDate(0, -5, -now.Day()+1)

	// GORM Query untuk agregasi per bulan di PostgreSQL
	dc.DB.Model(&models.PemeriksaanBalita{}).
		Select("TO_CHAR(tanggal_periksa, 'Mon YYYY') as bulan, ROUND(AVG(tinggi_badan)::numeric, 1) as rata_tinggi, ROUND(AVG(berat_badan)::numeric, 1) as rata_berat").
		Where("tanggal_periksa >= ?", enamBulanLalu).
		Group("TO_CHAR(tanggal_periksa, 'Mon YYYY'), EXTRACT(YEAR FROM tanggal_periksa), EXTRACT(MONTH FROM tanggal_periksa)").
		Order("EXTRACT(YEAR FROM tanggal_periksa) ASC, EXTRACT(MONTH FROM tanggal_periksa) ASC").
		Scan(&grafik)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil ringkasan dashboard",
		"data": gin.H{
			"total_balita":     totalBalita,
			"balita_ditimbang": balitaDitimbang,
			"ibu_hamil_aktif":  ibuHamilAktif,
			"total_lansia":     totalLansia,
			"grafik_pertumbuhan": grafik,
		},
	})
}