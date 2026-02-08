package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockMotorcycleService is a mock for input.MotorcycleService
type MockMotorcycleService struct {
	mock.Mock
}

func (m *MockMotorcycleService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockMotorcycleService) BeginPermissionTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockMotorcycleService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) ValidateReferenceExists(ctx context.Context, referenceID string) error {
	args := m.Called(ctx, referenceID)
	return args.Error(0)
}

func (m *MockMotorcycleService) ValidateLicensePlateUnique(ctx context.Context, licensePlate string) error {
	args := m.Called(ctx, licensePlate)
	return args.Error(0)
}

func (m *MockMotorcycleService) CreateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	args := m.Called(ctx, tx, motorcycle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, licensePlate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) ApplyMotorcycleUpdates(ctx context.Context, motorcycle *domain.Motorcycle, updates *domain.Motorcycle) error {
	args := m.Called(ctx, motorcycle, updates)
	return args.Error(0)
}

func (m *MockMotorcycleService) UpdateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	args := m.Called(ctx, tx, motorcycle)
	return args.Error(0)
}

func (m *MockMotorcycleService) CheckServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	args := m.Called(ctx, motorcycleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMotorcycleService) DeleteMotorcycle(ctx context.Context, tx output.Tx, motorcycleID string, hasHistory bool) error {
	args := m.Called(ctx, tx, motorcycleID, hasHistory)
	return args.Error(0)
}

func (m *MockMotorcycleService) DeleteProfileImage(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

func (m *MockMotorcycleService) DeleteStorageFile(ctx context.Context, imageURL string) {
	m.Called(ctx, imageURL)
}

func (m *MockMotorcycleService) GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

func (m *MockMotorcycleService) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx, brandID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

func (m *MockMotorcycleService) GrantDiagnosticPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string) (*domain.DiagnosticPermission, error) {
	args := m.Called(ctx, tx, motorcycleID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiagnosticPermission), args.Error(1)
}

func (m *MockMotorcycleService) RevokeDiagnosticPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	args := m.Called(ctx, tx, motorcycleID, branchID)
	return args.Error(0)
}

func (m *MockMotorcycleService) ListDiagnosticPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticPermission), args.Error(1)
}

func (m *MockMotorcycleService) WithStorageClient(client output.StorageClient) {
	m.Called(client)
}
