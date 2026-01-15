package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockFranchiseInteractor is a mock implementation of FranchiseInteractor for handler testing
type MockFranchiseInteractor struct {
	mock.Mock
}

func (m *MockFranchiseInteractor) CreateFranchiseWithBranches(ctx context.Context, franchise domain.Franchise, branchIDs []string, representativeID string) (*domain.Franchise, error) {
	args := m.Called(ctx, franchise, branchIDs, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseInteractor) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	args := m.Called(ctx, franchiseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Franchise), args.Error(1)
}

func (m *MockFranchiseInteractor) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	args := m.Called(ctx, representativeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Franchise), args.Error(1)
}

func (m *MockFranchiseInteractor) UpdateFranchise(ctx context.Context, franchise domain.Franchise, representativeID string) error {
	args := m.Called(ctx, franchise, representativeID)
	return args.Error(0)
}

func (m *MockFranchiseInteractor) DeleteFranchise(ctx context.Context, franchiseID, representativeID string) error {
	args := m.Called(ctx, franchiseID, representativeID)
	return args.Error(0)
}

func (m *MockFranchiseInteractor) AddBranchToFranchise(ctx context.Context, franchiseID, branchID, representativeID string) error {
	args := m.Called(ctx, franchiseID, branchID, representativeID)
	return args.Error(0)
}

func (m *MockFranchiseInteractor) RemoveBranchFromFranchise(ctx context.Context, franchiseID, branchID, representativeID string) error {
	args := m.Called(ctx, franchiseID, branchID, representativeID)
	return args.Error(0)
}
