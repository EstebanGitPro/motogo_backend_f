package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/platform/identity_provider/keycloak"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/databases/mysql"

	repo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/person"
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

	

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		return nil, err
	}

	personRepo, err := repo.NewClientRepository(db,keycloakClient)
	if err != nil {
		return nil, err
	}


	personService := services.NewService(personRepo,keycloakClient)

	return &Dependencies{
		PersonService:  personService,
		PersonRepo:     personRepo,
		KeycloakClient: keycloakClient,
		Config:         cfg,
	}, nil
}
