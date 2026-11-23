package service

import (
	"context"
	"errors"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/repository"
)

type UserDashboardService struct {
	dashboardRepo *repository.UserDashboardRepository
}

func NewUserDashboardService(
	dashboardRepo *repository.UserDashboardRepository,
) *UserDashboardService {
	return &UserDashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *UserDashboardService) GetStats(
	ctx context.Context,
	userID int64,
) (*dto.UserDashboardResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	stats, err := s.dashboardRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserDashboardResponse{
		Total:           stats.Total,
		TotalTodo:       stats.TotalTodo,
		TotalInProgress: stats.TotalInProgress,
		TotalDone:       stats.TotalDone,
	}, nil
}
