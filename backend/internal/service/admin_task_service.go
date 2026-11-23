package service

import (
	"context"
	"errors"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"
)

var ErrInvalidAdminTaskFilter = errors.New("filter task tidak valid")

type AdminTaskService struct {
	taskRepo *repository.TaskRepository
}

func NewAdminTaskService(
	taskRepo *repository.TaskRepository,
) *AdminTaskService {
	return &AdminTaskService{
		taskRepo: taskRepo,
	}
}

// ListTasks mengambil seluruh task atau task berdasarkan user_id.
func (s *AdminTaskService) ListTasks(
	ctx context.Context,
	req dto.AdminTaskListRequest,
) (*dto.AdminTaskListResponse, error) {
	if req.UserID != nil && *req.UserID < 1 {
		return nil, ErrInvalidAdminTaskFilter
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 || page > 1_000_000 {
		return nil, ErrInvalidAdminTaskFilter
	}

	tasks, total, err := s.taskRepo.ListAllTasks(
		ctx,
		req.UserID,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminTaskResponse, 0, len(tasks))

	for i := range tasks {
		responses = append(
			responses,
			toAdminTaskResponse(&tasks[i]),
		)
	}

	var totalPages int64
	if total > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return &dto.AdminTaskListResponse{
		Tasks:      responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(totalPages),
	}, nil
}

func toAdminTaskResponse(task *model.Task) dto.AdminTaskResponse {
	return dto.AdminTaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		CategoryID:  task.CategoryID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
