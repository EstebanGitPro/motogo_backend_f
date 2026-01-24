package interactor_test

import (
	"context"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateException Tests (HU20)
// ============================================

func TestCreateException_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	scheduleID := "schedule-123"
	representativeID := "rep-123"
	startDate := time.Now().AddDate(0, 0, 7) // 7 días en el futuro

	// Branch pertenece al representante
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}, nil)

	exception := domain.ScheduleDetail{
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
		ExceptionEndDate:   &startDate,
		IsClosed:           true,
	}

	createdEx := &domain.ScheduleDetail{
		ID:                 "exception-123",
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
		ExceptionEndDate:   &startDate,
		IsClosed:           true,
		Active:             true,
	}

	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("CreateException", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(createdEx, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := exceptionInteractor.CreateException(ctx, exception, representativeID, branchID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "exception-123", result.ID)
	assert.True(t, result.Active)

	mockDetailService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestCreateException_BranchNotOwned(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	representativeID := "rep-other" // Diferente al dueño
	startDate := time.Now().AddDate(0, 0, 7)

	// Branch pertenece a otro representante
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}, nil)

	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &startDate,
		IsClosed:           true,
	}

	// Act
	result, err := exceptionInteractor.CreateException(ctx, exception, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateException_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "non-existent"
	representativeID := "rep-123"
	startDate := time.Now().AddDate(0, 0, 7)

	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	exception := domain.ScheduleDetail{
		ExceptionStartDate: &startDate,
		IsClosed:           true,
	}

	// Act
	result, err := exceptionInteractor.CreateException(ctx, exception, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// ListExceptions Tests (HU23)
// ============================================

func TestListExceptions_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	scheduleID := "schedule-123"
	startDate1 := time.Now().AddDate(0, 0, 7)
	startDate2 := time.Now().AddDate(0, 0, 14)

	expectedExceptions := []domain.ScheduleDetail{
		{ID: "ex-1", ScheduleID: scheduleID, ExceptionStartDate: &startDate1, IsClosed: true, Active: true},
		{ID: "ex-2", ScheduleID: scheduleID, ExceptionStartDate: &startDate2, IsClosed: false, Active: true},
	}

	mockDetailService.On("GetExceptionsByScheduleID", ctx, scheduleID).Return(expectedExceptions, nil)

	// Act
	result, err := exceptionInteractor.ListExceptions(ctx, scheduleID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockDetailService.AssertExpectations(t)
}

func TestListExceptions_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	scheduleID := "schedule-no-exceptions"

	mockDetailService.On("GetExceptionsByScheduleID", ctx, scheduleID).Return([]domain.ScheduleDetail{}, nil)

	// Act
	result, err := exceptionInteractor.ListExceptions(ctx, scheduleID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockDetailService.AssertExpectations(t)
}

// ============================================
// DeleteException Tests (HU22)
// ============================================

func TestDeleteException_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	exceptionID := "exception-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"
	startDate := time.Now().AddDate(0, 0, 7)

	// Get existing exception
	mockDetailService.On("GetExceptionByID", ctx, exceptionID).Return(&domain.ScheduleDetail{
		ID:                 exceptionID,
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
	}, nil)

	// Get schedule
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(&domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}, nil)

	// Get branch and verify ownership
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}, nil)

	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("DeleteException", ctx, mockTx, exceptionID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := exceptionInteractor.DeleteException(ctx, exceptionID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteException_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	exceptionID := "non-existent"
	representativeID := "rep-123"

	mockDetailService.On("GetExceptionByID", ctx, exceptionID).Return(nil, domain.ErrScheduleExceptionNotFound)

	// Act
	err := exceptionInteractor.DeleteException(ctx, exceptionID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleExceptionNotFound, err)

	mockDetailService.AssertExpectations(t)
}

func TestDeleteException_NotOwned(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	exceptionID := "exception-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other" // Different from owner
	startDate := time.Now().AddDate(0, 0, 7)

	mockDetailService.On("GetExceptionByID", ctx, exceptionID).Return(&domain.ScheduleDetail{
		ID:                 exceptionID,
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
	}, nil)

	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(&domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}, nil)

	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner", // Original owner
	}, nil)

	// Act
	err := exceptionInteractor.DeleteException(ctx, exceptionID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// ActivateException Tests (HU24)
// ============================================

func TestActivateException_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	exceptionID := "exception-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"
	startDate := time.Now().AddDate(0, 0, 7)

	mockDetailService.On("GetExceptionByID", ctx, exceptionID).Return(&domain.ScheduleDetail{
		ID:                 exceptionID,
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
		Active:             false, // Currently inactive
	}, nil)

	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(&domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}, nil)

	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}, nil)

	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("UpdateException", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := exceptionInteractor.ActivateException(ctx, exceptionID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

// ============================================
// DeactivateException Tests (HU25)
// ============================================

func TestDeactivateException_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	exceptionInteractor := interactor.NewScheduleExceptionInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	exceptionID := "exception-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"
	startDate := time.Now().AddDate(0, 0, 7)

	mockDetailService.On("GetExceptionByID", ctx, exceptionID).Return(&domain.ScheduleDetail{
		ID:                 exceptionID,
		ScheduleID:         scheduleID,
		ExceptionStartDate: &startDate,
		Active:             true, // Currently active
	}, nil)

	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(&domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}, nil)

	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}, nil)

	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("UpdateException", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := exceptionInteractor.DeactivateException(ctx, exceptionID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}
