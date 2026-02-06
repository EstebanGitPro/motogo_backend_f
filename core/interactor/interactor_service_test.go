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
// GetServiceTypes Tests (HU75)
// ============================================

func TestGetServiceTypes_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	expectedTypes := []domain.ServiceType{
		domain.ServiceTypeMaintenance,
		domain.ServiceTypeRepair,
		domain.ServiceTypeElectrical,
	}

	// Mock expectations
	mockService.On("GetAllServiceTypes").Return(expectedTypes)

	// Act
	result := serviceInteractor.GetServiceTypes(ctx)

	// Assert
	assert.Len(t, result, 3)
	assert.Contains(t, result, domain.ServiceTypeMaintenance)

	mockService.AssertExpectations(t)
}

// ============================================
// GetAllServices Tests (HU63)
// ============================================

func TestGetAllServices_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	expectedServices := []domain.Service{
		{ID: "svc-1", Name: "Cambio de Aceite", ServiceType: domain.ServiceTypeMaintenance, IsActive: true},
		{ID: "svc-2", Name: "Reparación de Motor", ServiceType: domain.ServiceTypeRepair, IsActive: true},
	}

	// Mock expectations
	mockService.On("GetAllServices", ctx).Return(expectedServices, nil)

	// Act
	result, err := serviceInteractor.GetAllServices(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Cambio de Aceite", result[0].Name)

	mockService.AssertExpectations(t)
}

func TestGetAllServices_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	dbError := errors.New("database connection failed")

	// Mock expectations
	mockService.On("GetAllServices", ctx).Return(nil, dbError)

	// Act
	result, err := serviceInteractor.GetAllServices(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockService.AssertExpectations(t)
}

// ============================================
// GetServicesByType Tests (HU63)
// ============================================

func TestGetServicesByType_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceType := domain.ServiceTypeMaintenance
	expectedServices := []domain.Service{
		{ID: "svc-1", Name: "Cambio de Aceite", ServiceType: serviceType, IsActive: true},
		{ID: "svc-2", Name: "Revisión General", ServiceType: serviceType, IsActive: true},
	}

	// Mock expectations
	mockService.On("GetServicesByType", ctx, serviceType).Return(expectedServices, nil)

	// Act
	result, err := serviceInteractor.GetServicesByType(ctx, serviceType)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, svc := range result {
		assert.Equal(t, serviceType, svc.ServiceType)
	}

	mockService.AssertExpectations(t)
}

func TestGetServicesByType_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceType := domain.ServiceTypeRepair
	dbError := errors.New("query failed")

	// Mock expectations
	mockService.On("GetServicesByType", ctx, serviceType).Return(nil, dbError)

	// Act
	result, err := serviceInteractor.GetServicesByType(ctx, serviceType)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockService.AssertExpectations(t)
}

// ============================================
// GetServiceByID Tests (HU68)
// ============================================

func TestGetServiceByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceID := "svc-123"
	expectedService := &domain.Service{
		ID:          serviceID,
		Name:        "Cambio de Aceite",
		Description: "Cambio de aceite y filtro",
		ServiceType: domain.ServiceTypeMaintenance,
		IsActive:    true,
	}

	// Mock expectations
	mockService.On("GetServiceByID", ctx, serviceID).Return(expectedService, nil)

	// Act
	result, err := serviceInteractor.GetServiceByID(ctx, serviceID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, serviceID, result.ID)
	assert.Equal(t, "Cambio de Aceite", result.Name)

	mockService.AssertExpectations(t)
}

func TestGetServiceByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceID := "non-existent"

	// Mock expectations
	mockService.On("GetServiceByID", ctx, serviceID).Return(nil, domain.ErrServiceNotFound)

	// Act
	result, err := serviceInteractor.GetServiceByID(ctx, serviceID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrServiceNotFound, err)

	mockService.AssertExpectations(t)
}

// ============================================
// GetServicesByBranch Tests
// ============================================

func TestGetServicesByBranch_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	branchID := "branch-123"
	expectedServices := []domain.BranchServiceInfo{
		{
			Service: domain.Service{ID: "svc-1", Name: "Cambio de Aceite"},
			AddedAt: "2024-01-15T10:00:00Z",
			Active:  true,
		},
		{
			Service: domain.Service{ID: "svc-2", Name: "Lavado"},
			AddedAt: "2024-01-10T08:00:00Z",
			Active:  true,
		},
	}

	// Mock expectations
	mockService.On("GetServicesByBranch", ctx, branchID).Return(expectedServices, nil)

	// Act
	result, err := serviceInteractor.GetServicesByBranch(ctx, branchID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockService.AssertExpectations(t)
}

func TestGetServicesByBranch_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	branchID := "branch-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockService.On("GetServicesByBranch", ctx, branchID).Return(nil, dbError)

	// Act
	result, err := serviceInteractor.GetServicesByBranch(ctx, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockService.AssertExpectations(t)
}

// ============================================
// UpdateService Tests (HU68 - Admin)
// ============================================

func TestUpdateService_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	service := domain.Service{
		ID:          "svc-123",
		Name:        "Cambio de Aceite Premium",
		Description: "Cambio de aceite sintético",
		ServiceType: domain.ServiceTypeMaintenance,
		IsActive:    true,
	}

	// Mock expectations
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdateService", ctx, mockTx, service).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := serviceInteractor.UpdateService(ctx, service)

	// Assert
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateService_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	service := domain.Service{
		ID:          "non-existent",
		Name:        "Service",
		ServiceType: domain.ServiceTypeMaintenance,
	}

	// Mock expectations
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdateService", ctx, mockTx, service).Return(domain.ErrServiceNotFound)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := serviceInteractor.UpdateService(ctx, service)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFound, err)

	mockService.AssertExpectations(t)
}

func TestUpdateService_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	service := domain.Service{
		ID:          "svc-123",
		Name:        "Service",
		ServiceType: domain.ServiceTypeMaintenance,
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockService.On("BeginTx", ctx).Return(nil, txError)

	// Act
	err := serviceInteractor.UpdateService(ctx, service)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInternalServer, err)

	mockService.AssertExpectations(t)
}

// ============================================
// ValidateServiceIDs Tests
// ============================================

func TestValidateServiceIDs_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceIDs := []string{"svc-1", "svc-2", "svc-3"}

	// Mock expectations
	mockService.On("ValidateServiceIDs", ctx, serviceIDs).Return(nil)

	// Act
	err := serviceInteractor.ValidateServiceIDs(ctx, serviceIDs)

	// Assert
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}

func TestValidateServiceIDs_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	serviceIDs := []string{"svc-1", "invalid-id"}

	// Mock expectations
	mockService.On("ValidateServiceIDs", ctx, serviceIDs).Return(domain.ErrServiceNotFound)

	// Act
	err := serviceInteractor.ValidateServiceIDs(ctx, serviceIDs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFound, err)

	mockService.AssertExpectations(t)
}

// ============================================
// AssociateBranchServices Tests
// ============================================

func TestAssociateBranchServices_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	branchID := "branch-123"
	serviceIDs := []string{"svc-1", "svc-2"}

	// Mock expectations
	mockService.On("AssociateBranchServices", ctx, mockTx, branchID, serviceIDs).Return(nil)

	// Act
	err := serviceInteractor.AssociateBranchServices(ctx, mockTx, branchID, serviceIDs)

	// Assert
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}

// ============================================
// DissociateBranchService Tests
// ============================================

func TestDissociateBranchService_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	branchID := "branch-123"
	serviceID := "svc-1"

	// Mock expectations
	mockService.On("DissociateBranchService", ctx, mockTx, branchID, serviceID).Return(nil)

	// Act
	err := serviceInteractor.DissociateBranchService(ctx, mockTx, branchID, serviceID)

	// Assert
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}

func TestDissociateBranchService_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	branchID := "branch-123"
	serviceID := "non-existent"

	// Mock expectations
	mockService.On("DissociateBranchService", ctx, mockTx, branchID, serviceID).Return(domain.ErrServiceNotFound)

	// Act
	err := serviceInteractor.DissociateBranchService(ctx, mockTx, branchID, serviceID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFound, err)

	mockService.AssertExpectations(t)
}

// ============================================
// BeginTx Tests
// ============================================

func TestBeginTx_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	mockTx := new(mocks.MockTx)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	// Mock expectations
	mockService.On("BeginTx", ctx).Return(mockTx, nil)

	// Act
	tx, err := serviceInteractor.BeginTx(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	mockService.AssertExpectations(t)
}

func TestBeginTx_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockServiceCatalogService)

	serviceInteractor := interactor.NewServiceInteractor(mockService)

	txError := errors.New("connection failed")

	// Mock expectations
	mockService.On("BeginTx", ctx).Return(nil, txError)

	// Act
	tx, err := serviceInteractor.BeginTx(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tx)

	mockService.AssertExpectations(t)
}
