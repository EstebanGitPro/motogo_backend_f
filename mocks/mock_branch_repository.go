package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockBranchRepository is a mock implementation of output.BranchRepository
type MockBranchRepository struct {
	mock.Mock
}

func (m *MockBranchRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockBranchRepository) SaveBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	args := m.Called(ctx, tx, branch)
	return args.Error(0)
}

func (m *MockBranchRepository) UpdateBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	args := m.Called(ctx, tx, branch)
	return args.Error(0)
}

func (m *MockBranchRepository) DeleteBranch(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}

func (m *MockBranchRepository) SaveBranchBrands(ctx context.Context, tx output.Tx, branchID string, brands []string) error {
	args := m.Called(ctx, tx, branchID, brands)
	return args.Error(0)
}

func (m *MockBranchRepository) DeleteBranchBrands(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}

func (m *MockBranchRepository) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) GetBranchByFranchiseAndName(ctx context.Context, franchiseID, name string) (*domain.Branch, error) {
	args := m.Called(ctx, franchiseID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	args := m.Called(ctx, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) HasBranchesByRepresentative(ctx context.Context, representativeID string) (bool, error) {
	args := m.Called(ctx, representativeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBranchRepository) ValidateBrands(ctx context.Context, brands []string) error {
	args := m.Called(ctx, brands)
	return args.Error(0)
}

func (m *MockBranchRepository) GetBranchesNearby(ctx context.Context, params domain.NearbySearchParams) ([]domain.NearbyBranch, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.NearbyBranch), args.Error(1)
}

func (m *MockBranchRepository) SaveBranchDisplacementRanges(ctx context.Context, tx output.Tx, branchID string, ranges []string) error {
	args := m.Called(ctx, tx, branchID, ranges)
	return args.Error(0)
}

func (m *MockBranchRepository) DeleteBranchDisplacementRanges(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}
