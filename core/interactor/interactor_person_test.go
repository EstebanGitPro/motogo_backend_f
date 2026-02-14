package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// RegisterPerson (existing tests)
// ============================================

func TestRegisterPerson_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
		Role:      "customer",
	}

	registrationResult := &dto.RegistrationResult{
		Message: "Validaciones iniciales completadas",
	}

	keycloakUserID := "keycloak-user-123"

	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return(keycloakUserID, nil)
	mockService.On("SetUserPassword", ctx, keycloakUserID, person.Password).Return(nil)
	mockService.On("AssignUserRole", ctx, keycloakUserID, string(person.Role)).Return(nil)
	mockService.On("UpdatePersonKeycloakID", ctx, mockTx, mock.AnythingOfType("string"), keycloakUserID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockService.On("SendVerificationEmail", ctx, keycloakUserID).Return(nil)

	result, err := personInteractor.RegisterPerson(ctx, person)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, keycloakUserID, result.Person.KeycloakUserID)
	assert.NotEmpty(t, result.Person.ID)
	assert.Equal(t, "Usuario registrado exitosamente", result.Message)

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterPerson_FailsAtKeycloakCreation_RollsBack(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
		Role:      "customer",
	}

	registrationResult := &dto.RegistrationResult{
		Message: "Validaciones iniciales completadas",
	}

	keycloakError := errors.New("keycloak service unavailable")

	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return("", keycloakError)
	mockTx.On("Rollback").Return(nil)

	result, err := personInteractor.RegisterPerson(ctx, person)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.ErrKeycloakUserCreationFailed, err)

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterPerson_FailsAtSaveDB_RollsBackKeycloak(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
		Role:      "customer",
	}

	registrationResult := &dto.RegistrationResult{
		Message: "Validaciones iniciales completadas",
	}

	keycloakUserID := "keycloak-user-123"
	dbError := errors.New("update keycloak ID failed")

	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return(keycloakUserID, nil)
	mockService.On("SetUserPassword", ctx, keycloakUserID, person.Password).Return(nil)
	mockService.On("AssignUserRole", ctx, keycloakUserID, string(person.Role)).Return(nil)
	mockService.On("UpdatePersonKeycloakID", ctx, mockTx, mock.AnythingOfType("string"), keycloakUserID).Return(dbError)
	mockTx.On("Rollback").Return(nil)
	mockService.On("RollbackKeycloakUser", ctx, keycloakUserID).Return(nil)

	result, err := personInteractor.RegisterPerson(ctx, person)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dbError, err)
	mockService.AssertCalled(t, "RollbackKeycloakUser", ctx, keycloakUserID)
	mockTx.AssertCalled(t, "Rollback")

	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================
// Login
// ============================================

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	jwt := &gocloak.JWT{
		AccessToken:  "access-token-123",
		TokenType:    "Bearer",
		ExpiresIn:    300,
		RefreshToken: "refresh-token-456",
	}

	mockService.On("Login", ctx, "user@test.com", "password123").Return(jwt, nil)

	result, err := personInteractor.Login(ctx, "user@test.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access-token-123", result.AccessToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, 300, result.ExpiresIn)
	assert.Equal(t, "refresh-token-456", result.RefreshToken)
	mockService.AssertExpectations(t)
}

func TestLogin_Error(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("Login", ctx, "user@test.com", "wrong").Return(nil, domain.ErrInvalidCredentials)

	result, err := personInteractor.Login(ctx, "user@test.com", "wrong")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	mockService.AssertExpectations(t)
}

// ============================================
// RefreshToken
// ============================================

func TestRefreshToken_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	jwt := &gocloak.JWT{
		AccessToken:  "new-access-token",
		TokenType:    "Bearer",
		ExpiresIn:    300,
		RefreshToken: "new-refresh-token",
	}

	mockService.On("RefreshToken", ctx, "old-refresh-token").Return(jwt, nil)

	result, err := personInteractor.RefreshToken(ctx, "old-refresh-token")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Equal(t, "new-refresh-token", result.RefreshToken)
	mockService.AssertExpectations(t)
}

func TestRefreshToken_Error(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("RefreshToken", ctx, "expired-token").Return(nil, errors.New("token expired"))

	result, err := personInteractor.RefreshToken(ctx, "expired-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockService.AssertExpectations(t)
}

// ============================================
// ChangePassword
// ============================================

func TestChangePassword_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ChangePassword", ctx, "kc-user-1", "old-pass", "new-pass").Return(nil)

	err := personInteractor.ChangePassword(ctx, "kc-user-1", "old-pass", "new-pass")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ChangePassword", ctx, "kc-unknown", "old", "new").Return(domain.ErrUserNotFound)

	err := personInteractor.ChangePassword(ctx, "kc-unknown", "old", "new")

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	mockService.AssertExpectations(t)
}

func TestChangePassword_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ChangePassword", ctx, "kc-user-1", "wrong-pass", "new").Return(domain.ErrInvalidCredentials)

	err := personInteractor.ChangePassword(ctx, "kc-user-1", "wrong-pass", "new")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	mockService.AssertExpectations(t)
}

func TestChangePassword_UpdateFailed(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ChangePassword", ctx, "kc-user-1", "old", "new").Return(domain.ErrPasswordUpdateFailed)

	err := personInteractor.ChangePassword(ctx, "kc-user-1", "old", "new")

	assert.ErrorIs(t, err, domain.ErrPasswordUpdateFailed)
	mockService.AssertExpectations(t)
}

func TestChangePassword_KeycloakUnavailable(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ChangePassword", ctx, "kc-user-1", "old", "new").Return(domain.ErrKeycloakUnavailable)

	err := personInteractor.ChangePassword(ctx, "kc-user-1", "old", "new")

	assert.ErrorIs(t, err, domain.ErrKeycloakUnavailable)
	mockService.AssertExpectations(t)
}

// ============================================
// UpdateProfile
// ============================================

func TestUpdateProfile_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{ID: "person-1", FirstName: "Updated", LastName: "Name"}

	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdatePersonProfile", ctx, mockTx, person).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := personInteractor.UpdateProfile(ctx, person)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "person-1", result.ID)
	assert.Equal(t, "Updated", result.FirstName)
	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUpdateProfile_BeginTxError(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{ID: "person-1"}

	mockService.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := personInteractor.UpdateProfile(ctx, person)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockService.AssertExpectations(t)
}

func TestUpdateProfile_UpdateError_RollsBack(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{ID: "person-1"}

	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdatePersonProfile", ctx, mockTx, person).Return(errors.New("update failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := personInteractor.UpdateProfile(ctx, person)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUpdateProfile_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockService)

	person := domain.Person{ID: "person-1"}

	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("UpdatePersonProfile", ctx, mockTx, person).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := personInteractor.UpdateProfile(ctx, person)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
	mockService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================
// GetPublicContact
// ============================================

func TestGetPublicContact_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	expected := &domain.Person{ID: "person-1", PhoneNumber: "3001234567"}
	mockService.On("GetPersonByID", ctx, "person-1").Return(expected, nil)

	result, err := personInteractor.GetPublicContact(ctx, "person-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockService.AssertExpectations(t)
}

func TestGetPublicContact_NotFound(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("GetPersonByID", ctx, "unknown").Return(nil, domain.ErrUserNotFound)

	result, err := personInteractor.GetPublicContact(ctx, "unknown")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockService.AssertExpectations(t)
}

// ============================================
// VerifyEmailByToken
// ============================================

func TestVerifyEmailByToken_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("VerifyEmailByToken", ctx, "valid-token").Return("user@test.com", nil)

	email, err := personInteractor.VerifyEmailByToken(ctx, "valid-token")

	assert.NoError(t, err)
	assert.Equal(t, "user@test.com", email)
	mockService.AssertExpectations(t)
}

func TestVerifyEmailByToken_InvalidToken(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("VerifyEmailByToken", ctx, "bad-token").Return("", domain.ErrInvalidToken)

	email, err := personInteractor.VerifyEmailByToken(ctx, "bad-token")

	assert.ErrorIs(t, err, domain.ErrInvalidToken)
	assert.Empty(t, email)
	mockService.AssertExpectations(t)
}

func TestVerifyEmailByToken_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("VerifyEmailByToken", ctx, "token-orphan").Return("orphan@test.com", domain.ErrUserNotFound)

	email, err := personInteractor.VerifyEmailByToken(ctx, "token-orphan")

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.Equal(t, "orphan@test.com", email)
	mockService.AssertExpectations(t)
}

func TestVerifyEmailByToken_AlreadyVerified(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("VerifyEmailByToken", ctx, "token-verified").Return("verified@test.com", domain.ErrEmailAlreadyVerified)

	email, err := personInteractor.VerifyEmailByToken(ctx, "token-verified")

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyVerified)
	assert.Equal(t, "verified@test.com", email)
	mockService.AssertExpectations(t)
}

// ============================================
// ResetPasswordWithToken
// ============================================

func TestResetPasswordWithToken_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ResetPasswordWithToken", ctx, "valid-token", "newPass123").Return(nil)

	err := personInteractor.ResetPasswordWithToken(ctx, "valid-token", "newPass123")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestResetPasswordWithToken_InvalidToken(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ResetPasswordWithToken", ctx, "bad-token", "newPass").Return(domain.ErrInvalidToken)

	err := personInteractor.ResetPasswordWithToken(ctx, "bad-token", "newPass")

	assert.ErrorIs(t, err, domain.ErrInvalidToken)
	mockService.AssertExpectations(t)
}

func TestResetPasswordWithToken_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ResetPasswordWithToken", ctx, "token", "newPass").Return(domain.ErrUserNotFound)

	err := personInteractor.ResetPasswordWithToken(ctx, "token", "newPass")

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	mockService.AssertExpectations(t)
}

func TestResetPasswordWithToken_PasswordUpdateFailed(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("ResetPasswordWithToken", ctx, "token", "newPass").Return(domain.ErrPasswordUpdateFailed)

	err := personInteractor.ResetPasswordWithToken(ctx, "token", "newPass")

	assert.ErrorIs(t, err, domain.ErrPasswordUpdateFailed)
	mockService.AssertExpectations(t)
}

// ============================================
// RequestPasswordReset
// ============================================

func TestRequestPasswordReset_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("SendPasswordResetEmail", ctx, "user@test.com").Return(nil)

	// Always returns nil for security (never reveal if user exists)
	err := personInteractor.RequestPasswordReset(ctx, "user@test.com")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestRequestPasswordReset_EmailNotFound_StillReturnsNil(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	// Even if SendPasswordResetEmail fails, RequestPasswordReset returns nil
	mockService.On("SendPasswordResetEmail", ctx, "nonexistent@test.com").Return(errors.New("user not found"))

	err := personInteractor.RequestPasswordReset(ctx, "nonexistent@test.com")

	assert.NoError(t, err) // Security: always nil
	mockService.AssertExpectations(t)
}

// ============================================
// ResendVerificationEmail
// ============================================

func TestResendVerificationEmail_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	userID := "kc-user-1"
	user := &gocloak.User{
		ID:            &userID,
		EmailVerified: gocloak.BoolP(false),
	}

	mockService.On("GetUserByEmail", ctx, "user@test.com").Return(user, nil)
	mockService.On("SendVerificationEmail", ctx, userID).Return(nil)

	err := personInteractor.ResendVerificationEmail(ctx, "user@test.com")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestResendVerificationEmail_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("GetUserByEmail", ctx, "unknown@test.com").Return(nil, errors.New("not found"))

	err := personInteractor.ResendVerificationEmail(ctx, "unknown@test.com")

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	mockService.AssertExpectations(t)
}

func TestResendVerificationEmail_AlreadyVerified(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	userID := "kc-user-1"
	user := &gocloak.User{
		ID:            &userID,
		EmailVerified: gocloak.BoolP(true),
	}

	mockService.On("GetUserByEmail", ctx, "verified@test.com").Return(user, nil)

	err := personInteractor.ResendVerificationEmail(ctx, "verified@test.com")

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyVerified)
	mockService.AssertExpectations(t)
}

func TestResendVerificationEmail_SendFails(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	userID := "kc-user-1"
	user := &gocloak.User{
		ID:            &userID,
		EmailVerified: gocloak.BoolP(false),
	}
	sendErr := errors.New("smtp failed")

	mockService.On("GetUserByEmail", ctx, "user@test.com").Return(user, nil)
	mockService.On("SendVerificationEmail", ctx, userID).Return(sendErr)

	err := personInteractor.ResendVerificationEmail(ctx, "user@test.com")

	assert.Error(t, err)
	assert.Equal(t, sendErr, err)
	mockService.AssertExpectations(t)
}

// ============================================
// DeleteKeycloakUser
// ============================================

func TestDeleteKeycloakUser_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("RollbackKeycloakUser", ctx, "kc-user-1").Return(nil)

	err := personInteractor.DeleteKeycloakUser(ctx, "kc-user-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeleteKeycloakUser_Error(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("RollbackKeycloakUser", ctx, "kc-user-1").Return(errors.New("keycloak error"))

	err := personInteractor.DeleteKeycloakUser(ctx, "kc-user-1")

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

// ============================================
// DeletePersonFromDB
// ============================================

func TestDeletePersonFromDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("RollbackPerson", ctx, "person-1").Return(nil)

	err := personInteractor.DeletePersonFromDB(ctx, "person-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeletePersonFromDB_Error(t *testing.T) {
	ctx := context.Background()
	mockService := new(mocks.MockService)
	personInteractor := interactor.NewInteractor(mockService)

	mockService.On("RollbackPerson", ctx, "person-1").Return(errors.New("db error"))

	err := personInteractor.DeletePersonFromDB(ctx, "person-1")

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}
