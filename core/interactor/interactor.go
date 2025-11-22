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
	i.logger.Info("Iniciando proceso de registro",
		"email", person.Email,
		"role", person.Role)

	// PASO 1: Validaciones iniciales
	result, err = i.service.RegisterPerson(ctx, person)
	if err != nil {
		i.logger.Error("[PASO 1/7] Validaciones fallidas", "error", err)
		return
	}
	i.logger.Success("[PASO 1/7] Validaciones completadas")

	person.SetID()
	i.logger.Debug("ID generado para persona", "person_id", person.ID)

	// PASO 2: Iniciar transacción
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error("[PASO 2/7] Error iniciando transacción", "error", err)
		return
	}
	i.logger.Success("[PASO 2/7] Transacción iniciada")

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error("ROLLBACK BD FALLÓ - ALERTA CRÍTICA",
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn("Rollback BD ejecutado correctamente")
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					i.logger.Error("ROLLBACK KEYCLOAK FALLÓ - ALERTA CRÍTICA",
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
				} else {
					i.logger.Warn("Rollback Keycloak ejecutado correctamente")
				}
			}
		}
	}()

	// PASO 3: Guardarza persona en BD
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		i.logger.Error("[PASO 3/7] Error guardando persona", "error", err)
		return
	}
	i.logger.Success("[PASO 3/7] Persona guardada en BD")

	// PASO 4: Crear usuario en Keycloak
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		i.logger.Error("[PASO 4/7] Error creando usuario en Keycloak", "error", err)
		return
	}
	keycloakCreated = true // Marcar para compensación en defer
	i.logger.Success("[PASO 4/7] Usuario creado en Keycloak", "keycloak_user_id", keycloakUserID)

	if err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password); err != nil {
		i.logger.Error("[PASO 5/7] Error configurando password", "error", err)
		return
	}
	i.logger.Success("[PASO 5/7] Password configurado")

	if err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role); err != nil {
		i.logger.Error("[PASO 6/7] Error asignando rol", "error", err)
		return
	}
	i.logger.Success("[PASO 6/7] Rol asignado", "role", person.Role)

	// PASO 7: Actualizar BD con keycloak_user_id
	if err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID); err != nil {
		i.logger.Error("[PASO 7/7] Error actualizando Keycloak ID en BD", "error", err)
		return
	}
	i.logger.Success("[PASO 7/7] Keycloak_user_id actualizado en BD")

	// COMMIT: Confirmar toda la transacción
	if err = tx.Commit(); err != nil {
		i.logger.Error("COMMIT FALLÓ - ALERTA CRÍTICA", "error", err)
		return
	}
	i.logger.Success("Transacción confirmada exitosamente")

	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	i.logger.Success("Registro completado exitosamente",
		"email", person.Email,
		"person_id", person.ID,
		"keycloak_user_id", keycloakUserID)

	err = nil //asegurar que defer NO ejecute rollback
	return
}
