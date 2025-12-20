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

func (b *Builder) WithValidateRegister() gin.HandlerFunc {
	b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator)
}

func (b *Builder) WithValidateMessage() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.MessageValidator)
}

func (b *Builder) WithValidateResendVerification() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ResendVerificationValidator)
}

func (b *Builder) WithValidatePasswordReset() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.PasswordResetValidator)
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
			c.Error(json_schema.ErrBodyReadFailed)
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		var data map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			if log != nil {
				log.Error(logger.LogMiddlewareJSONParseError, "error", err, "path", c.Request.URL.Path)
			}
			c.Error(json_schema.ErrBadRequest)
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
			if len(fieldNames) > 0 {
				c.Set("validation_fields", fieldNames)
			}

			if log != nil {
				log.Warn(logger.LogMiddlewareValidationFailed, "path", c.Request.URL.Path, "fields", fieldNames)
			}
			c.Error(validationError)
			c.Abort()
			return
		}

		if log != nil {
			log.Debug(logger.LogMiddlewareValidationOK, "path", c.Request.URL.Path)
		}
		c.Next()
	}
}
