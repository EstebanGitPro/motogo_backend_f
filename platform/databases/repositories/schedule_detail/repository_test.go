package schedule_detail

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// First prepare should fail
	mock.ExpectPrepare("INSERT INTO schedule_details").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSaveDetail")
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
// GetDetailByID Tests
// ============================================

func TestGetDetailByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	openTime := "08:00:00"
	closeTime := "18:00:00"

	rows := sqlmock.NewRows([]string{
		"id", "schedule_id", "entry_type", "day_of_week",
		"exception_start_date", "exception_end_date",
		"opening_time", "closing_time", "is_closed", "active",
		"created_at", "updated_at",
	}).AddRow(
		"detail-123", "schedule-456", "REGULAR", 1,
		nil, nil,
		openTime, closeTime, false, true,
		now, now,
	)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("detail-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDetailByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	detail, err := repo.GetDetailByID(context.Background(), "detail-123")

	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "detail-123", detail.ID)
	assert.Equal(t, "schedule-456", detail.ScheduleID)
	assert.NotNil(t, detail.DayOfWeek)
	assert.Equal(t, 1, *detail.DayOfWeek)
	assert.False(t, detail.IsClosed)
}

func TestGetDetailByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetDetailByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	detail, err := repo.GetDetailByID(context.Background(), "not-found")

	// GetDetailByID returns nil, nil for not found
	assert.Nil(t, detail)
	assert.Nil(t, err)
}

func TestGetDetailByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("detail-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDetailByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	detail, err := repo.GetDetailByID(context.Background(), "detail-error")

	assert.Nil(t, detail)
	assert.Error(t, err)
}

// ============================================
// GetDetailsByScheduleAndDay Tests
// ============================================

func TestGetDetailsByScheduleAndDay_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "schedule_id", "entry_type", "day_of_week",
		"exception_start_date", "exception_end_date",
		"opening_time", "closing_time", "is_closed", "active",
		"created_at", "updated_at",
	}).AddRow(
		"detail-001", "schedule-123", "REGULAR", 1,
		nil, nil,
		"08:00:00", "12:00:00", false, true,
		now, now,
	).AddRow(
		"detail-002", "schedule-123", "REGULAR", 1,
		nil, nil,
		"14:00:00", "18:00:00", false, true,
		now, now,
	)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1).
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDetailsByScheduleDay, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ? AND day_of_week = ?")

	details, err := repo.GetDetailsByScheduleAndDay(context.Background(), "schedule-123", 1)

	assert.NoError(t, err)
	assert.Len(t, details, 2)
	assert.Equal(t, "detail-001", details[0].ID)
	assert.Equal(t, "detail-002", details[1].ID)
}

func TestGetDetailsByScheduleAndDay_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "schedule_id", "entry_type", "day_of_week",
		"exception_start_date", "exception_end_date",
		"opening_time", "closing_time", "is_closed", "active",
		"created_at", "updated_at",
	})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-empty", 7).
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDetailsByScheduleDay, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ? AND day_of_week = ?")

	details, err := repo.GetDetailsByScheduleAndDay(context.Background(), "schedule-empty", 7)

	assert.NoError(t, err)
	assert.Empty(t, details)
}

func TestGetDetailsByScheduleAndDay_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", 1).
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDetailsByScheduleDay, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ? AND day_of_week = ?")

	details, err := repo.GetDetailsByScheduleAndDay(context.Background(), "schedule-error", 1)

	assert.Nil(t, details)
	assert.Error(t, err)
}

// ============================================
// CheckTimeConflict Tests
// ============================================

func TestCheckTimeConflict_NoConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "00000000-0000-0000-0000-000000000000", "18:00", "08:00", "18:00", "08:00", "08:00", "18:00").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckTimeConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckTimeConflict(context.Background(), "schedule-123", 1, "08:00", "18:00", "")

	assert.NoError(t, err)
	assert.False(t, hasConflict)
}

func TestCheckTimeConflict_HasConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "detail-exclude", "18:00", "08:00", "18:00", "08:00", "08:00", "18:00").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckTimeConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckTimeConflict(context.Background(), "schedule-123", 1, "08:00", "18:00", "detail-exclude")

	assert.NoError(t, err)
	assert.True(t, hasConflict)
}

func TestCheckTimeConflict_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", 1, "00000000-0000-0000-0000-000000000000", "18:00", "08:00", "18:00", "08:00", "08:00", "18:00").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckTimeConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckTimeConflict(context.Background(), "schedule-error", 1, "08:00", "18:00", "")

	assert.False(t, hasConflict)
	assert.Error(t, err)
}

// ============================================
// GetDetailsByScheduleID Tests
// ============================================

func TestGetDetailsByScheduleID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "schedule_id", "entry_type", "day_of_week",
		"exception_start_date", "exception_end_date",
		"opening_time", "closing_time", "is_closed", "active",
		"created_at", "updated_at",
	}).
		AddRow("detail-001", "schedule-123", "regular", 1, nil, nil, "08:00", "18:00", false, true, now, now).
		AddRow("detail-002", "schedule-123", "regular", 2, nil, nil, "09:00", "17:00", false, true, now, now)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDetailsByScheduleID, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ?")

	details, err := repo.GetDetailsByScheduleID(context.Background(), "schedule-123")

	assert.NoError(t, err)
	assert.Len(t, details, 2)
	assert.Equal(t, "detail-001", details[0].ID)
}

func TestGetDetailsByScheduleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDetailsByScheduleID, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ?")

	details, err := repo.GetDetailsByScheduleID(context.Background(), "schedule-error")

	assert.Nil(t, details)
	assert.Error(t, err)
}
