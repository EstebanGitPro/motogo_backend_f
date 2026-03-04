package interactor

import (
	"context"
	"net/url"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// BranchInteractor handles branch-related use cases (HU59)
type BranchInteractor struct {
	branchService input.BranchService
	storageClient output.StorageFileDeleter // Optional: Firebase Storage for image deletion
}

// NewBranchInteractor creates a new BranchInteractor instance
func NewBranchInteractor(branchService input.BranchService) *BranchInteractor {
	return &BranchInteractor{
		branchService: branchService,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (i *BranchInteractor) WithStorageClient(client output.StorageFileDeleter) *BranchInteractor {
	i.storageClient = client
	return i
}

// RegisterBranch registers a new branch for an authenticated representative (HU59)
// This method orchestrates the branch registration with proper transaction management
// Returns: (branch, geocodingSucceeded, error)
// - geocodingSucceeded: true if coordinates were auto-generated, false otherwise
// NOTE: Business logic (ID generation, defaults, validation) is handled by the service layer
func (i *BranchInteractor) RegisterBranch(ctx context.Context, branch domain.Branch) (*domain.Branch, bool, error) {
	// Extract traceID from context and create logger with it
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorRegStart,
		"branch_name", branch.Name,
		"representative_id", branch.RepresentativeID,
		"establishment_type", branch.EstablishmentType)

	// STEP 1: Pre-transaction validation and geocoding
	if err := i.validateBranchInputs(ctx, branch, log); err != nil {
		return nil, false, err
	}

	geocodingSucceeded := i.geocodeBranchLocation(ctx, branch.Location, log)

	// STEP 2: Begin transaction
	tx, err := i.branchService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogBranchInteractorTxError, "error", err)
		return nil, false, err
	}
	log.Debug(logger.LogBranchInteractorTxStarted)

	defer func() { _ = tx.Rollback() }()

	// STEP 3: Delegate to service (handles ID generation, defaults, save operations)
	savedBranch, err := i.branchService.RegisterBranch(ctx, tx, branch)
	if err != nil {
		log.Error(logger.LogBranchInteractorRegError, "error", err)
		return nil, false, err
	}
	log.Debug(logger.LogBranchInteractorRegSaved, "branch_id", savedBranch.ID)

	// STEP 4: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogBranchInteractorCommitError, "error", err)
		return nil, false, err
	}

	log.Success(logger.LogBranchInteractorRegComplete,
		"branch_id", savedBranch.ID,
		"branch_name", savedBranch.Name,
		"representative_id", savedBranch.RepresentativeID,
		"geocoding_succeeded", geocodingSucceeded)

	return savedBranch, geocodingSucceeded, nil
}

// GetBranchByID retrieves a branch by its ID
func (i *BranchInteractor) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorGetByID, "branch_id", branchID)

	branch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchInteractorGetByIDError, "error", err, "branch_id", branchID)
		return nil, err
	}

	log.Success(logger.LogBranchInteractorGetByIDOK, "branch_id", branchID)
	return branch, nil
}

// GeocodeLocation wraps the service method for testing purposes
func (i *BranchInteractor) GeocodeLocation(ctx context.Context, location *domain.Location) (bool, error) {
	return i.branchService.GeocodeLocation(ctx, location)
}

// GetBranchesByRepresentative retrieves all branches for a representative (HU62 - List my branches)
func (i *BranchInteractor) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorListByRep, "representative_id", representativeID)

	branches, err := i.branchService.GetBranchesByRepresentative(ctx, representativeID)
	if err != nil {
		log.Error(logger.LogBranchInteractorListByRepErr, "error", err, "representative_id", representativeID)
		return nil, err
	}

	log.Success(logger.LogBranchInteractorListByRepOK, "representative_id", representativeID, "count", len(branches))
	return branches, nil
}

// UpdateBranch updates an existing branch with ownership validation (HU60)
// Returns: (updatedBranch, geocodingSucceeded, error)
func (i *BranchInteractor) UpdateBranch(ctx context.Context, branchID string, branch domain.Branch, personID string) (*domain.Branch, bool, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorUpdateStart, "branch_id", branchID, "person_id", personID)

	// 1. Get existing branch and validate ownership
	existingBranch, err := i.getOwnedBranch(ctx, branchID, personID, log)
	if err != nil {
		return nil, false, err
	}

	// 2. Validate input constraints (brands and displacement ranges)
	if err := i.validateBranchInputs(ctx, branch, log); err != nil {
		return nil, false, err
	}

	// 3. Geocode location if address changed and no coordinates provided
	geocodingSucceeded := i.geocodeBranchLocation(ctx, branch.Location, log)

	// 4. Begin transaction
	tx, err := i.branchService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogBranchInteractorTxError, "error", err)
		return nil, false, err
	}

	defer func() { _ = tx.Rollback() }()

	// 6. Set branch ID, representative ID, and status (preserve from existing)
	branch.ID = branchID
	branch.RepresentativeID = existingBranch.RepresentativeID
	branch.Status = existingBranch.Status // Preserve existing status

	// 7. Update branch via service
	if err = i.branchService.UpdateBranch(ctx, tx, branch); err != nil {
		log.Error(logger.LogBranchInteractorUpdateError, "error", err, "branch_id", branchID)
		return nil, false, err
	}

	// 8. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogBranchInteractorCommitError, "error", err)
		return nil, false, err
	}

	// 9. Delete old profile image from Firebase Storage if replaced
	i.cleanupOldProfileImage(ctx, branch, existingBranch, log)

	// 10. Fetch updated branch for response
	updatedBranch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Warn(logger.LogBranchRefetchFailed, "error", err, "branch_id", branchID)
		// Return original data, update was successful
		return &branch, geocodingSucceeded, nil
	}

	log.Success(logger.LogBranchInteractorUpdateComplete, "branch_id", branchID, "geocoding_succeeded", geocodingSucceeded)
	return updatedBranch, geocodingSucceeded, nil
}

// DeleteBranch deletes a branch with ownership validation (HU61)
// A branch cannot be deleted if it has diagnostics or completed_services (FK RESTRICT)
func (i *BranchInteractor) DeleteBranch(ctx context.Context, branchID string, personID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorDeleteStart, "branch_id", branchID, "person_id", personID)

	// 1. Get existing branch and validate ownership
	existingBranch, err := i.getOwnedBranch(ctx, branchID, personID, log)
	if err != nil {
		return err
	}

	// 2. Begin transaction
	tx, err := i.branchService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogBranchInteractorTxError, "error", err)
		return err
	}

	defer func() { _ = tx.Rollback() }()

	// 4. Delete branch via service
	// NOTE: FK RESTRICT on diagnostics/completed_services is interpreted by service layer
	// CASCADE handles branch_brands, locations, schedules, branch_services
	if err = i.branchService.DeleteBranch(ctx, tx, branchID); err != nil {
		log.Error(logger.LogBranchInteractorDeleteError, "error", err, "branch_id", branchID)
		return err
	}

	// 5. Delete profile image from Firebase Storage (if exists and client configured)
	if existingBranch.ProfileImageURL != nil && *existingBranch.ProfileImageURL != "" && i.storageClient != nil {
		if storageErr := i.storageClient.DeleteStorageFile(ctx, *existingBranch.ProfileImageURL); storageErr != nil {
			// Log warning but don't fail the delete - image cleanup is best effort
			log.Warn(logger.LogBranchInteractorDeleteError, "storage delete failed (continuing)", storageErr)
		} else {
			log.Info(logger.LogBranchInteractorDeleteComplete, "action", "storage_file_deleted", "url", *existingBranch.ProfileImageURL)
		}
	}

	// 6. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogBranchInteractorCommitError, "error", err)
		return err
	}

	log.Success(logger.LogBranchInteractorDeleteComplete, "branch_id", branchID)
	return nil
}

// GetBranchesNearby retrieves branches within radius of given coordinates (HU89)
// Default radius is 5km if not specified
func (i *BranchInteractor) GetBranchesNearby(ctx context.Context, lat, lng, radiusKm float64, establishmentType, brandID, displacementRange string) ([]domain.NearbyBranch, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorNearbyStart,
		"lat", lat,
		"lng", lng,
		"radius_km", radiusKm,
		"type", establishmentType,
		"brand_id", brandID,
		"displacement_range", displacementRange)

	branches, err := i.branchService.GetBranchesNearby(ctx, lat, lng, radiusKm, establishmentType, brandID, displacementRange)
	if err != nil {
		log.Error(logger.LogBranchInteractorNearbyError, "error", err)
		return nil, err
	}

	log.Success(logger.LogBranchInteractorNearbyComplete, "count", len(branches))
	return branches, nil
}

// geocodeBranchLocation geocodes the branch location if provided (non-fatal on error).
func (i *BranchInteractor) geocodeBranchLocation(ctx context.Context, location *domain.Location, log logger.Logger) bool {
	if location == nil {
		return false
	}
	succeeded, err := i.branchService.GeocodeLocation(ctx, location)
	if err != nil {
		log.Warn(logger.LogBranchGeocodingFailed, "error", err)
	}
	if succeeded && location.Latitude != nil && location.Longitude != nil {
		log.Info(logger.LogBranchGeocodingGenerated, "lat", *location.Latitude, "lng", *location.Longitude)
	}
	return succeeded
}

// getOwnedBranch retrieves a branch and validates that it belongs to the given person.
func (i *BranchInteractor) getOwnedBranch(ctx context.Context, branchID, personID string, log logger.Logger) (*domain.Branch, error) {
	branch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchInteractorGetByIDError, "error", err, "branch_id", branchID)
		return nil, err
	}
	if branch.RepresentativeID != personID {
		log.Warn(logger.LogBranchInteractorOwnershipError, "branch_id", branchID, "owner_id", branch.RepresentativeID, "person_id", personID)
		return nil, domain.ErrForbidden
	}
	return branch, nil
}

// validateBranchInputs validates brands and displacement ranges if provided.
func (i *BranchInteractor) validateBranchInputs(ctx context.Context, branch domain.Branch, log logger.Logger) error {
	if len(branch.Brands) > 0 {
		if err := i.branchService.ValidateBrands(ctx, branch.Brands); err != nil {
			log.Warn(logger.LogBranchInteractorValidationError, "error", err, "brands", branch.Brands)
			return err
		}
	}
	if len(branch.DisplacementRanges) > 0 {
		rangeStrs := make([]string, len(branch.DisplacementRanges))
		for j, r := range branch.DisplacementRanges {
			rangeStrs[j] = string(r)
		}
		if err := i.branchService.ValidateDisplacementRanges(rangeStrs); err != nil {
			log.Warn(logger.LogBranchInteractorValidationError, "error", err, "displacement_ranges", branch.DisplacementRanges)
			return err
		}
	}
	return nil
}

// cleanupOldProfileImage deletes the old profile image from Firebase Storage if it was replaced.
// It compares storage object paths (not full URLs) to avoid deleting an overwritten file
// whose download URL changed only by its token query parameter.
func (i *BranchInteractor) cleanupOldProfileImage(ctx context.Context, branch domain.Branch, existingBranch *domain.Branch, log logger.Logger) {
	if branch.ProfileImageURL == nil || *branch.ProfileImageURL == "" {
		return
	}
	if existingBranch.ProfileImageURL == nil || *existingBranch.ProfileImageURL == "" {
		return
	}
	// Compare by storage object path, not full URL (tokens change on overwrite)
	if sameStoragePath(*branch.ProfileImageURL, *existingBranch.ProfileImageURL) {
		log.Info(logger.LogBranchInteractorUpdateComplete, "action", "same_storage_path_skip_delete")
		return
	}
	if i.storageClient == nil {
		return
	}

	if storageErr := i.storageClient.DeleteStorageFile(ctx, *existingBranch.ProfileImageURL); storageErr != nil {
		log.Warn(logger.LogBranchInteractorUpdateError, "old image delete failed (continuing)", storageErr)
	} else {
		log.Info(logger.LogBranchInteractorUpdateComplete, "action", "old_storage_file_deleted", "url", *existingBranch.ProfileImageURL)
	}
}

// sameStoragePath compares two Firebase Storage download URLs by their object path,
// ignoring query parameters (like token) that change when a file is overwritten in place.
func sameStoragePath(url1, url2 string) bool {
	path1 := extractStoragePath(url1)
	path2 := extractStoragePath(url2)
	if path1 == "" || path2 == "" {
		// Fallback to full string comparison if paths can't be extracted
		return url1 == url2
	}
	return path1 == path2
}

// extractStoragePath extracts the object path from a Firebase Storage download URL.
// Input:  https://firebasestorage.googleapis.com/v0/b/bucket/o/branches%2Fid%2Fprofile.jpg?alt=media&token=xxx
// Output: branches/id/profile.jpg
func extractStoragePath(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	// Expected path format: /v0/b/{bucket}/o/{encoded-object-path}
	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 5 || parts[1] != "v0" || parts[2] != "b" || parts[4] != "o" {
		return ""
	}
	objectPathEncoded := strings.Join(parts[5:], "/")
	objectPath, err := url.PathUnescape(objectPathEncoded)
	if err != nil {
		return ""
	}
	return objectPath
}
