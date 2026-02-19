package completed_service

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

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services")
	mock.ExpectPrepare("DELETE FROM completed_services")
	mock.ExpectPrepare("UPDATE completed_services SET deleted_at")
	mock.ExpectPrepare("UPDATE completed_services")
	mock.ExpectPrepare("SELECT id, completed_service_id, previous_status")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsert")
}

func TestNewRepository_PrepareError_InsertItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsertItem")
}

func TestNewRepository_PrepareError_InsertStatusHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtInsertStatusHistory")
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id").
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

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByMotorcycleID")
}

func TestNewRepository_PrepareError_GetByBranchID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByBranchID")
}

func TestNewRepository_PrepareError_GetItemsByCSID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetItemsByCSID")
}

func TestNewRepository_PrepareError_HasActiveService(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtHasActiveService")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services")
	mock.ExpectPrepare("DELETE FROM completed_services").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDelete")
}

func TestNewRepository_PrepareError_SoftDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services")
	mock.ExpectPrepare("DELETE FROM completed_services")
	mock.ExpectPrepare("UPDATE completed_services SET deleted_at").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSoftDelete")
}

func TestNewRepository_PrepareError_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services")
	mock.ExpectPrepare("DELETE FROM completed_services")
	mock.ExpectPrepare("UPDATE completed_services SET deleted_at")
	mock.ExpectPrepare("UPDATE completed_services").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdateStatus")
}

func TestNewRepository_PrepareError_GetStatusHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO completed_services")
	mock.ExpectPrepare("INSERT INTO completed_service_items")
	mock.ExpectPrepare("INSERT INTO service_status_transitions")
	mock.ExpectPrepare("SELECT id, branch_id, motorcycle_id.*WHERE id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.motorcycle_id")
	mock.ExpectPrepare("SELECT cs.id.*WHERE cs.branch_id")
	mock.ExpectPrepare("SELECT csi.id, csi.completed_service_id, csi.service_id")
	mock.ExpectPrepare("SELECT COUNT.*FROM completed_services")
	mock.ExpectPrepare("DELETE FROM completed_services")
	mock.ExpectPrepare("UPDATE completed_services SET deleted_at")
	mock.ExpectPrepare("UPDATE completed_services")
	mock.ExpectPrepare("SELECT id, completed_service_id, previous_status").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetStatusHistory")
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
		"id", "branch_id", "motorcycle_id", "diagnostic_id",
		"request_date", "completion_date", "status",
		"quoted_price", "final_price", "representative_notes",
		"created_at", "updated_at",
	}).AddRow(
		"cs-001", "branch-001", "moto-001", nil,
		now, nil, "PENDIENTE",
		50000.0, nil, "Revisión general",
		now, now,
	)

	stmt := mock.ExpectPrepare("SELECT id, branch_id")
	stmt.ExpectQuery().
		WithArgs("cs-001").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, branch_id, motorcycle_id, diagnostic_id, request_date, completion_date, status, quoted_price, final_price, representative_notes, created_at, updated_at FROM completed_services WHERE id = ?")

	result, err := repo.GetByID(context.Background(), "cs-001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cs-001", result.ID)
	assert.Equal(t, "branch-001", result.BranchID)
	assert.Equal(t, "moto-001", result.MotorcycleID)
	assert.Nil(t, result.DiagnosticID)
	assert.Equal(t, domain.ServiceStatus("PENDIENTE"), result.Status)
	assert.NotNil(t, result.QuotedPrice)
	assert.Equal(t, 50000.0, *result.QuotedPrice)
	assert.Nil(t, result.FinalPrice)
	assert.NotNil(t, result.RepresentativeNotes)
	assert.Equal(t, "Revisión general", *result.RepresentativeNotes)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, branch_id, motorcycle_id FROM completed_services WHERE id = ?")

	result, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrCompletedServiceNotFound, err)
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id")
	stmt.ExpectQuery().
		WithArgs("error-id").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT id, branch_id, motorcycle_id FROM completed_services WHERE id = ?")

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
	rows := sqlmock.NewRows([]string{
		"id", "branch_id", "branch_name",
		"motorcycle_id", "diagnostic_id",
		"request_date", "completion_date", "status",
		"quoted_price", "final_price", "representative_notes",
		"created_at", "updated_at",
	}).
		AddRow("cs-1", "b-1", "Sede Norte", "moto-200", nil, now, nil, "PENDIENTE", 30000.0, nil, nil, now, now).
		AddRow("cs-2", "b-2", "Sede Sur", "moto-200", "diag-1", now, now, "FINALIZADO", 40000.0, 38000.0, "Todo bien", now, now)

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("moto-200").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT cs.id, cs.branch_id, b.name AS branch_name, cs.motorcycle_id, cs.diagnostic_id, cs.request_date, cs.completion_date, cs.status, cs.quoted_price, cs.final_price, cs.representative_notes, cs.created_at, cs.updated_at FROM completed_services cs LEFT JOIN branches b ON b.id = cs.branch_id WHERE cs.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "moto-200")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "cs-1", results[0].ID)
	assert.NotNil(t, results[0].BranchName)
	assert.Equal(t, "Sede Norte", *results[0].BranchName)
	assert.Equal(t, "cs-2", results[1].ID)
	assert.NotNil(t, results[1].DiagnosticID)
	assert.Equal(t, "diag-1", *results[1].DiagnosticID)
}

func TestGetByMotorcycleID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "branch_id", "branch_name",
		"motorcycle_id", "diagnostic_id",
		"request_date", "completion_date", "status",
		"quoted_price", "final_price", "representative_notes",
		"created_at", "updated_at",
	})

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("no-services").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT cs.id FROM completed_services cs WHERE cs.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "no-services")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetByMotorcycleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("error-moto").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByMotorcycleID, _ = db.Prepare("SELECT cs.id FROM completed_services cs WHERE cs.motorcycle_id = ?")

	results, err := repo.GetByMotorcycleID(context.Background(), "error-moto")

	assert.Nil(t, results)
	assert.Error(t, err)
}

// ============================================
// GetByBranchID Tests
// ============================================

func TestGetByBranchID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "branch_id", "branch_name",
		"motorcycle_id", "diagnostic_id",
		"request_date", "completion_date", "status",
		"quoted_price", "final_price", "representative_notes",
		"created_at", "updated_at",
	}).AddRow("cs-10", "branch-10", "Sede Principal", "moto-10", nil, now, nil, "EN_PROCESO", 25000.0, nil, "En revisión", now, now)

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("branch-10").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByBranchID, _ = db.Prepare("SELECT cs.id, cs.branch_id, b.name AS branch_name FROM completed_services cs LEFT JOIN branches b ON b.id = cs.branch_id WHERE cs.branch_id = ?")

	results, err := repo.GetByBranchID(context.Background(), "branch-10")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "cs-10", results[0].ID)
	assert.Equal(t, domain.ServiceStatus("EN_PROCESO"), results[0].Status)
}

func TestGetByBranchID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "branch_id", "branch_name",
		"motorcycle_id", "diagnostic_id",
		"request_date", "completion_date", "status",
		"quoted_price", "final_price", "representative_notes",
		"created_at", "updated_at",
	})

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("empty-branch").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByBranchID, _ = db.Prepare("SELECT cs.id FROM completed_services cs WHERE cs.branch_id = ?")

	results, err := repo.GetByBranchID(context.Background(), "empty-branch")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetByBranchID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT cs.id")
	stmt.ExpectQuery().
		WithArgs("error-branch").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByBranchID, _ = db.Prepare("SELECT cs.id FROM completed_services cs WHERE cs.branch_id = ?")

	results, err := repo.GetByBranchID(context.Background(), "error-branch")

	assert.Nil(t, results)
	assert.Error(t, err)
}

// ============================================
// GetItemsByCompletedServiceID Tests
// ============================================

func TestGetItemsByCompletedServiceID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "completed_service_id", "service_id",
		"service_name",
		"rating", "comment", "rated_at", "is_offensive_comment",
	}).
		AddRow("item-1", "cs-500", "svc-1", "Cambio de aceite", 5, "Excelente servicio", now, false).
		AddRow("item-2", "cs-500", "svc-2", "Revisión de frenos", nil, nil, nil, false)

	stmt := mock.ExpectPrepare("SELECT csi.id")
	stmt.ExpectQuery().
		WithArgs("cs-500").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetItemsByCSID, _ = db.Prepare("SELECT csi.id, csi.completed_service_id, csi.service_id, s.name AS service_name, csi.rating, csi.comment, csi.rated_at, csi.is_offensive_comment FROM completed_service_items csi LEFT JOIN services s ON s.id = csi.service_id WHERE csi.completed_service_id = ?")

	results, err := repo.GetItemsByCompletedServiceID(context.Background(), "cs-500")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "item-1", results[0].ID)
	assert.NotNil(t, results[0].ServiceName)
	assert.Equal(t, "Cambio de aceite", *results[0].ServiceName)
	assert.NotNil(t, results[0].Rating)
	assert.Equal(t, 5, *results[0].Rating)
	assert.NotNil(t, results[0].Comment)
	assert.Equal(t, "Excelente servicio", *results[0].Comment)
	assert.Equal(t, "item-2", results[1].ID)
	assert.NotNil(t, results[1].ServiceName)
	assert.Equal(t, "Revisión de frenos", *results[1].ServiceName)
	assert.Nil(t, results[1].Rating)
	assert.Nil(t, results[1].Comment)
}

func TestGetItemsByCompletedServiceID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "completed_service_id", "service_id",
		"service_name",
		"rating", "comment", "rated_at", "is_offensive_comment",
	})

	stmt := mock.ExpectPrepare("SELECT csi.id")
	stmt.ExpectQuery().
		WithArgs("cs-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetItemsByCSID, _ = db.Prepare("SELECT csi.id, csi.completed_service_id, csi.service_id, s.name AS service_name, csi.rating, csi.comment, csi.rated_at, csi.is_offensive_comment FROM completed_service_items csi LEFT JOIN services s ON s.id = csi.service_id WHERE csi.completed_service_id = ?")

	results, err := repo.GetItemsByCompletedServiceID(context.Background(), "cs-empty")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetItemsByCompletedServiceID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT csi.id")
	stmt.ExpectQuery().
		WithArgs("cs-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetItemsByCSID, _ = db.Prepare("SELECT csi.id, csi.completed_service_id, csi.service_id, s.name AS service_name FROM completed_service_items csi LEFT JOIN services s ON s.id = csi.service_id WHERE csi.completed_service_id = ?")

	results, err := repo.GetItemsByCompletedServiceID(context.Background(), "cs-error")

	assert.Nil(t, results)
	assert.Error(t, err)
}

// ============================================
// HasActiveService Tests
// ============================================

func TestHasActiveService_ReturnsTrue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-1", "moto-1").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtHasActiveService, _ = db.Prepare("SELECT COUNT(*) FROM completed_services WHERE branch_id = ? AND motorcycle_id = ? AND status IN ('PENDIENTE', 'EN_PROCESO')")

	result, err := repo.HasActiveService(context.Background(), "branch-1", "moto-1")

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestHasActiveService_ReturnsFalse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-2", "moto-2").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtHasActiveService, _ = db.Prepare("SELECT COUNT(*) FROM completed_services WHERE branch_id = ? AND motorcycle_id = ?")

	result, err := repo.HasActiveService(context.Background(), "branch-2", "moto-2")

	assert.NoError(t, err)
	assert.False(t, result)
}

func TestHasActiveService_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-err", "moto-err").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtHasActiveService, _ = db.Prepare("SELECT COUNT(*) FROM completed_services WHERE branch_id = ? AND motorcycle_id = ?")

	result, err := repo.HasActiveService(context.Background(), "branch-err", "moto-err")

	assert.False(t, result)
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
	mock.ExpectExec("INSERT INTO completed_services").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	price := 50000.0
	finalPrice := 48000.0
	notes := "Revisión completa"
	cs := &domain.CompletedService{
		ID:                  "cs-save-1",
		BranchID:            "branch-save",
		MotorcycleID:        "moto-save",
		RequestDate:         time.Now(),
		Status:              domain.ServiceStatusPending,
		QuotedPrice:         &price,
		FinalPrice:          &finalPrice,
		RepresentativeNotes: &notes,
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, cs)

	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	cs := &domain.CompletedService{ID: "cs-invalid-tx"}
	repo := &repository{db: db}

	err = repo.Save(context.Background(), nil, cs)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO completed_services").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	cs := &domain.CompletedService{
		ID:           "cs-save-err",
		BranchID:     "branch-err",
		MotorcycleID: "moto-err",
		RequestDate:  time.Now(),
		Status:       domain.ServiceStatusPending,
	}

	repo := &repository{db: db}

	err = repo.Save(context.Background(), tx, cs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error saving completed service")
}

// ============================================
// SaveItems Tests
// ============================================

func TestSaveItems_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO completed_service_items").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO completed_service_items").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(2, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	items := []domain.CompletedServiceItem{
		{ID: "item-1", CompletedServiceID: "cs-1", ServiceID: "svc-1"},
		{ID: "item-2", CompletedServiceID: "cs-1", ServiceID: "svc-2"},
	}

	repo := &repository{db: db}

	err = repo.SaveItems(context.Background(), tx, items)

	assert.NoError(t, err)
}

func TestSaveItems_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	items := []domain.CompletedServiceItem{
		{ID: "item-invalid"},
	}

	repo := &repository{db: db}

	err = repo.SaveItems(context.Background(), nil, items)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSaveItems_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO completed_service_items").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	items := []domain.CompletedServiceItem{
		{ID: "item-err", CompletedServiceID: "cs-err", ServiceID: "svc-err"},
	}

	repo := &repository{db: db}

	err = repo.SaveItems(context.Background(), tx, items)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error saving completed service item")
}

// ============================================
// SaveStatusHistory Tests
// ============================================

func TestSaveStatusHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO service_status_transitions").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	prevStatus := domain.ServiceStatusPending
	history := &domain.ServiceStatusHistory{
		ID:                 "hist-1",
		CompletedServiceID: "cs-1",
		PreviousStatus:     &prevStatus,
		NewStatus:          domain.ServiceStatusInProgress,
		CreatedBy:          "user-001",
	}

	repo := &repository{db: db}

	err = repo.SaveStatusHistory(context.Background(), tx, history)

	assert.NoError(t, err)
}

func TestSaveStatusHistory_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	history := &domain.ServiceStatusHistory{ID: "hist-invalid"}
	repo := &repository{db: db}

	err = repo.SaveStatusHistory(context.Background(), nil, history)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSaveStatusHistory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO service_status_transitions").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	history := &domain.ServiceStatusHistory{
		ID:                 "hist-err",
		CompletedServiceID: "cs-err",
		NewStatus:          domain.ServiceStatusPending,
		CreatedBy:          "user-err",
	}

	repo := &repository{db: db}

	err = repo.SaveStatusHistory(context.Background(), tx, history)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error saving status history")
}

// ============================================
// ValidateBranchServices Tests
// ============================================

func TestValidateBranchServices_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("branch-1", "svc-1", "svc-2").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBranchServices(context.Background(), "branch-1", []string{"svc-1", "svc-2"})

	assert.NoError(t, err)
}

func TestValidateBranchServices_CountMismatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1) // Only 1 of 2 found

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("branch-1", "svc-1", "svc-invalid").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBranchServices(context.Background(), "branch-1", []string{"svc-1", "svc-invalid"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "some services are not associated with this branch")
}

func TestValidateBranchServices_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("branch-err", "svc-err").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	err = repo.ValidateBranchServices(context.Background(), "branch-err", []string{"svc-err"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error validating branch services")
}

// ============================================
// Delete Tests
// ============================================

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM completed_services").
		WithArgs("cs-delete-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "cs-delete-1")

	assert.NoError(t, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), nil, "cs-invalid")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM completed_services").
		WithArgs("cs-err").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.Delete(context.Background(), tx, "cs-err")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting completed service")
}

// ============================================
// SoftDelete Tests
// ============================================

func TestSoftDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_services SET deleted_at").
		WithArgs("cs-soft-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.SoftDelete(context.Background(), tx, "cs-soft-1")

	assert.NoError(t, err)
}

func TestSoftDelete_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.SoftDelete(context.Background(), nil, "cs-invalid")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSoftDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_services SET deleted_at").
		WithArgs("cs-soft-err").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.SoftDelete(context.Background(), tx, "cs-soft-err")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error soft-deleting completed service")
}

// ============================================
// UpdateStatus Tests
// ============================================

func TestUpdateStatus_SuccessWithCompletionDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_services").
		WithArgs("FINALIZADO", sqlmock.AnyArg(), "cs-status-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	now := time.Now()
	repo := &repository{db: db}

	err = repo.UpdateStatus(context.Background(), tx, "cs-status-1", "FINALIZADO", &now)

	assert.NoError(t, err)
}

func TestUpdateStatus_SuccessWithoutCompletionDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_services").
		WithArgs("EN_PROCESO", sqlmock.AnyArg(), "cs-status-2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.UpdateStatus(context.Background(), tx, "cs-status-2", "EN_PROCESO", nil)

	assert.NoError(t, err)
}

func TestUpdateStatus_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.UpdateStatus(context.Background(), nil, "cs-invalid", "EN_PROCESO", nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestUpdateStatus_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_services").
		WithArgs("CANCELADO", sqlmock.AnyArg(), "cs-err").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.UpdateStatus(context.Background(), tx, "cs-err", "CANCELADO", nil)

	assert.Error(t, err)
}

// ============================================
// GetStatusHistory Tests
// ============================================

func TestGetStatusHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "completed_service_id", "previous_status", "new_status", "created_by", "created_at",
	}).
		AddRow("hist-1", "cs-hist", nil, "PENDIENTE", "user-1", now).
		AddRow("hist-2", "cs-hist", "PENDIENTE", "EN_PROCESO", "user-1", now)

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("cs-hist").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetStatusHistory, _ = db.Prepare("SELECT id, completed_service_id, previous_status, new_status, created_by, created_at FROM service_status_transitions WHERE completed_service_id = ?")

	results, err := repo.GetStatusHistory(context.Background(), "cs-hist")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "hist-1", results[0].ID)
	assert.Nil(t, results[0].PreviousStatus)
	assert.Equal(t, domain.ServiceStatus("PENDIENTE"), results[0].NewStatus)
	assert.Equal(t, "hist-2", results[1].ID)
	assert.NotNil(t, results[1].PreviousStatus)
	assert.Equal(t, domain.ServiceStatus("PENDIENTE"), *results[1].PreviousStatus)
	assert.Equal(t, domain.ServiceStatus("EN_PROCESO"), results[1].NewStatus)
}

func TestGetStatusHistory_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "completed_service_id", "previous_status", "new_status", "created_by", "created_at",
	})

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("cs-no-hist").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetStatusHistory, _ = db.Prepare("SELECT id, completed_service_id, previous_status, new_status, created_by, created_at FROM service_status_transitions WHERE completed_service_id = ?")

	results, err := repo.GetStatusHistory(context.Background(), "cs-no-hist")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetStatusHistory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("cs-hist-err").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetStatusHistory, _ = db.Prepare("SELECT id, completed_service_id, previous_status, new_status, created_by, created_at FROM service_status_transitions WHERE completed_service_id = ?")

	results, err := repo.GetStatusHistory(context.Background(), "cs-hist-err")

	assert.Nil(t, results)
	assert.Error(t, err)
}
