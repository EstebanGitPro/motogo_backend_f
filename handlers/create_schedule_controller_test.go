package handlers_test

import (
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

// TestCreateBranchSchedule_Integration_Success validates the full HTTP handler pipeline
// for the success path of CreateBranchSchedule (HU30).
//
// It exercises: auth → branch ID decoding → ownership check → ScheduleInteractor.CreateSchedule →
// schedule ID encoding → HATEOAS links → 201 response.
func TestCreateBranchSchedule_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies for ScheduleInteractor
	mockScheduleService := new(mocks.MockScheduleService)
	mockBranchService := new(mocks.MockBranchService)
	mockTx := new(mocks.MockTx)

	scheduleInteractor := interactor.NewScheduleInteractor(mockScheduleService, mockBranchService)

	// Handler only needs IDEncoder + Response
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	// Test data
	ownerID := "a1111111-1111-4000-8000-111111111111"
	branchID := "a2222222-2222-4000-8000-222222222222"
	scheduleID := "a5555555-5555-4000-8000-555555555555"

	encodedBranchID, _ := encoder.Encode(branchID)

	// Mock: branch ownership check
	mockBranchService.On("GetBranchByID", mock.Anything, branchID).Return(&domain.Branch{
		ID:               branchID,
		RepresentativeID: ownerID,
	}, nil)

	// Mock: tx + create schedule
	mockScheduleService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockScheduleService.On("CreateSchedule", mock.Anything, mockTx, branchID).Return(&domain.BranchSchedule{
		ID:        scheduleID,
		BranchID:  branchID,
		Active:    false,
		StartDate: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	// Setup router and execute (no request body needed)
	router := gin.New()
	router.POST("/branches/:id/schedules", func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{ID: ownerID})
	}, h.CreateBranchSchedule(scheduleInteractor))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/branches/"+encodedBranchID+"/schedules", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_H_CREATE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["branch_id"])
	assert.NotEmpty(t, data["_links"])

	mockScheduleService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
