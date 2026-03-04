package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// BuildFranchiseLinks Tests
// ============================================

func TestBuildFranchiseLinks_AllLinksPresent(t *testing.T) {
	links := handlers.BuildFranchiseLinks("http://localhost:8080", "encoded-franchise-123")

	assert.NotEmpty(t, links)
	assert.Len(t, links, 6, "Should have 6 HATEOAS links")

	// Check for expected link relations
	relFound := make(map[string]bool)
	for _, link := range links {
		relFound[link.Rel] = true
	}

	assert.True(t, relFound["self"], "Should have 'self' link")
	assert.True(t, relFound["update"], "Should have 'update' link")
	assert.True(t, relFound["delete"], "Should have 'delete' link")
	assert.True(t, relFound["list"], "Should have 'list' link")
	assert.True(t, relFound["branches"], "Should have 'branches' link")
	assert.True(t, relFound["add-branch"], "Should have 'add-branch' link")
}

func TestBuildFranchiseLinks_CorrectURLs(t *testing.T) {
	franchiseID := "abc123"

	links := handlers.BuildFranchiseLinks("", franchiseID)

	// Verify URLs contain expected patterns (relative paths)
	for _, link := range links {
		assert.Contains(t, link.Href, "/motogo/api/v1", "Link should contain API base path")

		switch link.Rel {
		case "self":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "update":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "PUT", link.Method)
		case "delete":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "DELETE", link.Method)
		case "list":
			assert.Equal(t, "GET", link.Method)
		case "branches":
			assert.Contains(t, link.Href, "franchise_id="+franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "add-branch":
			assert.Contains(t, link.Href, franchiseID+"/branches")
			assert.Equal(t, "POST", link.Method)
		}
	}
}

func TestBuildFranchiseLinks_EmptyID(t *testing.T) {
	links := handlers.BuildFranchiseLinks("http://localhost:8080", "")

	assert.NotEmpty(t, links)
	// Should still return links, just with empty ID in URL
}

func TestBuildFranchiseLinks_DifferentBaseURLs(t *testing.T) {
	// baseURL is now ignored — all links are relative paths
	links := handlers.BuildFranchiseLinks("", "test-id")
	assert.NotEmpty(t, links)
	for _, link := range links {
		assert.Contains(t, link.Href, "/motogo/api/v1")
	}
}

// ============================================
// Integration Tests
// ============================================

// TestRegisterFranchise_Integration_Success validates the full HTTP handler pipeline
// for the success path of RegisterFranchise (HU26).
//
// It exercises: JSON binding → sanitization → branch ID decoding → domain mapping →
// FranchiseInteractor.CreateFranchiseWithBranches → response serialization → HATEOAS links → 201 response.
func TestRegisterFranchise_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies for FranchiseInteractor
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	// Handler only needs IDEncoder + Response (interactor passed as method param)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	// Test data
	ownerID := "a1111111-1111-4000-8000-111111111111"
	branchID := "a2222222-2222-4000-8000-222222222222"
	franchiseID := "a4444444-4444-4000-8000-444444444444"
	description := "Red de talleres premium"

	encodedBranchID, _ := encoder.Encode(branchID)

	// Mock: branch validation (ownership) via franchise service
	mockFranchiseService.On("ValidateBranchesForFranchise", mock.Anything, []string{branchID}, ownerID).Return(nil)

	// Mock: transaction lifecycle
	mockFranchiseService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", mock.Anything, mockTx, mock.AnythingOfType("domain.Franchise")).
		Return(&domain.Franchise{
			ID:          franchiseID,
			Name:        "Red Motos Norte",
			Description: &description,
		}, nil)
	mockFranchiseService.On("AssociateBranches", mock.Anything, mockTx, franchiseID, []string{branchID}).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	// Request body
	reqBody := map[string]interface{}{
		"name":        "Red Motos Norte",
		"description": "Red de talleres premium",
		"branch_ids":  []string{encodedBranchID},
	}
	body, _ := json.Marshal(reqBody)

	// Setup router and execute
	router := gin.New()
	router.POST("/franchises", func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID})
	}, h.RegisterFranchise(franchiseInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/franchises", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_F_REG_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Red Motos Norte", data["name"])
	assert.Equal(t, float64(1), data["branch_count"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["_links"])

	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
