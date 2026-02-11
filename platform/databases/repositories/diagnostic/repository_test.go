package diagnostic

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
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

	// Expect all 9 prepared statements in order
	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*FROM diagnostics.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*FROM diagnostics.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id.*FROM diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*FROM diagnostics.*WHERE motorcycle_id.*AND branch_id")
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

func TestNewRepository_PrepareError_GetEvidence(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id.*FROM diagnostic_evidence").
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
	mock.ExpectPrepare("SELECT id, diagnostic_id.*FROM diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*AND branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleAndBranch")
}

func TestNewRepository_PrepareError_DeleteEvidence(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO diagnostics")
	mock.ExpectPrepare("UPDATE diagnostics")
	mock.ExpectPrepare("DELETE FROM diagnostics")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("INSERT INTO diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, diagnostic_id.*FROM diagnostic_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*AND branch_id")
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
	rows := sqlmock.NewRows([]string{
		"id", "motorcycle_id", "branch_id", "date",
		"problem_description", "possible_solution", "sent_via_whatsapp",
	}).AddRow(
		"diag-123", "moto-456", "branch-789", now,
		"No enciende", "Revisar bujía", true,
	)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("diag-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE id = ?")

	diagnostic, err := repo.GetByID(context.Background(), "diag-123")

	assert.NoError(t, err)
	assert.NotNil(t, diagnostic)
	assert.Equal(t, "diag-123", diagnostic.ID)
	assert.Equal(t, "moto-456", diagnostic.MotorcycleID)
	assert.Equal(t, "branch-789", diagnostic.BranchID)
	assert.NotNil(t, diagnostic.ProblemDescription)
	assert.Equal(t, "No enciende", *diagnostic.ProblemDescription)
	assert.NotNil(t, diagnostic.PossibleSolution)
	assert.Equal(t, "Revisar bujía", *diagnostic.PossibleSolution)
	assert.True(t, diagnostic.SentViaWhatsApp)
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
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE id = ?")

	diagnostic, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, diagnostic)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("diag-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE id = ?")

	diagnostic, err := repo.GetByID(context.Background(), "diag-error")

	assert.Nil(t, diagnostic)
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
	rows := sqlmock.NewRows([]string{
		"id", "motorcycle_id", "branch_id", "date",
		"problem_description", "possible_solution", "sent_via_whatsapp",
	}).
		AddRow("diag-1", "moto-123", "branch-1", now, "Frenos", nil, false).
		AddRow("diag-2", "moto-123", "branch-2", now, nil, "Cambiar aceite", true)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ?")

	diagnostics, err := repo.GetByMotorcycleID(context.Background(), "moto-123")

	assert.NoError(t, err)
	assert.Len(t, diagnostics, 2)
	assert.Equal(t, "diag-1", diagnostics[0].ID)
	assert.Equal(t, "diag-2", diagnostics[1].ID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "motorcycle_id", "branch_id", "date",
		"problem_description", "possible_solution", "sent_via_whatsapp",
	})

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ?")

	diagnostics, err := repo.GetByMotorcycleID(context.Background(), "moto-empty")

	assert.NoError(t, err)
	assert.Empty(t, diagnostics)
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
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ?")

	diagnostics, err := repo.GetByMotorcycleID(context.Background(), "moto-error")

	assert.Nil(t, diagnostics)
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
	rows := sqlmock.NewRows([]string{
		"id", "motorcycle_id", "branch_id", "date",
		"problem_description", "possible_solution", "sent_via_whatsapp",
	}).AddRow("diag-combo", "moto-123", "branch-456", now, "Ruido extraño", nil, false)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123", "branch-456").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	diagnostic, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-123", "branch-456")

	assert.NoError(t, err)
	assert.NotNil(t, diagnostic)
	assert.Equal(t, "diag-combo", diagnostic.ID)
	assert.Equal(t, "moto-123", diagnostic.MotorcycleID)
}

func TestGetByMotorcycleAndBranch_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123", "branch-none").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	diagnostic, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-123", "branch-none")

	// Returns nil, nil when no rows found (not an error)
	assert.Nil(t, diagnostic)
	assert.Nil(t, err)
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
	repo.stmtGetByMotorcycleAndBranch, _ = db.Prepare("SELECT id, motorcycle_id FROM diagnostics WHERE motorcycle_id = ? AND branch_id = ?")

	diagnostic, err := repo.GetByMotorcycleAndBranch(context.Background(), "moto-err", "branch-err")

	assert.Nil(t, diagnostic)
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
	rows := sqlmock.NewRows([]string{
		"id", "diagnostic_id", "image_url", "description", "created_at",
	}).
		AddRow("evid-1", "diag-123", "https://img.com/1.jpg", "Foto frontal", now).
		AddRow("evid-2", "diag-123", "https://img.com/2.jpg", nil, now)

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("diag-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id FROM diagnostic_evidence WHERE diagnostic_id = ?")

	evidences, err := repo.GetEvidenceByDiagnosticID(context.Background(), "diag-123")

	assert.NoError(t, err)
	assert.Len(t, evidences, 2)
	assert.Equal(t, "evid-1", evidences[0].ID)
	assert.Equal(t, "https://img.com/1.jpg", evidences[0].ImageURL)
	assert.NotNil(t, evidences[0].Description)
	assert.Equal(t, "Foto frontal", *evidences[0].Description)
	assert.Nil(t, evidences[1].Description)
}

func TestGetEvidenceByDiagnosticID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "diagnostic_id", "image_url", "description", "created_at",
	})

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("diag-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id FROM diagnostic_evidence WHERE diagnostic_id = ?")

	evidences, err := repo.GetEvidenceByDiagnosticID(context.Background(), "diag-empty")

	assert.NoError(t, err)
	assert.Empty(t, evidences)
}

func TestGetEvidenceByDiagnosticID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, diagnostic_id")
	stmt.ExpectQuery().
		WithArgs("diag-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetEvidenceByDiagnosticID, _ = db.Prepare("SELECT id, diagnostic_id FROM diagnostic_evidence WHERE diagnostic_id = ?")

	evidences, err := repo.GetEvidenceByDiagnosticID(context.Background(), "diag-error")

	assert.Nil(t, evidences)
	assert.Error(t, err)
}

// ============================================
// Save Tests
// ============================================

func TestSave_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO diagnostics").
		WithArgs("diag-123", "moto-456", "branch-789", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	desc := "Motor no enciende"
	diagnostic := &domain.Diagnostic{
		ID:                 "diag-123",
		MotorcycleID:       "moto-456",
		BranchID:           "branch-789",
		Date:               time.Now(),
		ProblemDescription: &desc,
	}

	err = repo.Save(context.Background(), tx, diagnostic)
	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	repo := &repository{}

	err := repo.Save(context.Background(), nil, &domain.Diagnostic{})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO diagnostics").
		WillReturnError(sql.ErrConnDone)

	err = repo.Save(context.Background(), tx, &domain.Diagnostic{ID: "diag-err"})
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

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO diagnostic_evidence").
		WithArgs("evid-123", "diag-456", "https://img.com/photo.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	evidence := &domain.DiagnosticEvidence{
		ID:           "evid-123",
		DiagnosticID: "diag-456",
		ImageURL:     "https://img.com/photo.jpg",
		CreatedAt:    time.Now(),
	}

	err = repo.SaveEvidence(context.Background(), tx, evidence)
	assert.NoError(t, err)
}

func TestSaveEvidence_InvalidTx(t *testing.T) {
	repo := &repository{}

	err := repo.SaveEvidence(context.Background(), nil, &domain.DiagnosticEvidence{})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSaveEvidence_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO diagnostic_evidence").
		WillReturnError(sql.ErrConnDone)

	err = repo.SaveEvidence(context.Background(), tx, &domain.DiagnosticEvidence{ID: "evid-err"})
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

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE diagnostics").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), false, "diag-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	desc := "Descripción actualizada"
	diagnostic := &domain.Diagnostic{
		ID:                 "diag-123",
		ProblemDescription: &desc,
	}

	err = repo.Update(context.Background(), tx, diagnostic)
	assert.NoError(t, err)
}

func TestUpdate_InvalidTx(t *testing.T) {
	repo := &repository{}

	err := repo.Update(context.Background(), nil, &domain.Diagnostic{})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestUpdate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE diagnostics").
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(context.Background(), tx, &domain.Diagnostic{ID: "diag-err"})
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

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM diagnostics").
		WithArgs("diag-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), tx, "diag-123")
	assert.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM diagnostics").
		WithArgs("diag-not-found").
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.Delete(context.Background(), tx, "diag-not-found")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	repo := &repository{}

	err := repo.Delete(context.Background(), nil, "diag-123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM diagnostics").
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(context.Background(), tx, "diag-err")
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

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM diagnostic_evidence").
		WithArgs("diag-123").
		WillReturnResult(sqlmock.NewResult(0, 3))

	err = repo.DeleteEvidenceByDiagnosticID(context.Background(), tx, "diag-123")
	assert.NoError(t, err)
}

func TestDeleteEvidenceByDiagnosticID_InvalidTx(t *testing.T) {
	repo := &repository{}

	err := repo.DeleteEvidenceByDiagnosticID(context.Background(), nil, "diag-123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDeleteEvidenceByDiagnosticID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM diagnostic_evidence").
		WillReturnError(sql.ErrConnDone)

	err = repo.DeleteEvidenceByDiagnosticID(context.Background(), tx, "diag-err")
	assert.Error(t, err)
}

// ============================================
// Model Mapping Tests
// ============================================

func TestDiagnostic_ToDomain_AllFields(t *testing.T) {
	d := &Diagnostic{
		ID:                 "diag-001",
		MotorcycleID:       "moto-001",
		BranchID:           "branch-001",
		Date:               time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		ProblemDescription: sql.NullString{String: "Problema de frenos", Valid: true},
		PossibleSolution:   sql.NullString{String: "Cambiar pastillas", Valid: true},
		SentViaWhatsApp:    true,
	}

	dm := d.ToDomain()

	assert.Equal(t, "diag-001", dm.ID)
	assert.Equal(t, "moto-001", dm.MotorcycleID)
	assert.Equal(t, "branch-001", dm.BranchID)
	assert.NotNil(t, dm.ProblemDescription)
	assert.Equal(t, "Problema de frenos", *dm.ProblemDescription)
	assert.NotNil(t, dm.PossibleSolution)
	assert.Equal(t, "Cambiar pastillas", *dm.PossibleSolution)
	assert.True(t, dm.SentViaWhatsApp)
}

func TestDiagnostic_ToDomain_NullFields(t *testing.T) {
	d := &Diagnostic{
		ID:                 "diag-002",
		MotorcycleID:       "moto-002",
		BranchID:           "branch-002",
		Date:               time.Now(),
		ProblemDescription: sql.NullString{Valid: false},
		PossibleSolution:   sql.NullString{Valid: false},
		SentViaWhatsApp:    false,
	}

	dm := d.ToDomain()

	assert.Equal(t, "diag-002", dm.ID)
	assert.Nil(t, dm.ProblemDescription)
	assert.Nil(t, dm.PossibleSolution)
	assert.False(t, dm.SentViaWhatsApp)
}

func TestFromDomain_AllFields(t *testing.T) {
	desc := "Cadena oxidada"
	sol := "Lubricar cadena"
	diagnostic := &domain.Diagnostic{
		ID:                 "diag-003",
		MotorcycleID:       "moto-003",
		BranchID:           "branch-003",
		Date:               time.Now(),
		ProblemDescription: &desc,
		PossibleSolution:   &sol,
		SentViaWhatsApp:    true,
	}

	d := FromDomain(diagnostic)

	assert.Equal(t, "diag-003", d.ID)
	assert.Equal(t, "moto-003", d.MotorcycleID)
	assert.True(t, d.ProblemDescription.Valid)
	assert.Equal(t, "Cadena oxidada", d.ProblemDescription.String)
	assert.True(t, d.PossibleSolution.Valid)
	assert.Equal(t, "Lubricar cadena", d.PossibleSolution.String)
	assert.True(t, d.SentViaWhatsApp)
}

func TestFromDomain_NilOptionalFields(t *testing.T) {
	diagnostic := &domain.Diagnostic{
		ID:           "diag-004",
		MotorcycleID: "moto-004",
		BranchID:     "branch-004",
		Date:         time.Now(),
	}

	d := FromDomain(diagnostic)

	assert.Equal(t, "diag-004", d.ID)
	assert.False(t, d.ProblemDescription.Valid)
	assert.False(t, d.PossibleSolution.Valid)
	assert.False(t, d.SentViaWhatsApp)
}

func TestEvidence_ToDomain_AllFields(t *testing.T) {
	now := time.Now()
	e := &DiagnosticEvidence{
		ID:           "evid-001",
		DiagnosticID: "diag-001",
		ImageURL:     "https://img.com/photo.jpg",
		Description:  sql.NullString{String: "Vista lateral", Valid: true},
		CreatedAt:    now,
	}

	dm := e.EvidenceToDomain()

	assert.Equal(t, "evid-001", dm.ID)
	assert.Equal(t, "diag-001", dm.DiagnosticID)
	assert.Equal(t, "https://img.com/photo.jpg", dm.ImageURL)
	assert.NotNil(t, dm.Description)
	assert.Equal(t, "Vista lateral", *dm.Description)
}

func TestEvidence_ToDomain_NullDescription(t *testing.T) {
	e := &DiagnosticEvidence{
		ID:           "evid-002",
		DiagnosticID: "diag-002",
		ImageURL:     "https://img.com/photo2.jpg",
		Description:  sql.NullString{Valid: false},
		CreatedAt:    time.Now(),
	}

	dm := e.EvidenceToDomain()

	assert.Equal(t, "evid-002", dm.ID)
	assert.Nil(t, dm.Description)
}

func TestEvidenceFromDomain_AllFields(t *testing.T) {
	desc := "Foto de detalle"
	evidence := &domain.DiagnosticEvidence{
		ID:           "evid-003",
		DiagnosticID: "diag-003",
		ImageURL:     "https://img.com/detail.jpg",
		Description:  &desc,
		CreatedAt:    time.Now(),
	}

	e := EvidenceFromDomain(evidence)

	assert.Equal(t, "evid-003", e.ID)
	assert.Equal(t, "diag-003", e.DiagnosticID)
	assert.True(t, e.Description.Valid)
	assert.Equal(t, "Foto de detalle", e.Description.String)
}

func TestEvidenceFromDomain_NilDescription(t *testing.T) {
	evidence := &domain.DiagnosticEvidence{
		ID:           "evid-004",
		DiagnosticID: "diag-004",
		ImageURL:     "https://img.com/noDesc.jpg",
		CreatedAt:    time.Now(),
	}

	e := EvidenceFromDomain(evidence)

	assert.Equal(t, "evid-004", e.ID)
	assert.False(t, e.Description.Valid)
}
