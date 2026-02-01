package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) SaveScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, querySaveDetail,
		detail.ID,
		detail.ScheduleID,
		detail.EntryType,
		detail.DayOfWeek,
		detail.ExceptionStartDate,
		detail.ExceptionEndDate,
		detail.OpeningTime,
		detail.ClosingTime,
		detail.IsClosed,
		detail.Active,
	)

	if err != nil {
		log.Error(logger.LogScheduleDetailRepoSaveError, "detail_id", detail.ID, "error", err)
		return err
	}

	return nil
}
