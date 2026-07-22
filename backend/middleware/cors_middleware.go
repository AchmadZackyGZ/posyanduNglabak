package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware menyuntikkan kepala HTTP standar untuk mengizinkan akses dari domain SvelteKit
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Dapatkan siapa yang sedang mencoba mengetuk pintu
		origin := c.Request.Header.Get("Origin")

		// 2. Daftar tamu VIP dasar (untuk development lokal)
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
		}

		// 3. Tambahkan origin frontend production dari env var
		//    Bisa lebih dari satu domain, dipisah koma di value env var-nya
		//    contoh: FRONTEND_URL=https://posyandu-nglabak.vercel.app,https://posyandu-nglabak-git-staging.vercel.app
		if frontendURLs := os.Getenv("FRONTEND_URL"); frontendURLs != "" {
			for _, url := range strings.Split(frontendURLs, ",") {
				allowedOrigins[strings.TrimSpace(url)] = true
			}
		}

		// 4. Jika tamu ada di daftar VIP, bukakan pintu khusus untuknya
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// 5. Jika peramban mengirimkan permintaan preflight (OPTIONS), kembalikan sukses 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}