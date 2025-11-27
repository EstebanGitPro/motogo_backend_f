package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
	logger     logger.Logger
}

func NewService(repository output.Repository, keycloak output.AuthClient, log logger.Logger) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
		logger:     log,
	}
}

func (s service) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	s.logger.Debug("Buscando persona por email", "email", email)
	person, err := s.repository.GetPersonByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Error buscando persona por email", "email", email, "error", err)
		return nil, err
	}
	s.logger.Debug("Persona encontrada por email", "email", email, "person_id", person.ID)
	return person, nil
}

func (s service) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	s.logger.Debug("Buscando persona por ID", "person_id", id)
	person, err := s.repository.GetPersonByID(ctx, id)
	if err != nil {
		s.logger.Error("Error buscando persona por ID", "person_id", id, "error", err)
		return nil, err
	}
	s.logger.Debug("Persona encontrada por ID", "person_id", id, "email", person.Email)
	return person, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	s.logger.Info("Iniciando validaciones de registro", person.ToLogger())

	existingPerson, err := s.repository.GetPersonByEmail(ctx, person.Email)
	if err == nil && existingPerson != nil {
		s.logger.Warn("Intento de registro con email duplicado", "email", person.Email)
		return nil, domain.ErrDuplicateUser
	}

	s.logger.Info("Validaciones de registro completadas", person.ToLogger())
	return &dto.RegistrationResult{
		Person:  person,
		Message: "Validaciones exitosas",
	}, nil
}

func (s service) SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error {
	s.logger.Info("Guardando persona en base de datos", person.ToLogger())
	err := s.repository.SavePerson(ctx, tx, person)
	if err != nil {
		s.logger.Error("Error guardando persona en BD", person.ToLogger(), "error", err)
		return err
	}
	s.logger.Success("Persona guardada en base de datos", person.ToLogger())
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	s.logger.Info("Creando usuario en Keycloak", person.ToLogger())
	userID, err := s.keycloak.CreateUser(ctx, person)
	if err != nil {
		s.logger.Error("Error creando usuario en Keycloak", person.ToLogger(), "error", err)
		return "", err
	}
	s.logger.Success("Usuario creado en Keycloak", person.ToLogger(), "keycloak_user_id", userID)
	return userID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	s.logger.Debug("Configurando password de usuario", "keycloak_user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, true)
	if err != nil {
		s.logger.Error("Error configurando password", "keycloak_user_id", userID, "error", err)
		return err
	}
	s.logger.Success("Password configurado", "keycloak_user_id", userID)
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	s.logger.Info("Asignando rol a usuario", "keycloak_user_id", userID, "role", role)
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		s.logger.Error("Error asignando rol", "keycloak_user_id", userID, "role", role, "error", err)
		return err
	}
	s.logger.Success("Rol asignado", "keycloak_user_id", userID, "role", role)
	return nil
}

func (s service) UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error {
	s.logger.Debug("Actualizando keycloak_user_id en BD", "person_id", personID, "keycloak_user_id", keycloakUserID)
	err := s.repository.PatchPerson(ctx, tx, personID, keycloakUserID)
	if err != nil {
		s.logger.Error("Error actualizando keycloak_user_id", "person_id", personID, "error", err)
		return err
	}
	s.logger.Success("Keycloak_user_id actualizado", "person_id", personID, "keycloak_user_id", keycloakUserID)
	return nil
}

func (s service) RollbackPerson(ctx context.Context, personID string) error {
	s.logger.Warn("Ejecutando rollback: eliminando persona de BD", "person_id", personID)
	err := s.repository.DeletePerson(ctx, nil, personID)
	if err != nil {
		s.logger.Error("Error en rollback de persona", "person_id", personID, "error", err)
		return err
	}
	s.logger.Info("Rollback de persona completado", "person_id", personID)
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	s.logger.Warn("Ejecutando rollback: eliminando usuario de Keycloak", "keycloak_user_id", keycloakUserID)
	err := s.keycloak.DeleteUser(ctx, keycloakUserID)
	if err != nil {
		s.logger.Error("Error en rollback de usuario Keycloak", "keycloak_user_id", keycloakUserID, "error", err)
		return err
	}
	s.logger.Info("Rollback de usuario Keycloak completado", "keycloak_user_id", keycloakUserID)
	return nil
}
