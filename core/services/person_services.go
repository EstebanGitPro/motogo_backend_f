package services

import (
	"context"
	"log/slog"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/EstebanGitPro/motogo-backend/config"
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


func (s service) RegisterPerson(person domain.Person) (domain.Person, error) {
	// 1. Validar email duplicado
	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return domain.Person{}, domain.ErrDuplicateUser
	}

	// 2. Generar ID y hashear contraseña
	person.SetID()
	if err := person.HashPassword(); err != nil {
		return domain.Person{}, err
	}

	// 3. Guardar en tu base de datos
	err = s.repository.Save(person)
	if err != nil {
		return domain.Person{}, err
	}

	// 4. Si tiene rol asignado, sincronizar con Keycloak
	if person.Role != "" && person.Role != "user" {
		ctx := context.Background()
		err = s.syncUserWithKeycloak(ctx, &person)
		if err != nil {
			// Log el error pero no fallar el registro
			slog.Warn("Failed to sync user with Keycloak", 
				"user_id", person.ID, 
				"email", person.Email,
				"role", person.Role,
				"error", err)
		}
	}

	return person, nil
}

// syncUserWithKeycloak sincroniza el usuario con Keycloak y asigna el rol
func (s service) syncUserWithKeycloak(ctx context.Context, person *domain.Person) error {
	if s.authzService == nil {
		return nil // Keycloak no configurado, continuar sin error
	}

	// Sincronizar usuario y asignar rol
	err := s.authzService.AssignRole(ctx, person.ID, person.Role)
	if err != nil {
		return err
	}

	slog.Info("User synced with Keycloak successfully", 
		"user_id", person.ID, 
		"email", person.Email,
		"role", person.Role)
	
	return nil
}

func (s service) GetPersonByEmail(email string) (*domain.Person, error) {
	return s.repository.GetPersonByEmail(email)
}