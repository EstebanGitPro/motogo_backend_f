package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info("Nueva solicitud de registro recibida",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			h.Logger.Error("Error parseando JSON del request",
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		h.Logger.Info("Procesando registro de usuario",
			"email", personRequest.Email,
			"role", personRequest.Role)

		result, err := h.Interactor.RegisterPerson(c, personRequest.ToDomain())
		if err != nil {
			h.Logger.Error("Error en proceso de registro",
				"email", personRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host
		links := BuildAccountLinks(baseURL, result.Person.ID)

		locationURL := baseURL + "/motogo/api/v1/accounts/" + result.Person.ID
		c.Header("Location", locationURL)

		response := RegistrationResponse{
			Message: result.Message,
			Links:   links,
		}

		h.Logger.Success("Registro completado exitosamente",
			"email", result.Person.Email,
			"person_id", result.Person.ID,
			"status", http.StatusCreated,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusCreated, response)
	}
}
