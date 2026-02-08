package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DiagnosticInteractor handles diagnostic-related use cases (HU11-14)
type DiagnosticInteractor struct {
	diagnosticService input.DiagnosticService
}

// NewDiagnosticInteractor creates a new DiagnosticInteractor instance
func NewDiagnosticInteractor(
	diagnosticService input.DiagnosticService,
) *DiagnosticInteractor {
	return &DiagnosticInteractor{
		diagnosticService: diagnosticService,
	}
}

// RegisterDiagnostic creates a new diagnostic or updates an existing one for the same motorcycle+branch (UPSERT)
// If a diagnostic already exists for the same motorcycle and branch, it updates the existing record
// instead of creating a duplicate. This is transparent to the caller.
func (i *DiagnosticInteractor) RegisterDiagnostic(ctx context.Context, motorcycleID, branchID, ownerID string, problemDescription *string, evidenceURLs []string) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorCreateStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// 1. Validate motorcycle ownership
	if err := i.diagnosticService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// 2. Validate branch exists
	if err := i.diagnosticService.ValidateBranchExists(ctx, branchID); err != nil {
		return nil, err
	}

	// 3. Begin transaction
	tx, err := i.diagnosticService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 4. Upsert diagnostic via service (encapsulates create-or-update + evidence logic)
	diagnostic, err := i.diagnosticService.UpsertDiagnostic(ctx, tx, motorcycleID, branchID, problemDescription, evidenceURLs)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorSaveError, "error", err)
		return nil, err
	}

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	log.Success(logger.LogDiagnosticInteractorCreateSuccess, "id", diagnostic.ID, "motorcycle_id", motorcycleID)
	err = nil
	return diagnostic, nil
}

// GetDiagnosticByID retrieves a diagnostic with its evidence (HU14)
func (i *DiagnosticInteractor) GetDiagnosticByID(ctx context.Context, diagnosticID, ownerID string) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorGetStart, "diagnostic_id", diagnosticID)

	// 1. Get diagnostic with evidence
	diagnostic, err := i.diagnosticService.GetDiagnosticByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, err
	}

	// 2. Validate ownership through motorcycle
	if err := i.diagnosticService.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	log.Success(logger.LogDiagnosticInteractorGetSuccess, "diagnostic_id", diagnosticID)
	return diagnostic, nil
}

// ListDiagnosticsByMotorcycle retrieves all diagnostics for a motorcycle with evidence (HU14)
func (i *DiagnosticInteractor) ListDiagnosticsByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorListStart, "motorcycle_id", motorcycleID)

	// 1. Validate motorcycle ownership
	if err := i.diagnosticService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// 2. Get diagnostics with evidence via service
	diagnostics, err := i.diagnosticService.GetDiagnosticsByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagnosticInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(diagnostics))
	return diagnostics, nil
}

// UpdateDiagnostic updates an existing diagnostic (HU12)
func (i *DiagnosticInteractor) UpdateDiagnostic(ctx context.Context, diagnosticID, ownerID string, updates *domain.Diagnostic) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorUpdateStart, "diagnostic_id", diagnosticID)

	// 1. Get existing diagnostic
	diagnostic, err := i.diagnosticService.GetDiagnosticByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// 2. Validate ownership through motorcycle
	if err := i.diagnosticService.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// 3. Apply updates via service (field-by-field patching)
	i.diagnosticService.ApplyDiagnosticUpdates(diagnostic, updates)

	// 4. Begin transaction
	tx, err := i.diagnosticService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 5. Update diagnostic via service
	if err = i.diagnosticService.UpdateDiagnostic(ctx, tx, diagnostic); err != nil {
		log.Error(logger.LogDiagnosticInteractorUpdateError, "error", err, "diagnostic_id", diagnosticID)
		return nil, err
	}

	// 6. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	log.Success(logger.LogDiagnosticInteractorUpdateSuccess, "diagnostic_id", diagnosticID)
	err = nil
	return diagnostic, nil
}

// DeleteDiagnostic deletes a diagnostic and its evidence (HU13)
func (i *DiagnosticInteractor) DeleteDiagnostic(ctx context.Context, diagnosticID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorDeleteStart, "diagnostic_id", diagnosticID)

	// 1. Get diagnostic to find motorcycle for ownership check
	diagnostic, err := i.diagnosticService.GetDiagnosticByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return err
	}

	// 2. Validate ownership through motorcycle
	if err := i.diagnosticService.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return domain.ErrDiagnosticNotFound
	}

	// 3. Begin transaction
	tx, err := i.diagnosticService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return domain.ErrDiagnosticCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 4. Delete diagnostic via service (FK CASCADE handles evidence)
	if err = i.diagnosticService.DeleteDiagnostic(ctx, tx, diagnosticID); err != nil {
		log.Error(logger.LogDiagnosticInteractorDeleteError, "error", err, "diagnostic_id", diagnosticID)
		return err
	}

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return domain.ErrDiagnosticCannotDelete
	}

	log.Success(logger.LogDiagnosticInteractorDeleteSuccess, "diagnostic_id", diagnosticID)
	err = nil
	return nil
}

// ListDiagnosticsByMotorcycleID retrieves diagnostics with evidence for a motorcycle (used by motorcycle lookup)
// This method does NOT validate ownership - it's designed for workshop representatives
func (i *DiagnosticInteractor) ListDiagnosticsByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorListStart, "motorcycle_id", motorcycleID)

	diagnostics, err := i.diagnosticService.GetDiagnosticsByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagnosticInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(diagnostics))
	return diagnostics, nil
}
