package schedule_detail

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	// Schedule Detail queries - maps to schedule_details table
	querySaveDetail = `
		INSERT INTO schedule_details (
			id, schedule_id, entry_type, day_of_week, exception_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	queryGetDetailByID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE id = ?
	`
	queryGetDetailsByScheduleID = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_date,
			opening_time, closing_time, is_closed, active, created_at, updated_at
		FROM schedule_details
		WHERE schedule_id = ?
		ORDER BY day_of_week, opening_time
	`
	queryGetDetailsByScheduleAndDay = `
		SELECT id, schedule_id, entry_type, day_of_week, exception_date,
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
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                          *sql.DB
	stmtSaveDetail              *sql.Stmt
	stmtGetDetailByID           *sql.Stmt
	stmtGetDetailsByScheduleID  *sql.Stmt
	stmtGetDetailsByScheduleDay *sql.Stmt
	stmtUpdateDetail            *sql.Stmt
	stmtDeleteDetail            *sql.Stmt
	stmtCheckTimeConflict       *sql.Stmt
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

	return &repository{
		db:                          db,
		stmtSaveDetail:              stmtSaveDetail,
		stmtGetDetailByID:           stmtGetDetailByID,
		stmtGetDetailsByScheduleID:  stmtGetDetailsByScheduleID,
		stmtGetDetailsByScheduleDay: stmtGetDetailsByScheduleDay,
		stmtUpdateDetail:            stmtUpdateDetail,
		stmtDeleteDetail:            stmtDeleteDetail,
		stmtCheckTimeConflict:       stmtCheckTimeConflict,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// SaveScheduleDetail saves a new schedule detail to the database (HU6)
func (r *repository) SaveScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtSaveDetail).ExecContext(ctx,
		detail.ID,
		detail.ScheduleID,
		detail.EntryType,
		detail.DayOfWeek,
		detail.ExceptionDate,
		detail.OpeningTime,
		detail.ClosingTime,
		detail.IsClosed,
		detail.Active,
	)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoSaveError, "detail_id", detail.ID, "error", err)
		return err
	}

	return nil
}

// GetDetailByID retrieves a schedule detail by its ID
func (r *repository) GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	var detail domain.ScheduleDetail
	var entryType string

	err := r.stmtGetDetailByID.QueryRowContext(ctx, detailID).Scan(
		&detail.ID,
		&detail.ScheduleID,
		&entryType,
		&detail.DayOfWeek,
		&detail.ExceptionDate,
		&detail.OpeningTime,
		&detail.ClosingTime,
		&detail.IsClosed,
		&detail.Active,
		&detail.CreatedAt,
		&detail.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(logger.LogScheduleDetailRepoGetByIDError, "detail_id", detailID, "error", err)
		return nil, err
	}

	detail.EntryType = domain.EntryType(entryType)
	return &detail, nil
}

// GetDetailsByScheduleID retrieves all schedule details for a schedule (HU9)
func (r *repository) GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	rows, err := r.stmtGetDetailsByScheduleID.QueryContext(ctx, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, "schedule_id", scheduleID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []domain.ScheduleDetail
	for rows.Next() {
		var detail domain.ScheduleDetail
		var entryType string

		if err := rows.Scan(
			&detail.ID,
			&detail.ScheduleID,
			&entryType,
			&detail.DayOfWeek,
			&detail.ExceptionDate,
			&detail.OpeningTime,
			&detail.ClosingTime,
			&detail.IsClosed,
			&detail.Active,
			&detail.CreatedAt,
			&detail.UpdatedAt,
		); err != nil {
			log.Error(logger.LogScheduleDetailRepoScanError, "error", err)
			return nil, err
		}

		detail.EntryType = domain.EntryType(entryType)
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}

// GetDetailsByScheduleAndDay retrieves schedule details for a specific schedule and day
func (r *repository) GetDetailsByScheduleAndDay(ctx context.Context, scheduleID string, dayOfWeek int) ([]domain.ScheduleDetail, error) {
	rows, err := r.stmtGetDetailsByScheduleDay.QueryContext(ctx, scheduleID, dayOfWeek)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, "schedule_id", scheduleID, "day", dayOfWeek, "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []domain.ScheduleDetail
	for rows.Next() {
		var detail domain.ScheduleDetail
		var entryType string

		if err := rows.Scan(
			&detail.ID,
			&detail.ScheduleID,
			&entryType,
			&detail.DayOfWeek,
			&detail.ExceptionDate,
			&detail.OpeningTime,
			&detail.ClosingTime,
			&detail.IsClosed,
			&detail.Active,
			&detail.CreatedAt,
			&detail.UpdatedAt,
		); err != nil {
			log.Error(logger.LogScheduleDetailRepoScanError, "error", err)
			return nil, err
		}

		detail.EntryType = domain.EntryType(entryType)
		details = append(details, detail)
	}

	return details, nil
}

// UpdateScheduleDetail updates an existing schedule detail (HU7)
func (r *repository) UpdateScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtUpdateDetail).ExecContext(ctx,
		detail.DayOfWeek,
		detail.OpeningTime,
		detail.ClosingTime,
		detail.IsClosed,
		detail.Active,
		detail.ID,
	)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoUpdateError, "detail_id", detail.ID, "error", err)
		return err
	}

	return nil
}

// DeleteScheduleDetail deletes a schedule detail (HU8)
func (r *repository) DeleteScheduleDetail(ctx context.Context, tx output.Tx, detailID string) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtDeleteDetail).ExecContext(ctx, detailID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoDeleteError, "detail_id", detailID, "error", err)
		return err
	}

	return nil
}

// CheckTimeConflict checks if a proposed time slot conflicts with existing slots
func (r *repository) CheckTimeConflict(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	openingTime, closingTime string,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	var count int
	err := r.stmtCheckTimeConflict.QueryRowContext(ctx,
		scheduleID,
		dayOfWeek,
		excludeDetailID,
		closingTime, openingTime, // Check if new closing is after existing opening AND new opening is before existing closing
		closingTime, openingTime,
		openingTime, closingTime,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "error", err)
		return false, err
	}

	return count > 0, nil
}
