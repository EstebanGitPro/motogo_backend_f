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
	SavePerson(ctx context.Context, tx Tx, person domain.Person) error
	GetPersonByEmail(ctx context.Context, tx Tx, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, tx Tx, id string) (*domain.Person, error)
	UpdatePerson(ctx context.Context, tx Tx, person domain.Person) error
	PatchPerson(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeletePerson(ctx context.Context, tx Tx, id string) error
}
