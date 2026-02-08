package services

import (
	"context"
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var log logger.Logger = logger.NewSlogLogger()

// franchiseService implements input.FranchiseService
type franchiseService struct {
	repository output.FranchiseRepository
}

// NewFranchiseService creates a new FranchiseService instance
func NewFranchiseService(repo output.FranchiseRepository) input.FranchiseService {
	return &franchiseService{
		repository: repo,
	}
}

// BeginTx starts a new database transaction
func (s *franchiseService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

// CreateFranchise creates a new franchise
func (s *franchiseService) CreateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) (*domain.Franchise, error) {
	// 1. Check for duplicate name
	existing, err := s.repository.GetFranchiseByName(ctx, franchise.Name)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error checking franchise name", err)
		return nil, err
	}
	if existing != nil {
		log.Warn(logger.LogFranchiseNameExists, "name", franchise.Name)
		return nil, domain.ErrFranchiseDuplicateName
	}

	// 2. Generate UUID
	if franchise.ID == "" {
		franchise.SetID()
	}

	// 3. Save franchise
	if err := s.repository.SaveFranchise(ctx, tx, franchise); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error saving franchise", err)
		return nil, domain.ErrFranchiseCannotSave
	}

	log.Info(logger.LogFranchiseCreated, "franchise_id", franchise.ID, "name", franchise.Name)
	return &franchise, nil
}

// GetFranchiseByID retrieves a franchise by ID
func (s *franchiseService) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	franchise, err := s.repository.GetFranchiseByID(ctx, franchiseID)
	if err != nil {
		if errors.Is(err, domain.ErrFranchiseNotFound) {
			return nil, err
		}
		log.Error(logger.LogDatabaseUnavailable, "error getting franchise", err, "franchise_id", franchiseID)
		return nil, err
	}
	return franchise, nil
}

// GetFranchisesByRepresentative lists franchises owned by a representative
func (s *franchiseService) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	return s.repository.GetFranchisesByRepresentative(ctx, representativeID)
}

// UpdateFranchise updates an existing franchise
func (s *franchiseService) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	// 1. Check if franchise exists
	existing, err := s.repository.GetFranchiseByID(ctx, franchise.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrFranchiseNotFound
	}

	// 2. Check for duplicate name (if changed)
	if existing.Name != franchise.Name {
		duplicate, err := s.repository.GetFranchiseByName(ctx, franchise.Name)
		if err != nil {
			return err
		}
		if duplicate != nil {
			return domain.ErrFranchiseDuplicateName
		}
	}

	// 3. Update franchise
	if err := s.repository.UpdateFranchise(ctx, tx, franchise); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error updating franchise", err)
		return domain.ErrFranchiseCannotUpdate
	}

	log.Info(logger.LogFranchiseUpdated, "franchise_id", franchise.ID, "name", franchise.Name)
	return nil
}

// DeleteFranchise deletes a franchise
// Note: Branch dissociation should be handled by the interactor before calling this method
func (s *franchiseService) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	// 1. Verify franchise exists
	existing, err := s.repository.GetFranchiseByID(ctx, franchiseID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrFranchiseNotFound
	}

	// 2. Delete franchise (interactor handles branch dissociation first)
	if err := s.repository.DeleteFranchise(ctx, tx, franchiseID); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error deleting franchise", err)
		return domain.ErrFranchiseCannotDelete
	}

	log.Info(logger.LogFranchiseDeleted, "franchise_id", franchiseID)
	return nil
}

// AssociateBranches associates branches to a franchise
func (s *franchiseService) AssociateBranches(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	if len(branchIDs) == 0 {
		return domain.ErrFranchiseNoBranches
	}
	return s.repository.AssociateBranchesToFranchise(ctx, tx, franchiseID, branchIDs)
}

// DissociateBranches removes all branches from a franchise
func (s *franchiseService) DissociateBranches(ctx context.Context, tx output.Tx, franchiseID string) error {
	return s.repository.DissociateBranchesFromFranchise(ctx, tx, franchiseID)
}

// DissociateSingleBranch removes franchise association from a single branch
func (s *franchiseService) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	return s.repository.DissociateSingleBranch(ctx, tx, branchID)
}

// CountBranches returns the number of branches in a franchise
func (s *franchiseService) CountBranches(ctx context.Context, franchiseID string) (int, error) {
	return s.repository.CountBranchesByFranchise(ctx, franchiseID)
}

// CanRemoveBranch checks if a branch can be removed from a franchise
// Business rule: a franchise must have at least 1 branch
func (s *franchiseService) CanRemoveBranch(ctx context.Context, franchiseID string) error {
	count, err := s.repository.CountBranchesByFranchise(ctx, franchiseID)
	if err != nil {
		return err
	}
	if count <= 1 {
		log.Warn(logger.LogFranchiseCannotRemoveLast, "franchise_id", franchiseID)
		return domain.ErrFranchiseMinBranches
	}
	return nil
}
