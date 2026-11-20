package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"taskmanagement/internal/dto"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(ctx context.Context, userID int64, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	name := strings.TrimSpace(req.Name)

	if len(name) < 2 || len(name) > 100 {
		return nil, errors.New("nama kategori harus 2 sampai 100 karakter")
	}

	category := &model.Category{
		UserID: userID,
		Name:   name,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("gagal membuat kategori: %w", err)
	}

	response := toCategoryResponse(category)
	return &response, nil
}

func (s *CategoryService) FindAll(ctx context.Context, userID int64) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		responses = append(responses, toCategoryResponse(&category))
	}

	return responses, nil
}

func (s *CategoryService) FindByID(ctx context.Context, userID, categoryID int64) (*dto.CategoryResponse, error) {
	if categoryID < 1 {
		return nil, errors.New("ID kategori tidak valid")
	}

	var category model.Category

	if err := s.categoryRepo.FindByID(ctx, userID, categoryID, &category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	response := toCategoryResponse(&category)
	return &response, nil
}

func (s *CategoryService) Update(ctx context.Context, userID, categoryID int64, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	if categoryID < 1 {
		return nil, errors.New("ID kategori tidak valid")
	}

	name := strings.TrimSpace(req.Name)

	if len(name) < 2 || len(name) > 100 {
		return nil, errors.New("nama kategori harus 2 sampai 100 karakter")
	}

	if err := s.categoryRepo.Update(ctx, userID, categoryID, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return s.FindByID(ctx, userID, categoryID)
}

func (s *CategoryService) Delete(ctx context.Context, userID, categoryID int64) error {
	if categoryID < 1 {
		return errors.New("ID kategori tidak valid")
	}

	return s.categoryRepo.Delete(ctx, userID, categoryID)
}

func toCategoryResponse(category *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
