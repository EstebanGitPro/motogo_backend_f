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
}

// NewFranchiseInteractor creates a new FranchiseInteractor
func NewFranchiseInteractor(franchiseService input.FranchiseService) *FranchiseInteractor {
	return &FranchiseInteractor{
		franchiseService: franchiseService,
	}
}

// CreateFranchiseWithBranches creates a franchise and associates existing branches (HU26)
// Business rule: A franchise must have at least 1 branch
func (i *FranchiseInteractor) CreateFranchiseWithBranches(ctx context.Context, franchise domain.Franchise, branchIDs []string, representativeID string) (result *domain.Franchise, err error) {
	// 1. Validate branches (ownership + not already associated + at least 1)
	if valErr := i.franchiseService.ValidateBranchesForFranchise(ctx, branchIDs, representativeID); valErr != nil {
		return nil, valErr
	}

	// 2. Begin transaction
	tx, txErr := i.franchiseService.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogFranchiseInteractorTxError, "error", txErr)
		return nil, txErr
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogFranchiseInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogFranchiseInteractorRollbackOK)
			}
		}
	}()

	// 3. Create franchise
	result, err = i.franchiseService.CreateFranchise(ctx, tx, franchise)
	if err != nil {
		return nil, err
	}

	// 4. Associate branches to franchise
	if err = i.franchiseService.AssociateBranches(ctx, tx, result.ID, branchIDs); err != nil {
		return nil, err
	}

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogFranchiseInteractorCommitError, "error", err)
		return nil, err
	}

	log.Info(logger.LogFranchiseInteractorCreateComplete, "franchise_id", result.ID, "branch_count", len(branchIDs))

	err = nil
	return result, nil
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
func (i *FranchiseInteractor) UpdateFranchise(ctx context.Context, franchise domain.Franchise, representativeID string) (err error) {
	// 1. Verify franchise exists and ownership via branches
	count, countErr := i.franchiseService.CountBranches(ctx, franchise.ID)
	if countErr != nil {
		return countErr
	}
	if count == 0 {
		return domain.ErrFranchiseNotFound
	}

	// 2. Begin transaction
	tx, txErr := i.franchiseService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogFranchiseInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogFranchiseInteractorRollbackOK)
			}
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

	err = nil
	return nil
}

// DeleteFranchise removes a franchise (HU28)
// This only removes the franchise record, branches remain but with franchise_id = NULL
func (i *FranchiseInteractor) DeleteFranchise(ctx context.Context, franchiseID, representativeID string) (err error) {
	// 1. Begin transaction
	tx, txErr := i.franchiseService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogFranchiseInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogFranchiseInteractorRollbackOK)
			}
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

	err = nil
	return nil
}

// AddBranchToFranchise associates an additional branch to an existing franchise
func (i *FranchiseInteractor) AddBranchToFranchise(ctx context.Context, franchiseID, branchID, representativeID string) (err error) {
	// 1. Validate branch ownership via service
	if valErr := i.franchiseService.ValidateBranchOwnership(ctx, branchID, representativeID); valErr != nil {
		return valErr
	}

	// 2. Begin transaction
	tx, txErr := i.franchiseService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogFranchiseInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogFranchiseInteractorRollbackOK)
			}
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

	err = nil
	return nil
}

// RemoveBranchFromFranchise removes a branch from a franchise
// If this is the last branch, returns an error (franchise must have at least 1 branch)
func (i *FranchiseInteractor) RemoveBranchFromFranchise(ctx context.Context, franchiseID, branchID, representativeID string) (err error) {
	// 1. Validate minimum branches - cannot remove last branch
	if canErr := i.franchiseService.CanRemoveBranch(ctx, franchiseID); canErr != nil {
		return canErr
	}

	// 2. Begin transaction
	tx, txErr := i.franchiseService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogFranchiseInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogFranchiseInteractorRollbackOK)
			}
		}
	}()

	// 3. Remove single branch
	if err = i.franchiseService.DissociateSingleBranch(ctx, tx, branchID); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Info(logger.LogFranchiseBranchRemoved, "franchise_id", franchiseID, "branch_id", branchID)

	err = nil
	return nil
}
