package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	json_schema "github.com/EstebanGitPro/motogo-backend/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
)

type Builder struct {
	Validators *json_schema.Validators
	isLogin    bool
}

func NewMiddlewareValidator(validators *json_schema.Validators) *Builder {
	return &Builder{
		Validators: validators,
	}
}

// fieldNameMapping maps JSON field names to Spanish labels for user-friendly error messages
var fieldNameMapping = map[string]string{
	// Schedule Exceptions (HU20-25)
	"exception_date": "Fecha de excepción",
	"opening_time":   "Hora de apertura",
	"closing_time":   "Hora de cierre",
	"is_closed":      "Cerrado",

	// Schedule Details (HU6-9)
	"day_of_week": "Día de la semana",
	"entry_type":  "Tipo de entrada",

	// Person/Register
	"identity_number":  "Número de identificación",
	"email":            "Correo electrónico",
	"password":         "Contraseña",
	"first_name":       "Nombre",
	"last_name":        "Apellido",
	"second_last_name": "Segundo apellido",
	"phone_number":     "Número de teléfono",
	"role":             "Rol",
	"current_password": "Contraseña actual",
	"new_password":     "Nueva contraseña",
	"confirm_password": "Confirmar contraseña",
	"token":            "Token",

	// Branch (HU59)
	"name":          "Nombre",
	"address":       "Dirección",
	"branch_type":   "Tipo de establecimiento",
	"latitude":      "Latitud",
	"longitude":     "Longitud",
	"department_id": "Departamento",
	"city_id":       "Ciudad",
	"brands":        "Marcas",

	// Franchise (HU26-29)
	"branches": "Sedes",

	// Schedule (HU30-35)
	"start_date": "Fecha de inicio",
	"end_date":   "Fecha de fin",
	"active":     "Activo",

	// Messages
	"code":     "Código",
	"title":    "Título",
	"content":  "Contenido",
	"module":   "Módulo",
	"category": "Categoría",
	"type":     "Tipo",

	// Motorcycles (HU43-47)
	"license_plate":   "Placa",
	"reference_id":    "Referencia de motocicleta",
	"year":            "Año del modelo",
	"current_mileage": "Kilometraje actual",
	"owner_notes":     "Notas del propietario",

	// Evidence (HU16-19)
	"angle":       "Ángulo de la foto",
	"image_url":   "URL de imagen",
	"description": "Descripción",

	// Completed Services (HU64)
	"motorcycle_id":        "Motocicleta",
	"service_ids":          "Servicios",
	"diagnostic_id":        "Diagnóstico",
	"quoted_price":         "Precio cotizado",
	"final_price":          "Precio final",
	"representative_notes": "Notas del representante",

	// Status Transitions (HU74)
	"status": "Estado del servicio",

	// Diagnostics (HU11-12)
	"branch_id":           "Sede",
	"problem_description": "Descripción del problema",
	"possible_solution":   "Posible solución",
	"evidence_urls":       "URLs de evidencia",
}

// translateFieldNames converts technical field names to Spanish labels
func translateFieldNames(fields []string) []string {
	translated := make([]string, len(fields))
	for i, field := range fields {
		if label, exists := fieldNameMapping[field]; exists {
			translated[i] = label
		} else {
			translated[i] = field // Keep original if no mapping
		}
	}
	return translated
}

func (b *Builder) WithValidateRegister() gin.HandlerFunc {
	b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator)
}

func (b *Builder) WithValidateMessage() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.MessageValidator)
}

func (b *Builder) WithValidateCreateMessage() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateMessageValidator)
}

func (b *Builder) WithValidateResendVerification() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ResendVerificationValidator)
}

func (b *Builder) WithValidatePasswordReset() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.PasswordResetValidator)
}

func (b *Builder) WithValidateUpdateProfile() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateProfileValidator)
}

func (b *Builder) WithValidateResetPasswordWithToken() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ResetPasswordWithTokenValidator)
}

func (b *Builder) WithValidateChangePassword() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ChangePasswordValidator)
}

// WithValidateRegisterBranch validates register branch request (HU59)
func (b *Builder) WithValidateRegisterBranch() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.RegisterBranchValidator)
}

// WithValidateScheduleDetail validates schedule detail (time slot) request (HU6-9)
func (b *Builder) WithValidateScheduleDetail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ScheduleDetailValidator)
}

// WithValidateUpdateSchedule validates update schedule request (HU31)
func (b *Builder) WithValidateUpdateSchedule() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateScheduleValidator)
}

// WithValidateScheduleException validates schedule exception request (HU20-25)
func (b *Builder) WithValidateScheduleException() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ScheduleExceptionValidator)
}

// WithValidateUpdateScheduleException validates update schedule exception request (HU21)
func (b *Builder) WithValidateUpdateScheduleException() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateScheduleExceptionValidator)
}

// WithValidateFranchise validates franchise creation/update request (HU26-29)
func (b *Builder) WithValidateFranchise() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.FranchiseValidator)
}

// WithValidateRegisterMotorcycle validates motorcycle registration request (HU43)
func (b *Builder) WithValidateRegisterMotorcycle() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.RegisterMotorcycleValidator)
}

// WithValidateUpdateMotorcycle validates motorcycle update request (HU44)
func (b *Builder) WithValidateUpdateMotorcycle() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateMotorcycleValidator)
}

// WithValidateEvidence validates evidence creation request (HU16)
func (b *Builder) WithValidateEvidence() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateEvidenceValidator)
}

// WithValidateCompletedService validates completed service registration request (HU64)
func (b *Builder) WithValidateCompletedService() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CompletedServiceValidator)
}

// WithValidateDiagnostic validates diagnostic creation request (HU11)
func (b *Builder) WithValidateDiagnostic() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.DiagnosticValidator)
}

// WithValidateUpdateDiagnostic validates diagnostic update request (HU12)
func (b *Builder) WithValidateUpdateDiagnostic() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateDiagnosticValidator)
}

// WithValidateDiagnosticSolution validates diagnostic solution patch request
func (b *Builder) WithValidateDiagnosticSolution() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.DiagnosticSolutionValidator)
}

// WithValidateUpdateScheduleDetail validates schedule detail update request (HU7)
func (b *Builder) WithValidateUpdateScheduleDetail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateScheduleDetailValidator)
}

// WithValidateDiagnosticPermission validates diagnostic permission grant request
func (b *Builder) WithValidateDiagnosticPermission() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.DiagnosticPermissionValidator)
}

// WithValidateBranchServices validates branch services association request
func (b *Builder) WithValidateBranchServices() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.BranchServicesValidator)
}

// WithValidateFranchiseBranch validates franchise branch association request
func (b *Builder) WithValidateFranchiseBranch() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.FranchiseBranchValidator)
}

// WithValidateUpdateStatus validates status update request (HU74)
func (b *Builder) WithValidateUpdateStatus() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateStatusValidator)
}

func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID for trace correlation
		traceID := GetRequestID(c)
		log := log.WithTraceID(traceID)

		data, ok := readRequestBody(c, log)
		if !ok {
			return
		}

		result := schema.Validate(data)
		if !result.IsValid() {
			handleValidationResult(c, log, result)
			return
		}

		if log != nil {
			log.Debug(logger.LogMiddlewareValidationOK, "path", c.Request.URL.Path)
		}
		c.Next()
	}
}

// readRequestBody reads and parses the JSON request body.
// Returns false and aborts if reading or parsing fails.
func readRequestBody(c *gin.Context, log logger.Logger) (map[string]interface{}, bool) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if log != nil {
			log.Error(logger.LogMiddlewareBodyReadError, "error", err, "path", c.Request.URL.Path)
		}
		_ = c.Error(json_schema.ErrBodyReadFailed)
		c.Abort()
		return nil, false
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		if log != nil {
			log.Error(logger.LogMiddlewareJSONParseError, "error", err, "path", c.Request.URL.Path)
		}
		_ = c.Error(json_schema.ErrBadRequest)
		c.Abort()
		return nil, false
	}

	return data, true
}

// handleValidationResult processes a failed schema validation: extracts field names,
// classifies the error, stores translated fields in context, and aborts.
func handleValidationResult(c *gin.Context, log logger.Logger, result *jsonschema.EvaluationResult) {
	fieldNames := extractFieldNames(result.Errors)
	validationError := classifyValidationError(fieldNames, result.Errors)

	if len(fieldNames) > 0 {
		c.Set("validation_fields", translateFieldNames(fieldNames))
	}

	if log != nil {
		log.Warn(logger.LogMiddlewareValidationFailed, "path", c.Request.URL.Path, "fields", fieldNames)
	}
	_ = c.Error(validationError)
	c.Abort()
}

// ============================================
// jsonValidator helpers (extracted to reduce cognitive complexity)
// ============================================

// extractFieldNames parses validation error params to collect field names.
// Handles both "properties" (plural, multiple fields) and "property" (singular).
func extractFieldNames(errors map[string]*jsonschema.EvaluationError) []string {
	var fieldNames []string
	for _, validationError := range errors {
		if validationError.Params == nil {
			continue
		}
		// Try "properties" (plural) first - for multiple fields
		if properties, exists := validationError.Params["properties"]; exists {
			fieldNames = append(fieldNames, parsePropertiesParam(properties)...)
		} else if property, exists := validationError.Params["property"]; exists {
			fieldNames = append(fieldNames, parsePropertiesParam(property)...)
		}
	}
	return fieldNames
}

// parsePropertiesParam splits a comma-separated property value into trimmed field names.
func parsePropertiesParam(value interface{}) []string {
	str := fmt.Sprintf("%v", value)
	var result []string
	for _, field := range strings.Split(str, ",") {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "'\"")
		if field != "" {
			result = append(result, field)
		}
	}
	return result
}

// classifyValidationError determines the specific validation error type
// based on field count and the first error code.
func classifyValidationError(fieldNames []string, errors map[string]*jsonschema.EvaluationError) error {
	if len(fieldNames) > 1 {
		return json_schema.ErrMultipleFields
	}

	// Single field error - determine specific error type
	var firstError *jsonschema.EvaluationError
	for _, err := range errors {
		firstError = err
		break
	}

	if firstError == nil {
		return json_schema.ErrValidationFailed
	}

	switch firstError.Code {
	case "property_mismatch":
		return json_schema.ErrFieldPropertyMismatch
	case "required":
		return json_schema.ErrFieldRequired
	case "type":
		return json_schema.ErrFieldTypeInvalid
	default:
		return json_schema.ErrValidationFailed
	}
}
