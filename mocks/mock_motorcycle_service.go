package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockMotorcycleService is a mock implementation of input.MotorcycleService
type MockMotorcycleService struct {
	mock.Mock
}

// Transactions

func (m *MockMotorcycleService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Validation

func (m *MockMotorcycleService) ValidateReferenceExists(ctx context.Context, referenceID string) error {
	args := m.Called(ctx, referenceID)
	return args.Error(0)
}

func (m *MockMotorcycleService) CheckLicensePlateUnique(ctx context.Context, licensePlate string) error {
	args := m.Called(ctx, licensePlate)
	return args.Error(0)
}

func (m *MockMotorcycleService) ValidateOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// Motorcycle CRUD (HU43-47)

func (m *MockMotorcycleService) RegisterMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	args := m.Called(ctx, tx, motorcycle)
	return args.Error(0)
}

func (m *MockMotorcycleService) GetByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) GetByOwnerID(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, licensePlate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleService) ApplyUpdates(existing *domain.Motorcycle, updates *domain.Motorcycle) error {
	args := m.Called(existing, updates)
	return args.Error(0)
}

func (m *MockMotorcycleService) UpdateMotorcycle(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	args := m.Called(ctx, tx, motorcycle)
	return args.Error(0)
}

func (m *MockMotorcycleService) DeleteMotorcycle(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

func (m *MockMotorcycleService) ClearProfileImageURL(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

// Hybrid delete strategy (HU45)

func (m *MockMotorcycleService) HasServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	args := m.Called(ctx, motorcycleID)
	return args.Bool(0), args.Error(1)
}

// Storage cleanup

func (m *MockMotorcycleService) DeleteStorageFile(ctx context.Context, url string) {
	m.Called(ctx, url)
}

// Reference catalog (HU50, HU40)

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

// Category catalog (HU41)

func (m *MockMotorcycleService) GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleCategory), args.Error(1)
}

func (m *MockMotorcycleService) GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error) {
	args := m.Called(ctx, categoryName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CategoryLine), args.Error(1)
}

// GetDistinctDisplacements mocks (HU49)
func (m *MockMotorcycleService) GetDistinctDisplacements(ctx context.Context) ([]domain.EngineDisplacementRange, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EngineDisplacementRange), args.Error(1)
}

// GetRatingRanges mocks (HU48)
func (m *MockMotorcycleService) GetRatingRanges(ctx context.Context) ([]domain.RatingRange, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RatingRange), args.Error(1)
}

// Diagnostic Permissions

func (m *MockMotorcycleService) GrantPermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string, active bool) (*domain.DiagnosticPermission, error) {
	args := m.Called(ctx, tx, motorcycleID, branchID, active)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiagnosticPermission), args.Error(1)
}

func (m *MockMotorcycleService) RevokePermission(ctx context.Context, tx output.Tx, motorcycleID, branchID string) error {
	args := m.Called(ctx, tx, motorcycleID, branchID)
	return args.Error(0)
}

func (m *MockMotorcycleService) ListPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticPermission), args.Error(1)
}
