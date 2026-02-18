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

	return scanScheduleDetails(rows)
}
