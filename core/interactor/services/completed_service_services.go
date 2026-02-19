package services

import (
	"context"
	"errors"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// completedServiceService implements input.CompletedServiceService
type completedServiceService struct {
	repository     output.CompletedServiceRepository
	diagnosticRepo output.DiagnosticRepository
}

// NewCompletedServiceService creates a new CompletedServiceService instance
func NewCompletedServiceService(
	repo output.CompletedServiceRepository,
	diagnosticRepo output.DiagnosticRepository,
) *completedServiceService {
	return &completedServiceService{
		repository:     repo,
		diagnosticRepo: diagnosticRepo,
	}
}

// BeginTx starts a new database transaction
func (s *completedServiceService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

// ValidateBranchServices validates that all service IDs belong to the branch's active services
func (s *completedServiceService) ValidateBranchServices(ctx context.Context, branchID string, serviceIDs []string) error {
	if len(serviceIDs) == 0 {
		log.Warn(logger.LogCSServiceNoServices)
		return errors.New("at least one service_id is required")
	}

	return s.repository.ValidateBranchServices(ctx, branchID, serviceIDs)
}

// ValidateDiagnosticForMotorcycle validates that the diagnostic belongs to the motorcycle
func (s *completedServiceService) ValidateDiagnosticForMotorcycle(ctx context.Context, diagnosticID, motorcycleID string) error {
	diagnostic, err := s.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Warn(logger.LogCSServiceDiagGetErr, "diagnostic_id", diagnosticID, "error", err)
		return errors.New("diagnostic not found")
	}

	if diagnostic.MotorcycleID != motorcycleID {
		log.Warn(logger.LogCSServiceDiagNotForMoto, "diagnostic_id", diagnosticID, "motorcycle_id", motorcycleID)
		return errors.New("diagnostic does not belong to this motorcycle")
	}

	return nil
}

// SaveCompletedService saves a completed service to the database
func (s *completedServiceService) SaveCompletedService(ctx context.Context, tx output.Tx, service *domain.CompletedService) error {
	return s.repository.Save(ctx, tx, service)
}

// SaveItems saves completed service items (pivot records)
func (s *completedServiceService) SaveItems(ctx context.Context, tx output.Tx, items []domain.CompletedServiceItem) error {
	return s.repository.SaveItems(ctx, tx, items)
}

// SaveStatusHistory saves a status history entry
func (s *completedServiceService) SaveStatusHistory(ctx context.Context, tx output.Tx, history *domain.ServiceStatusHistory) error {
	return s.repository.SaveStatusHistory(ctx, tx, history)
}

// GetByID retrieves a completed service by ID
func (s *completedServiceService) GetByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	cs, err := s.repository.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Hydrate with items
	items, err := s.repository.GetItemsByCompletedServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	cs.Services = items

	return cs, nil
}

// GetByMotorcycleID retrieves completed services for a motorcycle (with items)
func (s *completedServiceService) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error) {
	services, err := s.repository.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		return nil, err
	}
	return s.hydrateItems(ctx, services)
}

// GetByBranchID retrieves completed services for a branch (with items)
func (s *completedServiceService) GetByBranchID(ctx context.Context, branchID string) ([]domain.CompletedService, error) {
	services, err := s.repository.GetByBranchID(ctx, branchID)
	if err != nil {
		return nil, err
	}
	return s.hydrateItems(ctx, services)
}

// hydrateItems loads the associated items for each completed service.
func (s *completedServiceService) hydrateItems(ctx context.Context, services []domain.CompletedService) ([]domain.CompletedService, error) {
	for i := range services {
		items, err := s.repository.GetItemsByCompletedServiceID(ctx, services[i].ID)
		if err != nil {
			return nil, err
		}
		services[i].Services = items
	}
	return services, nil
}

// ValidateNoActiveService checks that there is no active service (PENDIENTE/EN_PROCESO) for the motorcycle at the branch
func (s *completedServiceService) ValidateNoActiveService(ctx context.Context, branchID, motorcycleID string) error {
	hasActive, err := s.repository.HasActiveService(ctx, branchID, motorcycleID)
	if err != nil {
		log.Error(logger.LogCSInteractorActiveCheckErr, "error", err)
		return err
	}
	if hasActive {
		log.Warn(logger.LogCSInteractorActiveExists,
			"branch_id", branchID,
			"motorcycle_id", motorcycleID)
		return domain.ErrActiveServiceExists
	}
	return nil
}

// DeleteCompletedService applies hybrid delete strategy (HU65):
// - FINALIZADO/CANCELADO → soft delete (preserves ratings/history)
// - PENDIENTE/EN_PROCESO → hard delete (unblocks new service creation)
func (s *completedServiceService) DeleteCompletedService(ctx context.Context, tx output.Tx, serviceID string, status domain.ServiceStatus) error {
	if status == domain.ServiceStatusCompleted || status == domain.ServiceStatusCancelled {
		log.Info(logger.LogCSServiceDelStrategy, "strategy", "soft_delete", "status", string(status), "service_id", serviceID)
		return s.repository.SoftDelete(ctx, tx, serviceID)
	}
	log.Info(logger.LogCSServiceDelStrategy, "strategy", "hard_delete", "status", string(status), "service_id", serviceID)
	return s.repository.Delete(ctx, tx, serviceID)
}

// UpdateStatus updates the status and optionally the completion date (HU74)
func (s *completedServiceService) UpdateStatus(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time) error {
	return s.repository.UpdateStatus(ctx, tx, serviceID, status, completionDate)
}

// UpdateStatusWithPrice updates the status, completion date, and final price
func (s *completedServiceService) UpdateStatusWithPrice(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time, finalPrice *float64) error {
	return s.repository.UpdateStatusWithPrice(ctx, tx, serviceID, status, completionDate, finalPrice)
}

// UpdateDetails updates the quoted price, final price, and representative notes
func (s *completedServiceService) UpdateDetails(ctx context.Context, tx output.Tx, serviceID string, quotedPrice, finalPrice *float64, notes *string) error {
	return s.repository.UpdateDetails(ctx, tx, serviceID, quotedPrice, finalPrice, notes)
}

// GetStatusHistory retrieves the status transition history for a completed service (HU73)
func (s *completedServiceService) GetStatusHistory(ctx context.Context, serviceID string) ([]domain.ServiceStatusHistory, error) {
	return s.repository.GetStatusHistory(ctx, serviceID)
}
