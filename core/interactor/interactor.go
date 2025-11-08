package services

import (
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
	// TODO: Implement execution logic
	// Note: BeginTransaction is not part of output.Repository interface
	// If transaction management is needed, it should be handled at the repository implementation level
	return nil
}