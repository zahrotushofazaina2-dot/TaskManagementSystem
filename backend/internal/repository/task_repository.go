package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"taskmanagement/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks
		(user_id, category_id, title, description, priority, status, due_date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.UserID,
		task.CategoryID,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.DueDate,
	)
	if err != nil {
		return fmt.Errorf("gagal membuat task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("gagal mengambil ID task: %w", err)
	}

	task.ID = id

	return r.FindByID(ctx, task.UserID, id, task)
}

func (r *TaskRepository) FindByID(
	ctx context.Context,
	userID, taskID int64,
	task *model.Task,
) error {
	query := `
		SELECT id, user_id, category_id, title, description,
		       priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`

	var categoryID sql.NullInt64
	var description sql.NullString
	var dueDate sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		taskID,
		userID,
	).Scan(
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
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}

		return fmt.Errorf("gagal mencari task: %w", err)
	}

	if categoryID.Valid {
		task.CategoryID = &categoryID.Int64
	} else {
		task.CategoryID = nil
	}

	if description.Valid {
		task.Description = &description.String
	} else {
		task.Description = nil
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	} else {
		task.DueDate = nil
	}

	return nil
}

func (r *TaskRepository) FindAll(
	ctx context.Context,
	userID int64,
) ([]model.Task, error) {
	query := `
		SELECT id, user_id, category_id, title, description,
		       priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil task: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		var categoryID sql.NullInt64
		var description sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
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
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca task: %w", err)
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
		return nil, fmt.Errorf("gagal membaca data task: %w", err)
	}

	return tasks, nil
}

// ListAllTasks mengambil seluruh task untuk kebutuhan Admin.
// userID nil berarti mengambil task semua pengguna.
func (r *TaskRepository) ListAllTasks(
	ctx context.Context,
	userID *int64,
	page int,
	limit int,
) ([]model.Task, int64, error) {
	whereClause := ""
	args := make([]any, 0)

	if userID != nil {
		whereClause = " WHERE user_id = ?"
		args = append(args, *userID)
	}

	countQuery := "SELECT COUNT(*) FROM tasks" + whereClause

	var total int64

	if err := r.db.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung task: %w", err)
	}

	offset := (page - 1) * limit

	query := `
		SELECT id, user_id, category_id, title, description,
		       priority, status, due_date, created_at, updated_at
		FROM tasks
	` + whereClause + `
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil seluruh task: %w", err)
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
			return nil, 0, fmt.Errorf("gagal membaca task Admin: %w", err)
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
		return nil, 0, fmt.Errorf("gagal membaca hasil task Admin: %w", err)
	}

	return tasks, total, nil
}

func (r *TaskRepository) Update(
	ctx context.Context,
	userID, taskID int64,
	task *model.Task,
) error {
	query := `
		UPDATE tasks
		SET category_id = ?, title = ?, description = ?,
		    priority = ?, status = ?, due_date = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.CategoryID,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.DueDate,
		taskID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal mengubah task: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa perubahan task: %w", err)
	}

	if affected == 0 {
		var exists int

		err := r.db.QueryRowContext(
			ctx,
			`SELECT 1 FROM tasks WHERE id = ? AND user_id = ?`,
			taskID,
			userID,
		).Scan(&exists)

		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}

		if err != nil {
			return fmt.Errorf("gagal memeriksa task: %w", err)
		}
	}

	return nil
}

func (r *TaskRepository) Delete(
	ctx context.Context,
	userID, taskID int64,
) error {
	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("gagal menghapus task: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa penghapusan task: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
