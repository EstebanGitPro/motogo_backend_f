package handlers

import (
	"errors"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	cookiepkg "github.com/EstebanGitPro/motogo-backend/platform/cookie"
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
// @Router       /persons [post]
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
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		// Sanitize input (trim whitespace)
		personRequest.Sanitize()

		log.Info(logger.LogRegProcessing,
			"email", personRequest.Email,
			"role", personRequest.Role)

		result, err := h.Interactor.RegisterPerson(c, personRequest.ToDomain())
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", personRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		// Ofuscar el ID antes de exponerlo en la API
		encodedID, err := h.EncodeID(result.Person.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Person.ID, err)
			return
		}

		log.Success(logger.LogPersonControllerRegComplete,
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
		if c.ShouldBindJSON(&req) != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input
		req.Sanitize()

		err := h.Interactor.ResendVerificationEmail(c, req.Email)
		if err != nil {
			// Manejar diferentes tipos de errores
			switch {
			case errors.Is(err, domain.ErrUserNotFound):
				h.Response.Error(c, "MOD_KC_USER_NOT_FOUND_ERR_00001")
			case errors.Is(err, domain.ErrEmailAlreadyVerified):
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
		if c.ShouldBindJSON(&req) != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input
		req.Sanitize()

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

		// Sanitize input
		req.Sanitize()

		log.Info(logger.LogKeycloakUserLogin, "email", req.Email, "client_ip", c.ClientIP())

		// Llamar al servicio de autenticación de Keycloak
		token, err := h.Interactor.Login(c, req.Email, req.Password)
		if err != nil {
			log.Error(logger.LogKeycloakUserLoginError, "email", req.Email, "error", err, "client_ip", c.ClientIP())

			// Map specific errors to appropriate messages
			switch {
			case errors.Is(err, domain.ErrEmailNotVerified):
				h.Response.Error(c, domain.MsgUserEmailNotVerified)
			case errors.Is(err, domain.ErrKeycloakUnavailable):
				h.Response.Error(c, domain.MsgKeycloakUnavailable)
			case errors.Is(err, domain.ErrUserNotFound):
				h.Response.Error(c, domain.MsgUnauthorized)
			default:
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

		// Escribir cookies HttpOnly (preparación para Flutter Web)
		cookiepkg.SetAccessToken(c, token.AccessToken, h.CookieConfig)
		cookiepkg.SetRefreshToken(c, token.RefreshToken, h.CookieConfig)

		h.Response.SuccessWithData(c, "MOD_AUTH_LOGIN_SUCCESS_EXI_00001", response)
	}
}

// @Summary Refrescar access token
// @Description Obtiene un nuevo access token usando el refresh token. El frontend debe llamar este endpoint cuando reciba 401 por token expirado.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token actual"
// @Success 200 {object} middleware.APIResponse{data=LoginResponse} "Token refrescado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Formato inválido"
// @Failure 401 {object} middleware.APIResponse "Refresh token inválido o expirado"
// @Router /auth/refresh [post]
func (h handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// ShouldBindJSON puede fallar con body vacío; es válido si el token viene por cookie
			log.Debug("RefreshToken: body vacío o inválido, intentando cookie fallback", "error", err)
		}

		// Fallback: si no vino en el body, intentar leerlo de la cookie
		if req.RefreshToken == "" {
			cookieToken, err := cookiepkg.GetRefreshToken(c)
			if err != nil || cookieToken == "" {
				log.Error(logger.LogRegJSONParseError, "error", "refresh_token no proporcionado en body ni cookie")
				h.Response.Error(c, domain.MsgValBadFormat)
				return
			}
			req.RefreshToken = cookieToken
		}

		log.Info(logger.LogPersonControllerRefreshRequest, "client_ip", c.ClientIP())

		// Call Keycloak to refresh the token
		token, err := h.Interactor.RefreshToken(c, req.RefreshToken)
		if err != nil {
			log.Error(logger.LogPersonControllerRefreshError, "error", err, "client_ip", c.ClientIP())
			// All refresh errors return 401 - user needs to login again
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Build HATEOAS links for refresh response (same as login)
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildLoginLinks(baseURL)

		response := LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
			Links:        hateoasLinks,
		}

		log.Success(logger.LogPersonControllerTokenRefreshOK, "client_ip", c.ClientIP())

		// Actualizar cookies HttpOnly con los nuevos tokens
		cookiepkg.SetAccessToken(c, token.AccessToken, h.CookieConfig)
		cookiepkg.SetRefreshToken(c, token.RefreshToken, h.CookieConfig)

		h.Response.SuccessWithData(c, "MOD_AUTH_REFRESH_SUCCESS_EXI_00001", response)
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
			switch {
			case errors.Is(err, domain.ErrInvalidToken):
				h.Response.Error(c, domain.MsgKCInvalidToken)
			case errors.Is(err, domain.ErrUserNotFound):
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case errors.Is(err, domain.ErrEmailAlreadyVerified):
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

// @Summary Actualizar contraseña con token
// @Description Actualiza la contraseña de un usuario usando un token JWT del email de recuperación
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body ResetPasswordWithTokenRequest true "Token y nueva contraseña"
// @Success 200 {object} middleware.APIResponse "Contraseña actualizada exitosamente"
// @Failure 400 {object} middleware.APIResponse "Token inválido o contraseña no cumple requisitos"
// @Failure 404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/password/reset [post]
func (h handler) ResetPasswordWithToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req ResetPasswordWithTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogPasswordResetStart, "client_ip", c.ClientIP())

		// Llamar al servicio para reset de contraseña
		err := h.Interactor.ResetPasswordWithToken(c, req.Token, req.NewPassword)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidToken):
				log.Error(logger.LogPasswordResetTokenError, "error", err, "client_ip", c.ClientIP())
				h.Response.Error(c, "MOD_P_RESET_ERR_00001")
			case errors.Is(err, domain.ErrUserNotFound):
				log.Error(logger.LogPasswordResetUserNotFound, "error", err, "client_ip", c.ClientIP())
				h.Response.Error(c, "MOD_P_RESET_ERR_00002")
			case errors.Is(err, domain.ErrPasswordUpdateFailed):
				log.Error(logger.LogPasswordResetUpdateError, "error", err, "client_ip", c.ClientIP())
			case errors.Is(err, domain.ErrPasswordPolicyViolation):
				h.Response.Error(c, domain.MsgChangePasswordPolicyError)
				h.Response.Error(c, "MOD_P_RESET_ERR_00003")
			default:
				log.Error(logger.LogPasswordResetUpdateError, "error", err, "client_ip", c.ClientIP())
				h.Response.Error(c, "MOD_P_RESET_ERR_00003")
			}
			return
		}

		log.Success(logger.LogPasswordResetSuccess, "client_ip", c.ClientIP())
		h.Response.Success(c, "MOD_P_RESET_EXI_00001")
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
// @Router       /persons/me [get]
func (h handler) GetAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated user from context (injected by JWT middleware)
		person, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogPersonAuthNotFoundInContext)
			_ = c.Error(domain.ErrUserNotFound)
			return
		}

		log.Debug(logger.LogPersonRetrievingProfile, "user_id", person.ID, "email", person.Email)

		// Encode ID for response
		encodedID, err := h.EncodeID(person.ID)
		if err != nil {
			log.Error(logger.LogPersonIDEncodeError, "error", err, "user_id", person.ID)
			_ = c.Error(err)
			return
		}

		// Build HATEOAS links for persons/me response
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildAuthMeLinks(baseURL)

		// Build response
		response := AuthMeResponse{
			ID:             encodedID,
			IdentityNumber: person.IdentityNumber,
			Email:          person.Email,
			FirstName:      person.FirstName,
			LastName:       person.LastName,
			SecondLastName: person.SecondLastName,
			PhoneNumber:    person.PhoneNumber,
			Role:           string(person.Role),
			Links:          hateoasLinks,
		}

		log.Success(logger.LogPersonControllerProfileGetOK, "user_id", encodedID, "email", person.Email)
		h.Response.SuccessWithData(c, domain.MsgAuthProfileRetrieved, response)
	}
}

// @Summary Cambiar contraseña del usuario autenticado
// @Description Permite al usuario autenticado cambiar su contraseña proporcionando la contraseña actual
// @Tags Autenticación
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Contraseña actual y nueva"
// @Success 200 {object} middleware.APIResponse{data=ChangePasswordResponse} "Contraseña actualizada exitosamente"
// @Failure 400 {object} middleware.APIResponse "Formato inválido o contraseña no cumple requisitos"
// @Failure 401 {object} middleware.APIResponse "Token inválido o contraseña actual incorrecta"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /persons/me/password [put]
func (h handler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated user from context (injected by JWT middleware)
		person, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogPersonControllerUserNotInContext)
			_ = c.Error(domain.ErrUserNotFound)
			return
		}

		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogChangePasswordStart, "user_id", person.KeycloakUserID, "client_ip", c.ClientIP())

		// Call interactor to change password
		err := h.Interactor.ChangePassword(c, person.KeycloakUserID, req.CurrentPassword, req.NewPassword)
		if err != nil {
			log.Error(logger.LogChangePasswordUpdateError, "user_id", person.KeycloakUserID, "error", err, "client_ip", c.ClientIP())

			switch {
			case errors.Is(err, domain.ErrInvalidCredentials):
				h.Response.Error(c, domain.MsgChangePasswordInvalidCurrent)
			case errors.Is(err, domain.ErrUserNotFound):
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case errors.Is(err, domain.ErrPasswordUpdateFailed):
				h.Response.Error(c, domain.MsgChangePasswordUpdateError)
			case errors.Is(err, domain.ErrPasswordPolicyViolation):
				h.Response.Error(c, domain.MsgChangePasswordPolicyError)
			case errors.Is(err, domain.ErrKeycloakUnavailable):
				h.Response.Error(c, domain.MsgKeycloakUnavailable)
			default:
				h.Response.Error(c, domain.MsgChangePasswordUpdateError)
			}
			return
		}

		// Build HATEOAS links for change password response
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildChangePasswordLinks(baseURL)

		response := ChangePasswordResponse{
			Message: "Password changed successfully",
			Links:   hateoasLinks,
		}

		log.Success(logger.LogChangePasswordSuccess, "user_id", person.KeycloakUserID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgChangePasswordSuccess, response)
	}
}

// @Summary Actualizar perfil del usuario autenticado
// @Description Permite al usuario autenticado actualizar su información de perfil (excepto email y contraseña)
// @Tags Autenticación
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Datos del perfil a actualizar"
// @Success 200 {object} middleware.APIResponse{data=UpdateProfileResponse} "Perfil actualizado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Formato inválido"
// @Failure 401 {object} middleware.APIResponse "Token inválido o ausente"
// @Failure 409 {object} middleware.APIResponse "Datos duplicados (teléfono o número de identidad)"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /persons/me [put]
func (h handler) UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated user from context (injected by JWT middleware)
		person, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogPersonControllerUserNotInContext)
			_ = c.Error(domain.ErrUserNotFound)
			return
		}

		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input
		req.Sanitize()

		log.Info(logger.LogUpdateProfileStart, "user_id", person.ID, "client_ip", c.ClientIP())

		// Merge request with existing person data (only update provided fields)
		updatedPerson := mergeProfileFields(*person, req)

		// Call interactor to update profile
		result, err := h.Interactor.UpdateProfile(c, updatedPerson)
		if err != nil {
			log.Error(logger.LogUpdateProfileError, "user_id", person.ID, "error", err, "client_ip", c.ClientIP())
			h.mapUpdateProfileError(c, err)
			return
		}

		// Encode ID for response
		encodedID, err := h.EncodeID(result.ID)
		if err != nil {
			log.Error(logger.LogPersonControllerIDEncodeError, "error", err, "user_id", result.ID)
			_ = c.Error(err)
			return
		}

		// Build HATEOAS links for profile update response
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildUpdateProfileLinks(baseURL)

		response := UpdateProfileResponse{
			ID:             encodedID,
			IdentityNumber: result.IdentityNumber,
			Email:          result.Email,
			FirstName:      result.FirstName,
			LastName:       result.LastName,
			SecondLastName: result.SecondLastName,
			PhoneNumber:    result.PhoneNumber,
			Role:           string(result.Role),
			Links:          hateoasLinks,
		}

		log.Success(logger.LogUpdateProfileSuccess, "user_id", encodedID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgPersonUpdated, response)
	}
}

// mergeProfileFields merges non-empty request fields into the existing person.
func mergeProfileFields(person domain.Person, req UpdateProfileRequest) domain.Person {
	if req.IdentityNumber != "" {
		person.IdentityNumber = req.IdentityNumber
	}
	if req.FirstName != "" {
		person.FirstName = req.FirstName
	}
	if req.LastName != "" {
		person.LastName = req.LastName
	}
	if req.SecondLastName != "" {
		person.SecondLastName = req.SecondLastName
	}
	if req.PhoneNumber != "" {
		person.PhoneNumber = req.PhoneNumber
	}
	return person
}

// mapUpdateProfileError maps profile update errors to API responses.
func (h handler) mapUpdateProfileError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrDuplicateUser):
		h.Response.Error(c, domain.MsgUserDuplicate)
	default:
		h.Response.Error(c, domain.MsgPersonUpdated)
	}
}

// @Summary Obtener información de contacto pública de un representante
// @Description Permite a motociclistas obtener el número de contacto de un representante de sede
// @Tags Personas
// @Produce json
// @Param id path string true "ID ofuscado del representante"
// @Success 200 {object} middleware.APIResponse{data=PublicContactResponse} "Info de contacto obtenida"
// @Failure 400 {object} middleware.APIResponse "ID inválido"
// @Failure 404 {object} middleware.APIResponse "Representante no encontrado"
// @Router /persons/{id}/contact [get]
func (h handler) GetPublicContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get and decode person ID from URL param
		encodedID := c.Param("id")
		if encodedID == "" {
			log.Error(logger.LogPersonMissingIDInURL)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		personID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Error(logger.LogPersonIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogPersonControllerGetContactStart, "person_id", personID, "client_ip", c.ClientIP())

		// Get person from interactor (Clean Architecture - HU55)
		person, err := h.Interactor.GetPublicContact(c, personID)
		if err != nil {
			log.Error(logger.LogPersonGetError, "person_id", personID, "error", err)
			_ = c.Error(domain.ErrPersonNotFound)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		hateoasLinks := BuildPublicContactLinks(baseURL, encodedID)

		// Return only phone number (non-sensitive contact info)
		response := PublicContactResponse{
			PhoneNumber: person.PhoneNumber,
			Links:       hateoasLinks,
		}

		log.Success(logger.LogPersonControllerContactGetOK, "person_id", personID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgPersonContactRetrieved, response)
	}
}

// DeleteSelf handles DELETE /persons/me - deletes authenticated user's account (HU53)
// This is a self-delete only - users can only delete their own account
func (h handler) DeleteSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated user from context
		person, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogPersonAuthNotFoundInContext)
			_ = c.Error(domain.ErrUserNotFound)
			return
		}

		log.Info(logger.LogPersonDeleteSelfRequest, "user_id", person.ID, "email", person.Email)

		// STEP 1: Check if user has active branches
		branches, err := h.BranchInteractor.GetBranchesByRepresentative(c, person.ID)
		if err != nil {
			log.Error(logger.LogPersonCheckBranchesError, "error", err, "user_id", person.ID)
			h.Response.Error(c, domain.MsgPersonCannotDelete)
			return
		}

		if len(branches) > 0 {
			log.Warn(logger.LogPersonHasActiveBranches, "user_id", person.ID, "branch_count", len(branches))
			h.Response.Error(c, domain.MsgPersonHasBranches)
			return
		}

		// STEP 2: Delete from Keycloak first
		if err := h.Interactor.DeleteKeycloakUser(c, person.KeycloakUserID); err != nil {
			log.Error(logger.LogPersonDeleteKeycloakFailed, "error", err, "keycloak_id", person.KeycloakUserID)
			h.Response.Error(c, domain.MsgPersonCannotDelete)
			return
		}

		// STEP 3: Delete from database
		if err := h.Interactor.DeletePersonFromDB(c, person.ID); err != nil {
			log.Error(logger.LogPersonDeleteDBFailed, "error", err, "user_id", person.ID)
			// Note: User already deleted from Keycloak - inconsistent state
			h.Response.Error(c, domain.MsgPersonCannotDelete)
			return
		}

		log.Success(logger.LogPersonControllerAccountDeleteOK, "user_id", person.ID, "email", person.Email)
		h.Response.Success(c, domain.MsgPersonDeleted)
	}
}
