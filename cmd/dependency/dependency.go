package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/identity_provider/keycloak"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/databases/mysql"

	repo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/person"
)

type Dependencies struct {
	PersonService  input.Service
	PersonRepo     output.Repository
	KeycloakClient output.AuthClient
	Interactor     *interactor.Interactor
	Config         *config.Config
	Logger         logger.Logger
	IDEncoder      *idencoder.HashidsEncoder
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

	log.Success("Dependencias inicializadas correctamente")

	return &Dependencies{
		PersonService:  personService,
		PersonRepo:     personRepo,
		KeycloakClient: keycloakClient,
		Interactor:     interactorFacade,
		Config:         cfg,
		Logger:         log,
		IDEncoder:      encoder,
	}, nil
}
