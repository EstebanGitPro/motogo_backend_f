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
	s.logger.Debug(logger.LogPersonServiceSearchByEmail, "email", email)
	person, err := s.repository.GetPersonByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceErrorByEmail, "email", email, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogPersonServiceFoundByEmail, "email", email, "person_id", person.ID)
	return person, nil
}

func (s service) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	s.logger.Debug(logger.LogPersonServiceSearchByID, "person_id", id)
	person, err := s.repository.GetPersonByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceErrorByID, "person_id", id, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogPersonServiceFoundByID, "person_id", id, "email", person.Email)
	return person, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	s.logger.Info(logger.LogPersonServiceValidationStart, person.ToLogger())

	existingPerson, err := s.repository.GetPersonByEmail(ctx, person.Email)
	if err == nil && existingPerson != nil {
		s.logger.Warn(logger.LogPersonServiceDuplicateEmail, "email", person.Email)
		return nil, domain.ErrDuplicateUser
	}

	s.logger.Info(logger.LogPersonServiceValidationComplete, person.ToLogger())
	return &dto.RegistrationResult{
		Person:  person,
		Message: "Validaciones exitosas",
	}, nil
}

func (s service) SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error {
	s.logger.Info(logger.LogPersonServiceSavingToDB, person.ToLogger())
	err := s.repository.SavePerson(ctx, tx, person)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceSaveError, person.ToLogger(), "error", err)
		return err
	}
	s.logger.Success(logger.LogPersonServiceSavedToDB, person.ToLogger())
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	s.logger.Info(logger.LogPersonServiceCreatingKeycloak, person.ToLogger())
	userID, err := s.keycloak.CreateUser(ctx, person)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceKeycloakError, person.ToLogger(), "error", err)
		return "", err
	}
	s.logger.Success(logger.LogPersonServiceCreatedKeycloak, person.ToLogger(), "keycloak_user_id", userID)
	return userID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	s.logger.Debug(logger.LogPersonServicePasswordSet, "keycloak_user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, true)
	if err != nil {
		s.logger.Error(logger.LogPersonServicePasswordError, "keycloak_user_id", userID, "error", err)
		return err
	}
	s.logger.Success(logger.LogPersonServicePasswordSetOK, "keycloak_user_id", userID)
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	s.logger.Info(logger.LogPersonServiceRoleAssigning, "keycloak_user_id", userID, "role", role)
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceRoleError, "keycloak_user_id", userID, "role", role, "error", err)
		return err
	}
	s.logger.Success(logger.LogPersonServiceRoleAssigned, "keycloak_user_id", userID, "role", role)
	return nil
}

func (s service) UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error {
	s.logger.Debug(logger.LogPersonServiceKeycloakIDUpdate, "person_id", personID, "keycloak_user_id", keycloakUserID)
	err := s.repository.PatchPerson(ctx, tx, personID, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceKeycloakIDUpdateError, "person_id", personID, "error", err)
		return err
	}
	s.logger.Success(logger.LogPersonServiceKeycloakIDUpdated, "person_id", personID, "keycloak_user_id", keycloakUserID)
	return nil
}

func (s service) RollbackPerson(ctx context.Context, personID string) error {
	s.logger.Warn(logger.LogPersonServiceRollbackPerson, "person_id", personID)
	err := s.repository.DeletePerson(ctx, nil, personID)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceRollbackPersonError, "person_id", personID, "error", err)
		return err
	}
	s.logger.Info(logger.LogPersonServiceRollbackPersonComplete, "person_id", personID)
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	s.logger.Warn(logger.LogPersonServiceRollbackKeycloak, "keycloak_user_id", keycloakUserID)
	err := s.keycloak.DeleteUser(ctx, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceRollbackKeycloakError, "keycloak_user_id", keycloakUserID, "error", err)
		return err
	}
	s.logger.Info(logger.LogPersonServiceRollbackKeycloakComplete, "keycloak_user_id", keycloakUserID)
	return nil
}

func (s service) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	s.logger.Debug(logger.LogPersonServiceSearchByEmail, "email", email)

	// Check if user exists in business DB
	personInDB, errDB := s.repository.GetPersonByEmail(ctx, email)

	// Check if user exists in Keycloak
	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, email)

	// Both exist or neither exist - consistent state
	if (errDB == nil && errKC == nil) || (errDB != nil && errKC != nil) {
		return nil
	}

	s.logger.Warn(logger.LogPersonServiceInconsistentStateDetected,
		"email", email,
		"in_db", errDB == nil,
		"in_keycloak", errKC == nil)

	// User exists only in Keycloak - clean it
	if errDB != nil && errKC == nil {
		s.logger.Info(logger.LogPersonServiceCleaningOrphan,
			"email", email,
			"source", "keycloak",
			"keycloak_user_id", *keycloakUser.ID)

		if err := s.keycloak.DeleteUser(ctx, *keycloakUser.ID); err != nil {
			s.logger.Error(logger.LogPersonServiceOrphanCleanError,
				"email", email,
				"keycloak_user_id", *keycloakUser.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogPersonServiceOrphanCleaned,
			"email", email,
			"source", "keycloak")
	}

	// User exists only in DB - clean it
	if errDB == nil && errKC != nil {
		s.logger.Info(logger.LogPersonServiceCleaningOrphan,
			"email", email,
			"source", "database",
			"person_id", personInDB.ID)

		if err := s.repository.DeletePerson(ctx, nil, personInDB.ID); err != nil {
			s.logger.Error(logger.LogPersonServiceOrphanCleanError,
				"email", email,
				"person_id", personInDB.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogPersonServiceOrphanCleaned,
			"email", email,
			"source", "database")
	}

	return domain.ErrKeycloakInconsistentState
}
