package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockLocationInteractor is a mock implementation of input.LocationInteractorInterface
type MockLocationInteractor struct {
	mock.Mock
}

// GetAllDepartments mocks the GetAllDepartments method
func (m *MockLocationInteractor) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Department), args.Error(1)
}

// GetCitiesByDepartment mocks the GetCitiesByDepartment method
func (m *MockLocationInteractor) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	args := m.Called(ctx, departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.City), args.Error(1)
}
