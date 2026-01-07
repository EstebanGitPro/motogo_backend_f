package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// SaveBranchBrands saves brands for a branch
func (r *repository) SaveBranchBrands(ctx context.Context, tx output.Tx, branchID string, brands []string) error {
	sqlTx := tx.(*common.SQLTx)

	for _, brand := range brands {
		brandID := uuid.Generate()
		_, err := sqlTx.StmtContext(ctx, r.stmtSaveBranchBrand).ExecContext(ctx, brandID, branchID, brand)
		if err != nil {
			log.Error(logger.LogBranchRepoBrandSaveError, "error", err, "branch_id", branchID, "brand", brand)
			return domain.ErrBranchCannotSave
		}
	}

	return nil
}
