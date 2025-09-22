package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/EstebanGitPro/motogo-backend/core/services"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/mysql"
	repo "github.com/EstebanGitPro/motogo-backend/repositories/person"
)

type Dependencies struct {
	PersonService ports.Service
	PersonRepo    ports.Repository
	Config        *config.Config
}

func Init() (*Dependencies, error) {
	cfg := config.MustLoadConfig()

	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	personRepo := repo.NewRepository(db)
	personService := services.NewService(personRepo)

	return &Dependencies{
		PersonService: personService,
		PersonRepo:    personRepo,
		Config:        cfg,
	}, nil
}