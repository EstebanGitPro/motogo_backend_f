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
// GetServiceTypes Tests (HU75)
// ============================================

func TestGetServiceTypes_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	mockSvc.On("GetAllServiceTypes").Return([]domain.ServiceType{
		domain.ServiceTypeMaintenance,
		domain.ServiceTypeTires,
	})

	router := gin.New()
	router.GET("/service-types", h.GetServiceTypes())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/service-types", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	types, ok := data["types"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, types, 2)

	mockSvc.AssertExpectations(t)
}

// ============================================
// GetServices Tests (HU63)
// ============================================

func TestGetServices_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockCatalogService := new(mocks.MockServiceCatalogService)
	serviceInteractor := interactor.NewServiceInteractor(mockCatalogService)
	h := handlers.NewForTestWithConcrete(nil, serviceInteractor, nil, nil, encoder, responseHandler)

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

func TestGetServices_WithTypeFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	services := []domain.Service{
		{ID: "a1111111-1111-4000-8000-111111111111", Name: "Cambio de aceite", ServiceType: domain.ServiceTypeMaintenance, IsActive: true},
	}
	mockSvc.On("GetServicesByType", mock.Anything, domain.ServiceTypeMaintenance).Return(services, nil)

	router := gin.New()
	router.GET("/services", h.GetServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services?type=Mantenimiento", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetServices_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/services", h.GetServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services?type=INVALID_TYPE", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetServices_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	mockSvc.On("GetAllServices", mock.Anything).Return(nil, assert.AnError)

	router := gin.New()
	router.GET("/services", h.GetServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetServices_FilterTypeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	mockSvc.On("GetServicesByType", mock.Anything, domain.ServiceTypeMaintenance).Return(nil, assert.AnError)

	router := gin.New()
	router.GET("/services", h.GetServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services?type=Mantenimiento", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// GetBranchServices Tests
// ============================================

func TestGetBranchServices_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	encodedBranch, _ := encoder.Encode(branchUUID)

	branchServices := []domain.BranchServiceInfo{
		{
			Service: domain.Service{
				ID:          "a2222222-2222-4000-8000-222222222222",
				Name:        "Cambio de aceite",
				ServiceType: domain.ServiceTypeMaintenance,
			},
			AddedAt: "2026-01-15",
			Active:  true,
		},
	}
	mockSvc.On("GetServicesByBranch", mock.Anything, branchUUID).Return(branchServices, nil)

	router := gin.New()
	router.GET("/branches/:id/services", h.GetBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranch+"/services", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	mockSvc.AssertExpectations(t)
}

func TestGetBranchServices_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/branches/:id/services", h.GetBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/INVALID!!!/services", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetBranchServices_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	encodedBranch, _ := encoder.Encode(branchUUID)

	mockSvc.On("GetServicesByBranch", mock.Anything, branchUUID).Return(nil, assert.AnError)

	router := gin.New()
	router.GET("/branches/:id/services", h.GetBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranch+"/services", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// AssociateBranchServices Tests
// ============================================

func TestAssociateBranchServices_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("ValidateServiceIDs", mock.Anything, []string{serviceUUID}).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("AssociateBranchServices", mock.Anything, mockTx, branchUUID, []string{serviceUUID}).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	reqBody := map[string]interface{}{
		"service_ids": []string{encodedService},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
}

func TestAssociateBranchServices_InvalidBranchID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/INVALID!!!/services", bytes.NewBuffer([]byte(`{"service_ids":["x"]}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	encodedBranch, _ := encoder.Encode(branchUUID)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_InvalidServiceID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	encodedBranch, _ := encoder.Encode(branchUUID)

	reqBody := map[string]interface{}{
		"service_ids": []string{"INVALID!!!"},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("ValidateServiceIDs", mock.Anything, []string{serviceUUID}).Return(domain.ErrServiceNotFound)

	reqBody := map[string]interface{}{
		"service_ids": []string{encodedService},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_TxError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("ValidateServiceIDs", mock.Anything, []string{serviceUUID}).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(nil, assert.AnError)

	reqBody := map[string]interface{}{
		"service_ids": []string{encodedService},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_AssociateError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("ValidateServiceIDs", mock.Anything, []string{serviceUUID}).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("AssociateBranchServices", mock.Anything, mockTx, branchUUID, []string{serviceUUID}).Return(assert.AnError)
	mockTx.On("Rollback").Return(nil)

	reqBody := map[string]interface{}{
		"service_ids": []string{encodedService},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAssociateBranchServices_CommitError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("ValidateServiceIDs", mock.Anything, []string{serviceUUID}).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("AssociateBranchServices", mock.Anything, mockTx, branchUUID, []string{serviceUUID}).Return(nil)
	mockTx.On("Commit").Return(assert.AnError)

	reqBody := map[string]interface{}{
		"service_ids": []string{encodedService},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/branches/:id/services", h.AssociateBranchServices())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranch+"/services", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// DissociateBranchService Tests
// ============================================

func TestDissociateBranchService_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranchService", mock.Anything, mockTx, branchUUID, serviceUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/"+encodedService, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDissociateBranchService_InvalidBranchID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/INVALID!!!/services/abc", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDissociateBranchService_InvalidServiceID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	encodedBranch, _ := encoder.Encode(branchUUID)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/INVALID!!!", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDissociateBranchService_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranchService", mock.Anything, mockTx, branchUUID, serviceUUID).Return(domain.ErrServiceNotFound)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/"+encodedService, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDissociateBranchService_TxError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("BeginTx", mock.Anything).Return(nil, assert.AnError)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/"+encodedService, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDissociateBranchService_ServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranchService", mock.Anything, mockTx, branchUUID, serviceUUID).Return(assert.AnError)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/"+encodedService, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDissociateBranchService_CommitError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	branchUUID := "a1111111-1111-4000-8000-111111111111"
	serviceUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranch, _ := encoder.Encode(branchUUID)
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranchService", mock.Anything, mockTx, branchUUID, serviceUUID).Return(nil)
	mockTx.On("Commit").Return(assert.AnError)

	router := gin.New()
	router.DELETE("/branches/:id/services/:serviceId", h.DissociateBranchService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranch+"/services/"+encodedService, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// UpdateService Tests (HU68)
// ============================================

func TestUpdateService_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	existingService := &domain.Service{
		ID:          serviceUUID,
		Name:        "Old Name",
		Description: "Old desc",
		ServiceType: domain.ServiceTypeMaintenance,
		IsActive:    true,
	}

	// Controller: GetServiceByID → UpdateService (interactor handles tx internally)
	mockSvc.On("GetServiceByID", mock.Anything, serviceUUID).Return(existingService, nil)
	// Interactor internally: BeginTx → UpdateService(ctx, tx, svc) → Commit
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("UpdateService", mock.Anything, mockTx, mock.AnythingOfType("domain.Service")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	isActive := true
	reqBody := map[string]interface{}{
		"name":         "New Name",
		"description":  "New description",
		"service_type": "Mantenimiento",
		"is_active":    isActive,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateService_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/INVALID!!!", bytes.NewBuffer([]byte(`{"name":"x","service_type":"Mantenimiento"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateService_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	// Missing required "name" and "service_type"
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateService_InvalidServiceType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	reqBody := map[string]interface{}{
		"name":         "Test",
		"service_type": "INVALID_TYPE",
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateService_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("GetServiceByID", mock.Anything, serviceUUID).Return(nil, domain.ErrServiceNotFound)

	reqBody := map[string]interface{}{
		"name":         "Test",
		"service_type": "Mantenimiento",
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateService_GetByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	mockSvc.On("GetServiceByID", mock.Anything, serviceUUID).Return(nil, assert.AnError)

	reqBody := map[string]interface{}{
		"name":         "Test",
		"service_type": "Mantenimiento",
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateService_UpdateError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockServiceCatalogService)
	mockTx := new(mocks.MockTx)
	svcInteractor := interactor.NewServiceInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, svcInteractor, nil, nil, encoder, responseHandler)

	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedService, _ := encoder.Encode(serviceUUID)

	existingService := &domain.Service{
		ID:          serviceUUID,
		Name:        "Old",
		ServiceType: domain.ServiceTypeMaintenance,
		IsActive:    true,
	}
	mockSvc.On("GetServiceByID", mock.Anything, serviceUUID).Return(existingService, nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("UpdateService", mock.Anything, mockTx, mock.AnythingOfType("domain.Service")).Return(assert.AnError)
	mockTx.On("Rollback").Return(nil)

	reqBody := map[string]interface{}{
		"name":         "New Name",
		"service_type": "Mantenimiento",
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/services/:id", h.UpdateService())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/services/"+encodedService, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}
