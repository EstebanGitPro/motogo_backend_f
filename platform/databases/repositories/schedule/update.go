package schedule

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateSchedule,
		schedule.Active,
		schedule.StartDate,
		schedule.EndDate, // nullable
		schedule.ID,
	)
	if err != nil {
		log.Error(logger.LogScheduleRepoUpdateError, "schedule_id", schedule.ID, "error", err)
		return err
	}
	return nil
}
