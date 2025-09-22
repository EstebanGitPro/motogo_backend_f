package services

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
)

type service struct {
	repository     ports.Repository
}

func NewService(repo ports.Repository) ports.Service {
	return &service{
		repository:     repo,
	}
}


func (s service) RegisterPerson(person domain.Person) (domain.Person, error) {

	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return domain.Person{},domain.ErrDuplicateUser
	}

	person.SetID()

	if err := person.HashPassword(); err != nil {
		return domain.Person{}, err
	}

	err = s.repository.Save(person)
	if err != nil {
		return domain.Person{}, err
	}

	return person, nil
}

func (s service) GetPersonByEmail(email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}