package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

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

	return details, nil
}
