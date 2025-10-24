package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

var errorStatusMap = map[string]int{
	// User Management Errors
	"MOD_U_USU_ERR_00001": http.StatusConflict,      // ErrDuplicateUser
	"MOD_U_USU_ERR_00002": http.StatusInternalServerError, // ErrUserCannotSave
	"MOD_U_USU_ERR_00003": http.StatusNotFound,      // ErrPersonNotFound
	"MOD_U_USU_ERR_00004": http.StatusInternalServerError, // ErrGettingUserByEmail
	"MOD_U_USU_ERR_00005": http.StatusNotFound,      // ErrNotFoundUserByEmail
	"MOD_U_USU_ERR_00006": http.StatusNotFound,      // ErrUserCannotFound
	"MOD_U_USU_ERR_00007": http.StatusInternalServerError, // ErrUserCannotGet
	"MOD_U_USU_ERR_00008": http.StatusForbidden,     // ErrorEmailNotVerified
	"MOD_U_USU_ERR_00009": http.StatusNotFound,      // ErrVerificationTokenNotFound
	"MOD_U_USU_ERR_00010": http.StatusGone,          // ErrTokenExpired
	"MOD_U_USU_ERR_00011": http.StatusConflict,      // ErrTokenAlreadyUsed

	// Request Validation Errors
	"MOD_V_VAL_ERR_00001": http.StatusBadRequest,    // ErrInvalidJSONFormat
	"MOD_V_VAL_ERR_00002": http.StatusBadRequest,    // ErrInvalidRequest

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
			
			var domainErr *domain.DomainError
			if errors.As(err, &domainErr) {
				statusCode, exists := errorStatusMap[domainErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown error code: %s", domainErr.Code)
				}

				response := ErrorResponse{
					Status:  statusCode,
					Code:    domainErr.Code,
					Message: domainErr.Message,
				}

				c.JSON(statusCode, response)
				return
			}

			log.Printf("Non-domain error: %v", err)
			response := ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "MOD_G_GEN_ERR_00001",
				Message: "Error interno del servidor",
			}
			c.JSON(http.StatusInternalServerError, response)
		}
	}
}