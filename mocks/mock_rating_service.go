package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockRatingService is a mock implementation of input.RatingService
type MockRatingService struct {
	mock.Mock
}

func (m *MockRatingService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockRatingService) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string) error {
	args := m.Called(ctx, tx, itemID, rating, comment)
	return args.Error(0)
}

func (m *MockRatingService) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedServiceItem), args.Error(1)
}

func (m *MockRatingService) GetCompletedServiceByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedService), args.Error(1)
}

func (m *MockRatingService) GetReviewsByServiceID(ctx context.Context, serviceID string) (*domain.ServiceReviewSummary, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ServiceReviewSummary), args.Error(1)
}
