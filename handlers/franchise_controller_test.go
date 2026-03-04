package handlers_test

import (
	"bytes"
	"context"
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
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// BuildFranchiseLinks Tests
// ============================================

func TestBuildFranchiseLinks_AllLinksPresent(t *testing.T) {
	links := handlers.BuildFranchiseLinks("http://localhost:8080", "encoded-franchise-123")

	assert.NotEmpty(t, links)
	assert.Len(t, links, 6, "Should have 6 HATEOAS links")

	// Check for expected link relations
	relFound := make(map[string]bool)
	for _, link := range links {
		relFound[link.Rel] = true
	}

	assert.True(t, relFound["self"], "Should have 'self' link")
	assert.True(t, relFound["update"], "Should have 'update' link")
	assert.True(t, relFound["delete"], "Should have 'delete' link")
	assert.True(t, relFound["list"], "Should have 'list' link")
	assert.True(t, relFound["branches"], "Should have 'branches' link")
	assert.True(t, relFound["add-branch"], "Should have 'add-branch' link")
}

func TestBuildFranchiseLinks_CorrectURLs(t *testing.T) {
	franchiseID := "abc123"

	links := handlers.BuildFranchiseLinks("", franchiseID)

	// Verify URLs contain expected patterns (relative paths)
	for _, link := range links {
		assert.Contains(t, link.Href, "/motogo/api/v1", "Link should contain API base path")

		switch link.Rel {
		case "self":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "update":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "PUT", link.Method)
		case "delete":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "DELETE", link.Method)
		case "list":
			assert.Equal(t, "GET", link.Method)
		case "branches":
			assert.Contains(t, link.Href, "franchise_id="+franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "add-branch":
			assert.Contains(t, link.Href, franchiseID+"/branches")
			assert.Equal(t, "POST", link.Method)
		}
	}
}

func TestBuildFranchiseLinks_EmptyID(t *testing.T) {
	links := handlers.BuildFranchiseLinks("http://localhost:8080", "")

	assert.NotEmpty(t, links)
	// Should still return links, just with empty ID in URL
}

func TestBuildFranchiseLinks_DifferentBaseURLs(t *testing.T) {
	// baseURL is now ignored — all links are relative paths
	links := handlers.BuildFranchiseLinks("", "test-id")
	assert.NotEmpty(t, links)
	for _, link := range links {
		assert.Contains(t, link.Href, "/motogo/api/v1")
	}
}

// ============================================
// Integration Tests
// ============================================

// TestRegisterFranchise_Integration_Success validates the full HTTP handler pipeline
// for the success path of RegisterFranchise (HU26).
//
// It exercises: JSON binding → sanitization → branch ID decoding → domain mapping →
// FranchiseInteractor.CreateFranchiseWithBranches → response serialization → HATEOAS links → 201 response.
func TestRegisterFranchise_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies for FranchiseInteractor
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	// Handler only needs IDEncoder + Response (interactor passed as method param)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	// Test data
	ownerID := "a1111111-1111-4000-8000-111111111111"
	branchID := "a2222222-2222-4000-8000-222222222222"
	franchiseID := "a4444444-4444-4000-8000-444444444444"
	description := "Red de talleres premium"

	encodedBranchID, _ := encoder.Encode(branchID)

	// Mock: branch validation (ownership) via franchise service
	mockFranchiseService.On("ValidateBranchesForFranchise", mock.Anything, []string{branchID}, ownerID).Return(nil)

	// Mock: transaction lifecycle
	mockFranchiseService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", mock.Anything, mockTx, mock.AnythingOfType("domain.Franchise")).
		Return(&domain.Franchise{
			ID:          franchiseID,
			Name:        "Red Motos Norte",
			Description: &description,
		}, nil)
	mockFranchiseService.On("AssociateBranches", mock.Anything, mockTx, franchiseID, []string{branchID}).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	// Request body
	reqBody := map[string]interface{}{
		"name":        "Red Motos Norte",
		"description": "Red de talleres premium",
		"branch_ids":  []string{encodedBranchID},
	}
	body, _ := json.Marshal(reqBody)

	// Setup router and execute
	router := gin.New()
	router.POST("/franchises", func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID})
	}, h.RegisterFranchise(franchiseInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/franchises", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_F_REG_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Red Motos Norte", data["name"])
	assert.Equal(t, float64(1), data["branch_count"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["_links"])

	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================
// Franchise controller integration helpers
// ============================================

func createFranchiseMessageCache() *messagingCache.MessageCache {
	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_F_REG_EXI_00001", Type: "EXITO", Title: "Franquicia registrada", Content: "Franquicia registrada exitosamente", Active: true},
		{Code: "MOD_F_GET_EXI_00001", Type: "EXITO", Title: "Franquicia encontrada", Content: "Franquicia encontrada exitosamente", Active: true},
		{Code: "MOD_F_LIST_EXI_00001", Type: "EXITO", Title: "Franquicias listadas", Content: "Franquicias listadas exitosamente", Active: true},
		{Code: "MOD_F_UPD_EXI_00001", Type: "EXITO", Title: "Franquicia actualizada", Content: "Franquicia actualizada exitosamente", Active: true},
		{Code: "MOD_F_DEL_EXI_00001", Type: "EXITO", Title: "Franquicia eliminada", Content: "Franquicia eliminada exitosamente", Active: true},
		{Code: "MOD_F_BRANCH_ADD_EXI_00001", Type: "EXITO", Title: "Sede agregada", Content: "Sede agregada a franquicia exitosamente", Active: true},
		{Code: "MOD_F_BRANCH_REM_EXI_00001", Type: "EXITO", Title: "Sede removida", Content: "Sede removida de franquicia exitosamente", Active: true},
		{Code: "MOD_F_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "No encontrada", Content: "Franquicia no encontrada", Active: true},
		{Code: "MOD_F_DUP_NAME_ERR_00001", Type: "ERROR", Title: "Nombre duplicado", Content: "Ya existe una franquicia con ese nombre", Active: true},
		{Code: "MOD_F_NO_BRANCHES_ERR_00001", Type: "ERROR", Title: "Sin sedes", Content: "Debe incluir al menos una sede", Active: true},
		{Code: "MOD_F_BRANCH_NOT_OWNED_ERR_00001", Type: "ERROR", Title: "Sede ajena", Content: "La sede no le pertenece", Active: true},
		{Code: "MOD_F_MIN_BRANCHES_ERR_00001", Type: "ERROR", Title: "Mínimo sedes", Content: "No se puede remover la última sede", Active: true},
		{Code: "MOD_B_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Sede no encontrada", Content: "La sede no fue encontrada", Active: true},
		{Code: "GEN_SRV_ERR_00001", Type: "ERROR", Title: "Error del servidor", Content: "Error interno del servidor", Active: true},
		{Code: "MOD_V_VAL_ERR_00001", Type: "ERROR", Title: "Formato inválido", Content: "Formato de solicitud inválido", Active: true},
		{Code: "MOD_V_JSON_ERR_00012", Type: "ERROR", Title: "JSON inválido", Content: "El formato JSON de la solicitud es inválido", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	return cache
}

func setupFranchiseRouter(t *testing.T, mockSvc *mocks.MockFranchiseService) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createFranchiseMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	franchiseInteractor := interactor.NewFranchiseInteractor(mockSvc)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, franchiseInteractor, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   "owner-1",
			Role: "REPRESENTANTE",
		})
		c.Next()
	})

	router.GET("/franchises", h.ListFranchises(franchiseInteractor))
	router.GET("/franchises/:id", h.GetFranchise(franchiseInteractor))
	router.PUT("/franchises/:id", h.UpdateFranchise(franchiseInteractor))
	router.DELETE("/franchises/:id", h.DeleteFranchise(franchiseInteractor))
	router.POST("/franchises/:id/branches", h.AddBranchToFranchise(franchiseInteractor))
	router.DELETE("/franchises/:id/branches/:branchId", h.RemoveBranchFromFranchise(franchiseInteractor))

	return router
}

// ============================================
// ListFranchises
// ============================================

func TestListFranchises_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)

	desc := "Red de talleres"
	mockSvc.On("GetFranchisesByRepresentative", mock.Anything, "owner-1").Return([]domain.Franchise{
		{ID: "a1111111-1111-4000-8000-111111111111", Name: "Franquicia 1", Description: &desc},
		{ID: "a2222222-2222-4000-8000-222222222222", Name: "Franquicia 2"},
	}, nil)

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/franchises", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestListFranchises_Error(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)

	mockSvc.On("GetFranchisesByRepresentative", mock.Anything, "owner-1").Return(nil, errors.New("db error"))

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/franchises", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetFranchise
// ============================================

func TestGetFranchise_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	mockSvc.On("GetFranchiseByID", mock.Anything, franchiseUUID).Return(&domain.Franchise{
		ID: franchiseUUID, Name: "Franquicia Test",
	}, nil)

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/franchises/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestGetFranchise_InvalidID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/franchises/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetFranchise_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	mockSvc.On("GetFranchiseByID", mock.Anything, franchiseUUID).Return(nil, domain.ErrFranchiseNotFound)

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/franchises/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// ============================================
// UpdateFranchise
// ============================================

func TestUpdateFranchise_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	// Interactor flow: CountBranches → BeginTx → UpdateFranchise → Commit
	mockSvc.On("CountBranches", mock.Anything, franchiseUUID).Return(2, nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("UpdateFranchise", mock.Anything, mockTx, mock.AnythingOfType("domain.Franchise")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := setupFranchiseRouter(t, mockSvc)

	reqBody := map[string]interface{}{"name": "Updated Franchise"}
	bodyJSON, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/franchises/"+encodedID, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestUpdateFranchise_InvalidID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	router := setupFranchiseRouter(t, mockSvc)

	reqBody := map[string]interface{}{"name": "Updated"}
	bodyJSON, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/franchises/invalid-id", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateFranchise_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	router := setupFranchiseRouter(t, mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/franchises/"+encodedID, bytes.NewBuffer([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// DeleteFranchise
// ============================================

func TestDeleteFranchise_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	// Interactor flow: BeginTx → DissociateBranches → DeleteFranchise → Commit
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranches", mock.Anything, mockTx, franchiseUUID).Return(nil)
	mockSvc.On("DeleteFranchise", mock.Anything, mockTx, franchiseUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/franchises/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestDeleteFranchise_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(franchiseUUID)

	// Interactor flow: BeginTx → DissociateBranches → DeleteFranchise (fails)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateBranches", mock.Anything, mockTx, franchiseUUID).Return(nil)
	mockSvc.On("DeleteFranchise", mock.Anything, mockTx, franchiseUUID).Return(domain.ErrFranchiseNotFound)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := setupFranchiseRouter(t, mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/franchises/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// ============================================
// AddBranchToFranchise
// ============================================

func TestAddBranchToFranchise_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	branchUUID := "a2222222-2222-4000-8000-222222222222"
	encodedFranchiseID, _ := encoder.Encode(franchiseUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	// Interactor flow: ValidateBranchOwnership → BeginTx → AssociateBranches → Commit
	mockSvc.On("ValidateBranchOwnership", mock.Anything, branchUUID, "owner-1").Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("AssociateBranches", mock.Anything, mockTx, franchiseUUID, []string{branchUUID}).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := setupFranchiseRouter(t, mockSvc)

	reqBody := map[string]interface{}{"branch_id": encodedBranchID}
	bodyJSON, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/franchises/"+encodedFranchiseID+"/branches", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestAddBranchToFranchise_InvalidFranchiseID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	router := setupFranchiseRouter(t, mockSvc)

	reqBody := map[string]interface{}{"branch_id": "some-id"}
	bodyJSON, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/franchises/invalid-id/branches", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAddBranchToFranchise_InvalidBranchID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedFranchiseID, _ := encoder.Encode(franchiseUUID)

	router := setupFranchiseRouter(t, mockSvc)

	reqBody := map[string]interface{}{"branch_id": "invalid-branch-id"}
	bodyJSON, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/franchises/"+encodedFranchiseID+"/branches", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// RemoveBranchFromFranchise
// ============================================

func TestRemoveBranchFromFranchise_Success(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	branchUUID := "a2222222-2222-4000-8000-222222222222"
	encodedFranchiseID, _ := encoder.Encode(franchiseUUID)
	encodedBranchID, _ := encoder.Encode(branchUUID)

	// Interactor flow: CanRemoveBranch → BeginTx → DissociateSingleBranch → Commit
	mockSvc.On("CanRemoveBranch", mock.Anything, franchiseUUID).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DissociateSingleBranch", mock.Anything, mockTx, branchUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := setupFranchiseRouter(t, mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/franchises/"+encodedFranchiseID+"/branches/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestRemoveBranchFromFranchise_InvalidFranchiseID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	branchUUID := "a2222222-2222-4000-8000-222222222222"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	router := setupFranchiseRouter(t, mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/franchises/invalid-id/branches/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestRemoveBranchFromFranchise_InvalidBranchID(t *testing.T) {
	mockSvc := new(mocks.MockFranchiseService)
	encoder := createTestEncoder()

	franchiseUUID := "a1111111-1111-4000-8000-111111111111"
	encodedFranchiseID, _ := encoder.Encode(franchiseUUID)

	router := setupFranchiseRouter(t, mockSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/franchises/"+encodedFranchiseID+"/branches/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}
