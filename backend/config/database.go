package config

import (
	"fmt"
	"log"
	"os"

	"posyandu-api/models" // Sesuaikan dengan nama modul di go.mod Anda

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database PostgreSQL:", err)
	}

	fmt.Println("✅ Berhasil terhubung ke database Posyandu di Port:", os.Getenv("DB_PORT"))

	// Jalankan AutoMigrate
	fmt.Println("⏳ Menjalankan AutoMigrate tabel database...")
	err = db.AutoMigrate(
		&models.User{},
		&models.JadwalKegiatan{},
		&models.Balita{},
		&models.IbuHamil{},
		&models.Lansia{},
		&models.PemeriksaanBalita{},
		&models.ImunisasiBalita{},
		&models.PemeriksaanIbuHamil{},
		&models.PemeriksaanLansia{},
	)
	if err != nil {
		log.Fatal("❌ Gagal menjalankan migrasi tabel:", err)
	}
	fmt.Println("✅ AutoMigrate Sukses!")

	return db
}