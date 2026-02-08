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
func (i *DiagnosticInteractor) RegisterDiagnostic(ctx context.Context, motorcycleID, branchID, ownerID string, problemDescription *string, evidenceURLs []string) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorCreateStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and ownership
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 3: Validate branch exists
	_, err = i.branchRepo.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBranchError, "error", err, "branch_id", branchID)
		return nil, domain.ErrBranchNotFound
	}

	// Step 4: Check if diagnostic already exists for this motorcycle+branch (UPSERT)
	existing, err := i.diagnosticRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID, "branch_id", branchID)
		return nil, domain.ErrDiagnosticCannotSave
	}

	// Step 5: Begin transaction
	tx, err := i.diagnosticRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	if existing != nil {
		// === UPSERT: Update existing diagnostic ===
		log.Info(logger.LogDiagnosticInteractorExistingFound, "existing_id", existing.ID, "motorcycle_id", motorcycleID, "branch_id", branchID)

		// Refresh diagnostic fields (business logic delegated to services layer)
		services.RefreshDiagnostic(existing, problemDescription)

		// Update diagnostic record
		err = i.diagnosticRepo.Update(ctx, tx, existing)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorUpsertUpdateErr, "error", err, "id", existing.ID)
			tx.Rollback()
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Delete old evidence
		err = i.diagnosticRepo.DeleteEvidenceByDiagnosticID(ctx, tx, existing.ID)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorEvidCleanupError, "error", err, "diagnostic_id", existing.ID)
			tx.Rollback()
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Save new evidence
		existing.Evidence = nil
		for _, url := range evidenceURLs {
			evidence := services.NewDiagnosticEvidence(existing.ID, url, nil)
			err = i.diagnosticRepo.SaveEvidence(ctx, tx, evidence)
			if err != nil {
				log.Error(logger.LogDiagnosticInteractorSaveEvidError, "error", err, "url", url)
				tx.Rollback()
				return nil, domain.ErrDiagnosticCannotSave
			}
			existing.Evidence = append(existing.Evidence, *evidence)
		}

		// Commit
		if err := tx.Commit(); err != nil {
			log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
			return nil, domain.ErrDiagnosticCannotSave
		}

		log.Success(logger.LogDiagnosticInteractorUpsertSuccess, "id", existing.ID, "motorcycle_id", motorcycleID)
		return existing, nil
	}

	// === CREATE: New diagnostic ===
	diagnostic := services.NewDiagnostic(motorcycleID, branchID, problemDescription)
	log.Debug(logger.LogDiagnosticInteractorIDGenerated, "id", diagnostic.ID)

	// Save diagnostic
	err = i.diagnosticRepo.Save(ctx, tx, diagnostic)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorSaveError, "error", err)
		tx.Rollback()
		return nil, domain.ErrDiagnosticCannotSave
	}

	// Save evidence photos
	for _, url := range evidenceURLs {
		evidence := services.NewDiagnosticEvidence(diagnostic.ID, url, nil)
		err = i.diagnosticRepo.SaveEvidence(ctx, tx, evidence)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorSaveEvidError, "error", err, "url", url)
			tx.Rollback()
			return nil, domain.ErrDiagnosticCannotSave
		}
		diagnostic.Evidence = append(diagnostic.Evidence, *evidence)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	log.Success(logger.LogDiagnosticInteractorCreateSuccess, "id", diagnostic.ID, "motorcycle_id", motorcycleID)
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
func (i *DiagnosticInteractor) UpdateDiagnostic(ctx context.Context, diagnosticID, ownerID string, updates *domain.Diagnostic) (*domain.Diagnostic, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorUpdateStart, "diagnostic_id", diagnosticID)

	// Step 1: Get existing diagnostic
	diagnostic, err := i.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, domain.ErrDiagnosticNotFound
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, diagnostic.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
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
	tx, err := i.diagnosticRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	// Step 5: Update diagnostic
	err = i.diagnosticRepo.Update(ctx, tx, diagnostic)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorUpdateError, "error", err, "diagnostic_id", diagnosticID)
		tx.Rollback()
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	// Step 6: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return nil, domain.ErrDiagnosticCannotUpdate
	}

	log.Success(logger.LogDiagnosticInteractorUpdateSuccess, "diagnostic_id", diagnosticID)
	return diagnostic, nil
}

// DeleteDiagnostic deletes a diagnostic and its evidence (HU13)
func (i *DiagnosticInteractor) DeleteDiagnostic(ctx context.Context, diagnosticID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagnosticInteractorDeleteStart, "diagnostic_id", diagnosticID)

	// Step 1: Get diagnostic
	diagnostic, err := i.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return err
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, diagnostic.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "diagnostic_id", diagnosticID, "owner_id", ownerID)
		return domain.ErrDiagnosticNotFound
	}

	// Step 3: Begin transaction
	tx, err := i.diagnosticRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBeginTxError, "error", err)
		return domain.ErrDiagnosticCannotDelete
	}

	// Step 4: Delete diagnostic (cascades to evidence via FK ON DELETE CASCADE)
	err = i.diagnosticRepo.Delete(ctx, tx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorDeleteError, "error", err, "diagnostic_id", diagnosticID)
		tx.Rollback()
		return domain.ErrDiagnosticCannotDelete
	}

	// Step 5: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDiagnosticInteractorCommitError, "error", err)
		return domain.ErrDiagnosticCannotDelete
	}

	log.Success(logger.LogDiagnosticInteractorDeleteSuccess, "diagnostic_id", diagnosticID)
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
