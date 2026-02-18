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
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanScheduleDetails(rows)
}
