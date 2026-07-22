package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware menyuntikkan kepala HTTP standar untuk mengizinkan akses dari domain SvelteKit
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Dapatkan siapa yang sedang mencoba mengetuk pintu
		origin := c.Request.Header.Get("Origin")

		// 2. Daftar tamu VIP (Masukkan URL Vercel Anda di sini)
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			// "https://posyandu-ngablak.vercel.app":   true, // <--- URL VERCEL ANDA
		}

		// Tambahkan origin frontend production dari env var (dipisah koma kalau lebih dari satu)
		if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
			allowedOrigins[frontendURL] = true
		}

		// 3. Jika tamu ada di daftar VIP, bukakan pintu khusus untuknya
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 4. Buka gerbang untuk paspor Ngrok (Tambahkan di akhir)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning")
		
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// Jika peramban mengirimkan permintaan preflight (OPTIONS), kembalikan sukses 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}