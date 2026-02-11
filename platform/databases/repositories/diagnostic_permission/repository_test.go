package diagnostic_permission

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewRepository Tests
// ============================================

func TestNewRepository_NilDB(t *testing.T) {
	repo, err := NewRepository(nil)

	assert.Nil(t, repo)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id.*AND branch_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id.*AND active")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsert")
}

func TestNewRepository_PrepareError_Deactivate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeactivate")
}

func TestNewRepository_PrepareError_GetByMotorcycleAndBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id.*AND branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleAndBranch")
}

func TestNewRepository_PrepareError_GetByMotorcycleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id.*AND branch_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id.*AND active").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleID")
}

// ============================================
// BeginTx Tests
// ============================================

func TestBeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	repo := &repository{db: db}
	tx, err := repo.BeginTx(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	tx, err := repo.BeginTx(context.Background())

	assert.Nil(t, tx)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// Save Tests
// ============================================

func TestSave_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO motorcycle_diagnostic_permissions").
		WithArgs("perm-123", "moto-123", "branch-456", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	permission := &domain.DiagnosticPermission{
		ID:           "perm-123",
		MotorcycleID: "moto-123",
		BranchID:     "branch-456",
		Active:       true,
	}

	repo := &repository{db: db}
	err = repo.Save(context.Background(), tx, permission)

	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	permission := &domain.DiagnosticPermission{ID: "perm-invalid-tx"}

	repo := &repository{db: db}
	err = repo.Save(context.Background(), nil, permission)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO motorcycle_diagnostic_permissions").
		WithArgs("perm-error", "moto-123", "branch-456", true).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	permission := &domain.DiagnosticPermission{
		ID:           "perm-error",
		MotorcycleID: "moto-123",
		BranchID:     "branch-456",
		Active:       true,
	}

	repo := &repository{db: db}
	err = repo.Save(context.Background(), tx, permission)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionCannotSave, err)
}

// ============================================
// Deactivate Tests
// ============================================

func TestDeactivate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_diagnostic_permissions").
		WithArgs("moto-123", "branch-456").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}
	err = repo.Deactivate(context.Background(), tx, "moto-123", "branch-456")

	assert.NoError(t, err)
}

func TestDeactivate_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}
	err = repo.Deactivate(context.Background(), nil, "moto-123", "branch-456")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDeactivate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_diagnostic_permissions").
		WithArgs("moto-123", "branch-456").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}
	err = repo.Deactivate(context.Background(), tx, "moto-123", "branch-456")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionCannotDelete, err)
}

// ============================================
// GetByMotorcycleAndBranch Tests
// ============================================

func TestGetByMotorcycleAndBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "active"}).
		AddRow("perm-001", "moto-123", "branch-456", true)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123", "branch-456").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ? AND active = TRUE")

	permission, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-123", "branch-456")

	assert.NoError(t, err)
	assert.NotNil(t, permission)
	assert.Equal(t, "perm-001", permission.ID)
	assert.Equal(t, "moto-123", permission.MotorcycleID)
	assert.Equal(t, "branch-456", permission.BranchID)
	assert.True(t, permission.Active)
}

func TestGetByMotorcycleAndBranch_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-999", "branch-999").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ? AND active = TRUE")

	permission, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-999", "branch-999")

	assert.Nil(t, permission)
	assert.Equal(t, domain.ErrPermissionNotFound, err)
}

func TestGetByMotorcycleAndBranch_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-err", "branch-err").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ? AND active = TRUE")

	permission, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-err", "branch-err")

	assert.Nil(t, permission)
	assert.Error(t, err)
}

// ============================================
// GetByMotorcycleID Tests
// ============================================

func TestGetByMotorcycleID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "active"}).
		AddRow("perm-001", "moto-123", "branch-A", true).
		AddRow("perm-002", "moto-123", "branch-B", true)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	permissions, err := repo.GetByMotorcycleID(context.Background(), "moto-123")

	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
	assert.Equal(t, "perm-001", permissions[0].ID)
	assert.Equal(t, "branch-A", permissions[0].BranchID)
	assert.Equal(t, "perm-002", permissions[1].ID)
	assert.Equal(t, "branch-B", permissions[1].BranchID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "active"})

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-no-perms").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	permissions, err := repo.GetByMotorcycleID(context.Background(), "moto-no-perms")

	assert.NoError(t, err)
	assert.Empty(t, permissions)
}

func TestGetByMotorcycleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	permissions, err := repo.GetByMotorcycleID(context.Background(), "moto-error")

	assert.Nil(t, permissions)
	assert.Error(t, err)
}

// ============================================
// Model Tests
// ============================================

func TestDiagnosticPermission_ToDomain(t *testing.T) {
	dbPerm := &DiagnosticPermission{
		ID:           "perm-001",
		MotorcycleID: "moto-001",
		BranchID:     "branch-001",
		Active:       true,
	}

	dp := dbPerm.ToDomain()

	assert.Equal(t, "perm-001", dp.ID)
	assert.Equal(t, "moto-001", dp.MotorcycleID)
	assert.Equal(t, "branch-001", dp.BranchID)
	assert.True(t, dp.Active)
}

func TestDiagnosticPermission_ToDomain_Inactive(t *testing.T) {
	dbPerm := &DiagnosticPermission{
		ID:           "perm-002",
		MotorcycleID: "moto-002",
		BranchID:     "branch-002",
		Active:       false,
	}

	dp := dbPerm.ToDomain()

	assert.Equal(t, "perm-002", dp.ID)
	assert.False(t, dp.Active)
}

func TestFromDomain(t *testing.T) {
	dp := &domain.DiagnosticPermission{
		ID:           "perm-003",
		MotorcycleID: "moto-003",
		BranchID:     "branch-003",
		Active:       true,
	}

	dbPerm := FromDomain(dp)

	assert.Equal(t, "perm-003", dbPerm.ID)
	assert.Equal(t, "moto-003", dbPerm.MotorcycleID)
	assert.Equal(t, "branch-003", dbPerm.BranchID)
	assert.True(t, dbPerm.Active)
}

func TestFromDomain_Inactive(t *testing.T) {
	dp := &domain.DiagnosticPermission{
		ID:           "perm-004",
		MotorcycleID: "moto-004",
		BranchID:     "branch-004",
		Active:       false,
	}

	dbPerm := FromDomain(dp)

	assert.Equal(t, "perm-004", dbPerm.ID)
	assert.False(t, dbPerm.Active)
}
