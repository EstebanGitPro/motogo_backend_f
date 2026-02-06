package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/firebase"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	Interactor           *interactor.Interactor
	MessageInteractor    *interactor.MessageInteractor
	BranchInteractor     *interactor.BranchInteractor        // HU59
	BrandInteractor      input.BrandInteractorInterface      // Brands catalog (interface for testing)
	LocationInteractor   input.LocationInteractorInterface   // Geographic catalogs (interface for testing)
	ServiceInteractor    *interactor.ServiceInteractor       // Services catalog (HU63, HU75)
	FranchiseInteractor  *interactor.FranchiseInteractor     // Franchise CRUD (HU26-29)
	MotorcycleInteractor input.MotorcycleInteractorInterface // Motorcycle CRUD (interface for testing)
	EvidenceInteractor   input.EvidenceInteractorInterface   // Evidence CRUD (HU16-19)
	FirebaseClient       *firebase.Client                    // Firebase Auth
	MessagingCache       *messagingCache.MessageCache
	IDEncoder            *idencoder.HashidsEncoder
	Response             *middleware.ResponseHandler
}

func New(
	personInteractor *interactor.Interactor,
	messageInteractor *interactor.MessageInteractor,
	branchInteractor *interactor.BranchInteractor,
	brandInteractor *interactor.BrandInteractor,
	locationInteractor *interactor.LocationInteractor,
	serviceInteractor *interactor.ServiceInteractor,
	franchiseInteractor *interactor.FranchiseInteractor,
	motorcycleInteractor *interactor.MotorcycleInteractor,
	evidenceInteractor *interactor.EvidenceInteractor, // HU16-19
	firebaseClient *firebase.Client,
	messageCache *messagingCache.MessageCache,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		Interactor:           personInteractor,
		MessageInteractor:    messageInteractor,
		BranchInteractor:     branchInteractor,
		BrandInteractor:      brandInteractor,
		LocationInteractor:   locationInteractor,
		ServiceInteractor:    serviceInteractor,
		FranchiseInteractor:  franchiseInteractor,
		MotorcycleInteractor: motorcycleInteractor,
		EvidenceInteractor:   evidenceInteractor, // HU16-19
		FirebaseClient:       firebaseClient,
		MessagingCache:       messageCache,
		IDEncoder:            encoder,
		Response:             responseHandler,
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
