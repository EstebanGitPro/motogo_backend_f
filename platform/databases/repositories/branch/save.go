package branch

import (
	"context"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) SaveBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, querySaveBranch,
		branch.ID,
		branch.RepresentativeID,
		branch.FranchiseID,
		branch.Name,
		branch.EstablishmentType,
		branch.ProfileImageURL,
		branch.Status,
	)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.ErrDuplicateBranchName
		}
		log.Error(logger.LogBranchRepoSaveError, "error", err, "branch_id", branch.ID)
		return domain.ErrBranchCannotSave
	}

	return nil
}
