package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/dto"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/Nerzal/gocloak/v13"
)

type service struct {
	repository     ports.Repository
	authzService   ports.AuthorizationService
	config         *config.Config
}

func NewService(repo ports.Repository, authzService ports.AuthorizationService, cfg *config.Config) ports.Service {
	return &service{
		repository:     repo,
		authzService:   authzService,
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
	
	plainPassword := person.Password
	
	if err := person.HashPassword(); err != nil {
		return nil, err
	}

	err = s.repository.SavePerson(person)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = s.syncUserWithKeycloak(ctx, &person, plainPassword)
	if err != nil {
		slog.Error("Failed to sync user with Keycloak, rolling back",
			"user_id", person.ID,
			"email", person.Email,
			"error", err)
		
		if deleteErr := s.repository.DeletePerson(person.ID); deleteErr != nil {
			slog.Error("Failed to rollback user creation",
				"user_id", person.ID,
				"error", deleteErr)
		}
		
		return nil, domain.ErrRegistrationFailed
	}

	slog.Info("User registered successfully",
		"user_id", person.ID,
		"email", person.Email,
		"role", person.Role)

	return &dto.RegistrationResult{
		Person:  person,
		Message: "Usuario registrado exitosamente. Por favor, inicie sesión para continuar.",
	}, nil
}


func (s service) syncUserWithKeycloak(ctx context.Context, person *domain.Person, plainPassword string) error {
	if s.authzService == nil {
		return fmt.Errorf("keycloak service not configured")
	}

	keycloakUserID, err := s.authzService.SyncUserToKeycloak(ctx, person)
	if err != nil {
		return fmt.Errorf("failed to sync user: %w", err)
	}

	err = s.authzService.SetUserPassword(ctx, keycloakUserID, plainPassword)
	if err != nil {
		_ = s.authzService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return fmt.Errorf("failed to set password: %w", err)
	}
	
	err = s.authzService.AssignRole(ctx, person.ID, person.Role)
	if err != nil {
		_ = s.authzService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	slog.Info("User synced with Keycloak successfully", 
		"user_id", person.ID, 
		"email", person.Email,
		"role", person.Role,
		"keycloak_user_id", keycloakUserID)
	
	return nil
}


func (s service) LoginPerson(email, password string) (*gocloak.JWT, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	person, err := s.repository.GetPersonByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	ctx := context.Background()
	token, err := s.authzService.LoginUser(ctx, person.Email, password)
	if err != nil {
		slog.Warn("Login failed",
			"email", email,
			"error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	slog.Info("User logged in successfully",
		"user_id", person.ID,
		"email", person.Email)

	return token, nil
}

func (s service) GetPersonByEmail(email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}