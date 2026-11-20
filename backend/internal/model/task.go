package model

import "time"

type Task struct {
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
