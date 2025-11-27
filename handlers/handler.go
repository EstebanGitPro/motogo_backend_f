package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
)

type handler struct {
	PersonService input.Service
	Interactor    *interactor.Interactor
	Logger        logger.Logger
	IDEncoder     *idencoder.HashidsEncoder
}

func New(service input.Service, interactor *interactor.Interactor, log logger.Logger, encoder *idencoder.HashidsEncoder) *handler {
	return &handler{
		PersonService: service,
		Interactor:    interactor,
		Logger:        log,
		IDEncoder:     encoder,
	}
}
