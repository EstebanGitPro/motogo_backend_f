package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateScheduleException Tests
// ============================================

func TestCreateScheduleException_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)
	futureDate := time.Now().Add(48 * time.Hour)

	// GetScheduleByBranchID (with ownership) calls branchService.GetBranchByID + scheduleService.GetScheduleByBranchID
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true}, nil)

	// CreateException calls
	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("CheckExceptionDateConflict", mock.Anything, schedUUID, "", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, nil)
	mockDetailSvc.On("CreateException", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).
		Return(&domain.ScheduleDetail{
			ID:                 excUUID,
			ScheduleID:         schedUUID,
			ExceptionStartDate: &futureDate,
			ExceptionEndDate:   &futureDate,
			IsClosed:           true,
			Active:             true,
		}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"exception_start_date": futureDate.Format(constants.DateFormat),
		"is_closed":            true,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules/exceptions", h.CreateScheduleException(excInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules/exceptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateScheduleException_InvalidBranchID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	excInteractor := interactor.NewScheduleExceptionInteractor(new(mocks.MockScheduleDetailService), new(mocks.MockScheduleService), new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules/exceptions", h.CreateScheduleException(excInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/INVALID/schedules/exceptions", bytes.NewBuffer([]byte(`{"exception_start_date":"2026-12-25","is_closed":true}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// ListScheduleExceptions Tests
// ============================================

func TestListScheduleExceptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)
	futureDate := time.Now().Add(48 * time.Hour)

	// GetScheduleByBranchIDPublic calls scheduleService.GetScheduleByBranchID directly
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true}, nil)

	// ListExceptions
	mockDetailSvc.On("GetExceptionsByScheduleID", mock.Anything, schedUUID).
		Return([]domain.ScheduleDetail{
			{ID: excUUID, ScheduleID: schedUUID, ExceptionStartDate: &futureDate, ExceptionEndDate: &futureDate, IsClosed: true, Active: true},
		}, nil)

	router := gin.New()
	router.GET("/branches/:id/schedules/exceptions", h.ListScheduleExceptions(excInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules/exceptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

func TestListScheduleExceptions_ScheduleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	excInteractor := interactor.NewScheduleExceptionInteractor(new(mocks.MockScheduleDetailService), mockSchedSvc, new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.GET("/branches/:id/schedules/exceptions", h.ListScheduleExceptions(excInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules/exceptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// UpdateScheduleException Tests
// ============================================

func TestUpdateScheduleException_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedExcID, _ := encoder.Encode(excUUID)
	futureDate := time.Now().Add(48 * time.Hour)

	// UpdateException internals: GetExceptionByID → GetScheduleByID → branchService.GetBranchByID → ownership → BeginTx → ValidateTimeRange → UpdateException → Commit
	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(&domain.ScheduleDetail{ID: excUUID, ScheduleID: schedUUID, ExceptionStartDate: &futureDate, ExceptionEndDate: &futureDate}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("ValidateTimeRange", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	mockDetailSvc.On("UpdateException", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	openTime := "08:00"
	closeTime := "18:00"
	body, _ := json.Marshal(map[string]interface{}{
		"opening_time": openTime,
		"closing_time": closeTime,
		"is_closed":    false,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-exceptions/:id", h.UpdateScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-exceptions/"+encodedExcID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateScheduleException_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	encodedExcID, _ := encoder.Encode(excUUID)

	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(nil, domain.ErrScheduleExceptionNotFound)

	body, _ := json.Marshal(map[string]interface{}{"is_closed": true})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-exceptions/:id", h.UpdateScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-exceptions/"+encodedExcID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// DeleteScheduleException Tests
// ============================================

func TestDeleteScheduleException_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedExcID, _ := encoder.Encode(excUUID)

	// DeleteException: GetExceptionByID → GetScheduleByID → GetBranchByID → BeginTx → DeleteException → Commit
	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(&domain.ScheduleDetail{ID: excUUID, ScheduleID: schedUUID}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("DeleteException", mock.Anything, mockTx, excUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/schedule-exceptions/:id", h.DeleteScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedule-exceptions/"+encodedExcID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteScheduleException_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	encodedExcID, _ := encoder.Encode(excUUID)

	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(nil, domain.ErrScheduleExceptionNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/schedule-exceptions/:id", h.DeleteScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedule-exceptions/"+encodedExcID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// ActivateScheduleException Tests
// ============================================

func TestActivateScheduleException_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedExcID, _ := encoder.Encode(excUUID)

	// setExceptionActive: GetExceptionByID → GetScheduleByID → GetBranchByID → BeginTx → SetExceptionActive → Commit
	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(&domain.ScheduleDetail{ID: excUUID, ScheduleID: schedUUID, Active: false}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("SetExceptionActive", mock.Anything, mockTx, excUUID, true).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-exceptions/:id/activate", h.ActivateScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-exceptions/"+encodedExcID+"/activate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActivateScheduleException_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	excInteractor := interactor.NewScheduleExceptionInteractor(new(mocks.MockScheduleDetailService), new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-exceptions/:id/activate", h.ActivateScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-exceptions/BAD_ID/activate", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// DeactivateScheduleException Tests
// ============================================

func TestDeactivateScheduleException_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	excInteractor := interactor.NewScheduleExceptionInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	excUUID := "c1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedExcID, _ := encoder.Encode(excUUID)

	mockDetailSvc.On("GetExceptionByID", mock.Anything, excUUID).
		Return(&domain.ScheduleDetail{ID: excUUID, ScheduleID: schedUUID, Active: true}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("SetExceptionActive", mock.Anything, mockTx, excUUID, false).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-exceptions/:id/deactivate", h.DeactivateScheduleException(excInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-exceptions/"+encodedExcID+"/deactivate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
