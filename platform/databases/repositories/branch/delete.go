package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DeleteBranch deletes a branch by ID
func (r *repository) DeleteBranch(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtDeleteBranch).ExecContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchRepoDeleteError, "error", err, "branch_id", branchID)
		return domain.ErrBranchCannotDelete
	}

	return nil
}
