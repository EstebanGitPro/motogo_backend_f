package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockBranchService is a mock implementation of input.BranchService
// Only implements methods needed for FranchiseInteractor testing
type MockBranchService struct {
	mock.Mock
}

func (m *MockBranchService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockBranchService) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchService) RegisterBranch(ctx context.Context, tx output.Tx, branch domain.Branch) (*domain.Branch, error) {
	args := m.Called(ctx, tx, branch)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchService) ValidateBrands(ctx context.Context, brands []string) error {
	args := m.Called(ctx, brands)
	return args.Error(0)
}

func (m *MockBranchService) GeocodeLocation(ctx context.Context, location *domain.Location) (bool, error) {
	args := m.Called(ctx, location)
	return args.Bool(0), args.Error(1)
}

func (m *MockBranchService) SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	args := m.Called(ctx, tx, location)
	return args.Error(0)
}

func (m *MockBranchService) SaveBranchBrands(ctx context.Context, tx output.Tx, branchID string, brands []string) error {
	args := m.Called(ctx, tx, branchID, brands)
	return args.Error(0)
}

func (m *MockBranchService) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	args := m.Called(ctx, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Branch), args.Error(1)
}

func (m *MockBranchService) UpdateBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	args := m.Called(ctx, tx, branch)
	return args.Error(0)
}

func (m *MockBranchService) DeleteBranch(ctx context.Context, tx output.Tx, branchID string) error {
	args := m.Called(ctx, tx, branchID)
	return args.Error(0)
}
