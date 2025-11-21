package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

type handler struct {
	PersonService input.Service
	Interactor    *interactor.Interactor
	Logger        logger.Logger
}

func New(service input.Service, interactor *interactor.Interactor, log logger.Logger) *handler {
	return &handler{
		PersonService: service,
		Interactor:    interactor,
		Logger:        log,
	}
}
