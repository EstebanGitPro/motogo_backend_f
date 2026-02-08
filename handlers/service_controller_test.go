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

// TestGetServices_Integration_Success validates the full HTTP pipeline
// for GET /services (HU63).
//
// Exercises: interactor → service catalog → ID encoding → HATEOAS → 200 response.
func TestGetServices_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Create mock service catalog service
	mockCatalogService := new(mocks.MockServiceCatalogService)
	serviceInteractor := interactor.NewServiceInteractor(mockCatalogService)

	h := handlers.NewForTestWithConcrete(
		nil,
		serviceInteractor,
		nil, nil,
		encoder,
		responseHandler,
	)

	// Test data
	services := []domain.Service{
		{
			ID:          "a1111111-1111-4000-8000-111111111111",
			Name:        "Cambio de aceite",
			Description: "Cambio de aceite completo",
			ServiceType: domain.ServiceTypeMaintenance,
			IsActive:    true,
		},
		{
			ID:          "a2222222-2222-4000-8000-222222222222",
			Name:        "Cambio de llanta",
			Description: "Cambio de llanta trasera",
			ServiceType: domain.ServiceTypeTires,
			IsActive:    true,
		},
	}
	mockCatalogService.On("GetAllServices", mock.Anything).Return(services, nil)

	router := gin.New()
	router.GET("/services", h.GetServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	serviceList, ok := data["services"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, serviceList, 2)

	first := serviceList[0].(map[string]interface{})
	assert.NotEmpty(t, first["id"])
	assert.Equal(t, "Cambio de aceite", first["name"])
	assert.Equal(t, "Mantenimiento", first["service_type"])

	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockCatalogService.AssertExpectations(t)
}
