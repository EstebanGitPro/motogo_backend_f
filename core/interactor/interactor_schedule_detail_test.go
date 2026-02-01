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

// Helper functions for pointer creation
func dayPtr(i int) *int {
	return &i
}

func timePtr(s string) *string {
	return &s
}

// ============================================
// CreateDetail Tests (HU6)
// ============================================

func TestCreateDetail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	representativeID := "rep-123"
	scheduleID := "schedule-123"

	detail := domain.ScheduleDetail{
		ScheduleID:  scheduleID,
		DayOfWeek:   dayPtr(1), // Monday
		OpeningTime: timePtr("08:00"),
		ClosingTime: timePtr("17:00"),
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	createdDetail := &domain.ScheduleDetail{
		ID:          "detail-123",
		ScheduleID:  scheduleID,
		DayOfWeek:   dayPtr(1),
		OpeningTime: timePtr("08:00"),
		ClosingTime: timePtr("17:00"),
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("CreateDetail", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(createdDetail, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := detailInteractor.CreateDetail(ctx, detail, representativeID, branchID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "detail-123", result.ID)

	mockBranchService.AssertExpectations(t)
	mockDetailService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestCreateDetail_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "non-existent"
	representativeID := "rep-123"

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  dayPtr(1),
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	result, err := detailInteractor.CreateDetail(ctx, detail, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateDetail_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	representativeID := "rep-other" // Different from owner

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  dayPtr(1),
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	result, err := detailInteractor.CreateDetail(ctx, detail, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateDetail_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	representativeID := "rep-123"

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  dayPtr(1),
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(nil, txError)

	// Act
	result, err := detailInteractor.CreateDetail(ctx, detail, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockBranchService.AssertExpectations(t)
	mockDetailService.AssertExpectations(t)
}

func TestCreateDetail_CreateError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	branchID := "branch-123"
	representativeID := "rep-123"

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  dayPtr(1),
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	createError := errors.New("create failed")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("CreateDetail", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil, createError)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := detailInteractor.CreateDetail(ctx, detail, representativeID, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockTx.AssertCalled(t, "Rollback")
	mockBranchService.AssertExpectations(t)
	mockDetailService.AssertExpectations(t)
}

// ============================================
// ListDetails Tests (HU9)
// ============================================

func TestListDetails_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	scheduleID := "schedule-123"

	expectedDetails := []domain.ScheduleDetail{
		{ID: "detail-1", ScheduleID: scheduleID, DayOfWeek: dayPtr(1)},
		{ID: "detail-2", ScheduleID: scheduleID, DayOfWeek: dayPtr(2)},
	}

	// Mock expectations
	mockDetailService.On("GetDetailsByScheduleID", ctx, scheduleID).Return(expectedDetails, nil)

	// Act
	result, err := detailInteractor.ListDetails(ctx, scheduleID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockDetailService.AssertExpectations(t)
}

func TestListDetails_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	scheduleID := "schedule-no-details"

	// Mock expectations
	mockDetailService.On("GetDetailsByScheduleID", ctx, scheduleID).Return([]domain.ScheduleDetail{}, nil)

	// Act
	result, err := detailInteractor.ListDetails(ctx, scheduleID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockDetailService.AssertExpectations(t)
}

func TestListDetails_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	scheduleID := "schedule-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockDetailService.On("GetDetailsByScheduleID", ctx, scheduleID).Return(nil, dbError)

	// Act
	result, err := detailInteractor.ListDetails(ctx, scheduleID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDetailService.AssertExpectations(t)
}

// ============================================
// GetDetail Tests
// ============================================

func TestGetDetail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"

	expectedDetail := &domain.ScheduleDetail{
		ID:          detailID,
		ScheduleID:  "schedule-123",
		DayOfWeek:   dayPtr(1),
		OpeningTime: timePtr("08:00"),
		ClosingTime: timePtr("17:00"),
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(expectedDetail, nil)

	// Act
	result, err := detailInteractor.GetDetail(ctx, detailID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, detailID, result.ID)

	mockDetailService.AssertExpectations(t)
}

func TestGetDetail_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "non-existent"

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(nil, domain.ErrScheduleDetailNotFound)

	// Act
	result, err := detailInteractor.GetDetail(ctx, detailID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDetailService.AssertExpectations(t)
}

// ============================================
// UpdateDetail Tests (HU7)
// ============================================

func TestUpdateDetail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	detail := domain.ScheduleDetail{
		ID:          detailID,
		ScheduleID:  scheduleID,
		DayOfWeek:   dayPtr(1),
		OpeningTime: timePtr("09:00"),
		ClosingTime: timePtr("18:00"),
	}

	existingDetail := &domain.ScheduleDetail{
		ID:         detailID,
		ScheduleID: scheduleID,
	}

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(existingDetail, nil)
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("UpdateDetail", ctx, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := detailInteractor.UpdateDetail(ctx, detail, representativeID)

	// Assert
	assert.NoError(t, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateDetail_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "non-existent"
	representativeID := "rep-123"

	detail := domain.ScheduleDetail{
		ID: detailID,
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(nil, domain.ErrScheduleDetailNotFound)

	// Act
	err := detailInteractor.UpdateDetail(ctx, detail, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailNotFound, err)

	mockDetailService.AssertExpectations(t)
}

func TestUpdateDetail_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other" // Different from owner

	detail := domain.ScheduleDetail{
		ID: detailID,
	}

	existingDetail := &domain.ScheduleDetail{
		ID:         detailID,
		ScheduleID: scheduleID,
	}

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(existingDetail, nil)
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := detailInteractor.UpdateDetail(ctx, detail, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeleteDetail Tests (HU8)
// ============================================

func TestDeleteDetail_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	existingDetail := &domain.ScheduleDetail{
		ID:         detailID,
		ScheduleID: scheduleID,
	}

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(existingDetail, nil)
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("DeleteDetail", ctx, mockTx, detailID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := detailInteractor.DeleteDetail(ctx, detailID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteDetail_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(nil, domain.ErrScheduleDetailNotFound)

	// Act
	err := detailInteractor.DeleteDetail(ctx, detailID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailNotFound, err)

	mockDetailService.AssertExpectations(t)
}

func TestDeleteDetail_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-other" // Different from owner

	existingDetail := &domain.ScheduleDetail{
		ID:         detailID,
		ScheduleID: scheduleID,
	}

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(existingDetail, nil)
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := detailInteractor.DeleteDetail(ctx, detailID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

func TestDeleteDetail_DeleteError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	detailInteractor := interactor.NewScheduleDetailInteractor(
		mockDetailService,
		mockScheduleService,
		mockBranchService,
	)

	detailID := "detail-123"
	scheduleID := "schedule-123"
	branchID := "branch-123"
	representativeID := "rep-123"

	existingDetail := &domain.ScheduleDetail{
		ID:         detailID,
		ScheduleID: scheduleID,
	}

	existingSchedule := &domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
	}

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: representativeID,
	}

	deleteError := errors.New("delete failed")

	// Mock expectations
	mockDetailService.On("GetDetailByID", ctx, detailID).Return(existingDetail, nil)
	mockScheduleService.On("GetScheduleByID", ctx, scheduleID).Return(existingSchedule, nil)
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockDetailService.On("BeginTx", ctx).Return(mockTx, nil)
	mockDetailService.On("DeleteDetail", ctx, mockTx, detailID).Return(deleteError)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := detailInteractor.DeleteDetail(ctx, detailID, representativeID)

	// Assert
	assert.Error(t, err)

	mockTx.AssertCalled(t, "Rollback")
	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}
