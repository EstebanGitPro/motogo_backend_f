package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var diagnosticLog logger.Logger = logger.NewSlogLogger()

// DiagnosticServiceImpl implements input.DiagnosticService
type DiagnosticServiceImpl struct {
	diagnosticRepo output.DiagnosticRepository
	motorcycleRepo output.MotorcycleRepository
	branchRepo     output.BranchRepository
}

// NewDiagnosticService creates a new DiagnosticService instance
func NewDiagnosticService(
	diagnosticRepo output.DiagnosticRepository,
	motorcycleRepo output.MotorcycleRepository,
	branchRepo output.BranchRepository,
) *DiagnosticServiceImpl {
	return &DiagnosticServiceImpl{
		diagnosticRepo: diagnosticRepo,
		motorcycleRepo: motorcycleRepo,
		branchRepo:     branchRepo,
	}
}

// BeginTx starts a new database transaction
func (s *DiagnosticServiceImpl) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.diagnosticRepo.BeginTx(ctx)
}

// ============================================
// Validations
// ============================================

// ValidateMotorcycleOwnership validates that the motorcycle exists and belongs to the given owner.
// Returns ErrMotorcycleNotFound if motorcycle doesn't exist or doesn't belong to the owner (security by obscurity).
func (s *DiagnosticServiceImpl) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	if motorcycle.OwnerID != ownerID {
		diagnosticLog.Warn(logger.LogDiagnosticInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	return motorcycle, nil
}

// ValidateBranchExists checks that the branch exists.
func (s *DiagnosticServiceImpl) ValidateBranchExists(ctx context.Context, branchID string) error {
	_, err := s.branchRepo.GetBranchByID(ctx, branchID)
	if err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorBranchError, "error", err, "branch_id", branchID)
		return domain.ErrBranchNotFound
	}
	return nil
}

// ============================================
// Diagnostic CRUD
// ============================================

// RegisterOrUpdateDiagnostic implements UPSERT logic:
// If a diagnostic already exists for the same motorcycle+branch, it updates it.
// Otherwise, it creates a new one.
func (s *DiagnosticServiceImpl) RegisterOrUpdateDiagnostic(ctx context.Context, tx output.Tx, motorcycleID, branchID string, problemDescription *string) (*domain.Diagnostic, error) {
	// Check if diagnostic already exists for this motorcycle+branch
	existing, existErr := s.diagnosticRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
	if existErr != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorMotoError, "error", existErr, "motorcycle_id", motorcycleID, "branch_id", branchID)
		return nil, domain.ErrDiagnosticCannotSave
	}

	if existing != nil {
		// === UPSERT: Update existing diagnostic ===
		diagnosticLog.Info(logger.LogDiagnosticInteractorExistingFound, "existing_id", existing.ID, "motorcycle_id", motorcycleID, "branch_id", branchID)

		// Refresh diagnostic fields
		RefreshDiagnostic(existing, problemDescription)

		// Update diagnostic record
		if err := s.diagnosticRepo.Update(ctx, tx, existing); err != nil {
			diagnosticLog.Error(logger.LogDiagnosticInteractorUpsertUpdateErr, "error", err, "id", existing.ID)
			return nil, domain.ErrDiagnosticCannotSave
		}

		diagnosticLog.Success(logger.LogDiagnosticInteractorUpsertSuccess, "id", existing.ID, "motorcycle_id", motorcycleID)
		return existing, nil
	}

	// === CREATE: New diagnostic ===
	diagnostic := NewDiagnostic(motorcycleID, branchID, problemDescription)
	diagnosticLog.Debug(logger.LogDiagnosticInteractorIDGenerated, "id", diagnostic.ID)

	// Save diagnostic
	if err := s.diagnosticRepo.Save(ctx, tx, diagnostic); err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorSaveError, "error", err)
		return nil, domain.ErrDiagnosticCannotSave
	}

	diagnosticLog.Success(logger.LogDiagnosticInteractorCreateSuccess, "id", diagnostic.ID, "motorcycle_id", motorcycleID)
	return diagnostic, nil
}

// GetByID retrieves a diagnostic by its ID
func (s *DiagnosticServiceImpl) GetByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	return s.diagnosticRepo.GetByID(ctx, diagnosticID)
}

// GetByMotorcycleID retrieves all diagnostics for a motorcycle
func (s *DiagnosticServiceImpl) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	return s.diagnosticRepo.GetByMotorcycleID(ctx, motorcycleID)
}

// ApplyDiagnosticUpdates merges partial updates into an existing diagnostic
func (s *DiagnosticServiceImpl) ApplyDiagnosticUpdates(existing, updates *domain.Diagnostic) {
	if updates.ProblemDescription != nil {
		existing.ProblemDescription = updates.ProblemDescription
	}
	if updates.PossibleSolution != nil {
		existing.PossibleSolution = updates.PossibleSolution
	}
}

// UpdateDiagnostic persists diagnostic changes
func (s *DiagnosticServiceImpl) UpdateDiagnostic(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	if err := s.diagnosticRepo.Update(ctx, tx, diagnostic); err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorUpdateError, "error", err, "diagnostic_id", diagnostic.ID)
		return domain.ErrDiagnosticCannotUpdate
	}
	return nil
}

// DeleteDiagnostic deletes a diagnostic (cascades to evidence via FK ON DELETE CASCADE)
func (s *DiagnosticServiceImpl) DeleteDiagnostic(ctx context.Context, tx output.Tx, diagnosticID string) error {
	if err := s.diagnosticRepo.Delete(ctx, tx, diagnosticID); err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorDeleteError, "error", err, "diagnostic_id", diagnosticID)
		return domain.ErrDiagnosticCannotDelete
	}
	return nil
}

// SetSolution sets the possible solution for a diagnostic
func (s *DiagnosticServiceImpl) SetSolution(ctx context.Context, tx output.Tx, diagnosticID, solution string) error {
	diagnostic, err := s.diagnosticRepo.GetByID(ctx, diagnosticID)
	if err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorSetSolutionError, "error", err, "diagnostic_id", diagnosticID)
		return domain.ErrDiagnosticNotFound
	}

	diagnostic.PossibleSolution = &solution

	if err := s.diagnosticRepo.Update(ctx, tx, diagnostic); err != nil {
		diagnosticLog.Error(logger.LogDiagnosticInteractorSetSolutionError, "error", err, "diagnostic_id", diagnosticID)
		return domain.ErrDiagnosticCannotUpdate
	}

	diagnosticLog.Success(logger.LogDiagnosticInteractorSetSolutionSuccess, "diagnostic_id", diagnosticID)
	return nil
}
