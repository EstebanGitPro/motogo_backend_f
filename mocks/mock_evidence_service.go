package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockEvidenceService mocks input.EvidenceService
type MockEvidenceService struct {
	mock.Mock
}

func (m *MockEvidenceService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockEvidenceService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) error {
	args := m.Called(ctx, motorcycleID, ownerID)
	return args.Error(0)
}

func (m *MockEvidenceService) CheckEvidenceLimit(ctx context.Context, motorcycleID string) error {
	args := m.Called(ctx, motorcycleID)
	return args.Error(0)
}

func (m *MockEvidenceService) CreateEvidence(ctx context.Context, tx output.Tx, motorcycleID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, tx, motorcycleID, evidence)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

func (m *MockEvidenceService) GetEvidenceByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, evidenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}

func (m *MockEvidenceService) GetEvidenceByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleEvidence), args.Error(1)
}

func (m *MockEvidenceService) ApplyEvidenceUpdates(ctx context.Context, existing *domain.MotorcycleEvidence, updates *domain.MotorcycleEvidence) {
	m.Called(ctx, existing, updates)
}

func (m *MockEvidenceService) UpdateEvidence(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	args := m.Called(ctx, tx, evidence)
	return args.Error(0)
}

func (m *MockEvidenceService) DeleteEvidence(ctx context.Context, tx output.Tx, evidenceID string) error {
	args := m.Called(ctx, tx, evidenceID)
	return args.Error(0)
}

func (m *MockEvidenceService) DeleteStorageFile(ctx context.Context, imageURL string) {
	m.Called(ctx, imageURL)
}

func (m *MockEvidenceService) WithStorageClient(client output.StorageClient) {
	m.Called(client)
}
