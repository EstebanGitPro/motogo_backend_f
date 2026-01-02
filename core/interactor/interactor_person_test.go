package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterPerson_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	personInteractor := interactor.NewInteractor(mockService, mockLogger) // Corrected: Pass mockService directly

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

	// Mock expectations

	// WithTraceID and Logging
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Step 1: RegisterPerson (validaciones)
	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)

	// Step 1.5: CheckAndCleanInconsistentState
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)

	// Step 2: BeginTx
	mockService.On("BeginTx", ctx).Return(mockTx, nil)

	// Step 3: SavePersonToDB
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)

	// Step 4: CreateUserInKeycloak
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return(keycloakUserID, nil)

	// Step 5: SetUserPassword
	mockService.On("SetUserPassword", ctx, keycloakUserID, person.Password).Return(nil)

	// Step 6: AssignUserRole
	mockService.On("AssignUserRole", ctx, keycloakUserID, person.Role).Return(nil)

	// Step 7: UpdatePersonKeycloakID
	mockService.On("UpdatePersonKeycloakID", ctx, mockTx, mock.AnythingOfType("string"), keycloakUserID).Return(nil)

	// Step 8: Commit
	mockTx.On("Commit").Return(nil)

	// Step 9: SendVerificationEmail (non-blocking)
	mockService.On("SendVerificationEmail", ctx, keycloakUserID).Return(nil)

	// Act
	result, err := personInteractor.RegisterPerson(ctx, person)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, keycloakUserID, result.Person.KeycloakUserID)
	assert.NotEmpty(t, result.Person.ID)
	assert.Equal(t, "Usuario registrado exitosamente", result.Message)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterPerson_FailsAtKeycloakCreation_RollsBack(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	personInteractor := interactor.NewInteractor(mockService, mockLogger) // Corrected: Pass mockService directly

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

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)

	// Keycloak falla aquí
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return("", keycloakError)

	// Rollback expectations
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := personInteractor.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result) // Result can be non-nil even on error
	assert.Equal(t, domain.ErrKeycloakUserCreationFailed, err)

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterPerson_FailsAtSaveDB_RollsBackKeycloak(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := new(mocks.MockService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	personInteractor := interactor.NewInteractor(mockService, mockLogger) // Corrected: Pass mockService directly

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

	// Mock expectations
	mockLogger.On("WithTraceID", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Success", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	mockService.On("RegisterPerson", ctx, mock.AnythingOfType("domain.Person")).Return(registrationResult, nil)
	mockService.On("CheckAndCleanInconsistentState", ctx, person.Email).Return(nil)
	mockService.On("BeginTx", ctx).Return(mockTx, nil)
	mockService.On("SavePersonToDB", ctx, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)
	mockService.On("CreateUserInKeycloak", ctx, mock.AnythingOfType("*domain.Person")).Return(keycloakUserID, nil)
	mockService.On("SetUserPassword", ctx, keycloakUserID, person.Password).Return(nil)
	mockService.On("AssignUserRole", ctx, keycloakUserID, person.Role).Return(nil)

	// UpdatePersonKeycloakID falla aquí
	mockService.On("UpdatePersonKeycloakID", ctx, mockTx, mock.AnythingOfType("string"), keycloakUserID).Return(dbError)

	// Rollback expectations
	mockTx.On("Rollback").Return(nil)
	mockService.On("RollbackKeycloakUser", ctx, keycloakUserID).Return(nil)

	// Act
	result, err := personInteractor.RegisterPerson(ctx, person)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result) // Result can be non-nil even on error
	assert.Equal(t, dbError, err)

	// Verify rollback fue llamado
	mockService.AssertCalled(t, "RollbackKeycloakUser", ctx, keycloakUserID)
	mockTx.AssertCalled(t, "Rollback")

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
