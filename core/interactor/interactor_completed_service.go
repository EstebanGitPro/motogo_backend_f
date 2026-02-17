package interactor

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// CompletedServiceInteractor handles completed service use cases (HU64)
type CompletedServiceInteractor struct {
	service input.CompletedServiceService
}

// NewCompletedServiceInteractor creates a new CompletedServiceInteractor instance
func NewCompletedServiceInteractor(service input.CompletedServiceService) *CompletedServiceInteractor {
	return &CompletedServiceInteractor{
		service: service,
	}
}

// RegisterCompletedService registers a new completed service with its associated service items (HU64)
// Creates the completed_service record, pivot items, and initial status history entry in a transaction
func (i *CompletedServiceInteractor) RegisterCompletedService(ctx context.Context, cs *domain.CompletedService, serviceIDs []string, personID string) (*domain.CompletedService, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorRegStart,
		"branch_id", cs.BranchID,
		"motorcycle_id", cs.MotorcycleID,
		"services_count", len(serviceIDs))

	// STEP 1: Validate that all services belong to the branch (before transaction)
	if err := i.service.ValidateBranchServices(ctx, cs.BranchID, serviceIDs); err != nil {
		log.Warn(logger.LogCSInteractorValidateSvcErr, "error", err)
		return nil, err
	}
	log.Debug(logger.LogCSInteractorSvcValidated, "count", len(serviceIDs))

	// STEP 2: Validate diagnostic belongs to motorcycle (if provided)
	if cs.DiagnosticID != nil && *cs.DiagnosticID != "" {
		if err := i.service.ValidateDiagnosticForMotorcycle(ctx, *cs.DiagnosticID, cs.MotorcycleID); err != nil {
			log.Warn(logger.LogCSInteractorValidateDiagErr, "error", err)
			return nil, err
		}
		log.Debug(logger.LogCSInteractorDiagValidated, "diagnostic_id", *cs.DiagnosticID)
	}

	// STEP 3: Validate no active service exists for the same motorcycle+branch
	if err := i.service.ValidateNoActiveService(ctx, cs.BranchID, cs.MotorcycleID); err != nil {
		log.Warn(logger.LogCSInteractorActiveCheckErr, "error", err)
		return nil, err
	}

	// STEP 4: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogCSInteractorTxError, "error", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogCSInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogCSInteractorRollbackOK)
			}
		}
	}()

	// STEP 4: Generate ID and set defaults
	cs.SetID()
	cs.Status = domain.ServiceStatusPending
	cs.RequestDate = time.Now()

	// STEP 5: Save the completed service header
	if err = i.service.SaveCompletedService(ctx, tx, cs); err != nil {
		log.Error(logger.LogCSInteractorSaveError, "error", err)
		return nil, err
	}
	log.Debug(logger.LogCSInteractorSaved, "id", cs.ID)

	// STEP 6: Create and save pivot items
	items := make([]domain.CompletedServiceItem, len(serviceIDs))
	for idx, svcID := range serviceIDs {
		items[idx] = domain.CompletedServiceItem{
			CompletedServiceID: cs.ID,
			ServiceID:          svcID,
		}
		items[idx].SetID()
	}

	if err = i.service.SaveItems(ctx, tx, items); err != nil {
		log.Error(logger.LogCSInteractorSaveItemsErr, "error", err)
		return nil, err
	}
	log.Debug(logger.LogCSInteractorItemsSaved, "count", len(items))

	// STEP 7: Create initial status history entry
	history := &domain.ServiceStatusHistory{
		CompletedServiceID: cs.ID,
		PreviousStatus:     nil, // First entry — no previous status
		NewStatus:          domain.ServiceStatusPending,
		CreatedBy:          personID, // Representative who registered the service
	}
	history.SetID()

	if err = i.service.SaveStatusHistory(ctx, tx, history); err != nil {
		log.Error(logger.LogCSInteractorSaveHistoryErr, "error", err)
		return nil, err
	}

	// STEP 8: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogCSInteractorCommitError, "error", err)
		return nil, err
	}

	cs.Services = items

	log.Success(logger.LogCSInteractorRegSuccess,
		"id", cs.ID,
		"branch_id", cs.BranchID,
		"motorcycle_id", cs.MotorcycleID,
		"services_count", len(items))

	err = nil // Ensure defer doesn't execute rollback
	return cs, nil
}

// GetCompletedServiceByID retrieves a completed service by its ID
func (i *CompletedServiceInteractor) GetCompletedServiceByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorGetByID, "id", serviceID)

	cs, err := i.service.GetByID(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogCSInteractorGetByIDErr, "error", err, "id", serviceID)
		return nil, err
	}

	log.Success(logger.LogCSInteractorGetByIDOK, "id", serviceID)
	return cs, nil
}

// GetCompletedServicesByBranch retrieves completed services for a branch
func (i *CompletedServiceInteractor) GetCompletedServicesByBranch(ctx context.Context, branchID string) ([]domain.CompletedService, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorGetByBranch, "branch_id", branchID)

	services, err := i.service.GetByBranchID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogCSInteractorGetByBranchErr, "error", err, "branch_id", branchID)
		return nil, err
	}

	log.Success(logger.LogCSInteractorGetByBranchOK, "branch_id", branchID, "count", len(services))
	return services, nil
}

// GetCompletedServicesByMotorcycle retrieves completed services for a motorcycle
func (i *CompletedServiceInteractor) GetCompletedServicesByMotorcycle(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorGetByMoto, "motorcycle_id", motorcycleID)

	services, err := i.service.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogCSInteractorGetByMotoErr, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogCSInteractorGetByMotoOK, "motorcycle_id", motorcycleID, "count", len(services))
	return services, nil
}

// DeleteCompletedService orchestrates the deletion of a completed service (HU65).
// The hybrid strategy (soft vs hard delete) is handled by the service layer.
func (i *CompletedServiceInteractor) DeleteCompletedService(ctx context.Context, serviceID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorDelStart, "service_id", serviceID)

	// STEP 1: Verify service exists
	cs, err := i.service.GetByID(ctx, serviceID)
	if err != nil {
		log.Warn(logger.LogCSInteractorDelNotFound, "service_id", serviceID, "error", err)
		return domain.ErrCompletedServiceNotFound
	}

	// STEP 2: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogCSInteractorDelTxErr, "error", err)
		return domain.ErrCompletedServiceCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogCSInteractorDelRollbackErr,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogCSInteractorDelRollbackOK)
			}
		}
	}()

	// STEP 3: Delete — service layer decides strategy based on status
	if err = i.service.DeleteCompletedService(ctx, tx, serviceID, cs.Status); err != nil {
		log.Error(logger.LogCSInteractorDelError, "service_id", serviceID, "error", err)
		return domain.ErrCompletedServiceCannotDelete
	}

	// STEP 4: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogCSInteractorDelCommitErr, "error", err)
		return domain.ErrCompletedServiceCannotDelete
	}

	log.Success(logger.LogCSInteractorDelSuccess, "service_id", serviceID)
	err = nil // Ensure defer doesn't execute rollback
	return nil
}

// TransitionStatus validates and applies a status transition for a completed service (HU74)
func (i *CompletedServiceInteractor) TransitionStatus(ctx context.Context, serviceID string, newStatus string, personID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorStatusStart,
		"service_id", serviceID,
		"new_status", newStatus)

	// STEP 1: Verify service exists and get current status
	cs, err := i.service.GetByID(ctx, serviceID)
	if err != nil {
		log.Warn(logger.LogCSInteractorStatusNotFound, "service_id", serviceID, "error", err)
		return domain.ErrCompletedServiceNotFound
	}

	// STEP 2: Validate transition
	if !domain.IsValidTransition(cs.Status, domain.ServiceStatus(newStatus)) {
		log.Warn(logger.LogCSInteractorStatusInvalid,
			"from", cs.Status,
			"to", newStatus,
			"service_id", serviceID)
		return domain.ErrInvalidStatusTransition
	}

	// STEP 3: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogCSInteractorStatusTxErr, "error", err)
		return domain.ErrInvalidStatusTransition
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogCSInteractorStatusRbErr,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogCSInteractorStatusRbOK)
			}
		}
	}()

	// STEP 4: Set completion_date if transitioning to FINALIZADO
	var completionDate *time.Time
	if domain.ServiceStatus(newStatus) == domain.ServiceStatusCompleted {
		now := time.Now()
		completionDate = &now
	}

	// STEP 5: Update status in DB
	if err = i.service.UpdateStatus(ctx, tx, serviceID, newStatus, completionDate); err != nil {
		log.Error(logger.LogCSInteractorStatusUpdErr, "error", err)
		return domain.ErrInvalidStatusTransition
	}

	// STEP 6: Insert transition history
	previousStatus := cs.Status
	history := &domain.ServiceStatusHistory{
		CompletedServiceID: serviceID,
		PreviousStatus:     &previousStatus,
		NewStatus:          domain.ServiceStatus(newStatus),
		CreatedBy:          personID,
	}
	history.SetID()

	if err = i.service.SaveStatusHistory(ctx, tx, history); err != nil {
		log.Error(logger.LogCSInteractorStatusHistErr, "error", err)
		return domain.ErrInvalidStatusTransition
	}

	// STEP 7: Commit
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogCSInteractorStatusCommErr, "error", err)
		return domain.ErrInvalidStatusTransition
	}

	log.Success(logger.LogCSInteractorStatusSuccess,
		"service_id", serviceID,
		"from", cs.Status,
		"to", newStatus)

	err = nil
	return nil
}

// GetStatusHistory retrieves the status transition history for a completed service (HU73)
func (i *CompletedServiceInteractor) GetStatusHistory(ctx context.Context, serviceID string) ([]domain.ServiceStatusHistory, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorTransStart, "service_id", serviceID)

	// STEP 1: Verify service exists
	_, err := i.service.GetByID(ctx, serviceID)
	if err != nil {
		log.Warn(logger.LogCSInteractorTransNotFound, "service_id", serviceID, "error", err)
		return nil, domain.ErrCompletedServiceNotFound
	}

	// STEP 2: Get history
	history, err := i.service.GetStatusHistory(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogCSInteractorTransError, "service_id", serviceID, "error", err)
		return nil, err
	}

	log.Success(logger.LogCSInteractorTransSuccess,
		"service_id", serviceID,
		"count", len(history))

	return history, nil
}
