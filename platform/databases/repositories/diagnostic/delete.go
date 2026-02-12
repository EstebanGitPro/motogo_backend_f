package diagnostic

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Delete removes a diagnostic record (cascades to evidence via FK)
func (r *repository) Delete(ctx context.Context, tx output.Tx, diagnosticID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryDelete, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoDeleteError, err)
		return domain.ErrDiagnosticCannotDelete
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrDiagnosticNotFound
	}

	return nil
}
