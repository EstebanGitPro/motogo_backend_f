package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/platform/keycloak"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/mysql"

	repo "github.com/EstebanGitPro/motogo-backend/repositories/person"
)

type Dependencies struct {
	PersonService  input.Service
	PersonRepo     output.Repository
	KeycloakClient output.AuthClient
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


	personService := services.NewService(personRepo, cfg)

	return &Dependencies{
		PersonService:  personService,
		PersonRepo:     personRepo,
		KeycloakClient: keycloakClient,
		Config:         cfg,
	}, nil
}
