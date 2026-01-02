package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetBaseURL_HTTP(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://localhost:8080/test", nil)
	c.Request.Host = "localhost:8080"

	// Act
	baseURL := handlers.GetBaseURL(c)

	// Assert
	assert.Equal(t, "http://localhost:8080", baseURL)
}

func TestBuildResourceURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		resource   string
		resourceID string
		expected   string
	}{
		{
			name:       "accounts resource",
			baseURL:    "http://localhost:8080",
			resource:   "accounts",
			resourceID: "xyz123",
			expected:   "http://localhost:8080/motogo/api/v1/accounts/xyz123",
		},
		{
			name:       "messages resource",
			baseURL:    "https://api.example.com",
			resource:   "messages",
			resourceID: "abc456",
			expected:   "https://api.example.com/motogo/api/v1/messages/abc456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.BuildResourceURL(tt.baseURL, tt.resource, tt.resourceID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildCollectionURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		resource string
		expected string
	}{
		{
			name:     "accounts collection",
			baseURL:  "http://localhost:8080",
			resource: "accounts",
			expected: "http://localhost:8080/motogo/api/v1/accounts",
		},
		{
			name:     "messages collection",
			baseURL:  "https://api.example.com",
			resource: "messages",
			expected: "https://api.example.com/motogo/api/v1/messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.BuildCollectionURL(tt.baseURL, tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildResourceLinks(t *testing.T) {
	// Arrange
	baseURL := "http://localhost:8080"
	resource := "accounts"
	resourceID := "xyz123"

	// Act
	links := handlers.BuildResourceLinks(baseURL, resource, resourceID)

	// Assert
	assert.Len(t, links, 4)

	// Check self link
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Equal(t, "http://localhost:8080/motogo/api/v1/accounts/xyz123", selfLink.Href)
	assert.Equal(t, "GET", selfLink.Method)

	// Check update link
	updateLink := findLinkByRel(links, "update")
	assert.NotNil(t, updateLink)
	assert.Equal(t, "PUT", updateLink.Method)

	// Check delete link
	deleteLink := findLinkByRel(links, "delete")
	assert.NotNil(t, deleteLink)
	assert.Equal(t, "DELETE", deleteLink.Method)

	// Check collection link
	collectionLink := findLinkByRel(links, "collection")
	assert.NotNil(t, collectionLink)
	assert.Equal(t, "http://localhost:8080/motogo/api/v1/accounts", collectionLink.Href)
}

func TestBuildAccountLinks(t *testing.T) {
	// Act
	links := handlers.BuildAccountLinks("http://localhost:8080", "account-123")

	// Assert
	assert.Len(t, links, 4)
	selfLink := findLinkByRel(links, "self")
	assert.Contains(t, selfLink.Href, "accounts/account-123")
}

func TestBuildMessageLinks(t *testing.T) {
	// Act
	links := handlers.BuildMessageLinks("http://localhost:8080", "msg-123")

	// Assert
	assert.Len(t, links, 4)
	selfLink := findLinkByRel(links, "self")
	assert.Contains(t, selfLink.Href, "messages/msg-123")
}

func TestBuildMessageCreatedLinks(t *testing.T) {
	// Act
	links := handlers.BuildMessageCreatedLinks("http://localhost:8080", "msg-123")

	// Assert
	assert.Len(t, links, 4)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Equal(t, "GET", selfLink.Method)

	listLink := findLinkByRel(links, "list")
	assert.NotNil(t, listLink)
	assert.Equal(t, "GET", listLink.Method)
}

func TestBuildMessageUpdatedLinks(t *testing.T) {
	// Act
	links := handlers.BuildMessageUpdatedLinks("http://localhost:8080", "msg-123")

	// Assert
	assert.Len(t, links, 3)
	assert.NotNil(t, findLinkByRel(links, "self"))
	assert.NotNil(t, findLinkByRel(links, "delete"))
	assert.NotNil(t, findLinkByRel(links, "list"))
}

func TestBuildMessageListLinks(t *testing.T) {
	// Act
	links := handlers.BuildMessageListLinks("http://localhost:8080")

	// Assert
	assert.Len(t, links, 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/messages")

	createLink := findLinkByRel(links, "create")
	assert.NotNil(t, createLink)
	assert.Equal(t, "POST", createLink.Method)
}

func TestBuildLoginLinks(t *testing.T) {
	// Act
	links := handlers.BuildLoginLinks("http://localhost:8080")

	// Assert
	assert.Len(t, links, 2)

	profileLink := findLinkByRel(links, "profile")
	assert.NotNil(t, profileLink)
	assert.Contains(t, profileLink.Href, "/auth/me")

	resetLink := findLinkByRel(links, "password-reset")
	assert.NotNil(t, resetLink)
	assert.Equal(t, "POST", resetLink.Method)
}

func TestBuildAuthMeLinks(t *testing.T) {
	// Act
	links := handlers.BuildAuthMeLinks("http://localhost:8080")

	// Assert
	assert.Len(t, links, 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/auth/me")

	loginLink := findLinkByRel(links, "login")
	assert.NotNil(t, loginLink)
}

func TestBuildAccountCreatedLinks(t *testing.T) {
	// Act
	links := handlers.BuildAccountCreatedLinks("http://localhost:8080", "acct-xyz")

	// Assert
	assert.Len(t, links, 3)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "accounts/acct-xyz")

	loginLink := findLinkByRel(links, "login")
	assert.NotNil(t, loginLink)

	verifyLink := findLinkByRel(links, "verify-email")
	assert.NotNil(t, verifyLink)
}

func TestSetLocationHeader(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", nil)

	// Act
	handlers.SetLocationHeader(c, "http://localhost:8080", "accounts", "xyz123")

	// Assert
	location := w.Header().Get("Location")
	assert.Equal(t, "http://localhost:8080/motogo/api/v1/accounts/xyz123", location)
}

func TestLink_Structure(t *testing.T) {
	// Arrange
	link := handlers.Link{
		Href:   "/api/v1/test",
		Rel:    "self",
		Method: http.MethodGet,
	}

	// Assert
	assert.Equal(t, "/api/v1/test", link.Href)
	assert.Equal(t, "self", link.Rel)
	assert.Equal(t, "GET", link.Method)
}

// Helper function
func findLinkByRel(links []handlers.Link, rel string) *handlers.Link {
	for i, link := range links {
		if link.Rel == rel {
			return &links[i]
		}
	}
	return nil
}
