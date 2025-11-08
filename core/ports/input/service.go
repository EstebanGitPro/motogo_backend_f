package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)


type Service interface {
	RegisterPerson(ctx context.Context,person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context,email string) (*domain.Person, error)
}


