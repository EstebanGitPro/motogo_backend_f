package ports

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/Nerzal/gocloak/v13"
)

// RegistrationResult contiene el resultado del registro con token de Keycloak
type RegistrationResult struct {
	Person domain.Person
	Token  *gocloak.JWT
}

type Service interface {
	RegisterPerson(person domain.Person) (*RegistrationResult, error)
	Login(email, password string) (*gocloak.JWT, error)
	GetPersonByEmail(email string) (*domain.Person, error)
}