package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ErrorResponse Tests
// ============================================

func TestErrorResponse_StructFields(t *testing.T) {
	response := ErrorResponse{
		Success: false,
		Code:    "ERR_TEST",
		Message: "Test message",
	}

	assert.False(t, response.Success)
	assert.Equal(t, "ERR_TEST", response.Code)
	assert.Equal(t, "Test message", response.Message)
}

// ============================================
// NewErrorHandler Tests
// ============================================

func TestNewErrorHandler_ReturnsHandler(t *testing.T) {
	handler := NewErrorHandler(nil)

	assert.NotNil(t, handler)
	assert.Nil(t, handler.cache)
}

func TestNewErrorHandler_WithNilCache(t *testing.T) {
	handler := NewErrorHandler(nil)

	assert.NotNil(t, handler)
}

// ============================================
// errorToMessageCode Mapping Tests
// ============================================

func TestErrorToMessageCode_ContainsUserErrors(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{domain.ErrDuplicateUser, domain.MsgUserDuplicate},
		{domain.ErrUserNotFound, domain.MsgUserNotFound},
		{domain.ErrUserCannotSave, domain.MsgUserCannotSave},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			code, exists := errorToMessageCode[tc.err]
			assert.True(t, exists, "Error should be mapped: %v", tc.err)
			assert.Equal(t, tc.expected, code)
		})
	}
}

func TestErrorToMessageCode_ContainsValidationErrors(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{domain.ErrInvalidJSONFormat, domain.MsgValJSONInvalid},
		{domain.ErrInvalidRequest, domain.MsgValInvalidReq},
		{domain.ErrInvalidID, domain.MsgValIDInvalid},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			code, exists := errorToMessageCode[tc.err]
			assert.True(t, exists, "Error should be mapped: %v", tc.err)
			assert.Equal(t, tc.expected, code)
		})
	}
}

func TestErrorToMessageCode_ContainsSchemaErrors(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{domain.ErrSchemaBadRequest, domain.MsgValBadFormat},
		{domain.ErrSchemaValidationFailed, domain.MsgValFailed},
		{domain.ErrSchemaFieldRequired, domain.MsgValFieldRequired},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			code, exists := errorToMessageCode[tc.err]
			assert.True(t, exists, "Error should be mapped: %v", tc.err)
			assert.Equal(t, tc.expected, code)
		})
	}
}

func TestErrorToMessageCode_ContainsMotorcycleErrors(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{domain.ErrMotorcycleNotFound, domain.MsgMotorcycleNotFound},
		{domain.ErrMotorcycleCannotSave, domain.MsgMotorcycleCannotSave},
		{domain.ErrDuplicateLicensePlate, domain.MsgDuplicateLicensePlate},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			code, exists := errorToMessageCode[tc.err]
			assert.True(t, exists, "Error should be mapped: %v", tc.err)
			assert.Equal(t, tc.expected, code)
		})
	}
}

func TestErrorToMessageCode_ContainsInternalServerError(t *testing.T) {
	code, exists := errorToMessageCode[domain.ErrInternalServer]
	assert.True(t, exists)
	assert.Equal(t, domain.MsgServerError, code)
}

func TestErrorToMessageCode_UnmappedError(t *testing.T) {
	unknownErr := errors.New("unknown error")
	_, exists := errorToMessageCode[unknownErr]
	assert.False(t, exists, "Unknown error should not be mapped")
}

// ============================================
// Handle() Tests - No Errors Case
// ============================================

func TestHandle_NoErrors_PassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewErrorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// Call the handler directly
	handlerFunc := handler.Handle()
	handlerFunc(c)

	// No errors means no response body modification
	assert.Equal(t, http.StatusOK, w.Code)
}

// ============================================
// Handle() Tests - With Errors Case (Fallback)
// ============================================

func TestHandle_UnmappedError_ReturnsFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewErrorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// Add an unmapped error
	unknownErr := errors.New("unknown error")
	_ = c.Error(unknownErr)

	// Call the handler
	handlerFunc := handler.Handle()
	handlerFunc(c)

	// Should return internal server error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), domain.MsgServerError)
}
