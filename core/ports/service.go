package ports

import "github.com/EstebanGitPro/motogo-backend/core/domain"

type Service interface {
	RegisterPerson(person domain.Person) (domain.Person, error)
	GetPersonByEmail(email string) (*domain.Person, error)
}