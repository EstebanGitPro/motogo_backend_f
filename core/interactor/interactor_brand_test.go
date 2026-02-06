package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
)

// ============================================
// GetAllBrands Tests
// ============================================

func TestGetAllBrands_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBrandService := new(mocks.MockBrandService)

	brandInteractor := interactor.NewBrandInteractor(mockBrandService)

	expectedBrands := []domain.Brand{
		{ID: "brand-1", Name: "Honda"},
		{ID: "brand-2", Name: "Yamaha"},
		{ID: "brand-3", Name: "Suzuki"},
		{ID: "brand-4", Name: "Kawasaki"},
	}

	// Mock expectations
	mockBrandService.On("GetAllBrands", ctx).Return(expectedBrands, nil)

	// Act
	result, err := brandInteractor.GetAllBrands(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 4)
	assert.Equal(t, "Honda", result[0].Name)

	mockBrandService.AssertExpectations(t)
}

func TestGetAllBrands_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBrandService := new(mocks.MockBrandService)

	brandInteractor := interactor.NewBrandInteractor(mockBrandService)

	dbError := errors.New("database connection failed")

	// Mock expectations
	mockBrandService.On("GetAllBrands", ctx).Return(nil, dbError)

	// Act
	result, err := brandInteractor.GetAllBrands(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockBrandService.AssertExpectations(t)
}

func TestGetAllBrands_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBrandService := new(mocks.MockBrandService)

	brandInteractor := interactor.NewBrandInteractor(mockBrandService)

	// Mock expectations
	mockBrandService.On("GetAllBrands", ctx).Return([]domain.Brand{}, nil)

	// Act
	result, err := brandInteractor.GetAllBrands(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockBrandService.AssertExpectations(t)
}
