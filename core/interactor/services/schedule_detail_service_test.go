package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Helper Functions for Schedule Detail Tests
// ============================================

func setupScheduleDetailMocks() (*mocks.MockScheduleDetailRepository, *mocks.MockScheduleRepository) {
	return new(mocks.MockScheduleDetailRepository), new(mocks.MockScheduleRepository)
}

func intPtr(i int) *int {
	return &i
}

// ============================================
// NewScheduleDetailService Tests
// ============================================

func TestNewScheduleDetailService_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	// Act
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Assert
	assert.NotNil(t, service)
}

// ============================================
// BeginTx Tests for ScheduleDetailService
// ============================================

func TestScheduleDetailService_BeginTx_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detailRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	tx, err := service.BeginTx(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	detailRepo.AssertExpectations(t)
}

func TestScheduleDetailService_BeginTx_Error(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	detailRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("connection error"))

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	tx, err := service.BeginTx(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tx)
}

// ============================================
// ValidateTimeRange Tests
// ============================================

func TestValidateTimeRange_Valid(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("08:00", "17:00")

	// Assert
	assert.NoError(t, err)
}

func TestValidateTimeRange_Valid_WithSeconds(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("08:00:00", "17:00:00")

	// Assert
	assert.NoError(t, err)
}

func TestValidateTimeRange_InvalidOpeningFormat(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("invalid", "17:00")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailInvalidTime, err)
}

func TestValidateTimeRange_InvalidClosingFormat(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("08:00", "invalid")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailInvalidTime, err)
}

func TestValidateTimeRange_ClosingBeforeOpening(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("17:00", "08:00")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailInvalidTime, err)
}

func TestValidateTimeRange_SameTime(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.ValidateTimeRange("08:00", "08:00")

	// Assert
	assert.Error(t, err) // Same time is not valid (closing must be AFTER opening)
}

// ============================================
// GetDetailByID Tests
// ============================================

func TestGetDetailByID_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	expectedDetail := &domain.ScheduleDetail{
		ID:         "detail-123",
		ScheduleID: "schedule-456",
		DayOfWeek:  intPtr(1),
	}

	detailRepo.On("GetDetailByID", mock.Anything, "detail-123").Return(expectedDetail, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetDetailByID(context.Background(), "detail-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "detail-123", result.ID)
	detailRepo.AssertExpectations(t)
}

func TestGetDetailByID_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("GetDetailByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetDetailByID(context.Background(), "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleDetailNotFound, err)
}

// ============================================
// GetDetailsByScheduleID Tests
// ============================================

func TestGetDetailsByScheduleID_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	expectedDetails := []domain.ScheduleDetail{
		{ID: "detail-1", ScheduleID: "schedule-456", DayOfWeek: intPtr(1)},
		{ID: "detail-2", ScheduleID: "schedule-456", DayOfWeek: intPtr(2)},
	}

	detailRepo.On("GetDetailsByScheduleID", mock.Anything, "schedule-456").Return(expectedDetails, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetDetailsByScheduleID(context.Background(), "schedule-456")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	detailRepo.AssertExpectations(t)
}

func TestGetDetailsByScheduleID_Empty(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("GetDetailsByScheduleID", mock.Anything, "schedule-empty").Return([]domain.ScheduleDetail{}, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetDetailsByScheduleID(context.Background(), "schedule-empty")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// CreateDetail Tests
// ============================================

func TestCreateDetail_ScheduleNotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detail := domain.ScheduleDetail{
		ScheduleID: "non-existent",
		DayOfWeek:  intPtr(1),
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "non-existent").Return(nil, domain.ErrScheduleNotFound)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

func TestCreateDetail_InvalidDayOfWeek(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  intPtr(99), // Invalid day
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").
		Return(&domain.BranchSchedule{ID: "schedule-123"}, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleDetailInvalidDay, err)
}

func TestCreateDetail_DayAlreadyClosed(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	openingTime := "08:00"
	closingTime := "17:00"
	detail := domain.ScheduleDetail{
		ScheduleID:  "schedule-123",
		DayOfWeek:   intPtr(1),
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").
		Return(&domain.BranchSchedule{ID: "schedule-123"}, nil)
	detailRepo.On("CheckDayIsClosed", mock.Anything, "schedule-123", 1, "").Return(true, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleDetailDayAlreadyClosed, err)
}

func TestCreateDetail_Success_Closed(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detail := domain.ScheduleDetail{
		ScheduleID: "schedule-123",
		DayOfWeek:  intPtr(7), // Sunday
		IsClosed:   true,
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").
		Return(&domain.BranchSchedule{ID: "schedule-123"}, nil)
	detailRepo.On("CheckDayIsClosed", mock.Anything, "schedule-123", 7, "").Return(false, nil)
	detailRepo.On("CheckDayHasTimeSlots", mock.Anything, "schedule-123", 7, "").Return(false, nil)
	detailRepo.On("SaveScheduleDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.True(t, result.Active)
	detailRepo.AssertExpectations(t)
}

func TestCreateDetail_Success_WithTimeSlot(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	openingTime := "08:00"
	closingTime := "17:00"
	detail := domain.ScheduleDetail{
		ScheduleID:  "schedule-123",
		DayOfWeek:   intPtr(1), // Monday
		IsClosed:    false,
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").
		Return(&domain.BranchSchedule{ID: "schedule-123"}, nil)
	detailRepo.On("CheckDayIsClosed", mock.Anything, "schedule-123", 1, "").Return(false, nil)
	detailRepo.On("CheckTimeConflict", mock.Anything, "schedule-123", 1, "08:00", "17:00", "").Return(false, nil)
	detailRepo.On("SaveScheduleDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.EntryTypeRegular, result.EntryType)
	detailRepo.AssertExpectations(t)
}

func TestCreateDetail_TimeConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	openingTime := "08:00"
	closingTime := "12:00"
	detail := domain.ScheduleDetail{
		ScheduleID:  "schedule-123",
		DayOfWeek:   intPtr(1),
		IsClosed:    false,
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
	}

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").
		Return(&domain.BranchSchedule{ID: "schedule-123"}, nil)
	detailRepo.On("CheckDayIsClosed", mock.Anything, "schedule-123", 1, "").Return(false, nil)
	detailRepo.On("CheckTimeConflict", mock.Anything, "schedule-123", 1, "08:00", "12:00", "").Return(true, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.CreateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleDetailTimeConflict, err)
}

// ============================================
// UpdateDetail Tests
// ============================================

func TestUpdateDetail_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detail := domain.ScheduleDetail{
		ID:       "non-existent",
		IsClosed: true,
	}

	detailRepo.On("GetDetailByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.UpdateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailNotFound, err)
}

func TestUpdateDetail_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	openingTime := "09:00"
	closingTime := "18:00"
	existingDetail := &domain.ScheduleDetail{
		ID:         "detail-123",
		ScheduleID: "schedule-456",
		DayOfWeek:  intPtr(1),
	}

	detail := domain.ScheduleDetail{
		ID:          "detail-123",
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
		IsClosed:    false,
	}

	detailRepo.On("GetDetailByID", mock.Anything, "detail-123").Return(existingDetail, nil)
	detailRepo.On("CheckTimeConflict", mock.Anything, "schedule-456", 1, "09:00", "18:00", "detail-123").Return(false, nil)
	detailRepo.On("UpdateScheduleDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.UpdateDetail(context.Background(), mockTx, detail)

	// Assert
	assert.NoError(t, err)
	detailRepo.AssertExpectations(t)
}

// ============================================
// DeleteDetail Tests
// ============================================

func TestDeleteDetail_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detailRepo.On("GetDetailByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.DeleteDetail(context.Background(), mockTx, "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleDetailNotFound, err)
}

func TestDeleteDetail_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	existingDetail := &domain.ScheduleDetail{ID: "detail-123"}

	detailRepo.On("GetDetailByID", mock.Anything, "detail-123").Return(existingDetail, nil)
	detailRepo.On("DeleteScheduleDetail", mock.Anything, mockTx, "detail-123").Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.DeleteDetail(context.Background(), mockTx, "detail-123")

	// Assert
	assert.NoError(t, err)
	detailRepo.AssertExpectations(t)
}

// ============================================
// GetExceptionsByScheduleID Tests
// ============================================

func TestGetExceptionsByScheduleID_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	expectedExceptions := []domain.ScheduleDetail{
		{ID: "exc-1", EntryType: domain.EntryTypeException},
		{ID: "exc-2", EntryType: domain.EntryTypeException},
	}

	detailRepo.On("GetExceptionsByScheduleID", mock.Anything, "schedule-123").Return(expectedExceptions, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetExceptionsByScheduleID(context.Background(), "schedule-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============================================
// GetExceptionByID Tests
// ============================================

func TestGetExceptionByID_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	expectedException := &domain.ScheduleDetail{
		ID:        "exc-123",
		EntryType: domain.EntryTypeException,
	}

	detailRepo.On("GetExceptionByID", mock.Anything, "exc-123").Return(expectedException, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetExceptionByID(context.Background(), "exc-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "exc-123", result.ID)
}

func TestGetExceptionByID_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("GetExceptionByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	result, err := service.GetExceptionByID(context.Background(), "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleExceptionNotFound, err)
}

// ============================================
// UpdateException Tests
// ============================================

func TestUpdateException_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	exception := domain.ScheduleDetail{
		ID:       "non-existent",
		IsClosed: true,
	}

	detailRepo.On("GetExceptionByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.UpdateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleExceptionNotFound, err)
}

func TestUpdateException_Success_Closed(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	startDate := time.Now().AddDate(0, 0, 1)
	existingException := &domain.ScheduleDetail{
		ID:                 "exc-123",
		ScheduleID:         "schedule-456",
		ExceptionStartDate: &startDate,
		ExceptionEndDate:   &startDate,
	}

	exception := domain.ScheduleDetail{
		ID:       "exc-123",
		IsClosed: true,
	}

	detailRepo.On("GetExceptionByID", mock.Anything, "exc-123").Return(existingException, nil)
	detailRepo.On("UpdateScheduleDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.UpdateException(context.Background(), mockTx, exception)

	// Assert
	assert.NoError(t, err)
	detailRepo.AssertExpectations(t)
}

// ============================================
// DeleteException Tests
// ============================================

func TestDeleteException_NotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	detailRepo.On("GetExceptionByID", mock.Anything, "non-existent").Return(nil, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.DeleteException(context.Background(), mockTx, "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrScheduleExceptionNotFound, err)
}

func TestDeleteException_Success(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	existingException := &domain.ScheduleDetail{ID: "exc-123"}

	detailRepo.On("GetExceptionByID", mock.Anything, "exc-123").Return(existingException, nil)
	detailRepo.On("DeleteScheduleDetail", mock.Anything, mockTx, "exc-123").Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	err := service.DeleteException(context.Background(), mockTx, "exc-123")

	// Assert
	assert.NoError(t, err)
	detailRepo.AssertExpectations(t)
}

// ============================================
// CheckTimeConflict Tests
// ============================================

func TestCheckTimeConflict_NoConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("CheckTimeConflict", mock.Anything, "schedule-123", 1, "08:00", "12:00", "").Return(false, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	hasConflict, err := service.CheckTimeConflict(context.Background(), "schedule-123", 1, "08:00", "12:00", "")

	// Assert
	assert.NoError(t, err)
	assert.False(t, hasConflict)
}

func TestCheckTimeConflict_HasConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("CheckTimeConflict", mock.Anything, "schedule-123", 1, "08:00", "12:00", "").Return(true, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	hasConflict, err := service.CheckTimeConflict(context.Background(), "schedule-123", 1, "08:00", "12:00", "")

	// Assert
	assert.NoError(t, err)
	assert.True(t, hasConflict)
}

// ============================================
// CheckExceptionDateConflict Tests
// ============================================

func TestCheckExceptionDateConflict_NoConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("CheckExceptionDateConflict", mock.Anything, "schedule-123", "", "2026-02-01", "2026-02-03").Return(false, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	hasConflict, err := service.CheckExceptionDateConflict(context.Background(), "schedule-123", "", "2026-02-01", "2026-02-03")

	// Assert
	assert.NoError(t, err)
	assert.False(t, hasConflict)
}

func TestCheckExceptionDateConflict_HasConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("CheckExceptionDateConflict", mock.Anything, "schedule-123", "", "2026-02-01", "2026-02-03").Return(true, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	hasConflict, err := service.CheckExceptionDateConflict(context.Background(), "schedule-123", "", "2026-02-01", "2026-02-03")

	// Assert
	assert.NoError(t, err)
	assert.True(t, hasConflict)
}

func TestCheckExceptionDateConflict_ExcludeExisting(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	// When updating an exception, exclude itself from conflict check
	detailRepo.On("CheckExceptionDateConflict", mock.Anything, "schedule-123", "exc-456", "2026-02-01", "2026-02-03").Return(false, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	hasConflict, err := service.CheckExceptionDateConflict(context.Background(), "schedule-123", "exc-456", "2026-02-01", "2026-02-03")

	// Assert
	assert.NoError(t, err)
	assert.False(t, hasConflict)
}

func TestCheckExceptionDateConflict_Error(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()

	detailRepo.On("CheckExceptionDateConflict", mock.Anything, "schedule-123", "", "2026-02-01", "2026-02-03").Return(false, errors.New("database error"))

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Act
	_, err := service.CheckExceptionDateConflict(context.Background(), "schedule-123", "", "2026-02-01", "2026-02-03")

	// Assert
	assert.Error(t, err)
}

// ============================================
// CreateException Tests
// ============================================

func TestCreateException_ScheduleNotFound(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(nil, errors.New("not found"))

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Create exception with future date
	futureDate := time.Now().AddDate(0, 0, 7)
	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &futureDate,
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

func TestCreateException_StartDateNil(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	schedule := &domain.BranchSchedule{ID: "schedule-123", BranchID: "branch-123"}
	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(schedule, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: nil, // Missing start date
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleExceptionDatePast, err)
}

func TestCreateException_StartDateInPast(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	schedule := &domain.BranchSchedule{ID: "schedule-123", BranchID: "branch-123"}
	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(schedule, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Past date
	pastDate := time.Now().AddDate(0, 0, -7)
	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &pastDate,
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleExceptionDatePast, err)
}

func TestCreateException_DateConflict(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	schedule := &domain.BranchSchedule{ID: "schedule-123", BranchID: "branch-123"}
	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(schedule, nil)

	// Future date
	futureDate := time.Now().AddDate(0, 0, 7)
	existingDate := time.Now().AddDate(0, 0, 7) // Same date - conflict

	// Existing exception with same date
	existingExceptions := []domain.ScheduleDetail{
		{
			ID:                 "existing-exception-123",
			ScheduleID:         "schedule-123",
			ExceptionStartDate: &existingDate,
			ExceptionEndDate:   &existingDate,
			IsClosed:           true,
		},
	}

	detailRepo.On("GetExceptionsByScheduleIDForUpdate", mock.Anything, mockTx, "schedule-123").Return(existingExceptions, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &futureDate,
		IsClosed:           true,
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleExceptionDateConflict, err)
}

func TestCreateException_InvalidTimeWhenNotClosed(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	schedule := &domain.BranchSchedule{ID: "schedule-123", BranchID: "branch-123"}
	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(schedule, nil)

	futureDate := time.Now().AddDate(0, 0, 7)
	detailRepo.On("GetExceptionsByScheduleIDForUpdate", mock.Anything, mockTx, "schedule-123").Return([]domain.ScheduleDetail{}, nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	// Exception that is not closed but has no opening/closing times
	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &futureDate,
		IsClosed:           false,
		OpeningTime:        nil, // Missing times
		ClosingTime:        nil,
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrScheduleExceptionInvalidTime, err)
}

func TestCreateException_ClosedSuccess(t *testing.T) {
	// Arrange
	detailRepo, scheduleRepo := setupScheduleDetailMocks()
	mockTx := new(mocks.MockTx)

	schedule := &domain.BranchSchedule{ID: "schedule-123", BranchID: "branch-123"}
	scheduleRepo.On("GetScheduleByID", mock.Anything, "schedule-123").Return(schedule, nil)

	futureDate := time.Now().AddDate(0, 0, 7)
	detailRepo.On("GetExceptionsByScheduleIDForUpdate", mock.Anything, mockTx, "schedule-123").Return([]domain.ScheduleDetail{}, nil)
	detailRepo.On("CheckExceptionIsRedundant", mock.Anything, "schedule-123", mock.Anything).Return(false, nil)
	detailRepo.On("SaveScheduleDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)

	service := services.NewScheduleDetailService(detailRepo, scheduleRepo)

	exception := domain.ScheduleDetail{
		ScheduleID:         "schedule-123",
		ExceptionStartDate: &futureDate,
		IsClosed:           true,
	}

	// Act
	result, err := service.CreateException(context.Background(), mockTx, exception)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.True(t, result.IsClosed)
	detailRepo.AssertExpectations(t)
}
