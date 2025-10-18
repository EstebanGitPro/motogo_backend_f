package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/EstebanGitPro/motogo-backend/config"
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


func (s service) RegisterPerson(person domain.Person) (*ports.RegistrationResult, error) {
	// 1. Validar email duplicado
	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return nil, domain.ErrDuplicateUser
	}

	// 2. Validar que el rol venga del frontend
	if person.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	// 3. Generar ID
	person.SetID()
	
	// 4. Guardar contraseña sin hashear temporalmente para Keycloak
	plainPassword := person.Password
	
	// 5. Hashear contraseña para base de datos local
	if err := person.HashPassword(); err != nil {
		return nil, err
	}

	// 6. Guardar en base de datos local PRIMERO
	err = s.repository.Save(person)
	if err != nil {
		return nil, err
	}

	// 7. SIEMPRE sincronizar con Keycloak (CRÍTICO para autenticación)
	ctx := context.Background()
	token, err := s.syncUserWithKeycloak(ctx, &person, plainPassword)
	if err != nil {
		// ROLLBACK: Eliminar usuario de BD local si Keycloak falla
		slog.Error("Failed to sync user with Keycloak, rolling back",
			"user_id", person.ID,
			"email", person.Email,
			"error", err)
		
		// Intentar rollback
		if deleteErr := s.repository.Delete(person.ID); deleteErr != nil {
			slog.Error("Failed to rollback user creation",
				"user_id", person.ID,
				"error", deleteErr)
		}
		
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	slog.Info("User registered successfully",
		"user_id", person.ID,
		"email", person.Email,
		"role", person.Role)

	return &ports.RegistrationResult{
		Person: person,
		Token:  token,
	}, nil
}

// syncUserWithKeycloak sincroniza el usuario con Keycloak, configura contraseña, asigna rol y devuelve token
func (s service) syncUserWithKeycloak(ctx context.Context, person *domain.Person, plainPassword string) (*gocloak.JWT, error) {
	if s.authzService == nil {
		return nil, fmt.Errorf("keycloak service not configured")
	}

	// 1. Sincronizar usuario con Keycloak (crea el usuario)
	keycloakUserID, err := s.authzService.SyncUserToKeycloak(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user: %w", err)
	}

	// 2. Configurar contraseña en Keycloak para autenticación
	err = s.authzService.SetUserPassword(ctx, keycloakUserID, plainPassword)
	if err != nil {
		// Intentar eliminar usuario de Keycloak si falla la contraseña
		_ = s.authzService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	// 3. Asignar rol en Keycloak para autorización
	err = s.authzService.AssignRole(ctx, person.ID, person.Role)
	if err != nil {
		// Intentar eliminar usuario de Keycloak si falla el rol
		_ = s.authzService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	// 4. Autenticar al usuario y obtener su token JWT
	token, err := s.authzService.LoginUser(ctx, person.Email, plainPassword)
	if err != nil {
		slog.Warn("User created but login failed", 
			"user_id", person.ID,
			"error", err)
		// No fallar aquí, el usuario puede iniciar sesión después
		return nil, fmt.Errorf("user created but failed to generate token: %w", err)
	}

	slog.Info("User synced with Keycloak successfully", 
		"user_id", person.ID, 
		"email", person.Email,
		"role", person.Role,
		"keycloak_user_id", keycloakUserID)
	
	return token, nil
}

// Login autentica un usuario y devuelve su token JWT de Keycloak
func (s service) Login(email, password string) (*gocloak.JWT, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Verificar que el usuario existe en la base de datos local
	person, err := s.repository.GetPersonByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Autenticar con Keycloak
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