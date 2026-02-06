package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var log logger.Logger = logger.NewSlogLogger()

// FranchiseInteractor orchestrates franchise operations (HU26-29)
type FranchiseInteractor struct {
	franchiseService input.FranchiseService
	branchService    input.BranchService
}

// NewFranchiseInteractor creates a new FranchiseInteractor
func NewFranchiseInteractor(franchiseService input.FranchiseService, branchService input.BranchService) *FranchiseInteractor {
	return &FranchiseInteractor{
		franchiseService: franchiseService,
		branchService:    branchService,
	}
}

// CreateFranchiseWithBranches creates a franchise and associates existing branches (HU26)
// Business rule: A franchise must have at least 1 branch
func (i *FranchiseInteractor) CreateFranchiseWithBranches(ctx context.Context, franchise domain.Franchise, branchIDs []string, representativeID string) (*domain.Franchise, error) {
	// 1. Validate at least 1 branch
	if len(branchIDs) == 0 {
		log.Warn(logger.LogFranchiseInteractorNoBranches)
		return nil, domain.ErrFranchiseNoBranches
	}

	// 2. Validate branches belong to representative
	for _, branchID := range branchIDs {
		branch, err := i.branchService.GetBranchByID(ctx, branchID)
		if err != nil {
			return nil, domain.ErrBranchNotFound
		}
		if branch.RepresentativeID != representativeID {
			log.Warn(logger.LogFranchiseInteractorBranchNotOwned, "branch_id", branchID, "representative_id", representativeID)
			return nil, domain.ErrFranchiseBranchNotOwned
		}
		// Check if branch is already associated with another franchise
		if branch.FranchiseID != nil && *branch.FranchiseID != "" {
			log.Warn(logger.LogFranchiseInteractorBranchNotOwned, "branch_id", branchID, "franchise_id", *branch.FranchiseID)
			return nil, domain.ErrFranchiseBranchNotOwned
		}
	}

	// 3. Begin transaction
	tx, err := i.franchiseService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogFranchiseInteractorTxError, "error", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Intentionally ignoring rollback error
		}
	}()

	// 4. Create franchise
	createdFranchise, err := i.franchiseService.CreateFranchise(ctx, tx, franchise)
	if err != nil {
		return nil, err
	}

	// 5. Associate branches to franchise
	if err = i.franchiseService.AssociateBranches(ctx, tx, createdFranchise.ID, branchIDs); err != nil {
		return nil, err
	}

	// 6. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogFranchiseInteractorCommitError, "error", err)
		return nil, err
	}

	log.Info(logger.LogFranchiseInteractorCreateComplete, "franchise_id", createdFranchise.ID, "branch_count", len(branchIDs))
	return createdFranchise, nil
}

// GetFranchiseByID retrieves a franchise by ID (HU29)
func (i *FranchiseInteractor) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	return i.franchiseService.GetFranchiseByID(ctx, franchiseID)
}

// GetFranchisesByRepresentative lists franchises for a representative (HU29)
func (i *FranchiseInteractor) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	return i.franchiseService.GetFranchisesByRepresentative(ctx, representativeID)
}

// UpdateFranchise updates franchise info (HU27)
func (i *FranchiseInteractor) UpdateFranchise(ctx context.Context, franchise domain.Franchise, representativeID string) error {
	// 1. Verify franchise exists and ownership via branches
	count, err := i.franchiseService.CountBranches(ctx, franchise.ID)
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrFranchiseNotFound
	}

	// 2. Begin transaction
	tx, err := i.franchiseService.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Intentionally ignoring rollback error
		}
	}()

	// 3. Update franchise
	if err = i.franchiseService.UpdateFranchise(ctx, tx, franchise); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogFranchiseInteractorCommitError, "error", err)
		return err
	}

	log.Info(logger.LogFranchiseInteractorUpdateComplete, "franchise_id", franchise.ID)
	return nil
}

// DeleteFranchise removes a franchise (HU28)
// This only removes the franchise record, branches remain but with franchise_id = NULL
func (i *FranchiseInteractor) DeleteFranchise(ctx context.Context, franchiseID, representativeID string) error {
	// 1. Begin transaction
	tx, err := i.franchiseService.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Intentionally ignoring rollback error
		}
	}()

	// 2. Dissociate all branches (set franchise_id = NULL)
	if err = i.franchiseService.DissociateBranches(ctx, tx, franchiseID); err != nil {
		return err
	}

	// 3. Delete franchise (service validates existence)
	if err = i.franchiseService.DeleteFranchise(ctx, tx, franchiseID); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogFranchiseInteractorCommitError, "error", err)
		return err
	}

	log.Info(logger.LogFranchiseInteractorDeleteComplete, "franchise_id", franchiseID)
	return nil
}

// AddBranchToFranchise associates an additional branch to an existing franchise
func (i *FranchiseInteractor) AddBranchToFranchise(ctx context.Context, franchiseID, branchID, representativeID string) error {
	// 1. Validate branch ownership
	branch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrFranchiseBranchNotOwned
	}

	// 2. Begin transaction
	tx, err := i.franchiseService.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Intentionally ignoring rollback error
		}
	}()

	// 3. Associate branch
	if err = i.franchiseService.AssociateBranches(ctx, tx, franchiseID, []string{branchID}); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Info(logger.LogFranchiseBranchAdded, "franchise_id", franchiseID, "branch_id", branchID)
	return nil
}

// RemoveBranchFromFranchise removes a branch from a franchise
// If this is the last branch, returns an error (franchise must have at least 1 branch)
func (i *FranchiseInteractor) RemoveBranchFromFranchise(ctx context.Context, franchiseID, branchID, representativeID string) error {
	// 1. Check branch count
	count, err := i.franchiseService.CountBranches(ctx, franchiseID)
	if err != nil {
		return err
	}

	// 2. Validate minimum branches - cannot remove last branch
	if count <= 1 {
		log.Warn(logger.LogFranchiseCannotRemoveLast, "franchise_id", franchiseID, "branch_id", branchID)
		return domain.ErrFranchiseMinBranches
	}

	// 3. Begin transaction
	tx, err := i.franchiseService.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Intentionally ignoring rollback error
		}
	}()

	// 4. Remove single branch
	if err = i.franchiseService.DissociateSingleBranch(ctx, tx, branchID); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Info(logger.LogFranchiseBranchRemoved, "franchise_id", franchiseID, "branch_id", branchID)
	return nil
}
