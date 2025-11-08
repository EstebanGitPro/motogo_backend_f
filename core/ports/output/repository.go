package output

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

type Repository interface {
	SavePerson(person domain.Person) error
	GetPersonByEmail(email string) (*domain.Person, error)
	GetPersonByID(id string) (*domain.Person, error)
	UpdatePerson(person domain.Person) error
	PatchPerson(id string, keycloakUserID string) error
	DeletePerson(id string) error
}