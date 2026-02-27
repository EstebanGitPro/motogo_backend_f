package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockDiagnosticRepository is a mock implementation of output.DiagnosticRepository
type MockDiagnosticRepository struct {
	mock.Mock
}

func (m *MockDiagnosticRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockDiagnosticRepository) Save(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	args := m.Called(ctx, tx, diagnostic)
	return args.Error(0)
}

func (m *MockDiagnosticRepository) Update(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	args := m.Called(ctx, tx, diagnostic)
	return args.Error(0)
}

func (m *MockDiagnosticRepository) Delete(ctx context.Context, tx output.Tx, diagnosticID string) error {
	args := m.Called(ctx, tx, diagnosticID)
	return args.Error(0)
}

func (m *MockDiagnosticRepository) GetByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, diagnosticID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticRepository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticRepository) GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, motorcycleID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}
