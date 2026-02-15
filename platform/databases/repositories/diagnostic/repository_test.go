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
		Date:               now,
		ProblemDescription: sql.NullString{String: desc, Valid: true},
		PossibleSolution:   sql.NullString{String: sol, Valid: true},
	}

	result := d.ToDomain()

	assert.Equal(t, "diag-001", result.ID)
	assert.Equal(t, "moto-001", result.MotorcycleID)
	assert.Equal(t, "branch-001", result.BranchID)
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
		Date:               now,
		ProblemDescription: sql.NullString{Valid: false},
		PossibleSolution:   sql.NullString{Valid: false},
	}

	result := d.ToDomain()

	assert.Equal(t, "diag-002", result.ID)
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

func TestDiagnosticEvidence_EvidenceToDomain_WithDescription(t *testing.T) {
	now := time.Now()

	e := &DiagnosticEvidence{
		ID:           "evid-001",
		DiagnosticID: "diag-001",
		ImageURL:     "https://storage/img.jpg",
		Description:  sql.NullString{String: "Foto frontal", Valid: true},
		CreatedAt:    now,
	}

	result := e.EvidenceToDomain()

	assert.Equal(t, "evid-001", result.ID)
	assert.Equal(t, "diag-001", result.DiagnosticID)
	assert.Equal(t, "https://storage/img.jpg", result.ImageURL)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Foto frontal", *result.Description)
	assert.Equal(t, now, result.CreatedAt)
}

func TestDiagnosticEvidence_EvidenceToDomain_NullDescription(t *testing.T) {
	now := time.Now()

	e := &DiagnosticEvidence{
		ID:           "evid-002",
		DiagnosticID: "diag-002",
		ImageURL:     "https://storage/img2.jpg",
		Description:  sql.NullString{Valid: false},
		CreatedAt:    now,
	}

	result := e.EvidenceToDomain()

	assert.Equal(t, "evid-002", result.ID)
	assert.Nil(t, result.Description)
}

func TestEvidenceFromDomain_WithDescription(t *testing.T) {
	now := time.Now()
	desc := "Evidencia lateral"

	e := &domain.DiagnosticEvidence{
		ID:           "evid-003",
		DiagnosticID: "diag-003",
		ImageURL:     "https://storage/img3.jpg",
		Description:  &desc,
		CreatedAt:    now,
	}

	result := EvidenceFromDomain(e)

	assert.Equal(t, "evid-003", result.ID)
	assert.Equal(t, "diag-003", result.DiagnosticID)
	assert.Equal(t, "https://storage/img3.jpg", result.ImageURL)
	assert.True(t, result.Description.Valid)
	assert.Equal(t, desc, result.Description.String)
	assert.Equal(t, now, result.CreatedAt)
}

func TestEvidenceFromDomain_NilDescription(t *testing.T) {
	e := &domain.DiagnosticEvidence{
		ID:           "evid-004",
		DiagnosticID: "diag-004",
		ImageURL:     "https://storage/img4.jpg",
		Description:  nil,
	}

	result := EvidenceFromDomain(e)

	assert.False(t, result.Description.Valid)
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
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id")
	mock.ExpectPrepare("DELETE FROM diagnostic_evidence")

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
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id").
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
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleID")
}

func TestNewRepository_PrepareError_InsertEvidence(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsertEvidence")
}

func TestNewRepository_PrepareError_GetEvidenceByDiagnosticID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetEvidenceByDiagnosticID")
}

func TestNewRepository_PrepareError_GetByMotorcycleAndBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleAndBranch")
}

func TestNewRepository_PrepareError_DeleteEvidenceByDiagnosticID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id .* AND branch_id")
	mock.ExpectPrepare("DELETE FROM diagnostic_evidence").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeleteEvidenceByDiagnosticID")
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
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "date", "problem_description", "possible_solution"}).
		AddRow("diag-100", "moto-100", "branch-100", now, "Ruido en motor", nil)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("diag-100").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE id = ?")

	result, err := repo.GetByID(context.Background(), "diag-100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-100", result.ID)
	assert.Equal(t, "moto-100", result.MotorcycleID)
	assert.NotNil(t, result.ProblemDescription)
	assert.Equal(t, "Ruido en motor", *result.ProblemDescription)
	assert.Nil(t, result.PossibleSolution)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE id = ?")

	result, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("error-id").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE id = ?")

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
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "date", "problem_description", "possible_solution"}).
		AddRow("diag-1", "moto-200", "branch-1", now, "Problema 1", nil).
		AddRow("diag-2", "moto-200", "branch-2", now, nil, "Solución 2")

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-200").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ?")

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

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "date", "problem_description", "possible_solution"})

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("no-diags").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "no-diags")

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
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ?")

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
	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "branch_id", "date", "problem_description", "possible_solution"}).
		AddRow("diag-300", "moto-300", "branch-300", now, "Problema combo", nil)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-300", "branch-300").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-300", "branch-300")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-300", result.ID)
}

func TestGetByMotorcycleAndBranch_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-none", "branch-none").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-none", "branch-none")

	assert.Nil(t, result)
	assert.Nil(t, err) // Returns nil, nil when no rows
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
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	result, err := repo.GetByMotorcycleAndBranch(context.Background(), "error-moto", "error-branch")

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// GetEvidenceByDiagnosticID Tests
// ============================================

func TestGetEvidenceByDiagnosticID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "diagnostic_id", "image_url", "description", "created_at"}).
		AddRow("evid-1", "diag-500", "https://storage/img1.jpg", "Foto frontal", now).
		AddRow("evid-2", "diag-500", "https://storage/img2.jpg", nil, now)

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("diag-500").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id, image_url, description, created_at FROM diagnostic_evidence WHERE diagnostic_id = ?")

	results, err := repo.GetEvidenceByDiagnosticID(context.Background(), "diag-500")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "evid-1", results[0].ID)
	assert.NotNil(t, results[0].Description)
	assert.Equal(t, "evid-2", results[1].ID)
	assert.Nil(t, results[1].Description)
}

func TestGetEvidenceByDiagnosticID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "diagnostic_id", "image_url", "description", "created_at"})

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("no-evid").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id, image_url, description, created_at FROM diagnostic_evidence WHERE diagnostic_id = ?")

	results, err := repo.GetEvidenceByDiagnosticID(context.Background(), "no-evid")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetEvidenceByDiagnosticID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("error-diag").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id, image_url, description, created_at FROM diagnostic_evidence WHERE diagnostic_id = ?")

	results, err := repo.GetEvidenceByDiagnosticID(context.Background(), "error-diag")

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
// SaveEvidence Tests
// ============================================

func TestSaveEvidence_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO diagnostic_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	evidence := &domain.DiagnosticEvidence{
		ID:           "evid-save-1",
		DiagnosticID: "diag-save",
		ImageURL:     "https://storage/evidence.jpg",
		CreatedAt:    time.Now(),
	}

	repo := &repository{db: db}

	err = repo.SaveEvidence(context.Background(), tx, evidence)

	assert.NoError(t, err)
}

func TestSaveEvidence_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	evidence := &domain.DiagnosticEvidence{ID: "evid-invalid-tx"}
	repo := &repository{db: db}

	err = repo.SaveEvidence(context.Background(), nil, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSaveEvidence_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO diagnostic_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	evidence := &domain.DiagnosticEvidence{
		ID:           "evid-save-err",
		DiagnosticID: "diag-err",
		ImageURL:     "https://storage/err.jpg",
		CreatedAt:    time.Now(),
	}

	repo := &repository{db: db}

	err = repo.SaveEvidence(context.Background(), tx, evidence)

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

// ============================================
// DeleteEvidenceByDiagnosticID Tests
// ============================================

func TestDeleteEvidenceByDiagnosticID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM diagnostic_evidence").
		WithArgs("diag-del-evid").
		WillReturnResult(sqlmock.NewResult(0, 3))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.DeleteEvidenceByDiagnosticID(context.Background(), tx, "diag-del-evid")

	assert.NoError(t, err)
}

func TestDeleteEvidenceByDiagnosticID_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.DeleteEvidenceByDiagnosticID(context.Background(), nil, "diag-invalid")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDeleteEvidenceByDiagnosticID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM diagnostic_evidence").
		WithArgs("error-del").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.DeleteEvidenceByDiagnosticID(context.Background(), tx, "error-del")

	assert.Error(t, err)
}
