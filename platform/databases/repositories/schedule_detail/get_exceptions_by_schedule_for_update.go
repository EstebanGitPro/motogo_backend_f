package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetExceptionsByScheduleIDForUpdate(
	ctx context.Context,
	tx output.Tx,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError,
			"schedule_id", scheduleID,
			"error", "invalid transaction type")
		return nil, domain.ErrInternalServer
	}

	rows, err := sqlTx.QueryContext(ctx, queryGetExceptionsForUpdate, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoGetBySchedError,
			"schedule_id", scheduleID,
			"error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanScheduleDetails(rows)
}
