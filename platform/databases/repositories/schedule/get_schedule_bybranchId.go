package schedule

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	var schedule domain.BranchSchedule

	err := r.stmtGetScheduleByBranchID.QueryRowContext(ctx, branchID).Scan(
		&schedule.ID,
		&schedule.BranchID,
		&schedule.Active,
		&schedule.StartDate,
		&schedule.EndDate,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(logger.LogScheduleRepoGetByBranchError, "branch_id", branchID, "error", err)
		return nil, err
	}

	return &schedule, nil
}
