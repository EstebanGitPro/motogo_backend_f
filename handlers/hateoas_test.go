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
	assert.Contains(t, profileLink.Href, "/auth/me") // Login still points to /auth/me as it's a public reference

	resetLink := findLinkByRel(links, "password-reset")
	assert.NotNil(t, resetLink)
	assert.Equal(t, "POST", resetLink.Method)
}

func TestBuildAuthMeLinks(t *testing.T) {
	// Act
	links := handlers.BuildAuthMeLinks("http://localhost:8080")

	// Assert - now returns 3 links: self, change-password, login
	assert.Len(t, links, 3)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/persons/me")

	changePasswordLink := findLinkByRel(links, "change-password")
	assert.NotNil(t, changePasswordLink)
	assert.Contains(t, changePasswordLink.Href, "/persons/me/password")

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

// ============================================
// Additional HATEOAS Builder Tests
// ============================================

func TestBuildChangePasswordLinks(t *testing.T) {
	// Act
	links := handlers.BuildChangePasswordLinks("http://localhost:8080")

	// Assert
	assert.Len(t, links, 2)

	profileLink := findLinkByRel(links, "profile")
	assert.NotNil(t, profileLink)
	assert.Contains(t, profileLink.Href, "/persons/me")

	loginLink := findLinkByRel(links, "login")
	assert.NotNil(t, loginLink)
}

func TestBuildUpdateProfileLinks(t *testing.T) {
	// Act
	links := handlers.BuildUpdateProfileLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/persons/me")
}

func TestBuildPublicContactLinks(t *testing.T) {
	// Act
	links := handlers.BuildPublicContactLinks("http://localhost:8080", "person-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
}

func TestBuildBranchCreatedLinks(t *testing.T) {
	// Act
	links := handlers.BuildBranchCreatedLinks("http://localhost:8080", "branch-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "branches/branch-123")
	assert.Equal(t, "GET", selfLink.Method)
}

func TestBuildBrandListLinks(t *testing.T) {
	// Act
	links := handlers.BuildBrandListLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/brands")
}

func TestBuildDepartmentListLinks(t *testing.T) {
	// Act
	links := handlers.BuildDepartmentListLinks()

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/departments")
}

func TestBuildCityListLinks(t *testing.T) {
	// Act
	links := handlers.BuildCityListLinks("dept-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/cities")
}

func TestBuildBranchTypesLinks(t *testing.T) {
	// Act
	links := handlers.BuildBranchTypesLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/branch-types")
}

func TestBuildBranchDetailLinks(t *testing.T) {
	// Act
	links := handlers.BuildBranchDetailLinks("http://localhost:8080", "branch-123", true)

	// Assert
	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "branches/branch-123")
}

func TestBuildBranchDetailLinks_NoEditNoDelete(t *testing.T) {
	// Act
	links := handlers.BuildBranchDetailLinks("http://localhost:8080", "branch-123", false)

	// Assert
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)

	// When canEdit=false and canDelete=false, update and delete links should not be present
	updateLink := findLinkByRel(links, "update")
	assert.Nil(t, updateLink)

	deleteLink := findLinkByRel(links, "delete")
	assert.Nil(t, deleteLink)
}

// ============================================
// HATEOAS Builders with 0% coverage
// ============================================

func TestBuildBranchListLinks(t *testing.T) {
	// Act
	links := handlers.BuildBranchListLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/branches")
}

func TestBuildBranchDeletedLinks(t *testing.T) {
	// Act
	links := handlers.BuildBranchDeletedLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	listLink := findLinkByRel(links, "list")
	assert.NotNil(t, listLink)
}

func TestBuildServiceTypeListLinks(t *testing.T) {
	// Act
	links := handlers.BuildServiceTypeListLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/service-types")
}

func TestBuildServiceListLinks(t *testing.T) {
	// Act
	links := handlers.BuildServiceListLinks("http://localhost:8080", "Mantenimiento")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/services")
}

func TestBuildServiceListLinks_NoFilter(t *testing.T) {
	// Act
	links := handlers.BuildServiceListLinks("http://localhost:8080", "")

	// Assert
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
}

func TestBuildMotorcycleDetailLinks(t *testing.T) {
	// Act
	links := handlers.BuildMotorcycleDetailLinks("http://localhost:8080", "moto-123", true)

	// Assert
	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "motorcycles/moto-123")
}

func TestBuildMotorcycleDetailLinks_ReadOnly(t *testing.T) {
	// Act
	links := handlers.BuildMotorcycleDetailLinks("http://localhost:8080", "moto-456", false)

	// Assert
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)

	updateLink := findLinkByRel(links, "update")
	assert.Nil(t, updateLink)
}

func TestBuildMotorcycleDeletedLinks(t *testing.T) {
	// Act
	links := handlers.BuildMotorcycleDeletedLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	listLink := findLinkByRel(links, "list")
	assert.NotNil(t, listLink)
}

func TestBuildNearbyBranchesLinks(t *testing.T) {
	// Act
	links := handlers.BuildNearbyBranchesLinks("http://localhost:8080", 4.60971, -74.08175, 5.0)

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/branches/nearby")
}

func TestBuildMotorcycleReferencesLinks(t *testing.T) {
	// Act
	links := handlers.BuildMotorcycleReferencesLinks("http://localhost:8080")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/motorcycle-references")
}

func TestBuildBrandLinesLinks(t *testing.T) {
	// Act
	links := handlers.BuildBrandLinesLinks("http://localhost:8080", "brand-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "brands/brand-123/lines")
}

func TestBuildEvidenceDetailLinks(t *testing.T) {
	// Act
	links := handlers.BuildEvidenceDetailLinks("http://localhost:8080", "moto-123", "ev-456", true)

	// Assert
	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "evidence/ev-456")
}

func TestBuildEvidenceDetailLinks_NoDelete(t *testing.T) {
	// Act
	links := handlers.BuildEvidenceDetailLinks("http://localhost:8080", "moto-123", "ev-789", false)

	// Assert
	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)

	deleteLink := findLinkByRel(links, "delete")
	assert.Nil(t, deleteLink)
}

func TestBuildEvidenceListLinks(t *testing.T) {
	// Act
	links := handlers.BuildEvidenceListLinks("http://localhost:8080", "moto-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "motorcycles/moto-123/evidence")
}

func TestBuildEvidenceDeletedLinks(t *testing.T) {
	// Act
	links := handlers.BuildEvidenceDeletedLinks("http://localhost:8080", "moto-123")

	// Assert
	assert.GreaterOrEqual(t, len(links), 1)

	listLink := findLinkByRel(links, "list")
	assert.NotNil(t, listLink)
}

// ============================================
// New Catalog HATEOAS Builder Tests (PR #52)
// ============================================

func TestBuildMotorcycleCategoryListLinks(t *testing.T) {
	links := handlers.BuildMotorcycleCategoryListLinks("http://localhost:8080")

	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/motorcycle-categories")

	refsLink := findLinkByRel(links, "references")
	assert.NotNil(t, refsLink)
	assert.Contains(t, refsLink.Href, "/motorcycle-references")
}

func TestBuildMotorcycleCategoryItemLinks(t *testing.T) {
	links := handlers.BuildMotorcycleCategoryItemLinks("http://localhost:8080", "Sport")

	assert.GreaterOrEqual(t, len(links), 1)

	linesLink := findLinkByRel(links, "lines")
	assert.NotNil(t, linesLink)
	assert.Contains(t, linesLink.Href, "/motorcycle-categories/Sport/lines")
}

func TestBuildCategoryLinesLinks(t *testing.T) {
	links := handlers.BuildCategoryLinesLinks("http://localhost:8080", "Sport")

	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/motorcycle-categories/Sport/lines")

	catLink := findLinkByRel(links, "categories")
	assert.NotNil(t, catLink)
	assert.Contains(t, catLink.Href, "/motorcycle-categories")
}

func TestBuildEngineDisplacementLinks(t *testing.T) {
	links := handlers.BuildEngineDisplacementLinks("http://localhost:8080")

	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/engine-displacements")

	catLink := findLinkByRel(links, "categories")
	assert.NotNil(t, catLink)
}

func TestBuildRatingRangeLinks(t *testing.T) {
	links := handlers.BuildRatingRangeLinks("http://localhost:8080")

	assert.GreaterOrEqual(t, len(links), 2)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "/rating-ranges")

	dispLink := findLinkByRel(links, "displacements")
	assert.NotNil(t, dispLink)
}

func TestBuildPermissionLinks(t *testing.T) {
	links := handlers.BuildPermissionLinks("http://localhost:8080", "moto-123")

	assert.GreaterOrEqual(t, len(links), 1)

	selfLink := findLinkByRel(links, "self")
	assert.NotNil(t, selfLink)
	assert.Contains(t, selfLink.Href, "motorcycles/moto-123")
}
