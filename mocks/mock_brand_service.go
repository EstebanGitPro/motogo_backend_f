package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockBrandService is a mock implementation of input.BrandService
type MockBrandService struct {
	mock.Mock
}

// GetAllBrands retrieves all brands from the catalog
func (m *MockBrandService) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Brand), args.Error(1)
}

// ValidateBrandIDs checks if all provided brand IDs exist
func (m *MockBrandService) ValidateBrandIDs(ctx context.Context, brandIDs []string) error {
	args := m.Called(ctx, brandIDs)
	return args.Error(0)
}
