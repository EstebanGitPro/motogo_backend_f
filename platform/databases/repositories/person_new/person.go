package personnew

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
)

type Person struct {
	ID                  string `db:"id"`
	IdentityNumber      string `db:"identity_number"`
	FirstName           string `db:"first_name"`
	LastName            string `db:"last_name"`
	SecondLastName      string `db:"second_last_name"`
	Email               string `db:"email"`
	PhoneNumber         string `db:"phone_number"`
	Role                string `db:"role"`
	KeycloakUserID      string `db:"keycloak_user_id"`
}

func (p Person) ToDomain() domain.Person {
	return domain.Person{
		ID:                  p.ID,
		IdentityNumber:      p.IdentityNumber,
		FirstName:           p.FirstName,
		LastName:            p.LastName,
		SecondLastName:      p.SecondLastName,
		Email:               p.Email,
		PhoneNumber:         p.PhoneNumber,
		Role:                p.Role,
		KeycloakUserID:      p.KeycloakUserID,
	}
}

func FromDomain(domainPerson domain.Person) Person {
	return Person{
		ID:                  domainPerson.ID,
		IdentityNumber:      domainPerson.IdentityNumber,
		FirstName:           domainPerson.FirstName,
		LastName:            domainPerson.LastName,
		SecondLastName:      domainPerson.SecondLastName,
		Email:               domainPerson.Email,
		PhoneNumber:         domainPerson.PhoneNumber,
		Role:                domainPerson.Role,
		KeycloakUserID:      domainPerson.KeycloakUserID,
	}
}