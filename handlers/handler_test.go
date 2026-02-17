package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestEncoder() *idencoder.HashidsEncoder {
	encoder, _ := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    "test-secret-key-for-testing",
		MinLength: 8,
	})
	return encoder
}

func TestNew_CreatesHandler(t *testing.T) {
	// Arrange
	mockCache := new(mocks.MockMessageCache)
	encoder := createTestEncoder()
	responseHandler := middleware.NewResponseHandler(nil)

	// Note: In a real test, we'd mock the interactors too
	// For now, we're testing that handler creation works with nil interactors

	// Act
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	// Assert
	assert.NotNil(t, handler)
	_ = mockCache // Silence unused warning
}

func TestEncodeID_Success(t *testing.T) {
	// Arrange
	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	testUUID := "a1234567-89ab-cdef-0123-456789abcdef"

	// Act
	encoded, err := handler.EncodeID(testUUID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)
	assert.NotEqual(t, testUUID, encoded) // Should be different from original
}

func TestEncodeID_InvalidUUID(t *testing.T) {
	// Arrange
	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	// Act
	encoded, err := handler.EncodeID("not-a-valid-uuid")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, encoded)
}

func TestDecodeID_Success(t *testing.T) {
	// Arrange
	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	testUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encoded, _ := handler.EncodeID(testUUID)

	// Act
	decoded, err := handler.DecodeID(encoded)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, testUUID, decoded)
}

func TestDecodeID_InvalidEncoded(t *testing.T) {
	// Arrange
	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	// Act
	decoded, err := handler.DecodeID("invalid-encoded-id")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, decoded)
}

func TestDecodeID_EmptyString(t *testing.T) {
	// Arrange
	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	// Act
	decoded, err := handler.DecodeID("")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, decoded)
}

func TestHandleIDDecodingError_SetsError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	// Act
	handler.HandleIDDecodingError(c, "bad-id", assert.AnError)

	// Assert
	assert.True(t, len(c.Errors) > 0)
}

func TestHandleIDEncodingError_SetsError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	encoder := createTestEncoder()
	handler := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, nil)

	// Act
	handler.HandleIDEncodingError(c, "some-uuid", assert.AnError)

	// Assert
	assert.True(t, len(c.Errors) > 0)
}
