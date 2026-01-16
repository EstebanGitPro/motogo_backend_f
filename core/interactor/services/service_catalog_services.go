package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

// serviceCatalogService implements input.ServiceCatalogService
type serviceCatalogService struct {
	serviceRepo output.ServiceRepository
}

// NewServiceCatalogService creates a new ServiceCatalogService instance
func NewServiceCatalogService(serviceRepo output.ServiceRepository) *serviceCatalogService {
	return &serviceCatalogService{
		serviceRepo: serviceRepo,
	}
}

// GetAllServiceTypes returns all available service types (HU75)
func (s *serviceCatalogService) GetAllServiceTypes() []domain.ServiceType {
	return domain.AllServiceTypes()
}

// GetAllServices retrieves all services from catalog (HU63)
func (s *serviceCatalogService) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	return s.serviceRepo.GetAllServices(ctx)
}

// GetServicesByType retrieves services filtered by type (HU63)
func (s *serviceCatalogService) GetServicesByType(ctx context.Context, serviceType domain.ServiceType) ([]domain.Service, error) {
	return s.serviceRepo.GetServicesByType(ctx, string(serviceType))
}

// GetServicesByBranch retrieves services associated with a specific branch
func (s *serviceCatalogService) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	return s.serviceRepo.GetServicesByBranch(ctx, branchID)
}

// BeginTx starts a new database transaction
func (s *serviceCatalogService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.serviceRepo.BeginTx(ctx)
}

// AssociateBranchServices associates services to a branch
func (s *serviceCatalogService) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	return s.serviceRepo.AssociateBranchServices(ctx, tx, branchID, serviceIDs)
}

// DissociateBranchService removes a service from a branch
func (s *serviceCatalogService) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	return s.serviceRepo.DissociateBranchService(ctx, tx, branchID, serviceID)
}

// ValidateServiceIDs checks if all provided service IDs exist
func (s *serviceCatalogService) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	return s.serviceRepo.ValidateServiceIDs(ctx, serviceIDs)
}

// CheckServiceAssociation checks if a service is already associated with a branch
func (s *serviceCatalogService) CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error) {
	return s.serviceRepo.CheckServiceAssociation(ctx, branchID, serviceID)
}
