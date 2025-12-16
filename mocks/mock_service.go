package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of input.Service
type MockService struct {
	mock.Mock
}

// Transacciones

func (m *MockService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Person - Validaciones y consultas

func (m *MockService) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	args := m.Called(ctx, person)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegistrationResult), args.Error(1)
}

func (m *MockService) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockService) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockService) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// Person - Operaciones transaccionales de BD

func (m *MockService) SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error {
	args := m.Called(ctx, tx, person)
	return args.Error(0)
}

func (m *MockService) UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error {
	args := m.Called(ctx, tx, personID, keycloakUserID)
	return args.Error(0)
}

// Person - Operaciones de Keycloak

func (m *MockService) CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	args := m.Called(ctx, person)
	return args.String(0), args.Error(1)
}

func (m *MockService) SetUserPassword(ctx context.Context, userID string, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockService) AssignUserRole(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

// Person - Compensaciones (rollback)

func (m *MockService) RollbackPerson(ctx context.Context, personID string) error {
	args := m.Called(ctx, personID)
	return args.Error(0)
}

func (m *MockService) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	args := m.Called(ctx, keycloakUserID)
	return args.Error(0)
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (m *MockService) SendVerificationEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockService) SendPasswordResetEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
