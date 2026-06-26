package controllers

import (
	"net/http"
	"time"

	"posyandu-api/models" // Sesuaikan dengan nama modul Anda

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
}

// --- DTO Update ---
type UpdatePemeriksaanBalitaInput struct {
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
}

// Struktur input pencatatan imunisasi
type ImunisasiBalitaInput struct {
	BalitaID       string `json:"balita_id" binding:"required"`
	JenisImunisasi string `json:"jenis_imunisasi" binding:"required"` // e.g., "BCG", "DPT", "Campak"
	Catatan        string `json:"catatan"`
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

// RegisterBalita mendaftarkan pasien sekaligus membuatkan akun USER untuk orang tua
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

	// Gunakan transaksi DB agar jika salah satu gagal, semuanya di-rollback
	tx := bc.DB.Begin()

	// 1. Cari atau buat akun Orang Tua (Role: USER) menggunakan NoHP sebagai username
	var ortu models.User
	if err := tx.Where("username = ?", input.NoHP).First(&ortu).Error; err != nil {
		// Jika belum ada, buat akun baru dengan password default "posyandu123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("posyandu123"), bcrypt.DefaultCost)
		ortu = models.User{
			Username:     input.NoHP,
			PasswordHash: string(hashedPassword),
			NamaLengkap:  input.NamaOrangTua,
			Role:         "USER",
			IsActive:     true,
		}
		if err := tx.Create(&ortu).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun orang tua"})
			return
		}
	}

	// 2. Simpan Data Balita dengan mengaitkannya ke UserID orang tua
	balita := models.Balita{
		UserID:       ortu.ID,
		NIK:          input.NIK,
		NamaBalita:   input.NamaBalita,
		TanggalLahir: tglLahir,
		JenisKelamin: input.JenisKelamin,
		NamaOrangTua: input.NamaOrangTua,
		Alamat:       input.Alamat,
	}

	if err := tx.Create(&balita).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal mendaftarkan balita. NIK mungkin sudah terdaftar."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran Balita berhasil",
		"data": gin.H{
			"balita_id": balita.ID,
			"nama":      balita.NamaBalita,
			"ortu_akun": ortu.Username,
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

// GetRiwayatImunisasi mengambil seluruh riwayat pemberian vaksin
func (bc *BalitaController) GetRiwayatImunisasi(c *gin.Context) {
	var riwayat []models.ImunisasiBalita
	
	query := bc.DB.Preload("Pencatat").Order("tanggal_imunisasi desc")

	if balitaID := c.Query("balita_id"); balitaID != "" {
		query = query.Where("balita_id = ?", balitaID)
	}

	if err := query.Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat imunisasi balita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat imunisasi",
		"data":    riwayat,
	})
}

// CatatImunisasi menyimpan data pemberian vaksin/imunisasi
func (bc *BalitaController) CatatImunisasi(c *gin.Context) {
	var input ImunisasiBalitaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petugasID, _ := c.Get("userID")

	imunisasi := models.ImunisasiBalita{
		BalitaID:         input.BalitaID,
		JenisImunisasi:   input.JenisImunisasi,
		TanggalImunisasi: time.Now(),
		Catatan:          input.Catatan,
		DicatatOleh:      petugasID.(string),
	}

	if err := bc.DB.Create(&imunisasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data imunisasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data imunisasi berhasil dicatat",
		"data":    imunisasi,
	})
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

	// Ambil ID petugas (Kader/Bidan) yang sedang login dari Context Middleware
	kaderID, _ := c.Get("userID")

	pemeriksaan := models.PemeriksaanBalita{
		BalitaID:       input.BalitaID,
		TanggalPeriksa: time.Now(),
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