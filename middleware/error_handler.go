package middleware

import (
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// errorToMessageCode maps domain errors to message codes
// Messages are loaded from cache (DB) - not hardcoded
var errorToMessageCode = map[error]string{
	// User Management Errors (MOD_U_*)
	domain.ErrDuplicateUser:             domain.MsgUserDuplicate,
	domain.ErrUserNotFound:              domain.MsgUserNotFound, // User authenticated in Keycloak but not in local DB
	domain.ErrUserCannotSave:            domain.MsgUserCannotSave,
	domain.ErrUserCannotFound:           domain.MsgUserNotFound,
	domain.ErrUserCannotGet:             domain.MsgUserNotFound,
	domain.ErrNotFoundUserByEmail:       domain.MsgUserEmailNotFound,
	domain.ErrGettingUserByEmail:        domain.MsgUserEmailError,
	domain.ErrorEmailNotVerified:        domain.MsgUserEmailNotVerified,
	domain.ErrVerificationTokenNotFound: domain.MsgUserTokenNotFound,
	domain.ErrTokenExpired:              domain.MsgUserTokenExpired,
	domain.ErrTokenAlreadyUsed:          domain.MsgUserTokenUsed,
	domain.ErrRegistrationFailed:        domain.MsgUserRegError,
	domain.ErrRoleRequired:              domain.MsgUserRoleRequired,
	domain.ErrUserCannotDelete:          domain.MsgUserCannotDelete,

	// Person errors
	domain.ErrPersonNotFound:     domain.MsgPersonNotFound,
	domain.ErrInvalidTransaction: domain.MsgPersonInvalidTx,

	// Validation errors
	domain.ErrInvalidJSONFormat: domain.MsgValJSONInvalid,
	domain.ErrInvalidRequest:    domain.MsgValInvalidReq,
	domain.ErrInvalidID:         domain.MsgValIDInvalid,

	// Schema validation errors
	domain.ErrSchemaBadRequest:       domain.MsgValBadFormat,
	domain.ErrSchemaInvalidRequest:   domain.MsgValInvalidReq,
	domain.ErrSchemaReadFailed:       domain.MsgValSchemaRead,
	domain.ErrSchemaEmpty:            domain.MsgValSchemaEmpty,
	domain.ErrSchemaCompileFailed:    domain.MsgValSchemaCompile,
	domain.ErrSchemaValidationFailed: domain.MsgValFailed,
	domain.ErrSchemaBodyReadFailed:   domain.MsgValBodyRead,
	domain.ErrSchemaFieldFormat:      domain.MsgValFieldFormat,
	domain.ErrSchemaFieldRequired:    domain.MsgValFieldRequired,
	domain.ErrSchemaFieldType:        domain.MsgValFieldType,
	domain.ErrSchemaMultipleFields:   domain.MsgValMultiple,

	// Authorization errors
	domain.ErrRoleAssignmentFailed: domain.MsgRoleAssignError,
	domain.ErrRoleRemovalFailed:    domain.MsgRoleRemoveError,
	domain.ErrRoleCheckFailed:      domain.MsgRoleCheckError,
	domain.ErrGetUserRolesFailed:   domain.MsgRoleGetError,

	// Message errors
	domain.ErrMessageNotFound:         domain.MsgMessageNotFound,
	domain.ErrMessageCodeRequired:     domain.MsgMessageCodeRequired,
	domain.ErrMessageTypeRequired:     domain.MsgMessageTypeRequired,
	domain.ErrMessageTitleRequired:    domain.MsgMessageTitleRequired,
	domain.ErrMessageContentRequired:  domain.MsgMessageContentReq,
	domain.ErrMessageModuleRequired:   domain.MsgMessageModuleRequired,
	domain.ErrMessageCategoryRequired: domain.MsgMessageCategoryReq,
	domain.ErrMessageCodeDuplicate:    domain.MsgMessageCodeDuplicate,
	domain.ErrMessageCannotSave:       domain.MsgMessageSaveError,
	domain.ErrMessageCannotUpdate:     domain.MsgMessageUpdateError,
	domain.ErrMessageCannotDelete:     domain.MsgMessageDeleteError,
	domain.ErrMessageInvalidType:      domain.MsgMessageInvalidType,
	domain.ErrMessageListFailed:       domain.MsgMessageListError,

	// Infrastructure errors (MOD_INFRA_*)
	domain.ErrKeycloakInconsistentState:  domain.MsgKeycloakInconsistentState,
	domain.ErrKeycloakUserCreationFailed: domain.MsgKeycloakCreateError,
	domain.ErrKeycloakCleanupFailed:      domain.MsgKeycloakCleanupError,
	// Dependency availability errors
	domain.ErrKeycloakUnavailable: domain.MsgKeycloakUnavailable,
	domain.ErrDatabaseUnavailable: domain.MsgDatabaseUnavailable,
	// Incomplete registration (cleanup in progress)
	domain.ErrIncompleteRegistration: domain.MsgIncompleteRegistration,

	// Authentication errors (401 Unauthorized)
	domain.ErrInvalidToken:       domain.MsgUnauthorized,
	domain.ErrInvalidCredentials: domain.MsgUnauthorized,
	domain.ErrTokenExpired:       domain.MsgUserTokenExpired,

	// Motorcycle errors (MOD_MOT_*)
	domain.ErrMotorcycleNotFound:     domain.MsgMotorcycleNotFound,
	domain.ErrMotorcycleCannotSave:   domain.MsgMotorcycleCannotSave,
	domain.ErrMotorcycleCannotUpdate: domain.MsgMotorcycleCannotUpdate,
	domain.ErrMotorcycleCannotDelete: domain.MsgMotorcycleCannotDelete,
	domain.ErrDuplicateLicensePlate:  domain.MsgDuplicateLicensePlate,
	domain.ErrReferenceNotFound:      domain.MsgMotorcycleReferenceNotFound,
	domain.ErrReferenceRequired:      domain.MsgReferenceRequired,
	domain.ErrBranchNotAuthorized:    domain.MsgBranchNotAuthorized,

	// Diagnostic Permission errors (MOD_DGP_*)
	domain.ErrPermissionNotFound:     domain.MsgPermissionNotFound,
	domain.ErrPermissionCannotSave:   domain.MsgPermissionCannotSave,
	domain.ErrPermissionCannotDelete: domain.MsgPermissionCannotDelete,

	// Person deletion errors (HU53)
	domain.ErrPersonHasBranches: domain.MsgPersonHasBranches,

	// Password management errors (HU57)
	domain.ErrEmailAlreadyVerified:    domain.MsgKCEmailAlreadyVerified,
	domain.ErrPasswordUpdateFailed:    domain.MsgChangePasswordUpdateError,
	domain.ErrPasswordPolicyViolation: domain.MsgChangePasswordPolicyError,
	domain.ErrForbidden:               domain.MsgForbidden,

	// General errors
	domain.ErrInternalServer: domain.MsgServerError,
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"` // HU36: Show which fields failed validation
}

var log logger.Logger = logger.NewSlogLogger()

type ErrorHandler struct {
	cache *messagingCache.MessageCache
}

func NewErrorHandler(cache *messagingCache.MessageCache) *ErrorHandler {
	return &ErrorHandler{
		cache: cache,
	}
}

func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get request ID for trace correlation
			traceID := GetRequestID(c)
			log := log.WithTraceID(traceID)

			// Extract validation field names from context if available
			var params []string
			var validationFieldsList []string
			if validationFields, exists := c.Get("validation_fields"); exists {
				if fields, ok := validationFields.([]string); ok {
					validationFieldsList = fields // Keep original list for response
					// For multiple fields error, concatenate all field names into one parameter
					if len(fields) > 1 {
						// Join fields with comma for multiple fields message
						fieldsStr := fields[0]
						for i := 1; i < len(fields); i++ {
							fieldsStr += ", " + fields[i]
						}
						params = []string{fieldsStr}
					} else {
						params = fields
					}
				}
			}

			// Try to map domain error to message code
			if messageCode, ok := errorToMessageCode[err]; ok {
				// Get message from cache (or DB if not cached) with field params
				msg := h.cache.GetMessageResponse(messageCode, params...)
				status := h.cache.GetHTTPStatus(messageCode)

				if msg != nil {
					log.Warn(logger.LogMiddlewareErrorCaught,
						"error", err.Error(),
						"code", msg.Code,
						"status", status,
						"fields", params,
						"path", c.Request.URL.Path,
						"method", c.Request.Method,
						"client_ip", c.ClientIP())

					c.JSON(status, ErrorResponse{
						Success: false,
						Code:    msg.Code,
						Message: msg.Content,
						Fields:  validationFieldsList, // Include fields that failed validation
					})
					return
				}
			}

			// Fallback for unmapped errors
			log.Error(logger.LogMiddlewareInternalErr,
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP())

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Code:    domain.MsgServerError,
				Message: "Error interno del servidor",
			})
		}
	}
}
