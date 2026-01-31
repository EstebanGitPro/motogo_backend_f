package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ResponseHandler Tests
// ============================================

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewResponseHandler(t *testing.T) {
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	assert.NotNil(t, handler)
	assert.Equal(t, cache, handler.cache)
}

func TestResponseHandler_DataOnly(t *testing.T) {
	router := setupTestRouter()
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	type testData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	router.GET("/test", func(c *gin.Context) {
		handler.DataOnly(c, testData{Name: "Test", Age: 25})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"name":"Test"`)
	assert.Contains(t, w.Body.String(), `"age":25`)
}

func TestResponseHandler_DataOnly_NilData(t *testing.T) {
	router := setupTestRouter()
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	router.GET("/test", func(c *gin.Context) {
		handler.DataOnly(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestResponseHandler_DataOnly_EmptySlice(t *testing.T) {
	router := setupTestRouter()
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	router.GET("/test", func(c *gin.Context) {
		handler.DataOnly(c, []string{})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"data":[]`)
}

// ============================================
// ErrorHandler Tests
// ============================================

func TestNewErrorHandler(t *testing.T) {
	cache := &messagingCache.MessageCache{}
	handler := NewErrorHandler(cache)

	assert.NotNil(t, handler)
	assert.Equal(t, cache, handler.cache)
}

// ============================================
// APIResponse Structure Tests
// ============================================

func TestAPIResponse_JSONMarshalling(t *testing.T) {
	router := setupTestRouter()
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	router.GET("/test", func(c *gin.Context) {
		handler.DataOnly(c, map[string]interface{}{
			"id":   "test-123",
			"name": "Test Item",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"test-123"`)
	assert.Contains(t, w.Body.String(), `"name":"Test Item"`)
}

func TestAPIResponse_NestedData(t *testing.T) {
	router := setupTestRouter()
	cache := &messagingCache.MessageCache{}
	handler := NewResponseHandler(cache)

	router.GET("/test", func(c *gin.Context) {
		handler.DataOnly(c, map[string]interface{}{
			"user": map[string]interface{}{
				"id":    "user-123",
				"email": "test@example.com",
			},
			"items": []string{"item1", "item2"},
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"user"`)
	assert.Contains(t, w.Body.String(), `"id":"user-123"`)
	assert.Contains(t, w.Body.String(), `"items"`)
}

// ============================================
// ErrorResponse Structure Tests
// ============================================

func TestErrorResponse_Structure(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Code:    "ERR_TEST",
		Message: "Test error message",
	}

	assert.False(t, resp.Success)
	assert.Equal(t, "ERR_TEST", resp.Code)
	assert.Equal(t, "Test error message", resp.Message)
}

// ============================================
// ResponseHandler Method Tests (Fallback Paths)
// Note: Full method tests require initialized MessageCache
// ============================================

func TestNewResponseHandler_WithNilCache(t *testing.T) {
	handler := NewResponseHandler(nil)

	assert.NotNil(t, handler)
	assert.Nil(t, handler.cache)
}

func TestAPIResponse_Success_Structure(t *testing.T) {
	resp := APIResponse{
		Success: true,
		Code:    "MSG_SUCCESS",
		Message: "Operation completed successfully",
		Data:    map[string]string{"key": "value"},
	}

	assert.True(t, resp.Success)
	assert.Equal(t, "MSG_SUCCESS", resp.Code)
	assert.Equal(t, "Operation completed successfully", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestAPIResponse_Error_Structure(t *testing.T) {
	resp := APIResponse{
		Success: false,
		Code:    "ERR_VALIDATION",
		Message: "Validation error occurred",
		Data:    nil,
	}

	assert.False(t, resp.Success)
	assert.Equal(t, "ERR_VALIDATION", resp.Code)
	assert.Nil(t, resp.Data)
}

func TestAPIResponse_WithNestedData(t *testing.T) {
	resp := APIResponse{
		Success: true,
		Code:    "MSG_DATA",
		Message: "Data retrieved",
		Data: map[string]interface{}{
			"user": map[string]string{
				"id":   "user-123",
				"name": "John",
			},
			"count": 5,
		},
	}

	assert.True(t, resp.Success)
	data := resp.Data.(map[string]interface{})
	assert.NotNil(t, data["user"])
	assert.Equal(t, 5, data["count"])
}
