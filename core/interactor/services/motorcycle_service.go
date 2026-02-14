package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var motorcycleLog logger.Logger = logger.NewSlogLogger()

// MotorcycleServiceImpl implements input.MotorcycleService
type MotorcycleServiceImpl struct {
	motorcycleRepo output.MotorcycleRepository
	diagPermRepo   output.DiagnosticPermissionRepository
	storageClient  output.StorageClient // Optional: Firebase Storage for image deletion
}

// NewMotorcycleService creates a new MotorcycleService instance
func NewMotorcycleService(
	motorcycleRepo output.MotorcycleRepository,
	diagPermRepo output.DiagnosticPermissionRepository,
) *MotorcycleServiceImpl {
	return &MotorcycleServiceImpl{
		motorcycleRepo: motorcycleRepo,
		diagPermRepo:   diagPermRepo,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (s *MotorcycleServiceImpl) WithStorageClient(client output.StorageClient) {
	s.storageClient = client
}

// BeginTx starts a new database transaction
func (s *MotorcycleServiceImpl) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.motorcycleRepo.BeginTx(ctx)
}

// ============================================
// Validation
// ============================================

// ValidateReferenceExists checks that the given reference_id exists in the catalog
func (s *MotorcycleServiceImpl) ValidateReferenceExists(ctx context.Context, referenceID string) error {
	if referenceID == "" {
		return domain.ErrReferenceRequired
	}
	exists, err := s.motorcycleRepo.ValidateReferenceExists(ctx, referenceID)
	if err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceRefError, "error", err, "reference_id", referenceID)
		return domain.ErrMotorcycleCannotSave
	}
	if !exists {
		return domain.ErrReferenceNotFound
	}
	return nil
}

// CheckLicensePlateUnique validates that the license plate is not already registered
func (s *MotorcycleServiceImpl) CheckLicensePlateUnique(ctx context.Context, licensePlate string) error {
	exists, err := s.motorcycleRepo.CheckLicensePlateExists(ctx, licensePlate)
	if err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceCheckPlateErr, "error", err, "license_plate", licensePlate)
		return domain.ErrMotorcycleCannotSave
	}
	if exists {
		return domain.ErrDuplicateLicensePlate
	}
	return nil
}

// ValidateOwnership retrieves a motorcycle and validates ownership.
// Returns ErrMotorcycleNotFound if motorcycle doesn't exist or doesn't belong to the owner (security by obscurity)
func (s *MotorcycleServiceImpl) ValidateOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		return nil, err
	}
	if motorcycle.OwnerID != ownerID {
		return nil, domain.ErrMotorcycleNotFound
	}
	return motorcycle, nil
}

// ============================================
// Motorcycle CRUD
// ============================================

// RegisterMotorcycle sets ID and saves a new motorcycle (HU43)
func (s *MotorcycleServiceImpl) RegisterMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	motorcycle.SetID()
	if err := s.motorcycleRepo.Save(ctx, tx, motorcycle); err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceSaveError, "error", err)
		return domain.ErrMotorcycleCannotSave
	}
	return nil
}

// GetByID retrieves a motorcycle by its ID (HU46)
func (s *MotorcycleServiceImpl) GetByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByID(ctx, motorcycleID)
}

// GetByOwnerID retrieves all motorcycles for an owner (HU47)
func (s *MotorcycleServiceImpl) GetByOwnerID(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByOwnerID(ctx, ownerID)
}

// GetByLicensePlate retrieves a motorcycle by license plate
func (s *MotorcycleServiceImpl) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	return s.motorcycleRepo.GetByLicensePlate(ctx, licensePlate)
}

// ApplyUpdates merges partial updates into an existing motorcycle (HU44)
// Validates new reference_id if changed
func (s *MotorcycleServiceImpl) ApplyUpdates(existing *domain.Motorcycle, updates *domain.Motorcycle) error {
	// Validate new reference_id if changed
	if updates.ReferenceID != "" && updates.ReferenceID != existing.ReferenceID {
		exists, err := s.motorcycleRepo.ValidateReferenceExists(context.Background(), updates.ReferenceID)
		if err != nil {
			motorcycleLog.Error(logger.LogMotorcycleServiceRefError, "error", err, "reference_id", updates.ReferenceID)
			return domain.ErrMotorcycleCannotUpdate
		}
		if !exists {
			return domain.ErrReferenceNotFound
		}
		existing.ReferenceID = updates.ReferenceID
	}

	// Apply optional fields
	if updates.Year != nil {
		existing.Year = updates.Year
	}
	if updates.CurrentMileage != nil {
		existing.CurrentMileage = updates.CurrentMileage
	}
	if updates.OwnerNotes != nil {
		existing.OwnerNotes = updates.OwnerNotes
	}
	if updates.ProfileImageURL != nil {
		existing.ProfileImageURL = updates.ProfileImageURL
	}

	return nil
}

// UpdateMotorcycle persists motorcycle changes (HU44)
func (s *MotorcycleServiceImpl) UpdateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	if err := s.motorcycleRepo.Update(ctx, tx, motorcycle); err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceUpdateError, "error", err, "motorcycle_id", motorcycle.ID)
		return domain.ErrMotorcycleCannotUpdate
	}
	return nil
}

// DeleteMotorcycle implements hybrid delete strategy (HU45):
// If motorcycle has service history → soft delete; otherwise → hard delete
func (s *MotorcycleServiceImpl) DeleteMotorcycle(ctx context.Context, tx output.Tx, motorcycleID string) error {
	// Determine strategy
	hasHistory, err := s.motorcycleRepo.HasServiceHistory(ctx, motorcycleID)
	if err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceDeleteError, "error checking history", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotDelete
	}

	if hasHistory {
		// Soft delete - preserves historical data integrity
		motorcycleLog.Info(logger.LogMotorcycleServiceDeleteStart, "strategy", "soft_delete", "motorcycle_id", motorcycleID)
		if err = s.motorcycleRepo.Delete(ctx, tx, motorcycleID); err != nil {
			motorcycleLog.Error(logger.LogMotorcycleServiceDeleteError, "error", err, "motorcycle_id", motorcycleID)
			return domain.ErrMotorcycleCannotDelete
		}
	} else {
		// Hard delete - no history to preserve, clean removal
		motorcycleLog.Info(logger.LogMotorcycleServiceDeleteStart, "strategy", "hard_delete", "motorcycle_id", motorcycleID)
		if err = s.motorcycleRepo.HardDelete(ctx, tx, motorcycleID); err != nil {
			motorcycleLog.Error(logger.LogMotorcycleServiceDeleteError, "error", err, "motorcycle_id", motorcycleID)
			return domain.ErrMotorcycleCannotDelete
		}
	}

	return nil
}

// ClearProfileImageURL clears the profile image URL in database (HU39)
func (s *MotorcycleServiceImpl) ClearProfileImageURL(ctx context.Context, tx output.Tx, motorcycleID string) error {
	if err := s.motorcycleRepo.ClearProfileImageURL(ctx, tx, motorcycleID); err != nil {
		motorcycleLog.Error(logger.LogMotorcycleServiceUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotUpdate
	}
	return nil
}

// HasServiceHistory checks if motorcycle has any service history (HU45 hybrid delete)
func (s *MotorcycleServiceImpl) HasServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	return s.motorcycleRepo.HasServiceHistory(ctx, motorcycleID)
}

// ============================================
// Storage cleanup
// ============================================

// DeleteStorageFile deletes a file from cloud storage (best-effort, never fails)
func (s *MotorcycleServiceImpl) DeleteStorageFile(ctx context.Context, url string) {
	if s.storageClient == nil || url == "" {
		return
	}
	if err := s.storageClient.DeleteStorageFile(ctx, url); err != nil {
		motorcycleLog.Warn(logger.LogMotorcycleServiceStorageDeleteErr, "error", err, "url", url)
	}
}

// ============================================
// Reference catalog
// ============================================

// GetAllReferences retrieves all motorcycle references from catalog (HU50)
func (s *MotorcycleServiceImpl) GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	return s.motorcycleRepo.GetAllReferences(ctx)
}

// GetReferencesByBrandID retrieves motorcycle references for a specific brand (HU40)
func (s *MotorcycleServiceImpl) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	return s.motorcycleRepo.GetReferencesByBrandID(ctx, brandID)
}

// GetDistinctCategories retrieves all distinct motorcycle categories with line counts (HU41)
func (s *MotorcycleServiceImpl) GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error) {
	return s.motorcycleRepo.GetDistinctCategories(ctx)
}

// GetLinesByCategory retrieves all motorcycle lines for a specific category (HU41)
func (s *MotorcycleServiceImpl) GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error) {
	return s.motorcycleRepo.GetLinesByCategory(ctx, categoryName)
}

// GetDistinctDisplacements returns the fixed ENUM displacement range values (HU49)
// Values are hardcoded since they are fixed ENUM: BAJO (50-200cc), MEDIO (201-400cc), ALTO (401-3000cc)
func (s *MotorcycleServiceImpl) GetDistinctDisplacements(_ context.Context) ([]domain.EngineDisplacementRange, error) {
	return []domain.EngineDisplacementRange{
		{Range: domain.DisplacementRangeLow},
		{Range: domain.DisplacementRangeMedium},
		{Range: domain.DisplacementRangeHigh},
	}, nil
}

// GetRatingRanges returns the fixed rating range values (HU48)
// Values are hardcoded since they are fixed: 1 (Very bad) to 5 (Excellent)
func (s *MotorcycleServiceImpl) GetRatingRanges(_ context.Context) ([]domain.RatingRange, error) {
	return []domain.RatingRange{
		{Value: 1, Label: "Very bad"},
		{Value: 2, Label: "Bad"},
		{Value: 3, Label: "Average"},
		{Value: 4, Label: "Good"},
		{Value: 5, Label: "Excellent"},
	}, nil
}

// ============================================
// Diagnostic Permissions
// ============================================

// GrantPermission creates or updates a diagnostic permission for a branch on a motorcycle
func (s *MotorcycleServiceImpl) GrantPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string, active bool) (*domain.DiagnosticPermission, error) {
	permission := &domain.DiagnosticPermission{
		MotorcycleID: motorcycleID,
		BranchID:     branchID,
		Active:       active,
	}

	// Check for existing row to reuse its ID (avoids duplicate rows on upsert)
	existing, _ := s.diagPermRepo.GetByMotorcycleAndBranch(ctx, motorcycleID, branchID)
	if existing != nil {
		permission.ID = existing.ID
	} else {
		permission.SetID()
	}

	if err := s.diagPermRepo.Save(ctx, tx, permission); err != nil {
		motorcycleLog.Error(logger.LogDiagPermServiceSaveError, "error", err)
		return nil, domain.ErrPermissionCannotSave
	}
	return permission, nil
}

// RevokePermission deletes a diagnostic permission
func (s *MotorcycleServiceImpl) RevokePermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	if err := s.diagPermRepo.Delete(ctx, tx, motorcycleID, branchID); err != nil {
		motorcycleLog.Error(logger.LogDiagPermServiceDeleteError, "error", err)
		return err
	}
	return nil
}

// ListPermissions retrieves all active diagnostic permissions for a motorcycle
func (s *MotorcycleServiceImpl) ListPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	return s.diagPermRepo.GetByMotorcycleID(ctx, motorcycleID)
}
