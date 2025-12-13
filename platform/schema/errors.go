package schema

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Re-export domain validation errors for backward compatibility
// These now integrate with the unified messaging system via error_handler
var (
	ErrBadRequest            = domain.ErrSchemaBadRequest
	ErrInvalidRequest        = domain.ErrSchemaInvalidRequest
	ErrSchemaReadFailed      = domain.ErrSchemaReadFailed
	ErrSchemaEmpty           = domain.ErrSchemaEmpty
	ErrSchemaCompileFailed   = domain.ErrSchemaCompileFailed
	ErrValidationFailed      = domain.ErrSchemaValidationFailed
	ErrBodyReadFailed        = domain.ErrSchemaBodyReadFailed
	ErrFieldPropertyMismatch = domain.ErrSchemaFieldFormat
	ErrFieldRequired         = domain.ErrSchemaFieldRequired
	ErrFieldTypeInvalid      = domain.ErrSchemaFieldType
	ErrMultipleFields        = domain.ErrSchemaMultipleFields
)
