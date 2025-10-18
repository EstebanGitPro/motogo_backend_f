package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Password            string `json:"-"`
	Role                string `json:"role"`
	KeycloakUserID      string `json:"keycloak_user_id,omitempty"` // ID del usuario en Keycloak
}

func (u *Person) SetID() {
	u.ID = uuid.New().String()
}

func (u *Person) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
