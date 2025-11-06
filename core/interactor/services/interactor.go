package services

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

type Interactor struct {
	Service input.Service
	Repo output.Repository
}


func NewInteractor(srv input.Service, repo output.Repository) *Interactor {
	return &Interactor{Service: srv, Repo: repo}
}

func (i *Interactor) Execute() error {

	//Hay que ounto falla
	//Iniciar transacción
	tx, err := i.Repo.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()




	
	

	return nil
}