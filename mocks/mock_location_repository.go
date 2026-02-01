package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockLocationRepository is a mock implementation of output.LocationRepository
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockLocationRepository) SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	args := m.Called(ctx, tx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) UpdateLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	args := m.Called(ctx, tx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) CheckAddressExists(ctx context.Context, address string) (bool, error) {
	args := m.Called(ctx, address)
	return args.Bool(0), args.Error(1)
}

func (m *MockLocationRepository) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Department), args.Error(1)
}

func (m *MockLocationRepository) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	args := m.Called(ctx, departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.City), args.Error(1)
}

func (m *MockLocationRepository) ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error {
	args := m.Called(ctx, cityID, departmentID)
	return args.Error(0)
}

func (m *MockLocationRepository) GetDepartmentByID(ctx context.Context, departmentID string) (*domain.Department, error) {
	args := m.Called(ctx, departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Department), args.Error(1)
}
