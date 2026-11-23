package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type AdminDashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	InactiveUsers    int64 `json:"inactive_users"`
	TotalTasks       int64 `json:"total_tasks"`
	TotalCategories  int64 `json:"total_categories"`
	TotalTodo        int64 `json:"total_todo"`
	TotalInProgress  int64 `json:"total_in_progress"`
	TotalDone        int64 `json:"total_done"`
	TotalLowPriority int64 `json:"total_low_priority"`
	TotalMedium      int64 `json:"total_medium_priority"`
	TotalHigh        int64 `json:"total_high_priority"`
}

type AdminDashboardRepository struct {
	db *sql.DB
}

func NewAdminDashboardRepository(db *sql.DB) *AdminDashboardRepository {
	return &AdminDashboardRepository{db: db}
}

func (r *AdminDashboardRepository) GetStats(
	ctx context.Context,
) (*AdminDashboardStats, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM users WHERE is_active = 1),
			(SELECT COUNT(*) FROM users WHERE is_active = 0),
			(SELECT COUNT(*) FROM tasks),
			(SELECT COUNT(*) FROM categories),
			(SELECT COUNT(*) FROM tasks WHERE status = 'TODO'),
			(SELECT COUNT(*) FROM tasks WHERE status = 'IN_PROGRESS'),
			(SELECT COUNT(*) FROM tasks WHERE status = 'DONE'),
			(SELECT COUNT(*) FROM tasks WHERE priority = 'LOW'),
			(SELECT COUNT(*) FROM tasks WHERE priority = 'MEDIUM'),
			(SELECT COUNT(*) FROM tasks WHERE priority = 'HIGH')
	`

	var stats AdminDashboardStats

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.InactiveUsers,
		&stats.TotalTasks,
		&stats.TotalCategories,
		&stats.TotalTodo,
		&stats.TotalInProgress,
		&stats.TotalDone,
		&stats.TotalLowPriority,
		&stats.TotalMedium,
		&stats.TotalHigh,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil statistik dashboard: %w", err)
	}

	return &stats, nil
}
