package handlers

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateMessage godoc
// @Summary      Crear nuevo mensaje del sistema
// @Description  Crea un nuevo mensaje de sistema para manejo de notificaciones, errores y mensajes informativos
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        message  body      MessageRequest  true  "Datos del mensaje"
// @Success      201  {object}  middleware.APIResponse{data=MessageCreatedResponse}  "Mensaje creado exitosamente"
// @Failure      400  {object}  middleware.APIResponse  "Error de validación"
// @Failure      409  {object}  middleware.APIResponse  "Mensaje con código duplicado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /messages [post]
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
			ID:    encodedID,
			Links: links,
		}

		h.Logger.Success("Mensaje creado exitosamente",
			"id", result.ID,
			"encoded_id", encodedID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for message creation
		middleware.RecordMessageCreated(result.Module, string(result.Type))

		h.Response.SuccessWithData(c, domain.MsgMessageCreated, response)
	}
}

// UpdateMessage godoc
// @Summary      Actualizar mensaje existente
// @Description  Actualiza un mensaje del sistema por su ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "ID del mensaje (hashid)"
// @Param        message  body      MessageRequest  true  "Datos actualizados del mensaje"
// @Success      200  {object}  middleware.APIResponse{data=MessageUpdatedResponse}  "Mensaje actualizado exitosamente"
// @Failure      400  {object}  middleware.APIResponse  "Error de validación"
// @Failure      404  {object}  middleware.APIResponse  "Mensaje no encontrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /messages/{id} [put]
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
			Links: links,
		}

		h.Logger.Success("Mensaje actualizado exitosamente",
			"id", result.ID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageUpdated, response)
	}
}

// DeleteMessage godoc
// @Summary      Eliminar mensaje
// @Description  Elimina un mensaje del sistema por su ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "ID del mensaje (hashid)"
// @Success      200  {object}  middleware.APIResponse  "Mensaje eliminado exitosamente"
// @Failure      400  {object}  middleware.APIResponse  "ID inválido"
// @Failure      404  {object}  middleware.APIResponse  "Mensaje no encontrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /messages/{id} [delete]
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

		h.Logger.Success("Mensaje eliminado exitosamente",
			"id", uuid,
			"client_ip", c.ClientIP())

		h.Response.Success(c, domain.MsgMessageDeleted)
	}
}

// GetMessageByID godoc
// @Summary      Obtener mensaje por ID
// @Description  Obtiene un mensaje del sistema específico por su ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "ID del mensaje (hashid)"
// @Success      200  {object}  MessageResponse  "Mensaje encontrado"
// @Failure      400  {object}  middleware.APIResponse  "ID inválido"
// @Failure      404  {object}  middleware.APIResponse  "Mensaje no encontrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /messages/{id} [get]
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
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ListMessages godoc
// @Summary      Listar mensajes
// @Description  Obtiene una lista de mensajes del sistema con filtros opcionales
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        module    query     string  false  "Filtrar por módulo (ej: users, orders)"
// @Param        type      query     string  false  "Filtrar por tipo (ERROR, SUCCESS, INFO, WARNING)"
// @Param        category  query     string  false  "Filtrar por categoría (ej: usuario_final, admin)"
// @Param        active    query     boolean  false  "Filtrar por estado activo (true/false)"
// @Success      200  {object}  MessageListResponse  "Lista de mensajes"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /messages [get]
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
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ReloadMessageCache godoc
// @Summary      Recargar caché de mensajes
// @Description  Fuerza la recarga del caché de mensajes desde la base de datos. Útil después de actualizaciones o eliminaciones manuales.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Success      200  {object}  middleware.APIResponse{data=CacheReloadResponse}  "Caché recargado exitosamente"
// @Failure      500  {object}  middleware.APIResponse  "Error al recargar el caché"
// @Router       /messages/cache/reload [post]
func (h handler) ReloadMessageCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.Logger.Info("Recargando caché de mensajes",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Obtener el conteo antes del reload
		beforeCount := h.MessagingCache.MessageCount()

		// Recargar el caché desde BD
		err := h.MessagingCache.ReloadMessages(c.Request.Context())
		if err != nil {
			h.Logger.Error("Error al recargar caché de mensajes",
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInternalServer)
			return
		}

		// Obtener el conteo después del reload
		afterCount := h.MessagingCache.MessageCount()

		response := CacheReloadResponse{
			Success:     true,
			BeforeCount: beforeCount,
			AfterCount:  afterCount,
			Message:     "Caché de mensajes recargado exitosamente desde la base de datos",
		}

		h.Logger.Success("Caché de mensajes recargado exitosamente",
			"before_count", beforeCount,
			"after_count", afterCount,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
