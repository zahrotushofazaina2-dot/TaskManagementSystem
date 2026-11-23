package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/repository"
	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	adminUserService *service.AdminUserService
}

func NewAdminUserHandler(
	adminUserService *service.AdminUserService,
) *AdminUserHandler {
	return &AdminUserHandler{
		adminUserService: adminUserService,
	}
}

// GET /admin/users
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter page harus berupa angka",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter limit harus berupa angka",
		})
		return
	}

	var isActive *bool

	if value := c.Query("is_active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter is_active harus true atau false",
			})
			return
		}
		isActive = &parsed
	}

	req := dto.AdminUserListRequest{
		Search:   c.Query("search"),
		Role:     c.Query("role"),
		IsActive: isActive,
		Page:     page,
		Limit:    limit,
	}

	result, err := h.adminUserService.ListUsers(
		c.Request.Context(),
		req,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAdminFilter) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil daftar pengguna",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar pengguna berhasil diambil",
		"data":    result,
	})
}

// GET /admin/users/:id
func (h *AdminUserHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID pengguna tidak valid",
		})
		return
	}

	result, err := h.adminUserService.FindByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) ||
			errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Pengguna tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil detail pengguna",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail pengguna berhasil diambil",
		"data":    result,
	})
}

// PUT /admin/users/:id/role
func (h *AdminUserHandler) UpdateRole(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID pengguna tidak valid",
		})
		return
	}

	var req dto.UpdateUserRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Role wajib diisi dan hanya boleh USER atau ADMIN",
		})
		return
	}

	requesterValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Identitas pengguna tidak ditemukan",
		})
		return
	}

	requesterID, ok := requesterValue.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Format identitas pengguna tidak valid",
		})
		return
	}

	err = h.adminUserService.UpdateRole(
		c.Request.Context(),
		requesterID,
		targetID,
		req.Role,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCannotModifySelf):
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		case errors.Is(err, service.ErrInvalidUserRole):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		case errors.Is(err, service.ErrUserNotFound),
			errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Pengguna tidak ditemukan",
			})
		case errors.Is(err, repository.ErrLastActiveAdmin):
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mengubah role pengguna",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role pengguna berhasil diperbarui",
	})
}

// PUT /admin/users/:id/status
func (h *AdminUserHandler) UpdateStatus(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID pengguna tidak valid",
		})
		return
	}

	var req dto.UpdateUserStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "is_active wajib diisi dan harus berupa boolean",
		})
		return
	}

	requesterValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Identitas pengguna tidak ditemukan",
		})
		return
	}

	requesterID, ok := requesterValue.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Format identitas pengguna tidak valid",
		})
		return
	}

	err = h.adminUserService.UpdateStatus(
		c.Request.Context(),
		requesterID,
		targetID,
		*req.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCannotModifySelf):
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		case errors.Is(err, service.ErrUserNotFound),
			errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Pengguna tidak ditemukan",
			})
		case errors.Is(err, repository.ErrLastActiveAdmin):
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal mengubah status pengguna",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pengguna berhasil diperbarui",
	})
}
