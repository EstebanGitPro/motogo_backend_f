package schedule

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
	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE id")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active")
	mock.ExpectPrepare("DELETE FROM branch_schedules")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSaveSchedule")
}

func TestNewRepository_PrepareError_GetByBranchID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetScheduleByBranchID")
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetScheduleByID")
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE id")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdateSchedule")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE id")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active")
	mock.ExpectPrepare("DELETE FROM branch_schedules").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeleteSchedule")
}

func TestNewRepository_PrepareError_SetActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branch_schedules")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE branch_id")
	mock.ExpectPrepare("SELECT id, branch_id, active.*WHERE id")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active")
	mock.ExpectPrepare("DELETE FROM branch_schedules")
	mock.ExpectPrepare("UPDATE branch_schedules.*SET active").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSetActive")
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
// GetScheduleByID Tests
// ============================================

func TestGetScheduleByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "branch_id", "active", "start_date", "end_date", "created_at", "updated_at"}).
		AddRow("schedule-123", "branch-456", true, now, now.Add(24*time.Hour), now, now)

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("schedule-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetScheduleByID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE id = ?")

	schedule, err := repo.GetScheduleByID(context.Background(), "schedule-123")

	assert.NoError(t, err)
	assert.NotNil(t, schedule)
	assert.Equal(t, "schedule-123", schedule.ID)
	assert.Equal(t, "branch-456", schedule.BranchID)
	assert.True(t, schedule.Active)
}

func TestGetScheduleByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetScheduleByID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE id = ?")

	schedule, err := repo.GetScheduleByID(context.Background(), "not-found")

	assert.Nil(t, schedule)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

func TestGetScheduleByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("schedule-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetScheduleByID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE id = ?")

	schedule, err := repo.GetScheduleByID(context.Background(), "schedule-error")

	assert.Nil(t, schedule)
	assert.Error(t, err)
}

// ============================================
// GetScheduleByBranchID Tests
// ============================================

func TestGetScheduleByBranchID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "branch_id", "active", "start_date", "end_date", "created_at", "updated_at"}).
		AddRow("schedule-789", "branch-123", true, now, now.Add(24*time.Hour), now, now)

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("branch-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetScheduleByBranchID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE branch_id = ?")

	schedule, err := repo.GetScheduleByBranchID(context.Background(), "branch-123")

	assert.NoError(t, err)
	assert.NotNil(t, schedule)
	assert.Equal(t, "schedule-789", schedule.ID)
	assert.Equal(t, "branch-123", schedule.BranchID)
}

func TestGetScheduleByBranchID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("no-schedule").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetScheduleByBranchID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE branch_id = ?")

	schedule, err := repo.GetScheduleByBranchID(context.Background(), "no-schedule")

	// GetScheduleByBranchID returns nil, nil for not found
	assert.Nil(t, schedule)
	assert.Nil(t, err)
}

func TestGetScheduleByBranchID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, branch_id, active")
	stmt.ExpectQuery().
		WithArgs("branch-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetScheduleByBranchID, _ = db.Prepare("SELECT id, branch_id, active FROM branch_schedules WHERE branch_id = ?")

	schedule, err := repo.GetScheduleByBranchID(context.Background(), "branch-error")

	assert.Nil(t, schedule)
	assert.Error(t, err)
}

// ============================================
// SaveSchedule Tests
// ============================================

func TestSaveSchedule_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	startDate := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_schedules").
		WithArgs("schedule-123", "branch-123", true, startDate, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	schedule := domain.BranchSchedule{
		ID:        "schedule-123",
		BranchID:  "branch-123",
		Active:    true,
		StartDate: startDate,
	}

	err = repo.SaveSchedule(context.Background(), sqlTx, schedule)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveSchedule_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_schedules").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveSchedule(context.Background(), sqlTx, domain.BranchSchedule{})

	assert.Error(t, err)
}

func TestSaveSchedule_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.SaveSchedule(context.Background(), nil, domain.BranchSchedule{})

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdateSchedule Tests
// ============================================

func TestUpdateSchedule_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	startDate := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branch_schedules").
		WithArgs(true, startDate, nil, "schedule-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	schedule := domain.BranchSchedule{
		ID:        "schedule-123",
		Active:    true,
		StartDate: startDate,
	}

	err = repo.UpdateSchedule(context.Background(), sqlTx, schedule)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSchedule_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branch_schedules").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.UpdateSchedule(context.Background(), sqlTx, domain.BranchSchedule{})

	assert.Error(t, err)
}

func TestUpdateSchedule_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.UpdateSchedule(context.Background(), nil, domain.BranchSchedule{})

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeleteSchedule Tests
// ============================================

func TestDeleteSchedule_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_schedules").
		WithArgs("schedule-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteSchedule(context.Background(), sqlTx, "schedule-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSchedule_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_schedules").
		WithArgs("schedule-err").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteSchedule(context.Background(), sqlTx, "schedule-err")

	assert.Error(t, err)
}

func TestDeleteSchedule_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.DeleteSchedule(context.Background(), nil, "schedule-123")

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// SetActive Tests
// ============================================

func TestSetActive_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branch_schedules").
		WithArgs(true, "schedule-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SetActive(context.Background(), sqlTx, "schedule-123", true)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetActive_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branch_schedules").
		WithArgs(false, "schedule-err").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SetActive(context.Background(), sqlTx, "schedule-err", false)

	assert.Error(t, err)
}

func TestSetActive_InvalidTransaction(t *testing.T) {
	repo := &repository{}

	err := repo.SetActive(context.Background(), nil, "schedule-123", true)

	assert.Equal(t, domain.ErrInvalidTransaction, err)
}
