package services

import (
	"context"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

type service struct {
	repository     output.Repository
	keycloak       output.AuthClient
}

func NewService(repository output.Repository, keycloak output.AuthClient ) input.Service {
	return &service{
		repository:     repository,
		keycloak:       keycloak,
	}			
}


func (s service) RegisterPerson(ctx context.Context,person domain.Person) (*dto.RegistrationResult, error) {
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

	keycloakUserID, err := s.keycloak.CreateUser(ctx,&person)

		if err != nil {

		existingUser, getErr := s.keycloak.GetUserByEmail(ctx,person.Email)
		if getErr != nil {
			return nil, fmt.Errorf("failed to create or get user in keycloak: %w", err)
		}
		keycloakUserID = *existingUser.ID
		}
	
		

		person.KeycloakUserID = keycloakUserID
		err = s.repository.UpdatePerson(person)
		if err != nil {
			_ = s.keycloak.DeleteUser(ctx,keycloakUserID)
			return nil,fmt.Errorf("failed to update person with keycloak user id: %w", err)
		}

	// commit transaction
	if err != nil {
		return nil,err
	}



		
		return nil, domain.ErrRegistrationFailed
	}



func (s service) GetPersonByEmail(ctx context.Context,email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}