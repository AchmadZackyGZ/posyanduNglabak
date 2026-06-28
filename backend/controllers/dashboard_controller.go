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

// 1. TAMBAHKAN PROPERTI 'Link' AGAR TOMBOL DI SVELTE MENYALA
type NotifikasiCerdas struct {
	Tipe  string `json:"tipe"`  
	Pesan string `json:"pesan"`
	Link  string `json:"link,omitempty"` 
}

func (dc *DashboardController) GetSummary(c *gin.Context) {
	var totalBalita, balitaDitimbang, balitaDitimbangBulanLalu, ibuHamilAktif, totalLansia int64
	
	now := time.Now()
	bulanIni := int(now.Month())
	tahunIni := now.Year()
	waktuBulanLalu := now.AddDate(0, -1, 0)
	bulanLalu := int(waktuBulanLalu.Month())
	tahunLalu := waktuBulanLalu.Year()

	// ==========================================
	// 1. HITUNG KPI
	// ==========================================
	dc.DB.Model(&models.Balita{}).Count(&totalBalita)
	dc.DB.Model(&models.PemeriksaanBalita{}).Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", bulanIni, tahunIni).Count(&balitaDitimbang)
	dc.DB.Model(&models.PemeriksaanBalita{}).Where("EXTRACT(MONTH FROM tanggal_periksa) = ? AND EXTRACT(YEAR FROM tanggal_periksa) = ?", bulanLalu, tahunLalu).Count(&balitaDitimbangBulanLalu)
	dc.DB.Model(&models.IbuHamil{}).Where("status_kehamilan = ?", "AKTIF").Count(&ibuHamilAktif)
	dc.DB.Model(&models.Lansia{}).Where("status = ?", "AKTIF").Count(&totalLansia)

	// ==========================================
	// 2. GRAFIK PERTUMBUHAN (Ditambah COALESCE anti-error)
	// ==========================================
	var grafik []GrafikPertumbuhan
	enamBulanLalu := now.AddDate(0, -5, -now.Day()+1)

	dc.DB.Model(&models.PemeriksaanBalita{}).
		Select("TO_CHAR(tanggal_periksa, 'Mon YYYY') as bulan, COALESCE(ROUND(AVG(tinggi_badan)::numeric, 1), 0) as rata_tinggi, COALESCE(ROUND(AVG(berat_badan)::numeric, 1), 0) as rata_berat").
		Where("tanggal_periksa >= ?", enamBulanLalu).
		Group("TO_CHAR(tanggal_periksa, 'Mon YYYY'), EXTRACT(YEAR FROM tanggal_periksa), EXTRACT(MONTH FROM tanggal_periksa)").
		Order("EXTRACT(YEAR FROM tanggal_periksa) ASC, EXTRACT(MONTH FROM tanggal_periksa) ASC").
		Scan(&grafik)

	// ==========================================
	// 3. JADWAL TERDEKAT
	// ==========================================
	var jadwalTerdekat []JadwalSingkat
	awalHariIni := time.Date(tahunIni, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dc.DB.Model(&models.JadwalKegiatan{}).
		Select("id, nama_kegiatan, tanggal, jam_mulai, lokasi").
		Where("tanggal >= ? AND UPPER(status) = 'AKAN DATANG'", awalHariIni).
		Order("tanggal ASC").Limit(2).Scan(&jadwalTerdekat)

	// ==========================================
	// 4. NOTIFIKASI CERDAS (SEKARANG DIGABUNG!)
	// ==========================================
	var alerts []NotifikasiCerdas

	var bumilHPLDekat int64
	dc.DB.Model(&models.IbuHamil{}).Where("status_kehamilan = 'AKTIF' AND hpl BETWEEN ? AND ?", now, now.AddDate(0, 0, 30)).Count(&bumilHPLDekat)
	if bumilHPLDekat > 0 {
		alerts = append(alerts, NotifikasiCerdas{
			Tipe:  "WARNING",
			Pesan: fmt.Sprintf("Ada %d Ibu Hamil yang mendekati waktu persalinan (HPL) dalam 30 hari ke depan.", bumilHPLDekat),
			Link:  "/dashboard/ibu-hamil",
		})
	}

	type BalitaKritis struct {
		ID         string `json:"id"`
		NamaBalita string `json:"nama_balita"`
		StatusGizi string `json:"status_gizi"`
	}
	var listKritis []BalitaKritis

	// Cari balita Gizi Buruk/Kurang/Stunting
	dc.DB.Table("pemeriksaan_balita pb").
		Select("DISTINCT ON (pb.balita_id) pb.balita_id as id, b.nama_balita, pb.status_gizi").
		Joins("JOIN balita b ON b.id = pb.balita_id").
		Where("UPPER(pb.status_gizi) LIKE '%BURUK%' OR UPPER(pb.status_gizi) LIKE '%KURANG%' OR UPPER(pb.status_gizi) LIKE '%STUNTING%'").
		Order("pb.balita_id, pb.tanggal_periksa DESC").
		Scan(&listKritis)

	// Masukkan nama anak spesifik ke dalam array alerts yang sama
	for _, b := range listKritis {
		alerts = append(alerts, NotifikasiCerdas{
			Tipe:  "DANGER",
			Pesan: "Terindikasi " + b.StatusGizi + ": Balita atas nama " + b.NamaBalita,
			Link:  "/dashboard/balita?cari=" + b.NamaBalita, // Link ajaib menuju tabel
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
			"notifikasi":         alerts, // Hanya ada SATU wadah notifikasi sekarang!
		},
	})
}