package person

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/EstebanGitPro/motogo-backend/handlers"
)

type PersonRequest struct {
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	Password            string `json:"password"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Role                string `json:"role"`
}

type PersonResponse struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Role                string `json:"role"`
}

// Handler embeds the general handlers functionality
type Handler struct {
	*handlers.Handler
}

// New creates a new person handler with embedded general handler functionality
func New(service ports.Service) *Handler {
	return &Handler{
		Handler: handlers.New(service),
	}
}

func (p PersonRequest) ToDomain() domain.Person {
	return domain.Person{
		IdentityNumber:      p.IdentityNumber,
		FirstName:           p.FirstName,
		LastName:            p.LastName,
		SecondLastName:      p.SecondLastName,
		Email:               p.Email,
		PhoneNumber:         p.PhoneNumber,
		Password:            p.Password,
		EmailVerified:       p.EmailVerified,
		PhoneNumberVerified: p.PhoneNumberVerified,
		Role:                p.Role,
	}
}