package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocationRepository is a mock for output.LocationRepository
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

// ============================================
// NewLocationService Tests
// ============================================

func TestNewLocationService(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)
	assert.NotNil(t, service)
}

// ============================================
// GetAllDepartments Tests
// ============================================

func TestGetAllDepartments_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	departments := []domain.Department{
		{ID: "dep-1", Name: "Antioquia"},
		{ID: "dep-2", Name: "Cundinamarca"},
	}
	mockRepo.On("GetAllDepartments", ctx).Return(departments, nil)

	result, err := service.GetAllDepartments(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Antioquia", result[0].Name)
}

func TestGetAllDepartments_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	dbError := errors.New("database error")
	mockRepo.On("GetAllDepartments", ctx).Return(nil, dbError)

	result, err := service.GetAllDepartments(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)
}

// ============================================
// GetCitiesByDepartment Tests
// ============================================

func TestGetCitiesByDepartment_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	cities := []domain.City{
		{ID: "city-1", Name: "Medellín"},
		{ID: "city-2", Name: "Bogotá"},
	}
	mockRepo.On("GetCitiesByDepartment", ctx, "dep-1").Return(cities, nil)

	result, err := service.GetCitiesByDepartment(ctx, "dep-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Medellín", result[0].Name)
}

func TestGetCitiesByDepartment_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	mockRepo.On("GetCitiesByDepartment", ctx, "dep-unknown").Return([]domain.City{}, nil)

	result, err := service.GetCitiesByDepartment(ctx, "dep-unknown")

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestGetCitiesByDepartment_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	dbError := errors.New("department not found")
	mockRepo.On("GetCitiesByDepartment", ctx, "invalid").Return(nil, dbError)

	result, err := service.GetCitiesByDepartment(ctx, "invalid")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// ValidateCityInDepartment Tests
// ============================================

func TestValidateCityInDepartment_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	mockRepo.On("ValidateCityInDepartment", ctx, "city-1", "dep-1").Return(nil)

	err := service.ValidateCityInDepartment(ctx, "city-1", "dep-1")

	assert.NoError(t, err)
}

func TestValidateCityInDepartment_Invalid(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	validationError := errors.New("city does not belong to department")
	mockRepo.On("ValidateCityInDepartment", ctx, "city-1", "dep-2").Return(validationError)

	err := service.ValidateCityInDepartment(ctx, "city-1", "dep-2")

	assert.Error(t, err)
	assert.Equal(t, validationError, err)
}
