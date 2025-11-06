package services

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

type service struct {
	repository     output.Repository
	config         *config.Config
}

func NewService(repo output.Repository, cfg *config.Config) input.Service {
	return &service{
		repository:     repo,
		config:         cfg,
	}			
}


func (s service) RegisterPerson(person domain.Person) (*dto.RegistrationResult, error) {
	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return nil, domain.ErrDuplicateUser
	}

	if person.Role == "" {
		return nil, domain.ErrRoleRequired
	}

	person.SetID()
	
	
	err = s.repository.SavePerson(person)
	if err != nil {
		return nil, err
	}


		
		return nil, domain.ErrRegistrationFailed
	}



func (s service) GetPersonByEmail(email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}