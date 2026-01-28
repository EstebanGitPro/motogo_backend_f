package schedule

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) SetActive(ctx context.Context, tx output.Tx, scheduleID string, active bool) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtSetActive)

	_, err := stmt.ExecContext(ctx, active, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleRepoActivateError, "schedule_id", scheduleID, "active", active, "error", err)
		return err
	}
	return nil
}
