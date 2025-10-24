package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/ports"
)

type handler struct {
	PersonService ports.Service
}

func New(service ports.Service) *handler {
	return &handler{
		PersonService: service,
	}
}