package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockScheduleService is a mock implementation of input.ScheduleService
type MockScheduleService struct {
	mock.Mock
}

// Transactions

func (m *MockScheduleService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Schedule CRUD (HU30-33)

func (m *MockScheduleService) CreateSchedule(ctx context.Context, tx output.Tx, branchID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, tx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleService) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleService) GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleService) UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	args := m.Called(ctx, tx, schedule)
	return args.Error(0)
}

func (m *MockScheduleService) DeleteSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	args := m.Called(ctx, tx, scheduleID)
	return args.Error(0)
}

// Activation (HU34, HU35)

func (m *MockScheduleService) ActivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	args := m.Called(ctx, tx, scheduleID)
	return args.Error(0)
}

func (m *MockScheduleService) DeactivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	args := m.Called(ctx, tx, scheduleID)
	return args.Error(0)
}
