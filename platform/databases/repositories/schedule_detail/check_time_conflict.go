package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) CheckTimeConflict(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	openingTime, closingTime string,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	var count int
	err := r.stmtCheckTimeConflict.QueryRowContext(ctx,
		scheduleID,
		dayOfWeek,
		excludeDetailID,
		closingTime, openingTime, // Check if new closing is after existing opening AND new opening is before existing closing
		closingTime, openingTime,
		openingTime, closingTime,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "error", err)
		return false, err
	}

	return count > 0, nil
}
