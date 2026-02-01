package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockScheduleDetailRepository is a mock implementation of output.ScheduleDetailRepository
type MockScheduleDetailRepository struct {
	mock.Mock
}

func (m *MockScheduleDetailRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockScheduleDetailRepository) SaveScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	args := m.Called(ctx, tx, detail)
	return args.Error(0)
}

func (m *MockScheduleDetailRepository) UpdateScheduleDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	args := m.Called(ctx, tx, detail)
	return args.Error(0)
}

func (m *MockScheduleDetailRepository) DeleteScheduleDetail(ctx context.Context, tx output.Tx, detailID string) error {
	args := m.Called(ctx, tx, detailID)
	return args.Error(0)
}

func (m *MockScheduleDetailRepository) GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, detailID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) GetDetailsByScheduleAndDay(ctx context.Context, scheduleID string, dayOfWeek int) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) CheckTimeConflict(ctx context.Context, scheduleID string, dayOfWeek int, openingTime, closingTime string, excludeDetailID string) (bool, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek, openingTime, closingTime, excludeDetailID)
	return args.Bool(0), args.Error(1)
}

func (m *MockScheduleDetailRepository) GetExceptionsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) GetExceptionsByScheduleIDForUpdate(ctx context.Context, tx output.Tx, scheduleID string) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, tx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, exceptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailRepository) CheckExceptionDateConflict(ctx context.Context, scheduleID string, excludeExceptionID string, startDate string, endDate string) (bool, error) {
	args := m.Called(ctx, scheduleID, excludeExceptionID, startDate, endDate)
	return args.Bool(0), args.Error(1)
}

func (m *MockScheduleDetailRepository) CheckDayIsClosed(ctx context.Context, scheduleID string, dayOfWeek int, excludeDetailID string) (bool, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek, excludeDetailID)
	return args.Bool(0), args.Error(1)
}

func (m *MockScheduleDetailRepository) CheckDayHasTimeSlots(ctx context.Context, scheduleID string, dayOfWeek int, excludeDetailID string) (bool, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek, excludeDetailID)
	return args.Bool(0), args.Error(1)
}

func (m *MockScheduleDetailRepository) CheckExceptionIsRedundant(ctx context.Context, scheduleID string, dayOfWeek int) (bool, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek)
	return args.Bool(0), args.Error(1)
}
