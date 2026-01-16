package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServiceRepository is a mock implementation of output.ServiceRepository
type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Service), args.Error(1)
}

func (m *MockServiceRepository) GetServicesByType(ctx context.Context, serviceType string) ([]domain.Service, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Service), args.Error(1)
}

func (m *MockServiceRepository) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BranchServiceInfo), args.Error(1)
}

func (m *MockServiceRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockServiceRepository) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	args := m.Called(ctx, tx, branchID, serviceIDs)
	return args.Error(0)
}

func (m *MockServiceRepository) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	args := m.Called(ctx, tx, branchID, serviceID)
	return args.Error(0)
}

func (m *MockServiceRepository) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	args := m.Called(ctx, serviceIDs)
	return args.Error(0)
}

func (m *MockServiceRepository) CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error) {
	args := m.Called(ctx, branchID, serviceID)
	return args.Bool(0), args.Error(1)
}

// MockTx is a mock transaction
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewServiceCatalogService_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockServiceRepository)

	// Act
	service := services.NewServiceCatalogService(mockRepo)

	// Assert
	assert.NotNil(t, service)
}

func TestGetAllServiceTypes_ReturnsAllTypes(t *testing.T) {
	// Arrange
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	// Act
	types := service.GetAllServiceTypes()

	// Assert
	assert.Len(t, types, 8) // All 8 service types
	assert.Contains(t, types, domain.ServiceTypeMaintenance)
	assert.Contains(t, types, domain.ServiceTypeRepair)
	assert.Contains(t, types, domain.ServiceTypeTires)
}

func TestGetAllServices_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	expectedServices := []domain.Service{
		{ID: "uuid-1", Name: "Cambio de aceite", ServiceType: domain.ServiceTypeMaintenance},
		{ID: "uuid-2", Name: "Reparación de motor", ServiceType: domain.ServiceTypeRepair},
	}
	mockRepo.On("GetAllServices", ctx).Return(expectedServices, nil)

	// Act
	result, err := service.GetAllServices(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Cambio de aceite", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllServices_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	mockRepo.On("GetAllServices", ctx).Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetAllServices(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetServicesByType_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	expectedServices := []domain.Service{
		{ID: "uuid-1", Name: "Cambio de aceite", ServiceType: domain.ServiceTypeMaintenance},
	}
	mockRepo.On("GetServicesByType", ctx, "Mantenimiento").Return(expectedServices, nil)

	// Act
	result, err := service.GetServicesByType(ctx, domain.ServiceTypeMaintenance)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetServicesByBranch_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	expectedServices := []domain.BranchServiceInfo{
		{
			Service: domain.Service{ID: "uuid-1", Name: "Cambio de aceite"},
			AddedAt: "2026-01-15T10:30:00-05:00",
			Active:  true,
		},
	}
	mockRepo.On("GetServicesByBranch", ctx, "branch-123").Return(expectedServices, nil)

	// Act
	result, err := service.GetServicesByBranch(ctx, "branch-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "2026-01-15T10:30:00-05:00", result[0].AddedAt)
	mockRepo.AssertExpectations(t)
}

func TestGetServicesByBranch_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	mockRepo.On("GetServicesByBranch", ctx, "branch-123").Return(nil, errors.New("not found"))

	// Act
	result, err := service.GetServicesByBranch(ctx, "branch-123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestValidateServiceIDs_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	serviceIDs := []string{"uuid-1", "uuid-2"}
	mockRepo.On("ValidateServiceIDs", ctx, serviceIDs).Return(nil)

	// Act
	err := service.ValidateServiceIDs(ctx, serviceIDs)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestValidateServiceIDs_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	serviceIDs := []string{"invalid-uuid"}
	mockRepo.On("ValidateServiceIDs", ctx, serviceIDs).Return(domain.ErrServiceNotFound)

	// Act
	err := service.ValidateServiceIDs(ctx, serviceIDs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestCheckServiceAssociation_Exists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	mockRepo.On("CheckServiceAssociation", ctx, "branch-123", "service-456").Return(true, nil)

	// Act
	exists, err := service.CheckServiceAssociation(ctx, "branch-123", "service-456")

	// Assert
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestCheckServiceAssociation_NotExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)

	mockRepo.On("CheckServiceAssociation", ctx, "branch-123", "service-456").Return(false, nil)

	// Act
	exists, err := service.CheckServiceAssociation(ctx, "branch-123", "service-456")

	// Assert
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestBeginTx_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)
	mockTx := new(MockTx)

	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)

	// Act
	tx, err := service.BeginTx(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	mockRepo.AssertExpectations(t)
}

func TestAssociateBranchServices_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)
	mockTx := new(MockTx)

	serviceIDs := []string{"svc-1", "svc-2"}
	mockRepo.On("AssociateBranchServices", ctx, mockTx, "branch-123", serviceIDs).Return(nil)

	// Act
	err := service.AssociateBranchServices(ctx, mockTx, "branch-123", serviceIDs)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDissociateBranchService_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)
	mockTx := new(MockTx)

	mockRepo.On("DissociateBranchService", ctx, mockTx, "branch-123", "svc-1").Return(nil)

	// Act
	err := service.DissociateBranchService(ctx, mockTx, "branch-123", "svc-1")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDissociateBranchService_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockServiceRepository)
	service := services.NewServiceCatalogService(mockRepo)
	mockTx := new(MockTx)

	mockRepo.On("DissociateBranchService", ctx, mockTx, "branch-123", "svc-invalid").Return(domain.ErrServiceNotFound)

	// Act
	err := service.DissociateBranchService(ctx, mockTx, "branch-123", "svc-invalid")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFound, err)
	mockRepo.AssertExpectations(t)
}
