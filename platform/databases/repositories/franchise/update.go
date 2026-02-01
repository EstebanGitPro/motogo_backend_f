package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateFranchise, franchise.Name, franchise.Description, franchise.ID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoUpdateError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}
