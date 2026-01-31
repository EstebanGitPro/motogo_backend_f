package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) SaveFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, querySaveFranchise, franchise.ID, franchise.Name, franchise.Description)
	if err != nil {
		log.Error(logger.LogFranchiseRepoSaveError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}
