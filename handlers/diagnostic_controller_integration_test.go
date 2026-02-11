package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// authMiddleware injects an authenticated user into the gin context for tests.
func authMiddleware(userID, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   userID,
			Role: role,
		})
		c.Next()
	}
}

// ============================================================
// TestListDiagnostics_Integration_Success (HU14)
// ============================================================
func TestListDiagnostics_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)

	desc := "No enciende"
	now := time.Now()

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, motorcycleUUID, ownerID).Return(nil)
	mockService.On("GetDiagnosticsByMotorcycleID", mock.Anything, motorcycleUUID).
		Return([]domain.Diagnostic{
			{
				ID:                 diagUUID,
				MotorcycleID:       motorcycleUUID,
				BranchID:           branchUUID,
				Date:               now,
				ProblemDescription: &desc,
			},
		}, nil)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.GET("/motorcycles/:id/diagnostics", h.ListDiagnostics())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)

	first := data[0].(map[string]interface{})
	assert.NotEmpty(t, first["id"])
	assert.Equal(t, encodedMotorcycleID, first["motorcycle_id"])

	mockService.AssertExpectations(t)
}

// ============================================================
// TestListDiagnostics_Integration_MotorcycleNotFound
// ============================================================
func TestListDiagnostics_Integration_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, motorcycleUUID, ownerID).
		Return(domain.ErrMotorcycleNotFound)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.GET("/motorcycles/:id/diagnostics", h.ListDiagnostics())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// TestGetDiagnostic_Integration_Success (HU14)
// ============================================================
func TestGetDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedDiagID, _ := encoder.Encode(diagUUID)

	desc := "Frenos desgastados"
	now := time.Now()

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, mock.Anything, ownerID).Return(nil).Maybe()
	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(&domain.Diagnostic{
			ID:                 diagUUID,
			MotorcycleID:       motorcycleUUID,
			BranchID:           branchUUID,
			Date:               now,
			ProblemDescription: &desc,
		}, nil)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.GET("/motorcycles/:id/diagnostics/:diagnosticId", h.GetDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, encodedDiagID, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])

	mockService.AssertExpectations(t)
}

// ============================================================
// TestGetDiagnostic_Integration_NotFound
// ============================================================
func TestGetDiagnostic_Integration_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedDiagID, _ := encoder.Encode(diagUUID)

	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(nil, domain.ErrDiagnosticNotFound)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.GET("/motorcycles/:id/diagnostics/:diagnosticId", h.GetDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// TestUpdateDiagnostic_Integration_Success (HU12)
// ============================================================
func TestUpdateDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedDiagID, _ := encoder.Encode(diagUUID)

	now := time.Now()

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, mock.Anything, ownerID).Return(nil).Maybe()
	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(&domain.Diagnostic{
			ID:           diagUUID,
			MotorcycleID: motorcycleUUID,
			BranchID:     branchUUID,
			Date:         now,
		}, nil)

	mockTx := new(mocks.MockTx)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("ApplyDiagnosticUpdates", mock.Anything, mock.Anything).Return()
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, mock.Anything).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	newDesc := "Cadena oxidada"
	reqBody := map[string]interface{}{
		"problem_description": newDesc,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.PUT("/motorcycles/:id/diagnostics/:diagnosticId", h.UpdateDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagID, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================================
// TestDeleteDiagnostic_Integration_Success (HU13)
// ============================================================
func TestDeleteDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedDiagID, _ := encoder.Encode(diagUUID)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, mock.Anything, ownerID).Return(nil).Maybe()
	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(&domain.Diagnostic{
			ID:           diagUUID,
			MotorcycleID: motorcycleUUID,
			Date:         time.Now(),
		}, nil)

	mockTx := new(mocks.MockTx)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteDiagnostic", mock.Anything, mockTx, diagUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.DELETE("/motorcycles/:id/diagnostics/:diagnosticId", h.DeleteDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================================
// TestDeleteDiagnostic_Integration_NotFound
// ============================================================
func TestDeleteDiagnostic_Integration_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedDiagID, _ := encoder.Encode(diagUUID)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, mock.Anything, ownerID).Return(nil).Maybe()
	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(nil, domain.ErrDiagnosticNotFound)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.DELETE("/motorcycles/:id/diagnostics/:diagnosticId", h.DeleteDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// TestSetDiagnosticSolution_Integration_Success
// ============================================================
func TestSetDiagnosticSolution_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedDiagID, _ := encoder.Encode(diagUUID)
	solution := "Cambiar la bujía y revisar cables"

	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(&domain.Diagnostic{
			ID:   diagUUID,
			Date: time.Now(),
		}, nil)

	mockTx := new(mocks.MockTx)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, mock.Anything).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	reqBody := map[string]interface{}{
		"possible_solution": solution,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(authMiddleware("rep-id", "REPRESENTANTE"))
	router.PATCH("/diagnostics/:id/solution", h.SetDiagnosticSolution())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/diagnostics/"+encodedDiagID+"/solution", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================================
// TestSetDiagnosticSolution_Integration_NotFound
// ============================================================
func TestSetDiagnosticSolution_Integration_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	diagUUID := "a4444444-4444-4000-8000-444444444444"
	encodedDiagID, _ := encoder.Encode(diagUUID)
	solution := "No aplica"

	mockService.On("GetDiagnosticByID", mock.Anything, diagUUID).
		Return(nil, domain.ErrDiagnosticNotFound)

	reqBody := map[string]interface{}{
		"possible_solution": solution,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(authMiddleware("rep-id", "REPRESENTANTE"))
	router.PATCH("/diagnostics/:id/solution", h.SetDiagnosticSolution())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/diagnostics/"+encodedDiagID+"/solution", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// TestCreateDiagnostic_Integration_Unauthorized
// ============================================================
func TestCreateDiagnostic_Integration_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockService := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockService)
	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)

	reqBody := map[string]interface{}{
		"branch_id": "some-branch",
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	// No auth middleware → no authenticated user
	router.POST("/motorcycles/:id/diagnostics", h.CreateDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Without auth middleware, the CreateDiagnostic controller returns MsgUnauthorized (401).
	// In production, the auth middleware would reject the request before reaching the controller.
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
