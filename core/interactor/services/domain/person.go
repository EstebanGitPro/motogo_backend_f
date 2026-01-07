package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// Role constants - matches database CHECK constraint
const (
	RoleUser           = "user"
	RoleRepresentative = "representative"
	RoleAdmin          = "admin"
)

type Person struct {
	ID             string `json:"id"`
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Password       string `json:"-"`
	Role           string `json:"role"`
	KeycloakUserID string `json:"keycloak_user_id"`
}

func (u *Person) SetID() {
	u.ID = uuid.Generate()
}

func (p *Person) ToLogger() []string {
	return []string{
		"id:" + p.ID,
		"email:" + p.Email,
		"role:" + p.Role,
	}
}
