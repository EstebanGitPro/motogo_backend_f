package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)


// Service - Use Cases atómicos que el Interactor orquesta
type Service interface {
	// Validaciones
	RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	
	// Transacciones atómicas
	SavePersonToDB(ctx context.Context, person domain.Person) error
	CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	UpdatePersonKeycloakID(ctx context.Context, personID string, keycloakUserID string) error
	
	// Compensaciones (rollback)
	RollbackPerson(ctx context.Context, personID string) error
	RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error
}
