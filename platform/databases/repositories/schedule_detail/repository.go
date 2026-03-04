package schedule_detail

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	// Schedule Detail queries - maps to schedule_details table
	querySaveDetail = `
		INSERT INTO schedule_details (
			id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	queryGetDetailByID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE id = ?
	`
	queryGetDetailsByScheduleID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE schedule_id = ?
		ORDER BY day_of_week, opening_time
	`
	queryGetDetailsByScheduleAndDay = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE schedule_id = ? AND day_of_week = ? AND entry_type = 'REGULAR'
		ORDER BY opening_time
	`
	queryUpdateDetail = `
		UPDATE schedule_details
		SET day_of_week = ?, opening_time = ?, closing_time = ?, is_closed = ?, active = ?, updated_at = NOW()
		WHERE id = ?
	`
	queryDeleteDetail = `
		DELETE FROM schedule_details WHERE id = ?
	`
	// Time conflict check query - checks for overlapping time ranges
	queryCheckTimeConflict = `
		SELECT COUNT(*) FROM schedule_details
		WHERE schedule_id = ?
		  AND day_of_week = ?
		  AND entry_type = 'REGULAR'
		  AND is_closed = FALSE
		  AND id != ?
		  AND (
			(opening_time < ? AND closing_time > ?) OR
			(opening_time < ? AND closing_time > ?) OR
			(opening_time >= ? AND closing_time <= ?)
		  )
	`

	// ============================================
	// Schedule Exception queries (HU20-25)
	// ============================================

	// GetExceptionsByScheduleID retrieves all exceptions for a schedule ordered by start date
	queryGetExceptionsByScheduleID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE schedule_id = ? AND entry_type = 'EXCEPTION'
		ORDER BY exception_start_date ASC
	`

	// GetExceptionByID retrieves a specific exception ensuring it's an EXCEPTION type
	queryGetExceptionByID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE id = ? AND entry_type = 'EXCEPTION'
	`

	// CheckExceptionDateConflict checks if an exception already exists with overlapping dates
	// Overlap condition: existing.start <= new.end AND existing.end >= new.start
	queryCheckExceptionDateConflict = `
		SELECT COUNT(*) FROM schedule_details
		WHERE schedule_id = ? 
		  AND entry_type = 'EXCEPTION' 
		  AND id != ?
		  AND exception_start_date <= ?
		  AND exception_end_date >= ?
	`

	// ============================================
	// Duplicate Detection queries (Validation Rules)
	// ============================================

	// CheckDayIsClosed checks if a day already has a closed entry (R1, R2)
	queryCheckDayIsClosed = `
		SELECT COUNT(*) FROM schedule_details
		WHERE schedule_id = ? 
		  AND day_of_week = ? 
		  AND entry_type = 'REGULAR'
		  AND is_closed = TRUE
		  AND id != ?
	`

	// CheckDayHasTimeSlots checks if a day has time slot entries (not closed) (R3)
	queryCheckDayHasTimeSlots = `
		SELECT COUNT(*) FROM schedule_details
		WHERE schedule_id = ? 
		  AND day_of_week = ? 
		  AND entry_type = 'REGULAR'
		  AND is_closed = FALSE
		  AND id != ?
	`

	// CheckExceptionIsRedundant checks if exception is redundant because day is already closed in REGULAR (E1)
	queryCheckExceptionIsRedundant = `
		SELECT COUNT(*) FROM schedule_details
		WHERE schedule_id = ? 
		  AND day_of_week = ? 
		  AND entry_type = 'REGULAR'
		  AND is_closed = TRUE
		  AND active = TRUE
	`
	queryGetExceptionsForUpdate = `
	SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
		opening_time, closing_time, is_closed, active, created_at, updated_at
	FROM schedule_details
	WHERE schedule_id = ? AND entry_type = 'EXCEPTION'
	FOR UPDATE
`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                             *sql.DB
	stmtSaveDetail                 *sql.Stmt
	stmtGetDetailByID              *sql.Stmt
	stmtGetDetailsByScheduleID     *sql.Stmt
	stmtGetDetailsByScheduleDay    *sql.Stmt
	stmtUpdateDetail               *sql.Stmt
	stmtDeleteDetail               *sql.Stmt
	stmtCheckTimeConflict          *sql.Stmt
	stmtGetExceptionsByScheduleID  *sql.Stmt
	stmtGetExceptionByID           *sql.Stmt
	stmtCheckExceptionDateConflict *sql.Stmt
	stmtCheckDayIsClosed           *sql.Stmt
	stmtCheckDayHasTimeSlots       *sql.Stmt
	stmtCheckExceptionIsRedundant  *sql.Stmt
}

// NewRepository creates a new ScheduleDetailRepository with prepared statements (fail-fast pattern)
func NewRepository(db *sql.DB) (output.ScheduleDetailRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtSaveDetail, err := db.Prepare(querySaveDetail)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "SaveDetail", "error", err)
		return nil, fmt.Errorf("error preparing stmtSaveDetail: %w", err)
	}

	stmtGetDetailByID, err := db.Prepare(queryGetDetailByID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "GetDetailByID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetDetailByID: %w", err)
	}

	stmtGetDetailsByScheduleID, err := db.Prepare(queryGetDetailsByScheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "GetDetailsByScheduleID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetDetailsByScheduleID: %w", err)
	}

	stmtGetDetailsByScheduleDay, err := db.Prepare(queryGetDetailsByScheduleAndDay)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "GetDetailsByScheduleAndDay", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetDetailsByScheduleAndDay: %w", err)
	}

	stmtUpdateDetail, err := db.Prepare(queryUpdateDetail)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "UpdateDetail", "error", err)
		return nil, fmt.Errorf("error preparing stmtUpdateDetail: %w", err)
	}

	stmtDeleteDetail, err := db.Prepare(queryDeleteDetail)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "DeleteDetail", "error", err)
		return nil, fmt.Errorf("error preparing stmtDeleteDetail: %w", err)
	}

	stmtCheckTimeConflict, err := db.Prepare(queryCheckTimeConflict)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "CheckTimeConflict", "error", err)
		return nil, fmt.Errorf("error preparing stmtCheckTimeConflict: %w", err)
	}

	stmtGetExceptionsByScheduleID, err := db.Prepare(queryGetExceptionsByScheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "GetExceptionsByScheduleID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetExceptionsByScheduleID: %w", err)
	}

	stmtGetExceptionByID, err := db.Prepare(queryGetExceptionByID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "GetExceptionByID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetExceptionByID: %w", err)
	}

	stmtCheckExceptionDateConflict, err := db.Prepare(queryCheckExceptionDateConflict)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "CheckExceptionDateConflict", "error", err)
		return nil, fmt.Errorf("error preparing stmtCheckExceptionDateConflict: %w", err)
	}

	stmtCheckDayIsClosed, err := db.Prepare(queryCheckDayIsClosed)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "CheckDayIsClosed", "error", err)
		return nil, fmt.Errorf("error preparing stmtCheckDayIsClosed: %w", err)
	}

	stmtCheckDayHasTimeSlots, err := db.Prepare(queryCheckDayHasTimeSlots)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "CheckDayHasTimeSlots", "error", err)
		return nil, fmt.Errorf("error preparing stmtCheckDayHasTimeSlots: %w", err)
	}

	stmtCheckExceptionIsRedundant, err := db.Prepare(queryCheckExceptionIsRedundant)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoPrepareError, "statement", "CheckExceptionIsRedundant", "error", err)
		return nil, fmt.Errorf("error preparing stmtCheckExceptionIsRedundant: %w", err)
	}

	return &repository{
		db:                             db,
		stmtSaveDetail:                 stmtSaveDetail,
		stmtGetDetailByID:              stmtGetDetailByID,
		stmtGetDetailsByScheduleID:     stmtGetDetailsByScheduleID,
		stmtGetDetailsByScheduleDay:    stmtGetDetailsByScheduleDay,
		stmtUpdateDetail:               stmtUpdateDetail,
		stmtDeleteDetail:               stmtDeleteDetail,
		stmtCheckTimeConflict:          stmtCheckTimeConflict,
		stmtGetExceptionsByScheduleID:  stmtGetExceptionsByScheduleID,
		stmtGetExceptionByID:           stmtGetExceptionByID,
		stmtCheckExceptionDateConflict: stmtCheckExceptionDateConflict,
		stmtCheckDayIsClosed:           stmtCheckDayIsClosed,
		stmtCheckDayHasTimeSlots:       stmtCheckDayHasTimeSlots,
		stmtCheckExceptionIsRedundant:  stmtCheckExceptionIsRedundant,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
