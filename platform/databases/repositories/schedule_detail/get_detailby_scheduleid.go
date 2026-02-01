package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	rows, err := r.stmtGetDetailsByScheduleID.QueryContext(ctx, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError, "schedule_id", scheduleID, "error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var details []domain.ScheduleDetail
	for rows.Next() {
		var detail domain.ScheduleDetail
		var entryType string

		if err := rows.Scan(
			&detail.ID,
			&detail.ScheduleID,
			&entryType,
			&detail.DayOfWeek,
			&detail.ExceptionStartDate,
			&detail.ExceptionEndDate,
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
