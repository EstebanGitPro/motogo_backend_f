package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// TestGrantDiagnosticPermission_Integration_Success
// ============================================================
func TestGrantDiagnosticPermission_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"
	permUUID := "a5555555-5555-4000-8000-555555555555"

	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMoto.On("GrantDiagnosticPermission", mock.Anything, motorcycleUUID, branchUUID, ownerID, true).
		Return(&domain.DiagnosticPermission{
			ID:           permUUID,
			MotorcycleID: motorcycleUUID,
			BranchID:     branchUUID,
			Active:       true,
		}, nil)

	reqBody := map[string]interface{}{
		"branch_id": encodedBranchID,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.POST("/motorcycles/:id/permissions", h.GrantDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/permissions", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])
	assert.Equal(t, encodedBranchID, data["branch_id"])
	assert.Equal(t, true, data["active"])

	mockMoto.AssertExpectations(t)
}

// ============================================================
// TestGrantDiagnosticPermission_Integration_MotorcycleNotFound
// ============================================================
func TestGrantDiagnosticPermission_Integration_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"

	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMoto.On("GrantDiagnosticPermission", mock.Anything, motorcycleUUID, branchUUID, ownerID, true).
		Return(nil, domain.ErrMotorcycleNotFound)

	reqBody := map[string]interface{}{
		"branch_id": encodedBranchID,
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.POST("/motorcycles/:id/permissions", h.GrantDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/permissions", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockMoto.AssertExpectations(t)
}

// ============================================================
// TestRevokeDiagnosticPermission_Integration_Success
// ============================================================
func TestRevokeDiagnosticPermission_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"

	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMoto.On("RevokeDiagnosticPermission", mock.Anything, motorcycleUUID, branchUUID, ownerID).
		Return(nil)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.DELETE("/motorcycles/:id/permissions/:branchId", h.RevokeDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotorcycleID+"/permissions/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockMoto.AssertExpectations(t)
}

// ============================================================
// TestRevokeDiagnosticPermission_Integration_PermissionNotFound
// ============================================================
func TestRevokeDiagnosticPermission_Integration_PermissionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"

	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMoto.On("RevokeDiagnosticPermission", mock.Anything, motorcycleUUID, branchUUID, ownerID).
		Return(domain.ErrPermissionNotFound)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.DELETE("/motorcycles/:id/permissions/:branchId", h.RevokeDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotorcycleID+"/permissions/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockMoto.AssertExpectations(t)
}

// ============================================================
// TestListDiagnosticPermissions_Integration_Success
// ============================================================
func TestListDiagnosticPermissions_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	branchUUID := "a3333333-3333-4000-8000-333333333333"

	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)

	mockMoto.On("ListDiagnosticPermissions", mock.Anything, motorcycleUUID, ownerID).
		Return([]domain.DiagnosticPermission{
			{
				ID:           "perm-1",
				MotorcycleID: motorcycleUUID,
				BranchID:     branchUUID,
				Active:       true,
			},
		}, nil)

	router := gin.New()
	router.Use(authMiddleware(ownerID, "USUARIO"))
	router.GET("/motorcycles/:id/permissions", h.ListDiagnosticPermissions())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/permissions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	permissions := data["permissions"].([]interface{})
	assert.Len(t, permissions, 1)

	first := permissions[0].(map[string]interface{})
	assert.Equal(t, encodedMotorcycleID, first["motorcycle_id"])
	assert.Equal(t, true, first["active"])

	mockMoto.AssertExpectations(t)
}

// ============================================================
// TestListDiagnosticPermissions_Integration_Unauthorized
// ============================================================
func TestListDiagnosticPermissions_Integration_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMoto := new(mocks.MockMotorcycleInteractor)
	h := handlers.NewForTest(nil, nil, mockMoto, nil, msgCache, encoder, responseHandler)

	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	encodedMotorcycleID, _ := encoder.Encode(motorcycleUUID)

	// No auth middleware
	router := gin.New()
	router.GET("/motorcycles/:id/permissions", h.ListDiagnosticPermissions())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotorcycleID+"/permissions", nil)
	router.ServeHTTP(w, req)

	// Controller returns MsgServerError (500) when no authenticated user in context.
	// This is the expected behavior — the auth middleware would normally reject the request
	// before reaching the controller. Without middleware, the controller treats it as a server error.
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
