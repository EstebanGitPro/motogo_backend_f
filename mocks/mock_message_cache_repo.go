package mocks

import (
	"context"

	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/stretchr/testify/mock"
)

// MockMessageCacheRepo is a mock implementation of types.MessageCacheRepository
type MockMessageCacheRepo struct {
	mock.Mock
}

func (m *MockMessageCacheRepo) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cachetypes.CachedMessage), args.Error(1)
}

func (m *MockMessageCacheRepo) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}

func (m *MockMessageCacheRepo) GetByCodeIncludingInactive(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}
