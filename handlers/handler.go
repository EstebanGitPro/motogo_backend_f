package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
)

type handler struct {
	PersonService     input.Service
	Interactor        *interactor.Interactor
	MessageInteractor *interactor.MessageInteractor
	Logger            logger.Logger
	IDEncoder         *idencoder.HashidsEncoder
	Response          *middleware.ResponseHandler
}

func New(service input.Service, personInteractor *interactor.Interactor, messageInteractor *interactor.MessageInteractor, log logger.Logger, encoder *idencoder.HashidsEncoder, responseHandler *middleware.ResponseHandler) *handler {
	return &handler{
		PersonService:     service,
		Interactor:        personInteractor,
		MessageInteractor: messageInteractor,
		Logger:            log,
		IDEncoder:         encoder,
		Response:          responseHandler,
	}
}
