package ports

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/dto"
)


type Service interface {
	RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
}