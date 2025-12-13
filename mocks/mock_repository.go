package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockTx is a mock implementation of output.Tx
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockRepository is a mock implementation of output.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockRepository) SavePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	args := m.Called(ctx, tx, person)
	return args.Error(0)
}

func (m *MockRepository) UpdatePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	args := m.Called(ctx, tx, person)
	return args.Error(0)
}

func (m *MockRepository) PatchPerson(ctx context.Context, tx output.Tx, id string, keycloakUserID string) error {
	args := m.Called(ctx, tx, id, keycloakUserID)
	return args.Error(0)
}

func (m *MockRepository) DeletePerson(ctx context.Context, tx output.Tx, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockRepository) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}

func (m *MockRepository) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Person), args.Error(1)
}
