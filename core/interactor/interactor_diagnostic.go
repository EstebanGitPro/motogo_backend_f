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
	diagnosticSvc input.DiagnosticService
}

// NewDiagnosticInteractor creates a new DiagnosticInteractor instance
func NewDiagnosticInteractor(diagnosticSvc input.DiagnosticService) *DiagnosticInteractor {
	return &DiagnosticInteractor{
		diagnosticSvc: diagnosticSvc,
	}
}

// RegisterDiagnostic creates a new diagnostic or updates an existing one for the same motorcycle+branch (UPSERT)
// If a diagnostic already exists for the same motorcycle and branch, it updates the existing record
// instead of creating a duplicate. This is transparent to the caller.
func (i *DiagnosticInteractor) RegisterDiagnostic(ctx context.Context, motorcycleID, branchID, ownerID string, problemDescription *string, evidenceURLs []string) (result *domain.Diagnostic, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorCreateStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and ownership
	if _, err = i.diagnosticSvc.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// Step 2: Validate branch exists
	if err = i.diagnosticSvc.ValidateBranchExists(ctx, branchID); err != nil {
		return nil, err
	}

	// Step 3: Begin transaction
	tx, txErr := i.diagnosticSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrDiagnosticCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagnosticInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Register or update diagnostic (UPSERT + evidence)
	result, err = i.diagnosticSvc.RegisterOrUpdateDiagnostic(ctx, tx, motorcycleID, branchID, problemDescription, evidenceURLs)
	if err != nil {
		return nil, err
	}

	// Step 5: Commit
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	log.Success(logger.LogDiagnosticInteractorCreateSuccess, "id", result.ID, "motorcycle_id", motorcycleID)

	err = nil
	return result, nil
}

// GetDiagnosticByID retrieves a diagnostic with its evidence (HU14)
func (i *DiagnosticInteractor) GetDiagnosticByID(ctx context.Context, diagnosticID, ownerID string) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorGetStart, "diagnostic_id", diagnosticID)

	// Step 1: Get diagnostic
	diagnostic, err := i.diagnosticSvc.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, err
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.diagnosticSvc.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 3: Load evidence
	evidence, err := i.diagnosticSvc.LoadEvidence(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error loading evidence", err, "diagnostic_id", diagnosticID)
		return nil, err
	}
	diagnostic.Evidence = evidence

	log.Success(logger.LogDiagnosticInteractorGetSuccess, "diagnostic_id", diagnosticID)
	return diagnostic, nil
}

// ListDiagnosticsByMotorcycle retrieves all diagnostics for a motorcycle with evidence (HU14)
func (i *DiagnosticInteractor) ListDiagnosticsByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorListStart, "motorcycle_id", motorcycleID)

	// Step 1: Validate motorcycle exists and ownership
	if _, err := i.diagnosticSvc.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// Step 2: Get all diagnostics
	diagnostics, err := i.diagnosticSvc.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 3: Load evidence for each diagnostic
	if err = i.diagnosticSvc.LoadEvidenceForDiagnostics(ctx, diagnostics); err != nil {
		return nil, err
	}

	log.Success(logger.LogDiagnosticInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(diagnostics))
	return diagnostics, nil
}

// UpdateDiagnostic updates an existing diagnostic (HU12)
func (i *DiagnosticInteractor) UpdateDiagnostic(ctx context.Context, diagnosticID, ownerID string, updates *domain.Diagnostic) (result *domain.Diagnostic, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorUpdateStart, "diagnostic_id", diagnosticID)

	// Step 1: Get existing diagnostic
	diagnostic, getErr := i.diagnosticSvc.GetByID(ctx, diagnosticID)
	if getErr != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", getErr, "diagnostic_id", diagnosticID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.diagnosticSvc.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 3: Apply updates
	i.diagnosticSvc.ApplyDiagnosticUpdates(diagnostic, updates)

	// Step 4: Begin transaction
	tx, txErr := i.diagnosticSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagnosticInteractorRollbackOK)
			}
		}
	}()

	// Step 5: Update diagnostic
	if err = i.diagnosticSvc.UpdateDiagnostic(ctx, tx, diagnostic); err != nil {
		return nil, err
	}

	// Step 6: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	log.Success(logger.LogDiagnosticInteractorUpdateSuccess, "diagnostic_id", diagnosticID)

	err = nil
	return diagnostic, nil
}

// DeleteDiagnostic deletes a diagnostic and its evidence (HU13)
func (i *DiagnosticInteractor) DeleteDiagnostic(ctx context.Context, diagnosticID, ownerID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorDeleteStart, "diagnostic_id", diagnosticID)

	// Step 1: Get diagnostic
	diagnostic, getErr := i.diagnosticSvc.GetByID(ctx, diagnosticID)
	if getErr != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", getErr, "diagnostic_id", diagnosticID)
		return getErr
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.diagnosticSvc.ValidateMotorcycleOwnership(ctx, diagnostic.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return domain.ErrDiagnosticNotFound
	}

	// Step 3: Begin transaction
	tx, txErr := i.diagnosticSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", txErr)
		return domain.ErrDiagnosticCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagnosticInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Delete diagnostic (cascades to evidence via FK ON DELETE CASCADE)
	if err = i.diagnosticSvc.DeleteDiagnostic(ctx, tx, diagnosticID); err != nil {
		return err
	}

	// Step 5: Commit transaction
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

	// Step 1: Get diagnostics
	diagnostics, err := i.diagnosticSvc.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 2: Load evidence for each diagnostic
	if err = i.diagnosticSvc.LoadEvidenceForDiagnostics(ctx, diagnostics); err != nil {
		return nil, err
	}

	log.Success(logger.LogDiagnosticInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(diagnostics))
	return diagnostics, nil
}

// SetSolution sets the possible solution for a diagnostic (representative use - no ownership check)
// Used by PATCH /diagnostics/:id/solution
func (i *DiagnosticInteractor) SetSolution(ctx context.Context, diagnosticID string, solution string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorSetSolutionStart, "diagnostic_id", diagnosticID)

	// Step 1: Begin transaction
	tx, txErr := i.diagnosticSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", txErr)
		return domain.ErrDiagnosticCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagnosticInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagnosticInteractorRollbackOK)
			}
		}
	}()

	// Step 2: Set solution (service handles fetch + update)
	if err = i.diagnosticSvc.SetSolution(ctx, tx, diagnosticID, solution); err != nil {
		return err
	}

	// Step 3: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return domain.ErrDiagnosticCannotUpdate
	}

	log.Success(logger.LogDiagnosticInteractorSetSolutionSuccess, "diagnostic_id", diagnosticID)

	err = nil
	return nil
}
