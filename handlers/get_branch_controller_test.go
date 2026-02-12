package handlers_test

import (
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

// TestGetBranch_Integration_Success validates the full HTTP pipeline
// for GET /branches/:id (HU62).
//
// Exercises: auth → ID decoding → interactor → service → brand encoding →
// location encoding → HATEOAS (owner links) → 200 response.
func TestGetBranch_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Create mock branch service for BranchInteractor
	mockBranchService := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	h := handlers.NewForTestWithConcrete(
		branchInteractor,
		nil, nil, nil,
		encoder,
		responseHandler,
	)

	// IDs
	representativeID := "a1111111-1111-4000-8000-111111111111"
	branchUUID := "a2222222-2222-4000-8000-222222222222"
	brandUUID := "a3333333-3333-4000-8000-333333333333"
	deptUUID := "a4444444-4444-4000-8000-444444444444"
	cityUUID := "a5555555-5555-4000-8000-555555555555"

	encodedBranchID, err := encoder.Encode(branchUUID)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedBranchID)

	lat := 6.2442
	lng := -75.5812

	returnedBranch := &domain.Branch{
		ID:                branchUUID,
		RepresentativeID:  representativeID,
		Name:              "Taller Norte",
		EstablishmentType: domain.EstablishmentTypeWorkshop,
		Status:            domain.BranchStatusActive,
		Brands:            []string{brandUUID},
		Location: &domain.Location{
			ID:             "loc-uuid",
			BranchID:       branchUUID,
			DepartmentID:   deptUUID,
			CityID:         cityUUID,
			Address:        "Calle 45 #12-34",
			Latitude:       &lat,
			Longitude:      &lng,
			CityName:       "Medellín",
			DepartmentName: "Antioquia",
		},
	}

	mockBranchService.On("GetBranchByID", mock.Anything, branchUUID).Return(returnedBranch, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simulating the owner is logged in → should get owner HATEOAS links
		c.Set("authenticated_user", &domain.Person{
			ID:   representativeID,
			Role: "TALLER",
		})
		c.Next()
	})
	router.GET("/branches/:id", h.GetBranch())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, encodedBranchID, data["id"])
	assert.Equal(t, "Taller Norte", data["name"])
	assert.Equal(t, "WORKSHOP", data["establishment_type"])
	assert.Equal(t, "ACTIVE", data["status"])

	// Verify brands are encoded
	brands, ok := data["brands"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, brands, 1)
	assert.NotEqual(t, brandUUID, brands[0]) // must be encoded, not raw UUID

	// Verify location IDs are encoded
	location, ok := data["location"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEqual(t, deptUUID, location["department_id"]) // encoded
	assert.NotEqual(t, cityUUID, location["city_id"])       // encoded
	assert.Equal(t, "Calle 45 #12-34", location["address"])

	// Verify HATEOAS links present (owner should see edit/delete links)
	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockBranchService.AssertExpectations(t)
}
