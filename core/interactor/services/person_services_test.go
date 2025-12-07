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
	mockAuthClient.On("SetPassword", ctx, userID, password, true).Return(nil)

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
	mockAuthClient.On("SetPassword", ctx, userID, password, true).Return(kcError)

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
