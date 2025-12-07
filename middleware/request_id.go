package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"

	// RequestIDKey is the key used to store request ID in gin context
	RequestIDKey = "request_id"

	// TraceIDKey is the key used for log correlation in Loki
	// Uses the same value as RequestIDKey but with a name that matches Loki's derivedField regex
	TraceIDKey = "traceID"
)

// RequestID is a middleware that generates a unique ID for each request
// and adds it to the context and response headers
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has an ID (from client or load balancer)
		requestID := c.GetHeader(RequestIDHeader)

		// If not, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for use in handlers and logs
		c.Set(RequestIDKey, requestID)

		// Add to response headers
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID extracts the request ID from gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
