package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
)

type Interactor struct {
	service input.Service
}

func NewInteractor(service input.Service) *Interactor {
	return &Interactor{
		service: service,
	}
}

func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	// 1. Validaciones iniciales (email duplicado)
	result, err := i.service.RegisterPerson(ctx, person)
	if err != nil {
		return nil, err
	}

	person.SetID()

	// 2. Iniciar transacción de BD
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Guardar persona en BD dentro de la transacción
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 4. Crear usuario en Keycloak
	keycloakUserID, err := i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 5. Configurar password en Keycloak
	err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 6. Asignar rol en Keycloak
	err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 7. Actualizar keycloakID en BD dentro de la transacción
	err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 8. Confirmar transacción de BD (commit final)
	if err = tx.Commit(); err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		return nil, err
	}

	// 9. Retornar resultado exitoso
	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	return result, nil
}
