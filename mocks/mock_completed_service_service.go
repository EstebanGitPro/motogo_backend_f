package mocks

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockCompletedServiceService is a mock implementation of input.CompletedServiceService
type MockCompletedServiceService struct {
	mock.Mock
}

// Transactions

func (m *MockCompletedServiceService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Validations

func (m *MockCompletedServiceService) ValidateBranchServices(ctx context.Context, branchID string, serviceIDs []string) error {
	args := m.Called(ctx, branchID, serviceIDs)
	return args.Error(0)
}

func (m *MockCompletedServiceService) ValidateDiagnosticForMotorcycle(ctx context.Context, diagnosticID, motorcycleID string) error {
	args := m.Called(ctx, diagnosticID, motorcycleID)
	return args.Error(0)
}

func (m *MockCompletedServiceService) ValidateNoActiveService(ctx context.Context, branchID, motorcycleID string) error {
	args := m.Called(ctx, branchID, motorcycleID)
	return args.Error(0)
}

// CRUD - Write

func (m *MockCompletedServiceService) SaveCompletedService(ctx context.Context, tx output.Tx, service *domain.CompletedService) error {
	args := m.Called(ctx, tx, service)
	return args.Error(0)
}

func (m *MockCompletedServiceService) SaveItems(ctx context.Context, tx output.Tx, items []domain.CompletedServiceItem) error {
	args := m.Called(ctx, tx, items)
	return args.Error(0)
}

func (m *MockCompletedServiceService) SaveStatusHistory(ctx context.Context, tx output.Tx, history *domain.ServiceStatusHistory) error {
	args := m.Called(ctx, tx, history)
	return args.Error(0)
}

// CRUD - Read

func (m *MockCompletedServiceService) GetByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedService), args.Error(1)
}

func (m *MockCompletedServiceService) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CompletedService), args.Error(1)
}

func (m *MockCompletedServiceService) GetByBranchID(ctx context.Context, branchID string) ([]domain.CompletedService, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CompletedService), args.Error(1)
}

// CRUD - Delete (HU65)

func (m *MockCompletedServiceService) DeleteCompletedService(ctx context.Context, tx output.Tx, serviceID string, status domain.ServiceStatus) error {
	args := m.Called(ctx, tx, serviceID, status)
	return args.Error(0)
}

// Status Transitions (HU73/HU74)

func (m *MockCompletedServiceService) UpdateStatus(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time) error {
	args := m.Called(ctx, tx, serviceID, status, completionDate)
	return args.Error(0)
}

func (m *MockCompletedServiceService) GetStatusHistory(ctx context.Context, serviceID string) ([]domain.ServiceStatusHistory, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ServiceStatusHistory), args.Error(1)
}

// UpdateStatusWithPrice (Feature A)

func (m *MockCompletedServiceService) UpdateStatusWithPrice(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time, finalPrice *float64) error {
	args := m.Called(ctx, tx, serviceID, status, completionDate, finalPrice)
	return args.Error(0)
}

// UpdateDetails (Feature B)

func (m *MockCompletedServiceService) UpdateDetails(ctx context.Context, tx output.Tx, serviceID string, quotedPrice, finalPrice *float64, notes *string) error {
	args := m.Called(ctx, tx, serviceID, quotedPrice, finalPrice, notes)
	return args.Error(0)
}
