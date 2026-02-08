package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockScheduleDetailService is a mock implementation of input.ScheduleDetailService
type MockScheduleDetailService struct {
	mock.Mock
}

// Transactions

func (m *MockScheduleDetailService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Schedule Detail CRUD (HU6-9)

func (m *MockScheduleDetailService) CreateDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, tx, detail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, detailID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) UpdateDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error {
	args := m.Called(ctx, tx, detail)
	return args.Error(0)
}

func (m *MockScheduleDetailService) DeleteDetail(ctx context.Context, tx output.Tx, detailID string) error {
	args := m.Called(ctx, tx, detailID)
	return args.Error(0)
}

// Validation

func (m *MockScheduleDetailService) ValidateTimeRange(openingTime, closingTime string) error {
	args := m.Called(openingTime, closingTime)
	return args.Error(0)
}

func (m *MockScheduleDetailService) CheckTimeConflict(ctx context.Context, scheduleID string, dayOfWeek int, openingTime, closingTime string, excludeDetailID string) (bool, error) {
	args := m.Called(ctx, scheduleID, dayOfWeek, openingTime, closingTime, excludeDetailID)
	return args.Bool(0), args.Error(1)
}

// Schedule Exception operations (HU20-25)

func (m *MockScheduleDetailService) CreateException(ctx context.Context, tx output.Tx, exception domain.ScheduleDetail) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, tx, exception)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) GetExceptionsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error) {
	args := m.Called(ctx, exceptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduleDetail), args.Error(1)
}

func (m *MockScheduleDetailService) UpdateException(ctx context.Context, tx output.Tx, exception domain.ScheduleDetail) error {
	args := m.Called(ctx, tx, exception)
	return args.Error(0)
}

func (m *MockScheduleDetailService) DeleteException(ctx context.Context, tx output.Tx, exceptionID string) error {
	args := m.Called(ctx, tx, exceptionID)
	return args.Error(0)
}

func (m *MockScheduleDetailService) SetExceptionActive(ctx context.Context, tx output.Tx, exceptionID string, active bool) error {
	args := m.Called(ctx, tx, exceptionID, active)
	return args.Error(0)
}

func (m *MockScheduleDetailService) CheckExceptionDateConflict(ctx context.Context, scheduleID, excludeExceptionID, startDate, endDate string) (bool, error) {
	args := m.Called(ctx, scheduleID, excludeExceptionID, startDate, endDate)
	return args.Bool(0), args.Error(1)
}
