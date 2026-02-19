package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createRatingMessageCache creates a MessageCache with rating success/error messages pre-loaded.
func createRatingMessageCache() *messagingCache.MessageCache {
	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		// Rating success
		{Code: "MOD_CS_RATE_EXI_00001", Type: "EXITO", Title: "Calificado", Content: "Servicio calificado exitosamente", Active: true},
		// Rating errors
		{Code: "MOD_CS_ITEM_NF_ERR_00001", Type: "ERROR", Title: "Item no encontrado", Content: "El item de servicio no fue encontrado", Active: true},
		{Code: "MOD_CS_RATE_DUP_ERR_00001", Type: "ERROR", Title: "Ya calificado", Content: "El item ya fue calificado", Active: true},
		{Code: "MOD_CS_RATE_RANGE_ERR_00001", Type: "ERROR", Title: "Rango inválido", Content: "La calificación debe estar entre 1 y 5", Active: true},
		{Code: "MOD_CS_RATE_STATUS_ERR_00001", Type: "ERROR", Title: "No finalizado", Content: "El servicio no está finalizado", Active: true},
		{Code: "MOD_CS_RATE_SAVE_ERR_00001", Type: "ERROR", Title: "Error al guardar", Content: "No se pudo guardar la calificación", Active: true},
		// Validation errors
		{Code: "MOD_V_VAL_ERR_00001", Type: "ERROR", Title: "Formato inválido", Content: "Formato de solicitud inválido", Active: true},
		{Code: "MOD_V_ID_ERR_00013", Type: "ERROR", Title: "ID inválido", Content: "El ID proporcionado no es válido", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	return cache
}

// setupRatingHandler creates the handler, router, and recorder for rating tests.
func setupRatingHandler(t *testing.T, mockSvc *mocks.MockRatingService) (*gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createRatingMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	ratingInteractor := interactor.NewRatingInteractor(mockSvc)
	h := handlers.NewForTestWithRating(ratingInteractor, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   "user-1",
			Role: "USUARIO",
		})
		c.Next()
	})
	router.POST("/completed-services/:id/items/:itemId/rating", h.RateServiceItem())

	return router, httptest.NewRecorder()
}

// ==============================================
// RateServiceItem — Success (full HTTP pipeline)
// ==============================================

func TestRateServiceItem_Controller_Success(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	// Setup valid encoded IDs
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedServiceID, _ := encoder.Encode(serviceUUID)
	encodedItemID, _ := encoder.Encode(itemUUID)

	// Mock the full RateServiceItem flow (interactor delegates to service)
	item := &domain.CompletedServiceItem{
		ID:                 itemUUID,
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	mockSvc.On("GetItemByID", mock.Anything, itemUUID).Return(item, nil)
	mockSvc.On("GetCompletedServiceByID", mock.Anything, "cs-1").Return(cs, nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	comment := "Excelente servicio"
	mockSvc.On("RateServiceItem", mock.Anything, mockTx, itemUUID, 5, &comment).Return(nil)
	mockTx.On("Commit").Return(nil)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{
		"rating":  5,
		"comment": "  Excelente servicio  ", // will be sanitized
	}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/"+encodedItemID+"/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ==============================================
// RateServiceItem — Invalid IDs
// ==============================================

func TestRateServiceItem_Controller_InvalidServiceID(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedItemID, _ := encoder.Encode(itemUUID)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{"rating": 5}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/bad-id/items/"+encodedItemID+"/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestRateServiceItem_Controller_InvalidItemID(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	encodedServiceID, _ := encoder.Encode(serviceUUID)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{"rating": 5}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/bad-id/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

// ==============================================
// RateServiceItem — Invalid JSON body
// ==============================================

func TestRateServiceItem_Controller_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedServiceID, _ := encoder.Encode(serviceUUID)
	encodedItemID, _ := encoder.Encode(itemUUID)

	router, w := setupRatingHandler(t, mockSvc)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/"+encodedItemID+"/rating",
		bytes.NewBuffer([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

// ==============================================
// RateServiceItem — Error mapping
// ==============================================

func TestRateServiceItem_Controller_ItemNotFound(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedServiceID, _ := encoder.Encode(serviceUUID)
	encodedItemID, _ := encoder.Encode(itemUUID)

	mockSvc.On("GetItemByID", mock.Anything, itemUUID).Return(nil, domain.ErrServiceItemNotFound)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{"rating": 5}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/"+encodedItemID+"/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestRateServiceItem_Controller_AlreadyRated(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedServiceID, _ := encoder.Encode(serviceUUID)
	encodedItemID, _ := encoder.Encode(itemUUID)

	existingRating := 3
	item := &domain.CompletedServiceItem{
		ID:                 itemUUID,
		CompletedServiceID: "cs-1",
		Rating:             &existingRating,
	}
	mockSvc.On("GetItemByID", mock.Anything, itemUUID).Return(item, nil)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{"rating": 5}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/"+encodedItemID+"/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestRateServiceItem_Controller_ServiceNotFinalized(t *testing.T) {
	mockSvc := new(mocks.MockRatingService)
	encoder := createTestEncoder()
	serviceUUID := "a1111111-1111-4000-8000-111111111111"
	itemUUID := "a2222222-2222-4000-8000-222222222222"
	encodedServiceID, _ := encoder.Encode(serviceUUID)
	encodedItemID, _ := encoder.Encode(itemUUID)

	item := &domain.CompletedServiceItem{
		ID:                 itemUUID,
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusInProgress}

	mockSvc.On("GetItemByID", mock.Anything, itemUUID).Return(item, nil)
	mockSvc.On("GetCompletedServiceByID", mock.Anything, "cs-1").Return(cs, nil)

	router, w := setupRatingHandler(t, mockSvc)

	reqBody := map[string]interface{}{"rating": 5}
	bodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST",
		"/completed-services/"+encodedServiceID+"/items/"+encodedItemID+"/rating",
		bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}
