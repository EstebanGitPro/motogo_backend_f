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

var testPersonMoto = &domain.Person{
	ID:    "a1111111-1111-4000-8000-111111111111",
	Email: "owner@example.com",
	Role:  "user",
}

var testMotorcycle = &domain.Motorcycle{
	ID:           "b2222222-2222-4000-8000-222222222222",
	LicensePlate: "ABC123",
	ReferenceID:  "ref-id-001",
	OwnerID:      "a1111111-1111-4000-8000-111111111111",
	Year:         intPtr(2023),
	Reference: &domain.MotorcycleReference{
		ID:                 "ref-id-001",
		BrandID:            "brand-id-001",
		BrandName:          "Honda",
		Model:              "CB 190R",
		Category:           "Sport",
		EngineDisplacement: 190,
	},
}

func intPtr(v int) *int { return &v }

// TestGetMotorcycle_Integration_Success validates GET /motorcycles/:id
func TestGetMotorcycle_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockSvc.On("GetByID", mock.Anything, testMotorcycle.ID).Return(testMotorcycle, nil)

	router := gin.New()
	router.GET("/motorcycles/:id", func(c *gin.Context) {
		c.Set("authenticated_user", testPersonMoto)
	}, h.GetMotorcycle())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_GET_EXI_00001", resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "ABC123", data["license_plate"])
	assert.NotNil(t, data["reference"])
	assert.NotEmpty(t, data["_links"])
	mockSvc.AssertExpectations(t)
}

// TestListMotorcycles_Integration_Success validates GET /motorcycles
func TestListMotorcycles_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	motorcycles := []domain.Motorcycle{
		{
			ID: "b2222222-2222-4000-8000-222222222222", LicensePlate: "ABC123",
			OwnerID: testPersonMoto.ID,
			Reference: &domain.MotorcycleReference{
				ID: "ref-id-001", BrandID: "brand-id-001", BrandName: "Honda", Model: "CB 190R",
			},
		},
		{
			ID: "c3333333-3333-4000-8000-333333333333", LicensePlate: "XYZ789",
			OwnerID: testPersonMoto.ID,
		},
	}
	mockSvc.On("GetByOwnerID", mock.Anything, testPersonMoto.ID).Return(motorcycles, nil)

	router := gin.New()
	router.GET("/motorcycles", func(c *gin.Context) {
		c.Set("authenticated_user", testPersonMoto)
	}, h.ListMotorcycles())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_LIST_EXI_00001", resp["code"])
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
	mockSvc.AssertExpectations(t)
}

// TestUpdateMotorcycle_Integration_Success validates PUT /motorcycles/:id
func TestUpdateMotorcycle_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockTx := new(mocks.MockTx)

	year := 2024
	updatedMoto := &domain.Motorcycle{
		ID: testMotorcycle.ID, LicensePlate: "ABC123", OwnerID: testPersonMoto.ID,
		ReferenceID: "ref-id-001", Year: &year,
		Reference: testMotorcycle.Reference,
	}

	mockSvc.On("ValidateOwnership", mock.Anything, testMotorcycle.ID, testPersonMoto.ID).Return(testMotorcycle, nil)
	mockSvc.On("ApplyUpdates", testMotorcycle, mock.AnythingOfType("*domain.Motorcycle")).Return(nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("UpdateMotorcycle", mock.Anything, mockTx, testMotorcycle).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockSvc.On("GetByID", mock.Anything, testMotorcycle.ID).Return(updatedMoto, nil)

	reqBody := map[string]interface{}{"year": 2024}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/motorcycles/:id", func(c *gin.Context) {
		c.Set("authenticated_user", testPersonMoto)
	}, h.UpdateMotorcycle())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/motorcycles/"+encodedID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_UPDATE_EXI_00001", resp["code"])
	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestDeleteMotorcycle_Integration_Success validates DELETE /motorcycles/:id
func TestDeleteMotorcycle_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockTx := new(mocks.MockTx)

	mockSvc.On("ValidateOwnership", mock.Anything, testMotorcycle.ID, testPersonMoto.ID).Return(testMotorcycle, nil)
	mockSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSvc.On("DeleteMotorcycle", mock.Anything, mockTx, testMotorcycle.ID).Return(nil)
	mockTx.On("Commit").Return(nil)

	router := gin.New()
	router.DELETE("/motorcycles/:id", func(c *gin.Context) {
		c.Set("authenticated_user", testPersonMoto)
	}, h.DeleteMotorcycle())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/motorcycles/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_DELETE_EXI_00001", resp["code"])
	mockSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestGetMotorcycleReferences_Integration_Success validates GET /motorcycle-references
func TestGetMotorcycleReferences_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	refs := []domain.MotorcycleReference{
		{ID: "d4444444-4444-4000-8000-444444444444", BrandID: "e5555555-5555-4000-8000-555555555555", BrandName: "Honda", Model: "CB 190R", Category: "Sport", EngineDisplacement: 190},
		{ID: "f6666666-6666-4000-8000-666666666666", BrandID: "d7777777-7777-4000-8000-777777777777", BrandName: "Yamaha", Model: "MT-07", Category: "Naked", EngineDisplacement: 689},
	}
	mockSvc.On("GetAllReferences", mock.Anything).Return(refs, nil)

	router := gin.New()
	router.GET("/motorcycle-references", h.GetMotorcycleReferences())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycle-references", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_REF_LIST_EXI_00001", resp["code"])
	data := resp["data"].(map[string]interface{})
	references := data["references"].([]interface{})
	assert.Len(t, references, 2)
	mockSvc.AssertExpectations(t)
}

// TestGetBrandLines_Integration_Success validates GET /admin/brands/:brandId/lines
func TestGetBrandLines_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	brandID := "e5555555-5555-4000-8000-555555555555"
	encodedBrandID, _ := encoder.Encode(brandID)

	refs := []domain.MotorcycleReference{
		{ID: "d4444444-4444-4000-8000-444444444444", BrandID: brandID, BrandName: "Honda", Model: "CB 190R"},
		{ID: "f6666666-6666-4000-8000-666666666666", BrandID: brandID, BrandName: "Honda", Model: "CB 300F"},
	}
	mockSvc.On("GetReferencesByBrandID", mock.Anything, brandID).Return(refs, nil)

	router := gin.New()
	router.GET("/admin/brands/:brandId/lines", h.GetBrandLines())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/brands/"+encodedBrandID+"/lines", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_BRAND_LINES_EXI_00001", resp["code"])
	data := resp["data"].(map[string]interface{})
	lines := data["lines"].([]interface{})
	assert.Len(t, lines, 2)
	firstLine := lines[0].(map[string]interface{})
	assert.Equal(t, "Honda", firstLine["brand_name"])
	assert.Equal(t, "CB 190R", firstLine["model"])
	mockSvc.AssertExpectations(t)
}

// ============================================
// LookupMotorcycleByPlate — Full Flow (covers buildPermittedBranches, fetchAndEncodeDiagnostics,
//   encodeDiagnosticIDs, fetchAndEncodeEvidence)
// ============================================

var testRepPerson = &domain.Person{
	ID:    "a1111111-1111-4000-8000-111111111111",
	Email: "rep@workshop.com",
	Role:  "representante",
}

func TestLookupMotorcycleByPlate_Integration_FullFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// — Mocks for concrete interactors —
	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)

	mockDiagSvc := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockDiagSvc)

	// — Mocks for interface interactors —
	mockMotoInt := new(mocks.MockMotorcycleInteractor)
	mockEvidInt := new(mocks.MockEvidenceInteractor)

	h := handlers.NewForTestWithLookup(branchInteractor, mockMotoInt, diagnosticInteractor, mockEvidInt, encoder, responseHandler)

	// --- Fixture data ---
	branchID1 := "c3333333-3333-4000-8000-333333333333"
	branchID2 := "d4444444-4444-4000-8000-444444444444"
	motoID := "b2222222-2222-4000-8000-222222222222"
	diagID := "e5555555-5555-4000-8000-555555555555"
	evidID := "f6666666-6666-4000-8000-666666666666"
	motoEvidID := "a7777777-7777-4000-8000-777777777777"

	motorcycle := &domain.Motorcycle{
		ID:           motoID,
		LicensePlate: "ABC123",
		OwnerID:      "other-owner-id",
	}

	repBranches := []domain.Branch{
		{ID: branchID1, Name: "Taller Norte"},
		{ID: branchID2, Name: "Taller Sur"},
	}

	permissions := []domain.DiagnosticPermission{
		{ID: "p1", MotorcycleID: motoID, BranchID: branchID1, Active: true},
		{ID: "p2", MotorcycleID: motoID, BranchID: branchID2, Active: false}, // inactive
	}

	diagnostics := []domain.Diagnostic{
		{
			ID:           diagID,
			MotorcycleID: motoID,
			BranchID:     branchID1,
			Evidence: []domain.DiagnosticEvidence{
				{ID: evidID, DiagnosticID: diagID, ImageURL: "https://example.com/photo.jpg"},
			},
		},
	}

	motoEvidences := []domain.MotorcycleEvidence{
		{ID: motoEvidID, MotorcycleID: motoID, ImageURL: "https://example.com/evid.jpg"},
	}

	// --- Setup mock expectations ---
	// 1) GetBranchesByRepresentative
	mockBranchSvc.On("GetBranchesByRepresentative", mock.Anything, testRepPerson.ID).Return(repBranches, nil)

	// 2) GetMotorcycleByLicensePlate
	mockMotoInt.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(motorcycle, nil)

	// 3) LookupPermissions
	mockMotoInt.On("LookupPermissions", mock.Anything, motoID).Return(permissions, nil)

	// 4) ListDiagnosticsByMotorcycleID (concrete DiagnosticInteractor delegates to service)
	mockDiagSvc.On("GetByMotorcycleID", mock.Anything, motoID).Return(diagnostics, nil)
	mockDiagSvc.On("LoadEvidenceForDiagnostics", mock.Anything, mock.AnythingOfType("[]domain.Diagnostic")).Return(nil)

	// 5) LookupEvidence (interface)
	mockEvidInt.On("LookupEvidence", mock.Anything, motoID).Return(motoEvidences, nil)

	// --- Execute ---
	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup?plate=ABC123", nil)
	router.ServeHTTP(w, req)

	// --- Assertions ---
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "MOD_MOT_GET_EXI_00001", resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "ABC123", data["license_plate"])

	// Permitted branches: only branchID1 (active), not branchID2 (inactive)
	permBranches := data["permitted_branches"].([]interface{})
	assert.Len(t, permBranches, 1)
	pb := permBranches[0].(map[string]interface{})
	assert.Equal(t, "Taller Norte", pb["name"])

	// Diagnostics: 1 diagnostic from branchID1
	diagList := data["diagnostics"].([]interface{})
	assert.Len(t, diagList, 1)
	d := diagList[0].(map[string]interface{})
	// IDs should be encoded (not raw UUIDs)
	assert.NotEqual(t, diagID, d["id"])
	assert.NotEqual(t, branchID1, d["branch_id"])

	// Evidence within diagnostic
	diagEvidence := d["evidence"].([]interface{})
	assert.Len(t, diagEvidence, 1)

	// Motorcycle evidence
	motoEvid := data["evidence"].([]interface{})
	assert.Len(t, motoEvid, 1)

	mockBranchSvc.AssertExpectations(t)
	mockMotoInt.AssertExpectations(t)
	mockDiagSvc.AssertExpectations(t)
	mockEvidInt.AssertExpectations(t)
}

// TestLookupMotorcycleByPlate_MissingPlate validates missing plate parameter → error response
func TestLookupMotorcycleByPlate_MissingPlate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	h := handlers.NewForTestWithLookup(nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup", nil) // no ?plate=
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "MOD_MOT_PLATE_REQ_ERR_00001", resp["code"])
}

// TestLookupMotorcycleByPlate_MotorcycleNotFound validates 404 when plate doesn't exist
func TestLookupMotorcycleByPlate_MotorcycleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	mockMotoInt := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTestWithLookup(branchInteractor, mockMotoInt, nil, nil, encoder, responseHandler)

	mockBranchSvc.On("GetBranchesByRepresentative", mock.Anything, testRepPerson.ID).Return([]domain.Branch{}, nil)
	mockMotoInt.On("GetMotorcycleByLicensePlate", mock.Anything, "XYZ999").Return((*domain.Motorcycle)(nil), domain.ErrMotorcycleNotFound)

	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup?plate=XYZ999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "MOD_MOT_NOT_FOUND_ERR_00001", resp["code"])
	mockBranchSvc.AssertExpectations(t)
	mockMotoInt.AssertExpectations(t)
}

// TestLookupMotorcycleByPlate_Forbidden validates 403 when rep has no intersection with permissions
func TestLookupMotorcycleByPlate_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	mockMotoInt := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTestWithLookup(branchInteractor, mockMotoInt, nil, nil, encoder, responseHandler)

	motoID := "b2222222-2222-4000-8000-222222222222"
	motorcycle := &domain.Motorcycle{ID: motoID, LicensePlate: "ABC123", OwnerID: "other-owner"}

	repBranches := []domain.Branch{
		{ID: "branch-A", Name: "Taller A"},
	}
	// Permission exists but for a different branch
	permissions := []domain.DiagnosticPermission{
		{ID: "p1", MotorcycleID: motoID, BranchID: "branch-B", Active: true},
	}

	mockBranchSvc.On("GetBranchesByRepresentative", mock.Anything, testRepPerson.ID).Return(repBranches, nil)
	mockMotoInt.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(motorcycle, nil)
	mockMotoInt.On("LookupPermissions", mock.Anything, motoID).Return(permissions, nil)

	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup?plate=ABC123", nil)
	router.ServeHTTP(w, req)

	assert.True(t, w.Code >= 400) // Error status code from ResponseHandler
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "MOD_MOT_NO_PERMISSION_ERR_00001", resp["code"])
	mockBranchSvc.AssertExpectations(t)
	mockMotoInt.AssertExpectations(t)
}

// TestLookupMotorcycleByPlate_DiagnosticError validates non-fatal diagnostic error (still 200)
func TestLookupMotorcycleByPlate_DiagnosticError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	mockDiagSvc := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockDiagSvc)
	mockMotoInt := new(mocks.MockMotorcycleInteractor)
	mockEvidInt := new(mocks.MockEvidenceInteractor)

	h := handlers.NewForTestWithLookup(branchInteractor, mockMotoInt, diagnosticInteractor, mockEvidInt, encoder, responseHandler)

	branchID := "c3333333-3333-4000-8000-333333333333"
	motoID := "b2222222-2222-4000-8000-222222222222"
	motorcycle := &domain.Motorcycle{ID: motoID, LicensePlate: "ABC123", OwnerID: "other-owner"}

	mockBranchSvc.On("GetBranchesByRepresentative", mock.Anything, testRepPerson.ID).Return(
		[]domain.Branch{{ID: branchID, Name: "Taller"}}, nil)
	mockMotoInt.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(motorcycle, nil)
	mockMotoInt.On("LookupPermissions", mock.Anything, motoID).Return(
		[]domain.DiagnosticPermission{{ID: "p1", MotorcycleID: motoID, BranchID: branchID, Active: true}}, nil)
	// Diagnostic service returns error — non-fatal
	mockDiagSvc.On("GetByMotorcycleID", mock.Anything, motoID).Return(
		([]domain.Diagnostic)(nil), errors.New("db error"))
	// Evidence succeeds
	mockEvidInt.On("LookupEvidence", mock.Anything, motoID).Return([]domain.MotorcycleEvidence{}, nil)

	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup?plate=ABC123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"]) // Non-fatal: returns 200 without diagnostics
	data := resp["data"].(map[string]interface{})
	assert.Nil(t, data["diagnostics"]) // Diagnostics null due to error
	mockBranchSvc.AssertExpectations(t)
	mockMotoInt.AssertExpectations(t)
	mockDiagSvc.AssertExpectations(t)
	mockEvidInt.AssertExpectations(t)
}

// TestLookupMotorcycleByPlate_EvidenceError validates non-fatal evidence error (still 200)
func TestLookupMotorcycleByPlate_EvidenceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchSvc := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchSvc)
	mockDiagSvc := new(mocks.MockDiagnosticService)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(mockDiagSvc)
	mockMotoInt := new(mocks.MockMotorcycleInteractor)
	mockEvidInt := new(mocks.MockEvidenceInteractor)

	h := handlers.NewForTestWithLookup(branchInteractor, mockMotoInt, diagnosticInteractor, mockEvidInt, encoder, responseHandler)

	branchID := "c3333333-3333-4000-8000-333333333333"
	motoID := "b2222222-2222-4000-8000-222222222222"
	motorcycle := &domain.Motorcycle{ID: motoID, LicensePlate: "ABC123", OwnerID: "other-owner"}

	mockBranchSvc.On("GetBranchesByRepresentative", mock.Anything, testRepPerson.ID).Return(
		[]domain.Branch{{ID: branchID, Name: "Taller"}}, nil)
	mockMotoInt.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(motorcycle, nil)
	mockMotoInt.On("LookupPermissions", mock.Anything, motoID).Return(
		[]domain.DiagnosticPermission{{ID: "p1", MotorcycleID: motoID, BranchID: branchID, Active: true}}, nil)
	// Diagnostics succeed (no diagnostics)
	mockDiagSvc.On("GetByMotorcycleID", mock.Anything, motoID).Return([]domain.Diagnostic{}, nil)
	mockDiagSvc.On("LoadEvidenceForDiagnostics", mock.Anything, mock.AnythingOfType("[]domain.Diagnostic")).Return(nil)
	// Evidence returns error — non-fatal
	mockEvidInt.On("LookupEvidence", mock.Anything, motoID).Return(
		([]domain.MotorcycleEvidence)(nil), errors.New("storage error"))

	router := gin.New()
	router.GET("/motorcycles/lookup", func(c *gin.Context) {
		c.Set("authenticated_user", testRepPerson)
	}, h.LookupMotorcycleByPlate())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycles/lookup?plate=ABC123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"]) // Non-fatal: returns 200 without evidence
	data := resp["data"].(map[string]interface{})
	assert.Nil(t, data["evidence"]) // Evidence null due to error
	mockBranchSvc.AssertExpectations(t)
	mockMotoInt.AssertExpectations(t)
	mockDiagSvc.AssertExpectations(t)
	mockEvidInt.AssertExpectations(t)
}
