package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware menyuntikkan kepala HTTP standar untuk mengizinkan akses dari domain SvelteKit
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengizinkan secara spesifik asal domain frontend Anda
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// Jika peramban mengirimkan permintaan preflight (OPTIONS), langsung kembalikan status sukses 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}