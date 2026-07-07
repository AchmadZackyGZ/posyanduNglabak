package controllers

import (
	"net/http"
	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PengaturanController struct {
	DB *gorm.DB
}

func NewPengaturanController(db *gorm.DB) *PengaturanController {
	return &PengaturanController{DB: db}
}

type UpdatePengaturanInput struct {
	NamaPosyandu   string `json:"nama_posyandu" binding:"required"`
	Alamat         string `json:"alamat" binding:"required"`
	NomorHp        string `json:"nomor_hp" binding:"required"`
	WilayahKerja   string `json:"wilayah_kerja" binding:"required"`
	NotifJadwal    bool   `json:"notif_jadwal"`
	NotifImunisasi bool   `json:"notif_imunisasi"`
	AlertStunting  bool   `json:"alert_stunting"`
	BackupOtomatis bool   `json:"backup_otomatis"`
	ModeGelap      bool   `json:"mode_gelap"`
}

// GetPengaturan mengambil data pengaturan. Jika kosong, buat data default.
func (pc *PengaturanController) GetPengaturan(c *gin.Context) {
	var pengaturan models.Pengaturan

	// Ambil data pertama (karena pengaturan sistem biasanya hanya 1 baris di DB)
	if err := pc.DB.First(&pengaturan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// AUTO-SEEDING: Jika belum ada data, buatkan nilai default
			pengaturan = models.Pengaturan{
				NamaPosyandu:   "Posyandu Sehat Bersama",
				Alamat:         "Jl. Melati No. 5, RW 04, Kelurahan Harapan",
				NomorHp:        "0812-3456-7890",
				WilayahKerja:   "Kelurahan Harapan, Kecamatan Sejahtera",
				NotifJadwal:    true,
				NotifImunisasi: true,
				AlertStunting:  true,
				BackupOtomatis: true,
				ModeGelap:      false,
			}
			pc.DB.Create(&pengaturan)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengaturan"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil pengaturan",
		"data":    pengaturan,
	})
}

// UpdatePengaturan menyimpan perubahan yang dilakukan user
func (pc *PengaturanController) UpdatePengaturan(c *gin.Context) {
	var input UpdatePengaturanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pengaturan models.Pengaturan
	if err := pc.DB.First(&pengaturan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pengaturan belum diinisialisasi"})
		return
	}

	// Timpa data lama dengan inputan baru
	pengaturan.NamaPosyandu = input.NamaPosyandu
	pengaturan.Alamat = input.Alamat
	pengaturan.NomorHp = input.NomorHp
	pengaturan.WilayahKerja = input.WilayahKerja
	pengaturan.NotifJadwal = input.NotifJadwal
	pengaturan.NotifImunisasi = input.NotifImunisasi
	pengaturan.AlertStunting = input.AlertStunting
	pengaturan.BackupOtomatis = input.BackupOtomatis
	pengaturan.ModeGelap = input.ModeGelap

	if err := pc.DB.Save(&pengaturan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengaturan berhasil diperbarui",
		"data":    pengaturan,
	})
}