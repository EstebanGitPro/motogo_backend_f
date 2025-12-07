package handlers

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterPerson godoc
// @Summary      Registrar nueva cuenta
// @Description  Crea una nueva cuenta de usuario en el sistema con sincronización a Keycloak
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      PersonRequest  true  "Datos de la cuenta"
// @Success      201  {object}  middleware.APIResponse{data=RegistrationResponse}  "Cuenta creada exitosamente"
// @Failure      400  {object}  middleware.APIResponse  "Error de validación"
// @Failure      409  {object}  middleware.APIResponse  "Email o número de identidad ya registrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /accounts [post]
func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := h.Logger.WithTraceID(traceID)

		log.Info(logger.LogRegRequestReceived,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			log.Error(logger.LogRegJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		log.Info(logger.LogRegProcessing,
			"email", personRequest.Email,
			"role", personRequest.Role)

		result, err := h.Interactor.RegisterPerson(c, personRequest.ToDomain())
		if err != nil {
			log.Error(logger.LogRegProcessError,
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
			Links: links,
		}

		log.Success("Registro completado exitosamente",
			result.Person.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for person registration
		middleware.RecordPersonRegistration()

		h.Response.SuccessWithData(c, domain.MsgPersonRegistered, response)
	}
}
