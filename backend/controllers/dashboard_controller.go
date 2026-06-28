package controllers

import (
	"fmt"
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

// --- DTO & Structs ---
type GrafikPertumbuhan struct {
	Bulan      string  `json:"bulan"`
	RataTinggi float64 `json:"rata_tinggi"`
	RataBerat  float64 `json:"rata_berat"`
}

type JadwalSingkat struct {
	ID           string    `json:"id"`
	NamaKegiatan string    `json:"nama_kegiatan"`
	Tanggal      time.Time `json:"tanggal"`
	JamMulai     string    `json:"jam_mulai"`
	Lokasi       string    `json:"lokasi"`
}

type NotifikasiCerdas struct {
	Tipe  string `json:"tipe"`  // "INFO", "WARNING", "DANGER"
	Pesan string `json:"pesan"`
}

// GetSummary mengambil total agregasi metrik untuk beranda dashboard
func (dc *DashboardController) GetSummary(c *gin.Context) {
	var totalBalita, balitaDitimbang, balitaDitimbangBulanLalu, ibuHamilAktif, totalLansia int64
	
	now := time.Now()
	bulanIni := int(now.Month())
	tahunIni := now.Year()

	// Mundur 1 bulan untuk komparasi data (MoM)
	waktuBulanLalu := now.AddDate(0, -1, 0)
	bulanLalu := int(waktuBulanLalu.Month())
	tahunLalu := waktuBulanLalu.Year()

	// ==========================================
	// 1. HITUNG KPI (Key Performance Indicators)
	// ==========================================
	
	dc.DB.Model(&models.Balita{}).Count(&totalBalita)
	
	dc.DB.Model(&models.PemeriksaanBalita{}).
		Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", bulanIni, tahunIni).
		Count(&balitaDitimbang)

	dc.DB.Model(&models.PemeriksaanBalita{}).
		Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", bulanLalu, tahunLalu).
		Count(&balitaDitimbangBulanLalu)

	dc.DB.Model(&models.IbuHamil{}).Where("status_kehamilan = ?", "AKTIF").Count(&ibuHamilAktif)
	dc.DB.Model(&models.Lansia{}).Where("status = ?", "AKTIF").Count(&totalLansia)

	// ==========================================
	// 2. GRAFIK PERTUMBUHAN BALITA (6 BULAN)
	// ==========================================
	
	var grafik []GrafikPertumbuhan
	enamBulanLalu := now.AddDate(0, -5, -now.Day()+1) // Tarik mundur 5 bulan dari awal bulan ini

	dc.DB.Model(&models.PemeriksaanBalita{}).
		Select("TO_CHAR(tanggal_periksa, 'Mon YYYY') as bulan, ROUND(AVG(tinggi_badan)::numeric, 1) as rata_tinggi, ROUND(AVG(berat_badan)::numeric, 1) as rata_berat").
		Where("tanggal_periksa >= ?", enamBulanLalu).
		Group("TO_CHAR(tanggal_periksa, 'Mon YYYY'), EXTRACT(YEAR FROM tanggal_periksa), EXTRACT(MONTH FROM tanggal_periksa)").
		Order("EXTRACT(YEAR FROM tanggal_periksa) ASC, EXTRACT(MONTH FROM tanggal_periksa) ASC").
		Scan(&grafik)

	// ==========================================
	// 3. JADWAL KEGIATAN TERDEKAT
	// ==========================================
	
	var jadwalTerdekat []JadwalSingkat
	awalHariIni := time.Date(tahunIni, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	dc.DB.Model(&models.JadwalKegiatan{}).
		Select("id, nama_kegiatan, tanggal, jam_mulai, lokasi").
		Where("tanggal >= ? AND status = 'AKAN DATANG'", awalHariIni).
		Order("tanggal ASC").
		Limit(2).
		Scan(&jadwalTerdekat)

	// ==========================================
	// 4. NOTIFIKASI CERDAS (ALERTS)
	// ==========================================
	
	var alerts []NotifikasiCerdas

	// Cek Bumil Mendekati HPL (Dalam 30 Hari)
	var bumilHPLDekat int64
	dc.DB.Model(&models.IbuHamil{}).
		Where("status_kehamilan = 'AKTIF' AND hpl BETWEEN ? AND ?", now, now.AddDate(0, 0, 30)).
		Count(&bumilHPLDekat)
	
	if bumilHPLDekat > 0 {
		alerts = append(alerts, NotifikasiCerdas{
			Tipe:  "WARNING",
			Pesan: fmt.Sprintf("Ada %d Ibu Hamil yang mendekati waktu persalinan (HPL) dalam 30 hari ke depan.", bumilHPLDekat),
		})
	}

	// Cek Balita Gizi Kurang/Buruk Bulan Ini
	var giziBermasalah int64
	dc.DB.Model(&models.PemeriksaanBalita{}).
		Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", bulanIni, tahunIni).
		Where("status_gizi IN ?", []string{"GIZI KURANG", "GIZI BURUK", "STUNTING"}).
		Count(&giziBermasalah)
	
	if giziBermasalah > 0 {
		alerts = append(alerts, NotifikasiCerdas{
			Tipe:  "DANGER",
			Pesan: fmt.Sprintf("Perhatian: Terdapat %d balita dengan indikasi gizi kurang/buruk pada pemeriksaan bulan ini.", giziBermasalah),
		})
	}

	// ==========================================
	// 5. RESPONSE BUILDER
	// ==========================================
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil ringkasan dashboard",
		"data": gin.H{
			"kpi": gin.H{
				"total_balita":                 totalBalita,
				"balita_ditimbang_bulan_ini":   balitaDitimbang,
				"balita_ditimbang_bulan_lalu":  balitaDitimbangBulanLalu,
				"ibu_hamil_aktif":              ibuHamilAktif,
				"total_lansia":                 totalLansia,
			},
			"grafik_pertumbuhan": grafik,
			"jadwal_terdekat":    jadwalTerdekat,
			"notifikasi":         alerts,
		},
	})
}