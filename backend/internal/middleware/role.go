package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi akses berdasarkan role pengguna.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")

		authorized := false

		for _, role := range allowedRoles {
			if userRole == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Anda tidak memiliki akses untuk fitur ini",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
