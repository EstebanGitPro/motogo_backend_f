package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info(logger.LogRegRequestReceived,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			h.Logger.Error(logger.LogRegJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		h.Logger.Info(logger.LogRegProcessing,
			"email", personRequest.Email,
			"role", personRequest.Role)

		result, err := h.Interactor.RegisterPerson(c, personRequest.ToDomain())
		if err != nil {
			h.Logger.Error(logger.LogRegProcessError,
				"email", personRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Ofuscar el ID antes de exponerlo en la API
		encodedID, err := h.EncodeID(result.Person.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Person.ID, err)
			return
		}

		// Construir respuesta con HATEOAS
		baseURL := GetBaseURL(c)
		links := BuildAccountLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "accounts", encodedID)

		response := RegistrationResponse{
			Message: result.Message,
			Links:   links,
		}

		h.Logger.Success("Registro completado exitosamente",
			result.Person.ToLogger(),
			"encoded_id", encodedID,
			"status", http.StatusCreated,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusCreated, response)
	}
}
