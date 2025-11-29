package messaging

import (
	"context"
	"net/http"
	"sync"
	"time"

	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Re-export types for backward compatibility
type MessageType = cachetypes.MessageType
type CachedMessage = cachetypes.CachedMessage
type MessageResponse = cachetypes.MessageResponse
type MessageCacheRepository = cachetypes.MessageCacheRepository

const (
	TypeError   = cachetypes.TypeError
	TypeSuccess = cachetypes.TypeSuccess
	TypeWarning = cachetypes.TypeWarning
	TypeInfo    = cachetypes.TypeInfo
	TypeDebug   = cachetypes.TypeDebug
)

// MessageCache handles message caching in memory (future: Redis)
type MessageCache struct {
	repo            MessageCacheRepository
	log             logger.Logger
	messages        map[string]*CachedMessage
	mu              sync.RWMutex
	refreshInterval time.Duration
	stopRefresh     chan bool
}

// NewMessageCache creates a new message cache instance
func NewMessageCache(repo MessageCacheRepository, log logger.Logger, refreshInterval time.Duration) *MessageCache {
	return &MessageCache{
		repo:            repo,
		log:             log,
		messages:        make(map[string]*CachedMessage),
		refreshInterval: refreshInterval,
		stopRefresh:     make(chan bool),
	}
}

// LoadMessages loads all active messages from DB into cache
// Must be called at API startup
// Future: This will load from Redis or fallback to DB
func (c *MessageCache) LoadMessages(ctx context.Context) error {
	messages, err := c.repo.GetAllActiveForCache(ctx)
	if err != nil {
		c.log.Error(logger.LogMsgCacheLoadError, "error", err.Error())
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = make(map[string]*CachedMessage)
	for i := range messages {
		c.messages[messages[i].Code] = &messages[i]
	}

	c.log.Info(logger.LogMsgCacheLoaded, "count", len(c.messages))
	return nil
}

// ReloadMessages reloads messages from DB
// Future: Will invalidate Redis cache and reload
func (c *MessageCache) ReloadMessages(ctx context.Context) error {
	return c.LoadMessages(ctx)
}

// StartAutoRefresh starts a background goroutine to refresh cache periodically
// Should be called after LoadMessages in dependency initialization
func (c *MessageCache) StartAutoRefresh(ctx context.Context) {
	if c.refreshInterval <= 0 {
		c.log.Info(logger.LogMsgCacheRefreshDisabled)
		return
	}

	c.log.Info(logger.LogMsgCacheRefreshStart, "interval", c.refreshInterval.String())

	go func() {
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.log.Debug(logger.LogMsgCacheRefreshing)
				if err := c.ReloadMessages(ctx); err != nil {
					c.log.Error(logger.LogMsgCacheRefreshError, "error", err.Error())
				} else {
					c.log.Debug(logger.LogMsgCacheRefreshOK, "count", c.MessageCount())
				}
			case <-c.stopRefresh:
				c.log.Info(logger.LogMsgCacheRefreshStop)
				return
			}
		}
	}()
}

// StopAutoRefresh stops the auto-refresh goroutine
func (c *MessageCache) StopAutoRefresh() {
	if c.refreshInterval > 0 {
		close(c.stopRefresh)
	}
}

// GetMessage retrieves a message by its code from cache
// If not found in cache, falls back to DB and caches it
func (c *MessageCache) GetMessage(code string) *CachedMessage {
	// Try cache first (read lock)
	c.mu.RLock()
	msg, found := c.messages[code]
	c.mu.RUnlock()

	if found {
		return msg
	}

	// Not in cache, try DB
	c.log.Debug(logger.LogMsgNotInCache, "code", code)
	dbMsg, err := c.repo.GetByCodeForCache(context.Background(), code)
	if err != nil || dbMsg == nil {
		c.log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)
		return nil
	}

	// Cache it for future use (write lock)
	c.mu.Lock()
	c.messages[code] = dbMsg
	c.mu.Unlock()

	c.log.Debug(logger.LogMsgCachedFromDB, "code", code)
	return dbMsg
}

// GetMessageResponse retrieves formatted message response
func (c *MessageCache) GetMessageResponse(code string, params ...string) *MessageResponse {
	msg := c.GetMessage(code)
	if msg == nil {
		return nil
	}

	// Replace placeholders in content
	content := msg.Content
	for i, param := range params {
		placeholder := "${" + string(rune('0'+i)) + "}"
		content = replaceAll(content, placeholder, param)
	}

	return &MessageResponse{
		Code:    msg.Code,
		Type:    msg.Type,
		Title:   msg.Title,
		Content: content,
	}
}

// replaceAll is a simple helper for placeholder replacement
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

// messageCodeToHTTPStatus maps message codes to HTTP status codes
var messageCodeToHTTPStatus = map[string]int{
	"MOD_U_DUP_ERR_00001": http.StatusConflict,

	"MOD_U_EMAIL_NF_ERR_00005":  http.StatusNotFound,
	"MOD_P_NOT_FOUND_ERR_00001": http.StatusNotFound,
	"MOD_U_GET_ERR_00003":       http.StatusNotFound,
	"MOD_U_TOKEN_NF_ERR_00007":  http.StatusNotFound,

	"MOD_V_VAL_ERR_00001":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00002":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00006":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00008":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00009":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00010":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00011":  http.StatusBadRequest,
	"MOD_V_JSON_ERR_00012": http.StatusBadRequest,
	"MOD_V_ID_ERR_00013":   http.StatusBadRequest,

	"MOD_U_EMAIL_NV_ERR_00006": http.StatusForbidden,

	"MOD_U_TOKEN_EXP_ERR_00008":  http.StatusUnauthorized,
	"MOD_U_TOKEN_USED_ERR_00009": http.StatusUnauthorized,

	"GEN_AUTH_ERR_00002":      http.StatusUnauthorized,
	"GEN_FORBIDDEN_ERR_00003": http.StatusForbidden,
}

// GetHTTPStatus returns the HTTP status for a message code
func (c *MessageCache) GetHTTPStatus(code string) int {
	if status, ok := messageCodeToHTTPStatus[code]; ok {
		return status
	}

	msg := c.GetMessage(code)
	if msg == nil {
		return http.StatusInternalServerError
	}

	switch msg.Type {
	case TypeSuccess:
		return http.StatusOK
	case TypeError:
		return http.StatusInternalServerError
	case TypeWarning:
		return http.StatusOK
	case TypeInfo:
		return http.StatusOK
	case TypeDebug:
		return http.StatusOK
	default:
		return http.StatusOK
	}
}

// MessageCount returns the number of loaded messages in cache
func (c *MessageCache) MessageCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.messages)
}
