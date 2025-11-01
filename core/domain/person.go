package domain

import (
	"github.com/google/uuid"
)

type Person struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	Password            string `json:"-"`
	Role                string `json:"role"`
	KeycloakUserID      string `json:"keycloak_user_id"`
}

func (u *Person) SetID() {
	u.ID = uuid.New().String()
}

