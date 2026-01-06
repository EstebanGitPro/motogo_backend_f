package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var brandLogger logger.Logger = logger.NewSlogLogger()

// brandService implements input.BrandService
type brandService struct {
	repository output.BrandRepository
}

// NewBrandService creates a new BrandService instance
func NewBrandService(repository output.BrandRepository) input.BrandService {
	return &brandService{
		repository: repository,
	}
}

// GetAllBrands retrieves all brands from the catalog
func (s *brandService) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	brands, err := s.repository.GetAllBrands(ctx)
	if err != nil {
		brandLogger.Error(logger.LogDatabaseUnavailable, "error", err)
		return nil, err
	}
	return brands, nil
}

// ValidateBrandIDs checks if all provided brand IDs exist
func (s *brandService) ValidateBrandIDs(ctx context.Context, brandIDs []string) error {
	if len(brandIDs) == 0 {
		return nil
	}
	return s.repository.ValidateBrandIDs(ctx, brandIDs)
}
