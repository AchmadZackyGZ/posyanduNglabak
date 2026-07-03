package controllers

import (
	"net/http"
	"posyandu-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventarisController struct {
	DB *gorm.DB
}

func NewInventarisController(db *gorm.DB) *InventarisController {
	return &InventarisController{DB: db}
}

// Input Structs
type InventarisInput struct {
	NamaBarang string `json:"nama_barang" binding:"required"`
	Kategori   string `json:"kategori" binding:"required"`
	Stok       int    `json:"stok"`
}

type UpdateStokInput struct {
	Stok int `json:"stok" binding:"required"`
}

// 1. Get All Inventaris
func (ic *InventarisController) GetListInventaris(c *gin.Context) {
	var barang []models.Inventaris
	if err := ic.DB.Order("created_at desc").Find(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data inventaris"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil",
		"data":    barang,
	})
}

// 2. Tambah Barang Baru
func (ic *InventarisController) CreateInventaris(c *gin.Context) {
	var input InventarisInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barang := models.Inventaris{
		NamaBarang: input.NamaBarang,
		Kategori:   input.Kategori,
		Stok:       input.Stok,
	}

	if err := ic.DB.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan barang ke inventaris"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Barang berhasil ditambahkan",
		"data":    barang,
	})
}

// 3. Update Barang (Nama / Kategori)
func (ic *InventarisController) UpdateInventaris(c *gin.Context) {
	id := c.Param("id")
	var input InventarisInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var barang models.Inventaris
	if err := ic.DB.First(&barang, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	barang.NamaBarang = input.NamaBarang
	barang.Kategori = input.Kategori
	// Note: Stok tidak diupdate di sini untuk menjaga integritas, ada fungsi khusus

	ic.DB.Save(&barang)
	c.JSON(http.StatusOK, gin.H{"message": "Data barang berhasil diperbarui", "data": barang})
}

// 4. Update Stok Khusus (+updateStok dari Class Diagram)
func (ic *InventarisController) UpdateStok(c *gin.Context) {
	id := c.Param("id")
	var input UpdateStokInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format stok tidak valid"})
		return
	}

	if err := ic.DB.Model(&models.Inventaris{}).Where("id = ?", id).Update("stok", input.Stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui stok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stok berhasil diperbarui"})
}

// 5. Hapus Barang
func (ic *InventarisController) DeleteInventaris(c *gin.Context) {
	id := c.Param("id")
	if err := ic.DB.Delete(&models.Inventaris{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus barang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus"})
}