package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"taskmanagement/internal/dto"
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

// FindAll menangani pengambilan seluruh task milik user.
func (h *TaskHandler) FindAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	tasks, err := h.taskService.FindAll(
		c.Request.Context(),
		userID,
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
		"data":   tasks,
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
