package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRepository is a mock implementation of output.ScheduleRepository
type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockScheduleRepository) SaveSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	args := m.Called(ctx, tx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	args := m.Called(ctx, tx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) DeleteSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	args := m.Called(ctx, tx, scheduleID)
	return args.Error(0)
}

func (m *MockScheduleRepository) SetActive(ctx context.Context, tx output.Tx, scheduleID string, active bool) error {
	args := m.Called(ctx, tx, scheduleID, active)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleRepository) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}
