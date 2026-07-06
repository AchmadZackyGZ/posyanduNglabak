package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models" // Sesuaikan dengan nama modul Anda

	"github.com/gin-gonic/gin"
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
	TanggalPeriksa string  `json:"tanggal_periksa" binding:"required"`
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
	UserID       *string `json:"user_id"` // <--- TAMBAHKAN INI (Gunakan *string agar bisa menerima null)
}

// --- DTO Update ---
type UpdatePemeriksaanBalitaInput struct {
	TanggalPeriksa string  `json:"tanggal_periksa"`
	BeratBadan    float64 `json:"berat_badan"`
	TinggiBadan   float64 `json:"tinggi_badan"`
	LingkarKepala float64 `json:"lingkar_kepala"`
	StatusGizi    string  `json:"status_gizi"`
	Catatan       string  `json:"catatan"`
}

// Struktur input untuk Update Balita (Tanpa NoHP karena NoHP terikat di tabel User)
type UpdateBalitaInput struct {
	NIK          string `json:"nik" binding:"required"`
	NamaBalita   string `json:"nama_balita" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required"`
	NamaOrangTua string `json:"nama_orang_tua" binding:"required"`
	Alamat       string `json:"alamat" binding:"required"`
	UserID       *string `json:"user_id"` // <--- TAMBAHKAN INI JUGA
}

// Struktur input pencatatan imunisasi
type ImunisasiInput struct {
	BalitaID         string `json:"balita_id" binding:"required"`
	JenisImunisasi   string `json:"jenis_imunisasi" binding:"required"`
	TanggalImunisasi string `json:"tanggal_imunisasi" binding:"required"`
	Catatan          string `json:"catatan"`
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

// RegisterBalita mendaftarkan pasien MURNI tanpa membuat akun USER untuk orang tua
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

	// FIX: Simpan Data Balita SECARA LANGSUNG tanpa transaksi pembuatan User
	// (Kolom UserID akan bernilai NULL sesuai perubahan di models.go sebelumnya)
	balita := models.Balita{
		NIK:          input.NIK,
		NamaBalita:   input.NamaBalita,
		TanggalLahir: tglLahir,
		JenisKelamin: input.JenisKelamin,
		NamaOrangTua: input.NamaOrangTua,
		Alamat:       input.Alamat,
		// NoHP saat ini tidak ada di struct models.Balita, jadi kita abaikan dulu dari input
	}

	if err := bc.DB.Create(&balita).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftarkan balita. NIK mungkin sudah terdaftar."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Balita berhasil",
		"data": gin.H{
			"balita_id": balita.ID,
			"nama":      balita.NamaBalita,
		},
	})
}

// GetRiwayatTimbang mengambil seluruh data rekam medis penimbangan balita
func (bc *BalitaController) GetRiwayatTimbang(c *gin.Context) {
	var riwayat []models.PemeriksaanBalita
	
	// Preload("Pemeriksa") berfungsi seperti JOIN untuk menarik nama Bidan/Kader dari tabel User
	query := bc.DB.Preload("Pemeriksa").Order("tanggal_periksa desc")

	// Fitur Canggih: Jika ada query ?balita_id=xxx, filter hanya untuk 1 balita itu saja
	if balitaID := c.Query("balita_id"); balitaID != "" {
		query = query.Where("balita_id = ?", balitaID)
	}

	if err := query.Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat pemeriksaan balita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat timbang",
		"data":    riwayat,
	})
}

// 2. Ambil Riwayat Imunisasi (GET)
func (bc *BalitaController) GetRiwayatImunisasi(c *gin.Context) {
	var imunisasi []models.ImunisasiBalita
	
	// Gunakan Preload untuk mengambil relasi nama Balita dan Pencatat
	if err := bc.DB.Preload("Balita").Preload("Pencatat").Order("tanggal_imunisasi DESC").Find(&imunisasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat imunisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": imunisasi})
}

// 1. Catat Imunisasi (POST)
func (bc *BalitaController) CatatImunisasi(c *gin.Context) {
	var input ImunisasiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID User yang sedang login dari middleware JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diotorisasi"})
		return
	}

	tanggalParsed, _ := time.Parse("2006-01-02", input.TanggalImunisasi)

	imunisasi := models.ImunisasiBalita{
		BalitaID:         input.BalitaID,
		JenisImunisasi:   input.JenisImunisasi,
		TanggalImunisasi: tanggalParsed,
		Catatan:          input.Catatan,
		DicatatOleh:      userID.(string), // ID Kader/Bidan
	}

	if err := bc.DB.Create(&imunisasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat imunisasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data imunisasi berhasil dicatat", "data": imunisasi})
}

// 3. Update Imunisasi (PUT)
func (bc *BalitaController) UpdateImunisasi(c *gin.Context) {
	id := c.Param("id")
	var input ImunisasiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	var imunisasi models.ImunisasiBalita
	if err := bc.DB.Where("id = ?", id).First(&imunisasi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data imunisasi tidak ditemukan"})
		return
	}

	tanggalParsed, _ := time.Parse("2006-01-02", input.TanggalImunisasi)
	imunisasi.JenisImunisasi = input.JenisImunisasi
	imunisasi.TanggalImunisasi = tanggalParsed
	imunisasi.Catatan = input.Catatan

	bc.DB.Save(&imunisasi)
	c.JSON(http.StatusOK, gin.H{"message": "Data imunisasi berhasil diperbarui"})
}

// 4. Hapus Imunisasi (DELETE)
func (bc *BalitaController) DeleteImunisasi(c *gin.Context) {
	id := c.Param("id")
	if err := bc.DB.Where("id = ?", id).Delete(&models.ImunisasiBalita{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data imunisasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data imunisasi berhasil dihapus"})
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

	// --- FIX: PEMUTUSAN RANTAI FOREIGN KEY ---
	// Hapus secara paksa semua riwayat Pemeriksaan dan Imunisasi milik balita ini terlebih dahulu
	bc.DB.Where("balita_id = ?", id).Delete(&models.PemeriksaanBalita{})
	bc.DB.Where("balita_id = ?", id).Delete(&models.ImunisasiBalita{})

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

	//parsing string tanggal YYYY-MM-DD dari frontend menjadi tipe Time Golang
	tanggalParsed, err := time.Parse("2006-01-02", input.TanggalPeriksa)
	if err != nil {
		// Fallback aman: jika gagal parsing, gunakan hari ini
		tanggalParsed = time.Now()
	}

	// Ambil ID petugas (Kader/Bidan) yang sedang login dari Context Middleware
	kaderID, _ := c.Get("userID")

	pemeriksaan := models.PemeriksaanBalita{
		BalitaID:       input.BalitaID,
		TanggalPeriksa: tanggalParsed,
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

func (bc *BalitaController) UpdatePemeriksaan(c *gin.Context) {
	id := c.Param("id")
	var input UpdatePemeriksaanBalitaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	var pemeriksaan models.PemeriksaanBalita
	if err := bc.DB.Where("id = ?", id).First(&pemeriksaan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data rekam medis tidak ditemukan"})
		return
	}

	// Parsing string tanggal
	if input.TanggalPeriksa != "" {
		if tanggalParsed, err := time.Parse("2006-01-02", input.TanggalPeriksa); err == nil {
			pemeriksaan.TanggalPeriksa = tanggalParsed // <--- UPDATE TANGGALNYA
		}
	}

	pemeriksaan.BeratBadan = input.BeratBadan
	pemeriksaan.TinggiBadan = input.TinggiBadan
	pemeriksaan.LingkarKepala = input.LingkarKepala
	pemeriksaan.StatusGizi = input.StatusGizi
	pemeriksaan.Catatan = input.Catatan

	bc.DB.Save(&pemeriksaan)
	c.JSON(http.StatusOK, gin.H{"message": "Rekam medis balita berhasil diperbarui"})
}

func (bc *BalitaController) DeletePemeriksaan(c *gin.Context) {
	id := c.Param("id")
	if err := bc.DB.Where("id = ?", id).Delete(&models.PemeriksaanBalita{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus rekam medis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rekam medis balita berhasil dihapus"})
}