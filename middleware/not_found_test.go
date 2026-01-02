package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler_ReturnsJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(middleware.NotFoundHandler())

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent-path", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "404_NOT_FOUND", response["code"])
	assert.Equal(t, "Endpoint no encontrado", response["message"])
	assert.Equal(t, "/nonexistent-path", response["path"])
}

func TestNotFoundHandler_IncludesPath(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(middleware.NotFoundHandler())

	testPaths := []string{
		"/api/v1/unknown",
		"/users/123/profile",
		"/does/not/exist",
	}

	for _, path := range testPaths {
		t.Run(path, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", path, nil)
			router.ServeHTTP(w, req)

			// Assert
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, path, response["path"])
		})
	}
}

func TestNotFoundHandler_WorksWithRequestID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())
	router.NoRoute(middleware.NotFoundHandler())

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/nonexistent", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	// Request ID should be in response headers
	assert.NotEmpty(t, w.Header().Get(middleware.RequestIDHeader))
}

func TestNotFoundHandler_DifferentMethods(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(middleware.NotFoundHandler())

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, "/not-found", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	}
}
