package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockMessageRepository) SaveMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockMessageRepository) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByCode(ctx context.Context, code string) (*domain.Message, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	args := m.Called(ctx, msgType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	args := m.Called(ctx, module)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

// Tests

func TestGetMessageByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	service := services.NewMessageService(mockRepo)

	expectedMessage := &domain.Message{
		ID:      "msg-123",
		Code:    "TEST_CODE_001",
		Title:   "Test Message",
		Content: "Test content",
		Active:  true,
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, "msg-123").Return(expectedMessage, nil)

	// Act
	message, err := service.GetMessageByID(ctx, "msg-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, expectedMessage.ID, message.ID)
	assert.Equal(t, expectedMessage.Code, message.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetMessageByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	service := services.NewMessageService(mockRepo)

	// Mock expectations
	mockRepo.On("GetByID", ctx, "not-found").Return(nil, nil)

	// Act
	message, err := service.GetMessageByID(ctx, "not-found")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMessageNotFound, err)
	assert.Nil(t, message)

	mockRepo.AssertExpectations(t)
}

func TestGetMessageByCode_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	service := services.NewMessageService(mockRepo)

	expectedMessage := &domain.Message{
		ID:    "msg-123",
		Code:  "TEST_CODE_001",
		Title: "Test Message",
	}

	// Mock expectations
	mockRepo.On("GetByCode", ctx, "TEST_CODE_001").Return(expectedMessage, nil)

	// Act
	message, err := service.GetMessageByCode(ctx, "TEST_CODE_001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, expectedMessage.Code, message.Code)

	mockRepo.AssertExpectations(t)
}

func TestListActiveMessages_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	service := services.NewMessageService(mockRepo)

	expectedMessages := []domain.Message{
		{ID: "msg-1", Code: "CODE_001", Active: true},
		{ID: "msg-2", Code: "CODE_002", Active: true},
	}

	// Mock expectations
	mockRepo.On("GetAllActive", ctx).Return(expectedMessages, nil)

	// Act
	messages, err := service.ListActiveMessages(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, expectedMessages[0].Code, messages[0].Code)

	mockRepo.AssertExpectations(t)
}

func TestListActiveMessages_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	service := services.NewMessageService(mockRepo)

	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetAllActive", ctx).Return(nil, dbError)

	// Act
	messages, err := service.ListActiveMessages(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
}

func TestSaveMessageToDB_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewMessageService(mockRepo)

	message := domain.Message{
		ID:      "msg-123",
		Code:    "TEST_CODE_001",
		Title:   "Test",
		Content: "Test content",
		Active:  true,
	}

	// Mock expectations
	mockRepo.On("SaveMessage", ctx, mockTx, message).Return(nil)

	// Act
	err := service.SaveMessageToDB(ctx, mockTx, message)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSaveMessageToDB_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewMessageService(mockRepo)

	message := domain.Message{
		ID:   "msg-123",
		Code: "TEST_CODE_001",
	}

	dbError := errors.New("save failed")

	// Mock expectations
	mockRepo.On("SaveMessage", ctx, mockTx, message).Return(dbError)

	// Act
	err := service.SaveMessageToDB(ctx, mockTx, message)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMessageInDB_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewMessageService(mockRepo)

	message := domain.Message{
		ID:      "msg-123",
		Code:    "TEST_CODE_001",
		Title:   "Updated",
		Content: "Updated content",
	}

	existingMessage := &domain.Message{
		ID:   "msg-123",
		Code: "TEST_CODE_001",
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, "msg-123").Return(existingMessage, nil)
	mockRepo.On("UpdateMessage", ctx, mockTx, message).Return(nil)

	// Act
	err := service.UpdateMessageInDB(ctx, mockTx, message)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteMessageFromDB_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewMessageService(mockRepo)

	// Mock expectations
	mockRepo.On("DeleteMessage", ctx, mockTx, "msg-123").Return(nil)

	// Act
	err := service.DeleteMessageFromDB(ctx, mockTx, "msg-123")

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteMessageFromDB_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockMessageRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewMessageService(mockRepo)

	dbError := errors.New("delete failed")

	// Mock expectations
	mockRepo.On("DeleteMessage", ctx, mockTx, "msg-123").Return(dbError)

	// Act
	err := service.DeleteMessageFromDB(ctx, mockTx, "msg-123")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
}
