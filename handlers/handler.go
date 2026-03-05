package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	Interactor                 *interactor.Interactor
	MessageInteractor          *interactor.MessageInteractor
	BranchInteractor           *interactor.BranchInteractor           // HU59
	BrandInteractor            input.BrandLister                      // Brands catalog (interface for testing)
	LocationInteractor         input.LocationInteractorInterface      // Geographic catalogs (interface for testing)
	ServiceInteractor          *interactor.ServiceInteractor          // Services catalog (HU63, HU75)
	FranchiseInteractor        *interactor.FranchiseInteractor        // Franchise CRUD (HU26-29)
	MotorcycleInteractor       input.MotorcycleInteractorInterface    // Motorcycle CRUD (interface for testing)
	EvidenceInteractor         input.EvidenceInteractorInterface      // Evidence CRUD (HU16-19)
	DiagnosticInteractor       *interactor.DiagnosticInteractor       // Diagnostic CRUD (HU11-14)
	CompletedServiceInteractor *interactor.CompletedServiceInteractor // Completed Service (HU64)
	RatingInteractor           *interactor.RatingInteractor           // Rating (HU48)
	FirebaseClient             output.CustomTokenProvider             // Firebase Auth
	MessagingCache             *messagingCache.MessageCache
	IDEncoder                  *idencoder.HashidsEncoder
	Response                   *middleware.ResponseHandler
	CookieConfig               config.CookieConfig
}

// HandlerConfig groups all dependencies required by the handler (SonarCloud S107: max 7 params).
type HandlerConfig struct {
	PersonInteractor           *interactor.Interactor
	MessageInteractor          *interactor.MessageInteractor
	BranchInteractor           *interactor.BranchInteractor
	BrandInteractor            input.BrandLister
	LocationInteractor         input.LocationInteractorInterface
	ServiceInteractor          *interactor.ServiceInteractor
	FranchiseInteractor        *interactor.FranchiseInteractor
	MotorcycleInteractor       input.MotorcycleInteractorInterface
	EvidenceInteractor         input.EvidenceInteractorInterface
	DiagnosticInteractor       *interactor.DiagnosticInteractor
	CompletedServiceInteractor *interactor.CompletedServiceInteractor
	RatingInteractor           *interactor.RatingInteractor
	FirebaseClient             output.CustomTokenProvider
	MessagingCache             *messagingCache.MessageCache
	IDEncoder                  *idencoder.HashidsEncoder
	ResponseHandler            *middleware.ResponseHandler
	CookieConfig               config.CookieConfig
}

func New(cfg HandlerConfig) *handler {
	return &handler{
		Interactor:                 cfg.PersonInteractor,
		MessageInteractor:          cfg.MessageInteractor,
		BranchInteractor:           cfg.BranchInteractor,
		BrandInteractor:            cfg.BrandInteractor,
		LocationInteractor:         cfg.LocationInteractor,
		ServiceInteractor:          cfg.ServiceInteractor,
		FranchiseInteractor:        cfg.FranchiseInteractor,
		MotorcycleInteractor:       cfg.MotorcycleInteractor,
		EvidenceInteractor:         cfg.EvidenceInteractor,
		DiagnosticInteractor:       cfg.DiagnosticInteractor,
		CompletedServiceInteractor: cfg.CompletedServiceInteractor,
		RatingInteractor:           cfg.RatingInteractor,
		FirebaseClient:             cfg.FirebaseClient,
		MessagingCache:             cfg.MessagingCache,
		IDEncoder:                  cfg.IDEncoder,
		Response:                   cfg.ResponseHandler,
		CookieConfig:               cfg.CookieConfig,
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
	_ = c.Error(domain.ErrInternalServer)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	Logger.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	_ = c.Error(domain.ErrInvalidID)
}
