package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/mock"
)

// MockPersonService is a mock implementation of input.Service
type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockPersonService) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	args := m.Called(ctx, person)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegistrationResult), args.Error(1)
}

func (m *MockPersonService) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockPersonService) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockPersonService) GetPersonByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Person, error) {
	args := m.Called(ctx, keycloakUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockPersonService) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockPersonService) SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error {
	args := m.Called(ctx, tx, person)
	return args.Error(0)
}

func (m *MockPersonService) UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error {
	args := m.Called(ctx, tx, personID, keycloakUserID)
	return args.Error(0)
}

func (m *MockPersonService) CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	args := m.Called(ctx, person)
	return args.String(0), args.Error(1)
}

func (m *MockPersonService) SetUserPassword(ctx context.Context, userID string, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockPersonService) AssignUserRole(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockPersonService) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (m *MockPersonService) SendVerificationEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPersonService) SendPasswordResetEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockPersonService) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

func (m *MockPersonService) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

func (m *MockPersonService) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockPersonService) ResetPasswordWithToken(ctx context.Context, token string, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func (m *MockPersonService) ChangePassword(ctx context.Context, keycloakUserID, currentPassword, newPassword string) error {
	args := m.Called(ctx, keycloakUserID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockPersonService) UpdatePersonProfile(ctx context.Context, tx output.Tx, person domain.Person) error {
	args := m.Called(ctx, tx, person)
	return args.Error(0)
}

func (m *MockPersonService) RollbackPerson(ctx context.Context, personID string) error {
	args := m.Called(ctx, personID)
	return args.Error(0)
}

func (m *MockPersonService) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	args := m.Called(ctx, keycloakUserID)
	return args.Error(0)
}
