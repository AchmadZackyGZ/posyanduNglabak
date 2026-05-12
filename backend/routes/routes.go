package routes

import (
	"posyandu-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Inisialisasi Controllers
	authController := controllers.NewAuthController(db)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "sukses", "message": "API Posyandu Aktif!"})
		})

		// Rute Otentikasi
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}
	}
}