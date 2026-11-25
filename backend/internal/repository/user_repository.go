package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"taskmanagement/internal/model"
)

var (
	ErrLastActiveAdmin = errors.New(
		"sistem harus memiliki minimal satu akun administrator aktif",
	)

	ErrProtectedSuperAdmin = errors.New(
		"akun SUPER_ADMIN tidak dapat diubah",
	)

	ErrProtectedAdmin = errors.New(
		"ADMIN tidak memiliki izin untuk mengubah akun ADMIN",
	)

	ErrCannotModifySelf = errors.New(
		"Admin tidak dapat mengubah akun sendiri",
	)

	ErrCannotAssignSuperAdmin = errors.New(
		"role SUPER_ADMIN tidak dapat diberikan melalui fitur ini",
	)

	ErrUnauthorizedUserManagement = errors.New(
		"Anda tidak memiliki izin untuk mengelola pengguna",
	)
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
func (r *UserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
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

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
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

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
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

// lockActiveAdmins mengunci seluruh akun administrator aktif selama transaction.
// SUPER_ADMIN ikut dihitung sebagai administrator aktif.
// lockActiveAdmins mengunci seluruh akun ADMIN aktif selama transaction.
func lockActiveAdmins(
	ctx context.Context,
	tx *sql.Tx,
) (int, error) {
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id FROM users
		 WHERE role = 'ADMIN'
		 AND is_active = 1
		 FOR UPDATE`,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"gagal mengunci data administrator aktif: %w",
			err,
		)
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf(
				"gagal membaca data administrator aktif: %w",
				err,
			)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf(
			"gagal membaca hasil administrator aktif: %w",
			err,
		)
	}

	return count, nil
}

// validateManagementActor memastikan pelaku masih aktif dan role aktualnya
// memiliki izin untuk mengelola pengguna.
func validateManagementActor(
	ctx context.Context,
	tx *sql.Tx,
	requesterID int64,
) (string, error) {
	if requesterID < 1 {
		return "", ErrUnauthorizedUserManagement
	}

	var role string
	var isActive bool

	err := tx.QueryRowContext(
		ctx,
		`SELECT role, is_active
		 FROM users
		 WHERE id = ?
		 FOR UPDATE`,
		requesterID,
	).Scan(&role, &isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUnauthorizedUserManagement
		}

		return "", fmt.Errorf(
			"gagal memverifikasi akun pelaku: %w",
			err,
		)
	}

	role = strings.ToUpper(strings.TrimSpace(role))

	if !isActive {
		return "", ErrUnauthorizedUserManagement
	}

	if role != "ADMIN" && role != "SUPER_ADMIN" {
		return "", ErrUnauthorizedUserManagement
	}

	return role, nil
}

// validateTarget memastikan akun target boleh dikelola oleh pelaku.
func validateTarget(
	ctx context.Context,
	tx *sql.Tx,
	requesterID int64,
	requesterRole string,
	targetID int64,
) (string, bool, error) {
	if targetID < 1 {
		return "", false, sql.ErrNoRows
	}

	if requesterID == targetID {
		return "", false, ErrCannotModifySelf
	}

	var targetRole string
	var targetActive bool

	err := tx.QueryRowContext(
		ctx,
		`SELECT role, is_active
		 FROM users
		 WHERE id = ?
		 FOR UPDATE`,
		targetID,
	).Scan(&targetRole, &targetActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, sql.ErrNoRows
		}

		return "", false, fmt.Errorf(
			"gagal mengambil data target: %w",
			err,
		)
	}

	targetRole = strings.ToUpper(strings.TrimSpace(targetRole))

	if targetRole == "SUPER_ADMIN" {
		return "", false, ErrProtectedSuperAdmin
	}

	if requesterRole == "ADMIN" && targetRole == "ADMIN" {
		return "", false, ErrProtectedAdmin
	}

	return targetRole, targetActive, nil
}

// UpdateRole mengubah role pengguna dengan pemeriksaan izin di database.
// UpdateRole mengubah role pengguna dengan pemeriksaan izin di database.
func (r *UserRepository) UpdateRole(
	ctx context.Context,
	requesterID int64,
	userID int64,
	newRole string,
) error {
	newRole = strings.ToUpper(strings.TrimSpace(newRole))

	if newRole == "SUPER_ADMIN" {
		return ErrCannotAssignSuperAdmin
	}

	if newRole != "USER" && newRole != "ADMIN" {
		return errors.New("role pengguna tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaction: %w", err)
	}
	defer tx.Rollback()

	activeAdminCount, err := lockActiveAdmins(ctx, tx)
	if err != nil {
		return err
	}

	requesterRole, err := validateManagementActor(
		ctx,
		tx,
		requesterID,
	)
	if err != nil {
		return err
	}

	currentRole, targetActive, err := validateTarget(
		ctx,
		tx,
		requesterID,
		requesterRole,
		userID,
	)
	if err != nil {
		return err
	}

	if requesterRole == "ADMIN" && newRole == "ADMIN" {
		return ErrUnauthorizedUserManagement
	}

	if currentRole == "ADMIN" &&
		targetActive &&
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

// UpdateStatus mengubah status pengguna dengan pemeriksaan izin di database.
func (r *UserRepository) UpdateStatus(
	ctx context.Context,
	requesterID int64,
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

	requesterRole, err := validateManagementActor(
		ctx,
		tx,
		requesterID,
	)
	if err != nil {
		return err
	}

	currentRole, isActive, err := validateTarget(
		ctx,
		tx,
		requesterID,
		requesterRole,
		userID,
	)
	if err != nil {
		return err
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
