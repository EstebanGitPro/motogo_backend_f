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

	// General errors
	domain.ErrInternalServer: domain.MsgServerError,
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorHandler struct {
	cache *messagingCache.MessageCache
	log   logger.Logger
}

func NewErrorHandler(cache *messagingCache.MessageCache, log logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		cache: cache,
		log:   log,
	}
}

func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Extract validation field names from context if available
			var params []string
			if validationFields, exists := c.Get("validation_fields"); exists {
				if fields, ok := validationFields.([]string); ok {
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
					h.log.Warn(logger.LogMiddlewareErrorCaught,
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
					})
					return
				}
			}

			// Fallback for unmapped errors
			h.log.Error(logger.LogMiddlewareInternalErr,
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
