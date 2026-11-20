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

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data kategori tidak valid",
		})
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}

func (h *CategoryHandler) FindAll(c *gin.Context) {
	userID := c.GetInt64("user_id")

	categories, err := h.categoryService.FindAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   categories,
	})
}

func (h *CategoryHandler) FindByID(c *gin.Context) {
	userID := c.GetInt64("user_id")
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || categoryID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid",
		})
		return
	}

	category, err := h.categoryService.FindByID(c.Request.Context(), userID, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kategori tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   category,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || categoryID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid",
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data kategori tidak valid",
		})
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), userID, categoryID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kategori tidak ditemukan",
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
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || categoryID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kategori tidak valid",
		})
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), userID, categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kategori tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kategori berhasil dihapus",
	})
}
