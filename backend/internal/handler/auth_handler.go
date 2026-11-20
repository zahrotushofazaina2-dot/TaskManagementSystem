package handler

import (
	"errors"
	"net/http"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler menangani HTTP request terkait autentikasi.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler membuat instance AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register menangani registrasi pengguna.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data registrasi tidak valid",
		})
		return
	}

	result, err := h.authService.Register(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   result,
	})
}

// Login menangani login pengguna.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email atau password tidak valid",
		})
		return
	}

	result, err := h.authService.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// handleAuthError mengubah error Service menjadi HTTP response.
func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data yang diberikan tidak valid",
		})

	case errors.Is(err, service.ErrEmailAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Email sudah terdaftar",
		})

	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Email atau password salah",
		})

	case errors.Is(err, service.ErrAccountInactive):
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Akun tidak aktif",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan pada server",
		})
	}
}

// Me menangani permintaan data pengguna yang sedang login.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")

	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Identitas pengguna tidak valid",
		})
		return
	}

	result, err := h.authService.Me(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data pengguna",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}
