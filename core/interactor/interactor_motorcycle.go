package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/google/uuid"
)

// MotorcycleInteractor handles motorcycle-related use cases (HU43-47)
type MotorcycleInteractor struct {
	motorcycleRepo output.MotorcycleRepository
	diagPermRepo   output.DiagnosticPermissionRepository // Diagnostic permissions pivot
	storageClient  output.StorageClient                  // Optional: Firebase Storage for image deletion
}

// NewMotorcycleInteractor creates a new MotorcycleInteractor instance
func NewMotorcycleInteractor(motorcycleRepo output.MotorcycleRepository, diagPermRepo output.DiagnosticPermissionRepository) *MotorcycleInteractor {
	return &MotorcycleInteractor{
		motorcycleRepo: motorcycleRepo,
		diagPermRepo:   diagPermRepo,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (i *MotorcycleInteractor) WithStorageClient(client output.StorageClient) *MotorcycleInteractor {
	i.storageClient = client
	return i
}

// RegisterMotorcycle registers a new motorcycle for the authenticated user (HU43)
func (i *MotorcycleInteractor) RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorRegStart, "license_plate", motorcycle.LicensePlate, "owner_id", motorcycle.OwnerID)

	// Step 1: Validate reference_id is provided (required field)
	if motorcycle.ReferenceID == "" {
		log.Warn(logger.LogMotorcycleInteractorRefRequired)
		return nil, domain.ErrReferenceRequired
	}

	// Step 2: Validate reference exists in catalog
	refExists, err := i.motorcycleRepo.ValidateReferenceExists(ctx, motorcycle.ReferenceID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorRefError, "error", err, "reference_id", motorcycle.ReferenceID)
		return nil, domain.ErrMotorcycleCannotSave
	}
	if !refExists {
		log.Warn(logger.LogMotorcycleInteractorRefNotFound, "reference_id", motorcycle.ReferenceID)
		return nil, domain.ErrReferenceNotFound
	}

	// Step 2: Validate license plate is unique
	plateExists, err := i.motorcycleRepo.CheckLicensePlateExists(ctx, motorcycle.LicensePlate)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorCheckPlateErr, "error", err, "license_plate", motorcycle.LicensePlate)
		return nil, domain.ErrMotorcycleCannotSave
	}
	if plateExists {
		log.Warn(logger.LogMotorcycleInteractorDupPlate, "license_plate", motorcycle.LicensePlate)
		return nil, domain.ErrDuplicateLicensePlate
	}

	// Step 3: Generate UUID
	motorcycle.ID = uuid.New().String()
	log.Debug(logger.LogMotorcycleInteractorIDGenerated, "id", motorcycle.ID)

	// Step 4: Begin transaction
	tx, err := i.motorcycleRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}

	// Step 5: Save motorcycle
	err = i.motorcycleRepo.Save(ctx, tx, motorcycle)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorSaveError, "error", err)
		tx.Rollback()
		return nil, domain.ErrMotorcycleCannotSave
	}

	// Step 6: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}

	log.Success(logger.LogMotorcycleInteractorRegSuccess, "id", motorcycle.ID, "license_plate", motorcycle.LicensePlate)
	return motorcycle, nil
}

// GetMotorcycleByID retrieves a motorcycle by its ID (HU46)
func (i *MotorcycleInteractor) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetStart, "motorcycle_id", motorcycleID)

	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetSuccess, "motorcycle_id", motorcycleID)
	return motorcycle, nil
}

// GetMotorcyclesByOwner retrieves all motorcycles owned by a person (HU47)
func (i *MotorcycleInteractor) GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetOwnerStart, "owner_id", ownerID)

	motorcycles, err := i.motorcycleRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetOwnerError, "error", err, "owner_id", ownerID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetOwnerSuccess, "owner_id", ownerID, "count", len(motorcycles))
	return motorcycles, nil
}

// GetMotorcycleByLicensePlate retrieves a motorcycle by license plate (HU47)
// This endpoint is accessible by representatives (workshops) to lookup motorcycle info
func (i *MotorcycleInteractor) GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetPlateStart, "license_plate", licensePlate)

	motorcycle, err := i.motorcycleRepo.GetByLicensePlate(ctx, licensePlate)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetPlateError, "error", err, "license_plate", licensePlate)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetPlateSuccess, "license_plate", licensePlate, "motorcycle_id", motorcycle.ID)
	return motorcycle, nil
}

// UpdateMotorcycle updates motorcycle information (HU44)
// Only owner can update their motorcycle - caller must validate ownership
func (i *MotorcycleInteractor) UpdateMotorcycle(ctx context.Context, motorcycleID string, ownerID string, updates *domain.Motorcycle) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorUpdateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Get existing motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 2: Validate ownership
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "reason", "not owner", "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 3: Validate new reference_id if changed
	if updates.ReferenceID != "" && updates.ReferenceID != motorcycle.ReferenceID {
		refExists, err := i.motorcycleRepo.ValidateReferenceExists(ctx, updates.ReferenceID)
		if err != nil {
			log.Error(logger.LogMotorcycleInteractorRefError, "error", err, "reference_id", updates.ReferenceID)
			return nil, domain.ErrMotorcycleCannotUpdate
		}
		if !refExists {
			log.Warn(logger.LogMotorcycleInteractorRefNotFound, "reference_id", updates.ReferenceID)
			return nil, domain.ErrReferenceNotFound
		}
		motorcycle.ReferenceID = updates.ReferenceID
	}

	// Step 4: Apply updates (only if provided)
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

	// Step 5: Begin transaction
	tx, err := i.motorcycleRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	// Step 6: Update motorcycle
	err = i.motorcycleRepo.Update(ctx, tx, motorcycle)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		tx.Rollback()
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	// Step 7: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	log.Success(logger.LogMotorcycleInteractorUpdateSuccess, "motorcycle_id", motorcycleID)

	// Step 8: Return updated motorcycle with reference info
	return i.motorcycleRepo.GetByID(ctx, motorcycleID)
}

// DeleteMotorcycle implements hybrid delete strategy (HU45):
// - Hard delete: If motorcycle has NO service history (diagnostics/completed_services)
// - Soft delete: If motorcycle HAS service history (preserves historical data)
// In both cases, profile image is removed from Firebase Storage
// Only owner can delete their motorcycle - returns 404 for non-owners
func (i *MotorcycleInteractor) DeleteMotorcycle(ctx context.Context, motorcycleID string, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorDeleteStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Get existing motorcycle to validate ownership
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogMotorcycleInteractorDeleteError, "reason", "not owner", "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return domain.ErrMotorcycleNotFound
	}

	// Step 3: Check if motorcycle has service history
	hasHistory, err := i.motorcycleRepo.HasServiceHistory(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorDeleteError, "error checking history", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotDelete
	}

	// Step 4: Delete profile image from Firebase Storage (if exists and client configured)
	if motorcycle.ProfileImageURL != nil && *motorcycle.ProfileImageURL != "" && i.storageClient != nil {
		if err := i.storageClient.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL); err != nil {
			// Log warning but don't fail the delete - image cleanup is best effort
			log.Warn(logger.LogMotorcycleInteractorDeleteError, "storage delete failed (continuing)", err)
		}
	}

	// Step 5: Begin transaction
	tx, err := i.motorcycleRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return domain.ErrMotorcycleCannotDelete
	}

	// Step 6: Choose delete strategy based on history
	if hasHistory {
		// Soft delete - preserves historical data integrity
		log.Info(logger.LogMotorcycleInteractorDeleteStart, "strategy", "soft_delete", "motorcycle_id", motorcycleID)
		if err = i.motorcycleRepo.Delete(ctx, tx, motorcycleID); err != nil {
			log.Error(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
			tx.Rollback()
			return domain.ErrMotorcycleCannotDelete
		}
	} else {
		// Hard delete - no history to preserve, clean removal
		log.Info(logger.LogMotorcycleInteractorDeleteStart, "strategy", "hard_delete", "motorcycle_id", motorcycleID)
		if err = i.motorcycleRepo.HardDelete(ctx, tx, motorcycleID); err != nil {
			log.Error(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
			tx.Rollback()
			return domain.ErrMotorcycleCannotDelete
		}
	}

	// Step 7: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return domain.ErrMotorcycleCannotDelete
	}

	log.Success(logger.LogMotorcycleInteractorDeleteSuccess, "motorcycle_id", motorcycleID, "strategy", map[bool]string{true: "soft_delete", false: "hard_delete"}[hasHistory])
	return nil
}

// DeleteProfileImage removes profile image from both Firebase Storage and database (HU39)
// Only owner can delete their motorcycle's image - returns 404 for non-owners
func (i *MotorcycleInteractor) DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorUpdateStart, "action", "delete_profile_image", "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Get existing motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "reason", "not owner", "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return domain.ErrMotorcycleNotFound
	}

	// Step 3: Check if there's an image to delete
	if motorcycle.ProfileImageURL == nil || *motorcycle.ProfileImageURL == "" {
		log.Info(logger.LogMotorcycleInteractorUpdateSuccess, "action", "delete_profile_image", "result", "no_image_to_delete")
		return nil // Nothing to delete
	}

	// Step 4: Delete from Firebase Storage (if client configured)
	if i.storageClient != nil {
		if err := i.storageClient.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL); err != nil {
			// Log warning but continue - storage cleanup is best effort
			log.Warn(logger.LogMotorcycleInteractorUpdateError, "storage delete failed (continuing)", err, "url", *motorcycle.ProfileImageURL)
		} else {
			log.Info(logger.LogMotorcycleInteractorUpdateSuccess, "action", "storage_file_deleted", "url", *motorcycle.ProfileImageURL)
		}
	}

	// Step 5: Begin transaction
	tx, err := i.motorcycleRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return domain.ErrMotorcycleCannotUpdate
	}

	// Step 6: Clear profile image URL in database
	if err := i.motorcycleRepo.ClearProfileImageURL(ctx, tx, motorcycleID); err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		tx.Rollback()
		return domain.ErrMotorcycleCannotUpdate
	}

	// Step 7: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return domain.ErrMotorcycleCannotUpdate
	}

	log.Success(logger.LogMotorcycleInteractorUpdateSuccess, "action", "delete_profile_image", "motorcycle_id", motorcycleID)
	return nil
}

// GetMotorcycleReferences retrieves all motorcycle references from catalog (HU50)
func (i *MotorcycleInteractor) GetMotorcycleReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetRefsStart)

	references, err := i.motorcycleRepo.GetAllReferences(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetRefsError, "error", err)
		return nil, err
	}

	return references, nil
}

// GetReferencesByBrandID retrieves motorcycle references for a specific brand (HU40 - Admin only)
func (i *MotorcycleInteractor) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorBrandLinesStart, "brand_id", brandID)

	references, err := i.motorcycleRepo.GetReferencesByBrandID(ctx, brandID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBrandLinesError, "error", err, "brand_id", brandID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorBrandLinesSuccess, "brand_id", brandID, "count", len(references))
	return references, nil
}

// GrantDiagnosticPermission grants a branch permission to view motorcycle diagnostic details
// Only the motorcycle owner can grant permissions
func (i *MotorcycleInteractor) GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) (*domain.DiagnosticPermission, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorGrantStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and belongs to owner
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Create permission entity
	permission := &domain.DiagnosticPermission{
		ID:           utils.Generate(),
		MotorcycleID: motorcycleID,
		BranchID:     branchID,
		Active:       true,
	}

	// Step 3: Begin transaction
	tx, err := i.diagPermRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", err)
		return nil, domain.ErrPermissionCannotSave
	}

	// Step 4: Save permission (upsert)
	if err := i.diagPermRepo.Save(ctx, tx, permission); err != nil {
		log.Error(logger.LogDiagPermInteractorSaveError, "error", err)
		tx.Rollback()
		return nil, domain.ErrPermissionCannotSave
	}

	// Step 5: Commit
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDiagPermInteractorCommitError, "error", err)
		return nil, domain.ErrPermissionCannotSave
	}

	log.Success(logger.LogDiagPermInteractorGrantSuccess, "motorcycle_id", motorcycleID, "branch_id", branchID)
	return permission, nil
}

// RevokeDiagnosticPermission revokes a branch's permission to view motorcycle diagnostic details
// Only the motorcycle owner can revoke permissions
func (i *MotorcycleInteractor) RevokeDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorRevokeStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and belongs to owner
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return domain.ErrMotorcycleNotFound
	}

	// Step 2: Begin transaction
	tx, err := i.diagPermRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", err)
		return domain.ErrPermissionCannotDelete
	}

	// Step 3: Delete permission
	if err := i.diagPermRepo.Delete(ctx, tx, motorcycleID, branchID); err != nil {
		log.Error(logger.LogDiagPermInteractorDeleteError, "error", err)
		tx.Rollback()
		return err // Return the specific error (could be ErrPermissionNotFound)
	}

	// Step 4: Commit
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDiagPermInteractorCommitError, "error", err)
		return domain.ErrPermissionCannotDelete
	}

	log.Success(logger.LogDiagPermInteractorRevokeSuccess, "motorcycle_id", motorcycleID, "branch_id", branchID)
	return nil
}

// ListDiagnosticPermissions retrieves all active diagnostic permissions for a motorcycle
// Only the motorcycle owner can list permissions
func (i *MotorcycleInteractor) ListDiagnosticPermissions(ctx context.Context, motorcycleID, ownerID string) ([]domain.DiagnosticPermission, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorListStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and belongs to owner
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorMotoError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Retrieve permissions
	permissions, err := i.diagPermRepo.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagPermInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(permissions))
	return permissions, nil
}
