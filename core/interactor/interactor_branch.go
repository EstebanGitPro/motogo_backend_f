package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// BranchInteractor handles branch-related use cases (HU59)
type BranchInteractor struct {
	branchService input.BranchService
	logger        logger.Logger
}

// NewBranchInteractor creates a new BranchInteractor instance
func NewBranchInteractor(branchService input.BranchService, log logger.Logger) *BranchInteractor {
	return &BranchInteractor{
		branchService: branchService,
		logger:        log,
	}
}

// RegisterBranch registers a new branch for an authenticated representative (HU59)
// This method orchestrates the branch registration with proper transaction management
// Returns: (branch, geocodingSucceeded, error)
// - geocodingSucceeded: true if coordinates were auto-generated, false otherwise
// NOTE: Business logic (ID generation, defaults, validation) is handled by the service layer
func (i *BranchInteractor) RegisterBranch(ctx context.Context, branch domain.Branch) (*domain.Branch, bool, error) {
	// Extract traceID from context and create logger with it
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogBranchInteractorRegStart,
		"branch_name", branch.Name,
		"representative_id", branch.RepresentativeID,
		"establishment_type", branch.EstablishmentType)

	// STEP 1: Validate brands if provided (before starting transaction)
	if len(branch.Brands) > 0 {
		if err := i.branchService.ValidateBrands(ctx, branch.Brands); err != nil {
			log.Warn(logger.LogBranchInteractorValidationError, "error", err, "brands", branch.Brands)
			return nil, false, err
		}
		log.Debug(logger.LogBranchInteractorBrandsValidated, "brands_count", len(branch.Brands))
	}

	// STEP 2: Geocode location if coordinates not provided
	// This is done before the transaction to avoid holding it open during external API call
	var geocodingSucceeded bool
	var err error
	if branch.Location != nil {
		geocodingSucceeded, err = i.branchService.GeocodeLocation(ctx, branch.Location)
		if err != nil {
			// Log but don't fail - geocoding errors are non-fatal
			log.Warn("geocoding_step_failed", "error", err)
		}
		if geocodingSucceeded {
			log.Info("geocoding_coordinates_generated",
				"lat", *branch.Location.Latitude,
				"lng", *branch.Location.Longitude)
		}
	}

	// STEP 3: Begin transaction
	tx, err := i.branchService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogBranchInteractorTxError, "error", err)
		return nil, false, err
	}
	log.Debug(logger.LogBranchInteractorTxStarted)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogBranchInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogBranchInteractorRollbackOK)
			}
		}
	}()

	// STEP 4: Delegate to service (handles ID generation, defaults, save operations)
	savedBranch, err := i.branchService.RegisterBranch(ctx, tx, branch)
	if err != nil {
		log.Error(logger.LogBranchInteractorRegError, "error", err)
		return nil, false, err
	}
	log.Debug(logger.LogBranchInteractorRegSaved, "branch_id", savedBranch.ID)

	// STEP 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogBranchInteractorCommitError, "error", err)
		return nil, false, err
	}

	log.Success(logger.LogBranchInteractorRegComplete,
		"branch_id", savedBranch.ID,
		"branch_name", savedBranch.Name,
		"representative_id", savedBranch.RepresentativeID,
		"geocoding_succeeded", geocodingSucceeded)

	err = nil // Ensure defer doesn't execute rollback
	return savedBranch, geocodingSucceeded, nil
}

// GetBranchByID retrieves a branch by its ID
func (i *BranchInteractor) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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
