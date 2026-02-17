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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateBranchSchedule Tests
// ============================================

func TestCreateBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	scheduleUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("CreateSchedule", mock.Anything, mockTx, branchUUID).
		Return(&domain.BranchSchedule{ID: scheduleUUID, BranchID: branchUUID, Active: false, StartDate: time.Now()}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules", h.CreateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

func TestCreateBranchSchedule_InvalidBranchID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules", h.CreateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/INVALID/schedules", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCreateBranchSchedule_AlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("CreateSchedule", mock.Anything, mockTx, branchUUID).
		Return(nil, domain.ErrScheduleAlreadyExists)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.POST("/branches/:id/schedules", h.CreateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// ============================================
// GetBranchSchedule Tests
// ============================================

func TestGetBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true, StartDate: time.Now()}, nil)

	router := gin.New()
	router.GET("/branches/:id/schedules", h.GetBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

func TestGetBranchSchedule_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.GET("/branches/:id/schedules", h.GetBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBranchSchedule_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/branches/:id/schedules", h.GetBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/branches/BAD_ID/schedules", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// DeleteBranchSchedule Tests
// ============================================

func TestDeleteBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("DeleteSchedule", mock.Anything, mockTx, schedUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/branches/:id/schedules", h.DeleteBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBranchSchedule_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.DELETE("/branches/:id/schedules", h.DeleteBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// ActivateBranchSchedule Tests
// ============================================

func TestActivateBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: false, StartDate: time.Now()}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: false, StartDate: time.Now()}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("ActivateSchedule", mock.Anything, mockTx, schedUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules/activate", h.ActivateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID+"/schedules/activate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActivateBranchSchedule_ScheduleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules/activate", h.ActivateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID+"/schedules/activate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================
// DeactivateBranchSchedule Tests
// ============================================

func TestDeactivateBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true, StartDate: time.Now()}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true, StartDate: time.Now()}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("DeactivateSchedule", mock.Anything, mockTx, schedUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules/deactivate", h.DeactivateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID+"/schedules/deactivate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeactivateBranchSchedule_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules/deactivate", h.DeactivateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/INVALID/schedules/deactivate", nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ============================================
// GetDaysOfWeek Tests
// ============================================

func TestGetDaysOfWeek_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.GET("/schedules/days", h.GetDaysOfWeek())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/schedules/days", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
}

// ============================================
// UpdateBranchSchedule Tests
// ============================================

func TestUpdateBranchSchedule_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	mockBranchSvc := new(mocks.MockBranchService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, mockBranchSvc)
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	schedUUID := "b1234567-89ab-cdef-0123-456789abcdef"
	repID := "rep-123"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockTx := new(mocks.MockTx)
	startDate := time.Now()

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true, StartDate: startDate}, nil)
	mockSchedSvc.On("GetScheduleByID", mock.Anything, schedUUID).
		Return(&domain.BranchSchedule{ID: schedUUID, BranchID: branchUUID, Active: true, StartDate: startDate}, nil)
	mockBranchSvc.On("GetBranchByID", mock.Anything, branchUUID).
		Return(&domain.Branch{ID: branchUUID, RepresentativeID: repID}, nil)
	mockSchedSvc.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockSchedSvc.On("UpdateSchedule", mock.Anything, mockTx, mock.AnythingOfType("domain.BranchSchedule")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	bodyJSON, _ := json.Marshal(map[string]interface{}{"active": true})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: repID, Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules", h.UpdateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID+"/schedules", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBranchSchedule_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	schedInteractor := interactor.NewScheduleInteractor(new(mocks.MockScheduleService), new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules", h.UpdateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/INVALID/schedules", bytes.NewBuffer([]byte(`{"active": true}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateBranchSchedule_ScheduleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	mockSchedSvc := new(mocks.MockScheduleService)
	schedInteractor := interactor.NewScheduleInteractor(mockSchedSvc, new(mocks.MockBranchService))
	h := handlers.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, encoder, responseHandler)

	branchUUID := "a1234567-89ab-cdef-0123-456789abcdef"
	encodedBranchID, _ := encoder.Encode(branchUUID)

	mockSchedSvc.On("GetScheduleByBranchID", mock.Anything, branchUUID).
		Return(nil, domain.ErrScheduleNotFound)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: "rep-1", Role: domain.Role("REPRESENTANTE")})
		c.Next()
	})
	router.PUT("/branches/:id/schedules", h.UpdateBranchSchedule(schedInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/branches/"+encodedBranchID+"/schedules", bytes.NewBuffer([]byte(`{"active": true}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
