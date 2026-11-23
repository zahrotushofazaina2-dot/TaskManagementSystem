package handler

import (
	"net/http"

	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	dashboardService *service.AdminDashboardService
}

func NewAdminDashboardHandler(
	dashboardService *service.AdminDashboardService,
) *AdminDashboardHandler {
	return &AdminDashboardHandler{
		dashboardService: dashboardService,
	}
}

// GET /admin/dashboard
func (h *AdminDashboardHandler) GetStats(c *gin.Context) {
	result, err := h.dashboardService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil statistik dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Statistik dashboard berhasil diambil",
		"data":    result,
	})
}
