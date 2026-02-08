package services

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// diagnosticService implements input.DiagnosticService
type diagnosticService struct {
	diagnosticRepo output.DiagnosticRepository
	motorcycleRepo output.MotorcycleRepository
	branchRepo     output.BranchRepository
}

// NewDiagnosticService creates a new DiagnosticService instance
func NewDiagnosticService(
	diagnosticRepo output.DiagnosticRepository,
	motorcycleRepo output.MotorcycleRepository,
	branchRepo output.BranchRepository,
) input.DiagnosticService {
	return &diagnosticService{
		diagnosticRepo: diagnosticRepo,
		motorcycleRepo: motorcycleRepo,
		branchRepo:     branchRepo,
	}
}

// BeginTx starts a new database transaction
func (s *diagnosticService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.diagnosticRepo.BeginTx(ctx)
}

// ValidateMotorcycleOwnership validates that the motorcycle exists and belongs to the given owner
// Returns ErrMotorcycleNotFound if motorcycle doesn't exist or doesn't belong to owner (security by obscurity)
func (s *diagnosticService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) error {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleNotFound
	}
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagnosticInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return domain.ErrMotorcycleNotFound
	}
	return nil
}

// ValidateBranchExists validates that the branch exists
func (s *diagnosticService) ValidateBranchExists(ctx context.Context, branchID string) error {
	_, err := s.branchRepo.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorBranchError, "error", err, "branch_id", branchID)
		return domain.ErrBranchNotFound
	}
	return nil
}

// UpsertDiagnostic creates a new diagnostic or updates an existing one for the same motorcycle+branch
// This encapsulates the entire UPSERT business logic: check existing, refresh/create, replace evidence
func (s *diagnosticService) UpsertDiagnostic(ctx context.Context, tx output.Tx, motorcycleID, branchID string, problemDescription *string, evidenceURLs []string) (*domain.Diagnostic, error) {
	// Check if diagnostic already exists for this motorcycle+branch
	existing, err := s.diagnosticRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID, "branch_id", branchID)
		return nil, domain.ErrDiagnosticCannotSave
	}

	if existing != nil {
		// === UPSERT: Update existing diagnostic ===
		log.Info(logger.LogDiagnosticInteractorExistingFound, "existing_id", existing.ID, "motorcycle_id", motorcycleID, "branch_id", branchID)

		// Refresh diagnostic fields
		existing.ProblemDescription = problemDescription
		existing.Date = time.Now()
		existing.SentViaWhatsApp = false

		// Update diagnostic record
		if err := s.diagnosticRepo.Update(ctx, tx, existing); err != nil {
			log.Error(logger.LogDiagnosticInteractorUpsertUpdateErr, "error", err, "id", existing.ID)
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Delete old evidence
		if err := s.diagnosticRepo.DeleteEvidenceByDiagnosticID(ctx, tx, existing.ID); err != nil {
			log.Error(logger.LogDiagnosticInteractorEvidCleanupError, "error", err, "diagnostic_id", existing.ID)
			return nil, domain.ErrDiagnosticCannotSave
		}

		// Save new evidence
		existing.Evidence = nil
		if err := s.saveEvidenceList(ctx, tx, existing.ID, evidenceURLs, &existing.Evidence); err != nil {
			return nil, err
		}

		log.Success(logger.LogDiagnosticInteractorUpsertSuccess, "id", existing.ID, "motorcycle_id", motorcycleID)
		return existing, nil
	}

	// === CREATE: New diagnostic ===
	diagnostic := &domain.Diagnostic{
		ID:                 utils.Generate(),
		MotorcycleID:       motorcycleID,
		BranchID:           branchID,
		Date:               time.Now(),
		ProblemDescription: problemDescription,
		SentViaWhatsApp:    false,
	}
	log.Debug(logger.LogDiagnosticInteractorIDGenerated, "id", diagnostic.ID)

	if err := s.diagnosticRepo.Save(ctx, tx, diagnostic); err != nil {
		log.Error(logger.LogDiagnosticInteractorSaveError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	// Save evidence photos
	if err := s.saveEvidenceList(ctx, tx, diagnostic.ID, evidenceURLs, &diagnostic.Evidence); err != nil {
		return nil, err
	}

	log.Success(logger.LogDiagnosticInteractorCreateSuccess, "id", diagnostic.ID, "motorcycle_id", motorcycleID)
	return diagnostic, nil
}

// GetDiagnosticByID retrieves a diagnostic by ID and enriches it with evidence
func (s *diagnosticService) GetDiagnosticByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	diagnostic, err := s.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error", err, "diagnostic_id", diagnosticID)
		return nil, err
	}

	// Enrich with evidence
	evidence, err := s.diagnosticRepo.GetEvidenceByDiagnosticID(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorGetError, "error loading evidence", err, "diagnostic_id", diagnosticID)
		return nil, err
	}
	diagnostic.Evidence = evidence

	return diagnostic, nil
}

// GetDiagnosticsByMotorcycleID retrieves all diagnostics for a motorcycle, enriched with evidence
func (s *diagnosticService) GetDiagnosticsByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	diagnostics, err := s.diagnosticRepo.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Enrich each diagnostic with evidence
	for idx := range diagnostics {
		evidence, err := s.diagnosticRepo.GetEvidenceByDiagnosticID(ctx, diagnostics[idx].ID)
		if err != nil {
			log.Error(logger.LogDiagnosticInteractorListError, "error loading evidence", err, "diagnostic_id", diagnostics[idx].ID)
			return nil, err
		}
		diagnostics[idx].Evidence = evidence
	}

	return diagnostics, nil
}

// ApplyDiagnosticUpdates applies partial updates to an existing diagnostic (field-by-field patching)
func (s *diagnosticService) ApplyDiagnosticUpdates(existing *domain.Diagnostic, updates *domain.Diagnostic) {
	if updates.ProblemDescription != nil {
		existing.ProblemDescription = updates.ProblemDescription
	}
	if updates.PossibleSolution != nil {
		existing.PossibleSolution = updates.PossibleSolution
	}
	if updates.LaborQuote != nil {
		existing.LaborQuote = updates.LaborQuote
	}
	if updates.PartsQuote != nil {
		existing.PartsQuote = updates.PartsQuote
	}
	if updates.EstimatedTime != nil {
		existing.EstimatedTime = updates.EstimatedTime
	}
}

// UpdateDiagnostic persists a diagnostic update to the database
func (s *diagnosticService) UpdateDiagnostic(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	if err := s.diagnosticRepo.Update(ctx, tx, diagnostic); err != nil {
		log.Error(logger.LogDiagnosticInteractorUpdateError, "error", err, "diagnostic_id", diagnostic.ID)
		return domain.ErrDiagnosticCannotUpdate
	}
	return nil
}

// DeleteDiagnostic deletes a diagnostic (cascades to evidence via FK ON DELETE CASCADE)
func (s *diagnosticService) DeleteDiagnostic(ctx context.Context, tx output.Tx, diagnosticID string) error {
	if err := s.diagnosticRepo.Delete(ctx, tx, diagnosticID); err != nil {
		log.Error(logger.LogDiagnosticInteractorDeleteError, "error", err, "diagnostic_id", diagnosticID)
		return domain.ErrDiagnosticCannotDelete
	}
	return nil
}

// saveEvidenceList saves a list of evidence URLs and appends them to the evidence slice
func (s *diagnosticService) saveEvidenceList(ctx context.Context, tx output.Tx, diagnosticID string, evidenceURLs []string, evidence *[]domain.DiagnosticEvidence) error {
	for _, url := range evidenceURLs {
		ev := &domain.DiagnosticEvidence{
			ID:           utils.Generate(),
			DiagnosticID: diagnosticID,
			ImageURL:     url,
			CreatedAt:    time.Now(),
		}
		if err := s.diagnosticRepo.SaveEvidence(ctx, tx, ev); err != nil {
			log.Error(logger.LogDiagnosticInteractorSaveEvidError, "error", err, "url", url)
			return domain.ErrDiagnosticCannotSave
		}
		*evidence = append(*evidence, *ev)
	}
	return nil
}
