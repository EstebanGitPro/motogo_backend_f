package messaging

import (
	"context"
	"net/http"
	"sync"
	"time"

	messageRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/message"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// CacheService handles message caching in memory (future: Redis)
// This is infrastructure layer - will be replaced/extended with Redis
type CacheService struct {
	repo            messageRepo.Repository
	log             logger.Logger
	messages        map[string]*messageRepo.SystemMessage
	mu              sync.RWMutex
	refreshInterval time.Duration
	stopRefresh     chan bool
}

// NewCacheService creates a new message cache service instance
// refreshInterval: 0 = no auto-refresh, > 0 = refresh every N duration
func NewCacheService(repo messageRepo.Repository, log logger.Logger, refreshInterval time.Duration) *CacheService {
	return &CacheService{
		repo:            repo,
		log:             log,
		messages:        make(map[string]*messageRepo.SystemMessage),
		refreshInterval: refreshInterval,
		stopRefresh:     make(chan bool),
	}
}

// LoadMessages loads all active messages from DB into cache
// Must be called at API startup
// Future: This will load from Redis or fallback to DB
func (s *CacheService) LoadMessages(ctx context.Context) error {
	messages, err := s.repo.GetAllActive(ctx)
	if err != nil {
		s.log.Error("Error loading system messages from DB", "error", err.Error())
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = make(map[string]*messageRepo.SystemMessage)
	for i := range messages {
		s.messages[messages[i].Code] = &messages[i]
	}

	s.log.Info("System messages loaded into cache", "count", len(s.messages))
	return nil
}

// ReloadMessages reloads messages from DB
// Future: Will invalidate Redis cache and reload
func (s *CacheService) ReloadMessages(ctx context.Context) error {
	return s.LoadMessages(ctx)
}

// StartAutoRefresh starts a background goroutine to refresh cache periodically
// Should be called after LoadMessages in dependency initialization
func (s *CacheService) StartAutoRefresh(ctx context.Context) {
	if s.refreshInterval <= 0 {
		s.log.Info("Auto-refresh disabled (interval = 0)")
		return
	}

	s.log.Info("Starting message cache auto-refresh", "interval", s.refreshInterval.String())

	go func() {
		ticker := time.NewTicker(s.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.log.Debug("Auto-refreshing message cache from DB")
				if err := s.ReloadMessages(ctx); err != nil {
					s.log.Error("Error during auto-refresh", "error", err.Error())
				} else {
					s.log.Debug("Message cache auto-refreshed successfully", "count", s.MessageCount())
				}
			case <-s.stopRefresh:
				s.log.Info("Stopping message cache auto-refresh")
				return
			}
		}
	}()
}

// StopAutoRefresh stops the auto-refresh goroutine
func (s *CacheService) StopAutoRefresh() {
	if s.refreshInterval > 0 {
		close(s.stopRefresh)
	}
}

// GetMessage retrieves a message by its code from cache
// If not found in cache, falls back to DB and caches it
func (s *CacheService) GetMessage(code string) *messageRepo.SystemMessage {
	// Try cache first (read lock)
	s.mu.RLock()
	msg, found := s.messages[code]
	s.mu.RUnlock()

	if found {
		return msg
	}

	// Not in cache, try DB
	s.log.Debug("Message not in cache, loading from DB", "code", code)
	dbMsg, err := s.repo.GetByCode(context.Background(), code)
	if err != nil || dbMsg == nil {
		s.log.Warn("Message not found in DB", "code", code, "error", err)
		return nil
	}

	// Cache it for future use (write lock)
	s.mu.Lock()
	s.messages[code] = dbMsg
	s.mu.Unlock()

	s.log.Debug("Message loaded from DB and cached", "code", code)
	return dbMsg
}

// GetMessageResponse retrieves formatted message response
func (s *CacheService) GetMessageResponse(code string, params ...string) *messageRepo.MessageResponse {
	msg := s.GetMessage(code)
	if msg == nil {
		return nil
	}
	resp := msg.ToResponse(params...)
	return &resp
}

// HTTPStatusForCode maps message codes to HTTP status
var httpStatusMap = map[string]int{
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
func (s *CacheService) GetHTTPStatus(code string) int {
	if status, ok := httpStatusMap[code]; ok {
		return status
	}

	msg := s.GetMessage(code)
	if msg == nil {
		return http.StatusInternalServerError
	}

	switch msg.Type {
	case messageRepo.TypeSuccess:
		return http.StatusOK
	case messageRepo.TypeError:
		return http.StatusInternalServerError
	case messageRepo.TypeWarning:
		return http.StatusOK
	case messageRepo.TypeInfo:
		return http.StatusOK
	case messageRepo.TypeDebug:
		return http.StatusOK
	default:
		return http.StatusOK
	}
}

// MessageCount returns the number of loaded messages in cache
func (s *CacheService) MessageCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages)
}
