package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupUserRepositoryTest(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat mock database: %v", err)
	}

	repo := NewUserRepository(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func expectAdminLocks(mock sqlmock.Sqlmock, adminIDs ...int64) {
	rows := sqlmock.NewRows([]string{"id"})

	for _, id := range adminIDs {
		rows.AddRow(id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id FROM users WHERE role = 'ADMIN' AND is_active = 1 FOR UPDATE",
	)).WillReturnRows(rows)
}

func expectTargetAdmin(mock sqlmock.Sqlmock, id int64) {
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT role, is_active FROM users WHERE id = ? FOR UPDATE",
	)).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"role", "is_active"}).
				AddRow("ADMIN", true),
		)
}

func TestUpdateRoleRejectsDemotingLastActiveAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()
	adminID := int64(4)

	mock.ExpectBegin()
	expectAdminLocks(mock, adminID)
	expectTargetAdmin(mock, adminID)
	mock.ExpectRollback()

	err := repo.UpdateRole(ctx, adminID, "USER")

	if !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf("diharapkan ErrLastActiveAdmin, mendapatkan: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateStatusRejectsDeactivatingLastActiveAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()
	adminID := int64(4)

	mock.ExpectBegin()
	expectAdminLocks(mock, adminID)
	expectTargetAdmin(mock, adminID)
	mock.ExpectRollback()

	err := repo.UpdateStatus(ctx, adminID, false)

	if !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf("diharapkan ErrLastActiveAdmin, mendapatkan: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateRoleAllowsDemotingAdminWhenAnotherActiveAdminExists(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()
	adminID := int64(4)

	mock.ExpectBegin()
	expectAdminLocks(mock, 4, 5)
	expectTargetAdmin(mock, adminID)

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE users SET role = ? WHERE id = ?",
	)).
		WithArgs("USER", adminID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := repo.UpdateRole(ctx, adminID, "USER")

	if err != nil {
		t.Fatalf("seharusnya perubahan role berhasil: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateStatusAllowsDeactivationWhenAnotherActiveAdminExists(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()
	adminID := int64(4)

	mock.ExpectBegin()
	expectAdminLocks(mock, 4, 5)
	expectTargetAdmin(mock, adminID)

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE users SET is_active = ? WHERE id = ?",
	)).
		WithArgs(false, adminID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := repo.UpdateStatus(ctx, adminID, false)

	if err != nil {
		t.Fatalf("seharusnya perubahan status berhasil: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}
