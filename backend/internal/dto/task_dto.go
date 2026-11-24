package dto

import "time"

type CreateTaskRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type TaskResponse struct {
	ID          int64      `json:"id"`
	CategoryID  *int64     `json:"category_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}
