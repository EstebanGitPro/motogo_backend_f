package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

func (r *repository) SaveBranchDisplacementRanges(ctx context.Context, tx output.Tx, branchID string, ranges []string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	for _, dr := range ranges {
		id := uuid.Generate()
		_, err := sqlTx.ExecContext(ctx, querySaveBranchDisplacementRange, id, branchID, dr)
		if err != nil {
			log.Error(logger.LogBranchRepoDisplRangeSaveError, "error", err, "branch_id", branchID, "displacement_range", dr)
			return domain.ErrBranchCannotSave
		}
	}

	return nil
}
