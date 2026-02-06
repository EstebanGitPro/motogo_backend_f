package middleware

import (
	"net/http"

	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/gin-gonic/gin"
)

// ResponseHandler handles unified system responses
type ResponseHandler struct {
	cache *messagingCache.MessageCache
}

// NewResponseHandler creates a new response handler
func NewResponseHandler(cache *messagingCache.MessageCache) *ResponseHandler {
	return &ResponseHandler{
		cache: cache,
	}
}

// APIResponse unified structure for all responses
type APIResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code,omitempty"`    // Business message code
	Message string      `json:"message,omitempty"` // Business message content
	Data    interface{} `json:"data,omitempty"`
}

// StandardResponse is an alias for APIResponse - used in Swagger annotations
// @Description Standard API response wrapper
type StandardResponse = APIResponse

// Error sends an error response
func (h *ResponseHandler) Error(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Code:    code,
			Message: "Unknown error",
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: false,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

// ErrorWithData sends an error response with additional data
func (h *ResponseHandler) ErrorWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Code:    code,
			Message: "Unknown error",
			Data:    data,
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: false,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

// Success sends a success response
func (h *ResponseHandler) Success(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "Operation successful",
		})
		return
	}

	// Determine success flag based on message type
	isSuccess := msg.Type != "ERROR"

	c.JSON(status, APIResponse{
		Success: isSuccess,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

// SuccessWithData sends a success response with data
func (h *ResponseHandler) SuccessWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "Operation successful",
			Data:    data,
		})
		return
	}

	// Determine success flag based on message type
	// If message type is ERROR, success should be false even if called from SuccessWithData
	// This handles the case when a requested message doesn't exist and fallback error message is used
	isSuccess := msg.Type != "ERROR"

	c.JSON(status, APIResponse{
		Success: isSuccess,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

// Warning sends a warning response
func (h *ResponseHandler) Warning(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System warning",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

// WarningWithData sends a warning response with data
func (h *ResponseHandler) WarningWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System warning",
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

// Info sends an informational response
func (h *ResponseHandler) Info(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System information",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

// InfoWithData sends an informational response with data
func (h *ResponseHandler) InfoWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System information",
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

// DataOnly sends only data without message (for listings, etc.)
func (h *ResponseHandler) DataOnly(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}
