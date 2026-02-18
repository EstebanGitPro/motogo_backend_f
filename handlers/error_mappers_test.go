package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// errorMapperTestCase defines a single test case for error mapper functions.
type errorMapperTestCase struct {
	name         string
	inputErr     error
	expectedCode string
}

// TestMapRegisterBranchError covers all branches of mapRegisterBranchError.
func TestMapRegisterBranchError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"invalid branch type", domain.ErrInvalidBranchType, domain.MsgBranchInvalidType},
		{"duplicate branch name", domain.ErrDuplicateBranchName, domain.MsgBranchDuplicateName},
		{"brand not found", domain.ErrBrandNotFound, domain.MsgBrandNotFound},
		{"duplicate address", domain.ErrDuplicateAddress, domain.MsgDuplicateAddress},
		{"default error", assert.AnError, domain.MsgBranchCannotSave},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapRegisterBranchError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapUpdateBranchError covers all branches of mapUpdateBranchError.
func TestMapUpdateBranchError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"branch not found", domain.ErrBranchNotFound, domain.MsgBranchNotFound},
		{"forbidden", domain.ErrForbidden, domain.MsgForbidden},
		{"invalid branch type", domain.ErrInvalidBranchType, domain.MsgBranchInvalidType},
		{"brand not found", domain.ErrBrandNotFound, domain.MsgBrandNotFound},
		{"default error", assert.AnError, domain.MsgBranchCannotUpdate},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapUpdateBranchError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapMotorcycleRegError covers all branches of mapMotorcycleRegError.
func TestMapMotorcycleRegError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"reference not found", domain.ErrReferenceNotFound, domain.MsgMotorcycleReferenceNotFound},
		{"reference required", domain.ErrReferenceRequired, domain.MsgReferenceRequired},
		{"duplicate plate", domain.ErrDuplicateLicensePlate, domain.MsgDuplicateLicensePlate},
		{"cannot save", domain.ErrMotorcycleCannotSave, domain.MsgMotorcycleCannotSave},
		{"default error", assert.AnError, domain.MsgServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapMotorcycleRegError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapMotorcycleUpdateError covers all branches of mapMotorcycleUpdateError.
func TestMapMotorcycleUpdateError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"motorcycle not found", domain.ErrMotorcycleNotFound, domain.MsgMotorcycleNotFound},
		{"reference not found", domain.ErrReferenceNotFound, domain.MsgMotorcycleReferenceNotFound},
		{"cannot update", domain.ErrMotorcycleCannotUpdate, domain.MsgMotorcycleCannotUpdate},
		{"default error", assert.AnError, domain.MsgServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapMotorcycleUpdateError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapUpdateProfileError covers all branches of mapUpdateProfileError.
func TestMapUpdateProfileError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"duplicate user", domain.ErrDuplicateUser, domain.MsgUserDuplicate},
		{"default error", assert.AnError, domain.MsgPersonUpdated},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapUpdateProfileError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapScheduleUpdateError covers all branches of mapScheduleUpdateError.
func TestMapScheduleUpdateError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"schedule not found", domain.ErrScheduleNotFound, domain.MsgScheduleNotFound},
		{"forbidden", domain.ErrForbidden, domain.MsgForbidden},
		{"default error", assert.AnError, domain.MsgServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapScheduleUpdateError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapDetailCreationError covers all branches of mapDetailCreationError.
func TestMapDetailCreationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"schedule not found", domain.ErrScheduleNotFound, domain.MsgScheduleNotFound},
		{"invalid day", domain.ErrScheduleDetailInvalidDay, domain.MsgScheduleDetailInvalidDay},
		{"invalid time", domain.ErrScheduleDetailInvalidTime, domain.MsgScheduleDetailInvalidTime},
		{"time conflict", domain.ErrScheduleDetailTimeConflict, domain.MsgScheduleDetailTimeConflict},
		{"day already closed", domain.ErrScheduleDetailDayAlreadyClosed, domain.MsgScheduleDetailDayAlreadyClosed},
		{"day has slots", domain.ErrScheduleDetailDayHasSlots, domain.MsgScheduleDetailDayHasSlots},
		{"branch not found", domain.ErrBranchNotFound, domain.MsgBranchNotFound},
		{"forbidden", domain.ErrForbidden, domain.MsgForbidden},
		{"default error", assert.AnError, domain.MsgServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapDetailCreationError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}

// TestMapExceptionCreationError covers all branches of mapExceptionCreationError.
func TestMapExceptionCreationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)
	h := handlers.NewForTestWithConcrete(nil, nil, nil, nil, encoder, responseHandler)

	tests := []errorMapperTestCase{
		{"schedule not found", domain.ErrScheduleNotFound, domain.MsgScheduleNotFound},
		{"date past", domain.ErrScheduleExceptionDatePast, domain.MsgScheduleExceptionDatePast},
		{"date conflict", domain.ErrScheduleExceptionDateConflict, domain.MsgScheduleExceptionDateConflict},
		{"invalid time", domain.ErrScheduleExceptionInvalidTime, domain.MsgScheduleExceptionInvalidTime},
		{"redundant", domain.ErrScheduleExceptionRedundant, domain.MsgScheduleExceptionRedundant},
		{"forbidden", domain.ErrForbidden, domain.MsgForbidden},
		{"default error", assert.AnError, domain.MsgServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			h.MapExceptionCreationError(c, tc.inputErr)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if msg, ok := resp["message"].(map[string]interface{}); ok {
				assert.Equal(t, tc.expectedCode, msg["code"], "case: %s", tc.name)
			}
		})
	}
}
