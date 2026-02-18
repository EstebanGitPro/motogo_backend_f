package interactor

import (
	"context"
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

type Interactor struct {
	service input.Service
}

func NewInteractor(service input.Service) *Interactor {
	return &Interactor{
		service: service,
	}
}

func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (result *dto.RegistrationResult, err error) {
	// Extract traceID from context and create logger with it
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonInteractorRegStart, person.ToLogger())

	// PASO 1: Validaciones iniciales
	result, err = i.service.RegisterPerson(ctx, person)
	if err != nil {
		return nil, i.handleIncompleteRegistration(ctx, err, person)
	}
	log.Success(logger.LogPersonInteractorStep1_OK)

	person.SetID()
	log.Debug(logger.LogPersonInteractorIDGenerated, "person_id", person.ID)

	// PASO 1.5: Verificar estado consistente (ya no debería haber inconsistencies)
	if err = i.service.CheckAndCleanInconsistentState(ctx, person.Email); err != nil {
		log.Error(logger.LogPersonInteractorStep15_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorStep15_OK)

	// PASO 2: Iniciar transacción
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogPersonInteractorStep2_Error, "error", err)
		return result, err
	}

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		i.rollbackRegistration(tx, keycloakUserID, keycloakCreated, &err, ctx)
	}()

	// PASO 3: Guardar persona en BD
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		log.Error(logger.LogPersonInteractorStep3_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorStep3_OK)

	// PASO 4: Crear usuario en Keycloak
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		log.Error(logger.LogPersonInteractorStep4_Error, "error", err)
		// Wrap error to indicate Keycloak creation failure
		err = domain.ErrKeycloakUserCreationFailed
		return result, err
	}
	keycloakCreated = true // Marcar para compensación en defer
	log.Success(logger.LogPersonInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	if err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password); err != nil {
		log.Error(logger.LogPersonInteractorStep5_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorStep5_OK)

	if err = i.service.AssignUserRole(ctx, keycloakUserID, string(person.Role)); err != nil {
		log.Error(logger.LogPersonInteractorStep6_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorStep6_OK, "role", string(person.Role))

	// PASO 7: Actualizar BD con keycloak_user_id
	if err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID); err != nil {
		log.Error(logger.LogPersonInteractorStep7_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorStep7_OK)

	// PASO 8: Confirmar toda la transacción
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogPersonInteractorCommit_Error, "error", err)
		return result, err
	}
	log.Success(logger.LogPersonInteractorCommit_OK)

	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	// PASO 9: Enviar email de verificación (no bloquea el registro si falla)
	i.sendVerificationEmailBestEffort(ctx, keycloakUserID, person.Email)

	log.Success(logger.LogPersonInteractorRegComplete,
		person.ToLogger(),
		"keycloak_user_id", keycloakUserID)

	err = nil // asegurar que defer NO ejecute rollback
	return result, err
}

// handleIncompleteRegistration handles the error case when initial registration validation fails.
func (i *Interactor) handleIncompleteRegistration(ctx context.Context, err error, person domain.Person) error {
	if !errors.Is(err, domain.ErrIncompleteRegistration) {
		log.Error(logger.LogPersonInteractorStep1_Error, "error", err)
		return err
	}

	log.Warn(logger.LogPersonInteractorIncompleteDetected, "email", person.Email)

	// Intentar limpiar el estado inconsistente
	if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, person.Email); cleanErr != nil {
		log.Error(logger.LogPersonInteractorCleanup_Error, "email", person.Email, "error", cleanErr)
		return cleanErr
	}

	log.Success(logger.LogPersonInteractorCleanup_OK, "email", person.Email)
	return err
}

// rollbackRegistration handles rollback of both DB and Keycloak on error.
func (i *Interactor) rollbackRegistration(tx output.Tx, keycloakUserID string, keycloakCreated bool, err *error, ctx context.Context) {
	if *err == nil {
		return
	}
	if rbErr := tx.Rollback(); rbErr != nil {
		log.Error(logger.LogPersonInteractorRollbackDB_Error,
			"rollback_error", rbErr,
			"original_error", *err)
	} else {
		log.Warn(logger.LogPersonInteractorRollbackDB_OK)
	}

	if keycloakCreated {
		if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
			log.Error(logger.LogPersonInteractorRollbackKeycloak_Err,
				"keycloak_error", kcErr,
				"keycloak_user_id", keycloakUserID)
		} else {
			log.Warn(logger.LogPersonInteractorRollbackKeycloak_OK)
		}
	}
}

// sendVerificationEmailBestEffort sends a verification email without failing the registration.
func (i *Interactor) sendVerificationEmailBestEffort(ctx context.Context, keycloakUserID, email string) {
	if sendErr := i.service.SendVerificationEmail(ctx, keycloakUserID); sendErr != nil {
		log.Warn(logger.LogKeycloakSendVerificationEmailError,
			"keycloak_user_id", keycloakUserID,
			"email", email,
			"error", sendErr)
	} else {
		log.Info(logger.LogKeycloakSendVerificationEmailOK,
			"keycloak_user_id", keycloakUserID,
			"email", email)
	}
}

// ResendVerificationEmail reenvía el email de verificación a un usuario
func (i *Interactor) ResendVerificationEmail(ctx context.Context, email string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogKeycloakSendVerificationEmail, "email", email)

	// Buscar usuario en Keycloak por email
	user, err := i.service.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return domain.ErrUserNotFound
	}

	// Verificar si el email ya está verificado
	if user.EmailVerified != nil && *user.EmailVerified {
		log.Warn(logger.LogKeycloakSendVerificationEmailError, "email", email, "reason", "email already verified")
		return domain.ErrEmailAlreadyVerified
	}

	// Enviar email de verificación
	if err = i.service.SendVerificationEmail(ctx, *user.ID); err != nil {
		log.Error(logger.LogKeycloakSendVerificationEmailError, "email", email, "error", err)
		return err
	}

	log.Success(logger.LogKeycloakSendVerificationEmailOK, "email", email, "user_id", *user.ID)
	return nil
}

// RequestPasswordReset envía un email de recuperación de contraseña
func (i *Interactor) RequestPasswordReset(ctx context.Context, email string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogKeycloakSendPasswordReset, "email", email)

	// Enviar email de reset (internamente busca el usuario y envía)
	// NOTA: Por seguridad, no revelamos si el email existe o no
	if err := i.service.SendPasswordResetEmail(ctx, email); err != nil {
		// Log el error pero responder con éxito genérico al cliente
		log.Warn(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
	} else {
		log.Success(logger.LogKeycloakSendPasswordResetOK, "email", email)
	}

	// Siempre retornar nil por seguridad (no revelar si el usuario existe)
	return nil
}

// VerifyEmailByToken verifica el email de un usuario extrayéndolo del token JWT
// Este método delega al Service que maneja la lógica de negocio (parsing del token y verificación en Keycloak)
func (i *Interactor) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogKeycloakEmailVerify)

	// Delegar toda la lógica al Service (parsing del token + verificación en Keycloak)
	email, err := i.service.VerifyEmailByToken(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			log.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "invalid token")
		case errors.Is(err, domain.ErrUserNotFound):
			log.Warn(logger.LogKeycloakUserNotFound, "email", email)
		case errors.Is(err, domain.ErrEmailAlreadyVerified):
			log.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email)
		default:
			log.Error(logger.LogKeycloakEmailVerifyError, "email", email, "error", err)
		}
		return email, err
	}

	log.Success(logger.LogKeycloakEmailVerifyOK, "email", email)
	return email, nil
}

// ResetPasswordWithToken actualiza la contraseña de un usuario extrayendo el email del token JWT
// Este método delega al Service que maneja la lógica de negocio (parsing del token y actualización en Keycloak)
func (i *Interactor) ResetPasswordWithToken(ctx context.Context, token string, newPassword string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPasswordResetStart)

	// Delegar toda la lógica al Service (parsing del token + actualización de contraseña en Keycloak)
	err := i.service.ResetPasswordWithToken(ctx, token, newPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			log.Error(logger.LogPasswordResetTokenError, "error", err)
		case errors.Is(err, domain.ErrUserNotFound):
			log.Error(logger.LogPasswordResetUserNotFound, "error", err)
		case errors.Is(err, domain.ErrPasswordUpdateFailed):
			log.Error(logger.LogPasswordResetUpdateError, "error", err)
		default:
			log.Error(logger.LogPasswordResetUpdateError, "error", err)
		}
		return err
	}

	log.Success(logger.LogPasswordResetSuccess)
	return nil
}

func (i *Interactor) Login(ctx context.Context, email, password string) (*dto.TokenResponse, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonInteractorLoginStart, "email", email)

	// Delegar al servicio la autenticación en Keycloak
	token, err := i.service.Login(ctx, email, password)
	if err != nil {
		log.Error(logger.LogPersonInteractorLoginError, "email", email, "error", err)
		return nil, err
	}

	log.Success(logger.LogPersonInteractorLoginOK, "email", email, "token", token)
	return &dto.TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		RefreshToken: token.RefreshToken,
	}, nil
}

// RefreshToken obtains a new access token using the refresh token
// This is called by the frontend when the access token expires
func (i *Interactor) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonRefreshTokenStart)

	// Delegate to service to refresh token via Keycloak
	token, err := i.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.Error(logger.LogPersonRefreshTokenError, "error", err)
		return nil, err
	}

	log.Success(logger.LogPersonInteractorRefreshOK)
	return &dto.TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		RefreshToken: token.RefreshToken,
	}, nil
}

// ChangePassword allows an authenticated user to change their password (HU57)
// Requires the current password for verification before setting a new one
func (i *Interactor) ChangePassword(ctx context.Context, keycloakUserID, currentPassword, newPassword string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogChangePasswordStart, "keycloak_user_id", keycloakUserID)

	// Delegate to service
	err := i.service.ChangePassword(ctx, keycloakUserID, currentPassword, newPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			log.Error(logger.LogChangePasswordUserNotFound, "error", err)
		case errors.Is(err, domain.ErrInvalidCredentials):
			log.Warn(logger.LogChangePasswordInvalidCurrent, "error", err)
		case errors.Is(err, domain.ErrPasswordUpdateFailed):
			log.Error(logger.LogChangePasswordUpdateError, "error", err)
		case errors.Is(err, domain.ErrKeycloakUnavailable):
			log.Error(logger.LogKeycloakUnavailable, "error", err)
		default:
			log.Error(logger.LogChangePasswordUpdateError, "error", err)
		}
		return err
	}

	log.Success(logger.LogChangePasswordSuccess, "keycloak_user_id", keycloakUserID)
	return nil
}

// UpdateProfile updates the authenticated user's profile (HU52)
// This method orchestrates the profile update with proper transaction management
func (i *Interactor) UpdateProfile(ctx context.Context, person domain.Person) (*domain.Person, error) {
	// Extract traceID from context and create logger with it
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogUpdateProfileStart, person.ToLogger())

	// STEP 1: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogPersonInteractorStep2_Error, "error", err)
		return nil, err
	}

	var profileUpdated bool

	defer func() {
		if err != nil && !profileUpdated {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogPersonInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogPersonInteractorRollbackDB_OK)
			}
		}
	}()

	// STEP 2: Update profile via service
	if err = i.service.UpdatePersonProfile(ctx, tx, person); err != nil {
		log.Error(logger.LogUpdateProfileError, "error", err, "person_id", person.ID)
		return nil, err
	}

	// STEP 3: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogPersonInteractorCommit_Error, "error", err)
		return nil, err
	}
	profileUpdated = true

	log.Success(logger.LogUpdateProfileSuccess, person.ToLogger())
	return &person, nil
}

// GetPublicContact retrieves public contact info for a person (HU55)
// Only returns phone_number for motorcyclists to contact representatives
func (i *Interactor) GetPublicContact(ctx context.Context, personID string) (*domain.Person, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonGetPublicContactStart, "person_id", personID)

	person, err := i.service.GetPersonByID(ctx, personID)
	if err != nil {
		log.Error(logger.LogPersonGetPublicContactError, "person_id", personID, "error", err)
		return nil, err
	}

	log.Success(logger.LogPersonInteractorContactGetOK, "person_id", personID)
	return person, nil
}

// DeleteKeycloakUser deletes a user from Keycloak (HU53)
// This is used as part of account deletion flow
func (i *Interactor) DeleteKeycloakUser(ctx context.Context, keycloakUserID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonDeleteKeycloakStart, "keycloak_user_id", keycloakUserID)

	// Use existing RollbackKeycloakUser which internally calls DeleteUser
	if err := i.service.RollbackKeycloakUser(ctx, keycloakUserID); err != nil {
		log.Error(logger.LogPersonDeleteKeycloakError, "error", err, "keycloak_user_id", keycloakUserID)
		return err
	}

	log.Success(logger.LogPersonInteractorKeycloakDeleteOK, "keycloak_user_id", keycloakUserID)
	return nil
}

// DeletePersonFromDB deletes a person from the database (HU53)
// This is used as part of account deletion flow
func (i *Interactor) DeletePersonFromDB(ctx context.Context, personID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogPersonDeleteDBStart, "person_id", personID)

	// Use existing RollbackPerson which internally calls DeletePerson
	if err := i.service.RollbackPerson(ctx, personID); err != nil {
		log.Error(logger.LogPersonDeleteDBError, "error", err, "person_id", personID)
		return err
	}

	log.Success(logger.LogPersonInteractorPersonDeleteOK, "person_id", personID)
	return nil
}
