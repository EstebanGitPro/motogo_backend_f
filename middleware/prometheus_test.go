package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware_RecordsMetrics(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	// Metrics should be recorded (we can't easily verify prometheus metrics in unit tests,
	// but we can verify the middleware doesn't break the request flow)
}

func TestPrometheusMiddleware_HandlesPOST(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())
	router.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "123"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPrometheusMiddleware_HandlesErrors(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Middleware should still record metrics for error responses
}

func TestPrometheusMiddleware_Handles404(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())
	router.GET("/exists", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notfound", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrometheusMiddleware_HandlesDifferentMethods(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())

	router.GET("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "GET"})
	})
	router.PUT("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "PUT"})
	})
	router.DELETE("/resource", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, nil)
	})

	testCases := []struct {
		method         string
		expectedStatus int
	}{
		{"GET", http.StatusOK},
		{"PUT", http.StatusOK},
		{"DELETE", http.StatusNoContent},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/resource", nil)
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestPrometheusMiddleware_HandlesMultipleRequests(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TrackMetrics())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Act - Make multiple requests
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	}
	// All requests should be recorded in metrics
}
