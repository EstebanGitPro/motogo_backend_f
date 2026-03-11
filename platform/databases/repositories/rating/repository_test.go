package rating

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
// ToDomain Tests
// ============================================

func TestToDomain_WithAllFields(t *testing.T) {
	now := time.Now()
	comment := "Buen servicio"

	item := &RatingItem{
		ID:                 "item-001",
		CompletedServiceID: "cs-001",
		ServiceID:          "svc-001",
		Rating:             sql.NullInt32{Int32: 5, Valid: true},
		Comment:            sql.NullString{String: comment, Valid: true},
		RatedAt:            sql.NullTime{Time: now, Valid: true},
		IsOffensiveComment: false,
	}

	result := item.ToDomain()

	assert.Equal(t, "item-001", result.ID)
	assert.Equal(t, "cs-001", result.CompletedServiceID)
	assert.Equal(t, "svc-001", result.ServiceID)
	assert.NotNil(t, result.Rating)
	assert.Equal(t, 5, *result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, comment, *result.Comment)
	assert.NotNil(t, result.RatedAt)
	assert.Equal(t, now, *result.RatedAt)
	assert.False(t, result.IsOffensiveComment)
}

func TestToDomain_NullFields(t *testing.T) {
	item := &RatingItem{
		ID:                 "item-002",
		CompletedServiceID: "cs-002",
		ServiceID:          "svc-002",
		Rating:             sql.NullInt32{Valid: false},
		Comment:            sql.NullString{Valid: false},
		RatedAt:            sql.NullTime{Valid: false},
		IsOffensiveComment: false,
	}

	result := item.ToDomain()

	assert.Equal(t, "item-002", result.ID)
	assert.Nil(t, result.Rating)
	assert.Nil(t, result.Comment)
	assert.Nil(t, result.RatedAt)
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

	mock.ExpectPrepare("UPDATE completed_service_items")
	mock.ExpectPrepare("SELECT id, completed_service_id, service_id")
	mock.ExpectPrepare("SELECT.*csi.rating")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_RateServiceItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("UPDATE completed_service_items").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtRateServiceItem")
}

func TestNewRepository_PrepareError_GetItemByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("UPDATE completed_service_items")
	mock.ExpectPrepare("SELECT id, completed_service_id, service_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetItemByID")
}

func TestNewRepository_PrepareError_GetReviewsByServiceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("UPDATE completed_service_items")
	mock.ExpectPrepare("SELECT id, completed_service_id, service_id")
	mock.ExpectPrepare("SELECT.*csi.rating").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetReviewsByServiceID")
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
// RateServiceItem Tests
// ============================================

func TestRateServiceItem_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_service_items").
		WithArgs(5, (*string)(nil), false, "item-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.RateServiceItem(context.Background(), tx, "item-1", 5, nil, false)

	assert.NoError(t, err)
}

func TestRateServiceItem_InvalidTx(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.RateServiceItem(context.Background(), nil, "item-1", 5, nil, false)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestRateServiceItem_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE completed_service_items").
		WithArgs(5, (*string)(nil), false, "item-1").
		WillReturnError(sql.ErrConnDone)

	sqlTx, _ := db.Begin()
	tx := common.NewSQLTx(sqlTx)

	repo := &repository{db: db}

	err = repo.RateServiceItem(context.Background(), tx, "item-1", 5, nil, false)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

// ============================================
// GetItemByID Tests
// ============================================

func TestGetItemByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "completed_service_id", "service_id",
		"rating", "comment", "rated_at", "is_offensive_comment",
	}).AddRow("item-1", "cs-1", "svc-1", 4, "Buen servicio", now, false)

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("item-1").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetItemByID, _ = db.Prepare("SELECT id, completed_service_id, service_id, rating, comment, rated_at, is_offensive_comment FROM completed_service_items WHERE id = ?")

	result, err := repo.GetItemByID(context.Background(), "item-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "item-1", result.ID)
	assert.Equal(t, "cs-1", result.CompletedServiceID)
	assert.NotNil(t, result.Rating)
	assert.Equal(t, 4, *result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, "Buen servicio", *result.Comment)
}

func TestGetItemByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetItemByID, _ = db.Prepare("SELECT id, completed_service_id, service_id FROM completed_service_items WHERE id = ?")

	result, err := repo.GetItemByID(context.Background(), "not-found")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrServiceItemNotFound, err)
}

func TestGetItemByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, completed_service_id")
	stmt.ExpectQuery().
		WithArgs("error-id").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetItemByID, _ = db.Prepare("SELECT id, completed_service_id, service_id FROM completed_service_items WHERE id = ?")

	result, err := repo.GetItemByID(context.Background(), "error-id")

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// formatReviewerName Tests
// ============================================

func TestFormatReviewerName_Normal(t *testing.T) {
	result := formatReviewerName("Carlos", "Martínez")
	assert.Equal(t, "Carlos M.", result)
}

func TestFormatReviewerName_EmptyLastName(t *testing.T) {
	result := formatReviewerName("Carlos", "")
	assert.Equal(t, "Carlos", result)
}

func TestFormatReviewerName_WhitespaceLastName(t *testing.T) {
	result := formatReviewerName("  Carlos  ", "  ")
	assert.Equal(t, "Carlos", result)
}

// ============================================
// GetReviewsByServiceID Tests
// ============================================

func TestGetReviewsByServiceID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"rating", "comment", "rated_at",
		"first_name", "last_name",
		"motorcycle_model", "service_name",
	}).
		AddRow(5, "Excelente", now, "Carlos", "Martínez", "Honda CB 160F", "Cambio de aceite").
		AddRow(4, nil, now, "Ana", "López", nil, "Cambio de aceite")

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("svc-1", "branch-1").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetReviewsByServiceID, _ = db.Prepare("SELECT csi.rating, csi.comment, csi.rated_at, p.first_name, p.last_name, motorcycle_model, s.name AS service_name FROM completed_service_items csi WHERE csi.service_id = ? AND cs.branch_id = ?")

	result, err := repo.GetReviewsByServiceID(context.Background(), "svc-1", "branch-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "svc-1", result.ServiceID)
	assert.Equal(t, "Cambio de aceite", result.ServiceName)
	assert.Equal(t, 2, result.TotalReviews)
	assert.Equal(t, 4.5, result.AverageRating)
	assert.Equal(t, 1, result.Breakdown[5])
	assert.Equal(t, 1, result.Breakdown[4])
	assert.Len(t, result.Reviews, 2)
	assert.Equal(t, "Carlos M.", result.Reviews[0].ReviewerName)
	assert.NotNil(t, result.Reviews[0].Comment)
	assert.Equal(t, "Excelente", *result.Reviews[0].Comment)
	assert.NotNil(t, result.Reviews[0].MotorcycleModel)
	assert.Equal(t, "Honda CB 160F", *result.Reviews[0].MotorcycleModel)
	assert.Equal(t, "Ana L.", result.Reviews[1].ReviewerName)
	assert.Nil(t, result.Reviews[1].Comment)
	assert.Nil(t, result.Reviews[1].MotorcycleModel)
}

func TestGetReviewsByServiceID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"rating", "comment", "rated_at",
		"first_name", "last_name",
		"motorcycle_model", "service_name",
	})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("svc-empty", "branch-1").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetReviewsByServiceID, _ = db.Prepare("SELECT csi.rating FROM completed_service_items csi WHERE csi.service_id = ? AND cs.branch_id = ?")

	result, err := repo.GetReviewsByServiceID(context.Background(), "svc-empty", "branch-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalReviews)
	assert.Equal(t, 0.0, result.AverageRating)
	assert.Empty(t, result.Reviews)
}

func TestGetReviewsByServiceID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("svc-err", "branch-1").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetReviewsByServiceID, _ = db.Prepare("SELECT csi.rating FROM completed_service_items csi WHERE csi.service_id = ? AND cs.branch_id = ?")

	result, err := repo.GetReviewsByServiceID(context.Background(), "svc-err", "branch-1")

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetReviewsByServiceID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Return wrong number of columns to cause scan error
	rows := sqlmock.NewRows([]string{"rating", "comment"}).
		AddRow(5, "test")

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("svc-scan-err", "branch-1").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetReviewsByServiceID, _ = db.Prepare("SELECT csi.rating FROM completed_service_items csi WHERE csi.service_id = ? AND cs.branch_id = ?")

	result, err := repo.GetReviewsByServiceID(context.Background(), "svc-scan-err", "branch-1")

	assert.Nil(t, result)
	assert.Error(t, err)
}
