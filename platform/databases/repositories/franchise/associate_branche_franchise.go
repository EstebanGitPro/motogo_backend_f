package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// AssociateBranchesToFranchise updates branches to link them to a franchise
func (r *repository) AssociateBranchesToFranchise(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtAssociateBranchToFranchise)

	for _, branchID := range branchIDs {
		_, err := stmt.ExecContext(ctx, franchiseID, branchID)
		if err != nil {
			log.Error(logger.LogFranchiseRepoAssociateError,
				"franchise_id", franchiseID, "branch_id", branchID, "error", err)
			return err
		}
	}
	return nil
}
