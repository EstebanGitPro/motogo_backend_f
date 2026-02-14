package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockMotorcycleRepository is a mock implementation of output.MotorcycleRepository
type MockMotorcycleRepository struct {
	mock.Mock
}

// Transactions

func (m *MockMotorcycleRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Motorcycle operations - write (HU43, HU44, HU45)

func (m *MockMotorcycleRepository) Save(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	args := m.Called(ctx, tx, motorcycle)
	return args.Error(0)
}

func (m *MockMotorcycleRepository) Update(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	args := m.Called(ctx, tx, motorcycle)
	return args.Error(0)
}

func (m *MockMotorcycleRepository) Delete(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

// HardDelete permanently removes a motorcycle (HU45 hybrid delete)
func (m *MockMotorcycleRepository) HardDelete(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

// ClearProfileImageURL sets the profile_image_url to NULL (HU39)
func (m *MockMotorcycleRepository) ClearProfileImageURL(ctx context.Context, tx output.Tx, motorcycleID string) error {
	args := m.Called(ctx, tx, motorcycleID)
	return args.Error(0)
}

// HasServiceHistory checks if motorcycle has any completed_services or diagnostics (HU45 hybrid)
func (m *MockMotorcycleRepository) HasServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	args := m.Called(ctx, motorcycleID)
	return args.Bool(0), args.Error(1)
}

// Motorcycle operations - read (HU46, HU47)

func (m *MockMotorcycleRepository) GetByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Motorcycle), args.Error(1)
}

func (m *MockMotorcycleRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, licensePlate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}

// Reference catalog (HU50, HU40)

func (m *MockMotorcycleRepository) GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

func (m *MockMotorcycleRepository) GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error) {
	args := m.Called(ctx, brandID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleReference), args.Error(1)
}

// Category catalog (HU41)

func (m *MockMotorcycleRepository) GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleCategory), args.Error(1)
}

func (m *MockMotorcycleRepository) GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error) {
	args := m.Called(ctx, categoryName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CategoryLine), args.Error(1)
}

func (m *MockMotorcycleRepository) GetDistinctDisplacements(ctx context.Context) ([]domain.EngineDisplacementRange, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EngineDisplacementRange), args.Error(1)
}

// Validation methods (HU43, HU44)

func (m *MockMotorcycleRepository) ValidateReferenceExists(ctx context.Context, referenceID string) (bool, error) {
	args := m.Called(ctx, referenceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMotorcycleRepository) CheckLicensePlateExists(ctx context.Context, licensePlate string) (bool, error) {
	args := m.Called(ctx, licensePlate)
	return args.Bool(0), args.Error(1)
}
