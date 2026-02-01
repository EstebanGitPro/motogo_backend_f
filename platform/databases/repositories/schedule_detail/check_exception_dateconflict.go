package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) CheckExceptionDateConflict(
	ctx context.Context,
	scheduleID string,
	excludeExceptionID string,
	startDate string,
	endDate string,
) (bool, error) {
	if excludeExceptionID == "" {
		excludeExceptionID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	log.Info("DEBUG_CheckExceptionDateConflict_QUERY_PARAMS",
		"scheduleID", scheduleID,
		"excludeExceptionID", excludeExceptionID,
		"startDate", startDate,
		"endDate", endDate)

	var count int
	err := r.stmtCheckExceptionDateConflict.QueryRowContext(ctx,
		scheduleID,
		excludeExceptionID,
		endDate,
		startDate,
	).Scan(&count)

	log.Info("DEBUG_CheckExceptionDateConflict_QUERY_RESULT",
		"count", count,
		"error", err)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "start_date", startDate, "end_date", endDate, "error", err)
		return false, err
	}

	return count > 0, nil
}
