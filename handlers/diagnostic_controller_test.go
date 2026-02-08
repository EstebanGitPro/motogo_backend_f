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

// TestCreateDiagnostic_Integration_Success validates the full HTTP pipeline
// for POST /motorcycles/:id/diagnostics (HU11).
//
// Exercises: auth → motorcycle ID decoding → bind → sanitize → branch ID decoding →
// interactor (validate motorcycle, ownership, branch, create, tx) → ID encoding → 201 response.
func TestCreateDiagnostic_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Create mocked repo layer for DiagnosticInteractor
	mockDiagnosticRepo := new(mocks.MockDiagnosticRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)

	diagnosticInteractor := interactor.NewDiagnosticInteractor(
		mockDiagnosticRepo,
		mockMotorcycleRepo,
		mockBranchRepo,
	)

	h := handlers.NewForTestWithConcrete(
		nil, nil,
		diagnosticInteractor,
		nil,
		encoder,
		responseHandler,
	)

	// IDs
	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"

	encodedMotorcycleID, err := encoder.Encode(motorcycleUUID)
	assert.NoError(t, err)

	encodedBranchID, err := encoder.Encode(branchUUID)
	assert.NoError(t, err)

	problemDesc := "La moto no enciende"

	// Mock: motorcycle exists and belongs to owner
	mockMotorcycleRepo.On("GetByID", mock.Anything, motorcycleUUID).Return(&domain.Motorcycle{
		ID:           motorcycleUUID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
	}, nil)

	// Mock: branch exists
	mockBranchRepo.On("GetBranchByID", mock.Anything, branchUUID).Return(&domain.Branch{
		ID:   branchUUID,
		Name: "Taller Norte",
	}, nil)

	// Mock: no existing diagnostic for this moto+branch (UPSERT → CREATE path)
	mockDiagnosticRepo.On("GetByMotorcycleAndBranch", mock.Anything, motorcycleUUID, branchUUID).Return(nil, nil)

	// Mock: tx
	mockTx := new(mocks.MockTx)
	mockDiagnosticRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDiagnosticRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Request body
	reqBody := map[string]interface{}{
		"branch_id":           encodedBranchID,
		"problem_description": problemDesc,
	}
	bodyJSON, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   ownerID,
			Role: "USUARIO",
		})
		c.Next()
	})
	router.POST("/motorcycles/:id/diagnostics", h.CreateDiagnostic())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/diagnostics", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])
	assert.Equal(t, encodedBranchID, data["branch_id"])

	mockMotorcycleRepo.AssertExpectations(t)
	mockBranchRepo.AssertExpectations(t)
	mockDiagnosticRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
