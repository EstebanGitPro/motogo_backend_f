package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
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

// createTestMessageCache creates a MessageCache with success messages pre-loaded
// so ResponseHandler can resolve the correct HTTP status codes.
func createTestMessageCache() *messagingCache.MessageCache {
	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)

	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_REG_EXI_00001", Type: "EXITO", Title: "Sede registrada", Content: "La sede ha sido registrada exitosamente", Active: true},
		{Code: "MOD_B_BRD_EXI_00001", Type: "EXITO", Title: "Marcas obtenidas", Content: "Catálogo de marcas obtenido", Active: true},
		{Code: "MOD_L_DEP_EXI_00001", Type: "EXITO", Title: "Departamentos", Content: "Departamentos obtenidos", Active: true},
		{Code: "MOD_MOT_CREATE_EXI_00001", Type: "EXITO", Title: "Motocicleta registrada", Content: "Motocicleta registrada exitosamente", Active: true},
		{Code: "MOD_EVD_CREATE_EXI_00001", Type: "EXITO", Title: "Evidencia creada", Content: "Evidencia fotográfica creada exitosamente", Active: true},
		{Code: "MOD_MOT_IMG_UPDATE_EXI_00001", Type: "EXITO", Title: "Imagen actualizada", Content: "Imagen de perfil actualizada exitosamente", Active: true},
		{Code: "MOD_B_GET_EXI_00001", Type: "EXITO", Title: "Sede encontrada", Content: "Sede encontrada exitosamente", Active: true},
		{Code: "MOD_S_LIST_EXI_00001", Type: "EXITO", Title: "Servicios obtenidos", Content: "Lista de servicios obtenida exitosamente", Active: true},
		{Code: "MOD_DGN_CREATE_EXI_00001", Type: "EXITO", Title: "Diagnóstico creado", Content: "Diagnóstico creado exitosamente", Active: true},
		{Code: "MOD_F_REG_EXI_00001", Type: "EXITO", Title: "Franquicia registrada", Content: "Franquicia registrada exitosamente", Active: true},
		{Code: "MOD_H_CREATE_EXI_00001", Type: "EXITO", Title: "Horario registrado", Content: "Horario registrado exitosamente", Active: true},
		{Code: "MOD_HD_CREATE_EXI_00001", Type: "EXITO", Title: "Detalle horario registrado", Content: "Detalle de horario registrado exitosamente", Active: true},
		{Code: "MOD_EH_CREATE_EXI_00001", Type: "EXITO", Title: "Excepción creada", Content: "Excepción de horario creada exitosamente", Active: true},
		{Code: "MOD_AUTH_PROFILE_EXI_00001", Type: "EXITO", Title: "Perfil obtenido", Content: "Perfil del usuario obtenido exitosamente", Active: true},
		{Code: "GEN_OPE_EXI_00001", Type: "EXITO", Title: "Operación exitosa", Content: "Operación completada exitosamente", Active: true},
		{Code: "MOD_M_CREATE_EXI_00001", Type: "EXITO", Title: "Mensaje creado", Content: "Mensaje del sistema creado exitosamente", Active: true},
		{Code: "MOD_U_REG_EXI_00001", Type: "EXITO", Title: "Usuario registrado", Content: "Usuario registrado exitosamente", Active: true},
		{Code: "MOD_AUTH_LOGIN_SUCCESS_EXI_00001", Type: "EXITO", Title: "Login exitoso", Content: "Autenticación exitosa", Active: true},
		{Code: "MOD_P_UPD_EXI_00002", Type: "EXITO", Title: "Perfil actualizado", Content: "Perfil actualizado exitosamente", Active: true},
		{Code: "MOD_P_CHANGE_EXI_00001", Type: "EXITO", Title: "Contraseña cambiada", Content: "Contraseña actualizada exitosamente", Active: true},
		{Code: "MOD_P_DEL_EXI_00001", Type: "EXITO", Title: "Cuenta eliminada", Content: "Cuenta eliminada exitosamente", Active: true},
		{Code: "MOD_P_CONTACT_EXI_00001", Type: "EXITO", Title: "Contacto obtenido", Content: "Información de contacto obtenida", Active: true},
		{Code: "MOD_MOT_GET_EXI_00001", Type: "EXITO", Title: "Moto obtenida", Content: "Motocicleta obtenida exitosamente", Active: true},
		{Code: "MOD_MOT_LIST_EXI_00001", Type: "EXITO", Title: "Motos listadas", Content: "Motocicletas listadas exitosamente", Active: true},
		{Code: "MOD_MOT_UPDATE_EXI_00001", Type: "EXITO", Title: "Moto actualizada", Content: "Motocicleta actualizada exitosamente", Active: true},
		{Code: "MOD_MOT_DELETE_EXI_00001", Type: "EXITO", Title: "Moto eliminada", Content: "Motocicleta eliminada exitosamente", Active: true},
		{Code: "MOD_MOT_REF_LIST_EXI_00001", Type: "EXITO", Title: "Referencias listadas", Content: "Referencias de motocicleta listadas", Active: true},
		{Code: "MOD_MOT_BRAND_LINES_EXI_00001", Type: "EXITO", Title: "Líneas de marca", Content: "Líneas de marca obtenidas", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())

	return cache
}

// TestRegisterBranch_Integration_Success validates the full HTTP handler pipeline
// for the success path of RegisterBranch (HU59).
//
// It exercises: JSON binding → sanitization → ID decoding → domain mapping →
// interactor delegation → response serialization → HATEOAS links → 201 response.
func TestRegisterBranch_Integration_Success(t *testing.T) {
	// ============================================
	// 1. Setup dependencies
	// ============================================
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()

	// MessageCache with branch success message pre-loaded
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Real BranchInteractor with mocked service
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	// Create handler (only branchInteractor and encoder/response are needed)
	h := handlers.New(
		nil,              // personInteractor
		nil,              // messageInteractor
		branchInteractor, // branchInteractor ← real with mocked service
		nil,              // brandInteractor
		nil,              // locationInteractor
		nil,              // serviceInteractor
		nil,              // franchiseInteractor
		nil,              // motorcycleInteractor
		nil,              // evidenceInteractor
		nil,              // diagnosticInteractor
		nil,              // firebaseClient
		nil,              // messageCache
		encoder,          // IDEncoder
		responseHandler,  // ResponseHandler
	)

	// ============================================
	// 2. Prepare test data
	// ============================================
	// Raw UUIDs (what the interactor/service layer sees)
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	brandUUID1 := "b1234567-89ab-cdef-0123-456789abcdef"
	brandUUID2 := "c1234567-89ab-cdef-0123-456789abcdef"
	departmentUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	cityUUID := "e1234567-89ab-cdef-0123-456789abcdef"
	representativeID := "rep-123"

	// Encode IDs for the request body (frontend sends encoded IDs)
	encodedBrand1, _ := encoder.Encode(brandUUID1)
	encodedBrand2, _ := encoder.Encode(brandUUID2)
	encodedDeptID, _ := encoder.Encode(departmentUUID)
	encodedCityID, _ := encoder.Encode(cityUUID)

	// Geocoded coordinates (set by the interactor's geocoding step)
	lat := 4.710989
	lng := -74.072092

	// Domain branch returned by the service after save
	savedBranch := &domain.Branch{
		ID:                branchUUID,
		Name:              "Taller Norte",
		RepresentativeID:  representativeID,
		EstablishmentType: domain.EstablishmentTypeWorkshop,
		Status:            "ACTIVE",
		Brands:            []string{brandUUID1, brandUUID2},
		Location: &domain.Location{
			DepartmentID: departmentUUID,
			CityID:       cityUUID,
			Address:      "Calle 100 #15-20",
			Latitude:     &lat,
			Longitude:    &lng,
		},
	}

	// ============================================
	// 3. Configure mock expectations
	// ============================================
	// The interactor calls these in order:
	// 1) ValidateBrands
	mockBranchService.On("ValidateBrands", mock.Anything, []string{brandUUID1, brandUUID2}).Return(nil)
	// 2) GeocodeLocation (sets coordinates on the location)
	mockBranchService.On("GeocodeLocation", mock.Anything, mock.AnythingOfType("*domain.Location")).
		Run(func(args mock.Arguments) {
			loc := args.Get(1).(*domain.Location)
			loc.Latitude = &lat
			loc.Longitude = &lng
		}).Return(true, nil)
	// 3) BeginTx
	mockBranchService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	// 4) RegisterBranch (inside tx)
	mockBranchService.On("RegisterBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(savedBranch, nil)
	// 5) Commit
	mockTx.On("Commit").Return(nil)
	// Safety: Rollback should not be called on success, but register it to avoid panics
	mockTx.On("Rollback").Return(nil)

	// ============================================
	// 4. Build HTTP request
	// ============================================
	reqBody := map[string]interface{}{
		"name":               "  Taller Norte  ", // Extra spaces → sanitized
		"establishment_type": "WORKSHOP",
		"brands":             []string{encodedBrand1, encodedBrand2},
		"location": map[string]interface{}{
			"department_id": encodedDeptID,
			"city_id":       encodedCityID,
			"address":       "Calle 100 #15-20",
		},
	}
	bodyJSON, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	// ============================================
	// 5. Setup router with auth middleware
	// ============================================
	router := gin.New()
	// Simulate authentication middleware (sets authenticated user in context)
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   representativeID,
			Role: "REPRESENTANTE",
		})
		c.Next()
	})
	router.POST("/branches", h.RegisterBranch())

	// ============================================
	// 6. Execute request
	// ============================================
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// ============================================
	// 7. Assert response
	// ============================================
	// MessageCache resolves MOD_B_REG_EXI_00001 → 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify top-level API response structure
	assert.True(t, response["success"].(bool))

	// Verify data payload
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "response should contain 'data' object")

	// Branch name (sanitized from "  Taller Norte  ")
	assert.Equal(t, "Taller Norte", data["name"])
	assert.Equal(t, "WORKSHOP", data["establishment_type"])
	assert.Equal(t, "ACTIVE", data["status"])

	// Encoded ID should be non-empty
	assert.NotEmpty(t, data["id"])

	// Geocoding status (no user coords + geocoding succeeded = SUCCESS)
	assert.Equal(t, "SUCCESS", data["geocoding_status"])

	// HATEOAS links should be present
	links, ok := data["_links"].([]interface{})
	assert.True(t, ok, "response should contain '_links' array")
	assert.NotEmpty(t, links, "_links should not be empty")

	// Location should be present with encoded IDs
	location, ok := data["location"].(map[string]interface{})
	assert.True(t, ok, "response should contain 'location' object")
	assert.Equal(t, "Calle 100 #15-20", location["address"])
	assert.NotEmpty(t, location["department_id"])
	assert.NotEmpty(t, location["city_id"])

	// Brands should be present (encoded)
	brands, ok := data["brands"].([]interface{})
	assert.True(t, ok, "response should contain 'brands' array")
	assert.Len(t, brands, 2)

	// Location header should be set
	assert.NotEmpty(t, w.Header().Get("Location"))

	// ============================================
	// 8. Verify mock expectations
	// ============================================
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}
