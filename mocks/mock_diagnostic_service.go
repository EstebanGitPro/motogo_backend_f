package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockDiagnosticService is a mock implementation of input.DiagnosticService
type MockDiagnosticService struct {
	mock.Mock
}

// Transactions

func (m *MockDiagnosticService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Validations

func (m *MockDiagnosticService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockDiagnosticService) ValidateBranchExists(ctx context.Context, branchID string) error {
	args := m.Called(ctx, branchID)
	return args.Error(0)
}

// Diagnostic CRUD

func (m *MockDiagnosticService) RegisterOrUpdateDiagnostic(ctx context.Context, tx output.Tx, motorcycleID, branchID string, problemDescription *string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, tx, motorcycleID, branchID, problemDescription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) GetByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, diagnosticID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) ApplyDiagnosticUpdates(existing, updates *domain.Diagnostic) {
	m.Called(existing, updates)
}

func (m *MockDiagnosticService) UpdateDiagnostic(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	args := m.Called(ctx, tx, diagnostic)
	return args.Error(0)
}

func (m *MockDiagnosticService) DeleteDiagnostic(ctx context.Context, tx output.Tx, diagnosticID string) error {
	args := m.Called(ctx, tx, diagnosticID)
	return args.Error(0)
}

func (m *MockDiagnosticService) SetSolution(ctx context.Context, tx output.Tx, diagnosticID string, solution string) error {
	args := m.Called(ctx, tx, diagnosticID, solution)
	return args.Error(0)
}
