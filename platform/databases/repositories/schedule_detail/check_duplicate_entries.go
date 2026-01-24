package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// CheckDayIsClosed checks if a day already has a closed entry (R1, R2)
func (r *repository) CheckDayIsClosed(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	var count int
	err := r.stmtCheckDayIsClosed.QueryRowContext(ctx,
		scheduleID,
		dayOfWeek,
		excludeDetailID,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "day_of_week", dayOfWeek, "error", err)
		return false, err
	}

	return count > 0, nil
}

// CheckDayHasTimeSlots checks if a day has time slot entries (not closed) (R3)
func (r *repository) CheckDayHasTimeSlots(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = "00000000-0000-0000-0000-000000000000" // UUID that won't match any real ID
	}

	var count int
	err := r.stmtCheckDayHasTimeSlots.QueryRowContext(ctx,
		scheduleID,
		dayOfWeek,
		excludeDetailID,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "day_of_week", dayOfWeek, "error", err)
		return false, err
	}

	return count > 0, nil
}

// CheckExceptionIsRedundant checks if exception is redundant because day is already closed in REGULAR (E1)
func (r *repository) CheckExceptionIsRedundant(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
) (bool, error) {
	var count int
	err := r.stmtCheckExceptionIsRedundant.QueryRowContext(ctx,
		scheduleID,
		dayOfWeek,
	).Scan(&count)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoConflictCheck, "schedule_id", scheduleID, "day_of_week", dayOfWeek, "error", err)
		return false, err
	}

	return count > 0, nil
}
