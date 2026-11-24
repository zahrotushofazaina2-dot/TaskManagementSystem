package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"taskmanagement/internal/model"
)

type TaskListFilter struct {
	Page       int
	Limit      int
	Search     string
	Status     string
	Priority   string
	CategoryID *int64
	SortBy     string
	SortOrder  string
}

func (r *TaskRepository) ListUserTasks(
	ctx context.Context,
	userID int64,
	filter TaskListFilter,
) ([]model.Task, int64, error) {
	where := []string{"user_id = ?"}
	args := make([]any, 0)

	args = append(args, userID)

	if filter.Search != "" {
		where = append(where, "title LIKE ?")
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		where = append(where, "priority = ?")
		args = append(args, filter.Priority)
	}

	if filter.CategoryID != nil {
		where = append(where, "category_id = ?")
		args = append(args, *filter.CategoryID)
	}

	whereClause := " WHERE " + strings.Join(where, " AND ")

	countQuery := "SELECT COUNT(*) FROM tasks" + whereClause

	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung task: %w", err)
	}

	// Whitelist kolom sorting agar input tidak bisa menyisipkan SQL.
	sortColumns := map[string]string{
		"created_at": "created_at",
		"due_date":   "due_date",
		"title":      "title",
	}

	sortColumn, ok := sortColumns[filter.SortBy]
	if !ok {
		sortColumn = "created_at"
	}

	sortOrder := strings.ToUpper(filter.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	offset := (filter.Page - 1) * filter.Limit

	query := `
		SELECT id, user_id, category_id, title, description,
		       priority, status, due_date, created_at, updated_at
		FROM tasks
	` + whereClause + `
		ORDER BY ` + sortColumn + ` ` + sortOrder + `, id DESC
		LIMIT ? OFFSET ?
	`

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar task: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		var categoryID sql.NullInt64
		var description sql.NullString
		var dueDate sql.NullTime

		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&categoryID,
			&task.Title,
			&description,
			&task.Priority,
			&task.Status,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal membaca task: %w", err)
		}

		if categoryID.Valid {
			task.CategoryID = &categoryID.Int64
		}

		if description.Valid {
			task.Description = &description.String
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("gagal membaca hasil task: %w", err)
	}

	return tasks, total, nil
}
