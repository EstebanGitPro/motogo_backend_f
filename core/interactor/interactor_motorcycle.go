package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// MotorcycleInteractor handles motorcycle-related use cases (HU43-47)
// Orchestrates ownership checks and transaction management, delegating business logic to MotorcycleService
type MotorcycleInteractor struct {
	motorcycleSvc input.MotorcycleService
}

// NewMotorcycleInteractor creates a new MotorcycleInteractor instance
func NewMotorcycleInteractor(motorcycleSvc input.MotorcycleService) *MotorcycleInteractor {
	return &MotorcycleInteractor{
		motorcycleSvc: motorcycleSvc,
	}
}

// RegisterMotorcycle registers a new motorcycle for the authenticated user (HU43)
func (i *MotorcycleInteractor) RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (result *domain.Motorcycle, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorRegStart, "license_plate", motorcycle.LicensePlate, "owner_id", motorcycle.OwnerID)

	// Step 1: Validate reference exists
	if err = i.motorcycleSvc.ValidateReferenceExists(ctx, motorcycle.ReferenceID); err != nil {
		log.Warn(logger.LogMotorcycleInteractorRefError, "error", err, "reference_id", motorcycle.ReferenceID)
		return nil, err
	}

	// Step 2: Validate license plate is unique
	if err = i.motorcycleSvc.CheckLicensePlateUnique(ctx, motorcycle.LicensePlate); err != nil {
		log.Warn(logger.LogMotorcycleInteractorCheckPlateErr, "error", err, "license_plate", motorcycle.LicensePlate)
		return nil, err
	}

	// Step 3: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrMotorcycleCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMotorcycleInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogMotorcycleInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Register motorcycle (SetID + Save)
	if err = i.motorcycleSvc.RegisterMotorcycle(ctx, tx, motorcycle); err != nil {
		log.Error(logger.LogMotorcycleInteractorSaveError, "error", err)
		return nil, err
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
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

	motorcycle, err := i.motorcycleSvc.GetByID(ctx, motorcycleID)
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

	motorcycles, err := i.motorcycleSvc.GetByOwnerID(ctx, ownerID)
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

	motorcycle, err := i.motorcycleSvc.GetByLicensePlate(ctx, licensePlate)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetPlateError, "error", err, "license_plate", licensePlate)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetPlateSuccess, "license_plate", licensePlate, "motorcycle_id", motorcycle.ID)
	return motorcycle, nil
}

// UpdateMotorcycle updates motorcycle information (HU44)
// Only owner can update their motorcycle
func (i *MotorcycleInteractor) UpdateMotorcycle(ctx context.Context, motorcycleID string, ownerID string, updates *domain.Motorcycle) (result *domain.Motorcycle, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorUpdateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate ownership
	motorcycle, err := i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID)
	if err != nil {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, err
	}

	// Step 2: Apply field updates (includes reference validation if changed)
	if err = i.motorcycleSvc.ApplyUpdates(motorcycle, updates); err != nil {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 3: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMotorcycleInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogMotorcycleInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Persist update
	if err = i.motorcycleSvc.UpdateMotorcycle(ctx, tx, motorcycle); err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotUpdate
	}

	log.Success(logger.LogMotorcycleInteractorUpdateSuccess, "motorcycle_id", motorcycleID)

	// Step 6: Return updated motorcycle with reference info
	result, err = i.motorcycleSvc.GetByID(ctx, motorcycleID)
	if err != nil {
		// Refetch failed but update was successful, return original
		log.Warn(logger.LogMotorcycleInteractorGetError, "error", err, "motorcycle_id", motorcycleID)
		return motorcycle, nil
	}

	return result, nil
}

// DeleteMotorcycle implements hybrid delete strategy (HU45):
// - Hard delete: If motorcycle has NO service history (diagnostics/completed_services)
// - Soft delete: If motorcycle HAS service history (preserves historical data)
// In both cases, profile image is removed from Firebase Storage
// Only owner can delete their motorcycle - returns 404 for non-owners
func (i *MotorcycleInteractor) DeleteMotorcycle(ctx context.Context, motorcycleID string, ownerID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorDeleteStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate ownership
	motorcycle, err := i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID)
	if err != nil {
		log.Warn(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return err
	}

	// Step 2: Delete profile image from Firebase Storage (best-effort)
	if motorcycle.ProfileImageURL != nil && *motorcycle.ProfileImageURL != "" {
		i.motorcycleSvc.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL)
	}

	// Step 3: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", txErr)
		return domain.ErrMotorcycleCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMotorcycleInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogMotorcycleInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Execute hybrid delete (service determines soft vs hard)
	if err = i.motorcycleSvc.DeleteMotorcycle(ctx, tx, motorcycleID); err != nil {
		log.Error(logger.LogMotorcycleInteractorDeleteError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return domain.ErrMotorcycleCannotDelete
	}

	log.Success(logger.LogMotorcycleInteractorDeleteSuccess, "motorcycle_id", motorcycleID)

	return nil
}

// DeleteProfileImage removes profile image from both Firebase Storage and database (HU39)
// Only owner can delete their motorcycle's image - returns 404 for non-owners
func (i *MotorcycleInteractor) DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorUpdateStart, "action", "delete_profile_image", "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate ownership
	motorcycle, ownerErr := i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID)
	if ownerErr != nil {
		log.Warn(logger.LogMotorcycleInteractorUpdateError, "error", ownerErr, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return ownerErr
	}

	// Step 2: Check if there's an image to delete
	if motorcycle.ProfileImageURL == nil || *motorcycle.ProfileImageURL == "" {
		log.Info(logger.LogMotorcycleInteractorUpdateSuccess, "action", "delete_profile_image", "result", "no_image_to_delete")
		return nil // Nothing to delete
	}

	// Step 3: Delete from Firebase Storage (best-effort)
	i.motorcycleSvc.DeleteStorageFile(ctx, *motorcycle.ProfileImageURL)

	// Step 4: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", txErr)
		return domain.ErrMotorcycleCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMotorcycleInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogMotorcycleInteractorRollbackOK)
			}
		}
	}()

	// Step 5: Clear profile image URL in database
	if err = i.motorcycleSvc.ClearProfileImageURL(ctx, tx, motorcycleID); err != nil {
		log.Error(logger.LogMotorcycleInteractorUpdateError, "error", err, "motorcycle_id", motorcycleID)
		return err
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

	references, err := i.motorcycleSvc.GetAllReferences(ctx)
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

	references, err := i.motorcycleSvc.GetReferencesByBrandID(ctx, brandID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBrandLinesError, "error", err, "brand_id", brandID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorBrandLinesSuccess, "brand_id", brandID, "count", len(references))
	return references, nil
}

// GetDistinctCategories retrieves all distinct motorcycle categories with line counts (HU41)
func (i *MotorcycleInteractor) GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleCatInteractorStart)

	categories, err := i.motorcycleSvc.GetDistinctCategories(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleCatInteractorError, "error", err)
		return nil, err
	}

	log.Success(logger.LogMotorcycleCatInteractorSuccess, "count", len(categories))
	return categories, nil
}

// GetLinesByCategory retrieves all motorcycle lines for a specific category (HU41)
func (i *MotorcycleInteractor) GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleCatLinesInteractorStart, "category", categoryName)

	lines, err := i.motorcycleSvc.GetLinesByCategory(ctx, categoryName)
	if err != nil {
		log.Error(logger.LogMotorcycleCatLinesInteractorError, "error", err, "category", categoryName)
		return nil, err
	}

	log.Success(logger.LogMotorcycleCatLinesInteractorSuccess, "category", categoryName, "count", len(lines))
	return lines, nil
}

// GetDistinctDisplacements retrieves all distinct engine displacement values with counts (HU49)
func (i *MotorcycleInteractor) GetDistinctDisplacements(ctx context.Context) ([]domain.EngineDisplacementRange, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleDispInteractorStart)

	displacements, err := i.motorcycleSvc.GetDistinctDisplacements(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleDispInteractorError, "error", err)
		return nil, err
	}

	log.Success(logger.LogMotorcycleDispInteractorSuccess, "count", len(displacements))
	return displacements, nil
}

// GetRatingRanges retrieves all valid rating range values (HU48)
func (i *MotorcycleInteractor) GetRatingRanges(ctx context.Context) ([]domain.RatingRange, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogRatingRangeInteractorStart)

	ranges, err := i.motorcycleSvc.GetRatingRanges(ctx)
	if err != nil {
		log.Error(logger.LogRatingRangeInteractorError, "error", err)
		return nil, err
	}

	log.Success(logger.LogRatingRangeInteractorSuccess, "count", len(ranges))
	return ranges, nil
}

// GrantDiagnosticPermission grants a branch permission to view motorcycle diagnostic details
// Only the motorcycle owner can grant permissions
func (i *MotorcycleInteractor) GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string, active bool) (result *domain.DiagnosticPermission, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorGrantStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle ownership
	if _, err = i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, err
	}

	// Step 2: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrPermissionCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagPermInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagPermInteractorRollbackOK)
			}
		}
	}()

	// Step 3: Grant permission (entity creation + save)
	result, err = i.motorcycleSvc.GrantPermission(ctx, tx, motorcycleID, branchID, active)
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

	return result, nil
}

// RevokeDiagnosticPermission revokes a branch's permission to view motorcycle diagnostic details
// Only the motorcycle owner can revoke permissions
func (i *MotorcycleInteractor) RevokeDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorRevokeStart, "motorcycle_id", motorcycleID, "branch_id", branchID, "owner_id", ownerID)

	// Step 1: Validate motorcycle ownership
	if _, err = i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return err
	}

	// Step 2: Begin transaction
	tx, txErr := i.motorcycleSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogDiagPermInteractorBeginTxError, "error", txErr)
		return domain.ErrPermissionCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDiagPermInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogDiagPermInteractorRollbackOK)
			}
		}
	}()

	// Step 3: Revoke permission
	if err = i.motorcycleSvc.RevokePermission(ctx, tx, motorcycleID, branchID); err != nil {
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
	if _, err := i.motorcycleSvc.ValidateOwnership(ctx, motorcycleID, ownerID); err != nil {
		log.Warn(logger.LogDiagPermInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, err
	}

	// Step 2: Retrieve permissions
	permissions, err := i.motorcycleSvc.ListPermissions(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagPermInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(permissions))
	return permissions, nil
}

// LookupPermissions retrieves active diagnostic permissions for a motorcycle (no ownership check)
// Used by representatives to verify access during plate lookup
func (i *MotorcycleInteractor) LookupPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDiagPermInteractorLookupStart, "motorcycle_id", motorcycleID)

	permissions, err := i.motorcycleSvc.ListPermissions(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermInteractorLookupError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogDiagPermInteractorLookupSuccess, "motorcycle_id", motorcycleID, "count", len(permissions))
	return permissions, nil
}
