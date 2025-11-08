package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
)

// ========================================================================
// INTERACTOR - FACADE
// ========================================================================

// Interactor - Facade que orquesta múltiples use cases (servicios)
// Maneja la "transacción lógica" completa con commit/rollback
type Interactor struct {
	service input.Service
}

func NewInteractor(service input.Service) *Interactor {
	return &Interactor{
		service: service,
	}
}

// RegisterPerson - FACADE: Orquesta el proceso completo de registro
// Siguiendo el patrón del diagrama: Iniciar tx -> Ejecutar -> Confirmar/Cancelar
func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	var (
		personSaved    bool
		keycloakUserID string
	)

	result, err := i.service.RegisterPerson(ctx, person)
	if err != nil {
		return nil, err
	}

	person.SetID()

	
	err = i.service.SavePersonToDB(ctx, person)
	if err != nil {
		return nil, err
	}
	personSaved = true

	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		if personSaved {
			_ = i.service.RollbackPerson(ctx, person.ID)
		}
		return nil, err
	}


	err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackPerson(ctx, person.ID)
		return nil, err
	}

	
	err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackPerson(ctx, person.ID)
		return nil, err
	}

	err = i.service.UpdatePersonKeycloakID(ctx, person.ID, keycloakUserID)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackPerson(ctx, person.ID)
		return nil, err
	}

	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	return result, nil
}
