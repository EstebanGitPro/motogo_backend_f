package ports

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/dto"
	"github.com/Nerzal/gocloak/v13"
)


type Service interface {
	RegisterPerson(person domain.Person) (*dto.RegistrationResult, error)
	LoginPerson(email, password string) (*gocloak.JWT, error)
	GetPersonByEmail(email string) (*domain.Person, error)
}