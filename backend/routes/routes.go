package routes

import (
	"posyandu-api/controllers"
	"posyandu-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// PASANG MIDDLEWARE CORS SECARA GLOBAL DI SINI
	r.Use(middleware.CORSMiddleware())
	
	authController := controllers.NewAuthController(db)
	balitaController := controllers.NewBalitaController(db)
	lansiaController := controllers.NewLansiaController(db)
	ibuHamilController := controllers.NewIbuHamilController(db)
	dashboardController := controllers.NewDashboardController(db)
	jadwalController := controllers.NewJadwalController(db)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "sukses", "message": "API Posyandu Aktif!"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		// Blok rute terproteksi (Memerlukan validasi token JWT)
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// Endpoint Dashboard
			dashboard := protected.Group("/dashboard")
			dashboard.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				dashboard.GET("/summary", dashboardController.GetSummary)
			}

			laporan := protected.Group("/laporan")
	
			// Hak akses laporan diberikan kepada ADMIN, BIDAN, dan KADER
			laporan.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				// Endpoint API Laporan Balita
				laporan.GET("/balita", controllers.GetLaporanBalita)
				
				// Endpoint API Laporan Ibu Hamil
				laporan.GET("/ibu-hamil", controllers.GetLaporanIbuHamil)
				
				// Endpoint API Laporan Lansia
				laporan.GET("/lansia", controllers.GetLaporanLansia)
			}

			// Endpoint Operasional Balita
			balita := protected.Group("/balita")
			balita.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				balita.POST("/register", balitaController.RegisterBalita)
				balita.POST("/timbang", balitaController.CatatPemeriksaan)
				balita.POST("/imunisasi", balitaController.CatatImunisasi)
				balita.GET("", balitaController.GetListBalita) // API Ambil Daftar Balita
				balita.PUT("/:id", balitaController.UpdateBalita)
				balita.DELETE("/:id", balitaController.DeleteBalita)
				balita.GET("/pemeriksaan", balitaController.GetRiwayatTimbang)
				balita.GET("/imunisasi", balitaController.GetRiwayatImunisasi)
				balita.PUT("/pemeriksaan/:id", balitaController.UpdatePemeriksaan)
				balita.DELETE("/pemeriksaan/:id", balitaController.DeletePemeriksaan)
			}

			// Endpoint Operasional Ibu Hamil
			ibuHamil := protected.Group("/ibu-hamil")
			ibuHamil.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				ibuHamil.POST("/register", ibuHamilController.RegisterIbuHamil)
				ibuHamil.POST("/periksa", ibuHamilController.CatatPemeriksaan)
				ibuHamil.GET("", ibuHamilController.GetListIbuHamil) // API Ambil Daftar Ibu Hamil
				ibuHamil.PUT("/:id", ibuHamilController.UpdateIbuHamil)
				ibuHamil.DELETE("/:id", ibuHamilController.DeleteIbuHamil)
				ibuHamil.GET("/pemeriksaan", ibuHamilController.GetRiwayatPeriksa)
				ibuHamil.PUT("/pemeriksaan/:id", ibuHamilController.UpdatePemeriksaan)
				ibuHamil.DELETE("/pemeriksaan/:id", ibuHamilController.DeletePemeriksaan)
			}

			// Endpoint Operasional Lansia
			lansia := protected.Group("/lansia")
			lansia.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				lansia.POST("/register", lansiaController.RegisterLansia)
				lansia.POST("/periksa", lansiaController.CatatPemeriksaan)
				lansia.GET("", lansiaController.GetListLansia) // API Ambil Daftar Lansia
				lansia.PUT("/:id", lansiaController.UpdateLansia)
				lansia.DELETE("/:id", lansiaController.DeleteLansia)
				lansia.GET("/pemeriksaan", lansiaController.GetRiwayatPeriksa)
				lansia.PUT("/pemeriksaan/:id", lansiaController.UpdatePemeriksaan)
				lansia.DELETE("/pemeriksaan/:id", lansiaController.DeletePemeriksaan)
			}

			// Endpoint Manajemen Jadwal
			jadwal := protected.Group("/jadwal")
			jadwal.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				jadwal.POST("", jadwalController.CreateJadwal)
				jadwal.GET("", jadwalController.GetListJadwal)
			}
		}
	}
}