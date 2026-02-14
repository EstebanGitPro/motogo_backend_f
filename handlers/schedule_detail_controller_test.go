package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateScheduleDetail Tests
// ============================================

func TestCreateScheduleDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)
	dayOfWeek := 1

	// GetScheduleByBranchID (with ownership) → branchService.GetBranchByID + scheduleService.GetScheduleByBranchID
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true}, nil)

	// CreateDetail: branchService.GetBranchByID (already mocked), BeginTx, CreateDetail, Commit
	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("CreateDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).
		Return(&domain.ScheduleDetail{
			ID:         detailUUID,
			ScheduleID: schedUUID,
			DayOfWeek:  &dayOfWeek,
			IsClosed:   false,
			Active:     true,
		}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	openTime := "08:00"
	closeTime := "18:00"
	body, _ := json.Marshal(map[string]interface{}{
		"day_of_week":  1,
		"opening_time": openTime,
		"closing_time": closeTime,
		"is_closed":    false,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules/details", h.CreateScheduleDetail(detailInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules/details", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateScheduleDetail_InvalidBranchID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	detailInteractor := interactor.NewScheduleDetailInteractor(new(mocks.MockScheduleDetailService), new(mocks.MockScheduleService), new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules/details", h.CreateScheduleDetail(detailInteractor, schedInteractor))

	body, _ := json.Marshal(map[string]interface{}{
		"day_of_week":  1,
		"opening_time": "08:00",
		"closing_time": "18:00",
		"is_closed":    false,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/BAD_ID/schedules/details", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCreateScheduleDetail_ScheduleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	detailInteractor := interactor.NewScheduleDetailInteractor(new(mocks.MockScheduleDetailService), mockSchedSvc, mockBranchSvc)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	body, _ := json.Marshal(map[string]interface{}{
		"day_of_week": 1, "opening_time": "08:00", "closing_time": "18:00", "is_closed": false,
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules/details", h.CreateScheduleDetail(detailInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules/details", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// ListScheduleDetails Tests
// ============================================

func TestListScheduleDetails_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, mockSchedSvc, new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)
	dayOfWeek := 1

	// GetScheduleByBranchIDPublic: scheduleService.GetScheduleByBranchID
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true}, nil)

	// ListDetails: detailService.GetDetailsByScheduleID
	mockDetailSvc.On("GetDetailsByScheduleID", mock.Anything, schedUUID).
		Return([]domain.ScheduleDetail{
			{ID: detailUUID, ScheduleID: schedUUID, DayOfWeek: &dayOfWeek, Active: true},
		}, nil)

	router := gin.New()
	router.GET("/branches/:id/schedules/details", h.ListScheduleDetails(detailInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules/details", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

func TestListScheduleDetails_ScheduleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	detailInteractor := interactor.NewScheduleDetailInteractor(new(mocks.MockScheduleDetailService), mockSchedSvc, new(mocks.MockBranchService))
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.GET("/branches/:id/schedules/details", h.ListScheduleDetails(detailInteractor, schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules/details", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// UpdateScheduleDetail Tests
// ============================================

func TestUpdateScheduleDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedDetailID, _ := encoder.Encode(detailUUID)
	dayOfWeek := 1

	// UpdateDetail: GetDetailByID → GetScheduleByID → GetBranchByID → BeginTx → UpdateDetail → Commit
	mockDetailSvc.On("GetDetailByID", mock.Anything, detailUUID).
		Return(&domain.ScheduleDetail{ID: detailUUID, ScheduleID: schedUUID, DayOfWeek: &dayOfWeek}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("UpdateDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	openTime := "09:00"
	closeTime := "17:00"
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
	router.PUT("/schedule-details/:id", h.UpdateScheduleDetail(detailInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-details/"+encodedDetailID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateScheduleDetail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	encodedDetailID, _ := encoder.Encode(detailUUID)

	mockDetailSvc.On("GetDetailByID", mock.Anything, detailUUID).
		Return(nil, domain.ErrScheduleDetailNotFound)

	body, _ := json.Marshal(map[string]interface{}{"is_closed": true})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/schedule-details/:id", h.UpdateScheduleDetail(detailInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedule-details/"+encodedDetailID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// DeleteScheduleDetail Tests
// ============================================

func TestDeleteScheduleDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedDetailID, _ := encoder.Encode(detailUUID)
	dayOfWeek := 1

	// DeleteDetail: GetDetailByID → GetScheduleByID → GetBranchByID → BeginTx → DeleteDetail → Commit
	mockDetailSvc.On("GetDetailByID", mock.Anything, detailUUID).
		Return(&domain.ScheduleDetail{ID: detailUUID, ScheduleID: schedUUID, DayOfWeek: &dayOfWeek}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)

	mockTx := new(mocks.MockTx)
	mockDetailSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailSvc.On("DeleteDetail", mock.Anything, mockTx, detailUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/schedule-details/:id", h.DeleteScheduleDetail(detailInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedule-details/"+encodedDetailID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteScheduleDetail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockDetailSvc := new(mocks.MockScheduleDetailService)
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailSvc, new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	detailUUID := "d1234567-89ab-cdef-0123-456789abcdef"
	encodedDetailID, _ := encoder.Encode(detailUUID)

	mockDetailSvc.On("GetDetailByID", mock.Anything, detailUUID).
		Return(nil, domain.ErrScheduleDetailNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/schedule-details/:id", h.DeleteScheduleDetail(detailInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedule-details/"+encodedDetailID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteScheduleDetail_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	detailInteractor := interactor.NewScheduleDetailInteractor(new(mocks.MockScheduleDetailService), new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/schedule-details/:id", h.DeleteScheduleDetail(detailInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedule-details/INVALID_ID", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}
