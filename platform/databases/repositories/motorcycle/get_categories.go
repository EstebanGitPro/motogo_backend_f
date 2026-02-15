package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetDistinctCategories retrieves all distinct motorcycle categories with their line counts (HU41)
func (r *repository) GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error) {
	log.Info(logger.LogMotorcycleCatRepoQuery)

	rows, err := r.stmtGetDistinctCategories.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleCatRepoError, "error", err.Error())
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var categories []domain.MotorcycleCategory
	for rows.Next() {
		var cat domain.MotorcycleCategory
		if err := rows.Scan(&cat.Name, &cat.LineCount); err != nil {
			log.Error(logger.LogMotorcycleCatRepoScanError, "error", err.Error())
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleCatRepoIterError, "error", err.Error())
		return nil, err
	}

	return categories, nil
}

// GetLinesByCategory retrieves all motorcycle lines (models) for a specific category (HU41)
func (r *repository) GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error) {
	log.Info(logger.LogMotorcycleCatLinesRepoQuery, "category", categoryName)

	rows, err := r.stmtGetLinesByCategory.QueryContext(ctx, categoryName)
	if err != nil {
		log.Error(logger.LogMotorcycleCatLinesRepoError, "error", err.Error(), "category", categoryName)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var lines []domain.CategoryLine
	for rows.Next() {
		var line domain.CategoryLine
		if err := rows.Scan(&line.Model, &line.BrandName, &line.EngineDisplacement); err != nil {
			log.Error(logger.LogMotorcycleCatLinesRepoScanError, "error", err.Error())
			return nil, err
		}
		lines = append(lines, line)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleCatLinesRepoIterError, "error", err.Error())
		return nil, err
	}

	return lines, nil
}

// GetDistinctDisplacements retrieves all distinct engine displacement values with their reference counts (HU49)
func (r *repository) GetDistinctDisplacements(ctx context.Context) ([]domain.EngineDisplacementRange, error) {
	log.Info(logger.LogMotorcycleDispRepoQuery)

	rows, err := r.stmtGetDistinctDisplacements.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleDispRepoError, "error", err.Error())
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var displacements []domain.EngineDisplacementRange
	for rows.Next() {
		var d domain.EngineDisplacementRange
		if err := rows.Scan(&d.Range); err != nil {
			log.Error(logger.LogMotorcycleDispRepoScanError, "error", err.Error())
			return nil, err
		}
		displacements = append(displacements, d)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleDispRepoIterError, "error", err.Error())
		return nil, err
	}

	return displacements, nil
}
