package input

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
)


type Service interface {
	RegisterPerson(person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(email string) (*domain.Person, error)
}


