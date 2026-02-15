package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
)

// ============================================
// Test helpers – mock repo + pre-loaded cache
// ============================================

type mockMsgRepo struct {
	mock.Mock
}

func (m *mockMsgRepo) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cachetypes.CachedMessage), args.Error(1)
}

func (m *mockMsgRepo) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}

func (m *mockMsgRepo) GetByCodeIncludingInactive(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cachetypes.CachedMessage), args.Error(1)
}

// createLoadedCache creates a MessageCache pre-loaded with the given messages.
// It also sets up wildcard mocks for DB fallback calls to avoid panics on unknown codes.
func createLoadedCache(msgs []cachetypes.CachedMessage) *messagingCache.MessageCache {
	repo := new(mockMsgRepo)
	repo.On("GetAllActiveForCache", mock.Anything).Return(msgs, nil)
	// Wildcard mocks for DB fallback – return nil (not found) for any code
	repo.On("GetByCodeForCache", mock.Anything, mock.Anything).Return(nil, nil)
	repo.On("GetByCodeIncludingInactive", mock.Anything, mock.Anything).Return(nil, nil)
	cache := messagingCache.NewMessageCache(repo, time.Minute)
	_ = cache.LoadMessages(context.Background())
	return cache
}

// testMessages returns a set of cached messages covering error/success/warning/info types
func testMessages() []cachetypes.CachedMessage {
	return []cachetypes.CachedMessage{
		{ID: "1", Code: "ERR_TEST_001", Type: cachetypes.TypeError, Title: "Error", Content: "Something went wrong", Active: true},
		{ID: "2", Code: "SUC_TEST_001", Type: cachetypes.TypeSuccess, Title: "Success", Content: "Operation completed", Active: true},
		{ID: "3", Code: "WARN_TEST_001", Type: cachetypes.TypeWarning, Title: "Warning", Content: "Please be careful", Active: true},
		{ID: "4", Code: "INFO_TEST_001", Type: cachetypes.TypeInfo, Title: "Info", Content: "Here is some information", Active: true},
		{ID: "5", Code: "ERR_PARAM_001", Type: cachetypes.TypeError, Title: "Error", Content: "Field ${0} is invalid", Active: true},
	}
}

// ============================================
// ResponseHandler.Error Tests
// ============================================

func TestResponseHandler_Error_WithCache_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Error(c, "ERR_TEST_001")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// TypeError falls back to 500 via cache.GetHTTPStatus since it's not in the status map
	assert.False(t, resp.Success)
	assert.Equal(t, "ERR_TEST_001", resp.Code)
	assert.Equal(t, "Something went wrong", resp.Message)
}

func TestResponseHandler_Error_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Error(c, "NONEXISTENT_CODE")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"message":"Unknown error"`)
}

func TestResponseHandler_Error_WithParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Error(c, "ERR_PARAM_001", "email")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "Field email is invalid", resp.Message)
}

// ============================================
// ResponseHandler.ErrorWithData Tests
// ============================================

func TestResponseHandler_ErrorWithData_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.ErrorWithData(c, "ERR_TEST_001", map[string]string{"field": "email"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "ERR_TEST_001", resp.Code)
	assert.Equal(t, "Something went wrong", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestResponseHandler_ErrorWithData_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.ErrorWithData(c, "NONEXISTENT_CODE", map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Unknown error"`)
	assert.Contains(t, w.Body.String(), `"data"`)
}

// ============================================
// ResponseHandler.Success Tests
// ============================================

func TestResponseHandler_Success_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Success(c, "SUC_TEST_001")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.True(t, resp.Success) // TypeSuccess != "ERROR" → Success = true
	assert.Equal(t, "SUC_TEST_001", resp.Code)
	assert.Equal(t, "Operation completed", resp.Message)
}

func TestResponseHandler_Success_ErrorType_SetsFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	// Call Success with an ERROR-type message code
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Success(c, "ERR_TEST_001")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// msg.Type == "ERROR" → isSuccess = false
	assert.False(t, resp.Success)
}

func TestResponseHandler_Success_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Success(c, "NONEXISTENT_CODE")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Operation successful"`)
}

// ============================================
// ResponseHandler.SuccessWithData Tests
// ============================================

func TestResponseHandler_SuccessWithData_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.SuccessWithData(c, "SUC_TEST_001", map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Equal(t, "SUC_TEST_001", resp.Code)
	assert.NotNil(t, resp.Data)
}

func TestResponseHandler_SuccessWithData_ErrorType_SetsFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.SuccessWithData(c, "ERR_TEST_001", []string{"data"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.False(t, resp.Success)
}

func TestResponseHandler_SuccessWithData_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.SuccessWithData(c, "NONEXISTENT_CODE", "some data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Operation successful"`)
	assert.Contains(t, w.Body.String(), `"data"`)
}

// ============================================
// ResponseHandler.Warning Tests
// ============================================

func TestResponseHandler_Warning_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Warning(c, "WARN_TEST_001")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, resp.Success)
	assert.Equal(t, "WARN_TEST_001", resp.Code)
	assert.Equal(t, "Please be careful", resp.Message)
}

func TestResponseHandler_Warning_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Warning(c, "NONEXISTENT_CODE")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"System warning"`)
}

// ============================================
// ResponseHandler.WarningWithData Tests
// ============================================

func TestResponseHandler_WarningWithData_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.WarningWithData(c, "WARN_TEST_001", map[string]int{"count": 5})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, resp.Success)
	assert.Equal(t, "WARN_TEST_001", resp.Code)
	assert.NotNil(t, resp.Data)
}

func TestResponseHandler_WarningWithData_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.WarningWithData(c, "NONEXISTENT_CODE", []int{1, 2, 3})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"System warning"`)
	assert.Contains(t, w.Body.String(), `"data"`)
}

// ============================================
// ResponseHandler.Info Tests
// ============================================

func TestResponseHandler_Info_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Info(c, "INFO_TEST_001")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, resp.Success)
	assert.Equal(t, "INFO_TEST_001", resp.Code)
	assert.Equal(t, "Here is some information", resp.Message)
}

func TestResponseHandler_Info_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.Info(c, "NONEXISTENT_CODE")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"System information"`)
}

// ============================================
// ResponseHandler.InfoWithData Tests
// ============================================

func TestResponseHandler_InfoWithData_KnownCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.InfoWithData(c, "INFO_TEST_001", map[string]string{"version": "1.0"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, resp.Success)
	assert.Equal(t, "INFO_TEST_001", resp.Code)
	assert.NotNil(t, resp.Data)
}

func TestResponseHandler_InfoWithData_UnknownCode_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := createLoadedCache(testMessages())
	handler := NewResponseHandler(cache)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handler.InfoWithData(c, "NONEXISTENT_CODE", "info data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"System information"`)
	assert.Contains(t, w.Body.String(), `"data"`)
}
