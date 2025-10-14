package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/EstebanGitPro/motogo-backend/core/services"
	"github.com/EstebanGitPro/motogo-backend/platform/keycloak"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/mysql"

	repo "github.com/EstebanGitPro/motogo-backend/repositories/person"
)

type Dependencies struct {
	PersonService  ports.Service
	PersonRepo     ports.Repository
	KeycloakClient ports.KeycloakClient
	AuthzService   ports.AuthorizationService
	Config         *config.Config
}

func Init() (*Dependencies, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	personRepo, err := repo.NewRepository(db)
	if err != nil {
		return nil, err
	}

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		return nil, err
	}

	authzService := services.NewAuthorizationService(keycloakClient, personRepo)

	personService := services.NewService(personRepo, authzService, cfg)

	return &Dependencies{
		PersonService:  personService,
		PersonRepo:     personRepo,
		KeycloakClient: keycloakClient,
		AuthzService:   authzService,
		Config:         cfg,
	}, nil
}
