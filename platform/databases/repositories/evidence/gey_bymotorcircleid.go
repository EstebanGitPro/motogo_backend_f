package evidence

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByMotorcycleID retrieves all evidence for a motorcycle
func (r *repository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	rows, err := r.stmtGetByMotorcycleID.QueryContext(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceRepoListByMotoError, err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var evidences []domain.MotorcycleEvidence
	for rows.Next() {
		var evidence Evidence
		if err := rows.Scan(
			&evidence.ID,
			&evidence.MotorcycleID,
			&evidence.Angle,
			&evidence.ImageURL,
			&evidence.Description,
			&evidence.CreatedAt,
		); err != nil {
			log.Error(logger.LogEvidenceRepoScanError, err)
			return nil, err
		}
		evidences = append(evidences, evidence.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return evidences, nil
}

// CountByMotorcycleID counts the number of evidence for a motorcycle
func (r *repository) CountByMotorcycleID(ctx context.Context, motorcycleID string) (int, error) {
	var count int
	err := r.stmtCountByMotorcycleID.QueryRowContext(ctx, motorcycleID).Scan(&count)
	if err != nil {
		log.Error(logger.LogEvidenceRepoCountError, err)
		return 0, err
	}
	return count, nil
}
