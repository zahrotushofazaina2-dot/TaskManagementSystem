package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"taskmanagement/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// Create menyimpan kategori baru milik user.
func (r *CategoryRepository) Create(
	ctx context.Context,
	category *model.Category,
) error {
	query := `
		INSERT INTO categories (user_id, name)
		VALUES (?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		category.UserID,
		category.Name,
	)
	if err != nil {
		return fmt.Errorf("gagal membuat kategori: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("gagal mengambil ID kategori: %w", err)
	}

	category.ID = id

	return r.FindByID(ctx, category.UserID, id, category)
}

// FindByID mencari kategori berdasarkan ID dan pemiliknya.
func (r *CategoryRepository) FindByID(
	ctx context.Context,
	userID int64,
	categoryID int64,
	category *model.Category,
) error {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM categories
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		categoryID,
		userID,
	).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}

		return fmt.Errorf("gagal mencari kategori: %w", err)
	}

	return nil
}

// FindAll mengambil seluruh kategori milik user.
func (r *CategoryRepository) FindAll(
	ctx context.Context,
	userID int64,
) ([]model.Category, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM categories
		WHERE user_id = ?
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil kategori: %w", err)
	}
	defer rows.Close()

	categories := make([]model.Category, 0)

	for rows.Next() {
		var category model.Category

		if err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal membaca kategori: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca data kategori: %w", err)
	}

	return categories, nil
}

// Update mengubah kategori milik user.
func (r *CategoryRepository) Update(
	ctx context.Context,
	userID int64,
	categoryID int64,
	name string,
) error {
	query := `
		UPDATE categories
		SET name = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		name,
		categoryID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal mengubah kategori: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa perubahan kategori: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete menghapus kategori milik user.
func (r *CategoryRepository) Delete(
	ctx context.Context,
	userID int64,
	categoryID int64,
) error {
	query := `
		DELETE FROM categories
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		categoryID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal menghapus kategori: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa penghapusan kategori: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
