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

// TestCreateScheduleException_Integration_Success validates the full HTTP handler pipeline
// for the success path of CreateScheduleException (HU20).
//
// It exercises: auth → branch ID decoding → get schedule by branch → date parsing →
// ScheduleExceptionInteractor.CreateException → ID encoding → HATEOAS links → 201 response.
func TestCreateScheduleException_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies
	mockDetailService := new(mocks.MockScheduleDetailService)
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)
	exceptionInteractor := interactor.NewScheduleExceptionInteractor(mockDetailService, mockScheduleService, mockBranchService)

	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	// Test data
	ownerID := "a1111111-1111-4000-8000-111111111111"
	branchID := "a2222222-2222-4000-8000-222222222222"
	scheduleID := "a5555555-5555-4000-8000-555555555555"
	exceptionID := "a7777777-7777-4000-8000-777777777777"
	openingTime := "09:00"
	closingTime := "14:00"

	encodedBranchID, _ := encoder.Encode(branchID)

	// Mock: GetScheduleByBranchID (ownership check + schedule lookup)
	mockBranchService.On("GetBranchByID", mock.Anything, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: ownerID,
	}, nil)
	mockScheduleService.On("GetScheduleByBranchID", mock.Anything, branchID).Return(&domain.BranchSchedule{
		ID:       scheduleID,
		BranchID: branchID,
		Active:   true,
	}, nil)

	// Mock: CreateException (exception interactor)
	exceptionStartDate, _ := time.ParseInLocation("2006-01-02", "2026-12-25", time.Local)
	mockDetailService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailService.On("CreateException", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).
		Return(&domain.ScheduleDetail{
			ID:                 exceptionID,
			ScheduleID:         scheduleID,
			EntryType:          domain.EntryTypeException,
			ExceptionStartDate: &exceptionStartDate,
			ExceptionEndDate:   &exceptionStartDate,
			OpeningTime:        &openingTime,
			ClosingTime:        &closingTime,
			IsClosed:           false,
			Active:             true,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}, nil)
	mockTx.On("Commit").Return(nil)

	// Request body
	reqBody := map[string]interface{}{
		"exception_start_date": "2026-12-25",
		"opening_time":         "09:00",
		"closing_time":         "14:00",
		"is_closed":            false,
	}
	body, _ := json.Marshal(reqBody)

	// Setup router
	router := gin.New()
	router.POST("/branches/:id/schedules/exceptions", func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID})
	}, h.CreateScheduleException(exceptionInteractor, scheduleInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules/exceptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_EH_CREATE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["schedule_id"])
	assert.NotEmpty(t, data["_links"])

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
