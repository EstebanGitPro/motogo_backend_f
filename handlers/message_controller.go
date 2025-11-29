package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateMessage handles POST requests to create a new system message
func (h handler) CreateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info(logger.LogMessageCreate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var messageRequest MessageRequest
		if err := c.ShouldBindJSON(&messageRequest); err != nil {
			h.Logger.Error(logger.LogMiddlewareJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		h.Logger.Info(logger.LogMessageCreateProcessing,
			"code", messageRequest.Code,
			"type", messageRequest.Type)

		message := messageRequest.ToDomain()
		message.SetID() // Generate UUID

		result, err := h.MessageInteractor.CreateMessage(c, message)
		if err != nil {
			h.Logger.Error(logger.LogMessageCreateError,
				"code", messageRequest.Code,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUID for public API
		encodedID, err := h.EncodeID(result.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.ID, err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildMessageCreatedLinks(baseURL, encodedID)

		// Set Location header
		SetLocationHeader(c, baseURL, "messages", encodedID)

		response := MessageCreatedResponse{
			Message: "Mensaje creado exitosamente",
			ID:      encodedID,
			Links:   links,
		}

		h.Logger.Success("Mensaje creado exitosamente",
			"id", result.ID,
			"encoded_id", encodedID,
			"code", result.Code,
			"status", http.StatusCreated,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusCreated, response)
	}
}

// UpdateMessage handles PUT requests to update an existing system message
func (h handler) UpdateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info(logger.LogMessageUpdate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.DecodeID(encodedID)
		if err != nil {
			h.Logger.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		var messageRequest MessageRequest
		if err := c.ShouldBindJSON(&messageRequest); err != nil {
			h.Logger.Error(logger.LogMiddlewareJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		h.Logger.Info(logger.LogMessageUpdateProcessing,
			"id", uuid,
			"code", messageRequest.Code)

		message := messageRequest.ToDomain()
		message.ID = uuid

		result, err := h.MessageInteractor.UpdateMessage(c, message)
		if err != nil {
			h.Logger.Error(logger.LogMessageUpdateError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildMessageUpdatedLinks(baseURL, encodedID)

		response := MessageUpdatedResponse{
			Message: "Mensaje actualizado exitosamente",
			Links:   links,
		}

		h.Logger.Success("Mensaje actualizado exitosamente",
			"id", result.ID,
			"code", result.Code,
			"status", http.StatusOK,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusOK, response)
	}
}

// DeleteMessage handles DELETE requests to delete a system message
func (h handler) DeleteMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info(logger.LogMessageDelete,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.DecodeID(encodedID)
		if err != nil {
			h.Logger.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		h.Logger.Info(logger.LogMessageDeleteProcessing, "id", uuid)

		err = h.MessageInteractor.DeleteMessage(c, uuid)
		if err != nil {
			h.Logger.Error(logger.LogMessageDeleteError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		response := MessageDeletedResponse{
			Message: "Mensaje eliminado exitosamente",
		}

		h.Logger.Success("Mensaje eliminado exitosamente",
			"id", uuid,
			"status", http.StatusOK,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusOK, response)
	}
}

// GetMessageByID handles GET requests to retrieve a message by ID
func (h handler) GetMessageByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Debug(logger.LogMessageGet,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.DecodeID(encodedID)
		if err != nil {
			h.Logger.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		message, err := h.MessageInteractor.GetMessageByID(c, uuid)
		if err != nil {
			h.Logger.Error(logger.LogMessageGetError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUID for response
		encodedIDForResponse, err := h.EncodeID(message.ID)
		if err != nil {
			h.HandleIDEncodingError(c, message.ID, err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToMessageResponse(message)
		response.ID = encodedIDForResponse // Use encoded ID in response
		response.Links = BuildMessageLinks(baseURL, encodedIDForResponse)

		h.Logger.Debug(logger.LogMessageGetOK,
			"id", uuid,
			"code", message.Code,
			"status", http.StatusOK,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusOK, response)
	}
}

// ListMessages handles GET requests to list messages with optional filters
func (h handler) ListMessages() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Debug(logger.LogMessageList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if module := c.Query("module"); module != "" {
			filters["module"] = module
		}
		if msgType := c.Query("type"); msgType != "" {
			filters["type"] = msgType
		}
		if category := c.Query("category"); category != "" {
			filters["category"] = category
		}
		if active := c.Query("active"); active != "" {
			if active == "true" || active == "1" {
				filters["active"] = true
			} else if active == "false" || active == "0" {
				filters["active"] = false
			}
		}

		messages, err := h.MessageInteractor.ListMessages(c, filters)
		if err != nil {
			h.Logger.Error(logger.LogMessageListError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUIDs for each message in response
		baseURL := GetBaseURL(c)
		response := ToMessageListResponse(messages)
		for i := range response.Messages {
			encodedID, err := h.EncodeID(messages[i].ID)
			if err != nil {
				h.HandleIDEncodingError(c, messages[i].ID, err)
				return
			}
			response.Messages[i].ID = encodedID
			response.Messages[i].Links = BuildMessageLinks(baseURL, encodedID)
		}
		response.Links = BuildMessageListLinks(baseURL)

		h.Logger.Debug(logger.LogMessageListOK,
			"count", len(messages),
			"status", http.StatusOK,
			"client_ip", c.ClientIP())

		c.JSON(http.StatusOK, response)
	}
}
