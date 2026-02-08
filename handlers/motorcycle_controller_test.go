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

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockRepo.On("GetByID", mock.Anything, testMotorcycle.ID).Return(testMotorcycle, nil)

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
	mockRepo.AssertExpectations(t)
}

// TestListMotorcycles_Integration_Success validates GET /motorcycles
func TestListMotorcycles_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
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
	mockRepo.On("GetByOwnerID", mock.Anything, testPersonMoto.ID).Return(motorcycles, nil)

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
	mockRepo.AssertExpectations(t)
}

// TestUpdateMotorcycle_Integration_Success validates PUT /motorcycles/:id
func TestUpdateMotorcycle_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockTx := new(mocks.MockTx)

	year := 2024
	updatedMoto := &domain.Motorcycle{
		ID: testMotorcycle.ID, LicensePlate: "ABC123", OwnerID: testPersonMoto.ID,
		ReferenceID: "ref-id-001", Year: &year,
		Reference: testMotorcycle.Reference,
	}

	mockRepo.On("GetByID", mock.Anything, testMotorcycle.ID).Return(testMotorcycle, nil).Once()
	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("Update", mock.Anything, mockTx, mock.AnythingOfType("*domain.Motorcycle")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockRepo.On("GetByID", mock.Anything, testMotorcycle.ID).Return(updatedMoto, nil).Once()

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
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestDeleteMotorcycle_Integration_Success validates DELETE /motorcycles/:id
func TestDeleteMotorcycle_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	encodedID, _ := encoder.Encode(testMotorcycle.ID)
	mockTx := new(mocks.MockTx)

	mockRepo.On("GetByID", mock.Anything, testMotorcycle.ID).Return(testMotorcycle, nil)
	mockRepo.On("HasServiceHistory", mock.Anything, testMotorcycle.ID).Return(false, nil)
	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("HardDelete", mock.Anything, mockTx, testMotorcycle.ID).Return(nil)
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
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestGetMotorcycleReferences_Integration_Success validates GET /motorcycle-references
func TestGetMotorcycleReferences_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	refs := []domain.MotorcycleReference{
		{ID: "d4444444-4444-4000-8000-444444444444", BrandID: "e5555555-5555-4000-8000-555555555555", BrandName: "Honda", Model: "CB 190R", Category: "Sport", EngineDisplacement: 190},
		{ID: "f6666666-6666-4000-8000-666666666666", BrandID: "d7777777-7777-4000-8000-777777777777", BrandName: "Yamaha", Model: "MT-07", Category: "Naked", EngineDisplacement: 689},
	}
	mockRepo.On("GetAllReferences", mock.Anything).Return(refs, nil)

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
	mockRepo.AssertExpectations(t)
}

// TestGetBrandLines_Integration_Success validates GET /admin/brands/:brandId/lines
func TestGetBrandLines_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockRepo := new(mocks.MockMotorcycleRepository)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)
	h := handlers.NewForTest(nil, nil, motorcycleInteractor, nil, msgCache, encoder, responseHandler)

	brandID := "e5555555-5555-4000-8000-555555555555"
	encodedBrandID, _ := encoder.Encode(brandID)

	refs := []domain.MotorcycleReference{
		{ID: "d4444444-4444-4000-8000-444444444444", BrandID: brandID, BrandName: "Honda", Model: "CB 190R"},
		{ID: "f6666666-6666-4000-8000-666666666666", BrandID: brandID, BrandName: "Honda", Model: "CB 300F"},
	}
	mockRepo.On("GetReferencesByBrandID", mock.Anything, brandID).Return(refs, nil)

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
	mockRepo.AssertExpectations(t)
}
