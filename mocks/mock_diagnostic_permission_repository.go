package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockDiagnosticPermissionRepository is a mock implementation of output.DiagnosticPermissionRepository
type MockDiagnosticPermissionRepository struct {
	mock.Mock
}

func (m *MockDiagnosticPermissionRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockDiagnosticPermissionRepository) Save(ctx context.Context, tx output.Tx, permission *domain.DiagnosticPermission) error {
	args := m.Called(ctx, tx, permission)
	return args.Error(0)
}

func (m *MockDiagnosticPermissionRepository) Deactivate(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	args := m.Called(ctx, tx, motorcycleID, branchID)
	return args.Error(0)
}

func (m *MockDiagnosticPermissionRepository) GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiagnosticPermission), args.Error(1)
}

func (m *MockDiagnosticPermissionRepository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticPermission), args.Error(1)
}
