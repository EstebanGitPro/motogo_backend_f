package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockFranchiseRepository is a mock implementation of output.FranchiseRepository
type MockFranchiseRepository struct {
	mock.Mock
}

func (m *MockFranchiseRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockFranchiseRepository) SaveFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	args := m.Called(ctx, tx, franchise)
	return args.Error(0)
}

func (m *MockFranchiseRepository) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	args := m.Called(ctx, tx, franchise)
	return args.Error(0)
}

func (m *MockFranchiseRepository) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	args := m.Called(ctx, tx, franchiseID)
	return args.Error(0)
}

func (m *MockFranchiseRepository) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	args := m.Called(ctx, franchiseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseRepository) GetFranchiseByName(ctx context.Context, name string) (*domain.Franchise, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseRepository) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	args := m.Called(ctx, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Franchise), args.Error(1)
}

func (m *MockFranchiseRepository) CountBranchesByFranchise(ctx context.Context, franchiseID string) (int, error) {
	args := m.Called(ctx, franchiseID)
	return args.Int(0), args.Error(1)
}

func (m *MockFranchiseRepository) AssociateBranchesToFranchise(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	args := m.Called(ctx, tx, franchiseID, branchIDs)
	return args.Error(0)
}

func (m *MockFranchiseRepository) DissociateBranchesFromFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	args := m.Called(ctx, tx, franchiseID)
	return args.Error(0)
}

func (m *MockFranchiseRepository) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}
