package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Pastikan file .env sudah ada di folder backend.")
	}

	// 2. Merakit kunci format koneksi (DSN) untuk PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// 3. Membuka koneksi ke Database menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database PostgreSQL:", err)
	}
	fmt.Println("✅ Berhasil terhubung ke database Posyandu di Port:", os.Getenv("DB_PORT"))

	// (Nanti kita akan menambahkan kode AutoMigrate tabel ERD di sini)

	// 4. Setup Gin Router (Kerangka API)
	r := gin.Default()

	// Route percobaan (Ping) untuk memastikan server hidup
	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "sukses",
			"message": "Server API Posyandu Berjalan Lancar!",
		})
	})

	// 5. Jalankan server Golang
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback jika PORT tidak ditemukan di .env
	}
	
	fmt.Println("🚀 Server Golang berjalan di http://localhost:" + port)
	r.Run(":" + port)
}