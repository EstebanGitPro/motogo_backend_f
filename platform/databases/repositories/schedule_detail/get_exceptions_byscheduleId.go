package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetExceptionsByScheduleID retrieves all exceptions for a schedule (HU23)
func (r *repository) GetExceptionsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	rows, err := r.stmtGetExceptionsByScheduleID.QueryContext(ctx, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, "schedule_id", scheduleID, "error", err)
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
