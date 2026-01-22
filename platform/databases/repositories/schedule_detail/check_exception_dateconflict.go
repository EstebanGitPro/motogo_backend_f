package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// CheckExceptionDateConflict checks if an exception already exists for the given date (HU20)
func (r *repository) CheckExceptionDateConflict(
	ctx context.Context,
	scheduleID string,
	exceptionDate string,
	excludeExceptionID string,
) (bool, error) {
	if excludeExceptionID == "" {
		excludeExceptionID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	var count int
	err := r.stmtCheckExceptionDateConflict.QueryRowContext(ctx,
		scheduleID,
		exceptionDate,
		excludeExceptionID,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "date", exceptionDate, "error", err)
		return false, err
	}

	return count > 0, nil
}