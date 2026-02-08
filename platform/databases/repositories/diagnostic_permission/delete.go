package diagnostic_permission

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Delete removes a diagnostic permission record by motorcycle and branch IDs
func (r *repository) Delete(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryDelete, motorcycleID, branchID)
	if err != nil {
		log.Error(logger.LogDiagPermRepoDeleteError, err)
		return domain.ErrPermissionCannotDelete
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrPermissionNotFound
	}

	return nil
}
