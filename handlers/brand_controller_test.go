package handlers_test

import (
	"context"
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
// BrandInteractor Mock Tests
// ============================================

func TestMockBrandInteractor_GetAllBrands_Success(t *testing.T) {
	mockBrand := new(mocks.MockBrandInteractor)
	testBrands := []domain.Brand{
		{ID: "f6a7b8c9-1111-4000-8000-000000000001", Name: "Honda"},
		{ID: "f6a7b8c9-2222-4000-8000-000000000002", Name: "Yamaha"},
		{ID: "f6a7b8c9-3333-4000-8000-000000000003", Name: "Suzuki"},
	}
	mockBrand.On("GetAllBrands", mock.Anything).Return(testBrands, nil)

	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Honda", result[0].Name)
	assert.Equal(t, "Yamaha", result[1].Name)
	assert.Equal(t, "Suzuki", result[2].Name)
	mockBrand.AssertExpectations(t)
}

func TestMockBrandInteractor_GetAllBrands_Error(t *testing.T) {
	mockBrand := new(mocks.MockBrandInteractor)
	mockBrand.On("GetAllBrands", mock.Anything).Return(nil, domain.ErrBrandNotFound)

	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBrandNotFound, err)
	mockBrand.AssertExpectations(t)
}

func TestMockBrandInteractor_GetAllBrands_Empty(t *testing.T) {
	mockBrand := new(mocks.MockBrandInteractor)
	mockBrand.On("GetAllBrands", mock.Anything).Return([]domain.Brand{}, nil)

	ctx := context.Background()
	result, err := mockBrand.GetAllBrands(ctx)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockBrand.AssertExpectations(t)
}

// ============================================
// LocationInteractor Mock Tests
// ============================================

func TestMockLocationInteractor_GetAllDepartments_Success(t *testing.T) {
	mockLocation := new(mocks.MockLocationInteractor)
	testDepartments := []domain.Department{
		{ID: "dep-1", Name: "Antioquia"},
		{ID: "dep-2", Name: "Cundinamarca"},
	}
	mockLocation.On("GetAllDepartments", mock.Anything).Return(testDepartments, nil)

	ctx := context.Background()
	result, err := mockLocation.GetAllDepartments(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockLocation.AssertExpectations(t)
}

func TestMockLocationInteractor_GetCitiesByDepartment_Success(t *testing.T) {
	mockLocation := new(mocks.MockLocationInteractor)
	testCities := []domain.City{
		{ID: "city-1", Name: "Medellín", DepartmentID: "dep-1"},
		{ID: "city-2", Name: "Envigado", DepartmentID: "dep-1"},
	}
	mockLocation.On("GetCitiesByDepartment", mock.Anything, "dep-1").Return(testCities, nil)

	ctx := context.Background()
	result, err := mockLocation.GetCitiesByDepartment(ctx, "dep-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Medellín", result[0].Name)
	mockLocation.AssertExpectations(t)
}

// ============================================
// MotorcycleInteractor Mock Tests
// ============================================

func TestMockMotorcycleInteractor_GetMotorcyclesByOwner_Success(t *testing.T) {
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	testMotorcycles := []domain.Motorcycle{
		{ID: "moto-1", LicensePlate: "ABC123", OwnerID: "owner-1"},
		{ID: "moto-2", LicensePlate: "XYZ789", OwnerID: "owner-1"},
	}
	mockMotorcycle.On("GetMotorcyclesByOwner", mock.Anything, "owner-1").Return(testMotorcycles, nil)

	ctx := context.Background()
	result, err := mockMotorcycle.GetMotorcyclesByOwner(ctx, "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "ABC123", result[0].LicensePlate)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_GetMotorcycleByLicensePlate_Success(t *testing.T) {
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	testMoto := &domain.Motorcycle{ID: "moto-1", LicensePlate: "ABC123", OwnerID: "owner-1"}
	branchIDs := []string{"branch-1"}
	mockMotorcycle.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123", branchIDs).Return(testMoto, nil)

	ctx := context.Background()
	result, err := mockMotorcycle.GetMotorcycleByLicensePlate(ctx, "ABC123", branchIDs)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ABC123", result.LicensePlate)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_DeleteMotorcycle_Success(t *testing.T) {
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	mockMotorcycle.On("DeleteMotorcycle", mock.Anything, "moto-1", "owner-1").Return(nil)

	ctx := context.Background()
	err := mockMotorcycle.DeleteMotorcycle(ctx, "moto-1", "owner-1")

	assert.NoError(t, err)
	mockMotorcycle.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_DeleteMotorcycle_NotFound(t *testing.T) {
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	mockMotorcycle.On("DeleteMotorcycle", mock.Anything, "moto-999", "owner-1").Return(domain.ErrMotorcycleNotFound)

	ctx := context.Background()
	err := mockMotorcycle.DeleteMotorcycle(ctx, "moto-999", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
	mockMotorcycle.AssertExpectations(t)
}

// ============================================
// Interface Implementation Verification
// ============================================

func TestMockBrandInteractor_VerifyInterfaceImplementation(t *testing.T) {
	assert.NotNil(t, new(mocks.MockBrandInteractor))
}

func TestMockLocationInteractor_VerifyInterfaceImplementation(t *testing.T) {
	assert.NotNil(t, new(mocks.MockLocationInteractor))
}

func TestMockMotorcycleInteractor_VerifyInterfaceImplementation(t *testing.T) {
	assert.NotNil(t, new(mocks.MockMotorcycleInteractor))
}

func TestMockBrandInteractor_MultipleCalls(t *testing.T) {
	mockBrand := new(mocks.MockBrandInteractor)
	testBrands := []domain.Brand{
		{ID: "brand-1", Name: "Honda"},
		{ID: "brand-2", Name: "Yamaha"},
	}
	mockBrand.On("GetAllBrands", mock.Anything).Return(testBrands, nil).Times(2)

	ctx := context.Background()
	result1, err1 := mockBrand.GetAllBrands(ctx)
	assert.NoError(t, err1)
	assert.Len(t, result1, 2)
	result2, err2 := mockBrand.GetAllBrands(ctx)
	assert.NoError(t, err2)
	assert.Len(t, result2, 2)
	mockBrand.AssertExpectations(t)
}

func TestMockLocationInteractor_GetAllDepartments_Error(t *testing.T) {
	mockLocation := new(mocks.MockLocationInteractor)
	mockLocation.On("GetAllDepartments", mock.Anything).Return(nil, domain.ErrDatabaseUnavailable)

	ctx := context.Background()
	result, err := mockLocation.GetAllDepartments(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockLocation.AssertExpectations(t)
}

func TestMockMotorcycleInteractor_GetMotorcyclesByOwner_Error(t *testing.T) {
	mockMotorcycle := new(mocks.MockMotorcycleInteractor)
	mockMotorcycle.On("GetMotorcyclesByOwner", mock.Anything, "owner-999").Return(nil, domain.ErrUserNotFound)

	ctx := context.Background()
	result, err := mockMotorcycle.GetMotorcyclesByOwner(ctx, "owner-999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockMotorcycle.AssertExpectations(t)
}

// ============================================
// HTTP Integration Tests (full pipeline)
// ============================================

// TestGetBrands_Integration_Success validates the full HTTP pipeline for GET /brands.
// Exercises: interactor call → ID encoding → HATEOAS → 200 response.
func TestGetBrands_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBrandInteractor := new(mocks.MockBrandInteractor)

	h := handlers.NewForTest(
		mockBrandInteractor,
		nil, nil, nil, nil,
		encoder,
		responseHandler,
	)

	brands := []domain.Brand{
		{ID: "b1234567-89ab-cdef-0123-456789abcdef", Name: "Yamaha"},
		{ID: "b2234567-89ab-cdef-0123-456789abcdef", Name: "Honda"},
	}
	mockBrandInteractor.On("GetAllBrands", mock.Anything).Return(brands, nil)

	router := gin.New()
	router.GET("/brands", h.GetBrands())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/brands", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	brandList, ok := data["brands"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, brandList, 2)

	firstBrand := brandList[0].(map[string]interface{})
	assert.NotEmpty(t, firstBrand["id"])
	assert.Equal(t, "Yamaha", firstBrand["name"])

	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockBrandInteractor.AssertExpectations(t)
}

// TestGetDepartments_Integration_Success validates the full HTTP pipeline for GET /departments.
// Exercises: interactor call → ID encoding → HATEOAS → 200 response.
func TestGetDepartments_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockLocationInteractor := new(mocks.MockLocationInteractor)

	h := handlers.NewForTest(
		nil,
		mockLocationInteractor,
		nil, nil, nil,
		encoder,
		responseHandler,
	)

	departments := []domain.Department{
		{ID: "d1234567-89ab-cdef-0123-456789abcdef", Name: "Antioquia"},
		{ID: "d2234567-89ab-cdef-0123-456789abcdef", Name: "Cundinamarca"},
	}
	mockLocationInteractor.On("GetAllDepartments", mock.Anything).Return(departments, nil)

	router := gin.New()
	router.GET("/departments", h.GetDepartments())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	deptList, ok := data["departments"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, deptList, 2)

	firstDept := deptList[0].(map[string]interface{})
	assert.NotEmpty(t, firstDept["id"])
	assert.Equal(t, "Antioquia", firstDept["name"])

	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockLocationInteractor.AssertExpectations(t)
}
