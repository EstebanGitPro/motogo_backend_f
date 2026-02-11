package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// MotorcycleInteractor handles motorcycle-related use cases (HU43-47)
type MotorcycleInteractor struct {
	motorcycleService input.MotorcycleService
}

// NewMotorcycleInteractor creates a new MotorcycleInteractor instance
func NewMotorcycleInteractor(motorcycleService input.MotorcycleService) *MotorcycleInteractor {
	return &MotorcycleInteractor{
		motorcycleService: motorcycleService,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (i *MotorcycleInteractor) WithStorageClient(client output.StorageClient) *MotorcycleInteractor {
	i.motorcycleService.WithStorageClient(client)
	return i
}

// RegisterMotorcycle registers a new motorcycle for the authenticated user (HU43)
func (i *MotorcycleInteractor) RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorRegStart, "license_plate", motorcycle.LicensePlate, "owner_id", motorcycle.OwnerID)

	// Step 1: Validate reference_id exists in catalog
	if err := i.motorcycleService.ValidateReferenceExists(ctx, motorcycle.ReferenceID); err != nil {
		log.Warn(logger.LogMotorcycleInteractorRefNotFound, "reference_id", motorcycle.ReferenceID)
		return nil, err
	}

	// Step 2: Validate license plate is unique
	if err := i.motorcycleService.ValidateLicensePlateUnique(ctx, motorcycle.LicensePlate); err != nil {
		log.Warn(logger.LogMotorcycleInteractorDupPlate, "license_plate", motorcycle.LicensePlate)
		return nil, err
	}

	// Step 3: Begin transaction
	tx, err := i.motorcycleService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 4: Create motorcycle (UUID generation + save)
	result, err := i.motorcycleService.CreateMotorcycle(ctx, tx, motorcycle)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorSaveError, "error", err)
		return nil, err
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}

	log.Success(logger.LogMotorcycleInteractorRegSuccess, "id", result.ID, "license_plate", result.LicensePlate)
	return result, nil
}

// GetMotorcycleByID retrieves a motorcycle by its ID (HU46)
func (i *MotorcycleInteractor) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetStart, "motorcycle_id", motorcycleID)

	motorcycle, err := i.motorcycleService.GetMotorcycleByID(ctx, motorcycleID)
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

	motorcycles, err := i.motorcycleService.GetMotorcyclesByOwner(ctx, ownerID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetOwnerError, "error", err, "owner_id", ownerID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetOwnerSuccess, "owner_id", ownerID, "count", len(motorcycles))
	return motorcycles, nil
}

// GetMotorcycleByLicensePlate retrieves a motorcycle by license plate (HU47)
// This endpoint is accessible by representatives (workshops) to lookup motorcycle info
// branchIDs are the branches associated with the representative — at least one must have permission
func (i *MotorcycleInteractor) GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string, branchIDs []string) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetPlateStart, "license_plate", licensePlate)

	motorcycle, err := i.motorcycleService.GetMotorcycleByLicensePlate(ctx, licensePlate)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetPlateError, "error", err, "license_plate", licensePlate)
		return nil, err
	}

	// Validate branch permission — at least one of the representative's branches must be authorized
	if err := i.motorcycleService.ValidateBranchPermission(ctx, motorcycle.ID, branchIDs); err != nil {
		log.Warn(logger.LogMotorcycleInteractorPermError, "motorcycle_id", motorcycle.ID, "license_plate", licensePlate)
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

	// Step 1: Validate ownership and get existing motorcycle
	motorcycle, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID)
	if err != nil {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 2: Apply updates (includes reference validation if changed)
	if err := i.motorcycleService.ApplyMotorcycleUpdates(ctx, motorcycle, updates); err != nil {
		log.Error(logger.LogMotorcycleInteractorRefError, "error", err)
		return nil, err
	}

	// Step 3: Begin transaction
	tx, err := i.motorcycleService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return nil, domain.ErrMotorcycleCannotUpdate
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 4: Update motorcycle
	err = i.motorcycleService.UpdateMotorcycle(ctx, tx, motorcycle)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	log.Success(logger.LogMotorcycleInteractorUpdateSuccess, "motorcycle_id", motorcycleID)

	// Step 6: Return updated motorcycle with reference info
	return i.motorcycleService.GetMotorcycleByID(ctx, motorcycleID)
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

	// Step 1: Validate ownership
	motorcycle, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID)
	if err != nil {
		log.Warn(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	// Step 2: Check service history
	hasHistory, err := i.motorcycleService.CheckServiceHistory(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorDeleteError, "error checking history", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotDelete
	}

	// Step 3: Delete profile image from storage (best effort)
	if motorcycle.ProfileImageURL != nil && *motorcycle.ProfileImageURL != "" {
		i.motorcycleService.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL)
	}

	// Step 4: Begin transaction
	tx, err := i.motorcycleService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return domain.ErrMotorcycleCannotDelete
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 5: Delete motorcycle (strategy determined by service)
	strategy := map[bool]string{true: "soft_delete", false: "hard_delete"}[hasHistory]
	log.Info(logger.LogMotorcycleInteractorDeleteStart, "strategy", strategy, "motorcycle_id", motorcycleID)
	err = i.motorcycleService.DeleteMotorcycle(ctx, tx, motorcycleID, hasHistory)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotDelete
	}

	// Step 6: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return domain.ErrMotorcycleCannotDelete
	}

	log.Success(logger.LogMotorcycleInteractorDeleteSuccess, "motorcycle_id", motorcycleID, "strategy", strategy)
	return nil
}

// DeleteProfileImage removes profile image from both Firebase Storage and database (HU39)
// Only owner can delete their motorcycle's image - returns 404 for non-owners
func (i *MotorcycleInteractor) DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorUpdateStart, "action", "delete_profile_image", "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate ownership
	motorcycle, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID)
	if err != nil {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	// Step 2: Check if there's an image to delete
	if motorcycle.ProfileImageURL == nil || *motorcycle.ProfileImageURL == "" {
		log.Info(logger.LogMotorcycleInteractorUpdateSuccess, "action", "delete_profile_image", "result", "no_image_to_delete")
		return nil // Nothing to delete
	}

	// Step 3: Delete from Firebase Storage (best effort)
	i.motorcycleService.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL)

	// Step 4: Begin transaction
	tx, err := i.motorcycleService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return domain.ErrMotorcycleCannotUpdate
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 5: Clear profile image URL in database
	err = i.motorcycleService.DeleteProfileImage(ctx, tx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleCannotUpdate
	}

	// Step 6: Commit transaction
	if err = tx.Commit(); err != nil {
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

	references, err := i.motorcycleService.GetAllReferences(ctx)
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

	references, err := i.motorcycleService.GetReferencesByBrandID(ctx, brandID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBrandLinesError, "error", err, "brand_id", brandID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorBrandLinesSuccess, "brand_id", brandID, "count", len(references))
	return references, nil
}

// GrantDiagnosticPermission grants a branch permission to view motorcycle diagnostic details
// Only the motorcycle owner can grant permissions
func (i *MotorcycleInteractor) GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string, active bool) (*domain.DiagnosticPermission, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorGrantStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle ownership
	if _, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, err
	}

	// Step 2: Begin transaction
	tx, err := i.motorcycleService.BeginPermissionTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", err)
		return nil, domain.ErrPermissionCannotSave
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 3: Grant permission
	permission, err := i.motorcycleService.GrantDiagnosticPermission(ctx, tx, motorcycleID, branchID, active)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorSaveError, "error", err)
		return nil, err
	}

	// Step 4: Commit
	if err = tx.Commit(); err != nil {
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

	// Step 1: Validate motorcycle ownership
	if _, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return err
	}

	// Step 2: Begin transaction
	tx, err := i.motorcycleService.BeginPermissionTx(ctx)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", err)
		return domain.ErrPermissionCannotDelete
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Step 3: Revoke permission
	err = i.motorcycleService.RevokeDiagnosticPermission(ctx, tx, motorcycleID, branchID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorDeleteError, "error", err)
		return err
	}

	// Step 4: Commit
	if err = tx.Commit(); err != nil {
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

	// Step 1: Validate motorcycle ownership
	if _, err := i.motorcycleService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, err
	}

	// Step 2: Retrieve permissions
	permissions, err := i.motorcycleService.ListDiagnosticPermissions(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagPermInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(permissions))
	return permissions, nil
}
