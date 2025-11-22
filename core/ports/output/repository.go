package output

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Operaciones transaccionales 
	SavePerson(ctx context.Context, tx Tx, person domain.Person) error
	UpdatePerson(ctx context.Context, tx Tx, person domain.Person) error
	PatchPerson(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeletePerson(ctx context.Context, tx Tx, id string) error

	// Operaciones de lectura 
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
}
