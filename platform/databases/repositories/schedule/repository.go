package schedule

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	querySaveSchedule = `
		INSERT INTO branch_schedules (id, branch_id, active, start_date, end_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	queryGetScheduleByBranchID = `
		SELECT id, branch_id, active, start_date, end_date, created_at, updated_at
		FROM branch_schedules
		WHERE branch_id = ?
	`
	queryGetScheduleByID = `
		SELECT id, branch_id, active, start_date, end_date, created_at, updated_at
		FROM branch_schedules
		WHERE id = ?
	`
	queryUpdateSchedule = `
		UPDATE branch_schedules
		SET active = ?, start_date = ?, end_date = ?, updated_at = NOW()
		WHERE id = ?
	`
	queryDeleteSchedule = `
		DELETE FROM branch_schedules WHERE id = ?
	`
	querySetActive = `
		UPDATE branch_schedules
		SET active = ?, updated_at = NOW()
		WHERE id = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                        *sql.DB
	stmtSaveSchedule          *sql.Stmt
	stmtGetScheduleByBranchID *sql.Stmt
	stmtGetScheduleByID       *sql.Stmt
	stmtUpdateSchedule        *sql.Stmt
	stmtDeleteSchedule        *sql.Stmt
	stmtSetActive             *sql.Stmt
}

func NewRepository(db *sql.DB) (output.ScheduleRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtSaveSchedule, err := db.Prepare(querySaveSchedule)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "SaveSchedule", "error", err)
		return nil, fmt.Errorf("error preparing stmtSaveSchedule: %w", err)
	}

	stmtGetScheduleByBranchID, err := db.Prepare(queryGetScheduleByBranchID)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "GetScheduleByBranchID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetScheduleByBranchID: %w", err)
	}

	stmtGetScheduleByID, err := db.Prepare(queryGetScheduleByID)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "GetScheduleByID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetScheduleByID: %w", err)
	}

	stmtUpdateSchedule, err := db.Prepare(queryUpdateSchedule)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "UpdateSchedule", "error", err)
		return nil, fmt.Errorf("error preparing stmtUpdateSchedule: %w", err)
	}

	stmtDeleteSchedule, err := db.Prepare(queryDeleteSchedule)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "DeleteSchedule", "error", err)
		return nil, fmt.Errorf("error preparing stmtDeleteSchedule: %w", err)
	}

	stmtSetActive, err := db.Prepare(querySetActive)
	if err != nil {
		log.Error(logger.LogScheduleRepoPrepareError, "statement", "SetActive", "error", err)
		return nil, fmt.Errorf("error preparing stmtSetActive: %w", err)
	}

	return &repository{
		db:                        db,
		stmtSaveSchedule:          stmtSaveSchedule,
		stmtGetScheduleByBranchID: stmtGetScheduleByBranchID,
		stmtGetScheduleByID:       stmtGetScheduleByID,
		stmtUpdateSchedule:        stmtUpdateSchedule,
		stmtDeleteSchedule:        stmtDeleteSchedule,
		stmtSetActive:             stmtSetActive,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
