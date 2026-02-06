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

func TestGetAllDepartments_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLocationService := new(mocks.MockLocationService)

	locationInteractor := interactor.NewLocationInteractor(mockLocationService)

	dbError := errors.New("database connection failed")

	// Mock expectations
	mockLocationService.On("GetAllDepartments", ctx).Return(nil, dbError)

	// Act
	result, err := locationInteractor.GetAllDepartments(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockLocationService.AssertExpectations(t)
}

func TestGetAllDepartments_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLocationService := new(mocks.MockLocationService)

	locationInteractor := interactor.NewLocationInteractor(mockLocationService)

	// Mock expectations
	mockLocationService.On("GetAllDepartments", ctx).Return([]domain.Department{}, nil)

	// Act
	result, err := locationInteractor.GetAllDepartments(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockLocationService.AssertExpectations(t)
}

// ============================================
// GetCitiesByDepartment Tests
// ============================================

func TestGetCitiesByDepartment_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLocationService := new(mocks.MockLocationService)

	locationInteractor := interactor.NewLocationInteractor(mockLocationService)

	departmentID := "dept-cundinamarca"
	expectedCities := []domain.City{
		{ID: "city-1", Name: "Bogotá", DepartmentID: departmentID},
		{ID: "city-2", Name: "Soacha", DepartmentID: departmentID},
		{ID: "city-3", Name: "Chía", DepartmentID: departmentID},
	}

	// Mock expectations
	mockLocationService.On("GetCitiesByDepartment", ctx, departmentID).Return(expectedCities, nil)

	// Act
	result, err := locationInteractor.GetCitiesByDepartment(ctx, departmentID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Bogotá", result[0].Name)

	mockLocationService.AssertExpectations(t)
}

func TestGetCitiesByDepartment_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLocationService := new(mocks.MockLocationService)

	locationInteractor := interactor.NewLocationInteractor(mockLocationService)

	departmentID := "dept-invalid"
	dbError := errors.New("database error")

	// Mock expectations
	mockLocationService.On("GetCitiesByDepartment", ctx, departmentID).Return(nil, dbError)

	// Act
	result, err := locationInteractor.GetCitiesByDepartment(ctx, departmentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockLocationService.AssertExpectations(t)
}

func TestGetCitiesByDepartment_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLocationService := new(mocks.MockLocationService)

	locationInteractor := interactor.NewLocationInteractor(mockLocationService)

	departmentID := "dept-empty"

	// Mock expectations
	mockLocationService.On("GetCitiesByDepartment", ctx, departmentID).Return([]domain.City{}, nil)

	// Act
	result, err := locationInteractor.GetCitiesByDepartment(ctx, departmentID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockLocationService.AssertExpectations(t)
}
