package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	PersonService     input.Service
	Interactor        *interactor.Interactor
	MessageInteractor *interactor.MessageInteractor
	MessagingCache    *messagingCache.MessageCache
	IDEncoder         *idencoder.HashidsEncoder
	Response          *middleware.ResponseHandler
}

func New(service input.Service, personInteractor *interactor.Interactor, messageInteractor *interactor.MessageInteractor, messageCache *messagingCache.MessageCache, encoder *idencoder.HashidsEncoder, responseHandler *middleware.ResponseHandler) *handler {
	return &handler{
		PersonService:     service,
		Interactor:        personInteractor,
		MessageInteractor: messageInteractor,
		MessagingCache:    messageCache,
		IDEncoder:         encoder,
		Response:          responseHandler,
	}
}

var Logger logger.Logger = logger.NewSlogLogger()

// EncodeID ofusca un UUID usando el encoder del handler
// Retorna el ID ofuscado o un error si falla
func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		Logger.Error(logger.LogMessageIDEncodeError,
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

// DecodeID desofusca un ID ofuscado a UUID usando el encoder del handler
// Retorna el UUID o un error si falla
func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		Logger.Error(logger.LogMessageIDDecodeError,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

// HandleIDEncodingError maneja errores de ofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	Logger.Error(logger.LogMessageIDEncodeError,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	Logger.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}
