package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
}

func NewService(repository output.Repository, keycloak output.AuthClient) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
	}
}

func (s service) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}

func (s service) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	return s.repository.GetPersonByID(id)
}

func (s service) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return nil, domain.ErrDuplicateUser
	}

	return &dto.RegistrationResult{
		Person:  person,
		Message: "Validaciones exitosas",
	}, nil
}


func (s service) SavePersonToDB(ctx context.Context, person domain.Person) error {
	return s.repository.SavePerson(person)
}


func (s service) CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	return s.keycloak.CreateUser(ctx, person)
}


func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	return s.keycloak.SetPassword(ctx, userID, password, true)
}


func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	return s.keycloak.AssignRole(ctx, userID, role)
}

func (s service) UpdatePersonKeycloakID(ctx context.Context, personID string, keycloakUserID string) error {
	return s.repository.PatchPerson(personID, keycloakUserID)
}


func (s service) RollbackPerson(ctx context.Context, personID string) error {
	return s.repository.DeletePerson(personID)
}

func (s service) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	return s.keycloak.DeleteUser(ctx, keycloakUserID)
}

