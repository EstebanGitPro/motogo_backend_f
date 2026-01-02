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
	PersonService     input.Service
	PersonRepo        output.Repository
	KeycloakClient    output.AuthClient
	Interactor        *interactor.Interactor
	MessageInteractor *interactor.MessageInteractor
	Config            *config.Config
	Logger            logger.Logger
	IDEncoder         *idencoder.HashidsEncoder
	MessagingCache    *messagingCache.MessageCache
	ResponseHandler   *middleware.ResponseHandler
}

func Init() (*Dependencies, error) {
	log := logger.NewSlogLogger()
	log.Info(logger.LogAppStarting)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(logger.LogAppConfigError, "error", err)
		return nil, err
	}
	log.Info(logger.LogAppConfigLoaded)

	middleware.PrometheusInit()
	log.Success(logger.LogPrometheusInitOK)

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
	if err != nil {
		log.Error(logger.LogKeycloakClientError, "error", err)
		return nil, err
	}
	log.Success(logger.LogKeycloakClientOK)

	personRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error(logger.LogPersonRepoInitError, "error", err)
		return nil, err
	}

	personService := services.NewService(personRepo, keycloakClient, log)

	interactorFacade := interactor.NewInteractor(personService, log)

	encoder, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    cfg.IDEncoder.Secret,
		MinLength: cfg.IDEncoder.MinLength,
	}, log)
	if err != nil {
		log.Error(logger.LogIDEncoderInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogIDEncodeOK)

	msgRepo, err := messageRepo.NewMessageRepository(db)
	if err != nil {
		log.Error(logger.LogRepoMsgInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDependencyMessageRepoInit)

	refreshInterval := 5 * time.Minute
	messagingCache := messagingCache.NewMessageCache(msgRepo, refreshInterval)

	if err := messagingCache.LoadMessages(context.Background()); err != nil {
		log.Warn(logger.LogMsgCacheLoadError, "error", err)
	}
	log.Success(logger.LogMsgCacheInit, "messages_loaded", messagingCache.MessageCount())

	messagingCache.StartAutoRefresh(context.Background())

	responseHandler := middleware.NewResponseHandler(messagingCache)

	messageService := services.NewMessageService(msgRepo, log)
	messageInteractor := interactor.NewMessageInteractor(messageService, log)
	log.Success(logger.LogDependencyMessageIntInit)

	return &Dependencies{
		PersonService:     personService,
		PersonRepo:        personRepo,
		KeycloakClient:    keycloakClient,
		Interactor:        interactorFacade,
		MessageInteractor: messageInteractor,
		Config:            cfg,
		Logger:            log,
		IDEncoder:         encoder,
		MessagingCache:    messagingCache,
		ResponseHandler:   responseHandler,
	}, nil
}
