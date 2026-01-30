package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetReferencesByBrandID retrieves all motorcycle references for a specific brand (HU40)
func (r *repository) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	log.Info(logger.LogMotorcycleRepoBrandLinesQuery, "brand_id", brandID)

	rows, err := r.stmtGetReferencesByBrandID.QueryContext(ctx, brandID)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoBrandLinesError, "error", err, "brand_id", brandID)
		return nil, err
	}
	defer rows.Close()

	var references []domain.MotorcycleReference
	for rows.Next() {
		var ref domain.MotorcycleReference
		err := rows.Scan(
			&ref.ID,
			&ref.BrandID,
			&ref.BrandName,
			&ref.Model,
			&ref.Category,
			&ref.EngineDisplacement,
		)
		if err != nil {
			log.Error(logger.LogMotorcycleRepoBrandLinesScanError, "error", err)
			return nil, err
		}
		references = append(references, ref)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleRepoBrandLinesIterError, "error", err)
		return nil, err
	}

	return references, nil
}
