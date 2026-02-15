package handlers_test

import (
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
// GetMotorcycleCategories Tests
// ============================================

func TestGetMotorcycleCategories_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	categories := []domain.MotorcycleCategory{
		{Name: "Sport", LineCount: 5},
		{Name: "Scooter", LineCount: 3},
		{Name: "Urban", LineCount: 8},
	}
	mockMotoInteractor.On("GetDistinctCategories", mock.Anything).Return(categories, nil)

	router := gin.New()
	router.GET("/motorcycle-categories", h.GetMotorcycleCategories())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycle-categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_MOT_CAT_LIST_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestGetMotorcycleCategories_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	mockMotoInteractor.On("GetDistinctCategories", mock.Anything).Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/motorcycle-categories", h.GetMotorcycleCategories())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycle-categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}

// ============================================
// GetCategoryLines Tests
// ============================================

func TestGetCategoryLines_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	lines := []domain.CategoryLine{
		{Model: "CB 190R", BrandName: "Honda", EngineDisplacement: 190},
		{Model: "Gixxer 250", BrandName: "Suzuki", EngineDisplacement: 249},
	}
	mockMotoInteractor.On("GetLinesByCategory", mock.Anything, "Sport").Return(lines, nil)

	router := gin.New()
	router.GET("/motorcycle-categories/:categoryName/lines", h.GetCategoryLines())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycle-categories/Sport/lines", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_MOT_CAT_LINES_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestGetCategoryLines_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	mockMotoInteractor.On("GetLinesByCategory", mock.Anything, "Unknown").Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/motorcycle-categories/:categoryName/lines", h.GetCategoryLines())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/motorcycle-categories/Unknown/lines", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}

// ============================================
// GetEngineDisplacements Tests
// ============================================

func TestGetEngineDisplacements_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	displacements := []domain.EngineDisplacementRange{
		{Range: domain.DisplacementRangeLow},
		{Range: domain.DisplacementRangeMedium},
		{Range: domain.DisplacementRangeHigh},
	}
	mockMotoInteractor.On("GetDistinctDisplacements", mock.Anything).Return(displacements, nil)

	router := gin.New()
	router.GET("/engine-displacements", h.GetEngineDisplacements())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/engine-displacements", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_MOT_DISP_LIST_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestGetEngineDisplacements_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	mockMotoInteractor.On("GetDistinctDisplacements", mock.Anything).Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/engine-displacements", h.GetEngineDisplacements())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/engine-displacements", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}

// ============================================
// GetRatingRanges Tests
// ============================================

func TestGetRatingRanges_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	ranges := []domain.RatingRange{
		{Value: 1, Label: "Muy malo"},
		{Value: 2, Label: "Malo"},
		{Value: 3, Label: "Regular"},
		{Value: 4, Label: "Bueno"},
		{Value: 5, Label: "Excelente"},
	}
	mockMotoInteractor.On("GetRatingRanges", mock.Anything).Return(ranges, nil)

	router := gin.New()
	router.GET("/rating-ranges", h.GetRatingRanges())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rating-ranges", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_MOT_RATE_LIST_EXI_00001", resp["code"])
	mockMotoInteractor.AssertExpectations(t)
}

func TestGetRatingRanges_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockMotoInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(nil, nil, mockMotoInteractor, nil, msgCache, encoder, responseHandler)

	mockMotoInteractor.On("GetRatingRanges", mock.Anything).Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/rating-ranges", h.GetRatingRanges())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rating-ranges", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockMotoInteractor.AssertExpectations(t)
}
