package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/dto"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
)

type service struct {
	repository ports.Repository
	config     *config.Config
}

func NewService(repo ports.Repository, cfg *config.Config) ports.Service {
	return &service{
		repository: repo,
		config:     cfg,
	}
}

func (s service) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	existingPerson, err := s.repository.GetPersonByEmail(ctx, person.Email)
	if err == nil && existingPerson != nil {
		return nil, domain.ErrDuplicateUser
	}

	if person.Role == "" {
		return nil, domain.ErrRoleRequired
	}

	person.SetID()

	err = s.repository.SavePerson(ctx, &person)
	if err != nil {
		return nil, err
	}

	return nil, domain.ErrRegistrationFailed
}

func (s service) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(ctx, email)
}
