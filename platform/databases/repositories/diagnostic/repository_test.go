package diagnostic

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Model Conversion Tests
// ============================================

func TestDiagnostic_ToDomain_WithNullableFields(t *testing.T) {
	now := time.Now()
	desc := "Motor hace ruido"
	sol := "Cambiar aceite"

	d := &Diagnostic{
		ID:                 "diag-001",
		MotorcycleID:       "moto-001",
		BranchID:           "branch-001",
		BranchName:         sql.NullString{String: "Taller Norte", Valid: true},
		Date:               now,
		ProblemDescription: sql.NullString{String: desc, Valid: true},
		PossibleSolution:   sql.NullString{String: sol, Valid: true},
	}

	result := d.ToDomain()

	assert.Equal(t, "diag-001", result.ID)
	assert.Equal(t, "moto-001", result.MotorcycleID)
	assert.Equal(t, "branch-001", result.BranchID)
	assert.Equal(t, "Taller Norte", result.BranchName)
	assert.Equal(t, now, result.Date)
	assert.NotNil(t, result.ProblemDescription)
	assert.Equal(t, desc, *result.ProblemDescription)
	assert.NotNil(t, result.PossibleSolution)
	assert.Equal(t, sol, *result.PossibleSolution)
}

func TestDiagnostic_ToDomain_NullFields(t *testing.T) {
	now := time.Now()

	d := &Diagnostic{
		ID:                 "diag-002",
		MotorcycleID:       "moto-002",
		BranchID:           "branch-002",
		BranchName:         sql.NullString{Valid: false},
		Date:               now,
		ProblemDescription: sql.NullString{Valid: false},
		PossibleSolution:   sql.NullString{Valid: false},
	}

	result := d.ToDomain()

	assert.Equal(t, "diag-002", result.ID)
	assert.Equal(t, "", result.BranchName)
	assert.Nil(t, result.ProblemDescription)
	assert.Nil(t, result.PossibleSolution)
}

func TestFromDomain_WithPointers(t *testing.T) {
	now := time.Now()
	desc := "Frenos desgastados"
	sol := "Cambiar pastillas"

	d := &domain.Diagnostic{
		ID:                 "diag-003",
		MotorcycleID:       "moto-003",
		BranchID:           "branch-003",
		Date:               now,
		ProblemDescription: &desc,
		PossibleSolution:   &sol,
	}

	result := FromDomain(d)

	assert.Equal(t, "diag-003", result.ID)
	assert.Equal(t, "moto-003", result.MotorcycleID)
	assert.Equal(t, "branch-003", result.BranchID)
	assert.True(t, result.ProblemDescription.Valid)
	assert.Equal(t, desc, result.ProblemDescription.String)
	assert.True(t, result.PossibleSolution.Valid)
	assert.Equal(t, sol, result.PossibleSolution.String)
}

func TestFromDomain_NilPointers(t *testing.T) {
	d := &domain.Diagnostic{
		ID:                 "diag-004",
		MotorcycleID:       "moto-004",
		BranchID:           "branch-004",
		ProblemDescription: nil,
		PossibleSolution:   nil,
	}

	result := FromDomain(d)

	assert.False(t, result.ProblemDescription.Valid)
	assert.False(t, result.PossibleSolution.Valid)
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

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.id")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.motorcycle_id")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.motorcycle_id .* AND d.branch_id")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsert")
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdate")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDelete")
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByID")
}

func TestNewRepository_PrepareError_GetByMotorcycleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.id")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.motorcycle_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleID")
}

func TestNewRepository_PrepareError_GetByMotorcycleAndBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.id")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.motorcycle_id")
	mock.ExpectPrepare("SELECT d.id, d.motorcycle_id.*WHERE d.motorcycle_id .* AND d.branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleAndBranch")
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
// GetByID Tests
// ============================================

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "branch_name", "date", "problem_description", "possible_solution"}).
		AddRow("diag-100", "moto-100", "branch-100", "Taller Norte", now, "Ruido en motor", nil)

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("diag-100").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.id = ?")

	result, err := repo.GetByID(context.Background(), "diag-100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-100", result.ID)
	assert.Equal(t, "moto-100", result.MotorcycleID)
	assert.Equal(t, "Taller Norte", result.BranchName)
	assert.NotNil(t, result.ProblemDescription)
	assert.Equal(t, "Ruido en motor", *result.ProblemDescription)
	assert.Nil(t, result.PossibleSolution)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.id = ?")

	result, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-id").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.id = ?")

	result, err := repo.GetByID(context.Background(), "error-id")

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

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "branch_name", "date", "problem_description", "possible_solution"}).
		AddRow("diag-1", "moto-200", "branch-1", "Taller Norte", now, "Problema 1", nil).
		AddRow("diag-2", "moto-200", "branch-2", "Taller Sur", now, nil, "Solución 2")

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-200").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "moto-200")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "diag-1", results[0].ID)
	assert.Equal(t, "diag-2", results[1].ID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "branch_name", "date", "problem_description", "possible_solution"})

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("no-diags").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "no-diags")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetByMotorcycleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-moto").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "error-moto")

	assert.Nil(t, results)
	assert.Error(t, err)
}

// ============================================
// GetByMotorcycleAndBranch Tests
// ============================================

func TestGetByMotorcycleAndBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "branch_name", "date", "problem_description", "possible_solution"}).
		AddRow("diag-300", "moto-300", "branch-300", "Taller Centro", now, "Problema combo", nil)

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-300", "branch-300").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ? AND d.branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-300", "branch-300")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-300", result.ID)
}

func TestGetByMotorcycleAndBranch_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-none", "branch-none").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ? AND d.branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-none", "branch-none")

	assert.Nil(t, result)
	assert.Nil(t, err) // Returns nil, nil when no rows
}

func TestGetByMotorcycleAndBranch_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT d.id, d.motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-moto", "error-branch").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution FROM diagnostics d LEFT JOIN branches b ON b.id = d.branch_id WHERE d.motorcycle_id = ? AND d.branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "error-moto", "error-branch")

	assert.Nil(t, result)
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
	mock.ExpectExec("INSERT INTO diagnostics").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	desc := "Motor hace ruido"
	diag := &domain.Diagnostic{
		ID:                 "diag-save-1",
		MotorcycleID:       "moto-save",
		BranchID:           "branch-save",
		Date:               time.Now(),
		ProblemDescription: &desc,
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, diag)

	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	diag := &domain.Diagnostic{ID: "diag-invalid-tx"}
	repo := &repository{db: db}

	err = repo.Save(context.Background(), nil, diag)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO diagnostics").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	diag := &domain.Diagnostic{
		ID:           "diag-save-err",
		MotorcycleID: "moto-err",
		BranchID:     "branch-err",
		Date:         time.Now(),
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, diag)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

// ============================================
// Update Tests
// ============================================

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE diagnostics").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	desc := "Actualizado"
	diag := &domain.Diagnostic{
		ID:                 "diag-update",
		ProblemDescription: &desc,
	}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), tx, diag)

	assert.NoError(t, err)
}

func TestUpdate_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	diag := &domain.Diagnostic{ID: "diag-invalid-tx"}
	repo := &repository{db: db}

	err = repo.Update(context.Background(), nil, diag)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestUpdate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE diagnostics").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	diag := &domain.Diagnostic{ID: "diag-update-err"}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), tx, diag)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

// ============================================
// Delete Tests
// ============================================

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM diagnostics").
		WithArgs("diag-delete").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "diag-delete")

	assert.NoError(t, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), nil, "diag-invalid-tx")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM diagnostics").
		WithArgs("not-exist").
		WillReturnResult(sqlmock.NewResult(0, 0))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "not-exist")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM diagnostics").
		WithArgs("error-delete").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "error-delete")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}
