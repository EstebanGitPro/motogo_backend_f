package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
// TestGeocoding Tests
// ============================================

func TestGeocoding_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	h := handlers.NewForTestWithConcrete(branchInteractor, nil, nil, nil, encoder, responseHandler)

	lat, lng := 6.2442, -75.5812
	mockBranchSvc.On("GeocodeLocation", mock.Anything, mock.AnythingOfType("*domain.Location")).
		Run(func(args mock.Arguments) {
			loc := args.Get(1).(*domain.Location)
			loc.Latitude = &lat
			loc.Longitude = &lng
		}).
		Return(true, nil)

	body, _ := json.Marshal(map[string]string{
		"address":         "Calle 10 # 43A-20",
		"city_name":       "Medellín",
		"department_name": "Antioquia",
	})

	router := gin.New()
	router.POST("/location/geocode", h.TestGeocoding())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/location/geocode", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "GEOCODING_TEST_COMPLETE", resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.True(t, data["geocoded"].(bool))
	assert.InDelta(t, 6.2442, data["latitude"].(float64), 0.001)
	assert.InDelta(t, -75.5812, data["longitude"].(float64), 0.001)
	assert.Contains(t, data["formatted_address"], "Medellín")
	mockBranchSvc.AssertExpectations(t)
}

func TestGeocoding_GeocodeFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	h := handlers.NewForTestWithConcrete(branchInteractor, nil, nil, nil, encoder, responseHandler)

	mockBranchSvc.On("GeocodeLocation", mock.Anything, mock.AnythingOfType("*domain.Location")).
		Return(false, errors.New("geocoding service unavailable"))

	body, _ := json.Marshal(map[string]string{
		"address":         "Dirección Desconocida",
		"city_name":       "XYZ",
		"department_name": "ABC",
	})

	router := gin.New()
	router.POST("/location/geocode", h.TestGeocoding())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/location/geocode", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.False(t, data["geocoded"].(bool))
	assert.Contains(t, data["error"], "geocoding service unavailable")
	mockBranchSvc.AssertExpectations(t)
}

func TestGeocoding_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	h := handlers.NewForTestWithConcrete(branchInteractor, nil, nil, nil, encoder, responseHandler)

	// Missing required fields
	body, _ := json.Marshal(map[string]string{
		"address": "Calle 10",
	})

	router := gin.New()
	router.POST("/location/geocode", h.TestGeocoding())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/location/geocode", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp["success"].(bool))
	assert.Equal(t, "ERR_INVALID_REQUEST", resp["code"])
}
