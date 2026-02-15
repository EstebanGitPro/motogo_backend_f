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
// Model Conversion Tests
// ============================================

func TestDiagnosticPermission_ToDomain(t *testing.T) {
	p := &DiagnosticPermission{
		ID:           "perm-001",
		MotorcycleID: "moto-001",
		BranchID:     "branch-001",
		Active:       true,
	}

	result := p.ToDomain()

	assert.Equal(t, "perm-001", result.ID)
	assert.Equal(t, "moto-001", result.MotorcycleID)
	assert.Equal(t, "branch-001", result.BranchID)
	assert.True(t, result.Active)
}

func TestDiagnosticPermission_ToDomain_Inactive(t *testing.T) {
	p := &DiagnosticPermission{
		ID:           "perm-002",
		MotorcycleID: "moto-002",
		BranchID:     "branch-002",
		Active:       false,
	}

	result := p.ToDomain()

	assert.Equal(t, "perm-002", result.ID)
	assert.False(t, result.Active)
}

func TestFromDomain(t *testing.T) {
	dp := &domain.DiagnosticPermission{
		ID:           "perm-003",
		MotorcycleID: "moto-003",
		BranchID:     "branch-003",
		Active:       true,
	}

	result := FromDomain(dp)

	assert.Equal(t, "perm-003", result.ID)
	assert.Equal(t, "moto-003", result.MotorcycleID)
	assert.Equal(t, "branch-003", result.BranchID)
	assert.True(t, result.Active)
}

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
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND active")

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

func TestNewRepository_PrepareError_Revoke(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtRevoke")
}

func TestNewRepository_PrepareError_GetByMotorcycleAndBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("UPDATE motorcycle_diagnostic_permissions")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id").
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
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND active").
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
// GetByMotorcycleAndBranch Tests
// ============================================

func TestGetByMotorcycleAndBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "active"}).
		AddRow("perm-100", "moto-100", "branch-100", true)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-100", "branch-100").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-100", "branch-100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "perm-100", result.ID)
	assert.True(t, result.Active)
}

func TestGetByMotorcycleAndBranch_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-none", "branch-none").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-none", "branch-none")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrPermissionNotFound, err)
}

func TestGetByMotorcycleAndBranch_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-moto", "error-branch").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "error-moto", "error-branch")

	assert.Nil(t, result)
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
		AddRow("perm-1", "moto-200", "branch-1", true).
		AddRow("perm-2", "moto-200", "branch-2", true)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-200").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	results, err := repo.GetByMotorcycleID(context.Background(), "moto-200")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "perm-1", results[0].ID)
	assert.Equal(t, "perm-2", results[1].ID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "active"})

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("no-perms").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	results, err := repo.GetByMotorcycleID(context.Background(), "no-perms")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetByMotorcycleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-moto").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, active FROM motorcycle_diagnostic_permissions WHERE motorcycle_id = ? AND active = TRUE")

	results, err := repo.GetByMotorcycleID(context.Background(), "error-moto")

	assert.Nil(t, results)
	assert.Error(t, err)
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
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	permission := &domain.DiagnosticPermission{
		ID:           "perm-save-1",
		MotorcycleID: "moto-save",
		BranchID:     "branch-save",
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
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	permission := &domain.DiagnosticPermission{
		ID:           "perm-save-err",
		MotorcycleID: "moto-err",
		BranchID:     "branch-err",
		Active:       true,
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, permission)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionCannotSave, err)
}

// ============================================
// Delete (Revoke) Tests
// ============================================

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_diagnostic_permissions").
		WithArgs("moto-del", "branch-del").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "moto-del", "branch-del")

	assert.NoError(t, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), nil, "moto-invalid", "branch-invalid")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_diagnostic_permissions").
		WithArgs("moto-none", "branch-none").
		WillReturnResult(sqlmock.NewResult(0, 0))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "moto-none", "branch-none")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionNotFound, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_diagnostic_permissions").
		WithArgs("error-moto", "error-branch").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "error-moto", "error-branch")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionCannotDelete, err)
}
