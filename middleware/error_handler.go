package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	json_schema "github.com/EstebanGitPro/motogo-backend/platform/schema"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

var errorStatusMap = map[string]int{
	// User Management Errors
	"MOD_U_USU_ERR_00001": http.StatusConflict,             // ErrDuplicateUser
	"MOD_U_USU_ERR_00002": http.StatusInternalServerError,  // ErrUserCannotSave
	"MOD_U_USU_ERR_00003": http.StatusNotFound,             // ErrPersonNotFound
	"MOD_U_USU_ERR_00004": http.StatusInternalServerError,  // ErrGettingUserByEmail
	"MOD_U_USU_ERR_00005": http.StatusNotFound,             // ErrNotFoundUserByEmail
	"MOD_U_USU_ERR_00006": http.StatusNotFound,             // ErrUserCannotFound
	"MOD_U_USU_ERR_00007": http.StatusInternalServerError,  // ErrUserCannotGet
	"MOD_U_USU_ERR_00008": http.StatusForbidden,            // ErrorEmailNotVerified
	"MOD_U_USU_ERR_00009": http.StatusNotFound,             // ErrVerificationTokenNotFound
	"MOD_U_USU_ERR_00010": http.StatusGone,                 // ErrTokenExpired
	"MOD_U_USU_ERR_00011": http.StatusConflict,             // ErrTokenAlreadyUsed
	"MOD_U_USU_ERR_00012": http.StatusInternalServerError,  // ErrRegistrationFailed
	"MOD_U_USU_ERR_00013": http.StatusBadRequest,           // ErrRoleRequired

	// Request Validation Errors
	"MOD_V_VAL_ERR_00001": http.StatusBadRequest,           // ErrInvalidJSONFormat
	"MOD_V_VAL_ERR_00002": http.StatusBadRequest,           // ErrInvalidRequest
	"MOD_V_VAL_ERR_00003": http.StatusInternalServerError,  // Error reading JSON schema
	"MOD_V_VAL_ERR_00004": http.StatusInternalServerError,  // Schema JSON is null or empty
	"MOD_V_VAL_ERR_00005": http.StatusInternalServerError,  // Error compiling schema
	"MOD_V_VAL_ERR_00006": http.StatusBadRequest,           // Schema validation failed
	"MOD_V_VAL_ERR_00007": http.StatusBadRequest,           // Error reading request body
	"MOD_V_VAL_ERR_00008": http.StatusBadRequest,           // Field property mismatch
	"MOD_V_VAL_ERR_00009": http.StatusBadRequest,           // Field required
	"MOD_V_VAL_ERR_00010": http.StatusBadRequest,           // Field type invalid
	"MOD_V_VAL_ERR_00011": http.StatusBadRequest,           // Multiple field errors

	// Authorization Errors
	"MOD_A_AUT_ERR_00001": http.StatusInternalServerError, // ErrRoleAssignmentFailed
	"MOD_A_AUT_ERR_00002": http.StatusInternalServerError, // ErrRoleRemovalFailed
	"MOD_A_AUT_ERR_00003": http.StatusInternalServerError, // ErrRoleCheckFailed
	"MOD_A_AUT_ERR_00004": http.StatusInternalServerError, // ErrGetUserRolesFailed
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() 

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// Try SchemaError first
			var schemaErr *json_schema.SchemaError
			if errors.As(err, &schemaErr) {
				statusCode, exists := errorStatusMap[schemaErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown schema error code: %s", schemaErr.Code)
				}

				response := ErrorResponse{
					Status:  statusCode,
					Code:    schemaErr.Code,
					Message: schemaErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}
			
			// Try DomainError second
			var domainErr *domain.DomainError
			if errors.As(err, &domainErr) {
				statusCode, exists := errorStatusMap[domainErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown domain error code: %s", domainErr.Code)
				}

				response := ErrorResponse{
					Status:  statusCode,
					Code:    domainErr.Code,
					Message: domainErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}

			log.Printf("Non-domain error: %v", err)
			response := ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "MOD_G_GEN_ERR_00001",
				Message: "Error interno del servidor",
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		}
	}
}