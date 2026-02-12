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
	"email":            "Correo electrónico",
	"password":         "Contraseña",
	"first_name":       "Nombre",
	"last_name":        "Apellido",
	"phone":            "Teléfono",
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

// WithValidateEvidence validates evidence creation request (HU16)
func (b *Builder) WithValidateEvidence() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateEvidenceValidator)
}

func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID for trace correlation
		traceID := GetRequestID(c)
		log := log.WithTraceID(traceID)

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if log != nil {
				log.Error(logger.LogMiddlewareBodyReadError, "error", err, "path", c.Request.URL.Path)
			}
			_ = c.Error(json_schema.ErrBodyReadFailed)
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		var data map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			if log != nil {
				log.Error(logger.LogMiddlewareJSONParseError, "error", err, "path", c.Request.URL.Path)
			}
			_ = c.Error(json_schema.ErrBadRequest)
			c.Abort()
			return
		}

		result := schema.Validate(data)
		if !result.IsValid() {
			// Collect all field names with errors for logging
			var fieldNames []string

			// Process all errors to extract field names
			for _, validationError := range result.Errors {
				if validationError.Params != nil {
					// Try "properties" (plural) first - for multiple fields
					if properties, exists := validationError.Params["properties"]; exists {
						propertiesStr := fmt.Sprintf("%v", properties)
						// Parse multiple properties: "'email', 'role'" -> ["email", "role"]
						fields := strings.Split(propertiesStr, ",")
						for _, field := range fields {
							field = strings.TrimSpace(field)
							// Remove quotes
							field = strings.Trim(field, "'\"")
							if field != "" {
								fieldNames = append(fieldNames, field)
							}
						}
					} else if property, exists := validationError.Params["property"]; exists {
						// Single property
						propertyName := fmt.Sprintf("%v", property)
						// Remove quotes if present
						propertyName = strings.Trim(propertyName, "'\"")
						if propertyName != "" {
							fieldNames = append(fieldNames, propertyName)
						}
					}
				}
			}

			// Determine which validation error to use based on validation type
			var validationError error

			// If multiple fields failed, use specific multiple fields error
			if len(fieldNames) > 1 {
				validationError = json_schema.ErrMultipleFields
			} else {
				// Single field error - determine specific error type
				var firstError *jsonschema.EvaluationError
				for _, err := range result.Errors {
					firstError = err
					break
				}

				if firstError != nil {
					switch firstError.Code {
					case "property_mismatch":
						validationError = json_schema.ErrFieldPropertyMismatch
					case "required":
						validationError = json_schema.ErrFieldRequired
					case "type":
						validationError = json_schema.ErrFieldTypeInvalid
					default:
						validationError = json_schema.ErrValidationFailed
					}
				} else {
					validationError = json_schema.ErrValidationFailed
				}
			}

			// Store field names in context for error_handler to use in message parameters
			// Translate technical names to Spanish labels
			if len(fieldNames) > 0 {
				c.Set("validation_fields", translateFieldNames(fieldNames))
			}

			if log != nil {
				log.Warn(logger.LogMiddlewareValidationFailed, "path", c.Request.URL.Path, "fields", fieldNames)
			}
			_ = c.Error(validationError)
			c.Abort()
			return
		}

		if log != nil {
			log.Debug(logger.LogMiddlewareValidationOK, "path", c.Request.URL.Path)
		}
		c.Next()
	}
}
