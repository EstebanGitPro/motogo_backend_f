package dependency

import (
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/platform/mysql"
)

type Dependencies struct {
	Config        *config.Config
}

func Init() (*Dependencies, error) {
	cfg := config.MustLoadConfig()

	_, err := mysql.GetDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		Config:        cfg,
	}, nil
}