package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockEvidenceInteractor is a mock implementation of input.EvidenceInteractorInterface
type MockEvidenceInteractor struct {
	mock.Mock
}

// CreateEvidence mocks the CreateEvidence method
func (m *MockEvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motorcycleID, ownerID, evidence)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

// GetEvidenceByID mocks the GetEvidenceByID method
func (m *MockEvidenceInteractor) GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, evidenceID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

// ListEvidenceByMotorcycle mocks the ListEvidenceByMotorcycle method
func (m *MockEvidenceInteractor) ListEvidenceByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motorcycleID, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleEvidence), args.Error(1)
}

// DeleteEvidence mocks the DeleteEvidence method
func (m *MockEvidenceInteractor) DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error {
	args := m.Called(ctx, evidenceID, ownerID)
	return args.Error(0)
}

// UpdateEvidence mocks the UpdateEvidence method (HU17)
func (m *MockEvidenceInteractor) UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, evidenceID, ownerID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

// LookupEvidence mocks the LookupEvidence method (no ownership check)
func (m *MockEvidenceInteractor) LookupEvidence(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleEvidence), args.Error(1)
}
