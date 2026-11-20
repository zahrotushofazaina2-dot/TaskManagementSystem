package middleware

import (
	"errors"
	"net/http"

	"taskmanagement/config"
	"taskmanagement/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims menyimpan informasi yang terdapat di dalam JWT.
type AuthClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware memeriksa token dan status pengguna.
type AuthMiddleware struct {
	userRepo *repository.UserRepository
	authCfg  *config.AuthConfig
}

// NewAuthMiddleware membuat instance AuthMiddleware.
func NewAuthMiddleware(
	userRepo *repository.UserRepository,
	authCfg *config.AuthConfig,
) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		authCfg:  authCfg,
	}
}

// Authenticate memverifikasi JWT sebelum mengizinkan akses.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengambil token dari Authorization header.
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			unauthorized(c, "Token autentikasi diperlukan")
			return
		}

		// Format header harus: Bearer <token>
		const bearerPrefix = "Bearer "

		if len(authHeader) <= len(bearerPrefix) ||
			authHeader[:len(bearerPrefix)] != bearerPrefix {
			unauthorized(c, "Format token tidak valid")
			return
		}

		tokenString := authHeader[len(bearerPrefix):]

		// Menyiapkan claims untuk membaca isi JWT.
		claims := &AuthClaims{}

		// Memverifikasi signature, algoritma, dan expiration token.
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("algoritma JWT tidak valid")
				}

				return []byte(m.authCfg.JWTSecret), nil
			},
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
			jwt.WithExpirationRequired(),
		)

		if err != nil || token == nil || !token.Valid {
			unauthorized(c, "Token tidak valid atau sudah kedaluwarsa")
			return
		}

		// Memastikan ID pengguna valid.
		if claims.UserID <= 0 {
			unauthorized(c, "Identitas pengguna tidak valid")
			return
		}

		// Mengambil data pengguna terbaru dari database.
		user, err := m.userRepo.FindByID(
			c.Request.Context(),
			claims.UserID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memverifikasi pengguna",
			})
			c.Abort()
			return
		}

		if user == nil || !user.IsActive {
			unauthorized(c, "Akun tidak ditemukan atau tidak aktif")
			return
		}

		// Menyimpan identitas dan role terbaru ke Gin context.
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)

		// Melanjutkan request ke handler berikutnya.
		c.Next()
	}
}

// unauthorized mengirim respons ketika autentikasi gagal.
func unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": message,
	})

	c.Abort()
}
