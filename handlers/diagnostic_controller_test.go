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

// === Common test IDs ===

const (
	diagTestOwnerID        = "a1111111-1111-4000-8000-111111111111"
	diagTestMotorcycleUUID = "a2222222-2222-4000-8000-222222222222"
	diagTestBranchUUID     = "a3333333-3333-4000-8000-333333333333"
	diagTestDiagnosticUUID = "a4444444-4444-4000-8000-444444444444"
)

// TestCreateDiagnostic_Integration_Success validates the full HTTP pipeline
// for POST /motorcycles/:id/diagnostics (HU11).
func TestCreateDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedMotorcycleID, _ := encoder.Encode(diagTestMotorcycleUUID)
	encodedBranchID, _ := encoder.Encode(diagTestBranchUUID)
	problemDesc := "La moto no enciende"

	// Interactor flow: ValidateOwnership → ValidateBranch → BeginTx → RegisterOrUpdate → Commit
	mockSvc.On("ValidateMotorcycleOwnership", mock.Anything, diagTestMotorcycleUUID, diagTestOwnerID).Return(&domain.Motorcycle{
		ID:      diagTestMotorcycleUUID,
		OwnerID: diagTestOwnerID,
	}, nil)
	mockSvc.On("ValidateBranchExists", mock.Anything, diagTestBranchUUID).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("RegisterOrUpdateDiagnostic", mock.Anything, mockTx, diagTestMotorcycleUUID, diagTestBranchUUID, &problemDesc).Return(&domain.Diagnostic{
		ID:           diagTestDiagnosticUUID,
		MotorcycleID: diagTestMotorcycleUUID,
		BranchID:     diagTestBranchUUID,
		Date:         time.Now(),
	}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"branch_id":           encodedBranchID,
		"problem_description": problemDesc,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: diagTestOwnerID, Role: "USUARIO"})
		c.Next()
	})
	router.POST("/motorcycles/:id/diagnostics", h.CreateDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])
	assert.Equal(t, encodedBranchID, data["branch_id"])

	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestListDiagnostics_Integration_Success validates GET /motorcycles/:id/diagnostics (HU14).
func TestListDiagnostics_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedMotorcycleID, _ := encoder.Encode(diagTestMotorcycleUUID)
	now := time.Now()

	// Interactor flow: ValidateOwnership → GetByMotorcycleID
	mockSvc.On("ValidateMotorcycleOwnership", mock.Anything, diagTestMotorcycleUUID, diagTestOwnerID).Return(&domain.Motorcycle{
		ID:      diagTestMotorcycleUUID,
		OwnerID: diagTestOwnerID,
	}, nil)
	mockSvc.On("GetByMotorcycleID", mock.Anything, diagTestMotorcycleUUID).Return([]domain.Diagnostic{
		{ID: diagTestDiagnosticUUID, MotorcycleID: diagTestMotorcycleUUID, BranchID: diagTestBranchUUID, Date: now},
		{ID: "a5555555-5555-4000-8000-555555555555", MotorcycleID: diagTestMotorcycleUUID, BranchID: diagTestBranchUUID, Date: now},
	}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: diagTestOwnerID, Role: "USUARIO"})
		c.Next()
	})
	router.GET("/motorcycles/:id/diagnostics", h.ListDiagnostics())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	// Verify IDs are encoded (not raw UUIDs)
	first := data[0].(map[string]interface{})
	assert.NotEqual(t, diagTestDiagnosticUUID, first["id"])
	assert.NotEmpty(t, first["id"])

	mockSvc.AssertExpectations(t)
}

// TestGetDiagnostic_Integration_Success validates GET /motorcycles/:id/diagnostics/:diagnosticId (HU14).
func TestGetDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedMotorcycleID, _ := encoder.Encode(diagTestMotorcycleUUID)
	encodedDiagnosticID, _ := encoder.Encode(diagTestDiagnosticUUID)
	now := time.Now()

	// Interactor flow: GetByID → ValidateOwnership
	mockSvc.On("GetByID", mock.Anything, diagTestDiagnosticUUID).Return(&domain.Diagnostic{
		ID:           diagTestDiagnosticUUID,
		MotorcycleID: diagTestMotorcycleUUID,
		BranchID:     diagTestBranchUUID,
		Date:         now,
	}, nil)
	mockSvc.On("ValidateMotorcycleOwnership", mock.Anything, diagTestMotorcycleUUID, diagTestOwnerID).Return(&domain.Motorcycle{
		ID:      diagTestMotorcycleUUID,
		OwnerID: diagTestOwnerID,
	}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: diagTestOwnerID, Role: "USUARIO"})
		c.Next()
	})
	router.GET("/motorcycles/:id/diagnostics/:diagnosticId", h.GetDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagnosticID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, encodedDiagnosticID, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])

	mockSvc.AssertExpectations(t)
}

// TestUpdateDiagnostic_Integration_Success validates PUT /motorcycles/:id/diagnostics/:diagnosticId (HU12).
func TestUpdateDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedMotorcycleID, _ := encoder.Encode(diagTestMotorcycleUUID)
	encodedDiagnosticID, _ := encoder.Encode(diagTestDiagnosticUUID)
	now := time.Now()
	newDescription := "Motor recalentado"

	// Interactor flow: GetByID → ValidateOwnership → ApplyUpdates → BeginTx → Update → Commit
	existingDiag := &domain.Diagnostic{
		ID:           diagTestDiagnosticUUID,
		MotorcycleID: diagTestMotorcycleUUID,
		BranchID:     diagTestBranchUUID,
		Date:         now,
	}
	mockSvc.On("GetByID", mock.Anything, diagTestDiagnosticUUID).Return(existingDiag, nil)
	mockSvc.On("ValidateMotorcycleOwnership", mock.Anything, diagTestMotorcycleUUID, diagTestOwnerID).Return(&domain.Motorcycle{
		ID:      diagTestMotorcycleUUID,
		OwnerID: diagTestOwnerID,
	}, nil)
	mockSvc.On("ApplyDiagnosticUpdates", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		existing := args.Get(0).(*domain.Diagnostic)
		updates := args.Get(1).(*domain.Diagnostic)
		if updates.ProblemDescription != nil {
			existing.ProblemDescription = updates.ProblemDescription
		}
	})
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("UpdateDiagnostic", mock.Anything, mockTx, mock.Anything).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"problem_description": newDescription,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: diagTestOwnerID, Role: "USUARIO"})
		c.Next()
	})
	router.PUT("/motorcycles/:id/diagnostics/:diagnosticId", h.UpdateDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagnosticID, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, encodedDiagnosticID, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])

	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestDeleteDiagnostic_Integration_Success validates DELETE /motorcycles/:id/diagnostics/:diagnosticId (HU13).
func TestDeleteDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedMotorcycleID, _ := encoder.Encode(diagTestMotorcycleUUID)
	encodedDiagnosticID, _ := encoder.Encode(diagTestDiagnosticUUID)

	// Interactor flow: GetByID → ValidateOwnership → BeginTx → Delete → Commit
	mockSvc.On("GetByID", mock.Anything, diagTestDiagnosticUUID).Return(&domain.Diagnostic{
		ID:           diagTestDiagnosticUUID,
		MotorcycleID: diagTestMotorcycleUUID,
		BranchID:     diagTestBranchUUID,
	}, nil)
	mockSvc.On("ValidateMotorcycleOwnership", mock.Anything, diagTestMotorcycleUUID, diagTestOwnerID).Return(&domain.Motorcycle{
		ID:      diagTestMotorcycleUUID,
		OwnerID: diagTestOwnerID,
	}, nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DeleteDiagnostic", mock.Anything, mockTx, diagTestDiagnosticUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: diagTestOwnerID, Role: "USUARIO"})
		c.Next()
	})
	router.DELETE("/motorcycles/:id/diagnostics/:diagnosticId", h.DeleteDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotorcycleID+"/diagnostics/"+encodedDiagnosticID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestSetDiagnosticSolution_Integration_Success validates PATCH /diagnostics/:id/solution.
func TestSetDiagnosticSolution_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockSvc)

	h := handlers.NewForTestWithConcrete(nil, nil, diagnosticInteractor, nil, encoder, responseHandler)

	encodedDiagnosticID, _ := encoder.Encode(diagTestDiagnosticUUID)
	solution := "Revisar el sistema de encendido"

	// Interactor flow: BeginTx → SetSolution → Commit
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("SetSolution", mock.Anything, mockTx, diagTestDiagnosticUUID, solution).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"possible_solution": solution,
	})

	router := gin.New()
	router.PATCH("/diagnostics/:id/solution", h.SetDiagnosticSolution())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/diagnostics/"+encodedDiagnosticID+"/solution", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
