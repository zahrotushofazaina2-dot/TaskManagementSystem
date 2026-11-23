package service

import (
	"context"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/repository"
)

type AdminDashboardService struct {
	dashboardRepo *repository.AdminDashboardRepository
}

func NewAdminDashboardService(
	dashboardRepo *repository.AdminDashboardRepository,
) *AdminDashboardService {
	return &AdminDashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *AdminDashboardService) GetStats(
	ctx context.Context,
) (*dto.AdminDashboardResponse, error) {
	stats, err := s.dashboardRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.AdminDashboardResponse{
		Users: dto.UserDashboardStats{
			Total:    stats.TotalUsers,
			Active:   stats.ActiveUsers,
			Inactive: stats.InactiveUsers,
		},
		Tasks: dto.TaskDashboardStats{
			Total: stats.TotalTasks,
			ByStatus: dto.TaskStatusStats{
				Todo:       stats.TotalTodo,
				InProgress: stats.TotalInProgress,
				Done:       stats.TotalDone,
			},
			ByPriority: dto.TaskPriorityStats{
				Low:    stats.TotalLowPriority,
				Medium: stats.TotalMedium,
				High:   stats.TotalHigh,
			},
		},
		Categories: dto.CategoryDashboardStats{
			Total: stats.TotalCategories,
		},
	}, nil
}
