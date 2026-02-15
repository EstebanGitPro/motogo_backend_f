package evidence

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

	// Expect all prepared statements
	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence")
	mock.ExpectPrepare("DELETE FROM motorcycle_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("SELECT COUNT")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence").
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

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence").
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

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence")
	mock.ExpectPrepare("DELETE FROM motorcycle_evidence").
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

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence")
	mock.ExpectPrepare("DELETE FROM motorcycle_evidence")
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

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence")
	mock.ExpectPrepare("DELETE FROM motorcycle_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleID")
}

func TestNewRepository_PrepareError_CountByMotorcycleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO motorcycle_evidence")
	mock.ExpectPrepare("UPDATE motorcycle_evidence")
	mock.ExpectPrepare("DELETE FROM motorcycle_evidence")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT id, motorcycle_id.*WHERE motorcycle_id")
	mock.ExpectPrepare("SELECT COUNT").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtCountByMotorcycleID")
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

	uploadDate := time.Now()
	angle := "FRONT"

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "angle", "image_url", "description", "created_at"}).
		AddRow("evidence-123", "moto-456", angle, "https://firebasestorage.googleapis.com/v0/b/test/image.jpg", nil, uploadDate)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("evidence-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE id = ?")

	evidence, err := repo.GetByID(context.Background(), "evidence-123")

	assert.NoError(t, err)
	assert.NotNil(t, evidence)
	assert.Equal(t, "evidence-123", evidence.ID)
	assert.Equal(t, "moto-456", evidence.MotorcycleID)
	assert.NotNil(t, evidence.Angle)
	assert.Equal(t, domain.EvidenceAngle("FRONT"), *evidence.Angle)
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
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE id = ?")

	evidence, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, evidence)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
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
	repo.stmtGetByID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE id = ?")

	evidence, err := repo.GetByID(context.Background(), "error-id")

	assert.Nil(t, evidence)
	assert.Error(t, err)
}

// ============================================
// GetByMotorcycleID Tests
// ============================================

func TestGetByMotorcycleID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	uploadDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "angle", "image_url", "description", "created_at"}).
		AddRow("evidence-1", "moto-123", "FRONT", "https://storage/img1.jpg", nil, uploadDate).
		AddRow("evidence-2", "moto-123", "SIDE", "https://storage/img2.jpg", nil, uploadDate)

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("moto-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE motorcycle_id = ?")

	evidences, err := repo.GetByMotorcycleID(context.Background(), "moto-123")

	assert.NoError(t, err)
	assert.Len(t, evidences, 2)
	assert.Equal(t, "evidence-1", evidences[0].ID)
	assert.Equal(t, "evidence-2", evidences[1].ID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "motorcycle_id", "angle", "image_url", "description", "created_at"})

	stmt := mock.ExpectPrepare("SELECT id, motorcycle_id")
	stmt.ExpectQuery().
		WithArgs("no-photos").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE motorcycle_id = ?")

	evidences, err := repo.GetByMotorcycleID(context.Background(), "no-photos")

	assert.NoError(t, err)
	assert.Empty(t, evidences)
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
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT id, motorcycle_id, angle, image_url, description, created_at FROM motorcycle_evidence WHERE motorcycle_id = ?")

	evidences, err := repo.GetByMotorcycleID(context.Background(), "error-moto")

	assert.Nil(t, evidences)
	assert.Error(t, err)
}

// ============================================
// CountByMotorcycleID Tests
// ============================================

func TestCountByMotorcycleID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("moto-count").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCountByMotorcycleID, _ = db.Prepare("SELECT COUNT(*) FROM motorcycle_evidence WHERE motorcycle_id = ?")

	count, err := repo.CountByMotorcycleID(context.Background(), "moto-count")

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestCountByMotorcycleID_Zero(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("moto-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCountByMotorcycleID, _ = db.Prepare("SELECT COUNT(*) FROM motorcycle_evidence WHERE motorcycle_id = ?")

	count, err := repo.CountByMotorcycleID(context.Background(), "moto-empty")

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCountByMotorcycleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("error-count").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCountByMotorcycleID, _ = db.Prepare("SELECT COUNT(*) FROM motorcycle_evidence WHERE motorcycle_id = ?")

	count, err := repo.CountByMotorcycleID(context.Background(), "error-count")

	assert.Equal(t, 0, count)
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
	mock.ExpectExec("INSERT INTO motorcycle_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	angle := domain.EvidenceAngle("FRONT")
	evidence := &domain.MotorcycleEvidence{
		ID:           "evidence-new",
		MotorcycleID: "moto-123",
		Angle:        &angle,
		ImageURL:     "https://storage/image.jpg",
		CreatedAt:    time.Now(),
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, evidence)

	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	evidence := &domain.MotorcycleEvidence{
		ID: "evidence-invalid-tx",
	}

	repo := &repository{db: db}

	// Pass a nil tx or invalid type
	err = repo.Save(context.Background(), nil, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO motorcycle_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	evidence := &domain.MotorcycleEvidence{
		ID:           "evidence-error",
		MotorcycleID: "moto-123",
		ImageURL:     "https://storage/image.jpg",
		CreatedAt:    time.Now(),
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

// ============================================
// Update Tests
// ============================================

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	angle := domain.EvidenceAngle("BACK")
	evidence := &domain.MotorcycleEvidence{
		ID:       "evidence-update",
		Angle:    &angle,
		ImageURL: "https://storage/updated.jpg",
	}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), tx, evidence)

	assert.NoError(t, err)
}

func TestUpdate_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	evidence := &domain.MotorcycleEvidence{
		ID: "evidence-invalid-tx",
	}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), nil, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestUpdate_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	evidence := &domain.MotorcycleEvidence{
		ID:       "not-exist",
		ImageURL: "https://storage/notfound.jpg",
	}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), tx, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestUpdate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE motorcycle_evidence").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	evidence := &domain.MotorcycleEvidence{
		ID:       "error-update",
		ImageURL: "https://storage/error.jpg",
	}

	repo := &repository{db: db}

	err = repo.Update(context.Background(), tx, evidence)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotUpdate, err)
}

// ============================================
// Delete Tests
// ============================================

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM motorcycle_evidence").
		WithArgs("evidence-delete").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "evidence-delete")

	assert.NoError(t, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), nil, "evidence-invalid-tx")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM motorcycle_evidence").
		WithArgs("not-exist").
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "not-exist")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM motorcycle_evidence").
		WithArgs("error-delete").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "error-delete")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}

// ============================================
// Evidence Model Tests
// ============================================

func TestEvidence_ToDomain(t *testing.T) {
	uploadDate := time.Now()

	e := &Evidence{
		ID:           "evidence-001",
		MotorcycleID: "moto-001",
		Angle:        sql.NullString{String: "FRONT", Valid: true},
		ImageURL:     "https://storage/image.jpg",
		CreatedAt:    uploadDate,
	}

	de := e.ToDomain()

	assert.Equal(t, "evidence-001", de.ID)
	assert.Equal(t, "moto-001", de.MotorcycleID)
	assert.NotNil(t, de.Angle)
	assert.Equal(t, domain.EvidenceAngle("FRONT"), *de.Angle)
	assert.Equal(t, "https://storage/image.jpg", de.ImageURL)
	assert.Equal(t, uploadDate, de.CreatedAt)
}

func TestEvidence_ToDomain_NullAngle(t *testing.T) {
	uploadDate := time.Now()

	e := &Evidence{
		ID:           "evidence-002",
		MotorcycleID: "moto-002",
		Angle:        sql.NullString{Valid: false},
		ImageURL:     "https://storage/image2.jpg",
		CreatedAt:    uploadDate,
	}

	de := e.ToDomain()

	assert.Equal(t, "evidence-002", de.ID)
	assert.Nil(t, de.Angle)
}

func TestFromDomain(t *testing.T) {
	uploadDate := time.Now()
	angle := domain.EvidenceAngle("SIDE")

	de := &domain.MotorcycleEvidence{
		ID:           "evidence-003",
		MotorcycleID: "moto-003",
		Angle:        &angle,
		ImageURL:     "https://storage/image3.jpg",
		CreatedAt:    uploadDate,
	}

	e := FromDomain(de)

	assert.Equal(t, "evidence-003", e.ID)
	assert.Equal(t, "moto-003", e.MotorcycleID)
	assert.True(t, e.Angle.Valid)
	assert.Equal(t, "SIDE", e.Angle.String)
	assert.Equal(t, "https://storage/image3.jpg", e.ImageURL)
	assert.Equal(t, uploadDate, e.CreatedAt)
}

func TestFromDomain_NilAngle(t *testing.T) {
	uploadDate := time.Now()

	de := &domain.MotorcycleEvidence{
		ID:           "evidence-004",
		MotorcycleID: "moto-004",
		Angle:        nil,
		ImageURL:     "https://storage/image4.jpg",
		CreatedAt:    uploadDate,
	}

	e := FromDomain(de)

	assert.Equal(t, "evidence-004", e.ID)
	assert.False(t, e.Angle.Valid)
}
