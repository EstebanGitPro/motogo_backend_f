package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetExceptionsByScheduleIDForUpdate(
	ctx context.Context,
	tx output.Tx,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError,
			"schedule_id", scheduleID,
			"error", "invalid transaction type")
		return nil, domain.ErrInternalServer
	}

	rows, err := sqlTx.QueryContext(ctx, queryGetExceptionsForUpdate, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError,
			"schedule_id", scheduleID,
			"error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

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
