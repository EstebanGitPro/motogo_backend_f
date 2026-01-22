package schedule_detail

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DeleteScheduleDetail deletes a schedule detail (HU8)
func (r *repository) DeleteScheduleDetail(ctx context.Context, tx output.Tx, detailID string) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtDeleteDetail).ExecContext(ctx, detailID)
	if err != nil {
		log.Error(logger.LogScheduleDetailRepoDeleteError, "detail_id", detailID, "error", err)
		return err
	}

	return nil
}