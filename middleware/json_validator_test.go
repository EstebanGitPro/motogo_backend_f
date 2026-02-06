package middleware

import (
	"testing"

	json_schema "github.com/EstebanGitPro/motogo-backend/platform/schema"
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
	input := []string{"email", "unknown_field", "phone"}
	expected := []string{"Correo electrónico", "unknown_field", "Teléfono"}

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
