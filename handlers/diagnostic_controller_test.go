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

	// Create mocked service layer for DiagnosticInteractor
	mockDiagnosticSvc := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)

	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockDiagnosticSvc)

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

	// Mock: ownership validation
	mockDiagnosticSvc.On("ValidateMotorcycleOwnership", mock.Anything, motorcycleUUID, ownerID).Return(&domain.Motorcycle{
		ID:           motorcycleUUID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
	}, nil)

	// Mock: branch exists
	mockDiagnosticSvc.On("ValidateBranchExists", mock.Anything, branchUUID).Return(nil)

	// Mock: tx
	mockDiagnosticSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)

	// Mock: RegisterOrUpdateDiagnostic returns created diagnostic
	createdDiag := &domain.Diagnostic{
		ID:           "a4444444-4444-4000-8000-444444444444",
		MotorcycleID: motorcycleUUID,
		BranchID:     branchUUID,
	}
	mockDiagnosticSvc.On("RegisterOrUpdateDiagnostic", mock.Anything, mockTx, motorcycleUUID, branchUUID, &problemDesc, []string(nil)).Return(createdDiag, nil)
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

	mockDiagnosticSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
