package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

func TestMessageRequest_ToDomain(t *testing.T) {
	// Arrange
	req := handlers.MessageRequest{
		Code:     "ERR_001",
		Type:     domain.TypeError,
		Category: "USER",
		Module:   "MOD_U",
		Title:    "Error Title",
		Content:  "Error content description",
		Active:   true,
	}

	// Act
	domainMsg := req.ToDomain()

	// Assert
	assert.Equal(t, req.Code, domainMsg.Code)
	assert.Equal(t, req.Type, domainMsg.Type)
	assert.Equal(t, req.Category, domainMsg.Category)
	assert.Equal(t, req.Module, domainMsg.Module)
	assert.Equal(t, req.Title, domainMsg.Title)
	assert.Equal(t, req.Content, domainMsg.Content)
	assert.Equal(t, req.Active, domainMsg.Active)
}

func TestToMessageResponse(t *testing.T) {
	// Arrange
	msg := &domain.Message{
		ID:       "msg-123",
		Code:     "SUC_001",
		Type:     domain.TypeSuccess,
		Category: "SYSTEM",
		Module:   "MOD_S",
		Title:    "Success Title",
		Content:  "Success content",
		Active:   true,
	}

	// Act
	response := handlers.ToMessageResponse(msg)

	// Assert
	assert.Equal(t, msg.ID, response.ID)
	assert.Equal(t, msg.Code, response.Code)
	assert.Equal(t, msg.Type, response.Type)
	assert.Equal(t, msg.Category, response.Category)
	assert.Equal(t, msg.Module, response.Module)
	assert.Equal(t, msg.Title, response.Title)
	assert.Equal(t, msg.Content, response.Content)
	assert.Equal(t, msg.Active, response.Active)
}

func TestToMessageListResponse(t *testing.T) {
	// Arrange
	messages := []domain.Message{
		{ID: "msg-1", Code: "CODE_001", Title: "Message 1"},
		{ID: "msg-2", Code: "CODE_002", Title: "Message 2"},
		{ID: "msg-3", Code: "CODE_003", Title: "Message 3"},
	}

	// Act
	response := handlers.ToMessageListResponse(messages)

	// Assert
	assert.Equal(t, 3, response.Count)
	assert.Len(t, response.Messages, 3)
	assert.Equal(t, "msg-1", response.Messages[0].ID)
	assert.Equal(t, "msg-2", response.Messages[1].ID)
	assert.Equal(t, "msg-3", response.Messages[2].ID)
}

func TestToMessageListResponse_Empty(t *testing.T) {
	// Arrange
	messages := []domain.Message{}

	// Act
	response := handlers.ToMessageListResponse(messages)

	// Assert
	assert.Equal(t, 0, response.Count)
	assert.Empty(t, response.Messages)
}

func TestMessageResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.MessageResponse{
		ID:       "msg-123",
		Code:     "ERR_001",
		Type:     domain.TypeError,
		Category: "USER",
		Module:   "MOD_U",
		Title:    "Test Title",
		Content:  "Test Content",
		Active:   true,
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "msg-123", result["id"])
	assert.Equal(t, "ERR_001", result["code"])
	assert.Equal(t, "ERROR", result["type"])
	assert.Equal(t, "USER", result["category"])
	assert.Equal(t, true, result["active"])
}

func TestMessageListResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.MessageListResponse{
		Messages: []handlers.MessageResponse{
			{ID: "msg-1", Code: "CODE_001"},
			{ID: "msg-2", Code: "CODE_002"},
		},
		Count: 2,
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, float64(2), result["count"])
	messages := result["messages"].([]interface{})
	assert.Len(t, messages, 2)
}

func TestMessageCreatedResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.MessageCreatedResponse{
		ID: "encoded-id-xyz",
		Links: []handlers.Link{
			{Href: "/messages/xyz", Rel: "self", Method: "GET"},
		},
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "encoded-id-xyz", result["id"])
	links := result["_links"].([]interface{})
	assert.Len(t, links, 1)
}

func TestCacheReloadResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.CacheReloadResponse{
		Success:     true,
		BeforeCount: 10,
		AfterCount:  15,
		Message:     "Cache reloaded successfully",
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, true, result["success"])
	assert.Equal(t, float64(10), result["before_count"])
	assert.Equal(t, float64(15), result["after_count"])
	assert.Equal(t, "Cache reloaded successfully", result["message"])
}
