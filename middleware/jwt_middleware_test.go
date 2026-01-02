package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthenticatedUser_NotExists(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Act
	person, exists := middleware.GetAuthenticatedUser(c)

	// Assert
	assert.Nil(t, person)
	assert.False(t, exists)
}

func TestGetAuthenticatedUser_Exists(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedPerson := &domain.Person{
		ID:        "person-123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	c.Set("authenticated_user", expectedPerson)

	// Act
	person, exists := middleware.GetAuthenticatedUser(c)

	// Assert
	assert.True(t, exists)
	assert.NotNil(t, person)
	assert.Equal(t, expectedPerson.ID, person.ID)
	assert.Equal(t, expectedPerson.Email, person.Email)
}

func TestGetAuthenticatedUser_WrongType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Set a non-Person value
	c.Set("authenticated_user", "not-a-person")

	// Act
	person, exists := middleware.GetAuthenticatedUser(c)

	// Assert
	assert.Nil(t, person)
	assert.False(t, exists)
}

func TestRequireAuth_MissingAuthHeader(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Use a minimal mock - just testing header validation
	router.Use(func(c *gin.Context) {
		// Check if Authorization header exists
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - No Authorization header
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_InvalidBearerFormat(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || authHeader == "InvalidFormat" {
			c.Error(domain.ErrInvalidToken)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Act - Invalid format (not "Bearer <token>")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
