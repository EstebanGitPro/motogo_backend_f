package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRepo is a mock for ScheduleRepository
type MockScheduleRepo struct {
	mock.Mock
}

func (m *MockScheduleRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockScheduleRepo) SaveSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	args := m.Called(ctx, tx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepo) GetScheduleByID(ctx context.Context, id string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleRepo) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchSchedule), args.Error(1)
}

func (m *MockScheduleRepo) UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	args := m.Called(ctx, tx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepo) DeleteSchedule(ctx context.Context, tx output.Tx, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockScheduleRepo) SetActive(ctx context.Context, tx output.Tx, id string, active bool) error {
	args := m.Called(ctx, tx, id, active)
	return args.Error(0)
}

// ============================================
// NewScheduleService Tests
// ============================================

func TestNewScheduleService(t *testing.T) {
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := NewScheduleService(mockSchedRepo, mockBranchRepo)
	assert.NotNil(t, service)
}

// ============================================
// BeginTx Tests
// ============================================

func TestScheduleService_BeginTx_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)
	mockSchedRepo.On("BeginTx", ctx).Return(mockTx, nil)

	tx, err := service.BeginTx(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

// ============================================
// CreateSchedule Tests
// ============================================

func TestCreateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	branchID := "branch-123"
	branch := &domain.Branch{ID: branchID, Name: "Test Branch"}

	mockBranchRepo.On("GetBranchByID", ctx, branchID).Return(branch, nil)
	mockSchedRepo.On("GetScheduleByBranchID", ctx, branchID).Return(nil, nil)
	mockSchedRepo.On("SaveSchedule", ctx, mockTx, mock.AnythingOfType("domain.BranchSchedule")).Return(nil)

	result, err := service.CreateSchedule(ctx, mockTx, branchID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, branchID, result.BranchID)
	assert.True(t, result.Active)
}

func TestCreateSchedule_BranchNotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockBranchRepo.On("GetBranchByID", ctx, "not-found").Return(nil, nil)

	result, err := service.CreateSchedule(ctx, mockTx, "not-found")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotFound, err)
	assert.Nil(t, result)
}

func TestCreateSchedule_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	branchID := "branch-123"
	branch := &domain.Branch{ID: branchID}
	existingSchedule := &domain.BranchSchedule{ID: "sched-123", BranchID: branchID}

	mockBranchRepo.On("GetBranchByID", ctx, branchID).Return(branch, nil)
	mockSchedRepo.On("GetScheduleByBranchID", ctx, branchID).Return(existingSchedule, nil)

	result, err := service.CreateSchedule(ctx, mockTx, branchID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleAlreadyExists, err)
	assert.Nil(t, result)
}

// ============================================
// GetScheduleByBranchID Tests
// ============================================

func TestGetScheduleByBranchID_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	expected := &domain.BranchSchedule{ID: "sched-123", BranchID: "branch-123"}
	mockSchedRepo.On("GetScheduleByBranchID", ctx, "branch-123").Return(expected, nil)

	result, err := service.GetScheduleByBranchID(ctx, "branch-123")

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestGetScheduleByBranchID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockSchedRepo.On("GetScheduleByBranchID", ctx, "not-found").Return(nil, nil)

	result, err := service.GetScheduleByBranchID(ctx, "not-found")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
	assert.Nil(t, result)
}

// ============================================
// GetScheduleByID Tests
// ============================================

func TestGetScheduleByID_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	expected := &domain.BranchSchedule{ID: "sched-123"}
	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(expected, nil)

	result, err := service.GetScheduleByID(ctx, "sched-123")

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestGetScheduleByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockSchedRepo.On("GetScheduleByID", ctx, "not-found").Return(nil, domain.ErrScheduleNotFound)

	result, err := service.GetScheduleByID(ctx, "not-found")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// UpdateSchedule Tests
// ============================================

func TestUpdateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	existing := &domain.BranchSchedule{ID: "sched-123"}
	updated := domain.BranchSchedule{ID: "sched-123", Active: false}

	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(existing, nil)
	mockSchedRepo.On("UpdateSchedule", ctx, mockTx, updated).Return(nil)

	err := service.UpdateSchedule(ctx, mockTx, updated)

	assert.NoError(t, err)
}

func TestUpdateSchedule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockSchedRepo.On("GetScheduleByID", ctx, "not-found").Return(nil, nil)

	err := service.UpdateSchedule(ctx, mockTx, domain.BranchSchedule{ID: "not-found"})

	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

// ============================================
// DeleteSchedule Tests
// ============================================

func TestDeleteSchedule_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	existing := &domain.BranchSchedule{ID: "sched-123"}
	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(existing, nil)
	mockSchedRepo.On("DeleteSchedule", ctx, mockTx, "sched-123").Return(nil)

	err := service.DeleteSchedule(ctx, mockTx, "sched-123")

	assert.NoError(t, err)
}

func TestDeleteSchedule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockSchedRepo.On("GetScheduleByID", ctx, "not-found").Return(nil, nil)

	err := service.DeleteSchedule(ctx, mockTx, "not-found")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

// ============================================
// ActivateSchedule Tests
// ============================================

func TestActivateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	existing := &domain.BranchSchedule{ID: "sched-123", Active: false}
	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(existing, nil)
	mockSchedRepo.On("SetActive", ctx, mockTx, "sched-123", true).Return(nil)

	err := service.ActivateSchedule(ctx, mockTx, "sched-123")

	assert.NoError(t, err)
}

func TestActivateSchedule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	mockSchedRepo.On("GetScheduleByID", ctx, "not-found").Return(nil, nil)

	err := service.ActivateSchedule(ctx, mockTx, "not-found")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

// ============================================
// DeactivateSchedule Tests
// ============================================

func TestDeactivateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	existing := &domain.BranchSchedule{ID: "sched-123", Active: true}
	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(existing, nil)
	mockSchedRepo.On("SetActive", ctx, mockTx, "sched-123", false).Return(nil)

	err := service.DeactivateSchedule(ctx, mockTx, "sched-123")

	assert.NoError(t, err)
}

func TestDeactivateSchedule_Error(t *testing.T) {
	ctx := context.Background()
	mockSchedRepo := new(MockScheduleRepo)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)

	service := NewScheduleService(mockSchedRepo, mockBranchRepo)

	existing := &domain.BranchSchedule{ID: "sched-123"}
	dbError := errors.New("database error")

	mockSchedRepo.On("GetScheduleByID", ctx, "sched-123").Return(existing, nil)
	mockSchedRepo.On("SetActive", ctx, mockTx, "sched-123", false).Return(dbError)

	err := service.DeactivateSchedule(ctx, mockTx, "sched-123")

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
}
