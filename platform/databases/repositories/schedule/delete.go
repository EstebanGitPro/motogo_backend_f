package schedule

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DeleteSchedule removes a branch schedule by ID
func (r *repository) DeleteSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDeleteSchedule)

	_, err := stmt.ExecContext(ctx, scheduleID)
	if err != nil {
		log.Error(logger.LogScheduleRepoDeleteError, "schedule_id", scheduleID, "error", err)
		return err
	}
	return nil
}
