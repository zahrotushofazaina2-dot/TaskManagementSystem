
package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/repository"
	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler menangani HTTP request terkait task.
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler membuat instance TaskHandler.
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// Create menangani pembuatan task.
func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data task tidak valid",
		})
		return
	}

	task, err := h.taskService.Create(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Task berhasil dibuat",
		"data":    task,
	})
}

// FindAll menangani pencarian, filter, sorting, dan pagination task milik user.
func (h *TaskHandler) FindAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Membaca parameter pagination.
	page := 1
	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Parameter page harus berupa angka minimal 1",
			})
			return
		}
		page = parsed
	}

	limit := 10
	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Parameter limit harus berupa angka antara 1 sampai 100",
			})
			return
		}
		limit = parsed
	}

	// Membaca category_id jika diberikan.
	var categoryID *int64
	if value := c.Query("category_id"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Parameter category_id tidak valid",
			})
			return
		}
		categoryID = &parsed
	}

	// Membaca dan memvalidasi parameter filter.
	search := strings.TrimSpace(c.Query("search"))
	status := strings.ToUpper(strings.TrimSpace(c.Query("status")))
	priority := strings.ToUpper(strings.TrimSpace(c.Query("priority")))
	sortBy := strings.ToLower(strings.TrimSpace(c.Query("sort_by")))
	sortOrder := strings.ToUpper(strings.TrimSpace(c.Query("sort_order")))

	if len(search) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Parameter search maksimal 200 karakter",
		})
		return
	}

	if status != "" &&
		status != "TODO" &&
		status != "IN_PROGRESS" &&
		status != "DONE" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Status harus TODO, IN_PROGRESS, atau DONE",
		})
		return
	}

	if priority != "" &&
		priority != "LOW" &&
		priority != "MEDIUM" &&
		priority != "HIGH" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Priority harus LOW, MEDIUM, atau HIGH",
		})
		return
	}

	if sortBy != "" &&
		sortBy != "created_at" &&
		sortBy != "due_date" &&
		sortBy != "title" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "sort_by harus created_at, due_date, atau title",
		})
		return
	}

	if sortOrder != "" &&
		sortOrder != "ASC" &&
		sortOrder != "DESC" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "sort_order harus ASC atau DESC",
		})
		return
	}

	// Menyusun filter untuk repository.
	filter := repository.TaskListFilter{
		Page:       page,
		Limit:      limit,
		Search:     search,
		Status:     status,
		Priority:   priority,
		CategoryID: categoryID,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	// Mengambil task berdasarkan filter.
	result, err := h.taskService.ListTasks(
		c.Request.Context(),
		userID,
		filter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// FindByID menangani pengambilan task berdasarkan ID.
func (h *TaskHandler) FindByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || taskID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID task tidak valid",
		})
		return
	}

	task, err := h.taskService.FindByID(
		c.Request.Context(),
		userID,
		taskID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Task tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

// Update menangani perubahan task.
func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || taskID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID task tidak valid",
		})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data task tidak valid",
		})
		return
	}

	task, err := h.taskService.Update(
		c.Request.Context(),
		userID,
		taskID,
		req,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Task tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task berhasil diperbarui",
		"data":    task,
	})
}

// Delete menangani penghapusan task.
func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || taskID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID task tidak valid",
		})
		return
	}

	if err := h.taskService.Delete(
		c.Request.Context(),
		userID,
		taskID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Task tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task berhasil dihapus",
	})
}
