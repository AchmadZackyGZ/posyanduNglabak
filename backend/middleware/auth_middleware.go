package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired memvalidasi token JWT dan menyimpan data user ke Context (c *gin.Context)
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Token otentikasi tidak ditemukan"})
			c.Abort()
			return
		}

		// Format standar header: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Format token harus 'Bearer <token>'"})
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Token tidak valid atau telah kedaluwarsa"})
			c.Abort()
			return
		}

		// Ekstrak klaim (data payload di dalam token)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Simpan ID dan Role ke dalam memori Context Gin agar bisa dibaca oleh Controller
			c.Set("userID", claims["sub"])
			c.Set("userRole", claims["role"])
		}

		c.Next() // Lanjutkan ke Controller
	}
}

// RoleRequired membatasi akses endpoint hanya untuk role tertentu (Contoh: hanya ADMIN atau BIDAN)
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Identitas role tidak ditemukan"})
			c.Abort()
			return
		}

		roleStr, _ := userRole.(string)
		isAllowed := false
		for _, role := range allowedRoles {
			if strings.EqualFold(roleStr, role) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Role Anda tidak memiliki izin untuk tindakan ini"})
			c.Abort()
			return
		}

		c.Next()
	}
}