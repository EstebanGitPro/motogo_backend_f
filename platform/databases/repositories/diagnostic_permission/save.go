package diagnostic_permission

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Save inserts or updates a diagnostic permission record (upsert)
func (r *repository) Save(ctx context.Context, tx output.Tx, permission *domain.DiagnosticPermission) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbPerm := FromDomain(permission)
	_, err := sqlTx.ExecContext(ctx, queryInsert,
		dbPerm.ID,
		dbPerm.MotorcycleID,
		dbPerm.BranchID,
		dbPerm.Active,
	)
	if err != nil {
		log.Error(logger.LogDiagPermRepoSaveError, err)
		return domain.ErrPermissionCannotSave
	}

	return nil
}
