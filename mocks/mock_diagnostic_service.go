package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockDiagnosticService mocks input.DiagnosticService
type MockDiagnosticService struct {
	mock.Mock
}

func (m *MockDiagnosticService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockDiagnosticService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) error {
	args := m.Called(ctx, motorcycleID, ownerID)
	return args.Error(0)
}

func (m *MockDiagnosticService) ValidateBranchExists(ctx context.Context, branchID string) error {
	args := m.Called(ctx, branchID)
	return args.Error(0)
}

func (m *MockDiagnosticService) UpsertDiagnostic(ctx context.Context, tx output.Tx, motorcycleID, branchID string, problemDescription *string, evidenceURLs []string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, tx, motorcycleID, branchID, problemDescription, evidenceURLs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) GetDiagnosticByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, diagnosticID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) GetDiagnosticsByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Diagnostic), args.Error(1)
}

func (m *MockDiagnosticService) ApplyDiagnosticUpdates(existing *domain.Diagnostic, updates *domain.Diagnostic) {
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
