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
		log := Logger.WithTraceID(traceID)

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

// @Summary Reenviar email de verificación
// @Description Reenvía el email de verificación a un usuario registrado que aún no ha verificado su correo
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body ResendVerificationEmailRequest true "Email del usuario"
// @Success 200 {object} middleware.APIResponse "Email reenviado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Email inválido o ya verificado"
// @Failure 404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/resend-verification [post]
func (h handler) ResendVerificationEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResendVerificationEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		err := h.Interactor.ResendVerificationEmail(c, req.Email)
		if err != nil {
			// Manejar diferentes tipos de errores
			switch err {
			case domain.ErrUserNotFound:
				h.Response.Error(c, "MOD_KC_USER_NOT_FOUND_ERR_00001")
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001")
			default:
				h.Response.Error(c, "MOD_KC_VERIF_EMAIL_ERROR_ERR_00001")
			}
			return
		}

		h.Response.Success(c, "MOD_KC_VERIF_EMAIL_SENT_EXI_00001", req.Email)
	}
}

// @Summary Solicitar recuperación de contraseña
// @Description Envía un email con instrucciones para restablecer la contraseña del usuario
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body PasswordResetRequest true "Email del usuario"
// @Success 200 {object} middleware.APIResponse "Solicitud procesada (siempre retorna éxito por seguridad)"
// @Failure 400 {object} middleware.APIResponse "Email inválido"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/password-reset [post]
func (h handler) RequestPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordResetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Este método SIEMPRE retorna nil por seguridad (no revela si el email existe)
		// El logging interno sí registra el resultado real
		_ = h.Interactor.RequestPasswordReset(c, req.Email)

		// Siempre responder con éxito genérico
		h.Response.Success(c, "MOD_KC_PWD_RESET_SENT_EXI_00001")
	}
}
