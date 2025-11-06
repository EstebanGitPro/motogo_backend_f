package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
)

type handler struct {
	PersonService input.Service
}

func New(service input.Service) *handler {
	return &handler{
		PersonService: service,
	}
}