package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

type Interactor struct {
	service input.Service
	logger  logger.Logger
}

func NewInteractor(service input.Service, log logger.Logger) *Interactor {
	return &Interactor{
		service: service,
		logger:  log,
	}
}

func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (result *dto.RegistrationResult, err error) {
	// Extract traceID from context and create logger with it
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogPersonInteractorRegStart, person.ToLogger())

	// PASO 1: Validaciones iniciales
	result, err = i.service.RegisterPerson(ctx, person)
	if err != nil {
		// Si es un registro incompleto, ejecutar limpieza automática
		if err == domain.ErrIncompleteRegistration {
			log.Warn(logger.LogPersonInteractorIncompleteDetected, "email", person.Email)

			// Intentar limpiar el estado inconsistente
			if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, person.Email); cleanErr != nil {
				log.Error(logger.LogPersonInteractorCleanup_Error, "email", person.Email, "error", cleanErr)
				// Si falla la limpieza, retornar error de limpieza
				return nil, cleanErr
			}

			log.Success(logger.LogPersonInteractorCleanup_OK, "email", person.Email)
			// Retornar el error de registro incompleto para que el cliente sepa que debe reintentar
			return nil, err
		}

		log.Error(logger.LogPersonInteractorStep1_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep1_OK)

	person.SetID()
	log.Debug(logger.LogPersonInteractorIDGenerated, "person_id", person.ID)

	// PASO 1.5: Verificar estado consistente (ya no debería haber inconsistencias)
	if err = i.service.CheckAndCleanInconsistentState(ctx, person.Email); err != nil {
		log.Error(logger.LogPersonInteractorStep15_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep15_OK)

	// PASO 2: Iniciar transacción
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogPersonInteractorStep2_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep2_OK)

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogPersonInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
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
	}()

	// PASO 3: Guardar persona en BD
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		log.Error(logger.LogPersonInteractorStep3_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep3_OK)

	// PASO 4: Crear usuario en Keycloak
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		log.Error(logger.LogPersonInteractorStep4_Error, "error", err)
		// Wrap error to indicate Keycloak creation failure
		err = domain.ErrKeycloakUserCreationFailed
		return
	}
	keycloakCreated = true // Marcar para compensación en defer
	log.Success(logger.LogPersonInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	if err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password); err != nil {
		log.Error(logger.LogPersonInteractorStep5_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep5_OK)

	if err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role); err != nil {
		log.Error(logger.LogPersonInteractorStep6_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep6_OK, "role", person.Role)

	// PASO 7: Actualizar BD con keycloak_user_id
	if err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID); err != nil {
		log.Error(logger.LogPersonInteractorStep7_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorStep7_OK)

	// PASO 8: Confirmar toda la transacción
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogPersonInteractorCommit_Error, "error", err)
		return
	}
	log.Success(logger.LogPersonInteractorCommit_OK)

	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	// PASO 9: Enviar email de verificación (no bloquea el registro si falla)
	if sendErr := i.service.SendVerificationEmail(ctx, keycloakUserID); sendErr != nil {
		// Log warning pero NO fallar el registro
		log.Warn(logger.LogKeycloakSendVerificationEmailError,
			"keycloak_user_id", keycloakUserID,
			"email", person.Email,
			"error", sendErr)
	} else {
		log.Info(logger.LogKeycloakSendVerificationEmailOK,
			"keycloak_user_id", keycloakUserID,
			"email", person.Email)
	}

	//TODO: preguntar si dejar info en el logger success
	log.Success(logger.LogPersonInteractorRegComplete,
		person.ToLogger(),
		"keycloak_user_id", keycloakUserID)

	err = nil //asegurar que defer NO ejecute rollback
	return
}

// ResendVerificationEmail reenvía el email de verificación a un usuario
func (i *Interactor) ResendVerificationEmail(ctx context.Context, email string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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
	log := i.logger.WithTraceID(traceID)

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
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakEmailVerify)

	// Delegar toda la lógica al Service (parsing del token + verificación en Keycloak)
	email, err := i.service.VerifyEmailByToken(ctx, token)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			log.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "invalid token")
		case domain.ErrUserNotFound:
			log.Warn(logger.LogKeycloakUserNotFound, "email", email)
		case domain.ErrEmailAlreadyVerified:
			log.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email)
		default:
			log.Error(logger.LogKeycloakEmailVerifyError, "email", email, "error", err)
		}
		return email, err
	}

	log.Success(logger.LogKeycloakEmailVerifyOK, "email", email)
	return email, nil
}

func (i *Interactor) Login(ctx context.Context, email, password string) (*dto.TokenResponse, error) {
    traceID := middleware.GetTraceIDFromContext(ctx)
    log := i.logger.WithTraceID(traceID)
    
    log.Info(logger.LogPersonInteractorLoginStart, "email", email)
    
    // Delegar al servicio la autenticación en Keycloak
    token, err := i.service.Login(ctx, email, password)
    if err != nil {
        log.Error(logger.LogPersonInteractorLoginError, "email", email, "error", err)
        return nil, err
    }
    
    log.Success(logger.LogPersonInteractorLoginOK, "email", email, "token", token)
    return &dto.TokenResponse{
        AccessToken: token.AccessToken,
        TokenType:   token.TokenType,
        ExpiresIn:   token.ExpiresIn,
        RefreshToken: token.RefreshToken,
    }, nil
}
