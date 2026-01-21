package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// UpdateFranchise updates an existing franchise
func (r *repository) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtUpdateFranchise)

	_, err := stmt.ExecContext(ctx, franchise.Name, franchise.Description, franchise.ID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoUpdateError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}