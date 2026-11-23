package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

// ListUsers mengambil daftar pengguna dengan pencarian, filter, dan pagination.
func (r *UserRepository) ListUsers(
	ctx context.Context,
	search string,
	role string,
	isActive *bool,
	page int,
	limit int,
) ([]model.User, int64, error) {
	var conditions []string
	var args []any

	if search != "" {
		conditions = append(
			conditions,
			"(name LIKE ? OR email LIKE ?)",
		)

		keyword := "%" + search + "%"
		args = append(args, keyword, keyword)
	}

	if role != "" {
		conditions = append(conditions, "role = ?")
		args = append(args, role)
	}

	if isActive != nil {
		conditions = append(conditions, "is_active = ?")
		args = append(args, *isActive)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Menghitung total data sesuai filter.
	countQuery := `
		SELECT COUNT(*)
		FROM users
	` + whereClause

	var total int64

	if err := r.db.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf(
			"gagal menghitung pengguna: %w",
			err,
		)
	}

	offset := (page - 1) * limit

	// Mengambil data pengguna tanpa kolom password.
	query := `
		SELECT id, name, email, role,
		       is_active, created_at, updated_at
		FROM users
	` + whereClause + `
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"gagal mengambil daftar pengguna: %w",
			err,
		)
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf(
				"gagal membaca data pengguna: %w",
				err,
			)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf(
			"gagal membaca hasil query pengguna: %w",
			err,
		)
	}

	return users, total, nil
}

var ErrLastActiveAdmin = errors.New("sistem harus memiliki minimal satu Admin aktif")

// lockActiveAdmins mengunci seluruh baris Admin aktif selama transaction.
func lockActiveAdmins(
	ctx context.Context,
	tx *sql.Tx,
) (int, error) {
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id FROM users WHERE role = 'ADMIN' AND is_active = 1 FOR UPDATE`,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal mengunci data Admin aktif: %w", err)
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("gagal membaca data Admin aktif: %w", err)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("gagal membaca hasil Admin aktif: %w", err)
	}

	return count, nil
}

// UpdateRole mengubah role pengguna dengan perlindungan Admin terakhir.
func (r *UserRepository) UpdateRole(
	ctx context.Context,
	userID int64,
	newRole string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaction: %w", err)
	}
	defer tx.Rollback()

	activeAdminCount, err := lockActiveAdmins(ctx, tx)
	if err != nil {
		return err
	}

	var currentRole string
	var isActive bool

	err = tx.QueryRowContext(
		ctx,
		`SELECT role, is_active FROM users WHERE id = ? FOR UPDATE`,
		userID,
	).Scan(&currentRole, &isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}

	if currentRole == "ADMIN" &&
		isActive &&
		newRole != "ADMIN" &&
		activeAdminCount <= 1 {
		return ErrLastActiveAdmin
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE users SET role = ? WHERE id = ?`,
		newRole,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal mengubah role pengguna: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal menyimpan perubahan role: %w", err)
	}

	return nil
}

// UpdateStatus mengubah status aktif pengguna dengan perlindungan Admin terakhir.
func (r *UserRepository) UpdateStatus(
	ctx context.Context,
	userID int64,
	newStatus bool,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaction: %w", err)
	}
	defer tx.Rollback()

	activeAdminCount, err := lockActiveAdmins(ctx, tx)
	if err != nil {
		return err
	}

	var currentRole string
	var isActive bool

	err = tx.QueryRowContext(
		ctx,
		`SELECT role, is_active FROM users WHERE id = ? FOR UPDATE`,
		userID,
	).Scan(&currentRole, &isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}

	if currentRole == "ADMIN" &&
		isActive &&
		!newStatus &&
		activeAdminCount <= 1 {
		return ErrLastActiveAdmin
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE users SET is_active = ? WHERE id = ?`,
		newStatus,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal mengubah status pengguna: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal menyimpan perubahan status: %w", err)
	}

	return nil
}
