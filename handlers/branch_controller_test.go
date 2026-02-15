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
		{Code: "MOD_B_UPD_EXI_00001", Type: "EXITO", Title: "Sede actualizada", Content: "La sede ha sido actualizada exitosamente", Active: true},
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
		// Schedule success messages
		{Code: "MOD_H_GET_EXI_00001", Type: "EXITO", Title: "Horario obtenido", Content: "Horario obtenido exitosamente", Active: true},
		{Code: "MOD_H_UPDATE_EXI_00001", Type: "EXITO", Title: "Horario actualizado", Content: "Horario actualizado exitosamente", Active: true},
		{Code: "MOD_H_DELETE_EXI_00001", Type: "EXITO", Title: "Horario eliminado", Content: "Horario eliminado exitosamente", Active: true},
		{Code: "MOD_H_ACTIV_EXI_00001", Type: "EXITO", Title: "Horario activado", Content: "Horario activado exitosamente", Active: true},
		{Code: "MOD_H_DEACT_EXI_00001", Type: "EXITO", Title: "Horario desactivado", Content: "Horario desactivado exitosamente", Active: true},
		{Code: "MOD_H_DAYS_EXI_00001", Type: "EXITO", Title: "Días obtenidos", Content: "Catálogo de días obtenido", Active: true},
		// Error messages
		{Code: "MOD_B_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Sede no encontrada", Content: "La sede no fue encontrada", Active: true},
		{Code: "GEN_SRV_ERR_00001", Type: "ERROR", Title: "Error del servidor", Content: "Error interno del servidor", Active: true},
		{Code: "GEN_FORBIDDEN_ERR_00003", Type: "ERROR", Title: "Acción no permitida", Content: "No tiene permiso para realizar esta acción", Active: true},
		{Code: "MOD_H_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Horario no encontrado", Content: "El horario no fue encontrado", Active: true},
		{Code: "MOD_H_EXISTS_ERR_00001", Type: "ERROR", Title: "Horario ya existe", Content: "Ya existe un horario para esta sede", Active: true},
		{Code: "MOD_H_DATE_FORMAT_ERR_00001", Type: "ERROR", Title: "Formato inválido", Content: "Formato de fecha inválido", Active: true},
		{Code: "MOD_H_DATE_RANGE_ERR_00001", Type: "ERROR", Title: "Rango inválido", Content: "Rango de fechas inválido", Active: true},
		// Diagnostic permission messages
		{Code: "MOD_DGN_PERM_GRANT_EXI_00001", Type: "EXITO", Title: "Permiso otorgado", Content: "Permiso de diagnóstico otorgado", Active: true},
		{Code: "MOD_DGN_PERM_REVOKE_EXI_00001", Type: "EXITO", Title: "Permiso revocado", Content: "Permiso de diagnóstico revocado", Active: true},
		{Code: "MOD_DGN_PERM_GET_EXI_00001", Type: "EXITO", Title: "Permisos obtenidos", Content: "Permisos de diagnóstico obtenidos", Active: true},
		{Code: "MOD_DGN_PERM_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Permiso no encontrado", Content: "Permiso no encontrado", Active: true},
		// Schedule exception messages
		{Code: "MOD_EH_CREATE_EXI_00001", Type: "EXITO", Title: "Excepción creada", Content: "Excepción de horario creada exitosamente", Active: true},
		{Code: "MOD_EH_GET_EXI_00001", Type: "EXITO", Title: "Excepción obtenida", Content: "Excepción obtenida exitosamente", Active: true},
		{Code: "MOD_EH_LIST_EXI_00001", Type: "EXITO", Title: "Excepciones listadas", Content: "Excepciones listadas exitosamente", Active: true},
		{Code: "MOD_EH_UPDATE_EXI_00001", Type: "EXITO", Title: "Excepción actualizada", Content: "Excepción actualizada exitosamente", Active: true},
		{Code: "MOD_EH_DELETE_EXI_00001", Type: "EXITO", Title: "Excepción eliminada", Content: "Excepción eliminada exitosamente", Active: true},
		{Code: "MOD_EH_ACTIVATE_EXI_00001", Type: "EXITO", Title: "Excepción activada", Content: "Excepción activada exitosamente", Active: true},
		{Code: "MOD_EH_DEACTIVATE_EXI_00001", Type: "EXITO", Title: "Excepción desactivada", Content: "Excepción desactivada exitosamente", Active: true},
		{Code: "MOD_EH_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Excepción no encontrada", Content: "Excepción no encontrada", Active: true},
		{Code: "MOD_EH_DATE_CONFLICT_ERR_00001", Type: "ERROR", Title: "Conflicto de fecha", Content: "Conflicto de fecha en excepción", Active: true},
		{Code: "MOD_EH_DATE_PAST_ERR_00001", Type: "ERROR", Title: "Fecha pasada", Content: "La fecha de excepción ya pasó", Active: true},
		{Code: "MOD_EH_TIME_ERR_00001", Type: "ERROR", Title: "Hora inválida", Content: "Hora de excepción inválida", Active: true},
		{Code: "MOD_EH_REDUNDANT_ERR_00001", Type: "ERROR", Title: "Excepción redundante", Content: "La excepción es redundante", Active: true},
		// Schedule detail messages
		{Code: "MOD_HD_CREATE_EXI_00001", Type: "EXITO", Title: "Detalle creado", Content: "Detalle de horario creado exitosamente", Active: true},
		{Code: "MOD_HD_GET_EXI_00001", Type: "EXITO", Title: "Detalle obtenido", Content: "Detalle de horario obtenido", Active: true},
		{Code: "MOD_HD_UPDATE_EXI_00001", Type: "EXITO", Title: "Detalle actualizado", Content: "Detalle actualizado exitosamente", Active: true},
		{Code: "MOD_HD_DELETE_EXI_00001", Type: "EXITO", Title: "Detalle eliminado", Content: "Detalle eliminado exitosamente", Active: true},
		{Code: "MOD_HD_LIST_EXI_00001", Type: "EXITO", Title: "Detalles listados", Content: "Detalles de horario listados", Active: true},
		{Code: "MOD_HD_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Detalle no encontrado", Content: "Detalle de horario no encontrado", Active: true},
		{Code: "MOD_HD_CONFLICT_ERR_00001", Type: "ERROR", Title: "Conflicto de horario", Content: "Conflicto de horario detectado", Active: true},
		{Code: "MOD_HD_TIME_ERR_00001", Type: "ERROR", Title: "Hora inválida", Content: "Formato de hora inválido", Active: true},
		{Code: "MOD_HD_DAY_ERR_00001", Type: "ERROR", Title: "Día inválido", Content: "Día de la semana inválido", Active: true},
		{Code: "MOD_HD_DAY_CLOSED_ERR_00001", Type: "ERROR", Title: "Día cerrado", Content: "El día ya está cerrado", Active: true},
		{Code: "MOD_HD_DAY_HAS_SLOTS_ERR_00001", Type: "ERROR", Title: "Día con franja", Content: "El día ya tiene franjas horarias", Active: true},
		// Validation error
		{Code: "MOD_V_VAL_ERR_00001", Type: "ERROR", Title: "Formato inválido", Content: "Formato de solicitud inválido", Active: true},
		{Code: "MOD_V_ID_ERR_00013", Type: "ERROR", Title: "ID inválido", Content: "El ID proporcionado no es válido", Active: true},
		// Location messages
		{Code: "MOD_L_DEP_EXI_00001", Type: "EXITO", Title: "Departamentos obtenidos", Content: "Departamentos obtenidos exitosamente", Active: true},
		{Code: "MOD_L_CIT_EXI_00001", Type: "EXITO", Title: "Ciudades obtenidas", Content: "Ciudades obtenidas exitosamente", Active: true},
		// Motorcycle category messages
		{Code: "MOD_MOT_CAT_LIST_EXI_00001", Type: "EXITO", Title: "Categorías obtenidas", Content: "Categorías de motos obtenidas", Active: true},
		{Code: "MOD_MOT_CAT_LINES_EXI_00001", Type: "EXITO", Title: "Líneas obtenidas", Content: "Líneas de categoría obtenidas", Active: true},
		{Code: "MOD_MOT_DISP_LIST_EXI_00001", Type: "EXITO", Title: "Cilindradas obtenidas", Content: "Rangos de cilindrada obtenidos", Active: true},
		{Code: "MOD_MOT_RATE_LIST_EXI_00001", Type: "EXITO", Title: "Rangos obtenidos", Content: "Rangos de calificación obtenidos", Active: true},
		// Diagnostic permission messages (actual codes used by controller)
		{Code: "MOD_DGP_GRANT_EXI_00001", Type: "EXITO", Title: "Permiso otorgado", Content: "Permiso de diagnóstico otorgado", Active: true},
		{Code: "MOD_DGP_REVOKE_EXI_00001", Type: "EXITO", Title: "Permiso revocado", Content: "Permiso de diagnóstico revocado", Active: true},
		{Code: "MOD_DGP_LIST_EXI_00001", Type: "EXITO", Title: "Permisos listados", Content: "Permisos de diagnóstico listados", Active: true},
		{Code: "MOD_DGP_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Permiso no encontrado", Content: "Permiso de diagnóstico no encontrado", Active: true},
		{Code: "MOD_DGP_SAVE_ERR_00001", Type: "ERROR", Title: "No se pudo guardar", Content: "No se pudo guardar el permiso", Active: true},
		{Code: "MOD_DGP_DELETE_ERR_00001", Type: "ERROR", Title: "No se pudo eliminar", Content: "No se pudo eliminar el permiso", Active: true},
		{Code: "MOD_MOT_NOT_FOUND_ERR_00001", Type: "ERROR", Title: "Moto no encontrada", Content: "La motocicleta no fue encontrada", Active: true},
		// Motorcycle lookup error messages
		{Code: "MOD_MOT_PLATE_REQ_ERR_00001", Type: "ERROR", Title: "Placa requerida", Content: "El parámetro de placa es requerido", Active: true},
		{Code: "MOD_MOT_NO_PERMISSION_ERR_00001", Type: "ERROR", Title: "Sin permiso", Content: "Sin permiso para esta motocicleta", Active: true},
		{Code: "GEN_SRV_ERR_00001", Type: "ERROR", Title: "Error de servidor", Content: "Error interno del servidor", Active: true},
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

// TestGetNearbyBranches_Integration_Success validates GET /branches/nearby
// Exercises: parseNearbyFilters → parseOptionalBrandFilter → interactor call → brand encoding → response building
func TestGetNearbyBranches_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	// Add nearby success message to cache
	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_NEARBY_EXI_00001", Type: "EXITO", Title: "Sedes cercanas", Content: "Sedes cercanas encontradas", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	// Domain data returned by interactor
	branchNearbyUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	brandUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	encodedBrand, _ := encoder.Encode(brandUUID)
	lat, lng := 4.710989, -74.072092

	nearbyBranches := []domain.NearbyBranch{
		{
			ID:                branchNearbyUUID,
			Name:              "Taller Norte",
			EstablishmentType: domain.EstablishmentTypeWorkshop,
			DistanceKm:        2.5,
			Brands:            []string{brandUUID},
			Location: &domain.NearbyLocation{
				Address:        "Calle 100 #15-20",
				CityName:       "Bogotá",
				DepartmentName: "Cundinamarca",
				Latitude:       lat,
				Longitude:      lng,
			},
		},
	}

	// Mock expectations: interactor.GetBranchesNearby delegates to service
	mockBranchService.On("GetBranchesNearby",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		"WORKSHOP",
		mock.AnythingOfType("string"),
		"BAJO",
	).Return(nearbyBranches, nil)

	// Build request with all query params including brand filter
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "user-1", Role: "USER"})
		c.Next()
	})
	router.GET("/branches/nearby", h.GetNearbyBranches())

	w := httptest.NewRecorder()
	url := "/branches/nearby?lat=4.710989&lng=-74.072092&radius=10&type=WORKSHOP&brand=" + encodedBrand + "&displacement_range=BAJO"
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), data["count"])

	branches := data["branches"].([]interface{})
	assert.Len(t, branches, 1)
	branch := branches[0].(map[string]interface{})
	assert.Equal(t, "Taller Norte", branch["name"])
	assert.NotEmpty(t, branch["id"])
}

// TestListBranches_Integration_Success validates GET /branches (my branches)
// Exercises: buildBranchListItem helper with franchise, brands, location encoding
func TestListBranches_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_LIST_EXI_00001", Type: "EXITO", Title: "Sedes listadas", Content: "Lista de sedes obtenida", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	repID := "rep-123"
	franchiseUUID := "f1234567-89ab-cdef-0123-456789abcdef"
	brandUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	lat, lng := 4.71, -74.07

	branches := []domain.Branch{
		{
			ID:                "a1234567-89ab-cdef-0123-456789abcdef",
			Name:              "Taller Norte",
			RepresentativeID:  repID,
			EstablishmentType: domain.EstablishmentTypeWorkshop,
			Status:            "ACTIVE",
			FranchiseID:       &franchiseUUID,
			Brands:            []string{brandUUID},
			Location: &domain.Location{
				DepartmentID: "d1234567-89ab-cdef-0123-456789abcdef",
				CityID:       "e1234567-89ab-cdef-0123-456789abcdef",
				Address:      "Calle 100",
				Latitude:     &lat,
				Longitude:    &lng,
			},
		},
	}

	mockBranchService.On("GetBranchesByRepresentative", mock.Anything, repID).Return(branches, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: "REPRESENTANTE"})
		c.Next()
	})
	router.GET("/branches", h.ListBranches())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["count"])
}

// TestUpdateBranch_Integration_Success validates PUT /branches/:id
// Exercises: decodeBranchRequestIDs, encodeBranchResponseIDs, geocoding status
func TestUpdateBranch_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	brandUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	departmentUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	cityUUID := "e1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	lat, lng := 4.71, -74.07

	encodedBranchID, _ := encoder.Encode(branchUUID)
	encodedBrand, _ := encoder.Encode(brandUUID)
	encodedDept, _ := encoder.Encode(departmentUUID)
	encodedCity, _ := encoder.Encode(cityUUID)

	updatedBranch := &domain.Branch{
		ID:                branchUUID,
		Name:              "Taller Actualizado",
		RepresentativeID:  repID,
		EstablishmentType: domain.EstablishmentTypeWorkshopStore,
		Status:            "ACTIVE",
		Brands:            []string{brandUUID},
		Location: &domain.Location{
			DepartmentID: departmentUUID,
			CityID:       cityUUID,
			Address:      "Calle 200",
			Latitude:     &lat,
			Longitude:    &lng,
		},
	}

	// Mock: validation, geocoding, tx, update, re-fetch
	mockBranchService.On("ValidateBrands", mock.Anything, mock.AnythingOfType("[]string")).Return(nil)
	mockBranchService.On("ValidateDisplacementRanges", mock.AnythingOfType("[]string")).Return(nil)
	mockBranchService.On("GeocodeLocation", mock.Anything, mock.AnythingOfType("*domain.Location")).
		Run(func(args mock.Arguments) {
			loc := args.Get(1).(*domain.Location)
			loc.Latitude = &lat
			loc.Longitude = &lng
		}).Return(true, nil)
	mockBranchService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockBranchService.On("GetBranchByID", mock.Anything, branchUUID).Return(updatedBranch, nil)
	mockBranchService.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	reqBody := map[string]interface{}{
		"name":               "Taller Actualizado",
		"establishment_type": "WORKSHOP_STORE",
		"brands":             []string{encodedBrand},
		"location": map[string]interface{}{
			"department_id": encodedDept,
			"city_id":       encodedCity,
			"address":       "Calle 200",
		},
	}
	bodyJSON, _ := json.Marshal(reqBody)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: "REPRESENTANTE"})
		c.Next()
	})
	router.PUT("/branches/:id", h.UpdateBranch())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Taller Actualizado", data["name"])
	assert.NotEmpty(t, data["id"])
}

// TestDeleteBranch_Integration_Success validates DELETE /branches/:id
func TestDeleteBranch_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_DEL_EXI_00001", Type: "EXITO", Title: "Sede eliminada", Content: "Sede eliminada exitosamente", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	existingBranch := &domain.Branch{
		ID:               branchUUID,
		Name:             "Taller Eliminar",
		RepresentativeID: repID,
		Status:           "ACTIVE",
	}

	mockBranchService.On("GetBranchByID", mock.Anything, branchUUID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", mock.Anything, mockTx, branchUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: "REPRESENTANTE"})
		c.Next()
	})
	router.DELETE("/branches/:id", h.DeleteBranch())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

// TestGetBranchTypes_Integration_Success validates GET /branch-types
func TestGetBranchTypes_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_TYPES_EXI_00001", Type: "EXITO", Title: "Tipos obtenidos", Content: "Tipos de sede obtenidos", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/branch-types", h.GetBranchTypes())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branch-types", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	types := data["types"].([]interface{})
	assert.Len(t, types, 3) // WORKSHOP, STORE, WORKSHOP_STORE
}

// ============================================
// parseNearbyFilters Error Path Tests
// ============================================

// setupNearbyHandler creates the handler and router for nearby branch error tests.
func setupNearbyHandler(t *testing.T) (*gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_V_GEO_LAT_REQ_ERR_00001", Type: "ERROR", Title: "Latitud requerida", Content: "La latitud es requerida", Active: true},
		{Code: "MOD_V_GEO_LAT_INV_ERR_00001", Type: "ERROR", Title: "Latitud inválida", Content: "El valor de latitud es inválido", Active: true},
		{Code: "MOD_V_GEO_LNG_REQ_ERR_00001", Type: "ERROR", Title: "Longitud requerida", Content: "La longitud es requerida", Active: true},
		{Code: "MOD_V_GEO_LNG_INV_ERR_00001", Type: "ERROR", Title: "Longitud inválida", Content: "El valor de longitud es inválido", Active: true},
		{Code: "MOD_V_GEO_RAD_INV_ERR_00001", Type: "ERROR", Title: "Radio inválido", Content: "El radio de búsqueda es inválido", Active: true},
		{Code: "MOD_B_INVALID_TYPE_ERR_00001", Type: "ERROR", Title: "Tipo inválido", Content: "El tipo de establecimiento es inválido", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "user-1", Role: "USER"})
		c.Next()
	})
	router.GET("/branches/nearby", h.GetNearbyBranches())

	return router, httptest.NewRecorder()
}

func TestGetNearbyBranches_MissingLat(t *testing.T) {
	router, w := setupNearbyHandler(t)
	req, _ := http.NewRequest("GET", "/branches/nearby?lng=-74.07", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestGetNearbyBranches_InvalidLat(t *testing.T) {
	router, w := setupNearbyHandler(t)
	req, _ := http.NewRequest("GET", "/branches/nearby?lat=999&lng=-74.07", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestGetNearbyBranches_MissingLng(t *testing.T) {
	router, w := setupNearbyHandler(t)
	req, _ := http.NewRequest("GET", "/branches/nearby?lat=4.71", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestGetNearbyBranches_InvalidRadius(t *testing.T) {
	router, w := setupNearbyHandler(t)
	req, _ := http.NewRequest("GET", "/branches/nearby?lat=4.71&lng=-74.07&radius=-5", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestGetNearbyBranches_InvalidType(t *testing.T) {
	router, w := setupNearbyHandler(t)
	req, _ := http.NewRequest("GET", "/branches/nearby?lat=4.71&lng=-74.07&type=INVALID_TYPE", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

func TestGetNearbyBranches_InvalidBrandFilter(t *testing.T) {
	// When brand filter is invalid, it's silently ignored (empty string passed to interactor)
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_NEARBY_EXI_00001", Type: "EXITO", Title: "Sedes cercanas", Content: "Sedes cercanas encontradas", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	// Expect the interactor to be called with empty brandID because the invalid brand is ignored
	mockBranchService.On("GetBranchesNearby",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"),
		"", // brandID — invalid brand becomes empty string
		mock.AnythingOfType("string"),
	).Return([]domain.NearbyBranch{}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "user-1", Role: "USER"})
		c.Next()
	})
	router.GET("/branches/nearby", h.GetNearbyBranches())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/nearby?lat=4.71&lng=-74.07&brand=INVALID!!!", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeleteBranch Error Path Tests
// ============================================

func TestDeleteBranch_Integration_CannotDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()

	mockRepo := new(mocks.MockMessageCacheRepo)
	cache := messagingCache.NewMessageCache(mockRepo, 0)
	mockRepo.On("GetAllActiveForCache", mock.Anything).Return([]messagingCache.CachedMessage{
		{Code: "MOD_B_HAS_ASSOC_ERR_00001", Type: "ERROR", Title: "Tiene asociaciones", Content: "La sede tiene diagnósticos o servicios asociados", Active: true},
	}, nil)
	_ = cache.LoadMessages(context.TODO())
	responseHandler := middleware.NewResponseHandler(cache)

	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.New(nil, nil, branchInteractor, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	existingBranch := &domain.Branch{
		ID:               branchUUID,
		Name:             "Taller con Diagnosticos",
		RepresentativeID: repID,
		Status:           "ACTIVE",
	}

	mockBranchService.On("GetBranchByID", mock.Anything, branchUUID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", mock.Anything, mockTx, branchUUID).Return(domain.ErrBranchCannotDelete)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: "REPRESENTANTE"})
		c.Next()
	})
	router.DELETE("/branches/:id", h.DeleteBranch())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranchID, nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}
