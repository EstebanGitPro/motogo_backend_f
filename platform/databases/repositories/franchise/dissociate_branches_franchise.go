package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DissociateBranchesFromFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDissociateBranchesFromFranchise, franchiseID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDissociateError, "franchise_id", franchiseID, "error", err)
		return err
	}
	return nil
}
