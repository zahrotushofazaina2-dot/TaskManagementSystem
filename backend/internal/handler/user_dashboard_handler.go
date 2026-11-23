package handler

import (
	"net/http"

	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

type UserDashboardHandler struct {
	dashboardService *service.UserDashboardService
}

func NewUserDashboardHandler(
	dashboardService *service.UserDashboardService,
) *UserDashboardHandler {
	return &UserDashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *UserDashboardHandler) GetStats(c *gin.Context) {
	userID := c.GetInt64("user_id")

	stats, err := h.dashboardService.GetStats(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil statistik dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}
