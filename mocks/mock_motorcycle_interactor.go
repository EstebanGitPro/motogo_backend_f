package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockMotorcycleInteractor is a mock implementation of input.MotorcycleInteractorInterface
type MockMotorcycleInteractor struct {
	mock.Mock
}

// RegisterMotorcycle mocks the RegisterMotorcycle method
func (m *MockMotorcycleInteractor) RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// GetMotorcycleByID mocks the GetMotorcycleByID method
func (m *MockMotorcycleInteractor) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// GetMotorcyclesByOwner mocks the GetMotorcyclesByOwner method
func (m *MockMotorcycleInteractor) GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Motorcycle), args.Error(1)
}

// GetMotorcycleByLicensePlate mocks the GetMotorcycleByLicensePlate method
func (m *MockMotorcycleInteractor) GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, licensePlate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// UpdateMotorcycle mocks the UpdateMotorcycle method
func (m *MockMotorcycleInteractor) UpdateMotorcycle(ctx context.Context, motorcycleID string, ownerID string, updates *domain.Motorcycle) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID, ownerID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// DeleteMotorcycle mocks the DeleteMotorcycle method
func (m *MockMotorcycleInteractor) DeleteMotorcycle(ctx context.Context, motorcycleID string, ownerID string) error {
	args := m.Called(ctx, motorcycleID, ownerID)
	return args.Error(0)
}

// DeleteProfileImage mocks the DeleteProfileImage method (HU39)
func (m *MockMotorcycleInteractor) DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) error {
	args := m.Called(ctx, motorcycleID, ownerID)
	return args.Error(0)
}

// GetMotorcycleReferences mocks the GetMotorcycleReferences method
func (m *MockMotorcycleInteractor) GetMotorcycleReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

// GetReferencesByBrandID mocks the GetReferencesByBrandID method
func (m *MockMotorcycleInteractor) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx, brandID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

// GrantDiagnosticPermission mocks the GrantDiagnosticPermission method
func (m *MockMotorcycleInteractor) GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string, active bool) (*domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID, branchID, ownerID, active)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiagnosticPermission), args.Error(1)
}

// RevokeDiagnosticPermission mocks the RevokeDiagnosticPermission method
func (m *MockMotorcycleInteractor) RevokeDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) error {
	args := m.Called(ctx, motorcycleID, branchID, ownerID)
	return args.Error(0)
}

// ListDiagnosticPermissions mocks the ListDiagnosticPermissions method
func (m *MockMotorcycleInteractor) ListDiagnosticPermissions(ctx context.Context, motorcycleID, ownerID string) ([]domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticPermission), args.Error(1)
}

// LookupPermissions mocks the LookupPermissions method (no ownership check)
func (m *MockMotorcycleInteractor) LookupPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticPermission), args.Error(1)
}
