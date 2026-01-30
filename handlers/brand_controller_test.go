package handlers_test

import (
	"context"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// BrandInteractor Mock Integration Tests
// These tests verify the mock implementations work correctly
// Controller integration tests require full server setup
// ============================================

func TestMockBrandInteractor_GetAllBrands_Success(t *testing.T) {
	// Arrange
	mockBrand := new(mocks.MockBrandInteractor)

	testBrands := []domain.Brand{
		{ID: "f6a7b8c9-1111-4000-8000-000000000001", Name: "Honda"},
		{ID: "f6a7b8c9-2222-4000-8000-000000000002", Name: "Yamaha"},
		{ID: "f6a7b8c9-3333-4000-8000-000000000003", Name: "Suzuki"},
	}

	mockBrand.On("GetAllBrands", mock.Anything).Return(testBrands, nil)

	// Act
	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Honda", result[0].Name)
	assert.Equal(t, "Yamaha", result[1].Name)
	assert.Equal(t, "Suzuki", result[2].Name)
	mockBrand.AssertExpectations(t)
}

func TestMockBrandInteractor_GetAllBrands_Error(t *testing.T) {
	// Arrange
	mockBrand := new(mocks.MockBrandInteractor)
	mockBrand.On("GetAllBrands", mock.Anything).Return(nil, domain.ErrBrandNotFound)

	// Act
	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBrandNotFound, err)
	mockBrand.AssertExpectations(t)
}

func TestMockBrandInteractor_GetAllBrands_Empty(t *testing.T) {
	// Arrange
	mockBrand := new(mocks.MockBrandInteractor)
	mockBrand.On("GetAllBrands", mock.Anything).Return([]domain.Brand{}, nil)

	// Act
	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockBrand.AssertExpectations(t)
}

// ============================================
// LocationInteractor Mock Tests
// ============================================

func TestMockLocationInteractor_GetAllDepartments_Success(t *testing.T) {
	// Arrange
	mockLocation := new(mocks.MockLocationInteractor)

	testDepartments := []domain.Department{
		{ID: "dep-1", Name: "Antioquia"},
		{ID: "dep-2", Name: "Cundinamarca"},
	}

	mockLocation.On("GetAllDepartments", mock.Anything).Return(testDepartments, nil)

	// Act
	ctx := context.Background()
	result, err := mockLocation.GetAllDepartments(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockLocation.AssertExpectations(t)
}

func TestMockLocationInteractor_GetCitiesByDepartment_Success(t *testing.T) {
	// Arrange
	mockLocation := new(mocks.MockLocationInteractor)

	testCities := []domain.City{
		{ID: "city-1", Name: "Medellín", DepartmentID: "dep-1"},
		{ID: "city-2", Name: "Envigado", DepartmentID: "dep-1"},
	}

	mockLocation.On("GetCitiesByDepartment", mock.Anything, "dep-1").Return(testCities, nil)

	// Act
	ctx := context.Background()
	result, err := mockLocation.GetCitiesByDepartment(ctx, "dep-1")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Medellín", result[0].Name)
	mockLocation.AssertExpectations(t)
}

// ============================================
// MotorcycleInteractor Mock Tests
// ============================================

func TestMockMotorcycleInteractor_GetMotorcyclesByOwner_Success(t *testing.T) {
	// Arrange
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)

	testMotorcycles := []domain.Motorcycle{
		{ID: "moto-1", LicensePlate: "ABC123", OwnerID: "owner-1"},
		{ID: "moto-2", LicensePlate: "XYZ789", OwnerID: "owner-1"},
	}

	mockMotorcycle.On("GetMotorcyclesByOwner", mock.Anything, "owner-1").Return(testMotorcycles, nil)

	// Act
	ctx := context.Background()
	result, err := mockMotorcycle.GetMotorcyclesByOwner(ctx, "owner-1")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "ABC123", result[0].LicensePlate)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_GetMotorcycleByLicensePlate_Success(t *testing.T) {
	// Arrange
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)

	testMoto := &domain.Motorcycle{
		ID:           "moto-1",
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
	}

	mockMotorcycle.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(testMoto, nil)

	// Act
	ctx := context.Background()
	result, err := mockMotorcycle.GetMotorcycleByLicensePlate(ctx, "ABC123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ABC123", result.LicensePlate)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_DeleteMotorcycle_Success(t *testing.T) {
	// Arrange
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	mockMotorcycle.On("DeleteMotorcycle", mock.Anything, "moto-1", "owner-1").Return(nil)

	// Act
	ctx := context.Background()
	err := mockMotorcycle.DeleteMotorcycle(ctx, "moto-1", "owner-1")

	// Assert
	assert.NoError(t, err)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_DeleteMotorcycle_NotFound(t *testing.T) {
	// Arrange
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	mockMotorcycle.On("DeleteMotorcycle", mock.Anything, "moto-999", "owner-1").Return(domain.ErrMotorcycleNotFound)

	// Act
	ctx := context.Background()
	err := mockMotorcycle.DeleteMotorcycle(ctx, "moto-999", "owner-1")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
	mockMotorcycle.AssertExpectations(t)
}
