package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateSchedule Tests (HU30)
// ============================================

func TestCreateSchedule_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		Name:             "Taller Norte",
		RepresentativeID: representativeID,
	}

	createdSchedule := &domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
		Active:   true,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("CreateSchedule", ctx, mockTx, branchID).Return(createdSchedule, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := scheduleInteractor.CreateSchedule(ctx, branchID, representativeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "schedule-123", result.ID)

	mockBranchService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestCreateSchedule_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	result, err := scheduleInteractor.CreateSchedule(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateSchedule_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-other" // Different from owner

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	result, err := scheduleInteractor.CreateSchedule(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateSchedule_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(nil, txError)

	// Act
	result, err := scheduleInteractor.CreateSchedule(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockBranchService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
}

func TestCreateSchedule_CreateError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	createError := errors.New("create failed")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("CreateSchedule", ctx, mockTx, branchID).Return(nil, createError)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := scheduleInteractor.CreateSchedule(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockTx.AssertCalled(t, "Rollback")
	mockBranchService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
}

// ============================================
// GetScheduleByBranchID Tests (HU32)
// ============================================

func TestGetScheduleByBranchID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	expectedSchedule := &domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
		Active:   true,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("GetScheduleByBranchID", ctx, branchID).Return(expectedSchedule, nil)

	// Act
	result, err := scheduleInteractor.GetScheduleByBranchID(ctx, branchID, representativeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "schedule-123", result.ID)

	mockBranchService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
}

func TestGetScheduleByBranchID_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	result, err := scheduleInteractor.GetScheduleByBranchID(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestGetScheduleByBranchID_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-other"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	result, err := scheduleInteractor.GetScheduleByBranchID(ctx, branchID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// GetScheduleByBranchIDPublic Tests
// ============================================

func TestGetScheduleByBranchIDPublic_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"

	expectedSchedule := &domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByBranchID", ctx, branchID).Return(expectedSchedule, nil)

	// Act
	result, err := scheduleInteractor.GetScheduleByBranchIDPublic(ctx, branchID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockScheduleService.AssertExpectations(t)
}

func TestGetScheduleByBranchIDPublic_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"

	// Mock expectations
	mockScheduleService.On("GetScheduleByBranchID", ctx, branchID).Return(nil, domain.ErrScheduleNotFound)

	// Act
	result, err := scheduleInteractor.GetScheduleByBranchIDPublic(ctx, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockScheduleService.AssertExpectations(t)
}

// ============================================
// GetScheduleByID Tests
// ============================================

func TestGetScheduleByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"

	expectedSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: "branch-123",
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(expectedSchedule, nil)

	// Act
	result, err := scheduleInteractor.GetScheduleByID(ctx, scheduleID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, scheduleID, result.ID)

	mockScheduleService.AssertExpectations(t)
}

func TestGetScheduleByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "non-existent"

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(nil, domain.ErrScheduleNotFound)

	// Act
	result, err := scheduleInteractor.GetScheduleByID(ctx, scheduleID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockScheduleService.AssertExpectations(t)
}

// ============================================
// UpdateSchedule Tests (HU31)
// ============================================

func TestUpdateSchedule_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-123"

	schedule := domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("UpdateSchedule", ctx, mockTx, mock.AnythingOfType("domain.BranchSchedule")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := scheduleInteractor.UpdateSchedule(ctx, schedule, representativeID)

	// Assert
	assert.NoError(t, err)

	mockBranchService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateSchedule_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "non-existent"
	representativeID := "rep-123"

	schedule := domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	err := scheduleInteractor.UpdateSchedule(ctx, schedule, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateSchedule_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	branchID := "branch-123"
	representativeID := "rep-other"

	schedule := domain.BranchSchedule{
		ID:       "schedule-123",
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := scheduleInteractor.UpdateSchedule(ctx, schedule, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeleteSchedule Tests (HU33)
// ============================================

func TestDeleteSchedule_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("DeleteSchedule", ctx, mockTx, scheduleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := scheduleInteractor.DeleteSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteSchedule_ScheduleNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(nil, domain.ErrScheduleNotFound)

	// Act
	err := scheduleInteractor.DeleteSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)

	mockScheduleService.AssertExpectations(t)
}

func TestDeleteSchedule_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := scheduleInteractor.DeleteSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// ActivateSchedule Tests (HU34)
// ============================================

func TestActivateSchedule_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
		Active:   false,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("ActivateSchedule", ctx, mockTx, scheduleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := scheduleInteractor.ActivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestActivateSchedule_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(nil, domain.ErrScheduleNotFound)

	// Act
	err := scheduleInteractor.ActivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleNotFound, err)

	mockScheduleService.AssertExpectations(t)
}

func TestActivateSchedule_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := scheduleInteractor.ActivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeactivateSchedule Tests (HU35)
// ============================================

func TestDeactivateSchedule_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
		Active:   true,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockScheduleService.On("BeginTx", ctx).Return(mockTx, nil)
	mockScheduleService.On("DeactivateSchedule", ctx, mockTx, scheduleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := scheduleInteractor.DeactivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeactivateSchedule_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(nil, domain.ErrScheduleNotFound)

	// Act
	err := scheduleInteractor.DeactivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)

	mockScheduleService.AssertExpectations(t)
}

func TestDeactivateSchedule_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other"

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := scheduleInteractor.DeactivateSchedule(ctx, scheduleID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}
