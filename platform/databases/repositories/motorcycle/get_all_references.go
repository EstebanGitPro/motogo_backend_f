package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetAllReferences retrieves all motorcycle references from the catalog (HU50)
func (r *repository) GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	rows, err := r.stmtGetAllReferences.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoGetAllRefQuery, "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var references []domain.MotorcycleReference
	for rows.Next() {
		var ref MotorcycleReference
		err := rows.Scan(
			&ref.ID,
			&ref.BrandID,
			&ref.BrandName,
			&ref.Model,
			&ref.Category,
			&ref.EngineDisplacement,
		)
		if err != nil {
			log.Error(logger.LogMotorcycleRepoGetAllRefScanError, "error", err.Error())
			return nil, err
		}

		// Convert to domain model
		domainRef := domain.MotorcycleReference{
			ID:      ref.ID.String,
			BrandID: ref.BrandID.String,
		}
		if ref.BrandName.Valid {
			domainRef.BrandName = ref.BrandName.String
		}
		if ref.Model.Valid {
			domainRef.Model = ref.Model.String
		}
		if ref.Category.Valid {
			domainRef.Category = ref.Category.String
		}
		if ref.EngineDisplacement.Valid {
			domainRef.EngineDisplacement = int(ref.EngineDisplacement.Int64)
		}

		references = append(references, domainRef)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleRepoGetAllRefIterError, "error", err.Error())
		return nil, err
	}

	return references, nil
}
