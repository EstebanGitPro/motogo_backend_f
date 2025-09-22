package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/ports"
)

type Handler struct {
	PersonService ports.Service
}

func New(service ports.Service) *Handler {
	return &Handler{
		PersonService: service,
	}
}