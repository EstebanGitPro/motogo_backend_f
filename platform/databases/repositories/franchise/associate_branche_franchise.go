package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) AssociateBranchesToFranchise(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	for _, branchID := range branchIDs {
		_, err := sqlTx.ExecContext(ctx, queryAssociateBranchToFranchise, franchiseID, branchID)
		if err != nil {
			log.Error(logger.LogFranchiseRepoAssociateError,
				"franchise_id", franchiseID, "branch_id", branchID, "error", err)
			return err
		}
	}
	return nil
}
