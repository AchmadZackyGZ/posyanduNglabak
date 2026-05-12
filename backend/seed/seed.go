package seed

import (
	"fmt"
	"log"

	"posyandu-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	
	// Validasi: Jika akun admin sudah ada, langsung skip
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		fmt.Println("⏩ Akun Admin sudah ada di database, skip otomatisasi pembuatan akun.")
		return
	}

	fmt.Println("⏳ Akun Admin belum ditemukan. Melakukan seeding akun Admin baru...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("adminrahasia"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Gagal melakukan hashing password admin:", err)
	}

	admin := models.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		NamaLengkap:  "Admin Posyandu Ngablak",
		Role:         "ADMIN",
		IsActive:     true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("❌ Gagal menyimpan akun Admin seed:", err)
	}

	fmt.Println("✅ Seed Admin berhasil dibuat! Gunakan username: admin & password: adminrahasia")
}