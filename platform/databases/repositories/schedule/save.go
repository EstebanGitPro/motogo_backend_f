package schedule

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) SaveSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, querySaveSchedule,
		schedule.ID,
		schedule.BranchID,
		schedule.Active,
		schedule.StartDate,
		schedule.EndDate, // nullable
	)
	if err != nil {
		log.Error(logger.LogScheduleRepoSaveError, "schedule_id", schedule.ID, "branch_id", schedule.BranchID, "error", err)
		return err
	}
	return nil
}
