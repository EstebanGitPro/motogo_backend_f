package schedule_detail

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

func TestNewRepository_PrepareError_GetDetailByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO schedule_details")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetDetailByID")
}

func TestNewRepository_PrepareError_GetDetailsByScheduleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO schedule_details")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE schedule_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetDetailsByScheduleID")
}

func TestNewRepository_PrepareError_UpdateDetail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO schedule_details")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE schedule_id")
	mock.ExpectPrepare("SELECT.*day_of_week.*AND entry_type = 'REGULAR'")
	mock.ExpectPrepare("UPDATE schedule_details").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdateDetail")
}

func TestNewRepository_PrepareError_DeleteDetail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO schedule_details")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM schedule_details.*WHERE schedule_id")
	mock.ExpectPrepare("SELECT.*day_of_week.*AND entry_type = 'REGULAR'")
	mock.ExpectPrepare("UPDATE schedule_details")
	mock.ExpectPrepare("DELETE FROM schedule_details").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeleteDetail")
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

// ============================================
// GetExceptionByID Tests
// ============================================

func TestGetExceptionByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "schedule_id", "entry_type", "day_of_week",
		"exception_start_date", "exception_end_date",
		"opening_time", "closing_time", "is_closed", "active",
		"created_at", "updated_at",
	}).AddRow("exc-001", "schedule-123", "EXCEPTION", nil, now, now.Add(24*time.Hour), nil, nil, true, true, now, now)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("exc-001").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetExceptionByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	exception, err := repo.GetExceptionByID(context.Background(), "exc-001")

	assert.NoError(t, err)
	assert.NotNil(t, exception)
	assert.Equal(t, "exc-001", exception.ID)
	assert.Equal(t, domain.EntryType("EXCEPTION"), exception.EntryType)
}

func TestGetExceptionByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetExceptionByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	exception, err := repo.GetExceptionByID(context.Background(), "not-found")

	assert.NoError(t, err)
	assert.Nil(t, exception)
}

func TestGetExceptionByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("exc-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetExceptionByID, _ = db.Prepare("SELECT * FROM schedule_details WHERE id = ?")

	exception, err := repo.GetExceptionByID(context.Background(), "exc-error")

	assert.Nil(t, exception)
	assert.Error(t, err)
}

// ============================================
// GetExceptionsByScheduleID Tests
// ============================================

func TestGetExceptionsByScheduleID_Success(t *testing.T) {
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
		AddRow("exc-001", "schedule-123", "EXCEPTION", nil, now, now.Add(24*time.Hour), nil, nil, true, true, now, now).
		AddRow("exc-002", "schedule-123", "EXCEPTION", nil, now.Add(48*time.Hour), now.Add(72*time.Hour), nil, nil, true, true, now, now)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetExceptionsByScheduleID, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ?")

	exceptions, err := repo.GetExceptionsByScheduleID(context.Background(), "schedule-123")

	assert.NoError(t, err)
	assert.Len(t, exceptions, 2)
	assert.Equal(t, "exc-001", exceptions[0].ID)
}

func TestGetExceptionsByScheduleID_Empty(t *testing.T) {
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
		WithArgs("schedule-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetExceptionsByScheduleID, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ?")

	exceptions, err := repo.GetExceptionsByScheduleID(context.Background(), "schedule-empty")

	assert.NoError(t, err)
	assert.Empty(t, exceptions)
}

func TestGetExceptionsByScheduleID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetExceptionsByScheduleID, _ = db.Prepare("SELECT * FROM schedule_details WHERE schedule_id = ?")

	exceptions, err := repo.GetExceptionsByScheduleID(context.Background(), "schedule-error")

	assert.Nil(t, exceptions)
	assert.Error(t, err)
}

// ============================================
// CheckDayIsClosed Tests
// ============================================

func TestCheckDayIsClosed_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "00000000-0000-0000-0000-000000000000").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckDayIsClosed, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isClosed, err := repo.CheckDayIsClosed(context.Background(), "schedule-123", 1, "")

	assert.NoError(t, err)
	assert.True(t, isClosed)
}

func TestCheckDayIsClosed_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "exclude-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckDayIsClosed, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isClosed, err := repo.CheckDayIsClosed(context.Background(), "schedule-123", 1, "exclude-123")

	assert.NoError(t, err)
	assert.False(t, isClosed)
}

func TestCheckDayIsClosed_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", 1, "00000000-0000-0000-0000-000000000000").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckDayIsClosed, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isClosed, err := repo.CheckDayIsClosed(context.Background(), "schedule-error", 1, "")

	assert.False(t, isClosed)
	assert.Error(t, err)
}

// ============================================
// CheckDayHasTimeSlots Tests
// ============================================

func TestCheckDayHasTimeSlots_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "00000000-0000-0000-0000-000000000000").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckDayHasTimeSlots, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasSlots, err := repo.CheckDayHasTimeSlots(context.Background(), "schedule-123", 1, "")

	assert.NoError(t, err)
	assert.True(t, hasSlots)
}

func TestCheckDayHasTimeSlots_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1, "exclude-456").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckDayHasTimeSlots, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasSlots, err := repo.CheckDayHasTimeSlots(context.Background(), "schedule-123", 1, "exclude-456")

	assert.NoError(t, err)
	assert.False(t, hasSlots)
}

func TestCheckDayHasTimeSlots_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", 2, "00000000-0000-0000-0000-000000000000").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckDayHasTimeSlots, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasSlots, err := repo.CheckDayHasTimeSlots(context.Background(), "schedule-error", 2, "")

	assert.False(t, hasSlots)
	assert.Error(t, err)
}

// ============================================
// CheckExceptionIsRedundant Tests
// ============================================

func TestCheckExceptionIsRedundant_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1).
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckExceptionIsRedundant, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isRedundant, err := repo.CheckExceptionIsRedundant(context.Background(), "schedule-123", 1)

	assert.NoError(t, err)
	assert.True(t, isRedundant)
}

func TestCheckExceptionIsRedundant_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", 1).
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckExceptionIsRedundant, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isRedundant, err := repo.CheckExceptionIsRedundant(context.Background(), "schedule-123", 1)

	assert.NoError(t, err)
	assert.False(t, isRedundant)
}

func TestCheckExceptionIsRedundant_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", 3).
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckExceptionIsRedundant, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	isRedundant, err := repo.CheckExceptionIsRedundant(context.Background(), "schedule-error", 3)

	assert.False(t, isRedundant)
	assert.Error(t, err)
}

// ============================================
// CheckExceptionDateConflict Tests
// ============================================

func TestCheckExceptionDateConflict_HasConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", "00000000-0000-0000-0000-000000000000", "2024-12-31", "2024-01-01").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckExceptionDateConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckExceptionDateConflict(context.Background(), "schedule-123", "", "2024-01-01", "2024-12-31")

	assert.NoError(t, err)
	assert.True(t, hasConflict)
}

func TestCheckExceptionDateConflict_NoConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-123", "exclude-789", "2024-12-31", "2024-01-01").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckExceptionDateConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckExceptionDateConflict(context.Background(), "schedule-123", "exclude-789", "2024-01-01", "2024-12-31")

	assert.NoError(t, err)
	assert.False(t, hasConflict)
}

func TestCheckExceptionDateConflict_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("schedule-error", "00000000-0000-0000-0000-000000000000", "2024-12-31", "2024-01-01").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckExceptionDateConflict, _ = db.Prepare("SELECT COUNT(*) FROM schedule_details WHERE schedule_id = ?")

	hasConflict, err := repo.CheckExceptionDateConflict(context.Background(), "schedule-error", "", "2024-01-01", "2024-12-31")

	assert.False(t, hasConflict)
	assert.Error(t, err)
}

// ============================================
// SaveScheduleDetail Tests
// ============================================

func TestSaveScheduleDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	monday := 1
	openingTime := "09:00:00"
	closingTime := "18:00:00"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO schedule_details").
		WithArgs("detail-123", "schedule-123", domain.EntryTypeRegular, &monday, nil, nil, &openingTime, &closingTime, false, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	detail := domain.ScheduleDetail{
		ID:          "detail-123",
		ScheduleID:  "schedule-123",
		EntryType:   domain.EntryTypeRegular,
		DayOfWeek:   &monday,
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
		IsClosed:    false,
		Active:      true,
	}

	err = repo.SaveScheduleDetail(context.Background(), sqlTx, detail)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveScheduleDetail_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO schedule_details").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveScheduleDetail(context.Background(), sqlTx, domain.ScheduleDetail{ID: "detail-err"})

	assert.Error(t, err)
}

func TestSaveScheduleDetail_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.SaveScheduleDetail(context.Background(), nil, domain.ScheduleDetail{})

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdateScheduleDetail Tests
// ============================================

func TestUpdateScheduleDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	monday := 1
	openingTime := "10:00:00"
	closingTime := "19:00:00"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE schedule_details").
		WithArgs(&monday, &openingTime, &closingTime, false, true, "detail-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	detail := domain.ScheduleDetail{
		ID:          "detail-123",
		DayOfWeek:   &monday,
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
		IsClosed:    false,
		Active:      true,
	}

	err = repo.UpdateScheduleDetail(context.Background(), sqlTx, detail)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScheduleDetail_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE schedule_details").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.UpdateScheduleDetail(context.Background(), sqlTx, domain.ScheduleDetail{ID: "detail-err"})

	assert.Error(t, err)
}

func TestUpdateScheduleDetail_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.UpdateScheduleDetail(context.Background(), nil, domain.ScheduleDetail{})

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeleteScheduleDetail Tests
// ============================================

func TestDeleteScheduleDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM schedule_details").
		WithArgs("detail-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteScheduleDetail(context.Background(), sqlTx, "detail-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteScheduleDetail_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM schedule_details").
		WithArgs("detail-err").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteScheduleDetail(context.Background(), sqlTx, "detail-err")

	assert.Error(t, err)
}

func TestDeleteScheduleDetail_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.DeleteScheduleDetail(context.Background(), nil, "detail-123")

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// GetExceptionsByScheduleIDForUpdate Tests
// ============================================

func TestGetExceptionsByScheduleIDForUpdate_Success(t *testing.T) {
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
		AddRow("exc-001", "schedule-123", "EXCEPTION", nil, now, now.Add(24*time.Hour), nil, nil, true, true, now, now).
		AddRow("exc-002", "schedule-123", "EXCEPTION", nil, now.Add(48*time.Hour), now.Add(72*time.Hour), nil, nil, true, true, now, now)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, schedule_id").
		WithArgs("schedule-123").
		WillReturnRows(rows)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	exceptions, err := repo.GetExceptionsByScheduleIDForUpdate(context.Background(), sqlTx, "schedule-123")

	assert.NoError(t, err)
	assert.Len(t, exceptions, 2)
	assert.Equal(t, "exc-001", exceptions[0].ID)
	assert.Equal(t, domain.EntryType("EXCEPTION"), exceptions[0].EntryType)
}

func TestGetExceptionsByScheduleIDForUpdate_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	exceptions, err := repo.GetExceptionsByScheduleIDForUpdate(context.Background(), nil, "schedule-123")

	assert.Nil(t, exceptions)
	assert.Equal(t, domain.ErrInternalServer, err)
}

func TestGetExceptionsByScheduleIDForUpdate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, schedule_id").
		WithArgs("schedule-error").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	exceptions, err := repo.GetExceptionsByScheduleIDForUpdate(context.Background(), sqlTx, "schedule-error")

	assert.Nil(t, exceptions)
	assert.Error(t, err)
}
