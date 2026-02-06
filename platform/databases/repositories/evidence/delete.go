package evidence

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Delete removes an evidence record
func (r *repository) Delete(ctx context.Context, tx output.Tx, evidenceID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryDelete, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceRepoDeleteError, err)
		return domain.ErrEvidenceCannotDelete
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEvidenceNotFound
	}

	return nil
}
