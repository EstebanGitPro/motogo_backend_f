package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Mock MessageRepository
// ============================================

type mockMessageRepo struct{ mock.Mock }

func (m *mockMessageRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *mockMessageRepo) SaveMessage(ctx context.Context, tx output.Tx, msg domain.Message) error {
	return m.Called(ctx, tx, msg).Error(0)
}
func (m *mockMessageRepo) UpdateMessage(ctx context.Context, tx output.Tx, msg domain.Message) error {
	return m.Called(ctx, tx, msg).Error(0)
}
func (m *mockMessageRepo) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	return m.Called(ctx, tx, id).Error(0)
}
func (m *mockMessageRepo) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}
func (m *mockMessageRepo) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}
func (m *mockMessageRepo) GetByCode(ctx context.Context, code string) (*domain.Message, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}
func (m *mockMessageRepo) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	args := m.Called(ctx, msgType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}
func (m *mockMessageRepo) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	args := m.Called(ctx, module)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

type mockMsgTx struct{ mock.Mock }

func (m *mockMsgTx) Commit() error   { return m.Called().Error(0) }
func (m *mockMsgTx) Rollback() error { return m.Called().Error(0) }

// ============================================
// BeginTx Tests
// ============================================

func TestMessageSvc_BeginTx_Success(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	repo.On("BeginTx", mock.Anything).Return(tx, nil)

	result, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestMessageSvc_BeginTx_Error(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	repo.On("BeginTx", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.BeginTx(context.Background())
	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// ValidateMessage Tests
// ============================================

func TestMessageSvc_ValidateMessage_NewNoDuplicate(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	msg := domain.Message{Code: "TEST_CODE"} // no ID = new
	repo.On("GetByCode", mock.Anything, "TEST_CODE").Return(nil, errors.New("not found"))

	err := svc.ValidateMessage(context.Background(), msg)
	assert.NoError(t, err)
}

func TestMessageSvc_ValidateMessage_NewDuplicate(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	existing := &domain.Message{ID: "existing-1", Code: "TEST_CODE"}
	msg := domain.Message{Code: "TEST_CODE"} // no ID = new
	repo.On("GetByCode", mock.Anything, "TEST_CODE").Return(existing, nil)

	err := svc.ValidateMessage(context.Background(), msg)
	assert.Equal(t, domain.ErrMessageCodeDuplicate, err)
}

func TestMessageSvc_ValidateMessage_UpdateSkipsDuplicateCheck(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	msg := domain.Message{ID: "existing-1", Code: "TEST_CODE"} // has ID = update

	err := svc.ValidateMessage(context.Background(), msg)
	assert.NoError(t, err)
	// GetByCode should NOT be called for updates
	repo.AssertNotCalled(t, "GetByCode")
}

// ============================================
// ListMessages Tests
// ============================================

func TestMessageSvc_ListMessages_Success(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	expected := []domain.Message{
		{ID: "m1", Code: "CODE_1"},
		{ID: "m2", Code: "CODE_2"},
	}
	repo.On("GetAllActive", mock.Anything).Return(expected, nil)

	result, err := svc.ListMessages(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMessageSvc_ListMessages_Error(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)

	repo.On("GetAllActive", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.ListMessages(context.Background(), nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// SaveMessageToDB Tests
// ============================================

func TestMessageSvc_SaveToDB_Success(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	msg := domain.Message{ID: "m1", Code: "CODE"}
	repo.On("SaveMessage", mock.Anything, tx, msg).Return(nil)

	err := svc.SaveMessageToDB(context.Background(), tx, msg)
	assert.NoError(t, err)
}

func TestMessageSvc_SaveToDB_Error(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	msg := domain.Message{ID: "m1", Code: "CODE"}
	repo.On("SaveMessage", mock.Anything, tx, msg).Return(errors.New("save error"))

	err := svc.SaveMessageToDB(context.Background(), tx, msg)
	assert.Error(t, err)
}

// ============================================
// DeleteMessageFromDB Tests
// ============================================

func TestMessageSvc_DeleteFromDB_Success(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	repo.On("DeleteMessage", mock.Anything, tx, "m1").Return(nil)

	err := svc.DeleteMessageFromDB(context.Background(), tx, "m1")
	assert.NoError(t, err)
}

func TestMessageSvc_DeleteFromDB_Error(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	repo.On("DeleteMessage", mock.Anything, tx, "m1").Return(errors.New("delete error"))

	err := svc.DeleteMessageFromDB(context.Background(), tx, "m1")
	assert.Error(t, err)
}

// ============================================
// UpdateMessageInDB Tests
// ============================================

func TestMessageSvc_UpdateInDB_Success(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	existing := &domain.Message{ID: "m1", Code: "CODE"}
	msg := domain.Message{ID: "m1", Code: "CODE"}
	repo.On("GetByID", mock.Anything, "m1").Return(existing, nil)
	repo.On("UpdateMessage", mock.Anything, tx, msg).Return(nil)

	err := svc.UpdateMessageInDB(context.Background(), tx, msg)
	assert.NoError(t, err)
}

func TestMessageSvc_UpdateInDB_GetByIDError(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	msg := domain.Message{ID: "m1", Code: "CODE"}
	repo.On("GetByID", mock.Anything, "m1").Return(nil, errors.New("db error"))

	err := svc.UpdateMessageInDB(context.Background(), tx, msg)
	assert.Equal(t, domain.ErrMessageCannotUpdate, err)
}

func TestMessageSvc_UpdateInDB_CodeMismatch(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	existing := &domain.Message{ID: "m1", Code: "OLD_CODE"}
	msg := domain.Message{ID: "m1", Code: "NEW_CODE"}
	repo.On("GetByID", mock.Anything, "m1").Return(existing, nil)

	err := svc.UpdateMessageInDB(context.Background(), tx, msg)
	assert.Equal(t, domain.ErrMessageNotFound, err)
}

func TestMessageSvc_UpdateInDB_NilExisting(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := services.NewMessageService(repo)
	tx := new(mockMsgTx)

	msg := domain.Message{ID: "m1", Code: "CODE"}
	repo.On("GetByID", mock.Anything, "m1").Return((*domain.Message)(nil), nil)

	err := svc.UpdateMessageInDB(context.Background(), tx, msg)
	assert.Equal(t, domain.ErrMessageNotFound, err)
}
