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
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Role           string `json:"role"`
}

type RegistrationResponse struct {
	Links []Link `json:"_links"`
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
	Links        []Link `json:"_links"`
}

// AuthMeResponse represents the authenticated user profile response
type AuthMeResponse struct {
	ID             string `json:"id"`
	IdentityNumber string `json:"identity_number"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Role           string `json:"role"`
	Links          []Link `json:"_links"`
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

// ResendVerificationEmailRequest - DTO para reenviar email de verificación
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetRequest - DTO para solicitar recuperación de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest - DTO para verificar email mediante token proxy
// Este token es un JWT que contiene el email del usuario
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse - Respuesta de verificación de email
type VerifyEmailResponse struct {
	Verified bool   `json:"verified"`
	Email    string `json:"email"`
}

// ResetPasswordWithTokenRequest - DTO para actualizar contraseña con token
// El token es un JWT que viene del enlace del email de recuperación de contraseña
type ResetPasswordWithTokenRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
