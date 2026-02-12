package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockFranchiseService is a mock implementation of input.FranchiseService
type MockFranchiseService struct {
	mock.Mock
}

func (m *MockFranchiseService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockFranchiseService) CreateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) (*domain.Franchise, error) {
	args := m.Called(ctx, tx, franchise)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseService) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	args := m.Called(ctx, franchiseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseService) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	args := m.Called(ctx, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Franchise), args.Error(1)
}

func (m *MockFranchiseService) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	args := m.Called(ctx, tx, franchise)
	return args.Error(0)
}

func (m *MockFranchiseService) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	args := m.Called(ctx, tx, franchiseID)
	return args.Error(0)
}

func (m *MockFranchiseService) AssociateBranches(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	args := m.Called(ctx, tx, franchiseID, branchIDs)
	return args.Error(0)
}

func (m *MockFranchiseService) DissociateBranches(ctx context.Context, tx output.Tx, franchiseID string) error {
	args := m.Called(ctx, tx, franchiseID)
	return args.Error(0)
}

func (m *MockFranchiseService) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}

func (m *MockFranchiseService) CountBranches(ctx context.Context, franchiseID string) (int, error) {
	args := m.Called(ctx, franchiseID)
	return args.Int(0), args.Error(1)
}

func (m *MockFranchiseService) CanRemoveBranch(ctx context.Context, franchiseID string) error {
	args := m.Called(ctx, franchiseID)
	return args.Error(0)
}

func (m *MockFranchiseService) ValidateBranchOwnership(ctx context.Context, branchID, representativeID string) error {
	args := m.Called(ctx, branchID, representativeID)
	return args.Error(0)
}

func (m *MockFranchiseService) ValidateBranchesForFranchise(ctx context.Context, branchIDs []string, representativeID string) error {
	args := m.Called(ctx, branchIDs, representativeID)
	return args.Error(0)
}
