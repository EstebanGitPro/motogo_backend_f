package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// queryGetExceptionsForUpdate retrieves all exceptions for a schedule with FOR UPDATE lock
// This prevents race conditions during concurrent exception creation
const queryGetExceptionsForUpdate = `
	SELECT id, schedule_id, entry_type, day_of_week, exception_start_date, exception_end_date,
		opening_time, closing_time, is_closed, active, created_at, updated_at
	FROM schedule_details
	WHERE schedule_id = ? AND entry_type = 'EXCEPTION'
	FOR UPDATE
`

// GetExceptionsByScheduleIDForUpdate retrieves all exceptions with a row-level lock.
// This must be called within a transaction to prevent race conditions during concurrent creation.
func (r *repository) GetExceptionsByScheduleIDForUpdate(
	ctx context.Context,
	tx output.Tx,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	// Get the underlying SQL transaction
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, 
			"schedule_id", scheduleID, 
			"error", "invalid transaction type")
		return nil, domain.ErrInternalServer
	}

	// Execute query within transaction with FOR UPDATE lock
	rows, err := sqlTx.QueryContext(ctx, queryGetExceptionsForUpdate, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, 
			"schedule_id", scheduleID, 
			"error", err)
		return nil, err
	}
	defer rows.Close()

	var exceptions []domain.ScheduleDetail
	for rows.Next() {
		var exception domain.ScheduleDetail
		var entryType string

		if err := rows.Scan(
			&exception.ID,
			&exception.ScheduleID,
			&entryType,
			&exception.DayOfWeek,
			&exception.ExceptionStartDate,
			&exception.ExceptionEndDate,
			&exception.OpeningTime,
			&exception.ClosingTime,
			&exception.IsClosed,
			&exception.Active,
			&exception.CreatedAt,
			&exception.UpdatedAt,
		); err != nil {
			log.Error(logger.LogScheduleDetailRepoScanError, "error", err)
			return nil, err
		}

		exception.EntryType = domain.EntryType(entryType)
		exceptions = append(exceptions, exception)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exceptions, nil
}
