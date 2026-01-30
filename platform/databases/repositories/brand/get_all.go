package brand

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	rows, err := r.stmtGetAllBrands.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return nil, err
	}
	defer rows.Close()

	var brands []domain.Brand
	for rows.Next() {
		var brand domain.Brand
		if err := rows.Scan(&brand.ID, &brand.Name); err != nil {
			log.Error(logger.LogDatabaseUnavailable, "error scanning brand", err)
			continue
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}
