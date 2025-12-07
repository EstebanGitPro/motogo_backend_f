package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
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
	i.logger.Info(logger.LogPersonInteractorRegStart, person.ToLogger())

	// PASO 1: Validaciones iniciales
	result, err = i.service.RegisterPerson(ctx, person)
	if err != nil {
		// Si es un registro incompleto, ejecutar limpieza automática
		if err == domain.ErrIncompleteRegistration {
			i.logger.Warn(logger.LogPersonInteractorIncompleteDetected, "email", person.Email)

			// Intentar limpiar el estado inconsistente
			if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, person.Email); cleanErr != nil {
				i.logger.Error(logger.LogPersonInteractorCleanup_Error, "email", person.Email, "error", cleanErr)
				// Si falla la limpieza, retornar error de limpieza
				return nil, cleanErr
			}

			i.logger.Success(logger.LogPersonInteractorCleanup_OK, "email", person.Email)
			// Retornar el error de registro incompleto para que el cliente sepa que debe reintentar
			return nil, err
		}

		i.logger.Error(logger.LogPersonInteractorStep1_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep1_OK)

	person.SetID()
	i.logger.Debug(logger.LogPersonInteractorIDGenerated, "person_id", person.ID)

	// PASO 1.5: Verificar estado consistente (ya no debería haber inconsistencias)
	if err = i.service.CheckAndCleanInconsistentState(ctx, person.Email); err != nil {
		i.logger.Error(logger.LogPersonInteractorStep15_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep15_OK)

	// PASO 2: Iniciar transacción
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error(logger.LogPersonInteractorStep2_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep2_OK)

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error(logger.LogPersonInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn(logger.LogPersonInteractorRollbackDB_OK)
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					i.logger.Error(logger.LogPersonInteractorRollbackKeycloak_Err,
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
				} else {
					i.logger.Warn(logger.LogPersonInteractorRollbackKeycloak_OK)
				}
			}
		}
	}()

	// PASO 3: Guardar persona en BD
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		i.logger.Error(logger.LogPersonInteractorStep3_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep3_OK)

	// PASO 4: Crear usuario en Keycloak
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		i.logger.Error(logger.LogPersonInteractorStep4_Error, "error", err)
		// Wrap error to indicate Keycloak creation failure
		err = domain.ErrKeycloakUserCreationFailed
		return
	}
	keycloakCreated = true // Marcar para compensación en defer
	i.logger.Success(logger.LogPersonInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	if err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password); err != nil {
		i.logger.Error(logger.LogPersonInteractorStep5_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep5_OK)

	if err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role); err != nil {
		i.logger.Error(logger.LogPersonInteractorStep6_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep6_OK, "role", person.Role)

	// PASO 7: Actualizar BD con keycloak_user_id
	if err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID); err != nil {
		i.logger.Error(logger.LogPersonInteractorStep7_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorStep7_OK)

	// PASO 8: Confirmar toda la transacción
	if err = tx.Commit(); err != nil {
		i.logger.Error(logger.LogPersonInteractorCommit_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogPersonInteractorCommit_OK)

	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	//TODO: preguntar si dejar info en el logger success
	i.logger.Success(logger.LogPersonInteractorRegComplete,
		person.ToLogger(),
		"keycloak_user_id", keycloakUserID)

	err = nil //asegurar que defer NO ejecute rollback
	return
}
