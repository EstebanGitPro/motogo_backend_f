package middleware

import (
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

//TODO: hacer esto 
var mapError = map[error]ErrorResponse{

	domain.ErrDuplicateUser: {
		Code:    "MOD_U_DUP_ERR_00001",
		Message: "Usuario duplicado",
		Status:  http.StatusConflict,
	},
	domain.ErrPersonNotFound: {
		Code:    "MOD_P_NOT_FOUND_ERR_00001",
		Message: "Persona no encontrada",
		Status:  http.StatusNotFound,
	},
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		
		if len(c.Errors) > 0 {
			
			err := c.Errors.Last().Err

			if response, ok := mapError[err]; ok {
				log.Warn("Error de negocio capturado",
					"error", err.Error(),
					"code", response.Code,
					"status", response.Status,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"client_ip", c.ClientIP())
				c.JSON(response.Status, response)
				return
			}

			log.Error("Error interno del servidor",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP())

			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}
