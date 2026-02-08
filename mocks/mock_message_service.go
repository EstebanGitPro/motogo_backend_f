package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockMessageService is a mock implementation of input.MessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockMessageService) ValidateMessage(ctx context.Context, message domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageService) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) GetMessageByCode(ctx context.Context, code string) (*domain.Message, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageService) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageService) SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func (m *MockMessageService) UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func (m *MockMessageService) DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}
