package evidence

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Save inserts a new evidence record
func (r *repository) Save(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbEvidence := FromDomain(evidence)
	_, err := sqlTx.ExecContext(ctx, queryInsert,
		dbEvidence.ID,
		dbEvidence.MotorcycleID,
		dbEvidence.Angle,
		dbEvidence.ImageURL,
		dbEvidence.Description,
		dbEvidence.CreatedAt,
	)
	if err != nil {
		log.Error(logger.LogEvidenceRepoSaveError, err)
		return domain.ErrEvidenceCannotSave
	}

	return nil
}
