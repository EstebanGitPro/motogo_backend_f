package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/google/uuid"
)

// motorcycleService implements input.MotorcycleService
type motorcycleService struct {
	motorcycleRepo output.MotorcycleRepository
	diagPermRepo   output.DiagnosticPermissionRepository
	storageClient  output.StorageClient
}

// NewMotorcycleService creates a new MotorcycleService instance
func NewMotorcycleService(motorcycleRepo output.MotorcycleRepository, diagPermRepo output.DiagnosticPermissionRepository) *motorcycleService {
	return &motorcycleService{
		motorcycleRepo: motorcycleRepo,
		diagPermRepo:   diagPermRepo,
	}
}

// WithStorageClient sets the optional storage client for image deletion
func (s *motorcycleService) WithStorageClient(client output.StorageClient) {
	s.storageClient = client
}

// BeginTx starts a new transaction for motorcycle operations
func (s *motorcycleService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.motorcycleRepo.BeginTx(ctx)
}

// BeginPermissionTx starts a new transaction for diagnostic permission operations
func (s *motorcycleService) BeginPermissionTx(ctx context.Context) (output.Tx, error) {
	return s.diagPermRepo.BeginTx(ctx)
}

// ValidateMotorcycleOwnership validates motorcycle exists and belongs to owner
// Returns the motorcycle if valid, error otherwise
func (s *motorcycleService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		return nil, err
	}
	if motorcycle.OwnerID != ownerID {
		return nil, domain.ErrMotorcycleNotFound
	}
	return motorcycle, nil
}

// ValidateReferenceExists validates that a reference_id exists in the catalog
func (s *motorcycleService) ValidateReferenceExists(ctx context.Context, referenceID string) error {
	if referenceID == "" {
		return domain.ErrReferenceRequired
	}
	refExists, err := s.motorcycleRepo.ValidateReferenceExists(ctx, referenceID)
	if err != nil {
		return domain.ErrMotorcycleCannotSave
	}
	if !refExists {
		return domain.ErrReferenceNotFound
	}
	return nil
}

// ValidateLicensePlateUnique validates license plate is not already registered
func (s *motorcycleService) ValidateLicensePlateUnique(ctx context.Context, licensePlate string) error {
	plateExists, err := s.motorcycleRepo.CheckLicensePlateExists(ctx, licensePlate)
	if err != nil {
		return domain.ErrMotorcycleCannotSave
	}
	if plateExists {
		return domain.ErrDuplicateLicensePlate
	}
	return nil
}

// CreateMotorcycle generates UUID and saves motorcycle to database
func (s *motorcycleService) CreateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	motorcycle.ID = uuid.New().String()

	if err := s.motorcycleRepo.Save(ctx, tx, motorcycle); err != nil {
		return nil, domain.ErrMotorcycleCannotSave
	}

	return motorcycle, nil
}

// GetMotorcycleByID retrieves a motorcycle by ID
func (s *motorcycleService) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByID(ctx, motorcycleID)
}

// GetMotorcyclesByOwner retrieves all motorcycles for an owner
func (s *motorcycleService) GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByOwnerID(ctx, ownerID)
}

// GetMotorcycleByLicensePlate retrieves a motorcycle by license plate
func (s *motorcycleService) GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByLicensePlate(ctx, licensePlate)
}

// ApplyMotorcycleUpdates applies field patches to motorcycle entity
// Validates reference_id if changed, returns error if reference validation fails
func (s *motorcycleService) ApplyMotorcycleUpdates(ctx context.Context, motorcycle *domain.Motorcycle, updates *domain.Motorcycle) error {
	// Validate new reference_id if changed
	if updates.ReferenceID != "" && updates.ReferenceID != motorcycle.ReferenceID {
		refExists, err := s.motorcycleRepo.ValidateReferenceExists(ctx, updates.ReferenceID)
		if err != nil {
			return domain.ErrMotorcycleCannotUpdate
		}
		if !refExists {
			return domain.ErrReferenceNotFound
		}
		motorcycle.ReferenceID = updates.ReferenceID
	}

	// Apply optional field patches
	if updates.Year != nil {
		motorcycle.Year = updates.Year
	}
	if updates.CurrentMileage != nil {
		motorcycle.CurrentMileage = updates.CurrentMileage
	}
	if updates.OwnerNotes != nil {
		motorcycle.OwnerNotes = updates.OwnerNotes
	}
	if updates.ProfileImageURL != nil {
		motorcycle.ProfileImageURL = updates.ProfileImageURL
	}

	return nil
}

// UpdateMotorcycle persists motorcycle updates via transaction
func (s *motorcycleService) UpdateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	return s.motorcycleRepo.Update(ctx, tx, motorcycle)
}

// CheckServiceHistory checks if motorcycle has diagnostic/service history
func (s *motorcycleService) CheckServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	return s.motorcycleRepo.HasServiceHistory(ctx, motorcycleID)
}

// DeleteMotorcycle performs soft or hard delete based on history flag
func (s *motorcycleService) DeleteMotorcycle(ctx context.Context, tx output.Tx, motorcycleID string, hasHistory bool) error {
	if hasHistory {
		return s.motorcycleRepo.Delete(ctx, tx, motorcycleID)
	}
	return s.motorcycleRepo.HardDelete(ctx, tx, motorcycleID)
}

// DeleteProfileImage clears profile image URL in database
func (s *motorcycleService) DeleteProfileImage(ctx context.Context, tx output.Tx, motorcycleID string) error {
	return s.motorcycleRepo.ClearProfileImageURL(ctx, tx, motorcycleID)
}

// DeleteStorageFile deletes a file from Firebase Storage (best effort, no error propagation)
func (s *motorcycleService) DeleteStorageFile(ctx context.Context, imageURL string) {
	if imageURL == "" || s.storageClient == nil {
		return
	}
	// Best effort - log warning but don't fail
	_ = s.storageClient.DeleteStorageFile(ctx, imageURL)
}

// GetAllReferences retrieves all motorcycle references from catalog
func (s *motorcycleService) GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	return s.motorcycleRepo.GetAllReferences(ctx)
}

// GetReferencesByBrandID retrieves references for a specific brand
func (s *motorcycleService) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	return s.motorcycleRepo.GetReferencesByBrandID(ctx, brandID)
}

// GrantDiagnosticPermission creates or updates a diagnostic permission with the given active state
func (s *motorcycleService) GrantDiagnosticPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string, active bool) (*domain.DiagnosticPermission, error) {
	permission := &domain.DiagnosticPermission{
		ID:           utils.Generate(),
		MotorcycleID: motorcycleID,
		BranchID:     branchID,
		Active:       active,
	}

	if err := s.diagPermRepo.Save(ctx, tx, permission); err != nil {
		return nil, domain.ErrPermissionCannotSave
	}

	return permission, nil
}

// RevokeDiagnosticPermission deactivates a diagnostic permission (sets active = false)
func (s *motorcycleService) RevokeDiagnosticPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	return s.diagPermRepo.Deactivate(ctx, tx, motorcycleID, branchID)
}

// ListDiagnosticPermissions retrieves all active permissions for a motorcycle
func (s *motorcycleService) ListDiagnosticPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	return s.diagPermRepo.GetByMotorcycleID(ctx, motorcycleID)
}

// ValidateBranchPermission checks if any of the given branches has an active diagnostic permission for the motorcycle
// Returns nil if at least one branch is authorized, ErrBranchNotAuthorized otherwise
func (s *motorcycleService) ValidateBranchPermission(ctx context.Context, motorcycleID string, branchIDs []string) error {
	for _, branchID := range branchIDs {
		_, err := s.diagPermRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
		if err == nil {
			return nil // At least one branch is authorized
		}
	}
	return domain.ErrBranchNotAuthorized
}
