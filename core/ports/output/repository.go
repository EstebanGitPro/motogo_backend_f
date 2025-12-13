package output

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Tx interface for database transactions
type Tx interface {
	Commit() error
	Rollback() error
}

// Repository interface for Person operations
type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Person operations - transactional
	SavePerson(ctx context.Context, tx Tx, person domain.Person) error
	UpdatePerson(ctx context.Context, tx Tx, person domain.Person) error
	PatchPerson(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeletePerson(ctx context.Context, tx Tx, id string) error

	// Person operations - read
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
}

// MessageRepository interface for Message operations
type MessageRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Message operations - transactional
	SaveMessage(ctx context.Context, tx Tx, message domain.Message) error
	UpdateMessage(ctx context.Context, tx Tx, message domain.Message) error
	DeleteMessage(ctx context.Context, tx Tx, id string) error

	// Message operations - read
	GetAllActive(ctx context.Context) ([]domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	GetByCode(ctx context.Context, code string) (*domain.Message, error)
	GetByType(ctx context.Context, msgType string) ([]domain.Message, error)
	GetByModule(ctx context.Context, module string) ([]domain.Message, error)

}
