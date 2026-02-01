package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDissociateSingleBranch, branchID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDissociateError, "branch_id", branchID, "error", err)
		return err
	}
	return nil
}
