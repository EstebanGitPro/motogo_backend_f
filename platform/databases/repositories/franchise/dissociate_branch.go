package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDissociateSingleBranch)

	_, err := stmt.ExecContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDissociateError, "branch_id", branchID, "error", err)
		return err
	}
	return nil
}
