package diagnostic_permission

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Deactivate sets active = FALSE for a diagnostic permission record by motorcycle and branch IDs.
// This operation is idempotent: if no matching row exists or it is already inactive, it succeeds silently.
func (r *repository) Deactivate(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDeactivate, motorcycleID, branchID)
	if err != nil {
		log.Error(logger.LogDiagPermRepoDeleteError, err)
		return domain.ErrPermissionCannotDelete
	}

	return nil
}
