package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockServiceCatalogService is a mock implementation of input.ServiceCatalogService
type MockServiceCatalogService struct {
	mock.Mock
}

// BeginTx mocks the transaction start
func (m *MockServiceCatalogService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// GetAllServiceTypes returns all available service types (HU75)
func (m *MockServiceCatalogService) GetAllServiceTypes() []domain.ServiceType {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]domain.ServiceType)
}

// GetAllServices retrieves all services from catalog (HU63)
func (m *MockServiceCatalogService) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Service), args.Error(1)
}

// GetServicesByType retrieves services filtered by type (HU63)
func (m *MockServiceCatalogService) GetServicesByType(ctx context.Context, serviceType domain.ServiceType) ([]domain.Service, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Service), args.Error(1)
}

// GetServiceByID retrieves a service by its UUID (HU68)
func (m *MockServiceCatalogService) GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Service), args.Error(1)
}

// GetServicesByBranch retrieves services associated with a specific branch
func (m *MockServiceCatalogService) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BranchServiceInfo), args.Error(1)
}

// AssociateBranchServices associates services to a branch
func (m *MockServiceCatalogService) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	args := m.Called(ctx, tx, branchID, serviceIDs)
	return args.Error(0)
}

// DissociateBranchService removes a service from a branch
func (m *MockServiceCatalogService) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	args := m.Called(ctx, tx, branchID, serviceID)
	return args.Error(0)
}

// ValidateServiceIDs checks if all provided service IDs exist
func (m *MockServiceCatalogService) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	args := m.Called(ctx, serviceIDs)
	return args.Error(0)
}

// CheckServiceAssociation checks if a service is already associated with a branch
func (m *MockServiceCatalogService) CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error) {
	args := m.Called(ctx, branchID, serviceID)
	return args.Bool(0), args.Error(1)
}

// UpdateService updates an existing service in the catalog (HU68 - Admin only)
func (m *MockServiceCatalogService) UpdateService(ctx context.Context, tx output.Tx, service domain.Service) error {
	args := m.Called(ctx, tx, service)
	return args.Error(0)
}
