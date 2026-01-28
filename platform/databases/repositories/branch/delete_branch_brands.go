package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DeleteBranchBrands(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtDeleteBranchBrands).ExecContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchRepoBrandDelError, "error", err, "branch_id", branchID)
		return domain.ErrBranchCannotDelete
	}

	return nil
}
