package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DeleteBranchDisplacementRanges(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDeleteBranchDisplacementRanges, branchID)
	if err != nil {
		log.Error(logger.LogBranchRepoDisplRangeDelError, "error", err, "branch_id", branchID)
		return domain.ErrBranchCannotDelete
	}

	return nil
}
