package interactor_test

import (
	"context"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for MessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageService) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

// Tests

func TestCreateMessage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	message := domain.Message{
		Code:    "TEST_001",
		Title:   "Test",
		Content: "Test content",
	}

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()

	mockService.On("ValidateMessage", ctx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SaveMessageToDB", ctx, mockTx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	result, err := messageInteractor.CreateMessage(ctx, message)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, message.Code, result.Code)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestCreateMessage_ValidationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	message := domain.Message{
		// Missing required Code
		Title:   "Test",
		Content: "Test content",
	}

	validationError := domain.ErrMessageCodeRequired

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	mockService.On("ValidateMessage", ctx, mock.AnythingOfType("domain.Message")).Return(validationError)

	// Act
	result, err := messageInteractor.CreateMessage(ctx, message)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, validationError, err)
	assert.Nil(t, result)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetMessageByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	expectedMessage := &domain.Message{
		ID:    "msg-123",
		Code:  "TEST_001",
		Title: "Test",
	}

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockService.On("GetMessageByID", ctx, "msg-123").Return(expectedMessage, nil)

	// Act
	message, err := messageInteractor.GetMessageByID(ctx, "msg-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, expectedMessage.ID, message.ID)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestListActiveMessages_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	expectedMessages := []domain.Message{
		{ID: "msg-1", Code: "CODE_001"},
		{ID: "msg-2", Code: "CODE_002"},
	}

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockService.On("ListActiveMessages", ctx).Return(expectedMessages, nil)

	// Act
	messages, err := messageInteractor.ListActiveMessages(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteMessage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	existingMessage := &domain.Message{
		ID:   "msg-123",
		Code: "TEST_001",
	}

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()

	mockService.On("GetMessageByID", ctx, "msg-123").Return(existingMessage, nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("DeleteMessageFromDB", ctx, mockTx, "msg-123").Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := messageInteractor.DeleteMessage(ctx, "msg-123")

	// Assert
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteMessage_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	notFoundError := domain.ErrMessageNotFound

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	mockService.On("GetMessageByID", ctx, "not-found").Return(nil, notFoundError)

	// Act
	err := messageInteractor.DeleteMessage(ctx, "not-found")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, notFoundError, err)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateMessage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(MockMessageService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockService, mockLogger)

	message := domain.Message{
		ID:      "msg-123",
		Code:    "TEST_001",
		Title:   "Updated",
		Content: "Updated content",
	}

	existingMessage := &domain.Message{
		ID:   "msg-123",
		Code: "TEST_001",
	}

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()

	mockService.On("GetMessageByID", ctx, message.ID).Return(existingMessage, nil)
	mockService.On("ValidateMessage", ctx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdateMessageInDB", ctx, mockTx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockTx.On("Commit").Return(nil)

	//Act
	result, err := messageInteractor.UpdateMessage(ctx, message)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, message.ID, result.ID)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
