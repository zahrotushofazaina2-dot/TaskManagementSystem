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
