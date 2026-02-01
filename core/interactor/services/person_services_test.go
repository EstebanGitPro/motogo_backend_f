package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPersonByEmail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	expectedPerson := &domain.Person{
		ID:        "person-123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, "test@example.com").Return(expectedPerson, nil)

	// Act
	person, err := service.GetPersonByEmail(ctx, "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, expectedPerson.Email, person.Email)
	assert.Equal(t, expectedPerson.ID, person.ID)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPersonByEmail_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("record not found")

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, "notfound@example.com").Return(nil, notFoundError)

	// Act
	person, err := service.GetPersonByEmail(ctx, "notfound@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, person)
	assert.Equal(t, notFoundError, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPersonByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	expectedPerson := &domain.Person{
		ID:        "person-123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByID", ctx, "person-123").Return(expectedPerson, nil)

	// Act
	person, err := service.GetPersonByID(ctx, "person-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, expectedPerson.ID, person.ID)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPersonByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("record not found")

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByID", ctx, "not-found-id").Return(nil, notFoundError)

	// Act
	person, err := service.GetPersonByID(ctx, "not-found-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, person)
	assert.Equal(t, notFoundError, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_Success_NeitherExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	notFoundError := errors.New("record not found")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Check DB - not found
	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(nil, notFoundError)

	// Check Keycloak - not found
	mockAuthClient.On("GetUserByEmail", ctx, person.Email).Return(nil, notFoundError)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Validaciones exitosas", result.Message)

	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_BothExist_ReturnsDuplicate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email:     "existing@example.com",
		FirstName: "Existing",
		LastName:  "User",
	}

	existingPerson := &domain.Person{
		ID:    "existing-123",
		Email: person.Email,
	}
	id := "some-id"
	existingKeycloakUser := &gocloak.User{
		ID: &id,
	}

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(existingPerson, nil)
	mockAuthClient.On("GetUserByEmail", ctx, person.Email).Return(existingKeycloakUser, nil)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateUser, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_OnlyDBExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email: "dbonly@example.com",
	}

	existingPerson := &domain.Person{
		ID:    "db-123",
		Email: person.Email,
	}

	notFoundError := errors.New("record not found")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(existingPerson, nil)
	mockAuthClient.On("GetUserByEmail", ctx, person.Email).Return(nil, notFoundError)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrIncompleteRegistration, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_OnlyKeycloakExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email: "kconly@example.com",
	}
	id := "some-id"
	existingKeycloakUser := &gocloak.User{
		ID: &id,
	}

	notFoundError := errors.New("record not found")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(nil, notFoundError)
	mockAuthClient.On("GetUserByEmail", ctx, person.Email).Return(existingKeycloakUser, nil)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrIncompleteRegistration, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_DBConnectionError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email: "test@example.com",
	}

	dbError := errors.New("connection refused")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(nil, dbError)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDatabaseUnavailable, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRegisterPerson_KeycloakConnectionError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		Email: "test@example.com",
	}

	notFoundError := errors.New("record not found")
	kcError := errors.New("timeout")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	mockRepo.On("GetPersonByEmail", ctx, person.Email).Return(nil, notFoundError)
	mockAuthClient.On("GetUserByEmail", ctx, person.Email).Return(nil, kcError)

	// Act
	result, err := service.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrKeycloakUnavailable, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestSavePersonToDB_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		ID:        "person-123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("SavePerson", ctx, mockTx, person).Return(nil)

	// Act
	err := service.SavePersonToDB(ctx, mockTx, person)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestSavePersonToDB_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		ID:    "person-123",
		Email: "test@example.com",
	}

	dbError := errors.New("database error")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("SavePerson", ctx, mockTx, person).Return(dbError)

	// Act
	err := service.SavePersonToDB(ctx, mockTx, person)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateUserInKeycloak_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := &domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	keycloakUserID := "kc-user-123"

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("CreateUser", ctx, person).Return(keycloakUserID, nil)

	// Act
	userID, err := service.CreateUserInKeycloak(ctx, person)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, keycloakUserID, userID)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateUserInKeycloak_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := &domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	keycloakError := errors.New("keycloak creation failed")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("CreateUser", ctx, person).Return("", keycloakError)

	// Act
	userID, err := service.CreateUserInKeycloak(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Equal(t, domain.ErrKeycloakUserCreationFailed, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRollbackKeycloakUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	keycloakUserID := "kc-user-123"

	// Mock expectations
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("DeleteUser", ctx, keycloakUserID).Return(nil)

	// Act
	err := service.RollbackKeycloakUser(ctx, keycloakUserID)

	// Assert
	assert.NoError(t, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRollbackKeycloakUser_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	keycloakUserID := "kc-user-123"
	kcError := errors.New("keycloak error")

	// Mock expectations
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("DeleteUser", ctx, keycloakUserID).Return(kcError)

	// Act
	err := service.RollbackKeycloakUser(ctx, keycloakUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, kcError, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestSetUserPassword_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	userID := "kc-user-123"
	password := "new-password"

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SetPassword", ctx, userID, password, false).Return(nil)

	// Act
	err := service.SetUserPassword(ctx, userID, password)

	// Assert
	assert.NoError(t, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestSetUserPassword_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	userID := "kc-user-123"
	password := "new-password"
	kcError := errors.New("keycloak error")

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SetPassword", ctx, userID, password, false).Return(kcError)

	// Act
	err := service.SetUserPassword(ctx, userID, password)

	// Assert
	assert.Error(t, err)
	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAssignUserRole_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	userID := "kc-user-123"
	role := "admin"

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("AssignRole", ctx, userID, role).Return(nil)

	// Act
	err := service.AssignUserRole(ctx, userID, role)

	// Assert
	assert.NoError(t, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAssignUserRole_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	userID := "kc-user-123"
	role := "admin"
	kcError := errors.New("keycloak error")

	// Mock expectations
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("AssignRole", ctx, userID, role).Return(kcError)

	// Act
	err := service.AssignUserRole(ctx, userID, role)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, kcError, err)

	mockAuthClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdatePersonKeycloakID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	personID := "person-123"
	keycloakUserID := "kc-user-123"

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("PatchPerson", ctx, mockTx, personID, keycloakUserID).Return(nil)

	// Act
	err := service.UpdatePersonKeycloakID(ctx, mockTx, personID, keycloakUserID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdatePersonKeycloakID_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	personID := "person-123"
	keycloakUserID := "kc-user-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("PatchPerson", ctx, mockTx, personID, keycloakUserID).Return(dbError)

	// Act
	err := service.UpdatePersonKeycloakID(ctx, mockTx, personID, keycloakUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRollbackPerson_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	personID := "person-123"

	// Mock expectations
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockRepo.On("DeletePerson", ctx, nil, personID).Return(nil)

	// Act
	err := service.RollbackPerson(ctx, personID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRollbackPerson_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	personID := "person-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("DeletePerson", ctx, nil, personID).Return(dbError)

	// Act
	err := service.RollbackPerson(ctx, personID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// ============================================
// BeginTx Tests
// ============================================

func TestPersonService_BeginTx_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)

	// Act
	tx, err := service.BeginTx(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_BeginTx_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	dbError := errors.New("connection error")
	mockRepo.On("BeginTx", ctx).Return(nil, dbError)

	// Act
	tx, err := service.BeginTx(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tx)
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetPersonByKeycloakID Tests
// ============================================

func TestGetPersonByKeycloakID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	expectedPerson := &domain.Person{
		ID:             "person-123",
		Email:          "test@example.com",
		KeycloakUserID: "kc-user-123",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByKeycloakID", ctx, "kc-user-123").Return(expectedPerson, nil)

	// Act
	person, err := service.GetPersonByKeycloakID(ctx, "kc-user-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, "person-123", person.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetPersonByKeycloakID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("record not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByKeycloakID", ctx, "non-existent").Return(nil, notFoundError)

	// Act
	person, err := service.GetPersonByKeycloakID(ctx, "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, person)
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetUserByEmail Tests
// ============================================

func TestGetUserByEmail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	id := "kc-user-123"
	expectedUser := &gocloak.User{ID: &id}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByEmail", ctx, "test@example.com").Return(expectedUser, nil)

	// Act
	user, err := service.GetUserByEmail(ctx, "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockAuthClient.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("user not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByEmail", ctx, "notfound@example.com").Return(nil, notFoundError)

	// Act
	user, err := service.GetUserByEmail(ctx, "notfound@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// SendVerificationEmail Tests
// ============================================

func TestSendVerificationEmail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SendVerificationEmail", ctx, "kc-user-123").Return(nil)

	// Act
	err := service.SendVerificationEmail(ctx, "kc-user-123")

	// Assert
	assert.NoError(t, err)
	mockAuthClient.AssertExpectations(t)
}

func TestSendVerificationEmail_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	kcError := errors.New("keycloak error")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SendVerificationEmail", ctx, "kc-user-123").Return(kcError)

	// Act
	err := service.SendVerificationEmail(ctx, "kc-user-123")

	// Assert
	assert.Error(t, err)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// SendPasswordResetEmail Tests
// ============================================

func TestSendPasswordResetEmail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SendPasswordResetEmail", ctx, "test@example.com").Return(nil)

	// Act
	err := service.SendPasswordResetEmail(ctx, "test@example.com")

	// Assert
	assert.NoError(t, err)
	mockAuthClient.AssertExpectations(t)
}

func TestSendPasswordResetEmail_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	kcError := errors.New("keycloak error")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("SendPasswordResetEmail", ctx, "test@example.com").Return(kcError)

	// Act
	err := service.SendPasswordResetEmail(ctx, "test@example.com")

	// Assert
	assert.Error(t, err)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// Login Tests
// ============================================

func TestLogin_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	id := "kc-user-123"
	email := "test@example.com"
	verified := true
	existingUser := &gocloak.User{
		ID:            &id,
		Email:         &email,
		EmailVerified: &verified,
	}
	expectedToken := &gocloak.JWT{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(existingUser, nil)
	mockAuthClient.On("LoginUser", ctx, email, "password123").Return(expectedToken, nil)

	// Act
	token, err := service.Login(ctx, email, "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "access-token", token.AccessToken)
	mockAuthClient.AssertExpectations(t)
}

func TestLogin_EmailNotVerified(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	id := "kc-user-123"
	email := "unverified@example.com"
	verified := false
	existingUser := &gocloak.User{
		ID:            &id,
		Email:         &email,
		EmailVerified: &verified,
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(existingUser, nil)
	mockAuthClient.On("SendVerificationEmail", ctx, id).Return(nil)

	// Act
	token, err := service.Login(ctx, email, "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, domain.ErrorEmailNotVerified, err)
	mockAuthClient.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("user not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByEmail", ctx, "notfound@example.com").Return(nil, notFoundError)

	// Act
	token, err := service.Login(ctx, "notfound@example.com", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// RefreshToken Tests
// ============================================

func TestRefreshToken_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	expectedToken := &gocloak.JWT{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("RefreshToken", ctx, "old-refresh-token").Return(expectedToken, nil)

	// Act
	token, err := service.RefreshToken(ctx, "old-refresh-token")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "new-access-token", token.AccessToken)
	mockAuthClient.AssertExpectations(t)
}

func TestRefreshToken_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	kcError := errors.New("invalid refresh token")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("RefreshToken", ctx, "invalid-token").Return(nil, kcError)

	// Act
	token, err := service.RefreshToken(ctx, "invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, token)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// UpdatePersonProfile Tests
// ============================================

func TestUpdatePersonProfile_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		ID:             "person-123",
		Email:          "test@example.com",
		FirstName:      "Updated",
		LastName:       "Name",
		KeycloakUserID: "kc-user-123",
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("UpdatePerson", ctx, mockTx, person).Return(nil)
	mockAuthClient.On("UpdateUser", ctx, mock.AnythingOfType("*gocloak.User")).Return(nil)

	// Act
	err := service.UpdatePersonProfile(ctx, mockTx, person)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestUpdatePersonProfile_DBError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		ID:    "person-123",
		Email: "test@example.com",
	}

	dbError := errors.New("database error")

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("UpdatePerson", ctx, mockTx, person).Return(dbError)

	// Act
	err := service.UpdatePersonProfile(ctx, mockTx, person)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePersonProfile_WithoutKeycloakSync(t *testing.T) {
	// Arrange - when KeycloakUserID is empty, skip Keycloak sync
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	person := domain.Person{
		ID:             "person-123",
		Email:          "test@example.com",
		FirstName:      "Updated",
		LastName:       "Name",
		KeycloakUserID: "", // Empty - skip Keycloak sync
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("UpdatePerson", ctx, mockTx, person).Return(nil)
	// Note: UpdateUser should NOT be called

	// Act
	err := service.UpdatePersonProfile(ctx, mockTx, person)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	// Verify UpdateUser was never called
	mockAuthClient.AssertNotCalled(t, "UpdateUser")
}

// ============================================
// ChangePassword Tests
// ============================================

func TestChangePassword_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	keycloakUserID := "kc-user-123"
	email := "test@example.com"
	existingUser := &gocloak.User{
		ID:    &keycloakUserID,
		Email: &email,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByID", ctx, keycloakUserID).Return(existingUser, nil)
	mockAuthClient.On("LoginUser", ctx, email, "current-password").Return(&gocloak.JWT{}, nil)
	mockAuthClient.On("SetPassword", ctx, keycloakUserID, "new-password", false).Return(nil)

	// Act
	err := service.ChangePassword(ctx, keycloakUserID, "current-password", "new-password")

	// Assert
	assert.NoError(t, err)
	mockAuthClient.AssertExpectations(t)
}

func TestChangePassword_InvalidCurrentPassword(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	keycloakUserID := "kc-user-123"
	email := "test@example.com"
	existingUser := &gocloak.User{
		ID:    &keycloakUserID,
		Email: &email,
	}

	loginError := errors.New("invalid credentials")

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByID", ctx, keycloakUserID).Return(existingUser, nil)
	mockAuthClient.On("LoginUser", ctx, email, "wrong-password").Return(nil, loginError)

	// Act
	err := service.ChangePassword(ctx, keycloakUserID, "wrong-password", "new-password")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	mockAuthClient.AssertExpectations(t)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	notFoundError := errors.New("user not found")

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByID", ctx, "non-existent").Return(nil, notFoundError)

	// Act
	err := service.ChangePassword(ctx, "non-existent", "current-password", "new-password")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockAuthClient.AssertExpectations(t)
}

func TestChangePassword_SetPasswordError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	keycloakUserID := "kc-user-123"
	email := "test@example.com"
	existingUser := &gocloak.User{
		ID:    &keycloakUserID,
		Email: &email,
	}

	setPasswordError := errors.New("keycloak error")

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockAuthClient.On("GetUserByID", ctx, keycloakUserID).Return(existingUser, nil)
	mockAuthClient.On("LoginUser", ctx, email, "current-password").Return(&gocloak.JWT{}, nil)
	mockAuthClient.On("SetPassword", ctx, keycloakUserID, "new-password", false).Return(setPasswordError)

	// Act
	err := service.ChangePassword(ctx, keycloakUserID, "current-password", "new-password")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPasswordUpdateFailed, err)
	mockAuthClient.AssertExpectations(t)
}

// ============================================
// CheckAndCleanInconsistentState Tests
// ============================================

func TestCheckAndCleanInconsistentState_BothExist_ConsistentState(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	email := "test@example.com"
	existingPerson := &domain.Person{ID: "person-123", Email: email}
	id := "kc-user-123"
	existingKeycloakUser := &gocloak.User{ID: &id}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, email).Return(existingPerson, nil)
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(existingKeycloakUser, nil)

	// Act
	err := service.CheckAndCleanInconsistentState(ctx, email)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestCheckAndCleanInconsistentState_NeitherExist_ConsistentState(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	email := "nonexistent@example.com"
	notFoundError := errors.New("not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, email).Return(nil, notFoundError)
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(nil, notFoundError)

	// Act
	err := service.CheckAndCleanInconsistentState(ctx, email)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestCheckAndCleanInconsistentState_OnlyKeycloakExists_CleansOrphan(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	email := "orphan@example.com"
	id := "kc-user-123"
	existingKeycloakUser := &gocloak.User{ID: &id}
	notFoundError := errors.New("not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, email).Return(nil, notFoundError)
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(existingKeycloakUser, nil)
	mockAuthClient.On("DeleteUser", ctx, id).Return(nil)

	// Act
	err := service.CheckAndCleanInconsistentState(ctx, email)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestCheckAndCleanInconsistentState_OnlyDBExists_CleansOrphan(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockRepository)
	mockAuthClient := new(mocks.MockAuthClient)
	mockLogger := new(mocks.MockLogger)

	service := services.NewService(mockRepo, mockAuthClient, mockLogger)

	email := "dbonly@example.com"
	existingPerson := &domain.Person{ID: "person-123", Email: email}
	notFoundError := errors.New("not found")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockRepo.On("GetPersonByEmail", ctx, email).Return(existingPerson, nil)
	mockAuthClient.On("GetUserByEmail", ctx, email).Return(nil, notFoundError)
	mockRepo.On("DeletePerson", ctx, nil, "person-123").Return(nil)

	// Act
	err := service.CheckAndCleanInconsistentState(ctx, email)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}
