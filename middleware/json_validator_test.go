package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	json_schema "github.com/EstebanGitPro/motogo-backend/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewMiddlewareValidator Tests
// ============================================

func TestNewMiddlewareValidator_ReturnsBuilder(t *testing.T) {
	validators := &json_schema.Validators{}

	builder := NewMiddlewareValidator(validators)

	assert.NotNil(t, builder)
	assert.Equal(t, validators, builder.Validators)
	assert.False(t, builder.isLogin)
}

func TestNewMiddlewareValidator_NilValidators(t *testing.T) {
	builder := NewMiddlewareValidator(nil)

	assert.NotNil(t, builder)
	assert.Nil(t, builder.Validators)
}

// ============================================
// translateFieldNames Tests
// ============================================

func TestTranslateFieldNames_KnownFields(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"email"},
			expected: []string{"Correo electrónico"},
		},
		{
			input:    []string{"password"},
			expected: []string{"Contraseña"},
		},
		{
			input:    []string{"first_name", "last_name"},
			expected: []string{"Nombre", "Apellido"},
		},
		{
			input:    []string{"license_plate", "year", "current_mileage"},
			expected: []string{"Placa", "Año del modelo", "Kilometraje actual"},
		},
		{
			input:    []string{"day_of_week", "opening_time", "closing_time"},
			expected: []string{"Día de la semana", "Hora de apertura", "Hora de cierre"},
		},
		{
			input:    []string{"department_id", "city_id"},
			expected: []string{"Departamento", "Ciudad"},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := translateFieldNames(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslateFieldNames_UnknownFields(t *testing.T) {
	input := []string{"unknown_field", "another_unknown"}
	result := translateFieldNames(input)

	// Unknown fields should be kept as-is
	assert.Equal(t, input, result)
}

func TestTranslateFieldNames_MixedFields(t *testing.T) {
	input := []string{"email", "unknown_field", "phone_number"}
	expected := []string{"Correo electrónico", "unknown_field", "Número de teléfono"}

	result := translateFieldNames(input)

	assert.Equal(t, expected, result)
}

func TestTranslateFieldNames_EmptySlice(t *testing.T) {
	result := translateFieldNames([]string{})

	assert.Empty(t, result)
}

func TestTranslateFieldNames_AllScheduleFields(t *testing.T) {
	input := []string{"exception_date", "is_closed", "start_date", "end_date", "active"}
	expected := []string{"Fecha de excepción", "Cerrado", "Fecha de inicio", "Fecha de fin", "Activo"}

	result := translateFieldNames(input)

	assert.Equal(t, expected, result)
}

func TestTranslateFieldNames_AllBranchFields(t *testing.T) {
	input := []string{"name", "address", "branch_type", "latitude", "longitude"}
	expected := []string{"Nombre", "Dirección", "Tipo de establecimiento", "Latitud", "Longitud"}

	result := translateFieldNames(input)

	assert.Equal(t, expected, result)
}

func TestTranslateFieldNames_PasswordFields(t *testing.T) {
	input := []string{"current_password", "new_password", "confirm_password"}
	expected := []string{"Contraseña actual", "Nueva contraseña", "Confirmar contraseña"}

	result := translateFieldNames(input)

	assert.Equal(t, expected, result)
}

// ============================================
// fieldNameMapping Tests
// ============================================

func TestFieldNameMapping_HasExpectedEntries(t *testing.T) {
	// Verify key mappings exist
	expectedMappings := map[string]string{
		"email":           "Correo electrónico",
		"password":        "Contraseña",
		"license_plate":   "Placa",
		"day_of_week":     "Día de la semana",
		"opening_time":    "Hora de apertura",
		"exception_date":  "Fecha de excepción",
		"department_id":   "Departamento",
		"city_id":         "Ciudad",
		"reference_id":    "Referencia de motocicleta",
		"current_mileage": "Kilometraje actual",
	}

	for key, expectedValue := range expectedMappings {
		t.Run(key, func(t *testing.T) {
			actualValue, exists := fieldNameMapping[key]
			assert.True(t, exists, "Expected key %s to exist in fieldNameMapping", key)
			assert.Equal(t, expectedValue, actualValue)
		})
	}
}

// ============================================
// WithValidate* Methods Tests
// These tests verify each validator method returns a valid handler
// ============================================

func TestWithValidateRegister_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateRegister()
	assert.NotNil(t, handler)
	assert.False(t, builder.isLogin)
}

func TestWithValidateMessage_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateMessage()
	assert.NotNil(t, handler)
}

func TestWithValidateCreateMessage_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateCreateMessage()
	assert.NotNil(t, handler)
}

func TestWithValidateResendVerification_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateResendVerification()
	assert.NotNil(t, handler)
}

func TestWithValidatePasswordReset_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidatePasswordReset()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateProfile_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateProfile()
	assert.NotNil(t, handler)
}

func TestWithValidateResetPasswordWithToken_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateResetPasswordWithToken()
	assert.NotNil(t, handler)
}

func TestWithValidateChangePassword_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateChangePassword()
	assert.NotNil(t, handler)
}

func TestWithValidateRegisterBranch_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateRegisterBranch()
	assert.NotNil(t, handler)
}

func TestWithValidateScheduleDetail_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateScheduleDetail()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateSchedule_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateSchedule()
	assert.NotNil(t, handler)
}

func TestWithValidateScheduleException_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateScheduleException()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateScheduleException_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateScheduleException()
	assert.NotNil(t, handler)
}

func TestWithValidateFranchise_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateFranchise()
	assert.NotNil(t, handler)
}

func TestWithValidateRegisterMotorcycle_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateRegisterMotorcycle()
	assert.NotNil(t, handler)
}

func TestWithValidateEvidence_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateEvidence()
	assert.NotNil(t, handler)
}

func TestWithValidateCompletedService_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateCompletedService()
	assert.NotNil(t, handler)
}

func TestWithValidateDiagnostic_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateDiagnostic()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateDiagnostic_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateDiagnostic()
	assert.NotNil(t, handler)
}

func TestWithValidateDiagnosticSolution_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateDiagnosticSolution()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateScheduleDetail_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateScheduleDetail()
	assert.NotNil(t, handler)
}

func TestWithValidateDiagnosticPermission_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateDiagnosticPermission()
	assert.NotNil(t, handler)
}

func TestWithValidateBranchServices_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateBranchServices()
	assert.NotNil(t, handler)
}

func TestWithValidateFranchiseBranch_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateFranchiseBranch()
	assert.NotNil(t, handler)
}

func TestWithValidateUpdateStatus_ReturnsHandler(t *testing.T) {
	builder := NewMiddlewareValidator(&json_schema.Validators{})
	handler := builder.WithValidateUpdateStatus()
	assert.NotNil(t, handler)
}

func TestTranslateFieldNames_StatusField(t *testing.T) {
	fields := []string{"status"}
	expected := []string{"Estado del servicio"}

	result := translateFieldNames(fields)
	assert.Equal(t, expected, result)
}

func TestTranslateFieldNames_CompletedServiceFields(t *testing.T) {
	fields := []string{"motorcycle_id", "service_ids", "diagnostic_id", "quoted_price", "final_price", "representative_notes"}
	expected := []string{"Motocicleta", "Servicios", "Diagnóstico", "Precio cotizado", "Precio final", "Notas del representante"}

	result := translateFieldNames(fields)
	assert.Equal(t, expected, result)
}

func TestTranslateFieldNames_DiagnosticFields(t *testing.T) {
	fields := []string{"branch_id", "problem_description", "possible_solution", "evidence_urls"}
	expected := []string{"Sede", "Descripción del problema", "Posible solución", "URLs de evidencia"}

	result := translateFieldNames(fields)
	assert.Equal(t, expected, result)
}

// ============================================
// extractFieldNames Tests
// ============================================

func TestExtractFieldNames_WithPropertiesParam(t *testing.T) {
	errors := map[string]*jsonschema.EvaluationError{
		"required": {
			Keyword: "required",
			Code:    "required",
			Message: "missing properties",
			Params: map[string]any{
				"properties": "email, password",
			},
		},
	}

	result := extractFieldNames(errors)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "email")
	assert.Contains(t, result, "password")
}

func TestExtractFieldNames_WithPropertyParam(t *testing.T) {
	errors := map[string]*jsonschema.EvaluationError{
		"type": {
			Keyword: "type",
			Code:    "type",
			Message: "invalid type",
			Params: map[string]any{
				"property": "phone_number",
			},
		},
	}

	result := extractFieldNames(errors)

	assert.Len(t, result, 1)
	assert.Equal(t, "phone_number", result[0])
}

func TestExtractFieldNames_NilParams(t *testing.T) {
	errors := map[string]*jsonschema.EvaluationError{
		"null": {
			Keyword: "required",
			Code:    "required",
			Message: "missing",
			Params:  nil,
		},
	}

	result := extractFieldNames(errors)

	assert.Empty(t, result)
}

func TestExtractFieldNames_EmptyErrors(t *testing.T) {
	errors := map[string]*jsonschema.EvaluationError{}

	result := extractFieldNames(errors)

	assert.Empty(t, result)
}

// ============================================
// classifyValidationError Tests
// ============================================

func TestClassifyValidationError_MultipleFields(t *testing.T) {
	fields := []string{"email", "password", "phone_number"}
	errors := map[string]*jsonschema.EvaluationError{}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrMultipleFields, result)
}

func TestClassifyValidationError_PropertyMismatch(t *testing.T) {
	fields := []string{"email"}
	errors := map[string]*jsonschema.EvaluationError{
		"property": {Code: "property_mismatch"},
	}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrFieldPropertyMismatch, result)
}

func TestClassifyValidationError_Required(t *testing.T) {
	fields := []string{"name"}
	errors := map[string]*jsonschema.EvaluationError{
		"required": {Code: "required"},
	}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrFieldRequired, result)
}

func TestClassifyValidationError_TypeInvalid(t *testing.T) {
	fields := []string{"age"}
	errors := map[string]*jsonschema.EvaluationError{
		"type": {Code: "type"},
	}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrFieldTypeInvalid, result)
}

func TestClassifyValidationError_DefaultCase(t *testing.T) {
	fields := []string{"field"}
	errors := map[string]*jsonschema.EvaluationError{
		"unknown": {Code: "unknown_code"},
	}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrValidationFailed, result)
}

func TestClassifyValidationError_NilFirstError(t *testing.T) {
	fields := []string{"field"}
	errors := map[string]*jsonschema.EvaluationError{}

	result := classifyValidationError(fields, errors)

	assert.Equal(t, json_schema.ErrValidationFailed, result)
}

// ============================================
// jsonValidator Handler Tests
// ============================================

// helper to compile an inline JSON schema for tests
func compileTestSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile([]byte(`{
		"type": "object",
		"properties": {
			"email": {"type": "string"}
		},
		"required": ["email"],
		"additionalProperties": false
	}`))
	assert.NoError(t, err)
	return schema
}

func TestJsonValidator_ValidJSON_PassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"email":"test@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.False(t, c.IsAborted())
	assert.Empty(t, c.Errors)
}

func TestJsonValidator_InvalidJSON_Aborts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{invalid json`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.True(t, c.IsAborted())
	assert.Len(t, c.Errors, 1)
	assert.Equal(t, json_schema.ErrBadRequest, c.Errors[0].Err)
}

func TestJsonValidator_MissingRequiredField_Aborts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.True(t, c.IsAborted())
	assert.NotEmpty(t, c.Errors)
}

func TestJsonValidator_WrongFieldType_Aborts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// email should be string, not number
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"email": 12345}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.True(t, c.IsAborted())
	assert.NotEmpty(t, c.Errors)
}

func TestJsonValidator_AdditionalProperties_Aborts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"email":"a@b.com","extra":"field"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.True(t, c.IsAborted())
	assert.NotEmpty(t, c.Errors)
}

func TestJsonValidator_EmptyBody_Aborts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schema := compileTestSchema(t)
	builder := &Builder{}
	handler := builder.jsonValidator(schema)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(``))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.True(t, c.IsAborted())
	assert.NotEmpty(t, c.Errors)
}
