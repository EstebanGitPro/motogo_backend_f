package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockEvidenceRepository is a mock implementation of output.EvidenceRepository (HU16-19)
type MockEvidenceRepository struct {
	mock.Mock
}

// Transactions

func (m *MockEvidenceRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Evidence operations - write (HU16, HU19)

func (m *MockEvidenceRepository) Save(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	args := m.Called(ctx, tx, evidence)
	return args.Error(0)
}

func (m *MockEvidenceRepository) Update(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	args := m.Called(ctx, tx, evidence)
	return args.Error(0)
}

func (m *MockEvidenceRepository) Delete(ctx context.Context, tx output.Tx, evidenceID string) error {
	args := m.Called(ctx, tx, evidenceID)
	return args.Error(0)
}

// Evidence operations - read (HU18)

func (m *MockEvidenceRepository) GetByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, evidenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

func (m *MockEvidenceRepository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleEvidence), args.Error(1)
}

func (m *MockEvidenceRepository) CountByMotorcycleID(ctx context.Context, motorcycleID string) (int, error) {
	args := m.Called(ctx, motorcycleID)
	return args.Int(0), args.Error(1)
}
