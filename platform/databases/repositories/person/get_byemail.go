package person

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, queryGetByEmail, email).Scan(
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
