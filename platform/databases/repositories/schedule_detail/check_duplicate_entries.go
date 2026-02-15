package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// nullUUID is a UUID that won't match any real record ID, used as default excludeID
const nullUUID = "00000000-0000-0000-0000-000000000000"

func (r *repository) CheckDayIsClosed(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = nullUUID
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

func (r *repository) CheckDayHasTimeSlots(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	excludeDetailID string,
) (bool, error) {
	if excludeDetailID == "" {
		excludeDetailID = nullUUID
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
