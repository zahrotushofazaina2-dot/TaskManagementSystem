package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"
)

type TaskService struct {
	taskRepo     *repository.TaskRepository
	categoryRepo *repository.CategoryRepository
}

func NewTaskService(
	taskRepo *repository.TaskRepository,
	categoryRepo *repository.CategoryRepository,
) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		categoryRepo: categoryRepo,
	}
}

// Create membuat task baru milik user.
func (s *TaskService) Create(
	ctx context.Context,
	userID int64,
	req dto.CreateTaskRequest,
) (*dto.TaskResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	title := strings.TrimSpace(req.Title)

	if utf8.RuneCountInString(title) < 2 ||
		utf8.RuneCountInString(title) > 200 {
		return nil, errors.New("judul task harus 2 sampai 200 karakter")
	}

	description := req.Description
	if description != nil {
		if utf8.RuneCountInString(*description) > 5000 {
			return nil, errors.New("deskripsi maksimal 5000 karakter")
		}
	}

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if priority == "" {
		priority = "MEDIUM"
	}

	if priority != "LOW" &&
		priority != "MEDIUM" &&
		priority != "HIGH" {
		return nil, errors.New("priority harus LOW, MEDIUM, atau HIGH")
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = "TODO"
	}

	if status != "TODO" &&
		status != "IN_PROGRESS" &&
		status != "DONE" {
		return nil, errors.New("status harus TODO, IN_PROGRESS, atau DONE")
	}

	dueDate, err := parseDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}

	if err := s.validateCategory(ctx, userID, req.CategoryID); err != nil {
		return nil, err
	}

	task := &model.Task{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      status,
		DueDate:     dueDate,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("gagal membuat task: %w", err)
	}

	response := toTaskResponse(task)
	return &response, nil
}

// FindAll mengambil seluruh task milik user.
// Method ini dipertahankan untuk kompatibilitas.
func (s *TaskService) FindAll(
	ctx context.Context,
	userID int64,
) ([]dto.TaskResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	tasks, err := s.taskRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TaskResponse, 0, len(tasks))

	for i := range tasks {
		responses = append(responses, toTaskResponse(&tasks[i]))
	}

	return responses, nil
}

// ListTasks mengambil task dengan search, filter, sorting, dan pagination.
func (s *TaskService) ListTasks(
	ctx context.Context,
	userID int64,
	filter repository.TaskListFilter,
) (*dto.TaskListResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	// Default pagination.
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Page < 1 {
		return nil, errors.New("page minimal 1")
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		return nil, errors.New("limit harus antara 1 sampai 100")
	}

	// Validasi search.
	filter.Search = strings.TrimSpace(filter.Search)
	if utf8.RuneCountInString(filter.Search) > 200 {
		return nil, errors.New("search maksimal 200 karakter")
	}

	// Validasi status.
	filter.Status = strings.ToUpper(strings.TrimSpace(filter.Status))
	if filter.Status != "" &&
		filter.Status != "TODO" &&
		filter.Status != "IN_PROGRESS" &&
		filter.Status != "DONE" {
		return nil, errors.New("status harus TODO, IN_PROGRESS, atau DONE")
	}

	// Validasi priority.
	filter.Priority = strings.ToUpper(strings.TrimSpace(filter.Priority))
	if filter.Priority != "" &&
		filter.Priority != "LOW" &&
		filter.Priority != "MEDIUM" &&
		filter.Priority != "HIGH" {
		return nil, errors.New("priority harus LOW, MEDIUM, atau HIGH")
	}

	// Validasi category ID.
	if filter.CategoryID != nil && *filter.CategoryID < 1 {
		return nil, errors.New("category_id harus lebih besar dari 0")
	}

	// Default dan validasi sorting.
	filter.SortBy = strings.ToLower(strings.TrimSpace(filter.SortBy))
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	if filter.SortBy != "created_at" &&
		filter.SortBy != "due_date" &&
		filter.SortBy != "title" {
		return nil, errors.New("sort_by harus created_at, due_date, atau title")
	}

	filter.SortOrder = strings.ToUpper(strings.TrimSpace(filter.SortOrder))
	if filter.SortOrder == "" {
		filter.SortOrder = "DESC"
	}

	if filter.SortOrder != "ASC" && filter.SortOrder != "DESC" {
		return nil, errors.New("sort_order harus ASC atau DESC")
	}

	tasks, total, err := s.taskRepo.ListUserTasks(
		ctx,
		userID,
		filter,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TaskResponse, 0, len(tasks))
	for i := range tasks {
		responses = append(responses, toTaskResponse(&tasks[i]))
	}

	var totalPages int64
	if total > 0 {
		totalPages = (total + int64(filter.Limit) - 1) / int64(filter.Limit)
	}

	return &dto.TaskListResponse{
		Tasks: responses,
		Pagination: dto.PaginationMeta{
			Total:      total,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
		},
	}, nil
}

// FindByID mengambil detail task berdasarkan ID dan pemiliknya.
func (s *TaskService) FindByID(
	ctx context.Context,
	userID, taskID int64,
) (*dto.TaskResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	if taskID < 1 {
		return nil, errors.New("ID task tidak valid")
	}

	var task model.Task

	if err := s.taskRepo.FindByID(ctx, userID, taskID, &task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	response := toTaskResponse(&task)
	return &response, nil
}

// Update mengubah task milik user.
func (s *TaskService) Update(
	ctx context.Context,
	userID, taskID int64,
	req dto.UpdateTaskRequest,
) (*dto.TaskResponse, error) {
	if userID < 1 {
		return nil, errors.New("user ID tidak valid")
	}

	if taskID < 1 {
		return nil, errors.New("ID task tidak valid")
	}

	title := strings.TrimSpace(req.Title)

	if utf8.RuneCountInString(title) < 2 ||
		utf8.RuneCountInString(title) > 200 {
		return nil, errors.New("judul task harus 2 sampai 200 karakter")
	}

	description := req.Description
	if description != nil {
		if utf8.RuneCountInString(*description) > 5000 {
			return nil, errors.New("deskripsi maksimal 5000 karakter")
		}
	}

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if priority == "" {
		priority = "MEDIUM"
	}

	if priority != "LOW" &&
		priority != "MEDIUM" &&
		priority != "HIGH" {
		return nil, errors.New("priority harus LOW, MEDIUM, atau HIGH")
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = "TODO"
	}

	if status != "TODO" &&
		status != "IN_PROGRESS" &&
		status != "DONE" {
		return nil, errors.New("status harus TODO, IN_PROGRESS, atau DONE")
	}

	dueDate, err := parseDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}

	if err := s.validateCategory(ctx, userID, req.CategoryID); err != nil {
		return nil, err
	}

	task := &model.Task{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      status,
		DueDate:     dueDate,
	}

	if err := s.taskRepo.Update(ctx, userID, taskID, task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return s.FindByID(ctx, userID, taskID)
}

// Delete menghapus task milik user.
func (s *TaskService) Delete(
	ctx context.Context,
	userID, taskID int64,
) error {
	if userID < 1 {
		return errors.New("user ID tidak valid")
	}

	if taskID < 1 {
		return errors.New("ID task tidak valid")
	}

	return s.taskRepo.Delete(ctx, userID, taskID)
}

// validateCategory memastikan kategori milik user yang sedang login.
func (s *TaskService) validateCategory(
	ctx context.Context,
	userID int64,
	categoryID *int64,
) error {
	if categoryID == nil {
		return nil
	}

	if *categoryID < 1 {
		return errors.New("ID kategori tidak valid")
	}

	var category model.Category

	err := s.categoryRepo.FindByID(
		ctx,
		userID,
		*categoryID,
		&category,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("kategori tidak ditemukan atau bukan milik user")
		}

		return fmt.Errorf("gagal memvalidasi kategori: %w", err)
	}

	return nil
}

// parseDueDate memvalidasi tanggal dengan format YYYY-MM-DD.
func parseDueDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	date := strings.TrimSpace(*value)

	if date == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("format due_date harus YYYY-MM-DD")
	}

	return &parsed, nil
}

// toTaskResponse mengubah model menjadi DTO response.
func toTaskResponse(task *model.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          task.ID,
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
