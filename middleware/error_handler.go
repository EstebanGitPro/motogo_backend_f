package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	json_schema "github.com/EstebanGitPro/motogo-backend/platform/schema"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// determineStatusCode determina el código de estado HTTP basándose en el código de error
// Esto mantiene una lógica simple y predecible sin necesidad de un mapa gigante
func determineStatusCode(code string) int {
	// Extraer el módulo del código (posición 4-5 después de "MOD_")
	if len(code) < 7 {
		return http.StatusInternalServerError
	}

	// MOD_V_VAL_ERR = Errores de Validación
	if strings.HasPrefix(code, "MOD_V_VAL_ERR") {
		return http.StatusBadRequest
	}

	// Analizar el tipo de error por el código específico
	switch {
	// User Management - Conflictos y duplicados
	case code == "MOD_U_USU_ERR_00001", // ErrDuplicateUser
		code == "MOD_U_USU_ERR_00011": // ErrTokenAlreadyUsed
		return http.StatusConflict

	// User Management - No encontrado
	case code == "MOD_U_USU_ERR_00003", // ErrPersonNotFound
		code == "MOD_U_USU_ERR_00005", // ErrNotFoundUserByEmail
		code == "MOD_U_USU_ERR_00006", // ErrUserCannotFound
		code == "MOD_U_USU_ERR_00009": // ErrVerificationTokenNotFound
		return http.StatusNotFound

	// User Management - Prohibido
	case code == "MOD_U_USU_ERR_00008": // ErrorEmailNotVerified
		return http.StatusForbidden

	// User Management - Token expirado
	case code == "MOD_U_USU_ERR_00010": // ErrTokenExpired
		return http.StatusGone

	// User Management - Bad Request
	case code == "MOD_U_USU_ERR_00013": // ErrRoleRequired
		return http.StatusBadRequest

	// Authorization Errors - Todos son errores internos del servidor
	case strings.HasPrefix(code, "MOD_A_AUT_ERR"):
		return http.StatusInternalServerError

	// User Management - Errores internos del servidor
	case strings.HasPrefix(code, "MOD_U_USU_ERR"):
		return http.StatusInternalServerError

	default:
		// Por defecto, error interno del servidor
		return http.StatusInternalServerError
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Manejar errores de validación de esquema
			var schemaErr *json_schema.SchemaError
			if errors.As(err, &schemaErr) {
				statusCode := determineStatusCode(schemaErr.Code)

				response := ErrorResponse{
					Code:    schemaErr.Code,
					Message: schemaErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}

			// Manejar errores de dominio
			var domainErr *domain.DomainError
			if errors.As(err, &domainErr) {
				statusCode := determineStatusCode(domainErr.Code)

				response := ErrorResponse{
					Code:    domainErr.Code,
					Message: domainErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}

			// Errores no controlados
			log.Printf("Non-domain error: %v", err)
			response := ErrorResponse{
				Code:    "MOD_G_GEN_ERR_00001",
				Message: "Error interno del servidor",
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		}
	}
}