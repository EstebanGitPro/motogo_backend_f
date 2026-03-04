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

// TestCreateScheduleDetail_Integration_Success validates the full HTTP handler pipeline
// for the success path of CreateScheduleDetail (HU6).
//
// It exercises: auth → branch ID decoding → get schedule by branch →
// ScheduleDetailInteractor.CreateDetail → ID encoding → HATEOAS links → 201 response.
func TestCreateScheduleDetail_Integration_Success(t *testing.T) {
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
	detailInteractor := interactor.NewScheduleDetailInteractor(mockDetailService, mockScheduleService, mockBranchService)

	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	// Test data
	ownerID := "a1111111-1111-4000-8000-111111111111"
	branchID := "a2222222-2222-4000-8000-222222222222"
	scheduleID := "a5555555-5555-4000-8000-555555555555"
	detailID := "a6666666-6666-4000-8000-666666666666"
	dayOfWeek := 1 // Monday
	openingTime := "08:00"
	closingTime := "17:00"

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

	// Mock: CreateDetail (detail interactor)
	mockDetailService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockDetailService.On("CreateDetail", mock.Anything, mockTx, mock.AnythingOfType("domain.ScheduleDetail")).
		Return(&domain.ScheduleDetail{
			ID:          detailID,
			ScheduleID:  scheduleID,
			EntryType:   domain.EntryTypeRegular,
			DayOfWeek:   &dayOfWeek,
			OpeningTime: &openingTime,
			ClosingTime: &closingTime,
			IsClosed:    false,
			Active:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	// Request body
	reqBody := map[string]interface{}{
		"day_of_week":  1,
		"opening_time": "08:00",
		"closing_time": "17:00",
		"is_closed":    false,
	}
	body, _ := json.Marshal(reqBody)

	// Setup router
	router := gin.New()
	router.POST("/branches/:id/schedules/details", func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID})
	}, h.CreateScheduleDetail(detailInteractor, scheduleInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules/details", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_HD_CREATE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, float64(1), data["day_of_week"])
	assert.Equal(t, "Lunes", data["day_name"])
	assert.Equal(t, "08:00", data["opening_time"])
	assert.Equal(t, "17:00", data["closing_time"])
	assert.NotEmpty(t, data["_links"])

	mockDetailService.AssertExpectations(t)
	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
