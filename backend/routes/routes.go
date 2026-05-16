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
			// Endpoint Operasional Balita
			balita := protected.Group("/balita")
			balita.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				balita.POST("/register", balitaController.RegisterBalita)
				balita.POST("/timbang", balitaController.CatatPemeriksaan)
			}

			ibuHamil := protected.Group("/ibu-hamil")
			ibuHamil.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				ibuHamil.POST("/register", ibuHamilController.RegisterIbuHamil)
				ibuHamil.POST("/periksa", ibuHamilController.CatatPemeriksaan)
			}

			// <-- 2. INJEKSI ENDPOINT LANSIA DI SINI
			lansia := protected.Group("/lansia")
			lansia.Use(middleware.RoleRequired("ADMIN", "BIDAN", "KADER"))
			{
				lansia.POST("/register", lansiaController.RegisterLansia)
				lansia.POST("/periksa", lansiaController.CatatPemeriksaan)
			}
		}
	}
}