package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_GeneratesNewID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())

	var capturedID string
	router.GET("/test", func(c *gin.Context) {
		capturedID = middleware.GetRequestID(c)
		c.Status(http.StatusOK)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, capturedID)

	// Verify it's a valid UUID
	_, err := uuid.Parse(capturedID)
	assert.NoError(t, err, "Generated ID should be a valid UUID")

	// Verify it's in response headers
	responseID := w.Header().Get(middleware.RequestIDHeader)
	assert.Equal(t, capturedID, responseID)
}

func TestRequestID_UsesExistingID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())

	existingID := "test-request-id-123"
	var capturedID string

	router.GET("/test", func(c *gin.Context) {
		capturedID = middleware.GetRequestID(c)
		c.Status(http.StatusOK)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set(middleware.RequestIDHeader, existingID)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, existingID, capturedID)

	// Verify it's in response headers
	responseID := w.Header().Get(middleware.RequestIDHeader)
	assert.Equal(t, existingID, responseID)
}

func TestRequestID_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())

	ids := make(map[string]bool)
	router.GET("/test", func(c *gin.Context) {
		id := middleware.GetRequestID(c)
		ids[id] = true
		c.Status(http.StatusOK)
	})

	// Act - make multiple requests
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
	}

	// Assert
	assert.Equal(t, 10, len(ids), "Should generate 10 unique IDs")
}

func TestGetRequestID_ReturnsEmptyWhenNotSet(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Act
	id := middleware.GetRequestID(c)

	// Assert
	assert.Empty(t, id)
}

func TestGetRequestID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedID := "test-id-456"
	c.Set(middleware.RequestIDKey, expectedID)

	// Act
	id := middleware.GetRequestID(c)

	// Assert
	assert.Equal(t, expectedID, id)
}

func TestGetRequestID_HandlesInvalidType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(middleware.RequestIDKey, 12345) // Set non-string value

	// Act
	id := middleware.GetRequestID(c)

	// Assert
	assert.Empty(t, id, "Should return empty string for non-string values")
}

func TestRequestIDConstants(t *testing.T) {
	// Verify constants are properly defined
	assert.Equal(t, "X-Request-ID", middleware.RequestIDHeader)
	assert.Equal(t, "request_id", middleware.RequestIDKey)
	assert.Equal(t, "traceID", middleware.TraceIDKey)
}
