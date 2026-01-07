package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// BrandInteractor handles brand catalog use cases
type BrandInteractor struct {
	brandService input.BrandService
	logger       logger.Logger
}

// NewBrandInteractor creates a new BrandInteractor instance
func NewBrandInteractor(brandService input.BrandService, log logger.Logger) *BrandInteractor {
	return &BrandInteractor{
		brandService: brandService,
		logger:       log,
	}
}

// GetAllBrands retrieves all brands from the catalog
func (i *BrandInteractor) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogBrandInteractorGetAll)

	brands, err := i.brandService.GetAllBrands(ctx)
	if err != nil {
		log.Error(logger.LogBrandInteractorGetAllError, "error", err)
		return nil, err
	}

	log.Success(logger.LogBrandInteractorGetAllOK, "brands_count", len(brands))
	return brands, nil
}
