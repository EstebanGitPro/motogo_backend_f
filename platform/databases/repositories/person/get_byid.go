package person

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetPersonByID(id string) (*domain.Person, error) {
	var p Person
	err := r.db.QueryRow(queryGetByID, id).Scan(
		&p.ID,
		&p.IdentityNumber,
		&p.FirstName,
		&p.LastName,
		&p.SecondLastName,
		&p.Email,
		&p.PhoneNumber,
		&p.Role,
		&p.KeycloakUserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}

	domainPerson := p.ToDomain()
	return &domainPerson, nil
}
