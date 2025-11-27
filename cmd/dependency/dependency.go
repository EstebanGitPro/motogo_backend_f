package dependency

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/identity_provider/keycloak"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/databases/mysql"

	messageRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/message"
	repo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/person"
)

type Dependencies struct {
	PersonService   input.Service
	PersonRepo      output.Repository
	KeycloakClient  output.AuthClient
	Interactor      *interactor.Interactor
	Config          *config.Config
	Logger          logger.Logger
	IDEncoder       *idencoder.HashidsEncoder
	MessagingCache  *messagingCache.CacheService
	ResponseHandler *middleware.ResponseHandler
}

func Init() (*Dependencies, error) {
	// Inicializar logger
	log := logger.NewSlogLogger()
	log.Info("Iniciando aplicación MotoGo Backend")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error cargando configuración", "error", err)
		return nil, err
	}
	log.Info("Configuración cargada exitosamente")

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error("Error conectando a base de datos", "error", err)
		return nil, err
	}
	log.Success("Conexión a base de datos establecida")

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
	if err != nil {
		log.Error("Error inicializando cliente Keycloak", "error", err)
		return nil, err
	}
	log.Success("Cliente Keycloak inicializado")

	personRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error("Error inicializando repositorio", "error", err)
		return nil, err
	}

	personService := services.NewService(personRepo, keycloakClient, log)

	interactorFacade := interactor.NewInteractor(personService, log)

	// Inicializar IDEncoder
	encoder, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    cfg.IDEncoder.Secret,
		MinLength: cfg.IDEncoder.MinLength,
	})
	if err != nil {
		log.Error("Error inicializando ID encoder", "error", err)
		return nil, err
	}
	log.Success("ID Encoder inicializado correctamente")

	// Inicializar repositorio y servicio de mensajería
	msgRepo, err := messageRepo.NewRepository(db)
	if err != nil {
		log.Error("Error inicializando repositorio de mensajes", "error", err)
		return nil, err
	}

	// Auto-refresh cada 5 minutos (ajustable según necesidad)
	// Para deshabilitar: usar 0
	refreshInterval := 5 * time.Minute
	messagingCache := messagingCache.NewCacheService(msgRepo, log, refreshInterval)

	if err := messagingCache.LoadMessages(context.Background()); err != nil {
		log.Warn("Error loading system messages into cache, using fallback", "error", err)
		// Don't return error, continue with fallback
	}
	log.Success("Message cache initialized", "messages_loaded", messagingCache.MessageCount())

	// Iniciar auto-refresh en background
	messagingCache.StartAutoRefresh(context.Background())

	responseHandler := middleware.NewResponseHandler(messagingCache)

	log.Success("Dependencias inicializadas correctamente")

	return &Dependencies{
		PersonService:   personService,
		PersonRepo:      personRepo,
		KeycloakClient:  keycloakClient,
		Interactor:      interactorFacade,
		Config:          cfg,
		Logger:          log,
		IDEncoder:       encoder,
		MessagingCache:  messagingCache,
		ResponseHandler: responseHandler,
	}, nil
}
