package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"taskmanagement/internal/model"
)

// UserRepository menangani operasi database untuk tabel users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create menyimpan pengguna baru ke database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email, password, role, is_active)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("gagal membuat pengguna: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("gagal mengambil ID pengguna: %w", err)
	}

	user.ID = id

	return nil
}

// FindByEmail mencari pengguna berdasarkan email.
func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role,
		       is_active, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mencari pengguna: %w", err)
	}

	return user, nil
}

// FindByID mencari pengguna berdasarkan ID.
func (r *UserRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role,
		       is_active, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("gagal mencari pengguna: %w", err)
	}

	return user, nil
}
