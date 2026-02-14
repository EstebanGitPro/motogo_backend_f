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

func (r *LoginRequest) Sanitize() {
	r.Email = TrimString(r.Email)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Links        []Link `json:"_links"`
}

// RefreshTokenRequest - DTO para refrescar access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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

// Sanitize trims whitespace from all string fields except Password
func (p *PersonRequest) Sanitize() {
	p.IdentityNumber = TrimString(p.IdentityNumber)
	p.FirstName = TrimString(p.FirstName)
	p.LastName = TrimString(p.LastName)
	p.SecondLastName = TrimString(p.SecondLastName)
	p.Email = TrimString(p.Email)
	p.PhoneNumber = TrimString(p.PhoneNumber)
	p.Role = TrimString(p.Role)
	// Password is intentionally NOT trimmed - spaces may be valid
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
		Role:           domain.Role(p.Role),
	}
}

// ResendVerificationEmailRequest - DTO para reenviar email de verificación
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Sanitize trims whitespace from email
func (r *ResendVerificationEmailRequest) Sanitize() {
	r.Email = TrimString(r.Email)
}

// PasswordResetRequest - DTO para solicitar recuperación de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Sanitize trims whitespace from email
func (r *PasswordResetRequest) Sanitize() {
	r.Email = TrimString(r.Email)
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

// ChangePasswordRequest - DTO para cambiar contraseña estando autenticado (HU57)
// El usuario debe proporcionar su contraseña actual para verificar identidad
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordResponse - Respuesta del cambio de contraseña (HU57)
type ChangePasswordResponse struct {
	Message string `json:"message"`
	Links   []Link `json:"_links"`
}

// UpdateProfileRequest - DTO for updating profile (HU52)
// All fields are optional (PATCH-like behavior)
type UpdateProfileRequest struct {
	IdentityNumber string `json:"identity_number,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	SecondLastName string `json:"second_last_name,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (u *UpdateProfileRequest) Sanitize() {
	u.IdentityNumber = TrimString(u.IdentityNumber)
	u.FirstName = TrimString(u.FirstName)
	u.LastName = TrimString(u.LastName)
	u.SecondLastName = TrimString(u.SecondLastName)
	u.PhoneNumber = TrimString(u.PhoneNumber)
}

// UpdateProfileResponse - Response for profile update (HU52)
type UpdateProfileResponse struct {
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

// PublicContactResponse - Solo info de contacto para motociclistas (HU55)
// Endpoint público para que puedan contactar al representante de sede
type PublicContactResponse struct {
	PhoneNumber string `json:"phone_number"`
	Links       []Link `json:"_links"`
}
