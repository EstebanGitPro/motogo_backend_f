package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// SaveFranchise inserts a new franchise
func (r *repository) SaveFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtSaveFranchise)

	_, err := stmt.ExecContext(ctx, franchise.ID, franchise.Name, franchise.Description)
	if err != nil {
		log.Error(logger.LogFranchiseRepoSaveError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}