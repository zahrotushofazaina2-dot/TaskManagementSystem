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
		`SELECT id FROM users
		 WHERE role = 'ADMIN'
		 AND is_active = 1
		 FOR UPDATE`,
	)).WillReturnRows(rows)
}

func expectActor(
	mock sqlmock.Sqlmock,
	id int64,
	role string,
) {
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT role, is_active
		 FROM users
		 WHERE id = ?
		 FOR UPDATE`,
	)).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"role", "is_active"}).
				AddRow(role, true),
		)
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

func expectTargetSuperAdmin(mock sqlmock.Sqlmock, id int64) {
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT role, is_active FROM users WHERE id = ? FOR UPDATE",
	)).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"role", "is_active"}).
				AddRow("SUPER_ADMIN", true),
		)
}

func TestUpdateRoleRejectsDemotingLastActiveAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	requesterID := int64(1)
	targetID := int64(4)

	mock.ExpectBegin()

	expectAdminLocks(mock, targetID)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetAdmin(mock, targetID)

	mock.ExpectRollback()

	err := repo.UpdateRole(
		ctx,
		requesterID,
		targetID,
		"USER",
	)

	if !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf(
			"diharapkan ErrLastActiveAdmin, mendapatkan: %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateStatusRejectsDeactivatingLastActiveAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	requesterID := int64(1)
	targetID := int64(4)

	mock.ExpectBegin()

	expectAdminLocks(mock, targetID)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetAdmin(mock, targetID)

	mock.ExpectRollback()

	err := repo.UpdateStatus(
		ctx,
		requesterID,
		targetID,
		false,
	)

	if !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf(
			"diharapkan ErrLastActiveAdmin, mendapatkan: %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateRoleAllowsDemotingAdminWhenAnotherActiveAdminExists(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	requesterID := int64(1)
	targetID := int64(5)

	mock.ExpectBegin()

	expectAdminLocks(mock, 5, 6)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetAdmin(mock, targetID)

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE users SET role = ? WHERE id = ?",
	)).
		WithArgs("USER", targetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := repo.UpdateRole(
		ctx,
		requesterID,
		targetID,
		"USER",
	)

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

	requesterID := int64(1)
	targetID := int64(5)

	mock.ExpectBegin()

	expectAdminLocks(mock, 5, 6)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetAdmin(mock, targetID)

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE users SET is_active = ? WHERE id = ?",
	)).
		WithArgs(false, targetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := repo.UpdateStatus(
		ctx,
		requesterID,
		targetID,
		false,
	)

	if err != nil {
		t.Fatalf("seharusnya perubahan status berhasil: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateRoleRejectsChangingSuperAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	requesterID := int64(1)
	targetID := int64(4)

	mock.ExpectBegin()

	expectAdminLocks(mock)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetSuperAdmin(mock, targetID)

	mock.ExpectRollback()

	err := repo.UpdateRole(
		ctx,
		requesterID,
		targetID,
		"USER",
	)

	if !errors.Is(err, ErrProtectedSuperAdmin) {
		t.Fatalf(
			"diharapkan ErrProtectedSuperAdmin, mendapatkan: %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}

func TestUpdateStatusRejectsDeactivatingSuperAdmin(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	requesterID := int64(1)
	targetID := int64(4)

	mock.ExpectBegin()

	expectAdminLocks(mock)
	expectActor(mock, requesterID, "SUPER_ADMIN")
	expectTargetSuperAdmin(mock, targetID)

	mock.ExpectRollback()

	err := repo.UpdateStatus(
		ctx,
		requesterID,
		targetID,
		false,
	)

	if !errors.Is(err, ErrProtectedSuperAdmin) {
		t.Fatalf(
			"diharapkan ErrProtectedSuperAdmin, mendapatkan: %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ekspektasi SQL tidak terpenuhi: %v", err)
	}
}
