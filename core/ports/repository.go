package ports

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
)

type Repository interface {
	SavePerson(ctx context.Context, person *domain.Person) error
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
	UpdatePerson(ctx context.Context, person *domain.Person) error
	DeletePerson(ctx context.Context, id string) error
}