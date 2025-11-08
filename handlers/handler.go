package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
)

type handler struct {
	PersonService input.Service
	Interactor    *interactor.Interactor
}

func New(service input.Service, interactor*interactor.Interactor) *handler {
	return &handler{
		PersonService: service,
		Interactor:    interactor,
	}
}