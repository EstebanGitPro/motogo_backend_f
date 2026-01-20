package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockLocationService is a mock implementation of input.LocationService
type MockLocationService struct {
	mock.Mock
}

// GetAllDepartments retrieves all departments
func (m *MockLocationService) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Department), args.Error(1)
}

// GetCitiesByDepartment retrieves all cities for a specific department
func (m *MockLocationService) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	args := m.Called(ctx, departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.City), args.Error(1)
}

// ValidateCityInDepartment checks if the city belongs to the specified department
func (m *MockLocationService) ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error {
	args := m.Called(ctx, cityID, departmentID)
	return args.Error(0)
}
