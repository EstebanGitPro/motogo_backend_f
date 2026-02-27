package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockRatingRepository is a mock implementation of output.RatingRepository
type MockRatingRepository struct {
	mock.Mock
}

func (m *MockRatingRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockRatingRepository) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string, isOffensive bool) error {
	args := m.Called(ctx, tx, itemID, rating, comment, isOffensive)
	return args.Error(0)
}

func (m *MockRatingRepository) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedServiceItem), args.Error(1)
}
