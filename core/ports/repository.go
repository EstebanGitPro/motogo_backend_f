package ports

import "github.com/EstebanGitPro/motogo-backend/core/domain"

type Repository interface {
	SavePerson(person domain.Person) error
	GetPersonByEmail(email string) (*domain.Person, error)
	GetPersonByID(id string) (*domain.Person, error)
	UpdatePerson(person domain.Person) error
	DeletePerson(id string) error
}