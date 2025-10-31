package schema

import (
	"errors"
	"fmt"
)

type SchemaError struct {
	Code    string
	Message string
	err     error
}

func (e *SchemaError) Error() string {
	return e.Message
}

func (e *SchemaError) Unwrap() error {
	return e.err
}


func NewSchemaError(code, message string) *SchemaError {
	return &SchemaError{
		Code:    code,
		Message: message,
		err:     errors.New(message),
	}
}

// NewFieldSchemaError creates a schema error with field-specific information
func NewFieldSchemaError(baseError *SchemaError, fieldName string) *SchemaError {
	message := fmt.Sprintf("%s: %s", baseError.Message, fieldName)
	return &SchemaError{
		Code:    baseError.Code,
		Message: message,
		err:     errors.New(message),
	}
}

// NewMultipleFieldSchemaError creates a schema error for multiple fields
func NewMultipleFieldSchemaError(fieldNames []string) *SchemaError {
	message := "Errores de validación en los campos: " + fieldNames[0]
	for i := 1; i < len(fieldNames); i++ {
		message += ", " + fieldNames[i]
	}
	return &SchemaError{
		Code:    "MOD_V_VAL_ERR_00011",
		Message: message,
		err:     errors.New(message),
	}
}


var (
	ErrBadRequest = NewSchemaError("MOD_V_VAL_ERR_00001", "Formato de solicitud no válido")
	ErrInvalidRequest = NewSchemaError("MOD_V_VAL_ERR_00002", "Solicitud inválida")
	ErrSchemaReadFailed = NewSchemaError("MOD_V_VAL_ERR_00003", "Error leyendo esquema JSON")
	ErrSchemaEmpty = NewSchemaError("MOD_V_VAL_ERR_00004", "El esquema JSON está vacío o es nulo")
	ErrSchemaCompileFailed = NewSchemaError("MOD_V_VAL_ERR_00005", "Error compilando esquema")
	ErrValidationFailed = NewSchemaError("MOD_V_VAL_ERR_00006", "Falló la validación del esquema")
	ErrBodyReadFailed = NewSchemaError("MOD_V_VAL_ERR_00007", "Error leyendo cuerpo de la solicitud")
	ErrFieldPropertyMismatch = NewSchemaError("MOD_V_VAL_ERR_00008", "El campo no cumple con el formato requerido")
	ErrFieldRequired = NewSchemaError("MOD_V_VAL_ERR_00009", "El campo es requerido")
	ErrFieldTypeInvalid = NewSchemaError("MOD_V_VAL_ERR_00010", "El campo tiene un tipo de dato incorrecto")
)


