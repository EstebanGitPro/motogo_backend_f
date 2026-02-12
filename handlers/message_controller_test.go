package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateMessage_Integration_Success validates the full HTTP handler pipeline
// for the success path of CreateMessage (POST /admin/messages).
//
// It exercises: JSON binding → sanitization → SetID → MessageInteractor.CreateMessage
// (validate → tx → save → commit) → ID encoding → HATEOAS links → 201 response.
func TestCreateMessage_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies for MessageInteractor
	mockMessageService := new(mocks.MockMessageService)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockMessageService)

	// Handler with MessageInteractor
	h := handlers.NewForTestWithMessage(messageInteractor, encoder, responseHandler)

	// Mock: validate message (passes)
	mockMessageService.On("ValidateMessage", mock.Anything, mock.AnythingOfType("domain.Message")).Return(nil)

	// Mock: transaction lifecycle
	mockMessageService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockMessageService.On("SaveMessageToDB", mock.Anything, mockTx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Request body
	reqBody := map[string]interface{}{
		"code":     "TEST_MSG_001",
		"type":     "EXITO",
		"category": "test",
		"module":   "testing",
		"title":    "Test Message",
		"content":  "This is a test message",
		"active":   true,
	}
	body, _ := json.Marshal(reqBody)

	// Setup router
	router := gin.New()
	router.POST("/admin/messages", h.CreateMessage())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_M_CREATE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"], "Message ID should be encoded")
	assert.NotEmpty(t, data["_links"])

	// Verify Location header was set
	assert.NotEmpty(t, w.Header().Get("Location"))

	mockMessageService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
