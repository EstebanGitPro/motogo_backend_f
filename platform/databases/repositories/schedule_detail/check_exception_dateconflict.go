package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// CheckExceptionDateConflict checks if an exception already exists with overlapping dates (HU20)
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
	// Overlap logic: existing.start <= new.end AND existing.end >= new.start
	err := r.stmtCheckExceptionDateConflict.QueryRowContext(ctx,
		scheduleID,
		excludeExceptionID,
		endDate,   // existing.start <= new.end
		startDate, // existing.end >= new.start
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
