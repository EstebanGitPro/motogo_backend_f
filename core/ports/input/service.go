package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

// Service - Use Cases atómicos que el Interactor orquesta (Person)
type Service interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Person - Validaciones y consultas
	RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)

	// Person - Operaciones transaccionales de BD
	SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error
	UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error

	// Person - Operaciones de Keycloak
	CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error

	// Person - Compensaciones (rollback)
	RollbackPerson(ctx context.Context, personID string) error
	RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error
}

// MessageService - Use Cases para Messages (separado de Person)
type MessageService interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Messages - Validaciones y operaciones
	ValidateMessage(ctx context.Context, message domain.Message) error
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessageByCode(ctx context.Context, code string) (*domain.Message, error)
	ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error)
	ListActiveMessages(ctx context.Context) ([]domain.Message, error)

	// Messages - Operaciones transaccionales de BD
	SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error
	UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error
	DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error
}
