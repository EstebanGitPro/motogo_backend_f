package evidence

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByID retrieves an evidence by its ID
func (r *repository) GetByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	var evidence Evidence

	err := r.stmtGetByID.QueryRowContext(ctx, evidenceID).Scan(
		&evidence.ID,
		&evidence.MotorcycleID,
		&evidence.Angle,
		&evidence.ImageURL,
		&evidence.Description,
		&evidence.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEvidenceNotFound
		}
		log.Error(logger.LogEvidenceRepoGetByIDError, err)
		return nil, err
	}

	result := evidence.ToDomain()
	return &result, nil
}
