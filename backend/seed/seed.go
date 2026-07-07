package seed

import (
	"fmt"
	"log"

	"posyandu-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAccount adalah struct pembantu murni untuk mendefinisikan daftar akun yang akan disemai
type SeedAccount struct {
	Username    string
	Password    string
	NamaLengkap string
	Role        string
}

func SeedUsers(db *gorm.DB) {
	// 1. Definisikan array daftar akun dummy di sini
	accounts := []SeedAccount{
		{
			Username:    "admin",
			Password:    "adminrahasia",
			NamaLengkap: "Admin Posyandu Ngablak",
			Role:        "ADMIN",
		},
		{
			Username:    "kader1",
			Password:    "password123",
			NamaLengkap: "Ibu Kader Mawar",
			Role:        "KADER",
		},
		{
			Username:    "bidan1",
			Password:    "password123",
			NamaLengkap: "Bidan Siti",
			Role:        "BIDAN",
		},
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("⏳ Memulai proses verifikasi dan seeding akun...")
	fmt.Println("--------------------------------------------------")

	// 2. Lakukan iterasi untuk setiap akun di dalam array
	for _, acc := range accounts {
		var count int64

		// Validasi Ketat: Cek apakah username sudah eksis di database
		db.Model(&models.User{}).Where("username = ?", acc.Username).Count(&count)

		if count > 0 {
			// Jika sudah ada, lewati eksekusi dan lanjut ke akun berikutnya
			fmt.Printf("⏩ Akun %s (%s) sudah ada di database, skip pembuatan.\n", acc.Role, acc.Username)
			continue 
		}

		fmt.Printf("⏳ Akun %s (%s) belum ditemukan. Melakukan hashing & pembuatan...\n", acc.Role, acc.Username)

		// Lakukan enkripsi password secara dinamis sesuai password di array
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Gagal melakukan hashing password untuk %s: %v", acc.Username, err)
		}

		// Siapkan model User GORM
		newUser := models.User{
			Username:     acc.Username,
			PasswordHash: string(hashedPassword),
			NamaLengkap:  acc.NamaLengkap,
			Role:         acc.Role,
			IsActive:     true,
		}

		// Eksekusi injeksi data ke tabel
		if err := db.Create(&newUser).Error; err != nil {
			log.Fatalf("❌ Gagal menyimpan akun %s: %v", acc.Username, err)
		}

		fmt.Printf("✅ Seed %s berhasil! (Username: %s | Pass: %s)\n", acc.Role, acc.Username, acc.Password)
	}
	
	fmt.Println("--------------------------------------------------")
	fmt.Println("🎉 Proses seeding selesai!")
	fmt.Println("--------------------------------------------------")
}