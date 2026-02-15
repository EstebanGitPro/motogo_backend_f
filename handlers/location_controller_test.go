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
// GetDepartments Tests
// ============================================

func TestGetDepartments_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockLocInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(nil, mockLocInteractor, nil, nil, msgCache, encoder, responseHandler)

	departments := []domain.Department{
		{ID: "a1234567-89ab-cdef-0123-456789abcdef", Name: "Antioquia"},
		{ID: "b1234567-89ab-cdef-0123-456789abcdef", Name: "Cundinamarca"},
	}
	mockLocInteractor.On("GetAllDepartments", mock.Anything).Return(departments, nil)

	router := gin.New()
	router.GET("/departments", h.GetDepartments())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_L_DEP_EXI_00001", resp["code"])
	mockLocInteractor.AssertExpectations(t)
}

func TestGetDepartments_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockLocInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(nil, mockLocInteractor, nil, nil, msgCache, encoder, responseHandler)

	mockLocInteractor.On("GetAllDepartments", mock.Anything).Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/departments", h.GetDepartments())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockLocInteractor.AssertExpectations(t)
}

// ============================================
// GetCitiesByDepartment Tests
// ============================================

func TestGetCitiesByDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockLocInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(nil, mockLocInteractor, nil, nil, msgCache, encoder, responseHandler)

	deptUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedDeptID, _ := encoder.Encode(deptUUID)

	cities := []domain.City{
		{ID: "c1234567-89ab-cdef-0123-456789abcdef", Name: "Medellín", DepartmentID: deptUUID},
		{ID: "d1234567-89ab-cdef-0123-456789abcdef", Name: "Envigado", DepartmentID: deptUUID},
	}
	mockLocInteractor.On("GetCitiesByDepartment", mock.Anything, deptUUID).Return(cities, nil)

	router := gin.New()
	router.GET("/departments/:id/cities", h.GetCitiesByDepartment())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/"+encodedDeptID+"/cities", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "MOD_L_CIT_EXI_00001", resp["code"])
	mockLocInteractor.AssertExpectations(t)
}

func TestGetCitiesByDepartment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockLocInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(nil, mockLocInteractor, nil, nil, msgCache, encoder, responseHandler)

	router := gin.New()
	router.GET("/departments/:id/cities", h.GetCitiesByDepartment())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/INVALID_ID/cities", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetCitiesByDepartment_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockLocInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(nil, mockLocInteractor, nil, nil, msgCache, encoder, responseHandler)

	deptUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedDeptID, _ := encoder.Encode(deptUUID)

	mockLocInteractor.On("GetCitiesByDepartment", mock.Anything, deptUUID).Return(nil, domain.ErrInternalServer)

	router := gin.New()
	router.GET("/departments/:id/cities", h.GetCitiesByDepartment())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/"+encodedDeptID+"/cities", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockLocInteractor.AssertExpectations(t)
}
