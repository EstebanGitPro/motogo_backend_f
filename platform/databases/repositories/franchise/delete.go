package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDeleteFranchise, franchiseID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDeleteError, "franchise_id", franchiseID, "error", err)
		return err
	}
	return nil
}
