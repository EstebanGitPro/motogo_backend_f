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
	messages        map[string]*CachedMessage
	mu              sync.RWMutex
	refreshInterval time.Duration
	stopRefresh     chan bool
}

var log logger.Logger = logger.NewSlogLogger()

// NewMessageCache creates a new message cache instance
func NewMessageCache(repo MessageCacheRepository, refreshInterval time.Duration) *MessageCache {
	return &MessageCache{
		repo:            repo,
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
		log.Error(logger.LogMsgCacheLoadError, "error", err.Error())
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = make(map[string]*CachedMessage)
	for i := range messages {
		c.messages[messages[i].Code] = &messages[i]
	}

	log.Info(logger.LogMsgCacheLoaded, "count", len(c.messages))
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
		log.Info(logger.LogMsgCacheRefreshDisabled)
		return
	}

	log.Info(logger.LogMsgCacheRefreshStart, "interval", c.refreshInterval.String())

	go func() {
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Debug(logger.LogMsgCacheRefreshing)
				if err := c.ReloadMessages(ctx); err != nil {
					log.Error(logger.LogMsgCacheRefreshError, "error", err.Error())
				} else {
					log.Debug(logger.LogMsgCacheRefreshOK, "count", c.MessageCount())
				}
			case <-c.stopRefresh:
				log.Info(logger.LogMsgCacheRefreshStop)
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
// If message exists but is inactive or doesn't exist, returns a specific error message
func (c *MessageCache) GetMessage(code string) *CachedMessage {
	// Try cache first (read lock)
	c.mu.RLock()
	msg, found := c.messages[code]
	c.mu.RUnlock()

	if found {
		return msg
	}

	// Not in cache, try DB (only active messages)
	log.Debug(logger.LogMsgNotInCache, "code", code)
	dbMsg, err := c.repo.GetByCodeForCache(context.Background(), code)
	if err != nil {
		log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)
		// Avoid infinite recursion if the fallback message itself doesn't exist
		if code == "GEN_MSG_INACTIVE_ERR_00002" {
			return nil
		}
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	if dbMsg != nil {
		// Cache it for future use (write lock)
		c.mu.Lock()
		c.messages[code] = dbMsg
		c.mu.Unlock()

		log.Debug(logger.LogMsgCachedFromDB, "code", code)
		return dbMsg
	}

	// Not found in active messages, check if it exists but is inactive
	inactiveMsg, err := c.repo.GetByCodeIncludingInactive(context.Background(), code)
	if err != nil {
		log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)
		// Avoid infinite recursion
		if code == "GEN_MSG_INACTIVE_ERR_00002" {
			return nil
		}
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	if inactiveMsg != nil && !inactiveMsg.Active {
		// Message exists but is inactive - return specific error message
		log.Warn(logger.LogMsgInactive, "code", code)
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	// Message truly doesn't exist (not even in DB)
	log.Warn(logger.LogMsgNotInDB, "code", code)
	// Avoid infinite recursion
	if code == "GEN_MSG_INACTIVE_ERR_00002" {
		return nil
	}
	return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
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
// Organized by modules: Users (MOD_U_*), Validation (MOD_V_*), Keycloak (MOD_KC_*),
// Infrastructure (MOD_INFRA_*), General (GEN_*), Messages (MOD_M_*)
var messageCodeToHTTPStatus = map[string]int{
	// ========================================
	// Users Module (MOD_U_*)
	// ========================================
	"MOD_U_DUP_ERR_00001":        http.StatusConflict,     // 409 - Usuario duplicado
	"MOD_U_EMAIL_NF_ERR_00005":   http.StatusNotFound,     // 404 - Email no encontrado
	"MOD_U_GET_ERR_00003":        http.StatusNotFound,     // 404 - Usuario no encontrado
	"MOD_U_TOKEN_NF_ERR_00007":   http.StatusNotFound,     // 404 - Token no encontrado
	"MOD_U_EMAIL_NV_ERR_00006":   http.StatusForbidden,    // 403 - Email no verificado
	"MOD_U_TOKEN_EXP_ERR_00008":  http.StatusUnauthorized, // 401 - Token expirado
	"MOD_U_TOKEN_USED_ERR_00009": http.StatusUnauthorized, // 401 - Token ya usado

	// ========================================
	// Person Module (MOD_P_*)
	// ========================================
	"MOD_P_NOT_FOUND_ERR_00001": http.StatusNotFound, // 404 - Persona no encontrada
	// Password Reset (HU56)
	"MOD_P_RESET_EXI_00001": http.StatusOK,                  // 200 - Password reset exitoso
	"MOD_P_RESET_ERR_00001": http.StatusBadRequest,          // 400 - Token inválido
	"MOD_P_RESET_ERR_00002": http.StatusNotFound,            // 404 - Usuario no encontrado
	"MOD_P_RESET_ERR_00003": http.StatusInternalServerError, // 500 - Error actualizando password
	// Change Password (HU57)
	"MOD_P_CHANGE_EXI_00001": http.StatusOK,                  // 200 - Password cambiado exitosamente
	"MOD_P_CHANGE_ERR_00001": http.StatusUnauthorized,        // 401 - Contraseña actual incorrecta
	"MOD_P_CHANGE_ERR_00002": http.StatusUnprocessableEntity, // 422 - Error procesando cambio de contraseña
	"MOD_P_CHANGE_ERR_00003": http.StatusBadRequest,          // 400 - Contraseña no cumple requisitos

	// ========================================
	// Validation Module (MOD_V_*)
	// ========================================
	"MOD_V_VAL_ERR_00001":  http.StatusBadRequest, // 400 - Formato inválido
	"MOD_V_VAL_ERR_00002":  http.StatusBadRequest, // 400 - Request inválido
	"MOD_V_VAL_ERR_00006":  http.StatusBadRequest, // 400 - Validación fallida
	"MOD_V_VAL_ERR_00008":  http.StatusBadRequest, // 400 - Formato de campo
	"MOD_V_VAL_ERR_00009":  http.StatusBadRequest, // 400 - Campo requerido
	"MOD_V_VAL_ERR_00010":  http.StatusBadRequest, // 400 - Tipo de campo
	"MOD_V_VAL_ERR_00011":  http.StatusBadRequest, // 400 - Múltiples errores
	"MOD_V_JSON_ERR_00012": http.StatusBadRequest, // 400 - JSON inválido
	"MOD_V_ID_ERR_00013":   http.StatusBadRequest, // 400 - ID inválido

	// ========================================
	// Keycloak Module (MOD_KC_*) - Email Verification & Auth
	// ========================================
	"MOD_KC_EMAIL_VERIFIED_EXI_00001":          http.StatusOK,                  // 200 - Email verificado exitosamente
	"MOD_KC_INVALID_TOKEN_ERR_00001":           http.StatusBadRequest,          // 400 - Token inválido/malformado
	"MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001":      http.StatusInternalServerError, // 500 - Error de verificación (falla en Keycloak)
	"MOD_KC_USER_NOT_FOUND_ERR_00001":          http.StatusNotFound,            // 404 - Usuario no encontrado
	"MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001": http.StatusOK,                  // 200 - Email ya verificado (warning)
	"MOD_KC_VERIF_EMAIL_SENT_EXI_00001":        http.StatusOK,                  // 200 - Email de verificación enviado
	"MOD_KC_VERIF_EMAIL_ERROR_ERR_00001":       http.StatusServiceUnavailable,  // 503 - Error enviando email
	"MOD_KC_PWD_RESET_SENT_EXI_00001":          http.StatusOK,                  // 200 - Email de reset enviado
	"MOD_KC_PWD_RESET_ERROR_ERR_00001":         http.StatusServiceUnavailable,  // 503 - Error enviando reset
	// Authentication Profile
	"MOD_AUTH_PROFILE_EXI_00001": http.StatusOK, // 200 - Perfil obtenido exitosamente

	// ========================================
	// Infrastructure Module (MOD_INFRA_*)
	// ========================================
	"MOD_INFRA_KC_UNAVAIL_ERR_00004":      http.StatusLocked,              // 423 - Keycloak no disponible
	"MOD_INFRA_DB_UNAVAIL_ERR_00005":      http.StatusLocked,              // 423 - Base de datos no disponible
	"MOD_INFRA_DEP_FAIL_ERR_00006":        http.StatusLocked,              // 423 - Falla de dependencia
	"MOD_INFRA_KC_CLEANUP_ERR_00003":      http.StatusLocked,              // 423 - Error limpieza Keycloak
	"MOD_INFRA_KC_CREATE_ERR_00002":       http.StatusLocked,              // 423 - Error creación en Keycloak
	"MOD_INFRA_INCOMPLETE_REG_ERR_00009":  http.StatusConflict,            // 409 - Registro incompleto
	"MOD_INFRA_KC_INCONSISTENT_ERR_00001": http.StatusInternalServerError, // 500 - Estado inconsistente

	// ========================================
	// General Module (GEN_*)
	// ========================================
	"GEN_AUTH_ERR_00002":         http.StatusUnauthorized,        // 401 - No autorizado
	"GEN_FORBIDDEN_ERR_00003":    http.StatusForbidden,           // 403 - Acceso denegado
	"GEN_MSG_INACTIVE_ERR_00002": http.StatusNotFound,            // 404 - Mensaje no disponible
	"GEN_SRV_ERR_00001":          http.StatusInternalServerError, // 500 - Error del servidor
	"GEN_OPE_EXI_00001":          http.StatusOK,                  // 200 - Operación exitosa
	"GEN_INFO_00001":             http.StatusOK,                  // 200 - Información
	"GEN_WARN_00001":             http.StatusOK,                  // 200 - Advertencia

	// ========================================
	// Messages Module (MOD_M_*)
	// ========================================
	"MOD_M_UPDATE_ERR_00010":    http.StatusBadRequest, // 400 - Error actualizando mensaje
	"MOD_M_NOT_FOUND_ERR_00001": http.StatusNotFound,   // 404 - Mensaje no encontrado
	"MOD_M_CREATE_EXI_00001":    http.StatusCreated,    // 201 - Mensaje creado exitosamente

	// ========================================
	// Success Messages - Resource Creation (201 Created)
	// ========================================
	// User Module
	"MOD_U_REG_EXI_00001": http.StatusCreated, // 201 - Usuario registrado exitosamente
	// Person Module
	"MOD_P_REG_EXI_00001": http.StatusCreated, // 201 - Persona registrada exitosamente

}

// GetHTTPStatus returns the HTTP status for a message code
func (c *MessageCache) GetHTTPStatus(code string) int {
	// First try direct lookup in the map
	if status, ok := messageCodeToHTTPStatus[code]; ok {
		return status
	}

	// If not in map, get the message (which may return a fallback)
	msg := c.GetMessage(code)
	if msg == nil {
		return http.StatusInternalServerError
	}

	// IMPORTANT: If we got a fallback message, use ITS code to lookup the status
	// This ensures that GEN_MSG_INACTIVE_ERR_00002 returns 404, not 500
	if msg.Code != code {
		// We got a fallback message, try to get its status from the map
		if status, ok := messageCodeToHTTPStatus[msg.Code]; ok {
			return status
		}
	}

	// Fallback to determining status by message type
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
