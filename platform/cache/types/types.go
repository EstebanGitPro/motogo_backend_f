package types

import "context"

// MessageType represents system message types
type MessageType string

const (
	TypeError   MessageType = "ERROR"
	TypeSuccess MessageType = "EXITO"
	TypeWarning MessageType = "WARNING"
	TypeInfo    MessageType = "INFO"
	TypeDebug   MessageType = "DEBUG"
)

// CachedMessage represents a message in the cache
// This is cache's own type, independent from domain
type CachedMessage struct {
	ID       string 
	Code     string
	Type     MessageType
	Category string
	Module   string
	Title    string
	Content  string
	Active   bool
}

// MessageResponse is the response structure for the frontend
type MessageResponse struct {
	Code    string      `json:"code"`
	Type    MessageType `json:"type"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
}

// MessageCacheRepository defines what the cache needs from a data source
type MessageCacheRepository interface {
	GetAllActiveForCache(ctx context.Context) ([]CachedMessage, error)
	GetByCodeForCache(ctx context.Context, code string) (*CachedMessage, error)
}
