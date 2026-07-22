package main

import (
	"fmt"
	"log"
	"os"

	"posyandu-api/config"
	"posyandu-api/routes"
	"posyandu-api/seed"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// // 1. Muat konfigurasi .env
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("❌ Error loading .env file")
	// }

	if err := godotenv.Load(); err != nil {
    log.Println("⚠️  .env file tidak ditemukan, menggunakan environment variable sistem")
}

	// 2. Inisialisasi Database & Migrasi
	db := config.ConnectDB()

	// 3. Jalankan Seeder Admin Otomatis
	seed.SeedUsers(db)

	// 4. Inisialisasi Gin Engine & Rute
	r := gin.Default()
	routes.SetupRoutes(r, db)

	// 5. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server Golang berjalan bersih di http://localhost:" + port)
	r.Run(":" + port)
}