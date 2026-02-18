package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockBrandInteractor is a mock implementation of input.BrandLister
type MockBrandInteractor struct {
	mock.Mock
}

// GetAllBrands mocks the GetAllBrands method
func (m *MockBrandInteractor) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Brand), args.Error(1)
}
