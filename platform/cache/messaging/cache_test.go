package messaging

import (
	"context"
	"errors"
	"testing"
	"time"

	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageCacheRepository is a mock implementation of MessageCacheRepository
type MockMessageCacheRepository struct {
	mock.Mock
}

func (m *MockMessageCacheRepository) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cachetypes.CachedMessage), args.Error(1)
}

func (m *MockMessageCacheRepository) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}

func (m *MockMessageCacheRepository) GetByCodeIncludingInactive(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}

// ============================================
// NewMessageCache Tests
// ============================================

func TestNewMessageCache_ReturnsInstance(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	cache := NewMessageCache(mockRepo, time.Minute)

	assert.NotNil(t, cache)
	assert.Equal(t, 0, cache.MessageCount())
}

// ============================================
// LoadMessages Tests
// ============================================

func TestLoadMessages_Success(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "ERR_001", Type: TypeError, Title: "Error 1", Content: "Content 1", Active: true},
		{ID: "2", Code: "SUC_001", Type: TypeSuccess, Title: "Success 1", Content: "Content 2", Active: true},
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	err := cache.LoadMessages(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 2, cache.MessageCount())
	mockRepo.AssertExpectations(t)
}

func TestLoadMessages_RepositoryError(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(nil, errors.New("database error"))

	cache := NewMessageCache(mockRepo, time.Minute)
	err := cache.LoadMessages(context.Background())

	assert.Error(t, err)
	assert.Equal(t, 0, cache.MessageCount())
	mockRepo.AssertExpectations(t)
}

func TestLoadMessages_EmptyMessages(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	err := cache.LoadMessages(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 0, cache.MessageCount())
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMessage Tests
// ============================================

func TestGetMessage_FromCache(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "ERR_001", Type: TypeError, Title: "Error", Content: "Error content", Active: true},
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	// Get from cache
	msg := cache.GetMessage("ERR_001")

	assert.NotNil(t, msg)
	assert.Equal(t, "ERR_001", msg.Code)
	assert.Equal(t, "Error", msg.Title)
}

func TestGetMessage_NotInCache_FetchesFromDB(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	dbMessage := &cachetypes.CachedMessage{
		ID: "2", Code: "SUC_002", Type: TypeSuccess, Title: "Success", Content: "Success content", Active: true,
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	mockRepo.On("GetByCodeForCache", mock.Anything, "SUC_002").Return(dbMessage, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	msg := cache.GetMessage("SUC_002")

	assert.NotNil(t, msg)
	assert.Equal(t, "SUC_002", msg.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetMessage_InactiveMessage_ReturnsFallback(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	// For this test we need GetByCodeForCache to return nil (not error) so flow goes to GetByCodeIncludingInactive
	inactiveMsg := &cachetypes.CachedMessage{
		ID: "2", Code: "INACTIVE_MSG", Type: TypeError, Title: "Inactive", Content: "Not available", Active: false,
	}
	fallbackMsg := &cachetypes.CachedMessage{
		ID: "3", Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Title: "Fallback", Content: "Message unavailable", Active: true,
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	// When GetByCodeForCache returns nil without error, GetMessage calls GetByCodeIncludingInactive
	mockRepo.On("GetByCodeForCache", mock.Anything, "INACTIVE_MSG").Return(nil, nil)
	// Check if it exists but is inactive - returns the inactive message
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "INACTIVE_MSG").Return(inactiveMsg, nil)
	// Since message is inactive, it calls GetMessage recursively for fallback
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(fallbackMsg, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	msg := cache.GetMessage("INACTIVE_MSG")

	assert.NotNil(t, msg)
	assert.Equal(t, "GEN_MSG_INACTIVE_ERR_00002", msg.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetMessage_FallbackMessage_NotFound(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	// When querying for the fallback message itself and GetByCodeForCache returns error:
	// The code checks if code == "GEN_MSG_INACTIVE_ERR_00002" and returns nil immediately
	// to avoid infinite recursion - it DOES NOT call GetByCodeIncludingInactive
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(nil, errors.New("not found"))

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	// When the fallback message itself isn't found, should return nil
	msg := cache.GetMessage("GEN_MSG_INACTIVE_ERR_00002")

	assert.Nil(t, msg)
	mockRepo.AssertExpectations(t)
}

func TestGetMessage_GetByCodeIncludingInactive_Error(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	fallbackMsg := &cachetypes.CachedMessage{
		ID: "1", Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Title: "Fallback", Content: "Message unavailable", Active: true,
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	mockRepo.On("GetByCodeForCache", mock.Anything, "SOME_CODE").Return(nil, nil)
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "SOME_CODE").Return(nil, errors.New("db error"))
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(fallbackMsg, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	msg := cache.GetMessage("SOME_CODE")

	assert.NotNil(t, msg)
	assert.Equal(t, "GEN_MSG_INACTIVE_ERR_00002", msg.Code)
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMessageResponse Tests
// ============================================

func TestGetMessageResponse_Success(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "MSG_001", Type: TypeInfo, Title: "Info", Content: "Hello world", Active: true},
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	response := cache.GetMessageResponse("MSG_001")

	assert.NotNil(t, response)
	assert.Equal(t, "MSG_001", response.Code)
	assert.Equal(t, "Hello world", response.Content)
	assert.Equal(t, TypeInfo, response.Type)
}

func TestGetMessageResponse_WithPlaceholders(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "GREET", Type: TypeSuccess, Title: "Greeting", Content: "Hello ${0}, welcome to ${1}!", Active: true},
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	response := cache.GetMessageResponse("GREET", "John", "MotoGo")

	assert.NotNil(t, response)
	assert.Equal(t, "Hello John, welcome to MotoGo!", response.Content)
}

func TestGetMessageResponse_NotFound(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	mockRepo.On("GetByCodeForCache", mock.Anything, "NONEXISTENT").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "NONEXISTENT").Return(nil, errors.New("not found"))
	// Fallback lookup
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(nil, errors.New("not found"))

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	response := cache.GetMessageResponse("NONEXISTENT")

	assert.Nil(t, response)
}

// ============================================
// GetHTTPStatus Tests
// ============================================

func TestGetHTTPStatus_KnownCode(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	testCases := []struct {
		code           string
		expectedStatus int
	}{
		{"MOD_U_DUP_ERR_00001", 409},            // Conflict
		{"MOD_V_VAL_ERR_00001", 400},            // BadRequest
		{"GEN_AUTH_ERR_00002", 401},             // Unauthorized
		{"GEN_FORBIDDEN_ERR_00003", 403},        // Forbidden
		{"MOD_P_NOT_FOUND_ERR_00001", 404},      // NotFound
		{"MOD_U_REG_EXI_00001", 201},            // Created
		{"GEN_OPE_EXI_00001", 200},              // OK
		{"MOD_INFRA_KC_UNAVAIL_ERR_00004", 423}, // Locked
	}

	for _, tc := range testCases {
		t.Run(tc.code, func(t *testing.T) {
			status := cache.GetHTTPStatus(tc.code)
			assert.Equal(t, tc.expectedStatus, status)
		})
	}
}

func TestGetHTTPStatus_UnknownCode_FallbackToMessageType_Success(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "CUSTOM_SUCCESS", Type: TypeSuccess, Title: "Success", Content: "It worked", Active: true},
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("CUSTOM_SUCCESS")
	assert.Equal(t, 200, status) // TypeSuccess should return 200
}

func TestGetHTTPStatus_UnknownCode_FallbackToMessageType_Error(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "CUSTOM_ERROR", Type: TypeError, Title: "Error", Content: "It failed", Active: true},
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("CUSTOM_ERROR")
	assert.Equal(t, 500, status) // TypeError should return 500
}

func TestGetHTTPStatus_UnknownCode_FallbackToMessageType_Warning(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "CUSTOM_WARNING", Type: TypeWarning, Title: "Warning", Content: "Be careful", Active: true},
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("CUSTOM_WARNING")
	assert.Equal(t, 200, status) // TypeWarning should return 200
}

func TestGetHTTPStatus_UnknownCode_FallbackToMessageType_Info(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "CUSTOM_INFO", Type: TypeInfo, Title: "Info", Content: "FYI", Active: true},
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("CUSTOM_INFO")
	assert.Equal(t, 200, status) // TypeInfo should return 200
}

func TestGetHTTPStatus_UnknownCode_FallbackToMessageType_Debug(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "CUSTOM_DEBUG", Type: TypeDebug, Title: "Debug", Content: "Details", Active: true},
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("CUSTOM_DEBUG")
	assert.Equal(t, 200, status) // TypeDebug should return 200
}

func TestGetHTTPStatus_UnknownCode_MessageNotFound_Returns500(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	mockRepo.On("GetByCodeForCache", mock.Anything, "TOTALLY_UNKNOWN").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "TOTALLY_UNKNOWN").Return(nil, errors.New("not found"))
	// Fallback lookup
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(nil, errors.New("not found"))

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("TOTALLY_UNKNOWN")
	assert.Equal(t, 500, status) // When message is nil, return 500
}

func TestGetHTTPStatus_FallbackMessage_UsesItsStatusFromMap(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	fallbackMsg := &cachetypes.CachedMessage{
		ID: "1", Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Title: "Inactive", Content: "Message not available", Active: true,
	}
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]cachetypes.CachedMessage{}, nil)
	mockRepo.On("GetByCodeForCache", mock.Anything, "SOME_CODE").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeIncludingInactive", mock.Anything, "SOME_CODE").Return(nil, errors.New("not found"))
	mockRepo.On("GetByCodeForCache", mock.Anything, "GEN_MSG_INACTIVE_ERR_00002").Return(fallbackMsg, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	status := cache.GetHTTPStatus("SOME_CODE")
	assert.Equal(t, 404, status) // GEN_MSG_INACTIVE_ERR_00002 maps to 404 in the map
}

// ============================================
// MessageCount Tests
// ============================================

func TestMessageCount_Empty(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	cache := NewMessageCache(mockRepo, time.Minute)

	assert.Equal(t, 0, cache.MessageCount())
}

func TestMessageCount_WithMessages(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	messages := []cachetypes.CachedMessage{
		{ID: "1", Code: "MSG_1", Type: TypeError, Active: true},
		{ID: "2", Code: "MSG_2", Type: TypeSuccess, Active: true},
		{ID: "3", Code: "MSG_3", Type: TypeInfo, Active: true},
	}

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return(messages, nil)

	cache := NewMessageCache(mockRepo, time.Minute)
	_ = cache.LoadMessages(context.Background())

	assert.Equal(t, 3, cache.MessageCount())
}

// ============================================
// replaceAll Tests
// ============================================

func TestReplaceAll_SingleReplacement(t *testing.T) {
	result := replaceAll("Hello ${0}", "${0}", "World")
	assert.Equal(t, "Hello World", result)
}

func TestReplaceAll_MultipleReplacements(t *testing.T) {
	result := replaceAll("${0} and ${0}", "${0}", "A")
	assert.Equal(t, "A and A", result)
}

func TestReplaceAll_NoMatch(t *testing.T) {
	result := replaceAll("Hello World", "${0}", "X")
	assert.Equal(t, "Hello World", result)
}

func TestReplaceAll_EmptyString(t *testing.T) {
	result := replaceAll("", "${0}", "X")
	assert.Equal(t, "", result)
}

// ============================================
// StopAutoRefresh Tests
// ============================================

func TestStopAutoRefresh_WithZeroInterval(t *testing.T) {
	mockRepo := new(MockMessageCacheRepository)
	cache := NewMessageCache(mockRepo, 0)

	// Should not panic
	cache.StopAutoRefresh()
}
