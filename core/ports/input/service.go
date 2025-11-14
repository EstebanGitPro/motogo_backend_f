package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

// Service - Use Cases atómicos que el Interactor orquesta
type Service interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Validaciones y consultas
	RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)

	// Operaciones transaccionales de BD
	SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error
	UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error

	// Operaciones de Keycloak
	CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error

	// Compensaciones (rollback)
	RollbackPerson(ctx context.Context, personID string) error
	RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error
}
