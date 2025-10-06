package ports

import "github.com/EstebanGitPro/motogo-backend/core/domain"

type Repository interface {
	Save(person domain.Person) error
	GetPersonByEmail(email string) (*domain.Person, error)
}