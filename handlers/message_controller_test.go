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
	mockTx.On("Rollback").Return(nil).Maybe()

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

// TestUpdateMessage_Integration_Success validates the full HTTP handler pipeline
// for the success path of UpdateMessage (PUT /admin/messages/:id).
func TestUpdateMessage_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockMessageService := new(mocks.MockMessageService)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockMessageService)
	h := handlers.NewForTestWithMessage(messageInteractor, encoder, responseHandler)

	// First: encode a UUID to use as the URL param
	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	encodedID, err := encoder.Encode(testUUID)
	assert.NoError(t, err)

	// Mock: GetMessageByID (step 1: verify exists)
	existingMessage := &domain.Message{
		ID:   testUUID,
		Code: "TEST_MSG_001",
		Type: domain.TypeSuccess,
	}
	mockMessageService.On("GetMessageByID", mock.Anything, testUUID).Return(existingMessage, nil)

	// Mock: ValidateMessage
	mockMessageService.On("ValidateMessage", mock.Anything, mock.AnythingOfType("domain.Message")).Return(nil)

	// Mock: transaction lifecycle
	mockMessageService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockMessageService.On("UpdateMessageInDB", mock.Anything, mockTx, mock.AnythingOfType("domain.Message")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	// Request body
	reqBody := map[string]interface{}{
		"code":     "TEST_MSG_UPDATED",
		"type":     "INFO",
		"category": "test",
		"module":   "testing",
		"title":    "Updated Message",
		"content":  "Updated content",
		"active":   true,
	}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/admin/messages/:id", h.UpdateMessage())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/admin/messages/"+encodedID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	mockMessageService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestDeleteMessage_Integration_Success validates the full HTTP handler pipeline
// for the success path of DeleteMessage (DELETE /admin/messages/:id).
func TestDeleteMessage_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockMessageService := new(mocks.MockMessageService)
	mockTx := new(mocks.MockTx)

	messageInteractor := interactor.NewMessageInteractor(mockMessageService)
	h := handlers.NewForTestWithMessage(messageInteractor, encoder, responseHandler)

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	encodedID, err := encoder.Encode(testUUID)
	assert.NoError(t, err)

	// Mock: GetMessageByID (step 1: verify exists)
	existingMessage := &domain.Message{
		ID:   testUUID,
		Code: "TEST_MSG_001",
	}
	mockMessageService.On("GetMessageByID", mock.Anything, testUUID).Return(existingMessage, nil)

	// Mock: transaction lifecycle
	mockMessageService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockMessageService.On("DeleteMessageFromDB", mock.Anything, mockTx, testUUID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	router := gin.New()
	router.DELETE("/admin/messages/:id", h.DeleteMessage())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/admin/messages/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	mockMessageService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestGetMessageByID_Integration_Success validates the full HTTP handler pipeline
// for the success path of GetMessageByID (GET /admin/messages/:id).
func TestGetMessageByID_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockMessageService := new(mocks.MockMessageService)
	messageInteractor := interactor.NewMessageInteractor(mockMessageService)
	h := handlers.NewForTestWithMessage(messageInteractor, encoder, responseHandler)

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	encodedID, err := encoder.Encode(testUUID)
	assert.NoError(t, err)

	// Mock: GetMessageByID returns the message
	expectedMessage := &domain.Message{
		ID:       testUUID,
		Code:     "TEST_MSG_001",
		Type:     domain.TypeSuccess,
		Category: "test",
		Module:   "testing",
		Title:    "Test Message",
		Content:  "Test content",
		Active:   true,
	}
	mockMessageService.On("GetMessageByID", mock.Anything, testUUID).Return(expectedMessage, nil)

	router := gin.New()
	router.GET("/admin/messages/:id", h.GetMessageByID())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/messages/"+encodedID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "response should contain 'data' object")
	assert.Equal(t, "TEST_MSG_001", data["code"])
	assert.NotEmpty(t, data["_links"])

	mockMessageService.AssertExpectations(t)
}

// TestListMessages_Integration_Success validates the full HTTP handler pipeline
// for the success path of ListMessages (GET /admin/messages).
func TestListMessages_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockMessageService := new(mocks.MockMessageService)
	messageInteractor := interactor.NewMessageInteractor(mockMessageService)
	h := handlers.NewForTestWithMessage(messageInteractor, encoder, responseHandler)

	// Mock: ListMessages returns a list
	messageList := []domain.Message{
		{
			ID:       "550e8400-e29b-41d4-a716-446655440001",
			Code:     "MSG_001",
			Type:     domain.TypeSuccess,
			Category: "test",
			Module:   "testing",
			Title:    "Message 1",
			Content:  "Content 1",
			Active:   true,
		},
		{
			ID:       "550e8400-e29b-41d4-a716-446655440002",
			Code:     "MSG_002",
			Type:     domain.TypeError,
			Category: "test",
			Module:   "testing",
			Title:    "Message 2",
			Content:  "Content 2",
			Active:   true,
		},
	}
	mockMessageService.On("ListMessages", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(messageList, nil)

	router := gin.New()
	router.GET("/admin/messages", h.ListMessages())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/messages?module=testing", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "response should contain 'data' object")
	assert.NotEmpty(t, data["messages"])
	assert.Equal(t, float64(2), data["count"])

	mockMessageService.AssertExpectations(t)
}
