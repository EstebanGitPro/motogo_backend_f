package ports

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/dto"
)


type Service interface {
	RegisterPerson(person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(email string) (*domain.Person, error)
}