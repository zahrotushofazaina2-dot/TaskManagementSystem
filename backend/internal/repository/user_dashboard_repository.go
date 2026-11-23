package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type UserDashboardStats struct {
	Total           int64 `json:"total"`
	TotalTodo       int64 `json:"todo"`
	TotalInProgress int64 `json:"in_progress"`
	TotalDone       int64 `json:"done"`
}

type UserDashboardRepository struct {
	db *sql.DB
}

func NewUserDashboardRepository(db *sql.DB) *UserDashboardRepository {
	return &UserDashboardRepository{db: db}
}

func (r *UserDashboardRepository) GetStats(
	ctx context.Context,
	userID int64,
) (*UserDashboardStats, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'TODO' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'IN_PROGRESS' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'DONE' THEN 1 ELSE 0 END), 0)
		FROM tasks
		WHERE user_id = ?
	`

	var stats UserDashboardStats

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.Total,
		&stats.TotalTodo,
		&stats.TotalInProgress,
		&stats.TotalDone,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil statistik dashboard user: %w", err)
	}

	return &stats, nil
}
