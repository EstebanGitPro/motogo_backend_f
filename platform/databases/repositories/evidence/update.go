package evidence

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Update modifies an existing evidence record
func (r *repository) Update(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbEvidence := FromDomain(evidence)
	result, err := sqlTx.ExecContext(ctx, queryUpdate,
		dbEvidence.Angle,
		dbEvidence.ImageURL,
		dbEvidence.ID,
	)
	if err != nil {
		log.Error(logger.LogEvidenceRepoUpdateError, err)
		return domain.ErrEvidenceCannotUpdate
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEvidenceNotFound
	}

	return nil
}
