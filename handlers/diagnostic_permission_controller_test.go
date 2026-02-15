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

// ============================================
// GrantDiagnosticPermission Tests
// ============================================

func TestGrantDiagnosticPermission_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMotoInteractor.On("GrantDiagnosticPermission",
		mock.Anything, motoUUID, branchUUID, ownerID, true).
		Return(&domain.DiagnosticPermission{
			ID:           "perm-1",
			MotorcycleID: motoUUID,
			BranchID:     branchUUID,
			Active:       true,
		}, nil)

	body, _ := json.Marshal(map[string]interface{}{
		"branch_id": encodedBranchID,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.POST("/motorcycles/:id/permissions", h.GrantDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotoID+"/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_DGP_GRANT_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestGrantDiagnosticPermission_InvalidMotoID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "owner-1", Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.POST("/motorcycles/:id/permissions", h.GrantDiagnosticPermission())

	body, _ := json.Marshal(map[string]interface{}{"branch_id": "enc-branch"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/BAD_ID/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGrantDiagnosticPermission_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMotoInteractor.On("GrantDiagnosticPermission",
		mock.Anything, motoUUID, branchUUID, ownerID, true).
		Return(nil, domain.ErrMotorcycleNotFound)

	body, _ := json.Marshal(map[string]interface{}{
		"branch_id": encodedBranchID,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.POST("/motorcycles/:id/permissions", h.GrantDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotoID+"/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}

// ============================================
// RevokeDiagnosticPermission Tests
// ============================================

func TestRevokeDiagnosticPermission_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMotoInteractor.On("RevokeDiagnosticPermission",
		mock.Anything, motoUUID, branchUUID, ownerID).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.DELETE("/motorcycles/:id/permissions/:branchId", h.RevokeDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotoID+"/permissions/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_DGP_REVOKE_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestRevokeDiagnosticPermission_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockMotoInteractor.On("RevokeDiagnosticPermission",
		mock.Anything, motoUUID, branchUUID, ownerID).Return(domain.ErrPermissionNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.DELETE("/motorcycles/:id/permissions/:branchId", h.RevokeDiagnosticPermission())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedMotoID+"/permissions/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}

// ============================================
// ListDiagnosticPermissions Tests
// ============================================

func TestListDiagnosticPermissions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)

	permissions := []domain.DiagnosticPermission{
		{ID: "perm-1", MotorcycleID: motoUUID, BranchID: branchUUID, Active: true},
	}
	mockMotoInteractor.On("ListDiagnosticPermissions",
		mock.Anything, motoUUID, ownerID).Return(permissions, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.GET("/motorcycles/:id/permissions", h.ListDiagnosticPermissions())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotoID+"/permissions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_DGP_LIST_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestListDiagnosticPermissions_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	motoUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	ownerID := "owner-123"
	encodedMotoID, _ := encoder.Encode(motoUUID)

	mockMotoInteractor.On("ListDiagnosticPermissions",
		mock.Anything, motoUUID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID, Role: domain.Role("USUARIO")})
		c.Next()
	})
	router.GET("/motorcycles/:id/permissions", h.ListDiagnosticPermissions())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedMotoID+"/permissions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}
