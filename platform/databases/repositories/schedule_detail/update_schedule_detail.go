package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) UpdateScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateDetail,
		detail.DayOfWeek,
		detail.OpeningTime,
		detail.ClosingTime,
		detail.IsClosed,
		detail.Active,
		detail.ID,
	)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoUpdateError, "detail_id", detail.ID, "error", err)
		return err
	}

	return nil
}
