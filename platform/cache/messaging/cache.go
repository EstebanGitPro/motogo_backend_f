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
	"MOD_U_ROLE_REQ_ERR_00011":   http.StatusUnauthorized, // 401 - Rol requerido

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
	// Branch Module (MOD_B_*) - HU59
	// ========================================
	"MOD_B_REG_EXI_00001":             http.StatusCreated,    // 201 - Sede registrada exitosamente
	"MOD_B_REG_ERR_00001":             http.StatusBadRequest, // 400 - Error al registrar (datos inválidos)
	"MOD_B_DUP_NAME_ERR_00001":        http.StatusConflict,   // 409 - Nombre duplicado en franquicia
	"MOD_B_INVALID_TYPE_ERR_00001":    http.StatusBadRequest, // 400 - Tipo de establecimiento inválido
	"MOD_B_NOT_FOUND_ERR_00001":       http.StatusNotFound,   // 404 - Sede no encontrada
	"MOD_B_UPD_EXI_00001":             http.StatusOK,         // 200 - Sede actualizada
	"MOD_B_DEL_EXI_00001":             http.StatusOK,         // 200 - Sede eliminada
	"MOD_B_BRAND_NOT_FOUND_ERR_00001": http.StatusBadRequest, // 400 - Marca no encontrada
	"MOD_B_DEL_ERR_00001":             http.StatusBadRequest, // 400 - Error al eliminar sede
	"MOD_B_HAS_ASSOC_ERR_00001":       http.StatusConflict,   // 409 - Sede con asociaciones

	// ========================================
	// Person Module - Delete (HU53)
	// ========================================
	"MOD_P_DEL_EXI_00001":          http.StatusOK,         // 200 - Persona eliminada exitosamente
	"MOD_P_DEL_ERR_00001":          http.StatusBadRequest, // 400 - Error al eliminar persona
	"MOD_P_HAS_BRANCHES_ERR_00001": http.StatusConflict,   // 409 - Persona tiene sedes activas

	// ========================================
	// Success Messages - Resource Creation (201 Created)
	// ========================================
	// User Module
	"MOD_U_REG_EXI_00001": http.StatusCreated, // 201 - Usuario registrado exitosamente
	// Person Module
	"MOD_P_REG_EXI_00001": http.StatusCreated, // 201 - Persona registrada exitosamente

	// ========================================
	// Franchise Module (MOD_F_*) - HU26-29
	// ========================================
	"MOD_F_REG_EXI_00001":              http.StatusCreated,    // 201 - Franquicia registrada
	"MOD_F_GET_EXI_00001":              http.StatusOK,         // 200 - Franquicia encontrada
	"MOD_F_LIST_EXI_00001":             http.StatusOK,         // 200 - Franquicias listadas
	"MOD_F_UPD_EXI_00001":              http.StatusOK,         // 200 - Franquicia actualizada
	"MOD_F_DEL_EXI_00001":              http.StatusOK,         // 200 - Franquicia eliminada
	"MOD_F_NOT_FOUND_ERR_00001":        http.StatusNotFound,   // 404 - Franquicia no encontrada
	"MOD_F_DUP_NAME_ERR_00001":         http.StatusConflict,   // 409 - Nombre duplicado
	"MOD_F_NO_BRANCHES_ERR_00001":      http.StatusBadRequest, // 400 - Debe asociar al menos una sede
	"MOD_F_BRANCH_NOT_OWNED_ERR_00001": http.StatusForbidden,  // 403 - Sede no pertenece al representante
	"MOD_F_HAS_BRANCHES_ERR_00001":     http.StatusConflict,   // 409 - Franquicia tiene sedes
	"MOD_F_BRANCH_ADD_EXI_00001":       http.StatusOK,         // 200 - Sede vinculada
	"MOD_F_BRANCH_REM_EXI_00001":       http.StatusOK,         // 200 - Sede desvinculada
	"MOD_F_MIN_BRANCHES_ERR_00001":     http.StatusBadRequest, // 400 - Mínimo una sede

	// ========================================
	// Service Catalog Module (MOD_S_*) - HU63, HU68, HU75
	// ========================================
	"MOD_S_TYPES_EXI_00001":         http.StatusOK,         // 200 - Tipos de servicio obtenidos (HU75)
	"MOD_S_LIST_EXI_00001":          http.StatusOK,         // 200 - Catálogo de servicios obtenido (HU63)
	"MOD_S_UPD_EXI_00001":           http.StatusOK,         // 200 - Servicio actualizado (HU68 - Admin)
	"MOD_S_ACTIVATED_EXI_00001":     http.StatusOK,         // 200 - Servicio activado (HU68 - Admin)
	"MOD_S_DEACTIVATED_EXI_00001":   http.StatusOK,         // 200 - Servicio desactivado (HU68 - Admin)
	"MOD_S_RES_ERR_00001":           http.StatusNotFound,   // 404 - Servicio no encontrado (HU68)
	"MOD_S_TYPE_ERR_00001":          http.StatusBadRequest, // 400 - Tipo de servicio inválido (HU68)
	"MOD_S_INVALID_TYPE_ERR_00001":  http.StatusBadRequest, // 400 - Tipo de servicio inválido
	"MOD_S_ASSOC_EXI_00001":         http.StatusOK,         // 200 - Servicios asociados a sede
	"MOD_S_DISSOC_EXI_00001":        http.StatusOK,         // 200 - Servicio desasociado de sede
	"MOD_S_NOT_FOUND_ERR_00001":     http.StatusNotFound,   // 404 - Servicio no encontrado
	"MOD_S_ALREADY_ASSOC_ERR_00001": http.StatusConflict,   // 409 - Servicio ya asociado

	// ========================================
	// Schedule Module (MOD_H_*) - HU30-35, HU10
	// ========================================
	// Success messages
	"MOD_H_CREATE_EXI_00001": http.StatusCreated, // 201 - Horario registrado
	"MOD_H_GET_EXI_00001":    http.StatusOK,      // 200 - Horario consultado
	"MOD_H_UPDATE_EXI_00001": http.StatusOK,      // 200 - Horario actualizado
	"MOD_H_DELETE_EXI_00001": http.StatusOK,      // 200 - Horario eliminado
	"MOD_H_ACTIV_EXI_00001":  http.StatusOK,      // 200 - Horario activado
	"MOD_H_DEACT_EXI_00001":  http.StatusOK,      // 200 - Horario desactivado
	"MOD_H_DAYS_EXI_00001":   http.StatusOK,      // 200 - Catálogo de días
	// Error messages
	"MOD_H_NOT_FOUND_ERR_00001":  http.StatusNotFound,   // 404 - Horario no encontrado
	"MOD_H_EXISTS_ERR_00001":     http.StatusConflict,   // 409 - Sede ya tiene horario
	"MOD_H_DAY_ERR_00001":        http.StatusBadRequest, // 400 - Día inválido
	"MOD_H_TIME_ERR_00001":       http.StatusBadRequest, // 400 - Formato hora inválido
	"MOD_H_TIME_ORDER_ERR_00001": http.StatusBadRequest, // 400 - Hora cierre antes apertura
	"MOD_H_INACTIVE_ERR_00001":   http.StatusBadRequest, // 400 - Horario desactivado

	// ========================================
	// Schedule Detail Module (MOD_HD_*) - HU6-9
	// ========================================
	// Success messages
	"MOD_HD_CREATE_EXI_00001": http.StatusCreated, // 201 - Detalle horario registrado (HU6)
	"MOD_HD_GET_EXI_00001":    http.StatusOK,      // 200 - Detalle horario consultado
	"MOD_HD_UPDATE_EXI_00001": http.StatusOK,      // 200 - Detalle horario actualizado (HU7)
	"MOD_HD_DELETE_EXI_00001": http.StatusOK,      // 200 - Detalle horario eliminado (HU8)
	"MOD_HD_LIST_EXI_00001":   http.StatusOK,      // 200 - Detalles horario listados (HU9)
	// Error messages
	"MOD_HD_NOT_FOUND_ERR_00001":     http.StatusNotFound,   // 404 - Detalle horario no encontrado
	"MOD_HD_CONFLICT_ERR_00001":      http.StatusConflict,   // 409 - Conflicto de horario
	"MOD_HD_TIME_ERR_00001":          http.StatusBadRequest, // 400 - Formato hora inválido
	"MOD_HD_DAY_ERR_00001":           http.StatusBadRequest, // 400 - Día de la semana inválido
	"MOD_HD_DAY_CLOSED_ERR_00001":    http.StatusConflict,   // 409 - Día ya cerrado (no duplicar)
	"MOD_HD_DAY_HAS_SLOTS_ERR_00001": http.StatusConflict,   // 409 - Día tiene franjas (no cerrar)

	// ========================================
	// Schedule Exception Module (MOD_EH_*) - HU20-25
	// ========================================
	// Success messages
	"MOD_EH_CREATE_EXI_00001":     http.StatusCreated, // 201 - Excepción creada
	"MOD_EH_GET_EXI_00001":        http.StatusOK,      // 200 - Excepción consultada
	"MOD_EH_LIST_EXI_00001":       http.StatusOK,      // 200 - Excepciones listadas
	"MOD_EH_UPDATE_EXI_00001":     http.StatusOK,      // 200 - Excepción actualizada
	"MOD_EH_DELETE_EXI_00001":     http.StatusOK,      // 200 - Excepción eliminada
	"MOD_EH_ACTIVATE_EXI_00001":   http.StatusOK,      // 200 - Excepción activada
	"MOD_EH_DEACTIVATE_EXI_00001": http.StatusOK,      // 200 - Excepción desactivada
	// Error messages
	"MOD_EH_NOT_FOUND_ERR_00001":     http.StatusNotFound,   // 404 - Excepción no encontrada
	"MOD_EH_DATE_CONFLICT_ERR_00001": http.StatusConflict,   // 409 - Fecha duplicada/solapada
	"MOD_EH_DATE_PAST_ERR_00001":     http.StatusBadRequest, // 400 - Fecha pasada
	"MOD_EH_TIME_ERR_00001":          http.StatusBadRequest, // 400 - Formato hora inválido
	"MOD_EH_REDUNDANT_ERR_00001":     http.StatusConflict,   // 409 - Excepción redundante

	// ========================================
	// Motorcycle Module (MOD_MOT_*) - HU43-47
	// ========================================
	// Success messages
	"MOD_MOT_CREATE_EXI_00001":   http.StatusCreated, // 201 - Motocicleta registrada (HU43)
	"MOD_MOT_GET_EXI_00001":      http.StatusOK,      // 200 - Motocicleta consultada (HU46)
	"MOD_MOT_UPDATE_EXI_00001":   http.StatusOK,      // 200 - Motocicleta actualizada (HU44)
	"MOD_MOT_DELETE_EXI_00001":   http.StatusOK,      // 200 - Motocicleta eliminada (HU45)
	"MOD_MOT_LIST_EXI_00001":     http.StatusOK,      // 200 - Motocicletas listadas (HU47)
	"MOD_MOT_REF_LIST_EXI_00001": http.StatusOK,      // 200 - Referencias listadas (HU50)
	// Error messages
	"MOD_MOT_NOT_FOUND_ERR_00001":     http.StatusNotFound,   // 404 - Motocicleta no encontrada
	"MOD_MOT_CREATE_ERR_00001":        http.StatusBadRequest, // 400 - Error al registrar
	"MOD_MOT_UPDATE_ERR_00001":        http.StatusBadRequest, // 400 - Error al actualizar
	"MOD_MOT_DELETE_ERR_00001":        http.StatusBadRequest, // 400 - Error al eliminar
	"MOD_MOT_DUP_PLATE_ERR_00001":     http.StatusConflict,   // 409 - Placa duplicada
	"MOD_MOT_REF_NOT_FOUND_ERR_00001": http.StatusNotFound,   // 404 - Referencia no encontrada
	"MOD_MOT_REF_REQ_ERR_00001":       http.StatusBadRequest, // 400 - Referencia requerida
	"MOD_MOT_HAS_ASSOC_ERR_00001":     http.StatusConflict,   // 409 - Tiene registros asociados
	"MOD_MOT_FORBIDDEN_ERR_00001":     http.StatusForbidden,  // 403 - Acceso denegado
	"MOD_MOT_LIST_ERR_00001":          http.StatusBadRequest, // 400 - Error al listar
	"MOD_MOT_PLATE_REQ_ERR_00001":     http.StatusBadRequest, // 400 - Placa requerida

	// ========================================
	// Geolocation Validation (MOD_V_GEO_*) - HU89
	// ========================================
	"MOD_V_GEO_LAT_REQ_ERR_00001": http.StatusBadRequest, // 400 - Latitud requerida
	"MOD_V_GEO_LAT_INV_ERR_00001": http.StatusBadRequest, // 400 - Latitud inválida (-90 a 90)
	"MOD_V_GEO_LNG_REQ_ERR_00001": http.StatusBadRequest, // 400 - Longitud requerida
	"MOD_V_GEO_LNG_INV_ERR_00001": http.StatusBadRequest, // 400 - Longitud inválida (-180 a 180)
	"MOD_V_GEO_RAD_INV_ERR_00001": http.StatusBadRequest, // 400 - Radio inválido (0-50km)

	// ========================================
	// Nearby Branches (MOD_B_NEARBY_*) - HU89
	// ========================================
	"MOD_B_NEARBY_EXI_00001": http.StatusOK, // 200 - Sedes cercanas encontradas

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
