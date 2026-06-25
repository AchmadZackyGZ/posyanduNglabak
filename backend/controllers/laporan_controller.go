package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"posyandu-api/config"
	"posyandu-api/models"
)

// --- DTO Responses ---

type RekapBalitaSummary struct {
	TotalBalita     int64 `json:"total_balita"`
	BalitaDitimbang int64 `json:"balita_ditimbang"`
	BalitaNaikBB    int64 `json:"balita_naik_bb"`
	BalitaTurunBB   int64 `json:"balita_turun_bb"`
}

type DetailLaporanBalita struct {
	No         int     `json:"no"`
	NamaBalita string  `json:"nama_balita"`
	Usia       string  `json:"usia"`
	JK         string  `json:"jk"`
	BeratKG    float64 `json:"berat_kg"`
	TinggiCM   float64 `json:"tinggi_cm"`
	StatusGizi string  `json:"status_gizi"`
}

type LaporanBalitaResponse struct {
	Summary RekapBalitaSummary    `json:"summary"`
	Data    []DetailLaporanBalita `json:"data"`
}

type RekapIbuHamilSummary struct {
	TotalIbuHamil     int64 `json:"total_ibu_hamil"`
	IbuHamilDiperiksa int64 `json:"ibu_hamil_diperiksa"`
}

type DetailLaporanIbuHamil struct {
	No             int     `json:"no"`
	NamaIbu        string  `json:"nama_ibu"`
	UsiaKandungan  int     `json:"usia_kandungan"`
	TekananDarah   string  `json:"tekanan_darah"`
	BeratKG        float64 `json:"berat_kg"`
	Catatan        string  `json:"catatan"`
}

type LaporanIbuHamilResponse struct {
	Summary RekapIbuHamilSummary    `json:"summary"`
	Data    []DetailLaporanIbuHamil `json:"data"`
}

type RekapLansiaSummary struct {
	TotalLansia     int64 `json:"total_lansia"`
	LansiaDiperiksa int64 `json:"lansia_diperiksa"`
}

type DetailLaporanLansia struct {
	No           int     `json:"no"`
	NamaLansia   string  `json:"nama_lansia"`
	Usia         string  `json:"usia"`
	TekananDarah string  `json:"tekanan_darah"`
	GulaDarah    float64 `json:"gula_darah"`
	Kolesterol   float64 `json:"kolesterol"`
	Catatan      string  `json:"catatan"`
}

type LaporanLansiaResponse struct {
	Summary RekapLansiaSummary    `json:"summary"`
	Data    []DetailLaporanLansia `json:"data"`
}

// --- Helper Functions ---

// HitungUsia menghitung selisih waktu menjadi format "X th Y bln"
func HitungUsia(tanggalLahir, tanggalPeriksa time.Time) string {
	tahun := tanggalPeriksa.Year() - tanggalLahir.Year()
	bulan := int(tanggalPeriksa.Month()) - int(tanggalLahir.Month())

	if bulan < 0 {
		tahun--
		bulan += 12
	}

	if tahun > 0 && bulan > 0 {
		return fmt.Sprintf("%d th %d bln", tahun, bulan)
	} else if tahun > 0 {
		return fmt.Sprintf("%d th", tahun)
	} else {
		return fmt.Sprintf("%d bln", bulan)
	}
}

// --- Controller Main Function ---
// laporan balita

func GetLaporanBalita(c *gin.Context) {
	periode := c.Query("periode") // Expected format: "2024-05"
	if periode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter periode (YYYY-MM) wajib diisi"})
		return
	}

	db := config.DB // Asumsikan Anda menggunakan global DB instance

	var response LaporanBalitaResponse

	// 1. Hitung Total Balita Keseluruhan
	db.Model(&models.Balita{}).Count(&response.Summary.TotalBalita)

	// 2. Ambil Data Pemeriksaan Bulan Ini beserta Data Balitanya
	var pemeriksaanBulanIni []models.PemeriksaanBalita
	// Menggunakan TO_CHAR PostgreSQL untuk filter "YYYY-MM"
	err := db.Preload("Balita").
		Where("TO_CHAR(tanggal_periksa, 'YYYY-MM') = ?", periode).
		Order("tanggal_periksa ASC").
		Find(&pemeriksaanBulanIni).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pemeriksaan"})
		return
	}

	response.Summary.BalitaDitimbang = int64(len(pemeriksaanBulanIni))

	// 3. Kalkulasi Bulan Sebelumnya untuk Komparasi Berat Badan
	parsedPeriode, _ := time.Parse("2006-01", periode)
	bulanLalu := parsedPeriode.AddDate(0, -1, 0).Format("2006-01")

	// Ambil semua pemeriksaan bulan lalu untuk mapping (menghindari query N+1 di dalam loop)
	var pemeriksaanBulanLalu []models.PemeriksaanBalita
	db.Where("TO_CHAR(tanggal_periksa, 'YYYY-MM') = ?", bulanLalu).Find(&pemeriksaanBulanLalu)

	// Buat map untuk pencarian cepat berat badan bulan lalu berdasarkan BalitaID
	mapBeratBulanLalu := make(map[string]float64)
	for _, p := range pemeriksaanBulanLalu {
		mapBeratBulanLalu[p.BalitaID] = p.BeratBadan // Asumsi p.BalitaID bertipe string (UUID)
	}

	// 4. Proses Array Data dan Kalkulasi Naik/Turun
	for i, p := range pemeriksaanBulanIni {
		// Hitung Usia
		usiaStr := HitungUsia(p.Balita.TanggalLahir, p.TanggalPeriksa)

		// Cek Komparasi Berat Badan
		if beratLalu, ada := mapBeratBulanLalu[p.BalitaID]; ada {
			if p.BeratBadan > beratLalu {
				response.Summary.BalitaNaikBB++
			} else if p.BeratBadan < beratLalu {
				response.Summary.BalitaTurunBB++
			}
			// Jika sama, tidak masuk naik atau turun
		}

		// Masukkan ke array detail data
		detail := DetailLaporanBalita{
			No:         i + 1,
			NamaBalita: p.Balita.NamaBalita,
			Usia:       usiaStr,
			JK:         p.Balita.JenisKelamin,
			BeratKG:    p.BeratBadan,
			TinggiCM:   p.TinggiBadan,
			StatusGizi: p.StatusGizi,
		}
		response.Data = append(response.Data, detail)
	}

	// 5. Kirim Respons
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil rekap laporan balita",
		"data":    response,
	})
}

// laporan ibu hamil

func GetLaporanIbuHamil(c *gin.Context) {
	periode := c.Query("periode")
	if periode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter periode (YYYY-MM) wajib diisi"})
		return
	}

	db := config.DB
	var response LaporanIbuHamilResponse

	// 1. Hitung Total Ibu Hamil Keseluruhan
	db.Model(&models.IbuHamil{}).Count(&response.Summary.TotalIbuHamil)

	// 2. Ambil Data Pemeriksaan Bulan Ini beserta Data Ibu Hamilnya
	var pemeriksaanBulanIni []models.PemeriksaanIbuHamil
	err := db.Preload("IbuHamil").
		Where("TO_CHAR(tanggal_periksa, 'YYYY-MM') = ?", periode).
		Order("tanggal_periksa ASC").
		Find(&pemeriksaanBulanIni).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pemeriksaan ibu hamil"})
		return
	}

	response.Summary.IbuHamilDiperiksa = int64(len(pemeriksaanBulanIni))

	// 3. Masukkan ke array detail data
	for i, p := range pemeriksaanBulanIni {
		detail := DetailLaporanIbuHamil{
			No:            i + 1,
			NamaIbu:       p.IbuHamil.NamaIbu,
			UsiaKandungan: p.UsiaKehamilan,
			TekananDarah:  p.TekananDarah,
			BeratKG:       p.BeratBadan,
			Catatan:       p.Catatan,
		}
		response.Data = append(response.Data, detail)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil rekap laporan ibu hamil",
		"data":    response,
	})
}

// laporan lansia

func GetLaporanLansia(c *gin.Context) {
	periode := c.Query("periode")
	if periode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter periode (YYYY-MM) wajib diisi"})
		return
	}

	db := config.DB
	var response LaporanLansiaResponse

	// 1. Hitung Total Lansia Keseluruhan
	db.Model(&models.Lansia{}).Count(&response.Summary.TotalLansia)

	// 2. Ambil Data Pemeriksaan Bulan Ini beserta Data Lansianya
	var pemeriksaanBulanIni []models.PemeriksaanLansia
	err := db.Preload("Lansia").
		Where("TO_CHAR(tanggal_periksa, 'YYYY-MM') = ?", periode).
		Order("tanggal_periksa ASC").
		Find(&pemeriksaanBulanIni).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pemeriksaan lansia"})
		return
	}

	response.Summary.LansiaDiperiksa = int64(len(pemeriksaanBulanIni))

	// 3. Masukkan ke array detail data
	for i, p := range pemeriksaanBulanIni {
		usiaStr := HitungUsia(p.Lansia.TanggalLahir, p.TanggalPeriksa)

		detail := DetailLaporanLansia{
			No:           i + 1,
			NamaLansia:   p.Lansia.NamaLengkap,
			Usia:         usiaStr,
			TekananDarah: p.TekananDarah,
			GulaDarah:    p.GulaDarah,
			Kolesterol:   p.Kolesterol,
			Catatan:      p.Catatan,
		}
		response.Data = append(response.Data, detail)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil rekap laporan lansia",
		"data":    response,
	})
}