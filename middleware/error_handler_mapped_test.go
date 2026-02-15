package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
)

// ============================================
// Handle() Tests - Mapped Domain Errors
// ============================================

// createErrorHandlerCache creates a cache with the error codes used in the error handler map
func createErrorHandlerCache() *messagingCache.MessageCache {
	msgs := []cachetypes.CachedMessage{
		// User errors
		{ID: "1", Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Title: "Error", Content: "Usuario duplicado", Active: true},
		{ID: "2", Code: domain.MsgUserNotFound, Type: cachetypes.TypeError, Title: "Error", Content: "Usuario no encontrado", Active: true},
		// Validation errors
		{ID: "3", Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Title: "Error", Content: "Formato JSON inválido", Active: true},
		{ID: "4", Code: domain.MsgValInvalidReq, Type: cachetypes.TypeError, Title: "Error", Content: "Solicitud inválida", Active: true},
		{ID: "5", Code: domain.MsgValFieldRequired, Type: cachetypes.TypeError, Title: "Error", Content: "Campo ${0} requerido", Active: true},
		{ID: "6", Code: domain.MsgValMultiple, Type: cachetypes.TypeError, Title: "Error", Content: "Errores en campos: ${0}", Active: true},
		// Motorcycle errors
		{ID: "7", Code: domain.MsgMotorcycleNotFound, Type: cachetypes.TypeError, Title: "Error", Content: "Motocicleta no encontrada", Active: true},
		// General errors
		{ID: "8", Code: domain.MsgServerError, Type: cachetypes.TypeError, Title: "Error", Content: "Error interno del servidor", Active: true},
		// Auth errors
		{ID: "9", Code: domain.MsgUnauthorized, Type: cachetypes.TypeError, Title: "Error", Content: "No autorizado", Active: true},
	}

	repo := new(mockMsgRepo)
	repo.On("GetAllActiveForCache", mock.Anything).Return(msgs, nil)
	cache := messagingCache.NewMessageCache(repo, time.Minute)
	_ = cache.LoadMessages(context.Background())
	return cache
}

func TestHandle_MappedDomainError_DuplicateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrDuplicateUser)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgUserDuplicate, resp.Code)
	assert.Equal(t, "Usuario duplicado", resp.Message)
}

func TestHandle_MappedDomainError_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrInvalidJSONFormat)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgValJSONInvalid, resp.Code)
}

func TestHandle_MappedDomainError_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.GET("/motorcycles/:id", func(c *gin.Context) {
		_ = c.Error(domain.ErrMotorcycleNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/123", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgMotorcycleNotFound, resp.Code)
}

func TestHandle_MappedDomainError_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.GET("/protected", func(c *gin.Context) {
		_ = c.Error(domain.ErrInvalidToken)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgUnauthorized, resp.Code)
}

func TestHandle_MappedDomainError_InternalServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrInternalServer)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgServerError, resp.Code)
}

// ============================================
// Handle() Tests - Validation Fields
// ============================================

func TestHandle_WithValidationFields_SingleField(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		c.Set("validation_fields", []string{"email"})
		_ = c.Error(domain.ErrSchemaFieldRequired)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgValFieldRequired, resp.Code)
	assert.Contains(t, resp.Message, "email")
	assert.Equal(t, []string{"email"}, resp.Fields)
}

func TestHandle_WithValidationFields_MultipleFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		c.Set("validation_fields", []string{"email", "password", "name"})
		_ = c.Error(domain.ErrSchemaMultipleFields)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, resp.Success)
	assert.Equal(t, domain.MsgValMultiple, resp.Code)
	assert.Contains(t, resp.Message, "email, password, name")
	assert.Equal(t, []string{"email", "password", "name"}, resp.Fields)
}

func TestHandle_WithValidationFields_EmptyFieldsList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		c.Set("validation_fields", []string{})
		_ = c.Error(domain.ErrInvalidRequest)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, resp.Success)
}

func TestHandle_WithValidationFields_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		// Set a non-[]string value; type assertion should fail silently
		c.Set("validation_fields", "not-a-slice")
		_ = c.Error(domain.ErrInvalidRequest)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, resp.Success)
}

// ============================================
// Handle() Tests - Multiple errors (only last processed)
// ============================================

func TestHandle_MultipleErrors_LastUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := createErrorHandlerCache()
	handler := NewErrorHandler(cache)

	router := gin.New()
	router.Use(handler.Handle())
	router.POST("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrInvalidJSONFormat)
		_ = c.Error(domain.ErrMotorcycleNotFound) // last error
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should use the last error (motorcycle not found)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, domain.MsgMotorcycleNotFound, resp.Code)
}
