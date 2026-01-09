package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/jwt"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/Nerzal/gocloak/v13"
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

func (s service) GetPersonByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Person, error) {
	s.logger.Debug(logger.LogPersonServiceSearchByKeycloakID, "keycloak_user_id", keycloakUserID)
	person, err := s.repository.GetPersonByKeycloakID(ctx, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogPersonServiceErrorByKeycloakID, "keycloak_user_id", keycloakUserID, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogPersonServiceFoundByKeycloakID, "keycloak_user_id", keycloakUserID, "person_id", person.ID)
	return person, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	s.logger.Info(logger.LogPersonServiceValidationStart, person.ToLogger())
	s.logger.Debug(logger.LogDualSystemCheck, "email", person.Email)

	// Check in business database - IMPORTANTE: detectar indisponibilidad
	existingPerson, errDB := s.repository.GetPersonByEmail(ctx, person.Email)

	// CRÍTICO: Si hay error de conexión/timeout, la base de datos está caída
	if errDB != nil {
		if isConnectionError(errDB) || isTimeoutError(errDB) {
			//TODO: Agregar mensaje de log aquí
			s.logger.Error(logger.LogDatabaseUnavailable,
				"email", person.Email,
				"error", errDB,
				"error_type", "connection")
			return nil, domain.ErrDatabaseUnavailable
		}
		// Si el error NO es de conexión, asumimos que el usuario no existe
		// (errores como "record not found" son normales)
	}

	dbExists := errDB == nil && existingPerson != nil

	//TODO: CRÍTICO: Si hay error de conexión/timeout, Keycloak está caído
	// Check in Keycloak - IMPORTANTE: detectar indisponibilidad
	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, person.Email)

	// CRÍTICO: Si hay error de conexión/timeout, Keycloak está caído
	if errKC != nil {
		if isConnectionError(errKC) || isTimeoutError(errKC) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				"email", person.Email,
				"error", errKC,
				"error_type", "connection")
			return nil, domain.ErrKeycloakUnavailable
		}
		// Si el error NO es de conexión, asumimos que el usuario no existe
		// (errores como 404 Not Found son normales)
	}

	kcExists := errKC == nil && keycloakUser != nil

	// Log where the user exists
	if dbExists && kcExists {
		s.logger.Warn(logger.LogUserExistsInBoth, "email", person.Email)
		return nil, domain.ErrDuplicateUser // Usuario ya registrado completamente
	}

	if dbExists && !kcExists {
		s.logger.Warn(logger.LogUserExistsOnlyInDB,
			"email", person.Email,
			"person_id", existingPerson.ID,
			"action", "will be cleaned")
		// Retornar error de registro incompleto (mensaje: intente más tarde)
		return nil, domain.ErrIncompleteRegistration
	}

	if !dbExists && kcExists {
		s.logger.Warn(logger.LogUserExistsOnlyInKeycloak,
			"email", person.Email,
			"keycloak_id", *keycloakUser.ID,
			"action", "will be cleaned")
		// Retornar error de registro incompleto (mensaje: intente más tarde)
		return nil, domain.ErrIncompleteRegistration
	}

	s.logger.Debug(logger.LogUserNotFoundInEither, "email", person.Email)
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
		// Distinguish between unavailability and other errors
		if isConnectionError(err) || isTimeoutError(err) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				person.ToLogger(),
				"error", err,
				"error_type", "connection")
			return "", domain.ErrKeycloakUnavailable
		}

		s.logger.Error(logger.LogPersonServiceKeycloakError, person.ToLogger(), "error", err)
		return "", domain.ErrKeycloakUserCreationFailed
	}

	s.logger.Success(logger.LogPersonServiceCreatedKeycloak, person.ToLogger(), "keycloak_user_id", userID)
	return userID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	s.logger.Debug(logger.LogPersonServicePasswordSet, "keycloak_user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, false)
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
	s.logger.Debug(logger.LogDualSystemCheck, "email", email)

	// Check if user exists in business DB
	personInDB, errDB := s.repository.GetPersonByEmail(ctx, email)
	dbExists := errDB == nil && personInDB != nil

	// Check if user exists in Keycloak
	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, email)
	kcExists := errKC == nil && keycloakUser != nil

	// Both exist or neither exist - consistent state
	if (dbExists && kcExists) || (!dbExists && !kcExists) {
		if dbExists && kcExists {
			s.logger.Debug(logger.LogUserExistsInBoth, "email", email)
		} else {
			s.logger.Debug(logger.LogUserNotFoundInEither, "email", email)
		}
		return nil
	}

	// Log inconsistent state with details
	s.logger.Warn(logger.LogInconsistentStateDetect,
		"email", email,
		"in_database", dbExists,
		"in_keycloak", kcExists,
		"db_person_id", func() string {
			if dbExists {
				return personInDB.ID
			}
			return "N/A"
		}(),
		"kc_user_id", func() string {
			if kcExists {
				return *keycloakUser.ID
			}
			return "N/A"
		}())

	// User exists only in Keycloak - clean it
	if !dbExists && kcExists {
		s.logger.Info(logger.LogPersonServiceCleaningOrphan,
			"email", email,
			"source", "keycloak",
			"keycloak_user_id", *keycloakUser.ID,
			"reason", "missing in business database")

		if err := s.keycloak.DeleteUser(ctx, *keycloakUser.ID); err != nil {
			s.logger.Error(logger.LogPersonServiceOrphanCleanError,
				"email", email,
				"source", "keycloak",
				"keycloak_user_id", *keycloakUser.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogPersonServiceOrphanCleaned,
			"email", email,
			"source", "keycloak",
			"action", "deleted from Keycloak")
		return nil // Limpiado exitosamente, puede reintentar
	}

	// User exists only in DB - clean it
	if dbExists && !kcExists {
		s.logger.Info(logger.LogPersonServiceCleaningOrphan,
			"email", email,
			"source", "database",
			"person_id", personInDB.ID,
			"reason", "missing in Keycloak")

		if err := s.repository.DeletePerson(ctx, nil, personInDB.ID); err != nil {
			s.logger.Error(logger.LogPersonServiceOrphanCleanError,
				"email", email,
				"source", "database",
				"person_id", personInDB.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogPersonServiceOrphanCleaned,
			"email", email,
			"source", "database",
			"action", "deleted from business database")
		return nil // Limpiado exitosamente, puede reintentar
	}

	return nil
}

// Helper functions to detect error types

// isConnectionError checks if an error is a connection-related error
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for common connection error patterns
	return contains(errStr, "connection refused") ||
		contains(errStr, "no such host") ||
		contains(errStr, "connection reset") ||
		contains(errStr, "network is unreachable") ||
		contains(errStr, "connect: connection refused")
}

// isTimeoutError checks if an error is a timeout-related error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "timeout") ||
		contains(errStr, "deadline exceeded") ||
		contains(errStr, "context deadline exceeded")
}

// contains is a case-insensitive substring check
func contains(s, substr string) bool {
	// Simple case-insensitive check
	for i := 0; i+len(substr) <= len(s); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			// Convert to lowercase for comparison
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// isPasswordPolicyError checks if an error is related to Keycloak password policy violation
// Keycloak returns 400 Bad Request with specific messages when password doesn't meet policy requirements
func isPasswordPolicyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Common Keycloak password policy error patterns
	return contains(errStr, "invalidPasswordMinLength") ||
		contains(errStr, "invalidPasswordMinUpperCase") ||
		contains(errStr, "invalidPasswordMinLowerCase") ||
		contains(errStr, "invalidPasswordMinDigits") ||
		contains(errStr, "invalidPasswordMinSpecialChars") ||
		contains(errStr, "invalidPasswordNotUsername") ||
		contains(errStr, "invalidPasswordRegex") ||
		contains(errStr, "passwordHistory") ||
		contains(errStr, "Password policy not met") ||
		contains(errStr, "password policy") ||
		contains(errStr, "400 Bad Request") // Keycloak typically returns 400 for policy violations
}

// GetUserByEmail retrieves a user from Keycloak by email
func (s service) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	s.logger.Debug(logger.LogKeycloakSearchUserByEmail, "email", email)
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogKeycloakSearchUserByEmailOK, "email", email, "user_id", *user.ID)
	return user, nil
}

// SendVerificationEmail sends a verification email to a user
func (s service) SendVerificationEmail(ctx context.Context, userID string) error {
	s.logger.Debug(logger.LogKeycloakSendVerificationEmail, "user_id", userID)
	err := s.keycloak.SendVerificationEmail(ctx, userID)
	if err != nil {
		s.logger.Error(logger.LogKeycloakSendVerificationEmailError, "user_id", userID, "error", err)
		return err
	}
	s.logger.Success(logger.LogKeycloakSendVerificationEmailOK, "user_id", userID)
	return nil
}

// SendPasswordResetEmail sends a password reset email to a user
func (s service) SendPasswordResetEmail(ctx context.Context, email string) error {
	s.logger.Debug(logger.LogKeycloakSendPasswordReset, "email", email)
	err := s.keycloak.SendPasswordResetEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return err
	}
	s.logger.Success(logger.LogKeycloakSendPasswordResetOK, "email", email)
	return nil
}

// Login authenticates a user with email and password
func (s service) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	s.logger.Debug(logger.LogKeycloakUserLogin, "email", email)

	// First, check if user's email is verified
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, domain.ErrUserNotFound
	}

	// Validate email is verified
	if user.EmailVerified == nil || !*user.EmailVerified {
		s.logger.Warn(logger.LogKeycloakEmailNotVerified, "email", email, "user_id", *user.ID)

		// Auto-resend verification email to help user complete verification
		s.logger.Info(logger.LogKeycloakResendingVerificationEmail, "email", email, "user_id", *user.ID)
		if resendErr := s.keycloak.SendVerificationEmail(ctx, *user.ID); resendErr != nil {
			// Log error but don't fail the login - the main error is still "email not verified"
			s.logger.Error(logger.LogKeycloakResendVerificationEmailError, "email", email, "user_id", *user.ID, "error", resendErr)
		} else {
			s.logger.Success(logger.LogKeycloakResendVerificationEmailOK, "email", email, "user_id", *user.ID)
		}

		return nil, domain.ErrorEmailNotVerified
	}

	// Email is verified, proceed with login
	token, err := s.keycloak.LoginUser(ctx, email, password)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserLoginError, "email", email, "error", err)
		return nil, err
	}
	s.logger.Success(logger.LogKeycloakUserLoginOK, "email", email)
	return token, nil
}

// RefreshToken obtains a new access token using the refresh token
// This is called by the frontend when the access token expires
func (s service) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	s.logger.Debug("RefreshToken called")

	// Delegate to Keycloak client
	token, err := s.keycloak.RefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("RefreshToken failed", "error", err)
		return nil, err
	}

	s.logger.Success("RefreshToken completed successfully")
	return token, nil
}

// VerifyEmailByToken receives a JWT token, extracts the email, and marks it as verified in Keycloak
// This is called when a user clicks on the verification link from the email
// Returns the extracted email on success
func (s service) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	s.logger.Info(logger.LogKeycloakEmailVerify)

	// Extract email from the JWT token
	tokenParser := jwt.NewTokenParser()
	email, err := tokenParser.ExtractEmailFromToken(token)
	if err != nil {
		s.logger.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "failed to extract email from token")
		return "", domain.ErrInvalidToken
	}

	s.logger.Debug("Email extracted from token", "email", email)

	// Get user from Keycloak by email
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return "", domain.ErrUserNotFound
	}

	// Check if already verified
	if user.EmailVerified != nil && *user.EmailVerified {
		s.logger.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email, "user_id", *user.ID)
		return email, domain.ErrEmailAlreadyVerified
	}

	// Verify the email in Keycloak
	if err := s.keycloak.VerifyEmail(ctx, *user.ID); err != nil {
		s.logger.Error(logger.LogKeycloakEmailVerifyError, "email", email, "user_id", *user.ID, "error", err)
		return "", err
	}

	s.logger.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "user_id", *user.ID)
	return email, nil
}

// ResetPasswordWithToken receives a JWT token, extracts the email, and updates the password in Keycloak
// This is called when a user submits the password reset form from the email link
// Returns nil on success, error otherwise
func (s service) ResetPasswordWithToken(ctx context.Context, token string, newPassword string) error {
	s.logger.Info(logger.LogPasswordResetStart)

	// Extract email from the JWT token
	tokenParser := jwt.NewTokenParser()
	email, err := tokenParser.ExtractEmailFromToken(token)
	if err != nil {
		s.logger.Error(logger.LogPasswordResetTokenError, "error", err, "reason", "failed to extract email from token")
		return domain.ErrInvalidToken
	}

	s.logger.Debug(logger.LogPasswordResetEmailExtracted, "email", email)

	// Get user from Keycloak by email
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogPasswordResetUserNotFound, "email", email, "error", err)
		return domain.ErrUserNotFound
	}

	s.logger.Debug(logger.LogPasswordResetUserFound, "email", email, "user_id", *user.ID)

	// Update password in Keycloak using existing method (reuses SetUserPassword)
	if err := s.SetUserPassword(ctx, *user.ID, newPassword); err != nil {
		s.logger.Error(logger.LogPasswordResetUpdateError, "email", email, "user_id", *user.ID, "error", err)
		return domain.ErrPasswordUpdateFailed
	}

	s.logger.Success(logger.LogPasswordResetSuccess, "email", email, "user_id", *user.ID)
	return nil
}

// ChangePassword verifies the current password by attempting a login
// If successful, updates the password in Keycloak (HU57)
func (s service) ChangePassword(ctx context.Context, keycloakUserID, currentPassword, newPassword string) error {
	s.logger.Info(logger.LogChangePasswordStart, "keycloak_user_id", keycloakUserID)

	// Get user to retrieve email for password verification
	user, err := s.keycloak.GetUserByID(ctx, keycloakUserID)
	if err != nil {
		// Check if Keycloak is unavailable
		if isConnectionError(err) || isTimeoutError(err) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				"keycloak_user_id", keycloakUserID,
				"error", err,
				"error_type", "connection")
			return domain.ErrKeycloakUnavailable
		}
		s.logger.Error(logger.LogChangePasswordUserNotFound, "keycloak_user_id", keycloakUserID, "error", err)
		return domain.ErrUserNotFound
	}

	// Verify current password by attempting login
	_, err = s.keycloak.LoginUser(ctx, *user.Email, currentPassword)
	if err != nil {
		// Check if Keycloak is unavailable during login attempt
		if isConnectionError(err) || isTimeoutError(err) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				"keycloak_user_id", keycloakUserID,
				"error", err,
				"error_type", "connection")
			return domain.ErrKeycloakUnavailable
		}
		s.logger.Warn(logger.LogChangePasswordInvalidCurrent, "keycloak_user_id", keycloakUserID)
		return domain.ErrInvalidCredentials
	}

	// Current password verified, set new password
	if err := s.keycloak.SetPassword(ctx, keycloakUserID, newPassword, false); err != nil {
		// Check if Keycloak is unavailable during password update
		if isConnectionError(err) || isTimeoutError(err) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				"keycloak_user_id", keycloakUserID,
				"error", err,
				"error_type", "connection")
			return domain.ErrKeycloakUnavailable
		}
		// Check if error is due to password policy violation
		if isPasswordPolicyError(err) {
			s.logger.Warn("Password policy violation", "keycloak_user_id", keycloakUserID, "error", err)
			return domain.ErrPasswordPolicyViolation
		}
		s.logger.Error(logger.LogChangePasswordUpdateError, "keycloak_user_id", keycloakUserID, "error", err)
		return domain.ErrPasswordUpdateFailed
	}

	s.logger.Success(logger.LogChangePasswordSuccess, "keycloak_user_id", keycloakUserID)
	return nil
}

// UpdatePersonProfile updates person data in DB and optionally syncs to Keycloak (HU52)
// This method updates the person's profile information excluding email and password
func (s service) UpdatePersonProfile(ctx context.Context, tx output.Tx, person domain.Person) error {
	s.logger.Info(logger.LogUpdateProfileStart, "person_id", person.ID, "email", person.Email)

	// Update person in database
	if err := s.repository.UpdatePerson(ctx, tx, person); err != nil {
		s.logger.Error(logger.LogUpdateProfileError, "person_id", person.ID, "error", err)
		return err
	}
	s.logger.Success(logger.LogUpdateProfileDBSuccess, "person_id", person.ID)

	// Optionally sync first_name and last_name to Keycloak for consistency
	if person.KeycloakUserID != "" {
		// Build gocloak.User structure for update
		keycloakUser := &gocloak.User{
			ID:        gocloak.StringP(person.KeycloakUserID),
			FirstName: gocloak.StringP(person.FirstName),
			LastName:  gocloak.StringP(person.LastName),
		}
		if err := s.keycloak.UpdateUser(ctx, keycloakUser); err != nil {
			// Log warning but don't fail - DB update is the primary source of truth
			s.logger.Warn(logger.LogUpdateProfileKeycloakSyncWarn,
				"person_id", person.ID,
				"keycloak_user_id", person.KeycloakUserID,
				"error", err)
		} else {
			s.logger.Success(logger.LogUpdateProfileKeycloakSyncOK,
				"person_id", person.ID,
				"keycloak_user_id", person.KeycloakUserID)
		}
	}

	s.logger.Success(logger.LogUpdateProfileSuccess, "person_id", person.ID)
	return nil
}
