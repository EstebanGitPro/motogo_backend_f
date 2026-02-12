package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// ServiceInteractor handles service catalog use cases
type ServiceInteractor struct {
	serviceCatalogService input.ServiceCatalogService
}

// NewServiceInteractor creates a new ServiceInteractor instance
func NewServiceInteractor(serviceCatalogService input.ServiceCatalogService) *ServiceInteractor {
	return &ServiceInteractor{
		serviceCatalogService: serviceCatalogService,
	}
}

// GetServiceTypes retrieves all available service types (HU75)
func (i *ServiceInteractor) GetServiceTypes(ctx context.Context) []domain.ServiceType {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogServiceInteractorGetTypes)

	types := i.serviceCatalogService.GetAllServiceTypes()

	log.Success(logger.LogServiceInteractorGetTypesOK, "types_count", len(types))
	return types
}

// GetAllServices retrieves all services from the catalog (HU63)
func (i *ServiceInteractor) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogServiceInteractorGetAll)

	services, err := i.serviceCatalogService.GetAllServices(ctx)
	if err != nil {
		log.Error(logger.LogServiceInteractorGetAllError, "error", err)
		return nil, err
	}

	log.Success(logger.LogServiceInteractorGetAllOK, "services_count", len(services))
	return services, nil
}

// GetServicesByType retrieves services filtered by type (HU63)
func (i *ServiceInteractor) GetServicesByType(ctx context.Context, serviceType domain.ServiceType) ([]domain.Service, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogServiceInteractorGetByType, "type", serviceType)

	services, err := i.serviceCatalogService.GetServicesByType(ctx, serviceType)
	if err != nil {
		log.Error(logger.LogServiceInteractorGetByTypeError, "error", err)
		return nil, err
	}

	log.Success(logger.LogServiceInteractorGetByTypeOK, "services_count", len(services), "type", serviceType)
	return services, nil
}

// GetServicesByBranch retrieves services associated with a specific branch
func (i *ServiceInteractor) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchServicesIntGetByBranch, "branch_id", branchID)

	services, err := i.serviceCatalogService.GetServicesByBranch(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchServicesIntGetByBranchErr, "error", err, "branch_id", branchID)
		return nil, err
	}

	log.Success(logger.LogBranchServicesIntGetByBranchOK, "branch_id", branchID, "services_count", len(services))
	return services, nil
}

// BeginTx starts a new database transaction
func (i *ServiceInteractor) BeginTx(ctx context.Context) (output.Tx, error) {
	return i.serviceCatalogService.BeginTx(ctx)
}

// ValidateServiceIDs checks if all provided service IDs exist
func (i *ServiceInteractor) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	return i.serviceCatalogService.ValidateServiceIDs(ctx, serviceIDs)
}

// AssociateBranchServices associates services to a branch
func (i *ServiceInteractor) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	return i.serviceCatalogService.AssociateBranchServices(ctx, tx, branchID, serviceIDs)
}

// DissociateBranchService removes a service from a branch
func (i *ServiceInteractor) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	return i.serviceCatalogService.DissociateBranchService(ctx, tx, branchID, serviceID)
}

// GetServiceByID retrieves a service by UUID (HU68)
func (i *ServiceInteractor) GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogServiceInteractorGetByID, "service_id", serviceID)

	service, err := i.serviceCatalogService.GetServiceByID(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogServiceInteractorGetByIDError, "error", err, "service_id", serviceID)
		return nil, err
	}

	log.Success(logger.LogServiceInteractorGetByIDOK, "service_id", serviceID)
	return service, nil
}

// UpdateService updates a service in the catalog (HU68 - Admin only)
func (i *ServiceInteractor) UpdateService(ctx context.Context, service domain.Service) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogServiceInteractorUpdate, "service_id", service.ID)

	// Begin transaction
	tx, txErr := i.serviceCatalogService.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogServiceInteractorUpdateError, "error", txErr)
		return domain.ErrInternalServer
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogServiceInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogServiceInteractorRollbackOK)
			}
		}
	}()

	// Update service
	if err = i.serviceCatalogService.UpdateService(ctx, tx, service); err != nil {
		log.Error(logger.LogServiceInteractorUpdateError, "error", err, "service_id", service.ID)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogServiceInteractorUpdateError, "error", err)
		return domain.ErrInternalServer
	}

	log.Success(logger.LogServiceInteractorUpdateOK, "service_id", service.ID)

	err = nil
	return nil
}
