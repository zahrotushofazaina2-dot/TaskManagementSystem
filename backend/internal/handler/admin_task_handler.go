package handler

import (
	"net/http"
	"strconv"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminTaskHandler struct {
	adminTaskService *service.AdminTaskService
}

func NewAdminTaskHandler(
	adminTaskService *service.AdminTaskService,
) *AdminTaskHandler {
	return &AdminTaskHandler{
		adminTaskService: adminTaskService,
	}
}

// GET /admin/tasks
func (h *AdminTaskHandler) ListTasks(c *gin.Context) {
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

	var userID *int64

	if value := c.Query("user_id"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Parameter user_id tidak valid",
			})
			return
		}
		userID = &parsed
	}

	req := dto.AdminTaskListRequest{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.adminTaskService.ListTasks(
		c.Request.Context(),
		req,
	)
	if err != nil {
		if err == service.ErrInvalidAdminTaskFilter {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil daftar task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar task berhasil diambil",
		"data":    result,
	})
}
