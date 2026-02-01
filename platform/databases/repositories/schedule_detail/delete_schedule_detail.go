package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DeleteScheduleDetail(ctx context.Context, tx output.Tx, detailID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDeleteDetail, detailID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoDeleteError, "detail_id", detailID, "error", err)
		return err
	}

	return nil
}
