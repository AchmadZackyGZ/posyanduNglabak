package routes

import (
	"posyandu-api/controllers"
	"posyandu-api/middleware" // <-- Import middleware Anda

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	balitaController := controllers.NewBalitaController(db) // <-- Inisialisasi controller balita

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "sukses", "message": "API Posyandu Aktif!"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		// GRUP RUTE TERPROTEKSI (Wajib melampirkan Token JWT)	
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired()) // <-- Pasang gembok utama
		{
			// Rute khusus operasional Balita
			balita := protected.Group("/balita")
			// Hanya ADMIN, BIDAN, dan KADER yang boleh mendaftar & menimbang
			balita.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				balita.POST("/register", balitaController.RegisterBalita)
				balita.POST("/timbang", balitaController.CatatPemeriksaan)
			}
		}
	}
}