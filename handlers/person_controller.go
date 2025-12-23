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

		
		

		log.Success("Registro completado exitosamente",
			result.Person.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for person registration
		middleware.RecordPersonRegistration()

		h.Response.Success(c, "MOD_U_REG_EXI_00001")
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

// @Summary Login de usuario
// @Description Autentica un usuario y retorna tokens de acceso
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciales de login"
// @Success 200 {object} middleware.APIResponse{data=LoginResponse} "Login exitoso"
// @Failure 400 {object} middleware.APIResponse "Credenciales inválidas"
// @Failure 401 {object} middleware.APIResponse "Email no verificado o credenciales incorrectas"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/login [post]
func (h handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogKeycloakUserLogin, "email", req.Email, "client_ip", c.ClientIP())

		// Llamar al servicio de autenticación de Keycloak
		token, err := h.Interactor.Login(c, req.Email, req.Password)
		if err != nil {
			log.Error(logger.LogKeycloakUserLoginError, "email", req.Email, "error", err, "client_ip", c.ClientIP())

			// Map specific errors to appropriate messages
			switch err {
			case domain.ErrorEmailNotVerified:
				// Email not verified - auto-resend was triggered
				h.Response.Error(c, domain.MsgUserEmailNotVerified)
			case domain.ErrUserNotFound:
				// User not found - use generic unauthorized message for security
				h.Response.Error(c, domain.MsgUnauthorized)
			default:
				// All other errors (wrong password, etc) - generic message
				h.Response.Error(c, domain.MsgUnauthorized)
			}
			return
		}

		// Build HATEOAS links for login response
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildLoginLinks(baseURL)

		response := LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
			Links:        hateoasLinks,
		}

		log.Success(logger.LogKeycloakUserLoginOK, "email", req.Email, "client_ip", c.ClientIP())
		middleware.RecordPersonRegistration() // Por ahora usamos el mismo metric
		h.Response.SuccessWithData(c, "MOD_AUTH_LOGIN_SUCCESS_EXI_00001", response)
	}
}

// @Summary Verificar email de usuario (Proxy)
// @Description Verifica el email de un usuario usando un token JWT. Este endpoint actúa como proxy para no exponer Keycloak directamente.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Token de verificación del email"
// @Success 200 {object} middleware.APIResponse{data=VerifyEmailResponse} "Email verificado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Token inválido o expirado"
// @Failure 404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure 409 {object} middleware.APIResponse "Email ya estaba verificado"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/verify-email [post]
func (h handler) VerifyEmailByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req VerifyEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogKeycloakEmailVerify, "client_ip", c.ClientIP())

		// Pasar el token al Interactor - la extracción del email se hace en la capa de negocio
		email, err := h.Interactor.VerifyEmailByToken(c, req.Token)
		if err != nil {
			switch err {
			case domain.ErrInvalidToken:
				h.Response.Error(c, domain.MsgKCInvalidToken)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, domain.MsgKCEmailAlreadyVerified)
			default:
				h.Response.Error(c, domain.MsgKCEmailVerifyError)
			}
			return
		}

		response := VerifyEmailResponse{
			Verified: true,
			Email:    email,
		}

		log.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCEmailVerified, response)
	}
}

// @Summary      Obtener perfil del usuario autenticado
// @Description  Retorna los datos completos del usuario autenticado usando el token JWT
// @Tags         Autenticación
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse{data=AuthMeResponse}  "Perfil obtenido exitosamente"
// @Failure      401  {object}  middleware.APIResponse  "Token inválido o ausente"
// @Failure      404  {object}  middleware.APIResponse  "Usuario no encontrado"
// @Router       /auth/me [get]
func (h handler) GetAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated user from context (injected by JWT middleware)
		person, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error("authenticated user not found in context")
			c.Error(domain.ErrUserNotFound)
			return
		}

		log.Debug("retrieving authenticated user profile", "user_id", person.ID, "email", person.Email)

		// Encode ID for response
		encodedID, err := h.EncodeID(person.ID)
		if err != nil {
			log.Error("error encoding user ID", "error", err, "user_id", person.ID)
			c.Error(err)
			return
		}

		// Build HATEOAS links for auth/me response
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildAuthMeLinks(baseURL)

		// Build response
		response := AuthMeResponse{
			ID:             encodedID,
			Email:          person.Email,
			FirstName:      person.FirstName,
			LastName:       person.LastName,
			SecondLastName: person.SecondLastName,
			PhoneNumber:    person.PhoneNumber,
			Role:           person.Role,
			Links:          hateoasLinks,
		}

		log.Success("user profile retrieved successfully", "user_id", encodedID, "email", person.Email)
		h.Response.SuccessWithData(c, domain.MsgAuthProfileRetrieved, response)
	}
}
