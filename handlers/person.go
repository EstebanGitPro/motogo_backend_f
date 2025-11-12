package handlers

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type PersonRequest struct {
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

type PersonResponse struct {
	ID             string `json:"id"`
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Role           string `json:"role"`
}


type RegistrationResponse struct {
	User    PersonResponse `json:"user"`
	Message string         `json:"message"`
	Links   []Link         `json:"_links"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (p PersonRequest) ToDomain() domain.Person {
	return domain.Person{
		IdentityNumber: p.IdentityNumber,
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		SecondLastName: p.SecondLastName,
		Email:          p.Email,
		PhoneNumber:    p.PhoneNumber,
		Password:       p.Password,
		Role:           p.Role,
	}
}
