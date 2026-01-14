package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// UpdateBranch updates an existing branch
func (r *repository) UpdateBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtUpdateBranch).ExecContext(ctx,
		branch.FranchiseID,
		branch.Name,
		branch.EstablishmentType,
		branch.ProfileImageURL,
		branch.Status,
		branch.ID,
	)

	if err != nil {
		log.Error(logger.LogBranchRepoUpdateError, "error", err, "branch_id", branch.ID)
		return domain.ErrBranchCannotUpdate
	}

	return nil
}
