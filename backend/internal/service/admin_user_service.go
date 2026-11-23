package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"
)

var (
	ErrInvalidAdminFilter = errors.New("filter pengguna tidak valid")
	ErrUserNotFound       = errors.New("pengguna tidak ditemukan")
	ErrCannotModifySelf   = errors.New("Admin tidak dapat mengubah akun sendiri")
	ErrInvalidUserRole    = errors.New("role pengguna tidak valid")
)

// AdminUserService menangani logika bisnis pengelolaan pengguna oleh Admin.
type AdminUserService struct {
	userRepo *repository.UserRepository
}

// NewAdminUserService membuat instance AdminUserService.
func NewAdminUserService(
	userRepo *repository.UserRepository,
) *AdminUserService {
	return &AdminUserService{
		userRepo: userRepo,
	}
}

// ListUsers mengambil daftar pengguna dengan pencarian, filter, dan pagination.
func (s *AdminUserService) ListUsers(
	ctx context.Context,
	req dto.AdminUserListRequest,
) (*dto.AdminUserListResponse, error) {
	search := strings.TrimSpace(req.Search)
	role := strings.ToUpper(strings.TrimSpace(req.Role))

	if role != "" && role != "USER" && role != "ADMIN" {
		return nil, ErrInvalidAdminFilter
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
		return nil, ErrInvalidAdminFilter
	}

	users, total, err := s.userRepo.ListUsers(
		ctx,
		search,
		role,
		req.IsActive,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminUserResponse, 0, len(users))

	for i := range users {
		responses = append(
			responses,
			toAdminUserResponse(&users[i]),
		)
	}

	var totalPages int64
	if total > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return &dto.AdminUserListResponse{
		Users:      responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(totalPages),
	}, nil
}

// FindByID mengambil detail pengguna berdasarkan ID.
func (s *AdminUserService) FindByID(
	ctx context.Context,
	userID int64,
) (*dto.AdminUserResponse, error) {
	if userID < 1 {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, sql.ErrNoRows
	}

	response := toAdminUserResponse(user)
	return &response, nil
}

// UpdateRole mengubah role pengguna oleh Admin.
func (s *AdminUserService) UpdateRole(
	ctx context.Context,
	requesterID int64,
	targetID int64,
	newRole string,
) error {
	newRole = strings.ToUpper(strings.TrimSpace(newRole))

	if targetID < 1 {
		return ErrUserNotFound
	}

	if requesterID == targetID {
		return ErrCannotModifySelf
	}

	if newRole != "USER" && newRole != "ADMIN" {
		return ErrInvalidUserRole
	}

	user, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.UpdateRole(ctx, targetID, newRole)
}

// UpdateStatus mengubah status aktif pengguna oleh Admin.
func (s *AdminUserService) UpdateStatus(
	ctx context.Context,
	requesterID int64,
	targetID int64,
	isActive bool,
) error {
	if targetID < 1 {
		return ErrUserNotFound
	}

	if requesterID == targetID {
		return ErrCannotModifySelf
	}

	user, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.UpdateStatus(ctx, targetID, isActive)
}

// toAdminUserResponse mengubah model User menjadi DTO Admin yang aman.
func toAdminUserResponse(user *model.User) dto.AdminUserResponse {
	return dto.AdminUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
