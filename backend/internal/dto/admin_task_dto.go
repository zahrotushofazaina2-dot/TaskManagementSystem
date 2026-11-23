package dto

import "time"

type AdminTaskListRequest struct {
	UserID *int64
	Page   int
	Limit  int
}

type AdminTaskResponse struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	CategoryID  *int64     `json:"category_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AdminTaskListResponse struct {
	Tasks      []AdminTaskResponse `json:"tasks"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}
