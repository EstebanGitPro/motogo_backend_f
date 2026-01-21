package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DeleteFranchise removes a franchise by ID
func (r *repository) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDeleteFranchise)

	_, err := stmt.ExecContext(ctx, franchiseID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDeleteError, "franchise_id", franchiseID, "error", err)
		return err
	}
	return nil
}