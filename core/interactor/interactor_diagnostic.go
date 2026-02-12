package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// DiagnosticInteractor handles diagnostic-related use cases (HU11-14)
type DiagnosticInteractor struct {
	diagnosticRepo output.DiagnosticRepository
	motorcycleRepo output.MotorcycleRepository
	branchRepo     output.BranchRepository
}

// NewDiagnosticInteractor creates a new DiagnosticInteractor instance
func NewDiagnosticInteractor(
	diagnosticRepo output.DiagnosticRepository,
	motorcycleRepo output.MotorcycleRepository,
	branchRepo output.BranchRepository,
) *DiagnosticInteractor {
	return &DiagnosticInteractor{
		diagnosticRepo: diagnosticRepo,
		motorcycleRepo: motorcycleRepo,
		branchRepo:     branchRepo,
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
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if motoErr != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", motoErr, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 3: Validate branch exists
	_, branchErr := i.branchRepo.GetBranchByID(ctx, branchID)
	if branchErr != nil {
		log.Error(logger.LogDiagnosticInteractorBranchError, "error", branchErr, "branch_id", branchID)
		return nil, domain.ErrBranchNotFound
	}

	// Step 4: Check if diagnostic already exists for this motorcycle+branch (UPSERT)
	existing, existErr := i.diagnosticRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
	if existErr != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", existErr, "motorcycle_id", motorcycleID, "branch_id", branchID)
		return nil, domain.ErrDiagnosticCannotSave
	}

	// Step 5: Begin transaction
	tx, txErr := i.diagnosticRepo.BeginTx(ctx)
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

	if existing != nil {
		// === UPSERT: Update existing diagnostic ===
		log.Info(logger.LogDiagnosticInteractorExistingFound, "existing_id", existing.ID, "motorcycle_id", motorcycleID, "branch_id", branchID)

		// Refresh diagnostic fields (business logic delegated to services layer)
		services.RefreshDiagnostic(existing, problemDescription)

		// Update diagnostic record
		if err = i.diagnosticRepo.Update(ctx, tx, existing); err != nil {
			log.Error(logger.LogDiagnosticInteractorUpsertUpdateErr, "error", err, "id", existing.ID)
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Delete old evidence
		if err = i.diagnosticRepo.DeleteEvidenceByDiagnosticID(ctx, tx, existing.ID); err != nil {
			log.Error(logger.LogDiagnosticInteractorEvidCleanupError, "error", err, "diagnostic_id", existing.ID)
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Save new evidence
		existing.Evidence = nil
		for _, url := range evidenceURLs {
			evidence := services.NewDiagnosticEvidence(existing.ID, url, nil)
			if err = i.diagnosticRepo.SaveEvidence(ctx, tx, evidence); err != nil {
				log.Error(logger.LogDiagnosticInteractorSaveEvidError, "error", err, "url", url)
				return nil, domain.ErrDiagnosticCannotSave
			}
			existing.Evidence = append(existing.Evidence, *evidence)
		}

		// Commit
		if err = tx.Commit(); err != nil {
			log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
			return nil, domain.ErrDiagnosticCannotSave
		}

		log.Success(logger.LogDiagnosticInteractorUpsertSuccess, "id", existing.ID, "motorcycle_id", motorcycleID)

		err = nil
		return existing, nil
	}

	// === CREATE: New diagnostic ===
	diagnostic := services.NewDiagnostic(motorcycleID, branchID, problemDescription)
	log.Debug(logger.LogDiagnosticInteractorIDGenerated, "id", diagnostic.ID)

	// Save diagnostic
	if err = i.diagnosticRepo.Save(ctx, tx, diagnostic); err != nil {
		log.Error(logger.LogDiagnosticInteractorSaveError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	// Save evidence photos
	for _, url := range evidenceURLs {
		evidence := services.NewDiagnosticEvidence(diagnostic.ID, url, nil)
		if err = i.diagnosticRepo.SaveEvidence(ctx, tx, evidence); err != nil {
			log.Error(logger.LogDiagnosticInteractorSaveEvidError, "error", err, "url", url)
			return nil, domain.ErrDiagnosticCannotSave
		}
		diagnostic.Evidence = append(diagnostic.Evidence, *evidence)
	}

	// Commit
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

	// Step 1: Get diagnostic
	diagnostic, err := i.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, err
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, diagnostic.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 3: Load evidence
	evidence, err := i.diagnosticRepo.GetEvidenceByDiagnosticID(ctx, diagnosticID)
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
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Get all diagnostics
	diagnostics, err := i.diagnosticRepo.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 3: Load evidence for each diagnostic
	for idx := range diagnostics {
		evidence, err := i.diagnosticRepo.GetEvidenceByDiagnosticID(ctx, diagnostics[idx].ID)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorListError, "error loading evidence", err, "diagnostic_id", diagnostics[idx].ID)
			return nil, err
		}
		diagnostics[idx].Evidence = evidence
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
	diagnostic, getErr := i.diagnosticRepo.GetByID(ctx, diagnosticID)
	if getErr != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", getErr, "diagnostic_id", diagnosticID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, diagnostic.MotorcycleID)
	if motoErr != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 3: Apply updates
	if updates.ProblemDescription != nil {
		diagnostic.ProblemDescription = updates.ProblemDescription
	}
	if updates.PossibleSolution != nil {
		diagnostic.PossibleSolution = updates.PossibleSolution
	}
	if updates.LaborQuote != nil {
		diagnostic.LaborQuote = updates.LaborQuote
	}
	if updates.PartsQuote != nil {
		diagnostic.PartsQuote = updates.PartsQuote
	}
	if updates.EstimatedTime != nil {
		diagnostic.EstimatedTime = updates.EstimatedTime
	}

	// Step 4: Begin transaction
	tx, txErr := i.diagnosticRepo.BeginTx(ctx)
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
	if err = i.diagnosticRepo.Update(ctx, tx, diagnostic); err != nil {
		log.Error(logger.LogDiagnosticInteractorUpdateError, "error", err, "diagnostic_id", diagnosticID)
		return nil, domain.ErrDiagnosticCannotUpdate
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
	diagnostic, getErr := i.diagnosticRepo.GetByID(ctx, diagnosticID)
	if getErr != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", getErr, "diagnostic_id", diagnosticID)
		return getErr
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, diagnostic.MotorcycleID)
	if motoErr != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return domain.ErrDiagnosticNotFound
	}

	// Step 3: Begin transaction
	tx, txErr := i.diagnosticRepo.BeginTx(ctx)
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
	if err = i.diagnosticRepo.Delete(ctx, tx, diagnosticID); err != nil {
		log.Error(logger.LogDiagnosticInteractorDeleteError, "error", err, "diagnostic_id", diagnosticID)
		return domain.ErrDiagnosticCannotDelete
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
	diagnostics, err := i.diagnosticRepo.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 2: Load evidence for each diagnostic
	for idx := range diagnostics {
		evidence, err := i.diagnosticRepo.GetEvidenceByDiagnosticID(ctx, diagnostics[idx].ID)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorListError, "error loading evidence", err, "diagnostic_id", diagnostics[idx].ID)
			return nil, err
		}
		diagnostics[idx].Evidence = evidence
	}

	log.Success(logger.LogDiagnosticInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(diagnostics))
	return diagnostics, nil
}
